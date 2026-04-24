package country_handler

import (
	"encoding/json"
	"main_service/helper"
	"main_service/middleware"
	country_dto "main_service/module/country_service/dto"
	country_service "main_service/module/country_service/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
)

type countryHandler struct {
	service country_service.CountryService
}

func NewCountryHandler(router *httprouter.Router, group string, db *pgxpool.Pool) {
	h := &countryHandler{service: country_service.NewCountryService(db)}

	routes := group + "/countries"
	{
		router.POST(routes, middleware.CheckRole(h.Create))
		router.GET(routes, h.List)
		router.GET(routes+"/:id", h.GetByID)
		router.PUT(routes+"/:id", middleware.CheckRole(h.Update))
		router.DELETE(routes+"/:id", middleware.CheckRole(h.Delete))
	}
}

func (h *countryHandler) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req country_dto.CreateCountryRequest
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

	resp, err := h.service.Create(r.Context(), req)
	{
		if err != nil {
			helper.WriteInternalError(w, err)
			return
		}
	}

	helper.WriteJSON(w, http.StatusCreated, resp)
}

func (h *countryHandler) List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	q := r.URL.Query()
	afterID, limit := helper.ParseCursorPage(r)

	name := q.Get("name")
	if len(name) > 100 {
		name = name[:100]
	}
	f := country_dto.CountryFilter{Name: name}

	if v := q.Get("is_active"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			f.IsActive = &b
		}
	}

	var parentID *int64
	if v := q.Get("parent_id"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			parentID = &n
		}
	}

	items, hasMore, err := h.service.ListTree(r.Context(), parentID, f, afterID, limit)
	if err != nil {
		helper.WriteInternalError(w, err)
		return
	}

	var lastID int64
	if len(items) > 0 {
		lastID = items[len(items)-1].ID
	}

	helper.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"data": items,
		"meta": helper.NewCursorMeta(limit, hasMore, lastID, 0),
	})
}

func (h *countryHandler) GetByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	{
		if err != nil || id <= 0 {
			helper.WriteError(w, http.StatusBadRequest, "invalid id")
			return
		}
	}

	resp, err := h.service.GetByID(r.Context(), id)
	{
		if err != nil {
			helper.WriteError(w, http.StatusNotFound, "country not found")
			return
		}
	}

	helper.WriteJSON(w, http.StatusOK, resp)
}

func (h *countryHandler) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	{
		if err != nil || id <= 0 {
			helper.WriteError(w, http.StatusBadRequest, "invalid id")
			return
		}
	}

	var req country_dto.UpdateCountryRequest
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

	resp, err := h.service.Update(r.Context(), id, req)
	{
		if err != nil {
			helper.WriteError(w, http.StatusNotFound, "country not found")
			return
		}
	}

	helper.WriteJSON(w, http.StatusOK, resp)
}

func (h *countryHandler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	{
		if err != nil || id <= 0 {
			helper.WriteError(w, http.StatusBadRequest, "invalid id")
			return
		}
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		helper.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}
