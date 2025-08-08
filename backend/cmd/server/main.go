package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"backend/internal/graph/generated"
	"backend/internal/graph/resolver"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Initialize resolver
	resolver := &resolver.Resolver{}

	// Create GraphQL server
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	// Create Gin router
	r := gin.Default()

	// Enable CORS for development
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

	// GraphQL Playground for development
	r.GET("/playground", gin.WrapH(playground.Handler("GraphQL playground", "/graphql")))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"service": "graphql-typescript-go-backend",
		})
	})

	log.Printf("GraphQL server ready at http://localhost:%s/graphql", port)
	log.Printf("GraphQL playground available at http://localhost:%s/playground", port)
	
	log.Fatal(r.Run(":" + port))
}