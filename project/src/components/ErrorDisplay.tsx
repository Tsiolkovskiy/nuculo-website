import React from 'react';
import { ApolloError } from '@apollo/client';
import { useApolloRetry } from '../hooks/useRetry';

interface ErrorDisplayProps {
  error: Error | ApolloError;
  onRetry?: () => void;
  showDetails?: boolean;
  className?: string;
}

export const ErrorDisplay: React.FC<ErrorDisplayProps> = ({
  error,
  onRetry,
  showDetails = false,
  className = '',
}) => {
  const { retryQuery, isRetrying } = useApolloRetry();

  const handleRetry = async () => {
    if (onRetry) {
      if (error instanceof ApolloError) {
        try {
          await retryQuery(onRetry);
        } catch (retryError) {
          console.error('Retry failed:', retryError);
        }
      } else {
        onRetry();
      }
    }
  };

  const getErrorMessage = (error: Error | ApolloError): string => {
    if (error instanceof ApolloError) {
      // Handle GraphQL errors
      if (error.graphQLErrors.length > 0) {
        const firstError = error.graphQLErrors[0];
        const code = firstError.extensions?.code as string;
        
        switch (code) {
          case 'UNAUTHENTICATED':
            return 'Please log in to continue.';
          case 'UNAUTHORIZED':
          case 'FORBIDDEN':
            return 'You do not have permission to perform this action.';
          case 'NOT_FOUND':
            return 'The requested resource was not found.';
          case 'VALIDATION_FAILED':
          case 'INVALID_INPUT':
            return firstError.message || 'Please check your input and try again.';
          case 'RATE_LIMIT_EXCEEDED':
            return 'Too many requests. Please wait a moment and try again.';
          default:
            return firstError.message || 'An error occurred while processing your request.';
        }
      }
      
      // Handle network errors
      if (error.networkError) {
        const networkError = error.networkError as any;
        
        if (networkError.statusCode) {
          switch (networkError.statusCode) {
            case 400:
              return 'Invalid request. Please check your input.';
            case 401:
              return 'Authentication required. Please log in.';
            case 403:
              return 'Access denied. You do not have permission.';
            case 404:
              return 'The requested resource was not found.';
            case 429:
              return 'Too many requests. Please wait and try again.';
            case 500:
              return 'Server error. Please try again later.';
            case 502:
            case 503:
            case 504:
              return 'Service temporarily unavailable. Please try again later.';
            default:
              return 'Network error occurred. Please check your connection.';
          }
        }
        
        if (networkError.code === 'NETWORK_ERROR') {
          return 'Network connection failed. Please check your internet connection.';
        }
        
        if (networkError.code === 'TIMEOUT') {
          return 'Request timed out. Please try again.';
        }
        
        return 'Network error occurred. Please try again.';
      }
      
      return 'An unexpected error occurred. Please try again.';
    }
    
    // Handle regular errors
    return error.message || 'An unexpected error occurred.';
  };

  const getErrorType = (error: Error | ApolloError): 'error' | 'warning' | 'info' => {
    if (error instanceof ApolloError) {
      if (error.graphQLErrors.length > 0) {
        const code = error.graphQLErrors[0].extensions?.code as string;
        
        switch (code) {
          case 'UNAUTHENTICATED':
          case 'UNAUTHORIZED':
          case 'FORBIDDEN':
            return 'warning';
          case 'VALIDATION_FAILED':
          case 'INVALID_INPUT':
          case 'NOT_FOUND':
            return 'info';
          default:
            return 'error';
        }
      }
      
      if (error.networkError) {
        const networkError = error.networkError as any;
        if (networkError.statusCode >= 400 && networkError.statusCode < 500) {
          return 'warning';
        }
      }
    }
    
    return 'error';
  };

  const errorType = getErrorType(error);
  const errorMessage = getErrorMessage(error);

  const getIconAndColors = (type: 'error' | 'warning' | 'info') => {
    switch (type) {
      case 'error':
        return {
          bgColor: 'bg-red-50',
          borderColor: 'border-red-200',
          iconColor: 'text-red-400',
          textColor: 'text-red-800',
          buttonColor: 'text-red-600 hover:text-red-500',
          icon: (
            <svg className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
              <path
                fillRule="evenodd"
                d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                clipRule="evenodd"
              />
            </svg>
          ),
        };
      case 'warning':
        return {
          bgColor: 'bg-yellow-50',
          borderColor: 'border-yellow-200',
          iconColor: 'text-yellow-400',
          textColor: 'text-yellow-800',
          buttonColor: 'text-yellow-600 hover:text-yellow-500',
          icon: (
            <svg className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
              <path
                fillRule="evenodd"
                d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z"
                clipRule="evenodd"
              />
            </svg>
          ),
        };
      case 'info':
        return {
          bgColor: 'bg-blue-50',
          borderColor: 'border-blue-200',
          iconColor: 'text-blue-400',
          textColor: 'text-blue-800',
          buttonColor: 'text-blue-600 hover:text-blue-500',
          icon: (
            <svg className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
              <path
                fillRule="evenodd"
                d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
                clipRule="evenodd"
              />
            </svg>
          ),
        };
    }
  };

  const { bgColor, borderColor, iconColor, textColor, buttonColor, icon } = getIconAndColors(errorType);

  return (
    <div className={`rounded-md p-4 ${bgColor} border ${borderColor} ${className}`}>
      <div className="flex">
        <div className="flex-shrink-0">
          <div className={iconColor}>{icon}</div>
        </div>
        <div className="ml-3 flex-1">
          <h3 className={`text-sm font-medium ${textColor}`}>
            {errorType === 'error' ? 'Error' : errorType === 'warning' ? 'Warning' : 'Information'}
          </h3>
          <div className={`mt-2 text-sm ${textColor}`}>
            <p>{errorMessage}</p>
            
            {showDetails && error instanceof ApolloError && (
              <details className="mt-2">
                <summary className="cursor-pointer text-xs opacity-75">
                  Technical Details
                </summary>
                <div className="mt-1 text-xs opacity-75">
                  {error.graphQLErrors.length > 0 && (
                    <div>
                      <strong>GraphQL Errors:</strong>
                      <ul className="list-disc list-inside ml-2">
                        {error.graphQLErrors.map((err, index) => (
                          <li key={index}>
                            {err.message}
                            {err.extensions?.code && ` (${err.extensions.code})`}
                          </li>
                        ))}
                      </ul>
                    </div>
                  )}
                  {error.networkError && (
                    <div className="mt-1">
                      <strong>Network Error:</strong> {error.networkError.message}
                    </div>
                  )}
                </div>
              </details>
            )}
          </div>
          
          {onRetry && (
            <div className="mt-3">
              <button
                onClick={handleRetry}
                disabled={isRetrying}
                className={`text-sm ${buttonColor} underline disabled:opacity-50 disabled:cursor-not-allowed`}
              >
                {isRetrying ? 'Retrying...' : 'Try again'}
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

// Simplified error display for inline use
export const InlineError: React.FC<{ message: string; className?: string }> = ({
  message,
  className = '',
}) => (
  <div className={`text-sm text-red-600 ${className}`}>
    {message}
  </div>
);

// Loading state with error fallback
export const LoadingWithError: React.FC<{
  loading: boolean;
  error?: Error | ApolloError;
  onRetry?: () => void;
  children: React.ReactNode;
}> = ({ loading, error, onRetry, children }) => {
  if (loading) {
    return (
      <div className="flex items-center justify-center p-8">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <span className="ml-2 text-gray-600">Loading...</span>
      </div>
    );
  }

  if (error) {
    return <ErrorDisplay error={error} onRetry={onRetry} className="m-4" />;
  }

  return <>{children}</>;
};