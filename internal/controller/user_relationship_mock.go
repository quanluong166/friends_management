package controller

import (
	"github.com/stretchr/testify/mock"
)

type MockUserRelationshipRepository struct {
	mock.Mock
}

func (m *MockUserRelationshipRepository) CreateFriendRelationship(email1, email2 string) error {
	args := m.Called(email1, email2)
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

func (m *MockUserRelationshipRepository) DeleteRelationship(requestorEmail, targetEmail string) error {
	args := m.Called(requestorEmail, targetEmail)
	return args.Error(0)
}
