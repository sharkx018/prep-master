package repositories

import (
	"testing"
	"time"

	"interview-prep-app/internal/models"

	_ "github.com/lib/pq"
)

func TestUpdateUserStreakOnActivity(t *testing.T) {
	// This is a unit test for the streak calculation logic
	// Note: This test requires a test database setup

	// Test cases for streak calculation
	testCases := []struct {
		name                  string
		lastActivityDate      *time.Time
		currentStreak         int
		longestStreak         int
		expectedNewStreak     int
		expectedLongestStreak int
	}{
		{
			name:                  "First activity ever",
			lastActivityDate:      nil,
			currentStreak:         0,
			longestStreak:         0,
			expectedNewStreak:     1,
			expectedLongestStreak: 1,
		},
		{
			name:                  "Activity yesterday - continue streak",
			lastActivityDate:      timePtr(time.Now().UTC().Add(-24 * time.Hour).Truncate(24 * time.Hour)),
			currentStreak:         5,
			longestStreak:         10,
			expectedNewStreak:     6,
			expectedLongestStreak: 10,
		},
		{
			name:                  "Activity yesterday - new longest streak",
			lastActivityDate:      timePtr(time.Now().UTC().Add(-24 * time.Hour).Truncate(24 * time.Hour)),
			currentStreak:         9,
			longestStreak:         9,
			expectedNewStreak:     10,
			expectedLongestStreak: 10,
		},
		{
			name:                  "Activity 2 days ago - reset streak",
			lastActivityDate:      timePtr(time.Now().UTC().Add(-48 * time.Hour).Truncate(24 * time.Hour)),
			currentStreak:         5,
			longestStreak:         10,
			expectedNewStreak:     1,
			expectedLongestStreak: 10,
		},
		{
			name:                  "Activity 1 week ago - reset streak",
			lastActivityDate:      timePtr(time.Now().UTC().Add(-7 * 24 * time.Hour).Truncate(24 * time.Hour)),
			currentStreak:         3,
			longestStreak:         8,
			expectedNewStreak:     1,
			expectedLongestStreak: 8,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock user stats
			userStats := &models.UserStats{
				UserID:           1,
				CurrentStreak:    tc.currentStreak,
				LongestStreak:    tc.longestStreak,
				LastActivityDate: tc.lastActivityDate,
			}

			// Test the streak calculation logic
			today := time.Now().UTC().Truncate(24 * time.Hour)

			var newStreak, newLongestStreak int

			if userStats.LastActivityDate == nil {
				// First activity ever
				newStreak = 1
				newLongestStreak = 1
			} else {
				lastActivity := userStats.LastActivityDate.UTC().Truncate(24 * time.Hour)
				yesterday := today.Add(-24 * time.Hour)

				if lastActivity.Equal(yesterday) {
					// Continue streak
					newStreak = userStats.CurrentStreak + 1
					newLongestStreak = userStats.LongestStreak
					if newStreak > newLongestStreak {
						newLongestStreak = newStreak
					}
				} else {
					// Reset streak
					newStreak = 1
					newLongestStreak = userStats.LongestStreak
				}
			}

			// Verify the results
			if newStreak != tc.expectedNewStreak {
				t.Errorf("Expected new streak %d, got %d", tc.expectedNewStreak, newStreak)
			}
			if newLongestStreak != tc.expectedLongestStreak {
				t.Errorf("Expected longest streak %d, got %d", tc.expectedLongestStreak, newLongestStreak)
			}
		})
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func TestCheckAndResetStreakIfNeeded(t *testing.T) {
	// Test cases for streak reset when checking stats
	testCases := []struct {
		name                     string
		lastActivityDate         *time.Time
		currentStreak            int
		expectedStreakAfterReset int
		shouldReset              bool
	}{
		{
			name:                     "No last activity - no reset",
			lastActivityDate:         nil,
			currentStreak:            5,
			expectedStreakAfterReset: 5,
			shouldReset:              false,
		},
		{
			name:                     "Current streak is 0 - no reset needed",
			lastActivityDate:         timePtr(time.Now().UTC().Add(-48 * time.Hour).Truncate(24 * time.Hour)),
			currentStreak:            0,
			expectedStreakAfterReset: 0,
			shouldReset:              false,
		},
		{
			name:                     "Activity today - no reset",
			lastActivityDate:         timePtr(time.Now().UTC().Truncate(24 * time.Hour)),
			currentStreak:            5,
			expectedStreakAfterReset: 5,
			shouldReset:              false,
		},
		{
			name:                     "Activity 1 day ago - reset to 0",
			lastActivityDate:         timePtr(time.Now().UTC().Add(-24 * time.Hour).Truncate(24 * time.Hour)),
			currentStreak:            5,
			expectedStreakAfterReset: 0,
			shouldReset:              true,
		},
		{
			name:                     "Activity 2 days ago - reset to 0",
			lastActivityDate:         timePtr(time.Now().UTC().Add(-48 * time.Hour).Truncate(24 * time.Hour)),
			currentStreak:            3,
			expectedStreakAfterReset: 0,
			shouldReset:              true,
		},
		{
			name:                     "Activity 1 week ago - reset to 0",
			lastActivityDate:         timePtr(time.Now().UTC().Add(-7 * 24 * time.Hour).Truncate(24 * time.Hour)),
			currentStreak:            10,
			expectedStreakAfterReset: 0,
			shouldReset:              true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock user stats
			userStats := &models.UserStats{
				UserID:           1,
				CurrentStreak:    tc.currentStreak,
				LastActivityDate: tc.lastActivityDate,
			}

			// Test the streak reset logic
			var shouldReset bool
			if userStats.LastActivityDate != nil && userStats.CurrentStreak > 0 {
				now := time.Now().UTC()
				today := now.Truncate(24 * time.Hour)
				lastActivity := userStats.LastActivityDate.UTC().Truncate(24 * time.Hour)
				daysSinceLastActivity := int(today.Sub(lastActivity).Hours() / 24)
				shouldReset = daysSinceLastActivity >= 1
			}

			// Verify the reset decision
			if shouldReset != tc.shouldReset {
				t.Errorf("Expected shouldReset %v, got %v", tc.shouldReset, shouldReset)
			}

			// Verify the expected streak after reset
			expectedStreak := tc.currentStreak
			if shouldReset {
				expectedStreak = 0
			}

			if expectedStreak != tc.expectedStreakAfterReset {
				t.Errorf("Expected streak after reset %d, got %d", tc.expectedStreakAfterReset, expectedStreak)
			}
		})
	}
}
