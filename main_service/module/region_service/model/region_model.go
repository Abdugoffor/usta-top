package region_model

import "time"

type Region struct {
	ID          int64
	ParentID    int64
	Name        string
	Description *string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
