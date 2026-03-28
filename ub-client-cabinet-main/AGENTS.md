# Client Cabinet — ub-client-cabinet-main

## Stack
- **React 17.0.1** / **TypeScript 5.4.5** / **Redux 4.0.5** + **Redux-Saga 1.1.3**
- **Material-UI 4.11.0** / **styled-components 5.2.0** / **AG Grid Enterprise 25.3.0**
- **redux-form 8.3.6** / **react-intl 2.9.0** (i18n) / **amcharts4** / TradingView charting_library
- **Custom Webpack 4.44.1 build** — two configs: `craConfig/` (CRA-ejected scripts) + `internals/webpack/`
- **Express 4.17.1** server for production SPA serving
- **Node >=12** (Docker: node:12.18.3, local: nvm lts/dubnium)

## Project Structure
```
app/
├── app.tsx              # Entry point — ReactDOM.render, redux Provider, ConnectedRouter
├── configureStore.ts    # Redux store with dynamic reducer/saga injection
├── reducers.ts          # Root combineReducers + dynamic injection merge
├── i18n.ts              # react-intl setup (en, de, ar + RTL)
├── bundle.js            # Webpack bundle entry
├── containers/          # 21 page containers (see Container Pattern below)
├── components/          # 41 reusable UI component folders
├── services/            # 17 singleton API/business service files
├── hooks/               # Custom React hooks (onlineStatusHook, useImmer)
├── types/               # TypeScript type definitions (index.d.ts)
├── utils/               # Utilities, injectors, formatters, validators
├── translations/        # i18n JSON translation files (en.json, de.json, ar.json)
├── styles/              # Global theme/style definitions
├── images/              # Static images and icons
├── tests/               # Root-level tests (i18n.test.ts, store.test.ts)
├── fonts/               # Custom fonts
├── datafeeds/           # TradingView datafeed implementation
└── charting_library/    # TradingView charting lib (vendored)

craConfig/               # CRA-ejected build scripts (start.js, build.js, webpack.config.js)
internals/
├── webpack/             # Webpack configs (base, dev, prod)
├── generators/          # plop code generators
├── mocks/               # Jest mocks (css, images)
├── scripts/             # Build helper scripts
├── templates/           # Code generation templates
└── testing/             # Test setup (test-bundler.js)

server/                  # Express production server (index.js)
public/                  # HTML template + vendor scripts
maintenancePage/         # Static maintenance page
```

## Architecture

### Container Pattern
Every page follows: `app/containers/<PageName>/`
- `index.tsx` — Main component with `useInjectReducer()` and `useInjectSaga()`
- `reducer.ts` — Redux reducer
- `saga.ts` — Redux-Saga side effects (API calls)
- `actions.ts` — Action creators
- `constants.ts` — Action type constants
- `selectors.ts` — Memoized reselect selectors
- `types.d.ts` — TypeScript interfaces
- `Loadable.tsx` — React.lazy code-splitting wrapper
- Optional: `components/`, `steps/`, `pages/` subfolders

### 21 Containers
AcountPage, AddressManagementPage, App, ChangePassword, ChangeUserInfoPage,
ContactUsPage, DocumentVerificationPage, EmailVerification, FundsPage,
GoogleAuthenticationPage, HomePage, LanguageProvider, LocaleToggle, LoginPage,
NotFoundPage, OrdersPage, PhoneVerificationPage, RecapchaContainer, SignupPage,
TradePage, UpdatePasswordPage

### Services (app/services/)
| Service | Purpose | HTTP Client |
|---------|---------|-------------|
| `api_service.ts` | Main REST API facade (singleton) | Fetch API |
| `upload_service.ts` | File uploads with progress | Axios |
| `security_service.ts` | Auth: login, register, 2FA, password reset | apiService |
| `user_acount_service.ts` | User profile, countries, 2FA QR | apiService |
| `funds_services.ts` | Balances, deposits, withdrawals | apiService |
| `orders_service.ts` | Order CRUD, trade history | apiService |
| `pairs_service.ts` | Currency pair data | apiService |
| `trade_chart_service.ts` | Chart/OHLC data | apiService |
| `marketTrade_services.ts` | Market trade data | apiService |
| `address_management_service.ts` | Withdrawal addresses | apiService |
| `MqttService2.ts` | Real-time MQTT (WSS) via mqtt.js | mqtt.js |
| `mqttService.ts` | Legacy Paho MQTT (deprecated) | paho-mqtt |
| `RegisteredMqttService.ts` | Additional MQTT service | mqtt.js |
| `message_service.ts` | RxJS pub-sub message bus (60+ event types) | RxJS |
| `cookie.ts` | JWT token storage in cookies | js-cookie |
| `constants.ts` | API routes, MQTT config, storage keys | — |
| `toastService.ts` | Toast notification formatter | — |

