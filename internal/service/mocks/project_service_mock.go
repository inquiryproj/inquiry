// Code generated by mockery v2.19.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	app "github.com/inquiryproj/inquiry/internal/app"
)

// Project is an autogenerated mock type for the Project type
type Project struct {
	mock.Mock
}

// CreateProject provides a mock function with given fields: ctx, createProjectRequest
func (_m *Project) CreateProject(ctx context.Context, createProjectRequest *app.CreateProjectRequest) (*app.Project, error) {
	ret := _m.Called(ctx, createProjectRequest)

	var r0 *app.Project
	if rf, ok := ret.Get(0).(func(context.Context, *app.CreateProjectRequest) *app.Project); ok {
		r0 = rf(ctx, createProjectRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*app.Project)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *app.CreateProjectRequest) error); ok {
		r1 = rf(ctx, createProjectRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListProjects provides a mock function with given fields: ctx, getProjectsRequest
func (_m *Project) ListProjects(ctx context.Context, getProjectsRequest *app.ListProjectsRequest) ([]*app.Project, error) {
	ret := _m.Called(ctx, getProjectsRequest)

	var r0 []*app.Project
	if rf, ok := ret.Get(0).(func(context.Context, *app.ListProjectsRequest) []*app.Project); ok {
		r0 = rf(ctx, getProjectsRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*app.Project)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *app.ListProjectsRequest) error); ok {
		r1 = rf(ctx, getProjectsRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewProject interface {
	mock.TestingT
	Cleanup(func())
}

// NewProject creates a new instance of Project. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewProject(t mockConstructorTestingTNewProject) *Project {
	mock := &Project{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
