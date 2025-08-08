package model

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Email        string     `json:"email" db:"email"`
	Name         string     `json:"name" db:"name"`
	PasswordHash string     `json:"-" db:"password_hash"` // Never expose password hash in JSON
	Avatar       *string    `json:"avatar" db:"avatar"`
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time  `json:"updatedAt" db:"updated_at"`
}

// Post represents a blog post
type Post struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	AuthorID  uuid.UUID `json:"authorId" db:"author_id"`
	Tags      []string  `json:"tags" db:"tags"`
	Published bool      `json:"published" db:"published"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// Comment represents a comment on a post
type Comment struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Content   string    `json:"content" db:"content"`
	AuthorID  uuid.UUID `json:"authorId" db:"author_id"`
	PostID    uuid.UUID `json:"postId" db:"post_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

// CreateUserInput represents input for creating a user
type CreateUserInput struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Password string `json:"password" validate:"required,min=8"`
}

// UpdateUserInput represents input for updating a user
type UpdateUserInput struct {
	Name   *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Avatar *string `json:"avatar,omitempty"`
}

// GraphQL Input Types
type CreatePostInput struct {
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Tags      []string `json:"tags"`
	Published *bool    `json:"published,omitempty"`
}

type UpdatePostInput struct {
	Title     *string  `json:"title,omitempty"`
	Content   *string  `json:"content,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	Published *bool    `json:"published,omitempty"`
}

type PostFilters struct {
	AuthorID   *string  `json:"authorId,omitempty"`
	Published  *bool    `json:"published,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	SearchTerm *string  `json:"searchTerm,omitempty"`
}

type PaginationInput struct {
	Page  *int `json:"page,omitempty"`
	Limit *int `json:"limit,omitempty"`
}

// GraphQL Response Types
type PostConnection struct {
	Edges      []*PostEdge `json:"edges"`
	PageInfo   *PageInfo   `json:"pageInfo"`
	TotalCount int         `json:"totalCount"`
}

type PostEdge struct {
	Node   *Post  `json:"node"`
	Cursor string `json:"cursor"`
}

type PageInfo struct {
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
	StartCursor     *string `json:"startCursor,omitempty"`
	EndCursor       *string `json:"endCursor,omitempty"`
}

type AuthPayload struct {
	Token     string    `json:"token"`
	User      *User     `json:"user"`
	ExpiresAt time.Time `json:"expiresAt"`
}