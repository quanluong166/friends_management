package usecase

import (
	"friendsManagement/internal/repository"
)

type UserUseCase interface {
	Create(emails []string) error
}

type userUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(repo repository.UserRepository) UserUseCase {
	return &userUseCase{userRepo: repo}
}

func (u *userUseCase) Create(emails []string) error {
	err := u.userRepo.Create(emails)
	if err != nil {
		return err
	}
	return nil
}
