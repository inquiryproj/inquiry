//go:build integration

package sqlite

import (
	"context"

	"github.com/google/uuid"

	"github.com/inquiryproj/inquiry/internal/repository/domain"
)

func (s *SQLiteIntegrationSuite) TestCreateInitialRun() {
	projectID := uuid.New()
	run, err := s.repository.RunRepository.Create(context.Background(), &domain.CreateRunRequest{
		ProjectID: projectID,
	})
	s.NoError(err)
	s.Equal(projectID, run.ProjectID)
	s.Equal(false, run.Success)
	s.Equal(domain.RunStatePending, run.State)
	s.Equal("", run.ErrorMessage)
	s.Equal([]*domain.ScenarioRunDetails{}, run.ScenarioRunDetails)
}

func (s *SQLiteIntegrationSuite) TestListRunsForProject() {
	projectID := uuid.New()
	_, err := s.repository.RunRepository.Create(context.Background(), &domain.CreateRunRequest{
		ProjectID: projectID,
	})
	s.NoError(err)

	runs, err := s.repository.RunRepository.ListForProject(context.Background(), &domain.ListRunsForProjectRequest{
		ProjectID: projectID,
		Limit:     1,
		Offset:    0,
	})
	s.NoError(err)
	s.Equal(1, len(runs))
	s.Equal(projectID, runs[0].ProjectID)
	s.Equal(false, runs[0].Success)
	s.Equal(domain.RunStatePending, runs[0].State)
	s.Equal("", runs[0].ErrorMessage)
	s.Equal([]*domain.ScenarioRunDetails{}, runs[0].ScenarioRunDetails)
}

func (s *SQLiteIntegrationSuite) TestUpdateAndGet() {
	projectID := uuid.New()
	run, err := s.repository.RunRepository.Create(context.Background(), &domain.CreateRunRequest{
		ProjectID: projectID,
	})
	s.NoError(err)

	_, err = s.repository.RunRepository.Update(context.Background(), &domain.UpdateRunRequest{
		ID:                 run.ID,
		Success:            true,
		State:              domain.RunStateCompleted,
		ScenarioRunDetails: testScenarioDetails(),
	})
	s.NoError(err)

	runGet, err := s.repository.RunRepository.Get(context.Background(), run.ID)
	s.NoError(err)

	s.Equal(projectID, runGet.ProjectID)
	s.Equal(true, runGet.Success)
	s.Equal(domain.RunStateCompleted, runGet.State)
	s.Equal("", runGet.ErrorMessage)
	s.Equal(testScenarioDetails(), runGet.ScenarioRunDetails)
}

func testScenarioDetails() []*domain.ScenarioRunDetails {
	return []*domain.ScenarioRunDetails{
		{
			Name: "foo",
			Steps: []*domain.StepRunDetails{
				{
					Name:            "bar",
					Assertions:      42,
					Duration:        1.0,
					Success:         true,
					URL:             "https://example.com",
					RequestDuration: 1.0,
					Retries:         42,
				},
			},
			Duration:   1.0,
			Success:    true,
			Assertions: 42,
		},
	}
}
