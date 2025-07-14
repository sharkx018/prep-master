# 🔥 Daily Streak Feature Implementation

## Overview

The daily streak feature encourages users to maintain consistent daily study habits by tracking consecutive days of item completion. Users earn streak points for completing at least one item per day, with their current streak resetting if they miss a day.

## Features

### 🎯 Core Functionality

- **Daily Streak Tracking**: Automatically tracks consecutive days of study activity
- **Streak Calculation**: Increments streak when users complete items on consecutive days
- **Streak Reset**: Resets to 0 when users miss a day (automatically detected when viewing stats)
- **Longest Streak**: Tracks the user's best streak achievement
- **Activity Deduplication**: Prevents multiple streak updates on the same day

### 🏆 Streak Logic

1. **First Activity**: When a user completes their first item ever, streak starts at 1
2. **Consecutive Days**: If user completed items yesterday and completes one today, streak increments
3. **Missed Days**: If user's last activity was 1+ days ago, streak resets to 0 (when viewing stats) or 1 (when completing an item)
4. **Same Day**: Multiple completions on the same day don't affect streak count
5. **Longest Streak**: Automatically updates when current streak exceeds previous best

## Streak Reset Behavior

The streak reset logic works differently depending on the context:

### When Viewing Stats (GetUserStats)
- **Automatic Reset**: Streak is automatically reset to 0 if there's a gap of 24+ hours since last activity
- **Real-time**: The reset happens immediately when stats are requested
- **Database Update**: The streak is updated in the database to reflect the reset

### When Completing Items (UpdateUserStreakOnActivity)
- **Activity-based Reset**: If user completes an item after missing days, streak resets to 1 (since they're completing an item that day)
- **Consecutive Days**: If user completed items yesterday and completes one today, streak increments
- **Same Day Protection**: Multiple completions on the same day don't affect streak count

This dual approach ensures users always see accurate streak information while rewarding them appropriately when they resume activity.

## Implementation Details

### Backend Implementation

#### Database Schema

The streak data is stored in the `user_stats` table with these fields:

```sql
current_streak INTEGER DEFAULT 0,
longest_streak INTEGER DEFAULT 0,
last_activity_date DATE,
```

#### Repository Methods

**`UpdateUserStreakOnActivity(userID int)`**
- Main method called when user completes an item
- Checks if user already has activity today to prevent duplicates
- Calculates new streak based on last activity date
- Updates both current and longest streak as needed

**`HasActivityToday(userID int)`**
- Checks if user has already completed an item today
- Prevents multiple streak updates on the same day

**`checkAndResetStreakIfNeeded(stats *UserStats)`**
- Automatically checks if streak should be reset when viewing stats
- Resets current streak to 0 if there's a gap of 24+ hours since last activity
- Updates the database and stats object in real-time

**`resetUserStreak(userID int)`**
- Helper method to reset user's current streak to 0 in the database
- Called when a streak gap is detected

**`GetUserStreakInfo(userID int)`**
- Retrieves current streak, longest streak, and last activity date
- Used for displaying streak information

#### Service Integration

The streak update is integrated into the `CompleteItemWithUserProgress` method in the item service:

```go
// Update user's daily streak
if err := s.statsRepo.UpdateUserStreakOnActivity(userID); err != nil {
    // Log error but don't fail the completion
    fmt.Printf("Warning: failed to update user streak for user %d: %v\n", userID, err)
}
```

### Frontend Implementation

#### API Interface

The `Stats` interface includes streak fields:

```typescript
export interface Stats {
  total_items: number;
  completed_items: number;
  pending_items: number;
  progress_percentage: number;
  completed_all_count: number;
  current_streak: number;
  longest_streak: number;
}
```

#### UI Components

**Dashboard Streak Section**
- Prominent streak display with fire emoji and animated flame icon
- Shows current streak with encouraging messages
- Displays longest streak achievement
- Only visible when user has an active streak

**Stats Page Streak Card**
- Detailed streak information with motivational messages
- Shows current streak, longest streak, and progress encouragement
- Integrated with overall statistics display

**Study Page Streak Display**
- Compact streak indicator during study sessions
- Shows current streak with badge-style counter
- Provides immediate feedback and motivation

#### Visual Design

- **Fire Theme**: Uses flame icons and orange/red color scheme
- **Animated Elements**: Pulsing flame icon for active streaks
- **Badge Counter**: Small circular badge showing current streak number
- **Motivational Messages**: Context-aware encouragement based on streak length

## Usage Examples

### Streak Progression

1. **Day 1**: User completes first item → Streak = 1, Longest = 1
2. **Day 2**: User completes another item → Streak = 2, Longest = 2
3. **Day 3**: User completes multiple items → Streak = 3, Longest = 3
4. **Day 5**: User completes item (missed day 4) → Streak = 1, Longest = 3

### Motivational Messages

- **Streak = 0**: "Complete an item today to start your streak!"
- **Streak = 1**: "Great start! Keep it up!"
- **Streak = 2-6**: "Keep the momentum going!"
- **Streak = 7+**: "Amazing consistency! You're on fire! 🔥"

## Testing

The streak calculation logic is thoroughly tested with unit tests covering:

- First activity ever
- Consecutive day streak continuation
- New longest streak achievements
- Streak resets after missed days
- Various gap scenarios (2 days, 1 week, etc.)

Run tests with:
```bash
cd backend && go test ./internal/repositories -v -run TestUpdateUserStreakOnActivity
```

## Technical Considerations

### Time Zone Handling

- All streak calculations use UTC time truncated to day boundaries
- Consistent 24-hour periods regardless of user's local time zone
- Prevents streak manipulation through time zone changes

### Performance

- Streak updates are non-blocking (errors are logged but don't fail item completion)
- Efficient database queries with proper indexing
- Minimal overhead on item completion flow

### Error Handling

- Graceful degradation if streak update fails
- Item completion always succeeds even if streak update fails
- Comprehensive error logging for debugging

## Future Enhancements

### Potential Features

1. **Streak Milestones**: Special achievements for 7, 30, 100+ day streaks
2. **Streak Recovery**: Grace period for missed days (e.g., 1-day buffer)
3. **Weekly Streaks**: Alternative streak tracking for weekly goals
4. **Streak Sharing**: Social features to share streak achievements
5. **Streak Reminders**: Notifications to maintain streaks
6. **Streak Freezes**: Allow users to "freeze" streaks during planned breaks

### Analytics

- Track average streak length across users
- Identify patterns in streak breaks
- Measure impact of streak feature on user engagement
- A/B test different streak incentives

## Conclusion

The daily streak feature successfully gamifies the learning experience by:

- Encouraging consistent daily study habits
- Providing immediate visual feedback and motivation
- Creating a sense of achievement and progress
- Maintaining user engagement through positive reinforcement

The implementation is robust, well-tested, and designed for scalability while maintaining excellent user experience across all platforms. 