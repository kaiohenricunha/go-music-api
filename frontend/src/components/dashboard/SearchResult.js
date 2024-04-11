import React from 'react';
import { useSelector } from 'react-redux';
import SongsList from './SongsList';

const SearchResult = () => {
  const songs = useSelector(state => state.songs.items);  // Fetching items from songs state

  return (
    <div>
      <SongsList songs={songs || []} />
    </div>
  );
};

export default SearchResult;
