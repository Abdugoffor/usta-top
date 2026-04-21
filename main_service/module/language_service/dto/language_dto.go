package language_dto

import "time"

// ─── Request DTOs ────────────────────────────────────────────────────────────

type CreateLanguageRequest struct {
	Name        string  `json:"name"        validate:"required,min=2,max=255"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

type UpdateLanguageRequest struct {
	Name        *string `json:"name"        validate:"omitempty,min=2,max=255"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

// ─── Filter ──────────────────────────────────────────────────────────────────

type LanguageFilter struct {
	Name     string
	IsActive *bool
}

// ─── Response DTOs ───────────────────────────────────────────────────────────

type LanguageResponse struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
