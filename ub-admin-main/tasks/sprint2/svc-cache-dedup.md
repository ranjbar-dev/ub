# Task: Cache Static Data + Add Request Deduplication

## Priority: 🟡 MEDIUM
## Files to Modify: `src/services/globalDataService.ts`, `src/services/apiService.ts`

## Problem
1. Countries, currencies, managers, and pairs are fetched on EVERY login but never cached with TTL
2. No request deduplication — rapid clicks or concurrent renders can trigger identical parallel API calls

## Current Code

### `src/services/globalDataService.ts` (read the full file):
Contains `GetCountriesAPI`, `GetCurrenciesAPI`, `GetManagersAPI` — all simple wrappers around `apiService.fetchData`.

### How they're called (`LoginPage/saga.ts` lines 24-51):
```typescript
const countriesResponse = yield call(GetCountriesAPI);
localStorage[LocalStorageKeys.COUNTRIES] = JSON.stringify(countriesResponse.data);
const currenciesResponse = yield call(GetCurrenciesAPI);
localStorage[LocalStorageKeys.CURRENCIES] = JSON.stringify(currenciesResponse.data);
const managersResponse = yield call(GetManagersAPI);
localStorage[LocalStorageKeys.Managers] = JSON.stringify(managersResponse.data);
```

## Required Changes

### 1. Add a simple in-memory cache to `src/services/apiService.ts`

Add a cache map and helper BEFORE the `ApiService` class:
```typescript
/** Simple in-memory cache for API responses */
interface CacheEntry<T = unknown> {
    data: StandardResponse<T>;
    expiresAt: number;
}

const apiCache = new Map<string, CacheEntry>();

/** In-flight request deduplication map */
const pendingRequests = new Map<string, Promise<StandardResponse<any>>>();

/** Default cache TTL: 1 hour */
const DEFAULT_CACHE_TTL_MS = 60 * 60 * 1000;
```

### 2. Add `cache` and `cacheKey` options to `RequestParameters` in `src/services/constants.ts`

```typescript
export interface RequestParameters<T = Record<string, unknown>> {
    requestType: RequestTypes;
    url: string;
    data: T;
    isRawUrl?: boolean;
    requestName?: string;
    signal?: AbortSignal;
    cacheTtlMs?: number;   // NEW: cache TTL in ms (0 = no cache)
}
```

### 3. Add cache check at the START of `fetchData` in `apiService.ts`

At the beginning of the `fetchData` method, after the URL construction (after line ~48), add:
```typescript
// Cache check (GET requests only)
const cacheKey = `${params.requestType}:${url}`;
if (params.requestType === RequestTypes.GET && params.cacheTtlMs) {
    const cached = apiCache.get(cacheKey);
    if (cached && cached.expiresAt > Date.now()) {
        return cached.data as StandardResponse<T>;
    }
}

// Request deduplication (GET requests only)
if (params.requestType === RequestTypes.GET) {
    const pending = pendingRequests.get(cacheKey);
    if (pending) {
        return pending as Promise<StandardResponse<T>>;
    }
}
```

### 4. Add cache storage AFTER successful response

After `return json as StandardResponse<T>;` (around line 82), wrap it to cache the result:
```typescript
if (response.status === 200) {
    const result = json as StandardResponse<T>;
    // Cache successful GET responses if cacheTtlMs specified
    if (params.requestType === RequestTypes.GET && params.cacheTtlMs) {
        apiCache.set(cacheKey, {
            data: result,
            expiresAt: Date.now() + params.cacheTtlMs,
        });
    }
    pendingRequests.delete(cacheKey);
    return result;
}
```

### 5. Wrap the entire fetch+response section in dedup tracking

Before the fetch attempt, register the promise:
```typescript
if (params.requestType === RequestTypes.GET) {
    const fetchPromise = this._doFetch<T>(url, content, params, cacheKey);
    pendingRequests.set(cacheKey, fetchPromise);
    return fetchPromise;
}
```

Actually — this is getting complex. Let's use a SIMPLER approach. Instead of modifying the core `fetchData`, add caching at the service function level:

### SIMPLIFIED APPROACH — Cache in `globalDataService.ts`

Replace the contents of `src/services/globalDataService.ts` with:
```typescript
import { apiService } from './apiService';
import { RequestTypes, StandardResponse } from './constants';

/** Simple TTL cache for static data */
const cache = new Map<string, { data: StandardResponse; expiresAt: number }>();
const CACHE_TTL = 60 * 60 * 1000; // 1 hour

function getCached(key: string): StandardResponse | null {
    const entry = cache.get(key);
    if (entry && entry.expiresAt > Date.now()) {
        return entry.data;
    }
    cache.delete(key);
    return null;
}

function setCache(key: string, data: StandardResponse): void {
    cache.set(key, { data, expiresAt: Date.now() + CACHE_TTL });
}

/** In-flight request dedup */
const pending = new Map<string, Promise<StandardResponse>>();

async function fetchWithCache(url: string, cacheKey: string): Promise<StandardResponse> {
    // Check cache
    const cached = getCached(cacheKey);
    if (cached) return cached;

    // Dedup concurrent requests
    const inflight = pending.get(cacheKey);
    if (inflight) return inflight;

    const promise = apiService.fetchData({
        data: {} as Record<string, unknown>,
        url,
        requestType: RequestTypes.GET,
    }).then(response => {
        setCache(cacheKey, response);
        pending.delete(cacheKey);
        return response;
    }).catch(error => {
        pending.delete(cacheKey);
        throw error;
    });

    pending.set(cacheKey, promise);
    return promise;
}

/** Fetches the list of countries. Cached for 1 hour. */
export const GetCountriesAPI = () => fetchWithCache('countries', 'countries');

/** Fetches the list of currencies. Cached for 1 hour. */
export const GetCurrenciesAPI = () => fetchWithCache('currencies', 'currencies');

/** Fetches the list of admin/manager users. Cached for 1 hour. */
export const GetManagersAPI = () => fetchWithCache('managers', 'managers');

/** Invalidates all cached static data. Call on relevant mutations. */
export function invalidateGlobalDataCache(): void {
    cache.clear();
}
```

This keeps the cache simple and isolated to global data services. No changes needed to `apiService.ts` for this task.

**DO NOT modify `RequestParameters` or `apiService.ts` for caching — keep it in the service layer.**

## Validation
- Build must pass: `$env:NODE_OPTIONS='--openssl-legacy-provider'; npm run build`
- Login flow should still work (cached data returned on subsequent calls)
- Concurrent calls to the same endpoint should return the same promise
