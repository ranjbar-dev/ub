# Client Cabinet — ub-client-cabinet-main

> **Spec-driven AI reference** — This document is the single source of truth for an AI agent
> to understand, navigate, and modify this codebase. Every pattern, service, and convention
> is documented here with enough detail to produce working code without asking questions.

---

## a) Stack (verified from package.json)

| Layer | Package | Version | Notes |
|-------|---------|---------|-------|
| **UI Framework** | react / react-dom | ^18.3.1 | Runtime is React 18 (createRoot), `@types/react` pinned to 17.0.2 |
| **Language** | typescript | ^5.4.5 | Local tarball `typescript-5.4.5.tgz`; target ES2017, strict mode |
| **State** | redux | 4.0.5 | Dynamic reducer/saga injection |
| **Side Effects** | redux-saga | 1.1.3 | Saga per container, 3 injection modes |
| **Routing** | react-router-dom / redux-first-history | 5.2.0 / ^5.2.0 | URL ↔ state sync |
| **Forms** | redux-form | ^8.3.6 | Deprecated; migrate to react-hook-form |
| **i18n** | react-intl | 2.9.0 | v2 with addLocaleData (en, de, ar) |
| **Design System** | @material-ui/core | ^4.11.0 | MUI 4 + icons, lab, styles packages |
| **Styled** | styled-components | 5.2.0 | Plus styled-theming for theme variants |
| **Data Grid** | @ag-grid-enterprise/all-modules | ^25.3.0 | AG Grid Enterprise + ag-grid-react |
| **Charts** | @amcharts/amcharts4 | ^4.10.9 | amCharts 4 |
| **TradingView** | charting_library (vendored) | ^1.0.2 | Local `app/charting_library/` + `app/datafeeds/` |
| **Real-time** | mqtt (via mqtt.js) | n/a | MQTT over WSS on port 8443 |
| **HTTP** | fetch (native) + axios | ^0.21.0 | Fetch for REST, Axios only for file uploads |
| **Reactive** | rxjs | ^6.6.3 | Subject/ReplaySubject pub-sub message bus |
| **Auth** | js-cookie / universal-cookie | ^2.2.1 / ^4.0.4 | JWT in cookies |
| **reCAPTCHA** | react-google-recaptcha / react-google-recaptcha-v3 | ^2.1.0 / ^1.7.1 | v2 + v3 |
| **Build** | webpack | 4.44.1 | CRA-ejected + custom internals configs |
| **Server** | express | 4.17.1 | Production SPA server with gzip |
| **Test** | jest / ts-jest / @testing-library/react | ^29.7.0 / ^29.2.0 / ^14.3.1 | 98% coverage threshold |
| **Node** | — | >=18 | `.nvmrc`: lts/hydrogen (18.x); Docker: node:18-slim |
| **Package Manager** | yarn | 1.22.22 | `yarn.lock` present |

---

## b) Project Structure

```
ub-client-cabinet-main/
├── app/
│   ├── app.tsx                  # Entry: createRoot, Provider, Router, LanguageProvider, OnlineStatusProvider
│   ├── configureStore.ts        # Redux store: sagaMiddleware + routerMiddleware + DevTools
│   ├── reducers.ts              # Root combineReducers (global, language, router) + dynamic merge
│   ├── i18n.ts                  # react-intl v2 setup: en/de/ar locales, addLocaleData
│   ├── global-styles.ts         # CSS custom properties (40+ vars), dark theme, Vaadin overrides
│   ├── bundle.js                # Webpack bundle entry
│   ├── index.css                # Base CSS reset
│   ├── polyfills.js             # Browser polyfills
│   │
│   ├── containers/              # 21 page-level containers (see § c, d)
│   │   ├── AcountPage/          #   each: index.tsx, reducer.ts, saga.ts, actions.ts,
│   │   ├── AddressManagementPage/  #   constants.ts, selectors.ts, types.d.ts, Loadable.tsx
│   │   ├── App/                 #   + optional: components/, steps/, pages/ subfolders
│   │   ├── ChangePassword/
│   │   ├── ChangeUserInfoPage/
│   │   ├── ContactUsPage/
│   │   ├── DocumentVerificationPage/
│   │   ├── EmailVerification/
│   │   ├── FundsPage/
│   │   ├── GoogleAuthenticationPage/
│   │   ├── HomePage/
│   │   ├── LanguageProvider/
│   │   ├── LocaleToggle/
│   │   ├── LoginPage/
│   │   ├── NotFoundPage/
│   │   ├── OrdersPage/
│   │   ├── PhoneVerificationPage/
│   │   ├── RecapchaContainer/
│   │   ├── SignupPage/
│   │   ├── TradePage/
│   │   └── UpdatePasswordPage/
│   │
│   ├── components/              # 38 reusable UI components (see § e)
│   │   ├── AnimateChildren/     # Fade-in animation wrapper
│   │   ├── AppHeader/           # Fixed navigation header
│   │   ├── BreadCrumb/          # Breadcrumb navigation
│   │   ├── ConfirmPopup/        # Confirmation modal dialog
│   │   ├── Customized/          # Vendored react-toastify (custom CSS/SCSS)
│   │   ├── dataRow/             # Labeled data row display
│   │   ├── filterableSelect/    # Virtualized filterable dropdown (react-window)
│   │   ├── Footer/              # App footer with version
│   │   ├── gridFilters/         # AG Grid filter bar (date, pair, type)
│   │   ├── gridTest/            # ag-grid-react wrapper
│   │   ├── grid_loading/        # CSS ellipsis loader
│   │   ├── icons/               # SVG icon components
│   │   ├── imageWithPlaceHolder/# Image with fallback on error
│   │   ├── Img/                 # Type-safe img element
│   │   ├── inputWithValidator/  # Validated text input (throttled, error display)
│   │   ├── isLoadingWithText/   # Loading/text toggle
│   │   ├── loadingInButton/     # Small circular spinner
│   │   ├── LoadingIndicator/    # General loading indicator
│   │   ├── materialModal/       # MUI Dialog with Zoom transition
│   │   ├── miniTitledComponent/ # Compact titled card
│   │   ├── noRows/              # Animated empty state
│   │   ├── pinInput/            # 6-digit PIN input
│   │   ├── pulsingButton/       # Button with pulse animation
│   │   ├── registeredToastContainer/ # WebSocket order notification toasts
│   │   ├── renderer/            # ag-grid React cell renderer bridge
│   │   ├── RXLoader/            # MessageService-driven page loader
│   │   ├── securityLevelBar/    # Color-coded security progress bar
│   │   ├── shimmer/             # Shimmer loading placeholder
│   │   ├── splash/              # Splash screen (Lottie animation)
│   │   ├── streamLoadingButton/ # Button with stream loading state
│   │   ├── tawk/                # Tawk.to live chat widget
│   │   ├── textLoader/          # Text loading placeholder
│   │   ├── titled/              # Titled card wrapper
│   │   ├── Toggle/              # Toggle switch
│   │   ├── ToggleOption/        # Radio-style option
│   │   ├── twoFaAndVerificationPopup/ # 2FA and email verification popup
│   │   ├── verticalSpacer/      # Vertical spacing element
│   │   ├── vertivalAlignedWrapper/ # Vertical centering wrapper
│   │   └── wrappers/            # Animated appearance wrappers
│   │
│   ├── services/                # 17 service files (see § f)
│   ├── hooks/                   # Custom React hooks
│   │   ├── onlineStatusHook/    #   Online/offline detection provider
│   │   └── useImmer.tsx         #   Immer-based state hook + reducer hook
│   ├── types/                   # TypeScript type definitions
│   │   ├── index.d.ts           #   InjectedStore, ApplicationRootState (21 container states)
│   │   ├── jest-dom.d.ts        #   Jest DOM matchers
│   │   ├── react-intl.d.ts      #   react-intl type extensions
│   │   └── styled-components-fix.d.ts
│   ├── utils/                   # Utilities
│   │   ├── injectReducer.tsx    #   HOC + useInjectReducer hook for dynamic injection
│   │   ├── injectSaga.tsx       #   HOC + useInjectSaga hook with lifecycle modes
│   │   ├── reducerInjectors.ts  #   Reducer registry management
│   │   ├── sagaInjectors.ts     #   Saga registry (RESTART_ON_REMOUNT, DAEMON, ONCE_TILL_UNMOUNT)
│   │   ├── loadable.tsx         #   React.lazy + Suspense code-splitting wrapper
│   │   ├── history.ts           #   Browser history + redux-first-history setup
│   │   ├── request.ts           #   Fetch API wrapper with ResponseError
│   │   ├── environment.ts       #   isDevelopment flag (NODE_ENV + IS_DEV_BUILD)
│   │   ├── constants.ts         #   Saga injection mode constants
│   │   ├── formatters.ts        #   Number/currency/date formatting
│   │   ├── validators/          #   Input validation (login, signup, password, phone)
│   │   ├── gridUtilities/       #   AG Grid helpers
│   │   ├── storage.ts           #   LocalStorage/SessionStorage wrapper
│   │   ├── sharedData.ts        #   Cross-component singleton data
│   │   ├── throttle.ts          #   Throttle/debounce utilities
│   │   ├── decimal.js           #   Precise decimal math
│   │   ├── immerUtils.ts        #   Immer state mutation helpers
│   │   └── logger.ts            #   Debug logging
│   ├── translations/            # i18n JSON files
│   │   ├── en.json              #   English (default)
│   │   ├── de.json              #   German
│   │   └── ar.json              #   Arabic (RTL)
│   ├── styles/                  # Global theme/style definitions
│   ├── images/                  # Static images, icons, SVGs
│   ├── tests/                   # Root-level tests (i18n.test.ts, store.test.ts)
│   ├── fonts/                   # Open Sans font family
│   ├── datafeeds/               # TradingView UDF datafeed implementation
│   └── charting_library/        # TradingView charting lib (vendored)
│
├── craConfig/                   # CRA-ejected build scripts
│   ├── config/                  #   env.js, paths.js, webpack.config.js, webpackDevServer.config.js
│   └── scripts/                 #   start.js, build.js, test.js
│
├── internals/
│   ├── webpack/                 #   webpack.base.babel.js, webpack.dev.babel.js, webpack.prod.babel.js
│   ├── generators/              #   Plop generators: component/, container/, language/
│   ├── mocks/                   #   Jest mocks: cssModule.js, image.js
│   ├── scripts/                 #   analyze.js, clean.js, extract-intl.js, npmcheckversion.js
│   ├── templates/               #   Snapshot/test templates
│   └── testing/                 #   test-bundler.js (Jest setup)
│
├── server/                      # Express production server
│   ├── index.js                 #   Entry: static serving, gzip, optional ngrok tunnel
│   ├── middlewares/             #   Frontend middleware (SPA routing)
│   ├── logger.js                #   Server logging
│   ├── port.js                  #   Port config
│   └── argv.js                  #   CLI argument parsing
│
├── public/                      # HTML template + vendor polyfills
├── maintenancePage/             # Static maintenance fallback page
├── Dockerfile                   # Dev image (node:18-slim + yarn build-dev)
├── DockerfileProd               # Prod image (node:18-slim + yarn build)
├── package.json                 # Dependencies, scripts, engine constraints
├── tsconfig.json                # TypeScript config (ES2017, strict, baseUrl: ./app)
├── jest.config.js               # Jest config (ts-jest, 98% threshold)
├── .eslintrc.json               # ESLint config (@typescript-eslint, react, prettier)
├── tslint.json                  # TSLint config (legacy, deprecated)
├── .prettierrc                  # Prettier: single quotes, semicolons, trailing commas
├── .nvmrc                       # Node version: lts/hydrogen (18.x)
├── .env.example                 # Environment variable template
├── UPGRADE_PLAN.md              # 8-phase upgrade roadmap
└── unused.md                    # Deprecated file tracking list
```

