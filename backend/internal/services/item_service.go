package services

import (
	"fmt"

	"interview-prep-app/internal/models"
	"interview-prep-app/internal/repositories"
)

// ItemService handles business logic for items
type ItemService struct {
	itemRepo  *repositories.ItemRepository
	statsRepo *repositories.StatsRepository
}

// NewItemService creates a new item service
func NewItemService(itemRepo *repositories.ItemRepository, statsRepo *repositories.StatsRepository) *ItemService {
	return &ItemService{
		itemRepo:  itemRepo,
		statsRepo: statsRepo,
	}
}

// CreateItem creates a new item with validation
func (s *ItemService) CreateItem(req *models.CreateItemRequest) (*models.Item, error) {
	// Validate category
	if !models.IsValidCategory(req.Category) {
		return nil, fmt.Errorf("invalid category: %s. Valid categories are: %v", req.Category, models.ValidCategories())
	}

	// Validate required fields
	if req.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if req.Link == "" {
		return nil, fmt.Errorf("link is required")
	}
	if req.Subcategory == "" {
		return nil, fmt.Errorf("subcategory is required")
	}

	return s.itemRepo.Create(req)
}

// GetItem retrieves an item by ID
func (s *ItemService) GetItem(id int) (*models.Item, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid item ID")
	}

	return s.itemRepo.GetByID(id)
}

// GetItemWithUserProgress retrieves an item by ID with user-specific progress data
func (s *ItemService) GetItemWithUserProgress(userID, itemID int) (*models.ItemWithProgress, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	if itemID <= 0 {
		return nil, fmt.Errorf("invalid item ID")
	}

	return s.itemRepo.GetByIDWithUserProgress(userID, itemID)
}

// GetItems retrieves items with filtering and validation
func (s *ItemService) GetItems(filter *models.ItemFilter) ([]*models.Item, error) {
	// Validate filter parameters
	if filter.Category != nil && !models.IsValidCategory(*filter.Category) {
		return nil, fmt.Errorf("invalid category: %s", *filter.Category)
	}

	if filter.Status != nil && !models.IsValidStatus(*filter.Status) {
		return nil, fmt.Errorf("invalid status: %s", *filter.Status)
	}

	if filter.Limit != nil && *filter.Limit < 0 {
		return nil, fmt.Errorf("limit cannot be negative")
	}

	if filter.Offset != nil && *filter.Offset < 0 {
		return nil, fmt.Errorf("offset cannot be negative")
	}

	return s.itemRepo.GetAll(filter)
}

// GetItemsWithUserProgress retrieves items with user-specific progress data
func (s *ItemService) GetItemsWithUserProgress(userID int, filter *models.ItemFilter) ([]*models.ItemWithProgress, error) {
	// Validate filter parameters
	if filter.Category != nil && !models.IsValidCategory(*filter.Category) {
		return nil, fmt.Errorf("invalid category: %s", *filter.Category)
	}

	if filter.Status != nil && !models.IsValidStatus(*filter.Status) {
		return nil, fmt.Errorf("invalid status: %s", *filter.Status)
	}

	if filter.Limit != nil && *filter.Limit < 0 {
		return nil, fmt.Errorf("limit cannot be negative")
	}

	if filter.Offset != nil && *filter.Offset < 0 {
		return nil, fmt.Errorf("offset cannot be negative")
	}

	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	return s.itemRepo.GetAllWithUserProgress(userID, filter)
}

