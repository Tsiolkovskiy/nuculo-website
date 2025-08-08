# GraphQL TypeScript-Go Backend

This is the Go backend for the GraphQL TypeScript-Go integration project, built with gqlgen.

## Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Main server entry point
├── internal/
│   ├── graph/
│   │   ├── generated/           # gqlgen generated code
│   │   │   └── generated.go
│   │   ├── model/               # GraphQL models
│   │   │   ├── models_gen.go    # Generated models
│   │   │   └── types.go         # Custom types
│   │   ├── resolver/            # Resolver implementations
│   │   │   ├── resolver.go      # Root resolver
│   │   │   ├── query.resolvers.go
│   │   │   ├── mutation.resolvers.go
│   │   │   ├── subscription.resolvers.go
│   │   │   └── field.resolvers.go
│   │   └── schema.graphql       # GraphQL schema definition
├── scripts/
│   └── setup.sh                # Setup script
├── gqlgen.yml                  # gqlgen configuration
├── go.mod                      # Go module definition
└── README.md                   # This file
```

## Prerequisites

- Go 1.21 or later
- gqlgen CLI tool

## Setup

1. **Install Go dependencies:**
   ```bash
   go mod tidy
   ```

2. **Install gqlgen:**
   ```bash
   go install github.com/99designs/gqlgen@latest
   ```

3. **Generate GraphQL code:**
   ```bash
   go run github.com/99designs/gqlgen generate
   ```

   Or use the setup script:
   ```bash
   chmod +x scripts/setup.sh
   ./scripts/setup.sh
   ```

## Running the Server

```bash
go run cmd/server/main.go
```

The server will start on port 8080 by default. You can access:

- **GraphQL endpoint**: http://localhost:8080/graphql
- **GraphQL Playground**: http://localhost:8080/playground
- **Health check**: http://localhost:8080/health

## GraphQL Schema

The GraphQL schema is defined in `internal/graph/schema.graphql` and includes:

### Types
- **User**: User account information
- **Post**: Blog posts with author, tags, and content
- **Comment**: Comments on posts
- **AuthPayload**: Authentication response with JWT token

### Operations
- **Queries**: Get users, posts, search functionality
- **Mutations**: Authentication, CRUD operations for posts and comments
- **Subscriptions**: Real-time updates for posts and comments

### Example Queries

**Get all posts:**
```graphql
query GetPosts {
  posts {
    edges {
      node {
        id
        title
        content
        author {
          name
          email
        }
        tags
        published
        createdAt
      }
    }
    pageInfo {
      hasNextPage
      totalCount
    }
  }
}
```

**Create a post:**
```graphql
mutation CreatePost {
  createPost(input: {
    title: "My New Post"
    content: "This is the content of my post"
    tags: ["graphql", "go"]
    published: true
  }) {
    id
    title
    author {
      name
    }
    createdAt
  }
}
```

## Development

### Code Generation

When you modify the GraphQL schema in `internal/graph/schema.graphql`, run:

```bash
go run github.com/99designs/gqlgen generate
```

This will update the generated code and resolver interfaces.

### Adding New Resolvers

1. Update the schema in `internal/graph/schema.graphql`
2. Run code generation
3. Implement the new resolver methods in the appropriate resolver files
4. Add business logic and database operations

### Environment Variables

- `PORT`: Server port (default: 8080)

## Next Steps

This is the foundation setup for Task 1. The following tasks will add:

- Database integration (PostgreSQL)
- Authentication and JWT middleware
- Business logic services
- Real-time subscriptions
- Performance optimizations
- Testing

## Current Status

✅ **Task 1 Complete**: Go GraphQL server foundation with gqlgen
- Project structure created
- GraphQL schema defined
- Basic resolvers implemented
- Server with CORS and playground configured

✅ **Task 2 Complete**: Database layer and models
- PostgreSQL connection with pgx driver and connection pooling
- Database migration files for users, posts, and comments tables
- Repository interfaces and concrete implementations
- Go struct models with proper JSON and database tags
- Unit tests for repository CRUD operations
- Database initialization and migration utilities

✅ **Task 3 Complete**: Authentication and JWT middleware
- JWT token generation and validation using golang-jwt/jwt
- Password hashing and verification with bcrypt
- HTTP authentication middleware for Gin
- User context extraction and injection for GraphQL resolvers
- Comprehensive authentication service with login/register/refresh
- Unit tests for all authentication functions

✅ **Task 4 Complete**: Basic GraphQL resolvers
- Resolver struct with service dependencies injection
- Query resolvers for users, posts with filtering and pagination
- Mutation resolvers for authentication and CRUD operations
- Field resolvers for handling relationships (post.author, comment.author)
- Authentication checks integrated into protected operations
- Comprehensive integration tests with mocked dependencies

✅ **Task 5 Complete**: GraphQL input validation and error handling
- Custom GraphQL scalars for email and datetime validation
- Comprehensive input validation for all GraphQL operations
- Structured error response format with error codes and extensions
- GraphQL error handling with proper error categorization
- Unit tests for all validation and error handling scenarios

The server now has complete GraphQL functionality with robust validation and error handling.
## Datab
ase Setup

### Prerequisites
- PostgreSQL 12+ installed and running
- Database user with CREATE DATABASE privileges

### Environment Variables
Set these environment variables for database connection:

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_NAME=graphql_typescript_go
export DB_SSL_MODE=disable
```

