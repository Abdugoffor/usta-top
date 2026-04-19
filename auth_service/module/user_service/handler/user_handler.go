package user_handler

import (
	"auth_service/helper"
	user_dto "auth_service/module/user_service/dto"
	user_service "auth_service/module/user_service/service"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type UserHandler struct {
	svc user_service.UserService
}

func NewUserHandler(svc user_service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// POST /api/v1/auth/register
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req user_dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if errs := helper.ValidateStruct(req); errs != nil {
		helper.WriteValidation(w, errs)
		return
	}

	resp, err := h.svc.Register(r.Context(), req)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusCreated, resp)
}

// POST /api/v1/auth/login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req user_dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if errs := helper.ValidateStruct(req); errs != nil {
		helper.WriteValidation(w, errs)
		return
	}

	resp, err := h.svc.Login(r.Context(), req)
	if err != nil {
		helper.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, resp)
}
