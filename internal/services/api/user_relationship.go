package api

import "github.com/labstack/echo/v4"

type UserRelationship interface {
	AddFriend(c echo.Context) error
	AddSubscriber(c echo.Context) error
	ListFriend(e echo.Context) error
	ListCommonFriends(c echo.Context) error
	AddBlock(c echo.Context) error
}

type AddFriendRequest struct {
	Friends []string `json:"friends"`
}

type CommonResponse struct {
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

type ListCommonFriendsRequest struct {
	Friends []string `json:"friends"`
}

type ListCommonFriendsResponse struct {
	Success bool     `json:"success"`
	Friends []string `json:"friends"`
	Count   int      `json:"count"`
}

type AddSubscriberRequest struct {
	Requestor string `json:"requestor"`
	Target    string `json:"target"`
}

type AddBlockRequest struct {
	Requestor string `json:"requestor"`
	Target    string `json:"target"`
}