---

## c) Container Pattern

Every page-level container follows a strict file convention in `app/containers/<PageName>/`:

| File | Purpose | Required |
|------|---------|----------|
| `index.tsx` | Main component; calls `useInjectReducer({ key, reducer })` and `useInjectSaga({ key, saga })` on mount; uses `useSelector` or `connect()` for state; dispatches actions | Yes |
| `reducer.ts` | Immer-based Redux reducer handling the container's slice of state | Yes |
| `saga.ts` | Redux-Saga generator functions making API calls via service modules; uses `call`, `put`, `takeLatest` | Yes |
| `actions.ts` | Action creator functions returning typed action objects | Yes |
| `constants.ts` | Action type string constants namespaced as `'ContainerName/ACTION_NAME'` | Yes |
| `selectors.ts` | Memoized `reselect` selectors accessing `ApplicationRootState` | Yes |
| `types.d.ts` | `ContainerState` interface + action payload types | Yes |
| `Loadable.tsx` | `React.lazy(() => import('./index'))` + Suspense fallback for code splitting | Yes |
| `messages.ts` | react-intl `defineMessages` for i18n strings used in this container | Most |
| `components/` | Sub-components specific to this container (e.g., TradePage/components/TradeChart/) | Optional |
| `steps/` or `pages/` | Multi-step flows (e.g., ChangePassword/steps/, FundsPage/pages/) | Optional |

### Dynamic Injection Flow
```
1. Route matches → Loadable.tsx lazy-loads container
2. Container mounts → useInjectReducer({ key: 'pageName', reducer })
   → Adds reducer to store.injectedReducers, calls store.replaceReducer()
3. Container mounts → useInjectSaga({ key: 'pageName', saga })
   → Runs saga via store.runSaga(), tracks in store.injectedSagas
4. Container renders → useSelector(makeSelectXxx()) reads from state.pageName
5. User interaction → dispatch(action()) → reducer updates state → saga handles side effects
```

### Saga Injection Modes (from `utils/constants.ts`)
- `RESTART_ON_REMOUNT` (default) — Saga cancelled on unmount, restarted on next mount
- `DAEMON` — Saga runs forever, never cancelled
- `ONCE_TILL_UNMOUNT` — Saga runs once, cancelled on unmount

### How to Add a New Container
```bash
# 1. Use the Plop generator
yarn generate   # Select "container", enter name, answer prompts

# 2. Or manually create all required files:
mkdir app/containers/MyNewPage
# Create: index.tsx, reducer.ts, saga.ts, actions.ts, constants.ts,
#         selectors.ts, types.d.ts, Loadable.tsx, messages.ts

# 3. Register the route in app/containers/App/index.tsx
# 4. Add the state type to app/types/index.d.ts → ApplicationRootState
# 5. Add i18n messages to app/translations/*.json
```

---

## d) All 21 Containers — Detailed Reference

### App (Root Container)
- **Path:** `app/containers/App/`
- **Purpose:** Root application shell — routing, auth guards, global setup
- **Key Features:** All Route definitions, PrivateRoute HOC, MUI theme provider, reCAPTCHA provider, AppHeader/Footer conditionally rendered, analytics init, mobile redirect
- **Redux State:** `global` (AppState) — loggedIn, theme, currencies, countries
- **MQTT:** `useConnectToAuthorizedMqtt2()` hook connects authenticated users to MQTT
- **API:** None directly (delegates to child containers)
- **Routes defined:**
  - `/` → HomePage, `/login` → LoginPage, `/signup` → SignupPage
  - `/trade` → TradePage, `/funds` → FundsPage, `/orders` → OrdersPage
  - `/account` → AcountPage, `/change-password` → ChangePassword
  - `/phone-verification` → PhoneVerificationPage, `/google-authentication` → GoogleAuthenticationPage
  - `/document-verification` → DocumentVerificationPage, `/userInfo` → ChangeUserInfoPage
  - `/address-management` → AddressManagementPage, `/contactus` → ContactUsPage
  - `/auth/verify` → EmailVerification, `/auth/forgot-password/update` → UpdatePasswordPage

### HomePage
- **Purpose:** Redirect-only — routes to `/trade` (logged in) or `/login` (not logged in)
- **Redux State:** None (checks cookie token)
- **API:** None

### LoginPage
- **Purpose:** User login form with email/password
- **Key Features:** Validates existing token → redirects if logged in; LoginBody sub-component handles form
- **Redux State:** `loginPage` — form state, loading, errors
- **API:** `POST auth/login` via `security_service.loginAPI()`

### SignupPage
- **Purpose:** Multi-step user registration flow
- **Key Features:** StepSelector sub-component manages step progression
- **Redux State:** `signupPage` — step, loading, form data
- **API:** `POST auth/register` via `security_service.registerAPI()`

