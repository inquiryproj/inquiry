package sqlite

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/inquiryproj/inquiry/internal/app"
)

// Scenario is the sqlite model for scenarios.
type Scenario struct {
	BaseModel
	Name      string `gorm:"index:idx_project_id_name_unique,unique"`
	SpecType  string
	Spec      string
	ProjectID uuid.UUID `gorm:"index:idx_project_id_name_unique,unique"`
}

// ScenarioRepository is the sqlite repository for projects.
type ScenarioRepository struct {
	conn *gorm.DB
}

// CreateScenario creates a new scenario in sqlite.
func (r *ScenarioRepository) CreateScenario(ctx context.Context, createScenarioRequest *app.CreateScenarioRequest) (*app.Scenario, error) {
	sqliteScenario := &Scenario{
		Name:      createScenarioRequest.Name,
		SpecType:  string(createScenarioRequest.SpecType),
		Spec:      createScenarioRequest.Spec,
		ProjectID: createScenarioRequest.ProjectID,
	}
	err := r.conn.WithContext(ctx).Create(sqliteScenario).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return nil, fmt.Errorf("%w %w", app.ErrScenarioAlreadyExists, err)
	} else if err != nil {
		return nil, err
	}
	return scenarioToAppScenario(sqliteScenario), nil
}

// GetForProject returns all scenarios for a given project.
func (r *ScenarioRepository) GetForProject(ctx context.Context, getForProjectRequest *app.GetScenariosForProjectRequest) ([]*app.Scenario, error) {
	scenarios := []*Scenario{}
	err := r.conn.WithContext(ctx).Where("project_id = ?", getForProjectRequest.ProjectID).Find(&scenarios).Error
	if err != nil {
		return nil, err
	}
	result := []*app.Scenario{}
	for _, scenario := range scenarios {
		result = append(result, scenarioToAppScenario(scenario))
	}
	return result, nil
}

func scenarioToAppScenario(scenario *Scenario) *app.Scenario {
	return &app.Scenario{
		ID:        scenario.ID,
		Name:      scenario.Name,
		SpecType:  app.ScenarioSpecType(scenario.SpecType),
		Spec:      scenario.Spec,
		ProjectID: scenario.ProjectID,
	}
}
