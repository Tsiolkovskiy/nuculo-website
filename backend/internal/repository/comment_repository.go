package repository

import (
	"context"
	"fmt"

	"backend/internal/database"
	"backend/internal/graph/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// commentRepository implements CommentRepository interface
type commentRepository struct {
	db *database.DB
}

// NewCommentRepository creates a new comment repository
func NewCommentRepository(db *database.DB) CommentRepository {
	return &commentRepository{db: db}
}

// Create creates a new comment
func (r *commentRepository) Create(ctx context.Context, comment *model.Comment) error {
	query := `
		INSERT INTO comments (id, content, author_id, post_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	
	_, err := r.db.Pool.Exec(ctx, query,
		comment.ID, comment.Content, comment.AuthorID, 
		comment.PostID, comment.CreatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}
	
	return nil
}

// GetByID retrieves a comment by ID
func (r *commentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Comment, error) {
	query := `
		SELECT id, content, author_id, post_id, created_at
		FROM comments 
		WHERE id = $1
	`
	
	var comment model.Comment
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&comment.ID, &comment.Content, &comment.AuthorID,
		&comment.PostID, &comment.CreatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("comment not found")
		}
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}
	
	return &comment, nil
}

// GetByPostID retrieves comments by post ID with pagination
func (r *commentRepository) GetByPostID(ctx context.Context, postID uuid.UUID, limit, offset int) ([]*model.Comment, error) {
	query := `
		SELECT id, content, author_id, post_id, created_at
		FROM comments 
		WHERE post_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.Pool.Query(ctx, query, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments by post: %w", err)
	}
	defer rows.Close()
	
	return r.scanComments(rows)
}

// Update updates an existing comment
func (r *commentRepository) Update(ctx context.Context, comment *model.Comment) error {
	query := `
		UPDATE comments 
		SET content = $2
		WHERE id = $1
	`
	
	result, err := r.db.Pool.Exec(ctx, query, comment.ID, comment.Content)
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("comment not found")
	}
	
	return nil
}

// Delete deletes a comment by ID
func (r *commentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM comments WHERE id = $1`
	
	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("comment not found")
	}
	
	return nil
}

// Count counts comments for a post
func (r *commentRepository) Count(ctx context.Context, postID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM comments WHERE post_id = $1`
	
	var count int
	err := r.db.Pool.QueryRow(ctx, query, postID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count comments: %w", err)
	}
	
	return count, nil
}

// scanComments is a helper function to scan comment rows
func (r *commentRepository) scanComments(rows pgx.Rows) ([]*model.Comment, error) {
	var comments []*model.Comment
	
	for rows.Next() {
		var comment model.Comment
		err := rows.Scan(
			&comment.ID, &comment.Content, &comment.AuthorID,
			&comment.PostID, &comment.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, &comment)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating comments: %w", err)
	}
	
	return comments, nil
}