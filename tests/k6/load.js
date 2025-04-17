import http from 'k6/http';
import { check } from 'k6';

export let options = {
    vus: 20,
    duration: '30s',
};

const BASE_URL = 'http://localhost:5000';
const JWT_TOKEN = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ3ODM1MzksImlhdCI6MTc0NDc3OTkzOSwic3ViIjoiMTAzIiwidG9rZW5fdHlwZSI6ImFjY2VzcyJ9.UWjxxg0r5L_sx3JuXR-qNKMUVg2hz5JLJoPcqctbZjMs';

export default function () {
    let headers = {
      'Authorization': `Bearer ${JWT_TOKEN}`,
      'Content-Type': 'application/json',
    };
  
    let res = http.get(`${BASE_URL}/api/users/23`, { headers: headers });
    // let res = http.get(`${BASE_URL}/api/users`, { headers: headers });
    // let res = http.get(`${BASE_URL}/api/chat/messages`, { headers: headers });

    check(res, {
      'is status 200 / 429': (r) => r.status === 200 || r.status === 429,
    });
  }