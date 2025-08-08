package cache

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"backend/graph/model"
	"backend/internal/repository"
	"github.com/google/uuid"
)

// CachedUserRepository wraps UserRepository with caching
type CachedUserRepository struct {
	repo  repository.UserRepository
	cache Cache
	keys  *CacheKey
	ttl   time.Duration
}

// NewCachedUserRepository creates a new cached user repository
func NewCachedUserRepository(repo repository.UserRepository, cache Cache, ttl time.Duration) *CachedUserRepository {
	return &CachedUserRepository{
		repo:  repo,
		cache: cache,
		keys:  NewCacheKey("graphql"),
		ttl:   ttl,
	}
}

// GetByID retrieves a user by ID with caching
func (r *CachedUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	key := r.keys.User(id.String())
	
	// Try to get from cache first
	var user model.User
	if err := r.cache.Get(ctx, key, &user); err == nil {
		return &user, nil
	}
	
	// Cache miss, get from repository
	userPtr, err := r.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Store in cache
	if err := r.cache.Set(ctx, key, userPtr, r.ttl); err != nil {
		// Log cache error but don't fail the request
		fmt.Printf("Failed to cache user %s: %v\n", id, err)
	}
	
	return userPtr, nil
}

// GetByIDs retrieves multiple users by IDs with caching
func (r *CachedUserRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*model.User, error) {
	if len(ids) == 0 {
		return []*model.User{}, nil
	}

	// Prepare cache keys
	keys := make([]string, len(ids))
	keyToID := make(map[string]uuid.UUID)
	for i, id := range ids {
		key := r.keys.User(id.String())
		keys[i] = key
		keyToID[key] = id
	}

	// Try to get from cache
	cached, err := r.cache.GetMultiple(ctx, keys)
	if err != nil {
		// If cache fails, fall back to repository
		return r.repo.GetByIDs(ctx, ids)
	}

	// Separate cached and missing IDs
	var missingIDs []uuid.UUID
	userMap := make(map[uuid.UUID]*model.User)
	
	for key, value := range cached {
		if value != nil {
			var user model.User
			// Convert the cached value back to User struct
			if userBytes, ok := value.([]byte); ok {
				// Handle byte array from cache
				if err := json.Unmarshal(userBytes, &user); err == nil {
					userMap[keyToID[key]] = &user
				}
			} else if userMap, ok := value.(map[string]interface{}); ok {
				// Handle map from cache (JSON unmarshaled)
				user := convertMapToUser(userMap)
				if user != nil {
					userMap[keyToID[key]] = user
				}
			}
		}
	}

	// Find missing IDs
	for _, id := range ids {
		if _, exists := userMap[id]; !exists {
			missingIDs = append(missingIDs, id)
		}
	}

	// Fetch missing users from repository
	if len(missingIDs) > 0 {
		missingUsers, err := r.repo.GetByIDs(ctx, missingIDs)
		if err != nil {
			return nil, err
		}

		// Add missing users to map and cache them
		cacheValues := make(map[string]interface{})
		for _, user := range missingUsers {
			userMap[user.ID] = user
			key := r.keys.User(user.ID.String())
			cacheValues[key] = user
		}

		// Cache the missing users
		if len(cacheValues) > 0 {
			if err := r.cache.SetMultiple(ctx, cacheValues, r.ttl); err != nil {
				fmt.Printf("Failed to cache users: %v\n", err)
			}
		}
	}

	// Build result in the same order as requested
	result := make([]*model.User, len(ids))
	for i, id := range ids {
		if user, exists := userMap[id]; exists {
			result[i] = user
		}
		// Note: missing users will be nil in the result
	}

	return result, nil
}

// Create creates a new user and invalidates related cache
func (r *CachedUserRepository) Create(ctx context.Context, user *model.User) error {
	if err := r.repo.Create(ctx, user); err != nil {
		return err
	}
	
	// Cache the new user
	key := r.keys.User(user.ID.String())
	if err := r.cache.Set(ctx, key, user, r.ttl); err != nil {
		fmt.Printf("Failed to cache new user %s: %v\n", user.ID, err)
	}
	
	return nil
}

// Update updates a user and invalidates cache
func (r *CachedUserRepository) Update(ctx context.Context, user *model.User) error {
	if err := r.repo.Update(ctx, user); err != nil {
		return err
	}
	
	// Update cache
	key := r.keys.User(user.ID.String())
	if err := r.cache.Set(ctx, key, user, r.ttl); err != nil {
		fmt.Printf("Failed to update cached user %s: %v\n", user.ID, err)
	}
	
	return nil
}

// Delete deletes a user and removes from cache
func (r *CachedUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.repo.Delete(ctx, id); err != nil {
		return err
	}
	
	// Remove from cache
	key := r.keys.User(id.String())
	if err := r.cache.Delete(ctx, key); err != nil {
		fmt.Printf("Failed to delete cached user %s: %v\n", id, err)
	}
	
	return nil
}

// GetByEmail retrieves a user by email (not cached for security)
func (r *CachedUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return r.repo.GetByEmail(ctx, email)
}

// List retrieves users with caching
func (r *CachedUserRepository) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
	// Create a cache key based on parameters
	key := fmt.Sprintf("%s:users:list:%d:%d", r.keys.Prefix, limit, offset)
	
	// Try cache first
	var users []*model.User
	if err := r.cache.Get(ctx, key, &users); err == nil {
		return users, nil
	}
	
	// Cache miss, get from repository
	users, err := r.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	
	// Cache the result with shorter TTL for lists
	listTTL := r.ttl / 2
	if err := r.cache.Set(ctx, key, users, listTTL); err != nil {
		fmt.Printf("Failed to cache user list: %v\n", err)
	}
	
	return users, nil
}

// Helper function to convert map to User struct
func convertMapToUser(m map[string]interface{}) *model.User {
	user := &model.User{}
	
	if id, ok := m["id"].(string); ok {
		if parsedID, err := uuid.Parse(id); err == nil {
			user.ID = parsedID
		}
	}
	
	if email, ok := m["email"].(string); ok {
		user.Email = email
	}
	
	if name, ok := m["name"].(string); ok {
		user.Name = name
	}
	
	if avatar, ok := m["avatar"].(string); ok {
		user.Avatar = &avatar
	}
	
	// Add other fields as needed...
	
	return user
}