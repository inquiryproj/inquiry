//go:build integration

package sqlite

import (
	"context"

	"github.com/inquiryproj/inquiry/internal/repository/domain"
)

func (s *SQLiteIntegrationSuite) TestCreateUserAlreadyExists() {
	_, err := s.repository.UserRepository.CreateUser(context.Background(), "default")

	s.Error(err)
	s.ErrorIs(err, domain.ErrUserAlreadyExists)
}

func (s *SQLiteIntegrationSuite) TestGetDefaultUser() {
	user, err := s.repository.UserRepository.GetUserByName(context.Background(), "default")

	s.NoError(err)
	s.Equal("default", user.Name)
}
