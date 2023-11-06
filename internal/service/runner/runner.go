// Package runner implements the runner service.
package runner

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"

	"github.com/inquiryproj/inquiry/internal/app"
	"github.com/inquiryproj/inquiry/internal/events"
	"github.com/inquiryproj/inquiry/internal/repository"
	"github.com/inquiryproj/inquiry/internal/repository/domain"
	serviceOptions "github.com/inquiryproj/inquiry/internal/service/options"
)

// Runner is the runner service.
type Runner struct {
	projectRepository  repository.Project
	scenarioRepository repository.Scenario
	runRepository      repository.Run
	runsProducer       events.Producer[uuid.UUID]

	logger *slog.Logger
}

// NewService initialises the runner service.
func NewService(
	projectRepository repository.Project,
	scenarioRepository repository.Scenario,
	runRepository repository.Run,
	runsProducer events.Producer[uuid.UUID],
	opts ...serviceOptions.Opts,
) *Runner {
	options := serviceOptions.DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return &Runner{
		projectRepository:  projectRepository,
		scenarioRepository: scenarioRepository,
		runRepository:      runRepository,
		logger:             options.Logger,
		runsProducer:       runsProducer,
	}
}

// RunProject runs all scenarios for a given project.
func (s *Runner) RunProject(ctx context.Context, runProjectRequest *app.RunProjectRequest) (*app.ProjectRunOutput, error) {
	return s.runProjectForID(ctx, runProjectRequest.ProjectID)
}

// RunProjectByName runs all scenarios for a given project with a given name.
func (s *Runner) RunProjectByName(ctx context.Context, run *app.RunProjectByNameRequest) (*app.ProjectRunOutput, error) {
	project, err := s.projectRepository.GetByName(ctx, run.ProjectName)
	if errors.Is(err, domain.ErrProjectNotFound) {
		return nil, app.ErrProjectNotFound
	} else if err != nil {
		s.logger.Error("failed to get project by name", slog.String("error", err.Error()))
		return nil, err
	}

	return s.runProjectForID(ctx, project.ID)
}

func (s *Runner) runProjectForID(ctx context.Context, projectID uuid.UUID) (*app.ProjectRunOutput, error) {
	run, err := s.runRepository.Create(ctx, &domain.CreateRunRequest{
		ProjectID: projectID,
	})
	if err != nil {
		s.logger.Error("failed to create run", slog.String("error", err.Error()))
		return nil, err
	}
	err = s.runsProducer.Produce(ctx, run.ID)
	if err != nil {
		s.logger.Error("failed to produce run for project", slog.String("error", err.Error()))
		return nil, err
	}
	return &app.ProjectRunOutput{
		ID:        run.ID,
		ProjectID: run.ProjectID,
		State:     app.RunState(run.State),
		Success:   false,
	}, nil
}

// ListRunsForProject returns runs for a given project, limit and offset.
func (s *Runner) ListRunsForProject(ctx context.Context, listRunsForProjectRequest *app.ListRunsForProjectRequest) (*app.ListRunsForProjectResponse, error) {
	res, err := s.runRepository.ListForProject(ctx, &domain.ListRunsForProjectRequest{
		ProjectID: listRunsForProjectRequest.ProjectID,
		Limit:     listRunsForProjectRequest.Limit,
		Offset:    listRunsForProjectRequest.Offset,
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
	return &app.ListRunsForProjectResponse{
		Runs: projectRunOutputs,
	}, nil
}

func scenarioRunDetailsToAppScenarioRunDetails(scenario []*domain.ScenarioRunDetails) []*app.ScenarioRunDetails {
	result := []*app.ScenarioRunDetails{}
	for _, detail := range scenario {
		result = append(result, &app.ScenarioRunDetails{
			Name:       detail.Name,
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
