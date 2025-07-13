package handlers

import (
	"interview-prep-app/internal/models"
	"interview-prep-app/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AdminHandler handles admin-only operations
type AdminHandler struct {
	userService *services.UserService
}

// NewAdminHandler creates a new AdminHandler
func NewAdminHandler(userService *services.UserService) *AdminHandler {
	return &AdminHandler{
		userService: userService,
	}
}

// GetAllUsers returns all users (admin only)
func (h *AdminHandler) GetAllUsers(c *gin.Context) {
	// This endpoint would need a new repository method to get all users
	// For now, just return a success message
	c.JSON(http.StatusOK, gin.H{
		"message": "Admin access granted - this would return all users",
		"role":    c.GetString("userRole"),
	})
}

// UpdateUserRole updates a user's role (admin only)
func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		Role models.Role `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Validate role
	if req.Role != models.RoleUser && req.Role != models.RoleAdmin {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role. Must be 'user' or 'admin'"})
		return
	}

	// For demonstration purposes, just return success
	// In a real implementation, you'd add a method to update user role
	c.JSON(http.StatusOK, gin.H{
		"message":    "User role would be updated",
		"user_id":    userID,
		"new_role":   req.Role,
		"admin_user": c.GetString("userRole"),
	})
}

// GetAdminStats returns admin-specific statistics
func (h *AdminHandler) GetAdminStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Admin statistics",
		"stats": gin.H{
			"total_users":   "This would show total user count",
			"admin_users":   "This would show admin user count",
			"regular_users": "This would show regular user count",
		},
	})
}
