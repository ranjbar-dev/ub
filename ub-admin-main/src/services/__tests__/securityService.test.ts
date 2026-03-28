import { loginAPI, refreshTokenAPI } from '../securityService';
import { apiService } from '../apiService';
import { RequestTypes, LocalStorageKeys, StandardResponse } from '../constants';

// jest.mock is hoisted — use jest.fn() inline; get the reference via the imported module below
jest.mock('services/apiService', () => ({
  apiService: {
    fetchData: jest.fn(),
  },
}));

const mockFetchData = apiService.fetchData as jest.Mock;

describe('securityService', () => {
  beforeEach(() => {
    mockFetchData.mockReset();
    localStorage.clear();
  });

  // ---------------------------------------------------------------------------
  // loginAPI
  // ---------------------------------------------------------------------------
  describe('loginAPI', () => {
    const credentials = { username: 'admin@example.com', password: 's3cret' };

    it('calls fetchData with the correct parameters', async () => {
      const mockResponse: StandardResponse<unknown> = {
        status: true,
        message: 'OK',
        data: { token: 'jwt-token', user: { id: 1, email: 'admin@example.com', name: 'Admin' } },
      };
      mockFetchData.mockResolvedValueOnce(mockResponse);

      await loginAPI(credentials);

      expect(mockFetchData).toHaveBeenCalledTimes(1);
      expect(mockFetchData).toHaveBeenCalledWith({
        data: credentials,
        url: 'auth/login',
        requestType: RequestTypes.POST,
        requestName: 'login',
      });
    });

    it('returns the StandardResponse on success', async () => {
      const mockResponse: StandardResponse<unknown> = {
        status: true,
        message: 'Logged in',
        data: { token: 'jwt-token' },
      };
      mockFetchData.mockResolvedValueOnce(mockResponse);

      const result = await loginAPI(credentials);

      expect(result).toEqual(mockResponse);
    });

    it('propagates errors thrown by fetchData', async () => {
      const error = new Error('Network failure');
      mockFetchData.mockRejectedValueOnce(error);

      await expect(loginAPI(credentials)).rejects.toThrow('Network failure');
    });
  });

  // ---------------------------------------------------------------------------
  // refreshTokenAPI
  // ---------------------------------------------------------------------------
  describe('refreshTokenAPI', () => {
    it('reads the refresh token from localStorage', async () => {
      const storedToken = 'stored-refresh-token';
      localStorage.setItem(LocalStorageKeys.REFRESH_TOKEN, storedToken);
      mockFetchData.mockResolvedValueOnce({ status: true, message: 'OK', data: {} });

      await refreshTokenAPI();

      expect(mockFetchData).toHaveBeenCalledWith(
        expect.objectContaining({
          data: { refresh_token: storedToken },
        }),
      );
    });

    it('calls fetchData with the correct URL, method, and requestName', async () => {
      localStorage.setItem(LocalStorageKeys.REFRESH_TOKEN, 'rt-abc');
      const mockResponse: StandardResponse<unknown> = { status: true, message: 'OK', data: {} };
      mockFetchData.mockResolvedValueOnce(mockResponse);

      await refreshTokenAPI();

      expect(mockFetchData).toHaveBeenCalledWith({
        data: { refresh_token: 'rt-abc' },
        url: 'auth/refresh',
        requestType: RequestTypes.POST,
        requestName: 'refresh-token',
      });
    });

    it('uses an empty string when no refresh token is stored', async () => {
      // localStorage is cleared in beforeEach
      mockFetchData.mockResolvedValueOnce({ status: true, message: 'OK', data: {} });

      await refreshTokenAPI();

      expect(mockFetchData).toHaveBeenCalledWith(
        expect.objectContaining({
          data: { refresh_token: '' },
        }),
      );
    });

    it('returns the response on success', async () => {
      localStorage.setItem(LocalStorageKeys.REFRESH_TOKEN, 'rt-xyz');
      const mockResponse: StandardResponse<unknown> = {
        status: true,
        message: 'Token refreshed',
        data: { token: 'new-jwt' },
      };
      mockFetchData.mockResolvedValueOnce(mockResponse);

      const result = await refreshTokenAPI();

      expect(result).toEqual(mockResponse);
    });

    it('uses requestName "refresh-token" (guards against infinite loop in apiService)', async () => {
      mockFetchData.mockResolvedValueOnce({ status: true, message: 'OK', data: {} });

      await refreshTokenAPI();

      const callArgs = mockFetchData.mock.calls[0][0] as { requestName: string };
      expect(callArgs.requestName).toBe('refresh-token');
    });
  });
});
