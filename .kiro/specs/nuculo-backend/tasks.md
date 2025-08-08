# Implementation Plan

## Current Status
The current backend directory contains a Go-based GraphQL implementation that doesn't match the Nuculo backend requirements. The Nuculo backend should be implemented in Node.js/TypeScript according to the design document. The existing Go code needs to be replaced with a proper Node.js/TypeScript implementation.

- [ ] 1. Replace Go implementation with Node.js/TypeScript foundation
  - Remove existing Go-based GraphQL server code
  - Create Node.js/TypeScript project structure for Nuculo backend
  - Set up package.json with required dependencies (Apollo Server, TypeScript, etc.)
  - Initialize TypeScript configuration with strict mode
  - Create directory structure for resolvers, services, repositories, and GraphQL schema
  - _Requirements: 1.1, 2.1, 3.1, 4.1, 5.1, 6.1_

- [ ] 2. Define Nuculo-specific GraphQL schema
  - Create GraphQL schema for Contact, Service, AdminUser, and AnalyticsEvent types
  - Define input types for contact submission, service management, and authentication
  - Set up queries for services (public), contacts, analytics (admin-only)
  - Define mutations for contact submission, service CRUD, authentication
  - Create TypeScript interfaces matching the GraphQL schema
  - _Requirements: 1.1, 2.1, 3.1, 4.1, 5.1, 6.1_

- [ ] 3. Set up Apollo Server with TypeScript configuration
  - Initialize Apollo Server with GraphQL schema
  - Configure Express.js integration for HTTP handling
  - Set up CORS for frontend integration
  - Add GraphQL Playground for development
  - Create health check endpoint
  - _Requirements: 1.1, 2.1, 3.1, 4.1, 5.1, 6.1_

- [ ] 4. Implement database layer and connection management
  - Set up PostgreSQL connection with connection pooling using pg library
  - Create database migration scripts for contacts, services, admin_users, analytics_events tables
  - Implement database connection utilities with error handling
  - Create repository pattern interfaces and base repository class
  - Add database connectivity verification on startup
  - _Requirements: 1.2, 2.3, 3.1, 5.1, 6.5_

- [ ] 5. Create GraphQL validation and custom scalars
  - Implement custom GraphQL scalars for email and datetime validation
  - Create input validation directives for GraphQL arguments
  - Set up GraphQL schema validation with custom validators
  - Add input sanitization to prevent XSS and injection attacks
  - Write unit tests for validation directives and scalars
  - _Requirements: 1.1, 1.5, 1.6, 3.6, 6.3_

- [ ] 6. Implement GraphQL security and authentication
  - Set up JWT authentication in GraphQL context using jsonwebtoken
  - Create authentication directive for protected resolvers
  - Implement rate limiting directive with Redis storage
  - Add GraphQL security plugins (query complexity, depth limiting)
  - Write unit tests for authentication and security directives
  - _Requirements: 2.1, 2.6, 6.1, 6.2, 6.6_

- [ ] 7. Build contact form submission GraphQL mutation
  - Create ContactService with contact submission business logic
  - Implement ContactRepository with database operations
  - Create submitContact GraphQL mutation with validation and rate limiting
  - Add request context tracking (IP, user agent, referrer)
  - Write unit and integration tests for contact submission resolver
  - _Requirements: 1.1, 1.2, 1.5, 1.6_

- [ ] 8. Implement email notification system
  - Set up EmailService with Nodemailer configuration
  - Create email templates for admin notifications and user auto-replies
  - Integrate email sending into contact submission flow
  - Add email service error handling and retry logic
  - Write tests for email service with mocked SMTP
  - _Requirements: 1.3, 1.4_

- [ ] 9. Create admin contact management GraphQL operations
  - Implement contacts query with pagination and filtering
  - Create updateContactStatus mutation for status updates
  - Add search functionality by name, email, and date range in resolver
  - Implement proper authorization checks for admin-only operations
  - Write integration tests for admin contact resolvers
  - _Requirements: 2.2, 2.3, 2.4, 2.5_

- [ ] 10. Build services management GraphQL system
  - Create ServicesService with CRUD operations
  - Implement ServicesRepository with database operations
  - Create public services query with DataLoader caching (5-minute TTL)
  - Add cache invalidation logic for service updates
  - Write unit tests for services business logic
  - _Requirements: 3.1, 3.4, 4.2, 4.3_

- [ ] 11. Implement admin services management GraphQL mutations
  - Create createService mutation for service creation
  - Implement updateService mutation for service updates
  - Add deleteService mutation for service removal
  - Create reorderServices mutation for display order management
  - Write integration tests for all admin services resolvers
  - _Requirements: 3.1, 3.2, 3.3, 3.5, 3.6_

- [ ] 12. Create analytics tracking and reporting system
  - Implement AnalyticsService with event recording and aggregation
  - Create AnalyticsRepository for database operations
  - Add analytics tracking to contact form submissions
  - Implement data aggregation for metrics calculation
  - Write unit tests for analytics service logic
  - _Requirements: 5.1, 5.4_

- [ ] 13. Build admin analytics GraphQL query
  - Create analytics query with date range filtering
  - Implement comprehensive metrics resolver for dashboard
  - Add aggregation logic for daily/weekly/monthly trends
  - Include top referrers and submission statistics
  - Write integration tests for analytics resolver
  - _Requirements: 5.2, 5.3, 5.5_

- [ ] 14. Implement authentication GraphQL mutations
  - Create login mutation with credential validation using bcrypt
  - Implement refreshToken mutation for token renewal
  - Add logout mutation for session invalidation
  - Create me query for current user information
  - Write integration tests for authentication resolvers
  - _Requirements: 2.1, 2.6_

- [ ] 15. Add comprehensive GraphQL error handling and logging
  - Implement GraphQL error formatting with consistent error structure
  - Set up Winston logger with structured logging for GraphQL operations
  - Add error logging without exposing sensitive information in GraphQL responses
  - Create custom GraphQL error classes for different error types
  - Write tests for GraphQL error handling scenarios
  - _Requirements: 4.4, 6.4_

- [ ] 16. Optimize GraphQL performance and add caching
  - Implement DataLoader for efficient database batching and caching
  - Add Redis caching for GraphQL query results with 5-minute TTL
  - Set up cache invalidation logic for mutations
  - Optimize database queries with proper indexing and query analysis
  - Write performance tests to verify 200ms response time requirement
  - _Requirements: 4.1, 4.3, 4.5_

- [ ] 17. Create comprehensive GraphQL test suite
  - Write end-to-end GraphQL tests for complete contact submission flow
  - Create integration tests for admin workflow (auth + CRUD + analytics) using GraphQL operations
  - Add security tests for GraphQL input validation and injection prevention
  - Implement load testing for concurrent GraphQL query/mutation handling
  - Set up test data fixtures and database seeding for GraphQL testing
  - _Requirements: All requirements for comprehensive coverage_

- [ ] 18. Set up React frontend with Apollo Client integration
  - Create React application structure with Apollo Client setup
  - Implement GraphQL client configuration with authentication
  - Create React components for contact form with GraphQL mutations
  - Build admin dashboard components using GraphQL queries
  - Set up GraphQL code generation for TypeScript types
  - _Requirements: 1.1, 2.2, 3.4, 4.2, 5.2_

- [ ] 19. Set up deployment configuration and monitoring
  - Create Docker configuration for containerized deployment
  - Set up environment variable validation on startup
  - Implement GraphQL playground and health check endpoints
  - Add database connectivity verification on startup
  - Create deployment scripts and documentation
  - _Requirements: 6.5_