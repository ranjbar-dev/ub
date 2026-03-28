# Task: Add AbortController Timeout + Retry Logic to apiService

## Priority: 🔴 CRITICAL
## Files to Modify: `src/services/apiService.ts`, `src/services/constants.ts`

## Problem
1. `fetch()` has no timeout — hung requests never resolve
2. No retry logic — transient failures (502, 503, network blips) are permanent
3. No AbortController for request cancellation

## Current Code (`src/services/apiService.ts`)

The full file is ~120 lines. The key section is the `fetchData` method (lines 38-106):
```typescript
public async fetchData<T = unknown>(
    params: RequestParameters,
): Promise<StandardResponse<T>> {
    // ... token setup, URL building, logging ...
    
    const content: RequestInit = {
        method: params.requestType,
        headers: this.setHeaders(),
        body: params.requestType !== RequestTypes.GET
            ? JSON.stringify(params.data)
            : undefined,
    };

    let response: Response;
    try {
        response = await fetch(url, content);  // ⚠️ No timeout, no retry
    } catch (networkError) {
        throw new ApiError(
            `Network error calling ${params.requestType} ${params.url}`,
            0,
        );
    }
    // ... response handling ...
}
```

## Required Changes

### 1. Add timeout and retry constants to `src/services/constants.ts`

Add at the end of the file (after `webAppAddress`):
```typescript
/** Default request timeout in milliseconds */
export const API_TIMEOUT_MS = 30000;

/** Maximum retry attempts for idempotent requests */
export const API_MAX_RETRIES = 3;

/** Retryable HTTP status codes */
export const RETRYABLE_STATUS_CODES = [408, 429, 502, 503, 504];
```

### 2. Add `signal` support to `RequestParameters` in `src/services/constants.ts`

Add an optional `signal` field to `RequestParameters`:
```typescript
export interface RequestParameters<T = Record<string, unknown>> {
    requestType: RequestTypes;
    url: string;
    data: T;
    isRawUrl?: boolean;
    requestName?: string;
    signal?: AbortSignal;  // NEW: for request cancellation
}
```

### 3. Rewrite the fetch section in `src/services/apiService.ts`

Replace the fetch try-catch block (the `let response: Response;` through the `catch (networkError)` block) with:

```typescript
let response: Response;
const isIdempotent = params.requestType === RequestTypes.GET || params.requestType === RequestTypes.PUT;
const maxAttempts = isIdempotent ? API_MAX_RETRIES : 1;

for (let attempt = 1; attempt <= maxAttempts; attempt++) {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), API_TIMEOUT_MS);

    try {
        response = await fetch(url, {
            ...content,
            signal: params.signal || controller.signal,
        });
        clearTimeout(timeoutId);

        // Retry on retryable status codes (only for idempotent requests)
        if (RETRYABLE_STATUS_CODES.includes(response.status) && attempt < maxAttempts) {
            const backoffMs = Math.min(1000 * Math.pow(2, attempt - 1), 8000);
            await new Promise(resolve => setTimeout(resolve, backoffMs));
            continue;
        }

        break; // Success or non-retryable error
    } catch (fetchError) {
        clearTimeout(timeoutId);

        if (fetchError instanceof DOMException && fetchError.name === 'AbortError') {
            if (attempt < maxAttempts) {
                const backoffMs = Math.min(1000 * Math.pow(2, attempt - 1), 8000);
                await new Promise(resolve => setTimeout(resolve, backoffMs));
                continue;
            }
            throw new ApiError(
                `Request timeout after ${API_TIMEOUT_MS}ms: ${params.requestType} ${params.url}`,
                408,
            );
        }

        if (attempt < maxAttempts) {
            const backoffMs = Math.min(1000 * Math.pow(2, attempt - 1), 8000);
            await new Promise(resolve => setTimeout(resolve, backoffMs));
            continue;
        }

        throw new ApiError(
            `Network error calling ${params.requestType} ${params.url}`,
            0,
        );
    }
}
```

**IMPORTANT:** 
- Import `API_TIMEOUT_MS`, `API_MAX_RETRIES`, `RETRYABLE_STATUS_CODES` from `./constants`
- The `response!` variable needs the non-null assertion after the loop since TypeScript can't prove it's assigned
- Keep ALL existing code after this block (the `json`, status checks, error handling) unchanged
- After the for loop, use `response!` for the first access (e.g., `const json = await response!.json()...`)

## Validation
- Build must pass: `$env:NODE_OPTIONS='--openssl-legacy-provider'; npm run build`
- The retry logic should not affect normal 200/401/422 flows
- Only GET and PUT requests retry; POST and DELETE do not
