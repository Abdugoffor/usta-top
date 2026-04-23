package translation_handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"main_service/helper"
	"main_service/middleware"
	translation_dto "main_service/module/translations_service/dto"
	translation_service "main_service/module/translations_service/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
)

type translationHandler struct {
	service translation_service.TranslationService
}

func NewTranslationHandler(router *httprouter.Router, group string, db *pgxpool.Pool) {
	h := &translationHandler{service: translation_service.NewTranslationService(db)}

	routes := group + "/translations"
	{
		router.POST(routes, middleware.CheckRole(h.Create))
		router.GET(routes, h.List)
		router.GET(routes+"/:id", h.GetByID)
		router.PUT(routes+"/:id", middleware.CheckRole(h.Update))
		router.DELETE(routes+"/:id", middleware.CheckRole(h.Delete))
	}

	// Frontend uchun public endpoint: GET /api/v1/t?key=slug&lang=uz
	router.GET(group+"/t", h.GetTranslation)
}

// GetTranslation godoc
// @Summary      Frontend uchun tarjima olish
// @Description  key (slug) va lang bo'yicha tarjimani qaytaradi. lang yo'q bo'lsa default ishlatiladi. Agar tarjima topilmasa key qaytariladi.
// @Tags         Translations
// @Produce      json
// @Param        key   query     string  true   "Translation slug (kalit)"
// @Param        lang  query     string  false  "Til kodi (uz, ru, en ...)"
// @Success      200   {object}  translation_dto.TranslationKeyResponse
// @Failure      400   {object}  map[string]string
// @Router       /t [get]
func (h *translationHandler) GetTranslation(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	key := strings.TrimSpace(r.URL.Query().Get("key"))
	if key == "" {
		helper.WriteError(w, http.StatusBadRequest, "key is required")
		return
	}

	lang := strings.TrimSpace(r.URL.Query().Get("lang"))
	if lang == "" {
		lang = "default"
	}

	value := h.service.GetTranslation(r.Context(), key, lang)

	helper.WriteJSON(w, http.StatusOK, translation_dto.TranslationKeyResponse{
		Key:   key,
		Value: value,
		Lang:  lang,
	})
}

// Create godoc
// @Summary      Yangi translation yaratish
// @Tags         Translations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      translation_dto.CreateTranslationRequest  true  "Translation ma'lumotlari"
// @Success      201   {object}  translation_dto.TranslationResponse
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /translations [post]
func (h *translationHandler) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req translation_dto.CreateTranslationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	req.Slug = strings.TrimSpace(req.Slug)

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
// @Summary      Translationlar ro'yxati
// @Tags         Translations
// @Produce      json
// @Param        slug        query     string  false  "Slug bo'yicha filter"
// @Param        name        query     string  false  "Name bo'yicha filter"
// @Param        is_active   query     boolean false  "Faol/faolsiz"
// @Param        page        query     integer false  "Sahifa" default(1)
// @Param        limit       query     integer false  "Limit" default(10)
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /translations [get]
func (h *translationHandler) List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	q := r.URL.Query()
	afterID, limit := helper.ParseCursorPage(r)

	slug := q.Get("slug")
	if len(slug) > 150 {
		slug = slug[:150]
	}

	name := q.Get("name")
	if len(name) > 100 {
		name = name[:100]
	}

	f := translation_dto.TranslationFilter{
		Slug: slug,
		Name: name,
	}

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

	helper.WriteJSON(w, http.StatusOK, map[string]any{
		"data": items,
		"meta": helper.NewCursorMeta(limit, hasMore, lastID, 0),
	})
}

// GetByID godoc
// @Summary      Translationni ID bo'yicha olish
// @Tags         Translations
// @Produce      json
// @Param        id   path      integer  true  "Translation ID"
// @Success      200  {object}  translation_dto.TranslationResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /translations/{id} [get]
func (h *translationHandler) GetByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil || id <= 0 {
		helper.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	resp, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		helper.WriteError(w, http.StatusNotFound, "translation not found")
		return
	}

	helper.WriteJSON(w, http.StatusOK, resp)
}

// Update godoc
// @Summary      Translationni yangilash
// @Tags         Translations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      integer                                true  "Translation ID"
// @Param        body  body      translation_dto.UpdateTranslationRequest  true  "Yangi ma'lumotlar"
// @Success      200   {object}  translation_dto.TranslationResponse
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Router       /translations/{id} [put]
func (h *translationHandler) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil || id <= 0 {
		helper.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req translation_dto.UpdateTranslationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if req.Slug != nil {
		trimmed := strings.TrimSpace(*req.Slug)
		req.Slug = &trimmed
	}

	if errs := helper.ValidateStruct(req); errs != nil {
		helper.WriteValidation(w, errs)
		return
	}

	resp, err := h.service.Update(r.Context(), id, req)
	if err != nil {
		helper.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, resp)
}

// Delete godoc
// @Summary      Translationni o'chirish
// @Tags         Translations
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      integer  true  "Translation ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /translations/{id} [delete]
func (h *translationHandler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
