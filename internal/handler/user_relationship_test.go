package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/quanluong166/friends_management/internal/handler"
	"github.com/quanluong166/friends_management/internal/handler/api"
	"github.com/quanluong166/friends_management/pkg/utils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUserRelationshipHandler_AddFriend(t *testing.T) {
	// Setup
	e := echo.New()
	tcs := map[string]struct {
		email1         string
		email2         string
		err            error
		mockOn         []string
		callArgument   [][]interface{}
		returnArgument [][]interface{}
	}{
		"Success": {
			email1: "friend1@example.com",
			email2: "friend2@example.com",
			err:    nil,
			mockOn: []string{"AddFriendship"},
			callArgument: [][]interface{}{
				{"friend1@example.com", "friend2@example.com"},
			},
			returnArgument: [][]interface{}{
				{nil},
			},
		},
		"Error_AtLeastTwoEmailsAreRequired": {
			email2:         "friend2@example.com",
			err:            errors.New("AT_LEAST_TWO_EMAILS_ARE_REQUIRED"),
			mockOn:         []string{},
			callArgument:   [][]interface{}{},
			returnArgument: [][]interface{}{},
		},
		"Error_InvalidEmail": {
			email1: "invalid-email",
			email2: "friend2@example.com",
			err:    errors.New("INVALID_EMAIL_INPUT"),
			mockOn: []string{},
			callArgument: [][]interface{}{
				{"invalid-email", "friend2@example.com"},
			},
			returnArgument: [][]interface{}{},
		},
		"Error_AddFriendshipFailed": {
			email1: "friend1@example.com",
			email2: "friend2@example.com",
			err:    errors.New("DATABASE_ERROR"),
			mockOn: []string{"AddFriendship"},
			callArgument: [][]interface{}{
				{"friend1@example.com", "friend2@example.com"},
			},
			returnArgument: [][]interface{}{
				{errors.New("DATABASE_ERROR")},
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			mockController := new(handler.MockUserRelationshipController)
			for i, method := range tc.mockOn {
				mockController.On(method, tc.callArgument[i]...).Return(tc.returnArgument[i]...)
			}
			svc := &handler.UserRelationshipHandler{
				Controller: mockController,
			}
			reqBody, _ := buildRequestBody(tc.email1, tc.email2)
			req := httptest.NewRequest(http.MethodPost, "/api/user/relationship/add-friend", strings.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if assert.NoError(t, svc.AddFriend(c)) {
				if tc.err != nil {
					var resp api.ErrorResponse
					err := json.Unmarshal(rec.Body.Bytes(), &resp)
					assert.NoError(t, err)
					assert.Equal(t, http.StatusBadRequest, rec.Code)
					assert.Equal(t, tc.err.Error(), resp.Message)
					assert.False(t, resp.Success)
				} else {
					var resp api.CommonResponse
					err := json.Unmarshal(rec.Body.Bytes(), &resp)
					assert.Equal(t, http.StatusOK, rec.Code)
					assert.NoError(t, err)
					assert.True(t, resp.Success)
					mockController.AssertExpectations(t)
				}
			}
		})
	}
}

