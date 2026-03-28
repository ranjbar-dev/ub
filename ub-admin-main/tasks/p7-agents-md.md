# Task: Create AGENTS.md for AI Agents

**ID:** p7-agents-md  
**Phase:** 7 — Documentation for AI Agents  
**Severity:** 🔴 CRITICAL  
**Dependencies:** None  

## Problem

The existing `AGENTS.md` (if any) doesn't provide actionable how-to guides for AI agents performing common tasks.

## File to Create or Update

**`AGENTS.md`** (project root)

### Content

```markdown
# AI Agent Guide — ub-admin-main

## Quick Reference

### Build & Test Commands
```bash
npm install                    # Install dependencies
npm start                      # Dev server (port 3000)
npm run build                  # Production build (needs NODE_OPTIONS=--openssl-legacy-provider)
npm test                       # Jest tests (90% coverage threshold)
npm run checkTs                # TypeScript type check (no emit)
npm run lint                   # ESLint check
npm run lint:fix               # ESLint auto-fix
```

### Key File Locations
| What | Where |
|------|-------|
| App entry | `src/app/index.tsx` |
| Routes | `src/app/constants.ts` (AppPages enum) |
| Redux store | `src/store/configureStore.ts` |
| Root state type | `src/types/RootState.ts` |
| API client | `src/services/api_service.ts` |
| API base URL | `src/services/constants.ts` (BaseUrl) |
| Global state | `src/store/slice.ts` |
| Pub/sub events | `src/services/message_service.ts` |

## How-To Guides

### Add a New Page/Container

1. Create directory: `src/app/containers/MyPage/`
2. Create 6 files:
   - `types.ts` — State interface + domain types
   - `slice.ts` — Redux Toolkit slice with reducers
   - `saga.ts` — Redux-Saga generators for API calls
   - `selectors.ts` — Memoized selectors
   - `index.tsx` — Main component with `useInjectReducer()` + `useInjectSaga()`
   - `Loadable.tsx` — Lazy-load wrapper

3. Add route in `src/app/constants.ts`:
   ```typescript
   MyPage = '/MyPage',
   ```

4. Add route in `src/app/index.tsx` router:
   ```tsx
   <Route path={AppPages.MyPage} component={MyPageLoadable} />
   ```

5. Add navigation link in `src/app/components/sideNav/mainCat.tsx`

6. Import container state in `src/types/RootState.ts`:
   ```typescript
   import { MyPageState } from 'app/containers/MyPage/types';
   // ...
   myPage?: MyPageState;
   ```

### Add a New API Endpoint

1. Choose the appropriate service file in `src/services/`
2. Add the function:
   ```typescript
   export const MyNewAPI = (parameters: MyParams) => {
     return apiService.fetchData({
       data: parameters,
       url: 'endpoint/path',
       requestType: RequestTypes.GET, // or POST
     });
   };
   ```

3. Call from a saga:
   ```typescript
   function* myNewSaga(action: PayloadAction<MyParams>) {
     const response = yield* safeApiCall(MyNewAPI, action.payload);
     if (response) {
       // Handle success
     }
   }
   ```

### Add a Column to an AG Grid Table

1. Find the container's grid component (look for `columnDefs` array)
2. Add column definition:
   ```typescript
   { headerName: 'Column Name', field: 'field_name', width: 150 }
   ```
3. If the column needs formatting, use helpers from `src/utils/stylers.ts` or `src/utils/formatters.ts`

### Modify Form Validation

Forms use inline validation. Look for `MessageNames.SET_INPUT_ERROR` handlers in the container's `useEffect`.

### Debugging Tips
- **Check MessageService events:** Add `console.log` in `Subscriber.subscribe()` callback
- **Check API calls:** Look at the saga's `yield call()` — the URL is in the service function
- **Check Redux state:** Use Redux DevTools (but note: most data flows through MessageService, not Redux)
- **Build failures:** Try `NODE_OPTIONS=--openssl-legacy-provider npm run build`

## ⚠️ Gotchas

1. **Data flows through MessageService, not Redux** — most sagas send data via `MessageService.send()`, not `yield put()`. Check the saga before looking at Redux state.
2. **Empty slice reducers are normal** — most reducers just trigger sagas, they don't store data.
3. **`StandardResponse.data` is `any`** — you must check saga code to know the actual shape.
4. **Two naming conventions coexist** — services use snake_case, containers use PascalCase.
5. **The `Deopsits` typo** — `AppPages.Deopsits` is intentionally misspelled (known issue).
```

## Validation

No compilation needed — this is documentation.
