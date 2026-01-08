package fake

import (
	"context"
	"fmt"
	"sync"

	"oms/backend/core/model"
	"oms/backend/core/types"
)

// userMap maintains user state for fake store
var userMap = struct {
	sync.RWMutex
	m map[int]*model.User
	usernameMap map[string]*model.User
}{m: make(map[int]*model.User), usernameMap: make(map[string]*model.User)}

var nextUserID = 1

func init() {
	// Initialize admin user
	hashedPassword, _ := model.HashPassword("1234")
	adminUser := &model.User{
		ID:       1,
		Username: "admin",
		Password: hashedPassword,
		Role:     model.UserRoleAdmin,
	}
	userMap.Lock()
	userMap.m[1] = adminUser
	userMap.usernameMap["admin"] = adminUser
	nextUserID = 2
	userMap.Unlock()
}

// UserStoreFake is a fake implementation of UserStore for testing
type UserStoreFake struct {
	CreateFunc     func(ctx context.Context, user *model.User) error
	GetByIDFunc    func(ctx context.Context, userID int) (*model.User, error)
	GetByUsernameFunc func(ctx context.Context, username string) (*model.User, error)
}

// Create implements types.UserStore
func (f *UserStoreFake) Create(ctx context.Context, user *model.User) error {
	if f.CreateFunc != nil {
		return f.CreateFunc(ctx, user)
	}
	
	userMap.Lock()
	defer userMap.Unlock()
	
	// Check if username already exists
	if _, exists := userMap.usernameMap[user.Username]; exists {
		return fmt.Errorf("username already exists")
	}
	
	// Assign ID
	user.ID = nextUserID
	nextUserID++
	
	// Store user
	userMap.m[user.ID] = user
	userMap.usernameMap[user.Username] = user
	
	return nil
}

// GetByID implements types.UserStore
func (f *UserStoreFake) GetByID(ctx context.Context, userID int) (*model.User, error) {
	if f.GetByIDFunc != nil {
		return f.GetByIDFunc(ctx, userID)
	}
	
	userMap.RLock()
	defer userMap.RUnlock()
	
	user, exists := userMap.m[userID]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	
	// Return a copy to prevent external modification
	copiedUser := *user
	return &copiedUser, nil
}

// GetByUsername implements types.UserStore
func (f *UserStoreFake) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	if f.GetByUsernameFunc != nil {
		return f.GetByUsernameFunc(ctx, username)
	}
	
	userMap.RLock()
	defer userMap.RUnlock()
	
	user, exists := userMap.usernameMap[username]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	
	// Return a copy to prevent external modification
	copiedUser := *user
	return &copiedUser, nil
}

// Ensure UserStoreFake implements types.UserStore
var _ types.UserStore = (*UserStoreFake)(nil)

