package domain

import "github.com/google/uuid"

// APIKey is the API key domain model.
type APIKey struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	IsDefault bool
}

// CreateAPIKeyRequest is the request to create an API key.
type CreateAPIKeyRequest struct {
	Name      string
	Key       string
	UserID    uuid.UUID
	IsDefault bool
}

// UpdateAPIKeyRequest is the request to update an API key.
type UpdateAPIKeyRequest struct {
	ID  uuid.UUID
	Key string
}

// GetAPIKeyForNameAndUserIDRequest is the request to get an API key for a given name and user id.
type GetAPIKeyForNameAndUserIDRequest struct {
	Name   string
	UserID uuid.UUID
}
