// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package api

import (
	"github.com/google/uuid"
)

const (
	ApiKeyAuthScopes = "ApiKeyAuth.Scopes"
)

// Defines values for ProjectRunOutputState.
const (
	Cancelled ProjectRunOutputState = "cancelled"
	Failure   ProjectRunOutputState = "failure"
	Pending   ProjectRunOutputState = "pending"
	Running   ProjectRunOutputState = "running"
	Success   ProjectRunOutputState = "success"
)

// ErrMsg defines model for ErrMsg.
type ErrMsg struct {
	Message string `json:"message"`
}

// Project defines model for Project.
type Project struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// ProjectArray defines model for ProjectArray.
type ProjectArray = []Project

// ProjectRunOutput defines model for ProjectRunOutput.
type ProjectRunOutput struct {
	ID                 uuid.UUID             `json:"id"`
	ProjectID          uuid.UUID             `json:"project_id"`
	ScenarioRunDetails []ScenarioRunDetails  `json:"scenario_run_details"`
	State              ProjectRunOutputState `json:"state"`
	Success            bool                  `json:"success"`
}

// ProjectRunOutputState defines model for ProjectRunOutput.State.
type ProjectRunOutputState string

// ProjectRunOutputArray defines model for ProjectRunOutputArray.
type ProjectRunOutputArray = []ProjectRunOutput

// Scenario defines model for Scenario.
type Scenario struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	ProjectID uuid.UUID `json:"project_id"`
	Spec      string    `json:"spec"`
	SpecType  string    `json:"spec_type"`
}

// ScenarioCreateRequest defines model for ScenarioCreateRequest.
type ScenarioCreateRequest struct {
	Name     string `json:"name"`
	Spec     string `json:"spec"`
	SpecType string `json:"spec_type"`
}

// ScenarioRunDetails defines model for ScenarioRunDetails.
type ScenarioRunDetails struct {
	Assertions   int              `json:"assertions"`
	DurationInMs int              `json:"duration_in_ms"`
	Name         string           `json:"name"`
	Steps        []StepRunDetails `json:"steps"`
	Success      bool             `json:"success"`
}

// StepRunDetails defines model for StepRunDetails.
type StepRunDetails struct {
	Assertions          int    `json:"assertions"`
	DurationInMs        int    `json:"duration_in_ms"`
	Name                string `json:"name"`
	RequestDurationInMs int    `json:"request_duration_in_ms"`
	Retries             int    `json:"retries"`
	Success             bool   `json:"success"`
	URL                 string `json:"url"`
}

// ListProjectsParams defines parameters for ListProjects.
type ListProjectsParams struct {
	// Limit The number of projects to return
	Limit *int `form:"limit,omitempty" json:"limit,omitempty"`

	// Offset The number of projects to skip
	Offset *int `form:"offset,omitempty" json:"offset,omitempty"`
}

// GetRunsForProjectParams defines parameters for GetRunsForProject.
type GetRunsForProjectParams struct {
	// Limit The number of runs to return
	Limit *int `form:"limit,omitempty" json:"limit,omitempty"`

	// Offset The number of runs to skip
	Offset *int `form:"offset,omitempty" json:"offset,omitempty"`
}

// CreateProjectJSONRequestBody defines body for CreateProject for application/json ContentType.
type CreateProjectJSONRequestBody = Project

// CreateScenarioJSONRequestBody defines body for CreateScenario for application/json ContentType.
type CreateScenarioJSONRequestBody = ScenarioCreateRequest
