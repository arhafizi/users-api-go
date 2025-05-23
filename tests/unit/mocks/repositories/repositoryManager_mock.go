// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package mocks

import (
	"context"

	"example.com/api/internal/repository"
	mock "github.com/stretchr/testify/mock"
)

// NewMockRepositoryManager creates a new instance of MockRepositoryManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRepositoryManager(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRepositoryManager {
	mock := &MockRepositoryManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockRepositoryManager is an autogenerated mock type for the IRepositoryManager type
type MockRepositoryManager struct {
	mock.Mock
}

type MockRepositoryManager_Expecter struct {
	mock *mock.Mock
}

func (_m *MockRepositoryManager) EXPECT() *MockRepositoryManager_Expecter {
	return &MockRepositoryManager_Expecter{mock: &_m.Mock}
}

// Chat provides a mock function for the type MockRepositoryManager
func (_mock *MockRepositoryManager) Chat() repository.IChatRepo {
	ret := _mock.Called()

	if len(ret) == 0 {
		panic("no return value specified for Chat")
	}

	var r0 repository.IChatRepo
	if returnFunc, ok := ret.Get(0).(func() repository.IChatRepo); ok {
		r0 = returnFunc()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(repository.IChatRepo)
		}
	}
	return r0
}

// MockRepositoryManager_Chat_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Chat'
type MockRepositoryManager_Chat_Call struct {
	*mock.Call
}

// Chat is a helper method to define mock.On call
func (_e *MockRepositoryManager_Expecter) Chat() *MockRepositoryManager_Chat_Call {
	return &MockRepositoryManager_Chat_Call{Call: _e.mock.On("Chat")}
}

func (_c *MockRepositoryManager_Chat_Call) Run(run func()) *MockRepositoryManager_Chat_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockRepositoryManager_Chat_Call) Return(iChatRepo repository.IChatRepo) *MockRepositoryManager_Chat_Call {
	_c.Call.Return(iChatRepo)
	return _c
}

func (_c *MockRepositoryManager_Chat_Call) RunAndReturn(run func() repository.IChatRepo) *MockRepositoryManager_Chat_Call {
	_c.Call.Return(run)
	return _c
}

// User provides a mock function for the type MockRepositoryManager
func (_mock *MockRepositoryManager) User() repository.IUserRepo {
	ret := _mock.Called()

	if len(ret) == 0 {
		panic("no return value specified for User")
	}

	var r0 repository.IUserRepo
	if returnFunc, ok := ret.Get(0).(func() repository.IUserRepo); ok {
		r0 = returnFunc()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(repository.IUserRepo)
		}
	}
	return r0
}

// MockRepositoryManager_User_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'User'
type MockRepositoryManager_User_Call struct {
	*mock.Call
}

// User is a helper method to define mock.On call
func (_e *MockRepositoryManager_Expecter) User() *MockRepositoryManager_User_Call {
	return &MockRepositoryManager_User_Call{Call: _e.mock.On("User")}
}

func (_c *MockRepositoryManager_User_Call) Run(run func()) *MockRepositoryManager_User_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockRepositoryManager_User_Call) Return(iUserRepo repository.IUserRepo) *MockRepositoryManager_User_Call {
	_c.Call.Return(iUserRepo)
	return _c
}

func (_c *MockRepositoryManager_User_Call) RunAndReturn(run func() repository.IUserRepo) *MockRepositoryManager_User_Call {
	_c.Call.Return(run)
	return _c
}

// WithTx provides a mock function for the type MockRepositoryManager
func (_mock *MockRepositoryManager) WithTx(context1 context.Context, fn func(repository.IRepositoryManager) error) error {
	ret := _mock.Called(context1, fn)

	if len(ret) == 0 {
		panic("no return value specified for WithTx")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, func(repository.IRepositoryManager) error) error); ok {
		r0 = returnFunc(context1, fn)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockRepositoryManager_WithTx_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithTx'
type MockRepositoryManager_WithTx_Call struct {
	*mock.Call
}

// WithTx is a helper method to define mock.On call
//   - context1
//   - fn
func (_e *MockRepositoryManager_Expecter) WithTx(context1 interface{}, fn interface{}) *MockRepositoryManager_WithTx_Call {
	return &MockRepositoryManager_WithTx_Call{Call: _e.mock.On("WithTx", context1, fn)}
}

func (_c *MockRepositoryManager_WithTx_Call) Run(run func(context1 context.Context, fn func(repository.IRepositoryManager) error)) *MockRepositoryManager_WithTx_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(func(repository.IRepositoryManager) error))
	})
	return _c
}

func (_c *MockRepositoryManager_WithTx_Call) Return(err error) *MockRepositoryManager_WithTx_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockRepositoryManager_WithTx_Call) RunAndReturn(run func(context1 context.Context, fn func(repository.IRepositoryManager) error) error) *MockRepositoryManager_WithTx_Call {
	_c.Call.Return(run)
	return _c
}
