package helpers

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Teardown is used to close db Connection and cleanup
type Teardown func()

// NewTestDB creates a test transaction and teardown logic for cleanup
func NewTestDB() (*gorm.DB, Teardown, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		"127.0.0.1",
		"5435",
		"user_challenge_svc",
		"user_challenge_svc",
		"user_challenge_svc",
		"prefer",
	)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	tx := db.Begin()
	if tx.Error != nil {
		return nil, nil, fmt.Errorf("error starting test db transaction: %w", tx.Error)
	}

	teardown := func() {
		_ = tx.Rollback()
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	}

	return tx, teardown, nil
}
