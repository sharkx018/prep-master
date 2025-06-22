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
	stats, err := h.statsService.GetOverallStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetDetailedStats handles GET /stats/detailed
func (h *StatsHandler) GetDetailedStats(c *gin.Context) {
	stats, err := h.statsService.GetDetailedStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetCategoryStats handles GET /stats/category/:category
func (h *StatsHandler) GetCategoryStats(c *gin.Context) {
	categoryStr := c.Param("category")
	category := models.Category(categoryStr)

	stats, err := h.statsService.GetCategoryStats(category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetSubcategoryStats handles GET /stats/category/:category/subcategory/:subcategory
func (h *StatsHandler) GetSubcategoryStats(c *gin.Context) {
	categoryStr := c.Param("category")
	category := models.Category(categoryStr)
	subcategory := c.Param("subcategory")

	stats, err := h.statsService.GetSubcategoryStats(category, subcategory)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// ResetCompletedAllCount handles POST /stats/reset-completed-all
func (h *StatsHandler) ResetCompletedAllCount(c *gin.Context) {
	err := h.statsService.ResetCompletedAllCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Completed all count reset to zero"})
}
