package scalars

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
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
			name:        "Valid email with numbers",
			email:       "user123@example.com",
			expectError: false,
		},
		{
			name:        "Empty email",
			email:       "",
			expectError: true,
		},
		{
			name:        "Email too long",
			email:       string(make([]byte, 255)) + "@example.com",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUnmarshalEmail(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    string
		expectError bool
	}{
		{
			name:        "Valid email string",
			input:       "user@example.com",
			expected:    "user@example.com",
			expectError: false,
		},
		{
			name:        "Valid email string pointer",
			input:       stringPtr("user@example.com"),
			expected:    "user@example.com",
			expectError: false,
		},
		{
			name:        "Invalid email string",
			input:       "invalid-email",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Nil string pointer",
			input:       (*string)(nil),
			expected:    "",
			expectError: true,
		},
		{
			name:        "Invalid type",
			input:       123,
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := UnmarshalEmail(tt.input)
			
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestValidateDateTime(t *testing.T) {
	tests := []struct {
		name        string
		datetime    string
		expectError bool
	}{
		{
			name:        "Valid RFC3339 datetime",
			datetime:    "2023-01-01T12:00:00Z",
			expectError: false,
		},
		{
			name:        "Valid RFC3339 with timezone",
			datetime:    "2023-01-01T12:00:00+02:00",
			expectError: false,
		},
		{
			name:        "Valid date only",
			datetime:    "2023-01-01",
			expectError: false,
		},
		{
			name:        "Valid datetime without timezone",
			datetime:    "2023-01-01T12:00:00",
			expectError: false,
		},
		{
			name:        "Empty datetime",
			datetime:    "",
			expectError: true,
		},
		{
			name:        "Invalid datetime format",
			datetime:    "invalid-datetime",
			expectError: true,
		},
		{
			name:        "Invalid date format",
			datetime:    "2023-13-01",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDateTime(tt.datetime)
			
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUnmarshalDateTime(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectError bool
	}{
		{
			name:        "Valid RFC3339 string",
			input:       "2023-01-01T12:00:00Z",
			expectError: false,
		},
		{
			name:        "Valid string pointer",
			input:       stringPtr("2023-01-01T12:00:00Z"),
			expectError: false,
		},
		{
			name:        "Valid Unix timestamp int",
			input:       1672574400,
			expectError: false,
		},
		{
			name:        "Valid Unix timestamp int64",
			input:       int64(1672574400),
			expectError: false,
		},
		{
			name:        "Valid Unix timestamp float64",
			input:       float64(1672574400),
			expectError: false,
		},
		{
			name:        "Invalid string",
			input:       "invalid-datetime",
			expectError: true,
		},
		{
			name:        "Nil string pointer",
			input:       (*string)(nil),
			expectError: true,
		},
		{
			name:        "Invalid type",
			input:       []string{"invalid"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := UnmarshalDateTime(tt.input)
			
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateDateTimeRange(t *testing.T) {
	now := time.Now()
	
	tests := []struct {
		name        string
		datetime    time.Time
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid current time",
			datetime:    now,
			expectError: false,
		},
		{
			name:        "Valid past time",
			datetime:    now.AddDate(-10, 0, 0),
			expectError: false,
		},
		{
			name:        "Valid future time",
			datetime:    now.AddDate(10, 0, 0),
			expectError: false,
		},
		{
			name:        "Too far in the past",
			datetime:    time.Date(1800, 1, 1, 0, 0, 0, 0, time.UTC),
			expectError: true,
			errorMsg:    "too far in the past",
		},
		{
			name:        "Too far in the future",
			datetime:    now.AddDate(150, 0, 0),
			expectError: true,
			errorMsg:    "too far in the future",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDateTimeRange(tt.datetime)
			
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

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}