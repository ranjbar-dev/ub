# Task: Add Try-Catch to All Saga Functions

**ID:** p3-saga-trycatch  
**Phase:** 3 — Saga Error Handling  
**Severity:** 🔴 CRITICAL  
**Dependencies:** p3-saga-utils  
**Blocks:** p6-redux-migration  

## Problem

Most saga generator functions lack try-catch blocks. When an API call fails (network error, 500, etc.), the saga crashes silently and the UI freezes in a loading state.

## Scope

All saga files in `src/app/containers/*/saga.ts` — approximately 15-20 files, each containing 1-5 generator functions.

## Known Saga Files

```
src/app/containers/UserAccounts/saga.ts
src/app/containers/UserDetails/saga.ts
src/app/containers/Billing/saga.ts
src/app/containers/Deposits/saga.ts
src/app/containers/Withdrawals/saga.ts
src/app/containers/OpenOrders/saga.ts
src/app/containers/FilledOrders/saga.ts
src/app/containers/ExternalOrders/saga.ts
src/app/containers/VerificationWindow/saga.ts
src/app/containers/Reports/saga.ts
src/app/containers/Balances/saga.ts
src/app/containers/MarketTicks/saga.ts
src/app/containers/CurrencyPairs/saga.ts
src/app/containers/FinanceMethods/saga.ts
src/app/containers/LiquidityOrders/saga.ts
src/app/containers/LoginHistory/saga.ts
src/app/containers/Admins/saga.ts
src/app/containers/ExternalExchange/saga.ts
src/app/containers/Orders/saga.ts
src/app/containers/LoginPage/saga.ts
src/app/containers/ScanBlock/saga.ts
```

## Pattern: Replace Direct Calls with `safeApiCall`

### Before (from `UserAccounts/saga.ts`)
```typescript
export function* getUsersSaga(action) {
  MessageService.send({ name: MessageNames.SETLOADING, value: true });
  let response: StandardResponse = yield call(GetUsersAPI, action.payload);
  MessageService.send({ name: MessageNames.SETLOADING, value: false });
  if (response.status) {
    MessageService.send({
      name: MessageNames.SET_USER_ACCOUNTS,
      value: response.data,
    });
  }
}
```

### After
```typescript
import { safeApiCall } from 'utils/sagaUtils';

export function* getUsersSaga(action: PayloadAction<GetUsersParams>) {
  const response = yield* safeApiCall(GetUsersAPI, action.payload, {
    loadingId: 'users',
  });
  if (response) {
    MessageService.send({
      name: MessageNames.SET_USER_ACCOUNTS,
      value: response.data,
    });
  }
}
```

### For sagas that DON'T use `safeApiCall` (manual approach)

If `safeApiCall` doesn't fit (e.g., multi-step sagas), wrap manually:
```typescript
export function* complexSaga(action: PayloadAction<SomeParams>) {
  try {
    MessageService.send({ name: MessageNames.SETLOADING, value: true });
    
    const response: StandardResponse = yield call(SomeAPI, action.payload);
    if (!response?.status) {
      MessageService.send({ name: MessageNames.TOAST, value: 'Failed', type: 'error' });
      return;
    }

    // ... complex multi-step logic ...

  } catch (error) {
    console.error('complexSaga failed:', error);
    MessageService.send({
      name: MessageNames.TOAST,
      value: 'An error occurred',
      type: 'error',
    });
  } finally {
    MessageService.send({ name: MessageNames.SETLOADING, value: false });
  }
}
```

## Execution Steps

1. For each saga file:
   a. Import `safeApiCall` from `utils/sagaUtils`
   b. Replace each `yield call(SomeAPI, params)` with `yield* safeApiCall(SomeAPI, params)`
   c. Remove manual loading MessageService calls (safeApiCall handles them)
   d. Keep the data dispatch MessageService call
   e. Add `PayloadAction<T>` type to the action parameter
2. Run `npm run checkTs` after each file
3. Run `npm test` after all files updated

## Validation

```bash
npm run checkTs   # Must pass
npm test          # Must pass
```
