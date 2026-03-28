# ub-admin-main Architecture

## Overview

React admin panel for the UnitedBit cryptocurrency exchange platform. Provides full back-office management of users, orders, deposits, withdrawals, balances, currencies, and verification workflows.

**Tech Stack:**
- **React** 17.0.2 with hooks
- **TypeScript** 5.4.5
- **Redux Toolkit** 1.9.7 + **Redux-Saga** 1.1.3 for state & side effects
- **Material-UI** 4.10.1 for base components
- **styled-components** 5.1.1 for custom styling
- **AG Grid** 23.2.0 (`ag-grid-react`) for data tables
- **RxJS** Subjects for `MessageService` pub/sub event bus
- **i18next** for internationalization
- **React Router** 5 with `connected-react-router`
- **redux-injectors** for dynamic reducer/saga code-splitting

---

## Directory Structure

```
src/
├── app/
│   ├── components/         # Shared reusable UI components
│   ├── containers/         # Page-level containers (one per route)
│   │   └── <Container>/
│   │       ├── index.tsx       # Main component (useInjectReducer/Saga)
│   │       ├── saga.ts         # Redux-Saga side effects (API calls)
│   │       ├── slice.ts        # Redux Toolkit slice (reducers + actions)
│   │       ├── selectors.ts    # Memoized selectors via createSelector
│   │       ├── types.ts        # TypeScript interfaces
│   │       ├── Loadable.tsx    # Lazy-load wrapper (React.lazy + Suspense)
│   │       └── components/     # Container-specific sub-components (optional)
│   ├── constants.ts        # AppPages enum, WindowTypes enum, rowHeight
│   ├── appSelectors.ts     # Root selectors (router location)
│   ├── index.tsx           # App root: router, routes, theme provider
│   ├── ForceStyles.tsx     # Global CSS injection
│   └── NewWindowContainer.tsx  # Multi-window (new tab) support
├── services/
│   ├── api_service.ts          # Singleton HTTP client (Fetch API wrapper)
│   ├── constants.ts            # RequestTypes enum, LocalStorageKeys, API base URLs
│   ├── message_service.ts      # RxJS pub/sub event bus (67+ event types)
│   ├── security_service.ts     # loginAPI
│   ├── user_management_service.ts  # User, billing, order, finance APIs
│   ├── orders_service.ts       # Order actions, deposit updates, balances
│   ├── order_management_service.ts # Liquidity/commission reports
│   ├── external_orders_service.ts  # External exchange orders & queue
│   ├── admin_reports_service.ts    # Admin comments, currency/pair updates
│   ├── global_data_service.ts  # Countries, currencies, managers (reference data)
│   ├── profile_image_service.ts    # Profile image approval/rejection
│   └── toastService.ts         # react-toastify error display helper
├── store/
│   ├── configureStore.ts   # Store setup with saga middleware + injectors enhancer
│   ├── reducers.ts         # Root reducer (router + global + injected slices)
│   └── slice.ts            # Global slice: loggedIn boolean
├── types/
│   └── RootState.ts        # Root Redux state interface (24 optional container states)
├── styles/                 # Global styles, theme definitions (light/dark)
├── locales/                # i18next translation files
├── utils/
│   ├── formatters.ts       # Currency, date, number formatters
│   ├── stylers.ts          # AG Grid cell style helpers
│   ├── loadable.tsx        # React.lazy + Suspense wrapper factory
│   └── fileDownload.js     # Browser file-download trigger
├── images/                 # Static image assets
├── index.tsx               # Entry point (ReactDOM.render)
└── serviceWorker.ts
```

---

## Routes & Containers

All routes are defined in `src/app/index.tsx`. The router uses `connected-react-router` with a Redux history object.

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
| `NotFoundPage` | `*` | 404 catch-all |
| `UserDetails` | (new window) | Per-user detail view: wallets, addresses, permissions, billing |
| `VerificationWindow` | (new window) | KYC profile-image review workflow |

> **Note:** `UserDetails` and `VerificationWindow` open in a separate browser window via `react-new-window` (triggered by `OPEN_NEW_WINDOW` message event). `NavBar`, `LanguageSwitch`, and `ThemeSwitch` are layout components, not routed containers.

---

## Data Flow

### API Call Lifecycle

```
UI event (click / mount)
  → dispatch(SliceAction)
  → Redux-Saga takeLatest
  → call(DomainService.SomeAPI, params)
  → ApiService.fetchData()  [Fetch + JWT header]
  → HTTP response
  → handleRawResponse()  [401 → logout, 422 → validation, 500 → toast]
  → StandardResponse { status: boolean, data, message }
  → MessageService.send({ name: MessageNames.SOME_EVENT, payload })
  → Container useEffect subscriber
  → local useState update → re-render
```

