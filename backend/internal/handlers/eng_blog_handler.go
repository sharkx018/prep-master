package handlers

import (
	"net/http"
	"strconv"

	"interview-prep-app/internal/models"
	"interview-prep-app/internal/repositories"

	"github.com/gin-gonic/gin"
)

// EngBlogHandler handles HTTP requests for engineering blogs
type EngBlogHandler struct {
	engBlogRepo *repositories.EngBlogRepository
}

// NewEngBlogHandler creates a new engineering blog handler
func NewEngBlogHandler(engBlogRepo *repositories.EngBlogRepository) *EngBlogHandler {
	return &EngBlogHandler{
		engBlogRepo: engBlogRepo,
	}
}

// GetEngBlogs handles GET /eng-blogs - Returns all engineering blogs
func (h *EngBlogHandler) GetEngBlogs(c *gin.Context) {
	// Get optional query parameters
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	var limit, offset int
	var err error

	if limitStr != "" {
		if limit, err = strconv.Atoi(limitStr); err != nil || limit < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
			return
		}
	}

	if offsetStr != "" {
		if offset, err = strconv.Atoi(offsetStr); err != nil || offset < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
			return
		}
	}

	// Get blogs from database
	blogs, total, err := h.engBlogRepo.GetAll(limit, offset)
	if err != nil {
		gin.DefaultErrorWriter.Write([]byte("Error loading engineering blogs from database: " + err.Error() + "\n"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load engineering blogs data"})
		return
	}

	response := models.EngBlogsResponse{
		Blogs: blogs,
		Total: total,
	}

	c.JSON(http.StatusOK, response)
}

// GetEngBlog handles GET /eng-blogs/:id - Returns a specific engineering blog
func (h *EngBlogHandler) GetEngBlog(c *gin.Context) {
	id := c.Param("id")

	blog, err := h.engBlogRepo.GetByID(id)
	if err != nil {
		gin.DefaultErrorWriter.Write([]byte("Error loading engineering blog by ID: " + err.Error() + "\n"))
		c.JSON(http.StatusNotFound, gin.H{"error": "Engineering blog not found"})
		return
	}

	c.JSON(http.StatusOK, blog)
}
