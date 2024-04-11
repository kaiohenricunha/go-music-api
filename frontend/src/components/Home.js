import React, { useEffect } from 'react';
import { Link } from 'react-router-dom';

function Home() {
//   useEffect(() => {
//     document.body.classList.add('bg-image-page3');

//     return () => {
//       document.body.classList.remove('bg-image-page3');
//     };
//   }, []);

  return (
    <div className="container">
      <h1 id="title" className="text-center">Musilist</h1>
      <p id="description" className="description text-center">
        Dive into the world of Musilist, where music curation meets community wisdom.
      </p>
      <div className="text-center">
        <Link to="/registration">Sign Up</Link> {/* Use Link component for SPA */}
      </div>
    </div>
  );
}

export default Home;
