package sqlite

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/inquiryproj/inquiry/internal/repository/domain"
)

// Project is the sqlite project model.
type Project struct {
	BaseModel
	Name      string `gorm:"uniqueIndex"`
	Scenarios []*Scenario
}

// ProjectRepository is the sqlite repository for projects.
type ProjectRepository struct {
	conn *gorm.DB
}

// NewProjectRepository initialises the sqlite project repository.
func NewProjectRepository(conn *gorm.DB) *ProjectRepository {
	return &ProjectRepository{
		conn: conn,
	}
}

// GetByID returns a project from sqlite by id.
func (r *ProjectRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	project := &Project{}

	err := r.conn.WithContext(ctx).Model(&Project{}).Where("id = ?", id).First(project).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w %w", domain.ErrProjectNotFound, err)
	} else if err != nil {
		return nil, err
	}
	return &domain.Project{
		ID:   project.ID,
		Name: project.Name,
	}, nil
}

// GetByName returns a project from sqlite by name.
func (r *ProjectRepository) GetByName(ctx context.Context, name string) (*domain.Project, error) {
	project := &Project{}

	err := r.conn.WithContext(ctx).Model(&Project{}).Where("name = ?", name).First(project).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w %w", domain.ErrProjectNotFound, err)
	} else if err != nil {
		return nil, err
	}
	return &domain.Project{
		ID:   project.ID,
		Name: project.Name,
	}, nil
}

// List returns all projects from sqlite.
func (r *ProjectRepository) List(ctx context.Context, getProjectsRequest *domain.ListProjectsRequest) ([]*domain.Project, error) {
	projects := []*Project{}
	err := r.conn.
		WithContext(ctx).
		Model(&Project{}).
		Offset(getProjectsRequest.Limit * getProjectsRequest.Offset).
		Limit(getProjectsRequest.Limit).
		Find(&projects).Error
	if err != nil {
		return nil, err
	}
	return toAppProjects(projects), nil
}

func toAppProjects(projects []*Project) []*domain.Project {
	appProjects := make([]*domain.Project, len(projects))
	for i, project := range projects {
		appProjects[i] = &domain.Project{
			ID:   project.ID,
			Name: project.Name,
		}
	}
	return appProjects
}

// Create creates a new project in sqlite.
func (r *ProjectRepository) Create(ctx context.Context, project *domain.CreateProjectRequest) (*domain.Project, error) {
	sqliteProject := &Project{
		Name: project.Name,
	}
	err := r.conn.WithContext(ctx).Create(sqliteProject).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return nil, fmt.Errorf("%w %w", domain.ErrProjectAlreadyExists, err)
	} else if err != nil {
		return nil, err
	}
	return &domain.Project{
		ID:   sqliteProject.ID,
		Name: sqliteProject.Name,
	}, nil
}
