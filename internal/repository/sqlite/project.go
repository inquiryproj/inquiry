package sqlite

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/inquiryproj/inquiry/internal/app"
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

// GetProjects returns all projects from sqlite.
func (r *ProjectRepository) GetProjects(ctx context.Context, getProjectsRequest *app.GetProjectsRequest) ([]*app.Project, error) {
	projects := []*Project{}
	err := r.conn.
		WithContext(ctx).
		Offset(getProjectsRequest.Limit * getProjectsRequest.Offset).
		Limit(getProjectsRequest.Limit).
		Find(&projects).Error
	if err != nil {
		return nil, err
	}
	return toAppProjects(projects), nil
}

func toAppProjects(projects []*Project) []*app.Project {
	appProjects := make([]*app.Project, len(projects))
	for i, project := range projects {
		appProjects[i] = &app.Project{
			ID:   project.ID,
			Name: project.Name,
		}
	}
	return appProjects
}

// CreateProject creates a new project in sqlite.
func (r *ProjectRepository) CreateProject(ctx context.Context, project *app.CreateProjectRequest) (*app.Project, error) {
	sqliteProject := &Project{
		Name: project.Name,
	}
	err := r.conn.WithContext(ctx).Create(sqliteProject).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return nil, fmt.Errorf("%w %w", app.ErrProjectAlreadyExists, err)
	} else if err != nil {
		return nil, err
	}
	return &app.Project{
		ID:   sqliteProject.ID,
		Name: sqliteProject.Name,
	}, nil
}
