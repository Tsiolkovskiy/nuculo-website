# Design Document

## Overview

The GraphQL TypeScript-Go Integration system provides a modern, type-safe communication layer between a React TypeScript frontend and a Go backend. The architecture leverages gqlgen for Go GraphQL server implementation and Apollo Client for TypeScript client integration, ensuring end-to-end type safety and optimal performance.

### Key Design Principles
- **End-to-End Type Safety**: Automatic type generation from GraphQL schema to both Go and TypeScript
- **Performance First**: Efficient resolvers, DataLoader pattern, and intelligent caching
- **Developer Experience**: Hot reloading, code generation, and comprehensive tooling
- **Real-time Capabilities**: WebSocket-based subscriptions for live data updates
- **Security by Design**: Authentication, authorization, and query complexity analysis
- **Scalable Architecture**: Modular design supporting future feature expansion

## Architecture

### System Architecture
The system follows a GraphQL-first architecture with automatic code generation:

```
┌─────────────────────┐
│   React Frontend    │
│   (TypeScript)      │
├─────────────────────┤
│   Apollo Client     │
│ (GraphQL Client)    │
└──────────┬──────────┘
           │ GraphQL over HTTP/WebSocket
┌──────────▼──────────┐
│   Go GraphQL Server │
│     (gqlgen)        │
├─────────────────────┤
│   GraphQL Schema    │
│  (Schema-First)     │
├─────────────────────┤
│    Resolvers        │
│  (Go Functions)     │
├─────────────────────┤
│  Business Logic     │
│   (Services)        │
├─────────────────────┤
│   Data Layer        │
│ (Repositories/DB)   │
└─────────────────────┘
```

### Technology Stack

#### Backend (Go)
- **GraphQL Server**: gqlgen (99designs/gqlgen) - Schema-first GraphQL server
- **HTTP Server**: Gin or Chi router for HTTP handling
- **Database**: PostgreSQL with pgx driver
- **Authentication**: JWT tokens with golang-jwt/jwt
- **WebSocket**: gorilla/websocket for subscriptions
- **Caching**: Redis with go-redis client
- **Code Generation**: gqlgen for resolvers and types

#### Frontend (TypeScript)
- **GraphQL Client**: Apollo Client with TypeScript support
- **Code Generation**: GraphQL Code Generator (@graphql-codegen)
- **React Framework**: React 18 with hooks
- **Build Tool**: Vite for fast development and building
- **Type Checking**: TypeScript strict mode
- **State Management**: Apollo Client cache + React state

#### Development Tools
- **Schema Management**: GraphQL schema-first approach
- **Hot Reloading**: Both Go (air) and React (Vite) hot reload
- **Testing**: Go testing + React Testing Library
- **Linting**: golangci-lint + ESLint
- **Documentation**: GraphQL Playground + generated docs

## Components and Interfaces

### GraphQL Schema Definition

The schema serves as the contract between frontend and backend:

