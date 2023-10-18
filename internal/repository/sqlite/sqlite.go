// Package sqlite implements the sqlite repository.
package sqlite

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// BaseModel is the base model for all sqlite models.
type BaseModel struct {
	gorm.Model
	ID uuid.UUID `gorm:"type:uuid;primarykey"`
}

// BeforeCreate ensures that the ID is set to a uuid value before inserting into the database.
func (u *BaseModel) BeforeCreate(_ *gorm.DB) (err error) {
	if uuid.Nil == u.ID {
		u.ID = uuid.New()
	}
	return nil
}

// Repository is the sqlite repository.
type Repository struct {
	*ProjectRepository
	*ScenarioRepository
	*RunRepository
}

// NewRepository initialises the sqlite repository.
func NewRepository(dsn string) (*Repository, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite connection: %w", err)
	}
	err = db.AutoMigrate(
		&Project{},
		&Scenario{},
		&Run{},
	)
	if err != nil {
		return nil, err
	}
	return &Repository{
		ProjectRepository: &ProjectRepository{
			conn: db,
		},
		ScenarioRepository: &ScenarioRepository{
			conn: db,
		},
		RunRepository: &RunRepository{
			conn: db,
		},
	}, nil
}
