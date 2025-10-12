package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"interview-prep-app/internal/database"
	"interview-prep-app/internal/models"
	"interview-prep-app/internal/repositories"
)

func main() {


	// Load configuration
	// DatabaseURL := "postgresql://interview_user:interview_pass@localhost:5432/interview_prep?sslmode=disable"
	DatabaseURL := "ggg"
	filePath := "./eng-blogs.json"

	// Initialize database
	db, err := database.NewConnection(DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repository
	engBlogRepo := repositories.NewEngBlogRepository(db)

	// Load JSON data
	jsonPath := filepath.Join(filePath)
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		log.Fatal("Failed to read JSON file:", err)
	}

	// Parse JSON data
	var blogs []models.EngBlog
	if err := json.Unmarshal(data, &blogs); err != nil {
		log.Fatal("Failed to parse JSON:", err)
	}

	log.Printf("Found %d engineering blogs to migrate", len(blogs))

	// Migrate each blog
	for _, blog := range blogs {
		log.Printf("Migrating blog: %s", blog.Name)

		// Create the blog
		blogDB, err := engBlogRepo.CreateBlog(blog.Name, blog.Link, blog.OrderIdx)
		if err != nil {
			log.Printf("Failed to create blog %s: %v", blog.Name, err)
			continue
		}

		// Create articles for this blog
		for _, article := range blog.PracticeProblems {
			_, err := engBlogRepo.CreateArticle(blogDB.ID, article.Title, article.ExternalLink, article.OrderIdx)
			if err != nil {
				log.Printf("Failed to create article %s for blog %s: %v", article.Title, blog.Name, err)
				continue
			}
		}

		log.Printf("Successfully migrated blog %s with %d articles", blog.Name, len(blog.PracticeProblems))
	}

	log.Printf("Migration completed! Migrated %d blogs", len(blogs))
}
