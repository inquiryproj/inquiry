package domain

import "github.com/google/uuid"

// User is the user domain model.
type User struct {
	ID   uuid.UUID
	Name string
}
