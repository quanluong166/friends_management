package repository

import "gorm.io/gorm"

type Repository struct {
	UserRelationshipRepo UserRelationshipRepository
	// UserUpdates      UserUpdatesRepository
}

func NewRepositoy(db *gorm.DB) *Repository {
	return &Repository{
		UserRelationshipRepo: NewUserRelationshipRepository(db),
		// UserUpdates:      NewUserUpdatesRepository(),
	}
}
