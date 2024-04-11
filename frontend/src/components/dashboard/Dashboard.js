import React, { useState } from 'react';
import { useAuth } from '../../authContext';
import SearchForm from './SearchForm';
import SearchResult from './SearchResult';
import { initiateGetSongs } from '../../actions/result'; 
import { useDispatch } from 'react-redux'; 
import axios from 'axios';

const Dashboard = () => {
  // No more need for local songs state
  const { token } = useAuth(); 
  const dispatch = useDispatch(); 

  const handleSearch = (searchTerm) => {
    console.log('searchTerm:', searchTerm); 
    // Dispatch the action to fetch songs asynchronously using Redux
    dispatch(initiateGetSongs(searchTerm, token)); // Pass the search term and token
  };  

  return (
    <div>
      <SearchForm onSearch={handleSearch} />
      <SearchResult /> {/* Songs will be taken from Redux */}
    </div>
  );
};

export default Dashboard;