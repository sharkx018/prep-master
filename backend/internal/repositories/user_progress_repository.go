package repositories

import (
	"database/sql"
	"fmt"
	"interview-prep-app/internal/models"
	"time"
)

// UserProgressRepository handles database operations for user progress
type UserProgressRepository struct {
	db *sql.DB
}

// NewUserProgressRepository creates a new UserProgressRepository
func NewUserProgressRepository(db *sql.DB) *UserProgressRepository {
	return &UserProgressRepository{db: db}
}

// Create creates a new user progress record
func (r *UserProgressRepository) Create(progress *models.UserProgress) error {
	query := `
		INSERT INTO user_progress (user_id, item_id, status, notes, started_at, completed_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	progress.CreatedAt = now
	progress.UpdatedAt = now

	err := r.db.QueryRow(
		query,
		progress.UserID,
		progress.ItemID,
		progress.Status,
		progress.Notes,
		progress.StartedAt,
		progress.CompletedAt,
		progress.CreatedAt,
		progress.UpdatedAt,
	).Scan(&progress.ID, &progress.CreatedAt, &progress.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user progress: %w", err)
	}

	return nil
}

// GetByUserAndItem retrieves user progress for a specific user and item
func (r *UserProgressRepository) GetByUserAndItem(userID, itemID int) (*models.UserProgress, error) {
	query := `
		SELECT id, user_id, item_id, status, notes, started_at, completed_at, created_at, updated_at
		FROM user_progress
		WHERE user_id = $1 AND item_id = $2
	`

	progress := &models.UserProgress{}
	err := r.db.QueryRow(query, userID, itemID).Scan(
		&progress.ID,
		&progress.UserID,
		&progress.ItemID,
		&progress.Status,
		&progress.Notes,
		&progress.StartedAt,
		&progress.CompletedAt,
		&progress.CreatedAt,
		&progress.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user progress not found")
		}
		return nil, fmt.Errorf("failed to get user progress: %w", err)
	}

	return progress, nil
}

// Update updates a user progress record
func (r *UserProgressRepository) Update(progress *models.UserProgress) error {
	query := `
		UPDATE user_progress
		SET status = $1, notes = $2, started_at = $3, completed_at = $4, updated_at = $5
		WHERE id = $6
	`

	progress.UpdatedAt = time.Now()

	_, err := r.db.Exec(
		query,
		progress.Status,
		progress.Notes,
		progress.StartedAt,
		progress.CompletedAt,
		progress.UpdatedAt,
		progress.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user progress: %w", err)
	}

	return nil
}

// GetByUserID retrieves all progress records for a user
func (r *UserProgressRepository) GetByUserID(userID int) ([]*models.UserProgress, error) {
	query := `
		SELECT id, user_id, item_id, status, notes, started_at, completed_at, created_at, updated_at
		FROM user_progress
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user progress: %w", err)
	}
	defer rows.Close()

	var progressList []*models.UserProgress
	for rows.Next() {
		progress := &models.UserProgress{}
		err := rows.Scan(
			&progress.ID,
			&progress.UserID,
			&progress.ItemID,
			&progress.Status,
			&progress.Notes,
			&progress.StartedAt,
			&progress.CompletedAt,
			&progress.CreatedAt,
			&progress.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user progress: %w", err)
		}
		progressList = append(progressList, progress)
	}

	return progressList, nil
}

// Delete deletes a user progress record
func (r *UserProgressRepository) Delete(id int) error {
	query := `DELETE FROM user_progress WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user progress: %w", err)
	}

	return nil
}
