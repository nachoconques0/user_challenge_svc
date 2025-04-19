package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(opts ...Option) (*gorm.DB, error) {
	// Default and apply options
	options := defaultOptions()
	for _, o := range opts {
		o(&options)
	}

	// Setup connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		options.Host,
		options.Port,
		options.Database,
		options.User,
		options.Password,
		options.SSLMode,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: newLogger,
	})

	return db, err
}

// Check if the passed db connection is already
// running a transaction
func IsTransaction(tx *gorm.DB) bool {
	_, ok := tx.Statement.ConnPool.(sqlTx)
	return ok
}

type sqlTx interface {
	Commit() error
	Rollback() error
}