### State Management: Dual Pattern

The app operates two parallel state systems simultaneously:

**Redux Store** (action dispatch & saga trigger):
- Slice reducers exist primarily to hold action creators
- `useInjectReducer` / `useInjectSaga` load each container's slice lazily
- `UserAccounts` slice is globally injected at app startup
- Redux DevTools shows dispatched actions; most data does **not** live in Redux

**MessageService** (actual data transport):
- RxJS `Subject` / `ReplaySubject` / `BehaviorSubject` based event bus
- Sagas publish data via `MessageService.send()`
- Containers subscribe in `useEffect` and store data in `useState`
- Keeps local UI data invisible to Redux DevTools

This hybrid means: **sagas are triggered by Redux actions, but the response data flows through MessageService into component state**, not back into Redux slices.

### Authentication

- JWT stored in `localStorage[ACCESS_TOKEN]`
- `ApiService.setHeaders()` reads token on every request
- 401 response → `MessageService.send(MessageNames.AUTH_ERROR_EVENT)` → redirect to `/login`
- Login stores token via `security_service.loginAPI` → `POST /auth/login`
- Logout clears `localStorage` via global `setIsLoggedIn(false)` action

---

## Key Patterns

### Container Pattern (6-file convention)

Every page container follows this structure:

```
containers/MyFeature/
  index.tsx      — Component: useInjectReducer/Saga, useDispatch, useSelector, JSX
  slice.ts       — createSlice: initialState, reducers (action triggers)
  saga.ts        — takeLatest watchers, API calls, MessageService.send()
  selectors.ts   — createSelector from RootState['myFeature']
  types.ts       — MyFeatureState interface + related types
  Loadable.tsx   — loadable(lazy(() => import('./')))
```

### Service Pattern

All HTTP calls follow this chain:

```typescript
// In *_service.ts
export const GetSomethingAPI = (parameters: RequestParameters): Promise<StandardResponse> =>
  ApiService.getInstance().fetchData({ requestType: RequestTypes.GET, url: 'endpoint/', data: parameters });
```

- Services are pure functions returning `Promise<StandardResponse>`
- `StandardResponse = { status: boolean, data: any, message: string }`
- Services never interact with Redux or MessageService directly
- Sagas consume services and handle the response

### Grid Pattern

Most list pages render data through **AG Grid** via the `SimpleGrid` shared component:

```
Container saga → MessageService.send(SET_GRID_DATA)
  → Container useEffect → setState(rows)
  → <SimpleGrid rowData={rows} colDefs={cols} />
```

`SimpleGrid` wraps AG Grid with: pagination (`PaginationComponent`), tab navigation (`GridTabs`), filter system (`GridFilter`), custom cell rendering (`renderer`), and row loading states.

### Multi-Window Pattern

`UserDetails` and `VerificationWindow` open in new browser windows:

1. Main app dispatches `OPEN_NEW_WINDOW` message with `windowType` payload
2. `NewWindowContainer` (in `App.tsx`) receives event, renders `react-new-window`
3. New window re-renders the target container with its own Redux injectors
4. Windows are sized 1175×745px by default

---

## Services Reference

### `api_service.ts` — HTTP Client

Singleton class. All HTTP traffic passes through here.

| Method | Description |
|--------|-------------|
| `getInstance()` | Returns singleton instance |
| `fetchData(params)` | Main request method (GET/POST/PUT/DELETE) |
| `setHeaders()` | Builds auth headers from `localStorage[ACCESS_TOKEN]` |
| `handleRawResponse()` | 401 → logout event; 422 → validation; 500 → toast |

- **Base URL (admin):** `https://admin.unitedbit.com/api/v1/`
- **Base URL (web app):** `https://[dev-]app.unitedbit.com/api/v1/` (for countries/currencies)
- Dev mode logs every request/response to console with emoji prefixes

### `message_service.ts` — Event Bus

Three subject types:
- `Subscriber` — plain `Subject` (emit only to current subscribers)
- `RepaySubscriber3` — `ReplaySubject(3)` (replays last 3 emissions to new subscribers)
- `BehaviorSubscriber` — `BehaviorSubject` (always has current value)

Key event categories:

