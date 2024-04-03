import React, { useState } from 'react';
import './RegistrationForm.css';

function RegistrationForm() {
    const [fullName, setFullName] = useState('');
    const [email, setEmail] = useState('');
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [role, setRole] = useState('');

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
        console.log("Submit button clicked");
        let obj = {
            full_name: fullName,
            email: email,
            username: username,
            password: password,
            role: role,
        };

        console.log("Submitting registration data:", obj);

        // use .env to store the API endpoint base URL
        const apiEndpoint = `${process.env.REACT_APP_API_URL}/api/v1/register`;
        try {
            const response = await fetch(apiEndpoint, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(obj),
            });

            if (response.ok) {
                const data = await response.json();
                console.log('Registration successful', data);
            } else {
                console.error('Registration failed', response.statusText);
            }
        } catch (error) {
            console.error('Error submitting registration', error);
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
            </form>
        </div>
    );
}

export default RegistrationForm;
