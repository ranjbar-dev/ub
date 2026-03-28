# Task: Implement Token Refresh Flow

## Priority: 🟠 HIGH
## Files to Modify: `src/services/apiService.ts`, `src/services/securityService.ts`, `src/app/containers/LoginPage/saga.ts`

## Problem
1. `REFRESH_TOKEN` key is defined in `LocalStorageKeys` but never used anywhere
2. On 401, the user is immediately kicked to login instead of attempting a token refresh
3. No refresh endpoint is called; sessions end abruptly

## Context
The backend likely has a refresh endpoint (since `REFRESH_TOKEN` was defined). The typical pattern:
- Login returns both `access_token` and `refresh_token`
- On 401, try refreshing with the refresh token
- If refresh succeeds, retry the original request
- If refresh fails, then redirect to login

**IMPORTANT:** Since we don't know the exact backend refresh endpoint, we need to:
1. Create the refresh infrastructure with a reasonable endpoint guess (`auth/refresh`)
2. Make it configurable so it can be adjusted when the backend is confirmed

## Required Changes

### 1. Add refresh API to `src/services/securityService.ts`

Add after the existing `loginAPI` function:
```typescript
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
```

Import `LocalStorageKeys` from `./constants` if not already imported.

### 2. Modify `src/services/apiService.ts` — Add token refresh on 401

The current 401 handler (around lines 84-88):
```typescript
if (response.status === 401) {
    MessageService.send({name: MessageNames.SETLOADING, payload: false});
    MessageService.send({name: MessageNames.AUTH_ERROR_EVENT});
    throw new ApiError('Authentication required', 401);
}
```

Replace with a token refresh attempt:
```typescript
if (response.status === 401) {
    // Don't try refresh for the refresh endpoint itself (avoid infinite loop)
    if (params.requestName === 'refresh-token' || params.url === 'auth/login') {
        MessageService.send({name: MessageNames.SETLOADING, payload: false});
        MessageService.send({name: MessageNames.AUTH_ERROR_EVENT});
        throw new ApiError('Authentication required', 401);
    }

    // Attempt token refresh
    try {
        const refreshToken = localStorage.getItem(LocalStorageKeys.REFRESH_TOKEN) || '';
        if (refreshToken) {
            const refreshResponse = await fetch(BaseUrl + 'admin/auth/refresh', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ refresh_token: refreshToken }),
            });

            if (refreshResponse.ok) {
                const refreshJson = await refreshResponse.json();
                if (refreshJson.token) {
                    // Store new token and retry original request
                    localStorage.setItem(LocalStorageKeys.ACCESS_TOKEN, refreshJson.token);
                    if (refreshJson.refresh_token) {
                        localStorage.setItem(LocalStorageKeys.REFRESH_TOKEN, refreshJson.refresh_token);
                    }
                    this.token = refreshJson.token;
                    // Retry original request with new token
                    return this.fetchData<T>(params);
                }
            }
        }
    } catch {
        // Refresh failed — fall through to auth error
    }

    MessageService.send({name: MessageNames.SETLOADING, payload: false});
    MessageService.send({name: MessageNames.AUTH_ERROR_EVENT});
    throw new ApiError('Authentication required', 401);
}
```

**IMPORTANT NOTES:**
- Import `BaseUrl` and `LocalStorageKeys` are already imported in apiService.ts
- The `params.requestName === 'refresh-token'` guard prevents infinite recursion
- The `params.url === 'auth/login'` guard prevents refresh attempts during login
- If no refresh token exists in localStorage, it skips straight to AUTH_ERROR_EVENT
- The recursive `this.fetchData<T>(params)` call retries with the new token

### 3. Store refresh token in `src/app/containers/LoginPage/saga.ts`

Find the login success block where `localStorage[LocalStorageKeys.ACCESS_TOKEN] = response.token` is set (around line 23). After it, add:
```typescript
if (response.data && (response.data as any).refresh_token) {
    localStorage[LocalStorageKeys.REFRESH_TOKEN] = (response.data as any).refresh_token;
}
```

This stores the refresh token if the login response includes one. The `as any` cast is needed since we don't know the exact response shape from the backend.

## Validation
- Build must pass: `$env:NODE_OPTIONS='--openssl-legacy-provider'; npm run build`
- Normal login flow should still work
- If backend doesn't have a refresh endpoint, the refresh attempt will fail silently and fall through to the existing 401 → AUTH_ERROR_EVENT → logout flow (backward compatible)
