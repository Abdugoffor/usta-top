package region_service

import (
	"context"
	"fmt"
	region_dto "main_service/module/region_service/dto"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ─── Interface ───────────────────────────────────────────────────────────────

type RegionService interface {
	List(ctx context.Context, f region_dto.RegionFilter) ([]region_dto.RegionResponse, int64, error)
	ListTree(ctx context.Context, parentID int64, f region_dto.RegionFilter) ([]*region_dto.RegionResponse, int64, error)
	GetByID(ctx context.Context, id int64) (*region_dto.RegionResponse, error)
	Create(ctx context.Context, req region_dto.CreateRegionRequest) (*region_dto.RegionResponse, error)
	Update(ctx context.Context, id int64, req region_dto.UpdateRegionRequest) (*region_dto.RegionResponse, error)
	Delete(ctx context.Context, id int64) error
}

// ─── Implementation ──────────────────────────────────────────────────────────

type regionService struct {
	db *pgxpool.Pool
}

func NewRegionService(db *pgxpool.Pool) RegionService {
	return &regionService{db: db}
}

// ─── Create ──────────────────────────────────────────────────────────────────

func (s *regionService) Create(ctx context.Context, req region_dto.CreateRegionRequest) (*region_dto.RegionResponse, error) {
	isActive := true
	{
		if req.IsActive != nil {
			isActive = *req.IsActive
		}
	}

	var desc *string
	{
		if req.Description != "" {
			desc = &req.Description
		}
	}

	var resp region_dto.RegionResponse
	err := s.db.QueryRow(ctx, `
		INSERT INTO regions (parent_id, name, description, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, parent_id, name, COALESCE(description, ''), is_active, created_at, updated_at, deleted_at
	`, req.ParentID, req.Name, desc, isActive).Scan(
		&resp.ID, &resp.ParentID, &resp.Name, &resp.Description,
		&resp.IsActive, &resp.CreatedAt, &resp.UpdatedAt, &resp.DeletedAt,
	)
	return &resp, err
}

// ─── GetByID ─────────────────────────────────────────────────────────────────

func (s *regionService) GetByID(ctx context.Context, id int64) (*region_dto.RegionResponse, error) {
	var resp region_dto.RegionResponse
	err := s.db.QueryRow(ctx, `
		SELECT r.id, r.parent_id, COALESCE(p.name, ''),
		       r.name, COALESCE(r.description, ''),
		       r.is_active, r.created_at, r.updated_at, r.deleted_at
		FROM regions r
		LEFT JOIN regions p ON p.id = r.parent_id AND p.deleted_at IS NULL
		WHERE r.id = $1 AND r.deleted_at IS NULL
	`, id).Scan(
		&resp.ID, &resp.ParentID, &resp.ParentName,
		&resp.Name, &resp.Description,
		&resp.IsActive, &resp.CreatedAt, &resp.UpdatedAt, &resp.DeletedAt,
	)
	return &resp, err
}

// ─── Update ──────────────────────────────────────────────────────────────────

func (s *regionService) Update(ctx context.Context, id int64, req region_dto.UpdateRegionRequest) (*region_dto.RegionResponse, error) {
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
	if req.Description != nil {
		var desc *string
		{
			if *req.Description != "" {
				desc = req.Description
			}
		}
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", idx))
		args = append(args, desc)
		idx++
	}
	if req.IsActive != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_active = $%d", idx))
		args = append(args, *req.IsActive)
		idx++
	}

	args = append(args, id)
	query := fmt.Sprintf(`
		UPDATE regions SET %s
		WHERE id = $%d AND deleted_at IS NULL
		RETURNING id, parent_id, name, COALESCE(description, ''), is_active, created_at, updated_at, deleted_at
	`, strings.Join(setClauses, ", "), idx)

	var resp region_dto.RegionResponse
	err := s.db.QueryRow(ctx, query, args...).Scan(
		&resp.ID, &resp.ParentID, &resp.Name, &resp.Description,
		&resp.IsActive, &resp.CreatedAt, &resp.UpdatedAt, &resp.DeletedAt,
	)
	return &resp, err
}

// ─── Delete (soft) ───────────────────────────────────────────────────────────

func (s *regionService) Delete(ctx context.Context, id int64) error {
	tag, err := s.db.Exec(ctx, `
		UPDATE regions SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL
	`, time.Now(), id)
	{
		if err != nil {
			return err
		}
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("region not found")
	}
	return nil
}

// ─── ListTree ────────────────────────────────────────────────────────────────

