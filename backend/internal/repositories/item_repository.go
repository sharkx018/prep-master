package repositories

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"interview-prep-app/internal/models"
)

// ItemRepository handles database operations for items
type ItemRepository struct {
	db *sql.DB
}

// NewItemRepository creates a new item repository
func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

// Create adds a new item to the database
func (r *ItemRepository) Create(req *models.CreateItemRequest) (*models.Item, error) {
	// Initialize attachments if nil
	attachments := req.Attachments
	if attachments == nil {
		attachments = make(models.Attachments)
	}

	query := `
		INSERT INTO items (title, link, category, subcategory, attachments) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, title, link, category, subcategory, attachments, created_at`

	var item models.Item
	err := r.db.QueryRow(query, req.Title, req.Link, req.Category, req.Subcategory, attachments).Scan(
		&item.ID, &item.Title, &item.Link, &item.Category, &item.Subcategory,
		&item.Attachments, &item.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	return &item, nil
}

// GetByID retrieves an item by its ID
func (r *ItemRepository) GetByID(id int) (*models.Item, error) {
	query := `
		SELECT id, title, link, category, subcategory, attachments, created_at 
		FROM items 
		WHERE id = $1`

	var item models.Item
	err := r.db.QueryRow(query, id).Scan(
		&item.ID, &item.Title, &item.Link, &item.Category, &item.Subcategory,
		&item.Attachments, &item.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("item not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	return &item, nil
}

// GetByIDWithUserProgress retrieves an item by its ID with user-specific progress data
func (r *ItemRepository) GetByIDWithUserProgress(userID, itemID int) (*models.ItemWithProgress, error) {
	query := `
		SELECT 
			i.id, i.title, i.link, i.category, i.subcategory, i.attachments, i.created_at,
			COALESCE(up.status, 'pending') as status,
			COALESCE(up.starred, false) as starred,
			COALESCE(up.notes, '') as notes,
			up.completed_at
		FROM items i
		LEFT JOIN user_progress up 
			ON i.id = up.item_id AND up.user_id = $1
		WHERE i.id = $2`

	var item models.ItemWithProgress
	err := r.db.QueryRow(query, userID, itemID).Scan(
		&item.ID, &item.Title, &item.Link, &item.Category, &item.Subcategory,
		&item.Attachments, &item.CreatedAt, &item.Status, &item.Starred,
		&item.Notes, &item.CompletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("item not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get item with user progress: %w", err)
	}

	return &item, nil
}

// GetAll retrieves items with optional filtering
func (r *ItemRepository) GetAll(filter *models.ItemFilter) ([]*models.Item, error) {
	query := "SELECT id, title, link, category, subcategory, attachments, created_at FROM items WHERE 1=1"
	args := []interface{}{}
	argCount := 0

	// Build dynamic query based on filters
	if filter.Category != nil {
		argCount++
		query += fmt.Sprintf(" AND category = $%d", argCount)
		args = append(args, *filter.Category)
	}

	if filter.Subcategory != nil {
		argCount++
		query += fmt.Sprintf(" AND subcategory = $%d", argCount)
		args = append(args, *filter.Subcategory)
	}

	// Note: Status filtering is no longer supported in this method
	// Use GetAllWithUserProgress for user-specific status filtering

	query += " ORDER BY created_at DESC"

	if filter.Limit != nil {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, *filter.Limit)

		if filter.Offset != nil {
			argCount++
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, *filter.Offset)
		}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}
	defer rows.Close()

	var items []*models.Item
	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&item.ID, &item.Title, &item.Link, &item.Category, &item.Subcategory,
			&item.Attachments, &item.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		items = append(items, &item)
	}

	return items, nil
}

