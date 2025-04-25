package services_test

import (
	"encoding/json"
	"errors"
	"friendsManagement/internal/services/api"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"friendsManagement/internal/services"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRelationshipUsecase struct {
	mock.Mock
}

func (m *MockUserRelationshipUsecase) AddFriendship(email1, email2 string) error {
	args := m.Called(email1, email2)
	return args.Error(0)
}

func (m *MockUserRelationshipUsecase) ListFriendships(email string) ([]string, int64, error) {
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

func (m *MockUserRelationshipUsecase) ListCommonFriends(email1, email2 string) ([]string, int64, error) {
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

func (m *MockUserRelationshipUsecase) AddSubscriber(requestor, target string) error {
	args := m.Called(requestor, target)
	return args.Error(0)
}

func (m *MockUserRelationshipUsecase) AddBlock(requestor, target string) error {
	args := m.Called(requestor, target)
	return args.Error(0)
}

func (m *MockUserRelationshipUsecase) GetListFriendshipEmail(email string) ([]string, error) {
	args := m.Called(email)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockUserRelationshipUsecase) GetListEmailCanReceiveUpdate(updaterEmail, text string) ([]string, error) {
	args := m.Called(updaterEmail, text)
	return args.Get(0).([]string), args.Error(1)
}

func TestUserRelationshipService_AddFriend(t *testing.T) {
	// Setup
	e := echo.New()
	t.Run("Success", func(t *testing.T) {
		mockUsecase := new(MockUserRelationshipUsecase)
		mockUsecase.On("AddFriendship", "friend1@example.com", "friend2@example.com").Return(nil)

		svc := &services.UserRelationshipService{
			Usecase: mockUsecase,
		}

		// Create request
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
			mockUsecase.AssertExpectations(t)
		}
	})

	t.Run("Error_LessThanTwoEmails", func(t *testing.T) {
		mockUsecase := new(MockUserRelationshipUsecase)
		svc := &services.UserRelationshipService{
			Usecase: mockUsecase,
		}

		reqBody := `{"friends":["friend1@example.com"]}`
		req := httptest.NewRequest(http.MethodPost, "/api/user/relationship/add-friend", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Assertions
		if assert.NoError(t, svc.AddFriend(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)

			var resp api.ErrorRespose
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.False(t, resp.Success)
			assert.Equal(t, "AT_LEAST_TWO_EMAILS_ARE_REQUIRED", resp.Message)
		}
	})
}

func TestUserRelationshipService_ListFriend(t *testing.T) {
	// Setup
	e := echo.New()
	t.Run("Success", func(t *testing.T) {
		mockUsecase := new(MockUserRelationshipUsecase)
		expectedFriends := []string{"friend1@example.com", "friend2@example.com"}
		expectedCount := int64(2)

		mockUsecase.On("ListFriendships", "test@example.com").Return(expectedFriends, expectedCount, nil)

		svc := &services.UserRelationshipService{
			Usecase: mockUsecase,
		}

		reqBody := `{"email":"test@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/api/user/relationship/list-friend", strings.NewReader(reqBody))
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

			mockUsecase.AssertExpectations(t)
		}
	})
}

func TestUserRelationshipService_ListCommonFriends(t *testing.T) {
	// Setup
	e := echo.New()
	t.Run("Success", func(t *testing.T) {
		mockUsecase := new(MockUserRelationshipUsecase)
		expectedCommonFriends := []string{"person1@example.com", "person2@example.com"}
		expectedCount := int64(2)
		mockUsecase.On("ListCommonFriends", "test@example.com", "test2@example.com").Return(expectedCommonFriends, expectedCount, nil)

		svc := &services.UserRelationshipService{
			Usecase: mockUsecase,
		}

		reqBody := `{"friends" : ["test@example.com", "test2@example.com"]}`
		req := httptest.NewRequest(http.MethodPost, "/api/user/relationship/list-common-friends", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Assertions
		if assert.NoError(t, svc.ListCommonFriends(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var resp api.ListCommonFriendsResponse
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.True(t, resp.Success)
			assert.Equal(t, expectedCommonFriends, resp.Friends)
			assert.Equal(t, int(expectedCount), resp.Count)
			mockUsecase.AssertExpectations(t)
		}
	})
}

func TestUserRelationshipService_AddSubscriber(t *testing.T) {
	// Setup
	e := echo.New()
	t.Run("Success", func(t *testing.T) {
		mockUsecase := new(MockUserRelationshipUsecase)
		mockUsecase.On("AddSubscriber", "test1@example.com", "test2@example.com").Return(nil)
		svc := &services.UserRelationshipService{
			Usecase: mockUsecase,
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
			mockUsecase.AssertExpectations(t)
		}
	})
	t.Run("Error_EmptyRequestorOrTarget", func(t *testing.T) {
		mockUsecase := new(MockUserRelationshipUsecase)
		mockUsecase.On("AddSubscriber", "", "").Return(errors.New("Requestor and target are required"))
		svc := &services.UserRelationshipService{
			Usecase: mockUsecase,
		}

		reqBody := `{"requestor": "","target": "trendy@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/api/user/relationship/add-subscriber", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		// Assertions
		if assert.NoError(t, svc.AddSubscriber(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)

			var resp api.ErrorRespose
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.False(t, resp.Success)
			assert.Equal(t, "REQUESTOR_AND_TARGET_ARE_REQUIRED", resp.Message)
		}
	})
}
