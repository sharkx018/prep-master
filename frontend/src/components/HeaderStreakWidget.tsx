import React, { useState, useEffect, useCallback } from 'react';
import { useTheme } from '../contexts/ThemeContext';
import { Flame, Trophy, RefreshCw, Star } from 'lucide-react';
import { statsApi, Stats } from '../services/api';

const HeaderStreakWidget: React.FC = () => {
  const { isDarkMode } = useTheme();
  const [stats, setStats] = useState<Stats | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isRefreshing, setIsRefreshing] = useState(false);

  const fetchStats = useCallback(async (manual = false) => {
    try {
      if (manual) setIsRefreshing(true);
      const data = await statsApi.getStats();
      setStats(data);
    } catch (err) {
      console.error('Failed to fetch stats:', err);
    } finally {
      setIsLoading(false);
      setIsRefreshing(false);
    }
  }, []);

  const handleManualRefresh = useCallback(() => {
    fetchStats(true);
  }, [fetchStats]);

  useEffect(() => {
    fetchStats();
    
    const interval = setInterval(() => fetchStats(), 30000);
    
    const handleItemCompletion = () => {
      setTimeout(() => fetchStats(), 500);
    };

    window.addEventListener('itemCompleted', handleItemCompletion);
    
    return () => {
      clearInterval(interval);
      window.removeEventListener('itemCompleted', handleItemCompletion);
    };
  }, [fetchStats]);

  if (isLoading || !stats) {
    return (
      <div className="flex items-center space-x-3">
        <div className={`px-4 py-3 rounded-xl ${
          isDarkMode ? 'bg-gray-800/50' : 'bg-gray-100/70'
        }`}>
          <span className="text-sm font-medium text-gray-500">Loading...</span>
        </div>
      </div>
    );
  }

  const completionPercentage = stats.total_items > 0 
    ? (stats.completed_items / stats.total_items) * 100 
    : 0;

  const getProgressColor = () => {
    if (completionPercentage >= 90) return 'bg-emerald-500';
    if (completionPercentage >= 75) return 'bg-blue-500';
    if (completionPercentage >= 50) return 'bg-indigo-500';
    if (completionPercentage >= 25) return 'bg-orange-500';
    return 'bg-gray-500';
  };

  const getProgressGradient = () => {
    if (completionPercentage >= 90) return 'from-emerald-500 to-emerald-600';
    if (completionPercentage >= 75) return 'from-blue-500 to-blue-600';
    if (completionPercentage >= 50) return 'from-indigo-500 to-indigo-600';
    if (completionPercentage >= 25) return 'from-orange-500 to-orange-600';
    return 'from-gray-500 to-gray-600';
  };

  const getRankInfo = (count: number) => {
    if (count === 0) return { title: 'Newbie', color: 'text-gray-400', bgColor: 'from-gray-900/50 to-gray-800/40', borderColor: 'border-gray-600/50' };
    if (count === 1) return { title: 'Pupil', color: 'text-green-400', bgColor: 'from-green-900/50 to-green-800/40', borderColor: 'border-green-600/50' };
    if (count === 2) return { title: 'Specialist', color: 'text-cyan-400', bgColor: 'from-cyan-900/50 to-cyan-800/40', borderColor: 'border-cyan-600/50' };
    if (count === 3) return { title: 'Expert', color: 'text-blue-400', bgColor: 'from-blue-900/50 to-blue-800/40', borderColor: 'border-blue-600/50' };
    if (count === 4) return { title: 'Candidate Master', color: 'text-violet-400', bgColor: 'from-violet-900/50 to-violet-800/40', borderColor: 'border-violet-600/50' };
    if (count === 5) return { title: 'Master', color: 'text-orange-400', bgColor: 'from-orange-900/50 to-orange-800/40', borderColor: 'border-orange-600/50' };
    if (count === 6) return { title: 'International Master', color: 'text-orange-400', bgColor: 'from-orange-900/50 to-orange-800/40', borderColor: 'border-orange-600/50' };
    if (count === 7) return { title: 'Grandmaster', color: 'text-red-400', bgColor: 'from-red-900/50 to-red-800/40', borderColor: 'border-red-600/50' };
    if (count >= 10 && count < 15) return { title: 'International Grandmaster', color: 'text-red-400', bgColor: 'from-red-900/50 to-red-800/40', borderColor: 'border-red-600/50' };
    if (count >= 15) return { title: 'Legendary Grandmaster', color: 'text-red-400', bgColor: 'from-red-900/50 to-red-800/40', borderColor: 'border-red-600/50', isLegendary: true };
    return { title: 'Grandmaster', color: 'text-red-400', bgColor: 'from-red-900/50 to-red-800/40', borderColor: 'border-red-600/50' };
  };

  const rankInfo = getRankInfo(stats.completed_all_count || 0);
  
  return (
    <div className="flex items-center space-x-4">
      {/* Bold Streak Section */}
      <div className={`group flex items-center space-x-3 px-4 py-2 rounded-xl transition-all duration-200 hover:scale-105 shadow-lg ${
        isDarkMode 
          ? 'bg-gradient-to-r from-orange-900/50 to-red-900/40 border-2 border-orange-600/50 hover:border-orange-500/70 shadow-orange-900/30' 
          : 'bg-gradient-to-r from-orange-100 to-red-100 border-2 border-orange-300 hover:border-orange-400 shadow-orange-200/40'
      }`}>
        <div className="relative">
          <Flame className={`h-5 w-5 transition-all duration-200 drop-shadow-sm ${
            stats.current_streak > 0 
              ? 'text-orange-500 group-hover:text-orange-400' 
              : 'text-orange-400 group-hover:text-orange-300'
          }`} />
          {stats.current_streak > 0 && (
            <div className="absolute -top-1.5 -right-1.5 bg-gradient-to-r from-orange-500 to-red-500 text-white text-xs rounded-full h-4 w-4 flex items-center justify-center font-bold text-[9px] animate-pulse shadow-lg">
              {stats.current_streak > 9 ? '9+' : stats.current_streak}
            </div>
          )}
        </div>
        <div className="flex flex-col">
          {/* Current Streak */}
          <div className="flex items-center space-x-1.5">
            <span className={`text-sm font-bold ${
              isDarkMode ? 'text-orange-100' : 'text-orange-900'
            }`}>
              {stats.current_streak}
            </span>
            <span className={`text-xs font-semibold ${
              isDarkMode ? 'text-orange-200' : 'text-orange-700'
            }`}>
              {stats.current_streak === 1 ? 'day' : 'days'}
            </span>
          </div>
          {/* Best Streak Indicator */}
          {stats.longest_streak > 0 && (
            <div className="flex items-center space-x-1">
              <span className={`text-[10px] font-medium ${
                isDarkMode ? 'text-orange-300/80' : 'text-orange-600/80'
              }`}>
                best:
              </span>
              <span className={`text-[10px] font-bold ${
                isDarkMode ? 'text-orange-200/90' : 'text-orange-700/90'
              }`}>
                {stats.longest_streak}
              </span>
              {stats.current_streak === stats.longest_streak && (
                <span className="text-[10px]">🏆</span>
              )}
            </div>
          )}
        </div>
      </div>

      {/* Bold Progress Section */}
      <div className={`group flex items-center space-x-3 px-4 py-2 rounded-xl transition-all duration-200 hover:scale-105 shadow-lg ${
        isDarkMode 
          ? 'bg-gradient-to-r from-blue-900/50 to-indigo-900/40 border-2 border-blue-600/50 hover:border-blue-500/70 shadow-blue-900/30' 
          : 'bg-gradient-to-r from-blue-100 to-indigo-100 border-2 border-blue-300 hover:border-blue-400 shadow-blue-200/40'
      }`}>
        <div className="relative">
          <Trophy className={`h-5 w-5 transition-all duration-200 drop-shadow-sm ${
            completionPercentage >= 90 
              ? 'text-yellow-500 group-hover:text-yellow-400' 
              : 'text-blue-500 group-hover:text-blue-400'
          }`} />
          {completionPercentage >= 90 && (
            <div className="absolute -inset-1.5 bg-yellow-400/30 rounded-full animate-ping"></div>
          )}
        </div>
        
        <div className="flex flex-col space-y-1.5">
          <div className="flex items-center space-x-2">
            <span className={`text-sm font-bold ${
              isDarkMode ? 'text-blue-100' : 'text-blue-900'
            }`}>
              {completionPercentage.toFixed(1)}%
            </span>
            <span className={`text-xs font-semibold ${
              isDarkMode ? 'text-blue-200' : 'text-blue-700'
            }`}>
              {stats.completed_items}/{stats.total_items}
            </span>
          </div>
          
          {/* Bold progress bar */}
          <div className={`w-20 h-2 rounded-full ${
            isDarkMode ? 'bg-gray-700' : 'bg-gray-300'
          } overflow-hidden shadow-inner`}>
            <div 
              className={`h-full rounded-full bg-gradient-to-r ${getProgressGradient()} transition-all duration-500 shadow-sm`}
              style={{ width: `${Math.min(completionPercentage, 100)}%` }}
            />
            {/* Bold progress milestones */}
            <div className="relative -mt-2 flex justify-between w-20">
              <div className={`w-1 h-1 rounded-full shadow-sm ${
                completionPercentage >= 25 ? getProgressColor() : 'bg-gray-400'
              }`}></div>
              <div className={`w-1 h-1 rounded-full shadow-sm ${
                completionPercentage >= 50 ? getProgressColor() : 'bg-gray-400'
              }`}></div>
              <div className={`w-1 h-1 rounded-full shadow-sm ${
                completionPercentage >= 75 ? getProgressColor() : 'bg-gray-400'
              }`}></div>
            </div>
          </div>
        </div>
      </div>

      {/* Completion Cycles Section */}
      <div className={`group flex items-center space-x-3 px-3 py-2 rounded-xl transition-all duration-200 hover:scale-105 shadow-lg min-w-fit ${
        isDarkMode 
          ? `bg-gradient-to-r ${rankInfo.bgColor} border-2 ${rankInfo.borderColor} hover:border-opacity-70 shadow-gray-900/30` 
          : 'bg-gradient-to-r from-indigo-100 to-purple-100 border-2 border-indigo-300 hover:border-indigo-400 shadow-indigo-200/40'
      }`}>
        <div className="relative">
          <Star className={`h-5 w-5 transition-all duration-200 drop-shadow-sm ${rankInfo.color}`} />
          {stats.completed_all_count > 0 && (
            <div className="absolute -top-1.5 -right-1.5 bg-gradient-to-r from-indigo-500 to-purple-500 text-white text-xs rounded-full h-4 w-4 flex items-center justify-center font-bold text-[9px] shadow-lg">
              {stats.completed_all_count > 9 ? '9+' : stats.completed_all_count}
            </div>
          )}
        </div>
        
        <div className="flex flex-col min-w-0">
          <div className="flex items-center space-x-1.5">
            <span className={`text-xs font-bold ${rankInfo.color} whitespace-nowrap`}>
              {rankInfo.title}
            </span>
          </div>
          
          <div className="flex items-center space-x-1">
            <span className={`text-[10px] font-medium ${
              isDarkMode ? 'text-gray-400' : 'text-gray-600'
            }`}>
              cycles:
            </span>
            <span className={`text-[10px] font-bold ${rankInfo.color}`}>
              {stats.completed_all_count}
            </span>
          </div>
        </div>
      </div>

      {/* Bold Refresh Button */}
      <button
        onClick={handleManualRefresh}
        disabled={isRefreshing}
        className={`group px-3 py-2 rounded-xl transition-all duration-200 hover:scale-110 shadow-lg ${
          isDarkMode
            ? 'bg-gray-800/50 hover:bg-gray-700/60 text-gray-300 hover:text-white border-2 border-gray-600/50 hover:border-gray-500/70 shadow-gray-900/30'
            : 'bg-gray-100/70 hover:bg-gray-200/80 text-gray-600 hover:text-gray-800 border-2 border-gray-300 hover:border-gray-400 shadow-gray-200/40'
        } disabled:opacity-50 disabled:cursor-not-allowed`}
        title="Refresh stats"
      >
        <RefreshCw className={`h-5 w-5 transition-all duration-200 drop-shadow-sm ${
          isRefreshing ? 'animate-spin' : 'group-hover:rotate-180'
        }`} />
      </button>
    </div>
  );
};

export default HeaderStreakWidget; 