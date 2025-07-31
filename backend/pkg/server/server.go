package server

import (
	"interview-prep-app/internal/config"
	"interview-prep-app/internal/handlers"
	"interview-prep-app/internal/middleware"
	"interview-prep-app/internal/repositories"

	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server
type Server struct {
	config           *config.Config
	router           *gin.Engine
	itemHandler      *handlers.ItemHandler
	statsHandler     *handlers.StatsHandler
	authHandler      *handlers.AuthHandler
	userProgressRepo *repositories.UserProgressRepository
}

// New creates a new server instance
func New(cfg *config.Config, itemHandler *handlers.ItemHandler, statsHandler *handlers.StatsHandler, authHandler *handlers.AuthHandler, userProgressRepo *repositories.UserProgressRepository) *Server {
	// Set Gin mode based on environment
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	return &Server{
		config:           cfg,
		router:           router,
		itemHandler:      itemHandler,
		statsHandler:     statsHandler,
		authHandler:      authHandler,
		userProgressRepo: userProgressRepo,
	}
}

// setupMiddleware configures middleware for the server
func (s *Server) setupMiddleware() {
	// CORS middleware
	s.router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Recovery middleware
	s.router.Use(gin.Recovery())

	// Logger middleware (only in development)
	if s.config.IsDevelopment() {
		s.router.Use(gin.Logger())
	}
}

// setupRoutes configures all routes for the server
func (s *Server) setupRoutes() {
	// Health check (public)
	s.router.GET("/health", s.healthCheck)

	// Authentication routes (public) - Updated
	auth := s.router.Group("/api/v1/auth")
	{
		auth.POST("/register", s.authHandler.Register)
		auth.POST("/login", s.authHandler.Login)
		auth.POST("/oauth/login", s.authHandler.OAuthLogin)
	}

	// LeetCode proxy route (public)
	s.router.POST("/api/v1/leetcode/proxy", func(c *gin.Context) {
		// Convert Gin context to http.ResponseWriter and http.Request
		handlers.LeetCodeProxyHandler(c.Writer, c.Request)
	})

	// Protected API v1 routes
	v1 := s.router.Group("/api/v1")
	v1.Use(middleware.AuthMiddleware(s.authHandler)) // Apply JWT middleware to all v1 routes
	{
		// User routes
		user := v1.Group("/user")
		{
			user.GET("/profile", s.authHandler.GetCurrentUser)
			user.PUT("/profile", s.authHandler.UpdateProfile)
		}

		// Item routes
		items := v1.Group("/items")
		{
			items.POST("", s.itemHandler.CreateItem)
			items.GET("", s.itemHandler.GetItems)
			items.GET("/paginated", s.itemHandler.GetItemsPaginated)
			items.GET("/next", s.itemHandler.GetNextItem)
			items.POST("/skip", s.itemHandler.SkipItem)
			items.GET("/subcategories/:category", s.itemHandler.GetSubcategories)
			items.GET("/:id", s.itemHandler.GetItem)
			items.PUT("/:id", s.itemHandler.UpdateItem)
			items.PUT("/:id/complete", s.itemHandler.CompleteItem)
			items.PUT("/:id/star", s.itemHandler.ToggleStar)
			items.PUT("/:id/status", s.itemHandler.UpdateStatus)
			items.DELETE("/:id", s.itemHandler.DeleteItem)
			items.POST("/reset", s.itemHandler.ResetItems)
		}

		// Stats routes
		stats := v1.Group("/stats")
		{
			stats.GET("", s.statsHandler.GetStats)
			stats.GET("/detailed", s.statsHandler.GetDetailedStats)
			stats.GET("/category/:category", s.statsHandler.GetCategoryStats)
			stats.GET("/category/:category/subcategory/:subcategory", s.statsHandler.GetSubcategoryStats)
			stats.POST("/reset-completed-all", s.statsHandler.ResetCompletedAllCount)
		}
	}

	// Legacy routes (for backward compatibility) - also protected
	legacyProtected := s.router.Group("")
	legacyProtected.Use(middleware.AuthMiddleware(s.authHandler))
	{
		legacyProtected.POST("/items", s.itemHandler.CreateItem)
		legacyProtected.GET("/items", s.itemHandler.GetItems)
		legacyProtected.GET("/items/next", s.itemHandler.GetNextItem)
		legacyProtected.POST("/items/skip", s.itemHandler.SkipItem)
		legacyProtected.PUT("/items/:id/complete", s.itemHandler.CompleteItem)
		legacyProtected.GET("/stats", s.statsHandler.GetStats)
		legacyProtected.POST("/reset", s.itemHandler.ResetItems)
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.setupMiddleware()
	s.setupRoutes()

	return s.router.Run(":" + s.config.Port)
}

// healthCheck handles health check requests
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "healthy",
		"message": "Interview Prep API is running",
		"version": "2.0",
	})
}
