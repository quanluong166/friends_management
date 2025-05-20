package main

import (
	"log"

	"github.com/quanluong166/friends_management/internal/config"
	"github.com/quanluong166/friends_management/internal/db"
)

func main() {
	// Initialize the database connection
	config := config.LoadConfig()
	db.InitDB(config)
	if err := db.MigrateUp(db.DB); err != nil {
		log.Fatal("Migration failed: ", err.Error())
	}

	log.Println("Migration completed successfully.")
}
