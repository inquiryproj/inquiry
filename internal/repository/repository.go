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

// RepositoryType is the type of repository.
type RepositoryType string

// Repository types.
const (
	RepositoryTypeSQLite RepositoryType = "sqlite"
)

// RepositoryTypeFromString converts a string to a repository type.
func RepositoryTypeFromString(repositoryType string) RepositoryType {
	switch repositoryType {
	case "sqlite":
		return RepositoryTypeSQLite
	default:
		return RepositoryTypeSQLite
	}
}

// Options represents the options for the handlers.
type Options struct {
	DSN            string
	RepositoryType RepositoryType
}

func defaultOptions() *Options {
	return &Options{
		DSN:            "inquiry.db",
		RepositoryType: RepositoryTypeSQLite,
	}
}

// Opts represents a function that modifies the options.
type Opts func(*Options)

// WithDSN sets the DSN.
func WithDSN(dsn string) Opts {
	return func(o *Options) {
		o.DSN = dsn
	}
}

// WithRepositoryType sets the repository type.
func WithRepositoryType(repositoryType string) Opts {
	return func(o *Options) {
		o.RepositoryType = RepositoryTypeFromString(repositoryType)
	}
}

func NewWrapper(opts ...Opts) (Wrapper, error) {
	options := defaultOptions()
	for _, o := range opts {
		o(options)
	}
	switch options.RepositoryType {
	case RepositoryTypeSQLite:
		return NewSQLiteWrapper(options.DSN)
	default:
		return NewSQLiteWrapper(options.DSN)
	}
}

// NewSQLiteWrapper initialises sqlite repository implementation.
func NewSQLiteWrapper(dsn string) (Wrapper, error) {
	sqliteRepository, err := sqlite.NewRepository(dsn)
	if err != nil {
		return nil, err
	}
	return sqliteRepository, nil
}
