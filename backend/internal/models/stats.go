package models

// Stats represents the progress statistics
type Stats struct {
	TotalItems         int     `json:"total_items"`
	CompletedItems     int     `json:"completed_items"`
	PendingItems       int     `json:"pending_items"`
	ProgressPercentage float64 `json:"progress_percentage"`
	CompletedAllCount  int     `json:"completed_all_count"`
	CurrentStreak      int     `json:"current_streak"`
	LongestStreak      int     `json:"longest_streak"`
}

// AppStats represents the application-level statistics stored in database
type AppStats struct {
	ID                int `json:"id" db:"id"`
	CompletedAllCount int `json:"completed_all_count" db:"completed_all_count"`
}

// CategoryStats represents statistics for a specific category
type CategoryStats struct {
	Category           Category `json:"category"`
	TotalItems         int      `json:"total_items"`
	CompletedItems     int      `json:"completed_items"`
	PendingItems       int      `json:"pending_items"`
	ProgressPercentage float64  `json:"progress_percentage"`
}

// SubcategoryStats represents statistics for a specific subcategory
type SubcategoryStats struct {
	Subcategory        string  `json:"subcategory"`
	TotalItems         int     `json:"total_items"`
	CompletedItems     int     `json:"completed_items"`
	PendingItems       int     `json:"pending_items"`
	ProgressPercentage float64 `json:"progress_percentage"`
}

// CategoryWithSubcategoryStats represents category statistics with subcategory breakdown
type CategoryWithSubcategoryStats struct {
	Category           Category           `json:"category"`
	TotalItems         int                `json:"total_items"`
	CompletedItems     int                `json:"completed_items"`
	PendingItems       int                `json:"pending_items"`
	ProgressPercentage float64            `json:"progress_percentage"`
	Subcategories      []SubcategoryStats `json:"subcategories"`
}

// DetailedStats represents comprehensive statistics including category breakdown
type DetailedStats struct {
	Overall    Stats                          `json:"overall"`
	Categories []CategoryWithSubcategoryStats `json:"categories"`
}
