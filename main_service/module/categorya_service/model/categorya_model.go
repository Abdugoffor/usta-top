package categorya_model

import "time"

type Category struct {
	ID        int64
	Name      *string
	IsActive  *bool
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
