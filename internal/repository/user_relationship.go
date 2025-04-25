package repository

import (
	"fmt"
	"friendsManagement/internal/model"
	"time"

	"gorm.io/gorm"
)

type userRelationshipRepository struct {
	db *gorm.DB
}

type UserRelationshipRepository interface {
	CreateFriendRelationship(email1, email2 string) error
	UpdateToFriendship(email1, email2 string) error
	GetListBlockEmail(target string) ([]string, error)
	GetListSubscriberEmail(target string) ([]string, error)
	GetListFriendshipEmail(requestor string) ([]string, error)
	AddSubscriber(requestor, target string) error
	CreateBlockRelationship(requestor, target string) error
	CheckTwoUsersBlockedEachOther(email1, email2 string) (bool, error)
	CheckTwoUsersAreFriends(email1, email2 string) (bool, error)
	GetRelationshipType(email1, email2 string) (string, error)
}

func NewUserRelationshipRepository(db *gorm.DB) UserRelationshipRepository {
	return &userRelationshipRepository{db}
}

func (r *userRelationshipRepository) CreateFriendRelationship(email1, email2 string) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	fristRelationship := &model.UserRelationship{
		RequestorEmail: email1,
		TargetEmail:    email2,
		Type:           "FRIEND",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := tx.Create(fristRelationship).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create friendship from %s to %s: %w", email1, email2, err)
	}

	secondRelationship := &model.UserRelationship{
		RequestorEmail: email2,
		TargetEmail:    email1,
		Type:           "FRIEND",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := tx.Create(secondRelationship).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create friendship from %s to %s: %w", email2, email1, err)
	}
	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (r *userRelationshipRepository) GetListBlockEmail(target string) ([]string, error) {
	var relationships []model.UserRelationship
	err := r.db.Where("target_email = ? AND type = ?", target, "BLOCK").Find(&relationships).Error
	if err != nil {
		return nil, err
	}

	var blockEmails []string
	for _, relationship := range relationships {
		blockEmails = append(blockEmails, relationship.RequestorEmail)
	}

	return blockEmails, nil
}

func (r *userRelationshipRepository) GetListSubscriberEmail(target string) ([]string, error) {
	var relationships []model.UserRelationship
	err := r.db.Where("target_email = ? AND type = ?", target, "SUBSCRIBER").Find(&relationships).Error
	if err != nil {
		return nil, err
	}

	var subscriberEmails []string
	for _, relationship := range relationships {
		subscriberEmails = append(subscriberEmails, relationship.RequestorEmail)
	}

	return subscriberEmails, nil
}

func (r *userRelationshipRepository) GetListFriendshipEmail(requestor string) ([]string, error) {
	var relationships []model.UserRelationship
	err := r.db.Where("requestor_email = ? AND type = ?", requestor, "FRIEND").Find(&relationships).Error
	if err != nil {
		return nil, err
	}

	var friendshipEmails []string
	for _, relationship := range relationships {
		friendshipEmails = append(friendshipEmails, relationship.TargetEmail)
	}

	return friendshipEmails, nil
}

func (r *userRelationshipRepository) CheckTwoUsersBlockedEachOther(email1, email2 string) (bool, error) {
	var relationships []model.UserRelationship
	err := r.db.Where(`
    (requestor_email = ? AND target_email = ? AND type = 'BLOCK') OR
    (requestor_email = ? AND target_email = ? AND type = 'BLOCK')
`, email1, email2, email2, email1).Find(&relationships).Error

	if err != nil {
		return false, err
	}

	if len(relationships) > 0 {
		return true, nil
	}
	return false, nil
}

func (r *userRelationshipRepository) CheckTwoUsersAreFriends(email1, email2 string) (bool, error) {
	//Since the relationship is bi-directional, we only need to check one direction
	var relationships []model.UserRelationship
	err := r.db.Where("requestor_email = ? AND target_email = ? AND type = ?", email1, email2, "FRIEND").Find(&relationships).Error
	if err != nil {
		return false, err
	}

	if len(relationships) > 0 {
		return true, nil
	}
	return false, nil
}

func (r *userRelationshipRepository) UpdateToFriendship(email1, email2 string) error {
	err := r.db.Model(&model.UserRelationship{}).
		Where("requestor_email = ? AND target_email = ? AND type = ?", email1, email2, "SUBSCRIBER").
		Updates(map[string]interface{}{
			"type":       "FRIEND",
			"updated_at": time.Now(),
		}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *userRelationshipRepository) AddSubscriber(requestor, target string) error {
	subscription := &model.UserRelationship{
		RequestorEmail: requestor,
		TargetEmail:    target,
		Type:           "SUBSCRIBER",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := r.db.Create(subscription).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRelationshipRepository) CreateBlockRelationship(requestor, target string) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	/*
		Delete any existed relationship between the two users
		This will delete both FRIEND and SUBSCRIBER relationships
	*/
	err := tx.Where("requestor_email = ? AND target_email = ?", target, requestor).Delete(&model.UserRelationship{}).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete friendship relationship: %w", err)
	}

	err = tx.Where("requestor_email = ? AND target_email = ?", requestor, target).Delete(&model.UserRelationship{}).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete friendship relationship: %w", err)
	}

	//Create the block relationship
	block := &model.UserRelationship{
		RequestorEmail: requestor,
		TargetEmail:    target,
		Type:           "BLOCK",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := tx.Create(block).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create block relationship: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (r *userRelationshipRepository) GetRelationshipType(requestor, target string) (string, error) {
	var relationship model.UserRelationship
	err := r.db.Model(&relationship).
		Where("requestor_email = ? AND target_email = ?", requestor, target).First(&relationship).Error
	if err != nil {
		return "", err
	}

	return relationship.Type, nil
}
