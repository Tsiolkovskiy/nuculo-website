package auth

import (
	"backend/internal/repository"
)

// Manager holds all authentication services and middleware
type Manager struct {
	Config          *Config
	JWTService      *JWTService
	PasswordService *PasswordService
	AuthService     *AuthService
	Middleware      *AuthMiddleware
}

// NewManager creates a new authentication manager with all services
func NewManager(config *Config, userRepo repository.UserRepository) *Manager {
	jwtService := NewJWTService(config.JWTSecret, config.TokenDuration)
	passwordService := NewPasswordServiceWithCost(config.BCryptCost)
	authService := NewAuthService(jwtService, passwordService, userRepo)
	middleware := NewAuthMiddleware(jwtService, userRepo)

	return &Manager{
		Config:          config,
		JWTService:      jwtService,
		PasswordService: passwordService,
		AuthService:     authService,
		Middleware:      middleware,
	}
}