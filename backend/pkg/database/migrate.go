package database

import (
	"fmt"

	"github.com/vibecode/ecommerce/backend/internal/domain"
	"gorm.io/gorm"
)

// AutoMigrate applies schema for all registered models.
// For production, prefer versioned SQL migrations in /migrations.
func AutoMigrate(db *gorm.DB) error {
	models := []any{
		&domain.User{},
		&domain.Product{},
	}
	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}
	return nil
}
