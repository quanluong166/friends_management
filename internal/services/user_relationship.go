package services

import (
	"friendsManagement/internal/services/api"
	"friendsManagement/internal/usecase"

	"github.com/labstack/echo/v4"
)

type UserRelationshipService struct {
	Usecase usecase.UserRelationshipUsecase
}

func NewUserRelationshipService(uc usecase.UserRelationshipUsecase) api.UserRelationship {
	return &UserRelationshipService{Usecase: uc}
}

func (sv *UserRelationshipService) AddFriend(c echo.Context) error {
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

	err := sv.Usecase.AddFriendship(req.Friends[0], req.Friends[1])
	if err != nil {
		return c.JSON(400, api.ErrorRespose{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(200, api.CommonResponse{Success: true})
}

func (sv *UserRelationshipService) ListFriend(c echo.Context) error {
	var req api.ListFriendRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, api.ErrorRespose{
			Success: false,
			Message: err.Error(),
		})
	}

	friends, count, err := sv.Usecase.ListFriendships(req.Email)
	if err != nil {
		return c.JSON(400, api.ErrorRespose{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(200, api.ListFriendResponse{Success: true, Friends: friends, Count: int(count)})
}

func (sv *UserRelationshipService) ListCommonFriends(c echo.Context) error {
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

	commonFriends, count, err := sv.Usecase.ListCommonFriends(req.Friends[0], req.Friends[1])
	if err != nil {
		return c.JSON(400, api.ErrorRespose{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(200, api.ListCommonFriendsResponse{Success: true, Friends: commonFriends, Count: int(count)})
}

func (sv *UserRelationshipService) AddSubscriber(c echo.Context) error {
	var req api.AddSubscriberRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, api.ErrorRespose{
			Success: false,
			Message: err.Error(),
		})
	}

	if len(req.Requestor) == 0 || len(req.Target) == 0 {
		return c.JSON(400, api.ErrorRespose{
			Success: false,
			Message: "REQUESTOR_AND_TARGET_ARE_REQUIRED",
		})
	}

	err := sv.Usecase.AddSubscriber(req.Requestor, req.Target)
	if err != nil {
		return c.JSON(400, api.ErrorRespose{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(200, api.CommonResponse{Success: true})
}

func (sv *UserRelationshipService) AddBlock(c echo.Context) error {
	var req api.AddBlockRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, api.ErrorRespose{
			Success: false,
			Message: err.Error(),
		})
	}

	if len(req.Requestor) == 0 || len(req.Target) == 0 {
		return c.JSON(400, api.ErrorRespose{
			Success: false,
			Message: "REQUESTOR_AND_TARGET_ARE_REQUIRED",
		})
	}

	err := sv.Usecase.AddBlock(req.Requestor, req.Target)
	if err != nil {
		return c.JSON(400, api.ErrorRespose{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(200, api.CommonResponse{Success: true})
}

func (sv *UserRelationshipService) GetListEmailCanReceiveUpdate(c echo.Context) error {
	var req api.GetListEmailCanReceiveUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, api.ErrorRespose{
			Success: false,
			Message: err.Error(),
		})
	}

	if len(req.Sender) == 0 {
		return c.JSON(400, api.ErrorRespose{
			Success: false,
			Message: "SENDER_IS_REQUIRED",
		})
	}

	recipients, err := sv.Usecase.GetListEmailCanReceiveUpdate(req.Sender, req.Text)
	if err != nil {
		return c.JSON(400, api.ErrorRespose{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(200, api.GetListEmailCanReceiveUpdateResponse{Success: true, Recipients: recipients})
}
