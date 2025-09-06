package models

import (
	"time"
)

// EngBlogProblem represents a practice problem/article within an engineering blog
type EngBlogProblem struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	OrderIdx     int    `json:"order_idx"`
	ExternalLink string `json:"external_link"`
}

// EngBlog represents an engineering blog company with its articles
type EngBlog struct {
	ID               string           `json:"id"`
	Name             string           `json:"name"`
	Link             string           `json:"link"`
	OrderIdx         int              `json:"order_idx"`
	PracticeProblems []EngBlogProblem `json:"practice_problems"`
}

// EngBlogsResponse represents the response structure for eng blogs API
type EngBlogsResponse struct {
	Blogs []EngBlog `json:"blogs"`
	Total int       `json:"total"`
}

// Database models for eng_blogs tables

// EngBlogDB represents an engineering blog in the database
type EngBlogDB struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Link      string    `json:"link" db:"link"`
	OrderIdx  int       `json:"order_idx" db:"order_idx"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// EngBlogArticleDB represents an engineering blog article in the database
type EngBlogArticleDB struct {
	ID           int       `json:"id" db:"id"`
	BlogID       int       `json:"blog_id" db:"blog_id"`
	Title        string    `json:"title" db:"title"`
	OrderIdx     int       `json:"order_idx" db:"order_idx"`
	ExternalLink string    `json:"external_link" db:"external_link"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// EngBlogWithArticles represents a blog with its articles from the database
type EngBlogWithArticles struct {
	EngBlogDB
	Articles []EngBlogArticleDB `json:"articles"`
}
