package handlers

import (
	"net/http"

	"interview-prep-app/internal/services"

	"github.com/gin-gonic/gin"
)

// TestHandler handles HTTP requests for tests
type TestHandler struct {
	testService *services.TestService
}

// NewTestHandler creates a new test handler
func NewTestHandler(testService *services.TestService) *TestHandler {
	return &TestHandler{
		testService: testService,
	}
}

// CreateTest creates a new test session
// POST /api/v1/tests
func (h *TestHandler) CreateTest(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	uid, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if user can create a test (has miscellaneous item in progress)
	canCreate, err := h.testService.CheckCanCreateTest(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !canCreate {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot create test: no miscellaneous item is currently in progress",
		})
		return
	}

	// Create the test
	response, err := h.testService.CreateTest(uid)
	if err != nil {
		if err.Error() == "user already has an active test" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetActiveTest retrieves the current active test for the user
// GET /api/v1/tests/active
func (h *TestHandler) GetActiveTest(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	uid, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	response, err := h.testService.GetActiveTest(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response == nil {
		c.JSON(http.StatusOK, gin.H{"message": "No active test found"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// CheckCanCreateTest checks if user can create a test
// GET /api/v1/tests/can-create
func (h *TestHandler) CheckCanCreateTest(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	uid, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	canCreate, err := h.testService.CheckCanCreateTest(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"can_create": canCreate,
		"reason":     getCanCreateReason(canCreate),
	})
}

// CompleteTest marks a test as completed
// PUT /api/v1/tests/:session_id/complete
func (h *TestHandler) CompleteTest(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	uid, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	sessionID := c.Param("session_id")
	itemId := c.Param("item_id")

	err := h.testService.CompleteTest(uid, sessionID, itemId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Test marked as completed",
		"session_id": sessionID,
	})
}

// AbandonTest marks a test as abandoned
// PUT /api/v1/tests/:session_id/abandon
func (h *TestHandler) AbandonTest(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	uid, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	sessionID := c.Param("session_id")
	itemId := c.Param("item_id")

	err := h.testService.AbandonTest(uid, sessionID, itemId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Test marked as abandoned",
		"session_id": sessionID,
	})
}

// DeleteTest deletes a test
// DELETE /api/v1/tests/:session_id
func (h *TestHandler) DeleteTest(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	uid, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	sessionID := c.Param("session_id")

	err := h.testService.DeleteTest(uid, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Test deleted successfully",
		"session_id": sessionID,
	})
}

// getCanCreateReason returns a user-friendly reason for can/cannot create test
func getCanCreateReason(canCreate bool) string {
	if canCreate {
		return "You have a miscellaneous item in progress"
	}
	return "No miscellaneous item is currently in progress"
}
