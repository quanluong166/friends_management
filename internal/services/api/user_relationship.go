package api

import "github.com/labstack/echo/v4"

type UserRelationship interface {
	AddFriend(c echo.Context) error
	ListFriend(e echo.Context) error
}

type AddFriendRequest struct {
	Friends []string `json:"friends"`
}

type AddFriendResponse struct {
	Success bool `json:"success"`
}

type ListFriendRequest struct {
	Email string `json:"email"`
}

type ListFriendResponse struct {
	Success bool     `json:"success"`
	Friends []string `json:"friends"`
	Count   int      `json:"count"`
}
