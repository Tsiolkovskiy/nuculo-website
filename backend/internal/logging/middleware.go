package logging

import (
	"context"
	"fmt"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// GraphQLMiddleware creates a GraphQL logging middleware
func GraphQLMiddleware(logger *Logger) graphql.HandlerExtension {
	return &graphqlLogger{logger: logger}
}

type graphqlLogger struct {
	logger *Logger
}

func (g *graphqlLogger) ExtensionName() string {
	return "GraphQLLogger"
}

func (g *graphqlLogger) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

// InterceptOperation logs GraphQL operations
func (g *graphqlLogger) InterceptOperation(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	oc := graphql.GetOperationContext(ctx)
	start := time.Now()

	// Extract operation details
	operationName := "unknown"
	operationType := "unknown"
	if oc.Operation != nil {
		if oc.Operation.Name != "" {
			operationName = oc.Operation.Name
		}
		operationType = string(oc.Operation.Operation)
	}

	// Add operation context
	ctx = WithOperationName(ctx, operationName)
	logger := g.logger.WithContext(ctx)

	// Log operation start
	logger.Info("GraphQL operation started",
		"operation_name", operationName,
		"operation_type", operationType,
		"query", oc.RawQuery,
		"variables", oc.Variables,
	)

	return next(ctx)
}

// InterceptResponse logs GraphQL responses
func (g *graphqlLogger) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	start := time.Now()
	oc := graphql.GetOperationContext(ctx)
	logger := g.logger.WithContext(ctx)

	// Execute the operation
	resp := next(ctx)
	duration := time.Since(start)

	// Extract operation details
	operationName := "unknown"
	operationType := "unknown"
	if oc.Operation != nil {
		if oc.Operation.Name != "" {
			operationName = oc.Operation.Name
		}
		operationType = string(oc.Operation.Operation)
	}

	// Log based on response status
	if len(resp.Errors) > 0 {
		// Log errors
		for _, err := range resp.Errors {
			logger.Error("GraphQL operation error",
				"operation_name", operationName,
				"operation_type", operationType,
				"duration_ms", duration.Milliseconds(),
				"error_message", err.Message,
				"error_path", err.Path,
				"error_locations", err.Locations,
				"error_extensions", err.Extensions,
			)
		}
	} else {
		// Log successful operation
		logger.Info("GraphQL operation completed",
			"operation_name", operationName,
			"operation_type", operationType,
			"duration_ms", duration.Milliseconds(),
			"has_data", resp.Data != nil,
		)
	}

	return resp
}

// InterceptField logs field resolution (optional, can be verbose)
func (g *graphqlLogger) InterceptField(ctx context.Context, next graphql.Resolver) (interface{}, error) {
	fc := graphql.GetFieldContext(ctx)
	logger := g.logger.WithContext(ctx)

	// Only log slow fields or errors
	start := time.Now()
	res, err := next(ctx)
	duration := time.Since(start)

	// Log slow fields (> 100ms) or errors
	if duration > 100*time.Millisecond || err != nil {
		if err != nil {
			logger.Error("Field resolution error",
				"field", fc.Field.Name,
				"path", fc.Path(),
				"duration_ms", duration.Milliseconds(),
				"error", err.Error(),
			)
		} else {
			logger.Warn("Slow field resolution",
				"field", fc.Field.Name,
				"path", fc.Path(),
				"duration_ms", duration.Milliseconds(),
			)
		}
	}

	return res, err
}

// GinMiddleware creates a Gin logging middleware
func GinMiddleware(logger *Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := GenerateRequestID()

		// Add request ID to context
		ctx := WithRequestID(c.Request.Context(), requestID)
		ctx = WithLogger(ctx, logger)
		c.Request = c.Request.WithContext(ctx)

		// Add request ID to response headers
		c.Header("X-Request-ID", requestID)

		// Log request start
		logger.WithContext(ctx).Info("HTTP request started",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"query", c.Request.URL.RawQuery,
			"user_agent", c.Request.UserAgent(),
			"remote_addr", c.ClientIP(),
		)

		// Process request
		c.Next()

		// Log request completion
		duration := time.Since(start)
		status := c.Writer.Status()

		logLevel := "info"
		if status >= 400 {
			logLevel = "error"
		} else if status >= 300 {
			logLevel = "warn"
		}

		logEntry := logger.WithContext(ctx).WithFields(map[string]interface{}{
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"status":      status,
			"duration_ms": duration.Milliseconds(),
			"size":        c.Writer.Size(),
		})

		message := fmt.Sprintf("HTTP request completed - %d %s", status, c.Request.Method)

		switch logLevel {
		case "error":
			logEntry.Error(message)
		case "warn":
			logEntry.Warn(message)
		default:
			logEntry.Info(message)
		}
	}
}

// ErrorLogger logs GraphQL errors with context
func ErrorLogger(logger *Logger) graphql.ErrorPresenterFunc {
	return func(ctx context.Context, err error) *gqlerror.Error {
		logger := logger.WithContext(ctx)

		// Convert to GraphQL error if needed
		var gqlErr *gqlerror.Error
		if e, ok := err.(*gqlerror.Error); ok {
			gqlErr = e
		} else {
			gqlErr = &gqlerror.Error{
				Message: err.Error(),
			}
		}

		// Log the error with context
		logger.Error("GraphQL error occurred",
			"error_message", gqlErr.Message,
			"error_path", gqlErr.Path,
			"error_locations", gqlErr.Locations,
			"error_extensions", gqlErr.Extensions,
		)

		return gqlErr
	}
}

// RecoveryLogger logs panics and recovers
func RecoveryLogger(logger *Logger) graphql.RecoverFunc {
	return func(ctx context.Context, err interface{}) error {
		logger := logger.WithContext(ctx)

		logger.Error("GraphQL panic recovered",
			"panic", fmt.Sprintf("%v", err),
		)

		return fmt.Errorf("internal server error")
	}
}