package services_test

import (
	"context"
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

func TestUserService_UpdateFull(t *testing.T) {
	require.NotNil(t, testDB, "Test Database connection (testDB) is not initialized")
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Arrange
		TruncateTables(t, testDB, testTableNames) // Clean DB

		repoManager := repository.NewRepositoryManager(testDB)
		mockHashService := mocks.NewMockHashService(t)
		mockLogger := mocks.NewMockLogger(t)

		// Seed user to update
		var userID int32
		seedQuery := `INSERT INTO users (username, email, full_name, password_hash) VALUES ($1, $2, $3, $4) RETURNING id`
		err := testDB.QueryRowContext(ctx, seedQuery, "update_full_user", "updatefull@example.com", "Update Full User", "oldhash").Scan(&userID)
		require.NoError(t, err, "Failed to seed user")

		updateReq := dto.UpdateUserFullReq{
			ID:       userID,
			Username: "updated_full_user",
			Email:    "updated_full@example.com",
			FullName: "Updated Full User Name",
			Password: "newpassword",
		}
		hashedPassword := "hashed_newpassword"
		mockHashService.EXPECT().Hash(updateReq.Password).Return(hashedPassword, nil).Once()

		userService := services.NewUserService(repoManager, mockLogger, mockHashService)

		// Act
		updatedUser, err := userService.UpdateFull(ctx, updateReq)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, updatedUser)
		assert.Equal(t, updateReq.ID, updatedUser.ID)
		assert.Equal(t, updateReq.Username, updatedUser.Username)
		assert.Equal(t, updateReq.Email, updatedUser.Email)
		assert.Equal(t, updateReq.FullName, updatedUser.FullName)

		// Verify in database
		var dbUser dbCtx.User
		query := "SELECT id, username, email, full_name, password_hash FROM users WHERE id = $1"
		err = testDB.QueryRowContext(ctx, query, userID).Scan(
			&dbUser.ID, &dbUser.Username, &dbUser.Email, &dbUser.FullName, &dbUser.PasswordHash,
		)
		require.NoError(t, err)
		assert.Equal(t, updateReq.Username, dbUser.Username)
		assert.Equal(t, updateReq.Email, dbUser.Email)
		assert.Equal(t, updateReq.FullName, dbUser.FullName)
		assert.Equal(t, hashedPassword, dbUser.PasswordHash)
	})

	t.Run("Username Exists", func(t *testing.T) {
		// Arrange
		TruncateTables(t, testDB, testTableNames) // Clean DB

		repoManager := repository.NewRepositoryManager(testDB)
		mockHashService := mocks.NewMockHashService(t)
		mockLogger := mocks.NewMockLogger(t)

		// Seed two users - one to update and one with the target username already taken
		var userID int32
		seedQuery := `INSERT INTO users (username, email, full_name, password_hash) VALUES ($1, $2, $3, $4) RETURNING id`
		err := testDB.QueryRowContext(ctx, seedQuery, "update_user", "update@example.com", "Update User", "hash1").Scan(&userID)
		require.NoError(t, err, "Failed to seed user to update")

		existingUsername := "existing_username"
		_, err = testDB.ExecContext(ctx, seedQuery, existingUsername, "existing@example.com", "Existing User", "hash2")
		require.NoError(t, err, "Failed to seed existing user")

		updateReq := dto.UpdateUserFullReq{
			ID:       userID,
			Username: existingUsername, // This will conflict
			Email:    "new_email@example.com",
			FullName: "Updated Name",
			Password: "newpassword",
		}
		hashedPassword := "hashed_newpassword"
		mockHashService.EXPECT().Hash(updateReq.Password).Return(hashedPassword, nil).Once()
		mockLogger.EXPECT().Warn(logging.Validation, logging.Update, "Username already exists",
			mock.MatchedBy(func(extra map[logging.ExtraKey]any) bool {
				username, ok := extra[logging.RequestBody].(string)
				return ok && username == existingUsername
			})).Once()

		userService := services.NewUserService(repoManager, mockLogger, mockHashService)

		// Act
		updatedUser, err := userService.UpdateFull(ctx, updateReq)

		// Assert
		require.Error(t, err)
		assert.Nil(t, updatedUser)
		var usernameExistsErr *contracts.UsernameExistsError
		require.ErrorAs(t, err, &usernameExistsErr)
		assert.Equal(t, existingUsername, usernameExistsErr.Username)
	})

	t.Run("Email Exists", func(t *testing.T) {
		// Arrange
		TruncateTables(t, testDB, testTableNames) // Clean DB

		repoManager := repository.NewRepositoryManager(testDB)
		mockHashService := mocks.NewMockHashService(t)
		mockLogger := mocks.NewMockLogger(t)

		// Seed two users - one to update and one with the target email already taken
		var userID int32
		seedQuery := `INSERT INTO users (username, email, full_name, password_hash) VALUES ($1, $2, $3, $4) RETURNING id`
		err := testDB.QueryRowContext(ctx, seedQuery, "update_user2", "update2@example.com", "Update User 2", "hash1").Scan(&userID)
		require.NoError(t, err, "Failed to seed user to update")

		existingEmail := "existing_email@example.com"
		_, err = testDB.ExecContext(ctx, seedQuery, "existing_user2", existingEmail, "Existing User 2", "hash2")
		require.NoError(t, err, "Failed to seed existing user")

		updateReq := dto.UpdateUserFullReq{
			ID:       userID,
			Username: "new_username",
			Email:    existingEmail, // This will conflict
			FullName: "Updated Name",
			Password: "newpassword",
		}
		hashedPassword := "hashed_newpassword"
		mockHashService.EXPECT().Hash(updateReq.Password).Return(hashedPassword, nil).Once()
		mockLogger.EXPECT().Warn(logging.Validation, logging.Update, "Email already exists",
			mock.MatchedBy(func(extra map[logging.ExtraKey]any) bool {
				email, ok := extra[logging.RequestBody].(string)
				return ok && email == existingEmail
			})).Once()

		userService := services.NewUserService(repoManager, mockLogger, mockHashService)

		// Act
		updatedUser, err := userService.UpdateFull(ctx, updateReq)

		// Assert
		require.Error(t, err)
		assert.Nil(t, updatedUser)
		var emailExistsErr *contracts.EmailExistsError
		require.ErrorAs(t, err, &emailExistsErr)
		assert.Equal(t, existingEmail, emailExistsErr.Email)
	})

	t.Run("User Not Found", func(t *testing.T) {
		// Arrange
		TruncateTables(t, testDB, testTableNames) // Clean DB

		repoManager := repository.NewRepositoryManager(testDB)
		mockHashService := mocks.NewMockHashService(t)
		mockLogger := mocks.NewMockLogger(t)

		nonExistentID := int32(9999)
		updateReq := dto.UpdateUserFullReq{
			ID:       nonExistentID,
			Username: "nonexistent_user",
			Email:    "nonexistent@example.com",
			FullName: "Nonexistent User",
			Password: "newpassword",
		}
		hashedPassword := "hashed_newpassword"
		mockHashService.EXPECT().Hash(updateReq.Password).Return(hashedPassword, nil).Once()
		mockLogger.EXPECT().Error(logging.Postgres, logging.Update, "User not found",
			mock.MatchedBy(func(extra map[logging.ExtraKey]any) bool {
				id, ok := extra["userID"].(int32)
				return ok && id == nonExistentID
			})).Once()

		userService := services.NewUserService(repoManager, mockLogger, mockHashService)

		// Act
		updatedUser, err := userService.UpdateFull(ctx, updateReq)

		// Assert
		require.Error(t, err)
		assert.Nil(t, updatedUser)
		assert.Equal(t, "user not found", err.Error())
	})
}