| Category | Example Events |
|----------|----------------|
| Loading states | `SETLOADING`, `SET_BUTTON_LOADING`, `SET_ROW_LOADING` |
| Grid data | `SET_USER_ACCOUNTS`, `SET_BILLING_DATA`, `REFRESH_GRID`, `UPDATE_GRID_ROW` |
| Modals | `OPEN_NEW_WINDOW`, `CLOSE_POPUP`, `CLOSE_REJECT_POPUP` |
| Auth | `AUTH_ERROR_EVENT` |
| Layout | `RESIZE`, `GRID_RESIZE` |
| File | `DOWNLOAD_FILE`, `DATASEND` |

### Domain Services → Endpoints

| Service File | Key Functions | Endpoints |
|---|---|---|
| `security_service.ts` | `loginAPI` | `POST /auth/login` |
| `user_management_service.ts` | `GetUserAccountsAPI`, `GetInitialUserDataAPI`, `GetUserBalancesAPI`, `GetUserWhiteAddressesAPI`, `GetUserPermissionsAPI`, `UpdateUserPermissionsAPI`, `GetBillingGridDataAPI`, `GetOpenOrdersAPI`, `GetTradeHistoryAPI`, `GetUserImagesAPI`, `GetLoginHistoryAPI`, `GetFinanceMethodsAPI`, `GetCurrencyPairsAPI`, `GetExternalExchangeAPI`, `GetMarketTicksAPI`, `UpdateUserDataAPI`, `GetWithdrawDetailAPI`, `AddPaymentCommentAPI`, `UpdateWithdrawAPI` | `/user/*`, `/payment/*`, `/order/*`, `/trade/*`, `/currency/*`, `/ohlc/*`, `/exchange/*` |
| `orders_service.ts` | `CancelOrderAPI`, `FullFillOrderAPI`, `UpdateDepositAPI`, `GetBalancesAPI`, `GetBalanceHistoryAPI`, `UpdateAllBalancesAPI`, `InternalTransferAPI` | `/order/*`, `/payment/*`, `/crypto-balance/*`, `/crypto-internal-transfer/*` |
| `external_orders_service.ts` | `GetExternalOrdersAPI`, `GetNetQueueAPI`, `GetAllQueueAPI`, `ChangeNetQueueStatusAPI`, `CancelNetQueueAPI`, `SubmitNetQueueAPI` | `/exchange/order/*`, `/exchange/aggregation/*` |
| `order_management_service.ts` | `GetLiquidityOrdersAPI`, `UpdateCommissionReportAPI` | `/exchange/order/commission-report`, `/exchange/order/update-commission-report` |
| `admin_reports_service.ts` | `AddAdminCommentAPI`, `DeleteAdminCommentAPI`, `EditAdminCommentAPI`, `UpdateFinancialMethodAPI`, `UpdateCurrencyPairAPI`, `GetCommitionsAPI` | `/user/admin-comment/*`, `/currency/update*`, `/statistic/*` |
| `global_data_service.ts` | `GetCountriesAPI`, `GetCurrenciesAPI`, `GetManagersAPI` | `[webApp]/main-data/country-list`, `[webApp]/currencies`, `/user/admins` |
| `profile_image_service.ts` | `UpdateProfileImageStatusAPI` | `/user/profile-image/update` |

---

## Shared Components Reference

Located in `src/app/components/`.

| Component | Purpose |
|-----------|---------|
| `SimpleGrid` | AG Grid wrapper with pagination, filters, tabs, cell renderers |
| `sideNav` | Left sidebar navigation with categorized menu links |
| `newWindow` / `TNewWindow` | Opens `UserDetails`/`VerificationWindow` in a new browser window |
| `UserDetailsWindow` | Content shell rendered inside new window |
| `GridFilter` | Filter bar: `CountryFilter`, `DateFilter`, `DropDown` |
| `GridTabs` | Tab navigation for segmented grid views |
| `grid_loading` | Loading spinner overlay for grids |
| `PaginationComponent` | Material-UI Pagination wrapper |
| `renderer` | AG Grid custom cell renderer (DOM node factory) |
| `materialModal` | Modal dialog with title, close button, Zoom transition |
| `ConstructiveModal` | Confirmation dialog for destructive actions |
| `inputWithValidator` | TextField with real-time validation, throttled onChange |
| `UbDropDown` | Material-UI Select with custom styling |
| `UBInput` | Material-UI TextField wrapper |
| `UbCheckBox` | Checkbox with custom icon |
| `UbTextAreaAutoHeight` | Auto-expanding textarea |
| `UbModal` | General-purpose modal wrapper |
| `RawInput` | Plain text input |
| `BackupDatePicker` | Date picker |
| `CountryDropDown` | Country selection dropdown |
| `titledContainer` | Card with header title |
| `PageWrapper` | Top-level page layout wrapper |
| `FullWidthWrapper` | Full-width layout container |
| `GridWrapper` | Layout wrapper for grid pages |
| `MainTabsWrapper` | Wrapper for tab-based layouts |
| `wrappers` | Assorted styled layout wrappers |
| `clickableEmail` | Renders email as an action button in grids |
| `LoadingIndicator` | Generic spinner |
| `isLoadingWithText` | Overlay with loading text |
| `loadingInButton` | Loading state inside a button |
| `PlaceHolder` | Empty-state placeholder |
| `FormLabel` | Form label component |
| `Link` | Custom link component |
| `Radio` | Radio button |
| `Customized` | Theme/style customization utilities |

