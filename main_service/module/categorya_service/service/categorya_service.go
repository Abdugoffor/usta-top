package categorya_service

import (
	"context"
	"fmt"
	categorya_dto "main_service/module/categorya_service/dto"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ─── Interface ───────────────────────────────────────────────────────────────

type CategoryService interface {
	Create(ctx context.Context, req categorya_dto.CreateCategoryRequest) (*categorya_dto.CategoryResponse, error)
	GetByID(ctx context.Context, id int64) (*categorya_dto.CategoryResponse, error)
	Update(ctx context.Context, id int64, req categorya_dto.UpdateCategoryRequest) (*categorya_dto.CategoryResponse, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, f categorya_dto.CategoryFilter, afterID int64, limit int) ([]*categorya_dto.CategoryResponse, bool, error)
}

// ─── Implementation ──────────────────────────────────────────────────────────

type categoryService struct {
	db *pgxpool.Pool
}

func NewCategoryService(db *pgxpool.Pool) CategoryService {
	return &categoryService{db: db}
}

var validSortCols = map[string]string{
	"id":         "c.id",
	"name":       "c.name",
	"is_active":  "c.is_active",
	"created_at": "c.created_at",
	"updated_at": "c.updated_at",
}

// ─── Create ──────────────────────────────────────────────────────────────────

func (s *categoryService) Create(ctx context.Context, req categorya_dto.CreateCategoryRequest) (*categorya_dto.CategoryResponse, error) {
	isActive := true
	{
		if req.IsActive != nil {
			isActive = *req.IsActive
		}
	}

	var r categorya_dto.CategoryResponse
	err := s.db.QueryRow(ctx, `
		INSERT INTO categories (name, is_active)
		VALUES ($1, $2)
		RETURNING id, name, is_active, created_at, updated_at, deleted_at
	`, req.Name, isActive).Scan(
		&r.ID, &r.Name, &r.IsActive, &r.CreatedAt, &r.UpdatedAt, &r.DeletedAt,
	)
	return &r, err
}

// ─── GetByID ─────────────────────────────────────────────────────────────────

func (s *categoryService) GetByID(ctx context.Context, id int64) (*categorya_dto.CategoryResponse, error) {
	var r categorya_dto.CategoryResponse
	err := s.db.QueryRow(ctx, `
		SELECT id, name, is_active, created_at, updated_at, deleted_at
		FROM categories
		WHERE id = $1 AND deleted_at IS NULL
	`, id).Scan(&r.ID, &r.Name, &r.IsActive, &r.CreatedAt, &r.UpdatedAt, &r.DeletedAt)
	return &r, err
}

// ─── Update ──────────────────────────────────────────────────────────────────

func (s *categoryService) Update(ctx context.Context, id int64, req categorya_dto.UpdateCategoryRequest) (*categorya_dto.CategoryResponse, error) {
	setClauses := []string{"updated_at = NOW()"}
	args := []interface{}{}
	idx := 1

	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", idx))
		args = append(args, *req.Name)
		idx++
	}
	if req.IsActive != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_active = $%d", idx))
		args = append(args, *req.IsActive)
		idx++
	}

	args = append(args, id)
	query := fmt.Sprintf(`
		UPDATE categories SET %s
		WHERE id = $%d AND deleted_at IS NULL
		RETURNING id, name, is_active, created_at, updated_at, deleted_at
	`, strings.Join(setClauses, ", "), idx)

	var r categorya_dto.CategoryResponse
	err := s.db.QueryRow(ctx, query, args...).Scan(
		&r.ID, &r.Name, &r.IsActive, &r.CreatedAt, &r.UpdatedAt, &r.DeletedAt,
	)
	return &r, err
}

// ─── Delete ──────────────────────────────────────────────────────────────────

func (s *categoryService) Delete(ctx context.Context, id int64) error {
	tag, err := s.db.Exec(ctx, `
		UPDATE categories SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL
	`, time.Now(), id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("category not found")
	}
	return nil
}

// ─── List ────────────────────────────────────────────────────────────────────

func (s *categoryService) List(ctx context.Context, f categorya_dto.CategoryFilter, afterID int64, limit int) ([]*categorya_dto.CategoryResponse, bool, error) {
	args := []interface{}{}
	conditions := []string{"c.deleted_at IS NULL"}
	idx := 1

	if afterID > 0 {
		conditions = append(conditions, fmt.Sprintf("c.id < $%d", idx))
		args = append(args, afterID)
		idx++
	}
	if f.Name != "" {
		conditions = append(conditions, fmt.Sprintf("c.name ILIKE $%d", idx))
		args = append(args, "%"+f.Name+"%")
		idx++
	}
	if f.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("c.is_active = $%d", idx))
		args = append(args, *f.IsActive)
		idx++
	}

	args = append(args, limit+1)
	query := fmt.Sprintf(`
		SELECT c.id, c.name, c.is_active, c.created_at, c.updated_at, c.deleted_at
		FROM categories c
		WHERE %s
		ORDER BY c.id DESC
		LIMIT $%d
	`, strings.Join(conditions, " AND "), idx)

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	var items []*categorya_dto.CategoryResponse
	for rows.Next() {
		var r categorya_dto.CategoryResponse
		if err := rows.Scan(&r.ID, &r.Name, &r.IsActive, &r.CreatedAt, &r.UpdatedAt, &r.DeletedAt); err != nil {
			return nil, false, err
		}
		items = append(items, &r)
	}
	if err := rows.Err(); err != nil {
		return nil, false, err
	}

	hasMore := len(items) > limit
	if hasMore {
		items = items[:limit]
	}
	return items, hasMore, nil
}
