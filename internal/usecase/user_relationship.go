package usecase

import (
	"friendsManagement/internal/repository"
)

type UserRelationshipUsecase interface {
	AddFriendship(email1, email2 string) error
}

type userRelationshipUsecase struct {
	userRelationshipRepo repository.UserRelationshipRepository
}

func NewUserRelationshipUsecase(repo *repository.UserRelationshipRepository) UserRelationshipUsecase {
	return &userRelationshipUsecase{
		userRelationshipRepo: *repo,
	}
}

func (u *userRelationshipUsecase) AddFriendship(email1, email2 string) error {
	err := u.userRelationshipRepo.AddFriendship(email1, email2)
	if err != nil {
		return err
	}
	return nil
}
