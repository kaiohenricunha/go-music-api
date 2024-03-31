import http from 'k6/http';
import { check, fail } from 'k6';
import { b64encode } from 'k6/encoding';
import { sleep } from 'k6';

export let options = {
    vus: 1, // Simulate 1 user
    iterations: 1, // Execute the sequence once
};

const BASE_URL = 'http://localhost:8081/api/v1';
const USERNAME = `testUser${Math.random().toString(36).substring(7)}`; // Ensure unique username
const PASSWORD = 'testPass';

function basicAuthHeader(user, pass) {
    const credentials = b64encode(`${user}:${pass}`);
    return `Basic ${credentials}`;
}

function createUser() {
    let res = http.post(`${BASE_URL}/register`, JSON.stringify({
        username: USERNAME,
        password: PASSWORD,
    }), {
        headers: {
            'Content-Type': 'application/json',
        },
    });
    if (res.status === 201) {
        console.log(`User ${USERNAME} created successfully for testing.`);
    } else {
        fail(`Failed to create test user ${USERNAME}: ${res.body}`);
    }
    return JSON.parse(res.body).ID;
}

function deleteUser(userId) {
    let res = http.del(`${BASE_URL}/users/${userId}`, null, {
        headers: { 'Authorization': basicAuthHeader(USERNAME, PASSWORD) },
    });
    if (res.status === 200) {
        console.log(`Test user ${USERNAME} deleted successfully.`);
    } else {
        console.log(`Failed to delete test user ${USERNAME}: ${res.body}`);
    }
}

export default function () {
    const userId = createUser();
    const authHeader = basicAuthHeader(USERNAME, PASSWORD);

    // Add a Song
    let res = http.post(`${BASE_URL}/music`, JSON.stringify({
        title: "New Song",
        artist: "New Artist",
    }), {
        headers: {
            'Content-Type': 'application/json',
            'Authorization': authHeader,
        },
    });
    check(res, { 'added song successfully': (r) => r.status === 200 }) || fail(`Failed to add song: ${res.body}`);
    const songId = JSON.parse(res.body).id;

    // Get the song
    res = http.get(`${BASE_URL}/music/${songId}`, {
        headers: { 'Authorization': authHeader },
    });
    check(res, { 'retrieved song successfully': (r) => r.status === 200 }) || fail(`Failed to retrieve song: ${res.body}`);

    // Update the song
    res = http.put(`${BASE_URL}/music/${songId}`, JSON.stringify({
        title: "Updated Song",
        artist: "Updated Artist",
    }), {
        headers: {
            'Content-Type': 'application/json',
            'Authorization': authHeader,
        },
    });
    check(res, { 'updated song successfully': (r) => r.status === 200 }) || fail(`Failed to update song: ${res.body}`);

    // Delete the song
    res = http.del(`${BASE_URL}/music/${songId}`, null, {
        headers: { 'Authorization': authHeader },
    });
    check(res, { 'deleted song successfully': (r) => r.status === 200 }) || fail(`Failed to delete song: ${res.body}`);

    // Clean up: Delete the test user
    deleteUser(userId);
}
