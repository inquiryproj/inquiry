package domain

import "fmt"

// ErrInvalidScenarioSpecType is returned when an invalid scenario spec type is provided.
var ErrInvalidScenarioSpecType = fmt.Errorf("invalid scenario spec type")

// ErrScenarioAlreadyExists is returned when a scenario already exists.
var ErrScenarioAlreadyExists = fmt.Errorf("scenario already exists")

// ErrProjectAlreadyExists is returned when a project already exists.
var ErrProjectAlreadyExists = fmt.Errorf("project already exists")

// ErrProjectNotFound is returned when a project is not found.
var ErrProjectNotFound = fmt.Errorf("project not found")
