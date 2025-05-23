// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package mocks

import (
	"context"

	dto "example.com/api/internal/contracts"
	dbCtx "example.com/api/internal/repository/db"
	"github.com/golang-jwt/jwt/v5"
	mock "github.com/stretchr/testify/mock"
)

// NewMockAuthService creates a new instance of MockAuthService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAuthService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAuthService {
	mock := &MockAuthService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockAuthService is an autogenerated mock type for the IAuthService type
type MockAuthService struct {
	mock.Mock
}

type MockAuthService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAuthService) EXPECT() *MockAuthService_Expecter {
	return &MockAuthService_Expecter{mock: &_m.Mock}
}

// Authenticate provides a mock function for the type MockAuthService
func (_mock *MockAuthService) Authenticate(ctx context.Context, email string, password string) (*dbCtx.User, error) {
	ret := _mock.Called(ctx, email, password)

	if len(ret) == 0 {
		panic("no return value specified for Authenticate")
	}

	var r0 *dbCtx.User
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, string, string) (*dbCtx.User, error)); ok {
		return returnFunc(ctx, email, password)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, string, string) *dbCtx.User); ok {
		r0 = returnFunc(ctx, email, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dbCtx.User)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = returnFunc(ctx, email, password)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockAuthService_Authenticate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Authenticate'
type MockAuthService_Authenticate_Call struct {
	*mock.Call
}

// Authenticate is a helper method to define mock.On call
//   - ctx
//   - email
//   - password
func (_e *MockAuthService_Expecter) Authenticate(ctx interface{}, email interface{}, password interface{}) *MockAuthService_Authenticate_Call {
	return &MockAuthService_Authenticate_Call{Call: _e.mock.On("Authenticate", ctx, email, password)}
}

func (_c *MockAuthService_Authenticate_Call) Run(run func(ctx context.Context, email string, password string)) *MockAuthService_Authenticate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockAuthService_Authenticate_Call) Return(user *dbCtx.User, err error) *MockAuthService_Authenticate_Call {
	_c.Call.Return(user, err)
	return _c
}

func (_c *MockAuthService_Authenticate_Call) RunAndReturn(run func(ctx context.Context, email string, password string) (*dbCtx.User, error)) *MockAuthService_Authenticate_Call {
	_c.Call.Return(run)
	return _c
}

// GenerateAccessToken provides a mock function for the type MockAuthService
func (_mock *MockAuthService) GenerateAccessToken(userID string) (string, error) {
	ret := _mock.Called(userID)

	if len(ret) == 0 {
		panic("no return value specified for GenerateAccessToken")
	}

	var r0 string
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(string) (string, error)); ok {
		return returnFunc(userID)
	}
	if returnFunc, ok := ret.Get(0).(func(string) string); ok {
		r0 = returnFunc(userID)
	} else {
		r0 = ret.Get(0).(string)
	}
	if returnFunc, ok := ret.Get(1).(func(string) error); ok {
		r1 = returnFunc(userID)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockAuthService_GenerateAccessToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GenerateAccessToken'
type MockAuthService_GenerateAccessToken_Call struct {
	*mock.Call
}

// GenerateAccessToken is a helper method to define mock.On call
//   - userID
func (_e *MockAuthService_Expecter) GenerateAccessToken(userID interface{}) *MockAuthService_GenerateAccessToken_Call {
	return &MockAuthService_GenerateAccessToken_Call{Call: _e.mock.On("GenerateAccessToken", userID)}
}

func (_c *MockAuthService_GenerateAccessToken_Call) Run(run func(userID string)) *MockAuthService_GenerateAccessToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockAuthService_GenerateAccessToken_Call) Return(s string, err error) *MockAuthService_GenerateAccessToken_Call {
	_c.Call.Return(s, err)
	return _c
}

func (_c *MockAuthService_GenerateAccessToken_Call) RunAndReturn(run func(userID string) (string, error)) *MockAuthService_GenerateAccessToken_Call {
	_c.Call.Return(run)
	return _c
}

