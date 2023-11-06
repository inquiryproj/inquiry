// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	uuid "github.com/google/uuid"
	echo "github.com/labstack/echo/v4"
	mock "github.com/stretchr/testify/mock"

	api "github.com/inquiryproj/inquiry/internal/http/api"
)

// ServerInterface is an autogenerated mock type for the ServerInterface type
type ServerInterface struct {
	mock.Mock
}

// CreateProject provides a mock function with given fields: ctx
func (_m *ServerInterface) CreateProject(ctx echo.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateScenario provides a mock function with given fields: ctx, projectId
func (_m *ServerInterface) CreateScenario(ctx echo.Context, projectId uuid.UUID) error {
	ret := _m.Called(ctx, projectId)

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, projectId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListProjects provides a mock function with given fields: ctx, params
func (_m *ServerInterface) ListProjects(ctx echo.Context, params api.ListProjectsParams) error {
	ret := _m.Called(ctx, params)

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context, api.ListProjectsParams) error); ok {
		r0 = rf(ctx, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListRunsForProject provides a mock function with given fields: ctx, id, params
func (_m *ServerInterface) ListRunsForProject(ctx echo.Context, id uuid.UUID, params api.ListRunsForProjectParams) error {
	ret := _m.Called(ctx, id, params)

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context, uuid.UUID, api.ListRunsForProjectParams) error); ok {
		r0 = rf(ctx, id, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListScenariosForProject provides a mock function with given fields: ctx, projectId, params
func (_m *ServerInterface) ListScenariosForProject(ctx echo.Context, projectId uuid.UUID, params api.ListScenariosForProjectParams) error {
	ret := _m.Called(ctx, projectId, params)

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context, uuid.UUID, api.ListScenariosForProjectParams) error); ok {
		r0 = rf(ctx, projectId, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RunProject provides a mock function with given fields: ctx
func (_m *ServerInterface) RunProject(ctx echo.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewServerInterface creates a new instance of ServerInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewServerInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *ServerInterface {
	mock := &ServerInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
