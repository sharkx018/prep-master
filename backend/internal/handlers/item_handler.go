package handlers

import (
	"net/http"
	"strconv"

	"interview-prep-app/internal/models"
	"interview-prep-app/internal/services"

	"github.com/gin-gonic/gin"
)

// ItemHandler handles HTTP requests for items
type ItemHandler struct {
	itemService *services.ItemService
	userService *services.UserService
}

// NewItemHandler creates a new item handler
func NewItemHandler(itemService *services.ItemService, userService *services.UserService) *ItemHandler {
	return &ItemHandler{
		itemService: itemService,
		userService: userService,
	}
}

// CreateItem handles POST /items - Admin only
func (h *ItemHandler) CreateItem(c *gin.Context) {
	// Check if user has admin role
	if err := h.requireAdminRole(c); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required to create items"})
		return
	}

	var req models.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.itemService.CreateItem(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// requireAdminRole checks if the current user has admin role
func (h *ItemHandler) requireAdminRole(c *gin.Context) error {
	userID, exists := c.Get("userID")
	if !exists {
		return gin.Error{Err: gin.Error{}, Type: gin.ErrorTypePublic, Meta: "User not authenticated"}
	}

	user, err := h.userService.GetByID(userID.(int))
	if err != nil {
		return err
	}

	if user.Role != models.RoleAdmin {
		return gin.Error{Err: gin.Error{}, Type: gin.ErrorTypePublic, Meta: "Admin role required"}
	}

	return nil
}

// GetItem handles GET /items/:id
func (h *ItemHandler) GetItem(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	// Use the new method that includes user progress
	item, err := h.itemService.GetItemWithUserProgress(userID.(int), id)
	if err != nil {
		if err.Error() == "item not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// GetItems handles GET /items
func (h *ItemHandler) GetItems(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	filter := &models.ItemFilter{}

	// Parse query parameters
	if categoryStr := c.Query("category"); categoryStr != "" {
		category := models.Category(categoryStr)
		filter.Category = &category
	}

	if subcategory := c.Query("subcategory"); subcategory != "" {
		filter.Subcategory = &subcategory
	}

	if statusStr := c.Query("status"); statusStr != "" {
		status := models.Status(statusStr)
		filter.Status = &status
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
			return
		}
		filter.Limit = &limit
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
			return
		}
		filter.Offset = &offset
	}

	// Use the new method that includes user progress
	items, err := h.itemService.GetItemsWithUserProgress(userID.(int), filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

// GetItemsPaginated handles GET /items/paginated
func (h *ItemHandler) GetItemsPaginated(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	filter := &models.ItemFilter{}

	// Parse query parameters
	if categoryStr := c.Query("category"); categoryStr != "" {
		category := models.Category(categoryStr)
		filter.Category = &category
	}

	if subcategory := c.Query("subcategory"); subcategory != "" {
		filter.Subcategory = &subcategory
	}

	if statusStr := c.Query("status"); statusStr != "" {
		status := models.Status(statusStr)
		filter.Status = &status
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
			return
		}
		filter.Limit = &limit
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
			return
		}
		filter.Offset = &offset
	}

	// Use the new method that includes user progress
	result, err := h.itemService.GetItemsPaginatedWithUserProgress(userID.(int), filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetNextItem handles GET /items/next
func (h *ItemHandler) GetNextItem(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Use the new method that includes user progress
	item, err := h.itemService.GetNextItemWithUserProgress(userID.(int))
	if err != nil {
		if err.Error() == "no pending items found" {
			c.JSON(http.StatusNotFound, gin.H{"message": "No pending items found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// SkipItem handles POST /items/skip
func (h *ItemHandler) SkipItem(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Use the new method that includes user progress
	item, err := h.itemService.SkipItemWithUserProgress(userID.(int))
	if err != nil {
		if err.Error() == "no pending items found" {
			c.JSON(http.StatusNotFound, gin.H{"message": "No pending items found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// CompleteItem handles PUT /items/:id/complete
func (h *ItemHandler) CompleteItem(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	// Use the new method that includes user progress
	item, err := h.itemService.CompleteItemWithUserProgress(userID.(int), id)
	if err != nil {
		if err.Error() == "item not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// UpdateItem handles PUT /items/:id - Admin only
func (h *ItemHandler) UpdateItem(c *gin.Context) {
	// Check if user has admin role
	if err := h.requireAdminRole(c); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required to edit items"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var req models.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.itemService.UpdateItem(id, &req)
	if err != nil {
		if err.Error() == "item not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// DeleteItem handles DELETE /items/:id - Admin only
func (h *ItemHandler) DeleteItem(c *gin.Context) {
	// Check if user has admin role
	if err := h.requireAdminRole(c); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required to delete items"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	err = h.itemService.DeleteItem(id)
	if err != nil {
		if err.Error() == "item not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}

// ResetItems handles POST /items/reset
func (h *ItemHandler) ResetItems(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Use the new method that resets user-specific progress
	rowsAffected, err := h.itemService.ResetAllItemsWithUserProgress(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Your progress has been reset to pending status",
		"items_updated": rowsAffected,
	})
}

// GetSubcategories handles GET /items/subcategories/:category
func (h *ItemHandler) GetSubcategories(c *gin.Context) {
	categoryStr := c.Param("category")
	category := models.Category(categoryStr)

	subcategories, err := h.itemService.GetCommonSubcategories(category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"category":      category,
		"subcategories": subcategories,
	})
}

// ToggleStar handles PUT /items/:id/star
func (h *ItemHandler) ToggleStar(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	// Use the new method that includes user progress
	item, err := h.itemService.ToggleStarWithUserProgress(userID.(int), id)
	if err != nil {
		if err.Error() == "item not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// UpdateStatus handles PUT /items/:id/status
func (h *ItemHandler) UpdateStatus(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := models.Status(req.Status)
	// Use the new method that includes user progress
	item, err := h.itemService.UpdateStatusWithUserProgress(userID.(int), id, status)
	if err != nil {
		if err.Error() == "item not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}
