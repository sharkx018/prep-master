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
	CategoryDSA           Category = "dsa"
	CategoryLLD           Category = "lld"
	CategoryHLD           Category = "hld"
	CategoryMiscellaneous Category = "miscellaneous"
)

// Status represents the completion status of an item
type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in-progress"
	StatusDone       Status = "done"
)

// Special subcategory constants
const (
	Test_n_revise = "test_n_revise"
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
	Attachments Attachments `json:"attachments" db:"attachments"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
}

// ItemWithProgress represents an item with user-specific progress data
type ItemWithProgress struct {
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
	Notes       string      `json:"notes,omitempty" db:"notes"`
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
	RandomOrder *bool     `json:"random_order,omitempty"`
}

// PaginatedItemsResponse represents a paginated response for items
type PaginatedItemsResponse struct {
	Items      []*ItemWithProgress `json:"items"`
	Pagination PaginationMeta      `json:"pagination"`
}

// PaginationMeta contains pagination metadata
type PaginationMeta struct {
	Total      int  `json:"total"`
	Limit      int  `json:"limit"`
	Offset     int  `json:"offset"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
	TotalPages int  `json:"total_pages"`
	Page       int  `json:"page"`
}

// ValidCategories returns a slice of all valid categories
func ValidCategories() []Category {
	return []Category{CategoryDSA, CategoryLLD, CategoryHLD, CategoryMiscellaneous}
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
		"arrays",
		"strings",
		"two-pointers",
		"sliding window - fixed size",
		"sliding window - dynamic size",
		"prefix-sum",
		"kadane's algorithm",
		"matrix (2d array)",
		"linked-lists",
		"linkedList in-place reversal",
		"fast and slow pointers",
		"stacks",
		"monotonic stack",
		"queues",
		"monotonic queue",
		"hashing",
		"bit-manipulation",
		"bucket sort",
		"recursion",
		"divide-conquer",
		"merge sort",
		"quickSort / quickSelect",
		"binary search",
		"backtracking",
		"tree traversal - level order",
		"tree traversal - pre order",
		"tree traversal - in order",
		"tree traversal - post-order",
		"bst / ordered set",
		"tries",
		"heaps",
		"two heaps",
		"top k elements",
		"intervals",
		"k-way merge",
		"data structure design",
		"graphs",
		"depth first search (dfs)",
		"breadth first search (bfs)",
		"topological sort",
		"union find",
		"minimum spanning tree",
		"shortest path",
		"eulerian circuit",
		"greedy",
		"1-d dp",
		"knapsack dp",
		"unbounded knapsack dp",
		"longest increasing subsequence dp",
		"2d (grid) dp",
		"string dp",
		"tree / graph dp",
		"bitmask dp",
		"digit dp",
		"probability dp",
		"state machine dp",
		"string matching",
		"binary indexed tree / segment tree",
		"maths / geometry",
		"line sweep",
		"suffix array",
		"other",
	},
	CategoryLLD: {
		"object-oriented-programming",
		"design-principles",
		"uml",
		"design-patterns-creational",
		"design-patterns-structural",
		"design-patterns-behavioral",
		"lld-interview-tips",
		"lld-interview-questions",
	},
	CategoryHLD: {
		"introduction",
		"core concepts",
		"databases and storage",
		"database scaling techniques",
		"caching",
		"networking",
		"api",
		"asynchronous communications",
		"tradeoffs",
		"distributed system concepts",
		"microservices",
		"big data processing",
		"architectural patterns",
		"observability",
		"security",
		"interview tips",
		"interview questions",
	},
	CategoryMiscellaneous: {
		"gre",
		"finance",
		"development",
		"sql",
		"books",
		Test_n_revise,
		"other",
	},
}
