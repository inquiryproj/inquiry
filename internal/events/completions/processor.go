package completions

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/inquiryproj/inquiry/internal/notifiers"
	"github.com/inquiryproj/inquiry/internal/repository"
	"github.com/inquiryproj/inquiry/internal/repository/domain"
)

// Notifier is a notifier which can send completions.
type Notifier interface {
	// SendCompletion sends a completion message.
	SendCompletion(ctx context.Context, projectRun notifiers.ProjectRun) error
}

// Processor processes runs.
type Processor interface {
	// Process processes a run.
	Process(runID uuid.UUID) (uuid.UUID, error)
}

type processor struct {
	completionNotifiers []Notifier
	runRepository       repository.Run
	projectRepository   repository.Project
}

// NewProcessor creates a new processor.
func NewProcessor(
	completionNotifiers []Notifier,
	runRepository repository.Run,
	projectRepository repository.Project,
) Processor {
	return &processor{
		completionNotifiers: completionNotifiers,
		runRepository:       runRepository,
		projectRepository:   projectRepository,
	}
}

// Process processes a run.
func (p *processor) Process(runID uuid.UUID) (uuid.UUID, error) {
	ctx := context.Background()
	run, err := p.runRepository.GetRun(ctx, runID)
	if err != nil {
		return runID, err
	}
	project, err := p.projectRepository.GetByID(context.Background(), run.ProjectID)
	if err != nil {
		return runID, err
	}

	projectRun := projectAndRunToProjectRun(project, run)

	for _, notifier := range p.completionNotifiers {
		err := notifier.SendCompletion(ctx, projectRun)
		if err != nil {
			return runID, err
		}
	}

	return runID, nil
}

func projectAndRunToProjectRun(project *domain.Project, run *domain.Run) notifiers.ProjectRun {
	projectRun := notifiers.ProjectRun{
		Name:    project.Name,
		Success: run.Success,
		// FIXME add version to runs
		Version: "latest",
		Time:    run.CreatedAt,
	}
	totalDuration := time.Duration(0)
	for _, s := range run.ScenarioRunDetails {
		totalDuration += s.Duration
		assertions := len(s.Steps)
		successfullAssertions := 0
		for _, step := range s.Steps {
			if step.Success {
				successfullAssertions++
			}
		}
		scenarioRun := &notifiers.ScenarioRunDetails{
			Name:                 s.Name,
			Success:              s.Success,
			Duration:             s.Duration,
			Assertions:           assertions,
			SuccessfulAssertions: successfullAssertions,
		}
		projectRun.ScenarioRuns = append(projectRun.ScenarioRuns, scenarioRun)
	}
	projectRun.Duration = totalDuration
	return projectRun
}
