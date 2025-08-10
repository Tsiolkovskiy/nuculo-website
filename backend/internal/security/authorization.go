package security

import (
	"context"
	"fmt"
	"strings"

	"github.com/99designs/gqlgen/graphql"
)

// Role represents user roles in the system
type Role string

const (
	RoleAdmin     Role = "admin"
	RoleModerator Role = "moderator"
	RoleUser      Role = "user"
	RoleGuest     Role = "guest"
)

// Permission represents specific permissions
type Permission string

const (
	PermissionReadPost      Permission = "read:post"
	PermissionWritePost     Permission = "write:post"
	PermissionDeletePost    Permission = "delete:post"
	PermissionReadUser      Permission = "read:user"
	PermissionWriteUser     Permission = "write:user"
	PermissionDeleteUser    Permission = "delete:user"
	PermissionReadComment   Permission = "read:comment"
	PermissionWriteComment  Permission = "write:comment"
	PermissionDeleteComment Permission = "delete:comment"
	PermissionModerate      Permission = "moderate"
	PermissionAdmin         Permission = "admin"
)

// User represents the authenticated user
type User struct {
	ID          string   `json:"id"`
	Email       string   `json:"email"`
	Username    string   `json:"username"`
	Role        Role     `json:"role"`
	Permissions []string `json:"permissions"`
	IsActive    bool     `json:"is_active"`
	IsVerified  bool     `json:"is_verified"`
}

// HasRole checks if user has a specific role
func (u *User) HasRole(role Role) bool {
	return u.Role == role || u.Role == RoleAdmin // Admin has all roles
}

// HasPermission checks if user has a specific permission
func (u *User) HasPermission(permission Permission) bool {
	// Admin has all permissions
	if u.Role == RoleAdmin {
		return true
	}
	
	for _, perm := range u.Permissions {
		if perm == string(permission) {
			return true
		}
	}
	
	return false
}

// CanAccessResource checks if user can access a specific resource
func (u *User) CanAccessResource(resourceType, resourceID, action string) bool {
	// Check if user is active and verified
	if !u.IsActive || !u.IsVerified {
		return false
	}
	
	// Admin can access everything
	if u.Role == RoleAdmin {
		return true
	}
	
	// Check specific permissions based on resource type and action
	switch resourceType {
	case "post":
		return u.canAccessPost(resourceID, action)
	case "user":
		return u.canAccessUser(resourceID, action)
	case "comment":
		return u.canAccessComment(resourceID, action)
	default:
		return false
	}
}

// canAccessPost checks post-specific access
func (u *User) canAccessPost(postID, action string) bool {
	switch action {
	case "read":
		return u.HasPermission(PermissionReadPost)
	case "write", "update":
		// Users can update their own posts, moderators can update any
		return u.HasPermission(PermissionWritePost) || u.Role == RoleModerator
	case "delete":
		// Users can delete their own posts, moderators can delete any
		return u.HasPermission(PermissionDeletePost) || u.Role == RoleModerator
	default:
		return false
	}
}

// canAccessUser checks user-specific access
func (u *User) canAccessUser(userID, action string) bool {
	switch action {
	case "read":
		// Users can read their own profile, others need permission
		return userID == u.ID || u.HasPermission(PermissionReadUser)
	case "write", "update":
		// Users can update their own profile, admins can update any
		return userID == u.ID || u.HasPermission(PermissionWriteUser)
	case "delete":
		// Only admins can delete users
		return u.HasPermission(PermissionDeleteUser)
	default:
		return false
	}
}

// canAccessComment checks comment-specific access
func (u *User) canAccessComment(commentID, action string) bool {
	switch action {
	case "read":
		return u.HasPermission(PermissionReadComment)
	case "write", "update":
		// Users can update their own comments, moderators can update any
		return u.HasPermission(PermissionWriteComment) || u.Role == RoleModerator
	case "delete":
		// Users can delete their own comments, moderators can delete any
		return u.HasPermission(PermissionDeleteComment) || u.Role == RoleModerator
	default:
		return false
	}
}

