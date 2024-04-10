import React from 'react';
import { ListGroup } from 'react-bootstrap';

const SongResults = ({ songs }) => {
  return (
    <ListGroup>
      {songs.map((song, index) => (
        <ListGroup.Item key={index}>
          {song.name} by {song.artist}
        </ListGroup.Item>
      ))}
    </ListGroup>
  );
};

export default SongResults;
