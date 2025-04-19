package services_test

import (
	"context"
	"fmt"
	"math"
	"strings"
	"testing"

	dto "example.com/api/internal/contracts"
	"example.com/api/internal/repository"
	"example.com/api/internal/services"
	"example.com/api/pkg/logging"
	mocks "example.com/api/tests/unit/mocks/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserService_GetAll(t *testing.T) {
	require.NotNil(t, testDB, "Test Database connection (testDB) is not initialized")
	ctx := context.Background()

	seedUsers := func(t *testing.T, count int) []string {
		emails := make([]string, count)
		for i := range count {
			email := fmt.Sprintf("user%d@example.com", i+1)
			seedUser(t, email)
			emails[i] = email
		}
		return emails
	}

	t.Run("Success - Returns Users with Pagination", func(t *testing.T) {
		TruncateTables(t, testDB, testTableNames)
		repoManager := repository.NewRepositoryManager(testDB)

		// mockery-generated mocks
		mockLogger := mocks.NewMockLogger(t)
		mockHashService := mocks.NewMockHashService(t)

		userService := services.NewUserService(repoManager, mockLogger, mockHashService)

		// Seed 3 users
		emails := seedUsers(t, 3)

		// Fetch first 2 users
		arg := dto.ListUsersParams{Limit: 2, Offset: 0}
		users, err := userService.GetAll(ctx, arg)

		require.NoError(t, err)
		require.Len(t, users, 2)
		assert.Equal(t, emails[0], users[0].Email)
		assert.Equal(t, emails[1], users[1].Email)

		// Fetch next user with offset
		arg = dto.ListUsersParams{Limit: 2, Offset: 2}
		users, err = userService.GetAll(ctx, arg)

		require.NoError(t, err)
		require.Len(t, users, 1)
		assert.Equal(t, emails[2], users[0].Email)
	})

	t.Run("Success - Limit 0 Returns Empty List", func(t *testing.T) {
		TruncateTables(t, testDB, testTableNames)
		repoManager := repository.NewRepositoryManager(testDB)

		mockLogger := mocks.NewMockLogger(t)
		mockHashService := mocks.NewMockHashService(t)

		userService := services.NewUserService(repoManager, mockLogger, mockHashService)

		seedUsers(t, 2) // Seed users, but limit 0 should return none

		arg := dto.ListUsersParams{Limit: 0, Offset: 0}
		users, err := userService.GetAll(ctx, arg)

		require.NoError(t, err)
		assert.Empty(t, users)
	})

	t.Run("Success - Offset Beyond Data Returns Empty List", func(t *testing.T) {
		TruncateTables(t, testDB, testTableNames)
		repoManager := repository.NewRepositoryManager(testDB)

		mockLogger := mocks.NewMockLogger(t)
		mockHashService := mocks.NewMockHashService(t)

		userService := services.NewUserService(repoManager, mockLogger, mockHashService)

		seedUsers(t, 2) // Only 2 users

		arg := dto.ListUsersParams{Limit: 10, Offset: 3} // Offset beyond data
		users, err := userService.GetAll(ctx, arg)

		require.NoError(t, err)
		assert.Empty(t, users)
	})

	t.Run("Error - Invalid Negative Limit", func(t *testing.T) {
		TruncateTables(t, testDB, testTableNames)
		repoManager := repository.NewRepositoryManager(testDB)

		mockLogger := mocks.NewMockLogger(t)
		mockHashService := mocks.NewMockHashService(t)

		userService := services.NewUserService(repoManager, mockLogger, mockHashService)

		// Expect error log with PostgreSQL's "LIMIT must not be negative" message
		mockLogger.EXPECT().Error(logging.Postgres, logging.Select, "Failed to fetch all users",
			mock.MatchedBy(func(extra map[logging.ExtraKey]any) bool {
				errMsg, ok := extra[logging.ErrorMessage].(string)
				return ok && strings.Contains(errMsg, "LIMIT must not be negative")
			})).Once()

		arg := dto.ListUsersParams{Limit: -1, Offset: 0}
		users, err := userService.GetAll(ctx, arg)

		require.Error(t, err)
		assert.Equal(t, "failed to fetch users", err.Error())
		assert.Nil(t, users)
	})

	t.Run("Success - Extremely Large Limit", func(t *testing.T) {
		TruncateTables(t, testDB, testTableNames)
		repoManager := repository.NewRepositoryManager(testDB)

		mockLogger := mocks.NewMockLogger(t)
		mockHashService := mocks.NewMockHashService(t)

		userService := services.NewUserService(repoManager, mockLogger, mockHashService)

		emails := seedUsers(t, 3)

		arg := dto.ListUsersParams{Limit: math.MaxInt32, Offset: 0}
		users, err := userService.GetAll(ctx, arg)

		require.NoError(t, err)
		assert.Len(t, users, 3)
		assert.Equal(t, emails, []string{users[0].Email, users[1].Email, users[2].Email})
	})
}
