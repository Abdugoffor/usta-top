package comment_handler

import (
	"encoding/json"
	"main_service/helper"
	"main_service/middleware"
	comment_dto "main_service/module/comment_service/dto"
	comment_service "main_service/module/comment_service/service"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
)

type commentHandler struct {
	service comment_service.CommentService
}

func NewCommentHandler(router *httprouter.Router, group string, db *pgxpool.Pool) {
	h := &commentHandler{service: comment_service.NewCommentService(db)}

	routes := group + "/comments"
	{
		router.POST(routes, middleware.Auth(h.Create))
		router.GET(routes, h.List)
		router.GET(routes+"/:id", h.GetByID)
		router.PUT(routes+"/:id", middleware.Auth(h.Update))
		router.DELETE(routes+"/:id", middleware.Auth(h.Delete))
	}
}

// Create godoc
// @Summary      Yangi izoh qo'shish
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      comment_dto.CreateCommentRequest  true  "Izoh ma'lumotlari"
// @Success      201   {object}  comment_dto.CommentResponse
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      422   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /comments [post]
func (h *commentHandler) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID := middleware.GetUserID(r)
	if userID == 0 {
		helper.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req comment_dto.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if req.VakansiyaID == nil && req.ResumeID == nil {
		helper.WriteError(w, http.StatusUnprocessableEntity, "vakansiya_id or resume_id required")
		return
	}
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
// @Summary      Izohlar ro'yxati
// @Tags         Comments
// @Produce      json
// @Param        vakansiya_id  query     integer false  "Vakansiya ID"
// @Param        resume_id     query     integer false  "Resume ID"
// @Param        user_id       query     integer false  "Foydalanuvchi ID"
// @Param        type          query     string  false  "comment yoki review"
// @Param        page          query     integer false  "Sahifa" default(1)
// @Param        limit         query     integer false  "Limit" default(20)
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /comments [get]
func (h *commentHandler) List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	f := comment_dto.CommentFilter{Type: q.Get("type")}
	if v := q.Get("vakansiya_id"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			f.VakansiyaID = &n
		}
	}
	if v := q.Get("resume_id"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			f.ResumeID = &n
		}
	}
	if v := q.Get("user_id"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			f.UserID = &n
		}
	}

	items, total, err := h.service.List(r.Context(), f, page, limit)
	{
		if err != nil {
			helper.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	helper.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"data": items,
		"meta": helper.NewPageMeta(total, page, limit),
	})
}

// GetByID godoc
// @Summary      Izohni ID bo'yicha olish
// @Tags         Comments
// @Produce      json
// @Param        id   path      integer  true  "Izoh ID"
// @Success      200  {object}  comment_dto.CommentResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /comments/{id} [get]
func (h *commentHandler) GetByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil || id <= 0 {
		helper.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}
	resp, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		helper.WriteError(w, http.StatusNotFound, "comment not found")
		return
	}
	helper.WriteJSON(w, http.StatusOK, resp)
}

// Update godoc
// @Summary      Izohni yangilash
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      integer                           true  "Izoh ID"
// @Param        body  body      comment_dto.UpdateCommentRequest  true  "Yangi ma'lumotlar"
// @Success      200   {object}  comment_dto.CommentResponse
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Router       /comments/{id} [put]
func (h *commentHandler) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	var req comment_dto.UpdateCommentRequest
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
// @Summary      Izohni o'chirish
// @Tags         Comments
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      integer  true  "Izoh ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /comments/{id} [delete]
func (h *commentHandler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