// GetAllWithUserProgress retrieves items with user-specific progress data using LEFT JOIN
func (r *ItemRepository) GetAllWithUserProgress(userID int, filter *models.ItemFilter) ([]*models.ItemWithProgress, error) {
	query := `
		SELECT 
			i.id, i.title, i.link, i.category, i.subcategory, i.attachments, i.created_at,
			COALESCE(up.status, 'pending') as status,
			COALESCE(up.starred, false) as starred,
			COALESCE(up.notes, '') as notes,
			up.completed_at
		FROM items i
		LEFT JOIN user_progress up 
			ON i.id = up.item_id AND up.user_id = $1
		WHERE 1=1`

	args := []interface{}{userID}
	argCount := 1

	// Build dynamic query based on filters
	if filter.Category != nil {
		argCount++
		query += fmt.Sprintf(" AND i.category = $%d", argCount)
		args = append(args, *filter.Category)
	}

	if filter.Subcategory != nil {
		argCount++
		query += fmt.Sprintf(" AND i.subcategory = $%d", argCount)
		args = append(args, *filter.Subcategory)
	}

	if filter.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND COALESCE(up.status, 'pending') = $%d", argCount)
		args = append(args, *filter.Status)
	}

	query += " ORDER BY i.created_at DESC"

	if filter.Limit != nil {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, *filter.Limit)

		if filter.Offset != nil {
			argCount++
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, *filter.Offset)
		}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get items with user progress: %w", err)
	}
	defer rows.Close()

	var items []*models.ItemWithProgress
	for rows.Next() {
		var item models.ItemWithProgress
		err := rows.Scan(
			&item.ID, &item.Title, &item.Link, &item.Category, &item.Subcategory,
			&item.Attachments, &item.CreatedAt, &item.Status, &item.Starred,
			&item.Notes, &item.CompletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item with progress: %w", err)
		}
		items = append(items, &item)
	}

	return items, nil
}

// GetRandomPending is deprecated - use GetRandomPendingWithUserProgress instead
func (r *ItemRepository) GetRandomPending() (*models.Item, error) {
	return nil, fmt.Errorf("GetRandomPending is deprecated - use GetRandomPendingWithUserProgress instead")
}

// GetInProgressItem is deprecated - use GetInProgressItemWithUserProgress instead
func (r *ItemRepository) GetInProgressItem() (*models.Item, error) {
	return nil, fmt.Errorf("GetInProgressItem is deprecated - use GetInProgressItemWithUserProgress instead")
}

// SetInProgress is deprecated - use UpsertUserProgressForItem instead
func (r *ItemRepository) SetInProgress(id int) (*models.Item, error) {
	return nil, fmt.Errorf("SetInProgress is deprecated - use UpsertUserProgressForItem instead")
}

// MarkComplete is deprecated - use CompleteItemForUser instead
func (r *ItemRepository) MarkComplete(id int) (*models.Item, error) {
	return nil, fmt.Errorf("MarkComplete is deprecated - use CompleteItemForUser instead")
}

