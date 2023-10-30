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
	"github.com/inquiryproj/inquiry/internal/http/api"
	httpMocks "github.com/inquiryproj/inquiry/internal/http/mocks"
	serviceMocks "github.com/inquiryproj/inquiry/internal/service/mocks"
)

func TestRunProject(t *testing.T) {
	projectID := uuid.New()
	runID := uuid.New()
	projectName := "default"
	tests := []struct {
		name          string
		setupMocks    func(echoMockContext *httpMocks.Context, runnerServiceMock *serviceMocks.Runner)
		expectErr     bool
		errStatusCode int
	}{
		{
			name: "success by project id",
			setupMocks: func(echoMockContext *httpMocks.Context, runnerServiceMock *serviceMocks.Runner) {
				echoMockContext.On("Request").Return(httpRequestForStruct(t, api.RunProjectJSONRequestBody{
					ProjectID: &projectID,
				}))
				echoMockContext.On("JSON", http.StatusOK, mock.Anything).Run(func(args mock.Arguments) {
					assert.Equal(t, api.ProjectRunOutput{
						ID:        runID,
						ProjectID: projectID,
						Success:   false,
						State:     api.Pending,
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
			name: "unable to run project by id",
			setupMocks: func(echoMockContext *httpMocks.Context, runnerServiceMock *serviceMocks.Runner) {
				echoMockContext.On("Request").Return(httpRequestForStruct(t, api.RunProjectJSONRequestBody{
					ProjectID: &projectID,
				}))
				runnerServiceMock.On("RunProject", mock.Anything, &app.RunProjectRequest{
					ProjectID: projectID,
				}).Return(nil, assert.AnError)
			},
			expectErr:     true,
			errStatusCode: http.StatusInternalServerError,
		},
		{
			name: "success",
			setupMocks: func(echoMockContext *httpMocks.Context, runnerServiceMock *serviceMocks.Runner) {
				echoMockContext.On("Request").Return(httpRequestForStruct(t, api.RunProjectJSONRequestBody{
					ProjectName: &projectName,
				}))
				echoMockContext.On("JSON", http.StatusOK, mock.Anything).Run(func(args mock.Arguments) {
					assert.Equal(t, api.ProjectRunOutput{
						ID:        runID,
						ProjectID: projectID,
						Success:   false,
						State:     api.Pending,
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
			name: "unable to run project by name",
			setupMocks: func(echoMockContext *httpMocks.Context, runnerServiceMock *serviceMocks.Runner) {
				echoMockContext.On("Request").Return(httpRequestForStruct(t, api.RunProjectJSONRequestBody{
					ProjectName: &projectName,
				}))
				runnerServiceMock.On("RunProjectByName", mock.Anything, &app.RunProjectByNameRequest{
					ProjectName: "default",
				}).Return(nil, assert.AnError)
			},
			expectErr:     true,
			errStatusCode: http.StatusInternalServerError,
		},
		{
			name: "incorrect payload",
			setupMocks: func(echoMockContext *httpMocks.Context, runnerServiceMock *serviceMocks.Runner) {
				echoMockContext.On("Request").Return(httpRequestForStruct(t, api.RunProjectJSONRequestBody{}))
			},
			expectErr:     true,
			errStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			echoMockContext := httpMocks.NewContext(t)
			runnerServiceMock := serviceMocks.NewRunner(t)

			tt.setupMocks(echoMockContext, runnerServiceMock)

			runHandler := newRunHandler(runnerServiceMock)
			err := runHandler.RunProject(echoMockContext)
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

func TestListRunsForProject(t *testing.T) {
	projectID := uuid.New()
	runID := uuid.New()

	tests := []struct {
		name                     string
		ListRunsForProjectParams api.ListRunsForProjectParams
		setupMocks               func(echoMockContext *httpMocks.Context, runnerServiceMock *serviceMocks.Runner)
		expectErr                bool
		errStatusCode            int
	}{
		{
			name: "success",
			setupMocks: func(echoMockContext *httpMocks.Context, runnerServiceMock *serviceMocks.Runner) {
				echoMockContext.On("Request").Return(&http.Request{})
				echoMockContext.On("JSON", http.StatusOK, mock.Anything).Run(func(args mock.Arguments) {
					assert.Equal(t, []api.ProjectRunOutput{
						{
							ID:        runID,
							ProjectID: projectID,
							Success:   true,
							State:     api.Completed,
							ScenarioRunDetails: []api.ScenarioRunDetails{
								{
									Name:         "Test Scenario 1",
									DurationInMs: 1000,
									Assertions:   42,
									Steps: []api.StepRunDetails{
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
			ListRunsForProjectParams: api.ListRunsForProjectParams{
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
			err := runHandler.ListRunsForProject(echoMockContext, projectID, tt.ListRunsForProjectParams)
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
		State:     app.RunStateCompleted,
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
