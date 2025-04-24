package services

import "friendsManagement/internal/usecase"

type UserService interface {
	Create(emails []string) error
}

type userService struct {
	userUC usecase.UserUseCase
}

func NewUserService(uc usecase.UserUseCase) UserService {
	return &userService{
		userUC: uc,
	}
}

func (sv *userService) Create(emails []string) error {
	err := sv.userUC.Create(emails)
	if err != nil {
		return err
	}
	return nil
}
