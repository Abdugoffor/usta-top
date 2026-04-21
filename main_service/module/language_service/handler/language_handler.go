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
	h := &languageHandler{service: language_service.NewLanguageService(db)}

	routes := group + "/languages"
	{
		router.POST(routes, h.Create)
		router.GET(routes, h.List)
		router.GET(routes+"/:id", h.GetByID)
		router.PUT(routes+"/:id", h.Update)
		router.DELETE(routes+"/:id", middleware.CheckRole(h.Delete))
	}
}

// Create godoc
// @Summary      Yangi til yaratish
// @Tags         Languages
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      language_dto.CreateLanguageRequest  true  "Til ma'lumotlari"
// @Success      201   {object}  language_dto.LanguageResponse
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /languages [post]
func (h *languageHandler) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	resp, err := h.service.Create(r.Context(), req)
	{
		if err != nil {
			helper.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	helper.WriteJSON(w, http.StatusCreated, resp)
}

// List godoc
// @Summary      Tillar ro'yxati
// @Tags         Languages
// @Produce      json
// @Param        name        query     string  false  "Nomi bo'yicha filter"
// @Param        is_active   query     boolean false  "Faol/faolsiz"
// @Param        page        query     integer false  "Sahifa" default(1)
// @Param        limit       query     integer false  "Limit" default(10)
// @Param        sort_by     query     string  false  "Saralash maydoni (id, name, is_active, created_at, updated_at)"
// @Param        sort_order  query     string  false  "asc yoki desc"
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /languages [get]
func (h *languageHandler) List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	q := r.URL.Query()
	pq := helper.ParsePage(r)

	f := language_dto.LanguageFilter{
		Name: q.Get("name"),
	}

	if v := q.Get("is_active"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			f.IsActive = &b
		}
	}

	items, total, err := h.service.List(r.Context(), f, pq.Page, pq.Limit, pq.SortCol, pq.SortOrder)
	{
		if err != nil {
			helper.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	helper.WriteJSON(w, http.StatusOK, map[string]any{
		"data": items,
		"meta": helper.NewPageMeta(total, pq.Page, pq.Limit),
	})
}

// GetByID godoc
// @Summary      Tilni ID bo'yicha olish
// @Tags         Languages
// @Produce      json
// @Param        id   path      integer  true  "Til ID"
// @Success      200  {object}  language_dto.LanguageResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /languages/{id} [get]
func (h *languageHandler) GetByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
			helper.WriteError(w, http.StatusNotFound, "language not found")
			return
		}
	}

	helper.WriteJSON(w, http.StatusOK, resp)
}

// Update godoc
// @Summary      Tilni yangilash
// @Tags         Languages
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      integer                             true  "Til ID"
// @Param        body  body      language_dto.UpdateLanguageRequest  true  "Yangi ma'lumotlar"
// @Success      200   {object}  language_dto.LanguageResponse
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Router       /languages/{id} [put]
func (h *languageHandler) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	resp, err := h.service.Update(r.Context(), id, req)
	{
		if err != nil {
			helper.WriteError(w, http.StatusNotFound, "language not found")
			return
		}
	}

	helper.WriteJSON(w, http.StatusOK, resp)
}

// Delete godoc
// @Summary      Tilni o'chirish
// @Tags         Languages
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      integer  true  "Til ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /languages/{id} [delete]
func (h *languageHandler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
