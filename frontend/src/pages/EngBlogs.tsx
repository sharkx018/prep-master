import React, { useEffect, useState } from 'react';
import { useTheme } from '../contexts/ThemeContext';
import { 
  BookOpen, 
  ExternalLink, 
  Building2, 
  Loader2, 
  Search,
  Filter,
  Globe,
  ArrowUpRight,
  Star,
  Users
} from 'lucide-react';
import { engBlogsApi, EngBlog, EngBlogProblem } from '../services/api';
import MotivationalQuote from '../components/MotivationalQuote';

const EngBlogs: React.FC = () => {
  const { isDarkMode } = useTheme();
  const [engBlogs, setEngBlogs] = useState<EngBlog[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [expandedBlogs, setExpandedBlogs] = useState<Set<string>>(new Set());


  useEffect(() => {
    fetchEngBlogs();
  }, []);

  const fetchEngBlogs = async () => {
    try {
      setLoading(true);
      setError(null); // Clear any previous errors
      const data = await engBlogsApi.getEngBlogs();
      
      if (data && data.blogs && Array.isArray(data.blogs)) {
        setEngBlogs(data.blogs);
      } else {
        setError('Invalid data received from server');
      }
    } catch (err) {
      setError('Failed to fetch engineering blogs');
      console.error('Error fetching blogs:', err);
    } finally {
      setLoading(false);
    }
  };

  const toggleExpanded = (blogId: string) => {
    const newExpanded = new Set(expandedBlogs);
    if (newExpanded.has(blogId)) {
      newExpanded.delete(blogId);
    } else {
      newExpanded.add(blogId);
    }
    setExpandedBlogs(newExpanded);
  };

  const filteredBlogs = engBlogs.filter(blog =>
    blog.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    blog.practice_problems.some(problem => 
      problem.title.toLowerCase().includes(searchTerm.toLowerCase())
    )
  );


  const handleExternalLink = (url: string) => {
    window.open(url, '_blank', 'noopener,noreferrer');
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="flex items-center space-x-2">
          <Loader2 className="h-8 w-8 animate-spin text-blue-500" />
          <span className={`text-lg ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>
            Loading engineering blogs...
          </span>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="text-red-500 text-6xl mb-4">⚠️</div>
          <h2 className={`text-2xl font-bold mb-2 ${isDarkMode ? 'text-white' : 'text-gray-900'}`}>
            Oops! Something went wrong
          </h2>
          <p className={`text-lg mb-4 ${isDarkMode ? 'text-gray-300' : 'text-gray-600'}`}>
            {error}
          </p>
          <button
            onClick={fetchEngBlogs}
            className="bg-blue-500 hover:bg-blue-600 text-white px-6 py-3 rounded-lg transition-colors"
          >
            Try Again
          </button>
        </div>
      </div>
    );
  }

  return (
    <div>
      {/* Motivational Quote */}
      <MotivationalQuote />
      
      <div className="mb-8">
        <div className="flex items-center space-x-3 mb-4">
          <div className="p-3 bg-gradient-to-r from-blue-500 to-purple-600 rounded-xl shadow-lg">
            <Building2 className="h-8 w-8 text-white" />
          </div>
          <div>
            <h2 className={`text-2xl font-bold ${isDarkMode ? 'text-gray-100' : 'text-gray-900'}`}>
              Curated Engineering Blogs
            </h2>
            <p className={`mt-1 text-sm ${isDarkMode ? 'text-gray-400' : 'text-gray-600'}`}>
              Learn from the best tech companies' engineering practices
            </p>
          </div>
        </div>

        {/* Stats Bar */}
        <div className={`grid grid-cols-1 md:grid-cols-2 gap-4 p-4 rounded-xl ${
          isDarkMode ? 'bg-gray-800 border border-gray-700' : 'bg-white border border-gray-200'
        } shadow-sm mb-6`}>
          <div className="flex items-center space-x-3">
            <Users className="h-5 w-5 text-blue-500" />
            <div>
              <p className={`text-sm ${isDarkMode ? 'text-gray-400' : 'text-gray-600'}`}>
                Companies
              </p>
              <p className={`text-lg font-semibold ${isDarkMode ? 'text-white' : 'text-gray-900'}`}>
                {engBlogs.length}
              </p>
            </div>
          </div>
          <div className="flex items-center space-x-3">
            <BookOpen className="h-5 w-5 text-green-500" />
            <div>
              <p className={`text-sm ${isDarkMode ? 'text-gray-400' : 'text-gray-600'}`}>
                Total Articles
              </p>
              <p className={`text-lg font-semibold ${isDarkMode ? 'text-white' : 'text-gray-900'}`}>
                {engBlogs.reduce((sum, blog) => sum + blog.practice_problems.length, 0)}
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Search Bar */}
      <div className="mb-6">
        <div className="relative">
          <Search className={`absolute left-4 top-1/2 transform -translate-y-1/2 h-5 w-5 ${
            isDarkMode ? 'text-gray-400' : 'text-gray-500'
          }`} />
          <input
            type="text"
            placeholder="Search companies or articles..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className={`w-full pl-12 pr-4 py-3 rounded-xl border ${
              isDarkMode 
                ? 'bg-gray-800 border-gray-700 text-white placeholder-gray-400 focus:border-blue-500' 
                : 'bg-white border-gray-300 text-gray-900 placeholder-gray-500 focus:border-blue-500'
            } focus:outline-none focus:ring-2 focus:ring-blue-500/20 transition-all`}
          />
        </div>
      </div>

      {/* Engineering Blogs Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {filteredBlogs.map((blog) => (
          <div
            key={blog.id}
            className={`rounded-xl border ${
              isDarkMode ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-200'
            } shadow-sm hover:shadow-lg transition-all duration-300 overflow-hidden`}
          >
            {/* Blog Header */}
            <div className={`p-6 border-b ${isDarkMode ? 'border-gray-700' : 'border-gray-200'}`}>
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-3">
                  <div className="p-2 bg-gradient-to-r from-blue-500 to-purple-600 rounded-lg">
                    <Building2 className="h-6 w-6 text-white" />
                  </div>
                  <div>
                    <h3 className={`text-xl font-bold ${isDarkMode ? 'text-white' : 'text-gray-900'}`}>
                      {blog.name}
                    </h3>
                    <p className={`text-sm ${isDarkMode ? 'text-gray-400' : 'text-gray-600'}`}>
                      {blog.practice_problems.length} articles
                    </p>
                  </div>
                </div>
                <div className="flex items-center space-x-2">
                  <button
                    onClick={() => handleExternalLink(blog.link)}
                    className={`p-2 rounded-lg transition-colors ${
                      isDarkMode 
                        ? 'bg-gray-700 hover:bg-gray-600' 
                        : 'bg-gray-100 hover:bg-gray-200'
                    }`}
                    title="Visit blog"
                  >
                    <Globe className={`h-4 w-4 ${isDarkMode ? 'text-gray-400' : 'text-gray-600'}`} />
                  </button>
                  <button
                    onClick={() => toggleExpanded(blog.id)}
                    className={`px-3 py-1 rounded-lg text-sm font-medium transition-colors ${
                      expandedBlogs.has(blog.id)
                        ? isDarkMode 
                          ? 'bg-blue-900/30 text-blue-400'
                          : 'bg-blue-100 text-blue-700'
                        : isDarkMode
                          ? 'bg-gray-700 text-gray-300 hover:bg-gray-600'
                          : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                    }`}
                  >
                    {expandedBlogs.has(blog.id) ? 'Collapse' : 'Explore'}
                  </button>
                </div>
              </div>
            </div>

            {/* Articles List */}
            {expandedBlogs.has(blog.id) && (
              <div className="p-6 pt-0">
                <div className="space-y-3 max-h-96 overflow-y-auto">
                  {blog.practice_problems.map((problem) => (
                    <div
                      key={problem.id}
                      className={`p-4 rounded-lg border ${
                        isDarkMode ? 'bg-gray-700/50 border-gray-600' : 'bg-gray-50 border-gray-200'
                      } hover:shadow-md transition-all duration-200 group cursor-pointer`}
                      onClick={() => handleExternalLink(problem.external_link)}
                    >
                      <div className="flex items-start justify-between">
                        <div className="flex-1">
                          <h4 className={`text-sm font-medium ${
                            isDarkMode ? 'text-white' : 'text-gray-900'
                          } group-hover:text-blue-600 transition-colors line-clamp-2`}>
                            {problem.title}
                          </h4>
                          <div className="flex items-center space-x-2 mt-2">
                            <span className={`text-xs px-2 py-1 rounded ${
                              isDarkMode ? 'bg-gray-600 text-gray-300' : 'bg-gray-200 text-gray-600'
                            }`}>
                              #{problem.order_idx}
                            </span>
                          </div>
                        </div>
                        <ArrowUpRight className={`h-4 w-4 ${
                          isDarkMode ? 'text-gray-400' : 'text-gray-500'
                        } group-hover:text-blue-600 transition-colors flex-shrink-0 ml-2`} />
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Preview (when collapsed) */}
            {!expandedBlogs.has(blog.id) && (
              <div className="p-6 pt-0">
                <div className="space-y-2">
                  {blog.practice_problems.slice(0, 3).map((problem) => (
                    <div
                      key={problem.id}
                      className={`text-sm ${isDarkMode ? 'text-gray-400' : 'text-gray-600'} truncate`}
                    >
                      • {problem.title}
                    </div>
                  ))}
                  {blog.practice_problems.length > 3 && (
                    <div className={`text-sm ${isDarkMode ? 'text-gray-500' : 'text-gray-500'}`}>
                      ... and {blog.practice_problems.length - 3} more articles
                    </div>
                  )}
                </div>
              </div>
            )}
          </div>
        ))}
      </div>

      {filteredBlogs.length === 0 && searchTerm && (
        <div className="text-center py-12">
          <Search className={`h-12 w-12 mx-auto mb-4 ${isDarkMode ? 'text-gray-600' : 'text-gray-400'}`} />
          <h3 className={`text-lg font-medium mb-2 ${isDarkMode ? 'text-white' : 'text-gray-900'}`}>
            No blogs found
          </h3>
          <p className={`${isDarkMode ? 'text-gray-400' : 'text-gray-600'}`}>
            Try adjusting your search terms
          </p>
        </div>
      )}
    </div>
  );
};

export default EngBlogs;
