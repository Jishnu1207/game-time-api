package main

import (
	"game-time-api/config"
	"game-time-api/migrations"
	"log"
)

func main() {
	// Load environment variables
	if err := config.LoadEnv(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Initialize database
	if err := config.InitDB(); err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	// Run migrations
	if err := migrations.RunMigrations(); err != nil {
		log.Fatal("Error running migrations:", err)
	}

	log.Println("Migrations completed successfully")
}
