package sqlite

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/inquiryproj/inquiry/internal/repository/domain"
)

const defaultScenarioLimit = 100

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

// NewScenarioRepository initialises the sqlite scenario repository.
func NewScenarioRepository(conn *gorm.DB) *ScenarioRepository {
	return &ScenarioRepository{
		conn: conn,
	}
}

// Create creates a new scenario in sqlite.
func (r *ScenarioRepository) Create(ctx context.Context, createScenarioRequest *domain.CreateScenarioRequest) (*domain.Scenario, error) {
	sqliteScenario := &Scenario{
		Name:      createScenarioRequest.Name,
		SpecType:  string(createScenarioRequest.SpecType),
		Spec:      createScenarioRequest.Spec,
		ProjectID: createScenarioRequest.ProjectID,
	}
	err := r.conn.WithContext(ctx).Model(&Scenario{}).Create(sqliteScenario).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return nil, fmt.Errorf("%w %w", domain.ErrScenarioAlreadyExists, err)
	} else if err != nil {
		return nil, err
	}
	return scenarioToDomainScenario(sqliteScenario), nil
}

// GetForProject returns all scenarios for a given project.
func (r *ScenarioRepository) GetForProject(ctx context.Context, getForProjectRequest *domain.GetScenariosForProjectRequest) ([]*domain.Scenario, error) {
	if getForProjectRequest.Limit == 0 {
		getForProjectRequest.Limit = defaultScenarioLimit
	}
	scenarios := []*Scenario{}
	err := r.conn.WithContext(ctx).
		Model(&Scenario{}).
		Limit(getForProjectRequest.Limit).
		Offset(getForProjectRequest.Limit*getForProjectRequest.Offset).
		Where("project_id = ?", getForProjectRequest.ProjectID).
		Find(&scenarios).Error
	if err != nil {
		return nil, err
	}
	result := []*domain.Scenario{}
	for _, scenario := range scenarios {
		result = append(result, scenarioToDomainScenario(scenario))
	}
	return result, nil
}

func scenarioToDomainScenario(scenario *Scenario) *domain.Scenario {
	return &domain.Scenario{
		ID:        scenario.ID,
		Name:      scenario.Name,
		SpecType:  domain.ScenarioSpecType(scenario.SpecType),
		Spec:      scenario.Spec,
		ProjectID: scenario.ProjectID,
	}
}
