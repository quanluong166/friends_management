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
	"github.com/stretchr/testify/mock"
)

type MockUserRelationshipController struct {
	mock.Mock
}

func (m *MockUserRelationshipController) AddFriendship(email1, email2 string) error {
	args := m.Called(email1, email2)
	return args.Error(0)
}

func (m *MockUserRelationshipController) ListFriendships(email string) ([]string, int64, error) {
	args := m.Called(email)
	var friendships []string
	if args.Get(0) != nil {
		friendships = args.Get(0).([]string)
	}

	count := args.Get(1).(int64)

	var err error
	if args.Get(2) != nil {
		err = args.Get(2).(error)
	}

	return friendships, count, err
}

func (m *MockUserRelationshipController) ListCommonFriends(email1, email2 string) ([]string, int64, error) {
	args := m.Called(email1, email2)
	var friendships []string
	if args.Get(0) != nil {
		friendships = args.Get(0).([]string)
	}

	count := args.Get(1).(int64)

	var err error
	if args.Get(2) != nil {
		err = args.Get(2).(error)
	}

	return friendships, count, err
}

func (m *MockUserRelationshipController) AddSubscriber(requestor, target string) error {
	args := m.Called(requestor, target)
	return args.Error(0)
}

func (m *MockUserRelationshipController) AddBlock(requestor, target string) error {
	args := m.Called(requestor, target)
	return args.Error(0)
}

func (m *MockUserRelationshipController) GetListEmailCanReceiveUpdate(senderEmail, text string) ([]string, error) {
	args := m.Called(senderEmail, text)
	var listEmails []string
	if args.Get(0) != nil {
		listEmails = args.Get(0).([]string)
	}

	var err error
	if args.Get(1) != nil {
		err = args.Get(1).(error)
	}

	return listEmails, err
}

