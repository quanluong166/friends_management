package routes

import (
	"friendsManagement/internal/services/api"

	"github.com/labstack/echo/v4"
)

func RegisterUserRelationshipRoutes(router *echo.Group, userRelationshipService api.UserRelationship) {
	router.POST("/relationship/add-friend", userRelationshipService.AddFriend)
	router.GET("/relationship/list-friend", userRelationshipService.ListFriend)
}
