// Package run implements the run processor.
package run

import (
	"bytes"
	"context"
	"encoding/base64"
	"log/slog"
	"os"

	"github.com/google/uuid"

	"github.com/inquiryproj/inquiry/internal/app"
	"github.com/inquiryproj/inquiry/internal/executor"
	"github.com/inquiryproj/inquiry/internal/repository"
	"github.com/inquiryproj/inquiry/internal/repository/domain"
)

// Processor processes runs.
type Processor interface {
	// Process processes a run.
	Process(runID uuid.UUID) (*app.ProjectRunOutput, error)
}

type processor struct {
	scenarioRepository repository.Scenario
	runRepository      repository.Run

	logger *slog.Logger
}

// NewProcessor creates a new run processor.
func NewProcessor(scenarioRepository repository.Scenario, runRepository repository.Run) Processor {
	return &processor{
		scenarioRepository: scenarioRepository,
		runRepository:      runRepository,

		logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})),
	}
}

// Process processes a run for a given project ID.
func (p *processor) Process(runID uuid.UUID) (*app.ProjectRunOutput, error) {
	run, err := p.runRepository.UpdateRun(context.Background(), &domain.UpdateRunRequest{
		ID:    runID,
		State: domain.RunStateRunning,
	})
	if err != nil {
		return nil, err
	}

	p.logger.Info("processing project", slog.String("project_id", run.ProjectID.String()))

	_, err = p.processProject(run.ProjectID)
	if err != nil {
		_, err = p.runRepository.UpdateRun(context.Background(), &domain.UpdateRunRequest{
			ID:           runID,
			State:        domain.RunStateFailure,
			ErrorMessage: err.Error(),
		})
		if err != nil {
			return nil, err
		}
		return nil, err
	}
	p.logger.Info("project processed", slog.String("project_id", run.ProjectID.String()))

	_, err = p.runRepository.UpdateRun(context.Background(), &domain.UpdateRunRequest{
		ID:    runID,
		State: domain.RunStateSuccess,
	})

	return nil, err
}

func (p *processor) processProject(projectID uuid.UUID) (*app.ProjectRunOutput, error) {
	scenarios, err := p.scenarioRepository.GetForProject(context.Background(), &domain.GetScenariosForProjectRequest{
		ProjectID: projectID,
	})
	if err != nil {
		return nil, err
	}

	for _, scenario := range scenarios {
		p.logger.Info("processing scenario", slog.String("scenario_id", scenario.ID.String()))
		b, err := base64.StdEncoding.DecodeString(scenario.Spec)
		if err != nil {
			return nil, err
		}
		runExecutor, err := executor.New(scenario.Name,
			executor.WithReader(bytes.NewBuffer(b)),
			executor.WithLogger(p.logger))
		if err != nil {
			return nil, err
		}
		err = runExecutor.Play()
		if err != nil {
			return nil, err
		}
	}

	return &app.ProjectRunOutput{}, nil
}
