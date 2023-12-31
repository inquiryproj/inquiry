package runs

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"

	"github.com/google/uuid"

	"github.com/inquiryproj/inquiry/internal/events"
	"github.com/inquiryproj/inquiry/internal/executor"
	"github.com/inquiryproj/inquiry/internal/executor/http"
	"github.com/inquiryproj/inquiry/internal/repository"
	"github.com/inquiryproj/inquiry/internal/repository/domain"
)

// Processor processes runs.
type Processor interface {
	// Process processes a run.
	Process(runID uuid.UUID) (uuid.UUID, error)
}

type processor struct {
	completionsProducer events.Producer[uuid.UUID]

	scenarioRepository repository.Scenario
	runRepository      repository.Run

	logger *slog.Logger
}

// NewProcessor creates a new run processor.
func NewProcessor(completionsProducer events.Producer[uuid.UUID], scenarioRepository repository.Scenario, runRepository repository.Run) Processor {
	return &processor{
		completionsProducer: completionsProducer,

		scenarioRepository: scenarioRepository,
		runRepository:      runRepository,

		logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})),
	}
}

// Process processes a run for a given project ID.
func (p *processor) Process(runID uuid.UUID) (uuid.UUID, error) {
	ctx := context.Background()
	run, err := p.runRepository.Update(ctx, &domain.UpdateRunRequest{
		ID:    runID,
		State: domain.RunStateRunning,
	})
	if err != nil {
		return runID, err
	}

	p.logger.Info("processing project", slog.String("project_id", run.ProjectID.String()), slog.String("run_id", runID.String()))

	scenarioResults, err := p.processProject(ctx, run.ProjectID)
	if err != nil {
		p.logger.Error("project failed", slog.String("project_id", run.ProjectID.String()), slog.String("run_id", runID.String()), slog.String("error", err.Error()))

		_, updateErr := p.runRepository.Update(ctx, &domain.UpdateRunRequest{
			ID:           runID,
			State:        domain.RunStateFailure,
			ErrorMessage: err.Error(),
		})
		if updateErr != nil {
			return runID, fmt.Errorf("%w %w", err, updateErr)
		}
		return runID, err
	}
	p.logger.Info("project processed", slog.String("project_id", run.ProjectID.String()), slog.String("run_id", runID.String()))

	success := true
	for _, scenarioResult := range scenarioResults {
		success = success && scenarioResult.Success
	}

	_, err = p.runRepository.Update(ctx, &domain.UpdateRunRequest{
		ID:                 runID,
		State:              domain.RunStateCompleted,
		Success:            success,
		ScenarioRunDetails: executeResultsToScenarioRunDetails(scenarioResults),
	})
	if err != nil {
		return runID, err
	}

	err = p.completionsProducer.Produce(ctx, runID)

	return runID, err
}

func executeResultsToScenarioRunDetails(executeResults []*http.ExecuteResult) []*domain.ScenarioRunDetails {
	scenarioRunDetails := []*domain.ScenarioRunDetails{}
	for _, executeResult := range executeResults {
		scenarioRunDetails = append(scenarioRunDetails, executeResultToScenarioRunDetails(executeResult))
	}
	return scenarioRunDetails
}

func executeResultToScenarioRunDetails(executeResult *http.ExecuteResult) *domain.ScenarioRunDetails {
	return &domain.ScenarioRunDetails{
		Name:       executeResult.Name,
		Duration:   executeResult.TotalExecutionTime,
		Assertions: executeResult.TotalAssertions,
		Steps:      executeStepResultsToStepRunDetails(executeResult.StepResults),
		Success:    executeResult.Success,
	}
}

func executeStepResultsToStepRunDetails(executeStepResult []*http.ExecuteStepResult) []*domain.StepRunDetails {
	stepRunDetails := []*domain.StepRunDetails{}
	for _, stepResult := range executeStepResult {
		stepRunDetails = append(stepRunDetails, executeStepResultToStepRunDetails(stepResult))
	}
	return stepRunDetails
}

func executeStepResultToStepRunDetails(executeStepResult *http.ExecuteStepResult) *domain.StepRunDetails {
	return &domain.StepRunDetails{
		Name:            executeStepResult.Name,
		Assertions:      executeStepResult.Assertions,
		URL:             executeStepResult.URL,
		RequestDuration: executeStepResult.RequestDuration,
		Duration:        executeStepResult.Duration,
		Retries:         executeStepResult.Retries,
		Success:         executeStepResult.Success,
	}
}

func (p *processor) processProject(ctx context.Context, projectID uuid.UUID) ([]*http.ExecuteResult, error) {
	scenarios, err := p.scenarioRepository.GetForProject(ctx, &domain.GetScenariosForProjectRequest{
		ProjectID: projectID,
	})
	if err != nil {
		return nil, err
	}
	scenarioResults := []*http.ExecuteResult{}
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
		executeResult, err := runExecutor.Play()
		if err != nil {
			return nil, err
		}
		scenarioResults = append(scenarioResults, executeResult)
	}

	return scenarioResults, nil
}
