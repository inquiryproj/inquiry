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
	Name       string        `json:"name"`
	Duration   time.Duration `json:"duration"`
	Assertions int           `json:"assertions"`
	Steps      []*Step       `json:"steps"`
	Success    bool          `json:"success"`
}

// Step is the json model for scenario step run details.
type Step struct {
	Name            string        `json:"name"`
	Assertions      int           `json:"assertions"`
	URL             string        `json:"url"`
	RequestDuration time.Duration `json:"request_duration"`
	Duration        time.Duration `json:"duration"`
	Retries         int           `json:"retries"`
	Success         bool          `json:"success"`
}

// Run is the sqlite model for runs.
type Run struct {
	BaseModel
	ProjectID       uuid.UUID `gorm:"type:uuid;index:idx_project_id"`
	Success         bool
	State           RunState
	ErrorMessage    string
	ScenarioDetails []byte
}

// RunRepository is the sqlite repository for runs.
type RunRepository struct {
	conn *gorm.DB
}

// GetRun returns a run from sqlite.
func (r *RunRepository) GetRun(ctx context.Context, id uuid.UUID) (*domain.Run, error) {
	run := Run{}
	err := r.conn.WithContext(ctx).Model(&Run{}).Where("id = ?", id).First(&run).Error
	if err != nil {
		return nil, err
	}
	return runToDomainRun(&run)
}

// CreateRun creates a new run in sqlite.
func (r *RunRepository) CreateRun(ctx context.Context, createRunRequest *domain.CreateRunRequest) (*domain.Run, error) {
	run := &Run{
		ProjectID:       createRunRequest.ProjectID,
		State:           RunStatePending,
		ScenarioDetails: []byte(`{}`),
	}
	err := r.conn.WithContext(ctx).Model(&Run{}).Create(run).Error
	if err != nil {
		return nil, err
	}
	return runToDomainRun(run)
}

// UpdateRun updates a run in sqlite.
func (r *RunRepository) UpdateRun(ctx context.Context, updateRunRequest *domain.UpdateRunRequest) (*domain.Run, error) {
	run := Run{}
	err := r.conn.WithContext(ctx).Model(&Run{}).Where("id = ?", updateRunRequest.ID).First(&run).Error
	if err != nil {
		return nil, err
	}
	run.Success = updateRunRequest.Success
	run.State = RunState(updateRunRequest.State)
	run.ErrorMessage = updateRunRequest.ErrorMessage

	scenarioDetails := domainScenariosToScenarios(updateRunRequest.ScenarioRunDetails)
	b, err := json.Marshal(scenarioDetails)
	if err != nil {
		return nil, err
	}
	run.ScenarioDetails = b

	err = r.conn.WithContext(ctx).Model(&Run{}).Save(&run).Error
	if err != nil {
		return nil, err
	}

	return runToDomainRun(&run)
}

func domainScenariosToScenarios(scenarios []*domain.ScenarioRunDetails) []*ScenarioDetails {
	result := []*ScenarioDetails{}
	for _, scenario := range scenarios {
		result = append(result, domainScenarioToScenario(scenario))
	}
	return result
}

func domainScenarioToScenario(scenario *domain.ScenarioRunDetails) *ScenarioDetails {
	if scenario == nil {
		return &ScenarioDetails{}
	}
	return &ScenarioDetails{
		Name:       scenario.Name,
		Duration:   scenario.Duration,
		Assertions: scenario.Assertions,
		Steps:      domainStepsToSteps(scenario.Steps),
		Success:    scenario.Success,
	}
}

func domainStepsToSteps(steps []*domain.StepRunDetails) []*Step {
	result := []*Step{}
	for _, step := range steps {
		result = append(result, &Step{
			Name:            step.Name,
			Assertions:      step.Assertions,
			URL:             step.URL,
			RequestDuration: step.RequestDuration,
			Duration:        step.Duration,
			Retries:         step.Retries,
			Success:         step.Success,
		})
	}
	return result
}

// ListForProject returns all runs for a given project list runs request.
func (r *RunRepository) ListForProject(ctx context.Context, getForProjectRequest *domain.ListRunsForProjectRequest) ([]*domain.Run, error) {
	runs := []*Run{}
	err := r.conn.
		WithContext(ctx).
		Model(&Run{}).
		Offset(getForProjectRequest.Limit*getForProjectRequest.Offset).
		Limit(getForProjectRequest.Limit).
		Where("project_id = ?", getForProjectRequest.ProjectID).
		Order("created_at desc").
		Find(&runs).
		Error
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
	scenarioRunDetails, err := scenarioRunDetailsToDomainScenarioRunDetails(run.ScenarioDetails)
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
		CreatedAt:          run.CreatedAt,
	}, nil
}

func scenarioRunDetailsToDomainScenarioRunDetails(scenario []byte) ([]*domain.ScenarioRunDetails, error) {
	details := []*ScenarioDetails{}
	err := json.Unmarshal(scenario, &details)
	if err != nil {
		return nil, err
	}
	result := []*domain.ScenarioRunDetails{}
	for _, detail := range details {
		result = append(result, &domain.ScenarioRunDetails{
			Name:       detail.Name,
			Duration:   detail.Duration,
			Assertions: detail.Assertions,
			Steps:      stepsRunDetailsToDomainStepRunDetails(detail.Steps),
			Success:    detail.Success,
		})
	}
	return result, nil
}

func stepsRunDetailsToDomainStepRunDetails(steps []*Step) []*domain.StepRunDetails {
	result := []*domain.StepRunDetails{}
	for _, step := range steps {
		result = append(result, &domain.StepRunDetails{
			Name:            step.Name,
			Assertions:      step.Assertions,
			URL:             step.URL,
			RequestDuration: step.RequestDuration,
			Duration:        step.Duration,
			Retries:         step.Retries,
			Success:         step.Success,
		})
	}
	return result
}
