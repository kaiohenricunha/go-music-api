import React from 'react';
import _ from 'lodash';
import SongsList from './SongsList';

const SearchResult = (props) => {
  const { result, setCategory, selectedCategory } = props;
  console.log('result', result); // debug

  const songs = result && result.songs ? result.songs : []; // Default to an empty array if songs is not present in the result object

  return (
    <React.Fragment>
      <div className="search-buttons">
        {!_.isEmpty(songs.items) && (
          <button
            className={`${
              selectedCategory === 'songs' ? 'btn active' : 'btn'
            }`}
            onClick={() => setCategory('songs')}
          >
            Songs
          </button>
        )}
      </div>
      <div className={`${selectedCategory === 'songs' ? '' : 'hide'}`}>
        {songs && <SongsList songs={songs} />}
      </div>
    </React.Fragment>
  );
};

export default SearchResult;
