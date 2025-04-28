package main

import (
	"friendsManagement/internal/db"
	"log"
)

func main() {
	// Initialize the database connection
	db.InitDB()
	if err := db.MigrateUp(); err != nil {
		log.Fatal("Migration failed: ", err.Error())
	}

	log.Println("Migration completed successfully.")
}
