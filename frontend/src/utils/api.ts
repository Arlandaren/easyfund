import axios, { AxiosError, AxiosInstance } from 'axios';

// Create axios instance with default config
const api: AxiosInstance = axios.create({
  baseURL: '/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor for adding auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor for handling errors
api.interceptors.response.use(
  (response) => response,
  (error: AxiosError) => {
    if (error.response?.status === 401) {
      // Clear token and redirect to login
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default api;

// API endpoints
export const authAPI = {
  login: (email: string, password: string) =>
    api.post('/auth/login', { email, password }),
  register: (data: any) => api.post('/auth/register', data),
  logout: () => api.post('/auth/logout'),
  me: () => api.get('/auth/me'),
};

export const usersAPI = {
  getAll: () => api.get('/users'),
  getById: (id: string) => api.get(`/users/${id}`),
  update: (id: string, data: any) => api.put(`/users/${id}`, data),
  delete: (id: string) => api.delete(`/users/${id}`),
};

export const loanApplicationsAPI = {
  getAll: () => api.get('/loan-applications'),
  getById: (id: string) => api.get(`/loan-applications/${id}`),
  create: (data: any) => api.post('/loan-applications', data),
  update: (id: string, data: any) => api.put(`/loan-applications/${id}`, data),
  delete: (id: string) => api.delete(`/loan-applications/${id}`),
};

export const banksAPI = {
  getAll: () => api.get('/banks'),
  getById: (id: string) => api.get(`/banks/${id}`),
};

