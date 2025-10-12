package services

import (
	"fmt"

	"interview-prep-app/internal/models"
	"interview-prep-app/internal/repositories"
)

// TestService handles business logic for tests
type TestService struct {
	testRepo *repositories.TestRepository
	itemRepo *repositories.ItemRepository
}

// NewTestService creates a new test service
func NewTestService(testRepo *repositories.TestRepository, itemRepo *repositories.ItemRepository) *TestService {
	return &TestService{
		testRepo: testRepo,
		itemRepo: itemRepo,
	}
}

// CreateTest creates a new test with random completed items from different categories
func (s *TestService) CreateTest(userID int) (*models.CreateTestResponse, error) {
	// Check if user already has an active test
	existingSessionID, _, err := s.testRepo.GetActiveTestByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check for existing test: %w", err)
	}

	if existingSessionID != "" {
		return nil, fmt.Errorf("user already has an active test")
	}

	// Get 2 random completed items from DSA
	dsaCategory := models.CategoryDSA
	doneStatus := models.StatusDone
	dsaLimit := 2
	dsaFilter := &models.ItemFilter{
		Category: &dsaCategory,
		Status:   &doneStatus,
		Limit:    &dsaLimit,
	}
	dsaItems, err := s.itemRepo.GetRandomItems(userID, dsaFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to get DSA items: %w", err)
	}
	if len(dsaItems) < 2 {
		return nil, fmt.Errorf("not enough completed DSA items (need 2, found %d)", len(dsaItems))
	}

	// Get 1 random completed item from LLD
	lldCategory := models.CategoryLLD
	lldLimit := 1
	lldFilter := &models.ItemFilter{
		Category: &lldCategory,
		Status:   &doneStatus,
		Limit:    &lldLimit,
	}
	lldItems, err := s.itemRepo.GetRandomItems(userID, lldFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to get LLD items: %w", err)
	}
	if len(lldItems) < 1 {
		return nil, fmt.Errorf("not enough completed LLD items (need 1, found %d)", len(lldItems))
	}

	// Get 1 random completed item from HLD with subcategory "interview questions"
	hldCategory := models.CategoryHLD
	hldSubcategory := "interview questions"
	hldLimit := 1
	hldFilter := &models.ItemFilter{
		Category:    &hldCategory,
		Subcategory: &hldSubcategory,
		Status:      &doneStatus,
		Limit:       &hldLimit,
	}
	hldItems, err := s.itemRepo.GetRandomItems(userID, hldFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to get HLD items: %w", err)
	}
	if len(hldItems) < 1 {
		return nil, fmt.Errorf("not enough completed HLD items with subcategory 'interview questions' (need 1, found %d)", len(hldItems))
	}

	// Combine all items
	allItems := append(dsaItems, lldItems...)
	allItems = append(allItems, hldItems...)

	// Extract item IDs
	itemIDs := make([]int, len(allItems))
	for i, item := range allItems {
		itemIDs[i] = item.ID
	}

	// Create test items in database
	sessionID, err := s.testRepo.CreateTestItems(userID, itemIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to create test items: %w", err)
	}

	return &models.CreateTestResponse{
		SessionID: sessionID,
		Items:     allItems,
		Message:   "Test created successfully with 4 items (2 DSA, 1 LLD, 1 HLD)",
	}, nil
}

// GetActiveTest retrieves the current active test for a user
func (s *TestService) GetActiveTest(userID int) (*models.ActiveTestResponse, error) {
	sessionID, itemIDs, err := s.testRepo.GetActiveTestByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active test: %w", err)
	}

	if sessionID == "" {
		return nil, nil // No active test
	}

	// Get items with user progress
	items := make([]models.ItemWithProgress, 0, len(itemIDs))
	for _, itemID := range itemIDs {
		item, err := s.itemRepo.GetByIDWithUserProgress(userID, itemID)
		if err != nil {
			return nil, fmt.Errorf("failed to get item %d: %w", itemID, err)
		}
		items = append(items, *item)
	}

	// Get created_at timestamp
	createdAt, err := s.testRepo.GetTestCreatedAt(userID, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get test created_at: %w", err)
	}

	return &models.ActiveTestResponse{
		SessionID: sessionID,
		Items:     items,
		CreatedAt: createdAt,
	}, nil
}

// CompleteTest marks a test as completed
func (s *TestService) CompleteTest(userID int, sessionID string, item_id string) error {
	return s.testRepo.UpdateTestStatus(userID, sessionID, item_id, models.TestStatusCompleted)
}

// AbandonTest marks a test as abandoned
func (s *TestService) AbandonTest(userID int, sessionID string, item_id string) error {
	return s.testRepo.UpdateTestStatus(userID, sessionID, item_id, models.TestStatusAbandoned)
}

// DeleteTest deletes a test
func (s *TestService) DeleteTest(userID int, sessionID string) error {
	return s.testRepo.DeleteTestsBySessionID(userID, sessionID)
}

// CheckCanCreateTest checks if a user can create a test (has miscellaneous item in progress)
func (s *TestService) CheckCanCreateTest(userID int) (bool, error) {
	// Get in-progress items
	inProgressStatus := models.StatusInProgress
	miscCategory := models.CategoryMiscellaneous
	subcategory := models.Test_n_revise

	filter := &models.ItemFilter{
		Status:      &inProgressStatus,
		Category:    &miscCategory,
		Subcategory: &subcategory,
	}

	items, err := s.itemRepo.GetAllWithUserProgress(userID, filter)
	if err != nil {
		return false, fmt.Errorf("failed to check for in-progress miscellaneous items: %w", err)
	}

	// If there's at least one in-progress miscellaneous item, user can create a test
	return len(items) > 0, nil
}