### Database Migration

**Run migrations:**
```bash
go run cmd/migrate/main.go -up
```

**Rollback migrations:**
```bash
go run cmd/migrate/main.go -down -steps=1
```

### Database Schema

The database includes three main tables:

#### Users Table
- `id` (UUID, Primary Key)
- `email` (VARCHAR, Unique)
- `name` (VARCHAR)
- `password_hash` (VARCHAR)
- `avatar` (TEXT, Optional)
- `created_at`, `updated_at` (TIMESTAMP)

#### Posts Table
- `id` (UUID, Primary Key)
- `title` (VARCHAR)
- `content` (TEXT)
- `author_id` (UUID, Foreign Key to users)
- `tags` (TEXT Array)
- `published` (BOOLEAN)
- `created_at`, `updated_at` (TIMESTAMP)

#### Comments Table
- `id` (UUID, Primary Key)
- `content` (TEXT)
- `author_id` (UUID, Foreign Key to users)
- `post_id` (UUID, Foreign Key to posts)
- `created_at` (TIMESTAMP)

### Repository Pattern

The database layer uses the repository pattern with interfaces:

```go
// Example usage
repos := repository.NewManager(db)

// Create a user
user := &model.User{
    ID:    uuid.New(),
    Email: "user@example.com",
    Name:  "John Doe",
}
err := repos.User.Create(ctx, user)

// Get posts with filters
filters := &repository.PostFilters{
    Published: &[]bool{true}[0],
    Tags:      []string{"golang", "graphql"},
}
posts, err := repos.Post.List(ctx, filters, 10, 0)
```

### Testing Database Integration

Run the database integration example:
```bash
go run cmd/db-server/main.go
```

This will demonstrate CRUD operations across all repositories.
## A
uthentication System

### JWT Configuration

Set these environment variables for JWT configuration:

```bash
export JWT_SECRET=your-super-secret-jwt-key-change-in-production
export JWT_TOKEN_DURATION=24h
export BCRYPT_COST=12
export JWT_REFRESH_WINDOW=2h
```

### Authentication Endpoints

The authentication system provides the following endpoints:

#### Register User
```bash
POST /auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "name": "John Doe",
  "password": "securepassword123"
}
```

#### Login User
```bash
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

#### Refresh Token
```bash
POST /auth/refresh
Authorization: Bearer <your-jwt-token>
```

#### Get Current User (Protected)
```bash
GET /api/me
Authorization: Bearer <your-jwt-token>
```

### Password Requirements

Passwords must meet the following criteria:
- At least 8 characters long
- Less than 128 characters long
- Contains at least one letter
- Contains at least one number

### Middleware Usage

The authentication system provides two middleware options:

#### Required Authentication
```go
// Requires valid JWT token, returns 401 if not authenticated
r.Use(authManager.Middleware.RequiredAuth())
```

#### Optional Authentication
```go
// Extracts user if token is present, continues without user if not
r.Use(authManager.Middleware.OptionalAuth())
```

### GraphQL Context Integration

In GraphQL resolvers, you can access the authenticated user:

```go
func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
    user, err := auth.RequireUser(ctx)
    if err != nil {
        return nil, err
    }
    return user, nil
}

