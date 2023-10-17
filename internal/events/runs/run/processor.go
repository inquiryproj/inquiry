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
)

// Processor processes runs.
type Processor interface {
	// Process processes a run.
	Process(projectID uuid.UUID) (*app.ProjectRunOutput, error)
}

type processor struct {
	scenarioRepository repository.Scenario

	logger *slog.Logger
}

// NewProcessor creates a new run processor.
func NewProcessor(scenarioRepository repository.Scenario) Processor {
	return &processor{
		scenarioRepository: scenarioRepository,
		logger:             slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})),
	}
}

// Process processes a run for a given project ID.
func (p *processor) Process(projectID uuid.UUID) (*app.ProjectRunOutput, error) {
	p.logger.Info("processing project", slog.String("project_id", projectID.String()))
	scenarios, err := p.scenarioRepository.GetForProject(context.Background(), &app.GetScenariosForProjectRequest{
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
	p.logger.Info("project processed", slog.String("project_id", projectID.String()))

	return nil, nil
}
