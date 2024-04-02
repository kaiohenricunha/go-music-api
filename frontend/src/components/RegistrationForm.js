import React, {useState,setState} from 'react';
import './style.css'
function RegistrationForm() {
    
    const [fullName, setFullName] = useState(null);
    const [email, setEmail] = useState(null);
    const [username, setUsername] = useState(null);
    const [password,setPassword] = useState(null);
    const [confirmPassword,setConfirmPassword] = useState(null);

    const handleInputChange = (e) => {
        const {id , value} = e.target;
        if(id === "fullName"){
            setFullName(value);
        }
        if(id === "email"){
            setEmail(value);
        }
        if(id === "username"){
            setUsername(value);
        }
        if(id === "password"){
            setPassword(value);
        }
        if(id === "confirmPassword"){
            setConfirmPassword(value);
        }

    }

    const handleSubmit = async (event) => {
        event.preventDefault();
        console.log("Submit button clicked"); // Confirm function is triggered
        let obj = {
            fullName: fullName,
            email: email,
            username: username,
            password: password,
            confirmPassword: confirmPassword,
        };

        console.log("Submitting registration data:", obj); // Log data being sent
    
        // Define the API endpoint    
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
    }

    return(
      <div className="form">
          <div className="form-body">
              <div className="name">
                  <label className="form__label" htmlFor="fullName">Full Name </label>
                  <input className="form__input" type="text" value={fullName} onChange={(e) => handleInputChange(e)} id="fullName" placeholder="Full Name"/>
              </div>
              <div className="email">
                  <label className="form__label" htmlFor="email">Email </label>
                    <input className="form__input" type="email" id="email" value={email} onChange = {(e) => handleInputChange(e)} placeholder="Email"/>
              </div>
                <div className="usernam">
                    <label className="form__label" htmlFor="username">Username </label>
                    <input className="form__input" type="text" id="username" value={username} onChange = {(e) => handleInputChange(e)} placeholder="Username"/>
                </div>
              <div className="password">
                  <label className="form__label" htmlFor="password">Password </label>
                    <input className="form__input" type="password" id="password" value={password} onChange = {(e) => handleInputChange(e)} placeholder="Password"/>
              </div>
              <div className="confirm-password">
                  <label className="form__label" htmlFor="confirmPassword">Confirm Password </label>
                  <input className="form__input" type="password" id="confirmPassword" value={confirmPassword} onChange = {(e) => handleInputChange(e)} placeholder="Confirm Password"/>
              </div>
          </div>
          <div className="footer">
                <button onClick={(e) => handleSubmit(e)} className="btn">Register</button>
                <button onClick={() => console.log('Test button clicked')}>Test Log</button>
          </div>
        </div>
       
    )       
}

export default RegistrationForm
