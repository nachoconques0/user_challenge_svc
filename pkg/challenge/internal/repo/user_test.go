package repo_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/helpers"
	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/entity"
	"github.com/nachoconques0/user_challenge_svc/pkg/challenge/internal/repo"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	password = "faceit123"
	user     = entity.User{
		FirstName: "nacho",
		LastName:  "calcagno",
		Nickname:  "faceitcsgo",
		Password:  password,
		Email:     "faceitnacho@gmail.com",
		Country:   "ES",
	}
)

func TestRepository_Create(t *testing.T) {
	db, teardown, err := helpers.NewTestDB()
	if err != nil {
		assert.Nil(t, err)
	}
	defer teardown()

	t.Run("it should create an user", func(t *testing.T) {
		res, err := repo.Create(&user, db)
		assert.Nil(t, err)
		assert.Equal(t, user.ID, res.ID)
		assert.Equal(t, user.FirstName, res.FirstName)
		assert.Equal(t, user.LastName, res.LastName)
		assert.Equal(t, user.Nickname, res.Nickname)
		assert.Equal(t, user.LastName, res.LastName)
		assert.Equal(t, user.Email, res.Email)
		assert.Equal(t, user.Country, res.Country)
		assert.NotEqual(t, res.CreatedAt, time.Time{})
		assert.NotEqual(t, res.UpdatedAt, time.Time{})
		t.Run("and user password encrypted", func(t *testing.T) {
			err := bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(password))
			assert.Nil(t, err)
		})
	})

	t.Run("it should fail if the user already exists", func(t *testing.T) {
		_, err = repo.Create(&user, db)
		assert.NotNil(t, err)
	})
}

func TestRepository_Find(t *testing.T) {
	db, teardown, err := helpers.NewTestDB()
	if err != nil {
		assert.Nil(t, err)
	}
	defer teardown()
	insertTestUsers(t, db)

	t.Run("should return paginated users filtered by country", func(t *testing.T) {
		users, err := repo.Find(db, "ES", 1, 0)
		assert.NoError(t, err)
		assert.Len(t, users, 2)
		assert.Equal(t, "ES", users[0].Country)
		assert.Equal(t, "ES", users[1].Country)
	})

	t.Run("should return all users when no country is specified", func(t *testing.T) {
		users, err := repo.Find(db, "", 1, 0)
		assert.NoError(t, err)
		assert.Len(t, users, 3)
	})

	t.Run("should return paginated results correctly", func(t *testing.T) {
		usersPage1, err := repo.Find(db, "", 1, 2)
		assert.NoError(t, err)
		assert.Len(t, usersPage1, 2)

		usersPage2, err := repo.Find(db, "", 2, 2)
		assert.NoError(t, err)
		assert.Len(t, usersPage2, 1)
	})
}

func TestRepository_GetUserForUpdate(t *testing.T) {
	db, teardown, err := helpers.NewTestDB()
	if err != nil {
		assert.Nil(t, err)
	}
	defer teardown()

	t.Run("should return error if tx is nil", func(t *testing.T) {
		user, err := repo.GetUserForUpdate(uuid.MustParse(uuid.NewString()), nil)
		assert.Nil(t, user)
		assert.Equal(t, repo.ErrMissingDB, err)
	})

	t.Run("should return error if ID is nil", func(t *testing.T) {
		user, err := repo.GetUserForUpdate(uuid.Nil, db)
		assert.Nil(t, user)
		assert.Equal(t, repo.ErrIDShouldNotBeEmpty, err)
	})

	t.Run("should return error if user not found", func(t *testing.T) {
		randomID := uuid.MustParse(uuid.NewString())
		user, err := repo.GetUserForUpdate(randomID, db)
		assert.Nil(t, user)
		assert.Equal(t, repo.ErrRecordNotFound, err)
	})

	t.Run("should return user", func(t *testing.T) {
		user := &entity.User{
			FirstName: "Test",
			LastName:  "User",
			Nickname:  "TU123",
			Password:  "somepassword",
			Email:     "test@user.com",
			Country:   "AR",
		}
		created, err := repo.Create(user, db)
		assert.NoError(t, err)

		result, err := repo.GetUserForUpdate(created.ID, db)
		assert.NoError(t, err)
		assert.Equal(t, created.ID, result.ID)
	})
}

