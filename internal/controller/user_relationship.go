package controller

import (
	"errors"

	"gorm.io/gorm"
	"github.com/quanluong166/friends_management/internal/repository"
	"github.com/quanluong166/friends_management/internal/utils"
)

type UserRelationshipController interface {
	AddFriendship(email1, email2 string) error
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
		err := uc.userRelationshipRepo.CreateFriendRelationship(tx, email1, email2)
		if err != nil {
			return errors.New("CREATE_FRIENDSHIP_FAILED: " + err.Error())
		}
		return nil
	})
}

func (uc *userRelationshipController) ListFriendships(email string) ([]string, int64, error) {
	friendships, err := uc.userRelationshipRepo.GetListFriendshipEmail(email)
	if err != nil {
		return nil, 0, err
	}
	return friendships, int64(len(friendships)), nil
}

func (uc *userRelationshipController) ListCommonFriends(email1, email2 string) ([]string, int64, error) {
	isBlock, err := uc.userRelationshipRepo.CheckTwoUsersBlockedEachOther(email1, email2)
	if err != nil {
		return nil, 0, err
	}

	if isBlock {
		return nil, 0, errors.New("ONE_OF_YOU_BLOCK_EACH_OTHER")
	}

	friendships1, err := uc.userRelationshipRepo.GetListFriendshipEmail(email1)
	if err != nil {
		return nil, 0, err
	}

	friendships2, err := uc.userRelationshipRepo.GetListFriendshipEmail(email2)
	if err != nil {
		return nil, 0, err
	}

	commonFriends := utils.FindCommon(friendships1, friendships2)
	return commonFriends, int64(len(commonFriends)), nil
}

func (uc *userRelationshipController) AddSubscriber(requestor, target string) error {
	//Check if user already subcribe
	isSubscribe, err := uc.userRelationshipRepo.CheckIfTheRequestorAlreadySubscribe(requestor, target)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if isSubscribe {
		return errors.New("YOU_ALREADY_SUBSCRIBED")
	}

	//Check if target is blocked by requestor or vice versa
	isBlock, err := uc.userRelationshipRepo.CheckTwoUsersBlockedEachOther(requestor, target)
	if err != nil {
		return err
	}

	if isBlock {
		return errors.New("ONE_OF_YOU_BLOCK_EACH_OTHER")
	}

	return uc.userRelationshipRepo.AddSubscriber(requestor, target)
}

// Delete any existed relationship between the two users and then create a block relationship
func (uc *userRelationshipController) AddBlock(requestor, target string) error {
	return uc.db.Transaction(func(tx *gorm.DB) error {
		err := uc.userRelationshipRepo.DeleteRelationship(tx, requestor, target)
		if err != nil {
			return errors.New("DELETE_FRIENDSHIP_FAILED: " + err.Error())
		}

		err = uc.userRelationshipRepo.CreateBlockRelationship(requestor, target)
		if err != nil {
			return errors.New("CREATE_BLOCK_RELATIONSHIP_FAILED: " + err.Error())
		}
		return nil
	})
}

func (uc *userRelationshipController) GetListEmailCanReceiveUpdate(updaterEmail, text string) ([]string, error) {
	friendships, err := uc.userRelationshipRepo.GetListFriendshipEmail(updaterEmail)
	if err != nil {
		return nil, err
	}

	subscribers, err := uc.userRelationshipRepo.GetListSubscriberEmail(updaterEmail)
	if err != nil {
		return nil, err
	}

	if len(text) == 0 {
		return utils.Combine(friendships, subscribers), nil
	}

	//Get email from text
	mentionedEmails := utils.FindEmails(text)
	return utils.Combine(friendships, subscribers, mentionedEmails), nil
}
