import React, { useState, useEffect, useCallback } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { useTheme } from '../contexts/ThemeContext';
import HeaderStreakWidget from './HeaderStreakWidget';
import { statsApi, Stats } from '../services/api';
import { 
  LayoutDashboard, 
  List, 
  BookOpen, 
  Plus, 
  BarChart3,
  Brain,
  Sparkles,
  Github,
  Linkedin,
  Code2,
  LogOut,
  Moon,
  Sun,
  Star
} from 'lucide-react';

// CompletionStars component
const CompletionStars: React.FC<{ completionCycles: number }> = ({ completionCycles }) => {
  const getRankInfo = (count: number) => {
    if (count === 0) return { title: 'Newbie', color: 'text-gray-400' };
    if (count === 1) return { title: 'Pupil', color: 'text-green-400' };
    if (count === 2) return { title: 'Specialist', color: 'text-cyan-400' };
    if (count === 3) return { title: 'Expert', color: 'text-blue-400' };
    if (count === 4) return { title: 'Candidate Master', color: 'text-violet-400' };
    if (count === 5) return { title: 'Master', color: 'text-orange-400' };
    if (count === 6) return { title: 'International Master', color: 'text-orange-400' };
    if (count === 7) return { title: 'Grandmaster', color: 'text-red-400' };
    if (count >= 10 && count < 15) return { title: 'International Grandmaster', color: 'text-red-400' };
    if (count >= 15) return { title: 'Legendary Grandmaster', color: 'text-red-400', isLegendary: true };
    // For counts 8-9, use Grandmaster
    return { title: 'Grandmaster', color: 'text-red-400' };
  };

  const rankInfo = getRankInfo(completionCycles);

  const renderTitle = () => {
    if (rankInfo.isLegendary) {
      return (
        <span className="text-sm font-bold text-white drop-shadow-lg">
          <span className="text-white">L</span>
          <span className="text-white">egendary Grandmaster</span>
        </span>
      );
    }
    return (
      <span className="text-sm font-bold text-white drop-shadow-lg">
        {rankInfo.title}
      </span>
    );
  };

  return (
    <div className="flex items-center" title={`${completionCycles} Completion Cycles - ${rankInfo.title}`}>
      <div className="flex items-center bg-white/10 px-3 py-1.5 rounded-lg border border-white/20 shadow-lg">
        {renderTitle()}
      </div>
    </div>
  );
};

interface LayoutProps {
  children: React.ReactNode;
}

