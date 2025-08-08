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

// userRepository implements UserRepository interface
type userRepository struct {
	db *database.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *database.DB) UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (id, email, name, password_hash, avatar, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	
	_, err := r.db.Pool.Exec(ctx, query,
		user.ID, user.Email, user.Name, user.PasswordHash, 
		user.Avatar, user.CreatedAt, user.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	
	return nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	query := `
		SELECT id, email, name, password_hash, avatar, created_at, updated_at
		FROM users 
		WHERE id = $1
	`
	
	var user model.User
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Name, &user.PasswordHash,
		&user.Avatar, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, email, name, password_hash, avatar, created_at, updated_at
		FROM users 
		WHERE email = $1
	`
	
	var user model.User
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Name, &user.PasswordHash,
		&user.Avatar, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	return &user, nil
}

// GetByIDs retrieves multiple users by their IDs (for DataLoader)
func (r *userRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*model.User, error) {
	if len(ids) == 0 {
		return []*model.User{}, nil
	}

	// Convert UUIDs to interface{} for the query
	args := make([]interface{}, len(ids))
	placeholders := make([]string, len(ids))
	for i, id := range ids {
		args[i] = id
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf(`
		SELECT id, email, name, password_hash, avatar, created_at, updated_at
		FROM users 
		WHERE id IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID, &user.Email, &user.Name, &user.PasswordHash,
			&user.Avatar, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

// Update updates an existing user
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users 
		SET name = $2, avatar = $3, updated_at = $4
		WHERE id = $1
	`
	
	result, err := r.db.Pool.Exec(ctx, query,
		user.ID, user.Name, user.Avatar, user.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	
	return nil
}

// Delete deletes a user by ID
func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	
	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	
	return nil
}

// List retrieves a list of users with pagination
func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
	query := `
		SELECT id, email, name, password_hash, avatar, created_at, updated_at
		FROM users 
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	
	rows, err := r.db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()
	
	var users []*model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID, &user.Email, &user.Name, &user.PasswordHash,
			&user.Avatar, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &user)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}
	
	return users, nil
}