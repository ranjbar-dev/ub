# Task: Remove PASSWORD/USERNAME from LocalStorageKeys + Selective Logout Clear

## Priority: 🔴 CRITICAL (SEC-C1) + MEDIUM (SEC-M3)
## Files to Modify: `src/services/constants.ts`, `src/store/slice.ts`

## Problem
1. `LocalStorageKeys` enum contains `USERNAME` and `PASSWORD` entries — signals architecture designed for plaintext credential storage. Any future code using these would create a critical vulnerability.
2. `localStorage.clear()` on logout is a "nuclear option" that clears ALL localStorage data.

## Current Code

### `src/services/constants.ts` lines 14-36:
```typescript
export enum LocalStorageKeys {
	ACCESS_TOKEN = 'access_token',
	REFRESH_TOKEN = 'refresh_token',
	USERNAME = 'username',          // ⚠️ Line 17 - REMOVE
	PASSWORD = 'password',          // ⚠️ Line 18 - REMOVE
	CURRENCIES = 'currencies',
	Theme = 'theme',
	COUNTRIES = 'countries',
	Managers = 'Managers',
	LAYOUT_NAME = 'ln',
	PAIRS = 'pairs',
	SELECTED_COIN = 'selectedCoin',
	FUND_PAGE = 'fp',
	SHOW_TOP_INFO = 'sti',
	TRADELAYOUT = 'tl',
	FAV_PAIRS = 'fps',
	FAV_COIN = 'fc',
	SHOW_FAVS = 'sf',
	TIME_FRAME = 'timeframe',
	CHANNEL = 'chan',
	VERIFICATION_WINDOW_TYPE = 'VERIFICATION_WINDOW_TYPE',
	VISIBLE_ORDER_SECTION = 'vos',
}
```

**Verified:** `LocalStorageKeys.PASSWORD` and `LocalStorageKeys.USERNAME` have ZERO usages in the codebase.

### `src/store/slice.ts` lines 18-22:
```typescript
setIsLoggedIn(state, action: PayloadAction<boolean>) {
    state.loggedIn = action.payload;
    if (action.payload === false) {
        localStorage.clear();  // ⚠️ Nuclear option
    }
},
```

## Required Changes

### 1. Remove dangerous keys from `src/services/constants.ts`:
Delete lines 17-18:
```
USERNAME = 'username',
PASSWORD = 'password',
```

### 2. Replace `localStorage.clear()` in `src/store/slice.ts`:
Replace:
```typescript
localStorage.clear();
```
With selective clearing that iterates over known keys:
```typescript
// Clear only app-specific keys, preserve unrelated localStorage data
Object.values(LocalStorageKeys).forEach((key) => {
    localStorage.removeItem(key);
});
```

**IMPORTANT:** You need to import `LocalStorageKeys` into slice.ts:
```typescript
import { LocalStorageKeys } from 'services/constants';
```

## Validation
- Build must pass: `$env:NODE_OPTIONS='--openssl-legacy-provider'; npm run build`
- Verify: `LocalStorageKeys.PASSWORD` and `LocalStorageKeys.USERNAME` no longer exist
- Verify: logout still clears auth tokens and app data
