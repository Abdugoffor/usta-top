package language_handler

import (
	"encoding/json"
	"main_service/helper"
	"main_service/middleware"
	language_dto "main_service/module/language_service/dto"
	language_service "main_service/module/language_service/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
)

type languageHandler struct {
	service language_service.LanguageService
}

func NewLanguageHandler(router *httprouter.Router, group string, db *pgxpool.Pool) {
	handler := &languageHandler{service: language_service.NewLanguageService(db)}

	routes := group + "/languages"
	{
		router.GET(routes, middleware.CheckRole(handler.All, "admin"))
		router.GET(routes+"/:id", middleware.CheckRole(handler.Show, "admin"))
		router.POST(routes, middleware.CheckRole(handler.Create, "admin"))
		router.PUT(routes+"/:id", middleware.CheckRole(handler.Update, "admin"))
		router.DELETE(routes+"/:id", middleware.CheckRole(handler.Delete, "admin"))
	}
}

func (handler *languageHandler) All(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	q := r.URL.Query()

	afterID, limit := helper.ParseCursorPage(r)

	name := q.Get("name")

	if len(name) > 100 {
		name = name[:100]
	}
	f := language_dto.LanguageFilter{Name: name}

	if v := q.Get("is_active"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			f.IsActive = &b
		}
	}

	items, hasMore, err := handler.service.All(r.Context(), f, afterID, limit)
	{
		if err != nil {
			helper.WriteInternalError(w, err)
			return
		}
	}

	var lastID int64
	{
		if len(items) > 0 {
			lastID = items[len(items)-1].ID
		}
	}

	helper.WriteJSON(w, http.StatusOK, map[string]any{
		"data": items,
		"meta": helper.NewCursorMeta(limit, hasMore, lastID, 0),
	})
}

func (handler *languageHandler) Show(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	{
		if err != nil || id <= 0 {
			helper.WriteError(w, http.StatusBadRequest, "invalid id")
			return
		}
	}

	resp, err := handler.service.Show(r.Context(), id)
	{
		if err != nil {
			helper.WriteError(w, http.StatusNotFound, "language not found")
			return
		}
	}

	helper.WriteJSON(w, http.StatusOK, resp)
}

func (handler *languageHandler) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req language_dto.CreateLanguageRequest
	{
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helper.WriteError(w, http.StatusBadRequest, "invalid JSON")
			return
		}
	}

	req.Name = strings.TrimSpace(req.Name)
	{
		if errs := helper.ValidateStruct(req); errs != nil {
			helper.WriteValidation(w, errs)
			return
		}
	}

	resp, err := handler.service.Create(r.Context(), req)
	{
		if err != nil {
			helper.WriteInternalError(w, err)
			return
		}
	}

	helper.WriteJSON(w, http.StatusCreated, resp)
}

func (handler *languageHandler) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	{
		if err != nil || id <= 0 {
			helper.WriteError(w, http.StatusBadRequest, "invalid id")
			return
		}
	}

	var req language_dto.UpdateLanguageRequest
	{
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helper.WriteError(w, http.StatusBadRequest, "invalid JSON")
			return
		}
	}

	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		req.Name = &trimmed
	}

	if errs := helper.ValidateStruct(req); errs != nil {
		helper.WriteValidation(w, errs)
		return
	}

	resp, err := handler.service.Update(r.Context(), id, req)
	{
		if err != nil {
			helper.WriteError(w, http.StatusNotFound, "language not found")
			return
		}
	}

	helper.WriteJSON(w, http.StatusOK, resp)
}

func (handler *languageHandler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	{
		if err != nil || id <= 0 {
			helper.WriteError(w, http.StatusBadRequest, "invalid id")
			return
		}
	}

	if err := handler.service.Delete(r.Context(), id); err != nil {
		helper.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}
