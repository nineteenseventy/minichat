import axios from 'axios';
import auth0 from '@/auth0';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
});

api.interceptors.request.use(async (config) => {
  const token = await auth0.getAccessTokenSilently();
  console.log('token', token);
  config.headers.Authorization = `Bearer ${token}`;
  return config;
});

export const useApi = () => {
  return api;
};
