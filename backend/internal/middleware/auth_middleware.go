package middleware

import (
	"interview-prep-app/internal/handlers"
	"interview-prep-app/internal/models"
	"interview-prep-app/internal/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates a middleware that validates JWT tokens
func AuthMiddleware(authHandler *handlers.AuthHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Validate token
		claims, err := authHandler.ValidateToken(bearerToken[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)
		c.Set("username", claims.Username) // For backward compatibility
		c.Next()
	}
}

// RequireRole creates a middleware that requires a specific role
func RequireRole(userService *services.UserService, requiredRole models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (should be set by AuthMiddleware)
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Get user from database to check role
		user, err := userService.GetByID(userID.(int))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Check if user has required role
		if user.Role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		// Set user role in context for convenience
		c.Set("userRole", user.Role)
		c.Next()
	}
}

// RequireAdmin creates a middleware that requires admin role
func RequireAdmin(userService *services.UserService) gin.HandlerFunc {
	return RequireRole(userService, models.RoleAdmin)
}
