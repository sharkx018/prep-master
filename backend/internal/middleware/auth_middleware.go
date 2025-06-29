package middleware

import (
	"interview-prep-app/internal/handlers"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates a JWT authentication middleware
func AuthMiddleware(authHandler *handlers.AuthHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token
		claims, err := authHandler.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user info in context for use in handlers
		c.Set("username", claims.Username)
		c.Next()
	}
}