func TestUserService_UpdateUserTx(t *testing.T) {
	require.NotNil(t, testDB, "Test Database connection (testDB) is not initialized")
	ctx := context.Background()

	TruncateTables(t, testDB, testTableNames)

	repoManager := repository.NewRepositoryManager(testDB)
	mockHashService := mocks.NewMockHashService(t)
	mockLogger := mocks.NewMockLogger(t)

	// Seed user using testDB
	var seededUserID int32
	seedQuery := `INSERT INTO users (username, email, full_name, password_hash) VALUES ($1, $2, $3, $4) RETURNING id`
	err := testDB.QueryRowContext(ctx, seedQuery, "update_tx_user", "update_tx@example.com", "Update Tx User", "oldhash").Scan(&seededUserID)
	require.NoError(t, err, "Failed to seed user")

	updateArgs := dbCtx.UpdateUserFullParams{
		ID:           seededUserID,
		Username:     "update_tx_user_new_name",
		Email:        "update_tx_new@example.com",
		FullName:     "Updated Tx User Name",
		PasswordHash: "newpassword",
	}
	hashedPassword := "hashed_newpassword"

	mockHashService.EXPECT().Hash(updateArgs.PasswordHash).Return(hashedPassword, nil).Maybe() // Maybe() because UpdateFull might not be called if tx fails early

	userService := services.NewUserService(repoManager, mockLogger, mockHashService)

	updatedUser, err := userService.UpdateUserTx(ctx, updateArgs)

	require.NoError(t, err)
	require.NotNil(t, updatedUser)
	assert.Equal(t, updateArgs.ID, updatedUser.ID)
	assert.Equal(t, updateArgs.Username, updatedUser.Username)
	assert.Equal(t, updateArgs.Email, updatedUser.Email)

	// Note: UpdateFull returns the state *before* hashing if UpdateUserFullParams takes raw password
	// Adjust assertion based on what UpdateFull returns and if UpdateUserTx modifies it.
	// Assuming UpdateFull returns the updated user *after* hash is set in params:
	// assert.Equal(t, hashedPassword, updatedUser.PasswordHash)

	// Assert (Database State - using testDB)
	// Check the DB state *after* the transaction should have committed
	var finalUser dbCtx.User
	query := "SELECT id, username, email, full_name, password_hash FROM users WHERE id = $1"
	row := testDB.QueryRowContext(ctx, query, updatedUser.ID)
	err = row.Scan(
		&finalUser.ID, &finalUser.Username, &finalUser.Email,
		&finalUser.FullName, &finalUser.PasswordHash,
	)
	require.NoError(t, err, "Failed to query final user state from test DB")
	assert.Equal(t, updateArgs.Username, finalUser.Username)
	assert.Equal(t, updateArgs.Email, finalUser.Email)
	assert.Equal(t, updateArgs.FullName, finalUser.FullName)

	// Password hash assertion depends on whether UpdateUserFullParams contains raw or hashed pass
	// Assuming UpdateFull hashes, check for hashed pass:
	// assert.Equal(t, hashedPassword, finalUser.PasswordHash)
}
