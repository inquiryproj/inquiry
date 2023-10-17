// Package project implements the project service.
package project

import (
	"context"
	"log/slog"

	"github.com/inquiryproj/inquiry/internal/app"
	"github.com/inquiryproj/inquiry/internal/repository"
	serviceOptions "github.com/inquiryproj/inquiry/internal/service/options"
)

// Project is the project service.
type Project struct {
	projectRepository repository.Project

	logger *slog.Logger
}

// NewService initialises the project service.
func NewService(projectRepository repository.Project, opts ...serviceOptions.Opts) *Project {
	options := serviceOptions.DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return &Project{
		projectRepository: projectRepository,
		logger:            options.Logger,
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
