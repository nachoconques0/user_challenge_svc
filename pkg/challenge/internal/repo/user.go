package repo

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/entity"
	"gorm.io/gorm/clause"
)

var (
	// ErrMissingDB used when DB is nil
	ErrMissingDB = errors.New("DB connection is missing")
	// ErrIDShouldNotBeEmpty used when entity ID is empty
	ErrIDShouldNotBeEmpty = errors.New("entity ID should not be empty")
	// ErrIDnotValid used when entity ID is not valid
	ErrIDnotValid = errors.New("entity ID not valid")
	// ErrRecordNotFound used when record is not found
	ErrRecordNotFound = errors.New("record not found")
	// ErrHashingPassword used when there was an error hashing the password
	ErrHashingPassword = errors.New("error hashing password")
)

// Create creates a new user in the DB
func Create(u *entity.User, tx *gorm.DB) (*entity.User, error) {
	if tx == nil {
		return nil, ErrMissingDB
	}
	u.ID = uuid.MustParse(uuid.NewString())

	err := u.HashPassword(u.Password)
	if err != nil {
		return nil, ErrHashingPassword
	}

	if res := tx.Create(&u); res.Error != nil {
		return nil, res.Error
	}
	return u, nil
}

// Find returns a list of users. It can be paginated and filtered by user country
func Find(tx *gorm.DB, country string, page, limit int) ([]entity.User, error) {
	var users []entity.User
	query := tx.Model(&entity.User{})

	if tx == nil {
		return nil, ErrMissingDB
	}

	if country != "" {
		query = query.Where("country = ?", country)
	}

	// Default limit if none provided
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	if err := query.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserForUpdate returns an user and will lock the row in order to update it
func GetUserForUpdate(id uuid.UUID, tx *gorm.DB) (*entity.User, error) {
	var u entity.User
	if tx == nil {
		return nil, ErrMissingDB
	}
	if id == uuid.Nil {
		return nil, ErrIDShouldNotBeEmpty
	}

	if res := tx.Where(entity.User{ID: id}).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&u); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
	}

	return &u, nil
}

// Update updates only the nickname of an existing user
func Update(u *entity.User, tx *gorm.DB) (*entity.User, error) {
	if tx == nil {
		return nil, ErrMissingDB
	}
	if u == nil || u.ID == uuid.Nil {
		return nil, ErrIDShouldNotBeEmpty
	}
	if strings.TrimSpace(u.Nickname) == "" {
		return nil, errors.New("nickname cannot be empty")
	}

	if err := tx.Model(&entity.User{}).
		Where("id = ?", u.ID).
		Select("nickname").
		Updates(&entity.User{
			Nickname: u.Nickname,
		}).Error; err != nil {
		return nil, err
	}

	return u, nil
}

// Delete soft deletes a user by ID
func Delete(id uuid.UUID, tx *gorm.DB) error {
	if tx == nil {
		return ErrMissingDB
	}
	if id == uuid.Nil {
		return ErrIDShouldNotBeEmpty
	}

	if err := tx.Delete(&entity.User{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
