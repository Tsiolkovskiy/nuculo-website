package resolver

import (
	"context"
	"testing"
	"time"

	"backend/internal/auth"
	"backend/internal/graph/model"
	"backend/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repositories for testing
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepo) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepo) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*model.User), args.Error(1)
}

type MockPostRepo struct {
	mock.Mock
}

func (m *MockPostRepo) Create(ctx context.Context, post *model.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *MockPostRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Post, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Post), args.Error(1)
}

func (m *MockPostRepo) GetByAuthorID(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]*model.Post, error) {
	args := m.Called(ctx, authorID, limit, offset)
	return args.Get(0).([]*model.Post), args.Error(1)
}

func (m *MockPostRepo) Update(ctx context.Context, post *model.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *MockPostRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPostRepo) List(ctx context.Context, filters *repository.PostFilters, limit, offset int) ([]*model.Post, error) {
	args := m.Called(ctx, filters, limit, offset)
	return args.Get(0).([]*model.Post), args.Error(1)
}

func (m *MockPostRepo) Search(ctx context.Context, query string, limit int) ([]*model.Post, error) {
	args := m.Called(ctx, query, limit)
	return args.Get(0).([]*model.Post), args.Error(1)
}

func (m *MockPostRepo) Count(ctx context.Context, filters *repository.PostFilters) (int, error) {
	args := m.Called(ctx, filters)
	return args.Int(0), args.Error(1)
}

type MockCommentRepo struct {
	mock.Mock
}

func (m *MockCommentRepo) Create(ctx context.Context, comment *model.Comment) error {
	args := m.Called(ctx, comment)
	return args.Error(0)
}

func (m *MockCommentRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Comment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Comment), args.Error(1)
}

func (m *MockCommentRepo) GetByPostID(ctx context.Context, postID uuid.UUID, limit, offset int) ([]*model.Comment, error) {
	args := m.Called(ctx, postID, limit, offset)
	return args.Get(0).([]*model.Comment), args.Error(1)
}

func (m *MockCommentRepo) Update(ctx context.Context, comment *model.Comment) error {
	args := m.Called(ctx, comment)
	return args.Error(0)
}

func (m *MockCommentRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCommentRepo) Count(ctx context.Context, postID uuid.UUID) (int, error) {
	args := m.Called(ctx, postID)
	return args.Int(0), args.Error(1)
}

// Test setup helper
func setupTestResolver() (*Resolver, *MockUserRepo, *MockPostRepo, *MockCommentRepo) {
	mockUserRepo := new(MockUserRepo)
	mockPostRepo := new(MockPostRepo)
	mockCommentRepo := new(MockCommentRepo)

	// Create auth manager
	authConfig := auth.NewConfig()
	authManager := auth.NewManager(authConfig, mockUserRepo)

	resolver := &Resolver{
		UserRepo:    mockUserRepo,
		PostRepo:    mockPostRepo,
		CommentRepo: mockCommentRepo,
		AuthManager: authManager,
	}

	return resolver, mockUserRepo, mockPostRepo, mockCommentRepo
}

// Helper to create authenticated context
func createAuthenticatedContext(user *model.User) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, auth.UserContextKey, user)
	return ctx
}

func TestQueryResolver_User(t *testing.T) {
	resolver, mockUserRepo, _, _ := setupTestResolver()
	queryResolver := &queryResolver{resolver}

	userID := uuid.New()
	expectedUser := &model.User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}

	mockUserRepo.On("GetByID", mock.Anything, userID).Return(expectedUser, nil)

	user, err := queryResolver.User(context.Background(), userID.String())

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockUserRepo.AssertExpectations(t)
}

func TestQueryResolver_Me_RequiresAuth(t *testing.T) {
	resolver, _, _, _ := setupTestResolver()
	queryResolver := &queryResolver{resolver}

	// Test without authentication
	_, err := queryResolver.Me(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "authentication required")
}

