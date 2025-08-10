import { renderHook, act } from '@testing-library/react';
import { ApolloError } from '@apollo/client';
import { GraphQLError } from 'graphql';
import { vi } from 'vitest';
import { useRetry, useApolloRetry, useNetworkStatus } from '../useRetry';

// Mock navigator.onLine
Object.defineProperty(navigator, 'onLine', {
  writable: true,
  value: true,
});

// Mock window.addEventListener and removeEventListener
const mockAddEventListener = vi.fn();
const mockRemoveEventListener = vi.fn();
Object.defineProperty(window, 'addEventListener', {
  value: mockAddEventListener,
});
Object.defineProperty(window, 'removeEventListener', {
  value: mockRemoveEventListener,
});

describe('useRetry', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should succeed on first attempt', async () => {
    const { result } = renderHook(() => useRetry());
    const mockOperation = vi.fn().mockResolvedValue('success');

    await act(async () => {
      const response = await result.current.retry(mockOperation);
      expect(response).toBe('success');
    });

    expect(mockOperation).toHaveBeenCalledTimes(1);
    expect(result.current.attemptCount).toBe(0);
    expect(result.current.isRetrying).toBe(false);
    expect(result.current.lastError).toBe(null);
  });

  it('should retry on failure and eventually succeed', async () => {
    const { result } = renderHook(() => useRetry({ maxAttempts: 3, delay: 10 }));
    const mockOperation = vi
      .fn()
      .mockRejectedValueOnce(new Error('First failure'))
      .mockRejectedValueOnce(new Error('Second failure'))
      .mockResolvedValue('success');

    await act(async () => {
      const response = await result.current.retry(mockOperation);
      expect(response).toBe('success');
    });

    expect(mockOperation).toHaveBeenCalledTimes(3);
    expect(result.current.attemptCount).toBe(0);
    expect(result.current.isRetrying).toBe(false);
    expect(result.current.lastError).toBe(null);
  });

  it('should fail after max attempts', async () => {
    const { result } = renderHook(() => useRetry({ maxAttempts: 2, delay: 10 }));
    const error = new Error('Persistent failure');
    const mockOperation = vi.fn().mockRejectedValue(error);

    await act(async () => {
      await expect(result.current.retry(mockOperation)).rejects.toThrow('Persistent failure');
    });

    expect(mockOperation).toHaveBeenCalledTimes(2);
    expect(result.current.attemptCount).toBe(2);
    expect(result.current.isRetrying).toBe(false);
    expect(result.current.lastError).toBe(error);
  });

  it('should respect retry condition', async () => {
    const { result } = renderHook(() => 
      useRetry({ 
        maxAttempts: 3, 
        delay: 10,
        retryCondition: (error) => error.message !== 'Do not retry'
      })
    );
    
    const error = new Error('Do not retry');
    const mockOperation = vi.fn().mockRejectedValue(error);

    await act(async () => {
      await expect(result.current.retry(mockOperation)).rejects.toThrow('Do not retry');
    });

    expect(mockOperation).toHaveBeenCalledTimes(1);
    expect(result.current.attemptCount).toBe(1);
  });

  it('should calculate exponential backoff delay', async () => {
    const { result } = renderHook(() => 
      useRetry({ maxAttempts: 3, delay: 100, backoff: 'exponential' })
    );
    
    const mockOperation = vi
      .fn()
      .mockRejectedValueOnce(new Error('First failure'))
      .mockRejectedValueOnce(new Error('Second failure'))
      .mockResolvedValue('success');

    const startTime = Date.now();
    
    await act(async () => {
      await result.current.retry(mockOperation);
    });

    const endTime = Date.now();
    const totalTime = endTime - startTime;
    
    // Should have waited at least 100ms + 200ms = 300ms
    expect(totalTime).toBeGreaterThan(250);
  });

  it('should calculate linear backoff delay', async () => {
    const { result } = renderHook(() => 
      useRetry({ maxAttempts: 3, delay: 100, backoff: 'linear' })
    );
    
    const mockOperation = vi
      .fn()
      .mockRejectedValueOnce(new Error('First failure'))
      .mockRejectedValueOnce(new Error('Second failure'))
      .mockResolvedValue('success');

    const startTime = Date.now();
    
    await act(async () => {
      await result.current.retry(mockOperation);
    });

    const endTime = Date.now();
    const totalTime = endTime - startTime;
    
    // Should have waited at least 100ms + 200ms = 300ms
    expect(totalTime).toBeGreaterThan(250);
  });

  it('should reset state', async () => {
    const { result } = renderHook(() => useRetry({ maxAttempts: 1 }));
    const error = new Error('Test error');
    const mockOperation = vi.fn().mockRejectedValue(error);

    await act(async () => {
      await expect(result.current.retry(mockOperation)).rejects.toThrow();
    });

    expect(result.current.lastError).toBe(error);

    act(() => {
      result.current.reset();
    });

    expect(result.current.attemptCount).toBe(0);
    expect(result.current.isRetrying).toBe(false);
    expect(result.current.lastError).toBe(null);
  });
});

