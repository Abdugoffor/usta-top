package region_handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"main_service/middleware"
	region_dto "main_service/module/region_service/dto"
	region_service "main_service/module/region_service/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
)

type regionHandler struct {
	db      *pgxpool.Pool
	service region_service.RegionService
}

func NewRegionHandler(router *httprouter.Router, group string, db *pgxpool.Pool) {
	h := &regionHandler{
		db:      db,
		service: region_service.NewRegionService(db),
	}

	routes := group + "/regions"
	{
		router.POST(routes, middleware.CheckRole(h.Create))
		router.GET(routes, h.List)
		router.GET(routes+"/:id", middleware.CheckRole(h.GetByID))
		router.PUT(routes+"/:id", middleware.CheckRole(h.Update))
		router.DELETE(routes+"/:id", middleware.CheckRole(h.Delete))
	}
}

// ─── Response helpers ────────────────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// ─── Validation ──────────────────────────────────────────────────────────────

func validateCreate(req region_dto.CreateRegionRequest) string {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return "name is required"
	}
	if len(name) < 2 {
		return "name must be at least 2 characters"
	}
	if len(name) > 255 {
		return "name must be at most 255 characters"
	}
	if req.ParentID < 0 {
		return "parent_id must be >= 0"
	}
	if len(req.Description) > 1000 {
		return "description must be at most 1000 characters"
	}
	return ""
}

func validateUpdate(req region_dto.UpdateRegionRequest) string {
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return "name cannot be empty"
		}
		if len(name) < 2 {
			return "name must be at least 2 characters"
		}
		if len(name) > 255 {
			return "name must be at most 255 characters"
		}
	}
	if req.ParentID != nil && *req.ParentID < 0 {
		return "parent_id must be >= 0"
	}
	if req.Description != nil && len(*req.Description) > 1000 {
		return "description must be at most 1000 characters"
	}
	return ""
}

// Create godoc
// @Summary      Yangi region yaratish
// @Tags         Regions
// @Accept       json
// @Produce      json
// @Param        body  body      region_dto.CreateRegionRequest  true  "Region ma'lumotlari"
// @Success      201   {object}  region_dto.RegionResponse
// @Failure      400   {object}  map[string]string
// @Failure      422   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /regions [post]
func (h *regionHandler) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req region_dto.CreateRegionRequest
	{
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
			return
		}
	}

	req.Name = strings.TrimSpace(req.Name)

	if msg := validateCreate(req); msg != "" {
		writeError(w, http.StatusUnprocessableEntity, msg)
		return
	}

	resp, err := h.service.Create(r.Context(), req)
	{
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	writeJSON(w, http.StatusCreated, resp)
}

// List godoc
// @Summary      Regionlar ro'yxati
// @Tags         Regions
// @Produce      json
// @Param        name        query     string  false  "Nomi bo'yicha filter"
// @Param        parent_id   query     integer false  "Parent ID"
// @Param        is_active   query     boolean false  "Faol/faolsiz"
// @Param        page        query     integer false  "Sahifa raqami" default(1)
// @Param        limit       query     integer false  "Sahifadagi elementlar soni" default(10)
// @Param        sort_by     query     string  false  "Saralash maydoni"
// @Param        sort_order  query     string  false  "asc yoki desc"
// @Success      200  {object}  region_dto.ListRegionsResponse
// @Failure      500  {object}  map[string]string
// @Router       /regions [get]
func (h *regionHandler) List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	q := r.URL.Query()

	f := region_dto.RegionFilter{
		Name:      q.Get("name"),
		SortBy:    q.Get("sort_by"),
		SortOrder: q.Get("sort_order"),
		Page:      1,
		Limit:     10,
	}

	parentID := int64(0)
	if v := q.Get("parent_id"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil && n >= 0 {
			parentID = n
		}
	}
	if v := q.Get("is_active"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			f.IsActive = &b
		}
	}
	if v := q.Get("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			f.Page = n
		}
	}
	if v := q.Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			f.Limit = n
		}
	}

	regions, total, err := h.service.ListTree(r.Context(), parentID, f)
	{
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	totalPages := total / int64(f.Limit)
	{
		if total%int64(f.Limit) != 0 {
			totalPages++
		}
	}

	writeJSON(w, http.StatusOK, region_dto.ListRegionsResponse{
		Regions:    regions,
		Total:      total,
		Page:       f.Page,
		Limit:      f.Limit,
		TotalPages: totalPages,
	})
}

// GetByID godoc
// @Summary      Regionni ID bo'yicha olish
// @Tags         Regions
// @Produce      json
// @Param        id   path      integer  true  "Region ID"
// @Success      200  {object}  region_dto.RegionResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /regions/{id} [get]
func (h *regionHandler) GetByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	{
		if err != nil || id <= 0 {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}
	}

	resp, err := h.service.GetByID(r.Context(), id)
	{
		if err != nil {
			writeError(w, http.StatusNotFound, "region not found")
			return
		}
	}

	writeJSON(w, http.StatusOK, resp)
}

// Update godoc
// @Summary      Regionni yangilash
// @Tags         Regions
// @Accept       json
// @Produce      json
// @Param        id    path      integer                         true  "Region ID"
// @Param        body  body      region_dto.UpdateRegionRequest  true  "Yangi ma'lumotlar"
// @Success      200   {object}  region_dto.RegionResponse
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      422   {object}  map[string]string
// @Router       /regions/{id} [put]
func (h *regionHandler) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	{
		if err != nil || id <= 0 {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}
	}

	var req region_dto.UpdateRegionRequest
	{
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
			return
		}
	}

	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		req.Name = &trimmed
	}

	if msg := validateUpdate(req); msg != "" {
		writeError(w, http.StatusUnprocessableEntity, msg)
		return
	}

	resp, err := h.service.Update(r.Context(), id, req)
	{
		if err != nil {
			writeError(w, http.StatusNotFound, "region not found or update failed")
			return
		}
	}

	writeJSON(w, http.StatusOK, resp)
}

// Delete godoc
// @Summary      Regionni o'chirish
// @Tags         Regions
// @Produce      json
// @Param        id   path      integer  true  "Region ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /regions/{id} [delete]
func (h *regionHandler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	{
		if err != nil || id <= 0 {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "deleted successfully"})
}
