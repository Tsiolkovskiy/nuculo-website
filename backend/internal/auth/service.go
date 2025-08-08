package auth

import (
	"context"
	"fmt"
	"time"

	"backend/internal/graph/model"
	"backend/internal/repository"
	"github.com/google/uuid"
)

// AuthService provides authentication operations
type AuthService struct {
	jwtService      *JWTService
	passwordService *PasswordService
	userRepo        repository.UserRepository
}

// NewAuthService creates a new authentication service
func NewAuthService(jwtService *JWTService, passwordService *PasswordService, userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		jwtService:      jwtService,
		passwordService: passwordService,
		userRepo:        userRepo,
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Password string `json:"password" validate:"required,min=8"`
}

// AuthResponse represents an authentication response
type AuthResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	User      *model.User `json:"user"`
}

// Login authenticates a user with email and password
func (a *AuthService) Login(ctx context.Context, req LoginRequest, clientIP string) (*AuthResponse, error) {
	// Get user by email
	user, err := a.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		LogAuthAttempt(req.Email, false, clientIP)
		return nil, fmt.Errorf("invalid email or password")
	}

	// Verify password
	if err := a.passwordService.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		LogAuthAttempt(req.Email, false, clientIP)
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate JWT token
	token, expiresAt, err := a.jwtService.GenerateToken(user)
	if err != nil {
		LogAuthAttempt(req.Email, false, clientIP)
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	LogAuthAttempt(req.Email, true, clientIP)

	return &AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user,
	}, nil
}

// Register creates a new user account
func (a *AuthService) Register(ctx context.Context, req RegisterRequest, clientIP string) (*AuthResponse, error) {
	// Validate password requirements
	if err := a.passwordService.IsValidPassword(req.Password); err != nil {
		return nil, fmt.Errorf("password validation failed: %w", err)
	}

	// Check if user already exists
	existingUser, err := a.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password
	hashedPassword, err := a.passwordService.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &model.User{
		ID:           uuid.New(),
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := a.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, expiresAt, err := a.jwtService.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	LogAuthAttempt(req.Email, true, clientIP)

	return &AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user,
	}, nil
}

// RefreshToken generates a new token for an authenticated user
func (a *AuthService) RefreshToken(ctx context.Context, currentToken string) (*AuthResponse, error) {
	// Validate current token
	claims, err := a.jwtService.ValidateToken(currentToken)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Get user from database
	user, err := a.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Generate new token
	token, expiresAt, err := a.jwtService.RefreshToken(currentToken, user)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return &AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user,
	}, nil
}

// ChangePassword changes a user's password
func (a *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error {
	// Get user
	user, err := a.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Verify current password
	if err := a.passwordService.VerifyPassword(user.PasswordHash, currentPassword); err != nil {
		return fmt.Errorf("current password is incorrect")
	}

	// Validate new password
	if err := a.passwordService.IsValidPassword(newPassword); err != nil {
		return fmt.Errorf("new password validation failed: %w", err)
	}

	// Hash new password
	hashedPassword, err := a.passwordService.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update user password
	user.PasswordHash = hashedPassword
	user.UpdatedAt = time.Now()

	if err := a.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}