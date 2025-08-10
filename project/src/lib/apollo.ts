import {
  ApolloClient,
  InMemoryCache,
  createHttpLink,
  split,
  from,
} from '@apollo/client';
import { setContext } from '@apollo/client/link/context';
import { onError } from '@apollo/client/link/error';
import { GraphQLWsLink } from '@apollo/client/link/subscriptions';
import { getMainDefinition } from '@apollo/client/utilities';
import { createClient } from 'graphql-ws';

// HTTP Link for queries and mutations
const httpLink = createHttpLink({
  uri: import.meta.env.VITE_GRAPHQL_HTTP_URL || 'http://localhost:8080/graphql',
});

// WebSocket Link for subscriptions
const wsLink = new GraphQLWsLink(
  createClient({
    url: import.meta.env.VITE_GRAPHQL_WS_URL || 'ws://localhost:8080/graphql',
    connectionParams: () => {
      const token = localStorage.getItem('authToken');
      return {
        Authorization: token ? `Bearer ${token}` : '',
      };
    },
    on: {
      connected: () => console.log('🔗 GraphQL WebSocket connected'),
      closed: () => console.log('🔌 GraphQL WebSocket disconnected'),
      error: (error) => console.error('❌ GraphQL WebSocket error:', error),
    },
  })
);

// Authentication Link
const authLink = setContext((_, { headers }) => {
  const token = localStorage.getItem('authToken');
  return {
    headers: {
      ...headers,
      authorization: token ? `Bearer ${token}` : '',
    },
  };
});

// Error Link with comprehensive error handling
const errorLink = onError(({ graphQLErrors, networkError, operation, forward }) => {
  const operationName = operation.operationName || 'Unknown';
  const operationType = operation.query.definitions[0]?.kind === 'OperationDefinition' 
    ? operation.query.definitions[0].operation 
    : 'unknown';

  if (graphQLErrors) {
    graphQLErrors.forEach(({ message, locations, path, extensions }) => {
      console.group(`🚨 GraphQL Error in ${operationName}`);
      console.error('Message:', message);
      console.error('Path:', path);
      console.error('Locations:', locations);
      console.error('Extensions:', extensions);
      console.error('Operation:', operationType);
      console.groupEnd();
      
      // Handle authentication errors
      if (extensions?.code === 'UNAUTHENTICATED' || extensions?.code === 'INVALID_TOKEN') {
        console.warn('🔐 Authentication error - redirecting to login');
        localStorage.removeItem('authToken');
        localStorage.removeItem('user');
        
        // Avoid infinite redirects
        if (window.location.pathname !== '/login') {
          window.location.href = '/login';
        }
      }
      
      // Handle authorization errors
      if (extensions?.code === 'UNAUTHORIZED' || extensions?.code === 'FORBIDDEN') {
        console.warn('🚫 Authorization error');
        // Could show a toast notification here
      }
      
      // Handle validation errors
      if (extensions?.code === 'VALIDATION_FAILED' || extensions?.code === 'INVALID_INPUT') {
        console.warn('⚠️ Validation error:', message);
        // These are usually handled by the UI components
      }
    });
  }

  if (networkError) {
    console.group(`🌐 Network Error in ${operationName}`);
    console.error('Network Error:', networkError);
    console.error('Operation:', operationType);
    
    // Type-safe network error handling
    if ('statusCode' in networkError) {
      console.error('Status Code:', networkError.statusCode);
      
      // Handle specific HTTP status codes
      switch (networkError.statusCode) {
        case 401:
          console.warn('🔐 401 Unauthorized - clearing auth and redirecting');
          localStorage.removeItem('authToken');
          localStorage.removeItem('user');
          if (window.location.pathname !== '/login') {
            window.location.href = '/login';
          }
          break;
          
        case 403:
          console.warn('🚫 403 Forbidden - insufficient permissions');
          break;
          
        case 404:
          console.warn('🔍 404 Not Found - resource does not exist');
          break;
          
        case 429:
          console.warn('🐌 429 Rate Limited - too many requests');
          break;
          
        case 500:
        case 502:
        case 503:
        case 504:
          console.error('🔥 Server Error - service may be down');
          break;
          
        default:
          console.error(`❓ Unexpected status code: ${networkError.statusCode}`);
      }
    }
    
    // Handle network connectivity issues
    if ('code' in networkError) {
      switch (networkError.code) {
        case 'NETWORK_ERROR':
          console.error('📡 Network connectivity issue');
          break;
          
        case 'TIMEOUT':
          console.error('⏰ Request timeout');
          break;
          
        default:
          console.error(`❓ Network error code: ${networkError.code}`);
      }
    }
    
    console.groupEnd();
  }
  
  // Log operation context for debugging
  if (process.env.NODE_ENV === 'development') {
    console.log('🔍 Operation Context:', {
      operationName,
      operationType,
      variables: operation.variables,
      query: operation.query.loc?.source.body,
    });
  }
});

