package handlers

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/inquiryproj/inquiry/internal/app"
	httpInternal "github.com/inquiryproj/inquiry/internal/http"
	httpMocks "github.com/inquiryproj/inquiry/internal/http/mocks"
	serviceMocks "github.com/inquiryproj/inquiry/internal/service/mocks"
)

func TestRunProject(t *testing.T) {
	projectID := uuid.New()
	runID := uuid.New()

	tests := []struct {
		name          string
		setupMocks    func(echoMockContext *httpMocks.Context, runnerServiceMock *serviceMocks.Runner)
		expectErr     bool
		errStatusCode int
	}{
		{
			name: "success",
			setupMocks: func(echoMockContext *httpMocks.Context, runnerServiceMock *serviceMocks.Runner) {
				echoMockContext.On("Request").Return(&http.Request{})
				echoMockContext.On("JSON", http.StatusOK, mock.Anything).Run(func(args mock.Arguments) {
					assert.Equal(t, httpInternal.ProjectRunOutput{
						ID:        runID,
						ProjectID: projectID,
						Success:   false,
						State:     httpInternal.Pending,
					}, args.Get(1))
				}).Return(nil)
				runnerServiceMock.On("RunProject", mock.Anything, &app.RunProjectRequest{
					ProjectID: projectID,
				}).Return(&app.ProjectRunOutput{
					ID:        runID,
					ProjectID: projectID,
					Success:   false,
					State:     app.RunStatePending,
				}, nil)
			},
		},
		{
			name: "unable to run project",
			setupMocks: func(echoMockContext *httpMocks.Context, runnerServiceMock *serviceMocks.Runner) {
				echoMockContext.On("Request").Return(&http.Request{})
				runnerServiceMock.On("RunProject", mock.Anything, &app.RunProjectRequest{
					ProjectID: projectID,
				}).Return(nil, assert.AnError)
			},
			expectErr:     true,
			errStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			echoMockContext := httpMocks.NewContext(t)
			runnerServiceMock := serviceMocks.NewRunner(t)

			tt.setupMocks(echoMockContext, runnerServiceMock)

			runHandler := newRunHandler(runnerServiceMock)
			err := runHandler.RunProject(echoMockContext, projectID)
			if tt.expectErr {
				assert.Error(t, err)
				httpError := &echo.HTTPError{}
				assert.ErrorAs(t, err, &httpError)
				assert.Equal(t, tt.errStatusCode, httpError.Code)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestRunProjectByName(t *testing.T) {
	projectID := uuid.New()
	runID := uuid.New()

	tests := []struct {
		name          string
		setupMocks    func(echoMockContext *httpMocks.Context, runnerServiceMock *serviceMocks.Runner)
		expectErr     bool
		errStatusCode int
	}{
		{
			name: "success",
			setupMocks: func(echoMockContext *httpMocks.Context, runnerServiceMock *serviceMocks.Runner) {
				echoMockContext.On("Request").Return(&http.Request{})
				echoMockContext.On("JSON", http.StatusOK, mock.Anything).Run(func(args mock.Arguments) {
					assert.Equal(t, httpInternal.ProjectRunOutput{
						ID:        runID,
						ProjectID: projectID,
						Success:   false,
						State:     httpInternal.Pending,
					}, args.Get(1))
				}).Return(nil)
				runnerServiceMock.On("RunProjectByName", mock.Anything, &app.RunProjectByNameRequest{
					ProjectName: "default",
				}).Return(&app.ProjectRunOutput{
					ID:        runID,
					ProjectID: projectID,
					Success:   false,
					State:     app.RunStatePending,
				}, nil)
			},
		},
		{
			name: "unable to run project",
			setupMocks: func(echoMockContext *httpMocks.Context, runnerServiceMock *serviceMocks.Runner) {
				echoMockContext.On("Request").Return(&http.Request{})
				runnerServiceMock.On("RunProjectByName", mock.Anything, &app.RunProjectByNameRequest{
					ProjectName: "default",
				}).Return(nil, assert.AnError)
			},
			expectErr:     true,
			errStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			echoMockContext := httpMocks.NewContext(t)
			runnerServiceMock := serviceMocks.NewRunner(t)

			tt.setupMocks(echoMockContext, runnerServiceMock)

			runHandler := newRunHandler(runnerServiceMock)
			err := runHandler.RunProjectByName(echoMockContext, "default")
			if tt.expectErr {
				assert.Error(t, err)
				httpError := &echo.HTTPError{}
				assert.ErrorAs(t, err, &httpError)
				assert.Equal(t, tt.errStatusCode, httpError.Code)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestGetRunsForProject(t *testing.T) {
	projectID := uuid.New()
	runID := uuid.New()

	tests := []struct {
		name                    string
		getRunsForProjectParams httpInternal.GetRunsForProjectParams
		setupMocks              func(echoMockContext *httpMocks.Context, runnerServiceMock *serviceMocks.Runner)
		expectErr               bool
		errStatusCode           int
	}{
		{
			name: "success",
			setupMocks: func(echoMockContext *httpMocks.Context, runnerServiceMock *serviceMocks.Runner) {
				echoMockContext.On("Request").Return(&http.Request{})
				echoMockContext.On("JSON", http.StatusOK, mock.Anything).Run(func(args mock.Arguments) {
					assert.Equal(t, []httpInternal.ProjectRunOutput{
						{
							ID:        runID,
							ProjectID: projectID,
							Success:   true,
							State:     httpInternal.Success,
							ScenarioRunDetails: []httpInternal.ScenarioRunDetails{
								{
									Name:         "Test Scenario 1",
									DurationInMs: 1000,
									Assertions:   42,
									Steps: []httpInternal.StepRunDetails{
										{
											Name:                "Test Step 1",
											Assertions:          42,
											URL:                 "https://example.com",
											RequestDurationInMs: 1000,
											DurationInMs:        1000,
											Retries:             1,
											Success:             true,
										},
									},
									Success: true,
								},
							},
						},
					}, args.Get(1))
				}).Return(nil)
				runnerServiceMock.On("GetRunsForProject", mock.Anything, &app.GetRunsForProjectRequest{
					Limit:     100,
					Offset:    0,
					ProjectID: projectID,
				}).Return(&app.GetRunsForProjectResponse{
					Runs: []*app.ProjectRunOutput{
						dummyProjectRunOutput(projectID, runID),
					},
				}, nil)
			},
		},
		{
			name: "unable to get runs for project",
			setupMocks: func(echoMockContext *httpMocks.Context, runnerServiceMock *serviceMocks.Runner) {
				echoMockContext.On("Request").Return(&http.Request{})
				runnerServiceMock.On("GetRunsForProject", mock.Anything, &app.GetRunsForProjectRequest{
					Limit:     50,
					Offset:    10,
					ProjectID: projectID,
				}).Return(nil, assert.AnError)
			},
			expectErr:     true,
			errStatusCode: http.StatusInternalServerError,
			getRunsForProjectParams: httpInternal.GetRunsForProjectParams{
				Limit:  newInt(50),
				Offset: newInt(10),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			echoMockContext := httpMocks.NewContext(t)
			runnerServiceMock := serviceMocks.NewRunner(t)

			tt.setupMocks(echoMockContext, runnerServiceMock)

			runHandler := newRunHandler(runnerServiceMock)
			err := runHandler.GetRunsForProject(echoMockContext, projectID, tt.getRunsForProjectParams)
			if tt.expectErr {
				assert.Error(t, err)
				httpError := &echo.HTTPError{}
				assert.ErrorAs(t, err, &httpError)
				assert.Equal(t, tt.errStatusCode, httpError.Code)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func dummyProjectRunOutput(projectID, runID uuid.UUID) *app.ProjectRunOutput {
	return &app.ProjectRunOutput{
		ID:        runID,
		ProjectID: projectID,
		Success:   true,
		State:     app.RunStateSuccess,
		ScenarioRunDetails: []*app.ScenarioRunDetails{
			{
				Name:       "Test Scenario 1",
				Duration:   time.Second,
				Assertions: 42,
				Steps: []*app.StepRunDetails{
					{
						Name:            "Test Step 1",
						Assertions:      42,
						URL:             "https://example.com",
						RequestDuration: time.Second,
						Duration:        time.Second,
						Retries:         1,
						Success:         true,
					},
				},
				Success: true,
			},
		},
	}
}

func newInt(i int) *int {
	return &i
}
