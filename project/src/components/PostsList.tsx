import React, { useState } from 'react';
import { useGetPostsQuery } from '../generated/graphql';
import { LoadingSpinner } from './LoadingSpinner';

interface PostsListProps {
  limit?: number;
  showPublishedOnly?: boolean;
  onPostClick?: (postId: string) => void;
}

export const PostsList: React.FC<PostsListProps> = ({ 
  limit = 10, 
  showPublishedOnly = true,
  onPostClick
}) => {
  const [currentPage, setCurrentPage] = useState(1);
  
  const { data, loading, error, refetch, fetchMore } = useGetPostsQuery({
    variables: {
      filters: {
        published: showPublishedOnly,
      },
      pagination: {
        page: currentPage,
        limit,
      },
    },
    errorPolicy: 'partial',
    notifyOnNetworkStatusChange: true,
  });

  const handleLoadMore = async () => {
    if (data?.posts?.pageInfo?.hasNextPage) {
      try {
        await fetchMore({
          variables: {
            filters: {
              published: showPublishedOnly,
            },
            pagination: {
              page: currentPage + 1,
              limit,
            },
          },
        });
        setCurrentPage(prev => prev + 1);
      } catch (error) {
        console.error('Failed to load more posts:', error);
      }
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center py-8">
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-md p-4">
        <div className="flex">
          <div className="flex-shrink-0">
            <svg className="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
            </svg>
          </div>
          <div className="ml-3">
            <h3 className="text-sm font-medium text-red-800">Error loading posts</h3>
            <p className="text-sm text-red-700 mt-1">{error.message}</p>
            <button
              onClick={() => refetch()}
              className="mt-2 text-sm text-red-600 hover:text-red-500 underline"
            >
              Try again
            </button>
          </div>
        </div>
      </div>
    );
  }

  const posts = data?.posts?.edges || [];

  if (posts.length === 0) {
    return (
      <div className="text-center py-8">
        <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <h3 className="mt-2 text-sm font-medium text-gray-900">No posts found</h3>
        <p className="mt-1 text-sm text-gray-500">
          {showPublishedOnly ? 'No published posts available.' : 'No posts available.'}
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {posts.map(({ node: post }) => (
        <article key={post.id} className="bg-white rounded-lg shadow-sm border p-6">
          <div className="flex items-center mb-4">
            <div className="flex-shrink-0">
              {post.author.avatar ? (
                <img
                  src={post.author.avatar}
                  alt={post.author.name}
                  className="w-10 h-10 rounded-full"
                />
              ) : (
                <div className="w-10 h-10 bg-gray-300 rounded-full flex items-center justify-center">
                  <span className="text-sm font-medium text-gray-600">
                    {post.author.name.charAt(0).toUpperCase()}
                  </span>
                </div>
              )}
            </div>
            <div className="ml-4">
              <h4 className="text-sm font-medium text-gray-900">{post.author.name}</h4>
              <p className="text-sm text-gray-500">
                {new Date(post.createdAt).toLocaleDateString()}
              </p>
            </div>
          </div>
          
          <h2 className="text-xl font-semibold text-gray-900 mb-3">
            {post.title}
          </h2>
          
          <p className="text-gray-600 mb-4 line-clamp-3">
            {post.content}
          </p>
          
          {post.tags.length > 0 && (
            <div className="flex flex-wrap gap-2 mb-4">
              {post.tags.map((tag, index) => (
                <span
                  key={index}
                  className="px-2 py-1 text-xs font-medium bg-blue-100 text-blue-800 rounded-full"
                >
                  {tag}
                </span>
              ))}
            </div>
          )}
          
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4">
              <span className={`px-2 py-1 text-xs font-medium rounded-full ${
                post.published 
                  ? 'bg-green-100 text-green-800' 
                  : 'bg-yellow-100 text-yellow-800'
              }`}>
                {post.published ? 'Published' : 'Draft'}
              </span>
            </div>
            
            <button 
              onClick={() => onPostClick?.(post.id)}
              className="text-blue-600 hover:text-blue-800 text-sm font-medium"
            >
              Read more →
            </button>
          </div>
        </article>
      ))}
      
      {data?.posts?.pageInfo?.hasNextPage && (
        <div className="text-center">
          <button 
            onClick={handleLoadMore}
            disabled={loading}
            className="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {loading ? (
              <>
                <LoadingSpinner size="sm" className="mr-2" />
                Loading...
              </>
            ) : (
              'Load more posts'
            )}
          </button>
        </div>
      )}
    </div>
  );
};