package handlers_tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/api/internal/api/handlers"
	"example.com/api/internal/api/responses"
	dto "example.com/api/internal/contracts"
	contracts "example.com/api/internal/contracts/errors"
	dbCtx "example.com/api/internal/repository/db"
	"example.com/api/pkg/logging"
	mocks "example.com/api/tests/unit/mocks/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserHandlerTestSuite struct {
	suite.Suite
	serviceManager *mocks.MockServiceManager
	userService    *mocks.MockUserService
	logger         *mocks.MockLogger
	handler        *handlers.UserHandler
	ctx            *gin.Context
	recorder       *httptest.ResponseRecorder
}

func (suite *UserHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	suite.serviceManager = mocks.NewMockServiceManager(suite.T())
	suite.userService = mocks.NewMockUserService(suite.T())
	suite.logger = mocks.NewMockLogger(suite.T())
	suite.handler = handlers.NewUserHandler(suite.serviceManager, suite.logger)
	suite.recorder = httptest.NewRecorder()
	suite.ctx, _ = gin.CreateTestContext(suite.recorder)
}

func (suite *UserHandlerTestSuite) TestCreate_InvalidRequestBody() {
	reqBody := []byte(`invalid json`)
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	suite.ctx.Request = req

	suite.logger.EXPECT().Error(
		logging.Validation,
		logging.Api,
		"Invalid request body",
		mock.MatchedBy(func(extra map[logging.ExtraKey]any) bool {
			return extra[logging.ErrorMessage] != nil &&
				extra[logging.Path] == "/users" &&
				extra[logging.Method] == http.MethodPost
		}),
	).Once()

	suite.handler.Create(suite.ctx)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *UserHandlerTestSuite) TestCreate_ValidationError() {
	// Intentionally missing required fields
	reqBody := []byte(`{"username": "user"}`)
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	suite.ctx.Request = req

	suite.logger.EXPECT().Error(
		logging.Validation,
		logging.Api,
		"Invalid request body",
		mock.MatchedBy(func(extra map[logging.ExtraKey]any) bool {
			return extra[logging.Path] == "/users" &&
				extra[logging.Method] == http.MethodPost
		}),
	).Once()

	suite.handler.Create(suite.ctx)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *UserHandlerTestSuite) TestCreate_UsernameConflict() {
	// Include all required fields based on your validation
	reqBody := []byte(`{
		"username": "existing",
		"email": "test@example.com",
		"password": "password",
		"fullName": "Test User"
	}`)
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	suite.ctx.Request = req

	expectedErr := &contracts.UsernameExistsError{Username: "existing"}
	suite.serviceManager.EXPECT().User().Return(suite.userService).Once()
	suite.userService.EXPECT().Create(mock.Anything, mock.Anything).Return(nil, expectedErr).Once()

	suite.handler.Create(suite.ctx)

	suite.Equal(http.StatusConflict, suite.recorder.Code)
}

func (suite *UserHandlerTestSuite) TestCreate_EmailConflict() {
	reqBody := []byte(`{
		"username": "newuser",
		"email": "existing@example.com",
		"password": "password",
		"fullName": "Test User"
	}`)
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	suite.ctx.Request = req

	expectedErr := &contracts.EmailExistsError{Email: "existing@example.com"}
	suite.serviceManager.EXPECT().User().Return(suite.userService).Once()
	suite.userService.EXPECT().Create(mock.Anything, mock.Anything).Return(nil, expectedErr).Once()

	suite.handler.Create(suite.ctx)

	suite.Equal(http.StatusConflict, suite.recorder.Code)
}

func (suite *UserHandlerTestSuite) TestCreate_InternalError() {
	reqBody := []byte(`{
		"username": "user",
		"email": "test@example.com",
		"password": "password",
		"fullName": "Test User"
	}`)
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	suite.ctx.Request = req

	expectedErr := errors.New("internal server error")
	suite.serviceManager.EXPECT().User().Return(suite.userService).Once()
	suite.userService.EXPECT().Create(mock.Anything, mock.Anything).Return(nil, expectedErr).Once()

	suite.logger.EXPECT().Error(
		logging.Internal,
		logging.FailedToCreateUser,
		"Failed to create user",
		mock.MatchedBy(func(extra map[logging.ExtraKey]any) bool {
			return extra[logging.ErrorMessage] == expectedErr.Error() &&
				extra[logging.Path] == "/users" &&
				extra[logging.Method] == http.MethodPost
		}),
	).Once()

	suite.handler.Create(suite.ctx)

	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
}

func (suite *UserHandlerTestSuite) TestCreate_Success() {
	reqBody := []byte(`{
		"username": "newuser",
		"email": "new@example.com",
		"password": "password",
		"fullName": "Test User"
	}`)
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	suite.ctx.Request = req

	expectedUser := &dto.UserResponse{
		ID:       1,
		Username: "newuser",
		Email:    "new@example.com",
		FullName: "Test User",
	}

	suite.serviceManager.EXPECT().User().Return(suite.userService).Once()
	suite.userService.EXPECT().Create(
		mock.Anything,
		mock.MatchedBy(func(req dto.CreateUserReq) bool {
			return req.Username == "newuser" &&
				req.Email == "new@example.com" &&
				req.Password == "password" &&
				req.FullName == "Test User"
		}),
	).Return(expectedUser, nil).Once()

	suite.logger.EXPECT().Info(
		logging.Internal,
		logging.Api,
		"User created successfully",
		mock.MatchedBy(func(extra map[logging.ExtraKey]any) bool {
			return extra[logging.Path] == "/users" &&
				extra[logging.Method] == http.MethodPost
		}),
	).Once()

	suite.handler.Create(suite.ctx)

	suite.Equal(http.StatusCreated, suite.recorder.Code)
}

