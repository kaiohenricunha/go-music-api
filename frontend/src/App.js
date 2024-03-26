import React from 'react';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import HomePage from './pages/HomePage';
import AddSongPage from './pages/AddSongPage';
import SongListPage from './pages/SongListPage';
import './App.css';

function App() {
  return (
    <Router>
      <Switch>
        <Route path="/" exact component={HomePage} />
        <Route path="/add-song" component={AddSongPage} />
        <Route path="/song-list" component={SongListPage} />
        {/* Add more routes as needed */}
      </Switch>
    </Router>
  );
}

export default App;
