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

func TestUserService_Create(t *testing.T) {
	require.NotNil(t, testDB, "Test Database connection (testDB) is not initialized")

	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Arrange
		TruncateTables(t, testDB, testTableNames) // Clean DB before test

		// global testDB connection
		repoManager := repository.NewRepositoryManager(testDB)
		mockHashService := mocks.NewMockHashService(t)
		mockLogger := mocks.NewMockLogger(t)

		req := dto.CreateUserReq{
			Username: "create_success_notx",
			Email:    "create_success_notx@example.com",
			FullName: "Create Success User NoTX",
			Password: "password123",
		}
		hashedPassword := "hashed_password123"
		mockHashService.EXPECT().Hash(req.Password).Return(hashedPassword, nil).Once()
		mockLogger.EXPECT().Error(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()
		mockLogger.EXPECT().Warn(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()

		userService := services.NewUserService(repoManager, mockLogger, mockHashService)

		// Act
		userResponse, err := userService.Create(ctx, req)

		// Assert (Service Response)
		require.NoError(t, err)
		require.NotNil(t, userResponse)
		assert.Equal(t, req.Username, userResponse.Username)
		assert.NotZero(t, userResponse.ID)

		// Assert (Database State - using testDB)
		var createdUser dbCtx.User
		query := "SELECT id, username, email, full_name, password_hash, created_at FROM users WHERE id = $1"
		row := testDB.QueryRowContext(ctx, query, userResponse.ID)
		err = row.Scan(
			&createdUser.ID, &createdUser.Username, &createdUser.Email,
			&createdUser.FullName, &createdUser.PasswordHash, &createdUser.CreatedAt,
		)
		require.NoError(t, err, "Failed to query created user from test DB")
		assert.Equal(t, userResponse.ID, createdUser.ID)
		assert.Equal(t, req.Username, createdUser.Username)
		assert.Equal(t, hashedPassword, createdUser.PasswordHash)

	})

	t.Run("Username Exists", func(t *testing.T) {
		// Arrange
		TruncateTables(t, testDB, testTableNames) // Clean DB

		repoManager := repository.NewRepositoryManager(testDB)
		mockHashService := mocks.NewMockHashService(t)
		mockLogger := mocks.NewMockLogger(t)

		// Seed existing user
		existingUsername := "existing_user_notx"
		email := "unique_email_notx@example.com"
		userID := seedUser(t, email, existingUsername, "unique_email_notx@example.com", "Existing NoTX", "hash")
		require.NotZero(t, userID, "Failed to seed existing user")

		req := dto.CreateUserReq{
			Username: existingUsername, // Conflicting username
			Email:    "new_email_notx@example.com",
			FullName: "New User NoTX",
			Password: "password123",
		}
		hashedPassword := "hashed_password_conflict_notx"
		mockHashService.EXPECT().Hash(req.Password).Return(hashedPassword, nil).Once()
		mockLogger.EXPECT().Warn(logging.Validation, logging.FailedToCreateUser, "Username already exists", mock.AnythingOfType("map[logging.ExtraKey]interface {}")).Once()
		mockLogger.EXPECT().Error(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()

		userService := services.NewUserService(repoManager, mockLogger, mockHashService)

		// Act
		userResponse, err := userService.Create(ctx, req)

		// Assert
		require.Error(t, err)
		assert.Nil(t, userResponse)
		var usernameExistsErr *contracts.UsernameExistsError
		require.ErrorAs(t, err, &usernameExistsErr, "Expected UsernameExistsError")
		assert.Equal(t, existingUsername, usernameExistsErr.Username)

	})

	t.Run("Email Exists", func(t *testing.T) {
		TruncateTables(t, testDB, testTableNames)
		repoManager := repository.NewRepositoryManager(testDB)
		mockHashService := mocks.NewMockHashService(t)
		mockLogger := mocks.NewMockLogger(t)

		// Seed
		existingEmail := "existing_email_notx@example.com"
		userID := seedUser(t, existingEmail, "unique_username_notx", "Existing NoTX", "hash")
		require.NotZero(t, userID, "Failed to seed existing user")

		req := dto.CreateUserReq{
			Username: "new_user_notx",
			Email:    existingEmail, // Conflicting email
			FullName: "New User NoTX",
			Password: "password123",
		}
		hashedPassword := "hashed_password_conflict_notx"
		mockHashService.EXPECT().Hash(req.Password).Return(hashedPassword, nil).Once()
		mockLogger.EXPECT().Warn(logging.Validation, logging.FailedToCreateUser, "Email already exists", mock.AnythingOfType("map[logging.ExtraKey]interface {}")).Once()

		userService := services.NewUserService(repoManager, mockLogger, mockHashService)

		userResponse, err := userService.Create(ctx, req)

		require.Error(t, err)
		assert.Nil(t, userResponse)
		var emailExistsErr *contracts.EmailExistsError
		require.ErrorAs(t, err, &emailExistsErr, "Expected EmailExistsError")
		assert.Equal(t, existingEmail, emailExistsErr.Email)
	})

	t.Run("Hashing Failed", func(t *testing.T) {
		TruncateTables(t, testDB, testTableNames)
		repoManager := repository.NewRepositoryManager(testDB)
		mockHashService := mocks.NewMockHashService(t)
		mockLogger := mocks.NewMockLogger(t)

		req := dto.CreateUserReq{
			Username: "hash_fail_user",
			Email:    "hash_fail_user@example.com",
			FullName: "Hash Fail User",
			Password: "password123",
		}
		hashingErr := errors.New("hashing failed")
		mockHashService.EXPECT().Hash(req.Password).Return("", hashingErr).Once()

		userService := services.NewUserService(repoManager, mockLogger, mockHashService)

		userResponse, err := userService.Create(ctx, req)

		require.Error(t, err)
		assert.Nil(t, userResponse)
		assert.Contains(t, err.Error(), "failed to hash password")
	})
}
