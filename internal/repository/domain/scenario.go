package domain

import "github.com/google/uuid"

// ScenarioSpecType is the type of the scenario spec.
type ScenarioSpecType string

// String returns the string representation of the scenario spec type.
func (s ScenarioSpecType) String() string {
	return string(s)
}

// ScenarioSpecType constants.
const (
	ScenarioSpecTypeYAML ScenarioSpecType = "yaml"
)

// ScenarioSpecTypeFromString returns the scenario spec type from string or an error if unknown.
func ScenarioSpecTypeFromString(s string) (ScenarioSpecType, error) {
	switch s {
	case "yaml":
		return ScenarioSpecTypeYAML, nil
	default:
		return "", ErrInvalidScenarioSpecType
	}
}

// CreateScenarioRequest requests model for creating a scenario.
type CreateScenarioRequest struct {
	Name      string
	SpecType  ScenarioSpecType
	Spec      string
	ProjectID uuid.UUID
}

// Scenario is the scenario domain model.
type Scenario struct {
	ID        uuid.UUID
	Name      string
	SpecType  ScenarioSpecType
	Spec      string
	ProjectID uuid.UUID
}

// GetScenariosForProjectRequest requests model for retrieving scenarios for a project.
type GetScenariosForProjectRequest struct {
	Limit     int
	Offset    int
	ProjectID uuid.UUID
}
