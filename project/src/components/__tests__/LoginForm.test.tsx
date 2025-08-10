import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MockedProvider } from '@apollo/client/testing';
import { LoginForm } from '../LoginForm';
import { LoginDocument } from '../../generated/graphql';

const mockSuccessResponse = {
  request: {
    query: LoginDocument,
    variables: {
      email: 'test@example.com',
      password: 'password123',
    },
  },
  result: {
    data: {
      login: {
        token: 'mock-jwt-token',
        user: {
          id: '1',
          email: 'test@example.com',
          name: 'Test User',
          avatar: null,
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z',
        },
        expiresAt: '2024-12-31T23:59:59Z',
      },
    },
  },
};

const mockErrorResponse = {
  request: {
    query: LoginDocument,
    variables: {
      email: 'invalid@example.com',
      password: 'wrongpassword',
    },
  },
  error: new Error('Invalid credentials'),
};

describe('LoginForm', () => {
  const mockOnSuccess = vi.fn();
  const mockOnError = vi.fn();

  beforeEach(() => {
    mockOnSuccess.mockClear();
    mockOnError.mockClear();
  });

  it('renders login form with all fields', () => {
    render(
      <MockedProvider mocks={[]} addTypename={false}>
        <LoginForm onSuccess={mockOnSuccess} onError={mockOnError} />
      </MockedProvider>
    );

    expect(screen.getByLabelText(/email address/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /sign in/i })).toBeInTheDocument();
  });

  it('shows validation error for empty fields', async () => {
    const user = userEvent.setup();
    
    render(
      <MockedProvider mocks={[]} addTypename={false}>
        <LoginForm onSuccess={mockOnSuccess} onError={mockOnError} />
      </MockedProvider>
    );

    const submitButton = screen.getByRole('button', { name: /sign in/i });
    expect(submitButton).toBeDisabled();

    // Try to submit with empty fields
    await user.click(submitButton);
    
    // Button should remain disabled
    expect(submitButton).toBeDisabled();
  });

  it('enables submit button when fields are filled', async () => {
    const user = userEvent.setup();
    
    render(
      <MockedProvider mocks={[]} addTypename={false}>
        <LoginForm onSuccess={mockOnSuccess} onError={mockOnError} />
      </MockedProvider>
    );

    const emailInput = screen.getByLabelText(/email address/i);
    const passwordInput = screen.getByLabelText(/password/i);
    const submitButton = screen.getByRole('button', { name: /sign in/i });

    await user.type(emailInput, 'test@example.com');
    await user.type(passwordInput, 'password123');

    expect(submitButton).toBeEnabled();
  });

  it('calls onSuccess when login is successful', async () => {
    const user = userEvent.setup();
    
    render(
      <MockedProvider mocks={[mockSuccessResponse]} addTypename={false}>
        <LoginForm onSuccess={mockOnSuccess} onError={mockOnError} />
      </MockedProvider>
    );

    const emailInput = screen.getByLabelText(/email address/i);
    const passwordInput = screen.getByLabelText(/password/i);
    const submitButton = screen.getByRole('button', { name: /sign in/i });

    await user.type(emailInput, 'test@example.com');
    await user.type(passwordInput, 'password123');
    await user.click(submitButton);

    await waitFor(() => {
      expect(mockOnSuccess).toHaveBeenCalledWith(
        'mock-jwt-token',
        expect.objectContaining({
          id: '1',
          email: 'test@example.com',
          name: 'Test User',
        })
      );
    });
  });

  it('shows loading state during submission', async () => {
    const user = userEvent.setup();
    
    render(
      <MockedProvider mocks={[mockSuccessResponse]} addTypename={false}>
        <LoginForm onSuccess={mockOnSuccess} onError={mockOnError} />
      </MockedProvider>
    );

    const emailInput = screen.getByLabelText(/email address/i);
    const passwordInput = screen.getByLabelText(/password/i);
    const submitButton = screen.getByRole('button', { name: /sign in/i });

    await user.type(emailInput, 'test@example.com');
    await user.type(passwordInput, 'password123');
    await user.click(submitButton);

    // Should show loading state
    expect(screen.getByText(/signing in/i)).toBeInTheDocument();
    expect(submitButton).toBeDisabled();
  });

  it('calls onError when login fails', async () => {
    const user = userEvent.setup();
    
    render(
      <MockedProvider mocks={[mockErrorResponse]} addTypename={false}>
        <LoginForm onSuccess={mockOnSuccess} onError={mockOnError} />
      </MockedProvider>
    );

    const emailInput = screen.getByLabelText(/email address/i);
    const passwordInput = screen.getByLabelText(/password/i);
    const submitButton = screen.getByRole('button', { name: /sign in/i });

    await user.type(emailInput, 'invalid@example.com');
    await user.type(passwordInput, 'wrongpassword');
    await user.click(submitButton);

    await waitFor(() => {
      expect(mockOnError).toHaveBeenCalledWith('Invalid credentials');
    });
  });

  it('displays error message in UI when login fails', async () => {
    const user = userEvent.setup();
    
    render(
      <MockedProvider mocks={[mockErrorResponse]} addTypename={false}>
        <LoginForm />
      </MockedProvider>
    );

    const emailInput = screen.getByLabelText(/email address/i);
    const passwordInput = screen.getByLabelText(/password/i);
    const submitButton = screen.getByRole('button', { name: /sign in/i });

    await user.type(emailInput, 'invalid@example.com');
    await user.type(passwordInput, 'wrongpassword');
    await user.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText(/invalid credentials/i)).toBeInTheDocument();
    });
  });

  it('shows demo credentials hint', () => {
    render(
      <MockedProvider mocks={[]} addTypename={false}>
        <LoginForm onSuccess={mockOnSuccess} onError={mockOnError} />
      </MockedProvider>
    );

    expect(screen.getByText(/demo credentials/i)).toBeInTheDocument();
    expect(screen.getByText(/admin@example.com/)).toBeInTheDocument();
  });
});