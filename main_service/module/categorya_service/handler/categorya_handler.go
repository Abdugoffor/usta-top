package categorya_handler

import (
	"encoding/json"
	"main_service/helper"
	"main_service/middleware"
	categorya_dto "main_service/module/categorya_service/dto"
	categorya_service "main_service/module/categorya_service/service"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
)

type categoryHandler struct {
	service categorya_service.CategoryService
}

func NewCategoryHandler(router *httprouter.Router, group string, db *pgxpool.Pool) {
	h := &categoryHandler{service: categorya_service.NewCategoryService(db)}

	routes := group + "/categories"
	{
		router.POST(routes, middleware.CheckRole(h.Create))
		router.GET(routes, h.List)
		router.GET(routes+"/:id", h.GetByID)
		router.PUT(routes+"/:id", middleware.CheckRole(h.Update))
		router.DELETE(routes+"/:id", middleware.CheckRole(h.Delete))
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

	if errs := helper.ValidateStruct(req); errs != nil {
		helper.WriteValidation(w, errs)
		return
	}

	resp, err := h.service.Create(r.Context(), req)
	if err != nil {
		helper.WriteInternalError(w, err)
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
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /categories [get]
func (h *categoryHandler) List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	q := r.URL.Query()
	afterID, limit := helper.ParseCursorPage(r)

	name := q.Get("name")
	if len(name) > 100 {
		name = name[:100]
	}
	f := categorya_dto.CategoryFilter{Name: name}

	if v := q.Get("is_active"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			f.IsActive = &b
		}
	}

	items, hasMore, err := h.service.List(r.Context(), f, afterID, limit)
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
