import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom'; // This will be used for redirection
import { useAuth } from '../../authContext';

function LoginForm() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const { login } = useAuth();
  const navigate = useNavigate(); // Initialized for redirection
  const apiEndpoint = `${process.env.REACT_APP_API_URL}/api/v1/login`;

  useEffect(() => {
    document.body.classList.add('bg-image-page2');

    return () => {
      document.body.classList.remove('bg-image-page2');
    };
  }, []);

  const handleLogin = async (event) => {
    event.preventDefault();
    try {
      const response = await fetch(apiEndpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Basic ' + btoa(username + ':' + password),
        },
        body: JSON.stringify({ username, password }),
      });
      
      if (response.ok) {
        const data = await response.json();
        login(data.token); // Use the login function from context
        navigate('/dashboard'); // Redirect to dashboard after login
      } else {
        alert('Login failed!');
      }
    } catch (error) {
      console.error('Login error:', error);
    }
  };

  return (
    <div className="container">
      <form onSubmit={handleLogin}>
        <div className="form-group">
          <label id="username-label" htmlFor="username">Username</label>
          <input
            type="text"
            name="username"
            id="username"
            className="form-control"
            placeholder="Type your username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
          />
        </div>
        <div className="form-group">
          <label id="password-label" htmlFor="password">Password</label>
          <input
            type="password"
            name="password"
            id="password"
            className="form-control"
            placeholder="Type your password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        <div className="form-group">
          <input id="submit" className="submit-button" type="submit" value="Login" />
        </div>
        <div className="text-center">
          <a href="/register">Register</a> {/* Use Link component for SPA */}
        </div>
      </form>
    </div>
  );
}

export default LoginForm;
