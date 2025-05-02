package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/quanluong166/friends_management/internal/config"
	"github.com/quanluong166/friends_management/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(c *config.AppConfig) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort, c.SSLMode, c.TimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db

	if err := db.AutoMigrate(&model.UserRelationship{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	return db
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func MigrateUp() error {
	exampleData := []model.UserRelationship{
		{
			RequestorEmail: "mandy@example.com",
			TargetEmail:    "trendy@example.com",
			Type:           "FRIEND",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			RequestorEmail: "trendy@example.com",
			TargetEmail:    "mandy@example.com",
			Type:           "FRIEND",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			RequestorEmail: "trendy@example.com",
			TargetEmail:    "alameda@example.com",
			Type:           "FRIEND",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			RequestorEmail: "alameda@example.com",
			TargetEmail:    "trendy@example.com",
			Type:           "FRIEND",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			RequestorEmail: "alameda@example.com",
			TargetEmail:    "bingo@example.com",
			Type:           "FRIEND",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			RequestorEmail: "bingo@example.com",
			TargetEmail:    "alameda@example.com",
			Type:           "FRIEND",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			RequestorEmail: "bingo@example.com",
			TargetEmail:    "trendy@example.com",
			Type:           "FRIEND",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			RequestorEmail: "trendy@example.com",
			TargetEmail:    "bingo@example.com",
			Type:           "FRIEND",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			RequestorEmail: "leo@example.com",
			TargetEmail:    "trendy@example.com",
			Type:           "BLOCK",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			RequestorEmail: "adison@example.com",
			TargetEmail:    "trendy@example.com",
			Type:           "SUBSCRIBER",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			RequestorEmail: "lucas@example.com",
			TargetEmail:    "trendy@example.com",
			Type:           "SUBSCRIBER",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	for _, data := range exampleData {
		if err := DB.Create(&data).Error; err != nil {
			log.Fatalf("failed to create example data: %v", err)
			return err
		}
	}
	return nil
}
