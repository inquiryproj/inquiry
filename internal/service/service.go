// Package service declares the service interfaces.
package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/inquiryproj/inquiry/internal/app"
	"github.com/inquiryproj/inquiry/internal/events"
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
	ListProjects(ctx context.Context, getProjectsRequest *app.ListProjectsRequest) ([]*app.Project, error)
	CreateProject(ctx context.Context, createProjectRequest *app.CreateProjectRequest) (*app.Project, error)
}

// Scenario is the scenario service.
type Scenario interface {
	ListScenarios(ctx context.Context, listScenariosRequest *app.ListScenariosRequest) ([]*app.Scenario, error)
	CreateScenario(ctx context.Context, createScenarioRequest *app.CreateScenarioRequest) (*app.Scenario, error)
}

// Runner is the runner service.
type Runner interface {
	RunProject(ctx context.Context, run *app.RunProjectRequest) (*app.ProjectRunOutput, error)
	RunProjectByName(ctx context.Context, run *app.RunProjectByNameRequest) (*app.ProjectRunOutput, error)
	ListRunsForProject(ctx context.Context, listRunsForProjectRequest *app.ListRunsForProjectRequest) (*app.ListRunsForProjectResponse, error)
}

// NewServiceWrapper initialises all services.
func NewServiceWrapper(
	repositoryWrapper *repository.Wrapper,
	runsProducer events.Producer[uuid.UUID],
	opts ...options.Opts,
) Wrapper {
	return &struct {
		*project.Project
		*scenario.Scenario
		*runner.Runner
	}{
		project.NewService(repositoryWrapper.Project, opts...),
		scenario.NewService(repositoryWrapper.Scenario, repositoryWrapper.Project, opts...),
		runner.NewService(repositoryWrapper.Project, repositoryWrapper.Scenario, repositoryWrapper.Run, runsProducer, opts...),
	}
}
