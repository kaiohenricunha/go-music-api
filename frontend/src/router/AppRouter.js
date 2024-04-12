import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { useAuth } from '../authContext';
import Sidebar from '../components/sidebar/Sidebar';
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
        <div className="app-layout">
          <Sidebar />  {/* This will show the sidebar on all pages */}
          <div className="main-content">
            <Routes>
              <Route path="/" element={<Home />} />
              <Route path="/redirect" element={<RedirectPage />} />
              <Route path="/registration" element={<RegistrationForm />} />
              <Route path="/login" element={<LoginForm />} />
              <Route
                path="/dashboard"
                element={isAuthenticated ? <Dashboard /> : <LoginForm />}
              />
              <Route path="*" element={<NotFoundPage />} />
            </Routes>
          </div>
        </div>
      </BrowserRouter>
    );
};

export default AppRouter;
