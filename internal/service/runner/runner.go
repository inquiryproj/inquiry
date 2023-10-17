// Package runner implements the runner service.
package runner

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/inquiryproj/inquiry/internal/app"
	"github.com/inquiryproj/inquiry/internal/events/runs"
	"github.com/inquiryproj/inquiry/internal/repository"
	serviceOptions "github.com/inquiryproj/inquiry/internal/service/options"
)

// Runner is the runner service.
type Runner struct {
	scenarioRepository repository.Scenario
	runsProducer       runs.Producer

	logger *slog.Logger
}

// NewService initialises the runner service.
func NewService(scenarioRepository repository.Scenario, runsProducer runs.Producer, opts ...serviceOptions.Opts) *Runner {
	options := serviceOptions.DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return &Runner{
		scenarioRepository: scenarioRepository,
		logger:             options.Logger,
		runsProducer:       runsProducer,
	}
}

// RunProject runs all scenarios for a given project.
func (s *Runner) RunProject(ctx context.Context, runProjectRequest *app.RunProjectRequest) (*app.ProjectRunOutput, error) {
	// FIXME: Store run created in the database.
	_ = ctx
	err := s.runsProducer.Produce(runProjectRequest.ProjectID)
	if err != nil {
		s.logger.Error("failed to produce run for project", slog.String("error", err.Error()))
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "failed to trigger run for project")
	}
	return &app.ProjectRunOutput{}, nil
}
