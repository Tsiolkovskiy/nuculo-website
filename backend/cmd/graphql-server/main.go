package main

import (
	"log"
	"net/http"

	"backend/graph"
	"backend/internal/auth"
	"backend/internal/database"
	"backend/internal/repository"
	"backend/internal/subscription"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func main() {
	log.Println("🚀 Starting GraphQL server with authentication...")

	// Initialize database
	db, err := database.Initialize()
	if err != nil {
		log.Printf("⚠️  Database connection failed (this is expected without PostgreSQL): %v", err)
		log.Println("📝 To run with real database:")
		log.Println("   1. Install PostgreSQL")
		log.Println("   2. Set environment variables (DB_HOST, DB_USER, etc.)")
		log.Println("   3. Run migrations: go run cmd/migrate/main.go -up")
		log.Println("   4. Run this server again")
		
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

	// Create subscription manager
	subManager := subscription.NewManager()

	// Create GraphQL resolver with dependencies
	graphqlResolver := &graph.Resolver{
		UserRepo:    repos.User,
		PostRepo:    repos.Post,
		CommentRepo: repos.Comment,
		AuthManager: authManager,
		SubManager:  subManager,
	}

	// Create GraphQL server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: graphqlResolver,
	}))

	// Add WebSocket transport for subscriptions
	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Allow all origins for development
				// In production, you should validate the origin
				return true
			},
		},
	})

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

	// Apply optional authentication middleware to GraphQL endpoint
	// This allows both authenticated and anonymous access
	r.Use(authManager.Middleware.OptionalAuth())

	// GraphQL endpoint
	r.POST("/graphql", gin.WrapH(srv))
	r.GET("/graphql", gin.WrapH(srv))

	// GraphQL Playground
	r.GET("/playground", gin.WrapH(playground.Handler("GraphQL playground", "/graphql")))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "graphql-server",
			"message": "GraphQL server with authentication is running",
		})
	})

	log.Println("🎮 GraphQL server ready at http://localhost:8080/graphql")
	log.Println("🎯 GraphQL playground available at http://localhost:8080/playground")
	log.Println("❤️  Health check at http://localhost:8080/health")
	log.Println("")
	log.Println("📋 Available GraphQL operations:")
	log.Println("   Queries:")
	log.Println("     - me (requires auth)")
	log.Println("     - user(id)")
	log.Println("     - posts(filters, pagination)")
	log.Println("     - post(id)")
	log.Println("     - searchPosts(query, limit)")
	log.Println("   Mutations:")
	log.Println("     - login(email, password)")
	log.Println("     - register(email, password, name)")
	log.Println("     - refreshToken (requires auth)")
	log.Println("     - createPost(input) (requires auth)")
	log.Println("     - updatePost(id, input) (requires auth)")
	log.Println("     - deletePost(id) (requires auth)")
	log.Println("     - addComment(postId, content) (requires auth)")
	log.Println("     - deleteComment(id) (requires auth)")
	log.Println("   Subscriptions:")
	log.Println("     - postAdded (real-time new posts)")
	log.Println("     - postUpdated(id) (real-time post updates)")
	log.Println("     - commentAdded(postId) (real-time new comments)")
	log.Println("")
	log.Printf("📡 WebSocket subscriptions enabled with %d active subscribers", subManager.GetSubscriberCount())

	r.Run(":8080")
}

func runMockDemo() {
	log.Println("🎭 Running GraphQL mock demo...")
	log.Println("📝 GraphQL schema is ready but requires database connection")
	log.Println("🔧 Set up PostgreSQL to test full GraphQL functionality")
	
	// Create subscription manager for demo
	subManager := subscription.NewManager()
	
	// Create mock resolver for demo
	mockResolver := &graph.Resolver{
		SubManager: subManager,
	}
	
	// Create GraphQL server with WebSocket support for demo
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: mockResolver,
	}))

	// Add WebSocket transport for subscriptions
	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})

	// Create Gin router for demo
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

	// GraphQL endpoint
	r.POST("/graphql", gin.WrapH(srv))
	r.GET("/graphql", gin.WrapH(srv))

	// GraphQL Playground
	r.GET("/playground", gin.WrapH(playground.Handler("GraphQL playground", "/graphql")))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "graphql-server-demo",
			"message": "GraphQL server demo with subscriptions is running",
		})
	})

	log.Println("🎮 GraphQL demo server ready at http://localhost:8080/graphql")
	log.Println("🎯 GraphQL playground available at http://localhost:8080/playground")
	log.Println("❤️  Health check at http://localhost:8080/health")
	log.Printf("📡 WebSocket subscriptions enabled with %d active subscribers", subManager.GetSubscriberCount())
	log.Println("")
	log.Println("Example GraphQL queries you can test:")
	log.Println("")
	log.Println("# Register a user")
	log.Println("mutation {")
	log.Println("  register(email: \"user@example.com\", password: \"password123\", name: \"John Doe\") {")
	log.Println("    token")
	log.Println("    user { id name email }")
	log.Println("    expiresAt")
	log.Println("  }")
	log.Println("}")
	log.Println("")
	log.Println("# Login")
	log.Println("mutation {")
	log.Println("  login(email: \"user@example.com\", password: \"password123\") {")
	log.Println("    token")
	log.Println("    user { id name email }")
	log.Println("  }")
	log.Println("}")
	log.Println("")
	log.Println("# Get current user (requires Authorization header)")
	log.Println("query {")
	log.Println("  me { id name email }")
	log.Println("}")
	log.Println("")
	log.Println("# Create a post (requires Authorization header)")
	log.Println("mutation {")
	log.Println("  createPost(input: {")
	log.Println("    title: \"My First Post\"")
	log.Println("    content: \"This is my first post!\"")
	log.Println("    tags: [\"demo\", \"graphql\"]")
	log.Println("    published: true")
	log.Println("  }) {")
	log.Println("    id title content")
	log.Println("    author { name }")
	log.Println("    tags published")
	log.Println("  }")
	log.Println("}")
	log.Println("")
	log.Println("# Get posts")
	log.Println("query {")
	log.Println("  posts {")
	log.Println("    edges {")
	log.Println("      node {")
	log.Println("        id title content")
	log.Println("        author { name }")
	log.Println("        tags published")
	log.Println("      }")
	log.Println("    }")
	log.Println("    totalCount")
	log.Println("  }")
	log.Println("}")
	log.Println("")
	log.Println("# Subscribe to new posts (WebSocket)")
	log.Println("subscription {")
	log.Println("  postAdded {")
	log.Println("    id title content")
	log.Println("    author { name }")
	log.Println("    tags published")
	log.Println("  }")
	log.Println("}")
	
	// Start the demo server
	r.Run(":8080")
}