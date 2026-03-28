# Task: Remove Unguarded console.log Statements

## Priority: 🟡 MEDIUM (SEC-M2)
## Files to Modify: 7 files

## Problem
7 `console.log` statements exist in production code without `process.env.NODE_ENV` guards. These leak implementation details, user data, and financial payloads in production browser consoles.

## Files and Exact Locations

### 1. `src/app/components/GridFilter/GridFilter.tsx` line 87:
```typescript
console.log(data.colId, ' has Filter!!!!!!!')
```
**Action:** DELETE this line entirely.

### 2. `src/app/components/GridFilter/DateFilter.tsx` line 45:
```typescript
console.log('open backup datePicker');
```
**Action:** DELETE this line entirely.

### 3. `src/app/components/CountryDropDown/CountryDropDown.tsx` line 36:
```typescript
console.log(SelectedCountry);
```
**Action:** DELETE this line entirely.

### 4. `src/app/components/SimpleGrid/SimpleGrid.tsx` line 305:
```typescript
console.log('rerender');
```
**Action:** DELETE this line entirely.

### 5. `src/app/containers/Balances/saga.ts` line 38:
```typescript
console.log(action.payload);
```
**Action:** DELETE this line entirely. This logs financial balance payloads!

### 6. `src/app/containers/ScanBlock/scanPage.tsx` line 41:
```typescript
console.log(sendingData.current)
```
**Action:** DELETE this line entirely.

### 7. `src/app/containers/UserDetails/components/EditableValue.tsx` line 24:
```typescript
console.log(edditingValue.current);
```
**Action:** DELETE this line entirely. This logs every user keystroke.

## Notes
- The `console.log` in `apiService.ts` is already properly wrapped in `process.env.NODE_ENV !== 'production'` — leave it.
- `console.error` calls in `fileDownload.ts`, `sagaUtils.ts`, and `LoginPage/saga.ts` are appropriate error logging — leave them.
- Only delete the specific `console.log(...)` line, not surrounding code.

## Validation
- Build must pass: `$env:NODE_OPTIONS='--openssl-legacy-provider'; npm run build`
- Verify: `grep -rn "console.log" src/ --include="*.ts" --include="*.tsx"` shows only the NODE_ENV-guarded one in apiService.ts