```graphql
# Scalars
scalar DateTime
scalar Upload

# Core Types
type User {
  id: ID!
  email: String!
  name: String!
  avatar: String
  createdAt: DateTime!
  updatedAt: DateTime!
}

type Post {
  id: ID!
  title: String!
  content: String!
  author: User!
  tags: [String!]!
  published: Boolean!
  createdAt: DateTime!
  updatedAt: DateTime!
}

type Comment {
  id: ID!
  content: String!
  author: User!
  post: Post!
  createdAt: DateTime!
}

# Input Types
input CreatePostInput {
  title: String!
  content: String!
  tags: [String!]!
  published: Boolean = false
}

input UpdatePostInput {
  title: String
  content: String
  tags: [String!]
  published: Boolean
}

input PostFilters {
  authorId: ID
  published: Boolean
  tags: [String!]
  searchTerm: String
}

input PaginationInput {
  page: Int = 1
  limit: Int = 20
}

# Response Types
type PostConnection {
  edges: [PostEdge!]!
  pageInfo: PageInfo!
  totalCount: Int!
}

type PostEdge {
  node: Post!
  cursor: String!
}

type PageInfo {
  hasNextPage: Boolean!
  hasPreviousPage: Boolean!
  startCursor: String
  endCursor: String
}

type AuthPayload {
  token: String!
  user: User!
  expiresAt: DateTime!
}

# Root Types
type Query {
  # User queries
  me: User
  user(id: ID!): User
  
  # Post queries
  posts(filters: PostFilters, pagination: PaginationInput): PostConnection!
  post(id: ID!): Post
  
  # Search
  searchPosts(query: String!, limit: Int = 10): [Post!]!
}

type Mutation {
  # Authentication
  login(email: String!, password: String!): AuthPayload!
  register(email: String!, password: String!, name: String!): AuthPayload!
  refreshToken: AuthPayload!
  
  # Post mutations
  createPost(input: CreatePostInput!): Post!
  updatePost(id: ID!, input: UpdatePostInput!): Post!
  deletePost(id: ID!): Boolean!
  
  # Comment mutations
  addComment(postId: ID!, content: String!): Comment!
  deleteComment(id: ID!): Boolean!
}

type Subscription {
  # Real-time updates
  postAdded: Post!
  postUpdated(id: ID!): Post!
  commentAdded(postId: ID!): Comment!
}
```

### Go Backend Implementation

#### Project Structure
```
backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── graph/
│   │   ├── generated/        # gqlgen generated code
│   │   ├── model/           # GraphQL models
│   │   ├── resolver/        # Resolver implementations
│   │   └── schema.graphql   # GraphQL schema
│   ├── auth/               # Authentication middleware
│   ├── database/           # Database connection and migrations
│   ├── service/            # Business logic services
│   └── repository/         # Data access layer
├── gqlgen.yml             # gqlgen configuration
├── go.mod
└── go.sum
```

#### gqlgen Configuration (gqlgen.yml)
```yaml
schema:
  - internal/graph/schema.graphql

exec:
  filename: internal/graph/generated/generated.go
  package: generated

model:
  filename: internal/graph/model/models_gen.go
  package: model

resolver:
  layout: follow-schema
  dir: internal/graph/resolver
  package: resolver
  filename_template: "{name}.resolvers.go"

autobind:
  - "github.com/your-org/your-project/internal/graph/model"

models:
  ID:
    model:
      - github.com/99designs/gqlgen/graphql.ID
      - github.com/99designs/gqlgen/graphql.Int
      - github.com/99designs/gqlgen/graphql.Int64
      - github.com/99designs/gqlgen/graphql.Int32
  DateTime:
    model: github.com/99designs/gqlgen/graphql.Time
```

#### Resolver Implementation
```go
package resolver

import (
    "context"
    "github.com/your-org/your-project/internal/graph/generated"
    "github.com/your-org/your-project/internal/graph/model"
    "github.com/your-org/your-project/internal/service"
)

type Resolver struct {
    userService    *service.UserService
    postService    *service.PostService
    commentService *service.CommentService
}

// Query resolver
func (r *queryResolver) Posts(ctx context.Context, filters *model.PostFilters, pagination *model.PaginationInput) (*model.PostConnection, error) {
    user := auth.GetUserFromContext(ctx)
    if user == nil {
        return nil, errors.New("authentication required")
    }
    
    return r.postService.GetPosts(ctx, filters, pagination)
}

func (r *queryResolver) Post(ctx context.Context, id string) (*model.Post, error) {
    return r.postService.GetPostByID(ctx, id)
}

// Mutation resolver
func (r *mutationResolver) CreatePost(ctx context.Context, input model.CreatePostInput) (*model.Post, error) {
    user := auth.GetUserFromContext(ctx)
    if user == nil {
        return nil, errors.New("authentication required")
    }
    
    post, err := r.postService.CreatePost(ctx, user.ID, input)
    if err != nil {
        return nil, err
    }
    
    // Trigger subscription
    r.subscriptionManager.PostAdded(post)
    
    return post, nil
}

// Subscription resolver
func (r *subscriptionResolver) PostAdded(ctx context.Context) (<-chan *model.Post, error) {
    user := auth.GetUserFromContext(ctx)
    if user == nil {
        return nil, errors.New("authentication required")
    }
    
    return r.subscriptionManager.SubscribeToPostAdded(ctx), nil
}

// Field resolvers for complex types
func (r *postResolver) Author(ctx context.Context, obj *model.Post) (*model.User, error) {
    return r.userService.GetUserByID(ctx, obj.AuthorID)
}

func (r *postResolver) Comments(ctx context.Context, obj *model.Post) ([]*model.Comment, error) {
    return r.commentService.GetCommentsByPostID(ctx, obj.ID)
}
```

