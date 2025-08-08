package dataloader

import (
	"context"
	"fmt"
	"time"

	"backend/graph/model"
	"backend/internal/repository"
	"github.com/google/uuid"
	"github.com/graph-gophers/dataloader/v7"
)

// UserLoader wraps the User repository with DataLoader functionality
type UserLoader struct {
	userRepo repository.UserRepository
	loader   *dataloader.Loader[uuid.UUID, *model.User]
}

// NewUserLoader creates a new UserLoader with DataLoader
func NewUserLoader(userRepo repository.UserRepository) *UserLoader {
	ul := &UserLoader{
		userRepo: userRepo,
	}

	// Create the DataLoader with batch function
	ul.loader = dataloader.NewBatchedLoader(
		ul.batchGetUsers,
		dataloader.WithWait[uuid.UUID, *model.User](time.Millisecond*10), // Wait 10ms to batch requests
		dataloader.WithBatchCapacity[uuid.UUID, *model.User](100),         // Max 100 items per batch
	)

	return ul
}

// Load loads a single user by ID using DataLoader
func (ul *UserLoader) Load(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	return ul.loader.Load(ctx, userID)
}

// LoadMany loads multiple users by IDs using DataLoader
func (ul *UserLoader) LoadMany(ctx context.Context, userIDs []uuid.UUID) ([]*model.User, []error) {
	return ul.loader.LoadMany(ctx, userIDs)
}

// Clear clears the cache for a specific user ID
func (ul *UserLoader) Clear(userID uuid.UUID) {
	ul.loader.Clear(userID)
}

// ClearAll clears all cached users
func (ul *UserLoader) ClearAll() {
	ul.loader.ClearAll()
}

// batchGetUsers is the batch function that loads multiple users at once
func (ul *UserLoader) batchGetUsers(ctx context.Context, userIDs []uuid.UUID) []*dataloader.Result[*model.User] {
	// Create a map to store results
	userMap := make(map[uuid.UUID]*model.User)
	
	// Batch load users from repository
	users, err := ul.userRepo.GetByIDs(ctx, userIDs)
	if err != nil {
		// If there's an error, return error for all requested IDs
		results := make([]*dataloader.Result[*model.User], len(userIDs))
		for i := range userIDs {
			results[i] = &dataloader.Result[*model.User]{
				Error: fmt.Errorf("failed to load users: %w", err),
			}
		}
		return results
	}

	// Create map for quick lookup
	for _, user := range users {
		userMap[user.ID] = user
	}

	// Create results in the same order as requested IDs
	results := make([]*dataloader.Result[*model.User], len(userIDs))
	for i, userID := range userIDs {
		if user, exists := userMap[userID]; exists {
			results[i] = &dataloader.Result[*model.User]{Data: user}
		} else {
			results[i] = &dataloader.Result[*model.User]{
				Error: fmt.Errorf("user not found: %s", userID),
			}
		}
	}

	return results
}