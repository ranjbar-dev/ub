# Task: Implement Auth Error Listener + Automatic Logout on 401

## Priority: 🔴 CRITICAL
## Files to Modify: `src/app/index.tsx`
## Files to Read First: `src/services/messageService.ts` (MessageNames.AUTH_ERROR_EVENT)

## Problem
`apiService.ts` emits `AUTH_ERROR_EVENT` on every 401 response, but **nothing listens for it**. Users with expired tokens see broken UI but are never redirected to login. The PrivateRoute only checks if a token EXISTS in localStorage, not whether it's valid.

## Current Code

### `src/app/index.tsx` — existing useEffect for DOWNLOAD_FILE (lines 78-87):
```tsx
useEffect(() => {
    const subscription = Subscriber.subscribe((message: any) => {
        if (message.name === MessageNames.DOWNLOAD_FILE) {
            downloadFile(message.payload)
        }
    })
    return () => {
        subscription.unsubscribe()
    }
}, [])
```

### What happens on 401 today:
1. `apiService.ts:85-89` → sends `SETLOADING: false` + `AUTH_ERROR_EVENT` → throws `ApiError(401)`
2. `sagaUtils.ts:65-68` → catches 401, returns `undefined` (silent)
3. **Nothing happens** — user stays on current page with broken state

## Required Changes

### Add AUTH_ERROR_EVENT listener in `src/app/index.tsx`

Find the existing `useEffect` that subscribes to `DOWNLOAD_FILE` (lines 78-87). MODIFY it to also handle `AUTH_ERROR_EVENT`:

```tsx
useEffect(() => {
    const subscription = Subscriber.subscribe((message: any) => {
        if (message.name === MessageNames.DOWNLOAD_FILE) {
            downloadFile(message.payload)
        }
        if (message.name === MessageNames.AUTH_ERROR_EVENT) {
            // Clear auth state and redirect to login
            dispatch(globalActions.setIsLoggedIn(false));
        }
    })
    return () => {
        subscription.unsubscribe()
    }
}, [])
```

**IMPORTANT STEPS:**

1. Check if `dispatch` is already available in the `App` component. Look for `useDispatch` import and usage. If not present, add:
```tsx
import { useDispatch } from 'react-redux';
// Inside App component:
const dispatch = useDispatch();
```

2. Import `globalActions` from `store/slice` if not already imported:
```tsx
import { globalActions } from 'store/slice';
```

3. The `setIsLoggedIn(false)` action already clears localStorage (via the selective removal we added in Sprint 1) and sets `loggedIn: false` in Redux state.

4. After localStorage is cleared, the PrivateRoute will automatically redirect to login on next render since `localStorage[LocalStorageKeys.ACCESS_TOKEN]` will be falsy.

5. **BUT** we also need to explicitly navigate to the login page to force immediate redirect. Add a router push:
```tsx
if (message.name === MessageNames.AUTH_ERROR_EVENT) {
    dispatch(globalActions.setIsLoggedIn(false));
    dispatch(replace(AppPages.RootPage));
}
```

Check if `replace` from `connected-react-router` is already imported. If not, add:
```tsx
import { replace } from 'connected-react-router';
```

Also check if `AppPages` is already imported. Look at the top of the file.

## Validation
- Build must pass: `$env:NODE_OPTIONS='--openssl-legacy-provider'; npm run build`
- On 401 from any API call, user should be redirected to login page
- Normal non-401 flows should be unaffected
