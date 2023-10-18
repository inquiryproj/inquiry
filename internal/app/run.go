// Package app declares the domain models.
package app

import "github.com/google/uuid"

// RunProjectRequest requests model for running a project.
type RunProjectRequest struct {
	ProjectID uuid.UUID
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
	ID        uuid.UUID
	ProjectID uuid.UUID
	Success   bool
	State     RunState
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
