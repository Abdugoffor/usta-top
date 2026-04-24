package language_service

import (
	"context"
	"fmt"
	"strings"
	"time"

	language_dto "main_service/module/language_service/dto"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LanguageService interface {
	All(ctx context.Context, f language_dto.LanguageFilter, afterID int64, limit int) ([]*language_dto.LanguageResponse, bool, error)
	Show(ctx context.Context, id int64) (*language_dto.LanguageResponse, error)
	Create(ctx context.Context, req language_dto.CreateLanguageRequest) (*language_dto.LanguageResponse, error)
	Update(ctx context.Context, id int64, req language_dto.UpdateLanguageRequest) (*language_dto.LanguageResponse, error)
	Delete(ctx context.Context, id int64) error
}

type languageService struct {
	db *pgxpool.Pool
}

func NewLanguageService(db *pgxpool.Pool) LanguageService {
	return &languageService{db: db}
}

func (service *languageService) All(ctx context.Context, f language_dto.LanguageFilter, afterID int64, limit int) ([]*language_dto.LanguageResponse, bool, error) {
	rows, err := service.db.Query(ctx, allQuery,
		f.Name,
		f.IsActive,
		afterID,
		limit+1,
	)

	if err != nil {
		return nil, false, err
	}

	defer rows.Close()

	items := make([]*language_dto.LanguageResponse, 0)

	for rows.Next() {
		var r language_dto.LanguageResponse
		if err := rows.Scan(
			&r.ID,
			&r.Name,
			&r.Description,
			&r.IsActive,
			&r.CreatedAt,
			&r.UpdatedAt,
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

func (service *languageService) Show(ctx context.Context, id int64) (*language_dto.LanguageResponse, error) {
	var r language_dto.LanguageResponse
	err := service.db.QueryRow(ctx, showQuery, id).Scan(
		&r.ID,
		&r.Name,
		&r.Description,
		&r.IsActive,
		&r.CreatedAt,
		&r.UpdatedAt,
	)
	return &r, err
}

func (service *languageService) Create(ctx context.Context, req language_dto.CreateLanguageRequest) (*language_dto.LanguageResponse, error) {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	var r language_dto.LanguageResponse
	err := service.db.QueryRow(ctx, createQuery,
		strings.ToLower(req.Name),
		req.Description,
		isActive,
	).Scan(
		&r.ID,
		&r.Name,
		&r.Description,
		&r.IsActive,
		&r.CreatedAt,
		&r.UpdatedAt,
	)
	return &r, err
}

func (service *languageService) Update(ctx context.Context, id int64, req language_dto.UpdateLanguageRequest) (*language_dto.LanguageResponse, error) {
	var name *string

	if req.Name != nil {
		lower := strings.ToLower(*req.Name)
		name = &lower
	}

	var r language_dto.LanguageResponse

	err := service.db.QueryRow(ctx, updateQuery,
		name,
		req.Description,
		req.IsActive,
		id,
	).Scan(
		&r.ID,
		&r.Name,
		&r.Description,
		&r.IsActive,
		&r.CreatedAt,
		&r.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("language not found")
	}
	return &r, nil
}

func (service *languageService) Delete(ctx context.Context, id int64) error {
	tag, err := service.db.Exec(ctx, deleteQuery, time.Now(), id)

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("language not found")
	}

	return nil
}
