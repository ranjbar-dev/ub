# Task: Fix apiService PUT/DELETE Body + Add Auth to File Downloads

## Priority: 🟠 HIGH (SEC-H3 + SEC-M7)
## Files to Modify: `src/services/apiService.ts`, `src/utils/fileDownload.ts`

## Problem
1. `apiService.ts` only sends JSON body for POST requests. PUT and DELETE requests have `body: undefined`, silently dropping data.
2. `fileDownload.ts` fetches files without Bearer token — secured endpoints reject with 401.

## Current Code

### `src/services/apiService.ts` lines 61-67:
```typescript
const content: RequestInit = {
    method: params.requestType,
    headers: this.setHeaders(),
    body: params.requestType === RequestTypes.POST
        ? JSON.stringify(params.data)
        : undefined,                               // ⚠️ PUT/DELETE have NO BODY
};
```

### `src/services/apiService.ts` lines 39-43 — setHeaders:
```typescript
setHeaders(): Record<string, string> {
    this.token = localStorage[LocalStorageKeys.ACCESS_TOKEN]
        ? localStorage[LocalStorageKeys.ACCESS_TOKEN]
        : '';
    return {
        Authorization: `Bearer ${this.token}`,
        'Content-Type': 'application/json',
    };
}
```

### `src/utils/fileDownload.ts` lines 35-43:
```typescript
const downloadFile = async (params: { url: string; filename: string }): Promise<void> => {
  try {
    const response = await fetch(params.url);    // ⚠️ NO AUTH HEADER
    const blob = await response.blob();
    fileDownloader(blob, params.filename);
  } catch (error) {
    console.error('File download failed:', error);
  }
};
```

### `RequestTypes` enum (from constants.ts):
```typescript
export enum RequestTypes {
    GET = 'GET',
    POST = 'POST',
    PUT = 'PUT',
    DELETE = 'DELETE',
}
```

## Required Changes

### 1. Fix body handling in `src/services/apiService.ts` lines 64-66:

Replace:
```typescript
body: params.requestType === RequestTypes.POST
    ? JSON.stringify(params.data)
    : undefined,
```

With:
```typescript
body: params.requestType !== RequestTypes.GET
    ? JSON.stringify(params.data)
    : undefined,
```

This sends the JSON body for POST, PUT, and DELETE — only GET omits the body.

### 2. Add Bearer token to `src/utils/fileDownload.ts`:

Replace the `downloadFile` function:
```typescript
const downloadFile = async (params: { url: string; filename: string }): Promise<void> => {
  try {
    const token = localStorage.getItem('access_token') || '';
    const response = await fetch(params.url, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    if (!response.ok) {
      throw new Error(`Download failed with status ${response.status}`);
    }
    const blob = await response.blob();
    fileDownloader(blob, params.filename);
  } catch (error) {
    console.error('File download failed:', error);
  }
};
```

Changes:
- Adds `Authorization: Bearer <token>` header from localStorage
- Checks `response.ok` before attempting blob conversion
- Uses `localStorage.getItem('access_token')` (the string value, matching `LocalStorageKeys.ACCESS_TOKEN = 'access_token'`)

**NOTE:** We use the string literal `'access_token'` rather than importing `LocalStorageKeys` to avoid circular dependency with the services layer. The value matches `LocalStorageKeys.ACCESS_TOKEN`.

## Validation
- Build must pass: `$env:NODE_OPTIONS='--openssl-legacy-provider'; npm run build`
- Verify: PUT requests now send JSON body to the server
- Verify: File downloads include Authorization header
