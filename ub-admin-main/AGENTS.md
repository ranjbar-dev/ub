# AI Agent Guide — ub-admin-main

## Quick Reference

### Build & Test Commands
```bash
npm install                    # Install dependencies
npm start                      # Dev server (port 3000, CRA default)
npm run build                  # Production build
NODE_OPTIONS=--openssl-legacy-provider npm run build  # If build fails on older Node
npm test                       # Jest + React Testing Library (90% coverage threshold)
npm run checkTs                # TypeScript type-check only (no emit)
npm run lint                   # ESLint check on src/
npm run lint:fix               # ESLint auto-fix on src/
npm run generate               # Plop code generator — scaffold new containers interactively
```

### Key File Locations
| What | Where |
|------|-------|
| App entry | `src/app/index.tsx` |
| Route enum | `src/app/constants.ts` (`AppPages` enum) |
| Route registrations | `src/app/index.tsx` (`<Switch>` block) |
| Side-nav categories | `src/app/index.tsx` (`Categories` array in `App()`) |
| Redux store setup | `src/store/configureStore.ts` |
| Root state type | `src/types/RootState.ts` |
| API singleton | `src/services/api_service.ts` |
| API base URL | `src/services/constants.ts` (`BaseUrl` / `appUrl`) |
| Service constants | `src/services/constants.ts` (`RequestTypes`, `LocalStorageKeys`, `StandardResponse`) |
| Pub/sub events | `src/services/message_service.ts` (`MessageService`, `MessageNames`, `Subscriber`) |
| Grid helper styles | `src/utils/stylers.ts` (`cellColorAndNameFormatter`, `stateStyler`) |
| Number/date formatters | `src/utils/formatters.ts` (`CurrencyFormater`, `queryStringer`) |
| Shared grid component | `src/app/components/SimpleGrid/SimpleGrid.tsx` |

---

## How-To Guides

### Add a New Page / Container

1. **Create the directory:**
   ```
   src/app/containers/MyPage/
   ```

2. **Create 6 files** (follow the pattern in `Deposits/` or `OpenOrders/`):

   **`types.ts`** — state interface + domain enums:
   ```typescript
   export interface MyPageState { /* fields if any */ }
   export type ContainerState = MyPageState;
   ```

   **`slice.ts`** — Redux Toolkit slice (reducers are intentionally empty — they only trigger sagas):
   ```typescript
   import { PayloadAction } from '@reduxjs/toolkit';
   import { createSlice } from 'utils/@reduxjs/toolkit';
   import { ContainerState } from './types';

   export const initialState: ContainerState = {};

   const myPageSlice = createSlice({
     name: 'myPage',
     initialState,
     reducers: {
       GetMyPageDataAction(state, action: PayloadAction<any>) {},
     },
   });

   export const {
     actions: MyPageActions,
     reducer: MyPageReducer,
     name: sliceKey,
   } = myPageSlice;
   ```

   **`saga.ts`** — API calls; send results via `MessageService`, not `yield put()`:
   ```typescript
   import { call, takeLatest } from 'redux-saga/effects';
   import { MyPageActions } from './slice';
   import { StandardResponse } from 'services/constants';
   import { MessageService, MessageNames } from 'services/message_service';
   import { MyNewAPI } from 'services/my_service';

   export function* GetMyPageData(action: { type: string; payload: any }) {
     const response: StandardResponse = yield call(MyNewAPI, action.payload);
     if (response.status === true) {
       MessageService.send({
         name: MessageNames.SET_MY_PAGE_DATA,
         payload: response.data,
       });
     }
   }

   export function* myPageSaga() {
     yield takeLatest(MyPageActions.GetMyPageDataAction.type, GetMyPageData);
   }
   ```

   **`selectors.ts`** — memoized selectors:
   ```typescript
   import { createSelector } from '@reduxjs/toolkit';
   import { RootState } from 'types';
   import { initialState } from './slice';

   const selectDomain = (state: RootState) => state.myPage || initialState;

   export const selectMyPage = createSelector(
     [selectDomain],
     myPageState => myPageState,
   );
   ```

   **`index.tsx`** — main component; always inject reducer + saga at the top:
   ```tsx
   import React, { useMemo } from 'react';
   import { useInjectReducer, useInjectSaga } from 'utils/redux-injectors';
   import { MyPageReducer, sliceKey, MyPageActions } from './slice';
   import { myPageSaga } from './saga';
   import { SimpleGrid } from 'app/components/SimpleGrid/SimpleGrid';
   import { MessageNames } from 'services/message_service';
   import { useDispatch } from 'react-redux';

   export function MyPage() {
     useInjectReducer({ key: sliceKey, reducer: MyPageReducer });
     useInjectSaga({ key: sliceKey, saga: myPageSaga });
     const dispatch = useDispatch();

     const columnDefs = useMemo(() => [
       { headerName: 'ID', field: 'id', maxWidth: 100 },
       { headerName: 'Name', field: 'name' },
     ], []);

     return (
       <SimpleGrid
         containerId="myPage"
         messageName={MessageNames.SET_MY_PAGE_DATA}
         initialAction={MyPageActions.GetMyPageDataAction}
         arrayFieldName="items"
         immutableId="id"
         additionalInitialParams={{}}
         staticRows={columnDefs}
       />
     );
   }
   ```

   **`Loadable.tsx`** — lazy-load wrapper:
   ```tsx
   import React from 'react';
   import { lazyLoad } from 'utils/loadable';
   import { GridLoading } from 'app/components/grid_loading/gridLoading';

   export const MyPage = lazyLoad(
     () => import('./index'),
     module => module.MyPage,
     { fallback: <GridLoading /> },
   );
   ```

