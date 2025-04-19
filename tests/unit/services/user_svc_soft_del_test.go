package services_test

import (
	"context"
	"database/sql"
	"math"
	"testing"

	"example.com/api/internal/repository"
	"example.com/api/internal/services"
	mocks "example.com/api/tests/unit/mocks/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserService_SoftDelete(t *testing.T) {
	require.NotNil(t, testDB, "Test Database connection (testDB) is not initialized")
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Arrange
		TruncateTables(t, testDB, testTableNames) // Clean DB

		repoManager := repository.NewRepositoryManager(testDB)
		mockHashService := mocks.NewMockHashService(t)
		mockLogger := mocks.NewMockLogger(t)

		// Seed user to delete
		userID := seedUser(t, "softdelete@example.com", "softdelete_user", "Soft Delete User", "hash")
		require.NotZero(t, userID, "Failed to seed user")

		userService := services.NewUserService(repoManager, mockLogger, mockHashService)

		// Act
		err := userService.SoftDelete(ctx, userID)

		// Assert
		require.NoError(t, err)

		// Verify user is soft-deleted (check deleted_at is set)
		var deletedAt sql.NullTime
		query := "SELECT deleted_at FROM users WHERE id = $1"
		err = testDB.QueryRowContext(ctx, query, userID).Scan(&deletedAt)
		require.NoError(t, err)
		assert.True(t, deletedAt.Valid, "Expected deleted_at to be set")
	})

	t.Run("User Not Found", func(t *testing.T) {
		TruncateTables(t, testDB, testTableNames)
		repoManager := repository.NewRepositoryManager(testDB)
		mockHashService := mocks.NewMockHashService(t)
		mockLogger := mocks.NewMockLogger(t)

		nonExistentID := int32(9999)
		userService := services.NewUserService(repoManager, mockLogger, mockHashService)

		// Act
		err := userService.SoftDelete(ctx, nonExistentID)

		// Assert
		require.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
		mockLogger.AssertNotCalled(t, "Error")
	})

	t.Run("Negative ID", func(t *testing.T) {
		// Arrange
		TruncateTables(t, testDB, testTableNames)
		repoManager := repository.NewRepositoryManager(testDB)
		mockLogger := mocks.NewMockLogger(t)
		userService := services.NewUserService(repoManager, mockLogger, nil)

		// Act
		err := userService.SoftDelete(ctx, -1) // Negative ID

		// Assert
		require.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
		mockLogger.AssertNotCalled(t, "Error") // Ensure no error logging
	})

	t.Run("Zero ID", func(t *testing.T) {
		// Arrange
		TruncateTables(t, testDB, testTableNames)
		repoManager := repository.NewRepositoryManager(testDB)
		mockLogger := mocks.NewMockLogger(t)
		userService := services.NewUserService(repoManager, mockLogger, nil)

		// Act
		err := userService.SoftDelete(ctx, 0) // Zero ID

		// Assert
		require.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
		mockLogger.AssertNotCalled(t, "Error")
	})

	t.Run("Extremely Large ID", func(t *testing.T) {
		// Arrange
		TruncateTables(t, testDB, testTableNames)
		repoManager := repository.NewRepositoryManager(testDB)
		mockLogger := mocks.NewMockLogger(t)
		userService := services.NewUserService(repoManager, mockLogger, nil)

		// Act
		err := userService.SoftDelete(ctx, math.MaxInt32) // Max int32 value

		// Assert
		require.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
		mockLogger.AssertNotCalled(t, "Error")
	})
}
