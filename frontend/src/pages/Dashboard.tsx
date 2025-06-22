import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { 
  BookOpen, 
  CheckCircle, 
  Clock, 
  TrendingUp,
  ArrowRight,
  Loader2,
  Sparkles,
  Target,
  Trophy,
  RefreshCw,
  AlertCircle
} from 'lucide-react';
import { statsApi, itemsApi, Stats } from '../services/api';

const Dashboard: React.FC = () => {
  const [stats, setStats] = useState<Stats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [resetting, setResetting] = useState(false);
  const [showResetConfirm, setShowResetConfirm] = useState(false);

  useEffect(() => {
    fetchStats();
  }, []);

  const fetchStats = async () => {
    try {
      setLoading(true);
      const data = await statsApi.getStats();
      setStats(data);
    } catch (err) {
      setError('Failed to fetch statistics');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleReset = async () => {
    try {
      setResetting(true);
      await itemsApi.resetAllItems();
      await fetchStats(); // Refresh stats after reset
      setShowResetConfirm(false);
    } catch (err) {
      setError('Failed to reset items');
      console.error(err);
    } finally {
      setResetting(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="h-8 w-8 animate-spin text-indigo-600" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
        {error}
      </div>
    );
  }

  if (!stats) return null;

  const statCards = [
    {
      title: 'Total Items',
      value: stats?.total_items || 0,
      icon: BookOpen,
      color: 'bg-blue-500',
      lightColor: 'bg-blue-100',
      textColor: 'text-blue-600',
    },
    {
      title: 'Completed',
      value: stats?.completed_items || 0,
      icon: CheckCircle,
      color: 'bg-green-500',
      lightColor: 'bg-green-100',
      textColor: 'text-green-600',
    },
    {
      title: 'Pending',
      value: stats?.pending_items || 0,
      icon: Clock,
      color: 'bg-yellow-500',
      lightColor: 'bg-yellow-100',
      textColor: 'text-yellow-600',
    },
    {
      title: 'Completion Cycles',
      value: stats?.completed_all_count || 0,
      icon: TrendingUp,
      color: 'bg-purple-500',
      lightColor: 'bg-purple-100',
      textColor: 'text-purple-600',
    },
  ];

  return (
    <div>
      <div className="mb-8">
        <div className="bg-gradient-to-r from-indigo-600 to-purple-600 rounded-xl p-8 text-white shadow-xl relative overflow-hidden">
          <div className="absolute top-0 right-0 -mt-4 -mr-4 opacity-20">
            <Trophy className="h-32 w-32" />
          </div>
          <div className="relative z-10">
            <div className="flex items-center mb-4">
              <Target className="h-8 w-8 mr-3" />
              <h2 className="text-3xl font-bold">PrepMaster Dashboard</h2>
            </div>
            <p className="text-indigo-100 text-lg mb-2">
              Welcome back, Mukul! Ready to ace your next interview?
            </p>
            <div className="flex items-center space-x-2">
              <Sparkles className="h-5 w-5 text-yellow-300" />
              <span className="text-sm font-medium text-indigo-200">
                {stats?.completed_all_count > 0 
                  ? `Amazing! You've completed ${stats.completed_all_count} full cycles!` 
                  : 'Start your journey to interview success!'}
              </span>
            </div>
          </div>
        </div>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        {statCards.map((stat) => {
          const Icon = stat.icon;
          return (
            <div key={stat.title} className="bg-white rounded-xl shadow-lg p-6 hover:shadow-xl transition-all duration-300 transform hover:-translate-y-1 border border-gray-100">
              <div className="flex items-center">
                <div className={`${stat.lightColor} rounded-lg p-3 shadow-sm`}>
                  <Icon className={`h-6 w-6 ${stat.textColor}`} />
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600">{stat.title}</p>
                  <p className="text-3xl font-bold bg-gradient-to-r from-gray-700 to-gray-900 bg-clip-text text-transparent">
                    {stat.value}
                  </p>
                </div>
              </div>
            </div>
          );
        })}
      </div>

      {/* Progress Bar */}
      <div className="bg-gradient-to-br from-white to-indigo-50 rounded-xl shadow-lg p-8 mb-8 border border-indigo-100">
        <div className="flex items-center justify-between mb-6">
          <h3 className="text-xl font-bold text-gray-900">Overall Progress</h3>
          <div className="flex items-center space-x-2">
            {(stats?.progress_percentage || 0) >= 80 && <Trophy className="h-6 w-6 text-yellow-500" />}
            {(stats?.progress_percentage || 0) >= 50 && (stats?.progress_percentage || 0) < 80 && <Target className="h-6 w-6 text-blue-500" />}
            {(stats?.progress_percentage || 0) < 50 && <Clock className="h-6 w-6 text-gray-500" />}
          </div>
        </div>
        <div className="relative">
          <div className="flex mb-3 items-center justify-between">
            <div>
              <span className="text-sm font-bold text-indigo-700">
                {stats?.completed_items || 0} of {stats?.total_items || 0} items completed
              </span>
            </div>
            <div className="text-right">
              <span className="text-2xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">
                {(stats?.progress_percentage || 0).toFixed(1)}%
              </span>
            </div>
          </div>
          <div className="overflow-hidden h-4 text-xs flex rounded-full bg-gray-200 shadow-inner">
            <div
              style={{ width: `${stats?.progress_percentage || 0}%` }}
              className="shadow-none flex flex-col text-center whitespace-nowrap text-white justify-center bg-gradient-to-r from-indigo-500 to-purple-500 transition-all duration-1000 ease-out rounded-full"
            />
          </div>
          {(stats?.progress_percentage || 0) === 100 && (
            <p className="mt-3 text-sm font-medium text-green-600 flex items-center">
              <CheckCircle className="h-4 w-4 mr-1" />
              Congratulations! You've completed all items!
            </p>
          )}
        </div>
      </div>

      {/* Reset Section */}
      <div className="mb-8">
        {!showResetConfirm ? (
          <div className="flex justify-end">
            <button
              onClick={() => setShowResetConfirm(true)}
              disabled={stats?.total_items === 0}
              className="group flex items-center px-4 py-2 text-sm font-medium rounded-lg bg-white border border-gray-300 text-gray-700 hover:bg-gray-50 hover:border-red-300 hover:text-red-600 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <RefreshCw className="h-4 w-4 mr-2 group-hover:text-red-600" />
              Reset All Progress
            </button>
          </div>
        ) : (
          <div className="bg-red-50 border border-red-200 rounded-lg p-4">
            <div className="flex items-start">
              <AlertCircle className="h-5 w-5 text-red-600 mt-0.5 mr-3 flex-shrink-0" />
              <div className="flex-1">
                <h4 className="text-sm font-medium text-red-800 mb-1">
                  Are you sure you want to reset all progress?
                </h4>
                <p className="text-sm text-red-700 mb-3">
                  This will mark all {stats?.total_items} items as pending. This action cannot be undone.
                </p>
                <div className="flex space-x-3">
                  <button
                    onClick={handleReset}
                    disabled={resetting}
                    className="inline-flex items-center px-3 py-1.5 text-sm font-medium rounded-md text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 disabled:opacity-50"
                  >
                    {resetting ? (
                      <>
                        <Loader2 className="h-4 w-4 mr-1.5 animate-spin" />
                        Resetting...
                      </>
                    ) : (
                      <>
                        <RefreshCw className="h-4 w-4 mr-1.5" />
                        Yes, Reset All
                      </>
                    )}
                  </button>
                  <button
                    onClick={() => setShowResetConfirm(false)}
                    disabled={resetting}
                    className="inline-flex items-center px-3 py-1.5 text-sm font-medium rounded-md text-gray-700 bg-white border border-gray-300 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500"
                  >
                    Cancel
                  </button>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>

      {/* Quick Actions */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <Link
          to="/study"
          className="group bg-gradient-to-br from-indigo-500 to-purple-600 rounded-xl shadow-lg p-6 hover:shadow-2xl transition-all duration-300 transform hover:-translate-y-1 text-white"
        >
          <div className="flex items-center justify-between">
            <div>
              <h3 className="text-lg font-bold">Start Studying</h3>
              <p className="mt-1 text-sm text-indigo-100">
                Get a random item to study
              </p>
            </div>
            <div className="bg-white/20 rounded-full p-2 group-hover:bg-white/30 transition-colors">
              <ArrowRight className="h-5 w-5 text-white" />
            </div>
          </div>
        </Link>

        <Link
          to="/add-item"
          className="group bg-white rounded-xl shadow-lg p-6 hover:shadow-2xl transition-all duration-300 transform hover:-translate-y-1 border-2 border-transparent hover:border-indigo-200"
        >
          <div className="flex items-center justify-between">
            <div>
              <h3 className="text-lg font-bold text-gray-900 group-hover:text-indigo-600 transition-colors">Add New Item</h3>
              <p className="mt-1 text-sm text-gray-600">
                Add a new problem or article
              </p>
            </div>
            <div className="bg-indigo-50 rounded-full p-2 group-hover:bg-indigo-100 transition-colors">
              <ArrowRight className="h-5 w-5 text-indigo-600 group-hover:translate-x-1 transition-transform" />
            </div>
          </div>
        </Link>

        <Link
          to="/items"
          className="group bg-white rounded-xl shadow-lg p-6 hover:shadow-2xl transition-all duration-300 transform hover:-translate-y-1 border-2 border-transparent hover:border-purple-200"
        >
          <div className="flex items-center justify-between">
            <div>
              <h3 className="text-lg font-bold text-gray-900 group-hover:text-purple-600 transition-colors">Browse Items</h3>
              <p className="mt-1 text-sm text-gray-600">
                View and manage all items
              </p>
            </div>
            <div className="bg-purple-50 rounded-full p-2 group-hover:bg-purple-100 transition-colors">
              <ArrowRight className="h-5 w-5 text-purple-600 group-hover:translate-x-1 transition-transform" />
            </div>
          </div>
        </Link>
      </div>
    </div>
  );
};

export default Dashboard; 