### TradePage
- **Purpose:** Main trading interface — market data, order placement, TradingView charts
- **Key Features:** StreamComponentsWrapper for real-time data, reuses OrdersPage reducer/saga
- **Redux State:** `ordersPage` (inherited), `tradePage` — currencies, pair details
- **MQTT:** Real-time via StreamComponentsWrapper (MarketWatch, OrderBook, MarketTrade, TradeChart topics)
- **API:** `GET currencies/pairs-list`, `POST order/create`, `POST order/stop-order-create`

### FundsPage
- **Purpose:** Balance, deposits, withdrawals, transaction history
- **Key Features:** Tabbed UI (Balance/Deposit/Withdrawals/Transactions), top info panel, payment event listener
- **Redux State:** `fundsPage` — balanceData, currencies, transactions
- **MQTT:** `useNewPaymentEvent()` hook for real-time payment notifications
- **API:** `GET user-balance/balance`, `GET user-balance/withdraw-deposit`, `GET crypto-payment`, `POST crypto-payment/withdraw`, `POST crypto-payment/pre-withdraw`

### OrdersPage
- **Purpose:** Open orders, order history, trade history
- **Key Features:** Tabbed UI (Open Orders/Order History/Trade History), real-time updates
- **Redux State:** `ordersPage` — openOrders, orderHistory, tradeHistory + loading flags
- **MQTT:** RegisteredUserSubscriber for real-time order updates
- **API:** `GET order/open-orders`, `GET order/full-history`, `GET trade/full-history`, `GET order/detail`, `POST order/cancel`

### AcountPage
- **Purpose:** User account and security settings dashboard
- **Key Features:** User info display (ID, email, phone), security level bar, KYC status, links to 2FA/email/phone/password settings
- **Redux State:** `acountPage` — userData; `global` — loggedIn
- **API:** `GET user/user-data`, `POST user/send-verification-email`

### ChangePassword
- **Purpose:** Multi-step password change flow
- **Redux State:** `changePassword` — step, loading
- **API:** `POST user/change-password` via `security_service.changePasswordAPI()`

### ChangeUserInfoPage
- **Purpose:** Update personal profile information (name, address, DOB)
- **Redux State:** `changeUserInfoPage` — userProfileData, isLoading
- **API:** `GET user/get-user-profile`, `POST user/set-user-profile`

### PhoneVerificationPage
- **Purpose:** 5-step phone verification: enter phone → SMS code → password → 2FA → done
- **Redux State:** `phoneVerificationPage` — step, phone, code, loading
- **API:** `POST user/sms-send`, `POST user/sms-enable`

### GoogleAuthenticationPage
- **Purpose:** Enable/disable Google Authenticator 2FA
- **Key Features:** QR code display, code verification, enable/disable toggle
- **Redux State:** `googleAuthenticationPage` — qrCode, isLoading, userData
- **API:** `GET user/google-2fa-barcode`, `POST user/google-2fa-enable`, `POST user/google-2fa-disable`

### DocumentVerificationPage
- **Purpose:** KYC document upload (identity proof + address proof)
- **Key Features:** Dual upload panels (ProofOfIdentity), modal alerts for feedback
- **Redux State:** `documentVerificationPage` — userProfileData, isLoading
- **MessageService:** OPEN_ALERT for modal popups
- **API:** `POST user-profile-image/upload`, `POST user-profile-image/multiple-upload`

### AddressManagementPage
- **Purpose:** Manage cryptocurrency withdrawal addresses
- **Key Features:** Create new addresses, filter/search, address grid display
- **Redux State:** `AddressManagementPage` — currencies, addresses, isLoading
- **API:** `GET currencies`, `GET withdraw-address`, `POST withdraw-address/new`, `POST withdraw-address/delete`, `POST withdraw-address/favorite`

### EmailVerification
- **Purpose:** Email verification via activation code from URL query params
- **Redux State:** `emailAuthentication` — uses router location for code
- **API:** `POST auth/verify`

### UpdatePasswordPage
- **Purpose:** Password reset via email link (forgot-password flow)
- **Redux State:** `updatePasswordPage` — uses router location for reset token
- **API:** `POST auth/forgot-password`, `POST auth/forgot-password/update`

### ContactUsPage
- **Purpose:** Demo/test container (counter logic, API test, OrderList display)
- **Redux State:** `contactUsPage` — counter, inputValue
- **MessageService:** SET_LOADING_TEST, SET_LOADING_END

### RecapchaContainer
- **Purpose:** Google reCAPTCHA v3 provider wrapper
- **Key Features:** Fetches site key from API, wraps app with GoogleReCaptchaProvider
- **Redux State:** `recapchaContainer`; site key stored in sessionStorage
- **MessageService:** RESET_SITE_KEY, SET_SITE_KEY
- **API:** `GET main-data/common` (for reCAPTCHA site key)

### LanguageProvider
- **Purpose:** Redux-connected IntlProvider bridging Redux locale state to react-intl
- **Redux State:** `language` — locale string
- **API:** None

### LocaleToggle
- **Purpose:** Language selector dropdown (en, de, ar)
- **Redux State:** `language` — locale
- **Actions:** `changeLocale(locale)` from LanguageProvider

### NotFoundPage
- **Purpose:** 404 error page with return-to-login button
- **Redux State:** None

---

## e) All Components — Reference

| Component | Purpose | Key Props |
|-----------|---------|-----------|
| `AnimateChildren` | Fade-in + scale animation wrapper; shows GridLoading while loading | `isLoading`, `children`, `memoize?` |
| `AppHeader` | Fixed navigation header with tabs (Trade, Funds, Orders, Account), theme toggle, layout selector, exit button | None (Redux internal) |
| `BreadCrumb` | Clickable breadcrumb trail with i18n labels | `links: {pageName, pageLink, last?}[]` |
| `ConfirmPopup` | Modal confirmation dialog with cancel/submit buttons | `isOpen`, `onClose`, `onCancelClick`, `onSubmitClick`, `title`, `submitTitle`, `cancelTitle` |
| `Customized` | Vendored react-toastify with custom CSS/SCSS styling | (library wrapper) |
| `dataRow` | Labeled row with title:value, loading skeleton support | `title`, `value`, `boldValue?`, `small?`, `isLoading?`, `dense?`, `clickAddress?` |
| `filterableSelect` | Dropdown with virtualized list (react-window FixedSizeList) | `list`, `fieldName`, `onSelect`, `hasImage?` |
| `Footer` | Copyright and version display | `...props` |
| `gridFilters` | Multi-filter bar for AG Grid: date range, currency pair, buy/sell, coin, DW type | `onSearchClick`, `onCancelClick`, `hideTick?`, `TimePeriod?`, `CurrencyPair?`, `BuySell?`, `Coin?` |
| `gridTest` | Thin wrapper around `AgGridReact` | All ag-grid props |
| `grid_loading` | CSS ellipsis loading animation (4 dots) | `...props` |
| `icons` | SVG icon components (expandMore, etc.) | Varies |
| `imageWithPlaceHolder` | Image with fallback src on load error | `src`, `fallbackSrc`, `...imgProps` |
| `Img` | Type-safe `<img>` enforcing alt text | `src`, `alt?`, `className?` |
| `inputWithValidator` | Text input with throttled onChange, error display via MessageService, password toggle | `label`, `uniqueName`, `onChange`, `throttleTime?`, `inputType?`, `onEnter?`, `isPickable?` |
| `isLoadingWithText` | Loading spinner or text with opacity transition | `isLoading`, `text` |
| `loadingInButton` | Small (14px) white CircularProgress for buttons | None |
| `LoadingIndicator` | Memoized GridLoading wrapper | None |
| `materialModal` | MUI Dialog with Zoom transition (300ms), close X button | `isOpen`, `children`, `onClose` |
| `miniTitledComponent` | Compact card with uppercase title header and body | `title`, `children` |
| `noRows` (AnimatedNoRows) | Animated empty-state display with icon/image + text lines | `image?`, `texts`, `icon?`, `isMini?` |
| `pinInput` | 6-digit numeric PIN input with paste support | `onComplete`, `onChange?`, `onEnter?` |
| `pulsingButton` | Button with 4-ring pulse animation, MessageService loading integration | `title`, `onClick` |
| `registeredToastContainer` | Subscribes to RegisteredUserSubscriber for order notification toasts | None |
| `renderer` | Bridges React components into AG Grid cell rendering via createRoot | `children`, `styles?` |
| `RXLoader` | MessageService-driven page loader (SET_PAGE_LOADING_WITH_ID) | `id`, `style?` |
| `securityLevelBar` | Color-coded progress bar for security level | Security data props |
| `shimmer` | CSS shimmer loading placeholder | Dimension props |
| `splash` | Splash screen with Lottie animation (lottie-web) | None |
| `streamLoadingButton` | Button with streaming/loading state integration | Action props |
| `tawk` | Tawk.to live chat widget integration | Config props |
| `textLoader` | Text placeholder while loading | Style props |
| `titled` | Titled card wrapper (larger than miniTitledComponent) | `title`, `children` |
| `Toggle` | Toggle switch component | Toggle state props |
| `ToggleOption` | Radio-style selectable option | Option props |
| `twoFaAndVerificationPopup` | 2FA and email code verification popup dialog | Verification props |
| `verticalSpacer` | Vertical spacing element | `height?` |
| `vertivalAlignedWrapper` | Vertical centering flexbox wrapper | `children` |
| `wrappers` | Animated appearance wrappers (animatedApear) | `children` |

