package sqlite

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/inquiryproj/inquiry/internal/repository/domain"
)

// RunState is the state of a run.
type RunState string

// different run states.
const (
	RunStatePending   RunState = "pending"
	RunStateRunning   RunState = "running"
	RunStateSuccess   RunState = "success"
	RunStateFailure   RunState = "failure"
	RunstateCancelled RunState = "cancelled"
)

// ScenarioDetails is the json model for scenario run details.
type ScenarioDetails struct {
	Duration   time.Duration `json:"duration"`
	Assertions int           `json:"assertions"`
	Steps      []*Step       `json:"steps"`
}

// Step is the json model for scenario step run details.
type Step struct {
	Name             string        `json:"name"`
	Assertions       int           `json:"assertions"`
	FailedAssertions int           `json:"failed_assertions"`
	URL              string        `json:"url"`
	RequestDuration  time.Duration `json:"request_duration"`
	Duration         time.Duration `json:"duration"`
	Retries          int           `json:"retries"`
}

// Run is the sqlite model for runs.
type Run struct {
	BaseModel
	ProjectID    uuid.UUID `gorm:"type:uuid;index:idx_project_id"`
	Success      bool
	State        RunState
	ErrorMessage string
	StepDetails  []byte
}

// RunRepository is the sqlite repository for runs.
type RunRepository struct {
	conn *gorm.DB
}

// GetRun returns a run from sqlite.
func (r *RunRepository) GetRun(ctx context.Context, id uuid.UUID) (*domain.Run, error) {
	run := Run{}
	err := r.conn.Model(&Run{}).WithContext(ctx).Where("id = ?", id).First(&run).Error
	if err != nil {
		return nil, err
	}
	return runToDomainRun(&run)
}

// CreateRun creates a new run in sqlite.
func (r *RunRepository) CreateRun(ctx context.Context, createRunRequest *domain.CreateRunRequest) (*domain.Run, error) {
	run := &Run{
		ProjectID: createRunRequest.ProjectID,
		State:     RunStatePending,
	}
	err := r.conn.WithContext(ctx).Create(run).Error
	if err != nil {
		return nil, err
	}
	return &domain.Run{
		ID:        run.ID,
		ProjectID: run.ProjectID,
	}, nil
}

// UpdateRun updates a run in sqlite.
func (r *RunRepository) UpdateRun(ctx context.Context, updateRunRequest *domain.UpdateRunRequest) (*domain.Run, error) {
	run := Run{}
	err := r.conn.Model(&Run{}).Where("id = ?", updateRunRequest.ID).First(&run).Error
	if err != nil {
		return nil, err
	}
	run.Success = updateRunRequest.Success
	run.State = RunState(updateRunRequest.State)
	run.ErrorMessage = updateRunRequest.ErrorMessage

	scenarioDetails := domainScenarioToScenario(updateRunRequest.ScenarioRunDetails)
	b, err := json.Marshal(scenarioDetails)
	if err != nil {
		return nil, err
	}
	run.StepDetails = b

	err = r.conn.WithContext(ctx).Save(&run).Error
	if err != nil {
		return nil, err
	}

	return runToDomainRun(&run)
}

func domainScenarioToScenario(scenario *domain.ScenarioRunDetails) *ScenarioDetails {
	if scenario == nil {
		return &ScenarioDetails{}
	}
	return &ScenarioDetails{
		Duration:   scenario.Duration,
		Assertions: scenario.Assertions,
		Steps:      domainStepsToSteps(scenario.Steps),
	}
}

func domainStepsToSteps(steps []*domain.StepRunDetails) []*Step {
	result := []*Step{}
	for _, step := range steps {
		result = append(result, &Step{
			Name:             step.Name,
			Assertions:       step.Assertions,
			FailedAssertions: step.FailedAssertions,
			URL:              step.URL,
			RequestDuration:  step.RequestDuration,
			Duration:         step.Duration,
			Retries:          step.Retries,
		})
	}
	return result
}

// GetForProject returns all runs for a given project.
func (r *RunRepository) GetForProject(ctx context.Context, getForProjectRequest *domain.GetRunsForProjectRequest) ([]*domain.Run, error) {
	runs := []*Run{}
	err := r.conn.WithContext(ctx).Where("project_id = ?", getForProjectRequest.ProjectID).Find(&runs).Error
	if err != nil {
		return nil, err
	}
	result := []*domain.Run{}
	for _, run := range runs {
		domainRun, err := runToDomainRun(run)
		if err != nil {
			return nil, err
		}
		result = append(result, domainRun)
	}
	return result, nil
}

func runToDomainRun(run *Run) (*domain.Run, error) {
	scenarioRunDetails, err := scenarioRunDetailsToDomainScenarioRunDetails(run.StepDetails)
	if err != nil {
		return nil, err
	}
	return &domain.Run{
		ID:                 run.ID,
		ProjectID:          run.ProjectID,
		Success:            run.Success,
		State:              domain.RunState(run.State),
		ErrorMessage:       run.ErrorMessage,
		ScenarioRunDetails: scenarioRunDetails,
	}, nil
}

func scenarioRunDetailsToDomainScenarioRunDetails(scenario []byte) (*domain.ScenarioRunDetails, error) {
	details := &ScenarioDetails{}
	err := json.Unmarshal(scenario, details)
	if err != nil {
		return nil, err
	}
	return &domain.ScenarioRunDetails{
		Duration:   details.Duration,
		Assertions: details.Assertions,
		Steps:      stepsRunDetailsToDomainStepRunDetails(details.Steps),
	}, nil
}

func stepsRunDetailsToDomainStepRunDetails(steps []*Step) []*domain.StepRunDetails {
	result := []*domain.StepRunDetails{}
	for _, step := range steps {
		result = append(result, &domain.StepRunDetails{
			Name:             step.Name,
			Assertions:       step.Assertions,
			FailedAssertions: step.FailedAssertions,
			URL:              step.URL,
			RequestDuration:  step.RequestDuration,
			Duration:         step.Duration,
			Retries:          step.Retries,
		})
	}
	return result
}
