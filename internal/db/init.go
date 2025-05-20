package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/quanluong166/friends_management/internal/config"
	"github.com/quanluong166/friends_management/internal/constant"
	"github.com/quanluong166/friends_management/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// User for init repository functions
var DB *gorm.DB

func InitDB(c config.AppConfig) *gorm.DB {
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

	sqlDB.SetMaxOpenConns(constant.DATABASE_MAX_OPEN_CONNECTION)
	sqlDB.SetMaxIdleConns(constant.DATABASE_MAX_IDLE_CONNECTION)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db

	if err := db.AutoMigrate(&model.UserRelationship{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	return DB
}

// Init sample data fo
func MigrateUp(db *gorm.DB) error {
	dir, _ := os.Getwd()
	sqlFilePath := filepath.Join(dir, "user_relationships_seed.sql")
	data, err := os.ReadFile(sqlFilePath)
	if err != nil {
		return err
	}
	if execErr := db.Exec(string(data)).Error; execErr != nil {
		return execErr
	}
	return nil
}