#### Authentication Middleware
```go
package auth

import (
    "context"
    "net/http"
    "strings"
    "github.com/golang-jwt/jwt/v4"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            token := extractToken(r)
            if token != "" {
                if user, err := validateToken(token, secret); err == nil {
                    ctx := context.WithValue(r.Context(), UserContextKey, user)
                    r = r.WithContext(ctx)
                }
            }
            next.ServeHTTP(w, r)
        })
    }
}

func GetUserFromContext(ctx context.Context) *User {
    user, ok := ctx.Value(UserContextKey).(*User)
    if !ok {
        return nil
    }
    return user
}
```

### TypeScript Frontend Implementation

#### Project Structure
```
frontend/
├── src/
│   ├── components/         # React components
│   ├── graphql/           # GraphQL operations
│   │   ├── generated/     # Generated types and hooks
│   │   ├── mutations/     # GraphQL mutations
│   │   ├── queries/       # GraphQL queries
│   │   └── subscriptions/ # GraphQL subscriptions
│   ├── hooks/             # Custom React hooks
│   ├── utils/             # Utility functions
│   └── types/             # Additional TypeScript types
├── codegen.yml           # GraphQL Code Generator config
├── apollo.config.js      # Apollo Client configuration
└── package.json
```

#### GraphQL Code Generator Configuration (codegen.yml)
```yaml
overwrite: true
schema: "http://localhost:8080/graphql"
documents: "src/graphql/**/*.graphql"
generates:
  src/graphql/generated/types.ts:
    plugins:
      - "typescript"
      - "typescript-operations"
  src/graphql/generated/hooks.ts:
    plugins:
      - "typescript"
      - "typescript-operations"
      - "typescript-react-apollo"
    config:
      withHooks: true
      withComponent: false
      withHOC: false
```

#### Apollo Client Setup
```typescript
// src/apollo/client.ts
import { ApolloClient, InMemoryCache, createHttpLink, split } from '@apollo/client';
import { setContext } from '@apollo/client/link/context';
import { GraphQLWsLink } from '@apollo/client/link/subscriptions';
import { getMainDefinition } from '@apollo/client/utilities';
import { createClient } from 'graphql-ws';

const httpLink = createHttpLink({
  uri: 'http://localhost:8080/graphql',
});

const wsLink = new GraphQLWsLink(createClient({
  url: 'ws://localhost:8080/graphql',
  connectionParams: () => ({
    Authorization: `Bearer ${localStorage.getItem('token')}`,
  }),
}));

const authLink = setContext((_, { headers }) => {
  const token = localStorage.getItem('token');
  return {
    headers: {
      ...headers,
      authorization: token ? `Bearer ${token}` : "",
    }
  };
});

const splitLink = split(
  ({ query }) => {
    const definition = getMainDefinition(query);
    return (
      definition.kind === 'OperationDefinition' &&
      definition.operation === 'subscription'
    );
  },
  wsLink,
  authLink.concat(httpLink),
);

export const apolloClient = new ApolloClient({
  link: splitLink,
  cache: new InMemoryCache({
    typePolicies: {
      Query: {
        fields: {
          posts: {
            keyArgs: ["filters"],
            merge(existing, incoming) {
              return {
                ...incoming,
                edges: [...(existing?.edges || []), ...incoming.edges],
              };
            },
          },
        },
      },
    },
  }),
});
```

