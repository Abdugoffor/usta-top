package vacancy_handler

import (
	"encoding/json"
	"main_service/helper"
	"main_service/middleware"
	vacancy_dto "main_service/module/vacancy_service/dto"
	vacancy_service "main_service/module/vacancy_service/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
)

type vacancyHandler struct {
	service vacancy_service.VacancyService
}

func NewVacancyHandler(router *httprouter.Router, group string, db *pgxpool.Pool) {
	h := &vacancyHandler{service: vacancy_service.NewVacancyService(db)}

	routes := group + "/vacancies"
	{
		router.POST(routes, middleware.CheckRole(h.Create))
		router.GET(routes, h.List)
		router.GET(routes+"/:slug", h.GetBySlug)
		router.PUT(routes+"/:id", middleware.CheckRole(h.Update))
		router.DELETE(routes+"/:id", middleware.CheckRole(h.Delete))
	}
}

// Create godoc
// @Summary      Yangi vakansiya yaratish
// @Tags         Vacancies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      vacancy_dto.CreateVacancyRequest  true  "Vakansiya ma'lumotlari"
// @Success      201   {object}  vacancy_dto.VacancyResponse
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /vacancies [post]
func (h *vacancyHandler) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID := middleware.GetUserID(r)
	if userID == 0 {
		helper.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req vacancy_dto.CreateVacancyRequest
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
		helper.WriteInternalError(w, err)
		return
	}

	helper.WriteJSON(w, http.StatusCreated, resp)
}

// List godoc
// @Summary      Vakansiyalar ro'yxati
// @Tags         Vacancies
// @Produce      json
// @Param        name         query     string  false  "Nomi"
// @Param        title        query     string  false  "Sarlavha"
// @Param        user_id      query     integer false  "Foydalanuvchi ID"
// @Param        region_id    query     integer false  "Region ID"
// @Param        district_id  query     integer false  "Tuman ID"
// @Param        mahalla_id   query     integer false  "Mahalla ID"
// @Param        category_id  query     integer false  "Kategoriya ID"
// @Param        category_ids query     string  false  "Kategoriya ID lar (vergul bilan)"
// @Param        is_active    query     boolean false  "Faol/faolsiz"
// @Param        min_price    query     integer false  "Minimal narx"
// @Param        max_price    query     integer false  "Maksimal narx"
// @Param        page         query     integer false  "Sahifa" default(1)
// @Param        limit        query     integer false  "Limit" default(10)
// @Param        sort_by      query     string  false  "Saralash maydoni"
// @Param        sort_order   query     string  false  "asc yoki desc"
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /vacancies [get]
func (h *vacancyHandler) List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	q := r.URL.Query()
	cursor, limit := helper.ParseCursorPayload(r)

	name := q.Get("name")
	if len(name) > 100 {
		name = name[:100]
	}
	title := q.Get("title")
	if len(title) > 100 {
		title = title[:100]
	}
	search := q.Get("search")
	if len(search) > 100 {
		search = search[:100]
	}

	f := vacancy_dto.VacancyFilter{
		Name:      name,
		Title:     title,
		Search:    search,
		SortBy:    strings.TrimSpace(q.Get("sort_by")),
		SortOrder: strings.TrimSpace(q.Get("sort_order")),
	}

	if catIDs := q.Get("category_ids"); catIDs != "" {
		for _, s := range strings.Split(catIDs, ",") {
			if len(f.CategoryIDs) >= 20 {
				break
			}
			s = strings.TrimSpace(s)
			if n, err := strconv.ParseInt(s, 10, 64); err == nil && n > 0 {
				f.CategoryIDs = append(f.CategoryIDs, n)
			}
		}
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
	if v := q.Get("category_id"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			f.CategoryID = &n
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

	items, hasMore, total, err := h.service.List(r.Context(), f, cursor, limit)
	if err != nil {
		helper.WriteInternalError(w, err)
		return
	}

	var lastID int64
	var lastValue string
	if len(items) > 0 {
		lastID = items[len(items)-1].ID
		if f.SortBy == "price" && items[len(items)-1].Price != nil {
			lastValue = strconv.FormatInt(*items[len(items)-1].Price, 10)
		}
	}

	helper.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"data": items,
		"meta": helper.NewCursorMetaWithValue(limit, hasMore, lastID, lastValue, total),
	})
}

// GetBySlug godoc
// @Summary      Vakansiyani slug bo'yicha olish
// @Tags         Vacancies
// @Produce      json
// @Param        slug  path      string  true  "Vakansiya slug"
// @Success      200   {object}  vacancy_dto.VacancyResponse
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Router       /vacancies/{slug} [get]
func (h *vacancyHandler) GetBySlug(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	slug := ps.ByName("slug")
	if slug == "" {
		helper.WriteError(w, http.StatusBadRequest, "invalid slug")
		return
	}

	resp, err := h.service.GetBySlug(r.Context(), slug)
	if err != nil {
		helper.WriteError(w, http.StatusNotFound, "vacancy not found")
		return
	}

	helper.WriteJSON(w, http.StatusOK, resp)
}

// Update godoc
// @Summary      Vakansiyani yangilash
// @Tags         Vacancies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      integer                           true  "Vakansiya ID"
// @Param        body  body      vacancy_dto.UpdateVacancyRequest  true  "Yangi ma'lumotlar"
// @Success      200   {object}  vacancy_dto.VacancyResponse
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Router       /vacancies/{id} [put]
func (h *vacancyHandler) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	var req vacancy_dto.UpdateVacancyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		req.Name = &trimmed
	}
	if req.Title != nil {
		trimmed := strings.TrimSpace(*req.Title)
		req.Title = &trimmed
	}
	if req.Adress != nil {
		trimmed := strings.TrimSpace(*req.Adress)
		req.Adress = &trimmed
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
// @Summary      Vakansiyani o'chirish
// @Tags         Vacancies
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      integer  true  "Vakansiya ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /vacancies/{id} [delete]
func (h *vacancyHandler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
