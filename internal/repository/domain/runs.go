package domain

import (
	"time"

	"github.com/google/uuid"
)

// RunState is the state of a run.
type RunState string

// different run states.
const (
	RunStatePending   RunState = "pending"
	RunStateRunning   RunState = "running"
	RunStateSuccess   RunState = "success"
	RunStateFailure   RunState = "failure"
	RunstateCancelled RunState = "cancelled"
)

// Run is the domain model for runs.
type Run struct {
	ID                 uuid.UUID
	ProjectID          uuid.UUID
	Success            bool
	State              RunState
	ErrorMessage       string
	ScenarioRunDetails []*ScenarioRunDetails
}

// ScenarioRunDetails is the domain model for scenario run details.
type ScenarioRunDetails struct {
	Duration   time.Duration
	Assertions int
	Steps      []*StepRunDetails
	Success    bool
}

// StepRunDetails is the domain model for scenario step run details.
type StepRunDetails struct {
	Name            string
	Assertions      int
	URL             string
	RequestDuration time.Duration
	Duration        time.Duration
	Retries         int
	Success         bool
}

// CreateRunRequest is the request to create a run.
type CreateRunRequest struct {
	ProjectID uuid.UUID
}

// UpdateRunRequest is the request to update a run.
type UpdateRunRequest struct {
	ID                 uuid.UUID
	Success            bool
	State              RunState
	ErrorMessage       string
	ScenarioRunDetails []*ScenarioRunDetails
}

// GetRunsForProjectRequest is the request to get runs for a project.
type GetRunsForProjectRequest struct {
	ProjectID uuid.UUID
	Limit     int
	Offset    int
}
