package language_service

import (
	"context"
	"fmt"
	language_dto "main_service/module/language_service/dto"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ─── Interface ───────────────────────────────────────────────────────────────

type LanguageService interface {
	Create(ctx context.Context, req language_dto.CreateLanguageRequest) (*language_dto.LanguageResponse, error)
	GetByID(ctx context.Context, id int64) (*language_dto.LanguageResponse, error)
	Update(ctx context.Context, id int64, req language_dto.UpdateLanguageRequest) (*language_dto.LanguageResponse, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, f language_dto.LanguageFilter, afterID int64, limit int) ([]*language_dto.LanguageResponse, bool, error)
}

// ─── Implementation ──────────────────────────────────────────────────────────

type languageService struct {
	db *pgxpool.Pool
}

func NewLanguageService(db *pgxpool.Pool) LanguageService {
	return &languageService{db: db}
}

var validSortCols = map[string]string{
	"id":         "l.id",
	"name":       "l.name",
	"is_active":  "l.is_active",
	"created_at": "l.created_at",
	"updated_at": "l.updated_at",
}

// ─── Create ──────────────────────────────────────────────────────────────────

func (s *languageService) Create(ctx context.Context, req language_dto.CreateLanguageRequest) (*language_dto.LanguageResponse, error) {
	isActive := true

	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	var id int64
	err := s.db.QueryRow(ctx, `
		INSERT INTO languages (name, description, is_active)
		VALUES ($1, $2, $3)
		RETURNING id
	`, strings.ToLower(strings.TrimSpace(req.Name)), req.Description, isActive).Scan(&id)
	if err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

// ─── GetByID ─────────────────────────────────────────────────────────────────

func (s *languageService) GetByID(ctx context.Context, id int64) (*language_dto.LanguageResponse, error) {
	var r language_dto.LanguageResponse

	err := s.db.QueryRow(ctx, `
		SELECT l.id, l.name, COALESCE(l.description, ''),
		       l.is_active, l.created_at, l.updated_at, l.deleted_at
		FROM languages l
		WHERE l.id = $1 AND l.deleted_at IS NULL
	`, id).Scan(
		&r.ID, &r.Name, &r.Description,
		&r.IsActive, &r.CreatedAt, &r.UpdatedAt, &r.DeletedAt,
	)
	return &r, err
}

// ─── Update ──────────────────────────────────────────────────────────────────

func (s *languageService) Update(ctx context.Context, id int64, req language_dto.UpdateLanguageRequest) (*language_dto.LanguageResponse, error) {
	setClauses := []string{"updated_at = NOW()"}
	args := []interface{}{}
	idx := 1

	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", idx))
		args = append(args, strings.ToLower(strings.TrimSpace(*req.Name)))
		idx++
	}

	if req.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", idx))
		args = append(args, *req.Description)
		idx++
	}

	if req.IsActive != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_active = $%d", idx))
		args = append(args, *req.IsActive)
		idx++
	}

	args = append(args, id)
	query := fmt.Sprintf(`
		UPDATE languages SET %s
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

func (s *languageService) Delete(ctx context.Context, id int64) error {
	tag, err := s.db.Exec(ctx, `
		UPDATE languages SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL
	`, time.Now(), id)

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("language not found")
	}

	return nil
}

// ─── List ────────────────────────────────────────────────────────────────────

func (s *languageService) List(ctx context.Context, f language_dto.LanguageFilter, afterID int64, limit int) ([]*language_dto.LanguageResponse, bool, error) {
	args := []interface{}{}
	conditions := []string{"l.deleted_at IS NULL"}
	idx := 1

	if afterID > 0 {
		conditions = append(conditions, fmt.Sprintf("l.id < $%d", idx))
		args = append(args, afterID)
		idx++
	}
	if f.Name != "" {
		conditions = append(conditions, fmt.Sprintf("l.name ILIKE $%d", idx))
		args = append(args, "%"+f.Name+"%")
		idx++
	}
	if f.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("l.is_active = $%d", idx))
		args = append(args, *f.IsActive)
		idx++
	}

	args = append(args, limit+1)
	query := fmt.Sprintf(`
		SELECT l.id, l.name, COALESCE(l.description, ''),
		       l.is_active, l.created_at, l.updated_at, l.deleted_at
		FROM languages l
		WHERE %s
		ORDER BY l.id DESC
		LIMIT $%d
	`, strings.Join(conditions, " AND "), idx)

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	var items []*language_dto.LanguageResponse
	for rows.Next() {
		var r language_dto.LanguageResponse
		if err := rows.Scan(
			&r.ID, &r.Name, &r.Description,
			&r.IsActive, &r.CreatedAt, &r.UpdatedAt, &r.DeletedAt,
		); err != nil {
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