func TestRepository_Update(t *testing.T) {
	db, teardown, err := helpers.NewTestDB()
	if err != nil {
		assert.Nil(t, err)
	}
	defer teardown()

	t.Run("should update nickname", func(t *testing.T) {
		user := &entity.User{
			FirstName: "nacho",
			LastName:  "calcagno",
			Nickname:  "nachin",
			Password:  "nachoc!",
			Email:     "nachoc@gmail.com",
			Country:   "VE",
		}
		created, err := repo.Create(user, db)
		assert.NoError(t, err)

		// update nickname
		created.Nickname = "letsplaycsgo"
		updated, err := repo.Update(created, db)
		assert.NoError(t, err)
		assert.Equal(t, "letsplaycsgo", updated.Nickname)

		fromDB, err := repo.Find(db, "", 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, "letsplaycsgo", fromDB[0].Nickname)
	})

	t.Run("should return error if tx is nil", func(t *testing.T) {
		user := &entity.User{
			ID:       uuid.New(),
			Nickname: "nacho",
		}
		updated, err := repo.Update(user, nil)
		assert.Nil(t, updated)
		assert.Equal(t, repo.ErrMissingDB, err)
	})

	t.Run("should return error if user is nil", func(t *testing.T) {
		updated, err := repo.Update(nil, db)
		assert.Nil(t, updated)
		assert.Equal(t, repo.ErrIDShouldNotBeEmpty, err)
	})

	t.Run("should return error if user ID is nil", func(t *testing.T) {
		user := &entity.User{
			ID:       uuid.Nil,
			Nickname: "nacho",
		}
		updated, err := repo.Update(user, db)
		assert.Nil(t, updated)
		assert.Equal(t, repo.ErrIDShouldNotBeEmpty, err)
	})

	t.Run("should return error if nickname is empty", func(t *testing.T) {
		user := &entity.User{
			ID:       uuid.New(),
			Nickname: "   ",
		}
		updated, err := repo.Update(user, db)
		assert.Nil(t, updated)
		assert.Error(t, err)
		assert.Equal(t, "nickname cannot be empty", err.Error())
	})
}

func TestRepository_Delete(t *testing.T) {
	db, teardown, err := helpers.NewTestDB()
	if err != nil {
		assert.Nil(t, err)
	}
	defer teardown()

	t.Run("should delete an user", func(t *testing.T) {
		user := &entity.User{
			FirstName: "Soft",
			LastName:  "Delete",
			Nickname:  "softie",
			Password:  "securepass",
			Email:     "soft@delete.com",
			Country:   "CL",
		}
		created, err := repo.Create(user, db)
		assert.NoError(t, err)

		err = repo.Delete(created.ID, db)
		assert.NoError(t, err)

		users, err := repo.Find(db, "", 1, 10)
		assert.NoError(t, err)

		found := false
		for _, u := range users {
			if u.ID == created.ID {
				found = true
			}
		}
		assert.False(t, found, "deleted user should not be returned in Find")
	})

	t.Run("should return error if tx is nil", func(t *testing.T) {
		err := repo.Delete(uuid.New(), nil)
		assert.Equal(t, repo.ErrMissingDB, err)
	})

	t.Run("should return error if ID is nil", func(t *testing.T) {
		err := repo.Delete(uuid.Nil, db)
		assert.Equal(t, repo.ErrIDShouldNotBeEmpty, err)
	})
}

func insertTestUsers(t *testing.T, db *gorm.DB) {
	users := []entity.User{
		{
			FirstName: "Nacho1",
			LastName:  "faceit1",
			Nickname:  "nacho1",
			Password:  "hashedpass1",
			Email:     "Nacho1@gmail.com",
			Country:   "ES",
		},
		{
			FirstName: "Juan",
			LastName:  "Calcagno",
			Nickname:  "Juan",
			Password:  "juancalcagnp",
			Email:     "jcalcagno@nacho.com",
			Country:   "ES",
		},
		{
			FirstName: "Nacho",
			LastName:  "Calcagno",
			Nickname:  "NachoCalcagno",
			Password:  "nachonacho",
			Email:     "nachoc@gmail.com",
			Country:   "IT",
		},
	}

	for _, u := range users {
		_, err := repo.Create(&u, db)
		assert.NoError(t, err)
	}
}
