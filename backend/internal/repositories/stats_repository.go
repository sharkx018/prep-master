package repositories

import (
	"database/sql"
	"fmt"

	"interview-prep-app/internal/models"
)

// StatsRepository handles database operations for app statistics
type StatsRepository struct {
	db *sql.DB
}

// NewStatsRepository creates a new stats repository
func NewStatsRepository(db *sql.DB) *StatsRepository {
	return &StatsRepository{db: db}
}

// GetAppStats retrieves the app-level statistics
func (r *StatsRepository) GetAppStats() (*models.AppStats, error) {
	query := "SELECT id, completed_all_count FROM app_stats WHERE id = 1"

	var stats models.AppStats
	err := r.db.QueryRow(query).Scan(&stats.ID, &stats.CompletedAllCount)
	if err == sql.ErrNoRows {
		// Initialize if not exists
		return r.initializeAppStats()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get app stats: %w", err)
	}

	return &stats, nil
}

// IncrementCompletedAllCount increments the completed_all_count
func (r *StatsRepository) IncrementCompletedAllCount() error {
	query := "UPDATE app_stats SET completed_all_count = completed_all_count + 1 WHERE id = 1"

	result, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to increment completed_all_count: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		// Initialize if not exists and try again
		if _, err := r.initializeAppStats(); err != nil {
			return fmt.Errorf("failed to initialize app stats: %w", err)
		}
		return r.IncrementCompletedAllCount()
	}

	return nil
}

// ResetCompletedAllCount resets the completed_all_count to 0
func (r *StatsRepository) ResetCompletedAllCount() error {
	query := "UPDATE app_stats SET completed_all_count = 0 WHERE id = 1"

	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to reset completed_all_count: %w", err)
	}

	return nil
}

// ResetUserCompletedAllCount resets the completed_all_count for a specific user
func (r *StatsRepository) ResetUserCompletedAllCount(userID int) error {
	query := `
		UPDATE user_stats 
		SET completed_all_count = 0, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $1`

	result, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to reset user completed_all_count: %w", err)
	}

	// Check if user stats record exists, if not create one
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	// If no rows were affected, the user doesn't have stats yet, so create them
	if rowsAffected == 0 {
		_, err = r.initializeUserStats(userID)
		if err != nil {
			return fmt.Errorf("failed to initialize user stats: %w", err)
		}
	}

	return nil
}

// IncrementUserCompletedAllCount increments the completed_all_count for a specific user
func (r *StatsRepository) IncrementUserCompletedAllCount(userID int) error {
	query := `
		INSERT INTO user_stats (user_id, completed_all_count, created_at, updated_at)
		VALUES ($1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (user_id) 
		DO UPDATE SET 
			completed_all_count = user_stats.completed_all_count + 1,
			updated_at = CURRENT_TIMESTAMP`

	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to increment user completed_all_count: %w", err)
	}

	return nil
}

// GetUserStats retrieves user-specific statistics
func (r *StatsRepository) GetUserStats(userID int) (*models.UserStats, error) {
	query := `
		SELECT user_id, total_items, completed_items, in_progress_items, pending_items,
			   dsa_completed, lld_completed, hld_completed, completed_all_count,
			   current_streak, longest_streak, last_activity_date, created_at, updated_at
		FROM user_stats 
		WHERE user_id = $1`

	var stats models.UserStats
	err := r.db.QueryRow(query, userID).Scan(
		&stats.UserID, &stats.TotalItems, &stats.CompletedItems, &stats.InProgressItems,
		&stats.PendingItems, &stats.DSACompleted, &stats.LLDCompleted, &stats.HLDCompleted,
		&stats.CompletedAllCount, &stats.CurrentStreak, &stats.LongestStreak,
		&stats.LastActivityDate, &stats.CreatedAt, &stats.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// Initialize user stats if not exists
		return r.initializeUserStats(userID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}

	return &stats, nil
}

// initializeUserStats creates initial user stats record
func (r *StatsRepository) initializeUserStats(userID int) (*models.UserStats, error) {
	query := `
		INSERT INTO user_stats (user_id, created_at, updated_at)
		VALUES ($1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING user_id, total_items, completed_items, in_progress_items, pending_items,
				  dsa_completed, lld_completed, hld_completed, completed_all_count,
				  current_streak, longest_streak, last_activity_date, created_at, updated_at`

	var stats models.UserStats
	err := r.db.QueryRow(query, userID).Scan(
		&stats.UserID, &stats.TotalItems, &stats.CompletedItems, &stats.InProgressItems,
		&stats.PendingItems, &stats.DSACompleted, &stats.LLDCompleted, &stats.HLDCompleted,
		&stats.CompletedAllCount, &stats.CurrentStreak, &stats.LongestStreak,
		&stats.LastActivityDate, &stats.CreatedAt, &stats.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to initialize user stats: %w", err)
	}

	return &stats, nil
}

// initializeAppStats creates the initial app stats record
func (r *StatsRepository) initializeAppStats() (*models.AppStats, error) {
	query := `
		INSERT INTO app_stats (id, completed_all_count) 
		VALUES (1, 0) 
		ON CONFLICT (id) DO NOTHING
		RETURNING id, completed_all_count`

	var stats models.AppStats
	err := r.db.QueryRow(query).Scan(&stats.ID, &stats.CompletedAllCount)
	if err != nil {
		// If the INSERT didn't return anything due to conflict, try to SELECT
		selectQuery := "SELECT id, completed_all_count FROM app_stats WHERE id = 1"
		err = r.db.QueryRow(selectQuery).Scan(&stats.ID, &stats.CompletedAllCount)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize app stats: %w", err)
		}
	}

	return &stats, nil
}
