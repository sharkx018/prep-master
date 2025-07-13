import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';

interface User {
  id: number;
  email: string;
  name: string;
  avatar?: string;
  auth_provider: 'email' | 'google' | 'facebook' | 'apple';
  created_at: string;
  last_login_at?: string;
}

interface AuthContextType {
  user: User | null;
  login: (email: string, password: string) => Promise<boolean>;
  register: (email: string, password: string, name: string) => Promise<boolean>;
  loginWithGoogle: () => Promise<boolean>;
  loginWithFacebook: () => Promise<boolean>;
  loginWithApple: () => Promise<boolean>;
  logout: () => void;
  isAuthenticated: boolean;
  isLoading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // Initialize OAuth SDKs
  useEffect(() => {
    // Initialize Facebook SDK
    if (window.FB) {
      window.FB.init({
        appId: process.env.REACT_APP_FACEBOOK_APP_ID || '',
        cookie: true,
        xfbml: true,
        version: 'v18.0'
      });
    }

    // Initialize Apple Sign In
    if (window.AppleID && process.env.REACT_APP_APPLE_CLIENT_ID) {
      window.AppleID.auth.init({
        clientId: process.env.REACT_APP_APPLE_CLIENT_ID,
        scope: 'name email',
        redirectURI: window.location.origin,
        usePopup: true
      });
    }
  }, []);

  // Check for existing token on app load
  useEffect(() => {
    const token = localStorage.getItem('auth_token');
    const userStr = localStorage.getItem('auth_user');
    
    if (token && userStr) {
      try {
        const userData = JSON.parse(userStr);
        setUser(userData);
      } catch (error) {
        console.error('Error parsing user data:', error);
        localStorage.removeItem('auth_token');
        localStorage.removeItem('auth_user');
      }
    }
    setIsLoading(false);
  }, []);

  const login = async (email: string, password: string): Promise<boolean> => {
    try {
      const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';
      const response = await fetch(`${API_BASE_URL}/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      });

      if (response.ok) {
        const data = await response.json();
        
        setUser(data.user);
        localStorage.setItem('auth_token', data.token);
        localStorage.setItem('auth_user', JSON.stringify(data.user));
        
        return true;
      } else {
        const errorData = await response.json();
        console.error('Login error:', errorData.error);
        return false;
      }
    } catch (error) {
      console.error('Login error:', error);
      return false;
    }
  };

  const register = async (email: string, password: string, name: string): Promise<boolean> => {
    try {
      const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';
      const response = await fetch(`${API_BASE_URL}/auth/register`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password, name }),
      });

      if (response.ok) {
        const data = await response.json();
        
        setUser(data.user);
        localStorage.setItem('auth_token', data.token);
        localStorage.setItem('auth_user', JSON.stringify(data.user));
        
        return true;
      } else {
        const errorData = await response.json();
        console.error('Registration error:', errorData.error);
        return false;
      }
    } catch (error) {
      console.error('Registration error:', error);
      return false;
    }
  };

  const loginWithGoogle = async (): Promise<boolean> => {
    try {
      const googleClientId = process.env.REACT_APP_GOOGLE_CLIENT_ID;
      
      if (!googleClientId) {
        console.error('Google Client ID not configured. Please set REACT_APP_GOOGLE_CLIENT_ID in your .env file');
        alert('Google login is not configured. Please contact the administrator.');
        return false;
      }

      // Initialize Google OAuth
      if (!window.google) {
        console.error('Google OAuth SDK not loaded');
        alert('Google OAuth SDK not loaded. Please refresh the page and try again.');
        return false;
      }

      return new Promise((resolve) => {
        window.google.accounts.oauth2.initTokenClient({
          client_id: googleClientId,
          scope: 'email profile',
          callback: async (response: any) => {
            if (response.access_token) {
              const success = await handleOAuthLogin('google', response.access_token);
              resolve(success);
            } else {
              console.error('Google OAuth failed: No access token received');
              resolve(false);
            }
          },
        }).requestAccessToken();
      });
    } catch (error) {
      console.error('Google login error:', error);
      return false;
    }
  };

  const loginWithFacebook = async (): Promise<boolean> => {
    try {
      const facebookAppId = process.env.REACT_APP_FACEBOOK_APP_ID;
      
      if (!facebookAppId) {
        console.error('Facebook App ID not configured. Please set REACT_APP_FACEBOOK_APP_ID in your .env file');
        alert('Facebook login is not configured. Please contact the administrator.');
        return false;
      }

      // Initialize Facebook SDK
      if (!window.FB) {
        console.error('Facebook SDK not loaded');
        alert('Facebook SDK not loaded. Please refresh the page and try again.');
        return false;
      }

      return new Promise((resolve) => {
        window.FB.login((response: any) => {
          if (response.authResponse) {
            handleOAuthLogin('facebook', response.authResponse.accessToken)
              .then(resolve)
              .catch(() => resolve(false));
          } else {
            console.error('Facebook login failed: No auth response received');
            resolve(false);
          }
        }, { scope: 'email' });
      });
    } catch (error) {
      console.error('Facebook login error:', error);
      return false;
    }
  };

  const loginWithApple = async (): Promise<boolean> => {
    try {
      const appleClientId = process.env.REACT_APP_APPLE_CLIENT_ID;
      
      if (!appleClientId) {
        console.error('Apple Client ID not configured. Please set REACT_APP_APPLE_CLIENT_ID in your .env file');
        alert('Apple login is not configured. Please contact the administrator.');
        return false;
      }

      // Initialize Apple Sign In
      if (!window.AppleID) {
        console.error('Apple Sign In SDK not loaded');
        alert('Apple Sign In SDK not loaded. Please refresh the page and try again.');
        return false;
      }

      return new Promise((resolve) => {
        window.AppleID.auth.signIn().then(async (response: any) => {
          if (response.authorization && response.authorization.id_token) {
            const success = await handleOAuthLogin('apple', response.authorization.id_token);
            resolve(success);
          } else {
            console.error('Apple Sign In failed: No authorization token received');
            resolve(false);
          }
        }).catch((error: any) => {
          console.error('Apple Sign In error:', error);
          resolve(false);
        });
      });
    } catch (error) {
      console.error('Apple login error:', error);
      return false;
    }
  };

  const handleOAuthLogin = async (provider: string, accessToken: string): Promise<boolean> => {
    try {
      const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';
      const response = await fetch(`${API_BASE_URL}/auth/oauth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          provider,
          access_token: accessToken,
        }),
      });

      if (response.ok) {
        const data = await response.json();
        
        setUser(data.user);
        localStorage.setItem('auth_token', data.token);
        localStorage.setItem('auth_user', JSON.stringify(data.user));
        
        return true;
      } else {
        const errorData = await response.json();
        console.error('OAuth login error:', errorData.error);
        return false;
      }
    } catch (error) {
      console.error('OAuth login error:', error);
      return false;
    }
  };

  const logout = () => {
    setUser(null);
    localStorage.removeItem('auth_token');
    localStorage.removeItem('auth_user');
  };

  const value: AuthContextType = {
    user,
    login,
    register,
    loginWithGoogle,
    loginWithFacebook,
    loginWithApple,
    logout,
    isAuthenticated: !!user,
    isLoading,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

// Extend Window interface for OAuth SDKs
declare global {
  interface Window {
    google: any;
    FB: any;
    AppleID: any;
  }
} 