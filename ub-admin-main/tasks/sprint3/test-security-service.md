# Task: Test Security Service

## Goal
Create `src/services/__tests__/securityService.test.ts` with tests for auth-related API functions.

## Context
- File to test: `src/services/securityService.ts`
- `loginAPI(params)` — POST to auth/login
- `refreshTokenAPI()` — POST to auth/refresh with stored refresh_token (Sprint 2 addition)
- Both use `ApiService.getInstance().fetchData()`

## File to Create: `src/services/__tests__/securityService.test.ts`

### Test Cases Required

#### loginAPI
- Calls fetchData with correct URL, method POST, body with credentials
- Returns StandardResponse on success
- Propagates errors from fetchData

#### refreshTokenAPI
- Reads refresh_token from localStorage
- Calls fetchData with correct URL, method POST, body with refresh_token
- Returns response on success
- Has requestName 'refresh-token' (important for infinite loop guard in apiService)

## IMPORTANT
- Read `src/services/securityService.ts` FIRST
- Mock `ApiService.getInstance().fetchData` — don't actually make HTTP calls
- Mock `localStorage.getItem` for refreshTokenAPI
- Check the exact URL paths used (may include BaseUrl prefix)
- Import `LocalStorageKeys` from constants for the correct key name
- The functions may use `ApiService.getInstance()` or call fetchData differently — read the code first

## Validation
- Run: `npx react-scripts test --watchAll=false --testPathPattern="securityService" --verbose`
- All tests pass
