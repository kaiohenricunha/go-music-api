import React from 'react';
import './IndexPage.css';
import { Link } from 'react-router-dom';

function IndexPage() {
  return (
    <div className="container">
      <header className="header">
        <h1 id="title" className="text-center">Musiverso Annual Song Contest</h1>
        <p id="description" className="description text-center">
          Get your songs in front of #1 songwriters and the music industry.<br />
          Enter for your chance to win $10,000.
        </p>
      </header>

      {/* Adjust the href to use React Router's Link component for SPA navigation */}
      <div className="text-center">
        <Link href="/register"><h3>Enter Now</h3></Link>
      </div>
    </div>
  );
}

export default IndexPage;
