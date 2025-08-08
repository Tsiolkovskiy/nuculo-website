package resolver

import (
	"context"
	"strings"
	"testing"

	"backend/internal/graph/model"
)

func TestQueryResolver_Me(t *testing.T) {
	resolver := &Resolver{}
	queryResolver := &queryResolver{resolver}

	// Test without authentication - should fail
	_, err := queryResolver.Me(context.Background())
	if err == nil {
		t.Fatal("Expected authentication error, got nil")
	}

	if !strings.Contains(err.Error(), "authentication required") {
		t.Errorf("Expected authentication error, got %v", err)
	}
}

func TestMutationResolver_CreatePost(t *testing.T) {
	resolver := &Resolver{}
	mutationResolver := &mutationResolver{resolver}

	input := model.CreatePostInput{
		Title:   "Test Post",
		Content: "Test content",
		Tags:    []string{"test", "demo"},
	}

	// Test without authentication - should fail
	_, err := mutationResolver.CreatePost(context.Background(), input)
	if err == nil {
		t.Fatal("Expected authentication error, got nil")
	}

	if !strings.Contains(err.Error(), "authentication required") {
		t.Errorf("Expected authentication error, got %v", err)
	}
}