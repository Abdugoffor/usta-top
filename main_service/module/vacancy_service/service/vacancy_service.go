package vacancy_service

import (
	"context"
	"fmt"
	"main_service/helper"
	vacancy_dto "main_service/module/vacancy_service/dto"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VacancyService interface {
	Create(ctx context.Context, userID int64, req vacancy_dto.CreateVacancyRequest) (*vacancy_dto.VacancyResponse, error)
	GetByID(ctx context.Context, id int64) (*vacancy_dto.VacancyResponse, error)
	GetBySlug(ctx context.Context, slug string) (*vacancy_dto.VacancyResponse, error)
	Update(ctx context.Context, id, userID int64, req vacancy_dto.UpdateVacancyRequest) (*vacancy_dto.VacancyResponse, error)
	Delete(ctx context.Context, id, userID int64) error
	List(ctx context.Context, f vacancy_dto.VacancyFilter, cursor helper.CursorPayload, limit int) ([]*vacancy_dto.VacancyResponse, bool, int64, error)
}

type vacancyService struct {
	db *pgxpool.Pool
}

func NewVacancyService(db *pgxpool.Pool) VacancyService {
	return &vacancyService{db: db}
}

var validSortCols = map[string]string{
	"id":          "v.id",
	"name":        "v.name",
	"title":       "v.title",
	"price":       "v.price",
	"views_count": "v.views_count",
	"is_active":   "v.is_active",
	"created_at":  "v.created_at",
	"updated_at":  "v.updated_at",
}

const vacancySelectJoin = `
	SELECT v.id, v.slug, v.user_id,
	       v.region_id,   COALESCE(r.name, ''),
	       v.district_id, COALESCE(d.name, ''),
	       v.mahalla_id,  COALESCE(m.name, ''),
	       v.adress, v.name, v.title, v.text, v.contact,
	       v.price, v.views_count, v.is_active,
	       v.created_at, v.updated_at, v.deleted_at
	FROM vacancies v
	LEFT JOIN countries r ON r.id = v.region_id   AND r.deleted_at IS NULL
	LEFT JOIN countries d ON d.id = v.district_id AND d.deleted_at IS NULL
	LEFT JOIN countries m ON m.id = v.mahalla_id  AND m.deleted_at IS NULL
`

func scanVacancy(row interface{ Scan(...interface{}) error }, v *vacancy_dto.VacancyResponse) error {
	return row.Scan(
		&v.ID, &v.Slug, &v.UserID,
		&v.RegionID, &v.RegionName,
		&v.DistrictID, &v.DistrictName,
		&v.MahallaID, &v.MahallaName,
		&v.Adress, &v.Name, &v.Title, &v.Text, &v.Contact,
		&v.Price, &v.ViewsCount, &v.IsActive,
		&v.CreatedAt, &v.UpdatedAt, &v.DeletedAt,
	)
}

func (s *vacancyService) fetchCategories(ctx context.Context, vacancyID int64) ([]vacancy_dto.CategoryShort, error) {
	rows, err := s.db.Query(ctx, `
		SELECT c.id, c.name, c.is_active
		FROM categories c
		JOIN category_vacancy cv ON cv.categorya_id = c.id
		WHERE cv.vacancy_id = $1 AND c.deleted_at IS NULL
		ORDER BY c.id
	`, vacancyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []vacancy_dto.CategoryShort
	for rows.Next() {
		var cat vacancy_dto.CategoryShort
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.IsActive); err != nil {
			return nil, err
		}
		cats = append(cats, cat)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if cats == nil {
		cats = []vacancy_dto.CategoryShort{}
	}
	return cats, nil
}

func normalizeCategoryIDs(ids []int64) []int64 {
	if len(ids) == 0 {
		return nil
	}

	seen := make(map[int64]struct{}, len(ids))
	normalized := make([]int64, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		normalized = append(normalized, id)
	}
	return normalized
}

func (s *vacancyService) attachCategories(ctx context.Context, tx pgx.Tx, vacancyID int64, categoryIDs []int64) error {
	for _, catID := range normalizeCategoryIDs(categoryIDs) {
		if _, err := tx.Exec(ctx, `
			INSERT INTO category_vacancy (categorya_id, vacancy_id)
			VALUES ($1, $2)
			ON CONFLICT (categorya_id, vacancy_id) DO NOTHING
		`, catID, vacancyID); err != nil {
			return err
		}
	}
	return nil
}

func (s *vacancyService) replaceCategories(ctx context.Context, tx pgx.Tx, vacancyID int64, categoryIDs []int64) error {
	if _, err := tx.Exec(ctx, `DELETE FROM category_vacancy WHERE vacancy_id = $1`, vacancyID); err != nil {
		return err
	}
	return s.attachCategories(ctx, tx, vacancyID, categoryIDs)
}

func (s *vacancyService) Create(ctx context.Context, userID int64, req vacancy_dto.CreateVacancyRequest) (*vacancy_dto.VacancyResponse, error) {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var id int64
	err = tx.QueryRow(ctx, `
		INSERT INTO vacancies (user_id, region_id, district_id, mahalla_id, adress, name, title, text, contact, price, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`, userID, req.RegionID, req.DistrictID, req.MahallaID,
		req.Adress, req.Name, req.Title, req.Text, req.Contact,
		req.Price, isActive).Scan(&id)
	if err != nil {
		return nil, err
	}

	if _, err := tx.Exec(ctx, `UPDATE vacancies SET slug = $1 WHERE id = $2`, helper.Slug(req.Name, id), id); err != nil {
		return nil, err
	}

	if err := s.attachCategories(ctx, tx, id, req.CategoryIDs); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, id)
}

func (s *vacancyService) GetByID(ctx context.Context, id int64) (*vacancy_dto.VacancyResponse, error) {
	var v vacancy_dto.VacancyResponse
	row := s.db.QueryRow(ctx, vacancySelectJoin+`WHERE v.id = $1 AND v.deleted_at IS NULL`, id)
	if err := scanVacancy(row, &v); err != nil {
		return nil, err
	}

	cats, err := s.fetchCategories(ctx, id)
	if err != nil {
		return nil, err
	}
	v.Categories = cats
	return &v, nil
}

func (s *vacancyService) GetBySlug(ctx context.Context, slug string) (*vacancy_dto.VacancyResponse, error) {
	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_, _ = s.db.Exec(bgCtx, `
			UPDATE vacancies SET views_count = COALESCE(views_count, 0) + 1 WHERE slug = $1
		`, slug)
	}()

	var v vacancy_dto.VacancyResponse
	row := s.db.QueryRow(ctx, vacancySelectJoin+`WHERE v.slug = $1 AND v.deleted_at IS NULL`, slug)
	if err := scanVacancy(row, &v); err != nil {
		return nil, err
	}

	cats, err := s.fetchCategories(ctx, v.ID)
	if err != nil {
		return nil, err
	}
	v.Categories = cats
	return &v, nil
}

