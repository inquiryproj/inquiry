// Package scenario implements the scenario service.
package scenario

import (
	"context"
	"errors"
	"log/slog"

	"github.com/inquiryproj/inquiry/internal/app"
	"github.com/inquiryproj/inquiry/internal/repository"
	"github.com/inquiryproj/inquiry/internal/repository/domain"
	serviceOptions "github.com/inquiryproj/inquiry/internal/service/options"
)

// Scenario is the scenario service.
type Scenario struct {
	scenarioRepository repository.Scenario
	projectRepository  repository.Project

	logger *slog.Logger
}

// NewService initialises the scenario service.
func NewService(scenarioRepository repository.Scenario, projectRepository repository.Project, opts ...serviceOptions.Opts) *Scenario {
	options := serviceOptions.DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return &Scenario{
		scenarioRepository: scenarioRepository,
		projectRepository:  projectRepository,
		logger:             options.Logger,
	}
}

// ListScenarios returns all scenarios for a given project.
func (s *Scenario) ListScenarios(ctx context.Context, listScenariosRequest *app.ListScenariosRequest) ([]*app.Scenario, error) {
	scenarios, err := s.scenarioRepository.GetForProject(ctx, &domain.GetScenariosForProjectRequest{
		Limit:     listScenariosRequest.Limit,
		Offset:    listScenariosRequest.Offset,
		ProjectID: listScenariosRequest.ProjectID,
	})
	if err != nil {
		return nil, err
	}
	result := []*app.Scenario{}
	for _, scenario := range scenarios {
		result = append(result, scenarioToAppScenario(scenario))
	}

	return result, nil
}

// CreateScenario creates a new scenario.
func (s *Scenario) CreateScenario(ctx context.Context, createScenarioRequest *app.CreateScenarioRequest) (*app.Scenario, error) {
	// FIXME validate spec
	// Validate payload
	_, err := s.projectRepository.GetByID(ctx, createScenarioRequest.ProjectID)
	if errors.Is(err, domain.ErrProjectNotFound) {
		return nil, app.ErrProjectNotFound
	} else if err != nil {
		return nil, err
	}

	scenario, err := s.scenarioRepository.Create(ctx, &domain.CreateScenarioRequest{
		Name:      createScenarioRequest.Name,
		SpecType:  domain.ScenarioSpecType(createScenarioRequest.SpecType),
		Spec:      createScenarioRequest.Spec,
		ProjectID: createScenarioRequest.ProjectID,
	})
	if errors.Is(err, domain.ErrScenarioAlreadyExists) {
		return nil, app.ErrScenarioAlreadyExists
	} else if err != nil {
		return nil, err
	}
	return scenarioToAppScenario(scenario), nil
}

func scenarioToAppScenario(scenario *domain.Scenario) *app.Scenario {
	return &app.Scenario{
		ID:        scenario.ID,
		Name:      scenario.Name,
		SpecType:  app.ScenarioSpecType(scenario.SpecType),
		Spec:      scenario.Spec,
		ProjectID: scenario.ProjectID,
	}
}
