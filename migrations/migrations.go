package migrations

import (
	"game-time-api/config"
	"game-time-api/models"
	"log"
)

// RunMigrations performs all database migrations
func RunMigrations() error {
	db := config.GetDB()

	// Auto-migrate the schema
	log.Println("Running auto-migrations...")

	// Add all models here for auto-migration
	err := db.AutoMigrate(
		&models.User{},
		// Add more models here as your application grows
	)

	if err != nil {
		return err
	}

	log.Println("Auto-migrations completed successfully")
	return nil
}
