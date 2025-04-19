package services_test

import (
	"context"
	"errors"
	"testing"

	dto "example.com/api/internal/contracts"
	contracts "example.com/api/internal/contracts/errors"
	"example.com/api/internal/repository"
	dbCtx "example.com/api/internal/repository/db"
	"example.com/api/internal/services"
	"example.com/api/pkg/logging"
	mocks "example.com/api/tests/unit/mocks/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserService_UpdatePartial(t *testing.T) {
	require.NotNil(t, testDB, "Test Database connection (testDB) is not initialized")
	ctx := context.Background()

	seedUserWithId := func(t *testing.T, id int32, username, email, fullName, passwordHash string) {
		query := `
            INSERT INTO users (id, username, email, full_name, password_hash) 
            VALUES ($1, $2, $3, $4, $5)
            ON CONFLICT (id) DO NOTHING`
		_, err := testDB.ExecContext(ctx, query, id, username, email, fullName, passwordHash)
		require.NoError(t, err)
	}

	t.Run("Success - Full Update", func(t *testing.T) {
		TruncateTables(t, testDB, testTableNames)
		mockHash := mocks.NewMockHashService(t)
		mockLogger := mocks.NewMockLogger(t)
		userService := services.NewUserService(repository.NewRepositoryManager(testDB), mockLogger, mockHash)

		// Seed user
		seedUserWithId(t, 1, "olduser", "old@example.com", "Old Name", "oldhash")

		// Setup mock for password hashing
		newPass := "newpass123"
		hashedPass := "hashednewpass123"
		mockHash.EXPECT().Hash(newPass).Return(hashedPass, nil).Once()

		arg := dto.UpdateUserPartialReq{
			ID:       1,
			Username: ptr("newuser"),
			Email:    ptr("new@example.com"),
			FullName: ptr("New Name"),
			Password: &newPass,
		}

		user, err := userService.UpdatePartial(ctx, arg)
		require.NoError(t, err)
		assert.Equal(t, "newuser", user.Username)
		assert.Equal(t, "new@example.com", user.Email)
		assert.Equal(t, "New Name", user.FullName)
		assert.Equal(t, hashedPass, user.PasswordHash)

		// Verify DB state
		var dbUser dbCtx.User
		err = testDB.QueryRowContext(ctx, "SELECT * FROM users WHERE id = $1", 1).Scan(
			&dbUser.ID,
			&dbUser.Username,
			&dbUser.Email,
			&dbUser.FullName,
			&dbUser.PasswordHash,
			&dbUser.CreatedAt,
			&dbUser.UpdatedAt,
			&dbUser.DeletedAt, // Fixed: Scan into sql.NullTime directly
		)
		require.NoError(t, err)
		assert.Equal(t, "newuser", dbUser.Username)
		assert.Equal(t, "new@example.com", dbUser.Email)
		assert.Equal(t, "New Name", dbUser.FullName)
		assert.Equal(t, hashedPass, dbUser.PasswordHash)
	})

	t.Run("Username Conflict", func(t *testing.T) {
		TruncateTables(t, testDB, testTableNames)
		mockLogger := mocks.NewMockLogger(t)
		userService := services.NewUserService(repository.NewRepositoryManager(testDB), mockLogger, nil)

		// Seed conflicting users
		seedUserWithId(t, 1, "existing", "user1@example.com", "User1", "hash1")
		seedUserWithId(t, 2, "user2", "user2@example.com", "User2", "hash2")

		arg := dto.UpdateUserPartialReq{
			ID:       2,
			Username: ptr("existing"), // Conflicts with user1
		}

		mockLogger.EXPECT().Warn(logging.Validation, logging.Update, "Username already exists",
			mock.MatchedBy(func(extra map[logging.ExtraKey]any) bool {
				username, ok := extra[logging.RequestBody].(*string)
				return ok && username != nil && *username == "existing"
			})).Once()

		user, err := userService.UpdatePartial(ctx, arg)
		require.Error(t, err)
		assert.Nil(t, user)
		var usernameErr *contracts.UsernameExistsError
		require.True(t, errors.As(err, &usernameErr))
		assert.Equal(t, "existing", usernameErr.Username)
	})

	t.Run("Email Conflict", func(t *testing.T) {
		TruncateTables(t, testDB, testTableNames)
		mockLogger := mocks.NewMockLogger(t)
		userService := services.NewUserService(repository.NewRepositoryManager(testDB), mockLogger, nil)

		// Seed conflicting users
		seedUserWithId(t, 1, "user1", "existing@example.com", "User1", "hash1")
		seedUserWithId(t, 2, "user2", "user2@example.com", "User2", "hash2")

		arg := dto.UpdateUserPartialReq{
			ID:    2,
			Email: ptr("existing@example.com"), // Conflicts with user1
		}

		mockLogger.EXPECT().Warn(logging.Validation, logging.Update, "Email already exists",
			mock.MatchedBy(func(extra map[logging.ExtraKey]any) bool {
				email, ok := extra[logging.RequestBody].(*string)
				return ok && email != nil && *email == "existing@example.com"
			})).Once()

		user, err := userService.UpdatePartial(ctx, arg)
		require.Error(t, err)
		assert.Nil(t, user)
		var emailErr *contracts.EmailExistsError
		require.True(t, errors.As(err, &emailErr))
		assert.Equal(t, "existing@example.com", emailErr.Email)
	})

	t.Run("User Not Found", func(t *testing.T) {
		TruncateTables(t, testDB, testTableNames)
		mockLogger := mocks.NewMockLogger(t)
		userService := services.NewUserService(repository.NewRepositoryManager(testDB), mockLogger, nil)

		arg := dto.UpdateUserPartialReq{ID: 999, Username: ptr("newuser")}

		mockLogger.EXPECT().Error(logging.Postgres, logging.Update, "User not found",
			mock.MatchedBy(func(extra map[logging.ExtraKey]any) bool {
				return extra["userID"].(int32) == 999
			})).Once()

		user, err := userService.UpdatePartial(ctx, arg)
		require.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "user not found", err.Error())
	})

	t.Run("Partial Update (Only FullName)", func(t *testing.T) {
		TruncateTables(t, testDB, testTableNames)
		userService := services.NewUserService(repository.NewRepositoryManager(testDB), mocks.NewMockLogger(t), nil)

		seedUserWithId(t, 1, "user", "user@example.com", "Old Name", "hash")

		arg := dto.UpdateUserPartialReq{
			ID:       1,
			FullName: ptr("New Name"),
		}

		user, err := userService.UpdatePartial(ctx, arg)
		require.NoError(t, err)
		assert.Equal(t, "New Name", user.FullName)
		assert.Equal(t, "user", user.Username)          // Unchanged
		assert.Equal(t, "user@example.com", user.Email) // Unchanged
	})

	t.Run("Empty Password", func(t *testing.T) {
		TruncateTables(t, testDB, testTableNames)
		mockHash := mocks.NewMockHashService(t)
		userService := services.NewUserService(repository.NewRepositoryManager(testDB), mocks.NewMockLogger(t), mockHash)

		seedUserWithId(t, 1, "user", "user@example.com", "User", "oldhash")

		emptyPass := ""
		arg := dto.UpdateUserPartialReq{
			ID:       1,
			Password: &emptyPass,
		}

		user, err := userService.UpdatePartial(ctx, arg)
		require.NoError(t, err)
		assert.Equal(t, "", user.PasswordHash) // Empty string passed directly

		var dbPass string
		err = testDB.QueryRowContext(ctx, "SELECT password_hash FROM users WHERE id = $1", 1).Scan(&dbPass)
		require.NoError(t, err)
		assert.Equal(t, "", dbPass)

		mockHash.AssertNotCalled(t, "Hash")
	})

	t.Run("No Changes (All Fields Nil)", func(t *testing.T) {
		TruncateTables(t, testDB, testTableNames)
		userService := services.NewUserService(repository.NewRepositoryManager(testDB), mocks.NewMockLogger(t), nil)

		seedUserWithId(t, 1, "user", "user@example.com", "User", "hash")

		arg := dto.UpdateUserPartialReq{ID: 1}

		user, err := userService.UpdatePartial(ctx, arg)
		require.NoError(t, err)
		assert.Equal(t, "user", user.Username)
		assert.Equal(t, "user@example.com", user.Email)
		assert.Equal(t, "User", user.FullName)
		assert.Equal(t, "hash", user.PasswordHash)
	})
}

// Helper to create string pointers
func ptr(s string) *string {
	return &s
}