// GenerateRefreshToken provides a mock function for the type MockAuthService
func (_mock *MockAuthService) GenerateRefreshToken(userID string) (string, error) {
	ret := _mock.Called(userID)

	if len(ret) == 0 {
		panic("no return value specified for GenerateRefreshToken")
	}

	var r0 string
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(string) (string, error)); ok {
		return returnFunc(userID)
	}
	if returnFunc, ok := ret.Get(0).(func(string) string); ok {
		r0 = returnFunc(userID)
	} else {
		r0 = ret.Get(0).(string)
	}
	if returnFunc, ok := ret.Get(1).(func(string) error); ok {
		r1 = returnFunc(userID)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockAuthService_GenerateRefreshToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GenerateRefreshToken'
type MockAuthService_GenerateRefreshToken_Call struct {
	*mock.Call
}

// GenerateRefreshToken is a helper method to define mock.On call
//   - userID
func (_e *MockAuthService_Expecter) GenerateRefreshToken(userID interface{}) *MockAuthService_GenerateRefreshToken_Call {
	return &MockAuthService_GenerateRefreshToken_Call{Call: _e.mock.On("GenerateRefreshToken", userID)}
}

func (_c *MockAuthService_GenerateRefreshToken_Call) Run(run func(userID string)) *MockAuthService_GenerateRefreshToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockAuthService_GenerateRefreshToken_Call) Return(s string, err error) *MockAuthService_GenerateRefreshToken_Call {
	_c.Call.Return(s, err)
	return _c
}

func (_c *MockAuthService_GenerateRefreshToken_Call) RunAndReturn(run func(userID string) (string, error)) *MockAuthService_GenerateRefreshToken_Call {
	_c.Call.Return(run)
	return _c
}

// RegisterUser provides a mock function for the type MockAuthService
func (_mock *MockAuthService) RegisterUser(ctx context.Context, args dto.Register) (*dto.UserResponse, string, string, error) {
	ret := _mock.Called(ctx, args)

	if len(ret) == 0 {
		panic("no return value specified for RegisterUser")
	}

	var r0 *dto.UserResponse
	var r1 string
	var r2 string
	var r3 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, dto.Register) (*dto.UserResponse, string, string, error)); ok {
		return returnFunc(ctx, args)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, dto.Register) *dto.UserResponse); ok {
		r0 = returnFunc(ctx, args)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.UserResponse)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, dto.Register) string); ok {
		r1 = returnFunc(ctx, args)
	} else {
		r1 = ret.Get(1).(string)
	}
	if returnFunc, ok := ret.Get(2).(func(context.Context, dto.Register) string); ok {
		r2 = returnFunc(ctx, args)
	} else {
		r2 = ret.Get(2).(string)
	}
	if returnFunc, ok := ret.Get(3).(func(context.Context, dto.Register) error); ok {
		r3 = returnFunc(ctx, args)
	} else {
		r3 = ret.Error(3)
	}
	return r0, r1, r2, r3
}

// MockAuthService_RegisterUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RegisterUser'
type MockAuthService_RegisterUser_Call struct {
	*mock.Call
}

// RegisterUser is a helper method to define mock.On call
//   - ctx
//   - args
func (_e *MockAuthService_Expecter) RegisterUser(ctx interface{}, args interface{}) *MockAuthService_RegisterUser_Call {
	return &MockAuthService_RegisterUser_Call{Call: _e.mock.On("RegisterUser", ctx, args)}
}

func (_c *MockAuthService_RegisterUser_Call) Run(run func(ctx context.Context, args dto.Register)) *MockAuthService_RegisterUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(dto.Register))
	})
	return _c
}

func (_c *MockAuthService_RegisterUser_Call) Return(userResponse *dto.UserResponse, s string, s1 string, err error) *MockAuthService_RegisterUser_Call {
	_c.Call.Return(userResponse, s, s1, err)
	return _c
}

