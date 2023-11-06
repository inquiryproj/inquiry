// Package app declares the domain models.
package app

import (
	"time"

	"github.com/google/uuid"
)

// RunProjectRequest requests model for running a project.
type RunProjectRequest struct {
	ProjectID uuid.UUID
}

// RunProjectByNameRequest requests model for running a project for a given name.
type RunProjectByNameRequest struct {
	ProjectName string
}

// RunState is the state of a run.
type RunState string

// different run states.
const (
	RunStatePending   RunState = "pending"
	RunStateRunning   RunState = "running"
	RunStateCompleted RunState = "completed"
	RunStateFailure   RunState = "failure"
	RunstateCancelled RunState = "cancelled"
)

// ProjectRunOutput is the output of a project run.
type ProjectRunOutput struct {
	ID                 uuid.UUID
	ProjectID          uuid.UUID
	Success            bool
	State              RunState
	ScenarioRunDetails []*ScenarioRunDetails
}

// ScenarioRunDetails is the output of a scenario run.
type ScenarioRunDetails struct {
	Name       string
	Duration   time.Duration
	Assertions int
	Steps      []*StepRunDetails
	Success    bool
}

// StepRunDetails is the output of a step run.
type StepRunDetails struct {
	Name            string
	Assertions      int
	URL             string
	RequestDuration time.Duration
	Duration        time.Duration
	Retries         int
	Success         bool
}

// ListRunsForProjectRequest requests model for getting runs for a project.
type ListRunsForProjectRequest struct {
	ProjectID uuid.UUID
	Limit     int
	Offset    int
}

// ListRunsForProjectResponse is the response model for getting runs for a project.
type ListRunsForProjectResponse struct {
	Runs []*ProjectRunOutput
}
