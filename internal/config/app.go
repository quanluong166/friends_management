package config

import (
	"os"
)

type AppConfig struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	TimeZone   string
	SSLMode    string
	PORT       string
}

type TestConfig struct {
	DBHostTest     string
	DBUserTest     string
	DBPasswordTest string
	DBNameTest     string
	DBPortTest     string
	TimeZone       string
	SSLMode        string
	PORT           string
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// LoadConfig support to get application config
func LoadConfig() AppConfig {
	return AppConfig{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "admin"),
		DBName:     getEnv("DB_NAME", "friends_management"),
		DBPort:     getEnv("DB_PORT", "5432"),
		TimeZone:   getEnv("DB_TIMEZONE", "UTC"),
		SSLMode:    getEnv("DB_SSLMODE", "disable"),
		PORT:       getEnv("PORT", ":8080"),
	}
}

// LoadTestDBConfig Support to get test database config for unit test
func LoadTestDBConfig() TestConfig {
	return TestConfig{
		DBHostTest:     getEnv("DB_HOST_TEST", "localhost"),
		DBUserTest:     getEnv("DB_USER_TEST", "postgres"),
		DBPasswordTest: getEnv("DB_PASSWORD_TEST", "admin"),
		DBNameTest:     getEnv("DB_NAME_TEST", "friends_management_test"),
		DBPortTest:     getEnv("DB_PORT_TEST", "5432"),
		TimeZone:       getEnv("DB_TIMEZONE", "UTC"),
		SSLMode:        getEnv("DB_SSLMODE", "disable"),
		PORT:           getEnv("PORT", "8080"),
	}
}