3. **Register the route enum** in `src/app/constants.ts`:
   ```typescript
   export enum AppPages {
     // ...existing entries...
     MyPage = '/MyPage',
   }
   ```

4. **Register the route** in the `<Switch>` block in `src/app/index.tsx`:
   ```tsx
   import { MyPage } from './containers/MyPage/Loadable';
   // ...
   <Route path={AppPages.MyPage} component={MyPage} />
   ```

5. **Add navigation link** in the `Categories` array inside `App()` in `src/app/index.tsx`:
   ```typescript
   {
     name: 'My Section',
     icon: <SomeIcon />,
     childs: [
       { name: 'My Page', page: AppPages.MyPage },
     ],
   }
   ```

6. **Register the state type** in `src/types/RootState.ts`:
   ```typescript
   import { MyPageState } from 'app/containers/MyPage/types';
   // ...
   export interface RootState {
     // ...existing keys...
     myPage?: MyPageState;
     // [INSERT NEW REDUCER KEY ABOVE]
   }
   ```

7. **Add a `MessageNames` event** in `src/services/message_service.ts`:
   ```typescript
   enum MessageNames {
     // ...
     SET_MY_PAGE_DATA = 'SET_MY_PAGE_DATA',
   }
   ```

---

### Add a New API Endpoint

1. **Choose or create** a service file in `src/services/`. Name it `<domain>_service.ts`.

2. **Add the exported function:**
   ```typescript
   import { apiService } from './api_service';
   import { RequestTypes } from './constants';

   export const MyNewAPI = (parameters: any) => {
     return apiService.fetchData({
       data: parameters,
       url: 'endpoint/path',        // appended to BaseUrl automatically
       requestType: RequestTypes.GET, // GET | POST | PUT | DELETE
     });
   };
   ```
   - For **GET** requests, `parameters` is serialised as a query string by `queryStringer()`.
   - For **POST/PUT/DELETE**, `parameters` is sent as `JSON.stringify(body)`.
   - If you need an absolute URL (e.g. `webAppAddress`), pass `isRawUrl: true`.

3. **Call it from a saga** using `yield call()`:
   ```typescript
   import { call, takeLatest } from 'redux-saga/effects';
   import { StandardResponse } from 'services/constants';
   import { MyNewAPI } from 'services/my_service';

   export function* myNewSaga(action: { type: string; payload: any }) {
     const response: StandardResponse = yield call(MyNewAPI, action.payload);
     if (response.status === true) {
       // distribute data via MessageService, not yield put()
       MessageService.send({ name: MessageNames.SET_MY_DATA, payload: response.data });
     }
   }
   ```

---

### Add a Column to an AG Grid Table

1. **Locate** the container's `index.tsx` and find the `useMemo` block that returns a `ColDef[]` array (look for `staticRows` or `columnDefs`).

2. **Add a column definition object:**
   ```typescript
   {
     headerName: t(translations.CommonTitles.MyColumn()),  // or plain string
     field: 'fieldName',         // must match the key in the API response row object
     maxWidth: 150,              // optional; omit to let grid auto-fit
   }
   ```

