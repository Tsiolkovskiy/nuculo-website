package complexity

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/ast"
)

// Config holds complexity analysis configuration
type Config struct {
	MaxComplexity   int
	MaxDepth        int
	IntrospectionOk bool
}

// DefaultConfig returns a default complexity configuration
func DefaultConfig() Config {
	return Config{
		MaxComplexity:   1000,
		MaxDepth:        15,
		IntrospectionOk: true,
	}
}

// Analyzer provides query complexity analysis
type Analyzer struct {
	config Config
}

// NewAnalyzer creates a new complexity analyzer
func NewAnalyzer(config Config) *Analyzer {
	return &Analyzer{config: config}
}

// AnalyzeComplexity calculates the complexity of a GraphQL query
func (a *Analyzer) AnalyzeComplexity(ctx context.Context, rc *graphql.OperationContext) (int, error) {
	complexity := 0
	
	for _, selection := range rc.Operation.SelectionSet {
		fieldComplexity, err := a.calculateFieldComplexity(selection, rc.Variables, 1)
		if err != nil {
			return 0, err
		}
		complexity += fieldComplexity
	}
	
	if complexity > a.config.MaxComplexity {
		return complexity, fmt.Errorf("query complexity %d exceeds maximum allowed complexity %d", complexity, a.config.MaxComplexity)
	}
	
	return complexity, nil
}

// AnalyzeDepth calculates the depth of a GraphQL query
func (a *Analyzer) AnalyzeDepth(ctx context.Context, rc *graphql.OperationContext) (int, error) {
	depth := a.calculateMaxDepth(rc.Operation.SelectionSet, 1)
	
	if depth > a.config.MaxDepth {
		return depth, fmt.Errorf("query depth %d exceeds maximum allowed depth %d", depth, a.config.MaxDepth)
	}
	
	return depth, nil
}

// calculateFieldComplexity calculates complexity for a single field
func (a *Analyzer) calculateFieldComplexity(selection ast.Selection, variables map[string]interface{}, depth int) (int, error) {
	switch sel := selection.(type) {
	case *ast.Field:
		// Base complexity for each field
		complexity := 1
		
		// Add complexity based on arguments
		for _, arg := range sel.Arguments {
			if arg.Name == "limit" || arg.Name == "first" {
				if value := a.getArgumentValue(arg.Value, variables); value != nil {
					if limit, ok := value.(int); ok {
						complexity += limit / 10 // Add 1 complexity per 10 items
					}
				}
			}
		}
		
		// Add complexity for nested selections
		if sel.SelectionSet != nil {
			for _, nestedSelection := range sel.SelectionSet {
				nestedComplexity, err := a.calculateFieldComplexity(nestedSelection, variables, depth+1)
				if err != nil {
					return 0, err
				}
				complexity += nestedComplexity
			}
		}
		
		// Multiply by depth factor for deeply nested queries
		if depth > 5 {
			complexity *= depth - 4
		}
		
		return complexity, nil
		
	case *ast.InlineFragment:
		complexity := 0
		for _, nestedSelection := range sel.SelectionSet {
			nestedComplexity, err := a.calculateFieldComplexity(nestedSelection, variables, depth)
			if err != nil {
				return 0, err
			}
			complexity += nestedComplexity
		}
		return complexity, nil
		
	case *ast.FragmentSpread:
		// For fragment spreads, we'd need access to the fragment definition
		// For now, return a base complexity
		return 10, nil
		
	default:
		return 1, nil
	}
}

// calculateMaxDepth calculates the maximum depth of a selection set
func (a *Analyzer) calculateMaxDepth(selectionSet ast.SelectionSet, currentDepth int) int {
	maxDepth := currentDepth
	
	for _, selection := range selectionSet {
		switch sel := selection.(type) {
		case *ast.Field:
			if sel.SelectionSet != nil {
				fieldDepth := a.calculateMaxDepth(sel.SelectionSet, currentDepth+1)
				if fieldDepth > maxDepth {
					maxDepth = fieldDepth
				}
			}
		case *ast.InlineFragment:
			fragmentDepth := a.calculateMaxDepth(sel.SelectionSet, currentDepth)
			if fragmentDepth > maxDepth {
				maxDepth = fragmentDepth
			}
		}
	}
	
	return maxDepth
}

// getArgumentValue extracts the actual value from an argument
func (a *Analyzer) getArgumentValue(value *ast.Value, variables map[string]interface{}) interface{} {
	switch value.Kind {
	case ast.IntValue:
		return value.Raw
	case ast.FloatValue:
		return value.Raw
	case ast.StringValue:
		return value.Raw
	case ast.BooleanValue:
		return value.Raw == "true"
	case ast.Variable:
		if variables != nil {
			return variables[value.Raw]
		}
	}
	return nil
}

// ComplexityMiddleware creates a middleware for complexity analysis
func (a *Analyzer) ComplexityMiddleware() graphql.HandlerExtension {
	return &complexityExtension{analyzer: a}
}

// complexityExtension implements graphql.HandlerExtension
type complexityExtension struct {
	analyzer *Analyzer
}

func (e *complexityExtension) ExtensionName() string {
	return "ComplexityAnalysis"
}

func (e *complexityExtension) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

func (e *complexityExtension) InterceptOperation(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	return func(ctx context.Context) *graphql.Response {
		rc := graphql.GetOperationContext(ctx)
		
		// Skip introspection queries if configured
		if !e.analyzer.config.IntrospectionOk && rc.Operation.Name == "IntrospectionQuery" {
			return graphql.ErrorResponse(ctx, "introspection disabled")
		}
		
		// Analyze complexity
		complexity, err := e.analyzer.AnalyzeComplexity(ctx, rc)
		if err != nil {
			return graphql.ErrorResponse(ctx, err.Error())
		}
		
		// Analyze depth
		depth, err := e.analyzer.AnalyzeDepth(ctx, rc)
		if err != nil {
			return graphql.ErrorResponse(ctx, err.Error())
		}
		
		// Add complexity and depth to context for logging
		ctx = context.WithValue(ctx, "query_complexity", complexity)
		ctx = context.WithValue(ctx, "query_depth", depth)
		
		return next(ctx)
	}
}

func (e *complexityExtension) InterceptField(ctx context.Context, next graphql.Resolver) (interface{}, error) {
	return next(ctx)
}

func (e *complexityExtension) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	return next(ctx)
}