func (_c *MockAuthService_RegisterUser_Call) RunAndReturn(run func(ctx context.Context, args dto.Register) (*dto.UserResponse, string, string, error)) *MockAuthService_RegisterUser_Call {
	_c.Call.Return(run)
	return _c
}

// RotateTokens provides a mock function for the type MockAuthService
func (_mock *MockAuthService) RotateTokens(ctx context.Context, refreshToken string) (string, string, error) {
	ret := _mock.Called(ctx, refreshToken)

	if len(ret) == 0 {
		panic("no return value specified for RotateTokens")
	}

	var r0 string
	var r1 string
	var r2 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, string) (string, string, error)); ok {
		return returnFunc(ctx, refreshToken)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = returnFunc(ctx, refreshToken)
	} else {
		r0 = ret.Get(0).(string)
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, string) string); ok {
		r1 = returnFunc(ctx, refreshToken)
	} else {
		r1 = ret.Get(1).(string)
	}
	if returnFunc, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = returnFunc(ctx, refreshToken)
	} else {
		r2 = ret.Error(2)
	}
	return r0, r1, r2
}

// MockAuthService_RotateTokens_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RotateTokens'
type MockAuthService_RotateTokens_Call struct {
	*mock.Call
}

// RotateTokens is a helper method to define mock.On call
//   - ctx
//   - refreshToken
func (_e *MockAuthService_Expecter) RotateTokens(ctx interface{}, refreshToken interface{}) *MockAuthService_RotateTokens_Call {
	return &MockAuthService_RotateTokens_Call{Call: _e.mock.On("RotateTokens", ctx, refreshToken)}
}

func (_c *MockAuthService_RotateTokens_Call) Run(run func(ctx context.Context, refreshToken string)) *MockAuthService_RotateTokens_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockAuthService_RotateTokens_Call) Return(s string, s1 string, err error) *MockAuthService_RotateTokens_Call {
	_c.Call.Return(s, s1, err)
	return _c
}

func (_c *MockAuthService_RotateTokens_Call) RunAndReturn(run func(ctx context.Context, refreshToken string) (string, string, error)) *MockAuthService_RotateTokens_Call {
	_c.Call.Return(run)
	return _c
}

// ValidateAccessToken provides a mock function for the type MockAuthService
func (_mock *MockAuthService) ValidateAccessToken(tokenString string) (jwt.MapClaims, error) {
	ret := _mock.Called(tokenString)

	if len(ret) == 0 {
		panic("no return value specified for ValidateAccessToken")
	}

	var r0 jwt.MapClaims
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(string) (jwt.MapClaims, error)); ok {
		return returnFunc(tokenString)
	}
	if returnFunc, ok := ret.Get(0).(func(string) jwt.MapClaims); ok {
		r0 = returnFunc(tokenString)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(jwt.MapClaims)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(string) error); ok {
		r1 = returnFunc(tokenString)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockAuthService_ValidateAccessToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ValidateAccessToken'
type MockAuthService_ValidateAccessToken_Call struct {
	*mock.Call
}

// ValidateAccessToken is a helper method to define mock.On call
//   - tokenString
func (_e *MockAuthService_Expecter) ValidateAccessToken(tokenString interface{}) *MockAuthService_ValidateAccessToken_Call {
	return &MockAuthService_ValidateAccessToken_Call{Call: _e.mock.On("ValidateAccessToken", tokenString)}
}

func (_c *MockAuthService_ValidateAccessToken_Call) Run(run func(tokenString string)) *MockAuthService_ValidateAccessToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockAuthService_ValidateAccessToken_Call) Return(mapClaims jwt.MapClaims, err error) *MockAuthService_ValidateAccessToken_Call {
	_c.Call.Return(mapClaims, err)
	return _c
}

