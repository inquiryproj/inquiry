//go:build integration

package sqlite

import (
	"context"

	"github.com/google/uuid"

	"github.com/inquiryproj/inquiry/internal/repository/domain"
)

func (s *SQLiteIntegrationSuite) TestCreateScenario() {
	projectID := uuid.New()
	scenario, err := s.repository.ScenarioRepository.Create(context.Background(), &domain.CreateScenarioRequest{
		Name:      "test scenario",
		SpecType:  domain.ScenarioSpecTypeYAML,
		Spec:      "Feature: test scenario",
		ProjectID: projectID,
	})
	s.NoError(err)
	s.Equal("test scenario", scenario.Name)
	s.Equal("Feature: test scenario", scenario.Spec)
	s.Equal(domain.ScenarioSpecTypeYAML, scenario.SpecType)
	s.Equal(projectID, scenario.ProjectID)
}

func (s *SQLiteIntegrationSuite) TestGetScenariosForProject() {
	projectID := uuid.New()
	scenario, err := s.repository.ScenarioRepository.Create(context.Background(), &domain.CreateScenarioRequest{
		Name:      "test scenario",
		SpecType:  domain.ScenarioSpecTypeYAML,
		Spec:      "Feature: test scenario",
		ProjectID: projectID,
	})
	s.NoError(err)

	scenarios, err := s.repository.ScenarioRepository.GetForProject(context.Background(), &domain.GetScenariosForProjectRequest{
		ProjectID: projectID,
	})
	s.NoError(err)

	s.Equal(1, len(scenarios))
	s.Equal(scenario.ID, scenarios[0].ID)
}
