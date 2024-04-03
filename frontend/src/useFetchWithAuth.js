import { useAuth } from './authContext';

export function useFetchWithAuth() {
  const { token } = useAuth();

  const fetchWithAuth = async (url, options = {}) => {
    if (token) {
      options.headers = {
        ...options.headers,
        Authorization: `Bearer ${token}`,
      };
    }
    return fetch(url, options);
  };

  return fetchWithAuth;
}
