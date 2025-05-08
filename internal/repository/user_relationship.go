package repository

import (
	"time"

	"github.com/quanluong166/friends_management/internal/constant"
	"github.com/quanluong166/friends_management/internal/model"
	"gorm.io/gorm"
)

type userRelationshipRepository struct {
	db *gorm.DB
}

// UserRelationshipController all the functions to support operate and manage user relationships
type UserRelationshipRepository interface {
	CreateFriendRelationship(tx *gorm.DB, email1, email2 string) error
	UpdateToFriendship(email1, email2 string) error
	GetListSubscriberEmail(target string) ([]string, error)
	GetListFriendshipEmail(requestor string) ([]string, error)
	AddSubscriber(requestor, target string) error
	CreateBlockRelationship(requestor, target string) error
	CheckTwoUsersBlockedEachOther(email1, email2 string) (bool, error)
	CheckTwoUsersAreFriends(email1, email2 string) (bool, error)
	CheckIfTheRequestorAlreadySubscribe(email1, email2 string) (bool, error)
	DeleteRelationship(tx *gorm.DB, requestorEmail, targetEmail string) error
}

func NewUserRelationshipRepository(db *gorm.DB) UserRelationshipRepository {
	return &userRelationshipRepository{db}
}

// CreateFriendRelationship support create friend connection
func (r *userRelationshipRepository) CreateFriendRelationship(tx *gorm.DB, email1, email2 string) error {
	fristRelationship := &model.UserRelationship{
		RequestorEmail: email1,
		TargetEmail:    email2,
		Type:           constant.FRIEND_RELATIONSHIP_TYPE,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := tx.Create(fristRelationship).Error; err != nil {
		return err
	}

	secondRelationship := &model.UserRelationship{
		RequestorEmail: email2,
		TargetEmail:    email1,
		Type:           constant.FRIEND_RELATIONSHIP_TYPE,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := tx.Create(secondRelationship).Error; err != nil {
		return err
	}
	return nil

}

// GetListSubscriberEmail support query all the subscriber connection of the target email
func (r *userRelationshipRepository) GetListSubscriberEmail(target string) ([]string, error) {
	var relationships []model.UserRelationship
	err := r.db.Where("target_email = ? AND type = ?", target, constant.SUBSCRIBER_RELATIONSHIOP_TYPE).Find(&relationships).Error
	if err != nil {
		return nil, err
	}

	var subscriberEmails []string
	for _, relationship := range relationships {
		subscriberEmails = append(subscriberEmails, relationship.RequestorEmail)
	}

	return subscriberEmails, nil
}

// GetListFriendshipEmail support query all the friend connection of the requestor email
func (r *userRelationshipRepository) GetListFriendshipEmail(requestor string) ([]string, error) {
	var relationships []model.UserRelationship
	err := r.db.Where("requestor_email = ? AND type = ?", requestor, constant.FRIEND_RELATIONSHIP_TYPE).Find(&relationships).Error
	if err != nil {
		return nil, err
	}

	var friendshipEmails []string
	for _, relationship := range relationships {
		friendshipEmails = append(friendshipEmails, relationship.TargetEmail)
	}

	return friendshipEmails, nil
}

// CheckTwoUsersBlockedEachOther support to check whether two email block the other
func (r *userRelationshipRepository) CheckTwoUsersBlockedEachOther(email1, email2 string) (bool, error) {
	var relationships []model.UserRelationship
	err := r.db.Where(`
    (requestor_email = ? AND target_email = ? AND type = ?) OR
    (requestor_email = ? AND target_email = ? AND type = ?)
`, email1, email2, constant.BLOCK_RELATIONSHIP_TYPE, email2, email1, constant.BLOCK_RELATIONSHIP_TYPE).Find(&relationships).Error

	if err != nil {
		return false, err
	}

	if len(relationships) > 0 {
		return true, nil
	}
	return false, nil
}

// CheckTwoUsersAreFriends support to check whether two email are already been friend
func (r *userRelationshipRepository) CheckTwoUsersAreFriends(email1, email2 string) (bool, error) {
	//Since the relationship is bi-directional, we only need to check one direction
	var relationships []model.UserRelationship
	err := r.db.Where("requestor_email = ? AND target_email = ? AND type = ?", email1, email2, constant.FRIEND_RELATIONSHIP_TYPE).Find(&relationships).Error
	if err != nil {
		return false, err
	}

	if len(relationships) > 0 {
		return true, nil
	}
	return false, nil
}

// CheckTwoUsersAreFriends support to update subscriber to friend connection
func (r *userRelationshipRepository) UpdateToFriendship(email1, email2 string) error {
	err := r.db.Model(&model.UserRelationship{}).
		Where("requestor_email = ? AND target_email = ? AND type = ?", email1, email2, constant.SUBSCRIBER_RELATIONSHIOP_TYPE).
		Updates(map[string]interface{}{
			"type":       constant.FRIEND_RELATIONSHIP_TYPE,
			"updated_at": time.Now(),
		}).Error
	if err != nil {
		return err
	}
	return nil
}

// AddSubscriber create subscriber connection
func (r *userRelationshipRepository) AddSubscriber(requestor, target string) error {
	subscription := &model.UserRelationship{
		RequestorEmail: requestor,
		TargetEmail:    target,
		Type:           constant.SUBSCRIBER_RELATIONSHIOP_TYPE,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := r.db.Create(subscription).Error; err != nil {
		return err
	}
	return nil
}

// CreateBlockRelationship create block connection
func (r *userRelationshipRepository) CreateBlockRelationship(requestor, target string) error {
	block := &model.UserRelationship{
		RequestorEmail: requestor,
		TargetEmail:    target,
		Type:           constant.BLOCK_RELATIONSHIP_TYPE,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := r.db.Create(block).Error; err != nil {
		return err
	}
	return nil
}

// CheckIfTheRequestorAlreadySubscribe support to check if the requestor email already a subscriber of the target email
func (r *userRelationshipRepository) CheckIfTheRequestorAlreadySubscribe(requestor, target string) (bool, error) {
	var relationship *model.UserRelationship
	err := r.db.Model(&relationship).
		Where("requestor_email = ? AND target_email = ? AND type = ?", requestor, target, constant.SUBSCRIBER_RELATIONSHIOP_TYPE).First(&relationship).Error
	if err != nil {
		return false, err
	}

	if relationship != nil {
		return true, nil
	}

	return false, nil
}

// DeleteRelationship delete all the connection between the two emails
func (r *userRelationshipRepository) DeleteRelationship(tx *gorm.DB, requestor, target string) error {
	err := tx.Where("requestor_email = ? AND target_email = ?", target, requestor).Delete(&model.UserRelationship{}).Error
	if err != nil {
		return err
	}

	err = tx.Where("requestor_email = ? AND target_email = ?", requestor, target).Delete(&model.UserRelationship{}).Error
	if err != nil {
		return err
	}
	return nil
}
