package repository

import (
	"gorm.io/gorm"
)

type userRelationshipRepository struct {
	db *gorm.DB
}

type UserRelationshipRepository interface {
	//Create(user []*model.User) error
}

func NewUserRelationshipRepository(db *gorm.DB) UserRelationshipRepository {
	return &userRelationshipRepository{db}
}
