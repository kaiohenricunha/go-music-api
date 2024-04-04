import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './authContext';
import LoginForm from './components/login/LoginForm';
import RegistrationForm from './components/registration/RegistrationForm';
import Sidebar from './components/sidebar/Sidebar';
import IndexPage from './components/index/IndexPage';

function App() {
  return (
    <AuthProvider>
      <Router>
        <Sidebar /> {/* This ensures the sidebar is shown on every page */}
        <Routes>
          <Route path="/" element={<IndexPage />} />
          <Route path="/login" element={<LoginForm />} />
          <Route path="/register" element={<RegistrationForm />} />
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
