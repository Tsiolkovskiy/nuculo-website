package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"backend/internal/graph/model"
	"backend/internal/repository"
	"github.com/gin-gonic/gin"
)

// ContextKey is the type for context keys to avoid collisions
type ContextKey string

const (
	// UserContextKey is the key for storing user in context
	UserContextKey ContextKey = "user"
	// ClaimsContextKey is the key for storing JWT claims in context
	ClaimsContextKey ContextKey = "claims"
)

// AuthMiddleware provides authentication middleware for HTTP requests
type AuthMiddleware struct {
	jwtService *JWTService
	userRepo   repository.UserRepository
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(jwtService *JWTService, userRepo repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
		userRepo:   userRepo,
	}
}

// OptionalAuth middleware that extracts user from JWT token if present
// Does not require authentication - continues even if no token or invalid token
func (a *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// No token provided, continue without user context
			c.Next()
			return
		}

		token, err := ExtractTokenFromHeader(authHeader)
		if err != nil {
			// Invalid header format, continue without user context
			c.Next()
			return
		}

		// Validate token
		claims, err := a.jwtService.ValidateToken(token)
		if err != nil {
			// Invalid token, continue without user context
			c.Next()
			return
		}

		// Get user from database
		user, err := a.userRepo.GetByID(c.Request.Context(), claims.UserID)
		if err != nil {
			// User not found, continue without user context
			c.Next()
			return
		}

		// Add user and claims to context
		ctx := context.WithValue(c.Request.Context(), UserContextKey, user)
		ctx = context.WithValue(ctx, ClaimsContextKey, claims)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// RequiredAuth middleware that requires valid JWT token
// Returns 401 if no token or invalid token
func (a *AuthMiddleware) RequiredAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort()
			return
		}

		token, err := ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		// Validate token
		claims, err := a.jwtService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Get user from database
		user, err := a.userRepo.GetByID(c.Request.Context(), claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			c.Abort()
			return
		}

		// Add user and claims to context
		ctx := context.WithValue(c.Request.Context(), UserContextKey, user)
		ctx = context.WithValue(ctx, ClaimsContextKey, claims)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// GetUserFromContext extracts the user from the request context
func GetUserFromContext(ctx context.Context) (*model.User, bool) {
	user, ok := ctx.Value(UserContextKey).(*model.User)
	return user, ok
}

// GetClaimsFromContext extracts the JWT claims from the request context
func GetClaimsFromContext(ctx context.Context) (*JWTClaims, bool) {
	claims, ok := ctx.Value(ClaimsContextKey).(*JWTClaims)
	return claims, ok
}

// RequireUser is a helper function for GraphQL resolvers to ensure user is authenticated
func RequireUser(ctx context.Context) (*model.User, error) {
	user, ok := GetUserFromContext(ctx)
	if !ok || user == nil {
		return nil, fmt.Errorf("authentication required")
	}
	return user, nil
}

// LogAuthAttempt logs authentication attempts for security monitoring
func LogAuthAttempt(email string, success bool, ip string) {
	status := "SUCCESS"
	if !success {
		status = "FAILED"
	}
	log.Printf("AUTH_ATTEMPT: email=%s status=%s ip=%s", email, status, ip)
}