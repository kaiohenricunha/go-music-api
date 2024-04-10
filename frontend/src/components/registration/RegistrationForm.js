import React, { useState, useEffect } from 'react';

function RegistrationForm() {
    const [fullName, setFullName] = useState('');
    const [email, setEmail] = useState('');
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [role, setRole] = useState('');
    const [message, setMessage] = useState('');
    const [isError, setIsError] = useState(false);

    useEffect(() => {
        document.body.classList.add('bg-image-page1');
    
        return () => {
          document.body.classList.remove('bg-image-page1');
        };
      }, []);

    const handleInputChange = (e) => {
        const { name, value } = e.target;
        switch (name) {
            case 'name':
                setFullName(value);
                break;
            case 'email':
                setEmail(value);
                break;
            case 'username':
                setUsername(value);
                break;
            case 'password':
                setPassword(value);
                break;
            case 'role':
                setRole(value);
                break;
            default:
                break;
        }
    };

    const handleSubmit = async (event) => {
        event.preventDefault();
        setMessage(''); // Clear previous message
        setIsError(false); // Reset error state

        const apiEndpoint = `${process.env.REACT_APP_GO_BACKEND_BASE_URL}/register`;
        try {
            const response = await fetch(apiEndpoint, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ fullName, email, username, password, role }),
            });

            const data = await response.json();
            if (!response.ok) {
                throw new Error(data.message || 'Registration failed. Please try again.');
            }
            setMessage('Registration successful. Please log in.');
        } catch (error) {
            setMessage(error.toString());
            setIsError(true); // Set error state to true
        }
    };

    return (
        <div className="container">
            <form onSubmit={handleSubmit}>
                <div className="form-group">
                    <label htmlFor="name">Full Name</label>
                    <input
                        type="text"
                        name="name"
                        id="name"
                        className="form-control"
                        placeholder="Type your full name"
                        value={fullName}
                        onChange={handleInputChange} 
                        required
                    />
                </div>
                <div className="form-group">
                    <label htmlFor="email">Email</label>
                    <input
                        type="email"
                        name="email"
                        id="email"
                        className="form-control"
                        placeholder="Type your email"
                        value={email}
                        onChange={handleInputChange} 
                        required
                    />
                </div>
                <div className="form-group">
                    <label htmlFor="username">Username</label>
                    <input
                        type="text"
                        name="username"
                        id="username"
                        className="form-control"
                        placeholder="Type your username"
                        value={username}
                        onChange={handleInputChange} 
                        required
                    />
                </div>
                <div className="form-group">
                    <label htmlFor="password">Password</label>
                    <input
                        type="password"
                        name="password"
                        id="password"
                        className="form-control"
                        placeholder="Type your password"
                        value={password}
                        onChange={handleInputChange} 
                        required
                    />
                </div>
                <div className="form-group">
                    <label htmlFor="role">What option best defines you?</label>
                    <select
                        id="role"
                        name="role"
                        className="form-control"
                        value={role}
                        onChange={handleInputChange} 
                        required
                    >
                        <option disabled value="">Select your current status</option>
                        <option value="student">Music Student</option>
                        <option value="amateur">Amateur Musician</option>
                        <option value="professional">Professional Musician</option>
                        <option value="songwriter">Songwriter</option>
                        <option value="other">Other</option>
                    </select>
                </div>
                <div className="form-group">
                    <button type="submit" className="submit-button">Register</button>
                </div>
                {/* Display message */}
                {message && <div className={isError ? "error-message" : "success-message"}>{message}</div>}
            </form>
        </div>
    );
}

export default RegistrationForm;
