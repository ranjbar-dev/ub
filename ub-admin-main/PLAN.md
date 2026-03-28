# PLAN.md — ub-admin-main Code Quality & AI-Readability Improvement Plan

> **Purpose**: This document is the single source of truth for all code quality improvements in ub-admin-main. AI agents should read this before making any changes.

---

## Table of Contents

1. [Project Context](#1-project-context)
2. [Architecture Overview](#2-architecture-overview)
3. [Domain Glossary](#3-domain-glossary)
4. [File Dependency Map](#4-file-dependency-map)
5. [Issues & Fixes — Phase 1: Type Safety Foundation](#5-phase-1-type-safety-foundation)
6. [Issues & Fixes — Phase 2: Service Layer](#6-phase-2-service-layer)
7. [Issues & Fixes — Phase 3: Saga Error Handling](#7-phase-3-saga-error-handling)
8. [Issues & Fixes — Phase 4: Component Props & Docs](#8-phase-4-component-props--docs)
9. [Issues & Fixes — Phase 5: Code Organization](#9-phase-5-code-organization)
10. [Issues & Fixes — Phase 6: State Management Cleanup](#10-phase-6-state-management-cleanup)
11. [Issues & Fixes — Phase 7: Documentation](#11-phase-7-documentation)
12. [Conventions Reference](#12-conventions-reference)
13. [Execution Dependencies](#13-execution-dependencies)

---

## 1. Project Context

| Field | Value |
|-------|-------|
| **Project** | ub-admin-main — Admin panel for UnitedBit cryptocurrency exchange |
| **Stack** | React 17 · TypeScript 4.9 · Redux Toolkit 1.3.6 · Redux-Saga 1.1.3 · Material-UI 4.10 · AG Grid 23 |
| **Build** | CRA (react-scripts 5.0.1) · Node ≥18 |
| **Dev** | `npm start` (port 3000) |
| **Build** | `npm run build` (needs `NODE_OPTIONS=--openssl-legacy-provider`) |
| **Test** | `npm test` (Jest + RTL, 90% coverage threshold) |
| **Typecheck** | `npm run checkTs` |
| **Lint** | `npm run lint` / `npm run lint:fix` |
| **Base URL config** | `src/services/constants.ts` → hardcoded `https://admin.unitedbit.com/api/v1/` |
| **Auth** | JWT in localStorage (`access_token` key) |

---

## 2. Architecture Overview

### Directory Structure

```
src/
├── app/
│   ├── containers/          # 26 page-level containers (see §2.1)
│   ├── components/          # 33 shared UI components (see §2.2)
│   ├── constants.ts         # AppPages enum, Buttons enum, WindowTypes
│   ├── appSelectors.ts      # App-level selectors
│   ├── index.tsx            # Root App — routing, global injections, MessageService subscriptions
│   ├── NewWindowContainer.tsx
│   └── ForceStyles.tsx
├── services/                # 12 API service files + constants (see §2.3)
├── store/                   # Redux store setup
│   ├── configureStore.ts    # Store creation with saga middleware + router middleware + injectors
│   ├── reducers.ts          # Combines injected reducers + connectRouter + globalReducer
│   └── slice.ts             # Global slice (loggedIn state)
├── types/                   # Shared types
│   ├── RootState.ts         # Root state interface (all 24 container state slices)
│   ├── Repo.d.ts
│   └── index.ts             # Re-exports RootState
├── utils/                   # Helper utilities
│   ├── commonUtils.ts       # omit()
│   ├── formatters.ts        # queryStringer, CurrencyFormater, Format, FormatDate, censor, etc.
│   ├── fileDownload.js      # ⚠️ JS file (not TS) — blob download utility
│   ├── stylers.ts           # stateStyler (status→color), cellColorAndNameFormatter
│   ├── loading.ts           # SetLoading helper via MessageService
│   ├── history.ts           # createBrowserHistory export
│   ├── loadable.tsx         # lazyLoad HOC for code splitting
│   ├── request.ts           # Unused fetch wrapper with ResponseError class
│   ├── redux-injectors.ts   # Re-export of redux-injectors
│   ├── gridUtilities/       # AG Grid DOM helpers (headerHider, RowWithShaddow, etc.)
│   └── hooks/               # Custom hooks (useDimensions, useOpenWithdrawWindow)
├── styles/                  # Theme + global styles
├── locales/                 # i18next (en/de)
└── images/                  # SVG icons + themed icons
```

### 2.1 Container Pattern

Each page container lives in `src/app/containers/<Name>/` with these files:

| File | Purpose | Required |
|------|---------|----------|
| `index.tsx` | React component with hooks, `useInjectReducer`, `useInjectSaga` | ✅ |
| `saga.ts` | Redux-Saga generators — API calls via `yield call()`, MessageService dispatches | ✅ |
| `slice.ts` | Redux Toolkit slice — action creators (reducers are mostly empty trigger-only) | ✅ |
| `selectors.ts` | Memoized selectors via `createSelector` | ✅ |
| `types.ts` | TypeScript interfaces for container state | ✅ |
| `Loadable.tsx` | Code-splitting wrapper via `lazyLoad()` | ✅ |
| `components/` | Sub-components (optional, used by 10 containers) | Optional |

**26 Containers:** Admins, Balances, Billing, CurrencyPairs, Deposits, ExternalExchange, ExternalOrders, FilledOrders, FinanceMethods, HomePage, LanguageSwitch, LiquidityOrders, LoginHistory, LoginPage, MarketTicks, NavBar, NotFoundPage, OpenOrders, Orders, Reports, ScanBlock, ThemeSwitch, UserAccounts, UserDetails, VerificationWindow, Withdrawals

**Compliance status:**
- ✅ **12 fully compliant** (all 6 files): Admins, Deposits, OpenOrders, Withdrawals, FilledOrders, MarketTicks, CurrencyPairs, LiquidityOrders, FinanceMethods, UserDetails, VerificationWindow, LoginHistory
- ⚠️ **1 missing Loadable.tsx**: ScanBlock
- ℹ️ **3 minimal containers** (presentational-only, expected): NavBar, LanguageSwitch, ThemeSwitch
- ℹ️ **1 non-standard**: NotFoundPage (has `P.ts` extra file)

### 2.2 Shared Components (33 directories in `src/app/components/`)

```
BackupDatePicker, clickableEmail, ConstructiveModal, CountryDropDown,
Customized (react-toastify), FormLabel, GridFilter, GridTabs, grid_loading,
inputWithValidator, isLoadingWithText, Link, loadingInButton, LoadingIndicator,
materialModal, newWindow, PageWrapper, PaginationComponent, PlaceHolder, Radio,
RawInput, renderer, sideNav, SimpleGrid, titledContainer, TNewWindow,
UbCheckBox, UbDropDown, UBInput, UbModal, UbTextAreaAutoHeight,
UserDetailsWindow, wrappers
```

### 2.3 Service Files (12 in `src/services/`)

| File | Functions | Purpose |
|------|-----------|---------|
| `api_service.ts` | 4 | **Core** — Singleton HTTP client (fetch-based), JWT from localStorage |
| `user_management_service.ts` | 29 | User CRUD, KYC, balances, billing, payments, orders |
| `orders_service.ts` | 7 | Cancel, fulfill, deposit update, balances, internal transfers |
| `order_management_service.ts` | 2 | Liquidity orders |
| `security_service.ts` | 1 active (14 lines dead code) | Login only |
| `global_data_service.ts` | 3 | Countries, currencies, admin managers list |
| `admin_reports_service.ts` | 6 | Reports, commissions, admin comments |
| `external_orders_service.ts` | 6 | External exchange order management |
| `profile_image_service.ts` | 1 | Profile image upload |
| `toastService.ts` | 1 | Toast message parser |
| `message_service.ts` | 3 services | RxJS Subject/ReplaySubject/BehaviorSubject pubsub (60+ events) |
| `constants.ts` | — | RequestTypes, RequestParameters, LocalStorageKeys, StandardResponse, BaseUrl |

### 2.4 State Flow (Current — Hybrid)

**Path A: Redux (intended pattern)**
```
Component dispatches action → Saga intercepts → Saga calls API service → 
Service calls ApiService.fetchData → Response returned → Saga dispatches result → 
Reducer updates store → Selector picks data → Component re-renders
```

**Path B: RxJS MessageService (legacy pattern — used in ~20 containers)**
```
Saga calls API → Saga calls MessageService.send({name, payload}) →
RxJS Subject.next() → Component's Subscriber.subscribe() callback fires →
Component sets local state via useState/useRef
```

**⚠️ Problem:** Most containers use Path B — data flows through RxJS instead of Redux store, making state invisible to Redux DevTools and unpredictable for AI agents.

---

## 3. Domain Glossary

### Core Entities

| Term | Definition |
|------|-----------|
| **User** | A registered customer of the exchange platform. Has email, KYC status, wallets. |
| **Admin** | System administrator. Can manage users, orders, transactions, KYC. |
| **Wallet** | User's cryptocurrency balance for a specific coin (BTC, ETH, etc.) |
| **Hot Wallet** | Exchange-connected wallet for active trading |
| **Cold Wallet** | Offline wallet for secure long-term storage |

### Transaction Types

| Term | Definition |
|------|-----------|
| **Deposit** | User sends crypto to the exchange wallet address |
| **Withdrawal** | Exchange sends crypto from wallet to user's external address |
| **Internal Transfer** | Movement between a user's own exchange wallets |
| **Payment** | Generic term for any deposit/withdrawal/transfer in the system |

### Order Types

| Term | Definition |
|------|-----------|
| **Open Order** | An order that has been placed but not yet fully executed |
| **Filled Order** | An order that has been fully executed (completed trade) |
| **Limit Order** | Buy/sell at a specific price or better |
| **Market Order** | Buy/sell at the current market price immediately |
| **External Order** | An order placed on an external exchange (liquidity bridging) |
| **Liquidity Order** | Orders that provide liquidity to the order book |
| **Net Queue** | Queue of external orders pending execution |

### KYC / Verification

| Term | Definition |
|------|-----------|
| **KYC** | Know Your Customer — identity verification process |
| **Identity Confirmation** | Verification of user's ID document (passport, national ID) |
| **Address Confirmation** | Verification of user's proof of address |
| **Phone Confirmation** | Verification of user's phone number |

**Confirmation status values:** `NotConfirmed`, `Confirmed`, `Incomplete`, `Rejected`

### Market Data

| Term | Definition |
|------|-----------|
| **Currency Pair** | A trading pair like BTC/USD, ETH/BTC |
| **OHLC** | Open, High, Low, Close — candlestick price data |
| **Market Ticks** | Real-time price update data points |
| **Sync** | Process of synchronizing OHLC data from external sources |

### Payment Status Values

```
created → in_progress → completed
                      → rejected (by admin)
                      → canceled (by system)
                      → user_canceled (by user)
                      → pending (awaiting action)
```

### API Endpoint Map

| Prefix | Domain |
|--------|--------|
| `/api/v1/user/` | User management, KYC, profile, permissions |
| `/api/v1/payment/` | Deposits, withdrawals, payment details |
| `/api/v1/order/` | Open orders, cancel, fulfill |
| `/api/v1/trade/` | Trade history |
| `/api/v1/crypto-balance/` | User wallet balances |
| `/api/v1/crypto-internal-transfer/` | Internal transfers between wallets |
| `/api/v1/currency/` | Finance methods, currency pairs |
| `/api/v1/exchange/` | External exchange data and order routing |
| `/api/v1/ohlc/` | Market ticks, OHLC data, sync |
| `/api/v1/wallet/block/scan` | Blockchain scanning |
| `/api/v1/auth/login` | Authentication |

---

## 4. File Dependency Map

### Critical Path (changes here affect everything)

```
src/services/api_service.ts          → Used by ALL 12 service files
src/services/constants.ts            → Used by ALL services + most containers
src/services/message_service.ts      → Used by ~20 containers + SimpleGrid + sagas
src/types/RootState.ts               → Used by ALL containers (state type)
src/store/configureStore.ts          → App initialization
src/store/reducers.ts                → All reducer injection
src/store/slice.ts                   → Global app state (loggedIn)
src/app/index.tsx                    → Routing, global injections, subscriptions
```

### Medium Risk (changes affect multiple files)

```
src/utils/formatters.ts              → Used by many containers for display formatting
src/utils/stylers.ts                 → Used by AG Grid cell renderers across containers
src/utils/gridUtilities/index.tsx    → Used by all containers with data grids
src/app/components/SimpleGrid/       → Core grid wrapper used by 15+ containers
src/app/components/GridFilter/       → Used in most grid-based pages
```

### Safe to Modify (isolated impact)

```
Individual container components/     → Only affect their own page
Individual service API functions     → Only used by their respective saga
src/utils/commonUtils.ts             → 1 function, few consumers
src/utils/loading.ts                 → Small helper, limited usage
```

---

## 5. Phase 1: Type Safety Foundation

### 5.1 — Enable `noImplicitAny` in tsconfig.json

**Severity:** 🔴 CRITICAL  
**File:** `tsconfig.json`

**Current:**
```json
{
  "compilerOptions": {
    "noImplicitAny": false,
    "strict": true
  }
}
```

**Problem:** `noImplicitAny: false` undermines the entire `"strict": true` setting. Allows implicit `any` everywhere.

**Fix:**
```json
{
  "compilerOptions": {
    "noImplicitAny": true,
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noImplicitReturns": true,
    "noFallthroughCasesInSwitch": true
  }
}
```

**⚠️ Execution note:** This will surface 100+ type errors. Fix those FIRST before enabling, or enable incrementally with `// @ts-expect-error` annotations that can be resolved later. Recommended approach:
1. Run `npm run checkTs` to get baseline
2. Enable `noImplicitAny: true`
3. Run again, count errors
4. Fix errors file-by-file (services first, then sagas, then components)

---

### 5.2 — Make StandardResponse Generic

**Severity:** 🔴 CRITICAL  
**File:** `src/services/constants.ts`

**Current:**
```typescript
export interface StandardResponse {
  status: boolean;
  message: string;
  data: any;          // ← always any, never typed
  token?: string;
}
```

**Fix:**
```typescript
export interface StandardResponse<T = unknown> {
  status: boolean;
  message: string;
  data: T;
  token?: string;
  errors?: Record<string, string[]>;
}
```

**Impact:** All consumers that use `StandardResponse` continue working (default `unknown`). New code can specify: `StandardResponse<UserListData>`.

---

### 5.3 — Type RequestParameters.data

**Severity:** 🔴 CRITICAL  
**File:** `src/services/constants.ts`

**Current:**
```typescript
export interface RequestParameters {
  requestType: RequestTypes;
  url: string;
  data: any;              // ← always any
  isRawUrl?: boolean;
  requestName?: string;
}
```

**Fix:**
```typescript
export interface RequestParameters<T = Record<string, unknown>> {
  requestType: RequestTypes;
  url: string;
  data: T;
  isRawUrl?: boolean;
  requestName?: string;
}
```

---

### 5.4 — Fix Inconsistent Enum Casing

**Severity:** 🟡 HIGH  
**Files:**
- `src/app/containers/Deposits/types.ts`
- `src/app/containers/OpenOrders/types.ts`

**Current (Deposits):**
```typescript
export enum DepositStatusStrings {
  Completed = 'Completed',
  InProgress = 'In Progress',
  Rejected = 'Rejected',
  Created = 'Created',
  COMPLETED = 'COMPLETED',     // ← duplicate meaning, different case
  REJECTED = 'REJECTED',       // ← duplicate meaning
  CONFIRMED = 'CONFIRMED',
}
```

**Current (OpenOrders):**
```typescript
export enum Sides {
  Buy = 'Buy',
  BUY = 'BUY',      // ← duplicate
  Sell = 'Sell',
  SELL = 'SELL',     // ← duplicate
}
```

**Fix:** Consolidate to single casing. Keep the values the API returns, but use single enum keys:
```typescript
export enum DepositStatus {
  Completed = 'Completed',
  InProgress = 'In Progress',
  Rejected = 'Rejected',
  Created = 'Created',
  Confirmed = 'CONFIRMED',
}

export enum OrderSide {
  Buy = 'Buy',
  Sell = 'Sell',
}
```

Then update `src/utils/stylers.ts` to handle both casings in the switch statement (normalize input to lowercase before matching).

---

### 5.5 — Type Empty Container States

**Severity:** 🟡 HIGH  
**Files:** All containers with `export interface XState {}` (empty)

**Affected:** Billing, Deposits, OpenOrders, FilledOrders, ExternalOrders, FinanceMethods, CurrencyPairs, ExternalExchange, MarketTicks, Admins, HomePage, Balances, LiquidityOrders, ScanBlock, LoginHistory, Reports, Orders

**Current pattern:**
```typescript
// Most containers — empty state because data flows via MessageService
export interface BillingState {}
export type ContainerState = BillingState;
```

**Fix:** Define actual state shapes even if data currently flows via MessageService. This prepares for the MessageService→Redux migration (Phase 6) and helps AI agents understand data shapes:

```typescript
// Example: src/app/containers/Billing/types.ts
export interface BillingState {
  billingData: Payment[] | null;
  depositsData: Payment[] | null;
  withdrawalsData: Payment[] | null;
  allTransactionsData: Payment[] | null;
  selectedPaymentDetails: PaymentDetails | null;
  isLoading: boolean;
  error: string | null;
}
```

---

### 5.6 — Type RootState.global

**Severity:** 🟡 HIGH  
**File:** `src/types/RootState.ts`

**Current:**
```typescript
export interface RootState {
  global?: any;       // ← untyped
  router?: RouterState;
  // ...
}
```

**Fix:**
```typescript
import { GlobalState } from 'store/slice';

export interface RootState {
  global?: GlobalState;
  router?: RouterState;
  // ...
}
```

And in `src/store/slice.ts`:
```typescript
export interface GlobalState {
  loggedIn: boolean;
}
```

---

### 5.7 — Convert fileDownload.js to TypeScript

**Severity:** 🟠 MEDIUM  
**File:** `src/utils/fileDownload.js` → `src/utils/fileDownload.ts`

**Current:** Plain JavaScript, no types, uses deprecated `msSaveBlob`.

**Fix:** Rename to `.ts`, add types:
```typescript
export const FileDownloader = (data: Blob, filename: string): void => {
  // implementation with types
};

export const downloadFile = async (params: {
  url: string;
  filename: string;
}): Promise<void> => {
  // implementation with types
};
```

Remove `msSaveBlob` IE11 code (no longer needed).

---

### 5.8 — Fix AppPages Typo

**Severity:** 🟠 MEDIUM  
**File:** `src/app/constants.ts`

**Current:**
```typescript
export enum AppPages {
  Deopsits = '/Deopsits',   // ← typo
}
```

**Fix:**
```typescript
export enum AppPages {
  Deposits = '/Deposits',
}
```

**⚠️ Impact:** Must also update the route path in `src/app/index.tsx` and any navigation links that reference this enum value. Search for `Deopsits` across the codebase.

---

## 6. Phase 2: Service Layer

### 6.1 — Type All Service Function Parameters

**Severity:** 🔴 CRITICAL  
**Files:** All 12 files in `src/services/`

**Current pattern (ALL 50+ functions):**
```typescript
export const GetUserAccountsAPI = (parameters: any) => {
  return apiService.fetchData({
    data: parameters,
    url: 'user/',
    requestType: RequestTypes.GET,
  });
};
```

**Fix pattern — create interfaces for each API call:**

```typescript
// Create src/services/types/ directory with domain-specific types

// src/services/types/userTypes.ts
export interface GetUserAccountsParams {
  page?: number;
  limit?: number;
  status?: string;
  search?: string;
}

export interface UserListResponse {
  count: number;
  users: User[];
}

// Then in src/services/user_management_service.ts:
import { GetUserAccountsParams, UserListResponse } from './types/userTypes';

export const GetUserAccountsAPI = (
  parameters: GetUserAccountsParams
): Promise<StandardResponse<UserListResponse>> => {
  return apiService.fetchData({
    data: parameters,
    url: 'user/',
    requestType: RequestTypes.GET,
  });
};
```

**Service function count requiring types:**
| Service File | Functions to Type |
|-------------|-------------------|
| `user_management_service.ts` | 29 |
| `orders_service.ts` | 7 |
| `external_orders_service.ts` | 6 |
| `admin_reports_service.ts` | 6 |
| `global_data_service.ts` | 3 |
| `order_management_service.ts` | 2 |
| `security_service.ts` | 1 |
| `profile_image_service.ts` | 1 |
| `toastService.ts` | 1 |
| **Total** | **56** |

**Execution approach:** Work through one service file at a time. For each function:
1. Determine what params the API expects (check saga callers for clues)
2. Create param interface in `src/services/types/`
3. Add return type annotation
4. Remove `: any`

---

### 6.2 — Add JSDoc to All Service Functions

**Severity:** 🟡 HIGH  
**Files:** All service files

**Current:** Zero JSDoc comments in any service file.

**Fix — add to every exported function:**
```typescript
/**
 * Fetches paginated list of user accounts.
 *
 * @param parameters.page - Page number (1-indexed)
 * @param parameters.limit - Records per page
 * @param parameters.status - Filter by user status
 * @returns Paginated user list with total count
 *
 * @endpoint GET /api/v1/user/
 */
export const GetUserAccountsAPI = (parameters: GetUserAccountsParams) => { ... };
```

**Minimum JSDoc per function:**
- One-line description of what the API does
- `@param` for each parameter field
- `@returns` describing the response shape
- `@endpoint` with HTTP method + URL path

---

### 6.3 — Remove Dead Code from security_service.ts

**Severity:** 🟡 HIGH  
**File:** `src/services/security_service.ts`

**Current:** 1 active function (`loginAPI`), 14 lines of commented-out functions (getUserDataAPI, changePasswordAPI, registerAPI, set2FaAPI, getRecapchaKeyAPI, acountActivationAPI, forgotPasswordAPI, resetPasswordAPI).

**Fix:** Delete all commented-out code. If these functions are needed later, they can be recreated from the API documentation.

---

### 6.4 — Fix Duplicate API Functions

**Severity:** 🟡 HIGH  

**Issue 1:** `GetWithdrawalCommentsAPI` exists in BOTH:
- `src/services/user_management_service.ts`
- `src/services/admin_reports_service.ts`

**Fix:** Keep in one file only (admin_reports_service.ts since it's report-related), update all imports.

**Issue 2:** In `src/services/external_orders_service.ts`, two functions use the same endpoint:
- `ChangeNetQueueStatusAPI` → `exchange/order/change-status`
- `CancelNetQueueAPI` → `exchange/order/change-status`

**Fix:** If they do the same thing, consolidate. If the payload differs, rename for clarity and add JSDoc explaining the difference.

---

### 6.5 — Refactor ApiService Error Handling

**Severity:** 🔴 CRITICAL  
**File:** `src/services/api_service.ts`

**Current problems:**
1. 422 errors return raw JSON silently (no error thrown)
2. 401 errors send MessageService event but return nothing → caller gets `undefined`
3. 500 errors show toast but return nothing → caller gets `undefined`
4. No try-catch around `fetch()` → network errors crash silently
5. In non-production mode, errors return `new Error(...)` but success returns `rawResponse.json()` → inconsistent types
6. Token read from localStorage on every request (no caching)

**Fix — complete ApiService rewrite:**

```typescript
// src/services/api_service.ts

export class ApiError extends Error {
  constructor(
    public statusCode: number,
    public responseData: Record<string, unknown> | null,
    message?: string
  ) {
    super(message || `API Error: ${statusCode}`);
    this.name = 'ApiError';
  }
}

export class ApiService {
  private static instance: ApiService;
  private constructor() {}

  public static getInstance(): ApiService {
    if (!ApiService.instance) {
      ApiService.instance = new ApiService();
    }
    return ApiService.instance;
  }

  public baseUrl = BaseUrl;

  private getToken(): string {
    return localStorage.getItem(LocalStorageKeys.ACCESS_TOKEN) || '';
  }

  public async fetchData<T = unknown>(
    params: RequestParameters
  ): Promise<StandardResponse<T>> {
    const url = params.isRawUrl ? params.url : this.baseUrl + params.url;
    const token = this.getToken();

    if (process.env.NODE_ENV !== 'production') {
      console.log(`🚀 ${params.requestType} → ${url}`, params.data);
    }

    try {
      let query = '';
      if (params.requestType === RequestTypes.GET && params.data && Object.keys(params.data).length > 0) {
        query = queryStringer(params.data);
      }

      const response = await fetch(
        params.requestType === RequestTypes.GET ? url + query : url,
        {
          method: params.requestType,
          headers: this.buildHeaders(token),
          ...(params.requestType !== RequestTypes.GET && {
            body: JSON.stringify(params.data),
          }),
        }
      );

      if (!response.ok) {
        return this.handleErrorResponse(response);
      }

      const json = await response.json();

      if (process.env.NODE_ENV !== 'production') {
        console.log(`✅ ${params.requestType} ← ${url}`, json);
      }

      return json as StandardResponse<T>;
    } catch (error) {
      // Network error (no response at all)
      if (process.env.NODE_ENV !== 'production') {
        console.error(`⛔ ${params.requestType} ${url}`, error);
      }
      throw new ApiError(0, null, 'Network error — check your connection');
    }
  }

  private async handleErrorResponse(response: Response): Promise<never> {
    let responseData: Record<string, unknown> | null = null;
    try {
      responseData = await response.json();
    } catch {
      // Response body not JSON-parseable
    }

    if (response.status === 401) {
      MessageService.send({ name: MessageNames.SETLOADING, payload: false });
      MessageService.send({ name: MessageNames.AUTH_ERROR_EVENT });
    }

    if (response.status === 500) {
      toast.error('Server error — please try again later');
    }

    // Always throw — let callers handle specific status codes
    throw new ApiError(response.status, responseData);
  }

  private buildHeaders(token: string): Record<string, string> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    };
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }
    return headers;
  }
}
```

**⚠️ Breaking change:** This changes error behavior. All saga callers must add try-catch (see Phase 3). For 422 responses, the ApiError will contain the validation errors in `responseData`.

---

### 6.6 — Move URLs to Environment Variables

**Severity:** 🟠 MEDIUM  
**File:** `src/services/constants.ts`

**Current:**
```typescript
export const appUrl = 'https://admin.unitedbit.com';
export const BaseUrl = appUrl + '/api/v1/';
const prefix = process.env.NODE_ENV === 'development' ? 'dev-' : '';
export const webAppAddress = `https://${prefix}app.unitedbit.com/api/v1/`;
```

**Fix:**
```typescript
export const BaseUrl = process.env.REACT_APP_API_BASE_URL || 'https://admin.unitedbit.com/api/v1/';
export const webAppAddress = process.env.REACT_APP_WEB_APP_URL || 'https://app.unitedbit.com/api/v1/';
```

Update `.env.example`:
```env
REACT_APP_API_BASE_URL=https://admin.unitedbit.com/api/v1/
REACT_APP_WEB_APP_URL=https://app.unitedbit.com/api/v1/
```

---

## 7. Phase 3: Saga Error Handling

### 7.1 — Create Standard Saga Error Handler

**Severity:** 🟡 HIGH  
**New file:** `src/utils/sagaUtils.ts`

```typescript
import { toast } from 'app/components/Customized/react-toastify';
import { MessageService, MessageNames } from 'services/message_service';
import { ApiError } from 'services/api_service';

/**
 * Standard error handler for saga catch blocks.
 * Shows a toast, stops loading indicator, and logs in dev mode.
 */
export const handleSagaError = (
  error: unknown,
  operation: string,
): void => {
  const message = error instanceof ApiError
    ? (error.responseData as any)?.message || `Failed to ${operation}`
    : `Failed to ${operation}`;

  toast.error(message);
  MessageService.send({ name: MessageNames.SETLOADING, payload: false });

  if (process.env.NODE_ENV !== 'production') {
    console.error(`[Saga] ${operation}:`, error);
  }
};
```

---

### 7.2 — Add try-catch to All Sagas

**Severity:** 🟡 HIGH  
**Files:** All 15+ `saga.ts` files across containers

**Current pattern (no error handling):**
```typescript
export function* GetBillingGridData(action: { type: string; payload: any }) {
  const response: StandardResponse = yield call(GetBillingGridDataAPI, action.payload);
  if (response.status === true) {
    MessageService.send({ name: MessageNames.SET_BILLING_DATA, payload: response.data });
  }
  // ← no else, no catch — errors silently swallowed
}
```

**Fix pattern (standard for ALL sagas):**
```typescript
import { handleSagaError } from 'utils/sagaUtils';
import { ApiError } from 'services/api_service';

export function* GetBillingGridData(action: {
  type: string;
  payload: { user_id: number; type?: string };
}) {
  try {
    const response: StandardResponse = yield call(
      GetBillingGridDataAPI,
      action.payload,
    );
    if (response.status === true) {
      MessageService.send({
        name: MessageNames.SET_BILLING_DATA,
        payload: response.data,
        userId: action.payload.user_id,
      });
    } else {
      toast.error(response.message || 'Failed to load billing data');
    }
  } catch (error) {
    handleSagaError(error, 'load billing data');
  }
}
```

**Sagas requiring this fix (15 files):**
1. `Admins/saga.ts`
2. `Balances/saga.ts`
3. `Billing/saga.ts` (8 generators — largest saga file, 333 lines)
4. `CurrencyPairs/saga.ts`
5. `Deposits/saga.ts`
6. `ExternalExchange/saga.ts`
7. `ExternalOrders/saga.ts`
8. `FilledOrders/saga.ts`
9. `FinanceMethods/saga.ts`
10. `LiquidityOrders/saga.ts`
11. `LoginHistory/saga.ts`
12. `MarketTicks/saga.ts`
13. `OpenOrders/saga.ts`
14. `Orders/saga.ts` (was UserDetails before)
15. `Reports/saga.ts`
16. `UserAccounts/saga.ts`
17. `UserDetails/saga.ts`
18. `VerificationWindow/saga.ts`
19. `Withdrawals/saga.ts`

---

## 8. Phase 4: Component Props & Docs

### 8.1 — Fix Component Props Typing

**Severity:** 🟡 HIGH  
**Files:** 20+ component files in `src/app/components/`

**Common violations to fix:**

| Bad Pattern | Fix |
|------------|-----|
| `children: any` | `children: React.ReactNode` |
| `onChange: Function` | `onChange: (value: string) => void` |
| `onClose: Function` | `onClose: () => void` |
| `style?: any` | `style?: React.CSSProperties` |
| `gridApi: any` | `gridApi: GridApi` (from ag-grid) |
| `props: any` | Define explicit interface |
| `event => {` (untyped) | `(event: React.ChangeEvent<HTMLInputElement>) => {` |

**Priority components to fix (most widely used):**
1. `SimpleGrid/SimpleGrid.tsx` — Used by 15+ pages
2. `GridFilter/GridFilter.tsx` — Used in most grid pages
3. `UbModal/UbModal.tsx` — Used for all modals
4. `UBInput/UBInput.tsx` — Used in forms
5. `sideNav/index.tsx` — Main navigation
6. `UserDetailsWindow/UserDetailsWindow.tsx` — User detail popup

---

### 8.2 — Add JSDoc to All Shared Components

**Severity:** 🟡 HIGH  
**Files:** All 33 component directories

**Minimum per component:**
```typescript
/**
 * SimpleGrid — AG Grid wrapper with pagination, filtering, and loading states.
 *
 * Subscribes to MessageService for data updates. Uses `initialAction` to trigger
 * the first data fetch via Redux saga on mount.
 *
 * @param messageName - MessageService event name to listen for data updates
 * @param initialAction - Redux action to dispatch on mount for initial data
 * @param staticRows - AG Grid ColDef array defining columns
 * @param arrayFieldName - Key name in response data containing the row array
 * @param pagination - Enable pagination controls (default: false)
 *
 * @example
 * <SimpleGrid
 *   messageName={MessageNames.SET_BILLING_DATA}
 *   initialAction={BillingActions.GetBillingGridDataAction}
 *   staticRows={billingColumns}
 *   arrayFieldName="payments"
 *   pagination
 * />
 */
```

---

## 9. Phase 5: Code Organization

### 9.1 — Standardize File Naming

**Severity:** 🟡 HIGH

**Current mess:**
- Services: `snake_case` (`user_management_service.ts`)
- Containers: `PascalCase` directories (`UserAccounts/`)
- Components: Mixed (`clickableEmail/`, `UBInput/`, `grid_loading/`)
- Sub-components: Mixed (`segmentSelector.tsx` vs `RejectModal.tsx`)

**Convention to adopt:**

| File Type | Convention | Example |
|-----------|-----------|---------|
| Container directories | PascalCase | `UserAccounts/` ✅ (already correct) |
| Container files | camelCase | `index.tsx`, `saga.ts`, `slice.ts` ✅ (already correct) |
| Component directories | PascalCase | `SimpleGrid/` → keep, `clickableEmail/` → `ClickableEmail/` |
| Component files (.tsx) | PascalCase | `SimpleGrid.tsx` ✅, `mainCat.tsx` → `MainCat.tsx` |
| Service files | camelCase | `user_management_service.ts` → `userManagementService.ts` |
| Utility files | camelCase | ✅ (already correct) |
| Type files | camelCase | ✅ (already correct) |

**Renames required (services):**
```
src/services/user_management_service.ts  → userManagementService.ts
src/services/order_management_service.ts → orderManagementService.ts
src/services/security_service.ts         → securityService.ts
src/services/global_data_service.ts      → globalDataService.ts
src/services/admin_reports_service.ts    → adminReportsService.ts
src/services/external_orders_service.ts  → externalOrdersService.ts
src/services/profile_image_service.ts    → profileImageService.ts
src/services/api_service.ts              → apiService.ts
src/services/message_service.ts          → messageService.ts
```

**Renames required (components):**
```
src/app/components/clickableEmail/       → ClickableEmail/
src/app/components/grid_loading/         → GridLoading/
src/app/components/inputWithValidator/   → InputWithValidator/
src/app/components/isLoadingWithText/    → IsLoadingWithText/
src/app/components/loadingInButton/      → LoadingInButton/
src/app/components/materialModal/        → MaterialModal/
src/app/components/newWindow/            → NewWindow/
src/app/components/titledContainer/      → TitledContainer/
src/app/components/wrappers/             → Wrappers/
```

**⚠️ Execution:** For each rename, update ALL import statements across the codebase. Use a find-and-replace across `src/` for each old path.

---

### 9.2 — Standardize Import Ordering

**Severity:** 🟠 MEDIUM  
**Files:** All `.tsx` and `.ts` files

**Convention (3 groups, separated by blank line):**

```typescript
// 1. External packages (npm modules)
import React, { memo, useEffect, useState } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { useTranslation } from 'react-i18next';
import styled from 'styled-components/macro';

// 2. Absolute imports from src/ (via baseUrl in tsconfig)
import { useInjectReducer, useInjectSaga } from 'utils/redux-injectors';
import { MessageService, MessageNames, Subscriber } from 'services/message_service';
import { selectUserAccounts } from 'app/containers/UserAccounts/selectors';
import { translations } from 'locales/i18n';

// 3. Relative imports (same container/component)
import { reducer, sliceKey, actions } from './slice';
import { selectBilling } from './selectors';
import { billingSaga } from './saga';
import { BillingState } from './types';
```

---

### 9.3 — Remove All eslint-disable Comments

**Severity:** 🟠 MEDIUM  
**Files:** ~40 occurrences across codebase

**Current:** 4 whole-file `/*eslint-disable*/` + ~36 inline `// eslint-disable-next-line`

**Fix approach:**
- **Whole-file disables** (message_service.ts, constants.ts, validators): Remove the disable, fix underlying lint issues
- **`@typescript-eslint/no-unused-vars`** (most common): Remove unused variables instead of disabling the rule
- **`react-hooks/exhaustive-deps`**: Fix dependency arrays or add explanatory comment if intentional
- **`react-hooks/rules-of-hooks`**: Fix the hook usage (move inside component body)

---

### 9.4 — Remove Commented-Out Code

**Severity:** 🟠 MEDIUM

**Files with dead commented code:**
1. `src/services/security_service.ts` — 55 lines of commented-out API functions
2. `src/services/api_service.ts` — `extractMessages` method (lines 114-124)
3. `src/utils/formatters.ts` — `columnResize` function (lines 86-98)
4. `src/utils/loadable.tsx` — Alternative implementation (lines 31-54)

**Fix:** Delete all commented-out code blocks. Git history preserves them if needed.

---

## 10. Phase 6: State Management Cleanup

### 10.1 — Document MessageService Events

**Severity:** 🟡 HIGH  
**File:** `src/services/message_service.ts`

**Current:** 67 enum values in `MessageNames` with zero documentation, organized by loose comment separators.

**Immediate fix (before full migration):** Add JSDoc grouping:

```typescript
export enum MessageNames {
  // === UI State ===
  /** Show/hide global loading spinner */
  SETLOADING = 'SETLOADING',
  /** Show/hide loading on a specific AG Grid row */
  SET_ROW_LOADING = 'SET_ROW_LOADING',
  /** Show/hide loading on a specific button */
  SET_BUTTON_LOADING = 'SET_BUTTON_LOADING',

  // === Navigation ===
  /** Open user details or verification in new window */
  OPEN_NEW_WINDOW = 'OPEN_NEW_WINDOW',
  OPEN_WINDOW = 'OPEN_WINDOW',
  CLOSE_WINDOW = 'CLOSE_WINDOW',

  // === Data Events (Grid Population) ===
  /** Grid data for user accounts list */
  SET_USER_ACCOUNTS = 'SET_USER_ACCOUNTS',
  // ... etc for all 67 events
}
```

---

### 10.2 — Plan MessageService → Redux Migration

**Severity:** 🟡 HIGH (architecture improvement, long-term)

**Current flow:** Saga → MessageService.send() → Component subscribes → setState

**Target flow:** Saga → dispatch(action) → Reducer updates store → Selector → Component

**Migration strategy per container:**
1. Fill in the empty state interface in `types.ts`
2. Add reducer cases to `slice.ts` for data storage
3. Add selectors to `selectors.ts` for data retrieval
4. Update saga to `yield put(actions.setData(response.data))` instead of `MessageService.send()`
5. Update component to use `useSelector()` instead of `Subscriber.subscribe()`
6. Remove the MessageService subscription from the component

**Migrate in order of simplicity (smallest containers first):**
1. LoginHistory (simple single-grid page)
2. Admins (simple single-grid page)
3. FilledOrders
4. Deposits
5. OpenOrders
6. ... up to complex ones (Billing, UserDetails, VerificationWindow)

**⚠️ Do NOT remove MessageService entirely** — some cross-component communication (OPEN_NEW_WINDOW, AUTH_ERROR_EVENT, RESIZE) may still need it. The goal is to eliminate it for data flow only.

---

## 11. Phase 7: Documentation

### 11.1 — Create ARCHITECTURE.md

**File:** `ARCHITECTURE.md` in project root

Content should cover:
- Directory structure (copy from §2 of this plan)
- Container pattern explanation
- State flow diagram (both Redux and MessageService paths)
- Service layer pattern
- How to add a new page (step by step)
- How to add a new API endpoint (step by step)

---

### 11.2 — Enhance AGENTS.md

**File:** `AGENTS.md` in project root

Add sections:
- Quick reference: which file does what
- Common tasks with file paths
- "Don't touch" list (critical path files that need extra care)
- Testing instructions
- Known gotchas (MessageService hybrid pattern, empty state types, etc.)

---

### 11.3 — Create docs/GLOSSARY.md

**File:** `docs/GLOSSARY.md`

Content: Copy §3 of this plan (Domain Glossary) into a standalone file.

---

## 12. Conventions Reference

### TypeScript

```typescript
// ✅ DO: Explicit types everywhere
export const fetchUsers = (params: FetchUsersParams): Promise<StandardResponse<UserList>> => { ... }

// ❌ DON'T: any types
export const fetchUsers = (params: any) => { ... }

// ✅ DO: Named interfaces for component props
interface UserCardProps {
  user: User;
  onEdit: (userId: number) => void;
  className?: string;
}

// ❌ DON'T: Inline or any props
const UserCard = (props: any) => { ... }

// ✅ DO: ReactNode for children
interface ModalProps {
  children: React.ReactNode;
  onClose: () => void;
}

// ❌ DON'T: any or Function
interface ModalProps {
  children: any;
  onClose: Function;
}
```

### Saga Pattern

```typescript
// ✅ STANDARD SAGA (use this for ALL new/updated sagas)
export function* FetchSomeData(action: {
  type: string;
  payload: FetchParams;
}) {
  try {
    const response: StandardResponse<DataType> = yield call(FetchAPI, action.payload);
    if (response.status === true) {
      // Prefer: yield put(actions.setData(response.data));
      // Legacy: MessageService.send({ name: MessageNames.SET_DATA, payload: response.data });
    } else {
      toast.error(response.message || 'Operation failed');
    }
  } catch (error) {
    handleSagaError(error, 'fetch some data');
  }
}
```

### Service Functions

```typescript
// ✅ STANDARD SERVICE FUNCTION
/**
 * Fetches user account list with pagination.
 * @endpoint GET /api/v1/user/
 */
export const GetUserAccountsAPI = (
  parameters: GetUserAccountsParams
): Promise<StandardResponse<UserListResponse>> => {
  return apiService.fetchData({
    data: parameters,
    url: 'user/',
    requestType: RequestTypes.GET,
  });
};
```

### File Naming

| Type | Convention | Example |
|------|-----------|---------|
| Service | camelCase | `userManagementService.ts` |
| Container dir | PascalCase | `UserAccounts/` |
| Component dir | PascalCase | `SimpleGrid/` |
| Component file | PascalCase | `SimpleGrid.tsx` |
| Util file | camelCase | `formatters.ts` |
| Type file | camelCase | `userTypes.ts` |

---

## 13. Execution Dependencies

```
Phase 1 (Type Safety)
  ├── 5.1 Enable noImplicitAny         ─┐
  ├── 5.2 Generic StandardResponse      ├── Must complete before Phase 2
  ├── 5.3 Type RequestParameters         │
  ├── 5.6 Type RootState.global         ─┘
  ├── 5.4 Fix enum casing              (independent)
  ├── 5.5 Fill empty container states   (independent, but helps Phase 6)
  ├── 5.7 Convert fileDownload.js       (independent)
  └── 5.8 Fix AppPages typo            (independent)

Phase 2 (Services)
  ├── 6.1 Type service params           ─── Depends on 5.2, 5.3
  ├── 6.2 Add JSDoc to services         (independent)
  ├── 6.3 Remove dead code              (independent)
  ├── 6.4 Fix duplicate APIs            (independent)
  ├── 6.5 Refactor ApiService errors    ─── BLOCKS Phase 3 (saga error handling)
  └── 6.6 Move URLs to env vars         (independent)

Phase 3 (Sagas)
  ├── 7.1 Create sagaUtils.ts           ─── Depends on 6.5
  └── 7.2 Add try-catch to all sagas    ─── Depends on 7.1

Phase 4 (Components)
  ├── 8.1 Fix component props           (independent)
  └── 8.2 Add JSDoc to components       (independent)

Phase 5 (Code Organization)
  ├── 9.1 Standardize file naming       (independent, but do LAST — many import changes)
  ├── 9.2 Standardize import ordering   (independent)
  ├── 9.3 Remove eslint-disable         (independent)
  └── 9.4 Remove commented-out code     (independent)

Phase 6 (State Management)
  ├── 10.1 Document MessageService      (independent)
  └── 10.2 Migrate to Redux             ─── Depends on 5.5, 7.2

Phase 7 (Documentation)
  ├── 11.1 ARCHITECTURE.md              (independent, do early)
  ├── 11.2 Enhance AGENTS.md            (independent, do early)
  └── 11.3 docs/GLOSSARY.md             (independent, do early)
```

### Quick Wins (can start immediately, no dependencies)

1. **§11.1–11.3** — Documentation (ARCHITECTURE.md, AGENTS.md, GLOSSARY.md)
2. **§6.3** — Remove dead code from security_service.ts
3. **§6.4** — Fix duplicate APIs
4. **§9.4** — Remove commented-out code blocks
5. **§5.8** — Fix AppPages.Deopsits typo
6. **§5.7** — Convert fileDownload.js to TypeScript
7. **§6.2** — Add JSDoc to service functions (no code changes)

### Critical Path (must be sequential)

```
5.2 (Generic StandardResponse) → 6.1 (Type service params) → 6.5 (ApiService errors) → 7.1 (sagaUtils) → 7.2 (all sagas)
```

---

*Last updated: 2026-03-27. Generated from comprehensive codebase audit.*
