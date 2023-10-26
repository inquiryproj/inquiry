package sqlite

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/inquiryproj/inquiry/internal/repository/domain"
	"github.com/inquiryproj/inquiry/pkg/crypto"
)

// APIKey is the sqlite model for storing API Keys.
type APIKey struct {
	BaseModel
	Name      string    `gorm:"index:idx_api_key_name_user_id_unique,unique"`
	Key       string    `gorm:"uniqueIndex"`
	UserID    uuid.UUID `gorm:"index:idx_api_key_name_user_id_unique,unique"`
	IsDefault bool
}

// APIKeyRepository is the sqlite repository for storing API Keys.
type APIKeyRepository struct {
	conn *gorm.DB
}

// GetForNameAndUserID returns an API Key for a given name and user id.
func (r *APIKeyRepository) GetForNameAndUserID(ctx context.Context, getAPIKeyForNameAndUserIDRequest *domain.GetAPIKeyForNameAndUserIDRequest) (*domain.APIKey, error) {
	apiKey := &APIKey{}
	err := r.conn.WithContext(ctx).Model(&APIKey{}).Where("name = ? AND user_id = ?",
		getAPIKeyForNameAndUserIDRequest.Name,
		getAPIKeyForNameAndUserIDRequest.UserID).
		First(apiKey).Error
	if err != nil {
		return nil, err
	}
	return &domain.APIKey{
		ID:        apiKey.ID,
		Name:      apiKey.Name,
		UserID:    apiKey.UserID,
		IsDefault: apiKey.IsDefault,
	}, nil
}

// CreateAPIKey creates an API Key in sqlite.
func (r *APIKeyRepository) CreateAPIKey(ctx context.Context, createAPIKeyRequest *domain.CreateAPIKeyRequest) (*domain.APIKey, error) {
	apiKey := &APIKey{
		Name:      createAPIKeyRequest.Name,
		Key:       crypto.HashSHA512(createAPIKeyRequest.Key),
		UserID:    createAPIKeyRequest.UserID,
		IsDefault: createAPIKeyRequest.IsDefault,
	}
	err := r.conn.WithContext(ctx).Model(&APIKey{}).Create(apiKey).Error
	if err != nil {
		return nil, err
	}
	return &domain.APIKey{
		ID:        apiKey.ID,
		Name:      apiKey.Name,
		UserID:    apiKey.UserID,
		IsDefault: apiKey.IsDefault,
	}, nil
}

// UpdateKey updates the key of an API Key.
func (r *APIKeyRepository) UpdateKey(ctx context.Context, updateAPIKeyRequest *domain.UpdateAPIKeyRequest) error {
	err := r.conn.WithContext(ctx).Model(APIKey{}).Select("key", "is_default").Where("id = ?", updateAPIKeyRequest.ID).Updates(&APIKey{
		Key:       crypto.HashSHA512(updateAPIKeyRequest.Key),
		IsDefault: false,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

// Validate returns true if the key is valid.
func (r *APIKeyRepository) Validate(ctx context.Context, s string) (uuid.UUID, error) {
	apiKey := &APIKey{}
	err := r.conn.WithContext(ctx).Model(&APIKey{}).Where("key = ?", crypto.HashSHA512(s)).First(apiKey).Error
	if err != nil {
		return uuid.Nil, err
	}
	return apiKey.UserID, nil
}

// DeleteAPIKey deletes an API Key from sqlite.
func (r *APIKeyRepository) DeleteAPIKey(ctx context.Context, id uuid.UUID) error {
	err := r.conn.WithContext(ctx).Model(&APIKey{}).Delete(&APIKey{}, id).Error
	if err != nil {
		return err
	}
	return nil
}
