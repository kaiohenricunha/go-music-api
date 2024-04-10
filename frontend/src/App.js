// App.js
import React from 'react';
import { Provider } from 'react-redux';
import store from './store/store';
import { AuthProvider } from './authContext';
import AppRouter from './router/AppRouter';

function App() {
  return (
    <Provider store={store}>
      <AuthProvider>  {/* Wrap with AuthProvider */}
        <AppRouter />
      </AuthProvider>
    </Provider>
  );
}

export default App;
