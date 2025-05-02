package handler

import (
	"github.com/quanluong166/friends_management/internal/controller"
	"github.com/quanluong166/friends_management/internal/handler/api"
)

type Handler struct {
	UserRelationshipHandler api.UserRelationship
}

func NewHandler(userRelationshipController controller.UserRelationshipController) *Handler {
	return &Handler{
		UserRelationshipHandler: NewUserRelationshipHandler(userRelationshipController),
	}
}