func (s *regionService) ListTree(ctx context.Context, parentID int64, f region_dto.RegionFilter) ([]*region_dto.RegionResponse, int64, error) {
	args := []interface{}{parentID}
	conditions := []string{"r.deleted_at IS NULL", "r.parent_id = $1"}
	idx := 2

	if f.Name != "" {
		conditions = append(conditions, fmt.Sprintf("r.name ILIKE $%d", idx))
		args = append(args, "%"+f.Name+"%")
		idx++
	}
	if f.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("r.is_active = $%d", idx))
		args = append(args, *f.IsActive)
		idx++
	}

	where := strings.Join(conditions, " AND ")

	var total int64
	if err := s.db.QueryRow(ctx, "SELECT COUNT(*) FROM regions r WHERE "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	sortCol := "r.id"
	if col, ok := validSortCols[f.SortBy]; ok {
		sortCol = col
	}
	sortOrder := "ASC"
	if strings.ToUpper(f.SortOrder) == "DESC" {
		sortOrder = "DESC"
	}

	rootArgs := append(args, f.Limit, (f.Page-1)*f.Limit)
	rootQuery := fmt.Sprintf(
		"SELECT r.id FROM regions r WHERE %s ORDER BY %s %s LIMIT $%d OFFSET $%d",
		where, sortCol, sortOrder, idx, idx+1,
	)

	rootRows, err := s.db.Query(ctx, rootQuery, rootArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rootRows.Close()

	var rootIDs []int64
	for rootRows.Next() {
		var id int64
		if err := rootRows.Scan(&id); err != nil {
			return nil, 0, err
		}
		rootIDs = append(rootIDs, id)
	}
	rootRows.Close()

	if len(rootIDs) == 0 {
		return []*region_dto.RegionResponse{}, total, nil
	}

	allRows, err := s.db.Query(ctx, `
		WITH RECURSIVE tree AS (
			SELECT r.id, r.parent_id, r.name, COALESCE(r.description,'') AS description,
			       r.is_active, r.created_at, r.updated_at, r.deleted_at
			FROM regions r
			WHERE r.id = ANY($1) AND r.deleted_at IS NULL
			UNION ALL
			SELECT r.id, r.parent_id, r.name, COALESCE(r.description,''),
			       r.is_active, r.created_at, r.updated_at, r.deleted_at
			FROM regions r
			JOIN tree t ON r.parent_id = t.id
			WHERE r.deleted_at IS NULL
		)
		SELECT id, parent_id, name, description, is_active, created_at, updated_at, deleted_at
		FROM tree
		ORDER BY parent_id, id
	`, rootIDs)
	if err != nil {
		return nil, 0, err
	}
	defer allRows.Close()

	nodeMap := map[int64]*region_dto.RegionResponse{}
	for allRows.Next() {
		var r region_dto.RegionResponse
		if err := allRows.Scan(&r.ID, &r.ParentID, &r.Name, &r.Description,
			&r.IsActive, &r.CreatedAt, &r.UpdatedAt, &r.DeletedAt); err != nil {
			return nil, 0, err
		}
		rr := r
		nodeMap[r.ID] = &rr
	}
	if err := allRows.Err(); err != nil {
		return nil, 0, err
	}

	rootIDSet := make(map[int64]bool, len(rootIDs))
	for _, id := range rootIDs {
		rootIDSet[id] = true
	}

	for _, node := range nodeMap {
		if !rootIDSet[node.ID] {
			if parent, ok := nodeMap[node.ParentID]; ok {
				parent.Children = append(parent.Children, node)
			}
		}
	}

	roots := make([]*region_dto.RegionResponse, 0, len(rootIDs))
	for _, id := range rootIDs {
		if node, ok := nodeMap[id]; ok {
			roots = append(roots, node)
		}
	}

	return roots, total, nil
}

// ─── List ────────────────────────────────────────────────────────────────────

var validSortCols = map[string]string{
	"id":          "r.id",
	"parent_id":   "r.parent_id",
	"name":        "r.name",
	"description": "r.description",
	"is_active":   "r.is_active",
	"created_at":  "r.created_at",
	"updated_at":  "r.updated_at",
}

func (s *regionService) List(ctx context.Context, f region_dto.RegionFilter) ([]region_dto.RegionResponse, int64, error) {
	args := []interface{}{}
	conditions := []string{"r.deleted_at IS NULL"}
	idx := 1

	if f.ID != nil {
		conditions = append(conditions, fmt.Sprintf("r.id = $%d", idx))
		args = append(args, *f.ID)
		idx++
	}
	if f.ParentID != nil {
		conditions = append(conditions, fmt.Sprintf("r.parent_id = $%d", idx))
		args = append(args, *f.ParentID)
		idx++
	}
	if f.Name != "" {
		conditions = append(conditions, fmt.Sprintf("r.name ILIKE $%d", idx))
		args = append(args, "%"+f.Name+"%")
		idx++
	}
	if f.Description != "" {
		conditions = append(conditions, fmt.Sprintf("r.description ILIKE $%d", idx))
		args = append(args, "%"+f.Description+"%")
		idx++
	}
	if f.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("r.is_active = $%d", idx))
		args = append(args, *f.IsActive)
		idx++
	}

	sortCol := "r.id"
	{
		if col, ok := validSortCols[f.SortBy]; ok {
			sortCol = col
		}
	}
	sortOrder := "ASC"
	{
		if strings.ToUpper(f.SortOrder) == "DESC" {
			sortOrder = "DESC"
		}
	}

	args = append(args, f.Limit, (f.Page-1)*f.Limit)

	query := fmt.Sprintf(`
		SELECT r.id, r.parent_id, COALESCE(p.name, ''),
		       r.name, COALESCE(r.description, ''),
		       r.is_active, r.created_at, r.updated_at, r.deleted_at,
		       COUNT(*) OVER() AS total
		FROM regions r
		LEFT JOIN regions p ON p.id = r.parent_id AND p.deleted_at IS NULL
		WHERE %s 
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, strings.Join(conditions, " AND "), sortCol, sortOrder, idx, idx+1)

	rows, err := s.db.Query(ctx, query, args...)
	{
		if err != nil {
			return nil, 0, err
		}
	}
	defer rows.Close()

	var regions []region_dto.RegionResponse
	var total int64

	for rows.Next() {
		var r region_dto.RegionResponse
		{
			if err := rows.Scan(
				&r.ID, &r.ParentID, &r.ParentName,
				&r.Name, &r.Description,
				&r.IsActive, &r.CreatedAt, &r.UpdatedAt, &r.DeletedAt,
				&total,
			); err != nil {
				return nil, 0, err
			}
		}
		regions = append(regions, r)
	}

	return regions, total, rows.Err()
}
