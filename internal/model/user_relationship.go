package model

import (
	"time"
)

type UserRelationship struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	RelationshipKey string    `gorm:"uniqueIndex;not null;type:text" json:"email"` //combine two email that need to make relationship
	Status          string    `gorm:"type:text;check:status IN ('FRIEND', 'BLOCKED', 'SUBSCRIER')" json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
