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
	return &app.Scenario{
		ID:        sqliteScenario.ID,
		Name:      sqliteScenario.Name,
		SpecType:  app.ScenarioSpecType(sqliteScenario.SpecType),
		Spec:      sqliteScenario.Spec,
		ProjectID: sqliteScenario.ProjectID,
	}, nil
}