// Update updates an existing item
func (r *ItemRepository) Update(id int, req *models.UpdateItemRequest) (*models.Item, error) {
	setParts := []string{}
	args := []interface{}{}
	argCount := 0

	if req.Title != nil {
		argCount++
		setParts = append(setParts, fmt.Sprintf("title = $%d", argCount))
		args = append(args, *req.Title)
	}

	if req.Link != nil {
		argCount++
		setParts = append(setParts, fmt.Sprintf("link = $%d", argCount))
		args = append(args, *req.Link)
	}

	if req.Category != nil {
		argCount++
		setParts = append(setParts, fmt.Sprintf("category = $%d", argCount))
		args = append(args, *req.Category)
	}

	if req.Subcategory != nil {
		argCount++
		setParts = append(setParts, fmt.Sprintf("subcategory = $%d", argCount))
		args = append(args, *req.Subcategory)
	}

	if req.Attachments != nil {
		argCount++
		setParts = append(setParts, fmt.Sprintf("attachments = $%d", argCount))
		args = append(args, *req.Attachments)
	}

	if len(setParts) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	argCount++
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE items 
		SET %s 
		WHERE id = $%d
		RETURNING id, title, link, category, subcategory, attachments, created_at`,
		strings.Join(setParts, ", "), argCount)

	var item models.Item
	err := r.db.QueryRow(query, args...).Scan(
		&item.ID, &item.Title, &item.Link, &item.Category, &item.Subcategory,
		&item.Attachments, &item.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("item not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	return &item, nil
}

// Delete removes an item from the database and cascades to user_progress
func (r *ItemRepository) Delete(id int) error {
	// Start a transaction to ensure atomicity
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// First, check if the item exists
	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM items WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if item exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("item not found")
	}

	// Delete user progress entries for this item (optional since CASCADE will handle this)
	// This is explicit for clarity and potential logging
	_, err = tx.Exec("DELETE FROM user_progress WHERE item_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete user progress entries: %w", err)
	}

	// Delete the item (this would also cascade delete user_progress due to FK constraint)
	result, err := tx.Exec("DELETE FROM items WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("item not found")
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ResetAll is deprecated - use ResetAllUserProgress instead
func (r *ItemRepository) ResetAll() (int64, error) {
	return 0, fmt.Errorf("ResetAll is deprecated - use ResetAllUserProgress instead")
}

// GetCounts is deprecated - use GetCountsForUser instead
func (r *ItemRepository) GetCounts() (total, completed, pending int, err error) {
	return 0, 0, 0, fmt.Errorf("GetCounts is deprecated - use GetCountsForUser instead")
}

// GetCountsByCategory is deprecated - use GetCountsByCategoryForUser instead
func (r *ItemRepository) GetCountsByCategory() (map[models.Category]map[models.Status]int, error) {
	return nil, fmt.Errorf("GetCountsByCategory is deprecated - use GetCountsByCategoryForUser instead")
}

// GetCountsBySubcategory is deprecated - use GetCountsBySubcategoryForUser instead
func (r *ItemRepository) GetCountsBySubcategory() (map[models.Category]map[string]map[models.Status]int, error) {
	return nil, fmt.Errorf("GetCountsBySubcategory is deprecated - use GetCountsBySubcategoryForUser instead")
}

// CountPending is deprecated - use CountPendingForUser instead
func (r *ItemRepository) CountPending() (int, error) {
	return 0, fmt.Errorf("CountPending is deprecated - use CountPendingForUser instead")
}

// ToggleStar is deprecated - use ToggleStarForUser instead
func (r *ItemRepository) ToggleStar(id int) (*models.Item, error) {
	return nil, fmt.Errorf("ToggleStar is deprecated - use ToggleStarForUser instead")
}

// UpdateStatus is deprecated - use UpdateStatusForUser instead
func (r *ItemRepository) UpdateStatus(id int, status models.Status) (*models.Item, error) {
	return nil, fmt.Errorf("UpdateStatus is deprecated - use UpdateStatusForUser instead")
}

// GetTotalCount returns the total count of items matching the filter
func (r *ItemRepository) GetTotalCount(filter *models.ItemFilter) (int, error) {
	query := "SELECT COUNT(*) FROM items WHERE 1=1"
	args := []interface{}{}
	argCount := 0

	// Build dynamic query based on filters (same logic as GetAll)
	if filter.Category != nil {
		argCount++
		query += fmt.Sprintf(" AND category = $%d", argCount)
		args = append(args, *filter.Category)
	}

	if filter.Subcategory != nil {
		argCount++
		query += fmt.Sprintf(" AND subcategory = $%d", argCount)
		args = append(args, *filter.Subcategory)
	}

	// Note: Status filtering is no longer supported in this method
	// Use GetTotalCountWithUserProgress for user-specific status filtering

	var count int
	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count items: %w", err)
	}
	return count, nil
}

// GetTotalCountWithUserProgress returns the total count of items matching the filter with user-specific progress
func (r *ItemRepository) GetTotalCountWithUserProgress(userID int, filter *models.ItemFilter) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM items i
		LEFT JOIN user_progress up ON i.id = up.item_id AND up.user_id = $1
		WHERE 1=1`

	args := []interface{}{userID}
	argCount := 1

	// Build dynamic query based on filters
	if filter.Category != nil {
		argCount++
		query += fmt.Sprintf(" AND i.category = $%d", argCount)
		args = append(args, *filter.Category)
	}

	if filter.Subcategory != nil {
		argCount++
		query += fmt.Sprintf(" AND i.subcategory = $%d", argCount)
		args = append(args, *filter.Subcategory)
	}

	if filter.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND COALESCE(up.status, 'pending') = $%d", argCount)
		args = append(args, *filter.Status)
	}

	var count int
	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count items with user progress: %w", err)
	}
	return count, nil
}