// GetItemsPaginated retrieves items with filtering, validation and pagination metadata
func (s *ItemService) GetItemsPaginated(filter *models.ItemFilter) (*models.PaginatedItemsResponse, error) {
	// Validate filter parameters
	if filter.Category != nil && !models.IsValidCategory(*filter.Category) {
		return nil, fmt.Errorf("invalid category: %s", *filter.Category)
	}

	if filter.Status != nil && !models.IsValidStatus(*filter.Status) {
		return nil, fmt.Errorf("invalid status: %s", *filter.Status)
	}

	if filter.Limit != nil && *filter.Limit < 0 {
		return nil, fmt.Errorf("limit cannot be negative")
	}

	if filter.Offset != nil && *filter.Offset < 0 {
		return nil, fmt.Errorf("offset cannot be negative")
	}

	// Set default limit if not provided
	limit := 10
	if filter.Limit != nil {
		limit = *filter.Limit
	}

	offset := 0
	if filter.Offset != nil {
		offset = *filter.Offset
	}

	// Get total count
	totalCount, err := s.itemRepo.GetTotalCount(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// Get items
	items, err := s.itemRepo.GetAll(filter)
	if err != nil {
		return nil, err
	}

	// Convert []*models.Item to []*models.ItemWithProgress for backward compatibility
	itemsWithProgress := make([]*models.ItemWithProgress, len(items))
	for i, item := range items {
		itemsWithProgress[i] = &models.ItemWithProgress{
			ID:          item.ID,
			Title:       item.Title,
			Link:        item.Link,
			Category:    item.Category,
			Subcategory: item.Subcategory,
			Status:      models.StatusPending, // Default status for non-user-specific queries
			Starred:     false,                // Default starred for non-user-specific queries
			Attachments: item.Attachments,
			CreatedAt:   item.CreatedAt,
			CompletedAt: nil, // Default completed_at for non-user-specific queries
			Notes:       "",  // Default empty notes for non-user-specific queries
		}
	}

	// Calculate pagination metadata
	totalPages := (totalCount + limit - 1) / limit // Ceiling division
	page := (offset / limit) + 1
	hasNext := offset+limit < totalCount
	hasPrev := offset > 0

	return &models.PaginatedItemsResponse{
		Items: itemsWithProgress,
		Pagination: models.PaginationMeta{
			Page:       page,
			Limit:      limit,
			Offset:     offset,
			Total:      totalCount,
			TotalPages: totalPages,
			HasNext:    hasNext,
			HasPrev:    hasPrev,
		},
	}, nil
}

// GetItemsPaginatedWithUserProgress retrieves items with user-specific progress data, filtering, validation and pagination metadata
func (s *ItemService) GetItemsPaginatedWithUserProgress(userID int, filter *models.ItemFilter) (*models.PaginatedItemsResponse, error) {
	// Validate filter parameters
	if filter.Category != nil && !models.IsValidCategory(*filter.Category) {
		return nil, fmt.Errorf("invalid category: %s", *filter.Category)
	}

	if filter.Status != nil && !models.IsValidStatus(*filter.Status) {
		return nil, fmt.Errorf("invalid status: %s", *filter.Status)
	}

	if filter.Limit != nil && *filter.Limit < 0 {
		return nil, fmt.Errorf("limit cannot be negative")
	}

	if filter.Offset != nil && *filter.Offset < 0 {
		return nil, fmt.Errorf("offset cannot be negative")
	}

	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	// Set default limit if not provided
	limit := 10
	if filter.Limit != nil {
		limit = *filter.Limit
	}

	offset := 0
	if filter.Offset != nil {
		offset = *filter.Offset
	}

	// Get total count with user progress
	totalCount, err := s.itemRepo.GetTotalCountWithUserProgress(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// Get items with user progress
	items, err := s.itemRepo.GetAllWithUserProgress(userID, filter)
	if err != nil {
		return nil, err
	}

	// Calculate pagination metadata
	totalPages := (totalCount + limit - 1) / limit // Ceiling division
	page := (offset / limit) + 1
	hasNext := offset+limit < totalCount
	hasPrev := offset > 0

	pagination := models.PaginationMeta{
		Total:      totalCount,
		Limit:      limit,
		Offset:     offset,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
		TotalPages: totalPages,
		Page:       page,
	}

	return &models.PaginatedItemsResponse{
		Items:      items,
		Pagination: pagination,
	}, nil
}

// GetNextItem retrieves the current in-progress item or a random pending item
func (s *ItemService) GetNextItem() (*models.Item, error) {
	return nil, fmt.Errorf("GetNextItem is deprecated - use GetNextItemWithUserProgress instead")
}

// GetNextItemWithUserProgress retrieves the current in-progress item or a random pending item for a user
func (s *ItemService) GetNextItemWithUserProgress(userID int) (*models.ItemWithProgress, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	// First check if there's already an in-progress item for this user
	inProgressItem, err := s.itemRepo.GetInProgressItemWithUserProgress(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check for in-progress item: %w", err)
	}

	// If there's an in-progress item, return it
	if inProgressItem != nil {
		return inProgressItem, nil
	}

	// Otherwise, get a random pending item for this user
	pendingItem, err := s.itemRepo.GetRandomPendingWithUserProgress(userID)
	if err != nil {
		return nil, err
	}

	// Reset any existing in-progress items for this user
	err = s.itemRepo.ResetInProgressItemsForUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to reset in-progress items: %w", err)
	}

	// Create or update user progress record to set it as in-progress
	err = s.itemRepo.UpsertUserProgressForItem(userID, pendingItem.ID, models.StatusInProgress)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert user progress: %w", err)
	}

	// Update the item status to in-progress and return it
	pendingItem.Status = models.StatusInProgress
	return pendingItem, nil
}

