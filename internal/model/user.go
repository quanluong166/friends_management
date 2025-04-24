package model

import (
	"time"

	"github.com/lib/pq"
)

type User struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Email       string         `gorm:"uniqueIndex;not null" json:"email"`
	Friends     pq.StringArray `gorm:"type:text[]" json:"friends"`     //
	Subscribers pq.StringArray `gorm:"type:text[]" json:"subscribers"` // "
	Blockers    pq.StringArray `gorm:"type:text[]" json:"blockers"`    // "
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
