package repository

import (
	"friendsManagement/internal/helper"
	"friendsManagement/internal/model"
	"time"

	"gorm.io/gorm"
)

type userRelationshipRepository struct {
	db *gorm.DB
}

type UserRelationshipRepository interface {
	AddFriendship(email1, email2 string) error
	UpdateToFriendship(email1, email2 string) error
	GetRelationship(email1, email2 string) (*model.UserRelationship, error)
	ListFriendships(email string) ([]string, error)
	AddSubscriber(requestor, target string) error
	UpdateBlock(requestor, target string) error
	AddBlock(requestor, target string) error
}

func NewUserRelationshipRepository(db *gorm.DB) UserRelationshipRepository {
	return &userRelationshipRepository{db}
}

func (r *userRelationshipRepository) AddFriendship(email1, email2 string) error {
	// Create a new friendship record
	relationshipKey := helper.GenerateRelationshipKey(email1, email2)
	friendship := &model.UserRelationship{
		RelationshipKey: relationshipKey,
		Status:          "FRIEND",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := r.db.Create(friendship).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRelationshipRepository) UpdateToFriendship(email1, email2 string) error {
	// Update the status of the relationship to "friend"
	relationshipKey := helper.GenerateRelationshipKey(email1, email2)
	if err := r.db.Model(&model.UserRelationship{}).
		Where("relationship_key = ?", relationshipKey).
		Update("status", "FRIEND").Error; err != nil {
		return err
	}
	return nil
}

func (r *userRelationshipRepository) ListFriendships(email string) ([]string, error) {
	var friendships []model.UserRelationship
	if err := r.db.Where("relationship_key LIKE ?", "%"+email+"%").Find(&friendships).Error; err != nil {
		return nil, err
	}

	var friendEmails []string
	for _, friendship := range friendships {
		friendEmail := helper.GetOtherEmailFromKey(friendship.RelationshipKey, email)
		friendEmails = append(friendEmails, friendEmail)
	}

	return friendEmails, nil
}

func (r *userRelationshipRepository) AddSubscriber(email1, email2 string) error {
	// Create a new subscription record
	relationshipKey := helper.GenerateRelationshipKey(email1, email2)
	subscription := &model.UserRelationship{
		RelationshipKey: relationshipKey,
		Status:          "SUBSCRIBER",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := r.db.Create(subscription).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRelationshipRepository) UpdateBlock(email1, email2 string) error {
	// Update the status of the relationship to "blocked"
	relationshipKey := helper.GenerateRelationshipKey(email1, email2)
	if err := r.db.Model(&model.UserRelationship{}).
		Where("relationship_key = ?", relationshipKey).
		Update("status", "BLOCKED").Error; err != nil {
		return err
	}
	return nil
}

func (r *userRelationshipRepository) GetRelationship(email1, email2 string) (*model.UserRelationship, error) {
	relationshipKey := helper.GenerateRelationshipKey(email1, email2)
	var relationship model.UserRelationship
	if err := r.db.Where("relationship_key = ?", relationshipKey).First(&relationship).Error; err != nil {
		return nil, err
	}
	return &relationship, nil
}

func (r *userRelationshipRepository) AddBlock(email1, email2 string) error {
	relationshipKey := helper.GenerateRelationshipKey(email1, email2)
	block := &model.UserRelationship{
		RelationshipKey: relationshipKey,
		Status:          "BLOCKED",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := r.db.Create(block).Error; err != nil {
		return err
	}
	return nil
}
