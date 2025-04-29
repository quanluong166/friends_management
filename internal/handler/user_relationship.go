package handler

import (
	"friendsManagement/internal/controller"
	"friendsManagement/internal/handler/api"

	"github.com/labstack/echo/v4"
)

type UserRelationshipHandler struct {
	Controller controller.UserRelationshipController
}

func NewUserRelationshipHandler(Controller controller.UserRelationshipController) api.UserRelationship {
	return &UserRelationshipHandler{Controller: Controller}
}

func (sv *UserRelationshipHandler) AddFriend(c echo.Context) error {
	var req api.AddFriendRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, api.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	if len(req.Friends) < 2 {
		return c.JSON(400, api.ErrorResponse{
			Success: false,
			Message: "AT_LEAST_TWO_EMAILS_ARE_REQUIRED",
		})
	}

	err := sv.Controller.AddFriendship(req.Friends[0], req.Friends[1])
	if err != nil {
		return c.JSON(400, api.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(200, api.CommonResponse{Success: true})
}

func (sv *UserRelationshipHandler) ListFriend(c echo.Context) error {
	var req api.ListFriendRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, api.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	friends, count, err := sv.Controller.ListFriendships(req.Email)
	if err != nil {
		return c.JSON(400, api.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(200, api.ListFriendResponse{Success: true, Friends: friends, Count: int(count)})
}

func (sv *UserRelationshipHandler) ListCommonFriends(c echo.Context) error {
	var req api.ListCommonFriendsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, api.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	if len(req.Friends) < 2 {
		return c.JSON(400, api.ErrorResponse{
			Success: false,
			Message: "AT_LEAST_TWO_EMAILS_ARE_REQUIRED",
		})
	}

	commonFriends, count, err := sv.Controller.ListCommonFriends(req.Friends[0], req.Friends[1])
	if err != nil {
		return c.JSON(400, api.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(200, api.ListCommonFriendsResponse{Success: true, Friends: commonFriends, Count: int(count)})
}

func (sv *UserRelationshipHandler) AddSubscriber(c echo.Context) error {
	var req api.AddSubscriberRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, api.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	if len(req.Requestor) == 0 || len(req.Target) == 0 {
		return c.JSON(400, api.ErrorResponse{
			Success: false,
			Message: "REQUESTOR_AND_TARGET_ARE_REQUIRED",
		})
	}

	err := sv.Controller.AddSubscriber(req.Requestor, req.Target)
	if err != nil {
		return c.JSON(400, api.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(200, api.CommonResponse{Success: true})
}

func (sv *UserRelationshipHandler) AddBlock(c echo.Context) error {
	var req api.AddBlockRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, api.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	if len(req.Requestor) == 0 || len(req.Target) == 0 {
		return c.JSON(400, api.ErrorResponse{
			Success: false,
			Message: "REQUESTOR_AND_TARGET_ARE_REQUIRED",
		})
	}

	err := sv.Controller.AddBlock(req.Requestor, req.Target)
	if err != nil {
		return c.JSON(400, api.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(200, api.CommonResponse{Success: true})
}

func (sv *UserRelationshipHandler) GetListEmailCanReceiveUpdate(c echo.Context) error {
	var req api.GetListEmailCanReceiveUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, api.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	if len(req.Sender) == 0 {
		return c.JSON(400, api.ErrorResponse{
			Success: false,
			Message: "SENDER_IS_REQUIRED",
		})
	}

	recipients, err := sv.Controller.GetListEmailCanReceiveUpdate(req.Sender, req.Text)
	if err != nil {
		return c.JSON(400, api.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(200, api.GetListEmailCanReceiveUpdateResponse{Success: true, Recipients: recipients})
}
