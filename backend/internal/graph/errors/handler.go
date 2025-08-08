package errors

import (
	"context"
	"fmt"
	"strings"

	"backend/internal/logging"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// ErrorHandler handles GraphQL errors and provides consistent error formatting
type ErrorHandler struct {
	logger *logging.Logger
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(logger *logging.Logger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

// HandleError processes an error and returns a properly formatted GraphQL error
func (h *ErrorHandler) HandleError(ctx context.Context, err error) *gqlerror.Error {
	if err == nil {
		return nil
	}

	// If it's already a GraphQLError, convert it
	if gqlErr, ok := err.(*GraphQLError); ok {
		h.logError(ctx, gqlErr)
		return gqlErr.ToGQLError()
	}

	// If it's already a gqlerror.Error, return as is
	if gqlErr, ok := err.(*gqlerror.Error); ok {
		h.logGQLError(ctx, gqlErr)
		return gqlErr
	}

	// Handle common error types
	gqlErr := h.categorizeError(err)
	h.logError(ctx, gqlErr)
	return gqlErr.ToGQLError()
}

// categorizeError categorizes common errors into GraphQL error types
func (h *ErrorHandler) categorizeError(err error) *GraphQLError {
	errMsg := err.Error()
	errMsgLower := strings.ToLower(errMsg)

	// Authentication errors
	if strings.Contains(errMsgLower, "authentication required") ||
		strings.Contains(errMsgLower, "unauthenticated") ||
		strings.Contains(errMsgLower, "invalid token") ||
		strings.Contains(errMsgLower, "token expired") {
		return NewUnauthenticatedError("Authentication required")
	}

	// Authorization errors
	if strings.Contains(errMsgLower, "unauthorized") ||
		strings.Contains(errMsgLower, "forbidden") ||
		strings.Contains(errMsgLower, "access denied") ||
		strings.Contains(errMsgLower, "permission denied") {
		return NewUnauthorizedError("Access denied")
	}

	// Not found errors
	if strings.Contains(errMsgLower, "not found") ||
		strings.Contains(errMsgLower, "does not exist") {
		return NewNotFoundError("Resource")
	}

	// Already exists errors
	if strings.Contains(errMsgLower, "already exists") ||
		strings.Contains(errMsgLower, "duplicate") ||
		strings.Contains(errMsgLower, "unique constraint") {
		return NewAlreadyExistsError("Resource")
	}

	// Validation errors
	if strings.Contains(errMsgLower, "validation") ||
		strings.Contains(errMsgLower, "invalid") ||
		strings.Contains(errMsgLower, "required") ||
		strings.Contains(errMsgLower, "format") {
		return NewValidationError(errMsg, "")
	}

	// Database errors
	if strings.Contains(errMsgLower, "database") ||
		strings.Contains(errMsgLower, "sql") ||
		strings.Contains(errMsgLower, "connection") {
		return NewDatabaseError("Database operation failed")
	}

	// Default to internal error
	return NewInternalError("An unexpected error occurred")
}

// logError logs a GraphQLError
func (h *ErrorHandler) logError(ctx context.Context, err *GraphQLError) {
	if h.logger == nil {
		return
	}

	// Extract request ID from context if available
	requestID := h.getRequestID(ctx)

	// Log based on error severity
	switch err.Code {
	case ErrorCodeInternal, ErrorCodeDatabaseError, ErrorCodeNetworkError:
		h.logger.Printf("ERROR [%s] %s: %s (field: %s)", requestID, err.Code, err.Message, err.Field)
	case ErrorCodeUnauthenticated, ErrorCodeUnauthorized, ErrorCodeForbidden:
		h.logger.Printf("WARN [%s] %s: %s", requestID, err.Code, err.Message)
	default:
		h.logger.Printf("INFO [%s] %s: %s (field: %s)", requestID, err.Code, err.Message, err.Field)
	}
}

// logGQLError logs a gqlerror.Error
func (h *ErrorHandler) logGQLError(ctx context.Context, err *gqlerror.Error) {
	if h.logger == nil {
		return
	}

	requestID := h.getRequestID(ctx)
	code := "UNKNOWN"
	
	if err.Extensions != nil {
		if c, ok := err.Extensions["code"].(string); ok {
			code = c
		}
	}

	h.logger.Printf("INFO [%s] %s: %s", requestID, code, err.Message)
}

// getRequestID extracts request ID from context
func (h *ErrorHandler) getRequestID(ctx context.Context) string {
	if ctx == nil {
		return "unknown"
	}
	
	// Try to get request ID from context
	if id := ctx.Value("request_id"); id != nil {
		if idStr, ok := id.(string); ok {
			return idStr
		}
	}
	
	return "unknown"
}

// WrapDatabaseError wraps database errors with appropriate GraphQL error
func WrapDatabaseError(err error, operation string) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()
	errMsgLower := strings.ToLower(errMsg)

	// Handle specific database errors
	if strings.Contains(errMsgLower, "unique constraint") ||
		strings.Contains(errMsgLower, "duplicate key") {
		return NewAlreadyExistsError("Resource")
	}

	if strings.Contains(errMsgLower, "foreign key constraint") {
		return NewValidationError("Referenced resource does not exist", "")
	}

	if strings.Contains(errMsgLower, "not null constraint") {
		return NewValidationError("Required field is missing", "")
	}

	// Generic database error
	return NewDatabaseError(fmt.Sprintf("Database %s failed", operation))
}

// WrapValidationErrors wraps multiple validation errors
func WrapValidationErrors(errors []error) error {
	if len(errors) == 0 {
		return nil
	}

	if len(errors) == 1 {
		return errors[0]
	}

	// Combine multiple validation errors
	messages := make([]string, len(errors))
	for i, err := range errors {
		messages[i] = err.Error()
	}

	return NewValidationError(
		fmt.Sprintf("Multiple validation errors: %s", strings.Join(messages, "; ")),
		"",
	)
}

// IsGraphQLError checks if an error is a GraphQLError
func IsGraphQLError(err error) bool {
	_, ok := err.(*GraphQLError)
	return ok
}

// GetErrorCode extracts the error code from an error
func GetErrorCode(err error) ErrorCode {
	if gqlErr, ok := err.(*GraphQLError); ok {
		return gqlErr.Code
	}
	return ErrorCodeInternal
}