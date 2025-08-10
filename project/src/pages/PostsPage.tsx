import React from 'react';
import { PostsList } from '../components/PostsList';

export const PostsPage: React.FC = () => {
  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">All Posts</h1>
        <p className="text-gray-600 mt-2">
          Browse through our collection of blog posts and articles.
        </p>
      </div>

      <PostsList limit={20} showPublishedOnly={true} />
    </div>
  );
};