//go:build integration

package sqlite

import (
	"context"

	"github.com/inquiryproj/inquiry/internal/repository/domain"
)

func (s *SQLiteIntegrationSuite) TestGetAPIKeyForNameAnd() {
	user, err := s.repository.UserRepository.GetUserByName(context.Background(), "default")
	s.NoError(err)
	apiKey, err := s.repository.APIKeyRepository.GetForNameAndUserID(context.Background(), &domain.GetAPIKeyForNameAndUserIDRequest{
		Name:   "default",
		UserID: user.ID,
	})
	s.NoError(err)
	s.Equal("default", apiKey.Name)
	s.Equal(user.ID, apiKey.UserID)
}

func (s *SQLiteIntegrationSuite) TestValidateAPIKey() {
	user, err := s.repository.UserRepository.GetUserByName(context.Background(), "default")
	s.NoError(err)
	// Create a new key
	apiKey, err := s.repository.APIKeyRepository.CreateAPIKey(context.Background(), &domain.CreateAPIKeyRequest{
		Name:   "test",
		Key:    "test",
		UserID: user.ID,
	})
	s.NoError(err)
	s.Equal("test", apiKey.Name)
	s.Equal(user.ID, apiKey.UserID)
	s.Equal(false, apiKey.IsDefault)

	// Validate correct key
	userID, err := s.repository.APIKeyRepository.Validate(context.Background(), "test")
	s.NoError(err)
	s.Equal(userID, apiKey.UserID)

	// Validate incorrect key
	_, err = s.repository.APIKeyRepository.Validate(context.Background(), "invalid")
	s.Error(err)

	// Delete key
	err = s.repository.APIKeyRepository.DeleteAPIKey(context.Background(), apiKey.ID)
	s.NoError(err)

	// Validate correct key doesn't work anymore
	_, err = s.repository.APIKeyRepository.Validate(context.Background(), "test")
	s.Error(err)
}

func (s *SQLiteIntegrationSuite) TestUpdateAPIKey() {
	user, err := s.repository.UserRepository.GetUserByName(context.Background(), "default")
	s.NoError(err)
	// Create a new key
	apiKey, err := s.repository.APIKeyRepository.CreateAPIKey(context.Background(), &domain.CreateAPIKeyRequest{
		Name:   "test",
		Key:    "test",
		UserID: user.ID,
	})
	s.NoError(err)
	s.Equal("test", apiKey.Name)
	s.Equal(user.ID, apiKey.UserID)
	s.Equal(false, apiKey.IsDefault)

	// Validate correct key
	userID, err := s.repository.APIKeyRepository.Validate(context.Background(), "test")
	s.NoError(err)
	s.Equal(userID, apiKey.UserID)

	// Update the key
	err = s.repository.APIKeyRepository.UpdateKey(context.Background(), &domain.UpdateAPIKeyRequest{
		ID:  apiKey.ID,
		Key: "new",
	})
	s.NoError(err)

	// Validate correct key doesn't work anymore
	_, err = s.repository.APIKeyRepository.Validate(context.Background(), "test")
	s.Error(err)
}
