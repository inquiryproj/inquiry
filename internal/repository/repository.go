// Package repository declares the repository interfaces.
package repository

import (
	"context"

	"github.com/inquiryproj/inquiry/internal/app"
	"github.com/inquiryproj/inquiry/internal/repository/sqlite"
)

// Wrapper wraps all repositories.
type Wrapper interface {
	Project
	Scenario
}

// Project is the project repository.
type Project interface {
	GetProjects(ctx context.Context, getProjectsRequest *app.GetProjectsRequest) ([]*app.Project, error)
	CreateProject(ctx context.Context, project *app.CreateProjectRequest) (*app.Project, error)
}

// Scenario is the scenario repository.
type Scenario interface {
	CreateScenario(ctx context.Context, scenario *app.CreateScenarioRequest) (*app.Scenario, error)
}

// NewSQLiteWrapper initialises sqlite repository implementation.
// FIXME change this to new and switch on options.
func NewSQLiteWrapper() (Wrapper, error) {
	sqliteRepository, err := sqlite.NewRepository("inquiry.db")
	if err != nil {
		return nil, err
	}
	return sqliteRepository, nil
}
