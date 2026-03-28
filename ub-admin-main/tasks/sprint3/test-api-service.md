# Task: Test ApiService

## Goal
Create `src/services/__tests__/apiService.test.ts` with tests for the core HTTP layer.

## Context
- File to test: `src/services/apiService.ts`
- Singleton pattern with `getInstance()`
- Core method: `fetchData<T>(params: RequestParameters): Promise<StandardResponse<T>>`
- Sprint 2 additions: AbortController timeout (30s), exponential backoff retry (3 attempts for GET/PUT), token refresh on 401
- Throws `ApiError` on non-200 responses
- Emits `AUTH_ERROR_EVENT` via MessageService when auth fails completely

## File to Create: `src/services/__tests__/apiService.test.ts`

### Test Cases Required

#### Singleton
- `getInstance()` returns same instance
- Instance has token from localStorage

#### Successful requests
- GET request → correct URL, headers, returns parsed JSON
- POST request → includes body, correct Content-Type
- PUT request → includes body
- DELETE request → correct method

#### Error handling
- 400 response → throws ApiError with status 400
- 422 response → throws ApiError with validation errors
- 500 response → throws ApiError with status 500

#### Retry logic (Sprint 2)
- GET with 503 → retries up to 3 times with backoff
- POST with 503 → does NOT retry (non-idempotent)
- GET with 408 → retries
- Network error → retries for GET
- All retries exhausted → throws last error

#### 401 Token refresh (Sprint 2)
- 401 response + valid refresh token → attempts refresh → retries original request
- 401 response + refresh fails → emits AUTH_ERROR_EVENT
- 401 response + no refresh token → emits AUTH_ERROR_EVENT
- Refresh request itself (requestName='refresh-token') gets 401 → does NOT try to refresh again (infinite loop guard)

#### Timeout (Sprint 2)
- Request exceeding 30s → aborted via AbortController

## IMPORTANT
- Read `src/services/apiService.ts` FIRST — the retry + refresh logic is complex
- Mock `window.fetch` with `jest.fn()` (existing pattern from request.test.ts)
- Mock `localStorage` for token access
- Mock `MessageService.send` to verify AUTH_ERROR_EVENT emission
- Use `jest.useFakeTimers()` for timeout/backoff tests — but be careful with async/await + fake timers
- The singleton stores an instance — you may need to reset it between tests. Check if there's a way to clear it, or mock the constructor.
- AbortController needs to be available in jsdom (it is in Node 16+)

## Validation
- Run: `npx react-scripts test --watchAll=false --testPathPattern="apiService" --verbose`
- All tests pass
