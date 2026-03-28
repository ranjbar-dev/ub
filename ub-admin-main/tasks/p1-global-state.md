# Task: Type RootState.global

**ID:** p1-global-state  
**Phase:** 1 — Type Safety Foundation  
**Severity:** 🟡 HIGH  
**Dependencies:** None  

## Problem

`RootState.global` is typed as `any`, hiding the shape of the global Redux state.

## Files to Modify

### 1. `src/store/slice.ts`

**Current:**
```typescript
import { PayloadAction } from '@reduxjs/toolkit';
import { createSlice } from 'utils/@reduxjs/toolkit';
interface ContainerState {
  loggedIn: boolean;
}
// The initial state of the LoginPage container
export const initialState: ContainerState = {
  loggedIn: false,
};

const globalSlice = createSlice({
  name: 'global',
  initialState,
  reducers: {
    setIsLoggedIn(state, action: PayloadAction<boolean>) {
      state.loggedIn = action.payload;
      if (action.payload === false) {
        localStorage.clear();
      }
    },
  },
});

export const {
  actions: globalActions,
  reducer: globalReducer,
  name: sliceKey,
} = globalSlice;
```

**Target — export the state type:**
```typescript
import { PayloadAction } from '@reduxjs/toolkit';
import { createSlice } from 'utils/@reduxjs/toolkit';

export interface GlobalState {
  loggedIn: boolean;
}

export const initialState: GlobalState = {
  loggedIn: false,
};

const globalSlice = createSlice({
  name: 'global',
  initialState,
  reducers: {
    setIsLoggedIn(state, action: PayloadAction<boolean>) {
      state.loggedIn = action.payload;
      if (action.payload === false) {
        localStorage.clear();
      }
    },
  },
});

export const {
  actions: globalActions,
  reducer: globalReducer,
  name: sliceKey,
} = globalSlice;
```

### 2. `src/types/RootState.ts`

**Current (line 36):**
```typescript
global?: any;
```

**Target:**
```typescript
import { GlobalState } from 'store/slice';
// ...
global?: GlobalState;
```

## Validation

```bash
npm run checkTs   # Must pass
npm test          # Must pass
```
