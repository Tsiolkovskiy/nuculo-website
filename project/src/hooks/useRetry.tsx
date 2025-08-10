import React, { useState, useCallback, useEffect } from 'react';
import { ApolloError } from '@apollo/client';

interface RetryOptions {
  maxAttempts?: number;
  delay?: number;
  backoff?: 'linear' | 'exponential';
  retryCondition?: (error: Error) => boolean;
}

interface RetryState {
  isRetrying: boolean;
  attemptCount: number;
  lastError: Error | null;
}

export const useRetry = (options: RetryOptions = {}) => {
  const {
    maxAttempts = 3,
    delay = 1000,
    backoff = 'exponential',
    retryCondition = (error) => {
      // Default: retry on network errors but not on GraphQL errors
      if (error instanceof ApolloError) {
        return !!error.networkError && !error.graphQLErrors.length;
      }
      return true;
    },
  } = options;

  const [retryState, setRetryState] = useState<RetryState>({
    isRetrying: false,
    attemptCount: 0,
    lastError: null,
  });

  const calculateDelay = useCallback(
    (attempt: number) => {
      if (backoff === 'exponential') {
        return delay * Math.pow(2, attempt - 1);
      }
      return delay * attempt;
    },
    [delay, backoff]
  );

  const retry = useCallback(
    async <T,>(operation: () => Promise<T>): Promise<T> => {
      let lastError: Error;
      
      for (let attempt = 1; attempt <= maxAttempts; attempt++) {
        try {
          setRetryState({
            isRetrying: attempt > 1,
            attemptCount: attempt,
            lastError: null,
          });

          const result = await operation();
          
          // Success - reset state
          setRetryState({
            isRetrying: false,
            attemptCount: 0,
            lastError: null,
          });
          
          return result;
        } catch (error) {
          lastError = error as Error;
          
          setRetryState({
            isRetrying: true,
            attemptCount: attempt,
            lastError,
          });

          // Check if we should retry
          if (attempt === maxAttempts || !retryCondition(lastError)) {
            break;
          }

          // Wait before retrying
          const retryDelay = calculateDelay(attempt);
          await new Promise(resolve => setTimeout(resolve, retryDelay));
        }
      }

      // All attempts failed
      setRetryState({
        isRetrying: false,
        attemptCount: maxAttempts,
        lastError,
      });

      throw lastError;
    },
    [maxAttempts, retryCondition, calculateDelay]
  );

  const reset = useCallback(() => {
    setRetryState({
      isRetrying: false,
      attemptCount: 0,
      lastError: null,
    });
  }, []);

  return {
    retry,
    reset,
    ...retryState,
  };
};

// Hook for retrying Apollo operations
export const useApolloRetry = (options: RetryOptions = {}) => {
  const { retry, ...retryState } = useRetry({
    ...options,
    retryCondition: (error) => {
      // Retry on network errors, timeouts, and 5xx server errors
      if (error instanceof ApolloError) {
        if (error.networkError) {
          const networkError = error.networkError as any;
          
          // Retry on network failures
          if (networkError.code === 'NETWORK_ERROR') {
            return true;
          }
          
          // Retry on timeout
          if (networkError.code === 'TIMEOUT') {
            return true;
          }
          
          // Retry on 5xx server errors
          if (networkError.statusCode >= 500) {
            return true;
          }
          
          // Don't retry on 4xx client errors
          if (networkError.statusCode >= 400 && networkError.statusCode < 500) {
            return false;
          }
        }
        
        // Don't retry on GraphQL errors (business logic errors)
        if (error.graphQLErrors.length > 0) {
          return false;
        }
      }
      
      return true;
    },
  });

  const retryQuery = useCallback(
    async (queryFn: () => Promise<any>) => {
      return retry(queryFn);
    },
    [retry]
  );

  const retryMutation = useCallback(
    async (mutationFn: () => Promise<any>) => {
      return retry(mutationFn);
    },
    [retry]
  );

  return {
    retryQuery,
    retryMutation,
    ...retryState,
  };
};

// Hook for handling network status
export const useNetworkStatus = () => {
  const [isOnline, setIsOnline] = useState(navigator.onLine);
  const [wasOffline, setWasOffline] = useState(false);

  useEffect(() => {
    const handleOnline = () => {
      setIsOnline(true);
      if (wasOffline) {
        // Trigger refetch of queries when coming back online
        window.location.reload();
      }
    };

    const handleOffline = () => {
      setIsOnline(false);
      setWasOffline(true);
    };

    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);

    return () => {
      window.removeEventListener('online', handleOnline);
      window.removeEventListener('offline', handleOffline);
    };
  }, [wasOffline]);

  return {
    isOnline,
    wasOffline,
  };
};

// Component for displaying network status
export const NetworkStatus: React.FC = () => {
  const { isOnline, wasOffline } = useNetworkStatus();

  if (isOnline && !wasOffline) {
    return null;
  }

  const statusClasses = isOnline
    ? 'bg-green-500 text-white'
    : 'bg-red-500 text-white';

  return (
    <div className={`fixed top-0 left-0 right-0 z-50 p-2 text-center text-sm font-medium ${statusClasses}`}>
      {isOnline ? (
        <span>✅ Connection restored</span>
      ) : (
        <span>⚠️ No internet connection</span>
      )}
    </div>
  );
};