import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Types
export interface Item {
  id: number;
  title: string;
  link: string;
  category: 'dsa' | 'lld' | 'hld';
  subcategory: string;
  status: 'done' | 'pending' | 'in-progress';
  created_at: string;
  completed_at?: string;
}

export interface CreateItemRequest {
  title: string;
  link: string;
  category: 'dsa' | 'lld' | 'hld';
  subcategory: string;
}

export interface UpdateItemRequest {
  title?: string;
  link?: string;
  category?: 'dsa' | 'lld' | 'hld';
  subcategory?: string;
}

export interface Stats {
  total_items: number;
  completed_items: number;
  pending_items: number;
  progress_percentage: number;
  completed_all_count: number;
}

export interface SubcategoryStats {
  subcategory: string;
  total_items: number;
  completed_items: number;
  pending_items: number;
  progress_percentage: number;
}

export interface CategoryStats {
  category: 'dsa' | 'lld' | 'hld';
  total_items: number;
  completed_items: number;
  pending_items: number;
  progress_percentage: number;
  subcategories: SubcategoryStats[];
}

export interface DetailedStats {
  overall: Stats;
  categories: CategoryStats[];
}

// API calls
export const itemsApi = {
  // Get all items with optional filters
  getItems: async (filters?: {
    category?: string;
    subcategory?: string;
    status?: string;
    limit?: number;
    offset?: number;
  }) => {
    const params = new URLSearchParams();
    if (filters) {
      Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined) {
          params.append(key, value.toString());
        }
      });
    }
    const response = await api.get<Item[]>(`/items?${params.toString()}`);
    return response.data;
  },

  // Get item by ID
  getItem: async (id: number) => {
    const response = await api.get<Item>(`/items/${id}`);
    return response.data;
  },

  // Create new item
  createItem: async (item: CreateItemRequest) => {
    const response = await api.post<Item>('/items', item);
    return response.data;
  },

  // Update item
  updateItem: async (id: number, updates: UpdateItemRequest) => {
    const response = await api.put<Item>(`/items/${id}`, updates);
    return response.data;
  },

  // Delete item
  deleteItem: async (id: number) => {
    const response = await api.delete(`/items/${id}`);
    return response.data;
  },

  // Mark item as complete
  completeItem: async (id: number): Promise<Item> => {
    const response = await api.put<Item>(`/items/${id}/complete`);
    return response.data;
  },

  // Get next random item
  getNextItem: async (): Promise<Item> => {
    const response = await api.get('/items/next');
    return response.data;
  },

  // Reset all items
  resetAllItems: async () => {
    const response = await api.post('/items/reset');
    return response.data;
  },

  // Get subcategories for a category
  getSubcategories: async (category: string) => {
    const response = await api.get<{ category: string; subcategories: string[] }>(
      `/items/subcategories/${category}`
    );
    return response.data;
  },

  skipItem: async (): Promise<Item> => {
    const response = await api.post('/items/skip');
    return response.data;
  },
};

export const statsApi = {
  // Get overall stats
  getStats: async () => {
    const response = await api.get<Stats>('/stats');
    return response.data;
  },

  // Get detailed stats with category breakdown
  getDetailedStats: async () => {
    const response = await api.get<DetailedStats>('/stats/detailed');
    return response.data;
  },

  // Get category stats
  getCategoryStats: async (category: string) => {
    const response = await api.get<CategoryStats>(`/stats/category/${category}`);
    return response.data;
  },

  // Get subcategory stats
  getSubcategoryStats: async (category: string, subcategory: string) => {
    const response = await api.get<SubcategoryStats>(
      `/stats/category/${category}/subcategory/${subcategory}`
    );
    return response.data;
  },

  // Reset completed all count
  resetCompletedAllCount: async () => {
    const response = await api.post('/stats/reset-completed-all');
    return response.data;
  },
};

export default api; 