package runner

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/inquiryproj/inquiry/internal/app"
	eventMocks "github.com/inquiryproj/inquiry/internal/events/mocks"
	"github.com/inquiryproj/inquiry/internal/repository/domain"
	repositoryMocks "github.com/inquiryproj/inquiry/internal/repository/mocks"
)

type mockWrapper struct {
	scenarioRepositoryMock *repositoryMocks.Scenario
	projectRepositoryMock  *repositoryMocks.Project
	runRepositoryMock      *repositoryMocks.Run
	runProducerMock        *eventMocks.Producer[uuid.UUID]
}

func newMockWrapper(t *testing.T) *mockWrapper {
	return &mockWrapper{
		scenarioRepositoryMock: repositoryMocks.NewScenario(t),
		projectRepositoryMock:  repositoryMocks.NewProject(t),
		runRepositoryMock:      repositoryMocks.NewRun(t),
		runProducerMock:        eventMocks.NewProducer[uuid.UUID](t),
	}
}

func TestListRunsForProject(t *testing.T) {
	projectID := uuid.New()
	runID := uuid.New()
	tests := []struct {
		name                      string
		listRunsForProjectRequest *app.ListRunsForProjectRequest
		setupMocks                func(*mockWrapper)
		validateOutput            func(*testing.T, *app.ListRunsForProjectResponse, error)
	}{
		{
			name: "success",
			setupMocks: func(wrapper *mockWrapper) {
				wrapper.runRepositoryMock.On("ListForProject", mock.Anything,
					&domain.ListRunsForProjectRequest{
						ProjectID: projectID,
						Limit:     10,
						Offset:    0,
					}).
					Return([]*domain.Run{
						{
							ID:        runID,
							ProjectID: projectID,
							State:     domain.RunStateCompleted,
							Success:   false,
							ScenarioRunDetails: []*domain.ScenarioRunDetails{
								testScenarioDetails(),
							},
						},
					}, nil)
			},
			validateOutput: func(t *testing.T, res *app.ListRunsForProjectResponse, err error) {
				assert.NoError(t, err)

				assert.Equal(t, 1, len(res.Runs))
				assert.Equal(t, runID, res.Runs[0].ID)
				assert.Equal(t, projectID, res.Runs[0].ProjectID)
				assert.Equal(t, app.RunStateCompleted, res.Runs[0].State)
				assert.Equal(t, false, res.Runs[0].Success)
				assert.Equal(t, 1, len(res.Runs[0].ScenarioRunDetails))
				assert.Equal(t, "scenario 1", res.Runs[0].ScenarioRunDetails[0].Name)
				assert.Equal(t, 1, res.Runs[0].ScenarioRunDetails[0].Assertions)
				assert.Equal(t, 1, len(res.Runs[0].ScenarioRunDetails[0].Steps))
				assert.Equal(t, "step 1", res.Runs[0].ScenarioRunDetails[0].Steps[0].Name)
				assert.Equal(t, 1, res.Runs[0].ScenarioRunDetails[0].Steps[0].Assertions)
				assert.Equal(t, "http://localhost:8080", res.Runs[0].ScenarioRunDetails[0].Steps[0].URL)
				assert.Equal(t, 1, res.Runs[0].ScenarioRunDetails[0].Steps[0].Retries)
				assert.Equal(t, true, res.Runs[0].ScenarioRunDetails[0].Steps[0].Success)
			},
		},
		{
			name: "no results",
			setupMocks: func(wrapper *mockWrapper) {
				wrapper.runRepositoryMock.On("ListForProject", mock.Anything,
					&domain.ListRunsForProjectRequest{
						ProjectID: projectID,
						Limit:     10,
						Offset:    0,
					}).
					Return([]*domain.Run{}, nil)
			},
			validateOutput: func(t *testing.T, res *app.ListRunsForProjectResponse, err error) {
				assert.NoError(t, err)

				assert.Equal(t, 0, len(res.Runs))
			},
		},
		{
			name: "no scenario details",
			setupMocks: func(wrapper *mockWrapper) {
				wrapper.runRepositoryMock.On("ListForProject", mock.Anything,
					&domain.ListRunsForProjectRequest{
						ProjectID: projectID,
						Limit:     10,
						Offset:    0,
					}).
					Return([]*domain.Run{
						{
							ID:        runID,
							ProjectID: projectID,
							State:     domain.RunStateCompleted,
							Success:   false,
						},
					}, nil)
			},
			validateOutput: func(t *testing.T, res *app.ListRunsForProjectResponse, err error) {
				assert.NoError(t, err)

				assert.Equal(t, 1, len(res.Runs))
			},
		},
		{
			name: "unable to get runs",
			setupMocks: func(wrapper *mockWrapper) {
				wrapper.runRepositoryMock.On("ListForProject", mock.Anything,
					&domain.ListRunsForProjectRequest{
						ProjectID: projectID,
						Limit:     10,
						Offset:    0,
					}).
					Return(nil, assert.AnError)
			},
			validateOutput: func(t *testing.T, res *app.ListRunsForProjectResponse, err error) {
				assert.Error(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaultInput := &app.ListRunsForProjectRequest{
				ProjectID: projectID,
				Limit:     10,
				Offset:    0,
			}
			if tt.listRunsForProjectRequest != nil {
				defaultInput = tt.listRunsForProjectRequest
			}
			wrapper := newMockWrapper(t)
			tt.setupMocks(wrapper)
			s := newRunnerService(wrapper)
			res, err := s.ListRunsForProject(context.Background(), defaultInput)
			tt.validateOutput(t, res, err)
		})
	}
}

func TestRunProject(t *testing.T) {
	projectID := uuid.New()
	runID := uuid.New()
	tests := []struct {
		name              string
		runProjectRequest *app.RunProjectRequest
		setupMocks        func(*mockWrapper)
		validateOutput    func(*testing.T, *app.ProjectRunOutput, error)
	}{
		{
			name: "success",
			setupMocks: func(wrapper *mockWrapper) {
				wrapper.runRepositoryMock.On("Create", mock.Anything,
					&domain.CreateRunRequest{
						ProjectID: projectID,
					}).
					Return(&domain.Run{
						ID:        runID,
						ProjectID: projectID,
						State:     domain.RunStatePending,
						Success:   false,
					}, nil)
				wrapper.runProducerMock.On("Produce", mock.Anything, runID).Return(nil)
			},
			validateOutput: func(t *testing.T, res *app.ProjectRunOutput, err error) {
				assert.NoError(t, err)

				assert.Equal(t, runID, res.ID)
				assert.Equal(t, projectID, res.ProjectID)
				assert.Equal(t, app.RunStatePending, res.State)
				assert.Equal(t, false, res.Success)
			},
		},
		{
			name: "unable to produce",
			setupMocks: func(wrapper *mockWrapper) {
				wrapper.runRepositoryMock.On("Create", mock.Anything,
					&domain.CreateRunRequest{
						ProjectID: projectID,
					}).
					Return(&domain.Run{
						ID:        runID,
						ProjectID: projectID,
						State:     domain.RunStatePending,
						Success:   false,
					}, nil)
				wrapper.runProducerMock.On("Produce", mock.Anything, runID).Return(assert.AnError)
			},
			validateOutput: func(t *testing.T, res *app.ProjectRunOutput, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "unable to create run",
			setupMocks: func(wrapper *mockWrapper) {
				wrapper.runRepositoryMock.On("Create", mock.Anything,
					&domain.CreateRunRequest{
						ProjectID: projectID,
					}).
					Return(nil, assert.AnError)
			},
			validateOutput: func(t *testing.T, res *app.ProjectRunOutput, err error) {
				assert.Error(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaultInput := &app.RunProjectRequest{
				ProjectID: projectID,
			}
			if tt.runProjectRequest != nil {
				defaultInput = tt.runProjectRequest
			}
			wrapper := newMockWrapper(t)
			tt.setupMocks(wrapper)
			s := newRunnerService(wrapper)
			res, err := s.RunProject(context.Background(), defaultInput)
			tt.validateOutput(t, res, err)
		})
	}
}

func TestRunProjectByName(t *testing.T) {
	projectID := uuid.New()
	runID := uuid.New()
	tests := []struct {
		name              string
		runProjectRequest *app.RunProjectByNameRequest
		setupMocks        func(*mockWrapper)
		validateOutput    func(*testing.T, *app.ProjectRunOutput, error)
	}{
		{
			name: "success",
			setupMocks: func(wrapper *mockWrapper) {
				wrapper.projectRepositoryMock.On("GetByName", mock.Anything, "default").
					Return(&domain.Project{
						ID:   projectID,
						Name: "default",
					}, nil)

				wrapper.runRepositoryMock.On("Create", mock.Anything,
					&domain.CreateRunRequest{
						ProjectID: projectID,
					}).
					Return(&domain.Run{
						ID:        runID,
						ProjectID: projectID,
						State:     domain.RunStatePending,
						Success:   false,
					}, nil)
				wrapper.runProducerMock.On("Produce", mock.Anything, runID).Return(nil)
			},
			validateOutput: func(t *testing.T, res *app.ProjectRunOutput, err error) {
				assert.NoError(t, err)

				assert.Equal(t, runID, res.ID)
				assert.Equal(t, projectID, res.ProjectID)
				assert.Equal(t, app.RunStatePending, res.State)
				assert.Equal(t, false, res.Success)
			},
		},
		{
			name: "unable to produce",
			setupMocks: func(wrapper *mockWrapper) {
				wrapper.projectRepositoryMock.On("GetByName", mock.Anything, "default").
					Return(&domain.Project{
						ID:   projectID,
						Name: "default",
					}, nil)

				wrapper.runRepositoryMock.On("Create", mock.Anything,
					&domain.CreateRunRequest{
						ProjectID: projectID,
					}).
					Return(&domain.Run{
						ID:        runID,
						ProjectID: projectID,
						State:     domain.RunStatePending,
						Success:   false,
					}, nil)
				wrapper.runProducerMock.On("Produce", mock.Anything, runID).Return(assert.AnError)
			},
			validateOutput: func(t *testing.T, res *app.ProjectRunOutput, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "unable to create run",
			setupMocks: func(wrapper *mockWrapper) {
				wrapper.projectRepositoryMock.On("GetByName", mock.Anything, "default").
					Return(&domain.Project{
						ID:   projectID,
						Name: "default",
					}, nil)

				wrapper.runRepositoryMock.On("Create", mock.Anything,
					&domain.CreateRunRequest{
						ProjectID: projectID,
					}).
					Return(nil, assert.AnError)
			},
			validateOutput: func(t *testing.T, res *app.ProjectRunOutput, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "project not found",
			setupMocks: func(wrapper *mockWrapper) {
				wrapper.projectRepositoryMock.On("GetByName", mock.Anything, "default").
					Return(nil, domain.ErrProjectNotFound)
			},
			validateOutput: func(t *testing.T, res *app.ProjectRunOutput, err error) {
				assert.Error(t, err)
				assert.ErrorIs(t, err, app.ErrProjectNotFound)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaultInput := &app.RunProjectByNameRequest{
				ProjectName: "default",
			}
			if tt.runProjectRequest != nil {
				defaultInput = tt.runProjectRequest
			}
			wrapper := newMockWrapper(t)
			tt.setupMocks(wrapper)
			s := newRunnerService(wrapper)
			res, err := s.RunProjectByName(context.Background(), defaultInput)
			tt.validateOutput(t, res, err)
		})
	}
}

func newRunnerService(mockWrapper *mockWrapper) *Runner {
	return NewService(
		mockWrapper.projectRepositoryMock,
		mockWrapper.scenarioRepositoryMock,
		mockWrapper.runRepositoryMock,
		mockWrapper.runProducerMock,
	)
}

func testScenarioDetails() *domain.ScenarioRunDetails {
	return &domain.ScenarioRunDetails{
		Name:       "scenario 1",
		Duration:   1,
		Assertions: 1,
		Steps: []*domain.StepRunDetails{
			{
				Name:            "step 1",
				Assertions:      1,
				URL:             "http://localhost:8080",
				RequestDuration: 1,
				Duration:        1,
				Retries:         1,
				Success:         true,
			},
		},
		Success: true,
	}
}
