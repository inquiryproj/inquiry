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
	GetForProject(ctx context.Context, getForProjectRequest *app.GetScenariosForProjectRequest) ([]*app.Scenario, error)
}

// Type is the type of repository.
type Type string

// Repository types.
const (
	TypeSQLite Type = "sqlite"
)

// TypeFromString converts a string to a repository type.
func TypeFromString(repositoryType string) Type {
	switch repositoryType {
	case "sqlite":
		return TypeSQLite
	default:
		return TypeSQLite
	}
}

// Options represents the options for the handlers.
type Options struct {
	DSN  string
	Type Type
}

func defaultOptions() *Options {
	return &Options{
		DSN:  "inquiry.db",
		Type: TypeSQLite,
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

// WithType sets the repository type.
func WithType(repositoryType string) Opts {
	return func(o *Options) {
		o.Type = TypeFromString(repositoryType)
	}
}

// NewWrapper initialises the repository wrapper.
func NewWrapper(opts ...Opts) (Wrapper, error) {
	options := defaultOptions()
	for _, o := range opts {
		o(options)
	}
	switch options.Type {
	case TypeSQLite:
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
