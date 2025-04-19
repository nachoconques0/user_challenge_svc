package entity

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmptyField     = errors.New("field cannot be empty")
	ErrInvalidEmail   = errors.New("invalid email format")
	ErrWeakPassword   = errors.New("password must be at least 8 characters")
	ErrInvalidCountry = errors.New("country must be specified")
)

const (
	// TableName define user table name for user entity
	TableName = "challenge.user"
)

// User entity representation in DB
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	FirstName string    `gorm:"not null"`
	LastName  string    `gorm:"not null"`
	Nickname  string    `gorm:"not null;unique"`
	Password  string    `gorm:"not null"`
	Email     string    `gorm:"not null;unique"`
	Country   string    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

// TableName returns table name
func (User) TableName() string {
	return TableName
}

func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// Valid checks that the entity meets the criteria for being persisted
func (u *User) Valid() error {
	if strings.TrimSpace(u.FirstName) == "" {
		return fmt.Errorf("first name: %w", ErrEmptyField)
	}
	if strings.TrimSpace(u.LastName) == "" {
		return fmt.Errorf("last name: %w", ErrEmptyField)
	}
	if strings.TrimSpace(u.Nickname) == "" {
		return fmt.Errorf("nickname: %w", ErrEmptyField)
	}
	if strings.TrimSpace(u.Password) == "" || len(u.Password) < 8 {
		return fmt.Errorf("password: %w", ErrWeakPassword)
	}
	if strings.TrimSpace(u.Email) == "" {
		return fmt.Errorf("email: %w", ErrEmptyField)
	}
	if _, err := mail.ParseAddress(u.Email); err != nil {
		return fmt.Errorf("email: %w", ErrInvalidEmail)
	}
	if strings.TrimSpace(u.Country) == "" {
		return fmt.Errorf("country: %w", ErrInvalidCountry)
	}
	return nil
}
