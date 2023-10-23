package runner

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/inquiryproj/inquiry/internal/app"
	eventMocks "github.com/inquiryproj/inquiry/internal/events/runs/mocks"
	"github.com/inquiryproj/inquiry/internal/repository/domain"
	repositoryMocks "github.com/inquiryproj/inquiry/internal/repository/mocks"
)

type mockWrapper struct {
	scenarioRepositoryMock *repositoryMocks.Scenario
	runRepositoryMock      *repositoryMocks.Run
	runProducerMock        *eventMocks.Producer
}

func newMockWrapper(t *testing.T) *mockWrapper {
	return &mockWrapper{
		scenarioRepositoryMock: repositoryMocks.NewScenario(t),
		runRepositoryMock:      repositoryMocks.NewRun(t),
		runProducerMock:        eventMocks.NewProducer(t),
	}
}

func TestGetRunsForProject(t *testing.T) {
	projectID := uuid.New()
	runID := uuid.New()
	tests := []struct {
		name                     string
		getRunsForProjectRequest *app.GetRunsForProjectRequest
		setupMocks               func(*mockWrapper)
		validateOutput           func(*testing.T, *app.GetRunsForProjectResponse, error)
	}{
		{
			name: "success",
			setupMocks: func(wrapper *mockWrapper) {
				wrapper.runRepositoryMock.On("GetForProject", mock.Anything,
					&domain.GetRunsForProjectRequest{
						ProjectID: projectID,
						Limit:     10,
						Offset:    0,
					}).
					Return([]*domain.Run{
						{
							ID:        runID,
							ProjectID: projectID,
							State:     domain.RunStateSuccess,
							Success:   false,
							ScenarioRunDetails: []*domain.ScenarioRunDetails{
								testScenarioDetails(),
							},
						},
					}, nil)
			},
			validateOutput: func(t *testing.T, res *app.GetRunsForProjectResponse, err error) {
				assert.NoError(t, err)

				assert.Equal(t, 1, len(res.Runs))
				assert.Equal(t, runID, res.Runs[0].ID)
				assert.Equal(t, projectID, res.Runs[0].ProjectID)
				assert.Equal(t, app.RunStateSuccess, res.Runs[0].State)
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
				wrapper.runRepositoryMock.On("GetForProject", mock.Anything,
					&domain.GetRunsForProjectRequest{
						ProjectID: projectID,
						Limit:     10,
						Offset:    0,
					}).
					Return([]*domain.Run{}, nil)
			},
			validateOutput: func(t *testing.T, res *app.GetRunsForProjectResponse, err error) {
				assert.NoError(t, err)

				assert.Equal(t, 0, len(res.Runs))
			},
		},
		{
			name: "no scenario details",
			setupMocks: func(wrapper *mockWrapper) {
				wrapper.runRepositoryMock.On("GetForProject", mock.Anything,
					&domain.GetRunsForProjectRequest{
						ProjectID: projectID,
						Limit:     10,
						Offset:    0,
					}).
					Return([]*domain.Run{
						{
							ID:        runID,
							ProjectID: projectID,
							State:     domain.RunStateSuccess,
							Success:   false,
						},
					}, nil)
			},
			validateOutput: func(t *testing.T, res *app.GetRunsForProjectResponse, err error) {
				assert.NoError(t, err)

				assert.Equal(t, 1, len(res.Runs))
			},
		},
		{
			name: "unable to get runs",
			setupMocks: func(wrapper *mockWrapper) {
				wrapper.runRepositoryMock.On("GetForProject", mock.Anything,
					&domain.GetRunsForProjectRequest{
						ProjectID: projectID,
						Limit:     10,
						Offset:    0,
					}).
					Return(nil, assert.AnError)
			},
			validateOutput: func(t *testing.T, res *app.GetRunsForProjectResponse, err error) {
				assert.Error(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaultInput := &app.GetRunsForProjectRequest{
				ProjectID: projectID,
				Limit:     10,
				Offset:    0,
			}
			if tt.getRunsForProjectRequest != nil {
				defaultInput = tt.getRunsForProjectRequest
			}
			wrapper := newMockWrapper(t)
			tt.setupMocks(wrapper)
			s := newRunnerService(wrapper)
			res, err := s.GetRunsForProject(context.Background(), defaultInput)
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
				wrapper.runRepositoryMock.On("CreateRun", mock.Anything,
					&domain.CreateRunRequest{
						ProjectID: projectID,
					}).
					Return(&domain.Run{
						ID:        runID,
						ProjectID: projectID,
						State:     domain.RunStatePending,
						Success:   false,
					}, nil)
				wrapper.runProducerMock.On("Produce", runID).Return(nil)
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
				wrapper.runRepositoryMock.On("CreateRun", mock.Anything,
					&domain.CreateRunRequest{
						ProjectID: projectID,
					}).
					Return(&domain.Run{
						ID:        runID,
						ProjectID: projectID,
						State:     domain.RunStatePending,
						Success:   false,
					}, nil)
				wrapper.runProducerMock.On("Produce", runID).Return(assert.AnError)
			},
			validateOutput: func(t *testing.T, res *app.ProjectRunOutput, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "unable to create run",
			setupMocks: func(wrapper *mockWrapper) {
				wrapper.runRepositoryMock.On("CreateRun", mock.Anything,
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

func newRunnerService(mockWrapper *mockWrapper) *Runner {
	return NewService(
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
