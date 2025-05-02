package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/quanluong166/friends_management/internal/handler/api"
)

func RegisterUserRelationshipRoutes(e *echo.Echo, userRelationshipService api.UserRelationship) {
	e.POST("/api/user/relationship/add-friend", userRelationshipService.AddFriend)
	e.POST("/api/user/relationship/add-subscriber", userRelationshipService.AddSubscriber)
	e.POST("/api/user/relationship/add-block", userRelationshipService.AddBlock)
	e.GET("/api/user/relationship/list-friend", userRelationshipService.ListFriend)
	e.GET("/api/user/relationship/list-common-friends", userRelationshipService.ListCommonFriends)
	e.GET("/api/user/relationship/list-email-can-receive-update", userRelationshipService.GetListEmailCanReceiveUpdate)
}