### State Management
- **Redux** with dynamic injection: each container registers its reducer/saga on mount
- **RxJS MessageService** for cross-component events (60+ message types)
- **connected-react-router** for route↔state sync
- **reselect** for memoized selectors

### Real-time
- **MQTT over WSS** (port 8443) with custom XOR cipher auth
- 4 independent RxJS subjects: MarketWatch, OrderBook, MarketTrade, TradeChart
- Paho MQTT service is deprecated, MqttService2 is active

### Auth
- JWT tokens stored in cookies (SameSite: Strict, Secure, 28-day expiry)
- Cookie keys: `ubt` (token), `rt` (refresh token), `ube` (email)
- Domain: `.unitedbit.com` (prod) / `localhost` (dev)
- Google reCAPTCHA v2 + v3 on auth endpoints
- Google Authenticator 2FA support

## Build and Test

```bash
# Install dependencies
yarn install

# Development server (localhost)
yarn start                    # IS_LOCAL=true, NODE_ENV=development

# Development server (remote dev)
yarn start-dev               # NODE_ENV=development

# Production build
yarn build                   # NODE_ENV=production, no sourcemaps

# Dev build (production mode but dev APIs)
yarn build-dev               # IS_DEV_BUILD=true

# Run tests with coverage
yarn test                    # Jest + ts-jest, 98% coverage threshold

# Lint
yarn lint                    # ESLint + stylelint + TSLint + typecheck

# Production serve
yarn start:prod              # Express server on port 3000

# Bundle analysis
yarn report                  # webpack-bundle-analyzer on port 4200
```

### Docker
```bash
# Development build
docker build -f Dockerfile .        # node:12.18.3 + yarn build-dev

# Production build
docker build -f DockerfileProd .    # node:12.18.3 + yarn build
```

## Conventions

### API Patterns
- All REST APIs under `https://{domain}/api/v1/`
- Auth header: `Authorization: Bearer {token}`
- Error handling: 401/403 → AUTH_ERROR_EVENT via MessageService → redirect to login
- 422 → validation errors; 500 → generic toast
- File uploads use Axios (multipart/form-data), everything else uses Fetch

### URLs/Environment
- URLs are hardcoded in `app/services/constants.ts` (no .env file)
- Production: `app.unitedbit.com`
- Dev: `dev-app.unitedbit.com`
- Mobile redirect: viewport < 1000px → `m.unitedbit.com`
- Environment detected via `NODE_ENV` + `IS_DEV_BUILD` + `IS_LOCAL` flags

### Code Style
- TSLint (deprecated) for TypeScript linting
- Prettier for code formatting
- Single quotes, semicolons, 2-space indent
- No unused variable enforcement (relaxed TSLint)

### i18n
- react-intl v2 with `<FormattedMessage>` components
- 3 locales: English (en), German (de), Arabic (ar)
- Arabic includes RTL support via `jss-rtl`

## Known Issues & Technical Debt

### Critical
- **Node 12 Docker image is outdated** — Dockerfile still references node:12.18.3, build works on Node 22
- **Hardcoded URLs** — should use .env / environment variables
- **MQTT XOR cipher** — weak auth, should use standard TLS + JWT
- **axios 0.21.4** — CVE-2021-3749 (ReDoS) fixed; CVE-2023-45857 (CSRF) requires 1.x upgrade
- **TSLint is deprecated** — migrate to ESLint + typescript-eslint
- **Pre-existing test failures** — 10/30 test suites fail due to TS type mismatch in `internals/templates/configureStore.ts` (fixed `initialState` type to `any`)

### High Priority Upgrades
- **React 17 → 18** (ReactDOM.render → createRoot)
- **Material-UI 4 → MUI 5** (breaking import changes)
- **react-intl 2 → 6** (major API changes)
- **Webpack 4 → 5** or migrate to Vite
- **redux-form → react-hook-form** (redux-form is abandoned)
- **Redux → Redux Toolkit** (modernize actions/reducers/slices)

### Medium Priority
- Only 9 test files for 21 containers — coverage threshold may be misleading
- Mixed HTTP clients (Fetch + Axios)
- No Error Boundaries in React tree
- AG Grid Enterprise 25 → 31+

### Low Priority
- Typo: "AcountPage" (missing 'c' in Account)
- amcharts4 → amcharts5
- lodash 4.17.20 → 4.17.21 (security patch)
