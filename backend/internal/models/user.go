package models

import (
	"time"
)

// AuthProvider represents different authentication providers
type AuthProvider string

const (
	AuthProviderEmail    AuthProvider = "email"
	AuthProviderGoogle   AuthProvider = "google"
	AuthProviderFacebook AuthProvider = "facebook"
	AuthProviderApple    AuthProvider = "apple"
)

// Role represents user roles in the system
type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// User represents a user in the system
type User struct {
	ID           int          `json:"id" db:"id"`
	Email        string       `json:"email" db:"email"`
	Name         string       `json:"name" db:"name"`
	Avatar       string       `json:"avatar,omitempty" db:"avatar"`
	Role         Role         `json:"role" db:"role"`
	AuthProvider AuthProvider `json:"auth_provider" db:"auth_provider"`
	ProviderID   string       `json:"provider_id,omitempty" db:"provider_id"`
	PasswordHash string       `json:"-" db:"password_hash"` // Never include in JSON
	IsActive     bool         `json:"is_active" db:"is_active"`
	CreatedAt    time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at" db:"updated_at"`
	LastLoginAt  *time.Time   `json:"last_login_at,omitempty" db:"last_login_at"`
}

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Email        string       `json:"email" binding:"required,email"`
	Name         string       `json:"name" binding:"required"`
	Password     string       `json:"password,omitempty" binding:"omitempty,min=6"`
	AuthProvider AuthProvider `json:"auth_provider,omitempty"`
	ProviderID   string       `json:"provider_id,omitempty"`
	Avatar       string       `json:"avatar,omitempty"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	Name   string `json:"name,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// OAuthLoginRequest represents OAuth login request
type OAuthLoginRequest struct {
	Provider    AuthProvider `json:"provider" binding:"required"`
	AccessToken string       `json:"access_token" binding:"required"`
	Email       string       `json:"email,omitempty"`
	Name        string       `json:"name,omitempty"`
	Avatar      string       `json:"avatar,omitempty"`
	ProviderID  string       `json:"provider_id,omitempty"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	User         *User     `json:"user"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// UserProgress represents user progress on an item
type UserProgress struct {
	ID          int        `json:"id" db:"id"`
	UserID      int        `json:"user_id" db:"user_id"`
	ItemID      int        `json:"item_id" db:"item_id"`
	Status      Status     `json:"status" db:"status"`
	Starred     bool       `json:"starred" db:"starred"`
	Notes       string     `json:"notes,omitempty" db:"notes"`
	StartedAt   time.Time  `json:"started_at" db:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// RefreshToken represents a refresh token
type RefreshToken struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	IsRevoked bool      `json:"is_revoked" db:"is_revoked"`
}

// UserStats represents user-specific statistics
type UserStats struct {
	UserID            int        `json:"user_id" db:"user_id"`
	TotalItems        int        `json:"total_items" db:"total_items"`
	CompletedItems    int        `json:"completed_items" db:"completed_items"`
	InProgressItems   int        `json:"in_progress_items" db:"in_progress_items"`
	PendingItems      int        `json:"pending_items" db:"pending_items"`
	DSACompleted      int        `json:"dsa_completed" db:"dsa_completed"`
	LLDCompleted      int        `json:"lld_completed" db:"lld_completed"`
	HLDCompleted      int        `json:"hld_completed" db:"hld_completed"`
	CompletedAllCount int        `json:"completed_all_count" db:"completed_all_count"`
	CurrentStreak     int        `json:"current_streak" db:"current_streak"`
	LongestStreak     int        `json:"longest_streak" db:"longest_streak"`
	LastActivityDate  *time.Time `json:"last_activity_date,omitempty" db:"last_activity_date"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}
