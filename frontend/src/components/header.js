import { Link } from 'react-router-dom';
import React from 'react';

function Header() {
  return (
    <header>
      {/* Other header content */}
      <nav>
        <ul>
          <li><Link to="/">Home</Link></li>
          <li><Link to="/register">Register</Link></li>
          <li><Link to="/login">Login</Link></li>
        </ul>
      </nav>
    </header>
  );
}

export default Header;
