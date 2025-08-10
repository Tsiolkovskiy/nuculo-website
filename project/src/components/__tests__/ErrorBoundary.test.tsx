import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { ApolloError } from '@apollo/client';
import { GraphQLError } from 'graphql';
import { ErrorBoundary, GraphQLErrorBoundary } from '../ErrorBoundary';

// Mock component that throws an error
const ThrowError: React.FC<{ shouldThrow?: boolean; error?: Error }> = ({ 
  shouldThrow = true, 
  error = new Error('Test error') 
}) => {
  if (shouldThrow) {
    throw error;
  }
  return <div>No error</div>;
};

describe('ErrorBoundary', () => {
import { vi } from 'vitest';

  // Suppress console.error for these tests
  const originalError = console.error;
  beforeAll(() => {
    console.error = vi.fn();
  });
  afterAll(() => {
    console.error = originalError;
  });

  it('renders children when there is no error', () => {
    render(
      <ErrorBoundary>
        <ThrowError shouldThrow={false} />
      </ErrorBoundary>
    );

    expect(screen.getByText('No error')).toBeInTheDocument();
  });

  it('renders error UI when there is an error', () => {
    render(
      <ErrorBoundary>
        <ThrowError />
      </ErrorBoundary>
    );

    expect(screen.getByText('Something went wrong')).toBeInTheDocument();
    expect(screen.getByText('We\'re sorry, but something unexpected happened. Please try again.')).toBeInTheDocument();
  });

  it('shows error details in development mode', () => {
    const originalEnv = process.env.NODE_ENV;
    process.env.NODE_ENV = 'development';

    render(
      <ErrorBoundary>
        <ThrowError error={new Error('Detailed test error')} />
      </ErrorBoundary>
    );

    expect(screen.getByText('Error Details (Development Only)')).toBeInTheDocument();
    expect(screen.getByText('Detailed test error')).toBeInTheDocument();

    process.env.NODE_ENV = originalEnv;
  });

  it('hides error details in production mode', () => {
    const originalEnv = process.env.NODE_ENV;
    process.env.NODE_ENV = 'production';

    render(
      <ErrorBoundary>
        <ThrowError error={new Error('Production test error')} />
      </ErrorBoundary>
    );

    expect(screen.queryByText('Error Details (Development Only)')).not.toBeInTheDocument();
    expect(screen.queryByText('Production test error')).not.toBeInTheDocument();

    process.env.NODE_ENV = originalEnv;
  });

  it('calls onError callback when error occurs', () => {
    const onError = vi.fn();

    render(
      <ErrorBoundary onError={onError}>
        <ThrowError />
      </ErrorBoundary>
    );

    expect(onError).toHaveBeenCalledWith(
      expect.any(Error),
      expect.objectContaining({
        componentStack: expect.any(String),
      })
    );
  });

  it('renders custom fallback UI when provided', () => {
    const customFallback = <div>Custom error message</div>;

    render(
      <ErrorBoundary fallback={customFallback}>
        <ThrowError />
      </ErrorBoundary>
    );

    expect(screen.getByText('Custom error message')).toBeInTheDocument();
    expect(screen.queryByText('Something went wrong')).not.toBeInTheDocument();
  });

  it('resets error state when retry button is clicked', () => {
    const { rerender } = render(
      <ErrorBoundary>
        <ThrowError />
      </ErrorBoundary>
    );

    expect(screen.getByText('Something went wrong')).toBeInTheDocument();

    fireEvent.click(screen.getByText('Try Again'));

    // Re-render with no error
    rerender(
      <ErrorBoundary>
        <ThrowError shouldThrow={false} />
      </ErrorBoundary>
    );

    expect(screen.getByText('No error')).toBeInTheDocument();
  });
});

describe('GraphQLErrorBoundary', () => {
  const originalError = console.error;
  beforeAll(() => {
    console.error = vi.fn();
  });
  afterAll(() => {
    console.error = originalError;
  });

  it('renders children when there is no error', () => {
    render(
      <GraphQLErrorBoundary>
        <ThrowError shouldThrow={false} />
      </GraphQLErrorBoundary>
    );

    expect(screen.getByText('No error')).toBeInTheDocument();
  });

  it('handles Apollo GraphQL errors', () => {
    const graphQLError = new GraphQLError('GraphQL test error');
    const apolloError = new ApolloError({
      graphQLErrors: [graphQLError],
    });

    render(
      <GraphQLErrorBoundary>
        <ThrowError error={apolloError} />
      </GraphQLErrorBoundary>
    );

    expect(screen.getByText('GraphQL Error')).toBeInTheDocument();
    expect(screen.getByText('GraphQL test error')).toBeInTheDocument();
  });

  it('handles Apollo network errors', () => {
    const apolloError = new ApolloError({
      networkError: new Error('Network test error'),
    });

    render(
      <GraphQLErrorBoundary>
        <ThrowError error={apolloError} />
      </GraphQLErrorBoundary>
    );

    expect(screen.getByText('GraphQL Error')).toBeInTheDocument();
    expect(screen.getByText('Network test error')).toBeInTheDocument();
  });

  it('renders custom fallback for Apollo errors', () => {
    const customFallback = (error: ApolloError) => (
      <div>Custom Apollo error: {error.message}</div>
    );

    const apolloError = new ApolloError({
      graphQLErrors: [new GraphQLError('Custom test error')],
    });

    render(
      <GraphQLErrorBoundary fallback={customFallback}>
        <ThrowError error={apolloError} />
      </GraphQLErrorBoundary>
    );

    expect(screen.getByText(/Custom Apollo error:/)).toBeInTheDocument();
  });

  it('does not catch non-Apollo errors', () => {
    // This should fall through to a parent error boundary
    expect(() => {
      render(
        <GraphQLErrorBoundary>
          <ThrowError error={new Error('Regular error')} />
        </GraphQLErrorBoundary>
      );
    }).toThrow('Regular error');
  });

  it('resets error state when retry is clicked', () => {
    const apolloError = new ApolloError({
      graphQLErrors: [new GraphQLError('Retry test error')],
    });

    const { rerender } = render(
      <GraphQLErrorBoundary>
        <ThrowError error={apolloError} />
      </GraphQLErrorBoundary>
    );

    expect(screen.getByText('GraphQL Error')).toBeInTheDocument();

    fireEvent.click(screen.getByText('Try again'));

    // Re-render with no error
    rerender(
      <GraphQLErrorBoundary>
        <ThrowError shouldThrow={false} />
      </GraphQLErrorBoundary>
    );

    expect(screen.getByText('No error')).toBeInTheDocument();
  });
});