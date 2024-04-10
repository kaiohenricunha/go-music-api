import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { useAuth } from '../authContext';
import Home from '../components/Home';
import RedirectPage from '../components/RedirectPage';
import RegistrationForm from '../components/registration/RegistrationForm';
import LoginForm from '../components/login/LoginForm';
import Dashboard from '../components/dashboard/Dashboard';
import NotFoundPage from '../components/NotFoundPage';

const AppRouter = () => {
    const { isAuthenticated } = useAuth();
  
    return (
      <BrowserRouter>
        <div className="main">
          <Routes>
            <Route path="/" element={<Home />} exact /> {/* Update syntax for element */}
            <Route path="/redirect" element={<RedirectPage />} />
            <Route path="/registration" element={<RegistrationForm />} />
            <Route path="/login" element={<LoginForm />} />
            {/* Protected Dashboard route with authentication check */}
            <Route
              path="/dashboard"
              element={isAuthenticated ? <Dashboard /> : <LoginForm />}
            />
            <Route path="*" element={<NotFoundPage />} /> {/* Use path="*" for not found */}
          </Routes>
        </div>
      </BrowserRouter>
    );
  };
  
  export default AppRouter;
  