---

## f) All 17 Services — Detailed Reference

### api_service.ts — REST API Facade (Singleton)
- **Pattern:** Singleton via `ApiService.getInstance()`
- **Methods:**
  - `fetchData(params: RequestParameters)` — Execute HTTP request (GET/PUT/POST/DELETE)
  - `handleRawResponse(response)` — Parse JSON, detect auth errors (401/403), log in dev
  - `setHeaders()` — Build `Authorization: Bearer {token}` header from cookie
  - `retryWithNewToken()` — Token refresh (commented out)
- **Auth:** Bearer token from `cookies.get(CookieKeys.Token)`
- **Error Flow:** 401/403 → publishes `AUTH_ERROR_EVENT` via MessageService → redirect to login
- **CORS:** `mode: 'cors'`, `credentials: 'omit'`

### security_service.ts — Authentication
- **Methods + Endpoints:**
  - `loginAPI(data)` → `POST auth/login`
  - `registerAPI(data)` → `POST auth/register`
  - `getUserDataAPI()` → `GET user/user-data`
  - `getNewVerificationEmailAPI()` → `POST user/send-verification-email`
  - `changePasswordAPI(data)` → `POST user/change-password`
  - `set2FaAPI(data)` → `POST user/google-2fa-enable` or `POST user/google-2fa-disable`
  - `getRecapchaKeyAPI()` → `GET main-data/common`
  - `acountActivationAPI(data)` → `POST auth/verify`
  - `forgotPasswordAPI(data)` → `POST auth/forgot-password`
  - `resetPasswordAPI(data)` → `POST auth/forgot-password/update`

### user_acount_service.ts — User Profile
- **Methods + Endpoints:**
  - `getCountriesAPI()` → `GET main-data/country-list`
  - `requestSMSAPI(data)` → `POST user/sms-send`
  - `verifyCodeAPI(data)` → `POST user/sms-enable`
  - `getUserProfileAPI()` → `GET user/get-user-profile`
  - `updateUserProfileAPI(data)` → `POST user/set-user-profile`
  - `get2faQrcodeAPIAPI()` → `GET user/google-2fa-barcode`
  - `deleteUserImageAPI()` → `POST user-profile-image/delete`

### funds_services.ts — Balance & Payments
- **Methods + Endpoints:**
  - `getBalancesAPI()` → `GET user-balance/balance` (sorted desc)
  - `getDepositAndWithdrawAPI(code)` → `GET user-balance/withdraw-deposit`
  - `getTransactionHistoryAPI(params)` → `GET crypto-payment`
  - `getOrderDetailAPI(id)` → `GET crypto-payment/detail`
  - `getFormerWithdrawAddressesAPI(code)` → `GET withdraw-address/former-addresses`
  - `withdrawAPI(data)` → `POST crypto-payment/withdraw`
  - `preWithdrawAPI(data)` → `POST crypto-payment/pre-withdraw`

### orders_service.ts — Orders & Trades
- **Methods + Endpoints:**
  - `getOpenOrdersAPI()` → `GET order/open-orders`
  - `getOrderHistoryAPI()` → `GET order/full-history`
  - `getPaginatedOrderHistoryAPI(page, size)` → `GET order/full-history` (with pagination)
  - `getPaginatedTradeHistoryAPI(page, size)` → `GET trade/full-history` (with pagination)
  - `getFilteredOrderHistoryAPI(filters)` → `GET order/full-history` (with query params)
  - `getTradeHistoryAPI()` → `GET trade/full-history`
  - `getOrderHistoryDetailAPI(id)` → `GET order/detail`
  - `getCurrencyPairDetailsAPI(pair)` → `GET user-balance/pair-balance`
  - `createNewOrderAPI(data)` → `POST order/create`
  - `createNewStopOrderAPI(data)` → `POST order/stop-order-create`
  - `cancelOrderAPI(id)` → `POST order/cancel`

### pairs_service.ts — Currency Pairs
- **Methods + Endpoints:**
  - `addRemoveFavPairAPI(data)` → `POST currencies/favorite`
  - `getFavPairAPI()` → `GET currencies/favorite-pairs`
  - `getPairsListAPI()` → `GET currencies/pairs-list`

### address_management_service.ts — Withdrawal Addresses
- **Methods + Endpoints:**
  - `getCurrenciesAPI()` → `GET currencies`
  - `getWithDrawAddressesAPI()` → `GET withdraw-address`
  - `addNewWithDrawAddressAPI(data)` → `POST withdraw-address/new`
  - `deleteWithDrawAddressAPI(ids)` → `POST withdraw-address/delete`
  - `setFavoriteWithDrawAddressAPI(data)` → `POST withdraw-address/favorite`

### marketTrade_services.ts — Market Trades
- **Methods + Endpoints:**
  - `getMarketTradesAPI(pair)` → `GET trade-book`

### trade_chart_service.ts — TradingView Chart
- **Methods + Endpoints:**
  - `getChartConfigAPI()` → `GET http://116.203.76.196/tv/api/v1/js/get-configuration` (external IP)

### upload_service.ts — File Uploads (Axios)
- **Methods:**
  - `UploadFile(data)` — Single file upload with progress via `onUploadProgress`
  - `UploadMultiFile(data)` — Front + back image upload
  - `mocUpload()` — Dev mock upload
- **Endpoints:** `POST user-profile-image/upload?need_id=true`, `POST user-profile-image/multiple-upload`
- **Features:** FormData multipart, Bearer token, 30s timeout, progress published via MessageService (SET_UPLOADER_STATE, UPLOAD_PERCENTAGE)

### MqttService2.ts — Real-time MQTT (Active)
- **Pattern:** Singleton via `MqttService.getInstance()`
- **Connection:** `wss://{mainUrl}:8443` with XOR cipher auth
- **Methods:**
  - `ConnectToSubject({subject})` — Subscribe to MQTT topic
  - `DisconnectFromSubject({subject})` — Unsubscribe
  - `ConnectToNewSubject({oldsubject, newSubject})` — Switch subscription
- **Topic Routing:** Messages routed by topic prefix to specific RxJS subjects:
  - `main/trade/ticker` → MarketWatchMessageService
  - `main/trade/order-book/` → OrderBookMessageService
  - `main/trade/trade-book/` → MarketTradeMessageService
  - `main/trade/kline/` → TradeChartMessageService
  - (other) → SideMessageService

### RegisteredMqttService.ts — Authenticated MQTT
- **Pattern:** Singleton with optional re-initialization on token refresh
- **Connection:** Same WSS endpoint, token-based auth
- **Methods:** Same as MqttService2 (ConnectToSubject, DisconnectFromSubject, ConnectToNewSubject)
- **Features:** Health check (publishes `/testing` every 5 seconds), token refresh support, all messages → RegisteredUserMessageService

### mqttService.ts — Legacy MQTT (Deprecated)
- **Pattern:** Hook-style (`useStartMQTTMessages()`)
- **Status:** Deprecated — use MqttService2 instead

