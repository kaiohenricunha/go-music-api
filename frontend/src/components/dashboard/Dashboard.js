import React, { useState } from 'react';
import { useAuth } from '../../authContext';
import SearchForm from './SearchForm';
import SongResults from './SongResults';
import axios from 'axios';

const Dashboard = () => {
  const [songs, setSongs] = useState([]);
  const { token } = useAuth(); // Use the token from the AuthContext

  const handleSearch = async (searchTerm) => {
    try {
      const [trackName, artistName] = searchTerm.split(' by ');
      const response = await axios.get(`${process.env.REACT_APP_GO_BACKEND_BASE_URL}/songs/search`, {
        headers: {
          Authorization: `Bearer ${token}`, // Include the JWT token in the request header
        },
        params: { songName: trackName.trim(), artistName: artistName.trim() },
      });
      setSongs(response.data);
    } catch (error) {
      console.error('Search error:', error);
      setSongs([]);
    }
  };

  return (
    <div>
      <SearchForm onSearch={handleSearch} />
      <SongResults songs={songs} />
    </div>
  );
};

export default Dashboard;
