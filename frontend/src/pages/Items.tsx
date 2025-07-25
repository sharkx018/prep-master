import React, { useState, useEffect, useCallback, useMemo } from 'react';
import { useTheme } from '../contexts/ThemeContext';
import { useAuth } from '../contexts/AuthContext';
import { 
  ExternalLink, 
  CheckCircle, 
  CheckCircle2,
  Trash2, 
  Filter,
  Loader2,
  Star,
  Edit2,
  Save,
  X,
  Plus,
  ChevronLeft,
  ChevronRight
} from 'lucide-react';
import { itemsApi, Item, UpdateItemRequest, PaginationMeta } from '../services/api';
import MotivationalQuote from '../components/MotivationalQuote';

const Items: React.FC = () => {
  const { isDarkMode } = useTheme();
  const { isAdmin } = useAuth();
  const [items, setItems] = useState<Item[]>([]);
  const [pagination, setPagination] = useState<PaginationMeta>({
    total: 0,
    limit: 10,
    offset: 0,
    has_next: false,
    has_prev: false,
    total_pages: 0,
    page: 1
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [deleting, setDeleting] = useState<number | null>(null);
  const [updatingStatus, setUpdatingStatus] = useState<number | null>(null);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [editForm, setEditForm] = useState<UpdateItemRequest & { attachments: { [key: string]: string } }>({
    title: '',
    link: '',
    category: 'dsa',
    subcategory: '',
    attachments: {}
  });
  const [attachmentKey, setAttachmentKey] = useState('');
  const [attachmentValue, setAttachmentValue] = useState('');
  const [saving, setSaving] = useState(false);
  
  // Filter states
  const [categoryFilter, setCategoryFilter] = useState<string>('');
  const [subcategoryFilter, setSubcategoryFilter] = useState<string>('');
  const [statusFilter, setStatusFilter] = useState<string>('');
  
  // Subcategories cache - single source of truth
  const [subcategoriesCache, setSubcategoriesCache] = useState<{ [key: string]: string[] }>({});
  const [loadingSubcategories, setLoadingSubcategories] = useState<{ [key: string]: boolean }>({});
  
  // Pagination states
  const [currentPage, setCurrentPage] = useState(1);
  const [itemsPerPage, setItemsPerPage] = useState(10);

  // Memoized subcategories for filter dropdown
  const filterSubcategories = useMemo(() => {
    return categoryFilter ? subcategoriesCache[categoryFilter] || [] : [];
  }, [categoryFilter, subcategoriesCache]);

  // Memoized subcategories for edit form
  const editSubcategories = useMemo(() => {
    return editForm.category ? subcategoriesCache[editForm.category] || [] : [];
  }, [editForm.category, subcategoriesCache]);

  // Optimized function to fetch subcategories with caching
  const fetchSubcategories = useCallback(async (category: string) => {
    if (!category) return;
    
    // Return cached data if available
    if (subcategoriesCache[category]) {
      return subcategoriesCache[category];
    }
    
    // Prevent duplicate requests
    if (loadingSubcategories[category]) {
      return;
    }

    try {
      setLoadingSubcategories(prev => ({ ...prev, [category]: true }));
      const response = await itemsApi.getSubcategories(category);
      const subs = response?.subcategories && Array.isArray(response.subcategories) 
        ? response.subcategories 
        : [];
      
      // Cache the result
      setSubcategoriesCache(prev => ({ ...prev, [category]: subs }));
      return subs;
    } catch (err) {
      console.error('Failed to fetch subcategories', err);
      setSubcategoriesCache(prev => ({ ...prev, [category]: [] }));
      return [];
    } finally {
      setLoadingSubcategories(prev => ({ ...prev, [category]: false }));
    }
  }, [subcategoriesCache, loadingSubcategories]);

  // Fetch subcategories for filter when category changes
  useEffect(() => {
    if (categoryFilter) {
      fetchSubcategories(categoryFilter);
      // Reset subcategory filter when category filter changes
      setSubcategoryFilter('');
    } else {
      setSubcategoryFilter('');
    }
  }, [categoryFilter, fetchSubcategories]);

  // Fetch subcategories for edit form when category changes
  useEffect(() => {
    if (editingId && editForm.category) {
      fetchSubcategories(editForm.category);
    }
  }, [editForm.category, editingId, fetchSubcategories]);

  useEffect(() => {
    const loadItems = async () => {
      try {
        setLoading(true);
        const filters: any = {
          limit: itemsPerPage,
          offset: (currentPage - 1) * itemsPerPage
        };
        if (categoryFilter) filters.category = categoryFilter;
        if (subcategoryFilter) filters.subcategory = subcategoryFilter;
        if (statusFilter) filters.status = statusFilter;
        
        const data = await itemsApi.getItemsPaginated(filters);
        setItems(data.items || []);
        setPagination(data.pagination);
      } catch (err) {
        setError('Failed to fetch items');
        console.error(err);
        setItems([]);
      } finally {
        setLoading(false);
      }
    };
    
    loadItems();
  }, [categoryFilter, subcategoryFilter, statusFilter, currentPage, itemsPerPage]);

  const fetchItems = async () => {
    try {
      setLoading(true);
      const filters: any = {
        limit: itemsPerPage,
        offset: (currentPage - 1) * itemsPerPage
      };
      if (categoryFilter) filters.category = categoryFilter;
      if (subcategoryFilter) filters.subcategory = subcategoryFilter;
      if (statusFilter) filters.status = statusFilter;
      
      const data = await itemsApi.getItemsPaginated(filters);
      setItems(data.items || []);
      setPagination(data.pagination);
    } catch (err) {
      setError('Failed to fetch items');
      console.error(err);
      setItems([]);
    } finally {
      setLoading(false);
    }
  };

  // Reset to first page when filters change
  useEffect(() => {
    setCurrentPage(1);
  }, [categoryFilter, subcategoryFilter, statusFilter]);

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
  };

  const handleItemsPerPageChange = (newItemsPerPage: number) => {
    setItemsPerPage(newItemsPerPage);
    setCurrentPage(1); // Reset to first page when changing items per page
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

  const handleToggleStar = async (id: number) => {
    try {
      await itemsApi.toggleStar(id);
      await fetchItems();
    } catch (err) {
      setError('Failed to toggle star');
      console.error(err);
    }
  };

  const handleToggleStatus = async (id: number, currentStatus: string) => {
    try {
      setUpdatingStatus(id);
      // Toggle between done and pending, or mark in-progress as done
      const newStatus = currentStatus === 'done' ? 'pending' : 'done';
      await itemsApi.updateStatus(id, newStatus);
      
      // Dispatch custom event for widget to refresh if marking as complete
      if (newStatus === 'done') {
        window.dispatchEvent(new CustomEvent('itemCompleted'));
      }
      
      await fetchItems();
    } catch (err) {
      setError('Failed to update status');
      console.error(err);
    } finally {
      setUpdatingStatus(null);
    }
  };

  const handleEditClick = (item: Item) => {
    setEditingId(item.id);
    setEditForm({
      title: item.title,
      link: item.link,
      category: item.category,
      subcategory: item.subcategory,
      attachments: item.attachments || {}
    });
    setAttachmentKey('');
    setAttachmentValue('');
  };

  const handleCancelEdit = () => {
    setEditingId(null);
    setEditForm({
      title: '',
      link: '',
      category: 'dsa',
      subcategory: '',
      attachments: {}
    });
  };

  const handleSaveEdit = async () => {
    if (!editingId) return;

    try {
      setSaving(true);
      await itemsApi.updateItem(editingId, {
        title: editForm.title,
        link: editForm.link,
        category: editForm.category as 'dsa' | 'lld' | 'hld',
        subcategory: editForm.subcategory,
        attachments: editForm.attachments
      });
      await fetchItems();
      handleCancelEdit();
    } catch (err) {
      setError('Failed to update item');
      console.error(err);
    } finally {
      setSaving(false);
    }
  };

  const addAttachmentToEdit = () => {
    if (attachmentKey.trim() && attachmentValue.trim()) {
      setEditForm(prev => ({
        ...prev,
        attachments: {
          ...prev.attachments,
          [attachmentKey.trim()]: attachmentValue.trim()
        }
      }));
      setAttachmentKey('');
      setAttachmentValue('');
    }
  };

  const removeAttachmentFromEdit = (key: string) => {
    setEditForm(prev => {
      const newAttachments = { ...prev.attachments };
      delete newAttachments[key];
      return { ...prev, attachments: newAttachments };
    });
  };

  const handleEditCategoryChange = (category: string) => {
    setEditForm(prev => ({
      ...prev,
      category: category as 'dsa' | 'lld' | 'hld',
      subcategory: '' // Reset subcategory when category changes
    }));
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
      if (lowerKey.includes('youtube') || lowerKey.includes('yt') || lowerKey.includes('video')) {
        return (
          <svg className="h-3 w-3 text-red-500 mr-1.5 flex-shrink-0" fill="currentColor" viewBox="0 0 24 24">
            <path d="M23.498 6.186a3.016 3.016 0 0 0-2.122-2.136C19.505 3.545 12 3.545 12 3.545s-7.505 0-9.377.505A3.017 3.017 0 0 0 .502 6.186C0 8.07 0 12 0 12s0 3.93.502 5.814a3.016 3.016 0 0 0 2.122 2.136c1.871.505 9.376.505 9.376.505s7.505 0 9.377-.505a3.015 3.015 0 0 0 2.122-2.136C24 15.93 24 12 24 12s0-3.93-.502-5.814zM9.545 15.568V8.432L15.818 12l-6.273 3.568z"/>
          </svg>
        );
      } else if (lowerKey.includes('git') || lowerKey.includes('github')) {
        return (
          <svg className="h-3 w-3 text-gray-800 dark:text-gray-300 mr-1.5 flex-shrink-0" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
          </svg>
        );
      } else {
        return <ExternalLink className="h-3 w-3 text-gray-500 mr-1.5 flex-shrink-0" />;
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
      <div className="flex flex-wrap gap-2 mt-2">
        {Object.entries(attachments).map(([key, value], index) => {
          const isDifficulty = key.toLowerCase().includes('diff');
          const difficultyClass = isDifficulty ? getDifficultyColor(value) : '';
          
          return (
            <div key={key} className="inline-flex items-center">
              {isValidUrl(value) ? (
                <a
                  href={value}
                  target="_blank"
                  rel="noopener noreferrer"
                  className={`inline-flex items-center px-2 py-1 rounded-md text-xs font-medium border transition-all duration-200 hover:scale-105 hover:shadow-sm ${
                    isDarkMode 
                      ? 'bg-gray-700 border-gray-600 text-gray-200 hover:bg-gray-600 hover:border-gray-500' 
                      : 'bg-white border-gray-200 text-gray-800 hover:bg-gray-50 hover:border-gray-300'
                  }`}
                >
                  {getIconForKey(key)}
                  <span className="capitalize">
                    {key.replace(/[-_]/g, ' ')}
                  </span>
                </a>
              ) : (
                <div className={`inline-flex items-center px-2 py-1 rounded-md text-xs font-medium border ${
                  difficultyClass || (isDarkMode ? 'bg-gray-700 border-gray-600 text-gray-300' : 'bg-gray-50 border-gray-200 text-gray-700')
                }`}>
                  {isDifficulty ? (
                    <>
                      <span className="capitalize mr-1">{key.replace(/[-_]/g, ' ')}:</span>
                      <span className={`font-bold capitalize ${difficultyClass.split(' ')[0]}`}>
                        {value}
                      </span>
                    </>
                  ) : (
                    <>
                      <span className="capitalize mr-1">{key.replace(/[-_]/g, ' ')}:</span>
                      <span className="font-medium">{value}</span>
                    </>
                  )}
                </div>
              )}
            </div>
          );
        })}
      </div>
    );
  };

  return (
    <div>
      {/* Motivational Quote */}
      <MotivationalQuote />
      
      <div className="mb-8">
        <h2 className={`text-2xl font-bold ${isDarkMode ? 'text-gray-100' : 'text-gray-900'}`}>All Items</h2>
        <p className={`mt-1 text-sm ${isDarkMode ? 'text-gray-400' : 'text-gray-600'}`}>
          Browse and manage your interview prep items
        </p>
      </div>

      {/* Filters */}
      <div className={`rounded-lg shadow p-4 mb-6 ${isDarkMode ? 'bg-gray-800' : 'bg-white'}`}>
        <div className="flex items-center mb-4">
          <Filter className="h-5 w-5 text-gray-400 mr-2" />
          <h3 className={`text-sm font-medium ${isDarkMode ? 'text-gray-100' : 'text-gray-900'}`}>Filters</h3>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label htmlFor="category" className={`block text-sm font-medium mb-1 ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>
              Category
            </label>
            <select
              id="category"
              value={categoryFilter}
              onChange={(e) => setCategoryFilter(e.target.value)}
              className={`block w-full rounded-md shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm ${
                isDarkMode ? 'bg-gray-700 border-gray-600 text-gray-100' : 'bg-white border-gray-300 text-gray-900'
              }`}
            >
              <option value="">All Categories</option>
              <option value="dsa">Data Structures & Algorithms</option>
              <option value="lld">Low Level Design</option>
              <option value="hld">High Level Design</option>
            </select>
          </div>
          
          <div>
            <label htmlFor="subcategory" className={`block text-sm font-medium mb-1 ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>
              Subcategory
            </label>
            <select
              id="subcategory"
              value={subcategoryFilter}
              onChange={(e) => setSubcategoryFilter(e.target.value)}
              disabled={!categoryFilter || loadingSubcategories[categoryFilter]}
              className={`block w-full rounded-md shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm ${
                isDarkMode ? 'bg-gray-700 border-gray-600 text-gray-100' : 'bg-white border-gray-300 text-gray-900'
              }`}
            >
              <option value="">All Subcategories</option>
              {loadingSubcategories[categoryFilter] ? (
                <option value="">Loading subcategories...</option>
              ) : (
                filterSubcategories.map((subcategory) => (
                  <option key={subcategory} value={subcategory}>
                    {subcategory.split('-').map(word => 
                      word.charAt(0).toUpperCase() + word.slice(1)
                    ).join(' ')}
                  </option>
                ))
              )}
            </select>
          </div>
          
          <div>
            <label htmlFor="status" className={`block text-sm font-medium mb-1 ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>
              Status
            </label>
            <select
              id="status"
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value)}
              className={`block w-full rounded-md shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm ${
                isDarkMode ? 'bg-gray-700 border-gray-600 text-gray-100' : 'bg-white border-gray-300 text-gray-900'
              }`}
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
        <div className={`mb-6 border px-4 py-3 rounded-lg ${
          isDarkMode 
            ? 'bg-red-900/20 border-red-800 text-red-300' 
            : 'bg-red-50 border-red-200 text-red-700'
        }`}>
          {error}
        </div>
      )}

      {loading ? (
        <div className="flex items-center justify-center h-64">
          <Loader2 className="h-8 w-8 animate-spin text-indigo-600" />
        </div>
      ) : !Array.isArray(items) || items.length === 0 ? (
        <div className={`rounded-lg shadow p-8 text-center ${isDarkMode ? 'bg-gray-800' : 'bg-white'}`}>
          <p className={isDarkMode ? 'text-gray-400' : 'text-gray-600'}>No items found. Try adjusting your filters or add new items.</p>
        </div>
      ) : (
        <div className={`shadow overflow-hidden rounded-lg ${isDarkMode ? 'bg-gray-800' : 'bg-white'}`}>
          <ul className={`divide-y ${isDarkMode ? 'divide-gray-700' : 'divide-gray-200'}`}>
            {items.map((item) => (
              <li key={item.id} className={`px-6 py-4 ${
                isDarkMode ? 'hover:bg-gray-700' : 'hover:bg-gray-50'
              }`}>
                {editingId === item.id && isAdmin ? (
                  // Edit mode
                  <div className="space-y-4">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div>
                        <label className={`block text-sm font-medium mb-1 ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Title</label>
                        <input
                          type="text"
                          value={editForm.title}
                          onChange={(e) => setEditForm({ ...editForm, title: e.target.value })}
                          className={`block w-full rounded-md shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm ${
                            isDarkMode ? 'bg-gray-700 border-gray-600 text-gray-100' : 'bg-white border-gray-300 text-gray-900'
                          }`}
                        />
                      </div>
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Link</label>
                        <input
                          type="url"
                          value={editForm.link}
                          onChange={(e) => setEditForm({ ...editForm, link: e.target.value })}
                          className="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                        />
                      </div>
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Category</label>
                        <select
                          value={editForm.category}
                          onChange={(e) => {
                            handleEditCategoryChange(e.target.value);
                          }}
                          className="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                        >
                          <option value="dsa">Data Structures & Algorithms</option>
                          <option value="lld">Low Level Design</option>
                          <option value="hld">High Level Design</option>
                        </select>
                      </div>
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Subcategory</label>
                        <select
                          value={editForm.subcategory}
                          onChange={(e) => setEditForm({ ...editForm, subcategory: e.target.value })}
                          className="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                        >
                          <option value="">Select a subcategory</option>
                          {editSubcategories.map((subcategory) => (
                            <option key={subcategory} value={subcategory}>
                              {subcategory.split('-').map(word => 
                                word.charAt(0).toUpperCase() + word.slice(1)
                              ).join(' ')}
                            </option>
                          ))}
                        </select>
                      </div>
                    </div>

                    {/* Attachments section */}
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">Attachments</label>
                      
                      {/* Existing attachments */}
                      {Object.keys(editForm.attachments).length > 0 && (
                        <div className="mb-3 space-y-2">
                          {Object.entries(editForm.attachments).map(([key, value]) => (
                            <div key={key} className="flex items-center justify-between bg-gray-50 rounded-md px-3 py-2">
                              <span className="text-sm">
                                <span className="font-medium">{key}:</span> {value}
                              </span>
                              <button
                                type="button"
                                onClick={() => removeAttachmentFromEdit(key)}
                                className="text-red-500 hover:text-red-700"
                              >
                                <X className="h-4 w-4" />
                              </button>
                            </div>
                          ))}
                        </div>
                      )}

                      {/* Add new attachment */}
                      <div className="flex gap-2">
                        <input
                          type="text"
                          value={attachmentKey}
                          onChange={(e) => setAttachmentKey(e.target.value)}
                          onKeyPress={(e) => e.key === 'Enter' && (e.preventDefault(), addAttachmentToEdit())}
                          placeholder="Key"
                          className="flex-1 rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                        />
                        <input
                          type="text"
                          value={attachmentValue}
                          onChange={(e) => setAttachmentValue(e.target.value)}
                          onKeyPress={(e) => e.key === 'Enter' && (e.preventDefault(), addAttachmentToEdit())}
                          placeholder="Value"
                          className="flex-1 rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                        />
                        <button
                          type="button"
                          onClick={addAttachmentToEdit}
                          disabled={!attachmentKey.trim() || !attachmentValue.trim()}
                          className="inline-flex items-center px-3 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50"
                        >
                          <Plus className="h-4 w-4" />
                        </button>
                      </div>
                    </div>

                    {/* Action buttons */}
                    <div className="flex justify-end space-x-2">
                      <button
                        onClick={handleCancelEdit}
                        className="inline-flex items-center px-3 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
                      >
                        <X className="h-4 w-4 mr-1" />
                        Cancel
                      </button>
                      <button
                        onClick={handleSaveEdit}
                        disabled={saving || !editForm.title || !editForm.link}
                        className="inline-flex items-center px-3 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50"
                      >
                        {saving ? (
                          <Loader2 className="h-4 w-4 mr-1 animate-spin" />
                        ) : (
                          <Save className="h-4 w-4 mr-1" />
                        )}
                        Save
                      </button>
                    </div>
                  </div>
                ) : (
                  // View mode
                  <div className="flex items-center justify-between">
                    <div className="flex-1 cursor-pointer" onClick={() => handleEditClick(item)}>
                      <div className="flex items-center space-x-3 mb-2">
                        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getCategoryColor(item.category)}`}>
                          {item.category.toUpperCase()}
                        </span>
                        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                          isDarkMode ? 'bg-gray-700 text-gray-300' : 'bg-gray-100 text-gray-800'
                        }`}>
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
                        {item.starred && (
                          <Star className="h-4 w-4 text-yellow-500 fill-current" />
                        )}
                      </div>
                      <a
                        href={item.link}
                        target="_blank"
                        rel="noopener noreferrer"
                        className={`text-sm font-medium hover:underline ${isDarkMode ? 'text-gray-100 hover:text-indigo-400' : 'text-gray-900 hover:text-indigo-600'}`}
                      >
                        {item.title}
                      </a>
                      <p className={`text-sm mt-1 ${isDarkMode ? 'text-gray-400' : 'text-gray-500'}`}>
                        Added on {formatDate(item.created_at)}
                        {item.completed_at && ` • Completed on ${formatDate(item.completed_at)}`}
                      </p>
                      {renderAttachments(item.attachments)}
                    </div>
                    <div className="flex items-center space-x-2 ml-4">
                      {isAdmin && (
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            handleEditClick(item);
                          }}
                          className="p-2 text-gray-400 hover:text-indigo-600"
                          title="Edit item"
                        >
                          <Edit2 className="h-5 w-5" />
                        </button>
                      )}
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          handleToggleStar(item.id);
                        }}
                        className={`p-2 ${item.starred ? 'text-yellow-500' : 'text-gray-400'} hover:text-yellow-600`}
                        title={item.starred ? "Remove star" : "Add star"}
                      >
                        <Star className={`h-5 w-5 ${item.starred ? 'fill-current' : ''}`} />
                      </button>
                      <a
                        href={item.link}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="p-2 text-gray-400 hover:text-gray-600"
                        title="Open link"
                      >
                        <ExternalLink className="h-5 w-5" />
                      </a>
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          handleToggleStatus(item.id, item.status);
                        }}
                        disabled={updatingStatus === item.id}
                        className={`p-2 transition-colors ${
                          item.status === 'done' 
                            ? 'text-green-600 hover:text-gray-600' 
                            : item.status === 'in-progress'
                            ? 'text-gray-400 hover:text-green-600'
                            : 'text-gray-400 hover:text-green-600'
                        } disabled:opacity-50 disabled:cursor-not-allowed`}
                        title={
                          item.status === 'done' 
                            ? "Mark as pending" 
                            : item.status === 'in-progress'
                            ? "Mark as done"
                            : "Mark as done"
                        }
                      >
                        {updatingStatus === item.id ? (
                          <Loader2 className="h-5 w-5 animate-spin" />
                        ) : item.status === 'done' ? (
                          <CheckCircle2 className="h-5 w-5 text-green-600" />
                        ) : (
                          <CheckCircle className="h-5 w-5" />
                        )}
                      </button>
                      {isAdmin && (
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            handleDelete(item.id);
                          }}
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
                      )}
                    </div>
                  </div>
                )}
              </li>
            ))}
          </ul>
        </div>
      )}

      {/* Pagination Info and Controls */}
      {!loading && pagination.total > 0 && (
        <div className={`px-4 py-3 border-t sm:px-6 ${
          isDarkMode 
            ? 'bg-gray-800 border-gray-700' 
            : 'bg-white border-gray-200'
        }`}>
          <div className="flex items-center justify-between">
            <div className="flex items-center">
              <p className={`text-sm ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>
                Showing{' '}
                <span className="font-medium">
                  {pagination.offset + 1}
                </span>{' '}
                to{' '}
                <span className="font-medium">
                  {Math.min(pagination.offset + pagination.limit, pagination.total)}
                </span>{' '}
                of{' '}
                <span className="font-medium">{pagination.total}</span>{' '}
                results
              </p>
              <div className="ml-4 flex items-center">
                <label htmlFor="items-per-page" className="mr-2 text-sm text-gray-700">
                  Items per page:
                </label>
                <select
                  id="items-per-page"
                  value={itemsPerPage}
                  onChange={(e) => handleItemsPerPageChange(Number(e.target.value))}
                  className="rounded-md border-gray-300 text-sm focus:border-indigo-500 focus:ring-indigo-500"
                >
                  <option value={5}>5</option>
                  <option value={10}>10</option>
                  <option value={20}>20</option>
                  <option value={50}>50</option>
                </select>
              </div>
            </div>

            <div className="flex items-center space-x-2">
              <button
                onClick={() => handlePageChange(currentPage - 1)}
                disabled={!pagination.has_prev}
                className="relative inline-flex items-center px-2 py-2 rounded-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <ChevronLeft className="h-5 w-5" />
              </button>

              {/* Page numbers */}
              <div className="flex items-center space-x-1">
                {/* First page */}
                {currentPage > 3 && (
                  <>
                    <button
                      onClick={() => handlePageChange(1)}
                      className="relative inline-flex items-center px-3 py-2 rounded-md border border-gray-300 bg-white text-sm font-medium text-gray-700 hover:bg-gray-50"
                    >
                      1
                    </button>
                    {currentPage > 4 && (
                      <span className="text-gray-500">...</span>
                    )}
                  </>
                )}

                {/* Current page and neighbors */}
                {Array.from({ length: Math.min(5, pagination.total_pages) }, (_, i) => {
                  let pageNum: number;
                  if (pagination.total_pages <= 5) {
                    pageNum = i + 1;
                  } else if (currentPage <= 3) {
                    pageNum = i + 1;
                  } else if (currentPage >= pagination.total_pages - 2) {
                    pageNum = pagination.total_pages - 4 + i;
                  } else {
                    pageNum = currentPage - 2 + i;
                  }

                  if (pageNum < 1 || pageNum > pagination.total_pages) return null;

                  return (
                    <button
                      key={pageNum}
                      onClick={() => handlePageChange(pageNum)}
                      className={`relative inline-flex items-center px-3 py-2 rounded-md border text-sm font-medium ${
                        pageNum === currentPage
                          ? 'border-indigo-500 bg-indigo-50 text-indigo-600'
                          : 'border-gray-300 bg-white text-gray-700 hover:bg-gray-50'
                      }`}
                    >
                      {pageNum}
                    </button>
                  );
                })}

                {/* Last page */}
                {currentPage < pagination.total_pages - 2 && pagination.total_pages > 5 && (
                  <>
                    {currentPage < pagination.total_pages - 3 && (
                      <span className="text-gray-500">...</span>
                    )}
                    <button
                      onClick={() => handlePageChange(pagination.total_pages)}
                      className="relative inline-flex items-center px-3 py-2 rounded-md border border-gray-300 bg-white text-sm font-medium text-gray-700 hover:bg-gray-50"
                    >
                      {pagination.total_pages}
                    </button>
                  </>
                )}
              </div>

              <button
                onClick={() => handlePageChange(currentPage + 1)}
                disabled={!pagination.has_next}
                className="relative inline-flex items-center px-2 py-2 rounded-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <ChevronRight className="h-5 w-5" />
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Items;