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

// GetOverallStats retrieves comprehensive statistics
func (s *StatsService) GetOverallStats() (*models.Stats, error) {
	// Get item counts
	total, completed, pending, err := s.itemRepo.GetCounts()
	if err != nil {
		return nil, err
	}

	// Calculate progress percentage
	var progressPercentage float64
	if total > 0 {
		progressPercentage = float64(completed) / float64(total) * 100
	}

	// Get completed all count
	appStats, err := s.statsRepo.GetAppStats()
	if err != nil {
		return nil, err
	}

	return &models.Stats{
		TotalItems:         total,
		CompletedItems:     completed,
		PendingItems:       pending,
		ProgressPercentage: progressPercentage,
		CompletedAllCount:  appStats.CompletedAllCount,
	}, nil
}

// GetDetailedStats retrieves comprehensive statistics including category and subcategory breakdown
func (s *StatsService) GetDetailedStats() (*models.DetailedStats, error) {
	// Get overall stats
	overallStats, err := s.GetOverallStats()
	if err != nil {
		return nil, err
	}

	// Get subcategory breakdown
	subcategoryCounts, err := s.itemRepo.GetCountsBySubcategory()
	if err != nil {
		return nil, err
	}

	// Build category stats with subcategories
	var categoryStats []models.CategoryWithSubcategoryStats
	for _, category := range models.ValidCategories() {
		categoryData := subcategoryCounts[category]

		var categoryTotal, categoryCompleted, categoryPending int
		var subcategoryStats []models.SubcategoryStats

		// Process each subcategory
		for subcategory, statusCounts := range categoryData {
			subTotal := statusCounts[models.StatusDone] + statusCounts[models.StatusPending] + statusCounts[models.StatusInProgress]
			subCompleted := statusCounts[models.StatusDone]
			subPending := statusCounts[models.StatusPending] + statusCounts[models.StatusInProgress]

			categoryTotal += subTotal
			categoryCompleted += subCompleted
			categoryPending += subPending

			var subProgressPercentage float64
			if subTotal > 0 {
				subProgressPercentage = float64(subCompleted) / float64(subTotal) * 100
			}

			subcategoryStats = append(subcategoryStats, models.SubcategoryStats{
				Subcategory:        subcategory,
				TotalItems:         subTotal,
				CompletedItems:     subCompleted,
				PendingItems:       subPending,
				ProgressPercentage: subProgressPercentage,
			})
		}

		var categoryProgressPercentage float64
		if categoryTotal > 0 {
			categoryProgressPercentage = float64(categoryCompleted) / float64(categoryTotal) * 100
		}

		categoryStats = append(categoryStats, models.CategoryWithSubcategoryStats{
			Category:           category,
			TotalItems:         categoryTotal,
			CompletedItems:     categoryCompleted,
			PendingItems:       categoryPending,
			ProgressPercentage: categoryProgressPercentage,
			Subcategories:      subcategoryStats,
		})
	}

	return &models.DetailedStats{
		Overall:    *overallStats,
		Categories: categoryStats,
	}, nil
}

// GetCategoryStats retrieves statistics for a specific category
func (s *StatsService) GetCategoryStats(category models.Category) (*models.CategoryStats, error) {
	// Validate category
	if !models.IsValidCategory(category) {
		return nil, fmt.Errorf("invalid category: %s", category)
	}

	// Get category counts
	categoryCounts, err := s.itemRepo.GetCountsByCategory()
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

// GetSubcategoryStats retrieves statistics for a specific category and subcategory
func (s *StatsService) GetSubcategoryStats(category models.Category, subcategory string) (*models.SubcategoryStats, error) {
	// Validate category
	if !models.IsValidCategory(category) {
		return nil, fmt.Errorf("invalid category: %s", category)
	}

	// Get subcategory counts
	subcategoryCounts, err := s.itemRepo.GetCountsBySubcategory()
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
	return s.statsRepo.ResetCompletedAllCount()
}
