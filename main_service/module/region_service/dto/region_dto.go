package region_dto

import "time"

// ─── Request DTOs ────────────────────────────────────────────────────────────

type CreateRegionRequest struct {
	ParentID    int64  `json:"parent_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

type UpdateRegionRequest struct {
	ParentID    *int64  `json:"parent_id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

// ─── Filter ──────────────────────────────────────────────────────────────────

type RegionFilter struct {
	ID          *int64
	ParentID    *int64
	Name        string
	Description string
	IsActive    *bool
	Page        int
	Limit       int
	SortBy      string
	SortOrder   string
}

// ─── Response DTOs ───────────────────────────────────────────────────────────

type RegionResponse struct {
	ID          int64             `json:"id"`
	ParentID    int64             `json:"parent_id"`
	ParentName  string            `json:"parent_name,omitempty"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	IsActive    bool              `json:"is_active"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	DeletedAt   *time.Time        `json:"deleted_at,omitempty"`
	Children    []*RegionResponse `json:"children,omitempty"`
}

type ListRegionsResponse struct {
	Regions    []*RegionResponse `json:"regions"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int64             `json:"total_pages"`
}
