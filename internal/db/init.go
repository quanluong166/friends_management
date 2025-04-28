package db

import (
	"fmt"
	"friendsManagement/internal/model"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "admin"),
		getEnv("DB_NAME", "friends_management"),
		getEnv("DB_PORT", "5432"),
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
