import React, { useState, useEffect } from 'react';
import { 
  ExternalLink, 
  CheckCircle, 
  RefreshCw, 
  Loader2,
  AlertCircle,
  Star,
  Link,
  Hash,
  Info
} from 'lucide-react';
import { itemsApi, statsApi, Item, Stats } from '../services/api';

const Study: React.FC = () => {
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

    return (
      <div className="mt-6 mb-6 bg-gray-50 rounded-lg p-4">
        <div className="flex items-center mb-3">
          <Info className="h-5 w-5 text-gray-600 mr-2" />
          <h4 className="text-sm font-semibold text-gray-900">Additional Info</h4>
        </div>
        <div className="space-y-2">
          {Object.entries(attachments).map(([key, value]) => (
            <div key={key} className="flex items-start">
              {isValidUrl(value) ? (
                <a
                  href={value}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex items-center bg-white rounded-md px-3 py-2 shadow-sm border border-gray-200 hover:border-indigo-300 hover:shadow-md transition-all group"
                >
                  <Link className="h-4 w-4 text-indigo-600 mr-2 flex-shrink-0 group-hover:text-indigo-700" />
                  <span className="text-sm font-medium text-indigo-600 group-hover:text-indigo-700 underline">
                    {key}
                  </span>
                </a>
              ) : (
                <div className="flex items-center bg-white rounded-md px-3 py-2 shadow-sm border border-gray-200">
                  <Hash className="h-4 w-4 text-gray-500 mr-2 flex-shrink-0" />
                  <span className="text-sm font-medium text-gray-700 mr-2">{key}:</span>
                  <span className="text-sm text-gray-900">{value}</span>
                </div>
              )}
            </div>
          ))}
        </div>
      </div>
    );
  };

  return (
    <div>
      <div className="mb-8">
        <h2 className="text-2xl font-bold text-gray-900">Study Mode</h2>
        <p className="mt-1 text-sm text-gray-600">
          Practice with random items from your collection
        </p>
      </div>

      {/* Progress Bar */}
      {stats && (
        <div className="mb-6 bg-white rounded-lg shadow p-4">
          <div className="flex items-center justify-between mb-2">
            <h3 className="text-sm font-medium text-gray-700">Overall Progress</h3>
            <span className="text-sm font-semibold text-gray-900">
              {stats?.completed_items || 0} / {stats?.total_items || 0} items ({(stats?.progress_percentage || 0).toFixed(1)}%)
            </span>
          </div>
          <div className="w-full bg-gray-200 rounded-full h-3">
            <div
              className="bg-gradient-to-r from-indigo-500 to-purple-600 h-3 rounded-full transition-all duration-300"
              style={{ width: `${stats?.progress_percentage || 0}%` }}
            />
          </div>
          {stats && stats.completed_all_count > 0 && (
            <p className="mt-2 text-xs text-gray-600">
              Completion cycles: {stats.completed_all_count}
            </p>
          )}
        </div>
      )}

      {error && (
        <div className="mb-6 bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg flex items-center">
          <AlertCircle className="h-5 w-5 mr-2" />
          {error}
        </div>
      )}

      {!currentItem && !loading && !noItems && (
        <div className="bg-white rounded-lg shadow-lg p-8 text-center">
          <h3 className="text-lg font-medium text-gray-900 mb-4">
            Ready to start studying?
          </h3>
          <p className="text-gray-600 mb-6">
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
        <div className="bg-white rounded-lg shadow-lg p-8 text-center">
          <Loader2 className="h-8 w-8 animate-spin text-indigo-600 mx-auto mb-4" />
          <p className="text-gray-600">Loading next item...</p>
        </div>
      )}

      {noItems && (
        <div className="bg-white rounded-lg shadow-lg p-8 text-center">
          <CheckCircle className="h-12 w-12 text-green-500 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">
            All caught up!
          </h3>
          <p className="text-gray-600">
            You've completed all items. Great job! 🎉
          </p>
        </div>
      )}

      {currentItem && !loading && (
        <div className="bg-white rounded-lg shadow-lg overflow-hidden">
          <div className="p-6">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center space-x-3">
                <span className={`inline-flex items-center px-3 py-0.5 rounded-full text-sm font-medium ${getCategoryColor(currentItem.category)}`}>
                  {currentItem.category.toUpperCase()}
                </span>
                <span className="inline-flex items-center px-3 py-0.5 rounded-full text-sm font-medium bg-gray-100 text-gray-800">
                  {currentItem.subcategory}
                </span>
                {currentItem.starred && (
                  <Star className="h-5 w-5 text-yellow-500 fill-current" />
                )}
              </div>
            </div>

            <h3 className="text-xl font-semibold text-gray-900 mb-2">
              {currentItem.title}
            </h3>

            <p className="text-sm text-gray-600 mb-4">
              Category: {getCategoryFullName(currentItem.category)}
            </p>

            {renderAttachments(currentItem.attachments)}

            <div className="flex items-center justify-between mt-6">
              <a
                href={currentItem.link}
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              >
                <ExternalLink className="h-4 w-4 mr-2" />
                Open Link
              </a>

              <div className="space-x-3">
                <button
                  onClick={skipItem}
                  disabled={loading || completing}
                  className="inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
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