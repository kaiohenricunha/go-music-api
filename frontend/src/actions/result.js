import {
    SET_SONGS,
    ADD_SONG,
    ADD_SONGS,
    SET_PLAYLIST,
    ADD_PLAYLIST,
  } from '../utils/constants';
    import { get } from '../utils/api';

    const API_URL = `${process.env.REACT_APP_GO_BACKEND_BASE_URL}/songs/search`;

    export const setSongs = (songs) => ({
    type: SET_SONGS,
    songs,
    });

    export const addSong = (song) => ({
    type: ADD_SONG,
    song,
    });

    export const addSongs = (songs) => ({
    type: ADD_SONGS,
    songs,
    });

    export const setPlaylist = (playlists) => ({
    type: SET_PLAYLIST,
    playlists,
    });

    export const addPlaylist = (playlists) => ({
    type: ADD_PLAYLIST,
    playlists,
    });

    export const initiateGetSongs = (searchTerm, token) => async (dispatch) => {
        try {
          const [songName, artistName] = searchTerm.split(' by ').map(s => s.trim());
          const songs = await get(`${process.env.REACT_APP_GO_BACKEND_BASE_URL}/songs/search`, { songName, artistName }, token);
          dispatch({ type: 'SET_SONGS', payload: songs });
        } catch (error) {
          console.error('Search error:', error);
        }
      };
