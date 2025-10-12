package models

import (
	"time"
)

// TestStatus represents the status of a test
type TestStatus string

const (
	TestStatusPending    TestStatus = "pending"
	TestStatusCompleted TestStatus = "completed"
	TestStatusAbandoned TestStatus = "abandoned"
)

// Test represents a test session with multiple items
type Test struct {
	ID        int        `json:"id" db:"id"`
	SessionID string     `json:"session_id" db:"session_id"`
	UserID    int        `json:"user_id" db:"user_id"`
	ItemID    int        `json:"item_id" db:"item_id"`
	Status    TestStatus `json:"status" db:"status"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

// TestWithItem represents a test with its associated item details
type TestWithItem struct {
	ID        int              `json:"id" db:"id"`
	SessionID string           `json:"session_id" db:"session_id"`
	UserID    int              `json:"user_id" db:"user_id"`
	ItemID    int              `json:"item_id" db:"item_id"`
	Status    TestStatus       `json:"status" db:"status"`
	CreatedAt time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt time.Time        `json:"updated_at" db:"updated_at"`
	Item      ItemWithProgress `json:"item"`
}

// CreateTestResponse represents the response when creating a test
type CreateTestResponse struct {
	SessionID string             `json:"session_id"`
	Items     []ItemWithProgress `json:"items"`
	Message   string             `json:"message"`
}

// ActiveTestResponse represents the current active test
type ActiveTestResponse struct {
	SessionID string             `json:"session_id"`
	Items     []ItemWithProgress `json:"items"`
	CreatedAt time.Time          `json:"created_at"`
}

// IsValidTestStatus checks if a test status is valid
func IsValidTestStatus(status TestStatus) bool {
	switch status {
	case TestStatusPending, TestStatusCompleted, TestStatusAbandoned:
		return true
	}
	return false
}

// ValidTestStatuses returns a slice of all valid test statuses
func ValidTestStatuses() []TestStatus {
	return []TestStatus{TestStatusPending, TestStatusCompleted, TestStatusAbandoned}
}