func (s *vacancyService) Update(ctx context.Context, id, userID int64, req vacancy_dto.UpdateVacancyRequest) (*vacancy_dto.VacancyResponse, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	setClauses := []string{"updated_at = NOW()"}
	args := []interface{}{}
	idx := 1

	if req.RegionID != nil {
		setClauses = append(setClauses, fmt.Sprintf("region_id = $%d", idx))
		args = append(args, *req.RegionID)
		idx++
	}
	if req.DistrictID != nil {
		setClauses = append(setClauses, fmt.Sprintf("district_id = $%d", idx))
		args = append(args, *req.DistrictID)
		idx++
	}
	if req.MahallaID != nil {
		setClauses = append(setClauses, fmt.Sprintf("mahalla_id = $%d", idx))
		args = append(args, *req.MahallaID)
		idx++
	}
	if req.Adress != nil {
		setClauses = append(setClauses, fmt.Sprintf("adress = $%d", idx))
		args = append(args, *req.Adress)
		idx++
	}
	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", idx))
		args = append(args, *req.Name)
		idx++
	}
	if req.Title != nil {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", idx))
		args = append(args, *req.Title)
		idx++
	}
	if req.Text != nil {
		setClauses = append(setClauses, fmt.Sprintf("text = $%d", idx))
		args = append(args, *req.Text)
		idx++
	}
	if req.Contact != nil {
		setClauses = append(setClauses, fmt.Sprintf("contact = $%d", idx))
		args = append(args, *req.Contact)
		idx++
	}
	if req.Price != nil {
		setClauses = append(setClauses, fmt.Sprintf("price = $%d", idx))
		args = append(args, *req.Price)
		idx++
	}
	if req.IsActive != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_active = $%d", idx))
		args = append(args, *req.IsActive)
		idx++
	}

	args = append(args, id, userID)
	query := fmt.Sprintf(`
		UPDATE vacancies SET %s
		WHERE id = $%d AND user_id = $%d AND deleted_at IS NULL
		RETURNING id
	`, strings.Join(setClauses, ", "), idx, idx+1)

	var retID int64
	if err := tx.QueryRow(ctx, query, args...).Scan(&retID); err != nil {
		return nil, fmt.Errorf("vacancy not found or access denied")
	}

	if req.CategoryIDs != nil {
		if err := s.replaceCategories(ctx, tx, retID, req.CategoryIDs); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, retID)
}

