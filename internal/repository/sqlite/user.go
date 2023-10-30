package sqlite

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/inquiryproj/inquiry/internal/repository/domain"
)

// User is the sqlite model for users.
type User struct {
	BaseModel
	Name string `gorm:"uniqueIndex"`
}

// UserRepository is the sqlite repository for users.
type UserRepository struct {
	conn *gorm.DB
}

// NewUserRepository initialises the sqlite user repository.
func NewUserRepository(conn *gorm.DB) *UserRepository {
	return &UserRepository{
		conn: conn,
	}
}

// CreateUser creates a user in sqlite.
func (r *UserRepository) CreateUser(ctx context.Context, name string) (*domain.User, error) {
	user := &User{
		Name: name,
	}
	err := r.conn.WithContext(ctx).Model(&User{}).Create(user).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return nil, fmt.Errorf("%w %w", domain.ErrUserAlreadyExists, err)
	} else if err != nil {
		return nil, err
	}
	return &domain.User{
		ID:   user.ID,
		Name: user.Name,
	}, nil
}

// GetUserByName returns a user from sqlite by name.
func (r *UserRepository) GetUserByName(ctx context.Context, name string) (*domain.User, error) {
	user := &User{}
	err := r.conn.WithContext(ctx).Model(&User{}).Where("name = ?", name).First(user).Error
	if err != nil {
		return nil, err
	}
	return &domain.User{
		ID:   user.ID,
		Name: user.Name,
	}, nil
}
