import axios from 'axios';

// Get API URL from runtime config or build-time env var or default
const getApiUrl = () => {
  // First try runtime config (injected by Docker)
  if (window._env_ && window._env_.REACT_APP_API_URL) {
    return window._env_.REACT_APP_API_URL;
  }
  
  // Then try build-time env var (Vite format)
  if (import.meta.env.VITE_API_URL) {
    return import.meta.env.VITE_API_URL;
  }
  
  // Default fallback
  return 'http://localhost:8080';
};

const API_BASE_URL = getApiUrl();

console.log('API URL configured as:', API_BASE_URL);

const client = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 10000, // 10 second timeout
});

// Add request interceptor for debugging
client.interceptors.request.use(
  (config) => {
    if (import.meta.env.DEV) {
      console.log('API Request:', config.method?.toUpperCase(), config.url);
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Add response interceptor for error handling
client.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response) {
      console.error('API Error:', error.response.status, error.response.data);
    } else if (error.request) {
      console.error('Network Error:', error.message);
    }
    return Promise.reject(error);
  }
);

const api = {
  // Session endpoints
  createSession: async () => {
    const response = await client.post('/sessions');
    return response.data;
  },
  
  getTodaySession: async () => {
    const response = await client.get('/sessions/today');
    return response.data;
  },
  
  getSession: async (id) => {
    const response = await client.get(`/sessions/${id}`);
    return response.data;
  },
  
  joinSession: async (id, username) => {
    const response = await client.post(`/sessions/${id}/join`, { username });
    return response.data;
  },
  
  orderMeal: async (id, username, meal) => {
    const response = await client.put(`/sessions/${id}/meal`, { username, meal });
    return response.data;
  },
  
  lockSession: async (id) => {
    const response = await client.post(`/sessions/${id}/lock`);
    return response.data;
  },
  
  // Restaurant endpoints
  getRestaurants: async (sessionId) => {
    const response = await client.get(`/sessions/${sessionId}/restaurants`);
    return response.data || [];
  },
  
  proposeRestaurant: async (sessionId, name, proposedBy) => {
    const response = await client.post(`/sessions/${sessionId}/restaurants`, {
      name,
      proposedBy,
    });
    return response.data;
  },
  
  voteRestaurant: async (sessionId, restaurantId, username) => {
    const response = await client.post(
      `/sessions/${sessionId}/restaurants/${restaurantId}/vote`,
      { username }
    );
    return response.data;
  },
  
  deleteRestaurant: async (sessionId, restaurantId) => {
    const response = await client.delete(
      `/sessions/${sessionId}/restaurants/${restaurantId}`
    );
    return response.data;
  },
  
  // Order placer endpoint
  selectOrderPlacer: async (sessionId, username) => {
    const response = await client.put(`/sessions/${sessionId}/order-placer`, {
      username,
    });
    return response.data;
  },
  
  // Heartbeat endpoint
  sendHeartbeat: async (sessionId, username) => {
    const response = await client.post(`/sessions/${sessionId}/heartbeat`, {
      username,
    });
    return response.data;
  },
  
  // Remove participant endpoint
  removeParticipant: async (sessionId, username) => {
    const response = await client.delete(`/sessions/${sessionId}/participants/${username}`);
    return response.data;
  },
};

export default api;