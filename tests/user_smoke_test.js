import http from 'k6/http';
import { check, fail } from 'k6';
import { b64encode } from 'k6/encoding';

export let options = {
    vus: 1, // 1 virtual user
    iterations: 1, // run the default function exactly once
};

const BASE_URL = 'http://localhost:8081/api/v1';
const USERNAME = `testUser${Math.random().toString(36).substring(7)}`; // Generate a random username
const PASSWORD = 'testPass';
const NEW_PASSWORD = 'newPassword'; // Correctly define the new password for update operation

function basicAuthHeader(user, pass) {
    const credentials = `${user}:${pass}`;
    const encodedCredentials = b64encode(credentials);
    return `Basic ${encodedCredentials}`;
}

export default function () {
    // Register a new user
    let res = http.post(`${BASE_URL}/register`, JSON.stringify({
        username: USERNAME,
        password: PASSWORD,
    }), {
        headers: { 'Content-Type': 'application/json' },
    });
    if (!check(res, { 'registered successfully': (r) => r.status === 201 })) {
        fail(`Failed to register user: ${res.body}`);
    }

    let userResponse = JSON.parse(res.body);
    let userIdToUpdate = userResponse.ID; // Use the correct field name based on your API's response
    if (userIdToUpdate === undefined) {
        fail('User ID is undefined. Check the API response structure.');
    }
    console.log(`Registered user ID: ${userIdToUpdate}`);

    // Create a new user with the same username
    res = http.post(`${BASE_URL}/register`, JSON.stringify({
        username: USERNAME,
        password: PASSWORD,
    }), {
        headers: { 'Content-Type': 'application/json' },
    });
    check(res, { 'failed to register the same user again': (r) => r.status === 409 });

    // List users (assumes the endpoint requires authentication)
    res = http.get(`${BASE_URL}/users`, {
        headers: { 'Authorization': basicAuthHeader(USERNAME, PASSWORD) },
    });
    check(res, { 'listed users successfully': (r) => r.status === 200 });

    // Try to list users with wrong credentials
    // Should fail with a 401 Unauthorized status
    res = http.get(`${BASE_URL}/users`, {
        headers: { 'Authorization': basicAuthHeader(USERNAME, 'wrongPassword') },
    });
    check(res, { 'failed to list users': (r) => r.status === 401 });

    // Find the newly registered user by username
    res = http.get(`${BASE_URL}/users/${USERNAME}`, {
        headers: { 'Authorization': basicAuthHeader(USERNAME, PASSWORD) },
    });
    check(res, { 'found user by username successfully': (r) => r.status === 200 });

    // Try to find the user with wrong credentials
    // Should fail with a 401 Unauthorized status
    res = http.get(`${BASE_URL}/users/${USERNAME}`, {
        headers: { 'Authorization': basicAuthHeader(USERNAME, 'wrongPassword') },
    });
    check(res, { 'failed to find user by username': (r) => r.status === 401 });

    // Update the user's information using the dynamically obtained userIdToUpdate
    res = http.put(`${BASE_URL}/users/${userIdToUpdate}`, JSON.stringify({
        password: NEW_PASSWORD,
    }), {
        headers: {
            'Authorization': basicAuthHeader(USERNAME, PASSWORD),
            'Content-Type': 'application/json',
        },
    });
    check(res, { 'updated user successfully': (r) => r.status === 200 });

    // Try to update a different user than the one authenticated
    // Should fail with a 403 Forbidden status
    res = http.put(`${BASE_URL}/users/1`, JSON.stringify({
        password: NEW_PASSWORD,
    }), {
        headers: {
            'Authorization': basicAuthHeader(USERNAME, PASSWORD),
            'Content-Type': 'application/json',
        },
    });
    check(res, { 'unauthorized to update this user': (r) => r.status === 401 });

    // Delete the user using the dynamically obtained userIdToUpdate and the new password
    res = http.del(`${BASE_URL}/users/${userIdToUpdate}`, null, {
        headers: { 'Authorization': basicAuthHeader(USERNAME, NEW_PASSWORD) }, // Using the new password here
    });
    check(res, { 'deleted user successfully': (r) => r.status === 200 });
}
