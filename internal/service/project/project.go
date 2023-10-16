// Package project implements the project service.
package project

import (
	"context"

	"github.com/inquiryproj/inquiry/internal/app"
	"github.com/inquiryproj/inquiry/internal/repository"
)

// Project is the project service.
type Project struct {
	projectRepository repository.Project
}

// NewService initialises the project service.
func NewService(projectRepository repository.Project) *Project {
	return &Project{
		projectRepository: projectRepository,
	}
}

// GetProjects returns all projects.
func (s *Project) GetProjects(ctx context.Context, getProjectsRequest *app.GetProjectsRequest) ([]*app.Project, error) {
	return s.projectRepository.GetProjects(ctx, getProjectsRequest)
}

// CreateProject creates a new project.
func (s *Project) CreateProject(ctx context.Context, project *app.CreateProjectRequest) (*app.Project, error) {
	return s.projectRepository.CreateProject(ctx, project)
}
