# Task: Make StandardResponse Generic

**ID:** p1-generic-response  
**Phase:** 1 — Type Safety Foundation  
**Severity:** 🔴 CRITICAL  
**Dependencies:** None  
**Blocks:** p2-type-services, p2-apiservice-errors  

## Problem

`StandardResponse.data` is typed as `any`, meaning every API response loses all type information. This is the root cause of type-unsafe data flow throughout the entire app.

## File to Modify

**`src/services/constants.ts`**

### Current Code (lines 37–42)
```typescript
export interface StandardResponse {
	status: boolean;
	message: string;
	data: any;
	token?: string;
}
```

### Target Code
```typescript
export interface StandardResponse<T = unknown> {
	status: boolean;
	message: string;
	data: T;
	token?: string;
	errors?: Record<string, string[]>;
}
```

## Why This Is Safe

The default generic parameter `T = unknown` means all existing code that uses `StandardResponse` without a type parameter continues to compile. New code can opt-in:

```typescript
// Existing code — still works (data is 'unknown')
const response: StandardResponse = yield call(SomeAPI, params);

// New code — typed (data is UserListResponse)
const response: StandardResponse<UserListResponse> = yield call(GetUsersAPI, params);
```

## Affected Consumers

Every saga file imports and uses `StandardResponse`. A grep for usage:
```
src/app/containers/Billing/saga.ts
src/app/containers/UserAccounts/saga.ts
src/app/containers/UserDetails/saga.ts
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
```

None of them need to change — they all work with the default `unknown` parameter.

## Validation

```bash
npm run checkTs   # Must pass — backward compatible change
npm test          # Must pass all existing tests
```
