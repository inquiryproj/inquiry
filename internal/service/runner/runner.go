// Package runner implements the runner service.
package runner

import (
	"context"
	"log/slog"

	"github.com/inquiryproj/inquiry/internal/app"
	"github.com/inquiryproj/inquiry/internal/repository"
	serviceOptions "github.com/inquiryproj/inquiry/internal/service/options"
)

// Runner is the runner service.
type Runner struct {
	scenarioRepository repository.Scenario

	logger *slog.Logger
}

// NewService initialises the runner service.
func NewService(scenarioRepository repository.Scenario, opts ...serviceOptions.Opts) *Runner {
	options := serviceOptions.DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return &Runner{
		scenarioRepository: scenarioRepository,
		logger:             options.Logger,
	}
}

// RunProject runs all scenarios for a given project.
func (s *Runner) RunProject(ctx context.Context, runProjectRequest *app.RunProjectRequest) (*app.ProjectRunOutput, error) {
	// FIXME: implement
	_ = ctx
	_ = runProjectRequest
	return nil, nil
}