func TestUserRelationshipHandler_ListFriend(t *testing.T) {
	// Setup
	e := echo.New()
	expectedFriends := []string{"friend1@example.com", "friend2@example.com"}

	tcs := map[string]struct {
		email          string
		err            error
		mockOn         []string
		callArgument   [][]interface{}
		returnArgument [][]interface{}
	}{
		"Success": {
			email:          "test@example.com",
			mockOn:         []string{"ListFriendships"},
			callArgument:   [][]interface{}{{"test@example.com"}},
			returnArgument: [][]interface{}{{expectedFriends, int64(2), nil}},
			err:            nil,
		},
		"Error_InvalidEmail": {
			email:          "invalid-email",
			mockOn:         []string{},
			callArgument:   [][]interface{}{},
			returnArgument: [][]interface{}{},
			err:            errors.New("INVALID_EMAIL_INPUT"),
		},
		"Error_DatabaseError": {
			email:          "test@example.com",
			mockOn:         []string{"ListFriendships"},
			callArgument:   [][]interface{}{{"test@example.com"}},
			returnArgument: [][]interface{}{{nil, int64(0), errors.New("DATABASE_ERROR")}},
			err:            errors.New("DATABASE_ERROR"),
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			mockController := new(handler.MockUserRelationshipController)
			for i, method := range tc.mockOn {
				mockController.On(method, tc.callArgument[i]...).Return(tc.returnArgument[i]...)
			}
			svc := &handler.UserRelationshipHandler{
				Controller: mockController,
			}
			reqBody := `{"email":"` + tc.email + `"}`
			req := httptest.NewRequest(http.MethodGet, "/api/user/relationship/list-friend", strings.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if assert.NoError(t, svc.ListFriend(c)) {
				if tc.err != nil {
					var resp api.ErrorResponse
					err := json.Unmarshal(rec.Body.Bytes(), &resp)
					assert.NoError(t, err)
					assert.Equal(t, http.StatusBadRequest, rec.Code)
					assert.Equal(t, tc.err.Error(), resp.Message)
					assert.False(t, resp.Success)
				} else {
					var resp api.ListFriendResponse
					err := json.Unmarshal(rec.Body.Bytes(), &resp)
					assert.NoError(t, err)
					assert.Equal(t, http.StatusOK, rec.Code)
					assert.True(t, resp.Success)
					assert.Equal(t, expectedFriends, resp.Friends)
					assert.Equal(t, int(2), resp.Count)
				}
			}
			mockController.AssertExpectations(t)
		})
	}
}

func TestUserRelationshipHandler_ListCommonFriends(t *testing.T) {
	// Setup
	e := echo.New()
	expectedCommonFriends := []string{"person1@example.com", "person2@example.com"}
	tcs := map[string]struct {
		email1         string
		email2         string
		err            error
		mockOn         []string
		callArgument   [][]interface{}
		returnArgument [][]interface{}
	}{
		"Success": {
			email1:         "test@example.com",
			email2:         "test2@example.com",
			mockOn:         []string{"ListCommonFriends"},
			callArgument:   [][]interface{}{{"test@example.com", "test2@example.com"}},
			returnArgument: [][]interface{}{{expectedCommonFriends, int64(2), nil}},
			err:            nil,
		},
		"Error_AtLeastTwoEmailsAreRequired": {
			email2:         "test2@example.com",
			mockOn:         []string{},
			callArgument:   [][]interface{}{},
			returnArgument: [][]interface{}{},
			err:            errors.New("AT_LEAST_TWO_EMAILS_ARE_REQUIRED"),
		},
		"Error_InvalidEmail": {
			email1:         "invalid-email",
			email2:         "test2@example.com",
			mockOn:         []string{},
			callArgument:   [][]interface{}{},
			returnArgument: [][]interface{}{},
			err:            errors.New("INVALID_EMAIL_INPUT"),
		},
		"Error_DatabaseError": {
			email1:         "test1@example.com",
			email2:         "test2@example.com",
			mockOn:         []string{"ListCommonFriends"},
			callArgument:   [][]interface{}{{"test1@example.com", "test2@example.com"}},
			returnArgument: [][]interface{}{{nil, int64(0), errors.New("DATABASE_ERROR")}},
			err:            errors.New("DATABASE_ERROR"),
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			mockController := new(handler.MockUserRelationshipController)
			for i, method := range tc.mockOn {
				mockController.On(method, tc.callArgument[i]...).Return(tc.returnArgument[i]...)
			}
			svc := &handler.UserRelationshipHandler{
				Controller: mockController,
			}
			reqBody, _ := buildRequestBody(tc.email1, tc.email2)
			req := httptest.NewRequest(http.MethodGet, "/api/user/relationship/list-common-friends", strings.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if assert.NoError(t, svc.ListCommonFriends(c)) {
				if tc.err != nil {
					var resp api.ErrorResponse
					err := json.Unmarshal(rec.Body.Bytes(), &resp)
					assert.NoError(t, err)
					assert.Equal(t, http.StatusBadRequest, rec.Code)
					assert.Equal(t, tc.err.Error(), resp.Message)
					assert.False(t, resp.Success)
				} else {
					var resp api.ListCommonFriendsResponse
					err := json.Unmarshal(rec.Body.Bytes(), &resp)
					assert.NoError(t, err)
					assert.Equal(t, http.StatusOK, rec.Code)
					assert.True(t, resp.Success)
					assert.Equal(t, expectedCommonFriends, resp.Friends)
					assert.Equal(t, int(2), resp.Count)
				}
			}
			mockController.AssertExpectations(t)
		})
	}
}

