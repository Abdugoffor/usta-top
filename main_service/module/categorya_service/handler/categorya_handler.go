package categorya_handler

import (
	"encoding/json"
	"main_service/helper"
	"main_service/middleware"
	categorya_dto "main_service/module/categorya_service/dto"
	categorya_service "main_service/module/categorya_service/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
)

type categoryHandler struct {
	service categorya_service.CategoryService
}

var sortCols = map[string]string{
	"id": "c.id", "name": "c.name",
	"is_active": "c.is_active", "created_at": "c.created_at", "updated_at": "c.updated_at",
}

func NewCategoryHandler(router *httprouter.Router, group string, db *pgxpool.Pool) {
	h := &categoryHandler{service: categorya_service.NewCategoryService(db)}

	routes := group + "/categories"
	{
		router.POST(routes, middleware.Auth(h.Create))
		router.GET(routes, h.List)
		router.GET(routes+"/:id", h.GetByID)
		router.PUT(routes+"/:id", middleware.Auth(h.Update))
		router.DELETE(routes+"/:id", middleware.Auth(h.Delete))
	}
}

// Create godoc
// @Summary      Yangi kategoriya yaratish
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      categorya_dto.CreateCategoryRequest  true  "Kategoriya ma'lumotlari"
// @Success      201   {object}  categorya_dto.CategoryResponse
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /categories [post]
func (h *categoryHandler) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req categorya_dto.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	if errs := helper.ValidateStruct(req); errs != nil {
		helper.WriteValidation(w, errs)
		return
	}

	resp, err := h.service.Create(r.Context(), req)
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.WriteJSON(w, http.StatusCreated, resp)
}

// List godoc
// @Summary      Kategoriyalar ro'yxati
// @Tags         Categories
// @Produce      json
// @Param        name        query     string  false  "Nomi bo'yicha filter"
// @Param        is_active   query     boolean false  "Faol/faolsiz"
// @Param        page        query     integer false  "Sahifa" default(1)
// @Param        limit       query     integer false  "Limit" default(10)
// @Param        sort_by     query     string  false  "Saralash maydoni"
// @Param        sort_order  query     string  false  "asc yoki desc"
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /categories [get]
func (h *categoryHandler) List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	q := r.URL.Query()
	pq := helper.ParsePage(r, sortCols, "c.id")

	f := categorya_dto.CategoryFilter{Name: q.Get("name")}
	if v := q.Get("is_active"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			f.IsActive = &b
		}
	}

	items, total, err := h.service.List(r.Context(), f, pq.Page, pq.Limit, pq.SortCol, pq.SortOrder)
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"data": items,
		"meta": helper.NewPageMeta(total, pq.Page, pq.Limit),
	})
}

// GetByID godoc
// @Summary      Kategoriyani ID bo'yicha olish
// @Tags         Categories
// @Produce      json
// @Param        id   path      integer  true  "Kategoriya ID"
// @Success      200  {object}  categorya_dto.CategoryResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /categories/{id} [get]
func (h *categoryHandler) GetByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil || id <= 0 {
		helper.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}
	resp, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		helper.WriteError(w, http.StatusNotFound, "category not found")
		return
	}
	helper.WriteJSON(w, http.StatusOK, resp)
}

// Update godoc
// @Summary      Kategoriyani yangilash
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      integer                              true  "Kategoriya ID"
// @Param        body  body      categorya_dto.UpdateCategoryRequest  true  "Yangi ma'lumotlar"
// @Success      200   {object}  categorya_dto.CategoryResponse
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Router       /categories/{id} [put]
func (h *categoryHandler) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil || id <= 0 {
		helper.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req categorya_dto.UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
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
	if err != nil {
		helper.WriteError(w, http.StatusNotFound, "category not found")
		return
	}
	helper.WriteJSON(w, http.StatusOK, resp)
}

// Delete godoc
// @Summary      Kategoriyani o'chirish
// @Tags         Categories
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      integer  true  "Kategoriya ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /categories/{id} [delete]
func (h *categoryHandler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil || id <= 0 {
		helper.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		helper.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}
