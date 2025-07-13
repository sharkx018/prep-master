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
		RETURNING id, title, link, category, subcategory, status, starred, attachments, created_at, completed_at`

	var item models.Item
	err := r.db.QueryRow(query, req.Title, req.Link, req.Category, req.Subcategory, attachments).Scan(
		&item.ID, &item.Title, &item.Link, &item.Category, &item.Subcategory,
		&item.Status, &item.Starred, &item.Attachments, &item.CreatedAt, &item.CompletedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	return &item, nil
}

// GetByID retrieves an item by its ID
func (r *ItemRepository) GetByID(id int) (*models.Item, error) {
	query := `
		SELECT id, title, link, category, subcategory, status, starred, attachments, created_at, completed_at 
		FROM items 
		WHERE id = $1`

	var item models.Item
	err := r.db.QueryRow(query, id).Scan(
		&item.ID, &item.Title, &item.Link, &item.Category, &item.Subcategory,
		&item.Status, &item.Starred, &item.Attachments, &item.CreatedAt, &item.CompletedAt,
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
	query := "SELECT id, title, link, category, subcategory, status, starred, attachments, created_at, completed_at FROM items WHERE 1=1"
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

	if filter.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filter.Status)
	}

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
			&item.Status, &item.Starred, &item.Attachments, &item.CreatedAt, &item.CompletedAt,
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

// GetRandomPending retrieves a random item with status 'pending'
func (r *ItemRepository) GetRandomPending() (*models.Item, error) {
	query := `
		SELECT id, title, link, category, subcategory, status, starred, attachments, created_at, completed_at 
		FROM items 
		WHERE status = 'pending' 
		ORDER BY RANDOM() 
		LIMIT 1`

	var item models.Item
	err := r.db.QueryRow(query).Scan(
		&item.ID, &item.Title, &item.Link, &item.Category, &item.Subcategory,
		&item.Status, &item.Starred, &item.Attachments, &item.CreatedAt, &item.CompletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no pending items found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get random pending item: %w", err)
	}

	return &item, nil
}

// GetInProgressItem retrieves the current in-progress item if any
func (r *ItemRepository) GetInProgressItem() (*models.Item, error) {
	query := `
		SELECT id, title, link, category, subcategory, status, starred, attachments, created_at, completed_at 
		FROM items 
		WHERE status = 'in-progress' 
		LIMIT 1`

	var item models.Item
	err := r.db.QueryRow(query).Scan(
		&item.ID, &item.Title, &item.Link, &item.Category, &item.Subcategory,
		&item.Status, &item.Starred, &item.Attachments, &item.CreatedAt, &item.CompletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // No in-progress item
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get in-progress item: %w", err)
	}

	return &item, nil
}

// SetInProgress sets an item as in-progress and ensures only one item can be in-progress
func (r *ItemRepository) SetInProgress(id int) (*models.Item, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// First, set any existing in-progress items back to pending
	_, err = tx.Exec(`UPDATE items SET status = 'pending' WHERE status = 'in-progress'`)
	if err != nil {
		return nil, fmt.Errorf("failed to reset in-progress items: %w", err)
	}

	// Then set the specified item as in-progress
	query := `
		UPDATE items 
		SET status = 'in-progress' 
		WHERE id = $1 AND status = 'pending'
		RETURNING id, title, link, category, subcategory, status, starred, attachments, created_at, completed_at`

	var item models.Item
	err = tx.QueryRow(query, id).Scan(
		&item.ID, &item.Title, &item.Link, &item.Category, &item.Subcategory,
		&item.Status, &item.Starred, &item.Attachments, &item.CreatedAt, &item.CompletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("item not found or not pending")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to set item in-progress: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &item, nil
}

// MarkComplete marks an item as completed
func (r *ItemRepository) MarkComplete(id int) (*models.Item, error) {
	query := `
		UPDATE items 
		SET status = 'done', completed_at = CURRENT_TIMESTAMP 
		WHERE id = $1 AND status = 'in-progress'
		RETURNING id, title, link, category, subcategory, status, starred, attachments, created_at, completed_at`

	var item models.Item
	err := r.db.QueryRow(query, id).Scan(
		&item.ID, &item.Title, &item.Link, &item.Category, &item.Subcategory,
		&item.Status, &item.Starred, &item.Attachments, &item.CreatedAt, &item.CompletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("item not found or not in-progress")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to mark item complete: %w", err)
	}

	return &item, nil
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
		RETURNING id, title, link, category, subcategory, status, starred, attachments, created_at, completed_at`,
		strings.Join(setParts, ", "), argCount)

	var item models.Item
	err := r.db.QueryRow(query, args...).Scan(
		&item.ID, &item.Title, &item.Link, &item.Category, &item.Subcategory,
		&item.Status, &item.Starred, &item.Attachments, &item.CreatedAt, &item.CompletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("item not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	return &item, nil
}

// Delete removes an item from the database
func (r *ItemRepository) Delete(id int) error {
	query := "DELETE FROM items WHERE id = $1"
	result, err := r.db.Exec(query, id)
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

	return nil
}

// ResetAll marks all items as pending
func (r *ItemRepository) ResetAll() (int64, error) {
	query := "UPDATE items SET status = 'pending', completed_at = NULL WHERE status IN ('done', 'in-progress')"
	result, err := r.db.Exec(query)
	if err != nil {
		return 0, fmt.Errorf("failed to reset items: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

// GetCounts returns item counts by status
func (r *ItemRepository) GetCounts() (total, completed, pending int, err error) {
	query := `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN status = 'done' THEN 1 END) as completed,
			COUNT(CASE WHEN status IN ('pending', 'in-progress') THEN 1 END) as pending
		FROM items`

	err = r.db.QueryRow(query).Scan(&total, &completed, &pending)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to get item counts: %w", err)
	}

	return total, completed, pending, nil
}

// GetCountsByCategory returns item counts by category and status
func (r *ItemRepository) GetCountsByCategory() (map[models.Category]map[models.Status]int, error) {
	query := `
		SELECT category, status, COUNT(*) 
		FROM items 
		GROUP BY category, status`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get counts by category: %w", err)
	}
	defer rows.Close()

	counts := make(map[models.Category]map[models.Status]int)

	for rows.Next() {
		var category models.Category
		var status models.Status
		var count int

		err := rows.Scan(&category, &status, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category counts: %w", err)
		}

		if counts[category] == nil {
			counts[category] = make(map[models.Status]int)
		}
		counts[category][status] = count
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating category counts: %w", err)
	}

	return counts, nil
}

// GetCountsBySubcategory returns item counts by category, subcategory and status
func (r *ItemRepository) GetCountsBySubcategory() (map[models.Category]map[string]map[models.Status]int, error) {
	query := `
		SELECT category, subcategory, status, COUNT(*) 
		FROM items 
		GROUP BY category, subcategory, status`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get counts by subcategory: %w", err)
	}
	defer rows.Close()

	counts := make(map[models.Category]map[string]map[models.Status]int)

	for rows.Next() {
		var category models.Category
		var subcategory string
		var status models.Status
		var count int

		err := rows.Scan(&category, &subcategory, &status, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subcategory counts: %w", err)
		}

		if counts[category] == nil {
			counts[category] = make(map[string]map[models.Status]int)
		}
		if counts[category][subcategory] == nil {
			counts[category][subcategory] = make(map[models.Status]int)
		}
		counts[category][subcategory][status] = count
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating subcategory counts: %w", err)
	}

	return counts, nil
}

// CountPending returns the number of pending items
func (r *ItemRepository) CountPending() (int, error) {
	query := "SELECT COUNT(*) FROM items WHERE status IN ('pending', 'in-progress')"
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count pending items: %w", err)
	}
	return count, nil
}

