package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/quanluong166/friends_management/internal/handler/api"
)

func RegisterUserRelationshipRoutes(e *echo.Echo, userRelationshipService api.UserRelationship) {
	e.POST("/api/user/relationship/friend", userRelationshipService.AddFriend)
	e.POST("/api/user/relationship/subscriber", userRelationshipService.AddSubscriber)
	e.POST("/api/user/relationship/block", userRelationshipService.AddBlock)
	e.POST("/api/user/relationship/list", userRelationshipService.ListFriend)
	e.POST("/api/user/relationship/common-friends", userRelationshipService.ListCommonFriends)
	e.POST("/api/user/relationship/recipients", userRelationshipService.GetListEmailCanReceiveUpdate)
}
