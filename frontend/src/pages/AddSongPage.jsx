import React, { useState } from 'react';
import axios from 'axios';

function AddSongPage() {
  // State to store input values
  const [song, setSong] = useState({
    name: '',
    artist: '',
  });

  // Handles input changes
  const handleChange = (e) => {
    const { name, value } = e.target;
    setSong((prevSong) => ({
      ...prevSong,
      [name]: value,
    }));
  };

  // Handles form submission
  const handleSubmit = (e) => {
    e.preventDefault();

    // apiBaseUrl: the URL provided by the `minikube service` command
    const apiBaseUrl = 'http://127.0.0.1:59643';
    axios.post(`${apiBaseUrl}/api/songs`, song)
      .then(response => {
        // Handle success, maybe clear form or show a success message
        console.log(response.data);
        // Optionally reset the form here if needed
        setSong({ name: '', artist: '' });
      })
      .catch(error => {
        // Handle error, maybe show an error message
        console.error("There was an error adding the song:", error);
      });
  };

  return (
    <div className="container">
      <h2>Add a New Song</h2>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="name">Song Name:</label>
          <input
            type="text"
            className="form-control"
            id="name"
            name="name"
            value={song.name}
            onChange={handleChange}
            required
          />
        </div>
        <div className="form-group">
          <label htmlFor="artist">Artist:</label>
          <input
            type="text"
            className="form-control"
            id="artist"
            name="artist"
            value={song.artist}
            onChange={handleChange}
            required
          />
        </div>
        <button type="submit" className="submit-button">
          Add Song
        </button>
      </form>
    </div>
  );
}

export default AddSongPage;