---

## Store & State Shape

### `src/store/configureStore.ts`

```
configureStore({
  reducer: rootReducer,          // router + global + dynamic container slices
  middleware: [sagaMiddleware, routerMiddleware],
  enhancers: [createInjectorsEnhancer({ runSaga, createReducer })],
  devTools: !isProduction,
})
```

### `src/store/slice.ts` — Global Slice

```typescript
interface GlobalState { loggedIn: boolean }
// setIsLoggedIn(false) → clears localStorage (logout)
```

### `src/types/RootState.ts` — Full State Shape

```typescript
interface RootState {
  theme?:              ThemeState;
  global?:             GlobalState;
  router?:             RouterState;        // connected-react-router
  loginPage?:          LoginPageState;
  userAccounts?:       UserAccountsState;  // NOTE: always present (globally injected)
  userDetails?:        UserDetailsState;
  billing?:            BillingState;
  reports?:            ReportsState;
  orders?:             OrdersState;
  verificationWindow?: VerificationWindowState;
  loginHistory?:       LoginHistoryState;
  openOrders?:         OpenOrdersState;
  filledOrders?:       FilledOrdersState;
  externalOrders?:     ExternalOrdersState;
  deposits?:           DepositsState;
  withdrawals?:        WithdrawalsState;
  financeMethods?:     FinanceMethodsState;
  currencyPairs?:      CurrencyPairsState;
  externalExchange?:   ExternalExchangeState;
  marketTicks?:        MarketTicksState;
  admins?:             AdminsState;
  homePage?:           HomePageState;
  balances?:           BalancesState;
  scanBlock?:          ScanBlockState;
  liquidityOrders?:    LiquidityOrdersState;
}
```

All container states are **optional** because slices are injected lazily when the container mounts.

---

## App Constants (`src/app/constants.ts`)

```typescript
enum AppPages {
  RootPage      = '/',
  LoginPage     = '/login',
  UserAccounts  = '/userAccounts/',
  HomePage      = '/home',
  LoginHistory  = '/loginHistory',
  OpenOrders    = '/OpenOrders',
  Withdrawals   = '/Withdrawals',
  Deposits      = '/Deposits',
  FinanceMethods = '/FinanceMethods',
  FilledOrders  = '/FilledOrders',
  ExternalOrders = '/ExternalOrders',
  ExternalExchange = '/ExternalExchange',
  MarketTicks   = '/MarketTicks',
  Balances      = '/Balances',
  ScanBlock     = '/ScanBlock',
  Admins        = '/Admins',
  CurrencyPairs = '/CurrencyPairs',
  LiquidityOrders = '/LiquidityOrders',
  PlaceHolder   = '/PlaceHolder',
}

enum WindowTypes { User = 'user', Verification = 'verification' }
enum Buttons { SubmitButton, RedButton, BlackButton, SkyBlueButton, GreenButton,
               LightGreenButton, VeryLightGreenButton, VeryLightBlueButton, LightYellowButton }

const rowHeight = 35;  // AG Grid default row height (px)
```

---

## Development Workflow

```bash
# Install dependencies
npm install

# Start dev server (hot reload)
npm start

# Production build
npm run build

# Run tests (Jest, 98% coverage threshold)
npm test
```

### Environment

- Dev API logs are printed to the browser console with color-coded prefixes
- `IS_LOCAL=true` activates local development mode
- `localStorage[ACCESS_TOKEN]` holds the JWT; clearing it logs out the user
- Redux DevTools extension shows dispatched actions in non-production builds

---

## Adding a New Container

1. Create `src/app/containers/MyFeature/` with: `index.tsx`, `slice.ts`, `saga.ts`, `selectors.ts`, `types.ts`, `Loadable.tsx`
2. Add `MyFeatureState` to `src/types/RootState.ts`
3. Add route to `src/app/index.tsx` using the `Loadable` import
4. Add `AppPages.MyFeature = '/MyFeature'` to `src/app/constants.ts`
5. Add sidebar link in `src/app/components/sideNav/index.tsx`
