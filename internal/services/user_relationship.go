package services

import (
	"friendsManagement/internal/services/api"
	"friendsManagement/internal/usecase"

	"github.com/labstack/echo/v4"
)

type userRelationshipService struct {
	usecase usecase.UserRelationshipUsecase
}

func NewUserRelationshipService(uc usecase.UserRelationshipUsecase) api.UserRelationship {
	return &userRelationshipService{usecase: uc}
}

func (h *userRelationshipService) AddFriend(c echo.Context) error {
	//Validate the request
	var req api.AddFriendRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if len(req.Friends) < 2 {
		return echo.NewHTTPError(400, "At least two friends are required")
	}

	// Call the usecase to add friendship
	err := h.usecase.AddFriendship(req.Friends[0], req.Friends[1])
	if err != nil {
		return err
	}
	return c.JSON(200, api.AddFriendResponse{Success: true})
}

func (h *userRelationshipService) ListFriend(c echo.Context) error {
	return nil
}
