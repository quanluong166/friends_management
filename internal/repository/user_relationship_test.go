package repository_test

import (
	"friendsManagement/internal/helper"
	"friendsManagement/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateFriendRelationship(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := helper.SetupTestDB(t)
		tx := db.Begin()
		defer tx.Rollback()

		repo := repository.NewUserRelationshipRepository(db)

		email1 := "a@example.com"
		email2 := "b@example.com"

		err := repo.CreateFriendRelationship(tx, email1, email2)
		assert.NoError(t, err)
	})
}
