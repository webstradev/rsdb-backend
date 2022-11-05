package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/webstradev/rsdb-backend/db"
)

func main() {
	// Load environment variables
	setupEnvironment()

	dbService, err := db.Setup(os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		log.Fatal(err)
	}

	dbService.Ping()
}

func setupEnvironment() {
	// If a database connection string is not yet set in environment variables load the .env file
	if os.Getenv("DB_CONNECTION_STRING") == "" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal(err)
		}
	}
}
