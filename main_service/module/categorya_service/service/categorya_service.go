package categorya_service

import (
	"context"
	"encoding/json"
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

// ─── Helpers ─────────────────────────────────────────────────────────────────

func validateName(name map[string]string) error {
	if len(name) == 0 {
		return fmt.Errorf("name is required")
	}
	if strings.TrimSpace(name["default"]) == "" {
		return fmt.Errorf("name.default is required")
	}
	for k, v := range name {
		name[k] = strings.TrimSpace(v)
	}
	return nil
}

func (s *categoryService) getActiveLangCodes(ctx context.Context) ([]string, error) {
	rows, err := s.db.Query(ctx, `SELECT name FROM languages WHERE is_active = true AND deleted_at IS NULL`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var codes []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		codes = append(codes, code)
	}
	return codes, rows.Err()
}

func filterName(name map[string]string, activeLangs []string) map[string]string {
	allowed := make(map[string]bool, len(activeLangs)+1)
	allowed["default"] = true
	for _, lang := range activeLangs {
		allowed[lang] = true
	}
	result := make(map[string]string)
	for k, v := range name {
		if allowed[k] {
			result[k] = v
		}
	}
	return result
}

func unmarshalName(b []byte) (map[string]string, error) {
	var m map[string]string
	if len(b) == 0 {
		return map[string]string{}, nil
	}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// ─── Create ──────────────────────────────────────────────────────────────────

func (s *categoryService) Create(ctx context.Context, req categorya_dto.CreateCategoryRequest) (*categorya_dto.CategoryResponse, error) {
	if err := validateName(req.Name); err != nil {
		return nil, err
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	nameJSON, err := json.Marshal(req.Name)
	if err != nil {
		return nil, err
	}

	var id int64
	err = s.db.QueryRow(ctx, `
		INSERT INTO categories (name, is_active)
		VALUES ($1, $2)
		RETURNING id
	`, string(nameJSON), isActive).Scan(&id)
	if err != nil {
		return nil, err
	}

	return s.GetByID(ctx, id)
}

// ─── GetByID ─────────────────────────────────────────────────────────────────

func (s *categoryService) GetByID(ctx context.Context, id int64) (*categorya_dto.CategoryResponse, error) {
	var r categorya_dto.CategoryResponse
	var nameBytes []byte
	var deletedAt *time.Time

	err := s.db.QueryRow(ctx, `
		SELECT id, name, is_active, created_at, updated_at, deleted_at
		FROM categories
		WHERE id = $1 AND deleted_at IS NULL
	`, id).Scan(&r.ID, &nameBytes, &r.IsActive, &r.CreatedAt, &r.UpdatedAt, &deletedAt)
	if err != nil {
		return nil, err
	}

	fullName, err := unmarshalName(nameBytes)
	if err != nil {
		return nil, err
	}

	activeLangs, err := s.getActiveLangCodes(ctx)
	if err != nil {
		return nil, err
	}
	r.Name = filterName(fullName, activeLangs)
	r.DeletedAt = deletedAt

	return &r, nil
}

// ─── Update ──────────────────────────────────────────────────────────────────

func (s *categoryService) Update(ctx context.Context, id int64, req categorya_dto.UpdateCategoryRequest) (*categorya_dto.CategoryResponse, error) {
	setClauses := []string{"updated_at = NOW()"}
	args := []interface{}{}
	idx := 1

	if req.Name != nil {
		if err := validateName(*req.Name); err != nil {
			return nil, err
		}
		nameJSON, err := json.Marshal(*req.Name)
		if err != nil {
			return nil, err
		}
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", idx))
		args = append(args, string(nameJSON))
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
		RETURNING id
	`, strings.Join(setClauses, ", "), idx)

	var retID int64
	if err := s.db.QueryRow(ctx, query, args...).Scan(&retID); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, retID)
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
		conditions = append(conditions, fmt.Sprintf(`(
			c.name->>'default' ILIKE $%d OR
			c.name->>'uz'      ILIKE $%d OR
			c.name->>'ru'      ILIKE $%d OR
			c.name->>'en'      ILIKE $%d
		)`, idx, idx, idx, idx))
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

	activeLangs, err := s.getActiveLangCodes(ctx)
	if err != nil {
		return nil, false, err
	}

	var items []*categorya_dto.CategoryResponse
	for rows.Next() {
		var r categorya_dto.CategoryResponse
		var nameBytes []byte
		var deletedAt *time.Time

		if err := rows.Scan(&r.ID, &nameBytes, &r.IsActive, &r.CreatedAt, &r.UpdatedAt, &deletedAt); err != nil {
			return nil, false, err
		}

		fullName, err := unmarshalName(nameBytes)
		if err != nil {
			return nil, false, err
		}
		r.Name = filterName(fullName, activeLangs)
		r.DeletedAt = deletedAt
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
