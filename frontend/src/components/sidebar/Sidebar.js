import React from 'react';
import { useLocation, Link } from 'react-router-dom';
import { useAuth } from '../../authContext';
import './Sidebar.css';

const Sidebar = () => {
  const { isAuthenticated, logout } = useAuth();
  const location = useLocation();

  const guestLinks = [
    { path: "/", name: "Home" },
    { path: "/registration", name: "Register" },
    { path: "/login", name: "Login" }
  ];

  const userLinks = [
    { path: "/logout", name: "Logout", action: logout }
  ];

  const commonLinks = isAuthenticated ? userLinks : guestLinks;

  return (
    <div className="sidebar">
      <nav>
        <ul>
          {commonLinks.map((link) => (
            link.path !== location.pathname && (
              <li key={link.name}>
                {link.path === "/logout" ? (
                  <button onClick={link.action}>{link.name}</button>
                ) : (
                  <Link to={link.path}>{link.name}</Link>
                )}
              </li>
            )
          ))}
        </ul>
      </nav>
    </div>
  );
};

export default Sidebar;
