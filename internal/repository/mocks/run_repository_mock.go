// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	context "context"

	uuid "github.com/google/uuid"
	mock "github.com/stretchr/testify/mock"

	domain "github.com/inquiryproj/inquiry/internal/repository/domain"
)

// Run is an autogenerated mock type for the Run type
type Run struct {
	mock.Mock
}

// CreateRun provides a mock function with given fields: ctx, createRunRequest
func (_m *Run) CreateRun(ctx context.Context, createRunRequest *domain.CreateRunRequest) (*domain.Run, error) {
	ret := _m.Called(ctx, createRunRequest)

	var r0 *domain.Run
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.CreateRunRequest) (*domain.Run, error)); ok {
		return rf(ctx, createRunRequest)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *domain.CreateRunRequest) *domain.Run); ok {
		r0 = rf(ctx, createRunRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Run)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *domain.CreateRunRequest) error); ok {
		r1 = rf(ctx, createRunRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRun provides a mock function with given fields: ctx, id
func (_m *Run) GetRun(ctx context.Context, id uuid.UUID) (*domain.Run, error) {
	ret := _m.Called(ctx, id)

	var r0 *domain.Run
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*domain.Run, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *domain.Run); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Run)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListForProject provides a mock function with given fields: ctx, listForProject
func (_m *Run) ListForProject(ctx context.Context, listForProject *domain.ListRunsForProjectRequest) ([]*domain.Run, error) {
	ret := _m.Called(ctx, listForProject)

	var r0 []*domain.Run
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.ListRunsForProjectRequest) ([]*domain.Run, error)); ok {
		return rf(ctx, listForProject)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *domain.ListRunsForProjectRequest) []*domain.Run); ok {
		r0 = rf(ctx, listForProject)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Run)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *domain.ListRunsForProjectRequest) error); ok {
		r1 = rf(ctx, listForProject)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateRun provides a mock function with given fields: ctx, updateRunRequest
func (_m *Run) UpdateRun(ctx context.Context, updateRunRequest *domain.UpdateRunRequest) (*domain.Run, error) {
	ret := _m.Called(ctx, updateRunRequest)

	var r0 *domain.Run
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.UpdateRunRequest) (*domain.Run, error)); ok {
		return rf(ctx, updateRunRequest)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *domain.UpdateRunRequest) *domain.Run); ok {
		r0 = rf(ctx, updateRunRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Run)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *domain.UpdateRunRequest) error); ok {
		r1 = rf(ctx, updateRunRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewRun creates a new instance of Run. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRun(t interface {
	mock.TestingT
	Cleanup(func())
}) *Run {
	mock := &Run{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
