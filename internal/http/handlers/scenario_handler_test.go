package handlers

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/inquiryproj/inquiry/internal/app"
	httpInternal "github.com/inquiryproj/inquiry/internal/http"
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
				echoMockContext.On("Request").Return(httpRequestForStruct(t, httpInternal.CreateScenarioJSONRequestBody{
					Name:     "test",
					Spec:     "base64yaml",
					SpecType: "yaml",
				}))
				echoMockContext.On("JSON", http.StatusCreated, mock.Anything).Run(func(args mock.Arguments) {
					assert.Equal(t, httpInternal.Scenario{
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
				echoMockContext.On("Request").Return(httpRequestForStruct(t, httpInternal.Project{
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
				echoMockContext.On("Request").Return(httpRequestForStruct(t, httpInternal.Scenario{
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