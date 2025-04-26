package main

import (
	"game-time-api/config"
	"game-time-api/migrations"
	"game-time-api/routes"
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

	// Run database migrations
	if err := migrations.RunMigrations(); err != nil {
		log.Fatal("Error running migrations:", err)
	}

	// Setup router
	router := routes.SetupRouter()

	// Start server
	port := config.GetEnv("SERVER_PORT", "8080")
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
