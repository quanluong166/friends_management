package usecase

import "friendsManagement/internal/repository"

type Usecase struct {
	UserRelationshipUC UserRelationshipUsecase
	UserUC             UserUseCase
}

func NewUsecase(userRepo repository.UserRepository, userRelationshipRepo repository.UserRelationshipRepository) *Usecase {
	return &Usecase{
		UserUC:             NewUserUseCase(userRepo),
		UserRelationshipUC: NewUserRelationshipUsecase(&userRelationshipRepo),
	}
}