// GetInProgressItemWithUserProgress retrieves the current in-progress item for a user
func (r *ItemRepository) GetInProgressItemWithUserProgress(userID int) (*models.ItemWithProgress, error) {
	query := `
		SELECT 
			i.id, i.title, i.link, i.category, i.subcategory, i.attachments, i.created_at,
			up.status, up.starred, up.notes, up.completed_at
		FROM items i
		INNER JOIN user_progress up ON i.id = up.item_id AND up.user_id = $1
		WHERE up.status = 'in-progress'
		LIMIT 1`

	var item models.ItemWithProgress
	err := r.db.QueryRow(query, userID).Scan(
		&item.ID, &item.Title, &item.Link, &item.Category, &item.Subcategory,
		&item.Attachments, &item.CreatedAt, &item.Status, &item.Starred,
		&item.Notes, &item.CompletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // No in-progress item
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get in-progress item with user progress: %w", err)
	}

	return &item, nil
}

// GetRandomPendingWithUserProgress retrieves a random pending item for a user
func (r *ItemRepository) GetRandomPendingWithUserProgress(userID int) (*models.ItemWithProgress, error) {
	query := `
		SELECT 
			i.id, i.title, i.link, i.category, i.subcategory, i.attachments, i.created_at,
			COALESCE(up.status, 'pending') as status,
			COALESCE(up.starred, false) as starred,
			COALESCE(up.notes, '') as notes,
			up.completed_at
		FROM items i
		LEFT JOIN user_progress up 
			ON i.id = up.item_id AND up.user_id = $1
		WHERE COALESCE(up.status, 'pending') = 'pending'
		ORDER BY RANDOM() 
		LIMIT 1`

	var item models.ItemWithProgress
	err := r.db.QueryRow(query, userID).Scan(
		&item.ID, &item.Title, &item.Link, &item.Category, &item.Subcategory,
		&item.Attachments, &item.CreatedAt, &item.Status, &item.Starred,
		&item.Notes, &item.CompletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no pending items found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get random pending item with user progress: %w", err)
	}

	return &item, nil
}

// CreateUserProgressForItem creates or updates a user progress record for an item
func (r *ItemRepository) CreateUserProgressForItem(userID, itemID int, status models.Status) error {
	now := time.Now()

	query := `
		INSERT INTO user_progress (user_id, item_id, status, starred, notes, started_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (user_id, item_id) 
		DO UPDATE SET 
			status = EXCLUDED.status,
			started_at = CASE 
				WHEN EXCLUDED.status = 'in-progress' THEN EXCLUDED.started_at
				ELSE user_progress.started_at
			END,
			updated_at = EXCLUDED.updated_at`

	_, err := r.db.Exec(
		query,
		userID,
		itemID,
		status,
		false, // starred defaults to false for new records
		"",    // notes defaults to empty for new records
		now,   // started_at
		now,   // created_at
		now,   // updated_at
	)

	if err != nil {
		return fmt.Errorf("failed to create/update user progress for item: %w", err)
	}

	return nil
}

// UpsertUserProgressForItem creates or updates a user progress record preserving existing data
func (r *ItemRepository) UpsertUserProgressForItem(userID, itemID int, status models.Status) error {
	now := time.Now()

	query := `
		INSERT INTO user_progress (user_id, item_id, status, starred, notes, started_at, created_at, updated_at)
		VALUES ($1, $2, $3, false, '', $4, $5, $6)
		ON CONFLICT (user_id, item_id) 
		DO UPDATE SET 
			status = EXCLUDED.status,
			started_at = CASE 
				WHEN EXCLUDED.status = 'in-progress' AND user_progress.status != 'in-progress' THEN EXCLUDED.started_at
				ELSE user_progress.started_at
			END,
			completed_at = CASE 
				WHEN EXCLUDED.status = 'done' THEN $7
				WHEN EXCLUDED.status != 'done' THEN NULL
				ELSE user_progress.completed_at
			END,
			updated_at = EXCLUDED.updated_at`

	_, err := r.db.Exec(
		query,
		userID,
		itemID,
		status,
		now, // started_at for new records
		now, // created_at
		now, // updated_at
		now, // completed_at for done status
	)

	if err != nil {
		return fmt.Errorf("failed to upsert user progress for item: %w", err)
	}

	return nil
}

