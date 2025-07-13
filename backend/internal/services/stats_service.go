package services

import (
	"fmt"
	"interview-prep-app/internal/models"
	"interview-prep-app/internal/repositories"
)

// StatsService handles business logic for statistics
type StatsService struct {
	itemRepo  *repositories.ItemRepository
	statsRepo *repositories.StatsRepository
}

// NewStatsService creates a new stats service
func NewStatsService(itemRepo *repositories.ItemRepository, statsRepo *repositories.StatsRepository) *StatsService {
	return &StatsService{
		itemRepo:  itemRepo,
		statsRepo: statsRepo,
	}
}

// GetOverallStats returns overall statistics
func (s *StatsService) GetOverallStats() (*models.Stats, error) {
	return nil, fmt.Errorf("GetOverallStats is deprecated - use GetOverallStatsForUser instead")
}

// GetOverallStatsForUser retrieves comprehensive statistics for a specific user
func (s *StatsService) GetOverallStatsForUser(userID int) (*models.Stats, error) {
	// Get user-specific item counts
	total, completed, pending, _, err := s.itemRepo.GetCountsForUser(userID)
	if err != nil {
		return nil, err
	}

	// Calculate progress percentage
	var progressPercentage float64
	if total > 0 {
		progressPercentage = float64(completed) / float64(total) * 100
	}

	// Get user-specific completed all count and streak info
	userStats, err := s.statsRepo.GetUserStats(userID)
	if err != nil {
		return nil, err
	}

	return &models.Stats{
		TotalItems:         total,
		CompletedItems:     completed,
		PendingItems:       pending,
		ProgressPercentage: progressPercentage,
		CompletedAllCount:  userStats.CompletedAllCount,
		CurrentStreak:      userStats.CurrentStreak,
		LongestStreak:      userStats.LongestStreak,
	}, nil
}

// GetDetailedStats returns detailed statistics with category breakdown
func (s *StatsService) GetDetailedStats() (*models.DetailedStats, error) {
	return nil, fmt.Errorf("GetDetailedStats is deprecated - use GetDetailedStatsForUser instead")
}

// GetDetailedStatsForUser retrieves comprehensive statistics for a specific user including category and subcategory breakdown
func (s *StatsService) GetDetailedStatsForUser(userID int) (*models.DetailedStats, error) {
	// Get overall user stats
	overall, err := s.GetOverallStatsForUser(userID)
	if err != nil {
		return nil, err
	}

	// Get user-specific category counts
	categoryCounts, err := s.itemRepo.GetCountsByCategoryForUser(userID)
	if err != nil {
		return nil, err
	}

	// Get user-specific subcategory counts
	subcategoryCounts, err := s.itemRepo.GetCountsBySubcategoryForUser(userID)
	if err != nil {
		return nil, err
	}

	// Build category stats with subcategory breakdown
	var categories []models.CategoryWithSubcategoryStats

	for category, statusCounts := range categoryCounts {
		total := statusCounts[models.StatusPending] + statusCounts[models.StatusInProgress] + statusCounts[models.StatusDone]
		completed := statusCounts[models.StatusDone]
		pending := statusCounts[models.StatusPending]

		var progressPercentage float64
		if total > 0 {
			progressPercentage = float64(completed) / float64(total) * 100
		}

		// Get subcategories for this category
		var subcategories []models.SubcategoryStats
		if subCats, exists := subcategoryCounts[category]; exists {
			for subcategory, subStatusCounts := range subCats {
				subTotal := subStatusCounts[models.StatusPending] + subStatusCounts[models.StatusInProgress] + subStatusCounts[models.StatusDone]
				subCompleted := subStatusCounts[models.StatusDone]
				subPending := subStatusCounts[models.StatusPending]

				var subProgressPercentage float64
				if subTotal > 0 {
					subProgressPercentage = float64(subCompleted) / float64(subTotal) * 100
				}

				subcategories = append(subcategories, models.SubcategoryStats{
					Subcategory:        subcategory,
					TotalItems:         subTotal,
					CompletedItems:     subCompleted,
					PendingItems:       subPending,
					ProgressPercentage: subProgressPercentage,
				})
			}
		}

		categories = append(categories, models.CategoryWithSubcategoryStats{
			Category:           category,
			TotalItems:         total,
			CompletedItems:     completed,
			PendingItems:       pending,
			ProgressPercentage: progressPercentage,
			Subcategories:      subcategories,
		})
	}

	return &models.DetailedStats{
		Overall:    *overall,
		Categories: categories,
	}, nil
}

