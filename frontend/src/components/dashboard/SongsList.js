import React from 'react';
import music from '../../images/music.jpeg'; // Placeholder image

const SongsList = ({ songs, onToggleView }) => {
  return (
    <React.Fragment>
      {songs.length > 0 && (
        <div className="songs">
          {songs.map((song, index) => (
            <div className="card" key={index}> {/* Use .card class here */}
              <a
                target="_blank"
                href={song.external_url}
                rel="noopener noreferrer"
                className="card-image-link" // Ensures padding and center alignment
              >
                {song.album_image_url ? (
                  <img src={song.album_image_url || music} alt={song.name} className="card-image" />
                ) : (
                  <img src={music} alt="Song" style={{width: '100%', height: 'auto'}} /> // Ensure placeholder image also fits
                )}
              </a>
              <div className="card-body"> {/* Use .card-body for textual content */}
                <h5 className="card-title">{song.name}</h5> {/* Song name as title */}
                <p className="card-text">{song.artist}</p> {/* Artist name as text */}
              </div>
            </div>
          ))}
        </div>
      )}
    </React.Fragment>
  );
};

export default SongsList;
