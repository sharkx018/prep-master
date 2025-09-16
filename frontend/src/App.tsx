import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate, Outlet } from 'react-router-dom';
import { AuthProvider, useAuth } from './contexts/AuthContext';
import { ThemeProvider } from './contexts/ThemeContext';
import ProtectedRoute from './components/ProtectedRoute';
import Layout from './components/Layout';
// import Dashboard from './pages/Dashboard';
import Items from './pages/Items';
import Practice from './pages/Practice';
import AddItem from './pages/AddItem';
import Stats from './pages/Stats';
import Contest from './pages/Contest';
import EngBlogs from './pages/EngBlogs';
import Login from './pages/Login';

// Component to protect admin-only routes
const AdminRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAdmin } = useAuth();
  
  if (!isAdmin) {
    return <Navigate to="/practice" replace />;
  }
  
  return <>{children}</>;
};

function App() {
  return (
    <ThemeProvider>
      <AuthProvider>
        <Router>
          <Routes>
            <Route path="/login" element={<Login />} />
            <Route path="/" element={<ProtectedRoute><Layout /></ProtectedRoute>}>
              <Route index element={<Navigate to="practice" replace />} />
              {/* <Route path="dashboard" element={<Dashboard />} /> */}
              <Route path="items" element={<Items />} />
              <Route path="practice" element={<Practice />} />
              <Route path="contests" element={<Contest />} />
              <Route path="eng-blogs" element={<EngBlogs />} />
              <Route path="add-item" element={
                <AdminRoute>
                  <AddItem />
                </AdminRoute>
              } />
              <Route path="stats" element={<Stats />} />
              {/* Catch-all route for unknown paths */}
              <Route path="*" element={<Navigate to="/practice" replace />} />
            </Route>
          </Routes>
        </Router>
      </AuthProvider>
    </ThemeProvider>
  );
}

export default App;
