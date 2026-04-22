package country_service

import (
	"context"
	"fmt"
	country_dto "main_service/module/country_service/dto"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ─── Interface ───────────────────────────────────────────────────────────────

type CountryService interface {
	Create(ctx context.Context, req country_dto.CreateCountryRequest) (*country_dto.CountryResponse, error)
	GetByID(ctx context.Context, id int64) (*country_dto.CountryResponse, error)
	Update(ctx context.Context, id int64, req country_dto.UpdateCountryRequest) (*country_dto.CountryResponse, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, f country_dto.CountryFilter, afterID int64, limit int) ([]*country_dto.CountryResponse, bool, error)
	ListTree(ctx context.Context, parentID *int64, f country_dto.CountryFilter, afterID int64, limit int) ([]*country_dto.CountryResponse, bool, error)
}

// ─── Implementation ──────────────────────────────────────────────────────────

type countryService struct {
	db *pgxpool.Pool
}

func NewCountryService(db *pgxpool.Pool) CountryService {
	return &countryService{db: db}
}

var validSortCols = map[string]string{
	"id":         "c.id",
	"parent_id":  "c.parent_id",
	"name":       "c.name",
	"is_active":  "c.is_active",
	"created_at": "c.created_at",
	"updated_at": "c.updated_at",
}

// ─── Create ──────────────────────────────────────────────────────────────────

func (s *countryService) Create(ctx context.Context, req country_dto.CreateCountryRequest) (*country_dto.CountryResponse, error) {
	isActive := true
	{
		if req.IsActive != nil {
			isActive = *req.IsActive
		}
	}

	var id int64
	err := s.db.QueryRow(ctx, `
		INSERT INTO countries (parent_id, name, is_active)
		VALUES ($1, $2, $3)
		RETURNING id
	`, req.ParentID, req.Name, isActive).Scan(&id)
	if err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

// ─── GetByID ─────────────────────────────────────────────────────────────────

func (s *countryService) GetByID(ctx context.Context, id int64) (*country_dto.CountryResponse, error) {
	var r country_dto.CountryResponse
	err := s.db.QueryRow(ctx, `
		SELECT c.id, c.parent_id, COALESCE(p.name, ''),
		       c.name, c.is_active, c.created_at, c.updated_at, c.deleted_at
		FROM countries c
		LEFT JOIN countries p ON p.id = c.parent_id AND p.deleted_at IS NULL
		WHERE c.id = $1 AND c.deleted_at IS NULL
	`, id).Scan(
		&r.ID, &r.ParentID, &r.ParentName,
		&r.Name, &r.IsActive, &r.CreatedAt, &r.UpdatedAt, &r.DeletedAt,
	)
	return &r, err
}

// ─── Update ──────────────────────────────────────────────────────────────────

func (s *countryService) Update(ctx context.Context, id int64, req country_dto.UpdateCountryRequest) (*country_dto.CountryResponse, error) {
	setClauses := []string{"updated_at = NOW()"}
	args := []interface{}{}
	idx := 1

	if req.ParentID != nil {
		setClauses = append(setClauses, fmt.Sprintf("parent_id = $%d", idx))
		args = append(args, *req.ParentID)
		idx++
	}
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
		UPDATE countries SET %s
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

func (s *countryService) Delete(ctx context.Context, id int64) error {
	tag, err := s.db.Exec(ctx, `
		UPDATE countries SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL
	`, time.Now(), id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("country not found")
	}
	return nil
}

// ─── List ────────────────────────────────────────────────────────────────────

func (s *countryService) List(ctx context.Context, f country_dto.CountryFilter, afterID int64, limit int) ([]*country_dto.CountryResponse, bool, error) {
	args := []interface{}{}
	conditions := []string{"c.deleted_at IS NULL"}
	idx := 1

	if afterID > 0 {
		conditions = append(conditions, fmt.Sprintf("c.id < $%d", idx))
		args = append(args, afterID)
		idx++
	}
	if f.ParentID != nil {
		conditions = append(conditions, fmt.Sprintf("c.parent_id = $%d", idx))
		args = append(args, *f.ParentID)
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
		SELECT c.id, c.parent_id, COALESCE(p.name, ''),
		       c.name, c.is_active, c.created_at, c.updated_at, c.deleted_at
		FROM countries c
		LEFT JOIN countries p ON p.id = c.parent_id AND p.deleted_at IS NULL
		WHERE %s
		ORDER BY c.id DESC
		LIMIT $%d
	`, strings.Join(conditions, " AND "), idx)

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	var items []*country_dto.CountryResponse
	for rows.Next() {
		var r country_dto.CountryResponse
		if err := rows.Scan(
			&r.ID, &r.ParentID, &r.ParentName,
			&r.Name, &r.IsActive, &r.CreatedAt, &r.UpdatedAt, &r.DeletedAt,
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

// ─── ListTree ────────────────────────────────────────────────────────────────

func (s *countryService) ListTree(ctx context.Context, parentID *int64, f country_dto.CountryFilter, afterID int64, limit int) ([]*country_dto.CountryResponse, bool, error) {
	args := []interface{}{}
	conditions := []string{"c.deleted_at IS NULL"}
	idx := 1

	if afterID > 0 {
		conditions = append(conditions, fmt.Sprintf("c.id < $%d", idx))
		args = append(args, afterID)
		idx++
	}
	if parentID != nil {
		conditions = append(conditions, fmt.Sprintf("c.parent_id = $%d", idx))
		args = append(args, *parentID)
		idx++
	} else {
		conditions = append(conditions, "c.parent_id = 0")
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

	rootArgs := append(args, limit+1)
	rootQuery := fmt.Sprintf(
		"SELECT c.id FROM countries c WHERE %s ORDER BY c.id DESC LIMIT $%d",
		strings.Join(conditions, " AND "), idx,
	)

	rootRows, err := s.db.Query(ctx, rootQuery, rootArgs...)
	if err != nil {
		return nil, false, err
	}
	defer rootRows.Close()

	var rootIDs []int64
	for rootRows.Next() {
		var id int64
		if err := rootRows.Scan(&id); err != nil {
			return nil, false, err
		}
		rootIDs = append(rootIDs, id)
	}
	rootRows.Close()

	hasMore := len(rootIDs) > limit
	if hasMore {
		rootIDs = rootIDs[:limit]
	}

	if len(rootIDs) == 0 {
		return []*country_dto.CountryResponse{}, hasMore, nil
	}

	allRows, err := s.db.Query(ctx, `
		WITH RECURSIVE tree AS (
			SELECT c.id, c.parent_id, c.name, c.is_active, c.created_at, c.updated_at, c.deleted_at
			FROM countries c
			WHERE c.id = ANY($1) AND c.deleted_at IS NULL
			UNION ALL
			SELECT c.id, c.parent_id, c.name, c.is_active, c.created_at, c.updated_at, c.deleted_at
			FROM countries c
			JOIN tree t ON c.parent_id = t.id
			WHERE c.deleted_at IS NULL
		)
		SELECT id, parent_id, name, is_active, created_at, updated_at, deleted_at
		FROM tree
		ORDER BY parent_id NULLS FIRST, id
	`, rootIDs)
	if err != nil {
		return nil, false, err
	}
	defer allRows.Close()

	nodeMap := map[int64]*country_dto.CountryResponse{}
	for allRows.Next() {
		var r country_dto.CountryResponse
		if err := allRows.Scan(&r.ID, &r.ParentID, &r.Name, &r.IsActive, &r.CreatedAt, &r.UpdatedAt, &r.DeletedAt); err != nil {
			return nil, false, err
		}
		rr := r
		nodeMap[r.ID] = &rr
	}
	if err := allRows.Err(); err != nil {
		return nil, false, err
	}

	rootSet := make(map[int64]bool, len(rootIDs))
	for _, id := range rootIDs {
		rootSet[id] = true
	}
	for _, node := range nodeMap {
		if !rootSet[node.ID] && node.ParentID != nil {
			if parent, ok := nodeMap[*node.ParentID]; ok {
				parent.Children = append(parent.Children, node)
			}
		}
	}

	roots := make([]*country_dto.CountryResponse, 0, len(rootIDs))
	for _, id := range rootIDs {
		if node, ok := nodeMap[id]; ok {
			roots = append(roots, node)
		}
	}
	return roots, hasMore, nil
}
