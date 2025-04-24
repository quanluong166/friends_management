package services

import (
	"friendsManagement/internal/services/api"
	"friendsManagement/internal/usecase"

	"github.com/labstack/echo/v4"
)

type friendHandler struct {
	usecase *usecase.FriendUsecase
}

func NewFriendHandler(uc *usecase.FriendUsecase) api.FriendHandler {
	return &friendHandler{usecase: uc}
}

func (h *friendHandler) AddFriend(c echo.Context) (*api.AddFriendResponse, error) {
	// implementation remains the same
	return nil, nil
}

func (h *friendHandler) ListFriend(c echo.Context) (*api.ListFriendResponse, error) {
	return nil, nil
}
