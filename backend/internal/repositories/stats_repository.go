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

// ResetCompletedAllCount resets the completed_all_count to zero
func (r *StatsRepository) ResetCompletedAllCount() error {
	query := "UPDATE app_stats SET completed_all_count = 0 WHERE id = 1"

	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to reset completed_all_count: %w", err)
	}

	return nil
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