3. **For formatted/styled cells** use the helpers in `src/utils/`:
   ```typescript
   import { CurrencyFormater } from 'utils/formatters';
   import { cellColorAndNameFormatter } from 'utils/stylers';

   // Number formatting
   { field: 'amount', valueFormatter: (p: any) => CurrencyFormater(p.data.amount) }

   // Status colour + text normalisation (converts underscores to spaces, capitalises)
   { field: 'status', ...cellColorAndNameFormatter('status') }
   ```

4. **Custom cell renderer** (action buttons, icons):
   ```typescript
   import { CellRendererType } from 'app/components/renderer/...';
   { field: 'actions', cellRenderer: 'myRenderer', cellRendererParams: { ... } }
   ```

---

### Listening to MessageService Events in a Component

Use `Subscriber.subscribe()` inside a `useEffect` and unsubscribe on cleanup:
```typescript
import { Subscriber, MessageNames } from 'services/message_service';

useEffect(() => {
  const sub = Subscriber.subscribe((message: any) => {
    if (message.name === MessageNames.SET_MY_PAGE_DATA) {
      setData(message.payload);
    }
  });
  return () => sub.unsubscribe();
}, []);
```

`SimpleGrid` does this internally — you only need manual subscriptions in components that live outside the grid.

---

### Modify Form Validation

Forms use component-local state and dispatch `MessageNames.SET_INPUT_ERROR` via `MessageService`:
```typescript
MessageService.send({
  name: MessageNames.SET_INPUT_ERROR,
  errorId: 'myFieldId',
  payload: 'Error message',
});
```
Listen in `useEffect` with `Subscriber.subscribe()` to update local error state.

---

### Debugging Tips

| Problem | Where to look |
|---------|---------------|
| API request not firing | Check saga — is the `takeLatest` watcher registered in the root `*Saga()` function? |
| Data not showing in grid | Check that `messageName` prop on `SimpleGrid` matches the `MessageNames` used in `MessageService.send()` in the saga |
| Auth 401 errors | `api_service.ts` fires `MessageNames.AUTH_ERROR_EVENT` on 401; check `localStorage['access_token']` |
| 422 validation errors | `api_service.ts` returns the full JSON body on 422 — saga receives it as the response |
| 500 errors | `api_service.ts` shows a toast `'connection failed'` and returns `undefined` — guard with `if (response?.status === true)` |
| Redux state empty | Expected — most containers store nothing in Redux (empty slices). Data lives in component state via MessageService |
| TypeScript errors | Run `npm run checkTs` to see all errors without a build |
| Build fails (OpenSSL) | Prefix: `NODE_OPTIONS=--openssl-legacy-provider npm run build` |
| API URL wrong in dev | `BaseUrl` is hardcoded to `https://admin.unitedbit.com/api/v1/` — there is no `.env` override for it yet |

**Browser DevTools tips:**
- In development, every `fetchData` call logs `🚀 METHOD request to: URL` and `✅ success …` to the console with full request/response objects.
- Use Redux DevTools to trace dispatched actions, but remember that **almost no data is stored in Redux** — actions are purely triggers for sagas.
- Add a `console.log` inside `Subscriber.subscribe()` to trace all MessageService events in real time.

---

## ⚠️ Gotchas & Warnings

1. **Data flows through `MessageService`, not Redux state.**
   Sagas call `MessageService.send(...)` instead of `yield put(...)`. Inspecting the Redux store will show empty or minimal state for most containers. Always check the saga first.

2. **Empty slice reducers are intentional.**
   `GetDepositsAction(state, action) {}` — the reducer body is empty. These actions exist only to be picked up by `takeLatest` in the saga. This is normal and correct.

3. **`StandardResponse.data` is typed `any`.**
   You must read the saga and the corresponding API endpoint documentation to know the actual shape. The type system will not catch mistakes here.

4. **`apiService` uses native `fetch`, not axios.**
   `axios` is only used in `utils/fileDownload.ts` (legacy). Do not import axios for new API calls.

5. **`BaseUrl` is hardcoded to the production URL.**
   `src/services/constants.ts` contains `export const BaseUrl = 'https://admin.unitedbit.com/api/v1/'`. There is no `.env` variable that overrides it. Proxy configuration in `package.json` is also absent.

6. **`UserAccounts` reducer/saga are injected globally at the `App` level.**
   In `src/app/index.tsx`, `useInjectReducer` and `useInjectSaga` are called directly inside `App()` for `UserAccounts`. This means UserAccounts state is always available — do not re-inject it inside the `UserAccounts` container's `index.tsx` (it would be redundant).

