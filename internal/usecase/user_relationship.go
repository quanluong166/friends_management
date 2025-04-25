package usecase

import (
	"errors"
	"friendsManagement/internal/constant"
	"friendsManagement/internal/repository"
	"friendsManagement/internal/utils"

	"gorm.io/gorm"
)

type UserRelationshipUsecase interface {
	AddFriendship(email1, email2 string) error
	ListFriendships(email string) ([]string, int64, error)
	ListCommonFriends(email1, email2 string) ([]string, int64, error)
	AddSubscriber(requestor, target string) error
	AddBlock(requestor, target string) error
	GetListEmailCanReceiveUpdate(updaterEmail, text string) ([]string, error)
}

type userRelationshipUsecase struct {
	userRelationshipRepo repository.UserRelationshipRepository
}

func NewUserRelationshipUsecase(repo *repository.UserRelationshipRepository) UserRelationshipUsecase {
	return &userRelationshipUsecase{
		userRelationshipRepo: *repo,
	}
}

func (uc *userRelationshipUsecase) AddFriendship(email1, email2 string) error {
	//Check if two users block each other
	isBlock, err := uc.userRelationshipRepo.CheckTwoUsersBlockedEachOther(email1, email2)
	if err != nil {
		return err
	}

	if isBlock {
		return errors.New("ONE_OF_YOU_BLOCK_EACH_OTHER")
	}

	//Check if two users are already friends
	isFriend, err := uc.userRelationshipRepo.CheckTwoUsersAreFriends(email1, email2)
	if err != nil {
		return err
	}

	if isFriend {
		return errors.New("YOU_ALREADY_FRIENDS")
	}

	err = uc.userRelationshipRepo.CreateFriendRelationship(email1, email2)
	return err
}

func (uc *userRelationshipUsecase) ListFriendships(email string) ([]string, int64, error) {
	friendships, err := uc.userRelationshipRepo.GetListFriendshipEmail(email)
	if err != nil {
		return nil, 0, err
	}
	return friendships, int64(len(friendships)), nil
}

func (uc *userRelationshipUsecase) ListCommonFriends(email1, email2 string) ([]string, int64, error) {
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

func (uc *userRelationshipUsecase) AddSubscriber(requestor, target string) error {
	//Check if user already subcribe
	relationshipType, err := uc.userRelationshipRepo.GetRelationshipType(requestor, target)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if relationshipType == constant.SUBSCRIBER_STATUS {
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

func (uc *userRelationshipUsecase) AddBlock(requestor, target string) error {
	return uc.userRelationshipRepo.CreateBlockRelationship(requestor, target)
}

func (uc *userRelationshipUsecase) GetListEmailCanReceiveUpdate(updaterEmail, text string) ([]string, error) {
	// Get all friends of the user
	friendships, err := uc.userRelationshipRepo.GetListFriendshipEmail(updaterEmail)
	if err != nil {
		return nil, err
	}

	// Get all subscribers of the user
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
