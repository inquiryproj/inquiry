// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	context "context"

	app "github.com/inquiryproj/inquiry/internal/app"

	mock "github.com/stretchr/testify/mock"
)

// Runner is an autogenerated mock type for the Runner type
type Runner struct {
	mock.Mock
}

// GetRunsForProject provides a mock function with given fields: ctx, getRunsForProjectRequest
func (_m *Runner) GetRunsForProject(ctx context.Context, getRunsForProjectRequest *app.GetRunsForProjectRequest) (*app.GetRunsForProjectResponse, error) {
	ret := _m.Called(ctx, getRunsForProjectRequest)

	var r0 *app.GetRunsForProjectResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *app.GetRunsForProjectRequest) (*app.GetRunsForProjectResponse, error)); ok {
		return rf(ctx, getRunsForProjectRequest)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *app.GetRunsForProjectRequest) *app.GetRunsForProjectResponse); ok {
		r0 = rf(ctx, getRunsForProjectRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*app.GetRunsForProjectResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *app.GetRunsForProjectRequest) error); ok {
		r1 = rf(ctx, getRunsForProjectRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RunProject provides a mock function with given fields: ctx, run
func (_m *Runner) RunProject(ctx context.Context, run *app.RunProjectRequest) (*app.ProjectRunOutput, error) {
	ret := _m.Called(ctx, run)

	var r0 *app.ProjectRunOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *app.RunProjectRequest) (*app.ProjectRunOutput, error)); ok {
		return rf(ctx, run)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *app.RunProjectRequest) *app.ProjectRunOutput); ok {
		r0 = rf(ctx, run)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*app.ProjectRunOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *app.RunProjectRequest) error); ok {
		r1 = rf(ctx, run)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RunProjectByName provides a mock function with given fields: ctx, run
func (_m *Runner) RunProjectByName(ctx context.Context, run *app.RunProjectByNameRequest) (*app.ProjectRunOutput, error) {
	ret := _m.Called(ctx, run)

	var r0 *app.ProjectRunOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *app.RunProjectByNameRequest) (*app.ProjectRunOutput, error)); ok {
		return rf(ctx, run)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *app.RunProjectByNameRequest) *app.ProjectRunOutput); ok {
		r0 = rf(ctx, run)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*app.ProjectRunOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *app.RunProjectByNameRequest) error); ok {
		r1 = rf(ctx, run)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewRunner creates a new instance of Runner. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRunner(t interface {
	mock.TestingT
	Cleanup(func())
}) *Runner {
	mock := &Runner{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
