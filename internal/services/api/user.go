package api

import "github.com/labstack/echo/v4"

type FriendHandler interface {
	AddFriend(c echo.Context) (*AddFriendResponse, error)
	ListFriend(e echo.Context) (*ListFriendResponse, error)
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
