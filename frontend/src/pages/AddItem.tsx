import React, { useState, useEffect, useCallback, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import { useTheme } from '../contexts/ThemeContext';
import { Save, Loader2, Plus, X } from 'lucide-react';
import { itemsApi, CreateItemRequest } from '../services/api';

const AddItem: React.FC = () => {
  const { isDarkMode } = useTheme();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [attachmentKey, setAttachmentKey] = useState('');
  const [attachmentValue, setAttachmentValue] = useState('');
  
  // Subcategories cache
  const [subcategoriesCache, setSubcategoriesCache] = useState<{ [key: string]: string[] }>({});
  const [loadingSubcategories, setLoadingSubcategories] = useState<{ [key: string]: boolean }>({});
  
  const [formData, setFormData] = useState<CreateItemRequest>({
    title: '',
    link: '',
    category: 'dsa',
    subcategory: '',
    attachments: {},
  });

  // Memoized subcategories for current category
  const subcategories = useMemo(() => {
    return formData.category ? subcategoriesCache[formData.category] || [] : [];
  }, [formData.category, subcategoriesCache]);

  // Optimized function to fetch subcategories with caching
  const fetchSubcategories = useCallback(async (category: string) => {
    if (!category) return [];
    
    // Return cached data if available
    if (subcategoriesCache[category]) {
      return subcategoriesCache[category];
    }
    
    // Prevent duplicate requests
    if (loadingSubcategories[category]) {
      return [];
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

  useEffect(() => {
    const loadSubcategories = async () => {
      const subs = await fetchSubcategories(formData.category);
      
      // Reset subcategory when category changes - only set to first item if we have subcategories
      // This ensures the subcategory gets properly reset and user sees the change
      setFormData(prev => ({ 
        ...prev, 
        subcategory: subs && subs.length > 0 ? subs[0] : '' 
      }));
    };

    loadSubcategories();
  }, [formData.category, fetchSubcategories]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    try {
      setLoading(true);
      setError(null);
      await itemsApi.createItem(formData);
      navigate('/items');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to create item');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    
    // If category is changing, immediately reset subcategory to show the change
    if (name === 'category') {
      setFormData(prev => ({ 
        ...prev, 
        category: value as CreateItemRequest['category'],
        subcategory: '' // Reset subcategory immediately when category changes
      }));
    } else {
      setFormData(prev => ({ ...prev, [name]: value }));
    }
  };

  const addAttachment = () => {
    if (attachmentKey.trim() && attachmentValue.trim()) {
      setFormData(prev => ({
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

  const removeAttachment = (key: string) => {
    setFormData(prev => {
      const newAttachments = { ...prev.attachments };
      delete newAttachments[key];
      return { ...prev, attachments: newAttachments };
    });
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      addAttachment();
    }
  };

  return (
    <div>
      <div className="mb-8">
        <h2 className={`text-2xl font-bold ${isDarkMode ? 'text-gray-100' : 'text-gray-900'}`}>Add New Item</h2>
        <p className={`mt-1 text-sm ${isDarkMode ? 'text-gray-400' : 'text-gray-600'}`}>
          Add a new problem or article to your practice list
        </p>
      </div>

      <div className={`shadow rounded-lg ${isDarkMode ? 'bg-gray-800' : 'bg-white'}`}>
        <form onSubmit={handleSubmit} className="p-6 space-y-6">
          {error && (
            <div className={`border px-4 py-3 rounded-lg ${
              isDarkMode 
                ? 'bg-red-900/20 border-red-800 text-red-300' 
                : 'bg-red-50 border-red-200 text-red-700'
            }`}>
              {error}
            </div>
          )}

          <div>
            <label htmlFor="title" className={`block text-sm font-medium ${
              isDarkMode ? 'text-gray-300' : 'text-gray-700'
            }`}>
              Title
            </label>
            <input
              type="text"
              name="title"
              id="title"
              required
              value={formData.title}
              onChange={handleChange}
              className={`mt-1 block w-full rounded-md shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm ${
                isDarkMode 
                  ? 'bg-gray-700 border-gray-600 text-gray-100 placeholder-gray-400' 
                  : 'bg-white border-gray-300 text-gray-900'
              }`}
              placeholder="e.g., Two Sum Problem"
            />
          </div>

          <div>
            <label htmlFor="link" className={`block text-sm font-medium ${
              isDarkMode ? 'text-gray-300' : 'text-gray-700'
            }`}>
              Link
            </label>
            <input
              type="url"
              name="link"
              id="link"
              required
              value={formData.link}
              onChange={handleChange}
              className={`mt-1 block w-full rounded-md shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm ${
                isDarkMode 
                  ? 'bg-gray-700 border-gray-600 text-gray-100 placeholder-gray-400' 
                  : 'bg-white border-gray-300 text-gray-900'
              }`}
              placeholder="https://leetcode.com/problems/two-sum/"
            />
          </div>

          <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
            <div>
              <label htmlFor="category" className={`block text-sm font-medium ${
                isDarkMode ? 'text-gray-300' : 'text-gray-700'
              }`}>
                Category
              </label>
              <select
                name="category"
                id="category"
                required
                value={formData.category}
                onChange={handleChange}
                className={`mt-1 block w-full rounded-md shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm ${
                  isDarkMode 
                    ? 'bg-gray-700 border-gray-600 text-gray-100' 
                    : 'bg-white border-gray-300 text-gray-900'
                }`}
              >
                <option value="dsa">Data Structures & Algorithms</option>
                <option value="lld">Low Level Design</option>
                <option value="hld">High Level Design</option>
              </select>
            </div>

            <div>
              <label htmlFor="subcategory" className={`block text-sm font-medium ${
                isDarkMode ? 'text-gray-300' : 'text-gray-700'
              }`}>
                Subcategory
              </label>
              <select
                name="subcategory"
                id="subcategory"
                required
                value={formData.subcategory}
                onChange={handleChange}
                disabled={loadingSubcategories[formData.category]}
                className={`mt-1 block w-full rounded-md shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm ${
                  isDarkMode 
                    ? 'bg-gray-700 border-gray-600 text-gray-100' 
                    : 'bg-white border-gray-300 text-gray-900'
                }`}
              >
                {loadingSubcategories[formData.category] ? (
                  <option value="">Loading subcategories...</option>
                ) : subcategories && subcategories.length > 0 ? (
                  subcategories.map((sub) => (
                    <option key={sub} value={sub}>
                      {sub.split('-').map(word => 
                        word.charAt(0).toUpperCase() + word.slice(1)
                      ).join(' ')}
                    </option>
                  ))
                ) : (
                  <option value="">No subcategories available</option>
                )}
              </select>
            </div>
          </div>

          <div>
            <label className={`block text-sm font-medium mb-2 ${
              isDarkMode ? 'text-gray-300' : 'text-gray-700'
            }`}>
              Attachments
            </label>
            
            {/* Display existing attachments */}
            {formData.attachments && Object.keys(formData.attachments).length > 0 && (
              <div className="mb-3 space-y-2">
                {Object.entries(formData.attachments).map(([key, value]) => (
                  <div key={key} className={`flex items-center justify-between rounded-md px-3 py-2 ${
                    isDarkMode ? 'bg-gray-700' : 'bg-gray-50'
                  }`}>
                    <span className={`text-sm ${isDarkMode ? 'text-gray-200' : 'text-gray-900'}`}>
                      <span className="font-medium">{key}:</span> {value}
                    </span>
                    <button
                      type="button"
                      onClick={() => removeAttachment(key)}
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
                onKeyPress={handleKeyPress}
                placeholder="Key (e.g., youtube)"
                className={`flex-1 rounded-md shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm ${
                  isDarkMode 
                    ? 'bg-gray-700 border-gray-600 text-gray-100 placeholder-gray-400' 
                    : 'bg-white border-gray-300 text-gray-900'
                }`}
              />
              <input
                type="text"
                value={attachmentValue}
                onChange={(e) => setAttachmentValue(e.target.value)}
                onKeyPress={handleKeyPress}
                placeholder="Value (e.g., https://youtube.com/...)"
                className={`flex-1 rounded-md shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm ${
                  isDarkMode 
                    ? 'bg-gray-700 border-gray-600 text-gray-100 placeholder-gray-400' 
                    : 'bg-white border-gray-300 text-gray-900'
                }`}
              />
              <button
                type="button"
                onClick={addAttachment}
                disabled={!attachmentKey.trim() || !attachmentValue.trim()}
                className={`inline-flex items-center px-3 py-2 border text-sm font-medium rounded-md focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed ${
                  isDarkMode 
                    ? 'border-gray-600 text-gray-300 bg-gray-700 hover:bg-gray-600' 
                    : 'border-gray-300 text-gray-700 bg-white hover:bg-gray-50'
                }`}
              >
                <Plus className="h-4 w-4" />
              </button>
            </div>
            <p className={`mt-1 text-xs ${isDarkMode ? 'text-gray-400' : 'text-gray-500'}`}>
              Add custom attributes like difficulty level, video links, etc.
            </p>
          </div>

          <div className="flex justify-end space-x-3">
            <button
              type="button"
              onClick={() => navigate('/items')}
              className={`inline-flex items-center px-4 py-2 border text-sm font-medium rounded-md focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 ${
                isDarkMode 
                  ? 'border-gray-600 text-gray-300 bg-gray-700 hover:bg-gray-600' 
                  : 'border-gray-300 text-gray-700 bg-white hover:bg-gray-50'
              }`}
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={loading}
              className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {loading ? (
                <Loader2 className="h-4 w-4 mr-2 animate-spin" />
              ) : (
                <Save className="h-4 w-4 mr-2" />
              )}
              Save Item
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default AddItem; 