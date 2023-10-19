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

// GetRunsForProject returns runs for a given project, limit and offset.
func (s *Runner) GetRunsForProject(ctx context.Context, getRunsForProjectRequest *app.GetRunsForProjectRequest) (*app.GetRunsForProjectResponse, error) {
	res, err := s.runRepository.GetForProject(ctx, &domain.GetRunsForProjectRequest{
		ProjectID: getRunsForProjectRequest.ProjectID,
		Limit:     getRunsForProjectRequest.Limit,
		Offset:    getRunsForProjectRequest.Offset,
	})
	if err != nil {
		return nil, err
	}
	projectRunOutputs := make([]*app.ProjectRunOutput, len(res))
	for i, run := range res {
		projectRunOutputs[i] = &app.ProjectRunOutput{
			ID:                 run.ID,
			ProjectID:          run.ProjectID,
			Success:            run.Success,
			State:              app.RunState(run.State),
			ScenarioRunDetails: scenarioRunDetailsToAppScenarioRunDetails(run.ScenarioRunDetails),
		}
	}
	return &app.GetRunsForProjectResponse{
		Runs: projectRunOutputs,
	}, nil
}

func scenarioRunDetailsToAppScenarioRunDetails(scenario []*domain.ScenarioRunDetails) []*app.ScenarioRunDetails {
	result := []*app.ScenarioRunDetails{}
	for _, detail := range scenario {
		result = append(result, &app.ScenarioRunDetails{
			Duration:   detail.Duration,
			Assertions: detail.Assertions,
			Steps:      stepsRunDetailsToAppStepRunDetails(detail.Steps),
			Success:    detail.Success,
		})
	}
	return result
}

func stepsRunDetailsToAppStepRunDetails(steps []*domain.StepRunDetails) []*app.StepRunDetails {
	result := []*app.StepRunDetails{}
	for _, detail := range steps {
		result = append(result, &app.StepRunDetails{
			Name:            detail.Name,
			Assertions:      detail.Assertions,
			URL:             detail.URL,
			RequestDuration: detail.RequestDuration,
			Duration:        detail.Duration,
			Retries:         detail.Retries,
			Success:         detail.Success,
		})
	}
	return result
}
