# Task: Remove Commented-Out Code Blocks

**ID:** p5-remove-comments  
**Phase:** 5 — Code Organization  
**Severity:** 🟢 LOW  
**Dependencies:** None  

## Problem

Multiple files contain large blocks of commented-out code. This clutters files, confuses AI agents, and the code is always available in git history.

## Known Locations

### 1. `src/services/security_service.ts` (lines 12–67)
55 lines of commented-out functions. (Also covered by p2-dead-code)

### 2. `src/utils/loadable.tsx` (lines 31–54)
23 lines of commented-out alternative loadable implementation:
```typescript
/*

import React, { lazy, Suspense } from 'react';

interface Props {
  fallback: React.ReactNode | null;
}
const loadable = <T extends React.ComponentType<any>>(
  ...
*/
```

### 3. `src/utils/formatters.ts` (lines 86–98)
12 lines of commented-out `columnResize` function:
```typescript
// export const columnResize = (data: {
//   gridColumnApi: any;
//   resizeLimit: number;
// }) => {
//   ...
// };
```

### 4. `src/utils/stylers.ts` (lines 61–63)
Commented-out code in `cellColorAndNameFormatter`:
```typescript
// else if (tmpName.includes('reject') || tmpName.includes('cancel')) {
//  tmpName = tmpName + 'ed';
//}
```

### 5. Other files — search for large comment blocks
```bash
# Find blocks of 3+ consecutive commented lines
grep -n "^[ \t]*//" src/**/*.{ts,tsx} | ...
```

## Execution Steps

1. For each known location, delete the commented-out code block
2. Search for any other large comment blocks (3+ consecutive lines starting with `//`)
3. Preserve legitimate comments (JSDoc, explanatory notes, TODOs)
4. Run `npm run checkTs` and `npm test`

## Validation

```bash
npm run checkTs   # Must pass
npm test          # Must pass
```
