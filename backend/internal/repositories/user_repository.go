package repositories

import (
	"database/sql"
	"fmt"
	"interview-prep-app/internal/models"
	"time"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (email, name, avatar, auth_provider, provider_id, password_hash, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.IsActive = true

	// Handle empty provider_id as NULL for email users
	var providerID interface{}
	if user.ProviderID == "" {
		providerID = nil
	} else {
		providerID = user.ProviderID
	}

	err := r.db.QueryRow(
		query,
		user.Email,
		user.Name,
		user.Avatar,
		user.AuthProvider,
		providerID,
		user.PasswordHash,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id int) (*models.User, error) {
	query := `
		SELECT id, email, name, avatar, auth_provider, provider_id, password_hash, is_active, created_at, updated_at, last_login_at
		FROM users
		WHERE id = $1 AND is_active = true
	`

	user := &models.User{}
	var providerID sql.NullString
	var lastLoginAt sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Avatar,
		&user.AuthProvider,
		&providerID,
		&user.PasswordHash,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&lastLoginAt,
	)

	// Handle NULL values
	if providerID.Valid {
		user.ProviderID = providerID.String
	}
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, name, avatar, auth_provider, provider_id, password_hash, is_active, created_at, updated_at, last_login_at
		FROM users
		WHERE email = $1 AND is_active = true
	`

	user := &models.User{}
	var providerID sql.NullString
	var lastLoginAt sql.NullTime

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Avatar,
		&user.AuthProvider,
		&providerID,
		&user.PasswordHash,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&lastLoginAt,
	)

	// Handle NULL values
	if providerID.Valid {
		user.ProviderID = providerID.String
	}
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetByProviderID retrieves a user by provider and provider ID
func (r *UserRepository) GetByProviderID(provider models.AuthProvider, providerID string) (*models.User, error) {
	query := `
		SELECT id, email, name, avatar, auth_provider, provider_id, password_hash, is_active, created_at, updated_at, last_login_at
		FROM users
		WHERE auth_provider = $1 AND provider_id = $2 AND is_active = true
	`

	user := &models.User{}
	var providerIDResult sql.NullString
	var lastLoginAt sql.NullTime

	err := r.db.QueryRow(query, provider, providerID).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Avatar,
		&user.AuthProvider,
		&providerIDResult,
		&user.PasswordHash,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&lastLoginAt,
	)

	// Handle NULL values
	if providerIDResult.Valid {
		user.ProviderID = providerIDResult.String
	}
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users
		SET name = $2, avatar = $3, updated_at = $4
		WHERE id = $1 AND is_active = true
		RETURNING updated_at
	`

	user.UpdatedAt = time.Now()

	err := r.db.QueryRow(
		query,
		user.ID,
		user.Name,
		user.Avatar,
		user.UpdatedAt,
	).Scan(&user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// UpdateLastLogin updates the last login time for a user
func (r *UserRepository) UpdateLastLogin(userID int) error {
	query := `
		UPDATE users
		SET last_login_at = $2, updated_at = $2
		WHERE id = $1 AND is_active = true
	`

	now := time.Now()
	_, err := r.db.Exec(query, userID, now)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}

// EmailExists checks if an email already exists
func (r *UserRepository) EmailExists(email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = $1 AND is_active = true`

	var count int
	err := r.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return count > 0, nil
}

// ProviderUserExists checks if a provider user already exists
func (r *UserRepository) ProviderUserExists(provider models.AuthProvider, providerID string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE auth_provider = $1 AND provider_id = $2 AND is_active = true`

	var count int
	err := r.db.QueryRow(query, provider, providerID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check provider user existence: %w", err)
	}

	return count > 0, nil
}

// Deactivate deactivates a user (soft delete)
func (r *UserRepository) Deactivate(userID int) error {
	query := `
		UPDATE users
		SET is_active = false, updated_at = $2
		WHERE id = $1
	`

	_, err := r.db.Exec(query, userID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	return nil
}

// CreateRefreshToken creates a new refresh token
func (r *UserRepository) CreateRefreshToken(userID int, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at, created_at, is_revoked)
		VALUES ($1, $2, $3, $4, false)
	`

	_, err := r.db.Exec(query, userID, token, expiresAt, time.Now())
	if err != nil {
		return fmt.Errorf("failed to create refresh token: %w", err)
	}

	return nil
}

// GetRefreshToken retrieves a refresh token
func (r *UserRepository) GetRefreshToken(token string) (*models.RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at, is_revoked
		FROM refresh_tokens
		WHERE token = $1
	`

	refreshToken := &models.RefreshToken{}
	err := r.db.QueryRow(query, token).Scan(
		&refreshToken.ID,
		&refreshToken.UserID,
		&refreshToken.Token,
		&refreshToken.ExpiresAt,
		&refreshToken.CreatedAt,
		&refreshToken.IsRevoked,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("refresh token not found")
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	return refreshToken, nil
}

// RevokeRefreshToken revokes a refresh token
func (r *UserRepository) RevokeRefreshToken(token string) error {
	query := `
		UPDATE refresh_tokens
		SET is_revoked = true
		WHERE token = $1
	`

	_, err := r.db.Exec(query, token)
	if err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	return nil
}

// RevokeAllUserRefreshTokens revokes all refresh tokens for a user
func (r *UserRepository) RevokeAllUserRefreshTokens(userID int) error {
	query := `
		UPDATE refresh_tokens
		SET is_revoked = true
		WHERE user_id = $1
	`

	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke user refresh tokens: %w", err)
	}

	return nil
}

// CleanupExpiredRefreshTokens removes expired refresh tokens
func (r *UserRepository) CleanupExpiredRefreshTokens() error {
	query := `
		DELETE FROM refresh_tokens
		WHERE expires_at < $1 OR is_revoked = true
	`

	_, err := r.db.Exec(query, time.Now())
	if err != nil {
		return fmt.Errorf("failed to cleanup expired refresh tokens: %w", err)
	}

	return nil
}
