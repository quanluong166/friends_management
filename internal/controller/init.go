package controller

import (
	"github.com/quanluong166/friends_management/internal/repository"
	"gorm.io/gorm"
)

type Controller struct {
	UserRelationshipController UserRelationshipController
}

func NewController(db *gorm.DB, userRelationshipRepo repository.UserRelationshipRepository) Controller {
	return Controller{
		UserRelationshipController: NewUserRelationshipController(db, userRelationshipRepo),
	}
}
