# Task: Review and Remove eslint-disable Comments

**ID:** p5-eslint-disable  
**Phase:** 5 — Code Organization  
**Severity:** 🟠 MEDIUM  
**Dependencies:** None  

## Problem

~40 `eslint-disable` comments scattered across the codebase. Many were added to suppress warnings that should have been fixed instead. Some are blanket `/*eslint-disable*/` that suppress ALL rules for the entire file.

## Scope

Find all eslint-disable comments:
```bash
grep -rn "eslint-disable" src/ --include="*.ts" --include="*.tsx"
```

## Categories

### 1. Blanket file-wide disables — HIGH PRIORITY to remove
```typescript
/*eslint-disable*/  // Disables ALL rules for entire file
```

**Known locations:**
- `src/services/message_service.ts` (line 3)
- `src/app/constants.ts` (line 1)

These should be replaced with specific rule disables or the underlying issues should be fixed.

### 2. Specific rule disables — review case by case
```typescript
// eslint-disable-next-line @typescript-eslint/no-unused-vars
// eslint-disable-next-line no-console
```

Some of these are legitimate (e.g., unused parameters in callbacks required by framework signatures). Others should be fixed.

### 3. `@ts-ignore` comments — replace with `@ts-expect-error`
```typescript
//@ts-ignore   // Bad: silently ignores, stays even if the error is fixed
//@ts-expect-error — TODO: fix when X is typed   // Better: errors when no longer needed
```

## Execution Steps

1. Run `grep -rn "eslint-disable\|@ts-ignore" src/ --include="*.ts" --include="*.tsx"` to get full list
2. For each occurrence:
   a. Check if the underlying issue can be fixed (preferred)
   b. If not fixable now, replace blanket disable with specific rule name
   c. Replace `@ts-ignore` with `@ts-expect-error` + explanation
3. Remove any eslint-disable comments that are no longer needed (the rule no longer triggers)

## Validation

```bash
npm run lint      # Must pass
npm run checkTs   # Must pass
npm test          # Must pass

# Count remaining disables — should be significantly lower
grep -c "eslint-disable\|@ts-ignore" src/**/*.{ts,tsx}
```
