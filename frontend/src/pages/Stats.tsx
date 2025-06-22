import React, { useState, useEffect } from 'react';
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

  if (error || !stats) {
    return (
      <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
        {error || 'Failed to load statistics'}
      </div>
    );
  }

  const pieData = [
    { name: 'Completed', value: stats.overall.completed_items, color: '#10b981' },
    { name: 'Pending', value: stats.overall.pending_items, color: '#f59e0b' },
  ];

  const categoryData = stats.categories ? stats.categories.map(cat => ({
    name: cat.category.toUpperCase(),
    completed: cat.completed_items,
    pending: cat.pending_items,
    total: cat.total_items,
  })) : [];

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
        <h2 className="text-2xl font-bold text-gray-900">Statistics</h2>
        <p className="mt-1 text-sm text-gray-600">
          Detailed insights into your interview preparation progress
        </p>
      </div>

      {/* Overall Stats */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Overall Progress</h3>
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
              <Tooltip />
            </PieChart>
          </ResponsiveContainer>
          <div className="mt-4 text-center">
            <p className="text-2xl font-bold text-gray-900">
              {stats.overall.progress_percentage.toFixed(1)}%
            </p>
            <p className="text-sm text-gray-600">Overall Completion</p>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Category Breakdown</h3>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={categoryData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" />
              <YAxis />
              <Tooltip />
              <Legend />
              <Bar dataKey="completed" fill="#10b981" name="Completed" />
              <Bar dataKey="pending" fill="#f59e0b" name="Pending" />
            </BarChart>
          </ResponsiveContainer>
        </div>
      </div>

      {/* Category Details */}
      {stats.categories && stats.categories.length > 0 && (
        <div className="space-y-6">
          {stats.categories.map((category) => (
            <div key={category.category} className="bg-white rounded-lg shadow p-6">
            <div className="mb-4">
              <h3 className="text-lg font-medium text-gray-900">
                {category.category.toUpperCase()} - {
                  category.category === 'dsa' ? 'Data Structures & Algorithms' :
                  category.category === 'lld' ? 'Low Level Design' :
                  'High Level Design'
                }
              </h3>
              <div className="mt-2">
                <div className="flex justify-between text-sm text-gray-600 mb-1">
                  <span>Progress</span>
                  <span>{category.progress_percentage.toFixed(1)}%</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
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
              {category.subcategories && category.subcategories.map((sub) => (
                <div key={sub.subcategory} className="border rounded-lg p-4">
                  <h4 className="text-sm font-medium text-gray-900 mb-2">
                    {sub.subcategory.split('-').map(word => 
                      word.charAt(0).toUpperCase() + word.slice(1)
                    ).join(' ')}
                  </h4>
                  <div className="space-y-1 text-xs">
                    <div className="flex justify-between">
                      <span className="text-gray-600">Total:</span>
                      <span className="font-medium">{sub.total_items}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600">Completed:</span>
                      <span className="font-medium text-green-600">{sub.completed_items}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600">Pending:</span>
                      <span className="font-medium text-yellow-600">{sub.pending_items}</span>
                    </div>
                  </div>
                  <div className="mt-2">
                    <div className="w-full bg-gray-200 rounded-full h-1.5">
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
      <div className="mt-8 bg-indigo-50 rounded-lg p-6">
        <h3 className="text-lg font-medium text-indigo-900 mb-2">
          Completion Cycles: {stats.overall.completed_all_count}
        </h3>
        <p className="text-sm text-indigo-700">
          You've completed all items {stats.overall.completed_all_count} time{stats.overall.completed_all_count !== 1 ? 's' : ''}.
          {stats.overall.completed_all_count > 0 && ' Great job on your consistency!'}
        </p>
      </div>
    </div>
  );
};

export default Stats; 