func TestQueryResolver_Me_WithAuth(t *testing.T) {
	resolver, _, _, _ := setupTestResolver()
	queryResolver := &queryResolver{resolver}

	user := &model.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Name:  "Test User",
	}

	ctx := createAuthenticatedContext(user)
	result, err := queryResolver.Me(ctx)

	assert.NoError(t, err)
	assert.Equal(t, user, result)
}

func TestQueryResolver_Posts(t *testing.T) {
	resolver, _, mockPostRepo, _ := setupTestResolver()
	queryResolver := &queryResolver{resolver}

	expectedPosts := []*model.Post{
		{
			ID:        uuid.New(),
			Title:     "Test Post",
			Content:   "Test content",
			AuthorID:  uuid.New(),
			Tags:      []string{"test"},
			Published: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockPostRepo.On("List", mock.Anything, mock.AnythingOfType("*repository.PostFilters"), 20, 0).Return(expectedPosts, nil)
	mockPostRepo.On("Count", mock.Anything, mock.AnythingOfType("*repository.PostFilters")).Return(1, nil)

	result, err := queryResolver.Posts(context.Background(), nil, nil)

	assert.NoError(t, err)
	assert.Len(t, result.Edges, 1)
	assert.Equal(t, expectedPosts[0], result.Edges[0].Node)
	assert.Equal(t, 1, result.TotalCount)
	assert.False(t, result.PageInfo.HasNextPage)
	mockPostRepo.AssertExpectations(t)
}

func TestMutationResolver_CreatePost_RequiresAuth(t *testing.T) {
	resolver, _, _, _ := setupTestResolver()
	mutationResolver := &mutationResolver{resolver}

	input := model.CreatePostInput{
		Title:   "Test Post",
		Content: "Test content",
		Tags:    []string{"test"},
	}

	// Test without authentication
	_, err := mutationResolver.CreatePost(context.Background(), input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "authentication required")
}

func TestMutationResolver_CreatePost_WithAuth(t *testing.T) {
	resolver, _, mockPostRepo, _ := setupTestResolver()
	mutationResolver := &mutationResolver{resolver}

	user := &model.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Name:  "Test User",
	}

	input := model.CreatePostInput{
		Title:   "Test Post",
		Content: "Test content",
		Tags:    []string{"test"},
	}

	mockPostRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Post")).Return(nil)

	ctx := createAuthenticatedContext(user)
	result, err := mutationResolver.CreatePost(ctx, input)

	assert.NoError(t, err)
	assert.Equal(t, input.Title, result.Title)
	assert.Equal(t, input.Content, result.Content)
	assert.Equal(t, user.ID, result.AuthorID)
	assert.Equal(t, input.Tags, result.Tags)
	mockPostRepo.AssertExpectations(t)
}

func TestMutationResolver_Register(t *testing.T) {
	resolver, mockUserRepo, _, _ := setupTestResolver()
	mutationResolver := &mutationResolver{resolver}

	mockUserRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, assert.AnError)
	mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)

	result, err := mutationResolver.Register(context.Background(), "test@example.com", "password123", "Test User")

	assert.NoError(t, err)
	assert.NotEmpty(t, result.Token)
	assert.Equal(t, "test@example.com", result.User.Email)
	assert.Equal(t, "Test User", result.User.Name)
	mockUserRepo.AssertExpectations(t)
}

func TestPostResolver_Author(t *testing.T) {
	resolver, mockUserRepo, _, _ := setupTestResolver()
	postResolver := &postResolver{resolver}

	authorID := uuid.New()
	post := &model.Post{
		AuthorID: authorID,
	}

	expectedAuthor := &model.User{
		ID:    authorID,
		Email: "author@example.com",
		Name:  "Post Author",
	}

	mockUserRepo.On("GetByID", mock.Anything, authorID).Return(expectedAuthor, nil)

	author, err := postResolver.Author(context.Background(), post)

	assert.NoError(t, err)
	assert.Equal(t, expectedAuthor, author)
	mockUserRepo.AssertExpectations(t)
}