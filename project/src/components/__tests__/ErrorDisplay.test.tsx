import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { ApolloError } from '@apollo/client';
import { GraphQLError } from 'graphql';
import { ErrorDisplay, InlineError, LoadingWithError } from '../ErrorDisplay';

import { vi } from 'vitest';

// Mock the retry hook
vi.mock('../../hooks/useRetry', () => ({
  useApolloRetry: () => ({
    retryQuery: vi.fn(),
    isRetrying: false,
  }),
}));

describe('ErrorDisplay', () => {
  it('renders error message for regular errors', () => {
    const error = new Error('Test error message');
    
    render(<ErrorDisplay error={error} />);
    
    expect(screen.getByText('Error')).toBeInTheDocument();
    expect(screen.getByText('Test error message')).toBeInTheDocument();
  });

  it('renders authentication error with appropriate styling', () => {
    const apolloError = new ApolloError({
      graphQLErrors: [
        new GraphQLError('Authentication required', {
          extensions: { code: 'UNAUTHENTICATED' },
        }),
      ],
    });
    
    render(<ErrorDisplay error={apolloError} />);
    
    expect(screen.getByText('Warning')).toBeInTheDocument();
    expect(screen.getByText('Please log in to continue.')).toBeInTheDocument();
  });

  it('renders validation error with info styling', () => {
    const apolloError = new ApolloError({
      graphQLErrors: [
        new GraphQLError('Invalid input', {
          extensions: { code: 'VALIDATION_FAILED' },
        }),
      ],
    });
    
    render(<ErrorDisplay error={apolloError} />);
    
    expect(screen.getByText('Information')).toBeInTheDocument();
    expect(screen.getByText('Invalid input')).toBeInTheDocument();
  });

  it('renders network error message', () => {
    const apolloError = new ApolloError({
      networkError: {
        name: 'NetworkError',
        message: 'Network failure',
        statusCode: 500,
      } as any,
    });
    
    render(<ErrorDisplay error={apolloError} />);
    
    expect(screen.getByText('Error')).toBeInTheDocument();
    expect(screen.getByText('Server error. Please try again later.')).toBeInTheDocument();
  });

  it('shows retry button when onRetry is provided', () => {
    const error = new Error('Test error');
    const onRetry = vi.fn();
    
    render(<ErrorDisplay error={error} onRetry={onRetry} />);
    
    const retryButton = screen.getByText('Try again');
    expect(retryButton).toBeInTheDocument();
    
    fireEvent.click(retryButton);
    expect(onRetry).toHaveBeenCalled();
  });

  it('shows technical details when showDetails is true', () => {
    const apolloError = new ApolloError({
      graphQLErrors: [
        new GraphQLError('GraphQL error', {
          extensions: { code: 'INTERNAL_ERROR' },
        }),
      ],
      networkError: {
        name: 'NetworkError',
        message: 'Network failure',
      } as any,
    });
    
    render(<ErrorDisplay error={apolloError} showDetails={true} />);
    
    expect(screen.getByText('Technical Details')).toBeInTheDocument();
  });

  it('handles 404 network errors correctly', () => {
    const apolloError = new ApolloError({
      networkError: {
        name: 'NetworkError',
        message: 'Not found',
        statusCode: 404,
      } as any,
    });
    
    render(<ErrorDisplay error={apolloError} />);
    
    expect(screen.getByText('The requested resource was not found.')).toBeInTheDocument();
  });

  it('handles rate limit errors correctly', () => {
    const apolloError = new ApolloError({
      networkError: {
        name: 'NetworkError',
        message: 'Too many requests',
        statusCode: 429,
      } as any,
    });
    
    render(<ErrorDisplay error={apolloError} />);
    
    expect(screen.getByText('Too many requests. Please wait and try again.')).toBeInTheDocument();
  });

  it('handles network connectivity errors', () => {
    const apolloError = new ApolloError({
      networkError: {
        name: 'NetworkError',
        message: 'Network error',
        code: 'NETWORK_ERROR',
      } as any,
    });
    
    render(<ErrorDisplay error={apolloError} />);
    
    expect(screen.getByText('Network connection failed. Please check your internet connection.')).toBeInTheDocument();
  });

  it('handles timeout errors', () => {
    const apolloError = new ApolloError({
      networkError: {
        name: 'NetworkError',
        message: 'Timeout',
        code: 'TIMEOUT',
      } as any,
    });
    
    render(<ErrorDisplay error={apolloError} />);
    
    expect(screen.getByText('Request timed out. Please try again.')).toBeInTheDocument();
  });
});

describe('InlineError', () => {
  it('renders inline error message', () => {
    render(<InlineError message="Inline error message" />);
    
    expect(screen.getByText('Inline error message')).toBeInTheDocument();
    expect(screen.getByText('Inline error message')).toHaveClass('text-red-600');
  });

  it('applies custom className', () => {
    render(<InlineError message="Test message" className="custom-class" />);
    
    const errorElement = screen.getByText('Test message');
    expect(errorElement).toHaveClass('custom-class');
  });
});

describe('LoadingWithError', () => {
  it('shows loading state', () => {
    render(
      <LoadingWithError loading={true}>
        <div>Content</div>
      </LoadingWithError>
    );
    
    expect(screen.getByText('Loading...')).toBeInTheDocument();
    expect(screen.queryByText('Content')).not.toBeInTheDocument();
  });

  it('shows error state', () => {
    const error = new Error('Test error');
    
    render(
      <LoadingWithError loading={false} error={error}>
        <div>Content</div>
      </LoadingWithError>
    );
    
    expect(screen.getByText('Error')).toBeInTheDocument();
    expect(screen.queryByText('Content')).not.toBeInTheDocument();
  });

  it('shows content when not loading and no error', () => {
    render(
      <LoadingWithError loading={false}>
        <div>Content</div>
      </LoadingWithError>
    );
    
    expect(screen.getByText('Content')).toBeInTheDocument();
    expect(screen.queryByText('Loading...')).not.toBeInTheDocument();
  });

  it('calls onRetry when retry button is clicked', () => {
    const error = new Error('Test error');
    const onRetry = vi.fn();
    
    render(
      <LoadingWithError loading={false} error={error} onRetry={onRetry}>
        <div>Content</div>
      </LoadingWithError>
    );
    
    const retryButton = screen.getByText('Try again');
    fireEvent.click(retryButton);
    
    expect(onRetry).toHaveBeenCalled();
  });
});