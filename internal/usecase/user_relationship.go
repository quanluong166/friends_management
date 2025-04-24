package usecase

import (
	"errors"
	"friendsManagement/internal/constant"
	"friendsManagement/internal/helper"
	"friendsManagement/internal/repository"

	"gorm.io/gorm"
)

type UserRelationshipUsecase interface {
	AddFriendship(email1, email2 string) error
	CheckRelationshipStatus(email1, email2 string) (string, error)
	ListFriendships(email string) ([]string, int64, error)
	ListCommonFriends(email1, email2 string) ([]string, int64, error)
	AddSubscriber(requestor, target string) error
	AddBlock(requestor, target string) error
}

var (
	invalid_new_friend_status = []string{"BLOCK", "FRIEND"}
)

type userRelationshipUsecase struct {
	userRelationshipRepo repository.UserRelationshipRepository
}

func NewUserRelationshipUsecase(repo *repository.UserRelationshipRepository) UserRelationshipUsecase {
	return &userRelationshipUsecase{
		userRelationshipRepo: *repo,
	}
}

func (uc *userRelationshipUsecase) AddFriendship(email1, email2 string) error {
	//Check if users is blocked or already friends
	status, err := uc.CheckRelationshipStatus(email1, email2)
	if err != nil {
		return err
	}

	isInvalidStatus, status := helper.Contains(invalid_new_friend_status, status)
	if isInvalidStatus {
		return errors.New("you are already been: " + status)
	}

	if status == constant.SUBSCRIBER_STATUS {
		err = uc.userRelationshipRepo.UpdateToFriendship(email1, email2)
	} else {
		err = uc.userRelationshipRepo.AddFriendship(email1, email2)
	}

	return err
}

func (uc *userRelationshipUsecase) ListFriendships(email string) ([]string, int64, error) {
	friendships, err := uc.userRelationshipRepo.ListFriendships(email)
	if err != nil {
		return nil, 0, err
	}
	return friendships, int64(len(friendships)), nil
}

func (uc *userRelationshipUsecase) ListCommonFriends(email1, email2 string) ([]string, int64, error) {
	// Check if users are already friends
	status, err := uc.CheckRelationshipStatus(email1, email2)
	if err != nil {
		return nil, 0, err
	}

	if status != constant.FRIEND_STATUS {
		return nil, 0, errors.New("you are not friends")
	}

	// Get list friendships of two users
	friendships1, err := uc.userRelationshipRepo.ListFriendships(email1)
	if err != nil {
		return nil, 0, err
	}

	friendships2, err := uc.userRelationshipRepo.ListFriendships(email2)
	if err != nil {
		return nil, 0, err
	}

	commonFriends := helper.FindCommon(friendships1, friendships2)
	return commonFriends, int64(len(commonFriends)), nil
}

func (uc *userRelationshipUsecase) AddSubscriber(email1, email2 string) error {
	return uc.userRelationshipRepo.AddSubscriber(email1, email2)
}

func (uc *userRelationshipUsecase) CheckRelationshipStatus(email1, email2 string) (string, error) {
	relationship, err := uc.userRelationshipRepo.GetRelationship(email1, email2)
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}

	if relationship == nil {
		return "", nil
	}

	return relationship.Status, nil
}

func (uc *userRelationshipUsecase) AddBlock(requestor, target string) error {
	//Check if two users has relationship
	status, err := uc.CheckRelationshipStatus(requestor, target)
	if err != nil {
		return err
	}

	/* if exist update the status
	   if not exist create block relationship
	*/
	if status != "" {
		err = uc.userRelationshipRepo.UpdateBlock(requestor, target)
	} else {
		err = uc.userRelationshipRepo.AddBlock(requestor, target)
	}

	return err
}
