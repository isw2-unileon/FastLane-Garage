// Package database handles all database initialization, migrations, and seeding.
package database

import (
	"fmt"
	"log/slog"

	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Init initializes a connection to the SQLite database.
// It takes the path to the database file as a parameter and returns
// a GORM database instance or an error if the connection fails.
func Init(dbPath string) (*gorm.DB, error) {
	// Open SQLite database at the specified path.
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// Log successful connection.
	slog.Info("database connected successfully")
	return db, nil
}

// Migrate runs AutoMigrate for all models to create or update database tables.
// This function should be called once during application startup.
func Migrate(db *gorm.DB) error {
	// AutoMigrate creates the "parts" table if it doesn't exist
	// and synchronizes the schema with the Part struct.
	if err := db.AutoMigrate(&models.Part{}); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// Log successful migration completion.
	slog.Info("database migrations completed")
	return nil
}

// SeedParts inserts sample part data into the database.
// If the table already contains data, this function returns early (idempotent).
// This is useful for development and testing environments.
func SeedParts(db *gorm.DB) error {
	// Count existing parts in the database to check if seeding is needed.
	var count int64
	if err := db.Model(&models.Part{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check existing parts: %w", err)
	}

	// If parts already exist, skip seeding to avoid duplication.
	if count > 0 {
		slog.Info("parts table already seeded, skipping")
		return nil
	}

	// Sample parts data representing different car zones.
	// Each part has a name, car zone, image URL placeholder, and price.
	parts := []models.Part{
		{
			Name:     "Motor V6",
			CarZone:  "motor",
			ImageURL: "",
			Price:    2500.00,
		},
		{
			Name:     "Turbo",
			CarZone:  "motor",
			ImageURL: "",
			Price:    1200.00,
		},
		{
			Name:     "Neumático Michelin",
			CarZone:  "neumaticos",
			ImageURL: "",
			Price:    150.00,
		},
		{
			Name:     "Disco de freno",
			CarZone:  "frenos",
			ImageURL: "",
			Price:    120.00,
		},
		{
			Name:     "Pastillas de freno",
			CarZone:  "frenos",
			ImageURL: "",
			Price:    80.00,
		},
		{
			Name:     "Puerta delantera izquierda",
			CarZone:  "puertas",
			ImageURL: "",
			Price:    300.00,
		},
		{
			Name:     "Faro LED",
			CarZone:  "iluminacion",
			ImageURL: "",
			Price:    250.00,
		},
		{
			Name:     "Batería 12V",
			CarZone:  "electrico",
			ImageURL: "",
			Price:    180.00,
		},
	}

	// Insert all parts into the database in batches of 100.
	// Batch insertion is more efficient than inserting one by one.
	if err := db.CreateInBatches(parts, 100).Error; err != nil {
		return fmt.Errorf("failed to seed parts: %w", err)
	}

	// Log successful seeding with the number of parts inserted.
	slog.Info("database seeded successfully", "count", len(parts))
	return nil
}
