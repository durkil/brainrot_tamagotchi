package database

import (
	"brainrot-tamagotchi/internal/models"

	"gorm.io/gorm"
)

// AutoMigrate runs all database migrations
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.NFT{},
		&models.MarketListing{},
		&models.CaseOpening{},
	)
}

