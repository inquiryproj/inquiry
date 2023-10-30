// Package project implements the project service.
package project

import (
	"context"
	"errors"
	"log/slog"

	"github.com/inquiryproj/inquiry/internal/app"
	"github.com/inquiryproj/inquiry/internal/repository"
	"github.com/inquiryproj/inquiry/internal/repository/domain"
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

// ListProjects returns a list of projects.
func (s *Project) ListProjects(ctx context.Context, getProjectsRequest *app.ListProjectsRequest) ([]*app.Project, error) {
	projects, err := s.projectRepository.List(ctx, &domain.ListProjectsRequest{
		Limit:  getProjectsRequest.Limit,
		Offset: getProjectsRequest.Offset,
	})
	if err != nil {
		return nil, err
	}
	return toAppProjects(projects), nil
}

func toAppProjects(projects []*domain.Project) []*app.Project {
	appProjects := make([]*app.Project, len(projects))
	for i, project := range projects {
		appProjects[i] = &app.Project{
			ID:   project.ID,
			Name: project.Name,
		}
	}
	return appProjects
}

// CreateProject creates a new project.
func (s *Project) CreateProject(ctx context.Context, createProjectRequest *app.CreateProjectRequest) (*app.Project, error) {
	project, err := s.projectRepository.Create(ctx, &domain.CreateProjectRequest{
		Name: createProjectRequest.Name,
	})
	if errors.Is(err, domain.ErrProjectAlreadyExists) {
		return nil, app.ErrProjectAlreadyExists
	} else if err != nil {
		return nil, err
	}
	return &app.Project{
		ID:   project.ID,
		Name: project.Name,
	}, nil
}
