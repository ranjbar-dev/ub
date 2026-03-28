# Task: Create ARCHITECTURE.md

**ID:** p7-architecture-md  
**Phase:** 7 — Documentation for AI Agents  
**Severity:** 🔴 CRITICAL  
**Dependencies:** None  

## Problem

No architecture documentation exists. AI agents must read dozens of files to understand how the project is structured, how data flows, and what patterns are used.

## File to Create

**`src/ARCHITECTURE.md`** (or project root `ARCHITECTURE.md`)

### Content

```markdown
# ub-admin-main Architecture

## Overview
React 16 admin panel for the UnitedBit cryptocurrency exchange platform. Built with TypeScript, Redux-Saga, and Material-UI 4.

## Tech Stack
- **React** 16.13 with hooks
- **TypeScript** 3.9
- **Redux Toolkit** + **Redux-Saga** for state & side effects
- **Material-UI** 4 for base components
- **styled-components** for custom styling
- **AG Grid** (ag-grid-react) for data tables
- **RxJS** (Subjects) for MessageService pub/sub
- **i18next** for internationalization
- **React Router** 5 with connected-react-router

## Directory Structure
```
src/
├── app/
│   ├── components/         # Shared reusable components (SimpleGrid, sideNav, etc.)
│   ├── containers/         # Page-level containers (1 per route)
│   │   └── <Container>/
│   │       ├── index.tsx       # Main component
│   │       ├── saga.ts         # Redux-Saga side effects
│   │       ├── slice.ts        # Redux Toolkit slice (reducers + actions)
│   │       ├── selectors.ts    # Memoized Redux selectors
│   │       ├── types.ts        # TypeScript interfaces
│   │       ├── Loadable.tsx    # Lazy loading wrapper
│   │       └── components/     # Container-specific sub-components (optional)
│   ├── constants.ts        # AppPages enum (routes), WindowTypes
│   └── index.tsx           # App root, router, theme provider
├── services/
│   ├── api_service.ts      # Singleton HTTP client (fetch wrapper)
│   ├── constants.ts        # RequestTypes, StandardResponse, BaseUrl
│   ├── message_service.ts  # RxJS pub/sub (67 events)
│   └── *_service.ts        # Domain API functions (12 files)
├── store/
│   ├── configureStore.ts   # Redux store setup with saga middleware
│   └── slice.ts            # Global slice (loggedIn state)
├── types/
│   └── RootState.ts        # Root Redux state interface (24 container states)
├── utils/
│   ├── formatters.ts       # Currency, date, string formatters
│   ├── stylers.ts          # AG Grid cell style helpers
│   ├── loadable.tsx        # Lazy-load utility (React.lazy + Suspense)
│   └── fileDownload.js     # Browser file download helper
└── locales/                # i18next translation files
```

## Data Flow

### API Call Lifecycle
1. User interacts with UI → dispatches Redux action
2. Redux-Saga `takeLatest` catches the action
3. Saga calls service function (e.g., `GetUsersAPI(params)`)
4. Service function calls `apiService.fetchData()`
5. ApiService makes HTTP request with JWT auth header
6. Response is sent via `MessageService.send()` (⚠️ legacy pattern)
7. Container's `useEffect` subscribes to MessageService, receives data

### State Management (Dual Pattern)
The app has two parallel state systems:

**Redux Store** (intended but underutilized):
- Slice reducers exist but most are empty (just action triggers)
- Selectors exist but most containers don't use them
- Redux DevTools shows dispatched actions but not data

**MessageService** (legacy, actually holds the data):
- RxJS Subject-based pub/sub with 67 event types
- Sagas send data via `MessageService.send()`
- Containers subscribe in `useEffect` with `Subscriber.subscribe()`
- Data stored in local `useState`, invisible to Redux

### Authentication
- JWT token stored in `localStorage.token`
- ApiService reads token via `setHeaders()` on every request
- 401 responses trigger `MessageNames.AUTH_ERROR_EVENT`
- LoginPage container handles login flow

## Containers (26 total)
| Container | Route | Purpose |
|-----------|-------|---------|
| LoginPage | /Login | Admin authentication |
| HomePage | / | Dashboard overview |
| UserAccounts | /UserAccounts | User list & management |
| UserDetails | /Users/:id | Single user detail view |
| Billing | /Billing | Payment operations |
| Deposits | /Deposits | Deposit transactions |
| Withdrawals | /Withdrawals | Withdrawal transactions |
| OpenOrders | /OpenOrders | Active trading orders |
| FilledOrders | /FilledOrders | Completed trades |
| ExternalOrders | /ExternalOrders | External exchange orders |
| FinanceMethods | /FinanceMethods | Payment method config |
| CurrencyPairs | /CurrencyPairs | Trading pair settings |
| Reports | /Reports | Admin reports |
| Balances | /Balances | User balance overview |
| MarketTicks | /MarketTicks | Market data display |
| Admins | /Admins | Admin user management |
| VerificationWindow | /Verification | KYC document review |
| LoginHistory | /LoginHistory | User login audit trail |
| ExternalExchange | /ExternalExchange | External exchange config |
| LiquidityOrders | /LiquidityOrders | Liquidity management |
| ScanBlock | /ScanBlock | Blockchain scanner |
| Orders | /Orders | Order history |
| NavBar | (component) | Top navigation bar |
| LanguageSwitch | (component) | i18n language selector |
| ThemeSwitch | (component) | Light/dark theme toggle |
| NotFoundPage | /404 | 404 error page |

## Key Patterns

### Container Pattern (6-file convention)
Every page container follows this structure. See any container for reference.

### Service Pattern
All API calls go through `src/services/*_service.ts` → `apiService.fetchData()`.
Services are stateless functions that return `Promise<StandardResponse>`.

### Grid Pattern
Most list pages use AG Grid via the `SimpleGrid` component.
Grid data is set via MessageService events after API calls.
```

## Validation

The file should be readable by both humans and AI agents. No code compilation needed.