func TestUserRelationshipHandler_AddFriend(t *testing.T) {
	// Setup
	e := echo.New()
	t.Run("Success", func(t *testing.T) {
		mockController := new(MockUserRelationshipController)
		mockController.On("AddFriendship", "friend1@example.com", "friend2@example.com").Return(nil)

		svc := &handler.UserRelationshipHandler{
			Controller: mockController,
		}

		reqBody := `{"friends":["friend1@example.com","friend2@example.com"]}`
		req := httptest.NewRequest(http.MethodPost, "/api/user/relationship/add-friend", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Assertions
		if assert.NoError(t, svc.AddFriend(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var resp api.CommonResponse
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.True(t, resp.Success)
			mockController.AssertExpectations(t)
		}
	})

	t.Run("Error_OneOfTwoEmailBlockTheOther", func(t *testing.T) {
		mockController := new(MockUserRelationshipController)
		mockController.On("AddFriendship", "friend1@example.com", "friend2@example.com").Return(errors.New("ONE_OF_YOU_BLOCK_EACH_OTHER"))

		svc := &handler.UserRelationshipHandler{
			Controller: mockController,
		}

		reqBody := `{"friends":["friend1@example.com","friend2@example.com"]}`
		req := httptest.NewRequest(http.MethodPost, "/api/user/relationship/add-friend", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Assertions
		if assert.NoError(t, svc.AddFriend(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)

			var resp api.ErrorResponse
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.False(t, resp.Success)
			assert.Equal(t, "ONE_OF_YOU_BLOCK_EACH_OTHER", resp.Message)
		}
	})

	t.Run("Error_LessThanTwoEmails", func(t *testing.T) {
		mockController := new(MockUserRelationshipController)
		svc := &handler.UserRelationshipHandler{
			Controller: mockController,
		}

		reqBody := `{"friends":["friend1@example.com"]}`
		req := httptest.NewRequest(http.MethodPost, "/api/user/relationship/add-friend", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Assertions
		if assert.NoError(t, svc.AddFriend(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)

			var resp api.ErrorResponse
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.False(t, resp.Success)
			assert.Equal(t, "AT_LEAST_TWO_EMAILS_ARE_REQUIRED", resp.Message)
		}
	})
}

func TestUserRelationshipHandler_ListFriend(t *testing.T) {
	// Setup
	e := echo.New()
	t.Run("Success", func(t *testing.T) {
		mockController := new(MockUserRelationshipController)
		expectedFriends := []string{"friend1@example.com", "friend2@example.com"}
		expectedCount := int64(2)

		mockController.On("ListFriendships", "test@example.com").Return(expectedFriends, expectedCount, nil)

		svc := &handler.UserRelationshipHandler{
			Controller: mockController,
		}

		reqBody := `{"email":"test@example.com"}`
		req := httptest.NewRequest(http.MethodGet, "/api/user/relationship/list-friend", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Assertions
		if assert.NoError(t, svc.ListFriend(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var resp api.ListFriendResponse
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.True(t, resp.Success)
			assert.Equal(t, expectedFriends, resp.Friends)
			assert.Equal(t, int(expectedCount), resp.Count)

			mockController.AssertExpectations(t)
		}
	})

	t.Run("Error_DatabaseError", func(t *testing.T) {
		mockController := new(MockUserRelationshipController)
		expectedFriends := []string{""}
		expectedCount := int64(0)

		mockController.On("ListFriendships", "test@example.com").Return(expectedFriends, expectedCount, errors.New("DATABASE_ERROR"))

		svc := &handler.UserRelationshipHandler{
			Controller: mockController,
		}

		reqBody := `{"email":"test@example.com"}`
		req := httptest.NewRequest(http.MethodGet, "/api/user/relationship/list-friend", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, svc.ListFriend(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)

			var resp api.ErrorResponse
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.False(t, resp.Success)
			assert.Equal(t, "DATABASE_ERROR", resp.Message)
		}
	})
}

func TestUserRelationshipHandler_ListCommonFriends(t *testing.T) {
	// Setup
	e := echo.New()
	t.Run("Success", func(t *testing.T) {
		mockController := new(MockUserRelationshipController)
		expectedCommonFriends := []string{"person1@example.com", "person2@example.com"}
		expectedCount := int64(2)
		mockController.On("ListCommonFriends", "test@example.com", "test2@example.com").Return(expectedCommonFriends, expectedCount, nil)

		svc := &handler.UserRelationshipHandler{
			Controller: mockController,
		}

		reqBody := `{"friends" : ["test@example.com", "test2@example.com"]}`
		req := httptest.NewRequest(http.MethodGet, "/api/user/relationship/list-common-friends", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, svc.ListCommonFriends(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var resp api.ListCommonFriendsResponse
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.True(t, resp.Success)
			assert.Equal(t, expectedCommonFriends, resp.Friends)
			assert.Equal(t, int(expectedCount), resp.Count)
			mockController.AssertExpectations(t)
		}
	})

	t.Run("Error_LessThanTwoEmails", func(t *testing.T) {
		mockController := new(MockUserRelationshipController)
		svc := &handler.UserRelationshipHandler{
			Controller: mockController,
		}

		reqBody := `{"friends" : ["test@example.com"]}`
		req := httptest.NewRequest(http.MethodGet, "/api/user/relationship/list-common-friends", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, svc.ListCommonFriends(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)

			var resp api.ErrorResponse
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.False(t, resp.Success)
			assert.Equal(t, "AT_LEAST_TWO_EMAILS_ARE_REQUIRED", resp.Message)
		}
	})

	t.Run("Error_DatabaseError", func(t *testing.T) {
		mockController := new(MockUserRelationshipController)
		expectedCommonFriends := []string{""}
		expectedCount := int64(0)
		mockController.On("ListCommonFriends", "test@example.com", "test2@example.com").Return(expectedCommonFriends, expectedCount, errors.New("DATABASE_ERROR"))

		svc := &handler.UserRelationshipHandler{
			Controller: mockController,
		}

		reqBody := `{"friends" : ["test@example.com", "test2@example.com"]}`
		req := httptest.NewRequest(http.MethodGet, "/api/user/relationship/list-common-friends", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Assertions
		if assert.NoError(t, svc.ListCommonFriends(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)

			var resp api.ErrorResponse
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.False(t, resp.Success)
			assert.Equal(t, "DATABASE_ERROR", resp.Message)
		}
	})
}

func TestUserRelationshipHandler_AddSubscriber(t *testing.T) {
	// Setup
	e := echo.New()
	t.Run("Success", func(t *testing.T) {
		mockController := new(MockUserRelationshipController)
		mockController.On("AddSubscriber", "test1@example.com", "test2@example.com").Return(nil)
		svc := &handler.UserRelationshipHandler{
			Controller: mockController,
		}

		reqBody := `{"requestor":"test1@example.com","target":"test2@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/api/user/relationship/add-subscriber", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		// Assertions
		if assert.NoError(t, svc.AddSubscriber(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var resp api.CommonResponse
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.True(t, resp.Success)
			mockController.AssertExpectations(t)
		}
	})
	t.Run("Error_EmptyRequestorOrTarget", func(t *testing.T) {
		mockController := new(MockUserRelationshipController)
		mockController.On("AddSubscriber", "", "").Return(errors.New("REQUESTOR_AND_TARGET_ARE_REQUIRED"))
		svc := &handler.UserRelationshipHandler{
			Controller: mockController,
		}

		reqBody := `{"requestor": "","target": "trendy@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/api/user/relationship/add-subscriber", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		// Assertions
		if assert.NoError(t, svc.AddSubscriber(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)

			var resp api.ErrorResponse
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.False(t, resp.Success)
			assert.Equal(t, "REQUESTOR_AND_TARGET_ARE_REQUIRED", resp.Message)
		}
	})
}

func TestUserRelationshipHandler_AddBlock(t *testing.T) {
	// Setup
	e := echo.New()
	t.Run("Success", func(t *testing.T) {
		mockController := new(MockUserRelationshipController)
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
		mockController := new(MockUserRelationshipController)
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
		mockController := new(MockUserRelationshipController)
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
		mockController := new(MockUserRelationshipController)
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
		mockController := new(MockUserRelationshipController)
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
		mockController := new(MockUserRelationshipController)
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
