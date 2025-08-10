# GraphQL TypeScript-Go Frontend

A modern React TypeScript frontend application with Apollo Client integration for the GraphQL TypeScript-Go demo project.

## Features

- **React 18** with TypeScript for type-safe development
- **Apollo Client** for GraphQL data management with caching
- **React Router** for client-side routing
- **Tailwind CSS** for responsive styling
- **Authentication Context** with JWT token management
- **WebSocket Support** for real-time GraphQL subscriptions
- **Error Handling** with comprehensive error boundaries
- **Loading States** and user feedback components

## Project Structure

```
src/
├── components/          # Reusable UI components
│   ├── Layout.tsx      # Main application layout
│   ├── LoadingSpinner.tsx
│   └── ProtectedRoute.tsx
├── contexts/           # React contexts
│   └── AuthContext.tsx # Authentication state management
├── lib/               # Utility libraries
│   └── apollo.ts      # Apollo Client configuration
├── pages/             # Page components
│   ├── HomePage.tsx
│   ├── LoginPage.tsx
│   ├── RegisterPage.tsx
│   ├── DashboardPage.tsx
│   ├── PostsPage.tsx
│   ├── CreatePostPage.tsx
│   └── ProfilePage.tsx
├── App.tsx            # Main application component
└── main.tsx          # Application entry point
```

## Getting Started

### Prerequisites

- Node.js 18+ and npm
- Go backend server running on port 8080

### Installation

1. Install dependencies:
```bash
npm install
```

2. Copy environment variables:
```bash
cp .env.example .env
```

3. Update environment variables in `.env`:
```env
VITE_GRAPHQL_HTTP_URL=http://localhost:8080/graphql
VITE_GRAPHQL_WS_URL=ws://localhost:8080/graphql
```

### Development

Start the development server:
```bash
npm run dev
```

The application will be available at `http://localhost:5173`

### Build

Build for production:
```bash
npm run build
```

### Preview Production Build

Preview the production build:
```bash
npm run preview
```

## Apollo Client Configuration

The Apollo Client is configured with:

- **HTTP Link** for queries and mutations
- **WebSocket Link** for subscriptions
- **Authentication Link** for JWT token injection
- **Error Link** for centralized error handling
- **Cache Policies** for optimized data management

## Authentication

The application uses JWT-based authentication with:

- Login and registration forms
- Protected routes requiring authentication
- Automatic token refresh handling
- Secure token storage in localStorage
- Context-based authentication state management

## Available Routes

- `/` - Home page (public)
- `/login` - Login page (public)
- `/register` - Registration page (public)
- `/posts` - Browse posts (public)
- `/dashboard` - User dashboard (protected)
- `/create-post` - Create new post (protected)
- `/profile` - User profile (protected)

## GraphQL Integration

The frontend is ready to integrate with the Go GraphQL backend with:

- Type-safe GraphQL operations
- Automatic error handling
- Real-time subscriptions support
- Optimistic updates
- Cache management

## Styling

The application uses Tailwind CSS for styling with:

- Responsive design patterns
- Consistent color scheme
- Accessible form components
- Loading and error states
- Modern UI components

## GraphQL Code Generation

The project uses GraphQL Code Generator to create type-safe React hooks and TypeScript types from the GraphQL schema.

### Generate Types

```bash
# Generate types once
npm run codegen

# Watch for changes and regenerate
npm run codegen:watch
```

### Generated Files

- `src/generated/graphql.ts` - TypeScript types and React hooks
- `src/generated/introspection.json` - Schema introspection data

### Usage Example

```typescript
import { useGetPostsQuery, useCreatePostMutation } from '../generated/graphql';

const PostsList = () => {
  const { data, loading, error } = useGetPostsQuery({
    variables: { filters: { published: true } }
  });
  
  const [createPost] = useCreatePostMutation();
  
  // Fully type-safe GraphQL operations
};
```

See [CODEGEN.md](./CODEGEN.md) for detailed documentation.

## Next Steps

This frontend now has complete GraphQL integration ready. The next tasks will:

1. ✅ Generate TypeScript types from GraphQL schema
2. ✅ Implement type-safe GraphQL queries and mutations  
3. Add real-time subscriptions
4. Connect all forms to backend operations
5. Add comprehensive error handling

## Development Notes

- All components are written in TypeScript for type safety
- Apollo Client is configured for both HTTP and WebSocket connections
- Authentication state is managed through React Context
- Protected routes automatically redirect to login
- Error boundaries handle GraphQL and network errors
- The build process generates optimized production assets