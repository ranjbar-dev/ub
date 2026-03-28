# Task: Plan Redux Migration from MessageService

**ID:** p6-redux-migration  
**Phase:** 6 — State Management Improvement  
**Severity:** 🔴 CRITICAL  
**Dependencies:** p1-empty-states, p3-saga-trycatch  

## Problem

The app uses a dual state management pattern:
1. **Redux** (intended) — but slices have empty reducers and state types
2. **RxJS MessageService** (legacy) — 67 event types, used by ~20 containers as the actual data channel

Sagas call the API, then send data via `MessageService.send()` instead of `yield put(actions.setData())`. Containers subscribe to MessageService in `useEffect` instead of using Redux `useSelector()`. This means:
- Data is invisible to Redux DevTools
- No time-travel debugging
- State is not centralized
- Data flows are impossible to trace without reading every useEffect

## Architecture: Current vs Target

### Current Data Flow
```
User Action → dispatch(actions.fetchX()) → saga → call(API) → MessageService.send({SET_X_DATA})
                                                                        ↓
                                                               Container useEffect
                                                               Subscriber.subscribe()
                                                               → if msg.name === SET_X_DATA
                                                                 → setState(msg.value)
```

### Target Data Flow
```
User Action → dispatch(actions.fetchX()) → saga → call(API) → yield put(actions.setXData(response.data))
                                                                        ↓
                                                               Redux Store updated
                                                                        ↓
                                                               Container useSelector(selectXData)
                                                               → re-renders with new data
```

## Migration Pattern (Per Container)

### Step 1: Fill the slice reducer
```typescript
// Before (src/app/containers/UserAccounts/slice.ts)
const slice = createSlice({
  name: 'userAccounts',
  initialState,
  reducers: {
    getUsersAction() {},  // empty — just triggers saga
  },
});

// After
const slice = createSlice({
  name: 'userAccounts',
  initialState,
  reducers: {
    getUsersAction(state) {
      state.isLoading = true;
    },
    setUsersData(state, action: PayloadAction<User[]>) {
      state.users = action.payload;
      state.isLoading = false;
    },
    setUsersError(state, action: PayloadAction<string>) {
      state.error = action.payload;
      state.isLoading = false;
    },
  },
});
```

### Step 2: Update the saga to dispatch Redux actions
```typescript
// Before (saga.ts)
function* getUsersSaga(action) {
  let response = yield call(GetUsersAPI, action.payload);
  if (response.status) {
    MessageService.send({
      name: MessageNames.SET_USER_ACCOUNTS,
      value: response.data,
    });
  }
}

// After
function* getUsersSaga(action: PayloadAction<GetUsersParams>) {
  const response = yield* safeApiCall(GetUsersAPI, action.payload);
  if (response) {
    yield put(actions.setUsersData(response.data));
  }
}
```

### Step 3: Update the container to use selectors
```typescript
// Before (index.tsx)
const [users, setUsers] = useState([]);
useEffect(() => {
  const sub = Subscriber.subscribe((msg: any) => {
    if (msg.name === MessageNames.SET_USER_ACCOUNTS) {
      setUsers(msg.value);
    }
  });
  return () => sub.unsubscribe();
}, []);

// After
const users = useSelector(selectUsers);
const isLoading = useSelector(selectIsLoading);
```

### Step 4: Create selectors
```typescript
// selectors.ts
import { createSelector } from '@reduxjs/toolkit';
import { RootState } from 'types';
import { initialState } from './slice';

const selectDomain = (state: RootState) =>
  state.userAccounts || initialState;

export const selectUsers = createSelector(
  [selectDomain],
  (state) => state.users,
);

export const selectIsLoading = createSelector(
  [selectDomain],
  (state) => state.isLoading,
);
```

## Containers to Migrate (Priority Order)

| Priority | Container | MessageNames Used | Complexity |
|----------|-----------|-------------------|------------|
| 1 | UserAccounts | SET_USER_ACCOUNTS | Low |
| 2 | Deposits | SET_DEPOSITS_DATA | Low |
| 3 | Withdrawals | SET_WITHDRAWALS_DATA | Low |
| 4 | OpenOrders | SET_OPEN_ORDERS_DATA, SET_OPEN_ORDERS_PAGE_DATA | Medium |
| 5 | Billing | SET_BILLING_DATA + 4 more | High |
| 6 | UserDetails | SET_WALLETS_DATA + 8 more | High |
| 7 | ExternalOrders | SET_EXTERNAL_ORDERS_DATA + 2 more | Medium |
| 8 | All remaining | Various | Varies |

## ⚠️ Important: Incremental Migration

Do NOT attempt a big-bang migration. Migrate one container at a time:
1. Both patterns (Redux + MessageService) can coexist during migration
2. Migrate a container, verify it works, commit, then move to the next
3. Only remove MessageService imports from a container after it's fully on Redux
4. The MessageService itself should remain until ALL containers are migrated

## Validation (Per Container)

```bash
npm run checkTs   # Must pass
npm test          # Must pass
# Manual: verify the container loads data correctly in browser
```