func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
    // Optional authentication
    user, ok := auth.GetUserFromContext(ctx)
    if ok {
        // User is authenticated, show their posts
        return r.postService.GetUserPosts(ctx, user.ID)
    }
    // User not authenticated, show public posts only
    return r.postService.GetPublicPosts(ctx)
}
```

### Testing Authentication

Run the authentication server example:
```bash
go run cmd/auth-server/main.go
```

This demonstrates:
- User registration and login
- JWT token generation and validation
- Password hashing and verification
- Protected and public endpoints
- Middleware integration## Graph
QL Resolvers

### Resolver Implementation

The GraphQL resolvers are fully implemented with authentication and database integration:

#### Query Resolvers
- `me` - Get current authenticated user (requires auth)
- `user(id)` - Get user by ID
- `posts(filters, pagination)` - Get posts with filtering and pagination
- `post(id)` - Get single post by ID
- `searchPosts(query, limit)` - Search posts by title/content

#### Mutation Resolvers
- `login(email, password)` - Authenticate user
- `register(email, password, name)` - Register new user
- `refreshToken` - Refresh JWT token (requires auth)
- `createPost(input)` - Create new post (requires auth)
- `updatePost(id, input)` - Update post (requires auth, owner only)
- `deletePost(id)` - Delete post (requires auth, owner only)
- `addComment(postId, content)` - Add comment (requires auth)
- `deleteComment(id)` - Delete comment (requires auth, owner only)

#### Field Resolvers
- `Post.author` - Resolve post author from user repository
- `Comment.author` - Resolve comment author from user repository
- `Comment.post` - Resolve comment's post from post repository

### Authentication Integration

Resolvers use the authentication system seamlessly:

```go
func (r *mutationResolver) CreatePost(ctx context.Context, input model.CreatePostInput) (*model.Post, error) {
    // Require authentication
    user, err := auth.RequireUser(ctx)
    if err != nil {
        return nil, fmt.Errorf("authentication required: %w", err)
    }
    
    // Use authenticated user
    post := &model.Post{
        AuthorID: user.ID,
        // ... other fields
    }
    
    return r.PostRepo.Create(ctx, post)
}
```

### Authorization Checks

Protected operations include ownership verification:

```go
func (r *mutationResolver) UpdatePost(ctx context.Context, id string, input model.UpdatePostInput) (*model.Post, error) {
    user, err := auth.RequireUser(ctx)
    if err != nil {
        return nil, err
    }
    
    post, err := r.PostRepo.GetByID(ctx, postID)
    if err != nil {
        return nil, err
    }
    
    // Check ownership
    if post.AuthorID != user.ID {
        return nil, fmt.Errorf("unauthorized: you can only update your own posts")
    }
    
    // Proceed with update...
}
```

### Testing GraphQL Resolvers

Run the GraphQL resolver tests:
```bash
go test -v ./internal/graph/resolver/
```

Run the simple GraphQL server:
```bash
go run cmd/simple-graphql-server/main.go
```

### Resolver Dependencies

The resolver struct includes all necessary dependencies:

```go
type Resolver struct {
    // Repository dependencies
    UserRepo    repository.UserRepository
    PostRepo    repository.PostRepository
    CommentRepo repository.CommentRepository
    
    // Authentication service
    AuthManager *auth.Manager
}
```

This allows for easy dependency injection and testing with mocks.## Gr
aphQL Validation and Error Handling

### Input Validation

The GraphQL API includes comprehensive input validation:

#### Post Validation
- **Title**: 3-200 characters, no HTML tags
- **Content**: 10-50,000 characters minimum
- **Tags**: Max 10 tags, 2-30 characters each, alphanumeric + hyphens/underscores only
- **No duplicate tags allowed**

#### User Validation
- **Email**: Valid email format, max 254 characters
- **Name**: 2-100 characters, letters/spaces/hyphens/apostrophes only
- **Password**: 8-128 characters, must contain letters and numbers

#### Pagination Validation
- **Page**: 1-1000 range
- **Limit**: 1-100 range

#### Search Validation
- **Query**: 2-100 characters
- **Limit**: 1-50 range for search results

### Custom GraphQL Scalars

#### Email Scalar
```graphql
scalar Email
```
- Validates email format using regex
- Supports marshaling/unmarshaling with validation
- Returns proper GraphQL errors for invalid emails

#### DateTime Scalar
```graphql
scalar DateTime
```
- Supports multiple datetime formats (RFC3339, Unix timestamps)
- Validates datetime ranges (1900-2125)
- Handles timezone conversions properly

### Error Response Format

All GraphQL errors follow a structured format:

```json
{
  "errors": [
    {
      "message": "Title must be at least 3 characters long",
      "extensions": {
        "code": "VALIDATION_ERROR",
        "field": "title"
      },
      "path": ["createPost"]
    }
  ]
}
```

### Error Codes

The API uses standardized error codes:

- `VALIDATION_ERROR` - Input validation failures
- `INVALID_INPUT` - Malformed input data
- `INVALID_FORMAT` - Format validation errors (UUID, email, etc.)
- `UNAUTHENTICATED` - Authentication required
- `UNAUTHORIZED` - Access denied
- `FORBIDDEN` - Insufficient permissions
- `NOT_FOUND` - Resource not found
- `ALREADY_EXISTS` - Duplicate resource
- `CONFLICT` - Resource conflict
- `INTERNAL_ERROR` - Server errors
- `DATABASE_ERROR` - Database operation failures
- `RATE_LIMIT_EXCEEDED` - Rate limiting

### Validation Usage Example

```go
// In resolvers
func (r *mutationResolver) CreatePost(ctx context.Context, input model.CreatePostInput) (*model.Post, error) {
    // Validate input
    validator := validation.NewValidator()
    if err := validator.ValidateCreatePostInput(input); err != nil {
        return nil, err // Returns structured GraphQL error
    }
    
    // Continue with business logic...
}
```

### Error Handling Features

- **Automatic Error Categorization**: Common errors are automatically categorized
- **Database Error Wrapping**: Database errors are wrapped with appropriate GraphQL errors
- **Validation Error Aggregation**: Multiple validation errors can be combined
- **Security-Safe Errors**: Internal errors don't expose sensitive information
- **Structured Logging**: All errors are logged with appropriate severity levels

### Testing Validation

Run validation tests:
```bash
go test -v ./internal/graph/validation/
go test -v ./internal/graph/scalars/
```

The validation system ensures data integrity and provides clear, actionable error messages to clients.