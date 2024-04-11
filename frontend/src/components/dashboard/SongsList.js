import React from 'react';
import _ from 'lodash';
import music from '../../images/music.jpeg'; // Placeholder image

const SongsList = ({ songs }) => {
  return (
    <React.Fragment>
      {Object.keys(songs).length > 0 && (
        <div className="songs">
          {songs.map((song, index) => {
            return (
              <React.Fragment key={index}>
                <div className="song-item"> {/* Adjust class name as needed */}
                  <a
                    target="_blank"
                    href={song.external_url}
                    rel="noopener noreferrer"
                    className="song-image-link"
                  >
                    {!_.isEmpty(song.album_image_url) ? (
                      <img src={song.album_image_url} alt={song.name} />
                    ) : (
                      <img src={music} alt="Song" />
                    )}
                  </a>
                  <div className="song-info">
                    <p>{song.name}</p>
                    <small>{song.artist}</small>
                  </div>
                </div>
              </React.Fragment>
            );
          })}
        </div>
      )}
    </React.Fragment>
  );
};

export default SongsList;
