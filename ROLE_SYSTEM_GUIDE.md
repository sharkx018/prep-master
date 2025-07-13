# Role System Implementation Guide

## Overview

This guide explains the role-based access control (RBAC) system that has been implemented in the interview preparation application. The system supports two roles: `user` and `admin`, with all new signups automatically assigned the `user` role.

## Role Types

### Available Roles

```go
type Role string

const (
    RoleUser  Role = "user"   // Default role for all new users
    RoleAdmin Role = "admin"  // Administrative role with elevated permissions
)
```

### Default Behavior

- **All new signups** (both email and OAuth) are automatically assigned the `user` role
- The `user` role has access to all standard application features
- The `admin` role has additional administrative capabilities

## Database Schema

### Migration

A new migration has been added to include the `role` column in the `users` table:

```sql
ALTER TABLE users ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'user' 
CHECK (role IN ('user', 'admin'));
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
```

### User Model Updates

The `User` struct now includes the role field:

```go
type User struct {
    ID           int          `json:"id" db:"id"`
    Email        string       `json:"email" db:"email"`
    Name         string       `json:"name" db:"name"`
    Avatar       string       `json:"avatar,omitempty" db:"avatar"`
    Role         Role         `json:"role" db:"role"`           // New field
    AuthProvider AuthProvider `json:"auth_provider" db:"auth_provider"`
    // ... other fields
}
```

## Implementation Details

### 1. Model Changes

- Added `Role` type with `user` and `admin` constants
- Updated `User` struct to include `Role` field
- Modified all user-related database queries to include the role column

### 2. Repository Updates

The `UserRepository` has been updated to:
- Include role in `CREATE` operations with default `user` role
- Include role in all `SELECT` operations (`GetByID`, `GetByEmail`, `GetByProviderID`)
- Automatically set default role for new users if not specified

### 3. Service Updates

The `UserService` has been updated to:
- Assign `models.RoleUser` to all new email registrations
- Assign `models.RoleUser` to all new OAuth registrations

### 4. Middleware Implementation

Two new middleware functions have been added:

#### `RequireRole(userService, requiredRole)`
Generic middleware that requires a specific role:

```go
// Usage example
v1.Use(middleware.RequireRole(userService, models.RoleAdmin))
```

#### `RequireAdmin(userService)`
Convenience middleware that requires admin role:

```go
// Usage example
adminRoutes.Use(middleware.RequireAdmin(userService))
```

## Usage Examples

### Protecting Admin Routes

Here's how to add admin-only routes to your server:

```go
// In server.go setupRoutes() method
func (s *Server) setupRoutes() {
    // ... existing routes ...

    // Admin routes (require admin role)
    admin := v1.Group("/admin")
    admin.Use(middleware.RequireAdmin(s.userService)) // Require admin role
    {
        admin.GET("/users", s.adminHandler.GetAllUsers)
        admin.PUT("/users/:id/role", s.adminHandler.UpdateUserRole)
        admin.GET("/stats", s.adminHandler.GetAdminStats)
    }
}
```

### Example Admin Handler

```go
// AdminHandler example
type AdminHandler struct {
    userService *services.UserService
}

func (h *AdminHandler) GetAllUsers(c *gin.Context) {
    // This endpoint is only accessible to admin users
    // The middleware ensures only admin users reach this point
    
    // Get current user's role from context (set by middleware)
    userRole := c.GetString("userRole") // Will be "admin"
    
    // Implement admin functionality here
}
```

### Checking Roles in Handlers

```go
func (h *SomeHandler) SomeMethod(c *gin.Context) {
    // Get user ID from context (set by AuthMiddleware)
    userID := c.GetInt("userID")
    
    // Get user to check role
    user, err := h.userService.GetByID(userID)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
        return
    }
    
    // Check role
    if user.Role == models.RoleAdmin {
        // Admin-specific logic
    } else {
        // Regular user logic
    }
}
```

## API Response Changes

### User Registration/Login Response

The user object in API responses now includes the role field:

```json
{
    "token": "jwt_token_here",
    "user": {
        "id": 1,
        "email": "user@example.com",
        "name": "John Doe",
        "role": "user",
        "auth_provider": "email",
        "is_active": true,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
    },
    "expires_at": "2024-01-02T00:00:00Z"
}
```

## Security Considerations

### 1. Role Assignment
- New users cannot self-assign admin role during registration
- Role changes should only be possible through admin endpoints
- Default role is always `user` for security

### 2. Middleware Chain
Always use the middleware in the correct order:
```go
// Correct order
routes.Use(middleware.AuthMiddleware(authHandler))  // First: authenticate
routes.Use(middleware.RequireAdmin(userService))    // Then: authorize
```

### 3. Database Constraints
- Role column has CHECK constraint to only allow valid values
- Default value ensures all users have a role
- Index on role column for efficient role-based queries

## Testing the Implementation

### 1. Create a New User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "name": "Test User",
    "password": "password123"
  }'
```

The response should show `"role": "user"`.

### 2. Test Admin Endpoint Access
```bash
# This should fail with 403 Forbidden for regular users
curl -X GET http://localhost:8080/api/v1/admin/users \
  -H "Authorization: Bearer YOUR_USER_TOKEN"
```

### 3. Manually Promote User to Admin
You would need to update the database directly or create an admin endpoint:
```sql
UPDATE users SET role = 'admin' WHERE email = 'test@example.com';
```

## Future Enhancements

### Possible Extensions
1. **More Granular Roles**: Add roles like `moderator`, `viewer`, etc.
2. **Permission System**: Implement permission-based access control
3. **Role Hierarchy**: Define role inheritance (admin inherits user permissions)
4. **Temporary Roles**: Time-limited role assignments
5. **Role Management UI**: Frontend interface for admin role management

### Additional Middleware Options
```go
// Example: Multiple role support
func RequireAnyRole(userService *services.UserService, roles ...models.Role) gin.HandlerFunc

// Example: Permission-based access
func RequirePermission(userService *services.UserService, permission string) gin.HandlerFunc
```

## Migration Notes

### For Existing Users
- All existing users will automatically get the `user` role due to the DEFAULT constraint
- No data migration is needed
- The migration is safe to run on production databases

### Backward Compatibility
- All existing API endpoints continue to work
- New role field is additive and doesn't break existing functionality
- JWT tokens don't need to be regenerated

## Troubleshooting

### Common Issues

1. **Build Errors**: Ensure all imports are updated after adding the role system
2. **Migration Errors**: Check database permissions for ALTER TABLE operations
3. **403 Forbidden**: Verify middleware order and user role assignment

### Debugging Role Issues
```go
// Add logging to debug role assignments
log.Printf("User %d has role: %s", userID, user.Role)
```

This role system provides a solid foundation for access control while maintaining simplicity and extensibility for future requirements. 