// AuthorizationMiddleware provides authorization checks for GraphQL operations
type AuthorizationMiddleware struct {
	rolePermissions map[Role][]Permission
}

// NewAuthorizationMiddleware creates a new authorization middleware
func NewAuthorizationMiddleware() *AuthorizationMiddleware {
	return &AuthorizationMiddleware{
		rolePermissions: map[Role][]Permission{
			RoleAdmin: {
				PermissionReadPost, PermissionWritePost, PermissionDeletePost,
				PermissionReadUser, PermissionWriteUser, PermissionDeleteUser,
				PermissionReadComment, PermissionWriteComment, PermissionDeleteComment,
				PermissionModerate, PermissionAdmin,
			},
			RoleModerator: {
				PermissionReadPost, PermissionWritePost, PermissionDeletePost,
				PermissionReadUser, PermissionReadComment, PermissionWriteComment,
				PermissionDeleteComment, PermissionModerate,
			},
			RoleUser: {
				PermissionReadPost, PermissionWritePost,
				PermissionReadUser, PermissionReadComment, PermissionWriteComment,
			},
			RoleGuest: {
				PermissionReadPost, PermissionReadComment,
			},
		},
	}
}

// ExtensionName returns the name of this extension
func (a *AuthorizationMiddleware) ExtensionName() string {
	return "AuthorizationMiddleware"
}

// Validate validates the schema (no-op for this extension)
func (a *AuthorizationMiddleware) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

// InterceptField intercepts field resolution to check authorization
func (a *AuthorizationMiddleware) InterceptField(ctx context.Context, next graphql.Resolver) (interface{}, error) {
	fc := graphql.GetFieldContext(ctx)
	user := GetUserFromContext(ctx)
	
	// Check if field requires authorization
	if a.requiresAuthorization(fc.Field.Name) {
		if user == nil {
			return nil, fmt.Errorf("authentication required to access %s", fc.Field.Name)
		}
		
		// Check specific field permissions
		if !a.checkFieldPermission(user, fc.Field.Name, fc.Args) {
			return nil, fmt.Errorf("insufficient permissions to access %s", fc.Field.Name)
		}
	}
	
	return next(ctx)
}

// requiresAuthorization checks if a field requires authorization
func (a *AuthorizationMiddleware) requiresAuthorization(fieldName string) bool {
	protectedFields := map[string]bool{
		// Mutations
		"createPost":    true,
		"updatePost":    true,
		"deletePost":    true,
		"createComment": true,
		"updateComment": true,
		"deleteComment": true,
		"updateUser":    true,
		"deleteUser":    true,
		
		// Sensitive queries
		"userProfile":   true,
		"adminUsers":    true,
		"moderatePost":  true,
		"userComments":  true,
	}
	
	return protectedFields[fieldName]
}

// checkFieldPermission checks if user has permission for specific field
func (a *AuthorizationMiddleware) checkFieldPermission(user *User, fieldName string, args map[string]interface{}) bool {
	switch fieldName {
	case "createPost", "updatePost":
		return user.HasPermission(PermissionWritePost)
	case "deletePost":
		// Check if user owns the post or has delete permission
		if postID, ok := args["id"].(string); ok {
			return user.CanAccessResource("post", postID, "delete")
		}
		return user.HasPermission(PermissionDeletePost)
	case "createComment", "updateComment":
		return user.HasPermission(PermissionWriteComment)
	case "deleteComment":
		if commentID, ok := args["id"].(string); ok {
			return user.CanAccessResource("comment", commentID, "delete")
		}
		return user.HasPermission(PermissionDeleteComment)
	case "updateUser":
		if userID, ok := args["id"].(string); ok {
			return user.CanAccessResource("user", userID, "update")
		}
		return user.HasPermission(PermissionWriteUser)
	case "deleteUser":
		return user.HasPermission(PermissionDeleteUser)
	case "userProfile":
		if userID, ok := args["id"].(string); ok {
			return user.CanAccessResource("user", userID, "read")
		}
		return true // Can read own profile
	case "adminUsers":
		return user.HasRole(RoleAdmin)
	case "moderatePost":
		return user.HasPermission(PermissionModerate)
	default:
		return true
	}
}

