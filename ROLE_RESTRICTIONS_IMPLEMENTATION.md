# Role-Based Access Control Implementation

## Overview

This document outlines the implementation of role-based restrictions where users with the `user` role are prevented from:
1. Accessing the "Add Items" page
2. Deleting items from the items list
3. Editing items from the items list

Only users with the `admin` role can perform these actions.

## Backend Implementation

### 1. Item Handler Updates

**File: `backend/internal/handlers/item_handler.go`**

#### Changes Made:
- Added `userService` dependency to `ItemHandler`
- Added `requireAdminRole()` helper method
- Protected `CreateItem()` endpoint with admin-only access
- Protected `UpdateItem()` endpoint with admin-only access
- Protected `DeleteItem()` endpoint with admin-only access

#### Key Code:
```go
// ItemHandler now includes userService for role checking
type ItemHandler struct {
    itemService *services.ItemService
    userService *services.UserService
}

// Helper method to check admin role
func (h *ItemHandler) requireAdminRole(c *gin.Context) error {
    userID, exists := c.Get("userID")
    if !exists {
        return errors.New("User not authenticated")
    }

    user, err := h.userService.GetByID(userID.(int))
    if err != nil {
        return err
    }

    if user.Role != models.RoleAdmin {
        return errors.New("Admin role required")
    }

    return nil
}

// CreateItem - Admin only
func (h *ItemHandler) CreateItem(c *gin.Context) {
    if err := h.requireAdminRole(c); err != nil {
        c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required to create items"})
        return
    }
    // ... rest of the method
}

// UpdateItem - Admin only
func (h *ItemHandler) UpdateItem(c *gin.Context) {
    if err := h.requireAdminRole(c); err != nil {
        c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required to edit items"})
        return
    }
    // ... rest of the method
}

// DeleteItem - Admin only  
func (h *ItemHandler) DeleteItem(c *gin.Context) {
    if err := h.requireAdminRole(c); err != nil {
        c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required to delete items"})
        return
    }
    // ... rest of the method
}
```

### 2. Main Server Updates

**File: `backend/cmd/server/main.go`**

#### Changes Made:
- Updated `ItemHandler` constructor to include `userService`

```go
// Initialize handlers
itemHandler := handlers.NewItemHandler(itemService, userService)
```

## Frontend Implementation

### 1. AuthContext Updates

**File: `frontend/src/contexts/AuthContext.tsx`**

#### Changes Made:
- Added `role` field to `User` interface
- Added `isAdmin` property to `AuthContextType`
- Added `isAdmin` computed property to context value

```typescript
interface User {
  id: number;
  email: string;
  name: string;
  avatar?: string;
  role: 'user' | 'admin';  // New field
  auth_provider: 'email' | 'google' | 'facebook' | 'apple';
  created_at: string;
  last_login_at?: string;
}

interface AuthContextType {
  // ... existing properties
  isAdmin: boolean;  // New property
}

const value: AuthContextType = {
  // ... existing properties
  isAdmin: user?.role === 'admin',
};
```

### 2. Navigation Updates

**File: `frontend/src/components/Layout.tsx`**

#### Changes Made:
- Conditionally show "Add Item" navigation link only for admin users

```typescript
const navigation = [
  { name: 'Dashboard', href: '/dashboard', icon: LayoutDashboard },
  { name: 'Items', href: '/items', icon: List },
  { name: 'Study', href: '/study', icon: BookOpen },
  ...(isAdmin ? [{ name: 'Add Item', href: '/add-item', icon: Plus }] : []),
  { name: 'Statistics', href: '/stats', icon: BarChart3 },
];
```

### 3. Dashboard Updates

**File: `frontend/src/pages/Dashboard.tsx`**

#### Changes Made:
- Conditionally show "Add New Item" card only for admin users
- Adjust grid layout based on admin status

```typescript
<div className={`grid grid-cols-1 gap-6 ${isAdmin ? 'md:grid-cols-3' : 'md:grid-cols-2'}`}>
  {/* Start Studying card - always visible */}
  
  {isAdmin && (
    <Link to="/add-item">
      {/* Add New Item card - admin only */}
    </Link>
  )}
  
  {/* Browse Items card - always visible */}
</div>
```

### 4. Items Page Updates

**File: `frontend/src/pages/Items.tsx`**

#### Changes Made:
- Hide edit button for non-admin users
- Hide delete button for non-admin users
- Hide edit form for non-admin users
- Added role-based conditional rendering

