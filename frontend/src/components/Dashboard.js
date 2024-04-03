// src/components/Dashboard.js
import React, { useState, useEffect } from 'react';
import { useAuth } from '../authContext'; // Adjust the import path as necessary

const Dashboard = () => {
  const { logout } = useAuth();
  const [playlists, setPlaylists] = useState([]);
  const [spotifySearch, setSpotifySearch] = useState('');
  const [searchResults, setSearchResults] = useState([]);
  const [otherPlaylists, setOtherPlaylists] = useState([]);

  useEffect(() => {
    // Placeholder function, replace with actual fetch call to get user's playlists
    const fetchMyPlaylists = async () => {
      // Implement fetching of user's playlists
    };
    fetchMyPlaylists();
  }, []);

  const handleSpotifySearch = async (e) => {
    e.preventDefault();
    if (!spotifySearch.trim()) return; // Avoid empty search queries
    
    try {
      const response = await fetch(`http://localhost:8081/api/v1/songs/search?query=${encodeURIComponent(spotifySearch)}`, {
        method: 'GET',
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`, // Ensure the request is authenticated
        },
      });
      
      if (response.ok) {
        const data = await response.json();
        setSearchResults(data.songs); // Assuming the backend returns an array of songs
      } else {
        console.error("Failed to fetch search results");
        setSearchResults([]);
      }
    } catch (error) {
      console.error("Error searching Spotify:", error);
      setSearchResults([]);
    }
  };

  useEffect(() => {
    // Placeholder function, replace with actual fetch call to get other users' playlists
    const fetchOtherPlaylists = async () => {
      // Implement fetching of other users' playlists
    };
    fetchOtherPlaylists();
  }, []);

  return (
    <div>
      <h1>Dashboard</h1>
      <button onClick={logout}>Logout</button>
      
      <section>
        <h2>My Playlists</h2>
        <ul>
          {playlists.map((playlist, index) => (
            <li key={index}>{playlist.name}</li>
          ))}
        </ul>
      </section>

      <section>
        <h2>Search Spotify</h2>
        <form onSubmit={handleSpotifySearch}>
          <input
            type="text"
            value={spotifySearch}
            onChange={(e) => setSpotifySearch(e.target.value)}
            placeholder="Search songs on Spotify"
          />
          <button type="submit">Search</button>
        </form>
        <ul>
          {searchResults.map((result, index) => (
            <li key={index}>{result.name} by {result.artist}</li>
          ))}
        </ul>
      </section>

      <section>
        <h2>Explore Playlists</h2>
        <ul>
          {otherPlaylists.map((playlist, index) => (
            <li key={index}>{playlist.name} - {playlist.rating}</li>
          ))}
        </ul>
      </section>
    </div>
  );
};

export default Dashboard;
