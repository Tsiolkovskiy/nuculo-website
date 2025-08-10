package security

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/redis/go-redis/v9"
)

// RateLimiter implements rate limiting for GraphQL operations
type RateLimiter struct {
	redis  *redis.Client
	config RateLimitConfig
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	// Global rate limits
	GlobalRequestsPerMinute int
	GlobalRequestsPerHour   int
	
	// Per-user rate limits
	UserRequestsPerMinute int
	UserRequestsPerHour   int
	
	// Per-IP rate limits
	IPRequestsPerMinute int
	IPRequestsPerHour   int
	
	// Operation-specific limits
	MutationRequestsPerMinute int
	QueryRequestsPerMinute    int
	
	// Burst allowance
	BurstSize int
}

// DefaultRateLimitConfig returns default rate limiting configuration
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		GlobalRequestsPerMinute:   1000,
		GlobalRequestsPerHour:     10000,
		UserRequestsPerMinute:     100,
		UserRequestsPerHour:       1000,
		IPRequestsPerMinute:       200,
		IPRequestsPerHour:        2000,
		MutationRequestsPerMinute: 20,
		QueryRequestsPerMinute:    200,
		BurstSize:                5,
	}
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(redisClient *redis.Client, config RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		redis:  redisClient,
		config: config,
	}
}

// ExtensionName returns the name of this extension
func (r *RateLimiter) ExtensionName() string {
	return "RateLimiter"
}

// Validate validates the schema (no-op for this extension)
func (r *RateLimiter) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

// InterceptOperation intercepts operations to apply rate limiting
func (r *RateLimiter) InterceptOperation(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	oc := graphql.GetOperationContext(ctx)
	
	// Extract client information
	clientIP := r.getClientIP(ctx)
	userID := r.getUserID(ctx)
	operationType := string(oc.Operation.Operation)
	
	// Check various rate limits
	if err := r.checkRateLimits(ctx, clientIP, userID, operationType); err != nil {
		return func(ctx context.Context) *graphql.Response {
			return &graphql.Response{
				Errors: []*graphql.Error{
					{
						Message: err.Error(),
						Extensions: map[string]interface{}{
							"code": "RATE_LIMITED",
							"retryAfter": 60, // seconds
						},
					},
				},
			}
		}
	}
	
	return next(ctx)
}

// checkRateLimits checks all applicable rate limits
func (r *RateLimiter) checkRateLimits(ctx context.Context, clientIP, userID, operationType string) error {
	now := time.Now()
	
	// Check global rate limits
	if err := r.checkLimit(ctx, "global", r.config.GlobalRequestsPerMinute, time.Minute, now); err != nil {
		return fmt.Errorf("global rate limit exceeded: %w", err)
	}
	
	if err := r.checkLimit(ctx, "global_hour", r.config.GlobalRequestsPerHour, time.Hour, now); err != nil {
		return fmt.Errorf("global hourly rate limit exceeded: %w", err)
	}
	
	// Check IP-based rate limits
	if clientIP != "" {
		ipKey := fmt.Sprintf("ip:%s", clientIP)
		if err := r.checkLimit(ctx, ipKey, r.config.IPRequestsPerMinute, time.Minute, now); err != nil {
			return fmt.Errorf("IP rate limit exceeded: %w", err)
		}
		
		ipHourKey := fmt.Sprintf("ip_hour:%s", clientIP)
		if err := r.checkLimit(ctx, ipHourKey, r.config.IPRequestsPerHour, time.Hour, now); err != nil {
			return fmt.Errorf("IP hourly rate limit exceeded: %w", err)
		}
	}
	
	// Check user-based rate limits
	if userID != "" {
		userKey := fmt.Sprintf("user:%s", userID)
		if err := r.checkLimit(ctx, userKey, r.config.UserRequestsPerMinute, time.Minute, now); err != nil {
			return fmt.Errorf("user rate limit exceeded: %w", err)
		}
		
		userHourKey := fmt.Sprintf("user_hour:%s", userID)
		if err := r.checkLimit(ctx, userHourKey, r.config.UserRequestsPerHour, time.Hour, now); err != nil {
			return fmt.Errorf("user hourly rate limit exceeded: %w", err)
		}
	}
	
	// Check operation-specific rate limits
	if operationType == "mutation" {
		mutationKey := fmt.Sprintf("mutation:%s:%s", clientIP, userID)
		if err := r.checkLimit(ctx, mutationKey, r.config.MutationRequestsPerMinute, time.Minute, now); err != nil {
			return fmt.Errorf("mutation rate limit exceeded: %w", err)
		}
	} else if operationType == "query" {
		queryKey := fmt.Sprintf("query:%s:%s", clientIP, userID)
		if err := r.checkLimit(ctx, queryKey, r.config.QueryRequestsPerMinute, time.Minute, now); err != nil {
			return fmt.Errorf("query rate limit exceeded: %w", err)
		}
	}
	
	return nil
}

