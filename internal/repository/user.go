package repository

import (
	"friendsManagement/internal/model"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

type UserRepository interface {
	Create(emails []string) error
	GetUserByEmail(email string) (*model.User, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(emails []string) error {
	for _, email := range emails {
		user := &model.User{
			Email: email,
		}
		if err := r.db.Create(user).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *userRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
