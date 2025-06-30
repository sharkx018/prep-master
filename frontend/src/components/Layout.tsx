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
      <header className="bg-gradient-to-r from-indigo-600 to-purple-600 shadow-lg flex-shrink-0">
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

      <div className="flex flex-1">
        {/* Sidebar */}
        <nav className="w-64 bg-white shadow-lg border-r border-gray-100 flex flex-col">
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
          <div className="flex-shrink-0 p-4 border-t border-gray-100">
            <div className="p-4 bg-gradient-to-r from-gray-50 to-gray-100 rounded-lg">
              <p className="text-xs text-gray-600 text-center">
                <span className="font-semibold text-indigo-600">Your data is protected with industry-standard encryption</span>
              </p>
            </div>
          </div>
        </nav>

        {/* Main content */}
        <main className="flex-1 p-8 overflow-y-auto">
          <div className="max-w-7xl mx-auto">
            {children}
          </div>
        </main>
      </div>
    </div>
  );
};

export default Layout; 