// Avatar component
const Avatar: React.FC<{ user: { name: string; avatar?: string } }> = ({ user }) => {
  const getInitials = (name: string) => {
    return name
      .split(' ')
      .map(word => word.charAt(0))
      .join('')
      .toUpperCase()
      .slice(0, 2);
  };

  if (user.avatar) {
    return (
      <div className="relative">
        <img
          src={user.avatar}
          alt={`${user.name}'s avatar`}
          className="w-8 h-8 rounded-full object-cover border-2 border-white/30 hover:border-white/50 transition-all duration-200"
          onError={(e) => {
            // Fallback to initials if image fails to load
            const target = e.target as HTMLImageElement;
            target.style.display = 'none';
            const parent = target.parentElement;
            if (parent) {
              const fallback = document.createElement('div');
              fallback.className = 'w-8 h-8 rounded-full bg-gradient-to-br from-indigo-400 to-purple-400 flex items-center justify-center text-white text-sm font-semibold border-2 border-white/30 hover:border-white/50 transition-all duration-200';
              fallback.textContent = getInitials(user.name);
              parent.appendChild(fallback);
            }
          }}
        />
        <div className="absolute -bottom-1 -right-1 w-3 h-3 bg-green-400 rounded-full border-2 border-white"></div>
      </div>
    );
  }

  return (
    <div className="relative">
      <div className="w-8 h-8 rounded-full bg-gradient-to-br from-indigo-400 to-purple-400 flex items-center justify-center text-white text-sm font-semibold border-2 border-white/30 hover:border-white/50 transition-all duration-200">
        {getInitials(user.name)}
      </div>
      <div className="absolute -bottom-1 -right-1 w-3 h-3 bg-green-400 rounded-full border-2 border-white"></div>
    </div>
  );
};

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const location = useLocation();
  const { user, logout, isAdmin } = useAuth();
  const { isDarkMode, toggleTheme } = useTheme();
  const [stats, setStats] = useState<Stats | null>(null);

  const fetchStats = useCallback(async () => {
    try {
      const data = await statsApi.getStats();
      setStats(data);
    } catch (err) {
      console.error('Failed to fetch stats:', err);
    }
  }, []);

  useEffect(() => {
    if (user) {
      fetchStats();
      
      const interval = setInterval(() => fetchStats(), 30000);
      
      const handleItemCompletion = () => {
        setTimeout(() => fetchStats(), 500);
      };

      window.addEventListener('itemCompleted', handleItemCompletion);
      
      return () => {
        clearInterval(interval);
        window.removeEventListener('itemCompleted', handleItemCompletion);
      };
    }
  }, [user, fetchStats]);

  const navigation = [
    { name: 'Dashboard', href: '/dashboard', icon: LayoutDashboard },
    { name: 'Study', href: '/study', icon: BookOpen },
    { name: 'Items', href: '/items', icon: List },
    ...(isAdmin ? [{ name: 'Add Item', href: '/add-item', icon: Plus }] : []),
    { name: 'Statistics', href: '/stats', icon: BarChart3 },
  ];

  const isActive = (path: string) => location.pathname === path;

  const handleLogout = () => {
    logout();
  };

  return (
    <div className={`min-h-screen flex flex-col transition-colors duration-300 ${
      isDarkMode 
        ? 'bg-gradient-to-br from-gray-900 to-gray-800' 
        : 'bg-gradient-to-br from-gray-50 to-gray-100'
    }`}>
      {/* Header */}
      <header className="bg-gradient-to-r from-indigo-600 to-purple-600 shadow-lg flex-shrink-0 fixed top-0 left-0 right-0 z-20">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center">
              <div className="relative">
                <Brain className="h-10 w-10 text-white" />
                <Sparkles className="h-4 w-4 text-yellow-300 absolute -top-1 -right-1" />
              </div>
              <div className="ml-4">
                <h1 className="text-xl font-bold text-white">PrepMaster Pro</h1>
              </div>
            </div>
            
            {/* Center - Streak Widget */}
            <div className="flex-1 flex justify-center">
              <HeaderStreakWidget />
            </div>
            
            <div className="flex items-center space-x-4">
              <a
                href="https://github.com/sharkx018"
                target="_blank"
                rel="noopener noreferrer"
                className="text-white/80 hover:text-white transition-colors"
                title="GitHub"
              >
                <Github className="h-5 w-5" />
              </a>
              <a
                href="https://www.linkedin.com/in/mukul-verma-03151a139/"
                target="_blank"
                rel="noopener noreferrer"
                className="text-white/80 hover:text-white transition-colors"
                title="LinkedIn"
              >
                <Linkedin className="h-5 w-5" />
              </a>
              <button
                onClick={toggleTheme}
                className="p-2 text-white/80 hover:text-white hover:bg-white/20 rounded-full transition-all duration-300 hover:scale-110"
                title={isDarkMode ? "Switch to Light Mode" : "Switch to Dark Mode"}
              >
                {isDarkMode ? <Sun className="h-5 w-5" /> : <Moon className="h-5 w-5" />}
              </button>
              <div className="flex items-center text-white bg-white/20 px-3 py-1 rounded-full">
                <Code2 className="h-4 w-4 mr-1" />
                <span className="text-sm font-medium">v2.0</span>
              </div>
              {user && (
                <div className="flex items-center space-x-3 text-white">
                  <div className="flex items-center bg-white/20 px-3 py-2 rounded-full hover:bg-white/30 transition-all duration-200">
                    <Avatar user={user} />
                    <span className="text-sm font-bold ml-2">
                      {(() => {
                        const count = stats?.completed_all_count || 0;
                        if (count >= 15) {
                          // Legendary Grandmaster: First letter black, rest red
                          return (
                            <>
                              <span className="text-black">{user.name.charAt(0)}</span>
                              <span className="text-red-500">{user.name.slice(1)}</span>
                            </>
                          );
                        }
                        // Regular color for other ranks
                        const colorClass = (() => {
                          if (count === 0) return 'text-gray-300';
                          if (count === 1) return 'text-green-300';
                          if (count === 2) return 'text-cyan-300';
                          if (count === 3) return 'text-blue-400';
                          if (count === 4) return 'text-purple-300';
                          if (count === 5) return 'text-orange-300';
                          if (count === 6) return 'text-orange-300';
                          if (count >= 7) return 'text-red-400';
                          return 'text-white';
                        })();
                        return <span className={colorClass}>{user.name}</span>;
                      })()}
                    </span>
                  </div>
                  <CompletionStars completionCycles={stats?.completed_all_count || 0} />
                  <button
                    onClick={handleLogout}
                    className="p-2 text-white/80 hover:text-white hover:bg-white/20 rounded-full transition-colors"
                    title="Logout"
                  >
                    <LogOut className="h-4 w-4" />
                  </button>
                </div>
              )}
            </div>
          </div>
        </div>
      </header>

      <div className="flex flex-1 relative">
        {/* Sidebar */}
        <nav className={`w-64 shadow-lg border-r flex flex-col fixed top-16 bottom-0 left-0 z-10 transition-colors duration-300 ${
          isDarkMode 
            ? 'bg-gray-800 border-gray-700' 
            : 'bg-white border-gray-100'
        }`}>
          <div className="flex-1 p-4 overflow-y-auto">
            <div className={`mb-6 p-4 rounded-lg border transition-colors duration-300 ${
              isDarkMode 
                ? 'bg-gradient-to-r from-indigo-900/50 to-purple-900/50 border-indigo-800' 
                : 'bg-gradient-to-r from-indigo-50 to-purple-50 border-indigo-100'
            }`}>
              <h3 className={`text-sm font-semibold mb-1 transition-colors duration-300 ${
                isDarkMode ? 'text-indigo-300' : 'text-indigo-900'
              }`}>Welcome back!</h3>
              <p className={`text-xs transition-colors duration-300 ${
                isDarkMode ? 'text-indigo-400' : 'text-indigo-700'
              }`}>Keep crushing 🚀</p>
            </div>
            <ul className="space-y-1">
              {navigation.map((item) => {
                const Icon = item.icon;
                return (
                  <li key={item.name}>
                    <Link
                      to={item.href}
                      className={`flex items-center px-4 py-3 text-sm font-medium rounded-lg transition-all duration-200 ${
                        isActive(item.href)
                          ? 'bg-gradient-to-r from-indigo-500 to-purple-500 text-white shadow-md transform scale-105'
                          : isDarkMode
                            ? 'text-gray-300 hover:bg-gray-700 hover:text-white hover:translate-x-1'
                            : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900 hover:translate-x-1'
                      }`}
                    >
                      <Icon className={`mr-3 h-5 w-5 ${
                        isActive(item.href) 
                          ? 'text-white' 
                          : isDarkMode 
                            ? 'text-gray-400' 
                            : 'text-gray-600'
                      }`} />
                      {item.name}
                    </Link>
                  </li>
                );
              })}
            </ul>
          </div>
          <div className="flex-shrink-0 p-6">
            <div className="flex items-center justify-center">
              <a 
                href="https://www.instagram.com/vrma018_/" 
                target="_blank" 
                rel="noopener noreferrer"
                className={`group relative px-4 py-3 border rounded-xl transition-all duration-300 hover:shadow-lg hover:-translate-y-0.5 ${
                  isDarkMode
                    ? 'bg-gradient-to-r from-purple-900/30 to-pink-900/30 border-purple-700 hover:from-purple-900/50 hover:to-pink-900/50 hover:border-purple-600'
                    : 'bg-gradient-to-r from-purple-50 to-pink-50 border-purple-200 hover:from-purple-100 hover:to-pink-100 hover:border-purple-300'
                }`}
              >
                <div className="flex items-center space-x-2">
                  <div className="w-2 h-2 bg-gradient-to-r from-purple-500 to-pink-500 rounded-full group-hover:animate-pulse"></div>
                  <span className="text-lg font-mono font-semibold bg-gradient-to-r from-purple-600 to-pink-600 bg-clip-text text-transparent tracking-wider">
                    @vrma018_
                  </span>
                </div>
                <div className={`absolute inset-0 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-300 ${
                  isDarkMode
                    ? 'bg-gradient-to-r from-purple-400/20 to-pink-400/20'
                    : 'bg-gradient-to-r from-purple-400/10 to-pink-400/10'
                }`}></div>
              </a>
            </div>
          </div>
        </nav>

        {/* Main content */}
        <main className="flex-1 p-8 overflow-y-auto ml-64 mt-16">
          <div className="max-w-7xl mx-auto">
            {children}
          </div>
        </main>
      </div>
    </div>
  );
};

export default Layout; 