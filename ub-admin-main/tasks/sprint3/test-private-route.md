# Task: Test PrivateRoute Component

## Goal
Create `src/app/components/PrivateRoute/__tests__/index.test.tsx` with tests for the auth guard.

## Context
- File to test: `src/app/components/PrivateRoute/index.tsx`
- Created in Sprint 1 — wraps routes requiring authentication
- Checks `localStorage[LocalStorageKeys.ACCESS_TOKEN]` existence
- If token exists → renders the Component with route props
- If no token → redirects to `AppPages.RootPage` (login page)
- Uses React Router v5 (`<Route>`, `<Redirect>`)

## File to Create: `src/app/components/PrivateRoute/__tests__/index.test.tsx`

### Test Cases Required

#### With valid token
- Set `localStorage[LocalStorageKeys.ACCESS_TOKEN] = 'mock-jwt'`
- Render `<PrivateRoute component={MockComponent} path="/test" />`
- Expect MockComponent to be rendered
- Expect no redirect

#### Without token
- Ensure localStorage has no ACCESS_TOKEN
- Render `<PrivateRoute component={MockComponent} path="/test" />`
- Expect redirect to login page
- Expect MockComponent NOT rendered

#### Token removal mid-session
- Start with token, verify renders
- Remove token, re-render
- Verify redirect

## IMPORTANT
- Read `PrivateRoute/index.tsx` FIRST to understand exact implementation
- Must wrap with `<Router>` (React Router v5) for Route/Redirect to work
- Use `createMemoryHistory` from 'history' with `initialEntries: ['/test']`
- Mock localStorage with `jest.spyOn(Storage.prototype, 'getItem')`
- Import `LocalStorageKeys` from actual constants to use the real key names
- Check `AppPages` enum for the actual login page path

## Validation
- Run: `npx react-scripts test --watchAll=false --testPathPattern="PrivateRoute" --verbose`
- All tests pass
