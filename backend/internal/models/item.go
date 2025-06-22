package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// Category represents the different types of interview prep categories
type Category string

const (
	CategoryDSA Category = "dsa"
	CategoryLLD Category = "lld"
	CategoryHLD Category = "hld"
)

// Status represents the completion status of an item
type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in-progress"
	StatusDone       Status = "done"
)

// Attachments represents a JSON map for dynamic attributes
type Attachments map[string]string

// Value implements the driver.Valuer interface for database storage
func (a Attachments) Value() (driver.Value, error) {
	if a == nil {
		return "{}", nil
	}
	return json.Marshal(a)
}

// Scan implements the sql.Scanner interface for database retrieval
func (a *Attachments) Scan(value interface{}) error {
	if value == nil {
		*a = make(Attachments)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into Attachments", value)
	}

	return json.Unmarshal(bytes, a)
}

// Item represents an interview preparation item
type Item struct {
	ID          int         `json:"id" db:"id"`
	Title       string      `json:"title" db:"title"`
	Link        string      `json:"link" db:"link"`
	Category    Category    `json:"category" db:"category"`
	Subcategory string      `json:"subcategory" db:"subcategory"`
	Status      Status      `json:"status" db:"status"`
	Starred     bool        `json:"starred" db:"starred"`
	Attachments Attachments `json:"attachments" db:"attachments"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	CompletedAt *time.Time  `json:"completed_at,omitempty" db:"completed_at"`
}

// CreateItemRequest represents the request payload for creating an item
type CreateItemRequest struct {
	Title       string      `json:"title" binding:"required"`
	Link        string      `json:"link" binding:"required"`
	Category    Category    `json:"category" binding:"required"`
	Subcategory string      `json:"subcategory" binding:"required"`
	Attachments Attachments `json:"attachments,omitempty"`
}

// UpdateItemRequest represents the request payload for updating an item
type UpdateItemRequest struct {
	Title       *string      `json:"title,omitempty"`
	Link        *string      `json:"link,omitempty"`
	Category    *Category    `json:"category,omitempty"`
	Subcategory *string      `json:"subcategory,omitempty"`
	Attachments *Attachments `json:"attachments,omitempty"`
}

// ItemFilter represents filters for querying items
type ItemFilter struct {
	Category    *Category `json:"category,omitempty"`
	Subcategory *string   `json:"subcategory,omitempty"`
	Status      *Status   `json:"status,omitempty"`
	Limit       *int      `json:"limit,omitempty"`
	Offset      *int      `json:"offset,omitempty"`
}

// ValidCategories returns a slice of all valid categories
func ValidCategories() []Category {
	return []Category{CategoryDSA, CategoryLLD, CategoryHLD}
}

// IsValidCategory checks if a category is valid
func IsValidCategory(category Category) bool {
	for _, valid := range ValidCategories() {
		if category == valid {
			return true
		}
	}
	return false
}

// ValidStatuses returns a slice of all valid statuses
func ValidStatuses() []Status {
	return []Status{StatusPending, StatusInProgress, StatusDone}
}

// IsValidStatus checks if a status is valid
func IsValidStatus(status Status) bool {
	for _, valid := range ValidStatuses() {
		if status == valid {
			return true
		}
	}
	return false
}

// Common subcategories for different categories
var CommonSubcategories = map[Category][]string{
	CategoryDSA: {
		"arrays", "strings", "linked-lists", "trees", "graphs",
		"dynamic-programming", "sorting", "searching", "hashing",
		"stack", "queue", "heap", "recursion", "backtracking",
		"greedy", "bit-manipulation", "math", "two-pointers",
		"sliding-window", "divide-conquer", "other",
	},
	CategoryLLD: {
		"design-patterns", "solid-principles", "object-modeling",
		"class-design", "database-design", "api-design",
		"system-components", "scalability", "caching", "other",
	},
	CategoryHLD: {
		"distributed-systems", "microservices", "load-balancing",
		"caching", "databases", "messaging", "storage", "cdn",
		"monitoring", "security", "scalability", "reliability", "other",
	},
}
