import React from 'react';
import { useNavigate } from 'react-router-dom';
import { CreatePostForm } from '../components/CreatePostForm';

export const CreatePostPage: React.FC = () => {
  const navigate = useNavigate();

  const handleSuccess = () => {
    navigate('/dashboard');
  };

  const handleCancel = () => {
    navigate('/dashboard');
  };

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">Create New Post</h1>
        <p className="text-gray-600 mt-2">
          Share your thoughts and ideas with the community.
        </p>
      </div>

      <CreatePostForm onSuccess={handleSuccess} onCancel={handleCancel} />
    </div>
  );
};