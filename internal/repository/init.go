package repository

import "gorm.io/gorm"

type Repository struct {
	UserRelationshipRepo UserRelationshipRepository
}

func NewRepositoy(db *gorm.DB) Repository {
	return Repository{
		UserRelationshipRepo: NewUserRelationshipRepository(db),
	}
}
