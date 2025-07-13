import React, { useState, useEffect } from 'react';
import { useTheme } from '../contexts/ThemeContext';
import { 
  BarChart, 
  Bar, 
  PieChart, 
  Pie, 
  Cell, 
  XAxis, 
  YAxis, 
  CartesianGrid, 
  Tooltip, 
  Legend,
  ResponsiveContainer 
} from 'recharts';
import { Loader2 } from 'lucide-react';
import { statsApi, DetailedStats } from '../services/api';

const Stats: React.FC = () => {
  const { isDarkMode } = useTheme();
  const [stats, setStats] = useState<DetailedStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchStats();
  }, []);

  const fetchStats = async () => {
    try {
      setLoading(true);
      const data = await statsApi.getDetailedStats();
      setStats(data);
    } catch (err) {
      setError('Failed to fetch statistics');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="h-8 w-8 animate-spin text-indigo-600" />
      </div>
    );
  }

  if (error || !stats || !stats.overall) {
    return (
      <div className={`border px-4 py-3 rounded-lg ${
        isDarkMode 
          ? 'bg-red-900/20 border-red-800 text-red-300' 
          : 'bg-red-50 border-red-200 text-red-700'
      }`}>
        {error || 'Failed to load statistics'}
      </div>
    );
  }

  const pieData = [
    { name: 'Completed', value: stats?.overall?.completed_items || 0, color: '#10b981' },
    { name: 'Pending', value: stats?.overall?.pending_items || 0, color: '#f59e0b' },
  ];

  const categoryData = stats?.categories && Array.isArray(stats.categories) 
    ? stats.categories.map(cat => ({
        name: cat.category.toUpperCase(),
        completed: cat.completed_items,
        pending: cat.pending_items,
        total: cat.total_items,
      })) 
    : [];

  const getCategoryColor = (category: string) => {
    switch (category.toLowerCase()) {
      case 'dsa':
        return '#3b82f6';
      case 'lld':
        return '#10b981';
      case 'hld':
        return '#8b5cf6';
      default:
        return '#6b7280';
    }
  };

  return (
    <div>
      <div className="mb-8">
        <h2 className={`text-2xl font-bold ${isDarkMode ? 'text-gray-100' : 'text-gray-900'}`}>Statistics</h2>
        <p className={`mt-1 text-sm ${isDarkMode ? 'text-gray-400' : 'text-gray-600'}`}>
          Detailed insights into your interview preparation progress
        </p>
      </div>

      {/* Overall Stats */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
        <div className={`rounded-lg shadow p-6 ${isDarkMode ? 'bg-gray-800' : 'bg-white'}`}>
          <h3 className={`text-lg font-medium mb-4 ${isDarkMode ? 'text-gray-100' : 'text-gray-900'}`}>Overall Progress</h3>
          <ResponsiveContainer width="100%" height={300}>
            <PieChart>
              <Pie
                data={pieData}
                cx="50%"
                cy="50%"
                labelLine={false}
                label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                outerRadius={80}
                fill="#8884d8"
                dataKey="value"
              >
                {pieData.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={entry.color} />
                ))}
              </Pie>
              <Tooltip 
                contentStyle={{
                  backgroundColor: isDarkMode ? '#374151' : '#ffffff',
                  border: isDarkMode ? '1px solid #4b5563' : '1px solid #e5e7eb',
                  borderRadius: '6px',
                  color: isDarkMode ? '#f3f4f6' : '#1f2937'
                }}
              />
            </PieChart>
          </ResponsiveContainer>
          <div className="mt-4 text-center">
            <p className={`text-2xl font-bold ${isDarkMode ? 'text-gray-100' : 'text-gray-900'}`}>
              {(stats?.overall?.progress_percentage || 0).toFixed(1)}%
            </p>
            <p className={`text-sm ${isDarkMode ? 'text-gray-400' : 'text-gray-600'}`}>Overall Completion</p>
          </div>
        </div>

        <div className={`rounded-lg shadow p-6 ${isDarkMode ? 'bg-gray-800' : 'bg-white'}`}>
          <h3 className={`text-lg font-medium mb-4 ${isDarkMode ? 'text-gray-100' : 'text-gray-900'}`}>Category Breakdown</h3>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={categoryData}>
              <CartesianGrid strokeDasharray="3 3" stroke={isDarkMode ? '#374151' : '#e5e7eb'} />
              <XAxis 
                dataKey="name" 
                tick={{ fill: isDarkMode ? '#9ca3af' : '#6b7280' }}
                axisLine={{ stroke: isDarkMode ? '#4b5563' : '#d1d5db' }}
              />
              <YAxis 
                tick={{ fill: isDarkMode ? '#9ca3af' : '#6b7280' }}
                axisLine={{ stroke: isDarkMode ? '#4b5563' : '#d1d5db' }}
              />
              <Tooltip 
                contentStyle={{
                  backgroundColor: isDarkMode ? '#374151' : '#ffffff',
                  border: isDarkMode ? '1px solid #4b5563' : '1px solid #e5e7eb',
                  borderRadius: '6px',
                  color: isDarkMode ? '#f3f4f6' : '#1f2937'
                }}
              />
              <Legend />
              <Bar dataKey="completed" fill="#10b981" name="Completed" />
              <Bar dataKey="pending" fill="#f59e0b" name="Pending" />
            </BarChart>
          </ResponsiveContainer>
        </div>
      </div>

      {/* Category Details */}
      {stats?.categories && Array.isArray(stats.categories) && stats.categories.length > 0 && (
        <div className="space-y-6">
          {stats.categories.map((category) => (
            <div key={category.category} className={`rounded-lg shadow p-6 ${isDarkMode ? 'bg-gray-800' : 'bg-white'}`}>
            <div className="mb-4">
              <h3 className={`text-lg font-medium ${isDarkMode ? 'text-gray-100' : 'text-gray-900'}`}>
                {category.category.toUpperCase()} - {
                  category.category === 'dsa' ? 'Data Structures & Algorithms' :
                  category.category === 'lld' ? 'Low Level Design' :
                  'High Level Design'
                }
              </h3>
              <div className="mt-2">
                <div className={`flex justify-between text-sm mb-1 ${isDarkMode ? 'text-gray-400' : 'text-gray-600'}`}>
                  <span>Progress</span>
                  <span className="font-medium">{(category.progress_percentage || 0).toFixed(1)}%</span>
                </div>
                <div className={`w-full rounded-full h-2 ${isDarkMode ? 'bg-gray-600' : 'bg-gray-200'}`}>
                  <div
                    className="h-2 rounded-full transition-all duration-500"
                    style={{
                      width: `${category.progress_percentage}%`,
                      backgroundColor: getCategoryColor(category.category),
                    }}
                  />
                </div>
              </div>
            </div>

            {/* Subcategory Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {category.subcategories && Array.isArray(category.subcategories) && category.subcategories.map((sub) => (
                <div key={sub.subcategory} className={`border rounded-lg p-4 ${
                  isDarkMode ? 'border-gray-600' : 'border-gray-200'
                }`}>
                  <h4 className={`text-sm font-medium mb-2 ${isDarkMode ? 'text-gray-100' : 'text-gray-900'}`}>
                    {sub.subcategory.split('-').map(word => 
                      word.charAt(0).toUpperCase() + word.slice(1)
                    ).join(' ')}
                  </h4>
                  <div className="space-y-1 text-xs">
                    <div className="flex justify-between">
                      <span className={isDarkMode ? 'text-gray-400' : 'text-gray-600'}>Total:</span>
                      <span className="font-medium">{sub.total_items}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className={isDarkMode ? 'text-gray-400' : 'text-gray-600'}>Completed:</span>
                      <span className="font-medium text-green-600">{sub.completed_items}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className={isDarkMode ? 'text-gray-400' : 'text-gray-600'}>Pending:</span>
                      <span className="font-medium text-yellow-600">{sub.pending_items}</span>
                    </div>
                  </div>
                  <div className="mt-2">
                    <div className={`w-full rounded-full h-1.5 ${isDarkMode ? 'bg-gray-600' : 'bg-gray-200'}`}>
                      <div
                        className="h-1.5 rounded-full bg-indigo-600 transition-all duration-500"
                        style={{ width: `${sub.progress_percentage}%` }}
                      />
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>
      )}

      {/* Summary Card */}
      {stats?.overall && (
        <div className={`mt-8 rounded-lg p-6 ${
          isDarkMode 
            ? 'bg-indigo-900/30 border border-indigo-800' 
            : 'bg-indigo-50'
        }`}>
          <h3 className={`text-lg font-medium mb-2 ${
            isDarkMode ? 'text-indigo-300' : 'text-indigo-900'
          }`}>
            Completion Cycles: {stats.overall.completed_all_count || 0}
          </h3>
          <p className={`text-sm ${
            isDarkMode ? 'text-indigo-400' : 'text-indigo-700'
          }`}>
            You've completed all items {stats.overall.completed_all_count || 0} time{(stats.overall.completed_all_count || 0) !== 1 ? 's' : ''}.
            {(stats.overall.completed_all_count || 0) > 0 && ' Great job on your consistency!'}
          </p>
        </div>
      )}

      {/* Streak Card */}
      {stats?.overall && (
        <div className={`mt-6 rounded-lg p-6 ${
          isDarkMode 
            ? 'bg-orange-900/30 border border-orange-800' 
            : 'bg-orange-50'
        }`}>
          <div className="flex items-center justify-between">
            <div>
              <h3 className={`text-lg font-medium mb-2 ${
                isDarkMode ? 'text-orange-300' : 'text-orange-900'
              }`}>
                Daily Streak: {stats.overall.current_streak || 0} day{(stats.overall.current_streak || 0) !== 1 ? 's' : ''}
              </h3>
              <p className={`text-sm ${
                isDarkMode ? 'text-orange-400' : 'text-orange-700'
              }`}>
                {stats.overall.current_streak === 0 
                  ? 'Complete an item today to start your streak!' 
                  : stats.overall.current_streak === 1 
                  ? 'Great start! Complete another item tomorrow to continue your streak.'
                  : `Amazing! You've been consistent for ${stats.overall.current_streak} days in a row.`
                }
              </p>
            </div>
            <div className="text-right">
              <div className={`text-2xl font-bold ${
                isDarkMode ? 'text-orange-300' : 'text-orange-600'
              }`}>
                🔥
              </div>
              <div className={`text-sm mt-1 ${
                isDarkMode ? 'text-gray-400' : 'text-gray-600'
              }`}>
                Best: {stats.overall.longest_streak || 0}
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Stats; 