func TestUserRelationshipHandler_AddSubscriber(t *testing.T) {
	// Setup
	e := echo.New()
	tcs := map[string]struct {
		requestor      string
		target         string
		err            error
		mockOn         []string
		callArgument   [][]interface{}
		returnArgument [][]interface{}
	}{
		"Success": {
			requestor:      "test1@example.com",
			target:         "test2@example.com",
			err:            nil,
			mockOn:         []string{"AddSubscriber"},
			callArgument:   [][]interface{}{{"test1@example.com", "test2@example.com"}},
			returnArgument: [][]interface{}{{nil}},
		},
		"Error_EmptyRequestorOrTarget": {
			requestor:      "",
			target:         "",
			err:            errors.New("REQUESTOR_AND_TARGET_ARE_REQUIRED"),
			mockOn:         []string{"AddSubscriber"},
			callArgument:   [][]interface{}{{"", ""}},
			returnArgument: [][]interface{}{{errors.New("REQUESTOR_AND_TARGET_ARE_REQUIRED")}},
		},
		"Error_InvalidEmail": {
			requestor:      "invalid-email",
			target:         "test2@example.com",
			err:            errors.New("INVALID_EMAIL_INPUT"),
			mockOn:         []string{},
			callArgument:   [][]interface{}{{}},
			returnArgument: [][]interface{}{},
		},
		"Error_AddSubscriberFailed": {
			requestor:      "test1@example.com",
			target:         "test2@example.com",
			err:            errors.New("DATABASE_ERROR"),
			mockOn:         []string{"AddSubscriber"},
			callArgument:   [][]interface{}{{"test1@example.com", "test2@example.com"}},
			returnArgument: [][]interface{}{{errors.New("DATABASE_ERROR")}},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			mockController := new(handler.MockUserRelationshipController)
			for i, method := range tc.mockOn {
				mockController.On(method, tc.callArgument[i]...).Return(tc.returnArgument[i]...)
			}
			svc := &handler.UserRelationshipHandler{
				Controller: mockController,
			}
			reqBody := `{"requestor":"` + tc.requestor + `","target":"` + tc.target + `"}`
			req := httptest.NewRequest(http.MethodPost, "/api/user/relationship/add-subscriber", strings.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if assert.NoError(t, svc.AddSubscriber(c)) {
				if tc.err != nil {
					var resp api.ErrorResponse
					err := json.Unmarshal(rec.Body.Bytes(), &resp)
					assert.NoError(t, err)
					assert.Equal(t, http.StatusBadRequest, rec.Code)
					assert.Equal(t, tc.err.Error(), resp.Message)
					assert.False(t, resp.Success)
				} else {
					var resp api.CommonResponse
					err := json.Unmarshal(rec.Body.Bytes(), &resp)
					assert.NoError(t, err)
					assert.Equal(t, http.StatusOK, rec.Code)
					assert.True(t, resp.Success)
					mockController.AssertExpectations(t)
				}
			}
		})
	}
}

