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

// GetNextItem retrieves the current in-progress item or a random pending item
func (s *ItemService) GetNextItem() (*models.Item, error) {
	// First check if there's already an in-progress item
	inProgressItem, err := s.itemRepo.GetInProgressItem()
	if err != nil {
		return nil, fmt.Errorf("failed to check for in-progress item: %w", err)
	}

	// If there's an in-progress item, return it
	if inProgressItem != nil {
		return inProgressItem, nil
	}

	// Otherwise, get a random pending item
	pendingItem, err := s.itemRepo.GetRandomPending()
	if err != nil {
		return nil, err
	}

	// Set it as in-progress
	return s.itemRepo.SetInProgress(pendingItem.ID)
}

// SkipItem moves the current in-progress item back to pending and gets a new random item
func (s *ItemService) SkipItem() (*models.Item, error) {
	// Get a random pending item (this will automatically reset any in-progress items)
	pendingItem, err := s.itemRepo.GetRandomPending()
	if err != nil {
		return nil, err
	}

	// Set it as in-progress (this will also reset the current in-progress item to pending)
	return s.itemRepo.SetInProgress(pendingItem.ID)
}

// CompleteItem marks an item as completed and handles completion logic
func (s *ItemService) CompleteItem(id int) (*models.Item, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid item ID")
	}

	// Mark item as complete
	item, err := s.itemRepo.MarkComplete(id)
	if err != nil {
		return nil, err
	}

	// Check if all items are now completed
	pendingCount, err := s.itemRepo.CountPending()
	if err != nil {
		// Log error but don't fail the completion
		// In a real app, you might want to use a proper logger here
		fmt.Printf("Warning: failed to count pending items: %v\n", err)
		return item, nil
	}

	// If all items are completed, increment the completed_all_count
	if pendingCount == 0 {
		if err := s.statsRepo.IncrementCompletedAllCount(); err != nil {
			// Log error but don't fail the completion
			fmt.Printf("Warning: failed to increment completed_all_count: %v\n", err)
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
	return s.itemRepo.ResetAll()
}

// GetItemCounts returns basic item statistics
func (s *ItemService) GetItemCounts() (total, completed, pending int, err error) {
	return s.itemRepo.GetCounts()
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
	if id <= 0 {
		return nil, fmt.Errorf("invalid item ID")
	}

	return s.itemRepo.ToggleStar(id)
}

// UpdateStatus updates the status of an item
func (s *ItemService) UpdateStatus(id int, status models.Status) (*models.Item, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid item ID")
	}

	// Validate status
	if !models.IsValidStatus(status) {
		return nil, fmt.Errorf("invalid status: %s", status)
	}

	// Don't allow setting status to in-progress through this method
	if status == models.StatusInProgress {
		return nil, fmt.Errorf("cannot set status to in-progress directly")
	}

	return s.itemRepo.UpdateStatus(id, status)
}
