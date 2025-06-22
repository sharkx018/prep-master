import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Save, Loader2, Plus, X } from 'lucide-react';
import { itemsApi, CreateItemRequest } from '../services/api';

const AddItem: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [subcategories, setSubcategories] = useState<string[]>([]);
  const [attachmentKey, setAttachmentKey] = useState('');
  const [attachmentValue, setAttachmentValue] = useState('');
  
  const [formData, setFormData] = useState<CreateItemRequest>({
    title: '',
    link: '',
    category: 'dsa',
    subcategory: '',
    attachments: {},
  });

  useEffect(() => {
    fetchSubcategories(formData.category);
  }, [formData.category]);

  const fetchSubcategories = async (category: string) => {
    try {
      const response = await itemsApi.getSubcategories(category);
      const subs = response?.subcategories && Array.isArray(response.subcategories) 
        ? response.subcategories 
        : [];
      setSubcategories(subs);
      // Set default subcategory when category changes
      if (subs.length > 0) {
        setFormData(prev => ({ ...prev, subcategory: subs[0] }));
      }
    } catch (err) {
      console.error('Failed to fetch subcategories', err);
      setSubcategories([]); // Set empty array on error
    }
  };

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
    setFormData(prev => ({ ...prev, [name]: value }));
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
        <h2 className="text-2xl font-bold text-gray-900">Add New Item</h2>
        <p className="mt-1 text-sm text-gray-600">
          Add a new problem or article to your study list
        </p>
      </div>

      <div className="bg-white shadow rounded-lg">
        <form onSubmit={handleSubmit} className="p-6 space-y-6">
          {error && (
            <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
              {error}
            </div>
          )}

          <div>
            <label htmlFor="title" className="block text-sm font-medium text-gray-700">
              Title
            </label>
            <input
              type="text"
              name="title"
              id="title"
              required
              value={formData.title}
              onChange={handleChange}
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
              placeholder="e.g., Two Sum Problem"
            />
          </div>

          <div>
            <label htmlFor="link" className="block text-sm font-medium text-gray-700">
              Link
            </label>
            <input
              type="url"
              name="link"
              id="link"
              required
              value={formData.link}
              onChange={handleChange}
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
              placeholder="https://leetcode.com/problems/two-sum/"
            />
          </div>

          <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
            <div>
              <label htmlFor="category" className="block text-sm font-medium text-gray-700">
                Category
              </label>
              <select
                name="category"
                id="category"
                required
                value={formData.category}
                onChange={handleChange}
                className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
              >
                <option value="dsa">Data Structures & Algorithms</option>
                <option value="lld">Low Level Design</option>
                <option value="hld">High Level Design</option>
              </select>
            </div>

            <div>
              <label htmlFor="subcategory" className="block text-sm font-medium text-gray-700">
                Subcategory
              </label>
              <select
                name="subcategory"
                id="subcategory"
                required
                value={formData.subcategory}
                onChange={handleChange}
                className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
              >
                {subcategories && subcategories.length > 0 ? (
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
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Attachments
            </label>
            
            {/* Display existing attachments */}
            {formData.attachments && Object.keys(formData.attachments).length > 0 && (
              <div className="mb-3 space-y-2">
                {Object.entries(formData.attachments).map(([key, value]) => (
                  <div key={key} className="flex items-center justify-between bg-gray-50 rounded-md px-3 py-2">
                    <span className="text-sm">
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
                className="flex-1 rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
              />
              <input
                type="text"
                value={attachmentValue}
                onChange={(e) => setAttachmentValue(e.target.value)}
                onKeyPress={handleKeyPress}
                placeholder="Value (e.g., https://youtube.com/...)"
                className="flex-1 rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
              />
              <button
                type="button"
                onClick={addAttachment}
                disabled={!attachmentKey.trim() || !attachmentValue.trim()}
                className="inline-flex items-center px-3 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <Plus className="h-4 w-4" />
              </button>
            </div>
            <p className="mt-1 text-xs text-gray-500">
              Add custom attributes like difficulty level, video links, etc.
            </p>
          </div>

          <div className="flex justify-end space-x-3">
            <button
              type="button"
              onClick={() => navigate('/items')}
              className="inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
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