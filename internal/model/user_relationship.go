package model

import (
	"time"
)

type UserRelationship struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	RequestorEmail string    `gorm:"type:varchar(255);not null" json:"requestor_email"`
	TargetEmail    string    `gorm:"type:varchar(255);not null" json:"target_email"`
	Type           string    `gorm:"type:text;check:type IN ('FRIEND', 'BLOCK', 'SUBSCRIBER')" json:"type"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
