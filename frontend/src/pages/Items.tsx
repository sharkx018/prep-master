import React, { useState, useEffect } from 'react';
import { 
  ExternalLink, 
  CheckCircle, 
  Trash2, 
  Filter,
  Loader2
} from 'lucide-react';
import { itemsApi, Item } from '../services/api';

const Items: React.FC = () => {
  const [items, setItems] = useState<Item[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [deleting, setDeleting] = useState<number | null>(null);
  
  // Filter states
  const [categoryFilter, setCategoryFilter] = useState<string>('');
  const [subcategoryFilter, setSubcategoryFilter] = useState<string>('');
  const [statusFilter, setStatusFilter] = useState<string>('');

  useEffect(() => {
    const loadItems = async () => {
      try {
        setLoading(true);
        const filters: any = {};
        if (categoryFilter) filters.category = categoryFilter;
        if (subcategoryFilter) filters.subcategory = subcategoryFilter;
        if (statusFilter) filters.status = statusFilter;
        
        const data = await itemsApi.getItems(filters);
        setItems(Array.isArray(data) ? data : []);
      } catch (err) {
        setError('Failed to fetch items');
        console.error(err);
        setItems([]);
      } finally {
        setLoading(false);
      }
    };
    
    loadItems();
  }, [categoryFilter, subcategoryFilter, statusFilter]);

  const fetchItems = async () => {
    try {
      setLoading(true);
      const filters: any = {};
      if (categoryFilter) filters.category = categoryFilter;
      if (subcategoryFilter) filters.subcategory = subcategoryFilter;
      if (statusFilter) filters.status = statusFilter;
      
      const data = await itemsApi.getItems(filters);
      setItems(Array.isArray(data) ? data : []);
    } catch (err) {
      setError('Failed to fetch items');
      console.error(err);
      setItems([]);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: number) => {
    if (!window.confirm('Are you sure you want to delete this item?')) return;
    
    try {
      setDeleting(id);
      await itemsApi.deleteItem(id);
      await fetchItems();
    } catch (err) {
      setError('Failed to delete item');
      console.error(err);
    } finally {
      setDeleting(null);
    }
  };

  const handleComplete = async (id: number) => {
    try {
      await itemsApi.completeItem(id);
      await fetchItems();
    } catch (err) {
      setError('Failed to mark item as complete');
      console.error(err);
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

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    });
  };

  return (
    <div>
      <div className="mb-8">
        <h2 className="text-2xl font-bold text-gray-900">All Items</h2>
        <p className="mt-1 text-sm text-gray-600">
          Browse and manage your interview prep items
        </p>
      </div>

      {/* Filters */}
      <div className="bg-white rounded-lg shadow p-4 mb-6">
        <div className="flex items-center mb-4">
          <Filter className="h-5 w-5 text-gray-400 mr-2" />
          <h3 className="text-sm font-medium text-gray-900">Filters</h3>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label htmlFor="category" className="block text-sm font-medium text-gray-700 mb-1">
              Category
            </label>
            <select
              id="category"
              value={categoryFilter}
              onChange={(e) => setCategoryFilter(e.target.value)}
              className="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
            >
              <option value="">All Categories</option>
              <option value="dsa">Data Structures & Algorithms</option>
              <option value="lld">Low Level Design</option>
              <option value="hld">High Level Design</option>
            </select>
          </div>
          
          <div>
            <label htmlFor="subcategory" className="block text-sm font-medium text-gray-700 mb-1">
              Subcategory
            </label>
            <input
              type="text"
              id="subcategory"
              value={subcategoryFilter}
              onChange={(e) => setSubcategoryFilter(e.target.value)}
              placeholder="e.g., arrays"
              className="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
            />
          </div>
          
          <div>
            <label htmlFor="status" className="block text-sm font-medium text-gray-700 mb-1">
              Status
            </label>
            <select
              id="status"
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value)}
              className="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
            >
              <option value="">All Statuses</option>
              <option value="pending">Pending</option>
              <option value="in-progress">In Progress</option>
              <option value="done">Done</option>
            </select>
          </div>
        </div>
      </div>

      {error && (
        <div className="mb-6 bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
          {error}
        </div>
      )}

      {loading ? (
        <div className="flex items-center justify-center h-64">
          <Loader2 className="h-8 w-8 animate-spin text-indigo-600" />
        </div>
      ) : !Array.isArray(items) || items.length === 0 ? (
        <div className="bg-white rounded-lg shadow p-8 text-center">
          <p className="text-gray-600">No items found. Try adjusting your filters or add new items.</p>
        </div>
      ) : (
        <div className="bg-white shadow overflow-hidden rounded-lg">
          <ul className="divide-y divide-gray-200">
            {items.map((item) => (
              <li key={item.id} className="px-6 py-4 hover:bg-gray-50">
                <div className="flex items-center justify-between">
                  <div className="flex-1">
                    <div className="flex items-center space-x-3 mb-2">
                      <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getCategoryColor(item.category)}`}>
                        {item.category.toUpperCase()}
                      </span>
                      <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                        {item.subcategory}
                      </span>
                      <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                        item.status === 'done' 
                          ? 'bg-green-100 text-green-800' 
                          : item.status === 'in-progress'
                          ? 'bg-yellow-100 text-yellow-800'
                          : 'bg-gray-100 text-gray-800'
                      }`}>
                        {item.status === 'done' ? 'Done' : item.status === 'in-progress' ? 'In Progress' : 'Pending'}
                      </span>
                    </div>
                    <h3 className="text-sm font-medium text-gray-900">{item.title}</h3>
                    <p className="text-sm text-gray-500 mt-1">
                      Added on {formatDate(item.created_at)}
                      {item.completed_at && ` • Completed on ${formatDate(item.completed_at)}`}
                    </p>
                  </div>
                  <div className="flex items-center space-x-2 ml-4">
                    <a
                      href={item.link}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="p-2 text-gray-400 hover:text-gray-600"
                      title="Open link"
                    >
                      <ExternalLink className="h-5 w-5" />
                    </a>
                    {item.status === 'pending' && (
                      <button
                        onClick={() => handleComplete(item.id)}
                        className="text-green-600 hover:text-green-900"
                        title="Mark as complete"
                      >
                        <CheckCircle className="h-5 w-5" />
                      </button>
                    )}
                    <button
                      onClick={() => handleDelete(item.id)}
                      disabled={deleting === item.id}
                      className="p-2 text-gray-400 hover:text-red-600 disabled:opacity-50"
                      title="Delete item"
                    >
                      {deleting === item.id ? (
                        <Loader2 className="h-5 w-5 animate-spin" />
                      ) : (
                        <Trash2 className="h-5 w-5" />
                      )}
                    </button>
                  </div>
                </div>
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
};

export default Items; 