describe('useApolloRetry', () => {
  it('should retry on network errors', async () => {
    const { result } = renderHook(() => useApolloRetry({ maxAttempts: 2, delay: 10 }));
    
    const networkError = new ApolloError({
      networkError: { name: 'NetworkError', message: 'Network failure' } as any,
    });
    
    const mockOperation = vi
      .fn()
      .mockRejectedValueOnce(networkError)
      .mockResolvedValue('success');

    await act(async () => {
      const response = await result.current.retryQuery(mockOperation);
      expect(response).toBe('success');
    });

    expect(mockOperation).toHaveBeenCalledTimes(2);
  });

  it('should not retry on GraphQL errors', async () => {
    const { result } = renderHook(() => useApolloRetry({ maxAttempts: 3, delay: 10 }));
    
    const graphQLError = new ApolloError({
      graphQLErrors: [new GraphQLError('Validation error')],
    });
    
    const mockOperation = vi.fn().mockRejectedValue(graphQLError);

    await act(async () => {
      await expect(result.current.retryQuery(mockOperation)).rejects.toThrow();
    });

    expect(mockOperation).toHaveBeenCalledTimes(1);
  });

  it('should retry on 5xx server errors', async () => {
    const { result } = renderHook(() => useApolloRetry({ maxAttempts: 2, delay: 10 }));
    
    const serverError = new ApolloError({
      networkError: { 
        name: 'ServerError', 
        message: 'Server error',
        statusCode: 500 
      } as any,
    });
    
    const mockOperation = vi
      .fn()
      .mockRejectedValueOnce(serverError)
      .mockResolvedValue('success');

    await act(async () => {
      const response = await result.current.retryQuery(mockOperation);
      expect(response).toBe('success');
    });

    expect(mockOperation).toHaveBeenCalledTimes(2);
  });

  it('should not retry on 4xx client errors', async () => {
    const { result } = renderHook(() => useApolloRetry({ maxAttempts: 3, delay: 10 }));
    
    const clientError = new ApolloError({
      networkError: { 
        name: 'ClientError', 
        message: 'Bad request',
        statusCode: 400 
      } as any,
    });
    
    const mockOperation = vi.fn().mockRejectedValue(clientError);

    await act(async () => {
      await expect(result.current.retryQuery(mockOperation)).rejects.toThrow();
    });

    expect(mockOperation).toHaveBeenCalledTimes(1);
  });
});

describe('useNetworkStatus', () => {
  beforeEach(() => {
    mockAddEventListener.mockClear();
    mockRemoveEventListener.mockClear();
    navigator.onLine = true;
  });

  it('should initialize with current online status', () => {
    navigator.onLine = false;
    const { result } = renderHook(() => useNetworkStatus());
    
    expect(result.current.isOnline).toBe(false);
    expect(result.current.wasOffline).toBe(false);
  });

  it('should add event listeners on mount', () => {
    renderHook(() => useNetworkStatus());
    
    expect(mockAddEventListener).toHaveBeenCalledWith('online', expect.any(Function));
    expect(mockAddEventListener).toHaveBeenCalledWith('offline', expect.any(Function));
  });

  it('should remove event listeners on unmount', () => {
    const { unmount } = renderHook(() => useNetworkStatus());
    
    unmount();
    
    expect(mockRemoveEventListener).toHaveBeenCalledWith('online', expect.any(Function));
    expect(mockRemoveEventListener).toHaveBeenCalledWith('offline', expect.any(Function));
  });

  it('should update state when going offline', () => {
    const { result } = renderHook(() => useNetworkStatus());
    
    expect(result.current.isOnline).toBe(true);
    expect(result.current.wasOffline).toBe(false);
    
    // Simulate going offline
    act(() => {
      const offlineHandler = mockAddEventListener.mock.calls.find(
        call => call[0] === 'offline'
      )?.[1];
      offlineHandler?.();
    });
    
    expect(result.current.isOnline).toBe(false);
    expect(result.current.wasOffline).toBe(true);
  });

  it('should update state when coming back online', () => {
    const { result } = renderHook(() => useNetworkStatus());
    
    // First go offline
    act(() => {
      const offlineHandler = mockAddEventListener.mock.calls.find(
        call => call[0] === 'offline'
      )?.[1];
      offlineHandler?.();
    });
    
    expect(result.current.isOnline).toBe(false);
    expect(result.current.wasOffline).toBe(true);
    
    // Then come back online
    act(() => {
      const onlineHandler = mockAddEventListener.mock.calls.find(
        call => call[0] === 'online'
      )?.[1];
      onlineHandler?.();
    });
    
    expect(result.current.isOnline).toBe(true);
    expect(result.current.wasOffline).toBe(true);
  });
});