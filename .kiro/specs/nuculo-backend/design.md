# Design Document

## Overview

The Nuculo Backend API is a React-based GraphQL API system designed to support the Nuculo consulting website. The system provides secure GraphQL queries and mutations for contact form submissions, content management, basic analytics, and admin functionality. The architecture prioritizes security, performance, and maintainability while leveraging GraphQL's type safety and efficient data fetching.

### Key Design Principles
- **Security First**: All GraphQL operations implement proper validation, sanitization, and authentication
- **Performance Optimized**: GraphQL's efficient data fetching with caching strategies ensure fast response times
- **Type Safety**: Strong typing throughout with GraphQL schema and TypeScript
- **Scalable Architecture**: Modular resolver design allows for future feature expansion
- **Data Protection**: Encryption and privacy-conscious data handling throughout

## Architecture

### System Architecture
The backend follows a GraphQL-first architecture pattern:

```
┌─────────────────┐
│   Frontend      │
│   (React)       │
└─────────┬───────┘
          │ GraphQL
┌─────────▼───────┐
│ GraphQL Layer   │
│ (Apollo Server) │
├─────────────────┤
│   Resolvers     │
│ (Query/Mutation)│
├─────────────────┤
│ Business Logic  │
│   (Services)    │
├─────────────────┤
│  Data Access    │
│ (Repositories)  │
├─────────────────┤
│   Database      │
│  (PostgreSQL)   │
└─────────────────┘
```

### Technology Stack
- **Runtime**: Node.js with TypeScript
- **GraphQL Server**: Apollo Server
- **Frontend Framework**: React with Apollo Client
- **Database**: PostgreSQL with connection pooling
- **Authentication**: JWT tokens with bcrypt password hashing
- **Email Service**: Nodemailer with SMTP configuration
- **Caching**: Redis for GraphQL response caching and DataLoader
- **Validation**: GraphQL schema validation with custom validators
- **Security**: GraphQL security middleware, query complexity analysis

### Deployment Architecture
- **Environment**: Docker containerized deployment
- **Database**: Managed PostgreSQL instance
- **Caching**: Redis instance for session and API caching
- **Email**: SMTP service integration (SendGrid/AWS SES)
- **Monitoring**: Structured logging with Winston

## Components and Interfaces

### GraphQL Schema

#### Type Definitions
```graphql
type Contact {
  id: ID!
  name: String!
  email: String!
  message: String!
  status: ContactStatus!
  createdAt: DateTime!
  updatedAt: DateTime!
  userAgent: String
  referrer: String
  ipAddress: String
}

type Service {
  id: ID!
  title: String!
  description: String!
  iconIdentifier: String
  displayOrder: Int!
  isActive: Boolean!
  createdAt: DateTime!
  updatedAt: DateTime!
}

type AdminUser {
  id: ID!
  email: String!
  name: String!
  isActive: Boolean!
  lastLogin: DateTime
  createdAt: DateTime!
  updatedAt: DateTime!
}

type AnalyticsMetrics {
  totalContacts: Int!
  contactsByPeriod: [ContactMetric!]!
  averageResponseTime: Float
  topReferrers: [ReferrerMetric!]!
}

type ContactMetric {
  date: String!
  count: Int!
}

type ReferrerMetric {
  referrer: String!
  count: Int!
}

enum ContactStatus {
  NEW
  CONTACTED
  RESOLVED
}

input ContactInput {
  name: String!
  email: String!
  message: String!
}

input ServiceInput {
  title: String!
  description: String!
  iconIdentifier: String
}

input ServiceUpdateInput {
  title: String
  description: String
  iconIdentifier: String
  isActive: Boolean
}

input ContactFilters {
  status: ContactStatus
  dateFrom: DateTime
  dateTo: DateTime
  searchTerm: String
}

input PaginationInput {
  page: Int = 1
  limit: Int = 20
}

type PaginatedContacts {
  contacts: [Contact!]!
  totalCount: Int!
  hasNextPage: Boolean!
  hasPreviousPage: Boolean!
}
```

#### Query Operations
```graphql
type Query {
  # Public queries
  services: [Service!]!
  
  # Admin queries (require authentication)
  contacts(filters: ContactFilters, pagination: PaginationInput): PaginatedContacts!
  contact(id: ID!): Contact
  adminServices: [Service!]!
  service(id: ID!): Service
  analytics(dateFrom: DateTime, dateTo: DateTime): AnalyticsMetrics!
  me: AdminUser
}
```

#### Mutation Operations
```graphql
type Mutation {
  # Public mutations
  submitContact(input: ContactInput!): Contact!
  
  # Admin mutations (require authentication)
  updateContactStatus(id: ID!, status: ContactStatus!): Contact!
  createService(input: ServiceInput!): Service!
  updateService(id: ID!, input: ServiceUpdateInput!): Service!
  deleteService(id: ID!): Boolean!
  reorderServices(serviceIds: [ID!]!): [Service!]!
  
  # Authentication mutations
  login(email: String!, password: String!): AuthPayload!
  refreshToken: AuthPayload!
  logout: Boolean!
}

type AuthPayload {
  token: String!
  user: AdminUser!
  expiresAt: DateTime!
}
```

