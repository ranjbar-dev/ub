# Task: Test Global Store Slice

## Goal
Create `src/store/__tests__/slice.test.ts` with tests for the global Redux slice.

## Context
- File to test: `src/store/slice.ts`
- Has `globalSlice` with `setIsLoggedIn` action
- `setIsLoggedIn(false)` clears auth-related localStorage keys (Sprint 1 change)
- Has `selectLoggedIn` selector using createSelector
- Initial state: `{ loggedIn: false }`
- LocalStorageKeys are imported from `services/constants`

## File to Create: `src/store/__tests__/slice.test.ts`

### Test Cases Required

#### Reducer: setIsLoggedIn
- Initial state has `loggedIn: false`
- `setIsLoggedIn(true)` → `{ loggedIn: true }`
- `setIsLoggedIn(false)` → `{ loggedIn: false }` AND clears localStorage keys
- Verify all LocalStorageKeys are removed from localStorage when set to false

#### Selector: selectLoggedIn
- Returns `state.global.loggedIn` value
- Works with `true` state
- Works with `false` state

## IMPORTANT
- Read `src/store/slice.ts` FIRST to see exact reducer logic
- Check which LocalStorageKeys are cleared — test each one
- Use `jest.spyOn(Storage.prototype, 'removeItem')` to verify localStorage clearing
- There may already be a `src/store/__tests__/` directory with other tests — check first and don't overwrite existing files

## Validation
- Run: `npx react-scripts test --watchAll=false --testPathPattern="store/__tests__/slice" --verbose`
- All tests pass
