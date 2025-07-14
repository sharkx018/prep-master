package repositories

import (
	"database/sql"
	"fmt"
	"time"

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

	// Check and reset streak if there's a gap of 24+ hours
	err = r.checkAndResetStreakIfNeeded(&stats)
	if err != nil {
		return nil, fmt.Errorf("failed to check and reset streak: %w", err)
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

// UpdateUserStreakOnActivity updates the user's streak when they complete an item
func (r *StatsRepository) UpdateUserStreakOnActivity(userID int) error {
	// First check if user already has activity today
	hasActivityToday, err := r.HasActivityToday(userID)
	if err != nil {
		return fmt.Errorf("failed to check today's activity: %w", err)
	}

	// If user already completed something today, don't update streak
	if hasActivityToday {
		return nil
	}

	// Get current user stats
	userStats, err := r.GetUserStats(userID)
	if err != nil {
		return fmt.Errorf("failed to get user stats: %w", err)
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)

	// If this is the first activity ever, start streak at 1
	if userStats.LastActivityDate == nil {
		return r.updateUserStreak(userID, 1, 1, today)
	}

	lastActivity := userStats.LastActivityDate.UTC().Truncate(24 * time.Hour)

	// If user completed something yesterday, increment streak
	yesterday := today.Add(-24 * time.Hour)
	if lastActivity.Equal(yesterday) {
		newStreak := userStats.CurrentStreak + 1
		longestStreak := userStats.LongestStreak
		if newStreak > longestStreak {
			longestStreak = newStreak
		}
		return r.updateUserStreak(userID, newStreak, longestStreak, today)
	}

	// If user missed days, reset streak to 1 (since they're completing an item today)
	return r.updateUserStreak(userID, 1, userStats.LongestStreak, today)
}

// updateUserStreak updates the streak fields in the database
func (r *StatsRepository) updateUserStreak(userID int, currentStreak int, longestStreak int, lastActivityDate time.Time) error {
	query := `
		INSERT INTO user_stats (user_id, current_streak, longest_streak, last_activity_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (user_id) 
		DO UPDATE SET 
			current_streak = EXCLUDED.current_streak,
			longest_streak = EXCLUDED.longest_streak,
			last_activity_date = EXCLUDED.last_activity_date,
			updated_at = CURRENT_TIMESTAMP`

	_, err := r.db.Exec(query, userID, currentStreak, longestStreak, lastActivityDate)
	if err != nil {
		return fmt.Errorf("failed to update user streak: %w", err)
	}

	return nil
}

// GetUserStreakInfo returns just the streak information for a user
func (r *StatsRepository) GetUserStreakInfo(userID int) (currentStreak int, longestStreak int, lastActivityDate *time.Time, err error) {
	query := `
		SELECT current_streak, longest_streak, last_activity_date
		FROM user_stats 
		WHERE user_id = $1`

	err = r.db.QueryRow(query, userID).Scan(&currentStreak, &longestStreak, &lastActivityDate)
	if err == sql.ErrNoRows {
		// User doesn't have stats yet, return defaults
		return 0, 0, nil, nil
	}
	if err != nil {
		return 0, 0, nil, fmt.Errorf("failed to get user streak info: %w", err)
	}

	return currentStreak, longestStreak, lastActivityDate, nil
}

// HasActivityToday checks if the user has already completed an item today
func (r *StatsRepository) HasActivityToday(userID int) (bool, error) {
	today := time.Now().UTC().Truncate(24 * time.Hour)

	query := `
		SELECT last_activity_date
		FROM user_stats 
		WHERE user_id = $1`

	var lastActivityDate *time.Time
	err := r.db.QueryRow(query, userID).Scan(&lastActivityDate)
	if err == sql.ErrNoRows {
		// User doesn't have stats yet, so no activity today
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check user activity: %w", err)
	}

	// If no activity date recorded, no activity today
	if lastActivityDate == nil {
		return false, nil
	}

	// Check if last activity was today
	lastActivity := lastActivityDate.UTC().Truncate(24 * time.Hour)
	return lastActivity.Equal(today), nil
}

// checkAndResetStreakIfNeeded checks if the user's streak should be reset to 0 due to inactivity
func (r *StatsRepository) checkAndResetStreakIfNeeded(stats *models.UserStats) error {
	// If no last activity date or current streak is already 0, nothing to check
	if stats.LastActivityDate == nil || stats.CurrentStreak == 0 {
		return nil
	}

	now := time.Now().UTC()
	today := now.Truncate(24 * time.Hour)
	lastActivity := stats.LastActivityDate.UTC().Truncate(24 * time.Hour)

	// Calculate days since last activity
	daysSinceLastActivity := int(today.Sub(lastActivity).Hours() / 24)

	// If there's a gap of 1 or more days, reset streak to 0
	if daysSinceLastActivity >= 1 {
		// Update the streak in the database
		err := r.resetUserStreak(stats.UserID)
		if err != nil {
			return fmt.Errorf("failed to reset user streak: %w", err)
		}

		// Update the stats object to reflect the reset
		stats.CurrentStreak = 0
	}

	return nil
}

// resetUserStreak resets the user's current streak to 0
func (r *StatsRepository) resetUserStreak(userID int) error {
	query := `
		UPDATE user_stats 
		SET current_streak = 0, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $1`

	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to reset user streak: %w", err)
	}

	return nil
}
