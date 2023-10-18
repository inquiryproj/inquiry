// Package domain declares the domain models for the repository.
package domain

import "github.com/google/uuid"

// Project is the project domain model.
type Project struct {
	ID   uuid.UUID
	Name string
}

// GetProjectsRequest repository model for retrieving projects.
type GetProjectsRequest struct {
	Limit  int
	Offset int
}

// CreateProjectRequest requests model for creating a project.
type CreateProjectRequest struct {
	Name string
}
