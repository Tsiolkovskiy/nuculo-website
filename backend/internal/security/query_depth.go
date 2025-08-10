package security

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/ast"
)

// QueryDepthLimiter limits the depth of GraphQL queries to prevent abuse
type QueryDepthLimiter struct {
	maxDepth int
}

// NewQueryDepthLimiter creates a new query depth limiter
func NewQueryDepthLimiter(maxDepth int) *QueryDepthLimiter {
	return &QueryDepthLimiter{
		maxDepth: maxDepth,
	}
}

// ExtensionName returns the name of this extension
func (q *QueryDepthLimiter) ExtensionName() string {
	return "QueryDepthLimiter"
}

// Validate validates the schema (no-op for this extension)
func (q *QueryDepthLimiter) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

// InterceptOperation intercepts operations to check query depth
func (q *QueryDepthLimiter) InterceptOperation(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	oc := graphql.GetOperationContext(ctx)
	
	// Calculate query depth
	depth := q.calculateDepth(oc.Operation.SelectionSet, 0)
	
	if depth > q.maxDepth {
		return func(ctx context.Context) *graphql.Response {
			return &graphql.Response{
				Errors: []*graphql.Error{
					{
						Message: fmt.Sprintf("Query depth %d exceeds maximum allowed depth %d", depth, q.maxDepth),
						Extensions: map[string]interface{}{
							"code": "QUERY_TOO_DEEP",
							"maxDepth": q.maxDepth,
							"actualDepth": depth,
						},
					},
				},
			}
		}
	}
	
	return next(ctx)
}

// calculateDepth recursively calculates the depth of a selection set
func (q *QueryDepthLimiter) calculateDepth(selectionSet ast.SelectionSet, currentDepth int) int {
	if len(selectionSet) == 0 {
		return currentDepth
	}
	
	maxDepth := currentDepth
	
	for _, selection := range selectionSet {
		switch sel := selection.(type) {
		case *ast.Field:
			if sel.SelectionSet != nil {
				depth := q.calculateDepth(sel.SelectionSet, currentDepth+1)
				if depth > maxDepth {
					maxDepth = depth
				}
			}
		case *ast.InlineFragment:
			depth := q.calculateDepth(sel.SelectionSet, currentDepth)
			if depth > maxDepth {
				maxDepth = depth
			}
		case *ast.FragmentSpread:
			// For fragment spreads, we would need access to the document
			// to resolve the fragment definition. For simplicity, we'll
			// assume a depth of 1 for fragments.
			depth := currentDepth + 1
			if depth > maxDepth {
				maxDepth = depth
			}
		}
	}
	
	return maxDepth
}

// QueryComplexityAnalyzer analyzes query complexity to prevent expensive operations
type QueryComplexityAnalyzer struct {
	maxComplexity int
	fieldWeights  map[string]int
}

// NewQueryComplexityAnalyzer creates a new query complexity analyzer
func NewQueryComplexityAnalyzer(maxComplexity int) *QueryComplexityAnalyzer {
	return &QueryComplexityAnalyzer{
		maxComplexity: maxComplexity,
		fieldWeights: map[string]int{
			// Default field weights
			"posts":     5,  // List queries are more expensive
			"users":     5,
			"comments":  3,
			"post":      1,  // Single item queries are cheaper
			"user":      1,
			"comment":   1,
			"createPost": 10, // Mutations are expensive
			"updatePost": 8,
			"deletePost": 5,
			"login":     3,
			"register":  5,
		},
	}
}

// SetFieldWeight sets the complexity weight for a specific field
func (q *QueryComplexityAnalyzer) SetFieldWeight(field string, weight int) {
	q.fieldWeights[field] = weight
}

// ExtensionName returns the name of this extension
func (q *QueryComplexityAnalyzer) ExtensionName() string {
	return "QueryComplexityAnalyzer"
}

// Validate validates the schema (no-op for this extension)
func (q *QueryComplexityAnalyzer) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

// InterceptOperation intercepts operations to check query complexity
func (q *QueryComplexityAnalyzer) InterceptOperation(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	oc := graphql.GetOperationContext(ctx)
	
	// Calculate query complexity
	complexity := q.calculateComplexity(oc.Operation.SelectionSet, 1)
	
	if complexity > q.maxComplexity {
		return func(ctx context.Context) *graphql.Response {
			return &graphql.Response{
				Errors: []*graphql.Error{
					{
						Message: fmt.Sprintf("Query complexity %d exceeds maximum allowed complexity %d", complexity, q.maxComplexity),
						Extensions: map[string]interface{}{
							"code": "QUERY_TOO_COMPLEX",
							"maxComplexity": q.maxComplexity,
							"actualComplexity": complexity,
						},
					},
				},
			}
		}
	}
	
	return next(ctx)
}

// calculateComplexity recursively calculates the complexity of a selection set
func (q *QueryComplexityAnalyzer) calculateComplexity(selectionSet ast.SelectionSet, multiplier int) int {
	if len(selectionSet) == 0 {
		return 0
	}
	
	totalComplexity := 0
	
	for _, selection := range selectionSet {
		switch sel := selection.(type) {
		case *ast.Field:
			fieldWeight := q.getFieldWeight(sel.Name)
			fieldComplexity := fieldWeight * multiplier
			
			// Add complexity for nested selections
			if sel.SelectionSet != nil {
				// For list fields, assume a multiplier based on potential result size
				nestedMultiplier := multiplier
				if q.isListField(sel.Name) {
					nestedMultiplier = multiplier * 10 // Assume up to 10 items in lists
				}
				fieldComplexity += q.calculateComplexity(sel.SelectionSet, nestedMultiplier)
			}
			
			totalComplexity += fieldComplexity
			
		case *ast.InlineFragment:
			totalComplexity += q.calculateComplexity(sel.SelectionSet, multiplier)
			
		case *ast.FragmentSpread:
			// For fragment spreads, assume a base complexity
			totalComplexity += 5 * multiplier
		}
	}
	
	return totalComplexity
}

// getFieldWeight returns the complexity weight for a field
func (q *QueryComplexityAnalyzer) getFieldWeight(fieldName string) int {
	if weight, exists := q.fieldWeights[fieldName]; exists {
		return weight
	}
	return 1 // Default weight
}

// isListField checks if a field typically returns a list
func (q *QueryComplexityAnalyzer) isListField(fieldName string) bool {
	listFields := map[string]bool{
		"posts":    true,
		"users":    true,
		"comments": true,
	}
	return listFields[fieldName]
}

// SecurityConfig holds security configuration
type SecurityConfig struct {
	MaxQueryDepth      int
	MaxQueryComplexity int
	EnableDepthLimit   bool
	EnableComplexity   bool
}

// DefaultSecurityConfig returns default security configuration
func DefaultSecurityConfig() SecurityConfig {
	return SecurityConfig{
		MaxQueryDepth:      10,
		MaxQueryComplexity: 1000,
		EnableDepthLimit:   true,
		EnableComplexity:   true,
	}
}

// CreateSecurityExtensions creates security extensions based on config
func CreateSecurityExtensions(config SecurityConfig) []graphql.HandlerExtension {
	var extensions []graphql.HandlerExtension
	
	if config.EnableDepthLimit {
		extensions = append(extensions, NewQueryDepthLimiter(config.MaxQueryDepth))
	}
	
	if config.EnableComplexity {
		extensions = append(extensions, NewQueryComplexityAnalyzer(config.MaxQueryComplexity))
	}
	
	return extensions
}