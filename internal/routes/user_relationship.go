package routes

import (
	"friendsManagement/internal/handler/api"

	"github.com/labstack/echo/v4"
)

func RegisterUserRelationshipRoutes(e *echo.Echo, userRelationshipService api.UserRelationship) {
	e.POST("/api/user/relationship/add-friend", userRelationshipService.AddFriend)
	e.POST("/api/user/relationship/add-subscriber", userRelationshipService.AddSubscriber)
	e.POST("/api/user/relationship/add-block", userRelationshipService.AddBlock)
	e.GET("/api/user/relationship/list-friend", userRelationshipService.ListFriend)
	e.GET("/api/user/relationship/list-common-friends", userRelationshipService.ListCommonFriends)
	e.GET("/api/user/relationship/list-email-can-receive-update", userRelationshipService.GetListEmailCanReceiveUpdate)
}
