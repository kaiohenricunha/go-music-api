import { createStore, combineReducers, applyMiddleware, compose } from 'redux';
import thunk from 'redux-thunk';
import songsReducer from '../reducers/song';
import playlistReducer from '../reducers/playlist';

const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;

const store = createStore(
  combineReducers({
    songs: songsReducer,
    playlist: playlistReducer
  }),
  composeEnhancers(applyMiddleware(thunk))
);

export default store;
