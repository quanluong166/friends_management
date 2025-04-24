package repository

import "gorm.io/gorm"

type Repository struct {
	UserRepo             UserRepository
	UserRelationshipRepo UserRelationshipRepository
	// UserUpdates      UserUpdatesRepository
}

func NewRepositoy(db *gorm.DB) *Repository {
	return &Repository{
		UserRepo:             NewUserRepository(db),
		UserRelationshipRepo: NewUserRelationshipRepository(db),
		// UserUpdates:      NewUserUpdatesRepository(),
	}
}