// SkipItem moves the current in-progress item back to pending and gets a new random item
func (s *ItemService) SkipItem() (*models.Item, error) {
	return nil, fmt.Errorf("SkipItem is deprecated - use SkipItemWithUserProgress instead")
}

// SkipItemWithUserProgress moves the current in-progress item back to pending and gets a new random item for a user
func (s *ItemService) SkipItemWithUserProgress(userID int) (*models.ItemWithProgress, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	// First, reset any existing in-progress items for this user back to pending
	err := s.itemRepo.ResetInProgressItemsForUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to reset in-progress items: %w", err)
	}

	// Get a new random pending item for this user
	pendingItem, err := s.itemRepo.GetRandomPendingWithUserProgress(userID)
	if err != nil {
		return nil, err
	}

	// Set the new item as in-progress
	err = s.itemRepo.UpsertUserProgressForItem(userID, pendingItem.ID, models.StatusInProgress)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert user progress: %w", err)
	}

	// Update the item status to in-progress and return it
	pendingItem.Status = models.StatusInProgress
	return pendingItem, nil
}

// CompleteItem marks an item as completed and handles completion logic
func (s *ItemService) CompleteItem(id int) (*models.Item, error) {
	return nil, fmt.Errorf("CompleteItem is deprecated - use CompleteItemWithUserProgress instead")
}

// CompleteItemWithUserProgress marks an item as completed for a specific user and handles user stats
func (s *ItemService) CompleteItemWithUserProgress(userID, itemID int) (*models.ItemWithProgress, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	if itemID <= 0 {
		return nil, fmt.Errorf("invalid item ID")
	}

	// Mark item as complete for the user
	item, err := s.itemRepo.CompleteItemForUser(userID, itemID)
	if err != nil {
		return nil, err
	}

	// Check if all items are now completed for this user
	pendingCount, err := s.itemRepo.CountPendingForUser(userID)
	if err != nil {
		// Log error but don't fail the completion
		fmt.Printf("Warning: failed to count pending items for user %d: %v\n", userID, err)
		return item, nil
	}

	// If all items are completed for this user, increment their completed_all_count
	if pendingCount == 0 {
		if err := s.statsRepo.IncrementUserCompletedAllCount(userID); err != nil {
			// Log error but don't fail the completion
			fmt.Printf("Warning: failed to increment user completed_all_count for user %d: %v\n", userID, err)
		}
	}

	return item, nil
}

