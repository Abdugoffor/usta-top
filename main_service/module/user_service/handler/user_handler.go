package user_handler

import (
	"encoding/json"
	"main_service/helper"
	user_dto "main_service/module/user_service/dto"
	user_service "main_service/module/user_service/service"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
)

type userHandler struct {
	service user_service.UserService
}

func NewUserHandler(router *httprouter.Router, group string, db *pgxpool.Pool) {
	h := &userHandler{service: user_service.NewUserService(db)}
	routes := group + "/auth"
	{
		router.POST(routes+"/register", h.Register)
		router.POST(routes+"/login", h.Login)
	}
}

// POST /api/v1/auth/register
func (h *userHandler) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req user_dto.RegisterRequest
	{
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helper.WriteError(w, http.StatusBadRequest, "invalid JSON")
			return
		}
	}

	if errs := helper.ValidateStruct(req); errs != nil {
		helper.WriteValidation(w, errs)
		return
	}

	resp, err := h.service.Register(r.Context(), req)
	{
		if err != nil {
			helper.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	helper.WriteJSON(w, http.StatusCreated, resp)
}

// POST /api/v1/auth/login
func (h *userHandler) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req user_dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if errs := helper.ValidateStruct(req); errs != nil {
		helper.WriteValidation(w, errs)
		return
	}

	resp, err := h.service.Login(r.Context(), req)
	if err != nil {
		helper.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, resp)
}