// ResetInProgressItemsForUser resets any in-progress items for a user back to pending
func (r *ItemRepository) ResetInProgressItemsForUser(userID int) error {
	query := `
		UPDATE user_progress 
		SET status = 'pending', updated_at = $1
		WHERE user_id = $2 AND status = 'in-progress'`

	_, err := r.db.Exec(query, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to reset in-progress items for user: %w", err)
	}

	return nil
}

// CountPendingForUser counts pending items for a specific user
func (r *ItemRepository) CountPendingForUser(userID int) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM items i
		LEFT JOIN user_progress up 
			ON i.id = up.item_id AND up.user_id = $1
		WHERE COALESCE(up.status, 'pending') = 'pending'`

	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count pending items for user: %w", err)
	}
	return count, nil
}

// CompleteItemForUser marks an item as completed for a specific user
func (r *ItemRepository) CompleteItemForUser(userID, itemID int) (*models.ItemWithProgress, error) {
	// First, ensure the item exists
	var itemExists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM items WHERE id = $1)", itemID).Scan(&itemExists)
	if err != nil {
		return nil, fmt.Errorf("failed to check if item exists: %w", err)
	}
	if !itemExists {
		return nil, fmt.Errorf("item not found")
	}

	// Update or insert user progress to mark as completed
	err = r.UpsertUserProgressForItem(userID, itemID, models.StatusDone)
	if err != nil {
		return nil, fmt.Errorf("failed to mark item as completed: %w", err)
	}

	// Get the completed item with user progress
	item, err := r.GetByIDWithUserProgress(userID, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get completed item: %w", err)
	}

	return item, nil
}

// ToggleStarForUser toggles the starred status of an item for a specific user
func (r *ItemRepository) ToggleStarForUser(userID, itemID int) (*models.ItemWithProgress, error) {
	// First, ensure the item exists
	var itemExists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM items WHERE id = $1)", itemID).Scan(&itemExists)
	if err != nil {
		return nil, fmt.Errorf("failed to check if item exists: %w", err)
	}
	if !itemExists {
		return nil, fmt.Errorf("item not found")
	}

	// Get current starred status or default to false
	var currentStarred bool
	query := `
		SELECT COALESCE(up.starred, false)
		FROM items i
		LEFT JOIN user_progress up 
			ON i.id = up.item_id AND up.user_id = $1
		WHERE i.id = $2`

	err = r.db.QueryRow(query, userID, itemID).Scan(&currentStarred)
	if err != nil {
		return nil, fmt.Errorf("failed to get current starred status: %w", err)
	}

	// Toggle the starred status
	newStarred := !currentStarred

	// Upsert user progress with the new starred status
	now := time.Now()
	upsertQuery := `
		INSERT INTO user_progress (user_id, item_id, status, starred, notes, created_at, updated_at)
		VALUES ($1, $2, 'pending', $3, '', $4, $5)
		ON CONFLICT (user_id, item_id) 
		DO UPDATE SET 
			starred = EXCLUDED.starred,
			updated_at = EXCLUDED.updated_at`

	_, err = r.db.Exec(upsertQuery, userID, itemID, newStarred, now, now)
	if err != nil {
		return nil, fmt.Errorf("failed to toggle star status: %w", err)
	}

	// Get the updated item with user progress
	item, err := r.GetByIDWithUserProgress(userID, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated item: %w", err)
	}

	return item, nil
}

// UpdateStatusForUser updates the status of an item for a specific user
func (r *ItemRepository) UpdateStatusForUser(userID, itemID int, status models.Status) (*models.ItemWithProgress, error) {
	// First, ensure the item exists
	var itemExists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM items WHERE id = $1)", itemID).Scan(&itemExists)
	if err != nil {
		return nil, fmt.Errorf("failed to check if item exists: %w", err)
	}
	if !itemExists {
		return nil, fmt.Errorf("item not found")
	}

	// Use the UpsertUserProgressForItem method to update status
	err = r.UpsertUserProgressForItem(userID, itemID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to update status: %w", err)
	}

	// Get the updated item with user progress
	item, err := r.GetByIDWithUserProgress(userID, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated item: %w", err)
	}

	return item, nil
}

// ResetAllUserProgress resets all user progress for a specific user back to pending
func (r *ItemRepository) ResetAllUserProgress(userID int) (int64, error) {
	query := `
		UPDATE user_progress 
		SET status = 'pending', completed_at = NULL, updated_at = $1
		WHERE user_id = $2 AND status IN ('done', 'in-progress')`

	result, err := r.db.Exec(query, time.Now(), userID)
	if err != nil {
		return 0, fmt.Errorf("failed to reset user progress: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

// GetCountsForUser returns item counts by status for a specific user
func (r *ItemRepository) GetCountsForUser(userID int) (total, completed, pending, inProgress int, err error) {
	query := `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN COALESCE(up.status, 'pending') = 'done' THEN 1 END) as completed,
			COUNT(CASE WHEN COALESCE(up.status, 'pending') = 'pending' THEN 1 END) as pending,
			COUNT(CASE WHEN COALESCE(up.status, 'pending') = 'in-progress' THEN 1 END) as in_progress
		FROM items i
		LEFT JOIN user_progress up 
			ON i.id = up.item_id AND up.user_id = $1`

	err = r.db.QueryRow(query, userID).Scan(&total, &completed, &pending, &inProgress)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("failed to get user counts: %w", err)
	}

	return total, completed, pending, inProgress, nil
}

// GetCountsByCategoryForUser returns item counts by category and status for a specific user
func (r *ItemRepository) GetCountsByCategoryForUser(userID int) (map[models.Category]map[models.Status]int, error) {
	query := `
		SELECT 
			i.category,
			COALESCE(up.status, 'pending') as status,
			COUNT(*) as count
		FROM items i
		LEFT JOIN user_progress up 
			ON i.id = up.item_id AND up.user_id = $1
		GROUP BY i.category, COALESCE(up.status, 'pending')
		ORDER BY i.category, status`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user category counts: %w", err)
	}
	defer rows.Close()

	result := make(map[models.Category]map[models.Status]int)

	for rows.Next() {
		var category models.Category
		var status models.Status
		var count int

		err := rows.Scan(&category, &status, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category count: %w", err)
		}

		if result[category] == nil {
			result[category] = make(map[models.Status]int)
		}
		result[category][status] = count
	}

	return result, nil
}

// GetCountsBySubcategoryForUser returns item counts by subcategory and status for a specific user
func (r *ItemRepository) GetCountsBySubcategoryForUser(userID int) (map[models.Category]map[string]map[models.Status]int, error) {
	query := `
		SELECT 
			i.category,
			i.subcategory,
			COALESCE(up.status, 'pending') as status,
			COUNT(*) as count
		FROM items i
		LEFT JOIN user_progress up 
			ON i.id = up.item_id AND up.user_id = $1
		GROUP BY i.category, i.subcategory, COALESCE(up.status, 'pending')
		ORDER BY i.category, i.subcategory, status`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user subcategory counts: %w", err)
	}
	defer rows.Close()

	result := make(map[models.Category]map[string]map[models.Status]int)

	for rows.Next() {
		var category models.Category
		var subcategory string
		var status models.Status
		var count int

		err := rows.Scan(&category, &subcategory, &status, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subcategory count: %w", err)
		}

		if result[category] == nil {
			result[category] = make(map[string]map[models.Status]int)
		}
		if result[category][subcategory] == nil {
			result[category][subcategory] = make(map[models.Status]int)
		}
		result[category][subcategory][status] = count
	}

	return result, nil
}
