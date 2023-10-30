//go:build integration

package sqlite

import (
	"context"

	"github.com/google/uuid"

	"github.com/inquiryproj/inquiry/internal/repository/domain"
)

func (s *SQLiteIntegrationSuite) TestListProjectsDefault() {
	projects, err := s.repository.ProjectRepository.List(context.Background(), &domain.ListProjectsRequest{
		Limit:  42,
		Offset: 0,
	})

	s.NoError(err)
	s.Equal(1, len(projects))
	s.Equal("default", projects[0].Name)
}

func (s *SQLiteIntegrationSuite) TestGetDefaultProject() {
	projectByName, err := s.repository.ProjectRepository.GetByName(context.Background(), "default")
	s.NoError(err)
	s.Equal(projectByName.Name, "default")

	projectByID, err := s.repository.ProjectRepository.GetByID(context.Background(), projectByName.ID)
	s.NoError(err)
	s.Equal(projectByName.Name, "default")
	s.Equal(projectByID.ID, projectByName.ID)
}

func (s *SQLiteIntegrationSuite) TestUnableToCreateDuplicateProject() {
	_, err := s.repository.ProjectRepository.Create(context.Background(), &domain.CreateProjectRequest{
		Name: "default",
	})

	s.Error(err)
	s.ErrorIs(err, domain.ErrProjectAlreadyExists)
}

func (s *SQLiteIntegrationSuite) TestGetProjectByIDNotFound() {
	_, err := s.repository.ProjectRepository.GetByID(context.Background(), uuid.New())

	s.Error(err)
	s.ErrorIs(err, domain.ErrProjectNotFound)
}
