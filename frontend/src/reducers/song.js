import { SET_SONGS, ADD_SONGS } from '../utils/constants';

const songsReducer = (state = {}, action) => {
  const { songs } = action;
  switch (action.type) {
    case SET_SONGS:
      return songs;
    case ADD_SONGS:
      return {
        ...state,
        next: songs.next,
        items: [...state.items, ...songs.items]
      };
    default:
      return state;
  }
};

export default songsReducer;
