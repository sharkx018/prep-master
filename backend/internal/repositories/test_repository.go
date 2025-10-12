package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"interview-prep-app/internal/models"

	"github.com/lib/pq"
)

// TestRepository handles database operations for tests
type TestRepository struct {
	db *sql.DB
}

// NewTestRepository creates a new test repository
func NewTestRepository(db *sql.DB) *TestRepository {
	return &TestRepository{db: db}
}

// CreateTestItems creates multiple test items with the same session ID
func (r *TestRepository) CreateTestItems(userID int, itemIDs []int) (string, error) {
	// Generate a UUID using PostgreSQL's gen_random_uuid() function
	var sessionID string
	err := r.db.QueryRow("SELECT gen_random_uuid()::text").Scan(&sessionID)
	if err != nil {
		return "", fmt.Errorf("failed to generate session ID: %w", err)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO tests (session_id, user_id, item_id, status)
		VALUES ($1, $2, $3, 'pending')`

	for _, itemID := range itemIDs {
		_, err := tx.Exec(query, sessionID, userID, itemID)
		if err != nil {
			return "", fmt.Errorf("failed to create test item: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return sessionID, nil
}

// GetTestByUserWithStatus retrieves a test session for a user filtered by status
func (r *TestRepository) GetTestByUserWithStatus(userID int, itemStatus []string) (string, []int, error) {
	query := `
		SELECT session_id
		FROM tests
		WHERE user_id = $1 AND status = ANY($2)
		ORDER BY created_at DESC
		LIMIT 1`

	var sessionID string
	err := r.db.QueryRow(query, userID, pq.Array(itemStatus)).Scan(&sessionID)
	if err == sql.ErrNoRows {
		return "", nil, nil // No test found
	}
	if err != nil {
		return "", nil, fmt.Errorf("failed to get test: %w", err)
	}

	// Get all item IDs for this session with the specified statuses
	itemQuery := `
		SELECT item_id
		FROM tests
		WHERE user_id = $1 AND session_id = $2 AND status = ANY($3)
		ORDER BY id`

	rows, err := r.db.Query(itemQuery, userID, sessionID, pq.Array(itemStatus))
	if err != nil {
		return "", nil, fmt.Errorf("failed to get test items: %w", err)
	}
	defer rows.Close()

	var itemIDs []int
	for rows.Next() {
		var itemID int
		if err := rows.Scan(&itemID); err != nil {
			return "", nil, fmt.Errorf("failed to scan item ID: %w", err)
		}
		itemIDs = append(itemIDs, itemID)
	}

	if err := rows.Err(); err != nil {
		return "", nil, fmt.Errorf("error iterating test items: %w", err)
	}

	return sessionID, itemIDs, nil
}

// GetTestsBySessionID retrieves all tests for a specific session
func (r *TestRepository) GetTestsBySessionID(userID int, sessionID string) ([]*models.Test, error) {
	query := `
		SELECT id, session_id, user_id, item_id, status, created_at, updated_at
		FROM tests
		WHERE user_id = $1 AND session_id = $2
		ORDER BY id`

	rows, err := r.db.Query(query, userID, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tests by session: %w", err)
	}
	defer rows.Close()

	var tests []*models.Test
	for rows.Next() {
		var test models.Test
		err := rows.Scan(
			&test.ID,
			&test.SessionID,
			&test.UserID,
			&test.ItemID,
			&test.Status,
			&test.CreatedAt,
			&test.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan test: %w", err)
		}
		tests = append(tests, &test)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tests: %w", err)
	}

	return tests, nil
}

// UpdateTestStatus updates the status of all tests in a session
func (r *TestRepository) UpdateTestStatus(userID int, sessionID string, item_id string, status models.TestStatus) error {
	query := `
		UPDATE tests
		SET status = $1, updated_at = $2
		WHERE user_id = $3 AND session_id = $4 AND item_id = $5`

	result, err := r.db.Exec(query, status, time.Now(), userID, sessionID, item_id)
	if err != nil {
		return fmt.Errorf("failed to update test status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no tests found for session")
	}

	return nil
}

// DeleteTestsBySessionID deletes all tests for a specific session
func (r *TestRepository) DeleteTestsBySessionID(userID int, sessionID string) error {
	query := `
		DELETE FROM tests
		WHERE user_id = $1 AND session_id = $2`

	result, err := r.db.Exec(query, userID, sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete tests: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no tests found for session")
	}

	return nil
}

// GetTestCreatedAt retrieves the created_at timestamp for a session
func (r *TestRepository) GetTestCreatedAt(userID int, sessionID string) (time.Time, error) {
	query := `
		SELECT created_at
		FROM tests
		WHERE user_id = $1 AND session_id = $2
		ORDER BY created_at
		LIMIT 1`

	var createdAt time.Time
	err := r.db.QueryRow(query, userID, sessionID).Scan(&createdAt)
	if err == sql.ErrNoRows {
		return time.Time{}, fmt.Errorf("no tests found for session")
	}
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get test created_at: %w", err)
	}

	return createdAt, nil
}

// IsItemInPendingTest checks if an item is part of an pending test for a user
func (r *TestRepository) IsItemInPendingTest(userID int) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM tests
			WHERE user_id = $1 AND status = 'pending'
		)`

	var exists bool
	err := r.db.QueryRow(query, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if item is in pending test: %w", err)
	}

	return exists, nil
}