#### GraphQL Operations
```graphql
# src/graphql/queries/posts.graphql
query GetPosts($filters: PostFilters, $pagination: PaginationInput) {
  posts(filters: $filters, pagination: $pagination) {
    edges {
      node {
        id
        title
        content
        published
        createdAt
        author {
          id
          name
          avatar
        }
        tags
      }
      cursor
    }
    pageInfo {
      hasNextPage
      hasPreviousPage
      startCursor
      endCursor
    }
    totalCount
  }
}

# src/graphql/mutations/createPost.graphql
mutation CreatePost($input: CreatePostInput!) {
  createPost(input: $input) {
    id
    title
    content
    published
    createdAt
    author {
      id
      name
    }
    tags
  }
}

# src/graphql/subscriptions/postAdded.graphql
subscription PostAdded {
  postAdded {
    id
    title
    content
    published
    createdAt
    author {
      id
      name
      avatar
    }
    tags
  }
}
```

#### React Component with Generated Hooks
```typescript
// src/components/PostList.tsx
import React from 'react';
import { useGetPostsQuery, usePostAddedSubscription } from '../graphql/generated/hooks';

interface PostListProps {
  filters?: PostFilters;
}

export const PostList: React.FC<PostListProps> = ({ filters }) => {
  const { data, loading, error, fetchMore } = useGetPostsQuery({
    variables: { filters, pagination: { page: 1, limit: 10 } },
    errorPolicy: 'partial',
  });

  // Subscribe to new posts
  usePostAddedSubscription({
    onSubscriptionData: ({ subscriptionData }) => {
      if (subscriptionData.data?.postAdded) {
        // Apollo Client automatically updates cache
        console.log('New post added:', subscriptionData.data.postAdded);
      }
    },
  });

  const loadMore = () => {
    if (data?.posts.pageInfo.hasNextPage) {
      fetchMore({
        variables: {
          pagination: {
            page: Math.floor(data.posts.edges.length / 10) + 1,
            limit: 10,
          },
        },
      });
    }
  };

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <div>
      {data?.posts.edges.map(({ node: post }) => (
        <div key={post.id}>
          <h3>{post.title}</h3>
          <p>{post.content}</p>
          <small>By {post.author.name} on {post.createdAt}</small>
        </div>
      ))}
      {data?.posts.pageInfo.hasNextPage && (
        <button onClick={loadMore}>Load More</button>
      )}
    </div>
  );
};
```

## Data Models

### Go Struct Definitions
```go
// internal/graph/model/user.go
type User struct {
    ID        string    `json:"id" db:"id"`
    Email     string    `json:"email" db:"email"`
    Name      string    `json:"name" db:"name"`
    Avatar    *string   `json:"avatar" db:"avatar"`
    CreatedAt time.Time `json:"createdAt" db:"created_at"`
    UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// internal/graph/model/post.go
type Post struct {
    ID        string    `json:"id" db:"id"`
    Title     string    `json:"title" db:"title"`
    Content   string    `json:"content" db:"content"`
    AuthorID  string    `json:"authorId" db:"author_id"`
    Tags      []string  `json:"tags" db:"tags"`
    Published bool      `json:"published" db:"published"`
    CreatedAt time.Time `json:"createdAt" db:"created_at"`
    UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
```

### Database Schema
```sql
-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    avatar TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Posts table
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(500) NOT NULL,
    content TEXT NOT NULL,
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tags TEXT[] DEFAULT '{}',
    published BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Comments table
CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content TEXT NOT NULL,
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_posts_author_id ON posts(author_id);
CREATE INDEX idx_posts_published ON posts(published);
CREATE INDEX idx_posts_created_at ON posts(created_at);
CREATE INDEX idx_comments_post_id ON comments(post_id);
```

