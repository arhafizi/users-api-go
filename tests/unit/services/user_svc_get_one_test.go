package services_test

import (
	"context"
	"testing"

	"example.com/api/internal/repository"
	"example.com/api/internal/services"
	"example.com/api/pkg/logging"
	mocks "example.com/api/tests/unit/mocks/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserService_GetByID(t *testing.T) {
	require.NotNil(t, testDB, "Test Database connection (testDB) is not initialized")
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Arrange
		TruncateTables(t, testDB, testTableNames) // Clean DB

		repoManager := repository.NewRepositoryManager(testDB)
		mockHashService := mocks.NewMockHashService(t)
		mockLogger := mocks.NewMockLogger(t)

		// Seed user using testDB
		email := "getbyid_notx@example.com"
		seededUserID := seedUser(t, email, "getbyid_user_notx", "Get By ID NoTX", "hash")
		require.NotZero(t, seededUserID, "Failed to seed user")

		userService := services.NewUserService(repoManager, mockLogger, mockHashService)

		// Act
		user, err := userService.GetByID(ctx, seededUserID)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, seededUserID, user.ID)
		assert.Equal(t, "getbyid_user_notx", user.Username)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Arrange
		TruncateTables(t, testDB, testTableNames) // Clean DB

		repoManager := repository.NewRepositoryManager(testDB)
		mockHashService := mocks.NewMockHashService(t)
		mockLogger := mocks.NewMockLogger(t)

		nonExistentID := int32(99999)
		mockLogger.EXPECT().Error(logging.Postgres, logging.Select, "Failed to fetch user by ID",
			mock.MatchedBy(func(extra map[logging.ExtraKey]any) bool {
				id, ok := extra["userID"].(int32)
				return ok && id == nonExistentID
			})).Once()

		userService := services.NewUserService(repoManager, mockLogger, mockHashService)

		// Act
		user, err := userService.GetByID(ctx, nonExistentID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "user not found", err.Error())
	})
}

func TestUserService_GetByUsername(t *testing.T) {
	require.NotNil(t, testDB, "Test Database connection (testDB) is not initialized")
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Arrange
		TruncateTables(t, testDB, testTableNames) // Clean DB

		repoManager := repository.NewRepositoryManager(testDB)
		mockHashService := mocks.NewMockHashService(t)
		mockLogger := mocks.NewMockLogger(t)

		// Seed user using testDB
		email := "getbyusername@example.com"
		username := "getbyusername_user"
		seededUserID := seedUser(t, email, username, "Get By Username", "hash")
		require.NotZero(t, seededUserID, "Failed to seed user")

		userService := services.NewUserService(repoManager, mockLogger, mockHashService)

		// Act
		user, err := userService.GetByUsername(ctx, username)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, seededUserID, user.ID)
		assert.Equal(t, username, user.Username)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Arrange
		TruncateTables(t, testDB, testTableNames) // Clean DB

		repoManager := repository.NewRepositoryManager(testDB)
		mockHashService := mocks.NewMockHashService(t)
		mockLogger := mocks.NewMockLogger(t)

		nonExistentUsername := "nonexistent_username"
		mockLogger.EXPECT().Error(logging.Postgres, logging.Select, "Failed to fetch user by username",
			mock.MatchedBy(func(extra map[logging.ExtraKey]any) bool {
				username, ok := extra["username"].(string)
				return ok && username == nonExistentUsername
			})).Once()

		userService := services.NewUserService(repoManager, mockLogger, mockHashService)

		// Act
		user, err := userService.GetByUsername(ctx, nonExistentUsername)

		// Assert
		require.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "user not found", err.Error())
	})
}

func TestUserService_GetByEmail(t *testing.T) {
	require.NotNil(t, testDB, "Test Database connection (testDB) is not initialized")
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Arrange
		TruncateTables(t, testDB, testTableNames) // Clean DB

		repoManager := repository.NewRepositoryManager(testDB)
		mockHashService := mocks.NewMockHashService(t)
		mockLogger := mocks.NewMockLogger(t)

		// Seed user using testDB
		email := "getbyemail@example.com"
		seededUserID := seedUser(t, email, "getbyemail_user", "Get By Email", "hash")
		require.NotZero(t, seededUserID, "Failed to seed user")

		userService := services.NewUserService(repoManager, mockLogger, mockHashService)

		// Act
		user, err := userService.GetByEmail(ctx, email)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, seededUserID, user.ID)
		assert.Equal(t, email, user.Email)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Arrange
		TruncateTables(t, testDB, testTableNames) // Clean DB

		repoManager := repository.NewRepositoryManager(testDB)
		mockHashService := mocks.NewMockHashService(t)
		mockLogger := mocks.NewMockLogger(t)

		nonExistentEmail := "nonexistent@example.com"
		mockLogger.EXPECT().Error(logging.Postgres, logging.Select, "Failed to fetch user by email",
			mock.MatchedBy(func(extra map[logging.ExtraKey]any) bool {
				email, ok := extra["email"].(string)
				return ok && email == nonExistentEmail
			})).Once()

		userService := services.NewUserService(repoManager, mockLogger, mockHashService)

		// Act
		user, err := userService.GetByEmail(ctx, nonExistentEmail)

		// Assert
		require.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "user not found", err.Error())
	})
}