### message_service.ts — RxJS Pub-Sub Bus
- **Pattern:** Multiple RxJS Subject/ReplaySubject instances
- **10 Message Channels:**
  1. `MessageService` / `Subscriber` — Main app-wide event bus
  2. `DataInjectMessageService` / `dataInjectSubscriber` — Data injection
  3. `ReplayMessageService3` / `RepaySubscriber3` — ReplaySubject(3) buffer
  4. `SideMessageService` / `SideSubscriber` — MQTT side messages
  5. `EventMessageService` / `EventSubscriber` — UI events
  6. `RegisteredUserMessageService` / `RegisteredUserSubscriber` — Auth user MQTT messages
  7. `MarketTradeMessageService` / `MarketTradeSubscriber` — Trade data stream
  8. `OrderBookMessageService` / `OrderBookSubscriber` — Order book stream
  9. `MarketWatchMessageService` / `MarketWatchSubscriber` — Market watch stream
  10. `TradeChartMessageService` / `TradeChartSubscriber` — Chart data stream
- **Usage:** `MessageService.send({ name: MessageNames.XXX, value?, payload?, id? })`
- **Subscribe:** `Subscriber.subscribe((msg) => { if (msg.name === MessageNames.XXX) ... })`

### cookie.ts — Cookie Management
- **Config:** 28-day expiry, SameSite: Strict, Secure (prod), domain: `.unitedbit.com` (prod) / `localhost` (dev)
- **Cookie Keys:** `ubt` (token), `rt` (refresh token), `ube` (email), `fl` (from landing)

### constants.ts — Configuration
- **Enums:** RequestTypes, LocalStorageKeys (22 keys), SessionStorageKeys, UploadUrls
- **URLs:** `BaseUrl = https://{mainUrl}/api/v1/`, `mqttServer = wss://{mainUrl}:8443`
- **Environment:** `mainUrl` = `dev-app.unitedbit.com` (dev) or `app.unitedbit.com` (prod)
- **MQTT Config:** protocol: wss, connectTimeout: 30min, reconnectPeriod: 2s, keepalive: 0

### toastService.ts — Error Toast Formatter
- **Methods:** `ToastMessages(errors)` — Iterates error object, cleans field names, displays individual toast notifications

---

## g) State Management — Full Architecture

### Redux Store
```
configureStore.ts
├── sagaMiddleware (redux-saga)
├── routerMiddleware (redux-first-history)
├── composeWithDevTools (dev only)
└── Store extensions:
    ├── store.runSaga         — sagaMiddleware.run reference
    ├── store.injectedReducers — Dynamic reducer registry {}
    └── store.injectedSagas   — Dynamic saga registry {}
```

### Root Reducer (reducers.ts)
```
combineReducers({
  global: globalReducer,          // App container state (always loaded)
  language: languageProviderReducer, // Locale state (always loaded)
  router: routerReducer,          // Route state (always loaded)
  ...injectedReducers             // Dynamically injected per container
})
```
**Special:** Dispatching `App/LOGGED_IN_ACTION` resets entire state to `undefined` (fresh start)

### ApplicationRootState (types/index.d.ts)
```typescript
interface ApplicationRootState {
  router: RouterState;
  global: AppState;
  language: LanguageProviderState;
  home: HomeState;
  changePassword: ChangePasswordState;
  phoneVerificationPage: PhoneVerificationPageState;
  AddressManagementPage: AddressManagementPageState;
  loginPage: LoginState;
  signupPage: SignupPageState;
  acountPage: AcountPageState;
  ordersPage: OrdersPageState;
  tradeChart: TradeChartState;
  tradePage: TradePageState;
  tradeHeader: TradeHeaderState;
  fundsPage: FundsPageState;
  changeUserInfoPage: ChangeUserInfoPageState;
  recapchaContainer: RecapchaContainerState;
  documentVerificationPage: DocumentVerificationPageState;
  emailAuthentication: EmailAuthenticationState;
  updatePasswordPage: UpdatePasswordState;
  googleAuthenticationPage: GoogleAuthenticationPageState;
  contactUsPage: ContactUsPageState;
}
```

### All 90+ MessageService Message Types
```
// Loading states
SETLOADING, SET_POPUP_LOADING, SETGRIDLOADING, SETWITHDRAWLOADING,
SET_PAGE_LOADING, SET_PAGE_LOADING_WITH_ID, IS_LOADING_BUY_SELL

// Authentication
LOGGED_IN, LOGGED_OUT, AUTH_ERROR_EVENT, OPEN_LOGIN_POPUP

// Modals & popups
CLOSE_MODAL, OPEN_G2FA, OPEN_WITHDRAW_VERIFICATION_POPUP,
OPEN_TWOFA_AND_EMAILCODE_POPUP, OPEN_ALERT

// reCAPTCHA
SET_RECAPTCHA, RESET_RECAPTCHA, RESET_SITE_KEY, SET_SITE_KEY

// Grid data management
SET_GRID_DATA, SET_GRID_FILTER, RESET_GRID_FILTER, ADD_DATA_ROW_TO_GRID,
SET_PAGE_FILTERS_WITH_ID, DELETE_GRID_ROW

// Orders
SET_OPEN_ORDERS_DATA, NEW_ORDER_NOTIFICATION, HIDE_ORDER_NOTIFICATION,
IS_CANCELING_ORDER, SET_ORDER_HISTORY_DATA, SET_PAGINATED_ORDER_HISTORY_DATA,
SET_PAGINATED_TRADE_HISTORY_DATA, SET_TRADE_HISTORY_DATA, ADD_ONE_ORDER_TO_HISTORY

// Balance / Funds
SET_BALANCE_PAGE_DATA, SET_DEPOSIT_PAGE_DATA, SET_FORMER_WITHDRAW_ADDRESSES,
SET_ORDER_DETAIL, ADD_DATA_ROW_TO_WITHDRAWS

// Infinite scroll
SET_INFINITE_DW_PAGE_DATA, SET_INITIAL_INFINITE_DW_PAGE_DATA,
SET_INITIAL_INFINITE_DW_PAGE_DATA_LOADING, SET_DATA_TO_INFINITE_BOTTOM,
RESET_INFINITE_SCROLL

// File upload
SET_UPLOADER_STATE, UPLOAD_PERCENTAGE, SET_UPLOADED_IMAGE, DELETE_UPLOADED_IMAGE,
SET_IMAGE_PREVIEW, PROFILE_FILE_LOADED, UNLOCK_SUBTYPE_SELECT, RESET_IMAGES,
TOGGLE_SEND_IMAGE_BUTTON

// Address management
SET_FAVIORITE_ADDRESS

// Navigation / UI
SET_STEP, SET_TAB, RESIZE, CHANGE_THEME, CHANGE_LAYOUT, LAYOUT_RESIZE,
LAYOUT_CHANGE, ADDITIONAL_ACTION

// Trade page
UNSUBSCRIBE_FORM_STREAM, SET_TRADE_PAGE_CURRENCY_PAIR, SELECT_ORDERBOOK_ROW,
SUBSCRIBE_TO_STREAM, SET_CURRENCY_PAIR_DETAILS, GET_CURRENCY_PAIR_DETAILS,
MAIN_CHART_SUMMARY, MAIN_CHART_LAST_PRICE

// Validation / errors
SET_ERROR, SET_INPUT_ERROR

// Document upload
SET_DOCUMENT_IMAGES, REFRESH_VISIBLE_SECTION

// MQTT
RECONNECT_EVENT

// Test
SET_LOADING_TEST, SET_LOADING_END

// DataInject messages
MARKET_TRADES_INITIAL_DATA

// Event messages
OPEN_TOOLTIP, REFRESH_ORDER_GRID, GOT_FAV_PAIRS
```

### Data Flow Diagram
```
User Action
    │
    ▼
Container dispatches Redux Action
    │
    ├──► Reducer updates state synchronously
    │       │
    │       ▼
    │    Selector (reselect) recomputes
    │       │
    │       ▼
    │    Component re-renders
    │
    └──► Saga catches action (takeLatest)
            │
            ├──► Calls service function (e.g., ordersService.getOpenOrdersAPI())
            │       │
            │       ▼
            │    ApiService.getInstance().fetchData() → Fetch API → Server
            │       │
            │       ▼
            │    Response → Saga puts success/failure action
            │
            └──► Saga publishes to MessageService
                    │
                    ▼
                 RxJS Subscriber.next(message)
                    │
                    ▼
                 Any subscribed component receives via .subscribe()

MQTT Real-time Flow:
    MqttService2 singleton ←──WSS──► MQTT Broker (port 8443)
        │
        ▼ (on message)
    Topic routing by prefix:
        main/trade/ticker     → MarketWatchMessageService
        main/trade/order-book → OrderBookMessageService
        main/trade/trade-book → MarketTradeMessageService
        main/trade/kline      → TradeChartMessageService
        (other)               → SideMessageService
        │
        ▼
    Subscribed components (e.g., StreamComponentsWrapper) update UI
```

