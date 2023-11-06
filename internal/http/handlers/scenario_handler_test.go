package handlers

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/inquiryproj/inquiry/internal/app"
	"github.com/inquiryproj/inquiry/internal/http/api"
	httpMocks "github.com/inquiryproj/inquiry/internal/http/mocks"
	serviceMocks "github.com/inquiryproj/inquiry/internal/service/mocks"
)

func TestCreateScenario(t *testing.T) {
	projectID := uuid.New()
	scenarioID := uuid.New()

	tests := []struct {
		name          string
		setupMocks    func(echoMockContext *httpMocks.Context, runnerServiceMock *serviceMocks.Scenario)
		expectErr     bool
		errStatusCode int
	}{
		{
			name: "success",
			setupMocks: func(echoMockContext *httpMocks.Context, scenarioServiceMock *serviceMocks.Scenario) {
				echoMockContext.On("Request").Return(httpRequestForStruct(t, api.CreateScenarioJSONRequestBody{
					Name:     "test",
					Spec:     "base64yaml",
					SpecType: "yaml",
				}))
				echoMockContext.On("JSON", http.StatusCreated, mock.Anything).Run(func(args mock.Arguments) {
					assert.Equal(t, api.Scenario{
						ID:        scenarioID,
						ProjectID: projectID,
						Name:      "test",
						Spec:      "base64yaml",
						SpecType:  "yaml",
					}, args.Get(1))
				}).Return(nil)
				scenarioServiceMock.On("CreateScenario", mock.Anything, &app.CreateScenarioRequest{
					Name:      "test",
					ProjectID: projectID,
					SpecType:  app.ScenarioSpecType("yaml"),
					Spec:      "base64yaml",
				}).Return(&app.Scenario{
					ID:        scenarioID,
					ProjectID: projectID,
					Name:      "test",
					Spec:      "base64yaml",
					SpecType:  app.ScenarioSpecType("yaml"),
				}, nil)
			},
		},
		{
			name: "unable to create scenario, internal",
			setupMocks: func(echoMockContext *httpMocks.Context, scenarioServiceMock *serviceMocks.Scenario) {
				echoMockContext.On("Request").Return(httpRequestForStruct(t, api.Project{
					Name: "test",
				}))
				scenarioServiceMock.On("CreateScenario", mock.Anything, mock.Anything).Return(nil, assert.AnError)
			},
			expectErr:     true,
			errStatusCode: http.StatusInternalServerError,
		},
		{
			name: "unable to create scenario, already exists",
			setupMocks: func(echoMockContext *httpMocks.Context, scenarioServiceMock *serviceMocks.Scenario) {
				echoMockContext.On("Request").Return(httpRequestForStruct(t, api.Scenario{
					Name:     "test",
					Spec:     "base64yaml",
					SpecType: "yaml",
				}))
				scenarioServiceMock.On("CreateScenario", mock.Anything, &app.CreateScenarioRequest{
					Name:      "test",
					ProjectID: projectID,
					SpecType:  app.ScenarioSpecType("yaml"),
					Spec:      "base64yaml",
				}).Return(nil, app.ErrScenarioAlreadyExists)
			},
			expectErr:     true,
			errStatusCode: http.StatusConflict,
		},
		{
			name: "unable to create scenario, already exists",
			setupMocks: func(echoMockContext *httpMocks.Context, scenarioServiceMock *serviceMocks.Scenario) {
				echoMockContext.On("Request").Return(httpRequestForStruct(t, api.Scenario{
					Name:     "test",
					Spec:     "base64yaml",
					SpecType: "yaml",
				}))
				scenarioServiceMock.On("CreateScenario", mock.Anything, &app.CreateScenarioRequest{
					Name:      "test",
					ProjectID: projectID,
					SpecType:  app.ScenarioSpecType("yaml"),
					Spec:      "base64yaml",
				}).Return(nil, app.ErrProjectNotFound)
			},
			expectErr:     true,
			errStatusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			echoMockContext := httpMocks.NewContext(t)
			scenarioServiceMock := serviceMocks.NewScenario(t)

			tt.setupMocks(echoMockContext, scenarioServiceMock)

			runHandler := newScenarioHandler(scenarioServiceMock)
			err := runHandler.CreateScenario(echoMockContext, projectID)
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

func TestListScenariosForProject(t *testing.T) {
	projectID := uuid.New()
	scenarioID := uuid.New()
	tests := []struct {
		name          string
		setupMocks    func(echoMockContext *httpMocks.Context, scenarioServiceMock *serviceMocks.Scenario)
		expectErr     bool
		errStatusCode int
	}{
		{
			name: "success",
			setupMocks: func(echoMockContext *httpMocks.Context, scenarioServiceMock *serviceMocks.Scenario) {
				echoMockContext.On("JSON", http.StatusOK, mock.Anything).Run(func(args mock.Arguments) {
					assert.Equal(t, []api.Scenario{
						{
							ID:        scenarioID,
							ProjectID: projectID,
							Name:      "test",
							Spec:      "base64yaml",
							SpecType:  "yaml",
						},
					}, args.Get(1))
				}).Return(nil)
				echoMockContext.On("Request").Return(httpRequestForStruct(t, api.ListScenariosForProjectParams{}))
				scenarioServiceMock.On("ListScenarios", mock.Anything, &app.ListScenariosRequest{
					Limit:     100,
					Offset:    0,
					ProjectID: projectID,
				}).Return([]*app.Scenario{
					{
						ID:        scenarioID,
						ProjectID: projectID,
						Name:      "test",
						Spec:      "base64yaml",
						SpecType:  app.ScenarioSpecType("yaml"),
					},
				}, nil)
			},
		},
		{
			name: "unable to list scenarios, internal",
			setupMocks: func(echoMockContext *httpMocks.Context, scenarioServiceMock *serviceMocks.Scenario) {
				echoMockContext.On("Request").Return(httpRequestForStruct(t, api.ListScenariosForProjectParams{}))
				scenarioServiceMock.On("ListScenarios", mock.Anything, &app.ListScenariosRequest{
					Limit:     100,
					Offset:    0,
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
			scenarioServiceMock := serviceMocks.NewScenario(t)

			tt.setupMocks(echoMockContext, scenarioServiceMock)

			runHandler := newScenarioHandler(scenarioServiceMock)
			err := runHandler.ListScenariosForProject(echoMockContext, projectID, api.ListScenariosForProjectParams{})
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
