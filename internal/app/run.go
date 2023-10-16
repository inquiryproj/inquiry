// Package app declares the domain models.
package app

import "github.com/google/uuid"

// RunProjectRequest requests model for running a project.
type RunProjectRequest struct {
	ProjectID uuid.UUID
}

// ProjectRunOutput is the output of a project run.
type ProjectRunOutput struct {
	ID        uuid.UUID
	ProjectID uuid.UUID
	Name      string
}
