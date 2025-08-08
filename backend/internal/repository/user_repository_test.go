package repository

import (
	"context"
	"testing"
	"time"

	"backend/internal/graph/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDB is a mock database for testing
type MockDB struct {
	mock.Mock
}

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*model.User), args.Error(1)
}

func TestUserRepository_Create(t *testing.T) {
	mockRepo := new(MockUserRepository)
	ctx := context.Background()
	
	user := &model.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: "hashed_password",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	
	mockRepo.On("Create", ctx, user).Return(nil)
	
	err := mockRepo.Create(ctx, user)
	
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserRepository_GetByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	ctx := context.Background()
	userID := uuid.New()
	
	expectedUser := &model.User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}
	
	mockRepo.On("GetByID", ctx, userID).Return(expectedUser, nil)
	
	user, err := mockRepo.GetByID(ctx, userID)
	
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestUserRepository_GetByEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	ctx := context.Background()
	email := "test@example.com"
	
	expectedUser := &model.User{
		ID:    uuid.New(),
		Email: email,
		Name:  "Test User",
	}
	
	mockRepo.On("GetByEmail", ctx, email).Return(expectedUser, nil)
	
	user, err := mockRepo.GetByEmail(ctx, email)
	
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestUserRepository_Update(t *testing.T) {
	mockRepo := new(MockUserRepository)
	ctx := context.Background()
	
	user := &model.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		Name:      "Updated User",
		UpdatedAt: time.Now(),
	}
	
	mockRepo.On("Update", ctx, user).Return(nil)
	
	err := mockRepo.Update(ctx, user)
	
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserRepository_Delete(t *testing.T) {
	mockRepo := new(MockUserRepository)
	ctx := context.Background()
	userID := uuid.New()
	
	mockRepo.On("Delete", ctx, userID).Return(nil)
	
	err := mockRepo.Delete(ctx, userID)
	
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserRepository_List(t *testing.T) {
	mockRepo := new(MockUserRepository)
	ctx := context.Background()
	limit, offset := 10, 0
	
	expectedUsers := []*model.User{
		{
			ID:    uuid.New(),
			Email: "user1@example.com",
			Name:  "User 1",
		},
		{
			ID:    uuid.New(),
			Email: "user2@example.com",
			Name:  "User 2",
		},
	}
	
	mockRepo.On("List", ctx, limit, offset).Return(expectedUsers, nil)
	
	users, err := mockRepo.List(ctx, limit, offset)
	
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
	assert.Len(t, users, 2)
	mockRepo.AssertExpectations(t)
}