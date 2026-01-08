package datastore

import (
	"context"
	"errors"

	"oms/server/core/model"
	"oms/server/core/types"
	"gorm.io/gorm"
)

// userStore implements types.UserStore
type userStore struct {
	db *gorm.DB
}

// NewUserStore creates a new UserStore
func NewUserStore(db *gorm.DB) types.UserStore {
	return &userStore{db: db}
}

// Create creates a new user
func (s *userStore) Create(ctx context.Context, user *model.User) error {
	return s.db.WithContext(ctx).Create(user).Error
}

// GetByID retrieves a user by ID
func (s *userStore) GetByID(ctx context.Context, userID int) (*model.User, error) {
	var user model.User
	err := s.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetByUsername retrieves a user by username
func (s *userStore) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := s.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