// ToggleStar toggles the starred status of an item
func (r *ItemRepository) ToggleStar(id int) (*models.Item, error) {
	query := `
		UPDATE items 
		SET starred = NOT starred 
		WHERE id = $1
		RETURNING id, title, link, category, subcategory, status, starred, attachments, created_at, completed_at`

	var item models.Item
	err := r.db.QueryRow(query, id).Scan(
		&item.ID, &item.Title, &item.Link, &item.Category, &item.Subcategory,
		&item.Status, &item.Starred, &item.Attachments, &item.CreatedAt, &item.CompletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("item not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to toggle star: %w", err)
	}

	return &item, nil
}

// UpdateStatus updates the status of an item
func (r *ItemRepository) UpdateStatus(id int, status models.Status) (*models.Item, error) {
	var query string
	var args []interface{}

	if status == models.StatusDone {
		query = `
			UPDATE items 
			SET status = $2, completed_at = CURRENT_TIMESTAMP 
			WHERE id = $1
			RETURNING id, title, link, category, subcategory, status, starred, attachments, created_at, completed_at`
	} else {
		query = `
			UPDATE items 
			SET status = $2, completed_at = NULL 
			WHERE id = $1
			RETURNING id, title, link, category, subcategory, status, starred, attachments, created_at, completed_at`
	}

	args = []interface{}{id, status}

	var item models.Item
	err := r.db.QueryRow(query, args...).Scan(
		&item.ID, &item.Title, &item.Link, &item.Category, &item.Subcategory,
		&item.Status, &item.Starred, &item.Attachments, &item.CreatedAt, &item.CompletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("item not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update status: %w", err)
	}

	return &item, nil
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

	if filter.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filter.Status)
	}

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
