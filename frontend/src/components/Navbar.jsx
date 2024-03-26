import React from 'react';
import { Link } from 'react-router-dom';

function Navbar() {
  return (
    <nav>
      <ul>
        <li><Link to="/">Home</Link></li>
        <li><Link to="/add-song">Add Song</Link></li>
        <li><Link to="/song-list">Song List</Link></li>
      </ul>
    </nav>
  );
}

export default Navbar;
