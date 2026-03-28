# Task: Create Saga Utility Functions

**ID:** p3-saga-utils  
**Phase:** 3 — Saga Error Handling  
**Severity:** 🔴 CRITICAL  
**Dependencies:** p2-apiservice-errors  
**Blocks:** p3-saga-trycatch  

## Problem

There are 15+ saga files, each duplicating the same patterns: call API → check response.status → send MessageService event → handle loading states. There's no shared error handling utility.

## File to Create

**`src/utils/sagaUtils.ts`** (new file)

### Content
```typescript
import { call, CallEffect } from 'redux-saga/effects';
import { MessageService, MessageNames } from 'services/message_service';
import { StandardResponse } from 'services/constants';
import { ApiError } from 'services/api_service';

/**
 * Wraps an API call with standard error handling:
 * - Shows loading state before the call
 * - Hides loading state after (success or failure)
 * - Shows toast on error
 * - Returns undefined on failure (caller can check)
 *
 * @param apiFunc - The API service function to call
 * @param params - Parameters to pass to the API function
 * @param options - Optional config (loadingId, toastOnError)
 * @returns The StandardResponse on success, or undefined on failure
 *
 * @example
 * ```typescript
 * function* fetchUsersSaga(action: PayloadAction<GetUsersParams>) {
 *   const response = yield* safeApiCall(GetUsersAPI, action.payload);
 *   if (response) {
 *     MessageService.send({
 *       name: MessageNames.SET_USER_ACCOUNTS,
 *       value: response.data,
 *     });
 *   }
 * }
 * ```
 */
export function* safeApiCall<T = unknown>(
  apiFunc: (params: any) => Promise<StandardResponse<T>>,
  params: unknown,
  options: {
    loadingId?: string;
    toastOnError?: boolean;
    errorMessage?: string;
  } = {},
): Generator<CallEffect, StandardResponse<T> | undefined, StandardResponse<T>> {
  const { loadingId, toastOnError = true, errorMessage } = options;

  // Show loading
  if (loadingId) {
    MessageService.send({
      name: MessageNames.SETLOADING,
      value: true,
      loadingId,
    });
  }

  try {
    const response: StandardResponse<T> = yield call(apiFunc, params);

    if (response && response.status) {
      return response;
    }

    // API returned success HTTP but status: false
    if (toastOnError) {
      MessageService.send({
        name: MessageNames.TOAST,
        value: response?.message || errorMessage || 'Operation failed',
        type: 'error',
      });
    }
    return undefined;
  } catch (error) {
    if (error instanceof ApiError) {
      if (error.statusCode === 401) {
        // Auth error already handled by ApiService
        return undefined;
      }
      if (error.statusCode === 422 && error.errors) {
        // Validation error — send input errors to form
        MessageService.send({
          name: MessageNames.SET_INPUT_ERROR,
          value: error.errors,
        });
        return undefined;
      }
    }

    // Generic error
    if (toastOnError) {
      MessageService.send({
        name: MessageNames.TOAST,
        value: errorMessage || 'An unexpected error occurred',
        type: 'error',
      });
    }
    console.error(`Saga API call failed:`, error);
    return undefined;
  } finally {
    // Hide loading
    if (loadingId) {
      MessageService.send({
        name: MessageNames.SETLOADING,
        value: false,
        loadingId,
      });
    }
  }
}

/**
 * Shows a success toast notification.
 */
export function showSuccessToast(message: string): void {
  MessageService.send({
    name: MessageNames.TOAST,
    value: message,
    type: 'success',
  });
}

/**
 * Shows an error toast notification.
 */
export function showErrorToast(message: string): void {
  MessageService.send({
    name: MessageNames.TOAST,
    value: message,
    type: 'error',
  });
}
```

## Usage Example

**Before (typical saga pattern):**
```typescript
function* getUsersSaga(action) {
  MessageService.send({ name: MessageNames.SETLOADING, value: true });
  let response = yield call(GetUsersAPI, action.payload);
  MessageService.send({ name: MessageNames.SETLOADING, value: false });
  if (response.status) {
    MessageService.send({ name: MessageNames.SET_USER_ACCOUNTS, value: response.data });
  }
}
```

**After:**
```typescript
function* getUsersSaga(action: PayloadAction<GetUsersParams>) {
  const response = yield* safeApiCall(GetUsersAPI, action.payload);
  if (response) {
    MessageService.send({
      name: MessageNames.SET_USER_ACCOUNTS,
      value: response.data,
    });
  }
}
```

## Validation

```bash
npm run checkTs   # Must pass
npm test          # Must pass
```
