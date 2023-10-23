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
	RunStateSuccess   RunState = "success"
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

// GetRunsForProjectRequest requests model for getting runs for a project.
type GetRunsForProjectRequest struct {
	ProjectID uuid.UUID
	Limit     int
	Offset    int
}

// GetRunsForProjectResponse is the response model for getting runs for a project.
type GetRunsForProjectResponse struct {
	Runs []*ProjectRunOutput
}
