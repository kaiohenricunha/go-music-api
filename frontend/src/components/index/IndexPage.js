import React from 'react';
import { Link } from 'react-router-dom';
import './IndexPage.css';

function IndexPage() {
  return (
    <div className="container">
        <h1 id="title" className="text-center">Musilist</h1>
        <p id="description" className="description text-center">
          Dive into the world of Musilist, where music curation meets community wisdom.
        </p>
      <div className="text-center">
        <Link to="/register" className="enter-now-link">Enter Now</Link>
      </div>
    </div>
  );
}

export default IndexPage;
