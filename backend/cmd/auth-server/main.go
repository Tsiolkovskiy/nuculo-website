package main

import (
	"log"
	"net/http"

	"backend/internal/auth"
	"backend/internal/database"
	"backend/internal/repository"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("🚀 Starting authentication server example...")

	// Initialize database (without migrations for now)
	db, err := database.Initialize()
	if err != nil {
		log.Printf("⚠️  Database connection failed (this is expected without PostgreSQL): %v", err)
		log.Println("📝 To run with real database:")
		log.Println("   1. Install PostgreSQL")
		log.Println("   2. Set environment variables (DB_HOST, DB_USER, etc.)")
		log.Println("   3. Run migrations: go run cmd/migrate/main.go -up")
		log.Println("   4. Run this example again")
		
		// Continue with mock demonstration
		runMockDemo()
		return
	}
	defer db.Close()

	// Create repository manager
	repos := repository.NewManager(db)

	// Create authentication manager
	authConfig := auth.NewConfig()
	authManager := auth.NewManager(authConfig, repos.User)

	// Create Gin router
	r := gin.Default()

	// Enable CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Public routes
	r.POST("/auth/register", func(c *gin.Context) {
		var req auth.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		response, err := authManager.AuthService.Register(c.Request.Context(), req, c.ClientIP())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, response)
	})

	r.POST("/auth/login", func(c *gin.Context) {
		var req auth.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		response, err := authManager.AuthService.Login(c.Request.Context(), req, c.ClientIP())
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	})

	r.POST("/auth/refresh", authManager.Middleware.RequiredAuth(), func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		token, err := auth.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
			return
		}

		response, err := authManager.AuthService.RefreshToken(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	})

	// Protected routes
	protected := r.Group("/api")
	protected.Use(authManager.Middleware.RequiredAuth())

	protected.GET("/me", func(c *gin.Context) {
		user, ok := auth.GetUserFromContext(c.Request.Context())
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user": user,
			"message": "Successfully authenticated!",
		})
	})

	// Optional auth route
	r.GET("/public", authManager.Middleware.OptionalAuth(), func(c *gin.Context) {
		user, ok := auth.GetUserFromContext(c.Request.Context())
		if ok {
			c.JSON(http.StatusOK, gin.H{
				"message": "Hello authenticated user!",
				"user":    user.Name,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "Hello anonymous user!",
			})
		}
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "auth-server",
			"message": "Authentication server is running",
		})
	})

	log.Println("🔐 Authentication server ready at http://localhost:8080")
	log.Println("📋 Available endpoints:")
	log.Println("   POST /auth/register - Register new user")
	log.Println("   POST /auth/login    - Login user")
	log.Println("   POST /auth/refresh  - Refresh token")
	log.Println("   GET  /api/me        - Get current user (protected)")
	log.Println("   GET  /public        - Public endpoint with optional auth")
	log.Println("   GET  /health        - Health check")

	r.Run(":8080")
}

func runMockDemo() {
	log.Println("🎭 Running mock authentication demo...")

	// Create auth services without database
	passwordService := auth.NewPasswordService()

	// Test password hashing
	password := "testpassword123"
	hashedPassword, err := passwordService.HashPassword(password)
	if err != nil {
		log.Printf("❌ Password hashing failed: %v", err)
		return
	}
	log.Printf("✅ Password hashed successfully")

	// Test password verification
	err = passwordService.VerifyPassword(hashedPassword, password)
	if err != nil {
		log.Printf("❌ Password verification failed: %v", err)
		return
	}
	log.Printf("✅ Password verification successful")

	// Test password validation
	err = passwordService.IsValidPassword(password)
	if err != nil {
		log.Printf("❌ Password validation failed: %v", err)
		return
	}
	log.Printf("✅ Password validation successful")

	log.Println("🎉 Mock authentication demo completed successfully!")
	log.Println("💡 Set up PostgreSQL to test full authentication flow")
}