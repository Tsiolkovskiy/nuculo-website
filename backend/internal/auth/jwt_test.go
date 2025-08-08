package auth

import (
	"testing"
	"time"

	"backend/internal/graph/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJWTService_GenerateToken(t *testing.T) {
	jwtService := NewJWTService("test-secret-key", 24*time.Hour)
	
	user := &model.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Name:  "Test User",
	}

	token, expiresAt, err := jwtService.GenerateToken(user)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.True(t, expiresAt.After(time.Now()))
	assert.True(t, expiresAt.Before(time.Now().Add(25*time.Hour)))
}

func TestJWTService_ValidateToken(t *testing.T) {
	jwtService := NewJWTService("test-secret-key", 24*time.Hour)
	
	user := &model.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Name:  "Test User",
	}

	// Generate token
	token, _, err := jwtService.GenerateToken(user)
	assert.NoError(t, err)

	// Validate token
	claims, err := jwtService.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.Email, claims.Email)
	assert.Equal(t, user.Name, claims.Name)
}

func TestJWTService_ValidateToken_InvalidToken(t *testing.T) {
	jwtService := NewJWTService("test-secret-key", 24*time.Hour)

	// Test with invalid token
	_, err := jwtService.ValidateToken("invalid-token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse token")
}

func TestJWTService_ValidateToken_ExpiredToken(t *testing.T) {
	jwtService := NewJWTService("test-secret-key", -1*time.Hour) // Expired immediately
	
	user := &model.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Name:  "Test User",
	}

	// Generate expired token
	token, _, err := jwtService.GenerateToken(user)
	assert.NoError(t, err)

	// Validate expired token
	_, err = jwtService.ValidateToken(token)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token is expired")
}

func TestJWTService_RefreshToken(t *testing.T) {
	jwtService := NewJWTService("test-secret-key", 24*time.Hour)
	
	user := &model.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Name:  "Test User",
	}

	// Generate original token
	originalToken, _, err := jwtService.GenerateToken(user)
	assert.NoError(t, err)

	// Wait a moment to ensure different timestamps
	time.Sleep(1 * time.Millisecond)

	// Refresh token
	newToken, expiresAt, err := jwtService.RefreshToken(originalToken, user)
	assert.NoError(t, err)
	assert.NotEmpty(t, newToken)
	// Note: tokens might be the same if generated at the same second, which is fine
	assert.True(t, expiresAt.After(time.Now()))
}

func TestExtractTokenFromHeader(t *testing.T) {
	tests := []struct {
		name        string
		authHeader  string
		expectedErr bool
		expected    string
	}{
		{
			name:        "Valid Bearer token",
			authHeader:  "Bearer abc123",
			expectedErr: false,
			expected:    "abc123",
		},
		{
			name:        "Empty header",
			authHeader:  "",
			expectedErr: true,
		},
		{
			name:        "Invalid format - no Bearer",
			authHeader:  "abc123",
			expectedErr: true,
		},
		{
			name:        "Invalid format - Bearer only",
			authHeader:  "Bearer",
			expectedErr: true,
		},
		{
			name:        "Invalid format - Bearer with space only",
			authHeader:  "Bearer ",
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := ExtractTokenFromHeader(tt.authHeader)
			
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, token)
			}
		})
	}
}