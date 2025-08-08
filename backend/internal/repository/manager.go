package repository

import "backend/internal/database"

// Manager holds all repository instances
type Manager struct {
	User    UserRepository
	Post    PostRepository
	Comment CommentRepository
}

// NewManager creates a new repository manager with all repositories
func NewManager(db *database.DB) *Manager {
	return &Manager{
		User:    NewUserRepository(db),
		Post:    NewPostRepository(db),
		Comment: NewCommentRepository(db),
	}
}