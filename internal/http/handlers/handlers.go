// Package handlers implements the handlers for the API server.
package handlers

import (
	"log/slog"
	"os"

	"github.com/inquiryproj/inquiry/internal/http/api"
	"github.com/inquiryproj/inquiry/internal/service"
)

// validate server interface implementation.
var _ api.ServerInterface = &struct {
	*ProjectHandler
	*ScenarioHandler
	*RunHandler
}{}

// Options represents the options for the handlers.
type Options struct {
	Logger *slog.Logger
}

func defaultOptions() *Options {
	return &Options{
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})),
	}
}

// Opts represents a function that modifies the options.
type Opts func(*Options)

// WithLogger sets the logger.
func WithLogger(logger *slog.Logger) Opts {
	return func(o *Options) {
		o.Logger = logger
	}
}

// HandlerWrapper wraps all handlers.
type HandlerWrapper struct {
	*ProjectHandler
	*ScenarioHandler
	*RunHandler
}

// NewHandlerWrapper initialises all handlers.
func NewHandlerWrapper(
	serviceWrapper service.Wrapper,
	opts ...Opts,
) *HandlerWrapper {
	return &HandlerWrapper{
		ProjectHandler:  newProjectHandler(serviceWrapper, opts...),
		ScenarioHandler: newScenarioHandler(serviceWrapper, opts...),
		RunHandler:      newRunHandler(serviceWrapper, opts...),
	}
}
