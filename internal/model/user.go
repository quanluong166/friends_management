package model

type User struct {
	ID        string `gorm:"type:uuid;primaryKey" json:"id"`
	Email     string `gorm:"uniqueIndex;not null" json:"email"`
	CreatedAt string `json:"created_at"`
}
