package validation

import (
	"testing"

	"backend/internal/graph/errors"
	"backend/internal/graph/model"
	"github.com/stretchr/testify/assert"
)

func TestValidator_ValidateTitle(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name        string
		title       string
		expectError bool
		errorCode   errors.ErrorCode
	}{
		{
			name:        "Valid title",
			title:       "This is a valid title",
			expectError: false,
		},
		{
			name:        "Empty title",
			title:       "",
			expectError: true,
			errorCode:   errors.ErrorCodeValidation,
		},
		{
			name:        "Title too short",
			title:       "Hi",
			expectError: true,
			errorCode:   errors.ErrorCodeValidation,
		},
		{
			name:        "Title too long",
			title:       string(make([]rune, 201)),
			expectError: true,
			errorCode:   errors.ErrorCodeValidation,
		},
		{
			name:        "Title with invalid characters",
			title:       "Title with <script>",
			expectError: true,
			errorCode:   errors.ErrorCodeValidation,
		},
		{
			name:        "Title with whitespace only",
			title:       "   ",
			expectError: true,
			errorCode:   errors.ErrorCodeValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateTitle(tt.title)
			
			if tt.expectError {
				assert.Error(t, err)
				if gqlErr, ok := err.(*errors.GraphQLError); ok {
					assert.Equal(t, tt.errorCode, gqlErr.Code)
					assert.Equal(t, "title", gqlErr.Field)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateEmail(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name        string
		email       string
		expectError bool
	}{
		{
			name:        "Valid email",
			email:       "user@example.com",
			expectError: false,
		},
		{
			name:        "Valid email with subdomain",
			email:       "user@mail.example.com",
			expectError: false,
		},
		{
			name:        "Empty email",
			email:       "",
			expectError: true,
		},
		{
			name:        "Invalid email format",
			email:       "invalid-email",
			expectError: true,
		},
		{
			name:        "Email without domain",
			email:       "user@",
			expectError: true,
		},
		{
			name:        "Email without @",
			email:       "userexample.com",
			expectError: true,
		},
		{
			name:        "Email with spaces",
			email:       "user @example.com",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateEmail(tt.email)
			
			if tt.expectError {
				assert.Error(t, err)
				if gqlErr, ok := err.(*errors.GraphQLError); ok {
					assert.Equal(t, errors.ErrorCodeValidation, gqlErr.Code)
					assert.Equal(t, "email", gqlErr.Field)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidatePassword(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name        string
		password    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid password",
			password:    "password123",
			expectError: false,
		},
		{
			name:        "Valid complex password",
			password:    "MySecureP@ssw0rd",
			expectError: false,
		},
		{
			name:        "Password too short",
			password:    "pass1",
			expectError: true,
			errorMsg:    "at least 8 characters",
		},
		{
			name:        "Password too long",
			password:    string(make([]byte, 129)),
			expectError: true,
			errorMsg:    "cannot exceed 128 characters",
		},
		{
			name:        "Password without letters",
			password:    "12345678",
			expectError: true,
			errorMsg:    "at least one letter",
		},
		{
			name:        "Password without numbers",
			password:    "password",
			expectError: true,
			errorMsg:    "at least one number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidatePassword(tt.password)
			
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				if gqlErr, ok := err.(*errors.GraphQLError); ok {
					assert.Equal(t, errors.ErrorCodeValidation, gqlErr.Code)
					assert.Equal(t, "password", gqlErr.Field)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateTags(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name        string
		tags        []string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid tags",
			tags:        []string{"golang", "graphql", "api"},
			expectError: false,
		},
		{
			name:        "Empty tags array",
			tags:        []string{},
			expectError: false,
		},
		{
			name:        "Too many tags",
			tags:        []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"},
			expectError: true,
			errorMsg:    "Cannot have more than 10 tags",
		},
		{
			name:        "Empty tag",
			tags:        []string{"valid", ""},
			expectError: true,
			errorMsg:    "cannot be empty",
		},
		{
			name:        "Tag too short",
			tags:        []string{"a"},
			expectError: true,
			errorMsg:    "at least 2 characters",
		},
		{
			name:        "Tag too long",
			tags:        []string{string(make([]rune, 31))},
			expectError: true,
			errorMsg:    "cannot exceed 30 characters",
		},
		{
			name:        "Tag with invalid characters",
			tags:        []string{"tag with spaces"},
			expectError: true,
			errorMsg:    "invalid characters",
		},
		{
			name:        "Duplicate tags",
			tags:        []string{"golang", "golang"},
			expectError: true,
			errorMsg:    "Duplicate tag",
		},
		{
			name:        "Valid tags with hyphens and underscores",
			tags:        []string{"go-lang", "graph_ql", "rest-api"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateTags(tt.tags)
			
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				if gqlErr, ok := err.(*errors.GraphQLError); ok {
					assert.Equal(t, errors.ErrorCodeValidation, gqlErr.Code)
					assert.Equal(t, "tags", gqlErr.Field)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateCreatePostInput(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name        string
		input       model.CreatePostInput
		expectError bool
	}{
		{
			name: "Valid input",
			input: model.CreatePostInput{
				Title:   "Valid Post Title",
				Content: "This is valid post content with enough characters.",
				Tags:    []string{"golang", "graphql"},
			},
			expectError: false,
		},
		{
			name: "Invalid title",
			input: model.CreatePostInput{
				Title:   "",
				Content: "Valid content here",
				Tags:    []string{"golang"},
			},
			expectError: true,
		},
		{
			name: "Invalid content",
			input: model.CreatePostInput{
				Title:   "Valid Title",
				Content: "Short",
				Tags:    []string{"golang"},
			},
			expectError: true,
		},
		{
			name: "Invalid tags",
			input: model.CreatePostInput{
				Title:   "Valid Title",
				Content: "Valid content with enough characters",
				Tags:    []string{"invalid tag with spaces"},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCreatePostInput(tt.input)
			
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidatePaginationInput(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name        string
		input       *model.PaginationInput
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Nil input",
			input:       nil,
			expectError: false,
		},
		{
			name: "Valid pagination",
			input: &model.PaginationInput{
				Page:  intPtr(1),
				Limit: intPtr(20),
			},
			expectError: false,
		},
		{
			name: "Invalid page (too low)",
			input: &model.PaginationInput{
				Page: intPtr(0),
			},
			expectError: true,
			errorMsg:    "Page must be at least 1",
		},
		{
			name: "Invalid page (too high)",
			input: &model.PaginationInput{
				Page: intPtr(1001),
			},
			expectError: true,
			errorMsg:    "Page cannot exceed 1000",
		},
		{
			name: "Invalid limit (too low)",
			input: &model.PaginationInput{
				Limit: intPtr(0),
			},
			expectError: true,
			errorMsg:    "Limit must be at least 1",
		},
		{
			name: "Invalid limit (too high)",
			input: &model.PaginationInput{
				Limit: intPtr(101),
			},
			expectError: true,
			errorMsg:    "Limit cannot exceed 100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidatePaginationInput(tt.input)
			
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Helper function to create int pointers
func intPtr(i int) *int {
	return &i
}