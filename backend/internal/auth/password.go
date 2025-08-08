package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// PasswordService handles password hashing and validation
type PasswordService struct {
	cost int
}

// NewPasswordService creates a new password service
func NewPasswordService() *PasswordService {
	return &PasswordService{
		cost: bcrypt.DefaultCost, // Cost of 10
	}
}

// NewPasswordServiceWithCost creates a new password service with custom cost
func NewPasswordServiceWithCost(cost int) *PasswordService {
	return &PasswordService{
		cost: cost,
	}
}

// HashPassword hashes a plain text password
func (p *PasswordService) HashPassword(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), p.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedBytes), nil
}

// VerifyPassword verifies a plain text password against a hashed password
func (p *PasswordService) VerifyPassword(hashedPassword, password string) error {
	if hashedPassword == "" {
		return fmt.Errorf("hashed password cannot be empty")
	}
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return fmt.Errorf("invalid password")
		}
		return fmt.Errorf("failed to verify password: %w", err)
	}

	return nil
}

// IsValidPassword checks if a password meets minimum requirements
func (p *PasswordService) IsValidPassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	if len(password) > 128 {
		return fmt.Errorf("password must be less than 128 characters long")
	}

	// Check for at least one letter and one number
	hasLetter := false
	hasNumber := false
	
	for _, char := range password {
		if char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' {
			hasLetter = true
		}
		if char >= '0' && char <= '9' {
			hasNumber = true
		}
		if hasLetter && hasNumber {
			break
		}
	}

	if !hasLetter {
		return fmt.Errorf("password must contain at least one letter")
	}
	if !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}

	return nil
}