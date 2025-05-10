package testing

import (
	"context"
	"server/internal/auth"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUserRepository(t *testing.T) {

	repo := NewMockAuthRepository()

	t.Run("create and get user", func(t *testing.T) {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
		assert.NoError(t, err)

		user := &auth.User{
			ID:           uuid.New().String(),
			Name:         "Test User",
			Email:        "test@example.com",
			PasswordHash: string(hashedPassword),
		}

		// Create user
		createdUser, err := repo.UpsertUser(context.Background(), user)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, createdUser.ID)

		// Get user by ID
		foundUser, err := repo.GetUserByID(context.Background(), user.ID)
		assert.NoError(t, err)
		assert.Equal(t, user.Name, foundUser.Name)
		assert.Equal(t, user.Email, foundUser.Email)
	})

	t.Run("get user by email", func(t *testing.T) {
		hashedPassword2, err := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
		assert.NoError(t, err)

		user := &auth.User{
			ID:           uuid.New().String(),
			Name:         "Email User",
			Email:        "email@example.com",
			PasswordHash: string(hashedPassword2),
		}

		// Create user
		_, err = repo.UpsertUser(context.Background(), user)
		assert.NoError(t, err)

		// Get user by email
		foundUser, err := repo.GetUserByEmail(context.Background(), user.Email)
		assert.NoError(t, err)
		assert.Equal(t, user.Name, foundUser.Name)
		assert.Equal(t, user.Email, foundUser.Email)
	})

	t.Run("get user by google id", func(t *testing.T) {
		user := &auth.User{
			ID:       uuid.New().String(),
			Name:     "Google User",
			Email:    "google@example.com",
			GoogleID: "google123",
		}

		// Create user
		_, err := repo.UpsertUser(context.Background(), user)
		assert.NoError(t, err)

		// Get user by Google ID
		foundUser, err := repo.GetUserByGoogleID(context.Background(), user.GoogleID)
		assert.NoError(t, err)
		assert.Equal(t, user.Name, foundUser.Name)
		assert.Equal(t, user.Email, foundUser.Email)
		assert.Equal(t, user.GoogleID, foundUser.GoogleID)
	})

	t.Run("get non-existent user", func(t *testing.T) {
		_, err := repo.GetUserByID(context.Background(), uuid.New().String())
		assert.Nil(t, err)
	})
}
