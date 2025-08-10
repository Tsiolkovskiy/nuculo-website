import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MockedProvider } from '@apollo/client/testing';
import { PostsList } from '../PostsList';
import { GetPostsDocument } from '../../generated/graphql';

const mockPostsResponse = {
  request: {
    query: GetPostsDocument,
    variables: {
      filters: { published: true },
      pagination: { page: 1, limit: 10 },
    },
  },
  result: {
    data: {
      posts: {
        edges: [
          {
            node: {
              id: '1',
              title: 'Test Post 1',
              content: 'This is the content of test post 1',
              tags: ['react', 'typescript'],
              published: true,
              createdAt: '2024-01-01T00:00:00Z',
              updatedAt: '2024-01-01T00:00:00Z',
              author: {
                id: '1',
                name: 'John Doe',
                avatar: null,
              },
            },
            cursor: 'cursor1',
          },
          {
            node: {
              id: '2',
              title: 'Test Post 2',
              content: 'This is the content of test post 2',
              tags: ['graphql', 'apollo'],
              published: true,
              createdAt: '2024-01-02T00:00:00Z',
              updatedAt: '2024-01-02T00:00:00Z',
              author: {
                id: '2',
                name: 'Jane Smith',
                avatar: 'https://example.com/avatar.jpg',
              },
            },
            cursor: 'cursor2',
          },
        ],
        pageInfo: {
          hasNextPage: true,
          hasPreviousPage: false,
          startCursor: 'cursor1',
          endCursor: 'cursor2',
        },
        totalCount: 5,
      },
    },
  },
};

const mockEmptyResponse = {
  request: {
    query: GetPostsDocument,
    variables: {
      filters: { published: true },
      pagination: { page: 1, limit: 10 },
    },
  },
  result: {
    data: {
      posts: {
        edges: [],
        pageInfo: {
          hasNextPage: false,
          hasPreviousPage: false,
          startCursor: null,
          endCursor: null,
        },
        totalCount: 0,
      },
    },
  },
};

const mockErrorResponse = {
  request: {
    query: GetPostsDocument,
    variables: {
      filters: { published: true },
      pagination: { page: 1, limit: 10 },
    },
  },
  error: new Error('Failed to fetch posts'),
};

describe('PostsList', () => {
  const mockOnPostClick = vi.fn();

  beforeEach(() => {
    mockOnPostClick.mockClear();
  });

  it('renders loading state initially', () => {
    render(
      <MockedProvider mocks={[]} addTypename={false}>
        <PostsList />
      </MockedProvider>
    );

    expect(screen.getByRole('status', { name: /loading/i })).toBeInTheDocument();
  });

  it('renders posts when data is loaded', async () => {
    render(
      <MockedProvider mocks={[mockPostsResponse]} addTypename={false}>
        <PostsList onPostClick={mockOnPostClick} />
      </MockedProvider>
    );

    await waitFor(() => {
      expect(screen.getByText('Test Post 1')).toBeInTheDocument();
      expect(screen.getByText('Test Post 2')).toBeInTheDocument();
    });

    // Check post content
    expect(screen.getByText('This is the content of test post 1')).toBeInTheDocument();
    expect(screen.getByText('This is the content of test post 2')).toBeInTheDocument();

    // Check authors
    expect(screen.getByText('John Doe')).toBeInTheDocument();
    expect(screen.getByText('Jane Smith')).toBeInTheDocument();

    // Check tags
    expect(screen.getByText('react')).toBeInTheDocument();
    expect(screen.getByText('typescript')).toBeInTheDocument();
    expect(screen.getByText('graphql')).toBeInTheDocument();
    expect(screen.getByText('apollo')).toBeInTheDocument();

    // Check published status
    expect(screen.getAllByText('Published')).toHaveLength(2);
  });

  it('renders empty state when no posts', async () => {
    render(
      <MockedProvider mocks={[mockEmptyResponse]} addTypename={false}>
        <PostsList />
      </MockedProvider>
    );

    await waitFor(() => {
      expect(screen.getByText('No posts found')).toBeInTheDocument();
      expect(screen.getByText('No published posts available.')).toBeInTheDocument();
    });
  });

  it('renders error state when query fails', async () => {
    render(
      <MockedProvider mocks={[mockErrorResponse]} addTypename={false}>
        <PostsList />
      </MockedProvider>
    );

    await waitFor(() => {
      expect(screen.getByText('Error loading posts')).toBeInTheDocument();
      expect(screen.getByText('Failed to fetch posts')).toBeInTheDocument();
    });

    // Check retry button
    expect(screen.getByText('Try again')).toBeInTheDocument();
  });

  it('calls onPostClick when read more button is clicked', async () => {
    const user = userEvent.setup();
    
    render(
      <MockedProvider mocks={[mockPostsResponse]} addTypename={false}>
        <PostsList onPostClick={mockOnPostClick} />
      </MockedProvider>
    );

    await waitFor(() => {
      expect(screen.getByText('Test Post 1')).toBeInTheDocument();
    });

    const readMoreButtons = screen.getAllByText('Read more →');
    await user.click(readMoreButtons[0]);

    expect(mockOnPostClick).toHaveBeenCalledWith('1');
  });

  it('shows load more button when there are more pages', async () => {
    render(
      <MockedProvider mocks={[mockPostsResponse]} addTypename={false}>
        <PostsList />
      </MockedProvider>
    );

    await waitFor(() => {
      expect(screen.getByText('Load more posts')).toBeInTheDocument();
    });
  });

  it('displays author avatars correctly', async () => {
    render(
      <MockedProvider mocks={[mockPostsResponse]} addTypename={false}>
        <PostsList />
      </MockedProvider>
    );

    await waitFor(() => {
      // Jane Smith has an avatar
      const avatarImg = screen.getByAltText('Jane Smith');
      expect(avatarImg).toBeInTheDocument();
      expect(avatarImg).toHaveAttribute('src', 'https://example.com/avatar.jpg');

      // John Doe doesn't have an avatar, should show initials
      expect(screen.getByText('J')).toBeInTheDocument(); // Initial for John
    });
  });

  it('formats dates correctly', async () => {
    render(
      <MockedProvider mocks={[mockPostsResponse]} addTypename={false}>
        <PostsList />
      </MockedProvider>
    );

    await waitFor(() => {
      // Check that dates are formatted (exact format may vary by locale)
      expect(screen.getByText(/1\/1\/2024|Jan.*2024|2024/)).toBeInTheDocument();
    });
  });

  it('handles different published states', async () => {
    const draftPostResponse = {
      ...mockPostsResponse,
      result: {
        data: {
          posts: {
            ...mockPostsResponse.result.data.posts,
            edges: [
              {
                ...mockPostsResponse.result.data.posts.edges[0],
                node: {
                  ...mockPostsResponse.result.data.posts.edges[0].node,
                  published: false,
                },
              },
            ],
          },
        },
      },
    };

    render(
      <MockedProvider mocks={[draftPostResponse]} addTypename={false}>
        <PostsList showPublishedOnly={false} />
      </MockedProvider>
    );

    await waitFor(() => {
      expect(screen.getByText('Draft')).toBeInTheDocument();
    });
  });
});