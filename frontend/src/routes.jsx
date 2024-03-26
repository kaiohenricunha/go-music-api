// src/routes.jsx
import React from 'react';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import HomePage from './pages/HomePage';
import AddSongPage from './pages/AddSongPage';
import SongListPage from './pages/SongListPage';

const Routes = () => (
  <Router>
    <Switch>
      <Route path="/" exact component={HomePage} />
      <Route path="/add-song" component={AddSongPage} />
      <Route path="/songs" component={SongListPage} />
    </Switch>
  </Router>
);

export default Routes;
