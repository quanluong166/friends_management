package api

import "github.com/labstack/echo/v4"

type UserRelationship interface {
	AddFriend(c echo.Context) error
	AddSubscriber(c echo.Context) error
	ListFriend(e echo.Context) error
	ListCommonFriends(c echo.Context) error
	AddBlock(c echo.Context) error
	GetListEmailCanReceiveUpdate(c echo.Context) error
}

// AddFriendRequest is the request body for add friend API
type AddFriendRequest struct {
	Friends []string `json:"friends"`
}

// CommonResponse is the response body all API
type CommonResponse struct {
	Success bool `json:"success"`
}

// ErrorResponse is the error response body for all API

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ListFriendRequest is the request body for list friend API
type ListFriendRequest struct {
	Email string `json:"email"`
}

// ListFriendResponse is the response body for list friend AP
type ListFriendResponse struct {
	Success bool     `json:"success"`
	Friends []string `json:"friends"`
	Count   int      `json:"count"`
}

// ListCommonFriendsRequest is the request body for list common friends API
type ListCommonFriendsRequest struct {
	Friends []string `json:"friends"`
}

// ListCommonFriendsResponse is the response body for list common friends API
type ListCommonFriendsResponse struct {
	Success bool     `json:"success"`
	Friends []string `json:"friends"`
	Count   int      `json:"count"`
}

// AddSubscriberRequest is the request body for add subscriber API
type AddSubscriberRequest struct {
	Requestor string `json:"requestor"`
	Target    string `json:"target"`
}

// AddBlockRequest is the request body for add block API
type AddBlockRequest struct {
	Requestor string `json:"requestor"`
	Target    string `json:"target"`
}

// GetListEmailCanReceiveUpdateRequest is the request body for get list recipient API
type GetListEmailCanReceiveUpdateRequest struct {
	Sender string `json:"sender"`
	Text   string `json:"text"`
}

// GetListEmailCanReceiveUpdateResponse is the response body for get list recipient API
type GetListEmailCanReceiveUpdateResponse struct {
	Success    bool     `json:"success"`
	Recipients []string `json:"recipients"`
}
