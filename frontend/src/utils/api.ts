import axios, { AxiosError, AxiosInstance } from 'axios';

// Create axios instance with default config
const api: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_URL ?? 'https://api.easyfund.aldar.space/api/v1',
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
  refresh: () => api.post('/auth/refresh'),
};

export const usersAPI = {
  getAll: () => api.get('/users'),
  getById: (id: string) => api.get(`/users/${id}`),
  update: (id: string, data: any) => api.put(`/users/${id}`, data),
  delete: (id: string) => api.delete(`/users/${id}`),
  getRandom: () => api.get('/users/random'),
};

export const accountsAPI = {
  getUserAccounts: (userId: string) => api.get(`/users/${userId}/accounts`),
  getUserBalance: (userId: string) => api.get(`/users/${userId}/balance`),
};

export const transactionsAPI = {
  getUserTransactions: (userId: string, params?: any) => 
    api.get(`/users/${userId}/transactions`, { params }),
  getUserBankTransactions: (userId: string, bankId: number) => 
    api.get(`/users/${userId}/banks/${bankId}/transactions`),
};

export const loansAPI = {
  getUserLoans: (userId: string) => api.get(`/users/${userId}/loans`),
  createLoan: (data: any) => api.post('/loans', data),
  getLoan: (loanId: string) => api.get(`/loans/${loanId}`),
  makePayment: (loanId: string, data: any) => api.post(`/loans/${loanId}/payment`, data),
  getUserDebt: (userId: string) => api.get(`/users/${userId}/debt`),
};

// Добавляем loanApplicationsAPI для обратной совместимости
export const loanApplicationsAPI = {
  getAll: () => api.get('/applications'),
  getById: (id: string) => api.get(`/applications/${id}`),
  create: (data: any) => api.post('/applications', data),
  update: (id: string, data: any) => api.put(`/applications/${id}`, data),
  delete: (id: string) => api.delete(`/applications/${id}`),
};

export const applicationsAPI = {
  getAll: () => api.get('/applications'),
  getById: (id: string) => api.get(`/applications/${id}`),
  create: (data: any) => api.post('/applications', data),
  update: (id: string, data: any) => api.put(`/applications/${id}`, data),
  delete: (id: string) => api.delete(`/applications/${id}`),
  getUserApplications: (userId: string) => api.get(`/users/${userId}/applications`),
  approve: (applicationId: string) => api.post(`/applications/${applicationId}/approve`),
  reject: (applicationId: string) => api.post(`/applications/${applicationId}/reject`),
};

export const banksAPI = {
  getAll: () => api.get('/banks'),
  getById: (id: string) => api.get(`/banks/${id}`),
};

export const healthAPI = {
  check: () => api.get('/health'),
};

// Dashboard API - собираем данные из разных эндпоинтов
export const dashboardAPI = {
  // Старый метод с ограниченными транзакциями (для обратной совместимости)
  getDashboardData: (userId: string) => 
    Promise.all([
      accountsAPI.getUserBalance(userId),
      loansAPI.getUserDebt(userId),
      loansAPI.getUserLoans(userId),
      transactionsAPI.getUserTransactions(userId, { limit: 5 }),
      applicationsAPI.getUserApplications(userId),
    ]),
  
  // ✅ НОВЫЕ МЕТОДЫ для получения всех данных без ограничений
  getBalanceSummary: (userId: string) => accountsAPI.getUserBalance(userId),
  getUserDebt: (userId: string) => loansAPI.getUserDebt(userId),
  getUserLoans: (userId: string) => loansAPI.getUserLoans(userId),
  getUserApplications: (userId: string) => applicationsAPI.getUserApplications(userId),
  
  // ✅ ВАЖНО: Метод для получения ВСЕХ транзакций без ограничений
  getUserTransactions: (userId: string, params?: any) => 
    transactionsAPI.getUserTransactions(userId, params),
    
  // Метод для получения всех данных с большим лимитом транзакций
  getFullDashboardData: (userId: string) => 
    Promise.all([
      accountsAPI.getUserBalance(userId),
      loansAPI.getUserDebt(userId),
      loansAPI.getUserLoans(userId),
      transactionsAPI.getUserTransactions(userId, { limit: 200 }), // Увеличили лимит
      applicationsAPI.getUserApplications(userId),
    ]),
};