// UpdateItem updates an existing item with validation
func (s *ItemService) UpdateItem(id int, req *models.UpdateItemRequest) (*models.Item, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid item ID")
	}

	// Validate category if provided
	if req.Category != nil && !models.IsValidCategory(*req.Category) {
		return nil, fmt.Errorf("invalid category: %s", *req.Category)
	}

	// Validate that at least one field is being updated
	if req.Title == nil && req.Link == nil && req.Category == nil && req.Subcategory == nil {
		return nil, fmt.Errorf("at least one field must be provided for update")
	}

	// Validate non-empty strings
	if req.Title != nil && *req.Title == "" {
		return nil, fmt.Errorf("title cannot be empty")
	}
	if req.Link != nil && *req.Link == "" {
		return nil, fmt.Errorf("link cannot be empty")
	}
	if req.Subcategory != nil && *req.Subcategory == "" {
		return nil, fmt.Errorf("subcategory cannot be empty")
	}

	return s.itemRepo.Update(id, req)
}

// DeleteItem removes an item
func (s *ItemService) DeleteItem(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid item ID")
	}

	return s.itemRepo.Delete(id)
}

// ResetAllItems marks all items as pending
func (s *ItemService) ResetAllItems() (int64, error) {
	return 0, fmt.Errorf("ResetAllItems is deprecated - use ResetAllItemsWithUserProgress instead")
}

// ResetAllItemsWithUserProgress resets all user progress for a specific user back to pending
func (s *ItemService) ResetAllItemsWithUserProgress(userID int) (int64, error) {
	if userID <= 0 {
		return 0, fmt.Errorf("invalid user ID")
	}

	return s.itemRepo.ResetAllUserProgress(userID)
}

// GetItemCounts returns basic item statistics
func (s *ItemService) GetItemCounts() (total, completed, pending int, err error) {
	return 0, 0, 0, fmt.Errorf("GetItemCounts is deprecated - use GetCountsForUser instead")
}

// GetCommonSubcategories returns the list of common subcategories for a given category
func (s *ItemService) GetCommonSubcategories(category models.Category) ([]string, error) {
	if !models.IsValidCategory(category) {
		return nil, fmt.Errorf("invalid category: %s", category)
	}

	subcategories, exists := models.CommonSubcategories[category]
	if !exists {
		return []string{"other"}, nil
	}

	return subcategories, nil
}

// ToggleStar toggles the starred status of an item
func (s *ItemService) ToggleStar(id int) (*models.Item, error) {
	return nil, fmt.Errorf("ToggleStar is deprecated - use ToggleStarWithUserProgress instead")
}

// ToggleStarWithUserProgress toggles the starred status of an item for a specific user
func (s *ItemService) ToggleStarWithUserProgress(userID, itemID int) (*models.ItemWithProgress, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	if itemID <= 0 {
		return nil, fmt.Errorf("invalid item ID")
	}

	return s.itemRepo.ToggleStarForUser(userID, itemID)
}

// UpdateStatus updates the status of an item
func (s *ItemService) UpdateStatus(id int, status models.Status) (*models.Item, error) {
	return nil, fmt.Errorf("UpdateStatus is deprecated - use UpdateStatusWithUserProgress instead")
}

// UpdateStatusWithUserProgress updates the status of an item for a specific user
func (s *ItemService) UpdateStatusWithUserProgress(userID, itemID int, status models.Status) (*models.ItemWithProgress, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	if itemID <= 0 {
		return nil, fmt.Errorf("invalid item ID")
	}

	// Validate status
	if !models.IsValidStatus(status) {
		return nil, fmt.Errorf("invalid status: %s", status)
	}

	// Don't allow setting status to in-progress through this method
	// Users should use GetNextItem or SkipItem for that
	if status == models.StatusInProgress {
		return nil, fmt.Errorf("cannot set status to in-progress directly. Use GetNextItem or SkipItem instead")
	}

	// If setting to done, check if all items will be completed and update stats
	if status == models.StatusDone {
		// Use the CompleteItemWithUserProgress method which handles the stats logic
		return s.CompleteItemWithUserProgress(userID, itemID)
	}

	// For other statuses (pending), just update the status
	return s.itemRepo.UpdateStatusForUser(userID, itemID, status)
}