## Error Handling

### GraphQL Error Extensions
```go
// internal/graph/error.go
type ErrorCode string

const (
    ErrorCodeValidation     ErrorCode = "VALIDATION_ERROR"
    ErrorCodeAuthentication ErrorCode = "AUTHENTICATION_ERROR"
    ErrorCodeAuthorization  ErrorCode = "AUTHORIZATION_ERROR"
    ErrorCodeNotFound       ErrorCode = "NOT_FOUND"
    ErrorCodeInternal       ErrorCode = "INTERNAL_ERROR"
)

func NewGraphQLError(message string, code ErrorCode, path []string) *gqlerror.Error {
    return &gqlerror.Error{
        Message: message,
        Path:    path,
        Extensions: map[string]interface{}{
            "code": string(code),
        },
    }
}
```

### TypeScript Error Handling
```typescript
// src/utils/errorHandling.ts
import { ApolloError } from '@apollo/client';

export interface GraphQLErrorExtensions {
  code: string;
  field?: string;
}

export const handleGraphQLError = (error: ApolloError) => {
  if (error.graphQLErrors.length > 0) {
    const graphQLError = error.graphQLErrors[0];
    const extensions = graphQLError.extensions as GraphQLErrorExtensions;
    
    switch (extensions.code) {
      case 'AUTHENTICATION_ERROR':
        // Redirect to login
        window.location.href = '/login';
        break;
      case 'VALIDATION_ERROR':
        // Show field-specific errors
        return {
          type: 'validation',
          message: graphQLError.message,
          field: extensions.field,
        };
      default:
        return {
          type: 'general',
          message: graphQLError.message,
        };
    }
  }
  
  if (error.networkError) {
    return {
      type: 'network',
      message: 'Network error occurred. Please try again.',
    };
  }
  
  return {
    type: 'unknown',
    message: 'An unexpected error occurred.',
  };
};
```

## Testing Strategy

### Go Backend Testing
```go
// internal/graph/resolver/post_test.go
func TestCreatePost(t *testing.T) {
    resolver := setupTestResolver(t)
    ctx := context.WithValue(context.Background(), auth.UserContextKey, &model.User{
        ID: "user-1",
        Email: "test@example.com",
    })
    
    input := model.CreatePostInput{
        Title:   "Test Post",
        Content: "Test content",
        Tags:    []string{"test"},
    }
    
    post, err := resolver.Mutation().CreatePost(ctx, input)
    assert.NoError(t, err)
    assert.Equal(t, "Test Post", post.Title)
    assert.Equal(t, "user-1", post.AuthorID)
}
```

### TypeScript Frontend Testing
```typescript
// src/components/__tests__/PostList.test.tsx
import { render, screen } from '@testing-library/react';
import { MockedProvider } from '@apollo/client/testing';
import { PostList } from '../PostList';
import { GET_POSTS } from '../../graphql/queries/posts';

const mocks = [
  {
    request: {
      query: GET_POSTS,
      variables: { filters: {}, pagination: { page: 1, limit: 10 } },
    },
    result: {
      data: {
        posts: {
          edges: [
            {
              node: {
                id: '1',
                title: 'Test Post',
                content: 'Test content',
                author: { id: '1', name: 'Test User', avatar: null },
                tags: ['test'],
                published: true,
                createdAt: '2023-01-01T00:00:00Z',
              },
              cursor: 'cursor1',
            },
          ],
          pageInfo: {
            hasNextPage: false,
            hasPreviousPage: false,
            startCursor: 'cursor1',
            endCursor: 'cursor1',
          },
          totalCount: 1,
        },
      },
    },
  },
];

test('renders post list', async () => {
  render(
    <MockedProvider mocks={mocks} addTypename={false}>
      <PostList />
    </MockedProvider>
  );
  
  expect(await screen.findByText('Test Post')).toBeInTheDocument();
  expect(screen.getByText('Test content')).toBeInTheDocument();
});
```