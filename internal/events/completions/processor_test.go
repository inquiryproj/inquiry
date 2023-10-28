package completions

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/inquiryproj/inquiry/internal/notifiers"
	notifiersDomain "github.com/inquiryproj/inquiry/internal/notifiers/domain"
	notifierMocks "github.com/inquiryproj/inquiry/internal/notifiers/mocks"
	"github.com/inquiryproj/inquiry/internal/repository/domain"
	repositoryMocks "github.com/inquiryproj/inquiry/internal/repository/mocks"
)

type mockWrapper struct {
	notifierMock          *notifierMocks.Notifier
	runRepositoryMock     *repositoryMocks.Run
	projectRepositoryMock *repositoryMocks.Project
}

func newMockWrapper(t *testing.T) *mockWrapper {
	return &mockWrapper{
		notifierMock:          notifierMocks.NewNotifier(t),
		runRepositoryMock:     repositoryMocks.NewRun(t),
		projectRepositoryMock: repositoryMocks.NewProject(t),
	}
}

func TestProcess(t *testing.T) {
	runID := uuid.New()
	projectID := uuid.New()
	tests := []struct {
		name       string
		setupMocks func(mockWrapper *mockWrapper)
		expectErr  bool
	}{
		{
			name: "success",
			setupMocks: func(mockWrapper *mockWrapper) {
				mockWrapper.runRepositoryMock.On("GetRun", mock.Anything, runID).Return(&domain.Run{
					ProjectID: projectID,
					Success:   true,
					ScenarioRunDetails: []*domain.ScenarioRunDetails{
						{
							Duration: time.Second * 42,
							Success:  true,
							Steps: []*domain.StepRunDetails{
								{
									Success: true,
								},
							},
						},
					},
				}, nil)

				mockWrapper.projectRepositoryMock.On("GetByID", mock.Anything, projectID).Return(&domain.Project{
					ID:   projectID,
					Name: "test",
				}, nil)

				mockWrapper.notifierMock.On("SendCompletion", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					projectRun, ok := args.Get(1).(*notifiersDomain.ProjectRun)
					assert.True(t, ok)
					assert.Equal(t, "test", projectRun.Name)
					assert.Equal(t, true, projectRun.Success)
					assert.Equal(t, 1, len(projectRun.ScenarioRuns))
					assert.Equal(t, 1, projectRun.ScenarioRuns[0].Assertions)
					assert.Equal(t, 1, projectRun.ScenarioRuns[0].SuccessfulAssertions)
					assert.Equal(t, true, projectRun.ScenarioRuns[0].Success)
				}).Return(nil)
			},
		},
		{
			name: "unable to send notifier",
			setupMocks: func(mockWrapper *mockWrapper) {
				mockWrapper.runRepositoryMock.On("GetRun", mock.Anything, runID).Return(&domain.Run{
					ProjectID: projectID,
				}, nil)

				mockWrapper.projectRepositoryMock.On("GetByID", mock.Anything, projectID).Return(&domain.Project{
					ID: projectID,
				}, nil)

				mockWrapper.notifierMock.On("SendCompletion", mock.Anything, mock.Anything).Return(assert.AnError)
			},
			expectErr: true,
		},
		{
			name: "unable to get project",
			setupMocks: func(mockWrapper *mockWrapper) {
				mockWrapper.runRepositoryMock.On("GetRun", mock.Anything, runID).Return(&domain.Run{
					ProjectID: projectID,
				}, nil)

				mockWrapper.projectRepositoryMock.On("GetByID", mock.Anything, projectID).Return(nil, assert.AnError)
			},
			expectErr: true,
		},
		{
			name: "unable to get run",
			setupMocks: func(mockWrapper *mockWrapper) {
				mockWrapper.runRepositoryMock.On("GetRun", mock.Anything, runID).Return(nil, assert.AnError)
			},
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWrapper := newMockWrapper(t)
			tt.setupMocks(mockWrapper)
			p := NewProcessor([]notifiers.Notifier{
				mockWrapper.notifierMock,
			}, mockWrapper.runRepositoryMock, mockWrapper.projectRepositoryMock)
			outID, err := p.Process(runID)
			assert.Equal(t, runID, outID)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