// checkLimit checks a specific rate limit using sliding window algorithm
func (r *RateLimiter) checkLimit(ctx context.Context, key string, limit int, window time.Duration, now time.Time) error {
	// Use Redis sorted sets for sliding window rate limiting
	windowStart := now.Add(-window)
	windowStartScore := float64(windowStart.UnixNano())
	nowScore := float64(now.UnixNano())
	
	pipe := r.redis.Pipeline()
	
	// Remove expired entries
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%.0f", windowStartScore))
	
	// Count current requests in window
	countCmd := pipe.ZCount(ctx, key, fmt.Sprintf("%.0f", windowStartScore), "+inf")
	
	// Add current request
	pipe.ZAdd(ctx, key, redis.Z{
		Score:  nowScore,
		Member: fmt.Sprintf("%d", now.UnixNano()),
	})
	
	// Set expiration
	pipe.Expire(ctx, key, window+time.Minute)
	
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("rate limit check failed: %w", err)
	}
	
	count := countCmd.Val()
	if count >= int64(limit) {
		return fmt.Errorf("rate limit of %d requests per %v exceeded", limit, window)
	}
	
	return nil
}

// getClientIP extracts client IP from context
func (r *RateLimiter) getClientIP(ctx context.Context) string {
	// Try to get IP from various context keys
	if ip, ok := ctx.Value("client_ip").(string); ok {
		return ip
	}
	if ip, ok := ctx.Value("remote_addr").(string); ok {
		return ip
	}
	return ""
}

// getUserID extracts user ID from context
func (r *RateLimiter) getUserID(ctx context.Context) string {
	if userID, ok := ctx.Value("user_id").(string); ok {
		return userID
	}
	if user, ok := ctx.Value("user").(map[string]interface{}); ok {
		if id, exists := user["id"]; exists {
			return fmt.Sprintf("%v", id)
		}
	}
	return ""
}

// GetRateLimitStatus returns current rate limit status for debugging
func (r *RateLimiter) GetRateLimitStatus(ctx context.Context, clientIP, userID string) (map[string]interface{}, error) {
	status := make(map[string]interface{})
	now := time.Now()
	
	// Check various limits
	limits := map[string]struct {
		key    string
		limit  int
		window time.Duration
	}{
		"global_minute": {"global", r.config.GlobalRequestsPerMinute, time.Minute},
		"global_hour":   {"global_hour", r.config.GlobalRequestsPerHour, time.Hour},
	}
	
	if clientIP != "" {
		limits["ip_minute"] = struct {
			key    string
			limit  int
			window time.Duration
		}{fmt.Sprintf("ip:%s", clientIP), r.config.IPRequestsPerMinute, time.Minute}
		limits["ip_hour"] = struct {
			key    string
			limit  int
			window time.Duration
		}{fmt.Sprintf("ip_hour:%s", clientIP), r.config.IPRequestsPerHour, time.Hour}
	}
	
	if userID != "" {
		limits["user_minute"] = struct {
			key    string
			limit  int
			window time.Duration
		}{fmt.Sprintf("user:%s", userID), r.config.UserRequestsPerMinute, time.Minute}
		limits["user_hour"] = struct {
			key    string
			limit  int
			window time.Duration
		}{fmt.Sprintf("user_hour:%s", userID), r.config.UserRequestsPerHour, time.Hour}
	}
	
	for name, limitInfo := range limits {
		windowStart := now.Add(-limitInfo.window)
		windowStartScore := float64(windowStart.UnixNano())
		
		count, err := r.redis.ZCount(ctx, limitInfo.key, fmt.Sprintf("%.0f", windowStartScore), "+inf").Result()
		if err != nil {
			status[name] = map[string]interface{}{
				"error": err.Error(),
			}
			continue
		}
		
		status[name] = map[string]interface{}{
			"current":   count,
			"limit":     limitInfo.limit,
			"remaining": limitInfo.limit - int(count),
			"window":    limitInfo.window.String(),
		}
	}
	
	return status, nil
}

// ResetRateLimit resets rate limit for a specific key (admin function)
func (r *RateLimiter) ResetRateLimit(ctx context.Context, key string) error {
	return r.redis.Del(ctx, key).Err()
}

// BanIP temporarily bans an IP address
func (r *RateLimiter) BanIP(ctx context.Context, ip string, duration time.Duration) error {
	banKey := fmt.Sprintf("banned_ip:%s", ip)
	return r.redis.Set(ctx, banKey, "banned", duration).Err()
}

// IsIPBanned checks if an IP is banned
func (r *RateLimiter) IsIPBanned(ctx context.Context, ip string) (bool, error) {
	banKey := fmt.Sprintf("banned_ip:%s", ip)
	exists, err := r.redis.Exists(ctx, banKey).Result()
	return exists > 0, err
}

// InterceptField can be used to apply field-level rate limiting
func (r *RateLimiter) InterceptField(ctx context.Context, next graphql.Resolver) (interface{}, error) {
	fc := graphql.GetFieldContext(ctx)
	
	// Apply field-specific rate limiting for expensive operations
	if r.isExpensiveField(fc.Field.Name) {
		clientIP := r.getClientIP(ctx)
		userID := r.getUserID(ctx)
		
		fieldKey := fmt.Sprintf("field:%s:%s:%s", fc.Field.Name, clientIP, userID)
		if err := r.checkLimit(ctx, fieldKey, 10, time.Minute, time.Now()); err != nil {
			return nil, fmt.Errorf("field rate limit exceeded for %s: %w", fc.Field.Name, err)
		}
	}
	
	return next(ctx)
}

// isExpensiveField checks if a field is considered expensive
func (r *RateLimiter) isExpensiveField(fieldName string) bool {
	expensiveFields := map[string]bool{
		"searchPosts":    true,
		"generateReport": true,
		"exportData":     true,
		"bulkUpdate":     true,
	}
	return expensiveFields[fieldName]
}