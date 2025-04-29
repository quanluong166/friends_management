package main

import (
	"friendsManagement/internal/config"
	"friendsManagement/internal/db"
	"log"
)

func main() {
	// Initialize the database connection
	config := config.LoadConfig()
	db.InitDB(config)
	if err := db.MigrateUp(); err != nil {
		log.Fatal("Migration failed: ", err.Error())
	}

	log.Println("Migration completed successfully.")
}
