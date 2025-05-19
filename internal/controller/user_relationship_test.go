package controller_test

import (
	"errors"
	"testing"

	"github.com/quanluong166/friends_management/internal/controller"
	"github.com/quanluong166/friends_management/pkg/helper"
	"github.com/quanluong166/friends_management/pkg/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var mockDB *gorm.DB

func TestUserRealtionshipController_AddFriend(t *testing.T) {
	email1 := "friend1@example.com"
	email2 := "friend2@example.com"
	tcs := map[string]struct {
		err            error
		mockOn         []string
		callArgument   [][]interface{}
		returnArgument [][]interface{}
	}{
		"Error_DatabaseError": {
			callArgument: [][]interface{}{
				{
					email1,
					email2,
				},
			},
			err: errors.New("DATABASE_ERROR"),
			mockOn: []string{
				"CheckTwoUsersBlockedEachOther",
			},
			returnArgument: [][]interface{}{
				{
					false,
					errors.New("DATABASE_ERROR"),
				},
			},
		},
		"Error_OneOfTwoEmailBlockEachOther": {
			callArgument: [][]interface{}{
				{
					email1,
					email2,
				},
			},
			err: errors.New("ONE_OF_YOU_BLOCK_EACH_OTHER"),
			mockOn: []string{
				"CheckTwoUsersBlockedEachOther",
			},
			returnArgument: [][]interface{}{
				{
					true,
					nil,
				},
			},
		},
		"Error_AlreadyFriend": {
			callArgument: [][]interface{}{
				{
					email1,
					email2,
				},
				{
					email1,
					email2,
				},
			},
			err: errors.New("YOU_ALREADY_FRIENDS"),
			mockOn: []string{
				"CheckTwoUsersBlockedEachOther",
				"CheckTwoUsersAreFriends",
			},
			returnArgument: [][]interface{}{
				{
					false,
					nil,
				},
				{
					true,
					nil,
				},
			},
		},
		"Error_CreateFirstFriendshipFailed": {
			callArgument: [][]interface{}{
				{
					email1,
					email2,
				},
				{
					email1,
					email2,
				},
				{
					email1,
					email2,
				},
			},
			err: errors.New("CREATE_FRIST_FRIENDSHIP_RELATION_FAILED: db insert error"),
			mockOn: []string{
				"CheckTwoUsersBlockedEachOther",
				"CheckTwoUsersAreFriends",
				"CreateFriendRelationship",
			},
			returnArgument: [][]interface{}{
				{
					false,
					nil,
				},
				{
					false,
					nil,
				},
				{
					errors.New("db insert error"),
				},
			},
		},
		"Error_CreateSecondFriendshipFailed": {
			callArgument: [][]interface{}{
				{
					email1,
					email2,
				},
				{
					email1,
					email2,
				},
				{
					email1,
					email2,
				},
				{
					email2,
					email1,
				},
			},
			err: errors.New("CREATE_SECOND_FRIENDSHIP_RELATION_FAILED: db insert error"),
			mockOn: []string{
				"CheckTwoUsersBlockedEachOther",
				"CheckTwoUsersAreFriends",
				"CreateFriendRelationship",
				"CreateFriendRelationship",
			},
			returnArgument: [][]interface{}{
				{
					false,
					nil,
				},
				{
					false,
					nil,
				},
				{
					nil,
				},
				{
					errors.New("db insert error"),
				},
			},
		},
		"Success": {
			callArgument: [][]interface{}{
				{
					email1,
					email2,
				},
				{
					email1,
					email2,
				},
				{
					email1,
					email2,
				},
				{
					email2,
					email1,
				},
			},
			err: nil,
			mockOn: []string{
				"CheckTwoUsersBlockedEachOther",
				"CheckTwoUsersAreFriends",
				"CreateFriendRelationship",
				"CreateFriendRelationship",
			},
			returnArgument: [][]interface{}{
				{
					false,
					nil,
				},
				{
					false,
					nil,
				},
				{
					nil,
				},
				{
					nil,
				},
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			mockDB = helper.SetupTestDB(t)
			tx := mockDB.Begin()
			defer tx.Rollback()
			mockRepo := new(controller.MockUserRelationshipRepository)
			for idx, mockName := range tc.mockOn {
				argument := tc.returnArgument[idx]
				callArgument := tc.callArgument[idx]
				mockRepo.On(mockName, callArgument...).Return(argument...)
			}
			ctrl := controller.NewUserRelationshipController(tx, mockRepo)
			err := ctrl.AddFriendship(email1, email2)
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserRealtionshipController_ListFriendships(t *testing.T) {
	input := "test1@example.com"
	expectedFriendEmails := []string{"friend1@example.com", "friend2@example.com"}

	tcs := map[string]struct {
		err            error
		mockOn         []string
		callArgument   [][]interface{}
		returnArgument [][]interface{}
	}{
		"Error_DatabaseError": {
			callArgument: [][]interface{}{
				{
					input,
				},
			},
			err: errors.New("DATABASE_ERROR"),
			mockOn: []string{
				"GetListFriendshipEmail",
			},
			returnArgument: [][]interface{}{
				{
					nil,
					errors.New("DATABASE_ERROR"),
				},
			},
		},
		"Success": {
			callArgument: [][]interface{}{
				{
					input,
				},
			},
			err: nil,
			mockOn: []string{
				"GetListFriendshipEmail",
			},
			returnArgument: [][]interface{}{
				{
					expectedFriendEmails,
					nil,
				},
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			mockRepo := new(controller.MockUserRelationshipRepository)
			for idx, mockName := range tc.mockOn {
				argument := tc.returnArgument[idx]
				callArgument := tc.callArgument[idx]
				mockRepo.On(mockName, callArgument...).Return(argument...)
			}
			ctrl := controller.NewUserRelationshipController(mockDB, mockRepo)
			actualList, actualCount, err := ctrl.ListFriendships(input)
			if tc.err != nil {
				assert.EqualError(t, err, "GET_LIST_FRIENDSHIP_FAIL: "+tc.err.Error())
				assert.Nil(t, actualList)
				assert.Equal(t, int64(0), actualCount)
			} else {
				assert.Equal(t, expectedFriendEmails, actualList)
				assert.Equal(t, int64(len(expectedFriendEmails)), actualCount)
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserRealtionshipController_ListCommonFriends(t *testing.T) {
	email1 := "user1@example.com"
	email2 := "user2@example.com"

	tcs := map[string]struct {
		err            error
		mockOn         []string
		callArgument   [][]interface{}
		returnArgument [][]interface{}
	}{
		"ErrorCheckUserBlockEacOther_DatabaseError": {
			callArgument: [][]interface{}{
				{
					email1,
					email2,
				},
			},
			err: errors.New("CHECK_TWO_USERS_BLOCK_EACH_OTHER_FAIL: DATABASE_ERROR"),
			mockOn: []string{
				"CheckTwoUsersBlockedEachOther",
			},
			returnArgument: [][]interface{}{
				{
					false,
					errors.New("DATABASE_ERROR"),
				},
			},
		},
		"Error_GetListFriendshipFirstEmail_DatabaseError": {
			callArgument: [][]interface{}{
				{
					email1,
					email2,
				},
				{
					email1,
				},
			},
			err: errors.New("GET_LIST_FRIENDSHIP_FOR_FIRST_EMAIL_FAIL: DATABASE_ERROR"),
			mockOn: []string{
				"CheckTwoUsersBlockedEachOther",
				"GetListFriendshipEmail",
			},
			returnArgument: [][]interface{}{
				{
					false,
					nil,
				},
				{
					nil,
					errors.New("DATABASE_ERROR"),
				},
			},
		},
		"Error_GetListFriendshipSecondEmail_DatabaseError": {
			callArgument: [][]interface{}{
				{
					email1,
					email2,
				},
				{
					email1,
				},
				{
					email2,
				},
			},
			err: errors.New("GET_LIST_FRIENDSHIP_FOR_SECOND_EMAIL_FAIL: DATABASE_ERROR"),
			mockOn: []string{
				"CheckTwoUsersBlockedEachOther",
				"GetListFriendshipEmail",
				"GetListFriendshipEmail",
			},
			returnArgument: [][]interface{}{
				{
					false,
					nil,
				},
				{
					[]string{"friend1@example.com", "friend2@example.com", "common@example.com"},
					nil,
				},
				{
					nil,
					errors.New("DATABASE_ERROR"),
				},
			},
		},
		"Error_OneOfYouBlockEachOther": {
			callArgument: [][]interface{}{
				{
					email1,
					email2,
				},
			},
			err: errors.New("ONE_OF_YOU_BLOCK_EACH_OTHER"),
			mockOn: []string{
				"CheckTwoUsersBlockedEachOther",
			},
			returnArgument: [][]interface{}{
				{
					true,
					nil,
				},
			},
		},
		"Success_CommonFriendsFound": {
			callArgument: [][]interface{}{
				{
					email1,
					email2,
				},
				{
					email1,
				},
				{
					email2,
				},
			},
			err: nil,
			mockOn: []string{
				"CheckTwoUsersBlockedEachOther",
				"GetListFriendshipEmail",
				"GetListFriendshipEmail",
			},
			returnArgument: [][]interface{}{
				{
					false,
					nil,
				},
				{
					[]string{"friend1@example.com", "friend2@example.com", "common@example.com"},
					nil,
				},
				{
					[]string{"common@example.com", "friend3@example.com"},
					nil,
				},
			},
		},
		"Success_NoCommonFriends": {
			callArgument: [][]interface{}{
				{
					email1,
					email2,
				},
				{
					email1,
				},
				{
					email2,
				},
			},
			err: nil,
			mockOn: []string{
				"CheckTwoUsersBlockedEachOther",
				"GetListFriendshipEmail",
				"GetListFriendshipEmail",
			},
			returnArgument: [][]interface{}{
				{
					false,
					nil,
				},
				{
					[]string{},
					nil,
				},
				{
					[]string{},
					nil,
				},
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			mockRepo := new(controller.MockUserRelationshipRepository)
			for idx, mockName := range tc.mockOn {
				argument := tc.returnArgument[idx]
				callArgument := tc.callArgument[idx]
				mockRepo.On(mockName, callArgument...).Return(argument...)
			}
			ctrl := controller.NewUserRelationshipController(mockDB, mockRepo)
			actualList, actualCount, err := ctrl.ListCommonFriends(email1, email2)
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
				assert.Nil(t, actualList)
				assert.Equal(t, int64(0), actualCount)
			} else {
				if name == "Success_CommonFriendsFound" {
					assert.Equal(t, []string{"common@example.com"}, actualList)
					assert.Equal(t, int64(1), actualCount)
				} else {
					assert.Empty(t, actualList)
					assert.Equal(t, int64(0), actualCount)
				}
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserRealtionshipController_AddSubscriber(t *testing.T) {
	requestor := "user1@example.com"
	target := "user2@example.com"

	tcs := map[string]struct {
		err            error
		mockOn         []string
		callArgument   [][]interface{}
		returnArgument [][]interface{}
	}{
		"Error_DatabaseError": {
			callArgument: [][]interface{}{
				{
					requestor,
					target,
				},
			},
			err: errors.New("CHECK_IF_THE_REQUESTOR_ALREADY_SUBSCRIBE_FAIL: DATABASE_ERROR"),
			mockOn: []string{
				"CheckIfTheRequestorAlreadySubscribe",
			},
			returnArgument: [][]interface{}{
				{
					false,
					errors.New("DATABASE_ERROR"),
				},
			},
		},
		"Error_AlreadySubscribed": {
			callArgument: [][]interface{}{
				{
					requestor,
					target,
				},
			},
			err: errors.New("YOU_ALREADY_SUBSCRIBED"),
			mockOn: []string{
				"CheckIfTheRequestorAlreadySubscribe",
			},
			returnArgument: [][]interface{}{
				{
					true,
					nil,
				},
			},
		},
		"Error_OneOfYouBlockEachOther": {
			callArgument: [][]interface{}{
				{
					requestor,
					target,
				},
				{
					requestor,
					target,
				},
			},
			err: errors.New("ONE_OF_YOU_BLOCK_EACH_OTHER"),
			mockOn: []string{
				"CheckIfTheRequestorAlreadySubscribe",
				"CheckTwoUsersBlockedEachOther",
			},
			returnArgument: [][]interface{}{
				{
					false,
					nil,
				},
				{
					true,
					nil,
				},
			},
		},
		"Error_CheckTwoUsersBlockedEachOtherFailed": {
			callArgument: [][]interface{}{
				{
					requestor,
					target,
				},
				{
					requestor,
					target,
				},
			},
			err: errors.New("CHECK_TWO_USERS_BLOCK_EACH_OTHER_FAIL: db error"),
			mockOn: []string{
				"CheckIfTheRequestorAlreadySubscribe",
				"CheckTwoUsersBlockedEachOther",
			},
			returnArgument: [][]interface{}{
				{
					false,
					nil,
				},
				{
					false,
					errors.New("db error"),
				},
			},
		},
		"Error_AddSubscriberFailed": {
			callArgument: [][]interface{}{
				{
					requestor,
					target,
				},
				{
					requestor,
					target,
				},
				{
					requestor,
					target,
				},
			},
			err: errors.New("db insert error"),
			mockOn: []string{
				"CheckIfTheRequestorAlreadySubscribe",
				"CheckTwoUsersBlockedEachOther",
				"AddSubscriber",
			},
			returnArgument: [][]interface{}{
				{
					false,
					nil,
				},
				{
					false,
					nil,
				},
				{
					errors.New("db insert error"),
				},
			},
		},
		"Success": {
			callArgument: [][]interface{}{
				{
					requestor,
					target,
				},
				{
					requestor,
					target,
				},
				{
					requestor,
					target,
				},
			},
			err: nil,
			mockOn: []string{
				"CheckIfTheRequestorAlreadySubscribe",
				"CheckTwoUsersBlockedEachOther",
				"AddSubscriber",
			},
			returnArgument: [][]interface{}{
				{
					false,
					nil,
				},
				{
					false,
					nil,
				},
				{
					nil,
				},
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			mockRepo := new(controller.MockUserRelationshipRepository)
			for idx, mockName := range tc.mockOn {
				argument := tc.returnArgument[idx]
				callArgument := tc.callArgument[idx]
				mockRepo.On(mockName, callArgument...).Return(argument...)
			}
			ctrl := controller.NewUserRelationshipController(mockDB, mockRepo)
			err := ctrl.AddSubscriber(requestor, target)
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserRealtionshipController_AddBlock(t *testing.T) {
	requestor := "user1@example.com"
	target := "user2@example.com"
	mockDB = helper.SetupTestDB(t)
	tcs := map[string]struct {
		err            error
		mockOn         []string
		callArgument   [][]interface{}
		returnArgument [][]interface{}
	}{
		"Error_CheckTwoUsersBlockedEachOther_DatabaseError": {
			callArgument: [][]interface{}{
				{
					requestor,
					target,
				},
			},
			err: errors.New("CHECK_TWO_USERS_BLOCK_EACH_OTHER_FAIL: DATABASE_ERROR"),
			mockOn: []string{
				"CheckTwoUsersBlockedEachOther",
			},
			returnArgument: [][]interface{}{
				{
					false,
					errors.New("DATABASE_ERROR"),
				},
			},
		},
		"Error_CheckTwoUsersBlockedEachOther": {
			callArgument: [][]interface{}{
				{
					requestor,
					target,
				},
			},
			err: errors.New("ALREADY_BEEN_BLOCKED"),
			mockOn: []string{
				"CheckTwoUsersBlockedEachOther",
			},
			returnArgument: [][]interface{}{
				{
					true,
					nil,
				},
			},
		},
		"Error_DeleteRelationshipFailed": {
			callArgument: [][]interface{}{
				{
					requestor,
					target,
				},
				{
					requestor,
					target,
				},
			},
			err: errors.New("DELETE_REQUESTOR_RELATIONSHIP_FAIL: db error"),
			mockOn: []string{
				"CheckTwoUsersBlockedEachOther",
				"DeleteRelationship",
			},
			returnArgument: [][]interface{}{
				{
					false,
					nil,
				},
				{
					errors.New("db error"),
				},
			},
		},
		"Error_DeleteTargetRelationshipFailed": {
			callArgument: [][]interface{}{
				{
					requestor,
					target,
				},
				{
					requestor,
					target,
				},
				{
					target,
					requestor,
				},
			},
			err: errors.New("DELETE_TARGET_RELATIONSHIP_FAIL: db error"),
			mockOn: []string{
				"CheckTwoUsersBlockedEachOther",
				"DeleteRelationship",
				"DeleteRelationship",
			},
			returnArgument: [][]interface{}{
				{
					false,
					nil,
				},
				{
					nil,
				},
				{
					errors.New("db error"),
				},
			},
		},
		"Error_CreateBlockRelationshipFailed": {
			callArgument: [][]interface{}{
				{
					requestor,
					target,
				},
				{
					requestor,
					target,
				},
				{
					target,
					requestor,
				},
				{
					requestor,
					target,
				},
			},
			err: errors.New("CREATE_BLOCK_RELATIONSHIP_FAILED: db error"),
			mockOn: []string{
				"CheckTwoUsersBlockedEachOther",
				"DeleteRelationship",
				"DeleteRelationship",
				"CreateBlockRelationship",
			},
			returnArgument: [][]interface{}{
				{
					false,
					nil,
				},
				{
					nil,
				},
				{
					nil,
				},
				{
					errors.New("db error"),
				},
			},
		},
		"Success": {
			callArgument: [][]interface{}{
				{
					requestor,
					target,
				},
				{
					requestor,
					target,
				},
				{
					target,
					requestor,
				},
				{
					requestor,
					target,
				},
			},
			err: nil,
			mockOn: []string{
				"CheckTwoUsersBlockedEachOther",
				"DeleteRelationship",
				"DeleteRelationship",
				"CreateBlockRelationship",
			},
			returnArgument: [][]interface{}{
				{
					false,
					nil,
				},
				{
					nil,
				},
				{
					nil,
				},
				{
					nil,
				},
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			mockRepo := new(controller.MockUserRelationshipRepository)
			tx := mockDB.Begin()
			defer tx.Rollback()
			for idx, mockName := range tc.mockOn {
				argument := tc.returnArgument[idx]
				callArgument := tc.callArgument[idx]
				mockRepo.On(mockName, callArgument...).Return(argument...)
			}
			ctrl := controller.NewUserRelationshipController(tx, mockRepo)
			err := ctrl.AddBlock(requestor, target)
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserRealtionshipController_GetListEmailCanReceiveUpdate(t *testing.T) {
	updaterEmail := "user1@example.com"
	text := "Hello, this is an update!"
	friendEmails := []string{"friend1@example.com", "friend2@example.com"}
	subscriberEmails := []string{"subscriber1@example.com", "subscriber2@example.com"}

	tcs := map[string]struct {
		err            error
		mockOn         []string
		callArgument   [][]interface{}
		returnArgument [][]interface{}
	}{
		"Error_GetListFriendshipEmail_DatabaseError": {
			callArgument: [][]interface{}{
				{
					updaterEmail,
				},
			},
			err: errors.New("GET_LIST_FRIENDSHIP_EMAIL_FAIL: DATABASE_ERROR"),
			mockOn: []string{
				"GetListFriendshipEmail",
			},
			returnArgument: [][]interface{}{
				{
					nil,
					errors.New("DATABASE_ERROR"),
				},
			},
		},
		"Error_GetListSubscriberEmail_DatabaseError": {
			callArgument: [][]interface{}{
				{
					updaterEmail,
				},
				{
					updaterEmail,
				},
			},
			err: errors.New("GET_LIST_SUBSCRIBER_EMAIL_FAIL: DATABASE_ERROR"),
			mockOn: []string{
				"GetListFriendshipEmail",
				"GetListSubscriberEmail",
			},
			returnArgument: [][]interface{}{
				{
					friendEmails,
					nil,
				},
				{
					[]string{},
					errors.New("DATABASE_ERROR"),
				},
			},
		},
		"Success": {
			callArgument: [][]interface{}{
				{
					updaterEmail,
				},
				{
					updaterEmail,
				},
			},
			err: nil,
			mockOn: []string{
				"GetListFriendshipEmail",
				"GetListSubscriberEmail",
			},
			returnArgument: [][]interface{}{
				{
					friendEmails,
					nil,
				},
				{
					subscriberEmails,
					nil,
				},
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			mockRepo := new(controller.MockUserRelationshipRepository)
			for idx, mockName := range tc.mockOn {
				argument := tc.returnArgument[idx]
				callArgument := tc.callArgument[idx]
				mockRepo.On(mockName, callArgument...).Return(argument...)
			}
			ctrl := controller.NewUserRelationshipController(mockDB, mockRepo)
			actualList, err := ctrl.GetListEmailCanReceiveUpdate(updaterEmail, text)
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
				assert.Nil(t, actualList)
			} else {
				assert.Equal(t, utils.Combine(friendEmails, subscriberEmails), actualList)
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

// 	t.Run("Error_DatabaseError", func(t *testing.T) {
// 		mockRepo := new(controller.MockUserRelationshipRepository)
// 		mockRepo.On("GetListFriendshipEmail", updaterEmail).Return(nil, errors.New("DATABASE_ERROR"))

// 		ctrl := controller.NewUserRelationshipController(mockDB, mockRepo)
// 		actualList, err := ctrl.GetListEmailCanReceiveUpdate(updaterEmail, text)
// 		assert.Nil(t, actualList)
// 		assert.EqualError(t, err, "DATABASE_ERROR")
// 		mockRepo.AssertExpectations(t)
// 	})

// 	t.Run("Success", func(t *testing.T) {
// 		mockRepo := new(controller.MockUserRelationshipRepository)
// 		subscriberEmails := []string{"subscriber1@example.com", "subscriber2@example.com"}
// 		friendEmails := []string{"friend1@example.com", "friend2@example.com"}
// 		expectedEmails := []string{"subscriber1@example.com", "subscriber2@example.com", "friend1@example.com", "friend2@example.com"}

// 		mockRepo.On("GetListFriendshipEmail", updaterEmail).Return(friendEmails, nil)
// 		mockRepo.On("GetListSubscriberEmail", updaterEmail).Return(subscriberEmails, nil)

// 		ctrl := controller.NewUserRelationshipController(mockDB, mockRepo)
// 		actualList, err := ctrl.GetListEmailCanReceiveUpdate(updaterEmail, text)
// 		assert.ElementsMatch(t, expectedEmails, actualList)
// 		assert.NoError(t, err)
// 		mockRepo.AssertExpectations(t)
// 	})

// 	t.Run("Success_NoEmailsReturned", func(t *testing.T) {
// 		mockRepo := new(controller.MockUserRelationshipRepository)
// 		mockRepo.On("GetListFriendshipEmail", updaterEmail).Return([]string{}, nil)
// 		mockRepo.On("GetListSubscriberEmail", updaterEmail).Return([]string{}, nil)

// 		ctrl := controller.NewUserRelationshipController(mockDB, mockRepo)
// 		actualList, err := ctrl.GetListEmailCanReceiveUpdate(updaterEmail, text)
// 		assert.Empty(t, actualList)
// 		assert.NoError(t, err)
// 		mockRepo.AssertExpectations(t)
// 	})
// }
