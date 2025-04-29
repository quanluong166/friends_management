package handler

import (
	"friendsManagement/internal/controller"
	"friendsManagement/internal/handler/api"
)

type Handler struct {
	UserRelationshipHandler api.UserRelationship
}

func NewHandler(userRelationshipController controller.UserRelationshipController) *Handler {
	return &Handler{
		UserRelationshipHandler: NewUserRelationshipHandler(userRelationshipController),
	}
}
