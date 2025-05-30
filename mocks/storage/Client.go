// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	context "context"
	io "io"

	mock "github.com/stretchr/testify/mock"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// BucketExists provides a mock function with given fields: ctx, bucketName
func (_m *Client) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	ret := _m.Called(ctx, bucketName)

	if len(ret) == 0 {
		panic("no return value specified for BucketExists")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (bool, error)); ok {
		return rf(ctx, bucketName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, bucketName)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, bucketName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetObject provides a mock function with given fields: ctx, bucketName, objectName
func (_m *Client) GetObject(ctx context.Context, bucketName string, objectName string) ([]byte, error) {
	ret := _m.Called(ctx, bucketName, objectName)

	if len(ret) == 0 {
		panic("no return value specified for GetObject")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) ([]byte, error)); ok {
		return rf(ctx, bucketName, objectName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []byte); ok {
		r0 = rf(ctx, bucketName, objectName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, bucketName, objectName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MakeBucket provides a mock function with given fields: ctx, bucketName
func (_m *Client) MakeBucket(ctx context.Context, bucketName string) error {
	ret := _m.Called(ctx, bucketName)

	if len(ret) == 0 {
		panic("no return value specified for MakeBucket")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, bucketName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PutObject provides a mock function with given fields: ctx, bucketName, objectName, reader, objectSize
func (_m *Client) PutObject(ctx context.Context, bucketName string, objectName string, reader io.Reader, objectSize int64) error {
	ret := _m.Called(ctx, bucketName, objectName, reader, objectSize)

	if len(ret) == 0 {
		panic("no return value specified for PutObject")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, io.Reader, int64) error); ok {
		r0 = rf(ctx, bucketName, objectName, reader, objectSize)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveObject provides a mock function with given fields: ctx, bucketName, objectName
func (_m *Client) RemoveObject(ctx context.Context, bucketName string, objectName string) error {
	ret := _m.Called(ctx, bucketName, objectName)

	if len(ret) == 0 {
		panic("no return value specified for RemoveObject")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, bucketName, objectName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewClient creates a new instance of Client. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *Client {
	mock := &Client{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