func (s *vacancyService) Delete(ctx context.Context, id, userID int64) error {
	tag, err := s.db.Exec(ctx, `
		UPDATE vacancies SET deleted_at = $1
		WHERE id = $2 AND user_id = $3 AND deleted_at IS NULL
	`, time.Now(), id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("vacancy not found or access denied")
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

func (s *vacancyService) List(ctx context.Context, f vacancy_dto.VacancyFilter, cursor helper.CursorPayload, limit int) ([]*vacancy_dto.VacancyResponse, bool, int64, error) {
	args := []interface{}{}
	conditions := []string{"v.deleted_at IS NULL"}
	idx := 1

	if f.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("v.user_id = $%d", idx))
		args = append(args, *f.UserID)
		idx++
	}
	if f.RegionID != nil {
		conditions = append(conditions, fmt.Sprintf("v.region_id = $%d", idx))
		args = append(args, *f.RegionID)
		idx++
	}
	if f.DistrictID != nil {
		conditions = append(conditions, fmt.Sprintf("v.district_id = $%d", idx))
		args = append(args, *f.DistrictID)
		idx++
	}
	if f.MahallaID != nil {
		conditions = append(conditions, fmt.Sprintf("v.mahalla_id = $%d", idx))
		args = append(args, *f.MahallaID)
		idx++
	}
	if f.Name != "" {
		conditions = append(conditions, fmt.Sprintf("v.name ILIKE $%d", idx))
		args = append(args, "%"+escapeLike(f.Name)+"%")
		idx++
	}
	if f.Title != "" {
		conditions = append(conditions, fmt.Sprintf("v.title ILIKE $%d", idx))
		args = append(args, "%"+escapeLike(f.Title)+"%")
		idx++
	}
	if f.Search != "" {
		conditions = append(conditions, fmt.Sprintf(
			"(v.name ILIKE $%d OR v.title ILIKE $%d OR v.text ILIKE $%d)",
			idx, idx, idx,
		))
		args = append(args, "%"+escapeLike(f.Search)+"%")
		idx++
	}
	if f.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("v.is_active = $%d", idx))
		args = append(args, *f.IsActive)
		idx++
	}
	if f.MinPrice != nil {
		conditions = append(conditions, fmt.Sprintf("v.price >= $%d", idx))
		args = append(args, *f.MinPrice)
		idx++
	}
	if f.MaxPrice != nil {
		conditions = append(conditions, fmt.Sprintf("v.price <= $%d", idx))
		args = append(args, *f.MaxPrice)
		idx++
	}
	if f.CategoryID != nil {
		conditions = append(conditions, fmt.Sprintf(`EXISTS (
			SELECT 1 FROM category_vacancy cv
			WHERE cv.vacancy_id = v.id AND cv.categorya_id = $%d
		)`, idx))
		args = append(args, *f.CategoryID)
		idx++
	}
	if len(f.CategoryIDs) > 0 {
		conditions = append(conditions, fmt.Sprintf(`EXISTS (
			SELECT 1 FROM category_vacancy cv
			WHERE cv.vacancy_id = v.id AND cv.categorya_id = ANY($%d)
		)`, idx))
		args = append(args, f.CategoryIDs)
		idx++
	}

	// Snapshot for COUNT before adding cursor conditions
	countConditions := make([]string, len(conditions))
	copy(countConditions, conditions)
	countArgs := make([]interface{}, len(args))
	copy(countArgs, args)

	orderCol, valueExpr, cursorFallback, valueKind := vacancyOrderConfig(f.SortBy, f.SortOrder)
	orderDir := normalizedOrder(f.SortOrder)
	conditions, args, idx = appendCursorCondition(conditions, args, idx, "v", orderCol, valueExpr, cursorFallback, valueKind, orderDir, cursor)

	args = append(args, limit+1)
	query := fmt.Sprintf(`
		SELECT v.id, v.slug, v.user_id,
		       v.region_id,   COALESCE(r.name, ''),
		       v.district_id, COALESCE(d.name, ''),
		       v.mahalla_id,  COALESCE(m.name, ''),
		       v.adress, v.name, v.title, v.text, v.contact,
		       v.price, v.views_count, v.is_active,
		       v.created_at, v.updated_at, v.deleted_at
		FROM vacancies v
		LEFT JOIN countries r ON r.id = v.region_id   AND r.deleted_at IS NULL
		LEFT JOIN countries d ON d.id = v.district_id AND d.deleted_at IS NULL
		LEFT JOIN countries m ON m.id = v.mahalla_id  AND m.deleted_at IS NULL
		WHERE %s
		ORDER BY %s %s, v.id %s
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
		countSQL := fmt.Sprintf(`SELECT COUNT(*) FROM vacancies v WHERE %s`, strings.Join(countConditions, " AND "))
		totalErr = s.db.QueryRow(ctx, countSQL, countArgs...).Scan(&total)
	}()

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		wg.Wait()
		return nil, false, 0, err
	}
	defer rows.Close()

	var items []*vacancy_dto.VacancyResponse
	for rows.Next() {
		var v vacancy_dto.VacancyResponse
		if err := rows.Scan(
			&v.ID, &v.Slug, &v.UserID,
			&v.RegionID, &v.RegionName,
			&v.DistrictID, &v.DistrictName,
			&v.MahallaID, &v.MahallaName,
			&v.Adress, &v.Name, &v.Title, &v.Text, &v.Contact,
			&v.Price, &v.ViewsCount, &v.IsActive,
			&v.CreatedAt, &v.UpdatedAt, &v.DeletedAt,
		); err != nil {
			wg.Wait()
			return nil, false, 0, err
		}
		v.Categories = []vacancy_dto.CategoryShort{}
		items = append(items, &v)
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
		indexMap := make(map[int64]int, len(items))
		for i, item := range items {
			ids[i] = item.ID
			indexMap[item.ID] = i
		}

		catRows, err := s.db.Query(ctx, `
			SELECT cv.vacancy_id, c.id, c.name, c.is_active
			FROM category_vacancy cv
			JOIN categories c ON c.id = cv.categorya_id
			WHERE cv.vacancy_id = ANY($1) AND c.deleted_at IS NULL
			ORDER BY cv.vacancy_id, c.id
		`, ids)
		if err == nil {
			defer catRows.Close()
			for catRows.Next() {
				var vacancyID int64
				var cat vacancy_dto.CategoryShort
				if err := catRows.Scan(&vacancyID, &cat.ID, &cat.Name, &cat.IsActive); err == nil {
					if i, ok := indexMap[vacancyID]; ok {
						items[i].Categories = append(items[i].Categories, cat)
					}
				}
			}
		}
	}

	return items, hasMore, total, nil
}

func vacancyOrderConfig(sortBy, sortOrder string) (string, string, string, string) {
	switch sortBy {
	case "price":
		if normalizedOrder(sortOrder) == "ASC" {
			return "COALESCE(v.price, 9223372036854775807)", "price", "9223372036854775807", "int64"
		}
		return "COALESCE(v.price, -1)", "price", "-1", "int64"
	default:
		return "v.id", "", "", ""
	}
}

func normalizedOrder(order string) string {
	if strings.EqualFold(order, "asc") {
		return "ASC"
	}
	return "DESC"
}

func appendCursorCondition(conditions []string, args []interface{}, idx int, idCol, orderCol, valueField, fallbackValue, valueKind, orderDir string, cursor helper.CursorPayload) ([]string, []interface{}, int) {
	if cursor.ID <= 0 {
		return conditions, args, idx
	}
	if orderCol == "v.id" {
		op := "<"
		if orderDir == "ASC" {
			op = ">"
		}
		conditions = append(conditions, fmt.Sprintf("%s.id %s $%d", idCol, op, idx))
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
	conditions = append(conditions, fmt.Sprintf("(%s %s $%d OR (%s = $%d AND %s.id %s $%d))", orderCol, op, idx, orderCol, idx, idCol, op, idx+1))
	switch valueKind {
	case "int64":
		n, err := strconv.ParseInt(cursorValue, 10, 64)
		if err != nil {
			n, _ = strconv.ParseInt(fallbackValue, 10, 64)
		}
		args = append(args, n, cursor.ID)
	default:
		args = append(args, cursorValue, cursor.ID)
	}
	return conditions, args, idx + 2
}
