# Task: Add TypeScript Types to Component Props

**ID:** p4-component-props  
**Phase:** 4 — Component Props & Documentation  
**Severity:** 🟡 HIGH  
**Dependencies:** None  

## Problem

Many shared components use `any`, `Function`, or untyped props. This prevents TypeScript from catching prop misuse and makes components opaque to AI agents.

## Scope

All shared components in `src/app/components/` plus any reusable components in container subdirectories.

## Common Patterns to Fix

### 1. `Function` type → specific callback signatures

**Before:**
```typescript
interface Props {
  onClick: Function;
  onDataChange: Function;
  onSubmit: Function;
}
```

**After:**
```typescript
interface Props {
  onClick: (event: React.MouseEvent<HTMLButtonElement>) => void;
  onDataChange: (data: Record<string, unknown>) => void;
  onSubmit: (values: FormValues) => void;
}
```

### 2. `any` props → specific types

**Before:**
```typescript
interface Props {
  data: any;
  params: any;
  gridApi: any;
}
```

**After:**
```typescript
interface Props {
  data: UserAccount[];
  params: GridFilterParams;
  gridApi: GridApi;  // from ag-grid-community
}
```

### 3. Missing `children` type

**Before:**
```typescript
const Card = ({ children, title }) => { ... }
```

**After:**
```typescript
interface CardProps {
  children: React.ReactNode;
  title: string;
}
const Card: React.FC<CardProps> = ({ children, title }) => { ... }
```

### 4. AG Grid callback props

Many components use AG Grid. Type the callbacks properly:

```typescript
import { GridReadyEvent, GridApi, ColumnApi, CellClickedEvent } from 'ag-grid-community';

interface GridProps {
  onGridReady: (event: GridReadyEvent) => void;
  onCellClicked?: (event: CellClickedEvent) => void;
}
```

## Files to Check

Search for components with `any` or `Function` props:
```bash
grep -rn ": any\|: Function" src/app/components/ --include="*.tsx" --include="*.ts"
```

Also check container-specific components:
```bash
grep -rn ": any\|: Function" src/app/containers/*/components/ --include="*.tsx" --include="*.ts"
```

## Execution Steps

1. Run grep to identify all files with `any` or `Function` props
2. For each component:
   a. Identify the actual data types from usage context
   b. Create or update the Props interface
   c. Apply types to the component function signature
3. Run `npm run checkTs` after each file batch

## Validation

```bash
npm run checkTs                                                # Must pass
npm test                                                        # Must pass
grep -c ": any\|: Function" src/app/components/**/*.{ts,tsx}   # Should be significantly reduced
```
