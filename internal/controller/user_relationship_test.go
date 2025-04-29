package controller_test

import (
	"errors"
	"friendsManagement/internal/controller"
	"friendsManagement/internal/helper"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockUserRelationshipRepository struct {
	mock.Mock
}

func (m *MockUserRelationshipRepository) CreateFriendRelationship(tx *gorm.DB, email1, email2 string) error {
	args := m.Called(tx, email1, email2)
	return args.Error(0)
}

func (m *MockUserRelationshipRepository) UpdateToFriendship(email1, email2 string) error {
	args := m.Called(email1, email2)
	return args.Error(0)
}

func (m *MockUserRelationshipRepository) GetListSubscriberEmail(target string) ([]string, error) {
	args := m.Called(target)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockUserRelationshipRepository) GetListFriendshipEmail(target string) ([]string, error) {
	args := m.Called(target)
	var friendships []string
	if args.Get(0) != nil {
		friendships = args.Get(0).([]string)
	}
	return friendships, args.Error(1)
}

func (m *MockUserRelationshipRepository) CheckTwoUsersBlockedEachOther(email1, email2 string) (bool, error) {
	args := m.Called(email1, email2)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRelationshipRepository) CheckTwoUsersAreFriends(email1, email2 string) (bool, error) {
	args := m.Called(email1, email2)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRelationshipRepository) AddSubscriber(requestor, target string) error {
	args := m.Called(requestor, target)
	return args.Error(0)
}

func (m *MockUserRelationshipRepository) CreateBlockRelationship(requestor, target string) error {
	args := m.Called(requestor, target)
	return args.Error(0)
}

func (m *MockUserRelationshipRepository) CheckIfTheRequestorAlreadySubscribe(email1, email2 string) (bool, error) {
	args := m.Called(email1, email2)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRelationshipRepository) DeleteRelationship(tx *gorm.DB, requestorEmail, targetEmail string) error {
	args := m.Called(tx, requestorEmail, targetEmail)
	return args.Error(0)
}

var mockDB *gorm.DB

func TestUserRealtionshipController_AddFriend(t *testing.T) {
	email1 := "friend1@example.com"
	email2 := "friend2@example.com"
	t.Run("Error_DatabaseError", func(t *testing.T) {
		mockRepo := new(MockUserRelationshipRepository)
		mockRepo.On("CheckTwoUsersBlockedEachOther", email1, email2).Return(false, errors.New("DATABASE_ERROR"))

		ctrl := controller.NewUserRelationshipController(mockDB, mockRepo)
		err := ctrl.AddFriendship(email1, email2)
		assert.EqualError(t, err, "DATABASE_ERROR")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error_OneOfTwoEmailBlockEachOther", func(t *testing.T) {
		mockRepo := new(MockUserRelationshipRepository)
		mockRepo.On("CheckTwoUsersBlockedEachOther", email1, email2).Return(true, nil)

		ctrl := controller.NewUserRelationshipController(mockDB, mockRepo)
		err := ctrl.AddFriendship(email1, email2)
		assert.EqualError(t, err, "ONE_OF_YOU_BLOCK_EACH_OTHER")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error_AlreadyFriend", func(t *testing.T) {
		mockRepo := new(MockUserRelationshipRepository)
		mockRepo.On("CheckTwoUsersBlockedEachOther", email1, email2).Return(false, nil)
		mockRepo.On("CheckTwoUsersAreFriends", email1, email2).Return(true, nil)

		ctrl := controller.NewUserRelationshipController(mockDB, mockRepo)
		err := ctrl.AddFriendship(email1, email2)
		assert.EqualError(t, err, "YOU_ALREADY_FRIENDS")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		db := helper.SetupTestDB(t)

		tx := db.Begin()
		defer tx.Rollback()
		mockRepo := new(MockUserRelationshipRepository)

		mockRepo.On("CheckTwoUsersBlockedEachOther", email1, email2).Return(false, nil)
		mockRepo.On("CheckTwoUsersAreFriends", email1, email2).Return(false, nil)
		mockRepo.On("CreateFriendRelationship", mock.Anything, email1, email2).Return(nil)

		ctrl := controller.NewUserRelationshipController(tx, mockRepo)
		err := ctrl.AddFriendship(email1, email2)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail", func(t *testing.T) {
		db := helper.SetupTestDB(t)

		tx := db.Begin()
		defer tx.Rollback()
		mockRepo := new(MockUserRelationshipRepository)

		mockRepo.On("CheckTwoUsersBlockedEachOther", email1, email2).Return(false, nil)
		mockRepo.On("CheckTwoUsersAreFriends", email1, email2).Return(false, nil)
		mockRepo.On("CreateFriendRelationship", mock.Anything, email1, email2).Return(errors.New("db insert error"))

		ctrl := controller.NewUserRelationshipController(tx, mockRepo)
		err := ctrl.AddFriendship(email1, email2)
		assert.EqualError(t, err, "CREATE_FRIENDSHIP_FAILED: db insert error")
		mockRepo.AssertExpectations(t)
	})
}

func TestUserRealtionshipController_ListFriendships(t *testing.T) {
	input := "test1@example.com"
	t.Run("Success", func(t *testing.T) {
		expectedFriendEmails := []string{"friend1@example.com", "friend2@example.com"}
		expectedCount := int64(2)
		mockRepo := new(MockUserRelationshipRepository)
		mockRepo.On("GetListFriendshipEmail", input).Return(expectedFriendEmails, nil)
		ctrl := controller.NewUserRelationshipController(mockDB, mockRepo)
		actualList, actualCount, err := ctrl.ListFriendships(input)
		assert.Equal(t, expectedFriendEmails, actualList)
		assert.Equal(t, expectedCount, actualCount)
		assert.NoError(t, err)
	})

	t.Run("Error_DatabaseError", func(t *testing.T) {
		mockRepo := new(MockUserRelationshipRepository)
		mockRepo.On("GetListFriendshipEmail", input).Return(nil, errors.New("DATABASE_ERROR"))
		ctrl := controller.NewUserRelationshipController(mockDB, mockRepo)
		actualList, actualCount, err := ctrl.ListFriendships(input)
		assert.Nil(t, actualList)
		assert.Equal(t, int64(0), actualCount)
		assert.EqualError(t, err, "DATABASE_ERROR")
	})
}

func TestUserRealtionshipController_ListCommonFriends(t *testing.T) {
	email1 := "user1@example.com"
	email2 := "user2@example.com"

	t.Run("Error_DatabaseError", func(t *testing.T) {
		mockRepo := new(MockUserRelationshipRepository)
		mockRepo.On("CheckTwoUsersBlockedEachOther", email1, email2).Return(false, errors.New("DATABASE_ERROR"))

		ctrl := controller.NewUserRelationshipController(mockDB, mockRepo)
		actualList, actualCount, err := ctrl.ListCommonFriends(email1, email2)
		assert.Nil(t, actualList)
		assert.Equal(t, int64(0), actualCount)
		assert.EqualError(t, err, "DATABASE_ERROR")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error_OneOfYouBlockEachOther", func(t *testing.T) {
		mockRepo := new(MockUserRelationshipRepository)
		mockRepo.On("CheckTwoUsersBlockedEachOther", email1, email2).Return(true, nil)

		ctrl := controller.NewUserRelationshipController(mockDB, mockRepo)
		actualList, actualCount, err := ctrl.ListCommonFriends(email1, email2)
		assert.Nil(t, actualList)
		assert.Equal(t, int64(0), actualCount)
		assert.EqualError(t, err, "ONE_OF_YOU_BLOCK_EACH_OTHER")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success_CommonFriendsFound", func(t *testing.T) {
		mockRepo := new(MockUserRelationshipRepository)
		friendships1 := []string{"friend1@example.com", "friend2@example.com", "common@example.com"}
		friendships2 := []string{"common@example.com", "friend3@example.com"}

		mockRepo.On("CheckTwoUsersBlockedEachOther", email1, email2).Return(false, nil)
		mockRepo.On("GetListFriendshipEmail", email1).Return(friendships1, nil)
		mockRepo.On("GetListFriendshipEmail", email2).Return(friendships2, nil)

		ctrl := controller.NewUserRelationshipController(mockDB, mockRepo)
		actualList, actualCount, err := ctrl.ListCommonFriends(email1, email2)
		assert.Equal(t, []string{"common@example.com"}, actualList)
		assert.Equal(t, int64(1), actualCount)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success_NoCommonFriends", func(t *testing.T) {
		mockRepo := new(MockUserRelationshipRepository)
		friendships1 := []string{"friend1@example.com", "friend2@example.com"}
		friendships2 := []string{"friend3@example.com", "friend4@example.com"}

		mockRepo.On("CheckTwoUsersBlockedEachOther", email1, email2).Return(false, nil)
		mockRepo.On("GetListFriendshipEmail", email1).Return(friendships1, nil)
		mockRepo.On("GetListFriendshipEmail", email2).Return(friendships2, nil)

		ctrl := controller.NewUserRelationshipController(mockDB, mockRepo)
		actualList, actualCount, err := ctrl.ListCommonFriends(email1, email2)
		assert.Empty(t, actualList)
		assert.Equal(t, int64(0), actualCount)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserRealtionshipController_AddSubscriber(t *testing.T) {
	requestor := "user1@example.com"
	target := "user2@example.com"

	t.Run("Error_DatabaseError", func(t *testing.T) {
		mockRepo := new(MockUserRelationshipRepository)
		mockRepo.On("CheckIfTheRequestorAlreadySubscribe", requestor, target).Return(false, errors.New("DATABASE_ERROR"))

		ctrl := controller.NewUserRelationshipController(mockDB, mockRepo)
		err := ctrl.AddSubscriber(requestor, target)
		assert.EqualError(t, err, "DATABASE_ERROR")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error_AlreadySubscribed", func(t *testing.T) {
		mockRepo := new(MockUserRelationshipRepository)
		mockRepo.On("CheckIfTheRequestorAlreadySubscribe", requestor, target).Return(true, nil)

		ctrl := controller.NewUserRelationshipController(mockDB, mockRepo)
		err := ctrl.AddSubscriber(requestor, target)
		assert.EqualError(t, err, "YOU_ALREADY_SUBSCRIBED")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error_OneOfYouBlockEachOther", func(t *testing.T) {
		mockRepo := new(MockUserRelationshipRepository)
		mockRepo.On("CheckIfTheRequestorAlreadySubscribe", requestor, target).Return(false, nil)
		mockRepo.On("CheckTwoUsersBlockedEachOther", requestor, target).Return(true, nil)

		ctrl := controller.NewUserRelationshipController(mockDB, mockRepo)
		err := ctrl.AddSubscriber(requestor, target)
		assert.EqualError(t, err, "ONE_OF_YOU_BLOCK_EACH_OTHER")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockUserRelationshipRepository)
		mockRepo.On("CheckIfTheRequestorAlreadySubscribe", requestor, target).Return(false, nil)
		mockRepo.On("CheckTwoUsersBlockedEachOther", requestor, target).Return(false, nil)
		mockRepo.On("AddSubscriber", requestor, target).Return(nil)

		ctrl := controller.NewUserRelationshipController(mockDB, mockRepo)
		err := ctrl.AddSubscriber(requestor, target)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserRealtionshipController_AddBlock(t *testing.T) {
	requestor := "user1@example.com"
	target := "user2@example.com"
	db := helper.SetupTestDB(t)

	tx := db.Begin()
	defer tx.Rollback()

	t.Run("Error_DeleteRelationshipFailed", func(t *testing.T) {
		mockRepo := new(MockUserRelationshipRepository)
		mockRepo.On("DeleteRelationship", mock.Anything, requestor, target).Return(errors.New("DELETE_RELATIONSHIP_FAILED"))

		ctrl := controller.NewUserRelationshipController(tx, mockRepo)
		err := ctrl.AddBlock(requestor, target)
		assert.EqualError(t, err, "DELETE_FRIENDSHIP_FAILED: DELETE_RELATIONSHIP_FAILED")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error_CreateBlockRelationshipFailed", func(t *testing.T) {
		mockRepo := new(MockUserRelationshipRepository)
		mockRepo.On("DeleteRelationship", mock.Anything, requestor, target).Return(nil)
		mockRepo.On("CreateBlockRelationship", requestor, target).Return(errors.New("CREATE_BLOCK_RELATIONSHIP_FAILED"))

		ctrl := controller.NewUserRelationshipController(tx, mockRepo)
		err := ctrl.AddBlock(requestor, target)
		assert.EqualError(t, err, "CREATE_BLOCK_RELATIONSHIP_FAILED: CREATE_BLOCK_RELATIONSHIP_FAILED")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockUserRelationshipRepository)
		mockRepo.On("DeleteRelationship", mock.Anything, requestor, target).Return(nil)
		mockRepo.On("CreateBlockRelationship", requestor, target).Return(nil)

		ctrl := controller.NewUserRelationshipController(tx, mockRepo)
		err := ctrl.AddBlock(requestor, target)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserRealtionshipController_GetListEmailCanReceiveUpdate(t *testing.T) {
	updaterEmail := "user1@example.com"
	text := "Hello, this is an update!"

	t.Run("Error_DatabaseError", func(t *testing.T) {
		mockRepo := new(MockUserRelationshipRepository)
		mockRepo.On("GetListFriendshipEmail", updaterEmail).Return(nil, errors.New("DATABASE_ERROR"))

		ctrl := controller.NewUserRelationshipController(mockDB, mockRepo)
		actualList, err := ctrl.GetListEmailCanReceiveUpdate(updaterEmail, text)
		assert.Nil(t, actualList)
		assert.EqualError(t, err, "DATABASE_ERROR")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockUserRelationshipRepository)
		subscriberEmails := []string{"subscriber1@example.com", "subscriber2@example.com"}
		friendEmails := []string{"friend1@example.com", "friend2@example.com"}
		expectedEmails := []string{"subscriber1@example.com", "subscriber2@example.com", "friend1@example.com", "friend2@example.com"}

		mockRepo.On("GetListFriendshipEmail", updaterEmail).Return(friendEmails, nil)
		mockRepo.On("GetListSubscriberEmail", updaterEmail).Return(subscriberEmails, nil)

		ctrl := controller.NewUserRelationshipController(mockDB, mockRepo)
		actualList, err := ctrl.GetListEmailCanReceiveUpdate(updaterEmail, text)
		assert.ElementsMatch(t, expectedEmails, actualList)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success_NoEmailsReturned", func(t *testing.T) {
		mockRepo := new(MockUserRelationshipRepository)
		mockRepo.On("GetListFriendshipEmail", updaterEmail).Return([]string{}, nil)
		mockRepo.On("GetListSubscriberEmail", updaterEmail).Return([]string{}, nil)

		ctrl := controller.NewUserRelationshipController(mockDB, mockRepo)
		actualList, err := ctrl.GetListEmailCanReceiveUpdate(updaterEmail, text)
		assert.Empty(t, actualList)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
