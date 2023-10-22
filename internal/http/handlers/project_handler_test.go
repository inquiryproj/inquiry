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

func TestCreateProject(t *testing.T) {
	projectID := uuid.New()

	tests := []struct {
		name          string
		setupMocks    func(echoMockContext *httpMocks.Context, runnerServiceMock *serviceMocks.Project)
		expectErr     bool
		errStatusCode int
	}{
		{
			name: "success",
			setupMocks: func(echoMockContext *httpMocks.Context, projectServiceMock *serviceMocks.Project) {
				echoMockContext.On("Request").Return(httpRequestForStruct(t, httpInternal.Project{
					Name: "test",
				}))
				echoMockContext.On("JSON", http.StatusCreated, mock.Anything).Run(func(args mock.Arguments) {
					assert.Equal(t, httpInternal.Project{
						ID:   projectID,
						Name: "test",
					}, args.Get(1))
				}).Return(nil)
				projectServiceMock.On("CreateProject", mock.Anything, &app.CreateProjectRequest{
					Name: "test",
				}).Return(&app.Project{
					ID:   projectID,
					Name: "test",
				}, nil)
			},
		},
		{
			name: "unable to create project, internal",
			setupMocks: func(echoMockContext *httpMocks.Context, projectServiceMock *serviceMocks.Project) {
				echoMockContext.On("Request").Return(httpRequestForStruct(t, httpInternal.Project{
					Name: "test",
				}))
				projectServiceMock.On("CreateProject", mock.Anything, &app.CreateProjectRequest{
					Name: "test",
				}).Return(nil, assert.AnError)
			},
			expectErr:     true,
			errStatusCode: http.StatusInternalServerError,
		},
		{
			name: "unable to create project, already exists",
			setupMocks: func(echoMockContext *httpMocks.Context, projectServiceMock *serviceMocks.Project) {
				echoMockContext.On("Request").Return(httpRequestForStruct(t, httpInternal.Project{
					Name: "test",
				}))
				projectServiceMock.On("CreateProject", mock.Anything, &app.CreateProjectRequest{
					Name: "test",
				}).Return(nil, app.ErrProjectAlreadyExists)
			},
			expectErr:     true,
			errStatusCode: http.StatusConflict,
		},
		{
			name: "unable to create project, already exists",
			setupMocks: func(echoMockContext *httpMocks.Context, projectServiceMock *serviceMocks.Project) {
				echoMockContext.On("Request").Return(httpRequestForStruct(t, httpInternal.Project{}))
			},
			expectErr:     true,
			errStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			echoMockContext := httpMocks.NewContext(t)
			projectServiceMock := serviceMocks.NewProject(t)

			tt.setupMocks(echoMockContext, projectServiceMock)

			runHandler := newProjectHandler(projectServiceMock)
			err := runHandler.CreateProject(echoMockContext)
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

func TestListProjects(t *testing.T) {
	projectID := uuid.New()

	tests := []struct {
		name              string
		setupMocks        func(echoMockContext *httpMocks.Context, runnerServiceMock *serviceMocks.Project)
		listProjectParams httpInternal.ListProjectsParams
		expectErr         bool
		errStatusCode     int
	}{
		{
			name: "success",
			setupMocks: func(echoMockContext *httpMocks.Context, projectServiceMock *serviceMocks.Project) {
				echoMockContext.On("Request").Return(&http.Request{})
				echoMockContext.On("JSON", http.StatusOK, mock.Anything).Run(func(args mock.Arguments) {
					assert.Equal(t, []httpInternal.Project{
						{
							ID:   projectID,
							Name: "test",
						},
					}, args.Get(1))
				}).Return(nil)
				projectServiceMock.On("ListProjects", mock.Anything, &app.ListProjectsRequest{
					Limit:  100,
					Offset: 0,
				}).Return([]*app.Project{
					{
						ID:   projectID,
						Name: "test",
					},
				}, nil)
			},
		},
		{
			name: "pagination success",
			setupMocks: func(echoMockContext *httpMocks.Context, projectServiceMock *serviceMocks.Project) {
				echoMockContext.On("Request").Return(&http.Request{})
				echoMockContext.On("JSON", http.StatusOK, mock.Anything).Return(nil)
				projectServiceMock.On("ListProjects", mock.Anything, &app.ListProjectsRequest{
					Limit:  42,
					Offset: 42,
				}).Return([]*app.Project{}, nil)
			},
			listProjectParams: httpInternal.ListProjectsParams{
				Limit:  newInt(42),
				Offset: newInt(42),
			},
		},
		{
			name: "unable to list projects, internal",
			setupMocks: func(echoMockContext *httpMocks.Context, projectServiceMock *serviceMocks.Project) {
				echoMockContext.On("Request").Return(&http.Request{})
				projectServiceMock.On("ListProjects", mock.Anything, &app.ListProjectsRequest{
					Limit:  100,
					Offset: 0,
				}).Return(nil, assert.AnError)
			},
			expectErr:     true,
			errStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			echoMockContext := httpMocks.NewContext(t)
			projectServiceMock := serviceMocks.NewProject(t)

			tt.setupMocks(echoMockContext, projectServiceMock)

			runHandler := newProjectHandler(projectServiceMock)
			err := runHandler.ListProjects(echoMockContext, tt.listProjectParams)
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