### GraphQL Resolvers

#### Query Resolvers
```typescript
const Query = {
  services: async () => servicesService.getActiveServices(),
  contacts: async (_, { filters, pagination }, { user }) => {
    requireAuth(user);
    return contactService.getContacts(filters, pagination);
  },
  contact: async (_, { id }, { user }) => {
    requireAuth(user);
    return contactService.getContactById(id);
  },
  adminServices: async (_, __, { user }) => {
    requireAuth(user);
    return servicesService.getAllServices();
  },
  service: async (_, { id }, { user }) => {
    requireAuth(user);
    return servicesService.getServiceById(id);
  },
  analytics: async (_, { dateFrom, dateTo }, { user }) => {
    requireAuth(user);
    return analyticsService.getMetrics(dateFrom, dateTo);
  },
  me: async (_, __, { user }) => {
    requireAuth(user);
    return user;
  }
};
```

#### Mutation Resolvers
```typescript
const Mutation = {
  submitContact: async (_, { input }, { req }) => {
    await rateLimiter.checkLimit(req.ip, 'contact_submission');
    const contact = await contactService.submitContact(input, req);
    await emailService.sendNotifications(contact);
    return contact;
  },
  updateContactStatus: async (_, { id, status }, { user }) => {
    requireAuth(user);
    return contactService.updateContactStatus(id, status);
  },
  createService: async (_, { input }, { user }) => {
    requireAuth(user);
    return servicesService.createService(input);
  },
  updateService: async (_, { id, input }, { user }) => {
    requireAuth(user);
    return servicesService.updateService(id, input);
  },
  deleteService: async (_, { id }, { user }) => {
    requireAuth(user);
    return servicesService.deleteService(id);
  },
  reorderServices: async (_, { serviceIds }, { user }) => {
    requireAuth(user);
    return servicesService.reorderServices(serviceIds);
  },
  login: async (_, { email, password }) => {
    return authService.login(email, password);
  },
  refreshToken: async (_, __, { refreshToken }) => {
    return authService.refreshToken(refreshToken);
  },
  logout: async (_, __, { user }) => {
    requireAuth(user);
    return authService.logout(user.id);
  }
};
```

### Service Layer Components

#### ContactService
- **Purpose**: Handle contact form business logic
- **Methods**:
  - `submitContact(data, req)`: Validate and store contact submission with request context
  - `getContacts(filters, pagination)`: Retrieve filtered and paginated contact list
  - `getContactById(id)`: Retrieve single contact by ID
  - `updateContactStatus(id, status)`: Update submission status

#### ServicesService
- **Purpose**: Manage services content
- **Methods**:
  - `getActiveServices()`: Retrieve public services list with caching
  - `getAllServices()`: Retrieve all services for admin
  - `getServiceById(id)`: Retrieve single service by ID
  - `createService(data)`: Create new service entry
  - `updateService(id, data)`: Modify existing service
  - `deleteService(id)`: Remove service
  - `reorderServices(serviceIds)`: Update display order

#### AnalyticsService
- **Purpose**: Generate analytics and metrics
- **Methods**:
  - `getMetrics(dateFrom, dateTo)`: Get comprehensive analytics
  - `recordContactEvent(contact, req)`: Track contact submission events
  - `getContactMetrics(dateRange)`: Contact submission analytics
  - `getReferrerMetrics(dateRange)`: Top referrer statistics

#### EmailService
- **Purpose**: Handle email communications
- **Methods**:
  - `sendNotifications(contact)`: Send both admin and user notifications
  - `sendContactNotification(contact)`: Notify admin of new submission
  - `sendAutoReply(contact)`: Send confirmation to user

#### AuthService
- **Purpose**: Handle authentication and authorization
- **Methods**:
  - `login(email, password)`: Authenticate user and return JWT
  - `refreshToken(token)`: Refresh expired JWT token
  - `logout(userId)`: Invalidate user session
  - `verifyToken(token)`: Validate JWT token

### GraphQL Middleware Components

#### AuthenticationDirective
- **Purpose**: Verify JWT tokens in GraphQL context
- **Implementation**: Custom GraphQL directive for protected resolvers
- **Error Handling**: Throw GraphQL authentication errors

#### RateLimitingDirective
- **Purpose**: Prevent GraphQL operation abuse
- **Implementation**: Custom directive with Redis-based rate limiting
- **Configuration**: Different limits for queries vs mutations

#### ValidationDirective
- **Purpose**: Validate GraphQL input arguments
- **Implementation**: Schema-based validation using custom scalars
- **Error Handling**: Throw GraphQL validation errors

#### SecurityPlugin
- **Purpose**: Implement GraphQL security measures
- **Implementation**: Query complexity analysis, depth limiting
- **Features**: Introspection disabling in production, query whitelisting

## Data Models

### Database Schema

