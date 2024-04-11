import { SET_SONGS, ADD_SONGS } from '../utils/constants';

const songsReducer = (state = { items: [] }, action) => {
    switch (action.type) {
      case SET_SONGS: 
        return { ...state, items: action.payload };
      case ADD_SONGS:
        return {
          ...state,
          items: [...state.items, ...action.payload]
        };
      default:
        return state;
    }
  };

export default songsReducer;