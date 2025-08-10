package playground

import (
	"html/template"
	"net/http"
)

// Config holds playground configuration
type Config struct {
	GraphQLEndpoint      string
	SubscriptionEndpoint string
	Title                string
	Version              string
	EnableInProduction   bool
	EnableIntrospection  bool
	Headers              map[string]string
	Tabs                 []Tab
}

// Tab represents a playground tab
type Tab struct {
	Name     string `json:"name"`
	Query    string `json:"query"`
	Variables string `json:"variables,omitempty"`
	Headers   string `json:"headers,omitempty"`
}

// DefaultConfig returns default playground configuration
func DefaultConfig() Config {
	return Config{
		GraphQLEndpoint:      "/graphql",
		SubscriptionEndpoint: "/graphql",
		Title:                "GraphQL Playground",
		Version:              "1.7.25",
		EnableInProduction:   false,
		EnableIntrospection:  true,
		Headers: map[string]string{
			"Authorization": "Bearer <your-token-here>",
		},
		Tabs: []Tab{
			{
				Name: "Welcome",
				Query: `# Welcome to GraphQL Playground
# GraphQL Playground is a powerful GraphQL IDE built by Prisma and based on GraphiQL.
#
# Here are some example queries to get you started:

query GetPosts {
  posts(first: 10) {
    edges {
      node {
        id
        title
        content
        author {
          id
          username
        }
        createdAt
      }
    }
    pageInfo {
      hasNextPage
      endCursor
    }
  }
}`,
				Variables: `{}`,
			},
			{
				Name: "Authentication",
				Query: `# Authentication Examples

mutation Login {
  login(input: {
    email: "user@example.com"
    password: "password123"
  }) {
    token
    user {
      id
      username
      email
    }
  }
}

mutation Register {
  register(input: {
    username: "newuser"
    email: "newuser@example.com"
    password: "password123"
  }) {
    token
    user {
      id
      username
      email
    }
  }
}`,
				Variables: `{}`,
			},
			{
				Name: "Mutations",
				Query: `# Mutation Examples

mutation CreatePost {
  createPost(input: {
    title: "My New Post"
    content: "This is the content of my new post."
    tags: ["graphql", "golang", "react"]
  }) {
    id
    title
    content
    author {
      username
    }
    createdAt
  }
}

mutation UpdatePost {
  updatePost(input: {
    id: "post-id-here"
    title: "Updated Post Title"
    content: "Updated content here."
  }) {
    id
    title
    content
    updatedAt
  }
}`,
				Variables: `{}`,
			},
			{
				Name: "Subscriptions",
				Query: `# Subscription Examples

subscription PostCreated {
  postCreated {
    id
    title
    content
    author {
      username
    }
    createdAt
  }
}

subscription CommentAdded {
  commentAdded(postId: "post-id-here") {
    id
    content
    author {
      username
    }
    createdAt
  }
}`,
				Variables: `{}`,
			},
		},
	}
}

// Handler creates an HTTP handler for GraphQL Playground
func Handler(config Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if playground should be enabled
		if !config.EnableInProduction && isProduction() {
			http.Error(w, "GraphQL Playground is disabled in production", http.StatusNotFound)
			return
		}

		// Set content type
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// Render playground HTML
		tmpl := template.Must(template.New("playground").Parse(playgroundTemplate))
		err := tmpl.Execute(w, config)
		if err != nil {
			http.Error(w, "Failed to render playground", http.StatusInternalServerError)
			return
		}
	}
}

// isProduction checks if we're running in production
func isProduction() bool {
	// This would typically check environment variables
	// For now, we'll assume development
	return false
}

// playgroundTemplate is the HTML template for GraphQL Playground
const playgroundTemplate = `
<!DOCTYPE html>
<html>
<head>
  <meta charset=utf-8/>
  <meta name="viewport" content="user-scalable=no, initial-scale=1.0, minimum-scale=1.0, maximum-scale=1.0, minimal-ui">
  <title>{{.Title}}</title>
  <link rel="stylesheet" href="//cdn.jsdelivr.net/npm/graphql-playground-react@{{.Version}}/build/static/css/index.css" />
  <link rel="shortcut icon" href="//cdn.jsdelivr.net/npm/graphql-playground-react@{{.Version}}/build/favicon.png" />
  <script src="//cdn.jsdelivr.net/npm/graphql-playground-react@{{.Version}}/build/static/js/middleware.js"></script>
  <style>
    body {
      margin: 0;
      padding: 0;
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue', sans-serif;
    }
    .loading {
      display: flex;
      align-items: center;
      justify-content: center;
      height: 100vh;
      font-size: 18px;
      color: #666;
    }
  </style>
</head>
<body>
  <div id="root">
    <div class="loading">Loading GraphQL Playground...</div>
  </div>
  <script>
    window.addEventListener('load', function (event) {
      GraphQLPlayground.init(document.getElementById('root'), {
        endpoint: '{{.GraphQLEndpoint}}',
        subscriptionEndpoint: '{{.SubscriptionEndpoint}}',
        settings: {
          'general.betaUpdates': false,
          'editor.theme': 'dark',
          'editor.cursorShape': 'line',
          'editor.reuseHeaders': true,
          'tracing.hideTracingResponse': true,
          'queryPlan.hideQueryPlanResponse': true,
          'editor.fontSize': 14,
          'editor.fontFamily': '"Source Code Pro", "Consolas", "Inconsolata", "Droid Sans Mono", "Monaco", monospace',
          'request.credentials': 'include',
        },
        tabs: {{.Tabs | toJSON}},
        {{if .Headers}}
        headers: {{.Headers | toJSON}},
        {{end}}
        introspection: {{.EnableIntrospection}},
        schema: undefined,
        workspaceName: '{{.Title}}',
      })
    })
  </script>
  <script>
    // Helper function to convert Go data to JSON
    function toJSON(data) {
      return JSON.stringify(data);
    }
    
    // Add toJSON to template functions
    if (typeof window !== 'undefined') {
      window.toJSON = toJSON;
    }
  </script>
</body>
</html>
`

// HealthHandler provides a health check endpoint
func HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"status": "healthy",
			"service": "graphql-server",
			"timestamp": "` + time.Now().Format(time.RFC3339) + `",
			"version": "1.0.0"
		}`))
	}
}

// IntrospectionHandler provides schema introspection
func IntrospectionHandler(schema string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if introspection is enabled
		if !shouldAllowIntrospection(r) {
			http.Error(w, "Schema introspection is disabled", http.StatusForbidden)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(schema))
	}
}

// shouldAllowIntrospection checks if introspection should be allowed
func shouldAllowIntrospection(r *http.Request) bool {
	// In production, you might want to restrict introspection
	// to authenticated users or disable it entirely
	return !isProduction()
}

// CORSMiddleware adds CORS headers for playground
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// SecurityHeadersMiddleware adds security headers
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// CSP for playground
		if r.URL.Path == "/playground" {
			w.Header().Set("Content-Security-Policy", 
				"default-src 'self'; "+
				"script-src 'self' 'unsafe-inline' 'unsafe-eval' cdn.jsdelivr.net; "+
				"style-src 'self' 'unsafe-inline' cdn.jsdelivr.net; "+
				"img-src 'self' data: cdn.jsdelivr.net; "+
				"connect-src 'self' ws: wss:; "+
				"font-src 'self' cdn.jsdelivr.net;")
		} else {
			w.Header().Set("Content-Security-Policy", "default-src 'self'")
		}

		next.ServeHTTP(w, r)
	})
}