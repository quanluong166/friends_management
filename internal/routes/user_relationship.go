package routes

import (
	"friendsManagement/internal/services/api"

	"github.com/labstack/echo/v4"
)

func RegisterUserRelationshipRoutes(router *echo.Group, userRelationshipService api.UserRelationship) {
	router.POST("/relationship/add-friend", userRelationshipService.AddFriend)
	router.POST("/relationship/add-subscriber", userRelationshipService.AddSubscriber)
	router.POST("/relationship/add-block", userRelationshipService.AddBlock)
	router.GET("/relationship/list-friend", userRelationshipService.ListFriend)
	router.GET("/relationship/list-common-friends", userRelationshipService.ListCommonFriends)
	router.GET("/relationship/get-list-email-can-receive-update", userRelationshipService.GetListEmailCanReceiveUpdate)
}
