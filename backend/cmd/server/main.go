package main

import (
	"log"

	"interview-prep-app/internal/config"
	"interview-prep-app/internal/database"
	"interview-prep-app/internal/handlers"
	"interview-prep-app/internal/repositories"
	"interview-prep-app/internal/services"
	"interview-prep-app/pkg/server"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("=====>>>>>>>>No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize repositories
	itemRepo := repositories.NewItemRepository(db)
	statsRepo := repositories.NewStatsRepository(db)

	// Initialize services
	itemService := services.NewItemService(itemRepo, statsRepo)
	statsService := services.NewStatsService(itemRepo, statsRepo)

	// Initialize handlers
	itemHandler := handlers.NewItemHandler(itemService)
	statsHandler := handlers.NewStatsHandler(statsService)

	// Initialize and start server
	srv := server.New(cfg, itemHandler, statsHandler)

	log.Printf("Server starting on port %s", cfg.Port)
	log.Printf("Server starting on port %s", cfg)
	if err := srv.Start(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