---

## h) Real-time — MQTT

### Connection Setup
- **Server:** `wss://{mainUrl}:8443` (WebSocket Secure)
- **Protocol:** MQTT over WSS
- **Config:** connectTimeout: 30min, reconnectPeriod: 2s, keepalive: 0
- **Auth:** XOR cipher (`mqttCipher('ubSalt')`) encrypts random client ID
  - `username`: JWT token from cookie (or encrypted client ID if not logged in)
  - `password`: encrypted client ID
  - `clientId`: encrypted client ID

### Two MQTT Services
1. **MqttService2** (public) — Singleton, routes messages to 5 RxJS subjects by topic prefix
2. **RegisteredMqttService** (authenticated) — Singleton, optional token refresh, health check every 5s, all messages → RegisteredUserMessageService

### Topic Prefixes (from `MqttTopicsPrefixes` enum)
| Prefix | Example | Target RxJS Subject |
|--------|---------|-------------------|
| `main/trade/ticker` | `main/trade/ticker` | MarketWatchSubscriber |
| `main/trade/order-book/` | `main/trade/order-book/BTC_USDT` | OrderBookSubscriber |
| `main/trade/trade-book/` | `main/trade/trade-book/BTC_USDT` | MarketTradeSubscriber |
| `main/trade/kline/` | `main/trade/kline/BTC_USDT` | TradeChartSubscriber |

### Adding MQTT Subscriptions
```typescript
// 1. Get service instance
import { mqttService2 } from 'services/MqttService2';

// 2. Subscribe to a topic
mqttService2.ConnectToSubject({ subject: 'main/trade/order-book/BTC_USDT' });

// 3. Listen for messages via the appropriate RxJS subject
import { OrderBookSubscriber } from 'services/message_service';
OrderBookSubscriber.subscribe((msg) => {
  console.log(msg.payload); // Parsed JSON from MQTT message
});

// 4. Switch subscription (e.g., user changes currency pair)
mqttService2.ConnectToNewSubject({
  oldsubject: 'main/trade/order-book/BTC_USDT',
  newSubject: 'main/trade/order-book/ETH_USDT',
});

// 5. Unsubscribe
mqttService2.DisconnectFromSubject({ subject: 'main/trade/order-book/ETH_USDT' });
```

---

## i) Authentication

### JWT Flow
1. User submits login form → `POST /api/v1/auth/login` with email + password + reCAPTCHA token
2. Server returns JWT access token + refresh token
3. Client stores in cookies: `ubt` (access token), `rt` (refresh token), `ube` (email)
4. All subsequent API calls include `Authorization: Bearer {token}` header
5. On 401/403 response → `AUTH_ERROR_EVENT` published → user redirected to `/login`

### Cookie Configuration
| Key | Cookie Name | Content |
|-----|-------------|---------|
| Token | `ubt` | JWT access token |
| RefreshToken | `rt` | JWT refresh token |
| Email | `ube` | User email |
| FromLanding | `fl` | Landing page referrer |

**Settings:** 28-day expiry, SameSite: Strict, Secure (production), domain: `.unitedbit.com` (prod) / `localhost` (dev)

### reCAPTCHA
- **v3:** GoogleReCaptchaProvider wraps entire app (RecapchaContainer)
- **v2:** react-google-recaptcha for visible challenges on specific forms
- Site key fetched from `GET main-data/common`, stored in sessionStorage

### 2FA (Google Authenticator)
- QR code fetched from `GET user/google-2fa-barcode`
- Enable: `POST user/google-2fa-enable` with TOTP code
- Disable: `POST user/google-2fa-disable` with TOTP code
- PIN input: 6-digit via `pinInput` component

---

## j) Internationalization (i18n)

### Setup
- **Library:** react-intl v2.9.0 (uses deprecated `addLocaleData` API)
- **Config:** `app/i18n.ts` — CommonJS syntax required by `extract-intl` script
- **Provider:** `LanguageProvider` container wraps app with `IntlProvider`
- **Default locale:** `en` (English)

### Supported Locales
| Code | Language | Direction | File |
|------|----------|-----------|------|
| `en` | English | LTR | `app/translations/en.json` |
| `de` | German | LTR | `app/translations/de.json` |
| `ar` | Arabic | RTL | `app/translations/ar.json` |

### RTL Support
- Arabic locale triggers RTL via `jss-rtl` plugin
- `document.body.classList.add('arabic')` applied when locale is `ar`

### Message Key Convention
- Global titles: `app.globalTitles.{key}` (e.g., `app.globalTitles.userId`)
- Container messages: `containers.{ContainerName}.{key}`
- Component messages: `components.{ComponentName}.{key}`
- Grid headers: `GridHeaderNames.{key}`

### Adding Translations
```typescript
// 1. Define messages in container messages.ts:
import { defineMessages } from 'react-intl';
export default defineMessages({
  myTitle: { id: 'containers.MyPage.myTitle', defaultMessage: 'My Title' },
});

// 2. Add to all 3 translation JSON files:
// en.json: { "containers.MyPage.myTitle": "My Title" }
// de.json: { "containers.MyPage.myTitle": "Mein Titel" }
// ar.json: { "containers.MyPage.myTitle": "عنواني" }

// 3. Use in component:
<FormattedMessage {...messages.myTitle} />

// 4. Or extract automatically:
yarn extract-intl
```

---

## k) Build System

### npm Scripts (all via `yarn`)
| Command | Purpose | Environment |
|---------|---------|-------------|
| `yarn start` | Dev server (localhost, CRA-ejected) | `NODE_ENV=development IS_LOCAL=true` |
| `yarn start-dev` | Dev server (remote APIs) | `NODE_ENV=development` |
| `yarn build` | Production build (no sourcemaps) | `NODE_ENV=production GENERATE_SOURCEMAP=false` |
| `yarn build-dev` | Dev API production build | `NODE_ENV=production IS_DEV_BUILD=true` |
| `yarn start:prod` | Express production server | `NODE_ENV=production` |
| `yarn start:tunnel` | Dev server with ngrok tunnel | `NODE_ENV=development ENABLE_TUNNEL=true` |
| `yarn test` | Jest with coverage | `NODE_ENV=test` |
| `yarn test:watch` | Jest watch mode | `NODE_ENV=test` |
| `yarn lint` | ESLint + stylelint + TSLint | — |
| `yarn lint:ts` | ESLint for TS/TSX only | — |
| `yarn typecheck` | `tsc --noEmit` | — |
| `yarn report` | webpack-bundle-analyzer (port 4200) | — |
| `yarn analyze` | Full bundle analysis with stats.json | — |
| `yarn generate` | Plop code generators (component/container) | — |
| `yarn extract-intl` | Extract i18n messages from source | — |
| `yarn prettify` | Prettier formatting | — |

### Webpack Configuration
- **CRA path:** `craConfig/config/webpack.config.js` — Used by `yarn start` and `yarn build`
- **Internals path:** `internals/webpack/` — Legacy configs (webpack.base.babel.js, .dev, .prod)
- **Loaders:** babel-loader, ts-loader, css-loader, sass-loader, svg-url-loader, file-loader
- **Plugins:** HtmlWebpackPlugin, MiniCssExtractPlugin, TerserWebpackPlugin, CompressionPlugin, WorkboxPlugin
- **Code splitting:** Dynamic imports via React.lazy + Loadable pattern
- **Path aliases:** `baseUrl: ./app` in tsconfig.json (e.g., `import X from 'services/api_service'`)

### Docker
```bash
# Development (dev APIs)
docker build -f Dockerfile .
# → node:18-slim, NODE_OPTIONS=--openssl-legacy-provider, yarn build-dev

# Production
docker build -f DockerfileProd .
# → node:18-slim, NODE_OPTIONS=--openssl-legacy-provider, yarn build
```
**Note:** `--openssl-legacy-provider` required until Webpack 5 migration

