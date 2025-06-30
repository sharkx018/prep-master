import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
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
  User
} from 'lucide-react';

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const location = useLocation();
  const { user, logout } = useAuth();

  const navigation = [
    { name: 'Dashboard', href: '/dashboard', icon: LayoutDashboard },
    { name: 'Items', href: '/items', icon: List },
    { name: 'Study', href: '/study', icon: BookOpen },
    { name: 'Add Item', href: '/add-item', icon: Plus },
    { name: 'Statistics', href: '/stats', icon: BarChart3 },
  ];

  const isActive = (path: string) => location.pathname === path;

  const handleLogout = () => {
    logout();
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 flex flex-col">
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
              <div className="flex items-center text-white bg-white/20 px-3 py-1 rounded-full">
                <Code2 className="h-4 w-4 mr-1" />
                <span className="text-sm font-medium">v2.0</span>
              </div>
              {user && (
                <div className="flex items-center space-x-2 text-white">
                  <div className="flex items-center bg-white/20 px-3 py-1 rounded-full">
                    <User className="h-4 w-4 mr-2" />
                    <span className="text-sm font-medium">{user.username}</span>
                  </div>
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
        <nav className="w-64 bg-white shadow-lg border-r border-gray-100 flex flex-col fixed top-16 bottom-0 left-0 z-10">
          <div className="flex-1 p-4 overflow-y-auto">
            <div className="mb-6 p-4 bg-gradient-to-r from-indigo-50 to-purple-50 rounded-lg border border-indigo-100">
              <h3 className="text-sm font-semibold text-indigo-900 mb-1">Welcome back!</h3>
              <p className="text-xs text-indigo-700">Keep crushing 🚀</p>
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
                          : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900 hover:translate-x-1'
                      }`}
                    >
                      <Icon className={`mr-3 h-5 w-5 ${isActive(item.href) ? 'text-white' : ''}`} />
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
                className="group relative px-4 py-3 bg-gradient-to-r from-purple-50 to-pink-50 border border-purple-200 rounded-xl hover:from-purple-100 hover:to-pink-100 hover:border-purple-300 transition-all duration-300 hover:shadow-lg hover:-translate-y-0.5"
              >
                <div className="flex items-center space-x-2">
                  <div className="w-2 h-2 bg-gradient-to-r from-purple-500 to-pink-500 rounded-full group-hover:animate-pulse"></div>
                  <span className="text-lg font-mono font-semibold bg-gradient-to-r from-purple-600 to-pink-600 bg-clip-text text-transparent tracking-wider">
                    @vrma018_
                  </span>
                </div>
                <div className="absolute inset-0 bg-gradient-to-r from-purple-400/10 to-pink-400/10 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
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