// GetUserFromContext extracts user from GraphQL context
func GetUserFromContext(ctx context.Context) *User {
	if user, ok := ctx.Value("user").(*User); ok {
		return user
	}
	return nil
}

// WithUser adds user to context
func WithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, "user", user)
}

// RequireAuth is a helper function for resolvers to require authentication
func RequireAuth(ctx context.Context) (*User, error) {
	user := GetUserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("authentication required")
	}
	if !user.IsActive {
		return nil, fmt.Errorf("account is inactive")
	}
	if !user.IsVerified {
		return nil, fmt.Errorf("account is not verified")
	}
	return user, nil
}

// RequireRole is a helper function to require specific role
func RequireRole(ctx context.Context, role Role) (*User, error) {
	user, err := RequireAuth(ctx)
	if err != nil {
		return nil, err
	}
	if !user.HasRole(role) {
		return nil, fmt.Errorf("insufficient role: required %s", role)
	}
	return user, nil
}

// RequirePermission is a helper function to require specific permission
func RequirePermission(ctx context.Context, permission Permission) (*User, error) {
	user, err := RequireAuth(ctx)
	if err != nil {
		return nil, err
	}
	if !user.HasPermission(permission) {
		return nil, fmt.Errorf("insufficient permission: required %s", permission)
	}
	return user, nil
}

// RequireOwnership is a helper function to require resource ownership
func RequireOwnership(ctx context.Context, resourceType, resourceID string) (*User, error) {
	user, err := RequireAuth(ctx)
	if err != nil {
		return nil, err
	}
	
	// Admin can access everything
	if user.HasRole(RoleAdmin) {
		return user, nil
	}
	
	// Check ownership based on resource type
	switch resourceType {
	case "post":
		// This would typically involve a database lookup
		// For now, we'll assume the resourceID format includes owner info
		if !strings.HasPrefix(resourceID, user.ID+":") {
			return nil, fmt.Errorf("access denied: not the owner of this %s", resourceType)
		}
	case "user":
		if resourceID != user.ID {
			return nil, fmt.Errorf("access denied: can only access own user data")
		}
	default:
		return nil, fmt.Errorf("unknown resource type: %s", resourceType)
	}
	
	return user, nil
}

// AuditLog represents an audit log entry
type AuditLog struct {
	UserID      string                 `json:"user_id"`
	Action      string                 `json:"action"`
	Resource    string                 `json:"resource"`
	ResourceID  string                 `json:"resource_id"`
	Timestamp   int64                  `json:"timestamp"`
	IPAddress   string                 `json:"ip_address"`
	UserAgent   string                 `json:"user_agent"`
	Success     bool                   `json:"success"`
	Error       string                 `json:"error,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// AuditLogger logs security-related events
type AuditLogger struct {
	// In a real implementation, this would write to a database or log service
}

// NewAuditLogger creates a new audit logger
func NewAuditLogger() *AuditLogger {
	return &AuditLogger{}
}

// LogAccess logs access attempts
func (a *AuditLogger) LogAccess(ctx context.Context, user *User, action, resource, resourceID string, success bool, err error) {
	log := AuditLog{
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Timestamp:  time.Now().Unix(),
		Success:    success,
	}
	
	if user != nil {
		log.UserID = user.ID
	}
	
	if err != nil {
		log.Error = err.Error()
	}
	
	// Extract IP and User-Agent from context
	if ip, ok := ctx.Value("client_ip").(string); ok {
		log.IPAddress = ip
	}
	if ua, ok := ctx.Value("user_agent").(string); ok {
		log.UserAgent = ua
	}
	
	// In a real implementation, this would be written to a persistent store
	fmt.Printf("AUDIT: %+v\n", log)
}