// GetCategoryStats returns statistics for a specific category
func (s *StatsService) GetCategoryStats(category models.Category) (*models.CategoryStats, error) {
	return nil, fmt.Errorf("GetCategoryStats is deprecated - use GetCategoryStatsForUser instead")
}

// GetCategoryStatsForUser retrieves statistics for a specific category and user
func (s *StatsService) GetCategoryStatsForUser(userID int, category models.Category) (*models.CategoryStats, error) {
	// Validate category
	if !models.IsValidCategory(category) {
		return nil, fmt.Errorf("invalid category: %s", category)
	}

	// Get user-specific category counts
	categoryCounts, err := s.itemRepo.GetCountsByCategoryForUser(userID)
	if err != nil {
		return nil, err
	}

	counts := categoryCounts[category]
	total := counts[models.StatusDone] + counts[models.StatusPending] + counts[models.StatusInProgress]
	completed := counts[models.StatusDone]
	pending := counts[models.StatusPending] + counts[models.StatusInProgress]

	var progressPercentage float64
	if total > 0 {
		progressPercentage = float64(completed) / float64(total) * 100
	}

	return &models.CategoryStats{
		Category:           category,
		TotalItems:         total,
		CompletedItems:     completed,
		PendingItems:       pending,
		ProgressPercentage: progressPercentage,
	}, nil
}

// GetSubcategoryStats returns statistics for a specific subcategory
func (s *StatsService) GetSubcategoryStats(category models.Category, subcategory string) (*models.SubcategoryStats, error) {
	return nil, fmt.Errorf("GetSubcategoryStats is deprecated - use GetSubcategoryStatsForUser instead")
}

// GetSubcategoryStatsForUser retrieves statistics for a specific category, subcategory, and user
func (s *StatsService) GetSubcategoryStatsForUser(userID int, category models.Category, subcategory string) (*models.SubcategoryStats, error) {
	// Validate category
	if !models.IsValidCategory(category) {
		return nil, fmt.Errorf("invalid category: %s", category)
	}

	// Get user-specific subcategory counts
	subcategoryCounts, err := s.itemRepo.GetCountsBySubcategoryForUser(userID)
	if err != nil {
		return nil, err
	}

	categoryData := subcategoryCounts[category]
	if categoryData == nil {
		return &models.SubcategoryStats{
			Subcategory:        subcategory,
			TotalItems:         0,
			CompletedItems:     0,
			PendingItems:       0,
			ProgressPercentage: 0,
		}, nil
	}

	statusCounts := categoryData[subcategory]
	total := statusCounts[models.StatusDone] + statusCounts[models.StatusPending] + statusCounts[models.StatusInProgress]
	completed := statusCounts[models.StatusDone]
	pending := statusCounts[models.StatusPending] + statusCounts[models.StatusInProgress]

	var progressPercentage float64
	if total > 0 {
		progressPercentage = float64(completed) / float64(total) * 100
	}

	return &models.SubcategoryStats{
		Subcategory:        subcategory,
		TotalItems:         total,
		CompletedItems:     completed,
		PendingItems:       pending,
		ProgressPercentage: progressPercentage,
	}, nil
}

// ResetCompletedAllCount resets the completed all count to zero
func (s *StatsService) ResetCompletedAllCount() error {
	return fmt.Errorf("ResetCompletedAllCount is deprecated - use ResetUserCompletedAllCount instead")
}

// ResetUserCompletedAllCount resets the completed all count for a specific user to zero
func (s *StatsService) ResetUserCompletedAllCount(userID int) error {
	if userID <= 0 {
		return fmt.Errorf("invalid user ID")
	}

	return s.statsRepo.ResetUserCompletedAllCount(userID)
}
