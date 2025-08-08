package resolver

import (
	"backend/internal/auth"
	"backend/internal/repository"
	"backend/internal/subscription"
)

// Resolver is the root resolver with service dependencies
type Resolver struct {
	// Repository dependencies
	UserRepo    repository.UserRepository
	PostRepo    repository.PostRepository
	CommentRepo repository.CommentRepository
	
	// Authentication service
	AuthManager *auth.Manager
	
	// Subscription manager for real-time updates
	SubManager *subscription.Manager
}