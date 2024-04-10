import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate, useLocation } from 'react-router-dom';
import { AuthProvider, useAuth } from './authContext';
import LoginForm from './components/login/LoginForm';
import Sidebar from './components/sidebar/Sidebar';
import IndexPage from './components/index/IndexPage';
import RegistrationForm from './components/registration/RegistrationForm';
import Dashboard from './components/dashboard/Dashboard';

const AppWrapper = () => {
  return (
    <AuthProvider>
      <Router>
        <App />
      </Router>
    </AuthProvider>
  );
};

const App = () => {
  const location = useLocation();
  const { isAuthenticated } = useAuth();

  return (
    <>
      {location.pathname !== "/dashboard" && <Sidebar />}
      <Routes>
        <Route path="/" element={<IndexPage />} />
        <Route path="/login" element={<LoginForm />} />
        <Route path="/register" element={<RegistrationForm />} />
        <Route path="/dashboard" element={<PrivateRoute><Dashboard /></PrivateRoute>} />
      </Routes>
    </>
  );
};

function PrivateRoute({ children }) {
  const auth = useAuth();

  return auth.isAuthenticated ? children : <Navigate to="/login" />;
};

export default AppWrapper;
