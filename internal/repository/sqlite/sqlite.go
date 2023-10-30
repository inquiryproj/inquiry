// Package sqlite implements the sqlite repository.
package sqlite

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	slogGorm "github.com/orandin/slog-gorm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/inquiryproj/inquiry/internal/repository/domain"
	"github.com/inquiryproj/inquiry/pkg/crypto"
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
	ProjectRepository  *ProjectRepository
	ScenarioRepository *ScenarioRepository
	RunRepository      *RunRepository
	APIKeyRepository   *APIKeyRepository
}

// NewRepository initialises the sqlite repository.
func NewRepository(dsn string, logger *slog.Logger, options *MigrationOptions) (*Repository, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		TranslateError: true,
		Logger: slogGorm.New(
			slogGorm.WithLogger(logger),
		),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite connection: %w", err)
	}
	err = MigrateAndSeed(db, logger, options)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate sqlite database: %w", err)
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
		APIKeyRepository: &APIKeyRepository{
			conn: db,
		},
	}, nil
}

// MigrationOptions represents the options for the sqlite migration.
type MigrationOptions struct {
	APIKey string
}

// MigrateAndSeed migrates the sqlite database and
// seeds it with initial data.
func MigrateAndSeed(db *gorm.DB, logger *slog.Logger, options *MigrationOptions) error {
	err := db.AutoMigrate(
		tableList()...,
	)
	if err != nil {
		return err
	}

	return transactionExecution(db, seedSQLiteDB(logger, options))
}

func tableList() []any {
	return []any{
		&Project{},
		&Scenario{},
		&Run{},
		&User{},
		&APIKey{},
	}
}

func seedSQLiteDB(logger *slog.Logger, options *MigrationOptions) func(tx *gorm.DB) error {
	return func(tx *gorm.DB) error {
		ctx := context.Background()
		defaultUser, err := createDefaultUserIfNotExists(ctx, tx)
		if err != nil {
			return err
		}
		err = createDefaultAPIKey(ctx, tx, options.APIKey, defaultUser.ID, logger)
		if err != nil {
			return err
		}
		err = createDefaultProjectIfNotExists(ctx, tx)
		if err != nil {
			return err
		}
		return nil
	}
}

func createDefaultUserIfNotExists(ctx context.Context, tx *gorm.DB) (*domain.User, error) {
	uRepo := &UserRepository{
		conn: tx,
	}
	user, err := uRepo.GetUserByName(ctx, "default")
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return uRepo.CreateUser(ctx, "default")
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

func createDefaultAPIKey(ctx context.Context, tx *gorm.DB, key string, defaultUserID uuid.UUID, logger *slog.Logger) error {
	apiKeyRepo := &APIKeyRepository{
		conn: tx,
	}
	apiKey, err := apiKeyRepo.GetForNameAndUserID(ctx, &domain.GetAPIKeyForNameAndUserIDRequest{
		Name:   "default",
		UserID: defaultUserID,
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		isDefault := false
		if key == "" {
			var err error
			key, err = crypto.NewAPIKey()
			if err != nil {
				return err
			}
			msg := fmt.Sprintf("no API key provided, we generated a random one for you. Please store it somewhere safe: %s", key)
			logger.Warn(msg)
			logger.Warn("we recommend to generate a new API key and provide it via the API_KEY environment variable")
			isDefault = true
		}

		_, err = apiKeyRepo.CreateAPIKey(ctx, &domain.CreateAPIKeyRequest{
			Name:      "default",
			Key:       key,
			UserID:    defaultUserID,
			IsDefault: isDefault,
		})
		return err
	} else if err != nil {
		return err
	}
	if key == "" || !apiKey.IsDefault {
		return nil
	}
	err = apiKeyRepo.UpdateKey(ctx, &domain.UpdateAPIKeyRequest{
		ID:  apiKey.ID,
		Key: key,
	})
	return err
}

func createDefaultProjectIfNotExists(ctx context.Context, tx *gorm.DB) error {
	projectRepo := &ProjectRepository{
		conn: tx,
	}
	_, err := projectRepo.GetByName(ctx, "default")
	if errors.Is(err, gorm.ErrRecordNotFound) {
		_, err = projectRepo.Create(ctx, &domain.CreateProjectRequest{
			Name: "default",
		})
		return err
	} else if err != nil {
		return err
	}
	return nil
}

//nolint:nakedret,revive
func transactionExecution(db *gorm.DB, fn func(tx *gorm.DB) error) (err error) {
	tx := db.Begin()
	defer func() {
		if err != nil {
			rollBackErr := tx.Rollback().Error
			if rollBackErr != nil {
				err = fmt.Errorf("%w %w", err, rollBackErr)
				return
			}
			return
		}
		err = tx.Commit().Error
	}()
	err = fn(tx)
	return
}