func (suite *UserHandlerTestSuite) TestGetByID_InvalidID() {
	// Setup invalid ID parameter
	suite.ctx.Params = []gin.Param{{Key: "id", Value: "invalid"}}
	req, _ := http.NewRequest(http.MethodGet, "/users/invalid", nil)
	suite.ctx.Request = req

	suite.handler.GetByID(suite.ctx)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)

	var response responses.BaseResponse
	err := json.Unmarshal(suite.recorder.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("fail", response.Status)
	suite.Equal("Invalid user ID, must be an integer", response.Message)
}

func (suite *UserHandlerTestSuite) TestGetByID_NotFound() {
	// Setup valid ID parameter
	suite.ctx.Params = []gin.Param{{Key: "id", Value: "123"}}
	req, _ := http.NewRequest(http.MethodGet, "/users/123", nil)
	suite.ctx.Request = req

	// Mock service response
	expectedErr := errors.New("user not found")
	suite.serviceManager.EXPECT().User().Return(suite.userService).Once()
	suite.userService.EXPECT().GetByID(mock.Anything, int32(123)).Return(nil, expectedErr).Once()

	suite.handler.GetByID(suite.ctx)

	suite.Equal(http.StatusNotFound, suite.recorder.Code)

	var response responses.BaseResponse
	err := json.Unmarshal(suite.recorder.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("fail", response.Status)
	suite.Equal("User not found", response.Message)
}

func (suite *UserHandlerTestSuite) TestGetByID_Success() {
	// Setup valid ID parameter
	suite.ctx.Params = []gin.Param{{Key: "id", Value: "123"}}
	req, _ := http.NewRequest(http.MethodGet, "/users/123", nil)
	suite.ctx.Request = req

	// Mock service response
	expectedUser := &dbCtx.User{
		ID:       123,
		Username: "testuser",
		Email:    "test@example.com",
	}

	suite.serviceManager.EXPECT().User().Return(suite.userService).Once()
	suite.userService.EXPECT().GetByID(mock.Anything, int32(123)).Return(expectedUser, nil).Once()

	suite.handler.GetByID(suite.ctx)

	suite.Equal(http.StatusOK, suite.recorder.Code)

	var response responses.BaseResponse
	err := json.Unmarshal(suite.recorder.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("success", response.Status)
	suite.Equal("User retrieved successfully", response.Message)

	// Convert response data to User struct
	var actualUser dbCtx.User
	dataBytes, _ := json.Marshal(response.Data)
	err = json.Unmarshal(dataBytes, &actualUser)
	suite.NoError(err)

	// Compare relevant fields ignoring internal database fields
	suite.Equal(expectedUser.ID, actualUser.ID)
	suite.Equal(expectedUser.Username, actualUser.Username)
	suite.Equal(expectedUser.Email, actualUser.Email)
}
func (suite *UserHandlerTestSuite) TestDelete_InvalidID() {
	// Setup invalid ID (non-integer)
	suite.ctx.Params = []gin.Param{{Key: "id", Value: "abc"}}
	req, _ := http.NewRequest(http.MethodDelete, "/users/abc", nil)
	suite.ctx.Request = req

	suite.handler.DeleteUser(suite.ctx)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)

	var response responses.BaseResponse
	err := json.Unmarshal(suite.recorder.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("fail", response.Status)
	suite.Equal("Invalid user ID, must be an integer", response.Message)
}

func (suite *UserHandlerTestSuite) TestDelete_UserNotFound() {
	// Setup valid ID
	suite.ctx.Params = []gin.Param{{Key: "id", Value: "123"}}
	req, _ := http.NewRequest(http.MethodDelete, "/users/123", nil)
	suite.ctx.Request = req

	// Mock service response
	expectedErr := errors.New("user not found")
	suite.serviceManager.EXPECT().User().Return(suite.userService).Once()
	suite.userService.EXPECT().SoftDelete(mock.Anything, int32(123)).Return(expectedErr).Once()

	suite.handler.DeleteUser(suite.ctx)

	suite.Equal(http.StatusNotFound, suite.recorder.Code)

	var response responses.BaseResponse
	err := json.Unmarshal(suite.recorder.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("fail", response.Status)
	suite.Equal("User not found", response.Message)
}

func (suite *UserHandlerTestSuite) TestDelete_InternalError() {
	// Setup valid ID
	suite.ctx.Params = []gin.Param{{Key: "id", Value: "123"}}
	req, _ := http.NewRequest(http.MethodDelete, "/users/123", nil)
	suite.ctx.Request = req

	// Mock service response
	expectedErr := errors.New("database error")
	suite.serviceManager.EXPECT().User().Return(suite.userService).Once()
	suite.userService.EXPECT().SoftDelete(mock.Anything, int32(123)).Return(expectedErr).Once()

	suite.handler.DeleteUser(suite.ctx)

	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)

	var response responses.BaseResponse
	err := json.Unmarshal(suite.recorder.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("error", response.Status)
	suite.Equal("Failed to delete user", response.Message)
}

func (suite *UserHandlerTestSuite) TestDelete_Success() {
	// Setup valid ID
	suite.ctx.Params = []gin.Param{{Key: "id", Value: "123"}}
	req, _ := http.NewRequest(http.MethodDelete, "/users/123", nil)
	suite.ctx.Request = req

	// Mock service response
	suite.serviceManager.EXPECT().User().Return(suite.userService).Once()
	suite.userService.EXPECT().SoftDelete(mock.Anything, int32(123)).Return(nil).Once()

	suite.handler.DeleteUser(suite.ctx)

	suite.Equal(http.StatusNoContent, suite.recorder.Code)
	suite.Empty(suite.recorder.Body.Bytes())
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}