```typescript
// Edit button - admin only
{isAdmin && (
  <button
    onClick={(e) => {
      e.stopPropagation();
      handleEditClick(item);
    }}
    className="p-2 text-gray-400 hover:text-indigo-600"
    title="Edit item"
  >
    <Edit2 className="h-5 w-5" />
  </button>
)}

// Delete button - admin only
{isAdmin && (
  <button
    onClick={(e) => {
      e.stopPropagation();
      handleDelete(item.id);
    }}
    disabled={deleting === item.id}
    className="p-2 text-gray-400 hover:text-red-600 disabled:opacity-50"
    title="Delete item"
  >
    {deleting === item.id ? (
      <Loader2 className="h-5 w-5 animate-spin" />
    ) : (
      <Trash2 className="h-5 w-5" />
    )}
  </button>
)}

// Edit form - admin only
{editingId === item.id && isAdmin ? (
  // Edit form JSX here
) : (
  // Normal item display
)}
```

### 5. Route Protection

**File: `frontend/src/App.tsx`**

#### Changes Made:
- Created `AdminRoute` component to protect admin-only routes
- Protected `/add-item` route with `AdminRoute`

```typescript
// Component to protect admin-only routes
const AdminRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAdmin } = useAuth();
  
  if (!isAdmin) {
    return <Navigate to="/dashboard" replace />;
  }
  
  return <>{children}</>;
};

// Route protection
<Route path="/add-item" element={
  <AdminRoute>
    <AddItem />
  </AdminRoute>
} />
```

## Security Features

### Backend Security
1. **Endpoint Protection**: Create, update, and delete endpoints check user role before proceeding
2. **User Verification**: Each request validates the user exists and has admin role
3. **Proper Error Responses**: Returns 403 Forbidden with descriptive messages
4. **Authentication Required**: All protected endpoints require valid JWT token

### Frontend Security
1. **UI Hiding**: Admin-only features (add, edit, delete) are hidden from regular users
2. **Route Protection**: Direct navigation to admin routes redirects non-admin users
3. **Conditional Rendering**: Admin features only render for admin users
4. **Navigation Updates**: Admin-only links removed from navigation for regular users
5. **Form Protection**: Edit forms only accessible to admin users

## API Response Changes

### Error Responses
When a regular user tries to access admin-only endpoints:

```json
{
  "error": "Admin access required to create items"
}
```

```json
{
  "error": "Admin access required to edit items"
}
```

```json
{
  "error": "Admin access required to delete items"
}
```

## Testing the Implementation

### Backend Testing

1. **Test Create Item as Regular User**:
```bash
curl -X POST http://localhost:8080/api/v1/items \
  -H "Authorization: Bearer USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title": "Test", "link": "http://example.com", "category": "dsa", "subcategory": "arrays"}'
```
Expected: `403 Forbidden` with admin access required message

2. **Test Update Item as Regular User**:
```bash
curl -X PUT http://localhost:8080/api/v1/items/1 \
  -H "Authorization: Bearer USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title": "Updated Title", "link": "http://example.com", "category": "dsa", "subcategory": "arrays"}'
```
Expected: `403 Forbidden` with admin access required message

3. **Test Delete Item as Regular User**:
```bash
curl -X DELETE http://localhost:8080/api/v1/items/1 \
  -H "Authorization: Bearer USER_TOKEN"
```
Expected: `403 Forbidden` with admin access required message

4. **Test as Admin User**:
Same requests with admin token should succeed.

### Frontend Testing

1. **Regular User Experience**:
   - "Add Item" link not visible in navigation
   - "Add New Item" card not visible on dashboard
   - Edit buttons not visible on items list
   - Delete buttons not visible on items list
   - Edit forms not accessible
   - Direct navigation to `/add-item` redirects to dashboard

2. **Admin User Experience**:
   - All features visible and functional
   - Can access add item page
   - Can edit items inline
   - Can delete items
   - Full administrative capabilities

## User Role Management

### Creating Admin Users

Currently, admin users must be created by:

1. **Database Update**:
```sql
UPDATE users SET role = 'admin' WHERE email = 'admin@example.com';
```

2. **Future Enhancement**: Create admin management interface for role assignments

### Default Behavior
- All new user registrations get `user` role by default
- No self-promotion to admin role possible
- Role changes require database access or admin interface

## Backward Compatibility

- All existing functionality preserved for admin users
- Regular users maintain access to all non-administrative features
- No breaking changes to existing API endpoints
- Existing JWT tokens continue to work

## Future Enhancements

### Planned Features
1. **Admin Management Interface**: UI for admins to manage user roles
2. **Granular Permissions**: More specific permissions beyond admin/user
3. **Audit Logging**: Track admin actions for security
4. **Bulk Operations**: Admin-only bulk item management
5. **User Management**: Admin interface to view/manage all users

### Additional Security Considerations
1. **Rate Limiting**: Implement rate limiting on admin endpoints
2. **Audit Trail**: Log all administrative actions
3. **Session Management**: Enhanced session security for admin users
4. **Two-Factor Authentication**: Additional security for admin accounts

This implementation provides a solid foundation for role-based access control while maintaining security best practices and user experience. 