7. **Navigation categories live in `src/app/index.tsx`, not in a sidebar component.**
   The `Categories` array defined inside `App()` drives the side-nav. To add a new nav entry, edit this array — not `sideNav/mainCat.tsx`.

8. **Three variants of `MessageService` exist.**
   - `MessageService` + `Subscriber` — plain RxJS `Subject` (fire-and-forget, no replay). **Use this in 99% of cases.**
   - `ReplayMessageService3` + `RepaySubscriber3` — buffers last 3 messages for late subscribers.
   - `BehaviorMessageService` + `BehaviorSubscriber` — always emits the last value to new subscribers.
   Mixing them unintentionally causes subtle bugs where events are replayed unexpectedly.

9. **`useMemo` with `[gridApi.current]` as dependency is incorrect React usage.**
   `gridApi.current` is a mutable ref and changing it does not trigger re-renders. This is a known issue in `SimpleGrid` — do not copy the pattern for new code.

10. **The `generate` script scaffolds containers with plop.**
    Running `npm run generate` is the fastest way to create a correctly structured container. Review and adjust the generated files — the templates may not include `MessageNames` setup.

11. **90% coverage threshold is enforced on `npm test`.**
    Tests run with `react-scripts test`. Coverage is collected with `--coverage` flag. Loadable files, `*.d.ts`, `types.ts`, `index.tsx`, and `serviceWorker.ts` are excluded from coverage collection.

---

# Admin Panel — ub-admin-main

## Stack
- **React 17.0.2** / **TypeScript 5.4.5** / CRA (react-scripts 5.0.1)
- **Redux Toolkit 1.3.6** + **Redux-Saga 1.1.3** + redux-injectors 1.3.0
- **Material-UI 4.10.1** + styled-components 5.1.1
- **AG Grid 23.2.0** (community) for all data tables
- **connected-react-router 6.8.0** + history 4.10.1
- i18next 19.4.5 (en/de), date-fns 2.14.0
- RxJS `Subject` for internal pubsub (`services/message_service.ts`)
- Docker: node:18-slim (production)

## Architecture

### Container Pattern
Each page follows `src/app/containers/<PageName>/`:
- `index.tsx` — main component with hooks
- `saga.ts` — side effects (API calls via generators)
- `slice.ts` — Redux Toolkit slice (actions + reducer)
- `selectors.ts` — memoized selectors via `createSelector`
- `types.ts` — TypeScript interfaces for container state
- `Loadable.tsx` — code splitting via `loadable()`

### 26 Page Containers
`Admins`, `Balances`, `Billing`, `CurrencyPairs`, `Deposits`, `ExternalExchange`, `ExternalOrders`, `FilledOrders`, `FinanceMethods`, `HomePage`, `LanguageSwitch`, `LiquidityOrders`, `LoginHistory`, `LoginPage`, `MarketTicks`, `NavBar`, `NotFoundPage`, `OpenOrders`, `Orders`, `Reports`, `ScanBlock`, `ThemeSwitch`, `UserAccounts`, `UserDetails`, `VerificationWindow`, `Withdrawals`

### Presentational Components
`src/app/components/` — shared UI (sideNav, ag-grid wrappers, popups, react-toastify customized)

### Store Architecture
- `store/configureStore.ts` — configures store with saga middleware + router middleware + injector enhancer
- `store/reducers.ts` — combines injected reducers with `connectRouter(history)` and `global` slice
- `store/slice.ts` — global application state
- Dynamic injection: `useInjectReducer()` / `useInjectSaga()` from `utils/redux-injectors`
- UserAccounts reducer/saga globally injected from App component (cross-page access)

## Key Services (in `src/services/`)
| Service | Purpose |
|---------|---------|
| `api_service.ts` | Singleton HTTP client (**fetch-based**, not axios). JWT from localStorage. |
| `user_management_service.ts` | User CRUD, verification, KYC |
| `orders_service.ts` | Open/filled/external order queries |
| `order_management_service.ts` | Order actions (cancel, modify) |
| `security_service.ts` | Auth (login, logout, 2FA) |
| `global_data_service.ts` | Currencies, pairs, system config |
| `admin_reports_service.ts` | Reports and analytics |
| `external_orders_service.ts` | External exchange order management |
| `profile_image_service.ts` | User profile image upload |
| `message_service.ts` | RxJS Subject pubsub (60+ event types) |
| `toastService.ts` | Toast notification helpers |

