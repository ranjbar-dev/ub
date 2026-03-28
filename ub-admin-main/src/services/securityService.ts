import { apiService } from './apiService';
import { RequestTypes, LocalStorageKeys } from './constants';

export interface LoginParams {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: {
    id: number;
    email: string;
    name: string;
  };
}

/**
 * Authenticates an admin user and returns a JWT token.
 *
 * @param parameters - Login credentials (email, password)
 * @returns Promise with JWT token and user info
 * @endpoint POST auth/login
 */
export const loginAPI = (parameters: LoginParams) => {
  return apiService.fetchData({
    data: parameters as unknown as Record<string, unknown>,
    url: 'auth/login',
    requestType: RequestTypes.POST,
    requestName: 'login',
  });
};

/**
 * Attempts to refresh the access token using the stored refresh token.
 * @endpoint POST auth/refresh
 */
export const refreshTokenAPI = () => {
  const refreshToken = localStorage.getItem(LocalStorageKeys.REFRESH_TOKEN) || '';
  return apiService.fetchData({
    data: { refresh_token: refreshToken } as unknown as Record<string, unknown>,
    url: 'auth/refresh',
    requestType: RequestTypes.POST,
    requestName: 'refresh-token',
  });
};
