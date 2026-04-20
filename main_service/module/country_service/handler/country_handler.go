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

var sortCols = map[string]string{
	"id": "c.id", "parent_id": "c.parent_id", "name": "c.name",
	"is_active": "c.is_active", "created_at": "c.created_at", "updated_at": "c.updated_at",
}

func NewCountryHandler(router *httprouter.Router, group string, db *pgxpool.Pool) {
	h := &countryHandler{service: country_service.NewCountryService(db)}

	routes := group + "/countries"
	{
		router.POST(routes, h.Create)
		router.GET(routes, h.List)
		router.GET(routes+"/:id", h.GetByID)
		router.PUT(routes+"/:id", h.Update)
		router.DELETE(routes+"/:id", middleware.CheckRole(h.Delete))
	}
}

// Create godoc
// @Summary      Yangi mamlakat yaratish
// @Tags         Countries
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      country_dto.CreateCountryRequest  true  "Mamlakat ma'lumotlari"
// @Success      201   {object}  country_dto.CountryResponse
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /countries [post]
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
			helper.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	helper.WriteJSON(w, http.StatusCreated, resp)
}

// List godoc
// @Summary      Mamlakatlar ro'yxati
// @Tags         Countries
// @Produce      json
// @Param        name        query     string  false  "Nomi bo'yicha filter"
// @Param        parent_id   query     integer false  "Parent ID"
// @Param        is_active   query     boolean false  "Faol/faolsiz"
// @Param        page        query     integer false  "Sahifa" default(1)
// @Param        limit       query     integer false  "Limit" default(10)
// @Param        sort_by     query     string  false  "Saralash maydoni"
// @Param        sort_order  query     string  false  "asc yoki desc"
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /countries [get]
func (h *countryHandler) List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	q := r.URL.Query()
	pq := helper.ParsePage(r, sortCols, "c.id")

	f := country_dto.CountryFilter{
		Name: q.Get("name"),
	}

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

	items, total, err := h.service.ListTree(r.Context(), parentID, f, pq.Page, pq.Limit, pq.SortCol, pq.SortOrder)
	{
		if err != nil {
			helper.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	helper.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"data": items,
		"meta": helper.NewPageMeta(total, pq.Page, pq.Limit),
	})
}

// GetByID godoc
// @Summary      Mamlakatni ID bo'yicha olish
// @Tags         Countries
// @Produce      json
// @Param        id   path      integer  true  "Mamlakat ID"
// @Success      200  {object}  country_dto.CountryResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /countries/{id} [get]
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

// Update godoc
// @Summary      Mamlakatni yangilash
// @Tags         Countries
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      integer                           true  "Mamlakat ID"
// @Param        body  body      country_dto.UpdateCountryRequest  true  "Yangi ma'lumotlar"
// @Success      200   {object}  country_dto.CountryResponse
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Router       /countries/{id} [put]
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

// Delete godoc
// @Summary      Mamlakatni o'chirish
// @Tags         Countries
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      integer  true  "Mamlakat ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /countries/{id} [delete]
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
