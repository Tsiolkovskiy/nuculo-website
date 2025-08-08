package repository

import (
	"context"
	"fmt"
	"strings"

	"backend/internal/database"
	"backend/internal/graph/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// postRepository implements PostRepository interface
type postRepository struct {
	db *database.DB
}

// NewPostRepository creates a new post repository
func NewPostRepository(db *database.DB) PostRepository {
	return &postRepository{db: db}
}

// Create creates a new post
func (r *postRepository) Create(ctx context.Context, post *model.Post) error {
	query := `
		INSERT INTO posts (id, title, content, author_id, tags, published, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	
	_, err := r.db.Pool.Exec(ctx, query,
		post.ID, post.Title, post.Content, post.AuthorID,
		post.Tags, post.Published, post.CreatedAt, post.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}
	
	return nil
}

// GetByID retrieves a post by ID
func (r *postRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Post, error) {
	query := `
		SELECT id, title, content, author_id, tags, published, created_at, updated_at
		FROM posts 
		WHERE id = $1
	`
	
	var post model.Post
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&post.ID, &post.Title, &post.Content, &post.AuthorID,
		&post.Tags, &post.Published, &post.CreatedAt, &post.UpdatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("post not found")
		}
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	
	return &post, nil
}

// GetByIDs retrieves multiple posts by their IDs (for DataLoader)
func (r *postRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*model.Post, error) {
	if len(ids) == 0 {
		return []*model.Post{}, nil
	}

	// Convert UUIDs to interface{} for the query
	args := make([]interface{}, len(ids))
	placeholders := make([]string, len(ids))
	for i, id := range ids {
		args[i] = id
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf(`
		SELECT id, title, content, author_id, tags, published, created_at, updated_at
		FROM posts 
		WHERE id IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}
	defer rows.Close()

	return r.scanPosts(rows)
}

// GetByAuthorID retrieves posts by author ID with pagination
func (r *postRepository) GetByAuthorID(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]*model.Post, error) {
	query := `
		SELECT id, title, content, author_id, tags, published, created_at, updated_at
		FROM posts 
		WHERE author_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.Pool.Query(ctx, query, authorID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts by author: %w", err)
	}
	defer rows.Close()
	
	return r.scanPosts(rows)
}

// Update updates an existing post
func (r *postRepository) Update(ctx context.Context, post *model.Post) error {
	query := `
		UPDATE posts 
		SET title = $2, content = $3, tags = $4, published = $5, updated_at = $6
		WHERE id = $1
	`
	
	result, err := r.db.Pool.Exec(ctx, query,
		post.ID, post.Title, post.Content, post.Tags, 
		post.Published, post.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("post not found")
	}
	
	return nil
}

// Delete deletes a post by ID
func (r *postRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM posts WHERE id = $1`
	
	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("post not found")
	}
	
	return nil
}

// List retrieves posts with filters and pagination
func (r *postRepository) List(ctx context.Context, filters *PostFilters, limit, offset int) ([]*model.Post, error) {
	query := `
		SELECT id, title, content, author_id, tags, published, created_at, updated_at
		FROM posts 
		WHERE 1=1
	`
	
	args := []interface{}{}
	argIndex := 1
	
	// Apply filters
	if filters != nil {
		if filters.AuthorID != nil {
			query += fmt.Sprintf(" AND author_id = $%d", argIndex)
			args = append(args, *filters.AuthorID)
			argIndex++
		}
		
		if filters.Published != nil {
			query += fmt.Sprintf(" AND published = $%d", argIndex)
			args = append(args, *filters.Published)
			argIndex++
		}
		
		if len(filters.Tags) > 0 {
			query += fmt.Sprintf(" AND tags && $%d", argIndex)
			args = append(args, filters.Tags)
			argIndex++
		}
		
		if filters.SearchTerm != nil && *filters.SearchTerm != "" {
			query += fmt.Sprintf(" AND (title ILIKE $%d OR content ILIKE $%d)", argIndex, argIndex+1)
			searchPattern := "%" + *filters.SearchTerm + "%"
			args = append(args, searchPattern, searchPattern)
			argIndex += 2
		}
	}
	
	query += " ORDER BY created_at DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)
	
	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts: %w", err)
	}
	defer rows.Close()
	
	return r.scanPosts(rows)
}

// Search searches posts by title and content
func (r *postRepository) Search(ctx context.Context, query string, limit int) ([]*model.Post, error) {
	searchQuery := `
		SELECT id, title, content, author_id, tags, published, created_at, updated_at
		FROM posts 
		WHERE published = true 
		AND (title ILIKE $1 OR content ILIKE $1)
		ORDER BY created_at DESC
		LIMIT $2
	`
	
	searchPattern := "%" + query + "%"
	rows, err := r.db.Pool.Query(ctx, searchQuery, searchPattern, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search posts: %w", err)
	}
	defer rows.Close()
	
	return r.scanPosts(rows)
}

// Count counts posts with filters
func (r *postRepository) Count(ctx context.Context, filters *PostFilters) (int, error) {
	query := `SELECT COUNT(*) FROM posts WHERE 1=1`
	args := []interface{}{}
	argIndex := 1
	
	// Apply filters
	if filters != nil {
		if filters.AuthorID != nil {
			query += fmt.Sprintf(" AND author_id = $%d", argIndex)
			args = append(args, *filters.AuthorID)
			argIndex++
		}
		
		if filters.Published != nil {
			query += fmt.Sprintf(" AND published = $%d", argIndex)
			args = append(args, *filters.Published)
			argIndex++
		}
		
		if len(filters.Tags) > 0 {
			query += fmt.Sprintf(" AND tags && $%d", argIndex)
			args = append(args, filters.Tags)
			argIndex++
		}
		
		if filters.SearchTerm != nil && *filters.SearchTerm != "" {
			query += fmt.Sprintf(" AND (title ILIKE $%d OR content ILIKE $%d)", argIndex, argIndex+1)
			searchPattern := "%" + *filters.SearchTerm + "%"
			args = append(args, searchPattern, searchPattern)
			argIndex += 2
		}
	}
	
	var count int
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count posts: %w", err)
	}
	
	return count, nil
}

// scanPosts is a helper function to scan post rows
func (r *postRepository) scanPosts(rows pgx.Rows) ([]*model.Post, error) {
	var posts []*model.Post
	
	for rows.Next() {
		var post model.Post
		err := rows.Scan(
			&post.ID, &post.Title, &post.Content, &post.AuthorID,
			&post.Tags, &post.Published, &post.CreatedAt, &post.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, &post)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating posts: %w", err)
	}
	
	return posts, nil
}