package helper

import (
	"fmt"
	"log"
	"testing"

	"github.com/quanluong166/friends_management/internal/config"
	"github.com/quanluong166/friends_management/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	c := config.LoadTestDBConfig()
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		c.DBHostTest, c.DBUserTest, c.DBPasswordTest, c.DBNameTest, c.DBPortTest, c.SSLMode, c.TimeZone,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to PostgreSQL: %v", err)
	}

	if err := db.AutoMigrate(&model.UserRelationship{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	return db
}
