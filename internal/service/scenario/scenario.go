// Package scenario implements the scenario service.
package scenario

import (
	"context"
	"log/slog"

	"github.com/inquiryproj/inquiry/internal/app"
	"github.com/inquiryproj/inquiry/internal/repository"
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
func (s *Scenario) CreateScenario(ctx context.Context, scenario *app.CreateScenarioRequest) (*app.Scenario, error) {
	// FIXME validate spec
	// Validate payload
	return s.scenarioRepository.CreateScenario(ctx, scenario)
}
