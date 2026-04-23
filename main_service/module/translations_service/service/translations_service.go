package translation_service

import (
	"context"
	"encoding/json"
	"fmt"
	translation_dto "main_service/module/translations_service/dto"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ─── Interface ───────────────────────────────────────────────────────────────

type TranslationService interface {
	Create(ctx context.Context, req translation_dto.CreateTranslationRequest) (*translation_dto.TranslationResponse, error)
	GetByID(ctx context.Context, id int64) (*translation_dto.TranslationResponse, error)
	Update(ctx context.Context, id int64, req translation_dto.UpdateTranslationRequest) (*translation_dto.TranslationResponse, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, f translation_dto.TranslationFilter, afterID int64, limit int) ([]*translation_dto.TranslationResponse, bool, error)
	// GetTranslation — frontend uchun: key+lang bo'yicha qidiradi, fallback logikasi bilan qaytaradi
	GetTranslation(ctx context.Context, slug, lang string) string
}

// ─── Implementation ──────────────────────────────────────────────────────────

type translationService struct {
	db *pgxpool.Pool
}

func NewTranslationService(db *pgxpool.Pool) TranslationService {
	return &translationService{db: db}
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func validateName(name map[string]string) error {
	if len(name) == 0 {
		return fmt.Errorf("name is required")
	}

	defaultValue := strings.TrimSpace(name["default"])
	if defaultValue == "" {
		return fmt.Errorf("name.default is required")
	}

	for key, value := range name {
		name[key] = strings.TrimSpace(value)
	}

	return nil
}

func (s *translationService) getActiveLangCodes(ctx context.Context) ([]string, error) {
	rows, err := s.db.Query(ctx, `
		SELECT name FROM languages WHERE is_active = true AND deleted_at IS NULL
	`)
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

// normalizeName — barcha kalitlarni kichik harfga keltiradi
func normalizeName(name map[string]string) map[string]string {
	result := make(map[string]string, len(name))
	for k, v := range name {
		result[strings.ToLower(k)] = v
	}
	return result
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

// ─── GetTranslation ──────────────────────────────────────────────────────────
// Laravel getTranslation() mantiqiga mos:
//  1. slug bo'yicha qidiradi
//  2. name[lang] bo'lsa — shu qiymatni qaytaradi
//  3. bo'lmasa name["default"] ni qaytaradi
//  4. default ham bo'lmasa — slug (keyni o'zini) qaytaradi

func (s *translationService) GetTranslation(ctx context.Context, slug, lang string) string {
	var nameBytes []byte
	err := s.db.QueryRow(ctx, `
		SELECT name FROM translations
		WHERE slug = $1 AND deleted_at IS NULL AND is_active = true
	`, slug).Scan(&nameBytes)
	if err != nil {
		// tarjima topilmasa — keyni o'zini qaytaradi
		return slug
	}

	var name map[string]string
	if err := json.Unmarshal(nameBytes, &name); err != nil {
		return slug
	}

	if v, ok := name[strings.ToLower(lang)]; ok && strings.TrimSpace(v) != "" {
		return v
	}

	if v, ok := name["default"]; ok && strings.TrimSpace(v) != "" {
		return v
	}

	return slug
}

// ─── Create ──────────────────────────────────────────────────────────────────

func (s *translationService) Create(ctx context.Context, req translation_dto.CreateTranslationRequest) (*translation_dto.TranslationResponse, error) {
	req.Slug = strings.TrimSpace(req.Slug)

	if err := validateName(req.Name); err != nil {
		return nil, err
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	nameJSON, err := json.Marshal(normalizeName(req.Name))
	if err != nil {
		return nil, err
	}

	var id int64
	err = s.db.QueryRow(ctx, `
		INSERT INTO translations (slug, name, is_active)
		VALUES ($1, $2, $3)
		RETURNING id
	`, req.Slug, string(nameJSON), isActive).Scan(&id)
	if err != nil {
		return nil, err
	}

	return s.GetByID(ctx, id)
}

// ─── GetByID ─────────────────────────────────────────────────────────────────

func (s *translationService) GetByID(ctx context.Context, id int64) (*translation_dto.TranslationResponse, error) {
	var r translation_dto.TranslationResponse
	var nameBytes []byte

	err := s.db.QueryRow(ctx, `
		SELECT t.id, t.slug, t.name,
		       t.is_active, t.created_at, t.updated_at
		FROM translations t
		WHERE t.id = $1 AND t.deleted_at IS NULL
	`, id).Scan(
		&r.ID, &r.Slug, &nameBytes,
		&r.IsActive, &r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	var fullName map[string]string
	if len(nameBytes) > 0 {
		if err := json.Unmarshal(nameBytes, &fullName); err != nil {
			return nil, err
		}
	} else {
		fullName = map[string]string{}
	}

	activeLangs, err := s.getActiveLangCodes(ctx)
	if err != nil {
		return nil, err
	}
	r.Name = filterName(fullName, activeLangs)

	return &r, nil
}

// ─── Update ──────────────────────────────────────────────────────────────────

func (s *translationService) Update(ctx context.Context, id int64, req translation_dto.UpdateTranslationRequest) (*translation_dto.TranslationResponse, error) {
	setClauses := []string{"updated_at = NOW()"}
	args := []interface{}{}
	idx := 1

	if req.Slug != nil {
		trimmed := strings.TrimSpace(*req.Slug)
		setClauses = append(setClauses, fmt.Sprintf("slug = $%d", idx))
		args = append(args, trimmed)
		idx++
	}

	if req.Name != nil {
		if err := validateName(*req.Name); err != nil {
			return nil, err
		}

		nameJSON, err := json.Marshal(normalizeName(*req.Name))
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
		UPDATE translations SET %s
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

func (s *translationService) Delete(ctx context.Context, id int64) error {
	tag, err := s.db.Exec(ctx, `
		UPDATE translations
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`, time.Now(), id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("translation not found")
	}

	return nil
}

// ─── List ────────────────────────────────────────────────────────────────────

func (s *translationService) List(ctx context.Context, f translation_dto.TranslationFilter, afterID int64, limit int) ([]*translation_dto.TranslationResponse, bool, error) {
	args := []interface{}{}
	conditions := []string{"t.deleted_at IS NULL"}
	idx := 1

	if afterID > 0 {
		conditions = append(conditions, fmt.Sprintf("t.id < $%d", idx))
		args = append(args, afterID)
		idx++
	}

	if f.Slug != "" {
		conditions = append(conditions, fmt.Sprintf("t.slug ILIKE $%d", idx))
		args = append(args, "%"+f.Slug+"%")
		idx++
	}

	if f.Name != "" {
		conditions = append(conditions, fmt.Sprintf(`(
				t.name->>'default' ILIKE $%d OR
				t.name->>'uz' ILIKE $%d OR
				t.name->>'ru' ILIKE $%d OR
				t.name->>'en' ILIKE $%d
			)`, idx, idx, idx, idx))
		args = append(args, "%"+f.Name+"%")
		idx++
	}

	if f.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("t.is_active = $%d", idx))
		args = append(args, *f.IsActive)
		idx++
	}

	args = append(args, limit+1)

	query := fmt.Sprintf(`
		SELECT t.id, t.slug, t.name,
		       t.is_active, t.created_at, t.updated_at
		FROM translations t
		WHERE %s
		ORDER BY t.id DESC
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

	var items []*translation_dto.TranslationResponse
	for rows.Next() {
		var r translation_dto.TranslationResponse
		var nameBytes []byte

		if err := rows.Scan(
			&r.ID, &r.Slug, &nameBytes,
			&r.IsActive, &r.CreatedAt, &r.UpdatedAt,
		); err != nil {
			return nil, false, err
		}

		var fullName map[string]string
		if len(nameBytes) > 0 {
			if err := json.Unmarshal(nameBytes, &fullName); err != nil {
				return nil, false, err
			}
		} else {
			fullName = map[string]string{}
		}
		r.Name = filterName(fullName, activeLangs)

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
