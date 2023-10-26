package http

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	httpMocks "github.com/inquiryproj/inquiry/internal/http/mocks"
)

func TestAPIKeyMiddleware(t *testing.T) {
	userID := uuid.New()
	tests := []struct {
		name        string
		setupMocks  func(apiKeyRepositoryMock *httpMocks.APIKeyRepository, contextMock *httpMocks.Context)
		next        echo.HandlerFunc
		validateErr func(t *testing.T, err error)
	}{
		{
			name: "valid key",
			setupMocks: func(apiKeyRepositoryMock *httpMocks.APIKeyRepository, contextMock *httpMocks.Context) {
				apiKeyRepositoryMock.On("Validate", mock.Anything, "foo").Return(userID, nil)
				req := requestWithHeaders(map[string]string{
					"Authorization": "foo",
				})
				contextMock.On("Request").Return(req)
				contextMock.On("Set", userIDContextKey, userID)
			},
			next: func(c echo.Context) error {
				return nil
			},
			validateErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "invalid key",
			setupMocks: func(apiKeyRepositoryMock *httpMocks.APIKeyRepository, contextMock *httpMocks.Context) {
				apiKeyRepositoryMock.On("Validate", mock.Anything, "foo").Return(nil, assert.AnError)
				req := requestWithHeaders(map[string]string{
					"Authorization": "foo",
				})
				contextMock.On("Request").Return(req)
			},
			next: func(c echo.Context) error {
				return nil
			},
			validateErr: func(t *testing.T, err error) {
				httpError := &echo.HTTPError{}
				assert.ErrorAs(t, err, &httpError)
				assert.Equal(t, http.StatusUnauthorized, httpError.Code)
				assert.Equal(t, "you are not authorized to make this request", httpError.Message)
			},
		},
		{
			name: "missing key",
			setupMocks: func(apiKeyRepositoryMock *httpMocks.APIKeyRepository, contextMock *httpMocks.Context) {
				apiKeyRepositoryMock.On("Validate", mock.Anything, "").Return(nil, assert.AnError)
				req := requestWithHeaders(map[string]string{})
				contextMock.On("Request").Return(req)
			},
			next: func(c echo.Context) error {
				return nil
			},
			validateErr: func(t *testing.T, err error) {
				httpError := &echo.HTTPError{}
				assert.ErrorAs(t, err, &httpError)
				assert.Equal(t, http.StatusUnauthorized, httpError.Code)
				assert.Equal(t, "you are not authorized to make this request", httpError.Message)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiKeyRepositoryMock := httpMocks.NewAPIKeyRepository(t)
			contextMock := httpMocks.NewContext(t)
			tt.setupMocks(apiKeyRepositoryMock, contextMock)

			apiKeyMiddleware := APIKeyMiddleware(apiKeyRepositoryMock)
			err := apiKeyMiddleware(tt.next)(contextMock)
			tt.validateErr(t, err)
		})
	}
}

func requestWithHeaders(headers map[string]string) *http.Request {
	req := &http.Request{
		Header: map[string][]string{},
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return req
}

func TestDenyAll(t *testing.T) {
	denyAllRepo := &apiKeyDenyAll{}
	userID, err := denyAllRepo.Validate(context.Background(), "")
	assert.Equal(t, uuid.Nil, userID)
	assert.Equal(t, ErrDenyAllKey, err)
}
