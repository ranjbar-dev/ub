# Task: Implement PrivateRoute Auth Guard

## Priority: 🔴 CRITICAL (SEC-C4)
## Files to Modify: `src/app/index.tsx`, `src/store/slice.ts`
## Files to Create: `src/app/components/PrivateRoute/index.tsx`

## Problem
All routes in the app render unconditionally. The only "auth check" is `ShowSideNav()` (line 217-225) which checks `localStorage[LocalStorageKeys.ACCESS_TOKEN]` to decide sidebar visibility, but every `<Route>` from lines 247-267 mounts its component regardless of auth status. An unauthenticated user can navigate to `/Balances`, `/Withdrawals`, etc.

## Current Code

### `src/app/index.tsx` lines 247-267:
```tsx
<Switch>
    <Route exact path={AppPages.RootPage} component={LoginPage} />
    <Route path={AppPages.HomePage} component={HomePage} />
    <Route path={AppPages.UserAccounts} component={UserAccounts} />
    <Route path={AppPages.LoginHistory} component={LoginHistory} />
    <Route path={AppPages.OpenOrders} component={OpenOrders} />
    <Route path={AppPages.FilledOrders} component={FilledOrders} />
    <Route path={AppPages.ExternalOrders} component={ExternalOrders} />
    <Route path={AppPages.Deposits} component={Deposits} />
    <Route path={AppPages.Withdrawals} component={Withdrawals} />
    <Route path={AppPages.FinanceMethods} component={FinanceMethods} />
    <Route path={AppPages.CurrencyPairs} component={CurrencyPairs} />
    <Route path={AppPages.ExternalExchange} component={ExternalExchange} />
    <Route path={AppPages.MarketTicks} component={MarketTicks} />
    <Route path={AppPages.Balances} component={Balances} />
    <Route path={AppPages.ScanBlock} component={ScanBlock} />
    <Route path={AppPages.LiquidityOrders} component={LiquidityOrders} />
    <Route path={AppPages.Admins} component={Admins} />
    <Route component={NotFoundPage} />
</Switch>
```

### `src/store/slice.ts` (full file, 31 lines):
```tsx
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

**Note:** There is NO selector for `loggedIn`. You need to create one.

### `src/app/index.tsx` lines 217-225 — ShowSideNav:
```tsx
const ShowSideNav = (): boolean => {
    if (
        router.location.pathname !== AppPages.RootPage &&
        localStorage[LocalStorageKeys.ACCESS_TOKEN]
    ) {
        return true;
    }
    return false;
};
```

## Required Changes

### 1. Create `src/app/components/PrivateRoute/index.tsx`:
```tsx
import React from 'react';
import { Route, Redirect, RouteProps } from 'react-router-dom';
import { LocalStorageKeys } from 'services/constants';

interface PrivateRouteProps extends RouteProps {
  component: React.ComponentType<any>;
}

/**
 * Route wrapper that redirects to login if user has no access token.
 * Checks localStorage for ACCESS_TOKEN — mirrors the auth check used
 * throughout the app (ShowSideNav, apiService headers).
 */
const PrivateRoute: React.FC<PrivateRouteProps> = ({ component: Component, ...rest }) => (
  <Route
    {...rest}
    render={(props) =>
      localStorage[LocalStorageKeys.ACCESS_TOKEN] ? (
        <Component {...props} />
      ) : (
        <Redirect to={AppPages.RootPage} />
      )
    }
  />
);

export default PrivateRoute;
```

**IMPORTANT:** Import `AppPages` from the appropriate location. Check the existing imports in `src/app/index.tsx` for the correct import path.

### 2. Update `src/app/index.tsx` routes (lines 247-267):
- Import `PrivateRoute` from `app/components/PrivateRoute`
- Replace all `<Route path={...} component={...} />` with `<PrivateRoute path={...} component={...} />`
- KEEP `<Route exact path={AppPages.RootPage} component={LoginPage} />` as a regular Route (login page must be public)
- KEEP `<Route component={NotFoundPage} />` as regular Route (catch-all)

### 3. Add selector to `src/store/slice.ts`:
After the `globalSlice` definition, add:
```tsx
import { createSelector } from '@reduxjs/toolkit';
import { RootState } from 'types';

const selectGlobalDomain = (state: RootState) => state.global || initialState;

export const selectLoggedIn = createSelector(
  [selectGlobalDomain],
  (globalState) => globalState.loggedIn,
);
```

## Validation
- Build must pass: `$env:NODE_OPTIONS='--openssl-legacy-provider'; npm run build`
- Verify: navigating to `/Balances` without a token should redirect to `/` (login)
- Verify: LoginPage route (`/`) still works without auth
