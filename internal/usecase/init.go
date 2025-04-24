package usecase

import "friendsManagement/internal/repository"

type Usecase struct {
	UserRelationshipUC UserRelationshipUsecase
}

func NewUsecase(userRelationshipRepo repository.UserRelationshipRepository) *Usecase {
	return &Usecase{
		// UserRepo:             userRepo,
		UserRelationshipUC: NewUserRelationshipUsecase(&userRelationshipRepo),
	}
}