### Production Server (server/index.js)
- Express static file serving from `/build`
- Automatic gzip: serves `.js.gz` files with `Content-Encoding: gzip`
- SPA routing: all non-file requests serve `index.html`
- Optional ngrok tunnel for remote access
- Default port: 3000

---

## l) Testing

### Jest Configuration (jest.config.js)
- **Preset:** `ts-jest/presets/js-with-babel`
- **Environment:** jsdom
- **Test pattern:** `tests/.*\.test\.(js|ts(x?))$`
- **Module directories:** `['node_modules', 'app']` (same path resolution as webpack)
- **Module mocks:** CSS/SCSS → cssModule.js, images → image.js, toastify → toastify.js
- **Setup:** `@testing-library/jest-dom` + `internals/testing/test-bundler.js`
- **Watch plugins:** jest-watch-typeahead (filename + testname search)

### Coverage Thresholds
| Metric | Threshold |
|--------|-----------|
| Statements | 98% |
| Branches | 91% |
| Functions | 98% |
| Lines | 98% |

### Test File Locations
Tests live alongside their source in `tests/` subdirectories:
- `app/tests/` — Root-level (i18n.test.ts, store.test.ts)
- `app/containers/*/tests/` — Container tests
- `app/components/*/tests/` — Component tests
- `app/utils/tests/` — Utility tests

### Running Tests
```bash
yarn test                    # Full suite with coverage
yarn test:watch             # Watch mode
yarn test -- --testPathPattern="LoginPage"  # Single container
```

### Current Status
- **30/30 test suites pass**, 94 tests, 3 skipped, 3 todo
- Only ~9 containers have tests (83% untested)
- Infrastructure tests (injectors, store, i18n) have good coverage

---

## m) API Integration

### Base URL
- **Production:** `https://app.unitedbit.com/api/v1/`
- **Development:** `https://dev-app.unitedbit.com/api/v1/`
- **TradingView:** `https://{mainUrl}/tv/api/v1/`
- Determined by `NODE_ENV` + `IS_DEV_BUILD` flag in `services/constants.ts`

### All REST Endpoints
| Method | Endpoint | Service | Purpose |
|--------|----------|---------|---------|
| POST | `auth/login` | security_service | User login |
| POST | `auth/register` | security_service | User registration |
| POST | `auth/verify` | security_service | Email verification |
| POST | `auth/forgot-password` | security_service | Request password reset |
| POST | `auth/forgot-password/update` | security_service | Complete password reset |
| GET | `user/user-data` | security_service | Fetch user account data |
| POST | `user/send-verification-email` | security_service | Resend verification email |
| POST | `user/change-password` | security_service | Change password |
| POST | `user/google-2fa-enable` | security_service | Enable 2FA |
| POST | `user/google-2fa-disable` | security_service | Disable 2FA |
| GET | `user/get-user-profile` | user_acount_service | Get profile data |
| POST | `user/set-user-profile` | user_acount_service | Update profile |
| GET | `user/google-2fa-barcode` | user_acount_service | Get 2FA QR code |
| POST | `user/sms-send` | user_acount_service | Send SMS verification |
| POST | `user/sms-enable` | user_acount_service | Verify SMS code |
| GET | `main-data/common` | security_service | reCAPTCHA site key |
| GET | `main-data/country-list` | user_acount_service | Country list |
| GET | `user-balance/balance` | funds_services | User balance |
| GET | `user-balance/withdraw-deposit` | funds_services | Deposit/withdraw info |
| GET | `user-balance/pair-balance` | orders_service | Pair-specific balance |
| GET | `crypto-payment` | funds_services | Transaction history |
| GET | `crypto-payment/detail` | funds_services | Transaction detail |
| POST | `crypto-payment/withdraw` | funds_services | Execute withdrawal |
| POST | `crypto-payment/pre-withdraw` | funds_services | Preview withdrawal |
| GET | `order/open-orders` | orders_service | Open orders |
| GET | `order/full-history` | orders_service | Order history |
| GET | `order/detail` | orders_service | Order detail |
| POST | `order/create` | orders_service | Create market order |
| POST | `order/stop-order-create` | orders_service | Create stop order |
| POST | `order/cancel` | orders_service | Cancel order |
| GET | `trade/full-history` | orders_service | Trade history |
| GET | `trade-book` | marketTrade_services | Market trade book |
| GET | `currencies` | address_management | Currency list |
| GET | `currencies/pairs-list` | pairs_service | All currency pairs |
| GET | `currencies/favorite-pairs` | pairs_service | User's favorite pairs |
| POST | `currencies/favorite` | pairs_service | Toggle favorite pair |
| GET | `withdraw-address` | address_management | Withdrawal addresses |
| POST | `withdraw-address/new` | address_management | Add address |
| POST | `withdraw-address/delete` | address_management | Delete address |
| POST | `withdraw-address/favorite` | address_management | Favorite address |
| GET | `withdraw-address/former-addresses` | funds_services | Former addresses |
| POST | `user-profile-image/upload` | upload_service | Upload single image |
| POST | `user-profile-image/multiple-upload` | upload_service | Upload front+back |
| POST | `user-profile-image/delete` | user_acount_service | Delete profile image |

### Error Handling
| HTTP Status | Behavior |
|-------------|----------|
| 200 | Success — `response.data` returned |
| 401/403 | `AUTH_ERROR_EVENT` → MessageService → redirect to `/login` |
| 422 | Validation errors → `ToastMessages(errors)` → individual field toasts |
| 500 | Generic error toast notification |

### Standard Response Shape
```typescript
interface StandardResponse {
  status: boolean;
  message: string;
  data: any;
}
```

---

## n) Conventions

### Code Style
- **Linting:** ESLint with `@typescript-eslint/recommended`, `plugin:react/recommended`, `plugin:react-hooks/recommended`, `prettier`
- **Legacy:** tslint.json still present (deprecated, not actively used in CI)
- **Formatting:** Prettier — single quotes, semicolons, trailing commas, 80 char width, 2-space indent
- **TSX:** Double quotes in JSX (`avoidEscape: true`)

### Naming Conventions
| Type | Convention | Example |
|------|-----------|---------|
| Containers | PascalCase directory | `LoginPage/`, `TradePage/` |
| Components | camelCase or PascalCase directory | `dataRow/`, `AppHeader/` |
| Services | snake_case file | `api_service.ts`, `funds_services.ts` |
| Actions | camelCase with Action suffix | `getUserDataAction()` |
| Action types | `'Container/SCREAMING_SNAKE'` | `'App/LOGGED_IN_ACTION'` |
| Selectors | `makeSelectXxx` prefix | `makeSelectLoggedIn()` |
| Reducers | Default export from `reducer.ts` | `export default reducer` |
| Types | `ContainerState` interface | `export interface ContainerState { ... }` |
| Message names | SCREAMING_SNAKE enum | `MessageNames.SET_LOADING` |
| CSS variables | `--camelCase` | `--primary`, `--cardBorderRadius` |

### File Organization
- One container = one directory with all 8+ required files
- Components are self-contained directories (index.tsx + optional tests/, styles)
- Services are flat files in `app/services/`
- Shared types in `app/types/index.d.ts`
- Path aliases via tsconfig `baseUrl: ./app` (import from `'services/...'`, `'components/...'`, etc.)

### API Integration Pattern
```typescript
// In service file (e.g., services/my_service.ts):
import { apiService } from './api_service';
import { RequestTypes } from './constants';

export const getDataAPI = async () => {
  return apiService.fetchData({
    requestType: RequestTypes.GET,
    url: 'my-endpoint',
    data: null,
  });
};

// In saga (e.g., containers/MyPage/saga.ts):
import { call, put, takeLatest } from 'redux-saga/effects';
import { getDataAPI } from 'services/my_service';
import ActionTypes from './constants';

function* getData() {
  try {
    const response = yield call(getDataAPI);
    yield put({ type: ActionTypes.SET_DATA, payload: response.data });
  } catch (err) {
    // Error handling
  }
}

export default function* mySaga() {
  yield takeLatest(ActionTypes.GET_DATA, getData);
}
```

---

## o) Known Issues & Technical Debt

