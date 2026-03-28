# Task: Fill Empty Container State Types

**ID:** p1-empty-states  
**Phase:** 1 — Type Safety Foundation  
**Severity:** 🟡 HIGH  
**Dependencies:** None  
**Blocks:** p6-redux-migration  

## Problem

17 containers have empty state interfaces (`export interface XState {}`). This means Redux state is effectively untyped and AI agents cannot determine what data each container manages.

## Files to Modify

Each container's `types.ts` file. Below is the pattern and suggested shapes based on the data they receive via MessageService.

### Pattern

Each empty state should be filled with fields matching the data the saga sends via `MessageService.send()`. Refer to each container's `saga.ts` to see what `MessageNames.SET_*` events carry.

### Example: `src/app/containers/Billing/types.ts`

**Current:**
```typescript
export interface BillingState { }
export type ContainerState = BillingState;
```

**Target:**
```typescript
export interface BillingState {
  billingData: Payment[] | null;
  depositsData: Payment[] | null;
  withdrawalsData: Payment[] | null;
  allTransactionsData: Payment[] | null;
  selectedPaymentDetails: PaymentDetails | null;
  commissions: unknown | null;
  isLoading: boolean;
  error: string | null;
}
export type ContainerState = BillingState;
```

### Containers to Update (17 files)

| Container | File | Data hints (from saga MessageNames) |
|-----------|------|-------------------------------------|
| Billing | `Billing/types.ts` | SET_BILLING_DATA, SET_BILLING_DEPOSITS_DATA, SET_BILLING_WITHDRAWALS_DATA, SET_BILLING_ALLTRANSACTIONS_DATA |
| Deposits | `Deposits/types.ts` | SET_DEPOSITS_DATA |
| OpenOrders | `OpenOrders/types.ts` | SET_OPEN_ORDERS_DATA, SET_OPEN_ORDERS_PAGE_DATA |
| FilledOrders | `FilledOrders/types.ts` | SET_FILLED_ORDERS_DATA, SET_FILLED_ORDERS_PAGE_DATA |
| ExternalOrders | `ExternalOrders/types.ts` | SET_EXTERNAL_ORDERS_DATA, SET_NET_QUEUE_DATA, SET_ALL_QUEUE_DATA |
| FinanceMethods | `FinanceMethods/types.ts` | SET_FINANCEMETHODS_DATA |
| CurrencyPairs | `CurrencyPairs/types.ts` | SET_CURRENCYPAIRS_DATA |
| ExternalExchange | `ExternalExchange/types.ts` | SET_EXTERNAL_EXCHANGE_DATA |
| MarketTicks | `MarketTicks/types.ts` | SET_MARKETTICKS_DATA, SET_SYNC_LIST_DATA |
| Admins | `Admins/types.ts` | (uses UserAccounts data) |
| HomePage | `HomePage/types.ts` | (dashboard cards data) |
| Balances | `Balances/types.ts` | SET_BALANCES_DATA, SET_BALANCES_HISTORY_DATA, SET_TRANSFER_MODAL_BALANCES_DATA |
| LiquidityOrders | `LiquidityOrders/types.ts` | SET_LIQUIDITY_ORDERS |
| ScanBlock | `ScanBlock/types.ts` | (scan results) |
| LoginHistory | `LoginHistory/types.ts` | SET_LOGIN_HISTORY_DATA |
| Reports | `Reports/types.ts` | SET_ADMINREPORTS_DATA, SET_REPORTS_WITHDRAWAL_COMMENTS |
| Orders | `Orders/types.ts` | SET_ORDER_HISTORY_DATA, SET_TRADE_HISTORY_DATA |
| Withdrawals | `Withdrawals/types.ts` | SET_WITHDRAWALS_DATA |

### Execution Steps

For each container:
1. Open the container's `saga.ts` — find all `MessageService.send()` calls to identify data shapes
2. Look at the `response.data` being sent to identify what the API returns
3. Define state fields matching each data payload
4. Always include `isLoading: boolean` and `error: string | null`
5. Keep existing types (Payment, PaymentDetails, etc.) and reference them

### Important Notes

- Do NOT change `slice.ts` or `saga.ts` in this task — just define the types
- The slices and sagas will be updated in Phase 6 (p6-redux-migration) to actually use these types
- This task is preparation — it makes the data shapes explicit for AI agents

## Validation

```bash
npm run checkTs   # Must pass — these are just interface definitions
```
