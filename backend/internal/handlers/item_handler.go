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
}

// NewItemHandler creates a new item handler
func NewItemHandler(itemService *services.ItemService) *ItemHandler {
	return &ItemHandler{itemService: itemService}
}

// CreateItem handles POST /items
func (h *ItemHandler) CreateItem(c *gin.Context) {
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

// GetItem handles GET /items/:id
func (h *ItemHandler) GetItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	item, err := h.itemService.GetItem(id)
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

	items, err := h.itemService.GetItems(filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

// GetItemsPaginated handles GET /items/paginated
func (h *ItemHandler) GetItemsPaginated(c *gin.Context) {
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

	result, err := h.itemService.GetItemsPaginated(filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetNextItem handles GET /items/next
func (h *ItemHandler) GetNextItem(c *gin.Context) {
	item, err := h.itemService.GetNextItem()
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
	item, err := h.itemService.SkipItem()
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
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	item, err := h.itemService.CompleteItem(id)
	if err != nil {
		if err.Error() == "item not found or not in-progress" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found or not in-progress"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// UpdateItem handles PUT /items/:id
func (h *ItemHandler) UpdateItem(c *gin.Context) {
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

// DeleteItem handles DELETE /items/:id
func (h *ItemHandler) DeleteItem(c *gin.Context) {
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
	rowsAffected, err := h.itemService.ResetAllItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "All items reset to pending status",
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
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	item, err := h.itemService.ToggleStar(id)
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
	item, err := h.itemService.UpdateStatus(id, status)
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
