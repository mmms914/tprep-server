// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Session is an autogenerated mock type for the Session type
type Session struct {
	mock.Mock
}

// EndSession provides a mock function with given fields: ctx
func (_m *Session) EndSession(ctx context.Context) {
	_m.Called(ctx)
}

// WithTransaction provides a mock function with given fields: ctx, fn
func (_m *Session) WithTransaction(ctx context.Context, fn func(context.Context) (interface{}, error)) (interface{}, error) {
	ret := _m.Called(ctx, fn)

	if len(ret) == 0 {
		panic("no return value specified for WithTransaction")
	}

	var r0 interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, func(context.Context) (interface{}, error)) (interface{}, error)); ok {
		return rf(ctx, fn)
	}
	if rf, ok := ret.Get(0).(func(context.Context, func(context.Context) (interface{}, error)) interface{}); ok {
		r0 = rf(ctx, fn)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, func(context.Context) (interface{}, error)) error); ok {
		r1 = rf(ctx, fn)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewSession creates a new instance of Session. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSession(t interface {
	mock.TestingT
	Cleanup(func())
}) *Session {
	mock := &Session{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
