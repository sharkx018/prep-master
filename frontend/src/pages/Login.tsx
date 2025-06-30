import React, { useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { useTheme } from '../contexts/ThemeContext';
import { Lock, User, Eye, EyeOff, Brain, Sparkles, Code2, Github, Linkedin, Moon, Sun } from 'lucide-react';

const Login: React.FC = () => {
  const { isDarkMode, toggleTheme } = useTheme();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const { login } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError('');

    const success = await login(username, password);
    
    if (!success) {
      setError('Invalid username or password');
    }
    
    setIsLoading(false);
  };

  return (
    <div className={`min-h-screen flex flex-col transition-colors duration-300 ${
      isDarkMode 
        ? 'bg-gradient-to-br from-gray-900 to-gray-800' 
        : 'bg-gradient-to-br from-gray-50 to-gray-100'
    }`}>
      {/* Header matching the main app */}
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
            </div>
          </div>
        </div>
      </header>

      {/* Main login content */}
      <div className="flex-1 flex items-center justify-center p-8 relative">
        <div className="max-w-md w-full space-y-8">
          <div className={`rounded-2xl shadow-xl p-8 border transition-colors duration-300 ${
            isDarkMode 
              ? 'bg-gray-800 border-gray-700' 
              : 'bg-white border-gray-100'
          }`}>
            <div className="text-center">
              <div className="mx-auto h-16 w-16 bg-gradient-to-r from-indigo-600 to-purple-600 rounded-full flex items-center justify-center shadow-lg relative">
                <Lock className="h-8 w-8 text-white" />
                <Sparkles className="h-4 w-4 text-yellow-300 absolute -top-1 -right-1" />
              </div>
              <h2 className={`mt-6 text-3xl font-extrabold bg-gradient-to-r bg-clip-text text-transparent ${
                isDarkMode 
                  ? 'from-gray-200 to-gray-400' 
                  : 'from-gray-700 to-gray-900'
              }`}>
                Welcome Back
              </h2>
              <p className={`mt-2 text-sm ${isDarkMode ? 'text-gray-400' : 'text-gray-600'}`}>
                Sign in to access your <span className="font-semibold text-indigo-600">PrepMaster Pro</span> dashboard
              </p>
            </div>
            
            <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
              {error && (
                <div className={`border px-4 py-3 rounded-lg text-sm flex items-center ${
                  isDarkMode 
                    ? 'bg-red-900/20 border-red-800 text-red-300' 
                    : 'bg-red-50 border-red-200 text-red-600'
                }`}>
                  <div className="flex-shrink-0 mr-2">
                    <div className="h-2 w-2 bg-red-400 rounded-full"></div>
                  </div>
                  {error}
                </div>
              )}
              
              <div className="space-y-4">
                <div>
                  <label htmlFor="username" className={`block text-sm font-medium mb-2 ${
                    isDarkMode ? 'text-gray-300' : 'text-gray-700'
                  }`}>
                    Username
                  </label>
                  <div className="relative">
                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                      <User className="h-5 w-5 text-gray-400" />
                    </div>
                    <input
                      id="username"
                      name="username"
                      type="text"
                      required
                      className={`block w-full pl-10 pr-3 py-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200 shadow-sm ${
                        isDarkMode 
                          ? 'bg-gray-700 border-gray-600 text-gray-100 placeholder-gray-400' 
                          : 'bg-white border-gray-300 text-gray-900'
                      }`}
                      placeholder="Enter your username"
                      value={username}
                      onChange={(e) => setUsername(e.target.value)}
                    />
                  </div>
                </div>
                
                <div>
                  <label htmlFor="password" className={`block text-sm font-medium mb-2 ${
                    isDarkMode ? 'text-gray-300' : 'text-gray-700'
                  }`}>
                    Password
                  </label>
                  <div className="relative">
                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                      <Lock className="h-5 w-5 text-gray-400" />
                    </div>
                    <input
                      id="password"
                      name="password"
                      type={showPassword ? 'text' : 'password'}
                      required
                      className={`block w-full pl-10 pr-10 py-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200 shadow-sm ${
                        isDarkMode 
                          ? 'bg-gray-700 border-gray-600 text-gray-100 placeholder-gray-400' 
                          : 'bg-white border-gray-300 text-gray-900'
                      }`}
                      placeholder="Enter your password"
                      value={password}
                      onChange={(e) => setPassword(e.target.value)}
                    />
                    <button
                      type="button"
                      className="absolute inset-y-0 right-0 pr-3 flex items-center"
                      onClick={() => setShowPassword(!showPassword)}
                    >
                      {showPassword ? (
                        <EyeOff className="h-5 w-5 text-gray-400 hover:text-gray-600 transition-colors" />
                      ) : (
                        <Eye className="h-5 w-5 text-gray-400 hover:text-gray-600 transition-colors" />
                      )}
                    </button>
                  </div>
                </div>
              </div>

              <div>
                <button
                  type="submit"
                  disabled={isLoading}
                  className="group relative w-full flex justify-center py-3 px-4 border border-transparent text-sm font-medium rounded-lg text-white bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 shadow-lg hover:shadow-xl transform hover:-translate-y-0.5"
                >
                  {isLoading ? (
                    <div className="flex items-center">
                      <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                      Signing in...
                    </div>
                  ) : (
                    <div className="flex items-center">
                      <Brain className="h-4 w-4 mr-2" />
                      Sign in to PrepMaster Pro
                    </div>
                  )}
                </button>
              </div>
            </form>
            
            <div className="mt-6">
              <div className={`p-4 rounded-lg border ${
                isDarkMode 
                  ? 'bg-gradient-to-r from-indigo-900/30 to-purple-900/30 border-indigo-800' 
                  : 'bg-gradient-to-r from-indigo-50 to-purple-50 border-indigo-100'
              }`}>
                <div className="flex items-center justify-center space-x-2 mb-2">
                  <div className="h-2 w-2 bg-indigo-400 rounded-full"></div>
                  <span className={`text-xs font-semibold ${
                    isDarkMode ? 'text-indigo-300' : 'text-indigo-900'
                  }`}>Enterprise Security</span>
                  <div className="h-2 w-2 bg-purple-400 rounded-full"></div>
                </div>
                <p className={`text-xs text-center ${
                  isDarkMode ? 'text-indigo-400' : 'text-indigo-700'
                }`}>
                  Your data is protected with industry-standard encryption
                </p>
              </div>
            </div>
          </div>
        </div>
        
        {/* Stable branding */}
        <div className="absolute bottom-6 left-1/2 transform -translate-x-1/2">
          <a 
            href="https://www.instagram.com/vrma018_/" 
            target="_blank" 
            rel="noopener noreferrer"
            className={`flex items-center space-x-2 px-4 py-2 border rounded-lg transition-all duration-300 hover:scale-105 ${
              isDarkMode
                ? 'bg-gradient-to-r from-purple-900/30 to-pink-900/30 border-purple-700'
                : 'bg-gradient-to-r from-purple-50 to-pink-50 border-purple-200'
            }`}
          >
            <div className="w-2 h-2 bg-gradient-to-r from-purple-500 to-pink-500 rounded-full"></div>
            <span className="text-lg font-mono font-medium bg-gradient-to-r from-purple-600 to-pink-600 bg-clip-text text-transparent tracking-widest">
              @vrma018_
            </span>
          </a>
        </div>
      </div>
    </div>
  );
};

export default Login; 