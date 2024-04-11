import React, { useState } from 'react';
import { Form, Button } from 'react-bootstrap';

const SearchForm = ({ onSearch }) => {
  const [searchTerm, setSearchTerm] = useState('');
  const [errorMsg, setErrorMsg] = useState('');

  const handleInputChange = (event) => {
    setSearchTerm(event.target.value); // Just assign the value directly
  };

  const handleSearch = (event) => {
    event.preventDefault();
    const trimmedTerm = searchTerm.trim();
    if (trimmedTerm !== '') {
      // Split the searchTerm by " by " to get songName and artistName
      const parts = trimmedTerm.split(' by ');
      if (parts.length !== 2) {
        setErrorMsg('Please enter search in format "SongName by ArtistName"');
        return;
      }
      const songName = parts[0];
      const artistName = parts[1];
      setErrorMsg(''); // Clear any previous error messages
      onSearch(`${songName} by ${artistName}`); // Send as a single string
    } else {
      setErrorMsg('Please enter a search term.');
    }
  };

  return (
    <div>
      <Form onSubmit={handleSearch}>
        {errorMsg && <p className="errorMsg">{errorMsg}</p>}
        <Form.Group controlId="formBasicEmail">
          <Form.Label>Enter a song search</Form.Label>
          <Form.Control
            type="search"
            name="searchTerm"
            value={searchTerm}
            placeholder="e.g., Glass Eyes by Radiohead"
            onChange={handleInputChange}
            autoComplete="off"
          />
        </Form.Group>
        <Button variant="info" type="submit">
          Search
        </Button>
      </Form>
    </div>
  );
};

export default SearchForm;