#### contacts table
```sql
CREATE TABLE contacts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    status VARCHAR(50) DEFAULT 'new',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_agent TEXT,
    referrer TEXT,
    ip_address INET
);

CREATE INDEX idx_contacts_created_at ON contacts(created_at);
CREATE INDEX idx_contacts_status ON contacts(status);
CREATE INDEX idx_contacts_email ON contacts(email);
```

#### services table
```sql
CREATE TABLE services (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    icon_identifier VARCHAR(100),
    display_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_services_display_order ON services(display_order);
CREATE INDEX idx_services_active ON services(is_active);
```

#### admin_users table
```sql
CREATE TABLE admin_users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_admin_users_email ON admin_users(email);
```

#### analytics_events table
```sql
CREATE TABLE analytics_events (
    id SERIAL PRIMARY KEY,
    event_type VARCHAR(100) NOT NULL,
    event_data JSONB,
    user_agent TEXT,
    ip_address INET,
    referrer TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_analytics_events_type_date ON analytics_events(event_type, created_at);
CREATE INDEX idx_analytics_events_created_at ON analytics_events(created_at);
```

### TypeScript Interfaces

```typescript
interface Contact {
    id: number;
    name: string;
    email: string;
    message: string;
    status: 'new' | 'contacted' | 'resolved';
    createdAt: Date;
    updatedAt: Date;
    userAgent?: string;
    referrer?: string;
    ipAddress?: string;
}

interface Service {
    id: number;
    title: string;
    description: string;
    iconIdentifier?: string;
    displayOrder: number;
    isActive: boolean;
    createdAt: Date;
    updatedAt: Date;
}

interface AdminUser {
    id: number;
    email: string;
    name: string;
    isActive: boolean;
    lastLogin?: Date;
    createdAt: Date;
    updatedAt: Date;
}

interface AnalyticsEvent {
    id: number;
    eventType: string;
    eventData?: Record<string, any>;
    userAgent?: string;
    ipAddress?: string;
    referrer?: string;
    createdAt: Date;
}

interface ContactSubmissionRequest {
    name: string;
    email: string;
    message: string;
}

interface ServiceCreateRequest {
    title: string;
    description: string;
    iconIdentifier?: string;
}

interface ServiceReorderRequest {
    serviceIds: number[];
}
```

## Error Handling

### Error Response Format
All API errors follow a consistent format:
```typescript
interface ErrorResponse {
    error: {
        code: string;
        message: string;
        details?: any;
        timestamp: string;
        requestId: string;
    };
}
```

### Error Categories

#### Validation Errors (400)
- **Trigger**: Invalid request data, missing required fields
- **Response**: Detailed field-level validation messages
- **Logging**: Info level with request details

#### Authentication Errors (401)
- **Trigger**: Missing, invalid, or expired JWT tokens
- **Response**: Generic authentication failure message
- **Logging**: Warning level with IP and user agent

#### Authorization Errors (403)
- **Trigger**: Valid token but insufficient permissions
- **Response**: Access denied message
- **Logging**: Warning level with user context

#### Not Found Errors (404)
- **Trigger**: Requested resource doesn't exist
- **Response**: Resource not found message
- **Logging**: Info level

#### Rate Limit Errors (429)
- **Trigger**: Exceeded rate limit thresholds
- **Response**: Rate limit exceeded with retry information
- **Logging**: Warning level with IP tracking

#### Server Errors (500)
- **Trigger**: Unexpected application errors
- **Response**: Generic server error message (no sensitive data)
- **Logging**: Error level with full stack trace

### Error Handling Strategy
1. **Input Validation**: Joi schemas validate all incoming data
2. **Database Errors**: Wrapped in try-catch with appropriate error mapping
3. **External Service Errors**: Email service failures handled gracefully
4. **Logging**: Winston logger with structured error information
5. **Monitoring**: Error rate tracking and alerting capabilities

## Testing Strategy

### Unit Testing
- **Framework**: Jest with TypeScript support
- **Coverage Target**: 90% code coverage minimum
- **Focus Areas**:
  - Service layer business logic
  - Validation schemas and middleware
  - Data transformation functions
  - Error handling scenarios

### Integration Testing
- **Database Testing**: Test database with Docker containers
- **API Testing**: Supertest for endpoint testing
- **Email Testing**: Mock email service for notification testing
- **Authentication Testing**: JWT token validation and expiration

### End-to-End Testing
- **Contact Flow**: Complete contact submission and notification process
- **Admin Workflow**: Authentication, CRUD operations, and analytics
- **Performance Testing**: Load testing for response time requirements
- **Security Testing**: Input validation, SQL injection, XSS prevention

### Test Data Management
- **Fixtures**: Predefined test data for consistent testing
- **Database Seeding**: Automated test data setup and teardown
- **Mock Services**: External service mocking for isolated testing
- **Environment Isolation**: Separate test database and configurations

### Continuous Integration
- **Automated Testing**: All tests run on every commit
- **Code Quality**: ESLint, Prettier, and TypeScript strict mode
- **Security Scanning**: Dependency vulnerability checks
- **Performance Monitoring**: Response time regression testing