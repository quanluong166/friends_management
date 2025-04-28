package usecase

import (
	"friendsManagement/internal/repository"

	"gorm.io/gorm"
)

type Usecase struct {
	db                 *gorm.DB
	UserRelationshipUC UserRelationshipUsecase
}

func NewUsecase(db *gorm.DB, userRelationshipRepo repository.UserRelationshipRepository) *Usecase {
	return &Usecase{
		UserRelationshipUC: NewUserRelationshipUsecase(db, &userRelationshipRepo),
	}
}
