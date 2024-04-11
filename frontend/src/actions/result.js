import {
    SET_SONGS,
    ADD_SONG,
    SET_PLAYLIST,
    ADD_PLAYLIST,
  } from '../utils/constants';
  import { get } from '../utils/api';
  
  // Update API URL using the environment variable
  const API_URL = `${process.env.REACT_APP_GO_BACKEND_BASE_URL}/songs/search`;
  
  export const setSongs = (songs) => ({
    type: SET_SONGS,
    songs,
  });
  
  export const addSong = (song) => ({
    type: ADD_SONG,
    song,
  });
  
  export const setPlaylist = (playlists) => ({
    type: SET_PLAYLIST,
    playlists,
  });
  
  export const addPlaylist = (playlists) => ({
    type: ADD_PLAYLIST,
    playlists,
  });
  
  export const initiateGetSongs = (searchTerm) => {
    return async (dispatch) => {
      try {
        const url = `${API_URL}?trackName=${encodeURIComponent(searchTerm)}`;
        const result = await get(url);
        console.log(result);
        const { songs } = result;
        dispatch(setSongs(songs));
      } catch (error) {
        console.log('error', error);
      }
    };
  };
  