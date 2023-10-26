package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const userIDContextKey string = "user_id"

// ErrDenyAllKey is defined by the default middleware denying all
// incoming requests.
var ErrDenyAllKey = fmt.Errorf("api key is not valid")

// APIKeyRepository is the interface for API key repositories.
type APIKeyRepository interface {
	// Validate validates an API key. If the API Key is valid,
	// a user ID is returned. Otherwise an error is returned.
	Validate(ctx context.Context, key string) (uuid.UUID, error)
}

// apiKeyDenyAll default middleware to deny all requests.
type apiKeyDenyAll struct{}

func (*apiKeyDenyAll) Validate(_ context.Context, _ string) (uuid.UUID, error) {
	return uuid.Nil, ErrDenyAllKey
}

// APIKeyMiddleware is the middleware for validating API keys.
func APIKeyMiddleware(apiKey APIKeyRepository) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, err := apiKey.Validate(c.Request().Context(), c.Request().Header.Get("Authorization"))
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "you are not authorized to make this request")
			}
			c.Set(userIDContextKey, userID)
			return next(c)
		}
	}
}
