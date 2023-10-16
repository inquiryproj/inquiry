// Package scenario implements the scenario service.
package scenario

import (
	"context"

	"github.com/inquiryproj/inquiry/internal/app"
	"github.com/inquiryproj/inquiry/internal/repository"
)

// Scenario is the scenario service.
type Scenario struct {
	scenarioRepository repository.Scenario
}

// NewService initialises the scenario service.
func NewService(scenarioRepository repository.Scenario) *Scenario {
	return &Scenario{
		scenarioRepository: scenarioRepository,
	}
}

// CreateScenario creates a new scenario.
func (s *Scenario) CreateScenario(ctx context.Context, scenario *app.CreateScenarioRequest) (*app.Scenario, error) {
	// FIXME validate spec
	// Validate payload
	return s.scenarioRepository.CreateScenario(ctx, scenario)
}