## Routes
`/login`, `/home`, `/userAccounts`, `/OpenOrders`, `/FilledOrders`, `/ExternalOrders`, `/Withdrawals`, `/Deposits`, `/FinanceMethods`, `/CurrencyPairs`, `/ExternalExchange`, `/Balances`, `/MarketTicks`, `/ScanBlock`, `/Admins`, `/LiquidityOrders`, `/Reports`, `/LoginHistory`

## Testing
- **Framework**: Jest + React Testing Library + react-test-renderer
- **Setup**: `src/setupTests.ts` (jest-dom, jest-styled-components, polyfills)
- **Coverage threshold**: 90% (branches, functions, lines, statements)
- **Existing tests**: ~25 test files in `__tests__/` directories
- **Run**: `npx react-scripts test` (or `npm test`)
- **Pattern**: Shallow rendering with `createRenderer()` + snapshot matching

## Build & Dev Commands
```bash
npm install            # Install dependencies
npm start              # CRA dev server (port 3000)
npm run build          # Production build
npm test               # Run Jest tests
npm run checkTs        # TypeScript type-check (no emit)
npm run lint           # ESLint check
npm run lint:fix       # ESLint auto-fix
npm run generate       # Plop code generator for new containers
```

## Conventions
- Dynamic reducer/saga injection via `useInjectReducer()` and `useInjectSaga()`
- Base URL configured in `src/services/constants.ts` (hardcoded `https://admin.unitedbit.com/api/v1/`)
- Auth token stored in localStorage (`access_token` key)
- AG Grid for all data tables with custom cell renderers
- Toast notifications via customized react-toastify
- SCSS used sparingly (VerificationWindow components, react-toastify)
- Sass (dart-sass) handles SCSS compilation (node-sass removed)
- `connected-react-router` syncs router state to Redux store

## Known Issues / Technical Debt
- **Hardcoded URLs** in `src/services/constants.ts` — should use environment variables
- **axios removed** — `fileDownload.js` rewritten to use `fetch()` API
- **msSaveBlob** usage in fileDownload.js — deprecated IE API, browsers have moved on
- **`let timeOut`** module-level variable in App — potential memory leak
- **90% coverage threshold** set but actual coverage may be lower
- Some containers import RxJS `Subject` for pubsub instead of using Redux (parallel state management)
- **RTK upgrade blocked** — redux-injectors 1.3.0 has incompatible StoreEnhancer types with RTK ≥1.8 (would need to replace redux-injectors entirely)
- **Yarn resolutions required** — `@types/react` and `@types/react-router` pinned via `resolutions` in package.json to prevent transitive type version conflicts
- **`internals/startingTemplate/`** excluded from tsconfig `include` — template code is for code generator, not real app code
- **npm ci broken** on this system — must use `yarn install` (npm has "Invalid Version" arborist bug)

## Upgrade Roadmap
### Phase 1 (done)
- ~~node-sass removed~~ (dart-sass 1.43.4 handles SCSS)
- ~~Docker node:12 → node:18-slim~~
- ~~axios 0.21.1 → removed~~ (fileDownload.js rewritten to use fetch())
- ~~engines updated to node >=18~~

### Phase 2 (done)
- ~~react-scripts 3.4.1 → 5.0.1~~ (webpack 5)
- ~~React 16.13.1 → 17.0.2~~
- ~~TypeScript 3.9.5 → 4.9.5~~
- ~~@testing-library/react 10 → 12.1.5~~
- ~~ajv → 8.18.0~~ (CJS module resolution fix)

### Phase 3 (done)
- ~~TypeScript 4.9.5 → 5.4.5~~
- ~~Yarn resolutions added~~ (`@types/react@17.0.91`, `@types/react-router@5.1.8`)
- ~~tsconfig: excluded `internals/startingTemplate/`~~ from compilation
- ~~Fixed import path~~ `external_ordersService` → `externalOrdersService`
- ~~Fixed WalletTypes cast~~ in FromToComponent.tsx
- RTK 1.3.6 upgrade **blocked** — redux-injectors incompatible with RTK ≥1.8

### Phase 4 (planned — moderate risk)
- @testing-library/react 12 → 14
- Update remaining @types packages

### Phase 5 (planned — high risk)
- React 17 → 18 (concurrent mode, createRoot, useId)
- Material-UI 4 → MUI 5 (@mui/material, styled-engine changes)
- connected-react-router → react-router 6 (native data APIs)
- ag-grid 23 → 31+ (major API changes)
- Replace redux-injectors to unblock RTK upgrade
