package handlers

import (
	"net/http"

	"interview-prep-app/internal/models"
	"interview-prep-app/internal/services"

	"github.com/gin-gonic/gin"
)

// StatsHandler handles HTTP requests for statistics
type StatsHandler struct {
	statsService *services.StatsService
}

// NewStatsHandler creates a new stats handler
func NewStatsHandler(statsService *services.StatsService) *StatsHandler {
	return &StatsHandler{statsService: statsService}
}

// GetStats handles GET /stats
func (h *StatsHandler) GetStats(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Use the new method that gets user-specific statistics
	stats, err := h.statsService.GetOverallStatsForUser(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetDetailedStats handles GET /stats/detailed
func (h *StatsHandler) GetDetailedStats(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Use the new method that gets user-specific detailed statistics
	stats, err := h.statsService.GetDetailedStatsForUser(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetCategoryStats handles GET /stats/category/:category
func (h *StatsHandler) GetCategoryStats(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	categoryStr := c.Param("category")
	category := models.Category(categoryStr)

	// Use the new method that gets user-specific category statistics
	stats, err := h.statsService.GetCategoryStatsForUser(userID.(int), category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetSubcategoryStats handles GET /stats/category/:category/subcategory/:subcategory
func (h *StatsHandler) GetSubcategoryStats(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	categoryStr := c.Param("category")
	category := models.Category(categoryStr)
	subcategory := c.Param("subcategory")

	// Use the new method that gets user-specific subcategory statistics
	stats, err := h.statsService.GetSubcategoryStatsForUser(userID.(int), category, subcategory)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// ResetCompletedAllCount handles POST /stats/reset-completed-all
func (h *StatsHandler) ResetCompletedAllCount(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Use the new method that resets user-specific completed all count
	err := h.statsService.ResetUserCompletedAllCount(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Your completed all count has been reset to zero"})
}
