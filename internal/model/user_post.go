package model

import (
	"time"
)

type UpdatePost struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UpdaterEmail string    `gorm:"type:varchar(255);not null" json:"updater_email"`
	Text         string    `gorm:"type:text;not null" json:"content"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}