// Split link to route queries/mutations to HTTP and subscriptions to WebSocket
const splitLink = split(
  ({ query }) => {
    const definition = getMainDefinition(query);
    return (
      definition.kind === 'OperationDefinition' &&
      definition.operation === 'subscription'
    );
  },
  wsLink,
  from([errorLink, authLink, httpLink])
);

// Apollo Client instance
export const apolloClient = new ApolloClient({
  link: splitLink,
  cache: new InMemoryCache({
    typePolicies: {
      Query: {
        fields: {
          posts: {
            keyArgs: ['filters'],
            merge(existing, incoming, { args }) {
              const merged = existing ? { ...existing } : { edges: [], pageInfo: {}, totalCount: 0 };
              
              if (args?.pagination?.page === 1) {
                // First page - replace existing data
                return incoming;
              } else {
                // Subsequent pages - append to existing data
                return {
                  ...incoming,
                  edges: [...merged.edges, ...incoming.edges],
                };
              }
            },
          },
          searchPosts: {
            keyArgs: ['query'],
            merge(existing, incoming, { args }) {
              // Always replace search results
              return incoming;
            },
          },
        },
      },
      User: {
        keyFields: ['id'],
        fields: {
          // Cache user data for 5 minutes
          __typename: {
            merge: true,
          },
        },
      },
      Post: {
        keyFields: ['id'],
        fields: {
          // Merge post updates intelligently
          author: {
            merge(existing, incoming) {
              return incoming || existing;
            },
          },
          tags: {
            merge(existing, incoming) {
              return incoming || existing;
            },
          },
        },
      },
      Comment: {
        keyFields: ['id'],
        fields: {
          author: {
            merge(existing, incoming) {
              return incoming || existing;
            },
          },
          post: {
            merge(existing, incoming) {
              return incoming || existing;
            },
          },
        },
      },
      PostConnection: {
        keyFields: false, // Don't cache connections by ID
        fields: {
          edges: {
            merge(existing = [], incoming = []) {
              return incoming;
            },
          },
        },
      },
      PostEdge: {
        keyFields: ['cursor'],
      },
    },
    // Garbage collection settings
    possibleTypes: {
      // Define possible types for interfaces/unions if any
    },
    // Cache data for 5 minutes by default
    dataIdFromObject: (object: any) => {
      switch (object.__typename) {
        case 'User':
        case 'Post':
        case 'Comment':
          return `${object.__typename}:${object.id}`;
        default:
          return null;
      }
    },
  }),
  defaultOptions: {
    watchQuery: {
      errorPolicy: 'partial',
      notifyOnNetworkStatusChange: true,
      fetchPolicy: 'cache-first', // Use cache first for better performance
    },
    query: {
      errorPolicy: 'partial',
      fetchPolicy: 'cache-first',
    },
    mutate: {
      errorPolicy: 'partial',
      // Update cache optimistically for mutations
      optimisticResponse: false,
    },
  },
  // Enable query deduplication
  queryDeduplication: true,
});

// Helper functions for token management
export const getAuthToken = (): string | null => {
  return localStorage.getItem('authToken');
};

export const setAuthToken = (token: string): void => {
  localStorage.setItem('authToken', token);
};

export const removeAuthToken = (): void => {
  localStorage.removeItem('authToken');
  localStorage.removeItem('user');
};

export const isAuthenticated = (): boolean => {
  return !!getAuthToken();
};

// Helper to clear Apollo cache on logout
export const clearApolloCache = async (): Promise<void> => {
  await apolloClient.clearStore();
};