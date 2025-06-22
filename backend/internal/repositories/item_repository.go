package repositories

import (
	"database/sql"
	"fmt"
	"strings"

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
	}

	if filter.Offset != nil {
		argCount++
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, *filter.Offset)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query items: %w", err)
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

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
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
