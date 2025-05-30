package controller

import (
	"errors"

	"github.com/quanluong166/friends_management/internal/repository"
	"github.com/quanluong166/friends_management/pkg/utils"
	"gorm.io/gorm"
)

// UserRelationshipController defines the business logic for managing user relationships
type UserRelationshipController interface {
	AddFriendship(requestor, target string) error
	ListFriendships(email string) ([]string, int64, error)
	ListCommonFriends(email1, email2 string) ([]string, int64, error)
	AddSubscriber(requestor, target string) error
	AddBlock(requestor, target string) error
	GetListEmailCanReceiveUpdate(updaterEmail, text string) ([]string, error)
}

type userRelationshipController struct {
	db                   *gorm.DB
	userRelationshipRepo repository.UserRelationshipRepository
}

func NewUserRelationshipController(db *gorm.DB, repo repository.UserRelationshipRepository) UserRelationshipController {
	return &userRelationshipController{
		userRelationshipRepo: repo,
		db:                   db,
	}
}

// AddFriendship support to create friend connection between two emails
func (uc *userRelationshipController) AddFriendship(email1, email2 string) error {
	isBlock, err := uc.userRelationshipRepo.CheckTwoUsersBlockedEachOther(email1, email2)
	if err != nil {
		return err
	}

	if isBlock {
		return errors.New("ONE_OF_YOU_BLOCK_EACH_OTHER")
	}

	isFriend, err := uc.userRelationshipRepo.CheckTwoUsersAreFriends(email1, email2)
	if err != nil {
		return err
	}

	if isFriend {
		return errors.New("YOU_ALREADY_FRIENDS")
	}

	return uc.db.Transaction(func(tx *gorm.DB) error {
		err := uc.userRelationshipRepo.CreateFriendRelationship(email1, email2)
		if err != nil {
			return errors.New("CREATE_FRIST_FRIENDSHIP_RELATION_FAILED: " + err.Error())
		}

		err = uc.userRelationshipRepo.CreateFriendRelationship(email2, email1)
		if err != nil {
			return errors.New("CREATE_SECOND_FRIENDSHIP_RELATION_FAILED: " + err.Error())
		}
		return nil
	})
}

// ListFriendships support get list friend by email
func (uc *userRelationshipController) ListFriendships(email string) ([]string, int64, error) {
	friendships, err := uc.userRelationshipRepo.GetListFriendshipEmail(email)
	if err != nil {
		return nil, 0, errors.New("GET_LIST_FRIENDSHIP_FAIL: " + err.Error())
	}
	return friendships, int64(len(friendships)), nil
}

// ListCommonFriends support get list common friend between two email
func (uc *userRelationshipController) ListCommonFriends(email1, email2 string) ([]string, int64, error) {
	isBlock, err := uc.userRelationshipRepo.CheckTwoUsersBlockedEachOther(email1, email2)
	if err != nil {
		return nil, 0, errors.New("CHECK_TWO_USERS_BLOCK_EACH_OTHER_FAIL: " + err.Error())
	}

	if isBlock {
		return nil, 0, errors.New("ONE_OF_YOU_BLOCK_EACH_OTHER")
	}

	friendships1, err := uc.userRelationshipRepo.GetListFriendshipEmail(email1)
	if err != nil {
		return nil, 0, errors.New("GET_LIST_FRIENDSHIP_FOR_FIRST_EMAIL_FAIL: " + err.Error())
	}

	friendships2, err := uc.userRelationshipRepo.GetListFriendshipEmail(email2)
	if err != nil {
		return nil, 0, errors.New("GET_LIST_FRIENDSHIP_FOR_SECOND_EMAIL_FAIL: " + err.Error())
	}

	commonFriends := utils.FindCommon(friendships1, friendships2)
	return commonFriends, int64(len(commonFriends)), nil
}

// AddSubscriber support to create and check if two email can make subscibe connection
func (uc *userRelationshipController) AddSubscriber(requestor, target string) error {
	//Check if user already subcribe
	isSubscribe, err := uc.userRelationshipRepo.CheckIfTheRequestorAlreadySubscribe(requestor, target)
	if err != nil && err != gorm.ErrRecordNotFound {
		return errors.New("CHECK_IF_THE_REQUESTOR_ALREADY_SUBSCRIBE_FAIL: " + err.Error())
	}

	if isSubscribe {
		return errors.New("YOU_ALREADY_SUBSCRIBED")
	}

	//Check if target is blocked by requestor or vice versa
	isBlock, err := uc.userRelationshipRepo.CheckTwoUsersBlockedEachOther(requestor, target)
	if err != nil {
		return errors.New("CHECK_TWO_USERS_BLOCK_EACH_OTHER_FAIL: " + err.Error())
	}

	if isBlock {
		return errors.New("ONE_OF_YOU_BLOCK_EACH_OTHER")
	}

	return uc.userRelationshipRepo.AddSubscriber(requestor, target)
}

// AddBlock support create block and delete the other connection between two emails
func (uc *userRelationshipController) AddBlock(requestor, target string) error {
	//Check if target is blocked by requestor or vice versa
	isBlock, err := uc.userRelationshipRepo.CheckTwoUsersBlockedEachOther(requestor, target)
	if err != nil {
		return errors.New("CHECK_TWO_USERS_BLOCK_EACH_OTHER_FAIL: " + err.Error())
	}

	if isBlock {
		return errors.New("ALREADY_BEEN_BLOCKED")
	}

	//Delete all relationship of the requestor and target and then create new block connection
	return uc.db.Transaction(func(tx *gorm.DB) error {
		err := uc.userRelationshipRepo.DeleteRelationship(requestor, target)
		if err != nil {
			return errors.New("DELETE_REQUESTOR_RELATIONSHIP_FAIL: " + err.Error())
		}

		err = uc.userRelationshipRepo.DeleteRelationship(target, requestor)
		if err != nil {
			return errors.New("DELETE_TARGET_RELATIONSHIP_FAIL: " + err.Error())
		}

		err = uc.userRelationshipRepo.CreateBlockRelationship(requestor, target)
		if err != nil {
			return errors.New("CREATE_BLOCK_RELATIONSHIP_FAILED: " + err.Error())
		}
		return nil
	})
}

// GetListEmailCanReceiveUpdate function to support get list of email can receive update from the updater
func (uc *userRelationshipController) GetListEmailCanReceiveUpdate(updaterEmail, text string) ([]string, error) {
	friendships, err := uc.userRelationshipRepo.GetListFriendshipEmail(updaterEmail)
	if err != nil {
		return nil, errors.New("GET_LIST_FRIENDSHIP_EMAIL_FAIL: " + err.Error())
	}

	subscribers, err := uc.userRelationshipRepo.GetListSubscriberEmail(updaterEmail)
	if err != nil {
		return nil, errors.New("GET_LIST_SUBSCRIBER_EMAIL_FAIL: " + err.Error())
	}

	if len(text) == 0 {
		return utils.Combine(friendships, subscribers), nil
	}

	//Get email from text
	mentionedEmails := utils.FindEmails(text)
	return utils.Combine(friendships, subscribers, mentionedEmails), nil
}
