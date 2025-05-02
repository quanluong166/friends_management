package controller

import (
	"gorm.io/gorm"
	"github.com/quanluong166/friends_management/internal/repository"
)

type Controller struct {
	UserRelationshipController UserRelationshipController
}

func NewController(db *gorm.DB, userRelationshipRepo repository.UserRelationshipRepository) *Controller {
	return &Controller{
		UserRelationshipController: NewUserRelationshipController(db, userRelationshipRepo),
	}
}
