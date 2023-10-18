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

	logger *slog.Logger
}

// NewService initialises the scenario service.
func NewService(scenarioRepository repository.Scenario, opts ...serviceOptions.Opts) *Scenario {
	options := serviceOptions.DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return &Scenario{
		scenarioRepository: scenarioRepository,
		logger:             options.Logger,
	}
}

// CreateScenario creates a new scenario.
func (s *Scenario) CreateScenario(ctx context.Context, createScenarioRequest *app.CreateScenarioRequest) (*app.Scenario, error) {
	// FIXME validate spec
	// Validate payload
	scenario, err := s.scenarioRepository.CreateScenario(ctx, &domain.CreateScenarioRequest{
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
