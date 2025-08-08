package errors

import (
	"fmt"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// ErrorCode represents different types of GraphQL errors
type ErrorCode string

const (
	// Validation errors
	ErrorCodeValidation     ErrorCode = "VALIDATION_ERROR"
	ErrorCodeInvalidInput   ErrorCode = "INVALID_INPUT"
	ErrorCodeInvalidFormat  ErrorCode = "INVALID_FORMAT"
	
	// Authentication and authorization errors
	ErrorCodeUnauthenticated ErrorCode = "UNAUTHENTICATED"
	ErrorCodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrorCodeForbidden       ErrorCode = "FORBIDDEN"
	
	// Resource errors
	ErrorCodeNotFound       ErrorCode = "NOT_FOUND"
	ErrorCodeAlreadyExists  ErrorCode = "ALREADY_EXISTS"
	ErrorCodeConflict       ErrorCode = "CONFLICT"
	
	// System errors
	ErrorCodeInternal       ErrorCode = "INTERNAL_ERROR"
	ErrorCodeDatabaseError  ErrorCode = "DATABASE_ERROR"
	ErrorCodeNetworkError   ErrorCode = "NETWORK_ERROR"
	
	// Rate limiting
	ErrorCodeRateLimit      ErrorCode = "RATE_LIMIT_EXCEEDED"
)

// GraphQLError represents a structured GraphQL error
type GraphQLError struct {
	Message    string                 `json:"message"`
	Code       ErrorCode              `json:"code"`
	Field      string                 `json:"field,omitempty"`
	Path       []string               `json:"path,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

// Error implements the error interface
func (e *GraphQLError) Error() string {
	return e.Message
}

// ToGQLError converts to gqlerror.Error
func (e *GraphQLError) ToGQLError() *gqlerror.Error {
	extensions := make(map[string]interface{})
	extensions["code"] = string(e.Code)
	
	if e.Field != "" {
		extensions["field"] = e.Field
	}
	
	// Add any additional extensions
	for k, v := range e.Extensions {
		extensions[k] = v
	}
	
	// Convert path to ast.Path format
	var astPath ast.Path
	for _, p := range e.Path {
		astPath = append(astPath, ast.PathName(p))
	}
	
	return &gqlerror.Error{
		Message:    e.Message,
		Path:       astPath,
		Extensions: extensions,
	}
}

// NewValidationError creates a validation error
func NewValidationError(message, field string) *GraphQLError {
	return &GraphQLError{
		Message: message,
		Code:    ErrorCodeValidation,
		Field:   field,
	}
}

// NewInvalidInputError creates an invalid input error
func NewInvalidInputError(message, field string) *GraphQLError {
	return &GraphQLError{
		Message: message,
		Code:    ErrorCodeInvalidInput,
		Field:   field,
	}
}

// NewInvalidFormatError creates an invalid format error
func NewInvalidFormatError(message, field string) *GraphQLError {
	return &GraphQLError{
		Message: message,
		Code:    ErrorCodeInvalidFormat,
		Field:   field,
	}
}

// NewUnauthenticatedError creates an authentication error
func NewUnauthenticatedError(message string) *GraphQLError {
	return &GraphQLError{
		Message: message,
		Code:    ErrorCodeUnauthenticated,
	}
}

// NewUnauthorizedError creates an authorization error
func NewUnauthorizedError(message string) *GraphQLError {
	return &GraphQLError{
		Message: message,
		Code:    ErrorCodeUnauthorized,
	}
}

// NewForbiddenError creates a forbidden error
func NewForbiddenError(message string) *GraphQLError {
	return &GraphQLError{
		Message: message,
		Code:    ErrorCodeForbidden,
	}
}

// NewNotFoundError creates a not found error
func NewNotFoundError(resource string) *GraphQLError {
	return &GraphQLError{
		Message: fmt.Sprintf("%s not found", resource),
		Code:    ErrorCodeNotFound,
	}
}

// NewAlreadyExistsError creates an already exists error
func NewAlreadyExistsError(resource string) *GraphQLError {
	return &GraphQLError{
		Message: fmt.Sprintf("%s already exists", resource),
		Code:    ErrorCodeAlreadyExists,
	}
}

// NewConflictError creates a conflict error
func NewConflictError(message string) *GraphQLError {
	return &GraphQLError{
		Message: message,
		Code:    ErrorCodeConflict,
	}
}

// NewInternalError creates an internal error
func NewInternalError(message string) *GraphQLError {
	return &GraphQLError{
		Message: message,
		Code:    ErrorCodeInternal,
	}
}

// NewDatabaseError creates a database error
func NewDatabaseError(message string) *GraphQLError {
	return &GraphQLError{
		Message: message,
		Code:    ErrorCodeDatabaseError,
	}
}

// NewRateLimitError creates a rate limit error
func NewRateLimitError(message string) *GraphQLError {
	return &GraphQLError{
		Message: message,
		Code:    ErrorCodeRateLimit,
	}
}