import React, { useState, useEffect } from 'react';
import { useTheme } from '../contexts/ThemeContext';
import { 
  ExternalLink, 
  CheckCircle, 
  RefreshCw, 
  Loader2,
  AlertCircle,
  Star,
  Link,
  Hash,
  Info,
  Flame
} from 'lucide-react';
import { itemsApi, statsApi, Item, Stats } from '../services/api';
import MotivationalQuote from '../components/MotivationalQuote';

const Study: React.FC = () => {
  const { isDarkMode } = useTheme();
  const [currentItem, setCurrentItem] = useState<Item | null>(null);
  const [loading, setLoading] = useState(false);
  const [completing, setCompleting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [noItems, setNoItems] = useState(false);
  const [stats, setStats] = useState<Stats | null>(null);

  const fetchStats = async () => {
    try {
      const data = await statsApi.getStats();
      setStats(data);
    } catch (err) {
      console.error('Failed to fetch stats:', err);
    }
  };

  useEffect(() => {
    fetchStats();
    fetchNextItem();
  }, []);

  const fetchNextItem = async () => {
    try {
      setLoading(true);
      setError(null);
      setNoItems(false);
      const item = await itemsApi.getNextItem();
      setCurrentItem(item);
      // Refresh stats when getting next item
      fetchStats();
    } catch (err: any) {
      if (err.response?.status === 404) {
        setNoItems(true);
        setCurrentItem(null);
      } else {
        setError('Failed to fetch next item');
        console.error(err);
      }
    } finally {
      setLoading(false);
    }
  };

  const markAsComplete = async () => {
    if (!currentItem) return;

    try {
      setCompleting(true);
      await itemsApi.completeItem(currentItem.id);
      
      // Dispatch custom event for widget to refresh
      window.dispatchEvent(new CustomEvent('itemCompleted'));
      
      // Fetch next item after marking complete
      await fetchNextItem();
      // Refresh stats after completing
      fetchStats();
    } catch (err) {
      setError('Failed to mark item as complete');
      console.error(err);
    } finally {
      setCompleting(false);
    }
  };

  const skipItem = async () => {
    try {
      setLoading(true);
      setError(null);
      const item = await itemsApi.skipItem();
      setCurrentItem(item);
      // Refresh stats when skipping
      fetchStats();
    } catch (err: any) {
      if (err.response?.status === 404) {
        setNoItems(true);
        setCurrentItem(null);
      } else {
        setError('Failed to skip item');
        console.error(err);
      }
    } finally {
      setLoading(false);
    }
  };

  const getCategoryColor = (category: string) => {
    switch (category) {
      case 'dsa':
        return 'bg-blue-100 text-blue-800';
      case 'lld':
        return 'bg-green-100 text-green-800';
      case 'hld':
        return 'bg-purple-100 text-purple-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const getCategoryFullName = (category: string) => {
    switch (category) {
      case 'dsa':
        return 'Data Structures & Algorithms';
      case 'lld':
        return 'Low Level Design';
      case 'hld':
        return 'High Level Design';
      default:
        return category.toUpperCase();
    }
  };

  const isValidUrl = (string: string) => {
    try {
      new URL(string);
      return true;
    } catch (_) {
      return false;
    }
  };

  const renderAttachments = (attachments: { [key: string]: string }) => {
    if (!attachments || Object.keys(attachments).length === 0) return null;

    // Function to get icon based on key name
    const getIconForKey = (key: string) => {
      const lowerKey = key.toLowerCase();
      if (lowerKey.includes('youtube') || lowerKey.includes('video')) {
        return (
          <svg className="h-4 w-4 text-red-500 mr-2 flex-shrink-0" fill="currentColor" viewBox="0 0 24 24">
            <path d="M23.498 6.186a3.016 3.016 0 0 0-2.122-2.136C19.505 3.545 12 3.545 12 3.545s-7.505 0-9.377.505A3.017 3.017 0 0 0 .502 6.186C0 8.07 0 12 0 12s0 3.93.502 5.814a3.016 3.016 0 0 0 2.122 2.136c1.871.505 9.376.505 9.376.505s7.505 0 9.377-.505a3.015 3.015 0 0 0 2.122-2.136C24 15.93 24 12 24 12s0-3.93-.502-5.814zM9.545 15.568V8.432L15.818 12l-6.273 3.568z"/>
          </svg>
        );
      } else if (lowerKey.includes('github') || lowerKey.includes('git')) {
        return (
          <svg className="h-4 w-4 text-gray-800 mr-2 flex-shrink-0" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
          </svg>
        );
      } else if (lowerKey.includes('link') || lowerKey.includes('url')) {
        return <Link className="h-4 w-4 text-indigo-600 mr-2 flex-shrink-0" />;
      } else {
        return <Link className="h-4 w-4 text-indigo-600 mr-2 flex-shrink-0" />;
      }
    };

    // Function to get difficulty color
    const getDifficultyColor = (value: string) => {
      const lowerValue = value.toLowerCase();
      if (lowerValue === 'easy') return 'text-green-600 bg-green-50 border-green-200';
      if (lowerValue === 'medium' || lowerValue === 'med') return 'text-yellow-600 bg-yellow-50 border-yellow-200';
      if (lowerValue === 'hard') return 'text-red-600 bg-red-50 border-red-200';
      return '';
    };

    return (
      <div className="mt-6 mb-6">
        <div className={`rounded-xl p-6 shadow-sm border ${
          isDarkMode 
            ? 'bg-gradient-to-r from-indigo-900/30 to-purple-900/30 border-indigo-800' 
            : 'bg-gradient-to-r from-indigo-50 to-purple-50 border-indigo-100'
        }`}>
          <div className="flex items-center mb-4">
            <div className={`p-2 rounded-lg shadow-sm mr-3 ${
              isDarkMode ? 'bg-gray-700' : 'bg-white'
            }`}>
              <Info className="h-5 w-5 text-indigo-600" />
            </div>
            <h4 className={`text-base font-bold ${isDarkMode ? 'text-gray-100' : 'text-gray-900'}`}>Additional Information</h4>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
            {Object.entries(attachments).map(([key, value], index) => {
              const isDifficulty = key.toLowerCase().includes('diff');
              const difficultyClass = isDifficulty ? getDifficultyColor(value) : '';
              
              return (
                <div
                  key={key}
                  className="transform transition-all duration-200 hover:scale-105"
                  style={{ animationDelay: `${index * 50}ms` }}
                >
                  {isValidUrl(value) ? (
                    <a
                      href={value}
                      target="_blank"
                      rel="noopener noreferrer"
                      className={`flex items-center rounded-lg px-4 py-3 shadow-sm border hover:border-indigo-300 hover:shadow-lg transition-all group cursor-pointer ${
                        isDarkMode 
                          ? 'bg-gray-700 border-gray-600' 
                          : 'bg-white border-gray-200'
                      }`}
                    >
                      {getIconForKey(key)}
                      <span className={`text-sm font-semibold group-hover:text-indigo-600 transition-colors capitalize ${
                        isDarkMode ? 'text-gray-200' : 'text-gray-800'
                      }`}>
                        {key.replace(/[-_]/g, ' ')}
                      </span>
                      <ExternalLink className="h-3 w-3 text-gray-400 ml-auto group-hover:text-indigo-600 transition-colors" />
                    </a>
                  ) : (
                    <div className={`flex items-center rounded-lg px-4 py-3 shadow-sm border ${
                      difficultyClass || (isDarkMode ? 'bg-gray-700 border-gray-600' : 'bg-white border-gray-200')
                    }`}>
                      {isDifficulty ? (
                        <>
                          <Hash className={`h-4 w-4 mr-2 flex-shrink-0 ${difficultyClass.split(' ')[0]}`} />
                          <span className={`text-sm font-medium mr-2 capitalize ${
                            isDarkMode ? 'text-gray-300' : 'text-gray-700'
                          }`}>{key.replace(/[-_]/g, ' ')}:</span>
                          <span className={`text-sm font-bold capitalize ${difficultyClass.split(' ')[0]}`}>
                            {value}
                          </span>
                        </>
                      ) : (
                        <>
                          <Hash className="h-4 w-4 text-gray-500 mr-2 flex-shrink-0" />
                          <span className={`text-sm font-medium mr-2 capitalize ${
                            isDarkMode ? 'text-gray-300' : 'text-gray-700'
                          }`}>{key.replace(/[-_]/g, ' ')}:</span>
                          <span className={`text-sm font-medium ${
                            isDarkMode ? 'text-gray-200' : 'text-gray-900'
                          }`}>{value}</span>
                        </>
                      )}
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        </div>
      </div>
    );
  };

  return (
    <div>
      {/* Motivational Quote */}
      <MotivationalQuote />
      
      <div className="mb-8">
        <h2 className={`text-2xl font-bold ${isDarkMode ? 'text-gray-100' : 'text-gray-900'}`}>Study Mode</h2>
        <p className={`mt-1 text-sm ${isDarkMode ? 'text-gray-400' : 'text-gray-600'}`}>
          Practice with random items from your collection
        </p>
      </div>

      {/* Progress Bar
      {stats && (
        <div className={`mb-6 rounded-lg shadow p-4 ${
          isDarkMode ? 'bg-gray-800' : 'bg-white'
        }`}>
          <div className="flex items-center justify-between mb-2">
            <h3 className={`text-sm font-medium ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Overall Progress</h3>
            <span className={`text-sm font-semibold ${isDarkMode ? 'text-gray-100' : 'text-gray-900'}`}>
              {stats?.completed_items || 0} / {stats?.total_items || 0} items ({(stats?.progress_percentage || 0).toFixed(1)}%)
            </span>
          </div>
          <div className={`w-full rounded-full h-3 ${isDarkMode ? 'bg-gray-600' : 'bg-gray-200'}`}>
            <div
              className="bg-gradient-to-r from-indigo-500 to-purple-600 h-3 rounded-full transition-all duration-300"
              style={{ width: `${stats?.progress_percentage || 0}%` }}
            />
          </div>
          {stats && stats.completed_all_count > 0 && (
            <p className={`mt-2 text-xs ${isDarkMode ? 'text-gray-400' : 'text-gray-600'}`}>
              Completion cycles: {stats.completed_all_count}
            </p>
          )}
        </div>
      )} */}

      {/* Streak Display
      {stats && (stats.current_streak > 0 || stats.longest_streak > 0) && (
        <div className={`mb-6 rounded-lg shadow p-4 ${
          isDarkMode ? 'bg-gradient-to-r from-orange-900/20 to-red-900/20' : 'bg-gradient-to-r from-orange-50 to-red-50'
        }`}>
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <div className="relative">
                <Flame className="h-6 w-6 text-orange-500" />
                {stats.current_streak > 0 && (
                  <div className="absolute -top-1 -right-1 bg-orange-500 text-white text-xs rounded-full h-4 w-4 flex items-center justify-center font-bold">
                    {stats.current_streak}
                  </div>
                )}
              </div>
              <div>
                <h3 className={`text-sm font-medium ${
                  isDarkMode ? 'text-orange-300' : 'text-orange-800'
                }`}>
                  {stats.current_streak > 0 
                    ? `${stats.current_streak} day streak!` 
                    : 'Start your streak today!'
                  }
                </h3>
                <p className={`text-xs ${
                  isDarkMode ? 'text-orange-400' : 'text-orange-600'
                }`}>
                  {stats.current_streak === 0 
                    ? 'Complete an item to begin' 
                    : stats.current_streak === 1 
                    ? 'Keep it up tomorrow!' 
                    : 'You\'re on fire! 🔥'
                  }
                </p>
              </div>
            </div>
            <div className="text-right">
              <p className={`text-xs ${
                isDarkMode ? 'text-gray-400' : 'text-gray-600'
              }`}>
                Best: {stats.longest_streak || 0}
              </p>
            </div>
          </div>
        </div>
      )} */}

      {error && (
        <div className={`mb-6 border px-4 py-3 rounded-lg flex items-center ${
          isDarkMode 
            ? 'bg-red-900/20 border-red-800 text-red-300' 
            : 'bg-red-50 border-red-200 text-red-700'
        }`}>
          <AlertCircle className="h-5 w-5 mr-2" />
          {error}
        </div>
      )}

      {!currentItem && !loading && !noItems && (
        <div className={`rounded-lg shadow-lg p-8 text-center ${
          isDarkMode ? 'bg-gray-800' : 'bg-white'
        }`}>
          <h3 className={`text-lg font-medium mb-4 ${isDarkMode ? 'text-gray-100' : 'text-gray-900'}`}>
            Ready to start studying?
          </h3>
          <p className={`mb-6 ${isDarkMode ? 'text-gray-400' : 'text-gray-600'}`}>
            Click the button below to get a random item from your pending list
          </p>
          <button
            onClick={fetchNextItem}
            className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
          >
            <RefreshCw className="h-4 w-4 mr-2" />
            Get Next Item
          </button>
        </div>
      )}

      {loading && (
        <div className={`rounded-lg shadow-lg p-8 text-center ${
          isDarkMode ? 'bg-gray-800' : 'bg-white'
        }`}>
          <Loader2 className="h-8 w-8 animate-spin text-indigo-600 mx-auto mb-4" />
          <p className={isDarkMode ? 'text-gray-400' : 'text-gray-600'}>Loading next item...</p>
        </div>
      )}

      {noItems && (
        <div className={`rounded-lg shadow-lg p-8 text-center ${
          isDarkMode ? 'bg-gray-800' : 'bg-white'
        }`}>
          <CheckCircle className="h-12 w-12 text-green-500 mx-auto mb-4" />
          <h3 className={`text-lg font-medium mb-2 ${isDarkMode ? 'text-gray-100' : 'text-gray-900'}`}>
            All caught up!
          </h3>
          <p className={isDarkMode ? 'text-gray-400' : 'text-gray-600'}>
            You've completed all items. Great job! 🎉
          </p>
        </div>
      )}

      {currentItem && !loading && (
        <div className={`rounded-lg shadow-lg overflow-hidden ${
          isDarkMode ? 'bg-gray-800' : 'bg-white'
        }`}>
          <div className="p-6">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center space-x-3">
                <span className={`inline-flex items-center px-3 py-0.5 rounded-full text-sm font-medium ${getCategoryColor(currentItem.category)}`}>
                  {currentItem.category.toUpperCase()}
                </span>
                <span className={`inline-flex items-center px-3 py-0.5 rounded-full text-sm font-medium ${
                  isDarkMode ? 'bg-gray-700 text-gray-300' : 'bg-gray-100 text-gray-800'
                }`}>
                  {currentItem.subcategory}
                </span>
                {currentItem.starred && (
                  <Star className="h-5 w-5 text-yellow-500 fill-current" />
                )}
              </div>
            </div>

            <h3 className={`text-xl font-semibold mb-2 ${isDarkMode ? 'text-gray-100' : 'text-gray-900'}`}>
              <a
                href={currentItem.link}
                target="_blank"
                rel="noopener noreferrer"
                className="hover:text-indigo-600 transition-colors"
              >
                {currentItem.title}
              </a>
            </h3>

            <p className={`text-sm mb-4 ${isDarkMode ? 'text-gray-400' : 'text-gray-600'}`}>
              Category: {getCategoryFullName(currentItem.category)}
            </p>

            {renderAttachments(currentItem.attachments)}

            <div className="flex items-center justify-between mt-6">
              <a
                href={currentItem.link}
                target="_blank"
                rel="noopener noreferrer"
                className={`inline-flex items-center px-4 py-2 border text-sm font-medium rounded-md focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 ${
                  isDarkMode 
                    ? 'border-gray-600 text-gray-300 bg-gray-700 hover:bg-gray-600' 
                    : 'border-gray-300 text-gray-700 bg-white hover:bg-gray-50'
                }`}
              >
                <ExternalLink className="h-4 w-4 mr-2" />
                Open Link
              </a>

              <div className="space-x-3">
                <button
                  onClick={skipItem}
                  disabled={loading || completing}
                  className={`inline-flex items-center px-4 py-2 border text-sm font-medium rounded-md focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed ${
                    isDarkMode 
                      ? 'border-gray-600 text-gray-300 bg-gray-700 hover:bg-gray-600' 
                      : 'border-gray-300 text-gray-700 bg-white hover:bg-gray-50'
                  }`}
                >
                  <RefreshCw className="h-4 w-4 mr-2" />
                  Skip
                </button>

                <button
                  onClick={markAsComplete}
                  disabled={completing}
                  className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {completing ? (
                    <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                  ) : (
                    <CheckCircle className="h-4 w-4 mr-2" />
                  )}
                  Mark as Complete
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Study; 