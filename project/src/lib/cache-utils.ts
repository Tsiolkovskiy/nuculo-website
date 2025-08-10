import { apolloClient } from './apollo';
import { GetPostsDocument, GetPostsQuery } from '../generated/graphql';

// Cache management utilities for Apollo Client

/**
 * Evict a specific item from the cache
 */
export const evictFromCache = (typename: string, id: string) => {
  apolloClient.cache.evict({
    id: apolloClient.cache.identify({ __typename: typename, id }),
  });
  apolloClient.cache.gc(); // Garbage collect
};

/**
 * Update a specific item in the cache
 */
export const updateCacheItem = (typename: string, id: string, data: any) => {
  const cacheId = apolloClient.cache.identify({ __typename: typename, id });
  if (cacheId) {
    apolloClient.cache.writeFragment({
      id: cacheId,
      fragment: gql`
        fragment Update${typename} on ${typename} {
          __typename
          ${Object.keys(data).join('\n          ')}
        }
      `,
      data: {
        __typename: typename,
        ...data,
      },
    });
  }
};

/**
 * Add a new post to the posts list cache
 */
export const addPostToCache = (newPost: any, filters?: any) => {
  try {
    const existingData = apolloClient.cache.readQuery<GetPostsQuery>({
      query: GetPostsDocument,
      variables: {
        filters: filters || { published: true },
        pagination: { page: 1, limit: 20 },
      },
    });

    if (existingData?.posts) {
      apolloClient.cache.writeQuery({
        query: GetPostsDocument,
        variables: {
          filters: filters || { published: true },
          pagination: { page: 1, limit: 20 },
        },
        data: {
          posts: {
            ...existingData.posts,
            edges: [
              { node: newPost, cursor: newPost.id, __typename: 'PostEdge' },
              ...existingData.posts.edges,
            ],
            totalCount: existingData.posts.totalCount + 1,
          },
        },
      });
    }
  } catch (error) {
    console.warn('Failed to update posts cache:', error);
  }
};

/**
 * Remove a post from the posts list cache
 */
export const removePostFromCache = (postId: string, filters?: any) => {
  try {
    const existingData = apolloClient.cache.readQuery<GetPostsQuery>({
      query: GetPostsDocument,
      variables: {
        filters: filters || { published: true },
        pagination: { page: 1, limit: 20 },
      },
    });

    if (existingData?.posts) {
      apolloClient.cache.writeQuery({
        query: GetPostsDocument,
        variables: {
          filters: filters || { published: true },
          pagination: { page: 1, limit: 20 },
        },
        data: {
          posts: {
            ...existingData.posts,
            edges: existingData.posts.edges.filter(
              (edge) => edge.node.id !== postId
            ),
            totalCount: Math.max(0, existingData.posts.totalCount - 1),
          },
        },
      });
    }
  } catch (error) {
    console.warn('Failed to remove post from cache:', error);
  }
};

/**
 * Clear all cached data
 */
export const clearAllCache = async () => {
  await apolloClient.clearStore();
};

/**
 * Reset the Apollo Client store
 */
export const resetStore = async () => {
  await apolloClient.resetStore();
};

/**
 * Get cache size information
 */
export const getCacheInfo = () => {
  const cache = apolloClient.cache as any;
  const data = cache.data?.data || {};
  const size = Object.keys(data).length;
  
  return {
    size,
    keys: Object.keys(data),
    data: data,
  };
};

/**
 * Prefetch data for better performance
 */
export const prefetchPosts = async (filters?: any) => {
  try {
    await apolloClient.query({
      query: GetPostsDocument,
      variables: {
        filters: filters || { published: true },
        pagination: { page: 1, limit: 20 },
      },
      fetchPolicy: 'cache-first',
    });
  } catch (error) {
    console.warn('Failed to prefetch posts:', error);
  }
};

/**
 * Cache warming utility
 */
export const warmCache = async () => {
  // Prefetch commonly accessed data
  await Promise.allSettled([
    prefetchPosts({ published: true }),
    // Add other prefetch operations as needed
  ]);
};

/**
 * Cache debugging utilities
 */
export const debugCache = () => {
  const info = getCacheInfo();
  console.group('🗄️ Apollo Cache Debug Info');
  console.log('Cache size:', info.size, 'items');
  console.log('Cache keys:', info.keys);
  console.log('Cache data:', info.data);
  console.groupEnd();
  return info;
};

/**
 * Monitor cache performance
 */
export const monitorCachePerformance = () => {
  const originalReadQuery = apolloClient.cache.readQuery;
  const originalWriteQuery = apolloClient.cache.writeQuery;
  
  let readCount = 0;
  let writeCount = 0;
  
  apolloClient.cache.readQuery = function(...args) {
    readCount++;
    console.log(`📖 Cache read #${readCount}`);
    return originalReadQuery.apply(this, args);
  };
  
  apolloClient.cache.writeQuery = function(...args) {
    writeCount++;
    console.log(`✏️ Cache write #${writeCount}`);
    return originalWriteQuery.apply(this, args);
  };
  
  return {
    getStats: () => ({ reads: readCount, writes: writeCount }),
    reset: () => {
      readCount = 0;
      writeCount = 0;
    },
  };
};