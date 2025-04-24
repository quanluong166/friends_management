package services

import (
	"friendsManagement/internal/services/api"
	"friendsManagement/internal/usecase"
)

type Service struct {
	userService             UserService
	UserRelationshipService api.UserRelationship
}

func NewService(userService UserService, userRelationshipUsecase usecase.UserRelationshipUsecase) *Service {
	return &Service{
		userService:             userService,
		UserRelationshipService: NewUserRelationshipService(userRelationshipUsecase),
	}
}