func TestUserRelationshipHandler_AddBlock(t *testing.T) {
	// Setup
	e := echo.New()
	t.Run("Success", func(t *testing.T) {
		mockController := new(handler.MockUserRelationshipController)
		mockController.On("AddBlock", "test1@example.com", "test2@example.com").Return(nil)
		svc := &handler.UserRelationshipHandler{
			Controller: mockController,
		}

		reqBody := `{"requestor":"test1@example.com","target":"test2@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/api/user/relationship/add-block", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		// Assertions
		if assert.NoError(t, svc.AddBlock(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var resp api.CommonResponse
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.True(t, resp.Success)
			mockController.AssertExpectations(t)
		}
	})
	t.Run("Error_EmptyRequestorOrTarget", func(t *testing.T) {
		mockController := new(handler.MockUserRelationshipController)
		svc := &handler.UserRelationshipHandler{
			Controller: mockController,
		}

		reqBody := `{"requestor": "trendy@example.com","target": ""}`
		req := httptest.NewRequest(http.MethodPost, "/api/user/relationship/add-block", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		// Assertions
		if assert.NoError(t, svc.AddBlock(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)

			var resp api.ErrorResponse
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.False(t, resp.Success)
			assert.Equal(t, "REQUESTOR_AND_TARGET_ARE_REQUIRED", resp.Message)
		}
	})

	t.Run("Error_DatabaseError", func(t *testing.T) {
		mockController := new(handler.MockUserRelationshipController)
		mockController.On("AddBlock", "test1@example.com", "test2@example.com").Return(errors.New("DATABASE_ERROR"))
		svc := &handler.UserRelationshipHandler{
			Controller: mockController,
		}

		reqBody := `{"requestor":"test1@example.com","target":"test2@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/api/user/relationship/add-block", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		// Assertions
		if assert.NoError(t, svc.AddBlock(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)

			var resp api.ErrorResponse
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.False(t, resp.Success)
			assert.Equal(t, "DATABASE_ERROR", resp.Message)
		}
	})
}

func TestUserRelationshipHandler_GetListEmailCanReceiveUpdate(t *testing.T) {
	// Setup
	e := echo.New()
	t.Run("Success", func(t *testing.T) {
		mockController := new(handler.MockUserRelationshipController)
		expectedFriendEmails := []string{"friend1@example.com", "friend2@example.com"}
		expectedSubscriberEmails := []string{"subscriber1@example.com", "subscriber2@example.com"}
		listEmailsFromMention := []string{"mention1@example.com", "mention2@example.com"}
		expectedListRecipients := utils.Combine(expectedFriendEmails, expectedSubscriberEmails, listEmailsFromMention)
		text := "Hello mention1@example.com mention2@example.com"
		senderEmail := "test1@example.com"

		mockController.On("GetListEmailCanReceiveUpdate", senderEmail, text).Return(expectedListRecipients, nil)
		svc := &handler.UserRelationshipHandler{
			Controller: mockController,
		}

		reqBody := `{"sender" : "test1@example.com","text": "Hello mention1@example.com mention2@example.com"}`
		req := httptest.NewRequest(http.MethodGet, "/api/user/relationship/get-list-email-receive-update", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, svc.GetListEmailCanReceiveUpdate(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var resp api.GetListEmailCanReceiveUpdateResponse
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.True(t, resp.Success)
			assert.Equal(t, expectedListRecipients, resp.Recipients)
			mockController.AssertExpectations(t)
		}
	})

	t.Run("Error_EmptySenderEmail", func(t *testing.T) {
		mockController := new(handler.MockUserRelationshipController)
		mockController.On("GetListEmailCanReceiveUpdate", "", "Hello person1@example.com").Return(errors.New("SENDER_IS_REQUIRED"))
		svc := &handler.UserRelationshipHandler{
			Controller: mockController,
		}

		reqBody := `{"sender": "","text": "Hello person1@example.com"}`
		req := httptest.NewRequest(http.MethodGet, "/api/user/relationship/get-list-email-receive-update", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		// Assertions
		if assert.NoError(t, svc.GetListEmailCanReceiveUpdate(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)

			var resp api.ErrorResponse
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.False(t, resp.Success)
			assert.Equal(t, "SENDER_IS_REQUIRED", resp.Message)
		}
	})

	t.Run("Error_DatabaseError", func(t *testing.T) {
		mockController := new(handler.MockUserRelationshipController)
		mockController.On("GetListEmailCanReceiveUpdate", "test1@example.com", "Hello person1@example.com").Return([]string{}, errors.New("DATABASE_ERROR"))
		svc := &handler.UserRelationshipHandler{
			Controller: mockController,
		}

		reqBody := `{"sender": "test1@example.com","text": "Hello person1@example.com"}`
		req := httptest.NewRequest(http.MethodGet, "/api/user/relationship/get-list-email-receive-update", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		// Assertions
		if assert.NoError(t, svc.GetListEmailCanReceiveUpdate(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)

			var resp api.ErrorResponse
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.False(t, resp.Success)
			assert.Equal(t, "DATABASE_ERROR", resp.Message)
		}
	})
}

type Request struct {
	Friends []string `json:"friends"`
}

func buildRequestBody(friend1, friend2 string) (string, error) {
	friends := []string{}
	if friend1 != "" {
		friends = append(friends, friend1)
	}
	if friend2 != "" {
		friends = append(friends, friend2)
	}

	req := Request{
		Friends: friends,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
