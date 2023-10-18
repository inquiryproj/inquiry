// Package runner implements the runner service.
package runner

import (
	"context"
	"log/slog"

	"github.com/inquiryproj/inquiry/internal/app"
	"github.com/inquiryproj/inquiry/internal/events/runs"
	"github.com/inquiryproj/inquiry/internal/repository"
	"github.com/inquiryproj/inquiry/internal/repository/domain"
	serviceOptions "github.com/inquiryproj/inquiry/internal/service/options"
)

// Runner is the runner service.
type Runner struct {
	scenarioRepository repository.Scenario
	runRepository      repository.Run
	runsProducer       runs.Producer

	logger *slog.Logger
}

// NewService initialises the runner service.
func NewService(
	scenarioRepository repository.Scenario,
	runRepository repository.Run,
	runsProducer runs.Producer,
	opts ...serviceOptions.Opts,
) *Runner {
	options := serviceOptions.DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return &Runner{
		scenarioRepository: scenarioRepository,
		runRepository:      runRepository,
		logger:             options.Logger,
		runsProducer:       runsProducer,
	}
}

// RunProject runs all scenarios for a given project.
func (s *Runner) RunProject(ctx context.Context, runProjectRequest *app.RunProjectRequest) (*app.ProjectRunOutput, error) {
	run, err := s.runRepository.CreateRun(ctx, &domain.CreateRunRequest{
		ProjectID: runProjectRequest.ProjectID,
	})
	if err != nil {
		s.logger.Error("failed to create run", slog.String("error", err.Error()))
		return nil, err
	}
	err = s.runsProducer.Produce(run.ID)
	if err != nil {
		s.logger.Error("failed to produce run for project", slog.String("error", err.Error()))
		return nil, err
	}
	return &app.ProjectRunOutput{}, nil
}
