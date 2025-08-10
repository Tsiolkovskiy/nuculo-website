import React, { useState, useEffect } from 'react';
import { useUpdatePostMutation, useGetPostQuery, GetPostDocument, type Post } from '../generated/graphql';
import { LoadingSpinner } from './LoadingSpinner';

interface EditPostFormProps {
  postId: string;
  onSuccess?: () => void;
  onCancel?: () => void;
}

export const EditPostForm: React.FC<EditPostFormProps> = ({ 
  postId,
  onSuccess, 
  onCancel 
}) => {
  const [formData, setFormData] = useState({
    title: '',
    content: '',
    tags: '',
    published: false,
  });

  // Fetch the existing post data
  const { data: postData, loading: postLoading } = useGetPostQuery({
    variables: { id: postId },
    skip: !postId,
  });

  // Update form when post data loads
  useEffect(() => {
    if (postData?.post) {
      const post = postData.post;
      setFormData({
        title: post.title,
        content: post.content,
        tags: post.tags.join(', '),
        published: post.published,
      });
    }
  }, [postData]);

  const [updatePost, { loading, error }] = useUpdatePostMutation({
    // Optimistic response for immediate UI updates
    optimisticResponse: {
      updatePost: {
        __typename: 'Post',
        id: postId,
        title: formData.title,
        content: formData.content,
        tags: formData.tags.split(',').map(tag => tag.trim()).filter(Boolean),
        published: formData.published,
        // Keep existing fields
        author: postData?.post?.author || {
          __typename: 'User',
          id: '',
          name: '',
          email: '',
          avatar: null,
          createdAt: '',
          updatedAt: '',
        },
        createdAt: postData?.post?.createdAt || '',
        updatedAt: new Date().toISOString(),
      },
    },
    // Update cache after successful mutation
    update: (cache, { data }) => {
      if (data?.updatePost) {
        const updatedPost = data.updatePost;
        
        // Update the individual post in cache
        cache.writeQuery({
          query: useGetPostQuery.getQuery?.() || GetPostDocument,
          variables: { id: postId },
          data: {
            post: updatedPost,
          },
        });

        // Update the post in any lists that might contain it
        cache.modify({
          fields: {
            posts(existingPosts = { edges: [], pageInfo: {}, totalCount: 0 }) {
              return {
                ...existingPosts,
                edges: existingPosts.edges.map((edge: any) => 
                  edge.node.id === updatedPost.id 
                    ? { ...edge, node: updatedPost }
                    : edge
                ),
              };
            },
          },
        });
      }
    },
  });

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value, type } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? (e.target as HTMLInputElement).checked : value,
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    try {
      const result = await updatePost({
        variables: {
          id: postId,
          input: {
            title: formData.title,
            content: formData.content,
            tags: formData.tags.split(',').map(tag => tag.trim()).filter(Boolean),
            published: formData.published,
          },
        },
      });

      if (result.data?.updatePost) {
        onSuccess?.();
      }
    } catch (error) {
      console.error('Failed to update post:', error);
    }
  };

  if (postLoading) {
    return (
      <div className="flex justify-center py-8">
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  if (!postData?.post) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-md p-4">
        <p className="text-sm text-red-800">Post not found</p>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="bg-white shadow-sm rounded-lg">
        <div className="px-6 py-4 border-b border-gray-200">
          <h3 className="text-lg font-medium text-gray-900">Edit Post</h3>
        </div>
        
        <div className="px-6 py-4 space-y-6">
          {error && (
            <div className="bg-red-50 border border-red-200 rounded-md p-4">
              <div className="flex">
                <div className="flex-shrink-0">
                  <svg className="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
                    <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
                  </svg>
                </div>
                <div className="ml-3">
                  <p className="text-sm text-red-800">
                    {error.graphQLErrors?.[0]?.message || error.message || 'Failed to update post'}
                  </p>
                </div>
              </div>
            </div>
          )}

          <div>
            <label htmlFor="title" className="block text-sm font-medium text-gray-700">
              Title
            </label>
            <input
              type="text"
              name="title"
              id="title"
              required
              value={formData.title}
              onChange={handleInputChange}
              className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              placeholder="Enter your post title"
              disabled={loading}
            />
          </div>

          <div>
            <label htmlFor="content" className="block text-sm font-medium text-gray-700">
              Content
            </label>
            <textarea
              name="content"
              id="content"
              rows={8}
              required
              value={formData.content}
              onChange={handleInputChange}
              className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              placeholder="Write your post content here..."
              disabled={loading}
            />
          </div>

          <div>
            <label htmlFor="tags" className="block text-sm font-medium text-gray-700">
              Tags
            </label>
            <input
              type="text"
              name="tags"
              id="tags"
              value={formData.tags}
              onChange={handleInputChange}
              className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              placeholder="Enter tags separated by commas (e.g., react, typescript, tutorial)"
              disabled={loading}
            />
            <p className="mt-1 text-sm text-gray-500">
              Add relevant tags to help others discover your post.
            </p>
          </div>

          <div className="flex items-center">
            <input
              id="published"
              name="published"
              type="checkbox"
              checked={formData.published}
              onChange={handleInputChange}
              className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
              disabled={loading}
            />
            <label htmlFor="published" className="ml-2 block text-sm text-gray-900">
              Published
            </label>
          </div>
        </div>
      </div>

      <div className="flex justify-end space-x-4">
        {onCancel && (
          <button
            type="button"
            onClick={onCancel}
            className="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            disabled={loading}
          >
            Cancel
          </button>
        )}
        <button
          type="submit"
          disabled={loading || !formData.title || !formData.content}
          className="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {loading ? (
            <>
              <LoadingSpinner size="sm" className="mr-2" />
              Updating...
            </>
          ) : (
            'Update Post'
          )}
        </button>
      </div>
    </form>
  );
};