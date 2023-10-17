// Package service declares the service interfaces.
package service

import (
	"context"

	"github.com/inquiryproj/inquiry/internal/app"
	"github.com/inquiryproj/inquiry/internal/repository"
	"github.com/inquiryproj/inquiry/internal/service/options"
	"github.com/inquiryproj/inquiry/internal/service/project"
	"github.com/inquiryproj/inquiry/internal/service/runner"
	"github.com/inquiryproj/inquiry/internal/service/scenario"
)

// Wrapper wraps all services.
type Wrapper interface {
	Project
	Scenario
	Runner
}

// Project is the project service.
type Project interface {
	GetProjects(ctx context.Context, getProjectsRequest *app.GetProjectsRequest) ([]*app.Project, error)
	CreateProject(ctx context.Context, createProjectRequest *app.CreateProjectRequest) (*app.Project, error)
}

// Scenario is the scenario service.
type Scenario interface {
	CreateScenario(ctx context.Context, createScenarioRequest *app.CreateScenarioRequest) (*app.Scenario, error)
}

// Runner is the runner service.
type Runner interface {
	RunProject(ctx context.Context, run *app.RunProjectRequest) (*app.ProjectRunOutput, error)
}

// NewServiceWrapper initialises all services.
func NewServiceWrapper(repositoryWrapper repository.Wrapper, opts ...options.Opts) Wrapper {
	return &struct {
		*project.Project
		*scenario.Scenario
		*runner.Runner
	}{
		project.NewService(repositoryWrapper, opts...),
		scenario.NewService(repositoryWrapper, opts...),
		runner.NewService(repositoryWrapper, opts...),
	}
}
