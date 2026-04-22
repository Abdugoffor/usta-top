package user_service

import (
	"context"
	"fmt"
	"main_service/helper"
	user_dto "main_service/module/user_service/dto"
	user_model "main_service/module/user_service/model"

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
	{
		if err != nil {
			return nil, err
		}
	}

	role := "user"
	{
		if req.Role != "" {
			role = req.Role
		}
	}

	var (
		user     user_model.User
		password string
	)
	err = s.db.QueryRow(ctx, `
		INSERT INTO users (full_name, phone, password, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, full_name, phone, password, role, is_active, created_at, updated_at, deleted_at
	`, req.FullName, req.Phone, string(hash), role).Scan(
		&user.ID, &user.FullName, &user.Phone, &password,
		&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("phone allaqachon ro'yxatdan o'tgan")
	}

	token, err := helper.GenerateToken(int(user.ID), user.Role)
	{
		if err != nil {
			return nil, err
		}
	}

	return &user_dto.AuthResponse{Token: token, User: user}, nil
}

// ─── Login ───────────────────────────────────────────────────────────────────

func (s *userService) Login(ctx context.Context, req user_dto.LoginRequest) (*user_dto.AuthResponse, error) {
	var (
		user     user_model.User
		password string
	)

	err := s.db.QueryRow(ctx, `
		SELECT id, full_name, phone, password, role, is_active, created_at, updated_at, deleted_at
		FROM users
		WHERE phone = $1 AND deleted_at IS NULL
	`, req.Phone).Scan(
		&user.ID, &user.FullName, &user.Phone, &password,
		&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("phone yoki parol noto'g'ri")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("account faol emas")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("phone yoki parol noto'g'ri")
	}

	token, err := helper.GenerateToken(int(user.ID), user.Role)
	{
		if err != nil {
			return nil, err
		}
	}

	return &user_dto.AuthResponse{Token: token, User: user}, nil
}
