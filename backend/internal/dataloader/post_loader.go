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

// PostLoader wraps the Post repository with DataLoader functionality
type PostLoader struct {
	postRepo repository.PostRepository
	loader   *dataloader.Loader[uuid.UUID, *model.Post]
}

// NewPostLoader creates a new PostLoader with DataLoader
func NewPostLoader(postRepo repository.PostRepository) *PostLoader {
	pl := &PostLoader{
		postRepo: postRepo,
	}

	// Create the DataLoader with batch function
	pl.loader = dataloader.NewBatchedLoader(
		pl.batchGetPosts,
		dataloader.WithWait[uuid.UUID, *model.Post](time.Millisecond*10), // Wait 10ms to batch requests
		dataloader.WithBatchCapacity[uuid.UUID, *model.Post](100),         // Max 100 items per batch
	)

	return pl
}

// Load loads a single post by ID using DataLoader
func (pl *PostLoader) Load(ctx context.Context, postID uuid.UUID) (*model.Post, error) {
	return pl.loader.Load(ctx, postID)
}

// LoadMany loads multiple posts by IDs using DataLoader
func (pl *PostLoader) LoadMany(ctx context.Context, postIDs []uuid.UUID) ([]*model.Post, []error) {
	return pl.loader.LoadMany(ctx, postIDs)
}

// Clear clears the cache for a specific post ID
func (pl *PostLoader) Clear(postID uuid.UUID) {
	pl.loader.Clear(postID)
}

// ClearAll clears all cached posts
func (pl *PostLoader) ClearAll() {
	pl.loader.ClearAll()
}

// batchGetPosts is the batch function that loads multiple posts at once
func (pl *PostLoader) batchGetPosts(ctx context.Context, postIDs []uuid.UUID) []*dataloader.Result[*model.Post] {
	// Create a map to store results
	postMap := make(map[uuid.UUID]*model.Post)
	
	// Batch load posts from repository
	posts, err := pl.postRepo.GetByIDs(ctx, postIDs)
	if err != nil {
		// If there's an error, return error for all requested IDs
		results := make([]*dataloader.Result[*model.Post], len(postIDs))
		for i := range postIDs {
			results[i] = &dataloader.Result[*model.Post]{
				Error: fmt.Errorf("failed to load posts: %w", err),
			}
		}
		return results
	}

	// Create map for quick lookup
	for _, post := range posts {
		postMap[post.ID] = post
	}

	// Create results in the same order as requested IDs
	results := make([]*dataloader.Result[*model.Post], len(postIDs))
	for i, postID := range postIDs {
		if post, exists := postMap[postID]; exists {
			results[i] = &dataloader.Result[*model.Post]{Data: post}
		} else {
			results[i] = &dataloader.Result[*model.Post]{
				Error: fmt.Errorf("post not found: %s", postID),
			}
		}
	}

	return results
}