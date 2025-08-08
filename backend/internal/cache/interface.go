package cache

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ErrCacheMiss is returned when a key is not found in cache
var ErrCacheMiss = errors.New("cache miss")

// Cache defines the interface for caching operations
type Cache interface {
	// Basic operations
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	
	// Pattern operations
	DeletePattern(ctx context.Context, pattern string) error
	
	// Atomic operations
	SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error)
	Increment(ctx context.Context, key string) (int64, error)
	IncrementWithTTL(ctx context.Context, key string, ttl time.Duration) (int64, error)
	
	// Batch operations
	GetMultiple(ctx context.Context, keys []string) (map[string]interface{}, error)
	SetMultiple(ctx context.Context, values map[string]interface{}, ttl time.Duration) error
	
	// Connection management
	Ping(ctx context.Context) error
	Close() error
}

// CacheKey generates cache keys with consistent formatting
type CacheKey struct {
	Prefix string
}

// NewCacheKey creates a new cache key generator
func NewCacheKey(prefix string) *CacheKey {
	return &CacheKey{Prefix: prefix}
}

// User generates a cache key for user data
func (ck *CacheKey) User(userID string) string {
	return ck.Prefix + ":user:" + userID
}

// Post generates a cache key for post data
func (ck *CacheKey) Post(postID string) string {
	return ck.Prefix + ":post:" + postID
}

// PostsByAuthor generates a cache key for posts by author
func (ck *CacheKey) PostsByAuthor(authorID string, limit, offset int) string {
	return fmt.Sprintf("%s:posts:author:%s:%d:%d", ck.Prefix, authorID, limit, offset)
}

// PostsList generates a cache key for posts list
func (ck *CacheKey) PostsList(filters string, limit, offset int) string {
	return fmt.Sprintf("%s:posts:list:%s:%d:%d", ck.Prefix, filters, limit, offset)
}

// SearchPosts generates a cache key for post search results
func (ck *CacheKey) SearchPosts(query string, limit int) string {
	return fmt.Sprintf("%s:posts:search:%s:%d", ck.Prefix, query, limit)
}

// Comment generates a cache key for comment data
func (ck *CacheKey) Comment(commentID string) string {
	return ck.Prefix + ":comment:" + commentID
}

// CommentsByPost generates a cache key for comments by post
func (ck *CacheKey) CommentsByPost(postID string, limit, offset int) string {
	return fmt.Sprintf("%s:comments:post:%s:%d:%d", ck.Prefix, postID, limit, offset)
}

// RateLimit generates a cache key for rate limiting
func (ck *CacheKey) RateLimit(identifier string) string {
	return ck.Prefix + ":ratelimit:" + identifier
}

// Session generates a cache key for session data
func (ck *CacheKey) Session(sessionID string) string {
	return ck.Prefix + ":session:" + sessionID
}