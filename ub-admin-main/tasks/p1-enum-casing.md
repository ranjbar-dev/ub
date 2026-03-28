# Task: Fix Inconsistent Enum Casing

**ID:** p1-enum-casing  
**Phase:** 1 — Type Safety Foundation  
**Severity:** 🟡 HIGH  
**Dependencies:** None  

## Problem

Multiple enums contain duplicate keys that represent the same value in different casings. This creates confusion and redundant switch cases.

## Files to Modify

### 1. `src/app/containers/Deposits/types.ts`

**Current:**
```typescript
export enum DepositStatusStrings {
  Completed = 'completed',
  InProgress = 'in progress',
  Rejected = 'reject',
  Rejectedd = 'REJECTED',      // ← typo in key name + duplicate meaning
  Confirmed = 'CONFIRMED',
  Created = 'created',
  COMPLETED = 'COMPLETED',     // ← duplicate of Completed
}
```

**Target:**
```typescript
export enum DepositStatusStrings {
  Completed = 'completed',
  CompletedUpper = 'COMPLETED',
  InProgress = 'in progress',
  Rejected = 'reject',
  RejectedUpper = 'REJECTED',
  Confirmed = 'CONFIRMED',
  Created = 'created',
}
```

> **Note:** We keep both casing values because the API may return either format. The key names are fixed to remove the typo (`Rejectedd`) and make duplicates distinguishable.

### 2. `src/app/containers/OpenOrders/types.ts`

**Current:**
```typescript
export enum Sides {
  Buy = 'buy',
  BUY = 'BUY',     // ← duplicate
  Sell = 'sell',
  SELL = 'SELL',    // ← duplicate
}
```

**Target:**
```typescript
export enum Sides {
  Buy = 'buy',
  BuyUpper = 'BUY',
  Sell = 'sell',
  SellUpper = 'SELL',
}
```

### 3. `src/utils/stylers.ts` — Update switch cases

**Current (problematic):**
```typescript
case Sides.Buy:
  return '#369452';
case Sides.BUY:
  return '#369452';
case DepositStatusStrings.Completed:
  return '#369452';
case DepositStatusStrings.COMPLETED:
  return '#369452';
case DepositStatusStrings.COMPLETED:  // ← actual duplicate case (line 21)
  return '#369452';
```

**Target — normalize input instead of listing every case:**
```typescript
export const stateStyler = (state: string): string | undefined => {
  const normalized = state.toLowerCase().replace(/_/g, ' ');

  // Green — success states
  if (['completed', 'confirmed', 'successful', 'buy'].includes(normalized)) {
    return '#369452';
  }
  // Gray — pending states
  if (['created', 'in progress', 'notconfirmed', 'incomplete'].includes(normalized)) {
    return '#B3B3B3';
  }
  // Red — failure states
  if (['rejected', 'reject', 'failed', 'sell'].includes(normalized)) {
    return '#B16567';
  }
  return undefined;
};
```

## Files That Import These Enums

Search for consumers before making changes:
```
DepositStatusStrings → src/utils/stylers.ts, src/app/containers/Deposits/*.tsx
Sides               → src/utils/stylers.ts, src/app/containers/OpenOrders/*.tsx, src/app/containers/Orders/components/*.tsx
```

## Validation

```bash
npm run checkTs   # Must pass
npm test          # Must pass
```
