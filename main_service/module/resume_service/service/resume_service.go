package resume_service

import (
	"context"
	"fmt"
	"main_service/helper"
	resume_dto "main_service/module/resume_service/dto"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ─── Interface ───────────────────────────────────────────────────────────────

type ResumeService interface {
	Create(ctx context.Context, userID int64, req resume_dto.CreateResumeRequest) (*resume_dto.ResumeResponse, error)
	GetByID(ctx context.Context, id int64) (*resume_dto.ResumeResponse, error)
	GetBySlug(ctx context.Context, slug string) (*resume_dto.ResumeResponse, error)
	Update(ctx context.Context, id, userID int64, req resume_dto.UpdateResumeRequest) (*resume_dto.ResumeResponse, error)
	Delete(ctx context.Context, id, userID int64) error
	List(ctx context.Context, f resume_dto.ResumeFilter, cursor helper.CursorPayload, limit int) ([]*resume_dto.ResumeResponse, bool, int64, error)
	AddCategory(ctx context.Context, resumeID, categoryID int64) error
	RemoveCategory(ctx context.Context, resumeID, categoryID int64) error
}

// ─── Implementation ──────────────────────────────────────────────────────────

type resumeService struct {
	db *pgxpool.Pool
}

func NewResumeService(db *pgxpool.Pool) ResumeService {
	return &resumeService{db: db}
}

var validSortCols = map[string]string{
	"id":              "rs.id",
	"name":            "rs.name",
	"title":           "rs.title",
	"price":           "rs.price",
	"experience_year": "rs.experience_year",
	"views_count":     "rs.views_count",
	"is_active":       "rs.is_active",
	"created_at":      "rs.created_at",
	"updated_at":      "rs.updated_at",
}

func (s *resumeService) fetchCategories(ctx context.Context, resumeID int64) ([]resume_dto.CategoryShort, error) {
	rows, err := s.db.Query(ctx, `
		SELECT c.id, c.name, c.is_active
		FROM categories c
		JOIN category_resume cr ON cr.categorya_id = c.id
		WHERE cr.resume_id = $1 AND c.deleted_at IS NULL
		ORDER BY c.id
	`, resumeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []resume_dto.CategoryShort
	for rows.Next() {
		var c resume_dto.CategoryShort
		if err := rows.Scan(&c.ID, &c.Name, &c.IsActive); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	if cats == nil {
		cats = []resume_dto.CategoryShort{}
	}
	return cats, rows.Err()
}

func scanResume(row interface{ Scan(...interface{}) error }, r *resume_dto.ResumeResponse) error {
	return row.Scan(
		&r.ID, &r.Slug, &r.UserID,
		&r.RegionID, &r.RegionName,
		&r.DistrictID, &r.DistrictName,
		&r.MahallaID, &r.MahallaName,
		&r.Adress, &r.Name, &r.Photo, &r.Title, &r.Text, &r.Contact,
		&r.Price, &r.ExperienceYear, &r.Skills,
		&r.ViewsCount, &r.IsActive,
		&r.CreatedAt, &r.UpdatedAt, &r.DeletedAt,
	)
}

const resumeSelectJoin = `
	SELECT rs.id, rs.slug, rs.user_id,
	       rs.region_id,   COALESCE(r.name, ''),
	       rs.district_id, COALESCE(d.name, ''),
	       rs.mahalla_id,  COALESCE(m.name, ''),
	       rs.adress, rs.name, rs.photo, rs.title, rs.text, rs.contact,
	       rs.price, rs.experience_year, rs.skills,
	       rs.views_count, rs.is_active,
	       rs.created_at, rs.updated_at, rs.deleted_at
	FROM resumes rs
	LEFT JOIN countries r ON r.id = rs.region_id   AND r.deleted_at IS NULL
	LEFT JOIN countries d ON d.id = rs.district_id AND d.deleted_at IS NULL
	LEFT JOIN countries m ON m.id = rs.mahalla_id  AND m.deleted_at IS NULL
`

// ─── Create ──────────────────────────────────────────────────────────────────

func (s *resumeService) Create(ctx context.Context, userID int64, req resume_dto.CreateResumeRequest) (*resume_dto.ResumeResponse, error) {
	isActive := true
	{
		if req.IsActive != nil {
			isActive = *req.IsActive
		}
	}

	var id int64
	err := s.db.QueryRow(ctx, `
		INSERT INTO resumes (user_id, region_id, district_id, mahalla_id, adress, name, photo, title, text, contact, price, experience_year, skills, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id
	`, userID, req.RegionID, req.DistrictID, req.MahallaID,
		req.Adress, req.Name, req.Photo, req.Title, req.Text, req.Contact,
		req.Price, req.ExperienceYear, req.Skills, isActive).Scan(&id)
	if err != nil {
		return nil, err
	}

	if _, err := s.db.Exec(ctx, `UPDATE resumes SET slug = $1 WHERE id = $2`, helper.Slug(req.Name, id), id); err != nil {
		return nil, err
	}

	if len(req.CategoryIDs) > 0 {
		for _, catID := range req.CategoryIDs {
			s.db.Exec(ctx, `
				INSERT INTO category_resume (categorya_id, resume_id)
				VALUES ($1, $2)
				ON CONFLICT (categorya_id, resume_id) DO NOTHING
			`, catID, id)
		}
	}

	return s.GetByID(ctx, id)
}

// ─── GetByID ─────────────────────────────────────────────────────────────────

func (s *resumeService) GetByID(ctx context.Context, id int64) (*resume_dto.ResumeResponse, error) {
	var r resume_dto.ResumeResponse
	row := s.db.QueryRow(ctx, resumeSelectJoin+`WHERE rs.id = $1 AND rs.deleted_at IS NULL`, id)
	if err := scanResume(row, &r); err != nil {
		return nil, err
	}
	cats, err := s.fetchCategories(ctx, id)
	if err != nil {
		return nil, err
	}
	r.Categories = cats
	return &r, nil
}

// ─── GetBySlug ───────────────────────────────────────────────────────────────

func (s *resumeService) GetBySlug(ctx context.Context, slug string) (*resume_dto.ResumeResponse, error) {
	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_, _ = s.db.Exec(bgCtx, `
			UPDATE resumes SET views_count = COALESCE(views_count, 0) + 1 WHERE slug = $1
		`, slug)
	}()

	var r resume_dto.ResumeResponse
	row := s.db.QueryRow(ctx, resumeSelectJoin+`WHERE rs.slug = $1 AND rs.deleted_at IS NULL`, slug)
	if err := scanResume(row, &r); err != nil {
		return nil, err
	}
	cats, err := s.fetchCategories(ctx, r.ID)
	if err != nil {
		return nil, err
	}
	r.Categories = cats
	return &r, nil
}

// ─── Update ──────────────────────────────────────────────────────────────────

func (s *resumeService) Update(ctx context.Context, id, userID int64, req resume_dto.UpdateResumeRequest) (*resume_dto.ResumeResponse, error) {
	setClauses := []string{"updated_at = NOW()"}
	args := []interface{}{}
	idx := 1

	fields := []struct {
		val interface{}
		col string
	}{
		{req.RegionID, "region_id"}, {req.DistrictID, "district_id"}, {req.MahallaID, "mahalla_id"},
		{req.Adress, "adress"}, {req.Name, "name"}, {req.Photo, "photo"},
		{req.Title, "title"}, {req.Text, "text"}, {req.Contact, "contact"},
		{req.Price, "price"}, {req.ExperienceYear, "experience_year"},
		{req.Skills, "skills"}, {req.IsActive, "is_active"},
	}

	for _, f := range fields {
		if isNilPtr(f.val) {
			continue
		}
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", f.col, idx))
		args = append(args, derefPtr(f.val))
		idx++
	}

	args = append(args, id, userID)
	query := fmt.Sprintf(`
		UPDATE resumes SET %s
		WHERE id = $%d AND user_id = $%d AND deleted_at IS NULL
		RETURNING id
	`, strings.Join(setClauses, ", "), idx, idx+1)

	var retID int64
	if err := s.db.QueryRow(ctx, query, args...).Scan(&retID); err != nil {
		return nil, fmt.Errorf("resume not found or access denied")
	}
	return s.GetByID(ctx, retID)
}

// ─── Delete ──────────────────────────────────────────────────────────────────

func (s *resumeService) Delete(ctx context.Context, id, userID int64) error {
	tag, err := s.db.Exec(ctx, `
		UPDATE resumes SET deleted_at = $1
		WHERE id = $2 AND user_id = $3 AND deleted_at IS NULL
	`, time.Now(), id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("resume not found or access denied")
	}
	return nil
}

// ─── AddCategory / RemoveCategory ────────────────────────────────────────────

func (s *resumeService) AddCategory(ctx context.Context, resumeID, categoryID int64) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO category_resume (categorya_id, resume_id)
		VALUES ($1, $2)
		ON CONFLICT (categorya_id, resume_id) DO NOTHING
	`, categoryID, resumeID)
	return err
}

func (s *resumeService) RemoveCategory(ctx context.Context, resumeID, categoryID int64) error {
	tag, err := s.db.Exec(ctx, `
		DELETE FROM category_resume WHERE categorya_id = $1 AND resume_id = $2
	`, categoryID, resumeID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("category not attached to this resume")
	}
	return nil
}

// escapeLike escapes ILIKE special characters to prevent wildcard injection.
func escapeLike(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}

// ─── List ────────────────────────────────────────────────────────────────────

func (s *resumeService) List(ctx context.Context, f resume_dto.ResumeFilter, cursor helper.CursorPayload, limit int) ([]*resume_dto.ResumeResponse, bool, int64, error) {
	args := []interface{}{}
	conditions := []string{"rs.deleted_at IS NULL"}
	idx := 1

	if f.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("rs.user_id = $%d", idx))
		args = append(args, *f.UserID)
		idx++
	}
	if f.RegionID != nil {
		conditions = append(conditions, fmt.Sprintf("rs.region_id = $%d", idx))
		args = append(args, *f.RegionID)
		idx++
	}
	if f.DistrictID != nil {
		conditions = append(conditions, fmt.Sprintf("rs.district_id = $%d", idx))
		args = append(args, *f.DistrictID)
		idx++
	}
	if f.MahallaID != nil {
		conditions = append(conditions, fmt.Sprintf("rs.mahalla_id = $%d", idx))
		args = append(args, *f.MahallaID)
		idx++
	}
	if f.Name != "" {
		conditions = append(conditions, fmt.Sprintf("rs.name ILIKE $%d", idx))
		args = append(args, "%"+escapeLike(f.Name)+"%")
		idx++
	}
	if f.Title != "" {
		conditions = append(conditions, fmt.Sprintf("rs.title ILIKE $%d", idx))
		args = append(args, "%"+escapeLike(f.Title)+"%")
		idx++
	}
	if f.Search != "" {
		conditions = append(conditions, fmt.Sprintf(
			"(rs.name ILIKE $%d OR rs.title ILIKE $%d OR rs.skills ILIKE $%d OR rs.text ILIKE $%d)",
			idx, idx, idx, idx,
		))
		args = append(args, "%"+escapeLike(f.Search)+"%")
		idx++
	}
	if f.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("rs.is_active = $%d", idx))
		args = append(args, *f.IsActive)
		idx++
	}
	if f.MinPrice != nil {
		conditions = append(conditions, fmt.Sprintf("rs.price >= $%d", idx))
		args = append(args, *f.MinPrice)
		idx++
	}
	if f.MaxPrice != nil {
		conditions = append(conditions, fmt.Sprintf("rs.price <= $%d", idx))
		args = append(args, *f.MaxPrice)
		idx++
	}
	if f.MinExperience != nil {
		conditions = append(conditions, fmt.Sprintf("rs.experience_year >= $%d", idx))
		args = append(args, *f.MinExperience)
		idx++
	}
	if f.CategoryID != nil {
		conditions = append(conditions, fmt.Sprintf(`EXISTS (
			SELECT 1 FROM category_resume cr
			WHERE cr.resume_id = rs.id AND cr.categorya_id = $%d
		)`, idx))
		args = append(args, *f.CategoryID)
		idx++
	}
	if len(f.CategoryIDs) > 0 {
		conditions = append(conditions, fmt.Sprintf(`EXISTS (
			SELECT 1 FROM category_resume cr
			WHERE cr.resume_id = rs.id AND cr.categorya_id = ANY($%d)
		)`, idx))
		args = append(args, f.CategoryIDs)
		idx++
	}

	// Snapshot args/conditions for COUNT (no cursor) before modifying for data query
	countConditions := make([]string, len(conditions))
	copy(countConditions, conditions)
	countArgs := make([]interface{}, len(args))
	copy(countArgs, args)

	orderCol, fallbackValue, valueKind := resumeOrderConfig(f.SortBy, f.SortOrder)
	orderDir := normalizedOrder(f.SortOrder)
	conditions, args, idx = appendResumeCursorCondition(conditions, args, idx, orderCol, fallbackValue, valueKind, orderDir, cursor)

	args = append(args, limit+1)
	query := fmt.Sprintf(`
		SELECT rs.id, rs.slug, rs.user_id,
		       rs.region_id,   COALESCE(r.name, ''),
		       rs.district_id, COALESCE(d.name, ''),
		       rs.mahalla_id,  COALESCE(m.name, ''),
		       rs.adress, rs.name, rs.photo, rs.title, rs.text, rs.contact,
		       rs.price, rs.experience_year, rs.skills,
		       rs.views_count, rs.is_active,
		       rs.created_at, rs.updated_at, rs.deleted_at
		FROM resumes rs
		LEFT JOIN countries r ON r.id = rs.region_id   AND r.deleted_at IS NULL
		LEFT JOIN countries d ON d.id = rs.district_id AND d.deleted_at IS NULL
		LEFT JOIN countries m ON m.id = rs.mahalla_id  AND m.deleted_at IS NULL
		WHERE %s
		ORDER BY %s %s, rs.id %s
		LIMIT $%d
	`, strings.Join(conditions, " AND "), orderCol, orderDir, orderDir, idx)

	// Run COUNT and data query in parallel
	var (
		total    int64
		totalErr error
		wg       sync.WaitGroup
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		countSQL := fmt.Sprintf(`SELECT COUNT(*) FROM resumes rs WHERE %s`, strings.Join(countConditions, " AND "))
		totalErr = s.db.QueryRow(ctx, countSQL, countArgs...).Scan(&total)
	}()

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		wg.Wait()
		return nil, false, 0, err
	}
	defer rows.Close()

	var items []*resume_dto.ResumeResponse
	for rows.Next() {
		var rs resume_dto.ResumeResponse
		if err := rows.Scan(
			&rs.ID, &rs.Slug, &rs.UserID,
			&rs.RegionID, &rs.RegionName,
			&rs.DistrictID, &rs.DistrictName,
			&rs.MahallaID, &rs.MahallaName,
			&rs.Adress, &rs.Name, &rs.Photo, &rs.Title, &rs.Text, &rs.Contact,
			&rs.Price, &rs.ExperienceYear, &rs.Skills,
			&rs.ViewsCount, &rs.IsActive,
			&rs.CreatedAt, &rs.UpdatedAt, &rs.DeletedAt,
		); err != nil {
			wg.Wait()
			return nil, false, 0, err
		}
		rs.Categories = []resume_dto.CategoryShort{}
		items = append(items, &rs)
	}
	if err := rows.Err(); err != nil {
		wg.Wait()
		return nil, false, 0, err
	}

	wg.Wait()
	if totalErr != nil {
		return nil, false, 0, totalErr
	}

	hasMore := len(items) > limit
	if hasMore {
		items = items[:limit]
	}

	if len(items) > 0 {
		ids := make([]int64, len(items))
		idxMap := make(map[int64]int, len(items))
		for i, item := range items {
			ids[i] = item.ID
			idxMap[item.ID] = i
		}
		catRows, err := s.db.Query(ctx, `
			SELECT cr.resume_id, c.id, c.name, c.is_active
			FROM category_resume cr
			JOIN categories c ON c.id = cr.categorya_id
			WHERE cr.resume_id = ANY($1) AND c.deleted_at IS NULL
			ORDER BY cr.resume_id, c.id
		`, ids)
		if err == nil {
			defer catRows.Close()
			for catRows.Next() {
				var resumeID int64
				var cat resume_dto.CategoryShort
				if err := catRows.Scan(&resumeID, &cat.ID, &cat.Name, &cat.IsActive); err == nil {
					if i, ok := idxMap[resumeID]; ok {
						items[i].Categories = append(items[i].Categories, cat)
					}
				}
			}
		}
	}

	return items, hasMore, total, nil
}

func resumeOrderConfig(sortBy, sortOrder string) (string, string, string) {
	switch sortBy {
	case "price":
		if normalizedOrder(sortOrder) == "ASC" {
			return "COALESCE(rs.price, 9223372036854775807)", "9223372036854775807", "int64"
		}
		return "COALESCE(rs.price, -1)", "-1", "int64"
	case "experience_year":
		if normalizedOrder(sortOrder) == "ASC" {
			return "COALESCE(rs.experience_year, 2147483647)", "2147483647", "int"
		}
		return "COALESCE(rs.experience_year, -1)", "-1", "int"
	default:
		return "rs.id", "", ""
	}
}

func normalizedOrder(order string) string {
	if strings.EqualFold(order, "asc") {
		return "ASC"
	}
	return "DESC"
}

func appendResumeCursorCondition(conditions []string, args []interface{}, idx int, orderCol, fallbackValue, valueKind, orderDir string, cursor helper.CursorPayload) ([]string, []interface{}, int) {
	if cursor.ID <= 0 {
		return conditions, args, idx
	}
	if orderCol == "rs.id" {
		op := "<"
		if orderDir == "ASC" {
			op = ">"
		}
		conditions = append(conditions, fmt.Sprintf("rs.id %s $%d", op, idx))
		args = append(args, cursor.ID)
		return conditions, args, idx + 1
	}

	cursorValue := fallbackValue
	if cursor.Value != "" {
		cursorValue = cursor.Value
	}
	op := "<"
	if orderDir == "ASC" {
		op = ">"
	}
	conditions = append(conditions, fmt.Sprintf("(%s %s $%d OR (%s = $%d AND rs.id %s $%d))", orderCol, op, idx, orderCol, idx, op, idx+1))
	switch valueKind {
	case "int64":
		n, err := strconv.ParseInt(cursorValue, 10, 64)
		if err != nil {
			n, _ = strconv.ParseInt(fallbackValue, 10, 64)
		}
		args = append(args, n, cursor.ID)
	case "int":
		n, err := strconv.Atoi(cursorValue)
		if err != nil {
			n, _ = strconv.Atoi(fallbackValue)
		}
		args = append(args, n, cursor.ID)
	default:
		args = append(args, cursorValue, cursor.ID)
	}
	return conditions, args, idx + 2
}

// ─── pointer helpers ─────────────────────────────────────────────────────────

func isNilPtr(v interface{}) bool {
	switch x := v.(type) {
	case *int64:
		return x == nil
	case *int:
		return x == nil
	case *string:
		return x == nil
	case *bool:
		return x == nil
	}
	return v == nil
}

func derefPtr(v interface{}) interface{} {
	switch x := v.(type) {
	case *int64:
		return *x
	case *int:
		return *x
	case *string:
		return *x
	case *bool:
		return *x
	}
	return v
}