### Critical
- **Hardcoded URLs** — `services/constants.ts` has hardcoded prod/dev URLs; should use `.env`
- **MQTT XOR cipher** — `mqttCipher('ubSalt')` is trivially reversible; should use standard TLS + JWT
- **axios 0.21.0** — CVE-2023-45857 (CSRF); requires upgrade to ^1.x
- **react-intl 2.9.0** — Deprecated; `addLocaleData` API removed in v3+

### High Priority
- **Material-UI 4 → MUI 5** — Major breaking import/theme changes
- **Webpack 4 → 5** — Needed to remove `--openssl-legacy-provider` flag
- **redux-form → react-hook-form** — redux-form is abandoned
- **Redux → Redux Toolkit** — Modernize action/reducer boilerplate
- **@types/react pinned to 17.0.2** — Runtime is React 18 but types don't match

### Medium Priority
- Only ~9 test files for 21 containers (83% untested)
- Mixed HTTP clients (Fetch API + Axios) — standardize on one
- No Error Boundaries in React tree
- AG Grid Enterprise 25 → 31+ (major version gap)
- `trade_chart_service.ts` uses hardcoded external IP (`116.203.76.196`)
- TSLint config still present (tslint.json, tslint-imports.json) — remove

### Low Priority
- Typo: `AcountPage` (missing 'c' in Account) — pervasive, changing is high-risk
- amcharts 4 → 5
- lodash 4.17.20 → 4.17.21 (security patch)
- `unused.md` references developer's local D:\ paths — not useful for team

---

## p) Upgrade Roadmap

See `UPGRADE_PLAN.md` for the full 8-phase roadmap. Summary:

| Phase | What | Risk | Status |
|-------|------|------|--------|
| 1 | Security fixes (axios, lodash, .env, Docker) | Low | Partially done (Docker updated to node:18-slim) |
| 2 | TSLint → ESLint | Low | **Done** (.eslintrc.json created, ESLint configured) |
| 3 | TypeScript 4 → 5.4.5 | Medium | **Done** (tsconfig updated to ES2017) |
| 4 | React 17 → 18 | Medium | **Done** (createRoot, react-dom 18.3.1) |
| 5 | Material-UI 4 → MUI 5 | High | Not started |
| 6 | Webpack 4 → 5 | Medium | Not started |
| 7 | react-intl 2 → 6 | Medium | Not started |
| 8 | Modernization (redux-form, AG Grid, Error Boundaries) | Varies | Not started |

**Constraints:** Tests must pass after each phase. Baseline: 30/30 suites, 94 tests.

---

## q) Deep Audit Findings & Corrections

> Results of a line-by-line code audit cross-referencing all AGENTS.md claims against actual source.

### AGENTS.md Accuracy Corrections

| Section | Original Claim | Correction |
|---------|---------------|------------|
| §f API Service | "422 → ToastMessages(errors) → individual field toasts" | 422 handling is in individual saga `catch` blocks, not in `api_service.ts`. The API service has no 422-specific code. |
| §f API Service | 4 methods listed | Actually 5 methods — `handleRawResponse()` is implicitly public (no `private` keyword) and was omitted |
| §f mqttService.ts | "Legacy MQTT (Deprecated) — Hook-style, uses Paho" | Uses **mqtt.js** (same lib as MqttService2), not Eclipse Paho. Correct that it's hook-style and deprecated. |
| §g MessageNames | "All 90+ MessageService Message Types" | Exactly **82 total** enum values: 78 `MessageNames` + 3 `EventMessageNames` + 1 `DataInjectMessageNames` |
| §h MQTT Topics | 4 topic prefixes listed | Actually **7 distinct topic patterns** (see below) |
| §i Auth | Token refresh described as active | Token refresh is **DISABLED** — `retryWithNewToken()` call is commented out in `api_service.ts:68-70`; refresh token storage is commented out in `LoginPage/saga.ts:52-57` |
| §i Auth | Logout flow | Missing critical detail: **NO server-side logout API call** — logout is entirely client-side (clear cookies + localStorage) |
| §k Build | "Internals path: Legacy configs" | Should explicitly state **DEAD CODE** — only `craConfig/` configs are used; `internals/webpack/` files are unreachable |

### Additional MQTT Topics (Missing from §h)

| # | Topic Pattern | Service | Auth | Purpose |
|---|---------------|---------|------|---------|
| 5 | `main/trade/user/{channel}/open-orders/` | RegisteredMqttService | Yes | Real-time user order updates |
| 6 | `main/trade/user/{channel}/crypto-payments/` | RegisteredMqttService | Yes | Real-time payment notifications |
| 7 | `/testing` (publish only) | RegisteredMqttService | Yes | Health check every 5 seconds |

### Dead Message Types (7 enum values never referenced)

| Enum | Safe to Delete |
|------|---------------|
| `SETGRIDLOADING` | ✅ |
| `SETWITHDRAWLOADING` | ✅ |
| `SET_INFINITE_DW_PAGE_DATA` | ✅ |
| `SET_IMAGE_PREVIEW` | ✅ |
| `UNSUBSCRIBE_FORM_STREAM` | ✅ |
| `SUBSCRIBE_TO_STREAM` | ✅ |
| `REFRESH_VISIBLE_SECTION` | ✅ |

Additionally, **16 message types** are subscribed to but never published (orphaned listeners), and **3 message types** are published but never subscribed to (`DELETE_UPLOADED_IMAGE`, `GET_CURRENCY_PAIR_DETAILS`, `SET_ERROR`).

### Dead Code Files

| File | Evidence | Action |
|------|----------|--------|
| `app/services/mqttService.ts` | `useStartMQTTMessages` exported but **never imported** anywhere | Delete |
| `internals/webpack/webpack.base.babel.js` | Only referenced by `yarn old_build` (never called in CI/CD) | Delete |
| `internals/webpack/webpack.dev.babel.js` | Same — legacy, unreachable | Delete |
| `internals/webpack/webpack.prod.babel.js` | Same — legacy, unreachable | Delete |
| `unused.md` | References developer's local `D:\mohsen\` paths | Delete or replace |
| `tslint.json` + `tslint-imports.json` | TSLint deprecated; ESLint is the active linter | Delete |

### Security Findings

| Severity | Issue | Location |
|----------|-------|----------|
| 🔴 Critical | **Hardcoded IP address** `http://116.203.76.196/tv/api/v1/js/get-configuration` (HTTP, not HTTPS) | `services/trade_chart_service.ts:6` |
| 🔴 Critical | **No HttpOnly flag** on JWT cookies — XSS can steal tokens via `document.cookie` | `services/cookie.ts` (js-cookie cannot set HttpOnly; requires server-side Set-Cookie) |
| 🟠 High | **Token refresh disabled** — `retryWithNewToken()` commented out; users must re-login on expiry | `services/api_service.ts:68-70`, `LoginPage/saga.ts:52-57` |
| 🟠 High | **71 console.log statements** in production code — `api_service.ts` logs request params | App-wide |
| 🟡 Medium | **reCAPTCHA token in localStorage** (`'recall'` key) — should use sessionStorage | `LoginPage/loginBody.tsx`, `RecapchaContainer/` |
| 🟡 Medium | **No security headers** in Express server (no helmet/CSP/HSTS/X-Frame-Options) | `server/index.js` |
| 🟢 Low | **Misleading key name** `LocalStorageKeys.SITEKEY = 'refreshToken'` | `services/constants.ts:21` |

### Undocumented Build Features

| Feature | Location | Notes |
|---------|----------|-------|
| **Brotli compression** | `craConfig/config/webpack.config.js` | `.js`, `.css`, `.html`, `.svg` files; threshold 10KB |
| **PWA manifest** | WebpackPwaManifest plugin | App name: "Client Cabinet", theme: #396DE0 |
| **Service Worker** | WorkboxWebpackPlugin.GenerateSW | SPA navigateFallback to `index.html` |
| **Bundle analyzer** | BundleAnalyzerPlugin | Generates `stats.json` (mode: disabled in prod) |

### Express Server Gaps (Not in AGENTS.md)

- **No security headers** — no helmet, CSP, HSTS, X-Frame-Options, X-Content-Type-Options
- **No request logging** — no morgan or equivalent
- **No global error handler** — `app.use((err, req, res, next))` missing
- **Gzip hack** — `app.get('*.js')` rewrites URL to `.js.gz` instead of using content negotiation
- **No CORS middleware** — code exists but is commented out
