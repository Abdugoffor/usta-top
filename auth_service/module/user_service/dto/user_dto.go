package user_dto

import user_model "auth_service/module/user_service/model"

type RegisterRequest struct {
	FullName string `json:"full_name" validate:"required,min=2,max=255"`
	Phone    string `json:"phone"     validate:"required,min=9,max=20"`
	Password string `json:"password"  validate:"required,min=6"`
	Role     string `json:"role"      validate:"omitempty,oneof=user employer"`
}

type LoginRequest struct {
	Phone    string `json:"phone"    validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string          `json:"token"`
	User  user_model.User `json:"user"`
}
