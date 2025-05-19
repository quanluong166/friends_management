package handler

import (
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
