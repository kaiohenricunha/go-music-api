import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './authContext';
import LoginForm from './components/LoginForm';
import Dashboard from './components/Dashboard';
import RegistrationForm from './components/RegistrationForm';
import Sidebar from './components/Sidebar';
import IndexPage from './components/IndexPage';

function App() {
  return (
    <AuthProvider>
      <Router>
        <Sidebar /> {/* This ensures the sidebar is shown on every page */}
        <Routes>
          <Route path="/" element={<IndexPage />} />
          <Route path="/login" element={<LoginForm />} />
          <Route path="/register" element={<RegistrationForm />} />
          <Route path="/dashboard" element={<PrivateRoute><Dashboard /></PrivateRoute>} />
        </Routes>
      </Router>
    </AuthProvider>
  );
}

function PrivateRoute({ children }) {
  const { isAuthenticated } = useAuth();

  return isAuthenticated ? children : <Navigate to="/login" />;
}

export default App;
