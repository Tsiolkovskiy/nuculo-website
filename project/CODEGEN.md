# GraphQL Code Generation Setup

This document explains the GraphQL code generation setup for the TypeScript frontend.

## Overview

We use GraphQL Code Generator to automatically generate TypeScript types and React hooks from our GraphQL schema and operations. This ensures type safety and reduces boilerplate code.

## Configuration

### Dependencies

The following packages are installed for code generation:

```json
{
  "@graphql-codegen/cli": "^5.0.7",
  "@graphql-codegen/typescript": "^4.1.6",
  "@graphql-codegen/typescript-operations": "^4.6.1",
  "@graphql-codegen/typescript-react-apollo": "^4.3.3",
  "@graphql-codegen/introspection": "^4.0.3"
}
```

### Configuration File

The `codegen.yml` file configures the code generation:

```yaml
overwrite: true
schema: "../backend/internal/graph/schema.graphql"
documents: "src/graphql/**/*.graphql"
generates:
  src/generated/graphql.ts:
    plugins:
      - "typescript"
      - "typescript-operations"
      - "typescript-react-apollo"
    config:
      withHooks: true
      withHOC: false
      withComponent: false
      # ... additional configuration
  src/generated/introspection.json:
    plugins:
      - "introspection"
```

## GraphQL Operations Structure

Operations are organized in the `src/graphql/` directory:

```
src/graphql/
├── fragments/
│   ├── user.graphql
│   ├── post.graphql
│   └── comment.graphql
├── queries/
│   ├── auth.graphql
│   └── posts.graphql
├── mutations/
│   ├── auth.graphql
│   ├── posts.graphql
│   └── comments.graphql
└── subscriptions/
    ├── posts.graphql
    └── comments.graphql
```

### Fragments

Fragments define reusable pieces of GraphQL queries:

```graphql
fragment UserInfo on User {
  id
  email
  name
  avatar
  createdAt
  updatedAt
}
```

### Queries

Query operations for fetching data:

```graphql
#import "../fragments/user.graphql"

query Me {
  me {
    ...UserInfo
  }
}
```

### Mutations

Mutation operations for modifying data:

```graphql
#import "../fragments/post.graphql"

mutation CreatePost($input: CreatePostInput!) {
  createPost(input: $input) {
    ...PostInfo
  }
}
```

### Subscriptions

Subscription operations for real-time updates:

```graphql
#import "../fragments/post.graphql"

subscription PostAdded {
  postAdded {
    ...PostInfo
  }
}
```

## Generated Code

The code generator creates:

### TypeScript Types

```typescript
export type User = {
  __typename?: 'User';
  id: Scalars['ID']['output'];
  email: Scalars['String']['output'];
  name: Scalars['String']['output'];
  avatar?: Maybe<Scalars['String']['output']>;
  createdAt: Scalars['DateTime']['output'];
  updatedAt: Scalars['DateTime']['output'];
};
```

### React Hooks

```typescript
export function useGetPostsQuery(
  baseOptions?: ApolloReactHooks.QueryHookOptions<GetPostsQuery, GetPostsQueryVariables>
) {
  const options = {...defaultOptions, ...baseOptions}
  return ApolloReactHooks.useQuery<GetPostsQuery, GetPostsQueryVariables>(GetPostsDocument, options);
}
```

### GraphQL Documents

```typescript
export const GetPostsDocument = gql`
  query GetPosts($filters: PostFilters, $pagination: PaginationInput) {
    posts(filters: $filters, pagination: $pagination) {
      edges {
        node {
          ...PostSummary
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
  ${PostSummaryFragmentDoc}
`;
```

## Usage Examples

### Using Query Hooks

```typescript
import { useGetPostsQuery } from '../generated/graphql';

const PostsList = () => {
  const { data, loading, error } = useGetPostsQuery({
    variables: {
      filters: { published: true },
      pagination: { page: 1, limit: 10 }
    }
  });

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <div>
      {data?.posts?.edges.map(({ node: post }) => (
        <div key={post.id}>{post.title}</div>
      ))}
    </div>
  );
};
```

### Using Mutation Hooks

```typescript
import { useCreatePostMutation } from '../generated/graphql';

const CreatePostForm = () => {
  const [createPost, { loading, error }] = useCreatePostMutation({
    refetchQueries: ['GetPosts']
  });

  const handleSubmit = async (formData) => {
    try {
      const result = await createPost({
        variables: {
          input: {
            title: formData.title,
            content: formData.content,
            tags: formData.tags,
            published: formData.published
          }
        }
      });
      console.log('Post created:', result.data?.createPost);
    } catch (error) {
      console.error('Failed to create post:', error);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      {/* form fields */}
    </form>
  );
};
```

### Using Subscription Hooks

```typescript
import { usePostAddedSubscription } from '../generated/graphql';

const RealTimePosts = () => {
  const { data, loading } = usePostAddedSubscription();

  useEffect(() => {
    if (data?.postAdded) {
      console.log('New post added:', data.postAdded);
      // Update UI or show notification
    }
  }, [data]);

  return <div>Listening for new posts...</div>;
};
```

## Scripts

### Available Commands

```bash
# Generate types once
npm run codegen

# Watch for changes and regenerate
npm run codegen:watch

# Watch with silent output
npm run codegen:dev
```

### Development Workflow

1. **Write GraphQL operations** in `src/graphql/` directory
2. **Run code generation** with `npm run codegen`
3. **Import generated hooks** in your components
4. **Use type-safe GraphQL operations** with full IntelliSense

## Benefits

### Type Safety

- **Compile-time validation** of GraphQL operations
- **IntelliSense support** for queries, mutations, and variables
- **Automatic type checking** for response data

### Developer Experience

- **Auto-generated React hooks** for all operations
- **Consistent API** across all GraphQL operations
- **Reduced boilerplate** code

### Maintainability

- **Single source of truth** for GraphQL schema
- **Automatic updates** when schema changes
- **Consistent naming** conventions

## Best Practices

### Operation Naming

- Use descriptive names: `GetUserPosts` instead of `Posts`
- Follow PascalCase convention for operations
- Use consistent prefixes: `Get`, `Create`, `Update`, `Delete`

### Fragment Usage

- Create reusable fragments for common data structures
- Import fragments in operations using `#import`
- Keep fragments focused and cohesive

### Error Handling

- Always handle loading and error states
- Use Apollo's error policies for partial data
- Implement retry mechanisms for failed operations

### Performance

- Use pagination for large datasets
- Implement proper cache policies
- Leverage Apollo's caching mechanisms

## Troubleshooting

### Common Issues

1. **Schema not found**: Ensure the schema path in `codegen.yml` is correct
2. **Import errors**: Check fragment import paths in GraphQL files
3. **Type errors**: Regenerate types after schema changes

### Debugging

- Check the generated `graphql.ts` file for type definitions
- Use Apollo DevTools for query debugging
- Enable GraphQL Code Generator verbose logging

## Integration with Backend

The code generator reads the GraphQL schema from:
- **Local file**: `../backend/internal/graph/schema.graphql`
- **Remote endpoint**: `http://localhost:8080/graphql` (when server is running)

This ensures the frontend types are always in sync with the backend schema.