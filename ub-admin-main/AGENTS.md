# AI Agent Guide — ub-admin-main

> **Last verified against codebase:** package.json, tsconfig.json, and all source files.

## Quick Reference

### Build & Test Commands
```bash
yarn install                   # Install dependencies (preferred — npm ci has arborist bugs)
npm start                      # Dev server (port 3000, CRA default)
npm run build                  # Production build
npm test                       # Jest + React Testing Library (90% coverage threshold)
npm run checkTs                # TypeScript type-check only (no emit)
npm run lint                   # ESLint check on src/
npm run lint:fix               # ESLint auto-fix on src/
npm run lint:css               # Stylelint check (styled-components)
npm run generate               # Plop code generator — scaffold new containers interactively
npm run start:prod             # Build + serve locally (production preview)
npm run test:generators        # Validate plop generator templates
```

### Key File Locations
| What | Where |
|------|-------|
| App entry | `src/app/index.tsx` |
| HTML entry | `src/index.tsx` |
| Route enum | `src/app/constants.ts` (`AppPages` enum) |
| Route registrations | `src/app/index.tsx` (`<Switch>` block) |
| Side-nav categories | `src/app/index.tsx` (`Categories` array in `App()`) |
| Redux store setup | `src/store/configureStore.ts` |
| Root reducer | `src/store/reducers.ts` |
| Global slice | `src/store/slice.ts` (loggedIn boolean) |
| Root state type | `src/types/RootState.ts` |
| API singleton | `src/services/apiService.ts` |
| API base URL | `src/services/constants.ts` (`BaseUrl` / `appUrl`) |
| Service constants | `src/services/constants.ts` (`RequestTypes`, `LocalStorageKeys`, `StandardResponse`) |
| Pub/sub events | `src/services/messageService.ts` (`MessageService`, `MessageNames`, `Subscriber`) |
| Grid helper styles | `src/utils/stylers.ts` (`cellColorAndNameFormatter`, `stateStyler`) |
| Number/date formatters | `src/utils/formatters.ts` (`CurrencyFormater`, `queryStringer`) |
| Safe API call wrapper | `src/utils/sagaUtils.ts` (`safeApiCall`, `showSuccessToast`, `showErrorToast`) |
| Shared grid component | `src/app/components/SimpleGrid/SimpleGrid.tsx` |
| i18n config | `src/locales/i18n.ts` |
| Translations | `src/locales/en/translation.json`, `src/locales/de/translation.json` |
| Theme definitions | `src/styles/theme/` |
| Global styles | `src/styles/global-styles.ts` |
| TypeScript config | `tsconfig.json` (baseUrl: `./src` for absolute imports) |

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
   import { MessageService, MessageNames } from 'services/messageService';
   import { MyNewAPI } from 'services/myService';

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

   **Alternative saga pattern using `safeApiCall`** (recommended for new code):
   ```typescript
   import { takeLatest } from 'redux-saga/effects';
   import { MyPageActions } from './slice';
   import { MessageService, MessageNames } from 'services/messageService';
   import { safeApiCall } from 'utils/sagaUtils';
   import { MyNewAPI } from 'services/myService';

   export function* GetMyPageData(action: { type: string; payload: any }) {
     const response = yield* safeApiCall(MyNewAPI, action.payload);
     if (response) {
       MessageService.send({
         name: MessageNames.SET_MY_PAGE_DATA,
         value: response.data,
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
   import { MessageNames } from 'services/messageService';
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
   <PrivateRoute path={AppPages.MyPage} component={MyPage} />
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

7. **Add a `MessageNames` event** in `src/services/messageService.ts`:
   ```typescript
   enum MessageNames {
     // ...
     SET_MY_PAGE_DATA = 'SET_MY_PAGE_DATA',
   }
   ```

---

### Add a New API Endpoint

1. **Choose or create** a service file in `src/services/`. Name it `<domain>Service.ts` (camelCase).

2. **Add the exported function:**
   ```typescript
   import { apiService } from './apiService';
   import { RequestTypes } from './constants';

   export const MyNewAPI = (parameters: any) => {
     return apiService.fetchData({
       data: parameters,
       url: 'endpoint/path',        // appended to BaseUrl + 'admin/' automatically
       requestType: RequestTypes.GET, // GET | POST | PUT | DELETE
     });
   };
   ```
   - For **GET** requests, `parameters` is serialised as a query string by `queryStringer()`.
   - For **POST/PUT/DELETE**, `parameters` is sent as `JSON.stringify(body)`.
   - If you need a raw URL without the `admin/` prefix, pass `isRawUrl: true` — URL becomes `BaseUrl + endpoint/path`.
   - For external app URLs (e.g. `webAppAddress`), build the full URL yourself and pass `isRawUrl: true`.

3. **Call it from a saga** using `yield call()` or `safeApiCall()`:
   ```typescript
   // Option A: Manual error handling
   import { call, takeLatest } from 'redux-saga/effects';
   import { StandardResponse } from 'services/constants';
   import { MyNewAPI } from 'services/myService';
   import { MessageService, MessageNames } from 'services/messageService';

   export function* myNewSaga(action: { type: string; payload: any }) {
     const response: StandardResponse = yield call(MyNewAPI, action.payload);
     if (response.status === true) {
       MessageService.send({ name: MessageNames.SET_MY_DATA, payload: response.data });
     }
   }

   // Option B: Using safeApiCall (handles loading, errors, toasts automatically)
   import { safeApiCall } from 'utils/sagaUtils';

   export function* myNewSaga(action: { type: string; payload: any }) {
     const response = yield* safeApiCall(MyNewAPI, action.payload, {
       loadingId: 'myButton',  // optional: shows loading on specific button
       toastOnError: true,     // default: shows toast on failure
     });
     if (response) {
       MessageService.send({ name: MessageNames.SET_MY_DATA, value: response.data });
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
import { Subscriber, MessageNames } from 'services/messageService';

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
| Auth 401 errors | `apiService.ts` fires `MessageNames.AUTH_ERROR_EVENT` on 401; check `localStorage['access_token']` |
| 422 validation errors | `apiService.ts` throws `ApiError` with `errors` field — saga catches it in `safeApiCall` and sends `SET_INPUT_ERROR` |
| 500 errors | `apiService.ts` throws `ApiError` — `safeApiCall` shows error toast and returns `undefined` |
| Redux state empty | Expected — most containers store nothing in Redux (empty slices). Data lives in component state via MessageService |
| TypeScript errors | Run `npm run checkTs` to see all errors without a build |
| Build fails (OpenSSL) | Prefix: `NODE_OPTIONS=--openssl-legacy-provider npm run build` (typically not needed with Node 18) |
| API URL wrong in dev | Set `REACT_APP_API_BASE_URL` in `.env.local` to override the default production URL |

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

3. **`StandardResponse.data` is typed `unknown` by default.**
   Generic parameter `StandardResponse<T>` allows typing, but most existing code uses untyped responses. You must read the saga and the corresponding API endpoint documentation to know the actual shape.

4. **`apiService` uses native `fetch`, not axios.**
   Axios has been removed from the project entirely. `fileDownload.ts` also uses `fetch`. Do not introduce axios.

5. **`BaseUrl` uses `REACT_APP_API_BASE_URL` env var with production fallback.**
   `src/services/constants.ts`: `BaseUrl = process.env.REACT_APP_API_BASE_URL || 'https://admin.unitedbit.com/api/v1/'`. Set this in `.env.local` for local development. Note: `fetchData` prepends `admin/` to the URL unless `isRawUrl: true`.

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

11. **90% coverage threshold is enforced on `npm test -- --coverage`.**
    Tests run with `react-scripts test`. Loadable files, `*.d.ts`, `types.ts`, `index.tsx`, and `serviceWorker.ts` are excluded from coverage collection.

12. **API requests auto-prepend `admin/` to URLs.**
    `apiService.fetchData()` builds URLs as `BaseUrl + 'admin/' + params.url`. To skip the `admin/` prefix, set `isRawUrl: true` in the request parameters.

13. **CSRF protection is enabled for non-GET requests.**
    `apiService.ts` reads a CSRF token from `<meta name="csrf-token">` or localStorage and adds `X-CSRF-Token` header to POST/PUT/DELETE requests.

---

# Admin Panel — ub-admin-main

## Stack (verified from package.json)

| Dependency | Version | Purpose |
|-----------|---------|---------|
| React | 17.0.2 | UI framework |
| TypeScript | 5.4.5 | Type system |
| react-scripts (CRA) | 5.0.1 | Build toolchain (Webpack 5) |
| @reduxjs/toolkit | 1.3.6 | State management |
| redux-saga | 1.1.3 | Side effects |
| redux-injectors | 1.3.0 | Dynamic reducer/saga injection |
| @material-ui/core | 4.12.4 | Component library |
| @material-ui/icons | 4.11.3 | Icon library |
| @material-ui/lab | 4.0.0-alpha.55 | Lab components |
| @material-ui/pickers | 3.2.10 | Date pickers |
| styled-components | 5.1.1 | CSS-in-JS |
| ag-grid-community | 23.2.0 | Data grid |
| ag-grid-react | 23.2.0 | React grid wrapper |
| connected-react-router | 6.8.0 | Router ↔ Redux sync |
| react-router-dom | 5.2.0 | Client routing |
| history | 4.10.1 | History API |
| i18next | 19.4.5 | Internationalization |
| react-i18next | 11.5.0 | React i18n bindings |
| date-fns | 2.14.0 | Date utilities |
| notistack | 0.9.17 | Snackbar notifications |
| react-cropper | 1.3.0 | Image cropping |
| react-viewer | 3.2.1 | Image viewer |
| react-helmet-async | 1.0.6 | Document head manager |
| react-new-window | 0.1.2 | New window portal |
| sass | 1.43.4 | SCSS compiler (dart-sass) |
| animejs | 3.2.0 | Animation library |
| classnames | 2.2.6 | Conditional CSS classes |
| plop | 2.6.0 | Code generator |

**Runtime:** Node ≥18 (`.nvmrc`: 18), Yarn 1.22 (preferred), Docker: `node:18-slim`

## Architecture

### Container Pattern (6-file convention)
Each page follows `src/app/containers/<PageName>/`:

| File | Purpose |
|------|---------|
| `index.tsx` | Main component — `useInjectReducer`/`useInjectSaga`, hooks, JSX |
| `saga.ts` | Side effects — API calls via generator functions, `MessageService.send()` |
| `slice.ts` | Redux Toolkit slice — action creators + (usually empty) reducers |
| `selectors.ts` | Memoized selectors via `createSelector` from RootState |
| `types.ts` | TypeScript interfaces for container state (`ContainerState` alias) |
| `Loadable.tsx` | Code splitting via `lazyLoad()` wrapper (React.lazy + Suspense) |

Some containers also have:
- `components/` — container-specific sub-components
- `constants.ts` — container-specific enums/config
- `validators/` — form validation logic
- `tabPages/` — tab sub-views
- `__tests__/` — container tests

### Data Flow (Action → Saga → API → MessageService → Component)
```
UI event (click / mount)
  → dispatch(SliceAction)                    // Redux action
  → Redux-Saga takeLatest                    // Saga picks it up
  → yield call(DomainService.SomeAPI, params) // or yield* safeApiCall(...)
  → ApiService.fetchData()                   // fetch() + JWT + CSRF + retry
  → HTTP response
  → ApiError handling (401→refresh, 422→validation, 5xx→throw)
  → StandardResponse { status: boolean, data: T, message: string }
  → MessageService.send({ name, value/payload })  // RxJS Subject broadcast
  → Component useEffect → Subscriber.subscribe()
  → local useState update → re-render
```

**Key insight:** Sagas are triggered by Redux actions, but response data flows through `MessageService` into component `useState` — NOT back into Redux slices. Redux DevTools will show dispatched actions but most data is invisible to it.

### 27 Containers (`src/app/containers/`)

#### Routed Page Containers (18 routes)
| Container | Route | Purpose |
|-----------|-------|---------|
| `LoginPage` | `/` (exact) | Admin authentication (email + password, JWT) |
| `HomePage` | `/home` | Dashboard with summary cards and search |
| `UserAccounts` | `/userAccounts/` | User list grid — **globally injected** (always in Redux) |
| `LoginHistory` | `/loginHistory` | Login attempt audit trail |
| `OpenOrders` | `/OpenOrders` | Active trading orders (cancel / fulfill actions) |
| `FilledOrders` | `/FilledOrders` | Completed trade history |
| `ExternalOrders` | `/ExternalOrders` | External exchange orders with tab pages & queue |
| `Deposits` | `/Deposits` | Deposit transaction grid |
| `Withdrawals` | `/Withdrawals` | Withdrawal transaction grid |
| `FinanceMethods` | `/FinanceMethods` | Payment method configuration |
| `CurrencyPairs` | `/CurrencyPairs` | Trading pair configuration |
| `ExternalExchange` | `/ExternalExchange` | External exchange integration config |
| `MarketTicks` | `/MarketTicks` | Market OHLC data with sync page |
| `Balances` | `/Balances` | Crypto balance management with transfer modal |
| `ScanBlock` | `/ScanBlock` | Blockchain transaction scanner |
| `LiquidityOrders` | `/LiquidityOrders` | Commission / liquidity report |
| `Admins` | `/Admins` | Admin user management |
| `NotFoundPage` | `*` (catch-all) | 404 page |

#### Window Containers (opened in new browser windows)
| Container | Trigger | Purpose |
|-----------|---------|---------|
| `UserDetails` | `OPEN_NEW_WINDOW` (WindowTypes.User) | Per-user detail view: wallets, addresses, permissions, billing |
| `VerificationWindow` | `OPEN_NEW_WINDOW` (WindowTypes.Verification) | KYC profile-image review workflow |

#### Layout / Utility Containers (not routed)
| Container | Purpose |
|-----------|---------|
| `NavBar` | Top navigation bar with theme switch, language switch, logout |
| `LanguageSwitch` | Language selector (EN/DE) |
| `ThemeSwitch` | Theme selector (System/Light/Dark) |

#### Sub-containers (used within other containers, not independently routed)
| Container | Used In | Purpose |
|-----------|---------|---------|
| `Billing` | UserDetails | User billing data grid |
| `Orders` | UserDetails | User orders with tabs (Open/History/Trades) |
| `Reports` | UserDetails | User reporting and admin comments |

### 34 Shared Components (`src/app/components/`)

#### Data Display & Grid
| Component | File(s) | Purpose |
|-----------|---------|---------|
| `SimpleGrid` | `SimpleGrid.tsx` | AG Grid wrapper — pagination, filters, tabs, cell renderers, message subscriptions |
| `GridFilter` | `GridFilter.tsx`, `CountryFilter.tsx`, `DateFilter.tsx`, `DropDown.tsx` | Horizontal filter bar above AG Grid, mirrors column widths |
| `GridTabs` | `GridTabs.tsx` | Tab navigation for segmented grid views |
| `PaginationComponent` | `PaginationComponent.tsx` | Material-UI Pagination wrapper |
| `renderer` | `index.tsx` | AG Grid custom cell renderer (DOM node factory) |
| `grid_loading` | `gridLoading.tsx`, `wrapper.tsx` | Loading spinner overlay for grids |

#### Form Inputs
| Component | File(s) | Purpose |
|-----------|---------|---------|
| `UBInput` | `UBInput.tsx` | Material-UI OutlinedInput wrapper with label and end-adornment |
| `inputWithValidator` | `index.tsx` | TextField with real-time validation, throttled onChange, error display |
| `RawInput` | `RawInput.tsx` | Minimal plain HTML text input |
| `UbDropDown` | `index.tsx` | Material-UI Select dropdown |
| `UbCheckBox` | `UbCheckbox.tsx` | Styled checkbox |
| `UbTextAreaAutoHeight` | `UbTextAreaAutoHeight.tsx` | Auto-expanding textarea |
| `Radio` | `index.tsx` | Radio button group |
| `CountryDropDown` | `CountryDropDown.tsx` | Country selector dropdown |
| `BackupDatePicker` | `datePick.tsx` | Date selection component |
| `FormLabel` | `index.ts` | Styled form label |

#### Modals & Windows
| Component | File(s) | Purpose |
|-----------|---------|---------|
| `UbModal` | `UbModal.tsx` | Animated overlay modal (anime.js scale transitions, ESC to close) |
| `materialModal` | `modal.tsx` | Material-UI Dialog wrapper with Zoom transition |
| `ConstructiveModal` | `ConstructiveModal.tsx` | Confirmation dialog for destructive actions |
| `newWindow` / `TNewWindow` | `NewWindow.tsx`, `TNewWindow.tsx` | Opens content in real browser windows (portal-based) |
| `UserDetailsWindow` | `UserDetailsWindow.tsx` | Shell rendered inside new user-detail windows |

#### Navigation & Layout
| Component | File(s) | Purpose |
|-----------|---------|---------|
| `sideNav` | `index.tsx`, `mainCat.tsx` | Left sidebar navigation with collapsible categories |
| `PrivateRoute` | `index.tsx` | Route guard — redirects to login if no ACCESS_TOKEN in localStorage |
| `PageWrapper` | `index.ts` | Top-level page layout (960px centered, padded) |
| `titledContainer` | `TitledContainer.tsx` | Card with header title |
| `Link` | `index.ts` | Styled navigation link component |

#### Loading & Status
| Component | File(s) | Purpose |
|-----------|---------|---------|
| `LoadingIndicator` | `index.tsx` | SVG spinning circle animation |
| `isLoadingWithText` | `isLoadingWithText.tsx`, `isLoadingWithTextAuto.tsx` | Loading overlay with text label |
| `loadingInButton` | `loadingInButton.tsx` | Inline button loading spinner |

#### Wrapper Components (`wrappers/`)
| Component | Purpose |
|-----------|---------|
| `FullWidthWrapper` | Full-width page container with sidebar offset |
| `GridWrapper` | Grid container wrapper |
| `MainTabsWrapper` | Tabs container wrapper |
| `SegmentButtonsWrapper` | Segment buttons container |

#### Other
| Component | Purpose |
|-----------|---------|
| `clickableEmail` | Email rendered as action button in grids |
| `PlaceHolder` | Demo/placeholder withdrawal detail layout |
| `Customized` | Custom react-toastify styles |

### Store Architecture
- `store/configureStore.ts` — configures store with saga middleware + router middleware + injector enhancer
- `store/reducers.ts` — combines dynamically injected reducers with `connectRouter(history)` and `global` slice
- `store/slice.ts` — global application state (`{ loggedIn: boolean }`)
- Dynamic injection: `useInjectReducer()` / `useInjectSaga()` from `utils/redux-injectors`
- UserAccounts reducer/saga globally injected from App component (cross-page access)

### RootState Shape (24 container states)
```typescript
interface RootState {
  theme?: ThemeState;
  global?: GlobalState;           // { loggedIn: boolean }
  router?: RouterState;           // connected-react-router
  loginPage?: LoginPageState;
  userAccounts?: UserAccountsState;  // always present (globally injected)
  userDetails?: UserDetailsState;
  billing?: BillingState;
  reports?: ReportsState;
  orders?: OrdersState;
  verificationWindow?: VerificationWindowState;
  loginHistory?: LoginHistoryState;
  openOrders?: OpenOrdersState;
  filledOrders?: FilledOrdersState;
  externalOrders?: ExternalOrdersState;
  deposits?: DepositsState;
  withdrawals?: WithdrawalsState;
  financeMethods?: FinanceMethodsState;
  currencyPairs?: CurrencyPairsState;
  externalExchange?: ExternalExchangeState;
  marketTicks?: MarketTicksState;
  admins?: AdminsState;
  homePage?: HomePageState;
  balances?: BalancesState;
  scanBlock?: ScanBlockState;
  liquidityOrders?: LiquidityOrdersState;
}
```
All container states are **optional** because slices are injected lazily when the container mounts.

## API Integration

### ApiService Singleton (`src/services/apiService.ts`)
- **Pattern:** Singleton via `ApiService.getInstance()`
- **Transport:** Native `fetch` API (NOT axios)
- **Base URL:** `process.env.REACT_APP_API_BASE_URL || 'https://admin.unitedbit.com/api/v1/'`
- **URL building:** `fetchData()` prepends `BaseUrl + 'admin/'` unless `isRawUrl: true`
- **Auth:** Bearer JWT token from `localStorage[ACCESS_TOKEN]`
- **CSRF:** `X-CSRF-Token` header on POST/PUT/DELETE (from meta tag or localStorage)
- **Timeout:** 30 seconds per request (`AbortController`)
- **Retry:** GET/PUT retry up to 3 times with exponential backoff (1s, 2s, 4s, max 8s) on status codes `[408, 429, 502, 503, 504]`. POST/DELETE never retry.
- **Token refresh:** On 401, attempts `POST admin/auth/refresh` with stored refresh token. On success, retries original request. On failure, fires `AUTH_ERROR_EVENT`.
- **Dev logging:** `🚀 METHOD request to: URL` with full request/response objects

### Error Handling Flow
| Status | Behavior |
|--------|----------|
| 200 | Returns `StandardResponse<T>` |
| 401 | Attempts token refresh → retry or `AUTH_ERROR_EVENT` + redirect to login |
| 422 | Throws `ApiError` with `errors` field (validation failures) |
| 408, 429, 502-504 | Retries (GET/PUT only) with backoff |
| Other (403, 500, etc.) | Throws `ApiError` with message |

### `safeApiCall` Wrapper (`src/utils/sagaUtils.ts`)
Recommended for new saga code. Handles:
- Loading state (show/hide via `SET_BUTTON_LOADING`)
- Error toasts (automatic on failure)
- 401 errors (silently handled — already redirected by ApiService)
- 422 errors (sends `SET_INPUT_ERROR` to form fields)
- Returns `undefined` on failure (caller checks truthiness)

### Service Files (`src/services/`)
| File | Purpose | Key Endpoints |
|------|---------|---------------|
| `apiService.ts` | Singleton HTTP client | (base layer — not called directly by sagas) |
| `securityService.ts` | Authentication | `POST auth/login`, `POST auth/refresh` |
| `userManagementService.ts` | User CRUD, KYC, billing, orders | `GET/POST user/*`, `payment/*`, `order/*`, `trade/*`, `currency/*`, `ohlc/*`, `exchange/*` (~30+ endpoints) |
| `ordersService.ts` | Order actions, balances | `POST order/cancel`, `POST order/fulfill`, `GET crypto-balance`, `POST crypto-internal-transfer/create` |
| `externalOrdersService.ts` | External exchange queue | `GET exchange/order`, `GET/POST exchange/order/queue/*`, `POST exchange/aggregation/change-status` |
| `orderManagementService.ts` | Liquidity/commission | `GET exchange/order/commission-report`, `POST exchange/order/update-commission-report` |
| `adminReportsService.ts` | Admin comments, config updates | `POST user/admin-comment/*`, `POST currency/update*`, `GET statistic/user-statistic` |
| `globalDataService.ts` | Reference data (cached 1hr) | `GET [webApp]/main-data/country-list`, `GET [webApp]/currencies`, `GET admin/user/admins` |
| `profileImageService.ts` | Profile image review | `POST user/profile-image/update` |
| `messageService.ts` | RxJS pub/sub event bus | (67 event types — not HTTP) |
| `toastService.ts` | Toast notification helpers | (UI-only — not HTTP) |

### `globalDataService.ts` Caching
- 1-hour TTL cache for countries, currencies, and managers
- In-flight request deduplication (prevents duplicate calls)
- `invalidateGlobalDataCache()` for manual cache flush

## Routing
All routes defined in `src/app/index.tsx` `<Switch>` block.
- `/` → `LoginPage` (exact, public `Route`)
- All other routes use `<PrivateRoute>` (checks `localStorage[ACCESS_TOKEN]`)
- Routes defined in `src/app/constants.ts` `AppPages` enum
- `connected-react-router` syncs URL to Redux store

### AppPages Enum
```typescript
enum AppPages {
  RootPage = '/', LoginPage = '/login', UserAccounts = '/userAccounts/',
  HomePage = '/home', LoginHistory = '/loginHistory', OpenOrders = '/OpenOrders',
  Withdrawals = '/Withdrawals', Deposits = '/Deposits', FinanceMethods = '/FinanceMethods',
  FilledOrders = '/FilledOrders', ExternalOrders = '/ExternalOrders',
  ExternalExchange = '/ExternalExchange', MarketTicks = '/MarketTicks',
  Balances = '/Balances', ScanBlock = '/ScanBlock', Admins = '/Admins',
  CurrencyPairs = '/CurrencyPairs', LiquidityOrders = '/LiquidityOrders',
  PlaceHolder = '/PlaceHolder',
}
```

### Navigation Categories (sidebar)
Defined in `App()` in `src/app/index.tsx`:
1. **User Management** — User Accounts, Verification, Groups, Login History
2. **Order Management** — Open Orders, Filled Orders, External Orders, Liquidity Orders
3. **Accounting** — Deposits, Withdrawals, Balances, Scan Block
4. **Configuration** — Finance Methods, Currency Pairs, External Exchange, Market Ticks
5. **Administration** — Admins, Admin Rules (placeholder), Logs (placeholder)

## Internationalization (i18n)
- **Framework:** i18next + react-i18next + i18next-browser-languagedetector
- **Languages:** English (`src/locales/en/translation.json`), German (`src/locales/de/translation.json`)
- **Config:** `src/locales/i18n.ts` — fallback language: `en`, debug in development
- **Type-safe translations:** `src/locales/types.ts` defines `ConvertedToFunctionsType<T>` — converts JSON keys to `() => string` getters
- **Usage:** `const { t } = useTranslation(); t(translations.PageNames.Deposits())`
- **Translation keys:** `LoginPage.*`, `HomePage.*`, `PageNames.*`, `CommonTitles.*`, `i18nFeature.*`

## Testing
- **Framework**: Jest + React Testing Library + react-test-renderer
- **Setup**: `src/setupTests.ts` (jest-dom matchers, jest-styled-components, polyfills)
- **Coverage threshold**: 90% (branches, functions, lines, statements) — enforced with `--coverage`
- **Coverage exclusions**: `Loadable.{js,tsx}`, `*.d.ts`, `types.ts`, `index.tsx` (entry), `serviceWorker.ts`
- **Mocks**: `src/__mocks__/vaadin-date-picker.js` (mapped in jest config)
- **Existing tests**: `__tests__/` directories in containers, components, services, utils, store, styles, locales
- **Pattern**: Shallow rendering with `createRenderer()` + snapshot matching
- **Run**: `npm test` (interactive) or `npx react-scripts test --watchAll=false` (CI)

## Build & Dev Commands
```bash
yarn install                   # Install dependencies (preferred)
npm start                      # CRA dev server (port 3000, BROWSER=none in .env.local)
npm run build                  # Production build (GENERATE_SOURCEMAP=false in .env.production)
npm test                       # Run Jest tests (interactive watch mode)
npm run checkTs                # TypeScript type-check (no emit)
npm run lint                   # ESLint check on src/
npm run lint:fix               # ESLint auto-fix on src/
npm run lint:css               # Stylelint (styled-components)
npm run generate               # Plop code generator for new containers/components
npm run start:prod             # Build + serve -s build (local production preview)
npm run test:generators        # Validate plop generator templates
```

### Environment Variables
| File | Variables |
|------|-----------|
| `.env.example` | `REACT_APP_API_BASE_URL`, `REACT_APP_API_URL`, `REACT_APP_WEB_APP_URL` |
| `.env.local` | `BROWSER=none`, `IS_DEV=true` |
| `.env.production` | `GENERATE_SOURCEMAP=false` |

### CI/CD (`.gitlab-ci.yml`)
- **Stages:** prod-build → prod-deploy → prod-notification
- **Trigger:** `master` branch only
- **Build:** `docker build -t prod_admin_image$CI_COMMIT_SHA -f DockerfileProd .`
- **Deploy:** Copies `/usr/src/app/build` from container to host
- **Notification:** Telegram bot (success 🟢 / failure 🔴)

## Conventions

### File Naming
- Service files: `camelCaseService.ts` (e.g., `userManagementService.ts`)
- Container directories: `PascalCase` (e.g., `UserAccounts/`)
- Component directories: `camelCase` (e.g., `inputWithValidator/`)
- Type files always named `types.ts` with `ContainerState` alias

### Code Style
- ESLint extends `react-app` with `import/order` plugin (alphabetical, grouped imports)
- Stylelint for styled-components
- Husky pre-commit: `checkTs` + `lint-staged`
- `save-exact = true` in `.npmrc` (exact dependency pinning)
- Absolute imports from `src/` (tsconfig `baseUrl: "./src"`)
- Strict TypeScript (`strict: true`, `noImplicitAny`, `noImplicitReturns`)

### Import Order (enforced by ESLint)
```typescript
// 1. External packages
import React from 'react';
import { useDispatch } from 'react-redux';

// 2. Internal modules (absolute from src/)
import { apiService } from 'services/apiService';
import { MessageService } from 'services/messageService';

// 3. Parent/sibling/index
import { MyActions } from './slice';
import { mySaga } from './saga';
```

### Patterns to Follow
- Dynamic reducer/saga injection via `useInjectReducer()` and `useInjectSaga()` at top of container
- `safeApiCall()` in sagas for automatic error handling (new code)
- `MessageService.send()` for data distribution from sagas to components
- `SimpleGrid` for all AG Grid data tables
- `cellColorAndNameFormatter()` for status columns
- `CurrencyFormater()` for financial amounts
- `queryStringer()` for building GET query params
- AG Grid `rowHeight = 35` (from `src/app/constants.ts`)

## Utility Functions Reference

### `src/utils/formatters.ts`
| Function | Purpose |
|----------|---------|
| `CurrencyFormater(value)` | Formats number with thousand separators and 2 decimals |
| `Format(value)` | German locale formatting, parentheses for negatives |
| `FormatDate(YYYYMMDD)` | Converts to `YYYY-MM-DD` |
| `safeFinancialAdd(a, b, decimals)` | Floating-point-safe addition |
| `queryStringer(obj)` | Builds URL query string with encoding |
| `under(str)` | camelCase → snake_case |
| `censor(str, start, end)` | Masks characters in a range |
| `CopyToClipboard(text)` | Copies to clipboard |
| `PairFormat(pair)` | Removes dashes from pairs (ETH-USD → ETHUSD) |

### `src/utils/stylers.ts`
| Function | Purpose |
|----------|---------|
| `stateStyler(state)` | Returns hex color for status strings (green/grey/red) |
| `cellColorAndNameFormatter(field)` | AG Grid valueFormatter + cellStyle combo for status columns |

### `src/utils/sagaUtils.ts`
| Function | Purpose |
|----------|---------|
| `safeApiCall(apiFunc, params, options)` | Wraps API call with loading, error handling, toast |
| `showSuccessToast(message)` | Sends success toast via MessageService |
| `showErrorToast(message)` | Sends error toast via MessageService |

### `src/utils/hooks/`
| Hook | Purpose |
|------|---------|
| `useDimensions(options)` | ResizeObserver-based responsive sizing with breakpoints |
| `useForceUpdate()` | Simple re-render trigger |
| `useOpenWithdrawWindow()` | Subscription-based modal opener for withdrawals |

### Other Utils
| File | Purpose |
|------|---------|
| `fileDownload.ts` | `downloadFile({ url, filename })` — fetch blob + trigger browser download with auth header |
| `loadable.tsx` | `lazyLoad(importFunc, selector, opts)` — React.lazy + Suspense wrapper |
| `loading.ts` | `SetLoading({ id, loading })` — dispatches SET_BUTTON_LOADING |
| `commonUtils.ts` | `omit(key, obj)` — removes a property from an object |
| `gridUtilities/` | AG Grid helpers: `headerHider`, `ToggleDetail`, `getPageSize`, `RowWithShaddow` |
| `NW/new-window.jsx` | PureComponent for opening child browser windows (portal-based) |
| `history.ts` | Creates browser history instance for connected-react-router |
| `testUtils.tsx` | Test utilities |

## Known Issues / Technical Debt
- **90% coverage threshold** set but actual coverage may be lower — threshold enforcement means `npm test -- --coverage` can fail
- **Module-level `let timeOut`** in App — potential memory leak (window.onresize handler)
- **`useMemo` with mutable ref dep** in SimpleGrid — `[gridApi.current]` doesn't trigger re-renders
- **RTK upgrade blocked** — redux-injectors 1.3.0 has incompatible StoreEnhancer types with RTK ≥1.8
- **Yarn resolutions required** — `@types/react` and `@types/react-router` pinned in package.json
- **npm ci broken** — must use `yarn install` (npm has "Invalid Version" arborist bug)
- **AG Grid 23.x** has known CVEs — upgrade path blocked by major API changes
- **`internals/startingTemplate/`** excluded from tsconfig — template code for generator only
- **ScanBlock** is the only container missing `Loadable.tsx`
- **Some known typos in codebase:** `commitions` (commissions), `Rejectedd` (Rejected), `ALLPY_PARAMS_TO_GRID` (APPLY)

## Upgrade Roadmap

### Phase 1 ✅ (done)
- ~~node-sass removed~~ (dart-sass 1.43.4 handles SCSS)
- ~~Docker node:12 → node:18-slim~~
- ~~axios 0.21.1 → removed~~ (fileDownload.ts rewritten to use fetch())
- ~~engines updated to node >=18~~

### Phase 2 ✅ (done)
- ~~react-scripts 3.4.1 → 5.0.1~~ (webpack 5)
- ~~React 16.13.1 → 17.0.2~~
- ~~TypeScript 3.9.5 → 4.9.5~~
- ~~@testing-library/react 10 → 12.1.5~~
- ~~ajv → 8.18.0~~ (CJS module resolution fix)

### Phase 3 ✅ (done)
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
