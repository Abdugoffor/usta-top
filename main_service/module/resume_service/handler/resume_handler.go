package resume_handler

import (
	"encoding/json"
	"main_service/helper"
	"main_service/middleware"
	resume_dto "main_service/module/resume_service/dto"
	resume_service "main_service/module/resume_service/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
)

type resumeHandler struct {
	service resume_service.ResumeService
}

var sortCols = map[string]string{
	"id": "rs.id", "name": "rs.name", "title": "rs.title",
	"price": "rs.price", "experience_year": "rs.experience_year",
	"views_count": "rs.views_count",
	"is_active":   "rs.is_active", "created_at": "rs.created_at", "updated_at": "rs.updated_at",
}

func NewResumeHandler(router *httprouter.Router, group string, db *pgxpool.Pool) {
	h := &resumeHandler{service: resume_service.NewResumeService(db)}

	routes := group + "/resumes"
	{
		router.POST(routes, middleware.Auth(h.Create))
		router.GET(routes, h.List)
		router.GET(routes+"/:slug", h.GetBySlug)
		router.PUT(routes+"/:id", middleware.Auth(h.Update))
		router.DELETE(routes+"/:id", middleware.Auth(h.Delete))
		router.POST(routes+"/:id/categories", middleware.Auth(h.AddCategory))
		router.DELETE(routes+"/:id/categories/:cat_id", middleware.Auth(h.RemoveCategory))
	}
}

// Create godoc
// @Summary      Yangi resume yaratish
// @Tags         Resumes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      resume_dto.CreateResumeRequest  true  "Resume ma'lumotlari"
// @Success      201   {object}  resume_dto.ResumeResponse
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /resumes [post]
func (h *resumeHandler) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID := middleware.GetUserID(r)
	if userID == 0 {
		helper.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req resume_dto.CreateResumeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Title = strings.TrimSpace(req.Title)
	req.Adress = strings.TrimSpace(req.Adress)

	if errs := helper.ValidateStruct(req); errs != nil {
		helper.WriteValidation(w, errs)
		return
	}

	resp, err := h.service.Create(r.Context(), int64(userID), req)
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.WriteJSON(w, http.StatusCreated, resp)
}

// List godoc
// @Summary      Resumelar ro'yxati
// @Tags         Resumes
// @Produce      json
// @Param        name            query     string  false  "Nomi"
// @Param        title           query     string  false  "Sarlavha"
// @Param        user_id         query     integer false  "Foydalanuvchi ID"
// @Param        region_id       query     integer false  "Region ID"
// @Param        district_id     query     integer false  "Tuman ID"
// @Param        mahalla_id      query     integer false  "Mahalla ID"
// @Param        category_id     query     integer false  "Kategoriya ID"
// @Param        is_active       query     boolean false  "Faol/faolsiz"
// @Param        min_price       query     integer false  "Minimal narx"
// @Param        max_price       query     integer false  "Maksimal narx"
// @Param        min_experience  query     integer false  "Minimal tajriba (yil)"
// @Param        page            query     integer false  "Sahifa" default(1)
// @Param        limit           query     integer false  "Limit" default(10)
// @Param        sort_by         query     string  false  "Saralash maydoni"
// @Param        sort_order      query     string  false  "asc yoki desc"
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /resumes [get]
func (h *resumeHandler) List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	q := r.URL.Query()
	pq := helper.ParsePage(r, sortCols, "rs.id")

	f := resume_dto.ResumeFilter{
		Name:  q.Get("name"),
		Title: q.Get("title"),
	}
	if v := q.Get("user_id"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			f.UserID = &n
		}
	}
	if v := q.Get("region_id"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			f.RegionID = &n
		}
	}
	if v := q.Get("district_id"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			f.DistrictID = &n
		}
	}
	if v := q.Get("mahalla_id"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			f.MahallaID = &n
		}
	}
	if v := q.Get("is_active"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			f.IsActive = &b
		}
	}
	if v := q.Get("min_price"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			f.MinPrice = &n
		}
	}
	if v := q.Get("max_price"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			f.MaxPrice = &n
		}
	}
	if v := q.Get("min_experience"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			f.MinExperience = &n
		}
	}
	if v := q.Get("category_id"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			f.CategoryID = &n
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

// GetBySlug godoc
// @Summary      Resumeni slug bo'yicha olish
// @Tags         Resumes
// @Produce      json
// @Param        slug  path      string  true  "Resume slug"
// @Success      200   {object}  resume_dto.ResumeResponse
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Router       /resumes/{slug} [get]
func (h *resumeHandler) GetBySlug(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	slug := ps.ByName("slug")
	if slug == "" {
		helper.WriteError(w, http.StatusBadRequest, "invalid slug")
		return
	}
	resp, err := h.service.GetBySlug(r.Context(), slug)
	if err != nil {
		helper.WriteError(w, http.StatusNotFound, "resume not found")
		return
	}
	helper.WriteJSON(w, http.StatusOK, resp)
}

// Update godoc
// @Summary      Resumeni yangilash
// @Tags         Resumes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      integer                         true  "Resume ID"
// @Param        body  body      resume_dto.UpdateResumeRequest  true  "Yangi ma'lumotlar"
// @Success      200   {object}  resume_dto.ResumeResponse
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Router       /resumes/{id} [put]
func (h *resumeHandler) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := middleware.GetUserID(r)
	if userID == 0 {
		helper.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil || id <= 0 {
		helper.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req resume_dto.UpdateResumeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if errs := helper.ValidateStruct(req); errs != nil {
		helper.WriteValidation(w, errs)
		return
	}

	resp, err := h.service.Update(r.Context(), id, int64(userID), req)
	if err != nil {
		helper.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	helper.WriteJSON(w, http.StatusOK, resp)
}

// Delete godoc
// @Summary      Resumeni o'chirish
// @Tags         Resumes
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      integer  true  "Resume ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /resumes/{id} [delete]
func (h *resumeHandler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := middleware.GetUserID(r)
	if userID == 0 {
		helper.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil || id <= 0 {
		helper.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.service.Delete(r.Context(), id, int64(userID)); err != nil {
		helper.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}

// AddCategory godoc
// @Summary      Resumega kategoriya qo'shish
// @Tags         Resumes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      integer                    true  "Resume ID"
// @Param        body  body      object{category_id=integer}  true  "Kategoriya ID"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /resumes/{id}/categories [post]
func (h *resumeHandler) AddCategory(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil || id <= 0 {
		helper.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var body struct {
		CategoryID int64 `json:"category_id" validate:"required,min=1"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if errs := helper.ValidateStruct(body); errs != nil {
		helper.WriteValidation(w, errs)
		return
	}
	if err := h.service.AddCategory(r.Context(), id, body.CategoryID); err != nil {
		helper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "category added"})
}

// RemoveCategory godoc
// @Summary      Resumedan kategoriyani olib tashlash
// @Tags         Resumes
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      integer  true  "Resume ID"
// @Param        cat_id  path      integer  true  "Kategoriya ID"
// @Success      200     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Failure      404     {object}  map[string]string
// @Router       /resumes/{id}/categories/{cat_id} [delete]
func (h *resumeHandler) RemoveCategory(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil || id <= 0 {
		helper.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}
	catID, err := strconv.ParseInt(ps.ByName("cat_id"), 10, 64)
	if err != nil || catID <= 0 {
		helper.WriteError(w, http.StatusBadRequest, "invalid cat_id")
		return
	}
	if err := h.service.RemoveCategory(r.Context(), id, catID); err != nil {
		helper.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "category removed"})
}
