package handlers

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/inquiryproj/inquiry/internal/app"
	httpMocks "github.com/inquiryproj/inquiry/internal/http/mocks"
	serviceMocks "github.com/inquiryproj/inquiry/internal/service/mocks"
)

func TestRunProjects(t *testing.T) {
	projectID := uuid.New()

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
				echoMockContext.On("JSON", http.StatusOK, mock.Anything).Return(nil)
				runnerServiceMock.On("RunProject", mock.Anything, &app.RunProjectRequest{
					ProjectID: projectID,
				}).Return(&app.ProjectRunOutput{}, nil)
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
