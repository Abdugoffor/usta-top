package vacancy_service

import (
	"context"
	"fmt"
	"main_service/helper"
	vacancy_dto "main_service/module/vacancy_service/dto"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ─── Interface ───────────────────────────────────────────────────────────────

type VacancyService interface {
	Create(ctx context.Context, userID int64, req vacancy_dto.CreateVacancyRequest) (*vacancy_dto.VacancyResponse, error)
	GetByID(ctx context.Context, id int64) (*vacancy_dto.VacancyResponse, error)
	GetBySlug(ctx context.Context, slug string) (*vacancy_dto.VacancyResponse, error)
	Update(ctx context.Context, id, userID int64, req vacancy_dto.UpdateVacancyRequest) (*vacancy_dto.VacancyResponse, error)
	Delete(ctx context.Context, id, userID int64) error
	List(ctx context.Context, f vacancy_dto.VacancyFilter, page, limit int, sortCol, sortOrder string) ([]*vacancy_dto.VacancyResponse, int64, error)
}

// ─── Implementation ──────────────────────────────────────────────────────────

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

// ─── Create ──────────────────────────────────────────────────────────────────

func (s *vacancyService) Create(ctx context.Context, userID int64, req vacancy_dto.CreateVacancyRequest) (*vacancy_dto.VacancyResponse, error) {
	isActive := true
	{
		if req.IsActive != nil {
			isActive = *req.IsActive
		}
	}

	var id int64
	err := s.db.QueryRow(ctx, `
		INSERT INTO vacancies (user_id, region_id, district_id, mahalla_id, adress, name, title, text, contact, price, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`, userID, req.RegionID, req.DistrictID, req.MahallaID,
		req.Adress, req.Name, req.Title, req.Text, req.Contact,
		req.Price, isActive).Scan(&id)
	if err != nil {
		return nil, err
	}

	slug := helper.Slug(req.Name, id)
	if _, err := s.db.Exec(ctx, `UPDATE vacancies SET slug = $1 WHERE id = $2`, slug, id); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, id)
}

// ─── GetByID ─────────────────────────────────────────────────────────────────

func (s *vacancyService) GetByID(ctx context.Context, id int64) (*vacancy_dto.VacancyResponse, error) {
	var v vacancy_dto.VacancyResponse
	row := s.db.QueryRow(ctx, vacancySelectJoin+`WHERE v.id = $1 AND v.deleted_at IS NULL`, id)
	return &v, scanVacancy(row, &v)
}

// ─── GetBySlug ───────────────────────────────────────────────────────────────

func (s *vacancyService) GetBySlug(ctx context.Context, slug string) (*vacancy_dto.VacancyResponse, error) {
	go s.db.Exec(context.Background(), `
		UPDATE vacancies SET views_count = COALESCE(views_count, 0) + 1 WHERE slug = $1
	`, slug)

	var v vacancy_dto.VacancyResponse
	row := s.db.QueryRow(ctx, vacancySelectJoin+`WHERE v.slug = $1 AND v.deleted_at IS NULL`, slug)
	return &v, scanVacancy(row, &v)
}

// ─── Update ──────────────────────────────────────────────────────────────────

func (s *vacancyService) Update(ctx context.Context, id, userID int64, req vacancy_dto.UpdateVacancyRequest) (*vacancy_dto.VacancyResponse, error) {
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
	if err := s.db.QueryRow(ctx, query, args...).Scan(&retID); err != nil {
		return nil, fmt.Errorf("vacancy not found or access denied")
	}
	return s.GetByID(ctx, retID)
}

// ─── Delete ──────────────────────────────────────────────────────────────────

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

// ─── List ────────────────────────────────────────────────────────────────────

func (s *vacancyService) List(ctx context.Context, f vacancy_dto.VacancyFilter, page, limit int, sortCol, sortOrder string) ([]*vacancy_dto.VacancyResponse, int64, error) {
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
		args = append(args, "%"+f.Name+"%")
		idx++
	}
	if f.Title != "" {
		conditions = append(conditions, fmt.Sprintf("v.title ILIKE $%d", idx))
		args = append(args, "%"+f.Title+"%")
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

	col := validSortCols[sortCol]
	if col == "" {
		col = "v.id"
	}
	if sortOrder != "DESC" {
		sortOrder = "ASC"
	}

	args = append(args, limit, (page-1)*limit)
	query := fmt.Sprintf(`
		SELECT v.id, v.slug, v.user_id,
		       v.region_id,   COALESCE(r.name, ''),
		       v.district_id, COALESCE(d.name, ''),
		       v.mahalla_id,  COALESCE(m.name, ''),
		       v.adress, v.name, v.title, v.text, v.contact,
		       v.price, v.views_count, v.is_active,
		       v.created_at, v.updated_at, v.deleted_at,
		       COUNT(*) OVER() AS total
		FROM vacancies v
		LEFT JOIN countries r ON r.id = v.region_id   AND r.deleted_at IS NULL
		LEFT JOIN countries d ON d.id = v.district_id AND d.deleted_at IS NULL
		LEFT JOIN countries m ON m.id = v.mahalla_id  AND m.deleted_at IS NULL
		WHERE %s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, strings.Join(conditions, " AND "), col, sortOrder, idx, idx+1)

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []*vacancy_dto.VacancyResponse
	var total int64
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
			&total,
		); err != nil {
			return nil, 0, err
		}
		items = append(items, &v)
	}
	return items, total, rows.Err()
}
