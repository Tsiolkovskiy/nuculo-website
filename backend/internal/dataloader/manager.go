package dataloader

import (
	"context"
	"fmt"
	"net/http"

	"backend/graph/model"
	"backend/internal/repository"
	"github.com/google/uuid"
)

// contextKey is used for storing DataLoader in context
type contextKey string

const (
	loadersKey contextKey = "dataloaders"
)

// Loaders contains all DataLoaders
type Loaders struct {
	UserLoader *UserLoader
	PostLoader *PostLoader
}

// NewLoaders creates a new set of DataLoaders
func NewLoaders(repos *repository.Manager) *Loaders {
	return &Loaders{
		UserLoader: NewUserLoader(repos.User),
		PostLoader: NewPostLoader(repos.Post),
	}
}

// Middleware creates a middleware that adds DataLoaders to the context
func Middleware(repos *repository.Manager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			loaders := NewLoaders(repos)
			ctx := context.WithValue(r.Context(), loadersKey, loaders)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// For returns the DataLoaders from the context
func For(ctx context.Context) *Loaders {
	loaders, ok := ctx.Value(loadersKey).(*Loaders)
	if !ok {
		// If no loaders in context, return nil
		// This should not happen in normal operation
		return nil
	}
	return loaders
}

// GetUser loads a user using DataLoader
func GetUser(ctx context.Context, userID string) (*model.User, error) {
	loaders := For(ctx)
	if loaders == nil {
		return nil, fmt.Errorf("no dataloaders in context")
	}
	
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	
	return loaders.UserLoader.Load(ctx, id)
}

// GetPost loads a post using DataLoader
func GetPost(ctx context.Context, postID string) (*model.Post, error) {
	loaders := For(ctx)
	if loaders == nil {
		return nil, fmt.Errorf("no dataloaders in context")
	}
	
	id, err := uuid.Parse(postID)
	if err != nil {
		return nil, fmt.Errorf("invalid post ID: %w", err)
	}
	
	return loaders.PostLoader.Load(ctx, id)
}