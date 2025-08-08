package graph

import (
	"backend/internal/auth"
	"backend/internal/repository"
	"backend/internal/subscription"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

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
