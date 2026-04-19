package country_dto

import "time"

// ─── Request DTOs ────────────────────────────────────────────────────────────

type CreateCountryRequest struct {
	ParentID *int64 `json:"parent_id"`
	Name     string `json:"name"      validate:"required,min=2,max=255"`
	IsActive *bool  `json:"is_active"`
}

type UpdateCountryRequest struct {
	ParentID *int64  `json:"parent_id"`
	Name     *string `json:"name"      validate:"omitempty,min=2,max=255"`
	IsActive *bool   `json:"is_active"`
}

// ─── Filter ──────────────────────────────────────────────────────────────────

type CountryFilter struct {
	ParentID *int64
	Name     string
	IsActive *bool
}

// ─── Response DTOs ───────────────────────────────────────────────────────────

type CountryResponse struct {
	ID         int64             `json:"id"`
	ParentID   *int64            `json:"parent_id"`
	ParentName string            `json:"parent_name,omitempty"`
	Name       string            `json:"name"`
	IsActive   bool              `json:"is_active"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	DeletedAt  *time.Time        `json:"deleted_at,omitempty"`
	Children   []*CountryResponse `json:"children,omitempty"`
}
