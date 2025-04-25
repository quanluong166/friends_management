package services

import (
	"friendsManagement/internal/services/api"
	"friendsManagement/internal/usecase"
)

type Service struct {
	UserRelationshipService api.UserRelationship
}

func NewService(userRelationshipUsecase usecase.UserRelationshipUsecase) *Service {
	return &Service{
		UserRelationshipService: NewUserRelationshipService(userRelationshipUsecase),
	}
}
