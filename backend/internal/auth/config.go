package auth

import (
	"os"
	"strconv"
	"time"
)

// Config holds authentication configuration
type Config struct {
	JWTSecret       string
	TokenDuration   time.Duration
	BCryptCost      int
	RefreshWindow   time.Duration
}

// NewConfig creates a new authentication configuration from environment variables
func NewConfig() *Config {
	return &Config{
		JWTSecret:     getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
		TokenDuration: getDurationEnv("JWT_TOKEN_DURATION", 24*time.Hour),
		BCryptCost:    getIntEnv("BCRYPT_COST", 12),
		RefreshWindow: getDurationEnv("JWT_REFRESH_WINDOW", 2*time.Hour),
	}
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getIntEnv gets an integer environment variable with a fallback value
func getIntEnv(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

// getDurationEnv gets a duration environment variable with a fallback value
func getDurationEnv(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallback
}