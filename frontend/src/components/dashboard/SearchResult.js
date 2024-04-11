import React, { useState } from 'react';
import { useSelector } from 'react-redux';
import SongsList from './SongsList';

const SearchResult = () => {
  const [selectedCategory, setSelectedCategory] = useState('songs'); // Default category
  const songs = useSelector(state => state.songs.items);  // Fetching items from songs state

  return (
    <div>
      <div className="search-buttons">
        <button
          className={`btn ${selectedCategory === 'songs' ? 'active' : ''}`}
          onClick={() => setSelectedCategory('songs')}
        >
          Songs
        </button>
      </div>
      {selectedCategory === 'songs' && <SongsList songs={songs || []} />}
   </div>
  );
};

export default SearchResult;
