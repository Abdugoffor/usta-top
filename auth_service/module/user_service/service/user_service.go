package user_service

import (
	"auth_service/helper"
	user_dto "auth_service/module/user_service/dto"
	user_model "auth_service/module/user_service/model"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// ─── Interface ───────────────────────────────────────────────────────────────

type UserService interface {
	Register(ctx context.Context, req user_dto.RegisterRequest) (*user_dto.AuthResponse, error)
	Login(ctx context.Context, req user_dto.LoginRequest) (*user_dto.AuthResponse, error)
}

// ─── Implementation ──────────────────────────────────────────────────────────

type userService struct {
	db *pgxpool.Pool
}

func NewUserService(db *pgxpool.Pool) UserService {
	return &userService{db: db}
}

// ─── Register ────────────────────────────────────────────────────────────────

func (s *userService) Register(ctx context.Context, req user_dto.RegisterRequest) (*user_dto.AuthResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	role := "user"
	if req.Role != "" {
		role = req.Role
	}

	var (
		u        user_model.User
		password string
	)
	err = s.db.QueryRow(ctx, `
		INSERT INTO users (full_name, phone, password, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, full_name, phone, password, role, is_active, created_at, updated_at, deleted_at
	`, req.FullName, req.Phone, string(hash), role).Scan(
		&u.ID, &u.FullName, &u.Phone, &password,
		&u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("phone allaqachon ro'yxatdan o'tgan yoki xatolik: %v", err)
	}

	token, err := helper.GenerateToken(int(u.ID), u.Role)
	if err != nil {
		return nil, err
	}

	return &user_dto.AuthResponse{Token: token, User: u}, nil
}

// ─── Login ───────────────────────────────────────────────────────────────────

func (s *userService) Login(ctx context.Context, req user_dto.LoginRequest) (*user_dto.AuthResponse, error) {
	var (
		u        user_model.User
		password string
	)
	err := s.db.QueryRow(ctx, `
		SELECT id, full_name, phone, password, role, is_active, created_at, updated_at, deleted_at
		FROM users
		WHERE phone = $1 AND deleted_at IS NULL
	`, req.Phone).Scan(
		&u.ID, &u.FullName, &u.Phone, &password,
		&u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("phone yoki parol noto'g'ri")
	}

	if !u.IsActive {
		return nil, fmt.Errorf("account faol emas")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("phone yoki parol noto'g'ri")
	}

	token, err := helper.GenerateToken(int(u.ID), u.Role)
	if err != nil {
		return nil, err
	}

	return &user_dto.AuthResponse{Token: token, User: u}, nil
}
