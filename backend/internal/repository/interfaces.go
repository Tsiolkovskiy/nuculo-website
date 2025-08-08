package repository

import (
	"context"

	"backend/internal/graph/model"
	"github.com/google/uuid"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*model.User, error)
}

// PostRepository defines the interface for post data operations
type PostRepository interface {
	Create(ctx context.Context, post *model.Post) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Post, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*model.Post, error)
	GetByAuthorID(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]*model.Post, error)
	Update(ctx context.Context, post *model.Post) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filters *PostFilters, limit, offset int) ([]*model.Post, error)
	Search(ctx context.Context, query string, limit int) ([]*model.Post, error)
	Count(ctx context.Context, filters *PostFilters) (int, error)
}

// CommentRepository defines the interface for comment data operations
type CommentRepository interface {
	Create(ctx context.Context, comment *model.Comment) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Comment, error)
	GetByPostID(ctx context.Context, postID uuid.UUID, limit, offset int) ([]*model.Comment, error)
	Update(ctx context.Context, comment *model.Comment) error
	Delete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context, postID uuid.UUID) (int, error)
}

// PostFilters represents filters for post queries
type PostFilters struct {
	AuthorID   *uuid.UUID
	Published  *bool
	Tags       []string
	SearchTerm *string
}