func (_c *MockAuthService_ValidateAccessToken_Call) RunAndReturn(run func(tokenString string) (jwt.MapClaims, error)) *MockAuthService_ValidateAccessToken_Call {
	_c.Call.Return(run)
	return _c
}

// ValidateRefreshToken provides a mock function for the type MockAuthService
func (_mock *MockAuthService) ValidateRefreshToken(ctx context.Context, tokenString string) (jwt.MapClaims, error) {
	ret := _mock.Called(ctx, tokenString)

	if len(ret) == 0 {
		panic("no return value specified for ValidateRefreshToken")
	}

	var r0 jwt.MapClaims
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, string) (jwt.MapClaims, error)); ok {
		return returnFunc(ctx, tokenString)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, string) jwt.MapClaims); ok {
		r0 = returnFunc(ctx, tokenString)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(jwt.MapClaims)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = returnFunc(ctx, tokenString)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockAuthService_ValidateRefreshToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ValidateRefreshToken'
type MockAuthService_ValidateRefreshToken_Call struct {
	*mock.Call
}

// ValidateRefreshToken is a helper method to define mock.On call
//   - ctx
//   - tokenString
func (_e *MockAuthService_Expecter) ValidateRefreshToken(ctx interface{}, tokenString interface{}) *MockAuthService_ValidateRefreshToken_Call {
	return &MockAuthService_ValidateRefreshToken_Call{Call: _e.mock.On("ValidateRefreshToken", ctx, tokenString)}
}

func (_c *MockAuthService_ValidateRefreshToken_Call) Run(run func(ctx context.Context, tokenString string)) *MockAuthService_ValidateRefreshToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockAuthService_ValidateRefreshToken_Call) Return(mapClaims jwt.MapClaims, err error) *MockAuthService_ValidateRefreshToken_Call {
	_c.Call.Return(mapClaims, err)
	return _c
}

func (_c *MockAuthService_ValidateRefreshToken_Call) RunAndReturn(run func(ctx context.Context, tokenString string) (jwt.MapClaims, error)) *MockAuthService_ValidateRefreshToken_Call {
	_c.Call.Return(run)
	return _c
}

// ValidateToken provides a mock function for the type MockAuthService
func (_mock *MockAuthService) ValidateToken(tokenString string, expectedType string) (jwt.MapClaims, error) {
	ret := _mock.Called(tokenString, expectedType)

	if len(ret) == 0 {
		panic("no return value specified for ValidateToken")
	}

	var r0 jwt.MapClaims
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(string, string) (jwt.MapClaims, error)); ok {
		return returnFunc(tokenString, expectedType)
	}
	if returnFunc, ok := ret.Get(0).(func(string, string) jwt.MapClaims); ok {
		r0 = returnFunc(tokenString, expectedType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(jwt.MapClaims)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = returnFunc(tokenString, expectedType)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockAuthService_ValidateToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ValidateToken'
type MockAuthService_ValidateToken_Call struct {
	*mock.Call
}

// ValidateToken is a helper method to define mock.On call
//   - tokenString
//   - expectedType
func (_e *MockAuthService_Expecter) ValidateToken(tokenString interface{}, expectedType interface{}) *MockAuthService_ValidateToken_Call {
	return &MockAuthService_ValidateToken_Call{Call: _e.mock.On("ValidateToken", tokenString, expectedType)}
}

func (_c *MockAuthService_ValidateToken_Call) Run(run func(tokenString string, expectedType string)) *MockAuthService_ValidateToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *MockAuthService_ValidateToken_Call) Return(mapClaims jwt.MapClaims, err error) *MockAuthService_ValidateToken_Call {
	_c.Call.Return(mapClaims, err)
	return _c
}

func (_c *MockAuthService_ValidateToken_Call) RunAndReturn(run func(tokenString string, expectedType string) (jwt.MapClaims, error)) *MockAuthService_ValidateToken_Call {
	_c.Call.Return(run)
	return _c
}
