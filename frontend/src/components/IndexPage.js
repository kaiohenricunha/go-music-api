import React from 'react';
import './IndexPage.css';
import { Link } from 'react-router-dom';

function IndexPage() {
  return (
    <div className="container">
      {/* Page content */}
      <div className="text-center">
        <Link to="/register"><h3>Enter Now</h3></Link>
      </div>
    </div>
  );
}

export default IndexPage;
