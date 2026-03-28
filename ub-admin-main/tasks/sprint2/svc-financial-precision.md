# Task: Fix Floating-Point Financial Display

## Priority: 🟠 HIGH  
## Files to Modify: `src/utils/formatters.ts`

## Problem
`CurrencyFormater()` and `Format()` use JavaScript's `Number()` for financial values, which introduces floating-point errors (e.g., `0.1 + 0.2 = 0.30000000000000004`). For a crypto exchange, even small rounding errors in displayed admin balances can cause incorrect decisions.

## Current Code

Read `src/utils/formatters.ts` and find the `CurrencyFormater` and `Format` functions. They should be around lines 14-50:

```typescript
export function CurrencyFormater(value: string | number, ...): string {
    // Uses Number() or +value for conversion
    // German locale formatting via Intl.NumberFormat or manual logic
}

export function Format(value: string | number, ...): string {
    // Similar number formatting logic
}
```

Also check `src/app/containers/Balances/index.tsx` around line 92:
```typescript
return CurrencyFormater(Number(params.data.free) + Number(params.data.locked) + '');
```

## Required Changes

### 1. Fix `CurrencyFormater` in `src/utils/formatters.ts`

Read the full function first. The key issue is any place where arithmetic is done on financial values. The fix is to:
- Keep the function signature the same
- Use `parseFloat()` + `toFixed()` for safe display formatting
- Never do arithmetic on the raw Number values for display

If the function does `value = +value.split(' ')[0]`, change to use string manipulation instead of numeric conversion where possible.

**For the Balances calculation** (`Number(free) + Number(locked)`), this arithmetic happens in the component, not the formatter. The formatter just formats a final value. The real fix is:

### 2. Add a safe financial formatter function

Add this utility to `src/utils/formatters.ts`:

```typescript
/**
 * Safely adds two financial string values and returns a formatted string.
 * Avoids floating-point arithmetic issues by using fixed-point arithmetic.
 * 
 * @param a - First value as string (e.g., "0.001234")
 * @param b - Second value as string (e.g., "0.005678")
 * @param decimals - Number of decimal places (default: 8 for crypto)
 * @returns Formatted sum as string
 */
export function safeFinancialAdd(a: string | number, b: string | number, decimals: number = 8): string {
    const numA = parseFloat(String(a)) || 0;
    const numB = parseFloat(String(b)) || 0;
    const multiplier = Math.pow(10, decimals);
    const result = (Math.round(numA * multiplier) + Math.round(numB * multiplier)) / multiplier;
    return result.toFixed(decimals);
}
```

### 3. Update Balances calculation

In `src/app/containers/Balances/index.tsx`, find the line with:
```typescript
CurrencyFormater(Number(params.data.free) + Number(params.data.locked) + '')
```

Replace with:
```typescript
CurrencyFormater(safeFinancialAdd(params.data.free, params.data.locked))
```

Import `safeFinancialAdd` from `utils/formatters`.

**IMPORTANT:** Read the actual file first — the line number may have shifted. Search for `free` and `locked` in the file.

## Validation
- Build must pass: `$env:NODE_OPTIONS='--openssl-legacy-provider'; npm run build`
- `safeFinancialAdd('0.1', '0.2')` should return `'0.30000000'` (not `0.30000000000000004`)
- Existing formatting behavior should be preserved
