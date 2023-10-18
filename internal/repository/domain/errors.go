package domain

import "fmt"

// ErrInvalidScenarioSpecType is returned when an invalid scenario spec type is provided.
var ErrInvalidScenarioSpecType = fmt.Errorf("invalid scenario spec type")

// ErrProjectAlreadyExists is returned when a project already exists.
var ErrProjectAlreadyExists = fmt.Errorf("project already exists")

// ErrScenarioAlreadyExists is returned when a scenario already exists.
var ErrScenarioAlreadyExists = fmt.Errorf("scenario already exists")
