// Package repository declares the repository interfaces.
package repository

import (
	"context"

	"github.com/inquiryproj/inquiry/internal/repository/domain"
	"github.com/inquiryproj/inquiry/internal/repository/sqlite"
)

// Wrapper wraps all repositories.
type Wrapper struct {
	Project  Project
	Scenario Scenario
	Run      Run
}

// Project is the project repository.
type Project interface {
	GetProjects(ctx context.Context, getProjectsRequest *domain.GetProjectsRequest) ([]*domain.Project, error)
	CreateProject(ctx context.Context, project *domain.CreateProjectRequest) (*domain.Project, error)
}

// Scenario is the scenario repository.
type Scenario interface {
	CreateScenario(ctx context.Context, scenario *domain.CreateScenarioRequest) (*domain.Scenario, error)
	GetForProject(ctx context.Context, getForProjectRequest *domain.GetScenariosForProjectRequest) ([]*domain.Scenario, error)
}

// Run is the run repository.
type Run interface {
	CreateRun(ctx context.Context, createRunRequest *domain.CreateRunRequest) (*domain.Run, error)
	UpdateRun(ctx context.Context, updateRunRequest *domain.UpdateRunRequest) (*domain.Run, error)
	GetForProject(ctx context.Context, getForProjectRequest *domain.GetRunsForProjectRequest) ([]*domain.Run, error)
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
func NewWrapper(opts ...Opts) (*Wrapper, error) {
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
func NewSQLiteWrapper(dsn string) (*Wrapper, error) {
	sqliteRepository, err := sqlite.NewRepository(dsn)
	if err != nil {
		return nil, err
	}
	return &Wrapper{
		Project:  sqliteRepository.ProjectRepository,
		Scenario: sqliteRepository.ScenarioRepository,
		Run:      sqliteRepository.RunRepository,
	}, nil
}
