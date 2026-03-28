# Task: Test Saga Utilities

## Goal
Create `src/utils/__tests__/sagaUtils.test.ts` with tests for the saga helper functions.

## Context
- File to test: `src/utils/sagaUtils.ts`
- `safeApiCall<T>` — Redux-saga generator wrapping API calls with error handling
- `showSuccessToast` / `showErrorToast` — send messages via MessageService
- safeApiCall handles: loading state, success/failure branching, 401/422 error types, toast notifications
- Uses `ApiError` from apiService (has `.status` property)

## File to Create: `src/utils/__tests__/sagaUtils.test.ts`

### Test Cases Required

#### safeApiCall — Success
- API returns `{ status: true, data: {...} }` → yields the data
- Shows success toast if `options.successMessage` provided
- Hides loading indicator after completion

#### safeApiCall — API Failure (status: false)
- API returns `{ status: false, message: 'error' }` → shows error toast
- Does not yield data

#### safeApiCall — 401 ApiError
- API throws ApiError with status 401 → emits AUTH_ERROR_EVENT
- Does not show error toast (auth errors handled separately)

#### safeApiCall — 422 ApiError
- API throws ApiError with status 422 → dispatches input errors
- Shows validation error toast

#### safeApiCall — Generic Error
- API throws regular Error → shows generic error toast
- Loading indicator still hidden

#### showSuccessToast
- Sends message via MessageService with success type
- Verify MessageService.send called with correct MessageName

#### showErrorToast
- Sends message via MessageService with error type
- Verify correct MessageName used

## IMPORTANT
- Read `src/utils/sagaUtils.ts` FIRST to understand exact generator flow
- Testing Redux-saga generators: use `redux-saga-test-plan` if available, OR manually step through with `.next()` / `.throw()`
- Check package.json for `redux-saga-test-plan` — if not installed, use manual generator testing:
  ```ts
  const gen = safeApiCall(mockApi, params);
  expect(gen.next().value).toEqual(put(someAction));
  ```
- Mock MessageService: `jest.mock('services/messageService')`
- Mock ApiService: the saga calls the API function directly (passed as parameter)
- Import `ApiError` from apiService to create test errors

## Validation
- Run: `npx react-scripts test --watchAll=false --testPathPattern="sagaUtils" --verbose`
- All tests pass
