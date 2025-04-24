package repository

import (
	"friendsManagement/internal/model"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

type UserRepository interface {
	Create(user []*model.User) error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user []*model.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}
