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

func (sv *userRelationshipService) AddFriend(c echo.Context) error {
	var req api.AddFriendRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, api.ErrorRespose{
			Success: false,
			Message: err.Error(),
		})
	}

	if len(req.Friends) < 2 {
		return c.JSON(400, api.ErrorRespose{
			Success: false,
			Message: "AT_LEAST_TWO_EMAILS_ARE_REQUIRED",
		})
	}

	err := sv.usecase.AddFriendship(req.Friends[0], req.Friends[1])
	if err != nil {
		return c.JSON(400, api.ErrorRespose{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(200, api.CommonResponse{Success: true})
}

func (sv *userRelationshipService) ListFriend(c echo.Context) error {
	var req api.ListFriendRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, api.ErrorRespose{
			Success: false,
			Message: err.Error(),
		})
	}

	friends, count, err := sv.usecase.ListFriendships(req.Email)
	if err != nil {
		return c.JSON(400, api.ErrorRespose{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(200, api.ListFriendResponse{Success: true, Friends: friends, Count: int(count)})
}

func (sv *userRelationshipService) ListCommonFriends(c echo.Context) error {
	var req api.ListCommonFriendsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, api.ErrorRespose{
			Success: false,
			Message: err.Error(),
		})
	}

	if len(req.Friends) < 2 {
		return c.JSON(400, api.ErrorRespose{
			Success: false,
			Message: "AT_LEAST_TWO_EMAILS_ARE_REQUIRED",
		})
	}

	commonFriends, count, err := sv.usecase.ListCommonFriends(req.Friends[0], req.Friends[1])
	if err != nil {
		return c.JSON(400, api.ErrorRespose{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(200, api.ListCommonFriendsResponse{Success: true, Friends: commonFriends, Count: int(count)})
}

func (sv *userRelationshipService) AddSubscriber(c echo.Context) error {
	var req api.AddSubscriberRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, api.ErrorRespose{
			Success: false,
			Message: err.Error(),
		})
	}

	if len(req.Requestor) == 0 || len(req.Target) == 0 {
		return echo.NewHTTPError(400, "Requestor and target are required")
	}

	err := sv.usecase.AddSubscriber(req.Requestor, req.Target)
	if err != nil {
		return c.JSON(400, api.ErrorRespose{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(200, api.CommonResponse{Success: true})
}

func (sv *userRelationshipService) AddBlock(c echo.Context) error {
	var req api.AddBlockRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, api.ErrorRespose{
			Success: false,
			Message: err.Error(),
		})
	}

	if len(req.Requestor) == 0 || len(req.Target) == 0 {
		return echo.NewHTTPError(400, "Requestor and target are required")
	}

	err := sv.usecase.AddBlock(req.Requestor, req.Target)
	if err != nil {
		return c.JSON(400, api.ErrorRespose{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(200, api.CommonResponse{Success: true})
}

func (sv *userRelationshipService) GetListEmailReceiveUpdate(c echo.Context) error {
	var req api.GetListEmailReceiveUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, api.ErrorRespose{
			Success: false,
			Message: err.Error(),
		})
	}

	if len(req.Sender) == 0 {
		return echo.NewHTTPError(400, "Sender and text are required")
	}

	recipients, err := sv.usecase.GetListEmailCanReceiveUpdate(req.Sender, req.Text)
	if err != nil {
		return c.JSON(400, api.ErrorRespose{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(200, api.GetListEmailReceiveUpdateResponse{Success: true, Recipients: recipients})
}
