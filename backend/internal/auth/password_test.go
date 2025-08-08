package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordService_HashPassword(t *testing.T) {
	passwordService := NewPasswordService()

	password := "testpassword123"
	hashedPassword, err := passwordService.HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
	assert.NotEqual(t, password, hashedPassword)
	assert.True(t, len(hashedPassword) > 50) // bcrypt hashes are typically 60 characters
}

func TestPasswordService_HashPassword_EmptyPassword(t *testing.T) {
	passwordService := NewPasswordService()

	_, err := passwordService.HashPassword("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "password cannot be empty")
}

func TestPasswordService_VerifyPassword(t *testing.T) {
	passwordService := NewPasswordService()

	password := "testpassword123"
	hashedPassword, err := passwordService.HashPassword(password)
	assert.NoError(t, err)

	// Test correct password
	err = passwordService.VerifyPassword(hashedPassword, password)
	assert.NoError(t, err)

	// Test incorrect password
	err = passwordService.VerifyPassword(hashedPassword, "wrongpassword")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid password")
}

func TestPasswordService_VerifyPassword_EmptyInputs(t *testing.T) {
	passwordService := NewPasswordService()

	// Test empty hashed password
	err := passwordService.VerifyPassword("", "password")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "hashed password cannot be empty")

	// Test empty password
	err = passwordService.VerifyPassword("hashedpassword", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "password cannot be empty")
}

func TestPasswordService_IsValidPassword(t *testing.T) {
	passwordService := NewPasswordService()

	tests := []struct {
		name        string
		password    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid password",
			password:    "password123",
			expectError: false,
		},
		{
			name:        "Valid complex password",
			password:    "MySecureP@ssw0rd",
			expectError: false,
		},
		{
			name:        "Too short",
			password:    "pass1",
			expectError: true,
			errorMsg:    "at least 8 characters",
		},
		{
			name:        "Too long",
			password:    "a1" + string(make([]byte, 127)), // 129 characters
			expectError: true,
			errorMsg:    "less than 128 characters",
		},
		{
			name:        "No letters",
			password:    "12345678",
			expectError: true,
			errorMsg:    "at least one letter",
		},
		{
			name:        "No numbers",
			password:    "password",
			expectError: true,
			errorMsg:    "at least one number",
		},
		{
			name:        "Minimum valid length",
			password:    "abcdefg1",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := passwordService.IsValidPassword(tt.password)
			
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPasswordService_WithCustomCost(t *testing.T) {
	customCost := 6 // Lower cost for faster testing
	passwordService := NewPasswordServiceWithCost(customCost)

	password := "testpassword123"
	hashedPassword, err := passwordService.HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)

	// Verify the password works
	err = passwordService.VerifyPassword(hashedPassword, password)
	assert.NoError(t, err)
}