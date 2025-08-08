package main

import (
	"context"
	"log"
	"time"

	"backend/internal/database"
	"backend/internal/graph/model"
	"backend/internal/repository"
	"github.com/google/uuid"
)

func main() {
	log.Println("🚀 Starting database integration example...")

	// Initialize database (without migrations for now)
	db, err := database.Initialize()
	if err != nil {
		log.Printf("⚠️  Database connection failed (this is expected without PostgreSQL): %v", err)
		log.Println("📝 To run with real database:")
		log.Println("   1. Install PostgreSQL")
		log.Println("   2. Set environment variables:")
		log.Println("      DB_HOST=localhost")
		log.Println("      DB_PORT=5432")
		log.Println("      DB_USER=postgres")
		log.Println("      DB_PASSWORD=your_password")
		log.Println("      DB_NAME=graphql_typescript_go")
		log.Println("   3. Run migrations: go run cmd/migrate/main.go -up")
		log.Println("   4. Run this example again")
		return
	}
	defer db.Close()

	// Create repository manager
	repos := repository.NewManager(db)

	// Example: Create a user
	ctx := context.Background()
	user := &model.User{
		ID:           uuid.New(),
		Email:        "demo@example.com",
		Name:         "Demo User",
		PasswordHash: "hashed_password_here",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := repos.User.Create(ctx, user); err != nil {
		log.Printf("Failed to create user: %v", err)
		return
	}

	log.Printf("✅ Created user: %s (%s)", user.Name, user.Email)

	// Example: Get user by ID
	retrievedUser, err := repos.User.GetByID(ctx, user.ID)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		return
	}

	log.Printf("✅ Retrieved user: %s (%s)", retrievedUser.Name, retrievedUser.Email)

	// Example: Create a post
	post := &model.Post{
		ID:        uuid.New(),
		Title:     "My First Post",
		Content:   "This is the content of my first post using the new database layer!",
		AuthorID:  user.ID,
		Tags:      []string{"demo", "database", "graphql"},
		Published: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := repos.Post.Create(ctx, post); err != nil {
		log.Printf("Failed to create post: %v", err)
		return
	}

	log.Printf("✅ Created post: %s", post.Title)

	// Example: Get posts by author
	posts, err := repos.Post.GetByAuthorID(ctx, user.ID, 10, 0)
	if err != nil {
		log.Printf("Failed to get posts: %v", err)
		return
	}

	log.Printf("✅ Retrieved %d posts by author", len(posts))

	// Example: Create a comment
	comment := &model.Comment{
		ID:        uuid.New(),
		Content:   "Great post! Thanks for sharing.",
		AuthorID:  user.ID,
		PostID:    post.ID,
		CreatedAt: time.Now(),
	}

	if err := repos.Comment.Create(ctx, comment); err != nil {
		log.Printf("Failed to create comment: %v", err)
		return
	}

	log.Printf("✅ Created comment on post")

	// Example: Get comments for post
	comments, err := repos.Comment.GetByPostID(ctx, post.ID, 10, 0)
	if err != nil {
		log.Printf("Failed to get comments: %v", err)
		return
	}

	log.Printf("✅ Retrieved %d comments for post", len(comments))

	log.Println("🎉 Database integration example completed successfully!")
}