package model

import (
	"time"
)

type UserRelationship struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	Email             string    `gorm:"type:text" json:"email"`
	RelationshipEmail string    `gorm:"type:text" json:"relationship_email"`
	Status            string    `gorm:"type:text;check:status IN ('friend', 'blocked', 'subcribe)" json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
