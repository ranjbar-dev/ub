# AGENTS.md — unitedbit Flutter App (`ub-app-main`)

> **Audience**: AI coding agents. This document is the single source of truth for the entire Flutter application. Every section is exhaustive. Use this before exploring the codebase.

---

## Table of Contents

1. [Project Identity](#1-project-identity)
2. [SDK & Framework Constraints](#2-sdk--framework-constraints)
3. [Complete Dependency Inventory](#3-complete-dependency-inventory)
4. [Architecture Overview](#4-architecture-overview)
5. [Directory Structure](#5-directory-structure)
6. [Application Entry Point](#6-application-entry-point)
7. [Routing System](#7-routing-system)
8. [Global State & Controllers](#8-global-state--controllers)
9. [API Service & HTTP Layer](#9-api-service--http-layer)
10. [Interceptor Chain](#10-interceptor-chain)
11. [Centrifugo Realtime System](#11-centrifugo-realtime-system)
12. [Storage System](#12-storage-system)
13. [Cryptography](#13-cryptography)
14. [Biometrics Service](#14-biometrics-service)
15. [Constants & Configuration](#15-constants--configuration)
16. [Environment System](#16-environment-system)
17. [Feature Modules (All 23)](#17-feature-modules-all-23)
18. [Popup System](#18-popup-system)
19. [UI Component Library (70 Widgets)](#19-ui-component-library-70-widgets)
20. [Custom Widgets (49 Files)](#20-custom-widgets-49-files)
21. [Theme System](#21-theme-system)
22. [Localization](#22-localization)
23. [Mixins & Utilities](#23-mixins--utilities)
24. [Code Generation](#24-code-generation)
25. [Complete API Endpoint Reference](#25-complete-api-endpoint-reference)
26. [Centrifugo Channel Reference](#26-centrifugo-channel-reference)
27. [Common Data Provider](#27-common-data-provider)
28. [Build Scripts](#28-build-scripts)
29. [Docker Configuration](#29-docker-configuration)
30. [GitLab CI/CD Pipeline](#30-gitlab-cicd-pipeline)
31. [Platform Configuration](#31-platform-configuration)
32. [Testing](#32-testing)
33. [Known Issues & Bugs](#33-known-issues--bugs)
34. [Deep Audit Findings Summary](#34-deep-audit-findings-summary)
35. [Upgrade Roadmap](#35-upgrade-roadmap)

---

## 1. Project Identity

| Field | Value |
|---|---|
| **Name** | unitedbit |
| **Version** | 1.1.7+10 |
| **Package ID** | `com.unitedbit.app` |
| **Repo directory** | `ub-app-main/` |
| **Primary language** | Dart (pre-null-safety) |
| **Framework** | Flutter 2.x |
| **State management** | GetX 4.3.4 (primary), Provider 6.0.1 (secondary) |
| **Architecture** | GetX MVC module pattern |

---

## 2. SDK & Framework Constraints

| Constraint | Value | Notes |
|---|---|---|
| **Dart SDK** | `>=2.11.0 <3.0.0` | **Pre-null-safety**. Cannot use Dart 3.x. |
| **Flutter** | 2.x | Docker pins Flutter 2.10.5 |
| **Null safety** | ❌ Disabled | Entire codebase is pre-null-safety |
| **Sound null safety** | ❌ Not available | Requires Phase 2 migration |

> **CRITICAL**: Do NOT use `?` nullable types, `late` keyword, or `required` named parameters. All code must be Dart 2.11 compatible.

---

## 3. Complete Dependency Inventory

### Core Dependencies

| Package | Version | Purpose |
|---|---|---|
| `get` | 4.3.4 | Routing, state management, DI (bindings, controllers, .obs) |
| `dio` | 4.0.1 | HTTP client with interceptor chain |
| `centrifuge` | 0.15.0 | Centrifugo WebSocket client |
| `rxdart` | 0.27.2 | Reactive streams (BehaviorSubject) |
| `provider` | 6.0.1 | Secondary state management |
| `get_storage` | 2.0.3 | Local key-value persistence |
| `flutter_secure_storage` | 4.2.1 | Encrypted credential storage |

### Cryptography

| Package | Version | Purpose |
|---|---|---|
| `pointycastle` | 3.3.5 | RSA/AES cryptographic primitives |
| `encrypt` | 5.0.1 | Encryption helpers |
| `asn1lib` | 1.0.2 | ASN.1 parsing for key formats |

### Finance

| Package | Version | Purpose |
|---|---|---|
| `decimal` | 1.3.0 | Precise monetary arithmetic |

### UI & Widgets

| Package | Version | Purpose |
|---|---|---|
| `rive` | 0.7.33 | Rive animations |
| `cached_network_image` | 3.1.0 | Image caching |
| `webview_flutter` | 2.3.0 | Embedded WebView |
| `animations` | 2.0.2 | Material motion transitions |
| `carousel_slider` | 4.0.0 | Carousel widget |
| `flutter_svg` | 0.23.0+1 | SVG rendering |
| `shimmer` | 2.0.0 | Shimmer loading effects |
| `percent_indicator` | 3.4.0 | Circular/linear progress indicators |
| `dotted_border` | 2.0.0+1 | Dotted border decoration |
| `flutter_switch` | 0.3.2 | Custom switch widget |
| `pull_to_refresh` | 2.0.0 | Pull-to-refresh |
| `flutter_staggered_animations` | 1.0.0 | Staggered list animations |
| `animator` | 3.1.0 | Animation helpers |
| `flutter_datetime_picker` | 1.5.1 | Date/time picker |

### Platform & Device

| Package | Version | Purpose |
|---|---|---|
| `url_strategy` | 0.2.0 | PathUrlStrategy (no hash URLs on web) |
| `android_intent_plus` | 3.0.2 | Android intent launching |
| `local_auth` | 1.1.8 | Biometric authentication |
| `connectivity` | 3.0.6 | Network connectivity monitoring (**DEPRECATED**) |
| `permission_handler` | 8.2.6 | Runtime permissions |
| `ai_barcode` | 3.0.1 | QR/barcode scanning |
| `vibration` | 1.7.3 | Haptic feedback |
| `flutter_appavailability` | 0.0.21 | Check installed apps (**UNMAINTAINED**) |
| `url_launcher` | 6.0.12 | Launch URLs |
| `file_picker` | 4.2.0 | File selection |
| `qr_code_tools` | 0.0.7 | QR code utilities |
| `fluttertoast` | 8.0.8 | Toast notifications |
| `universal_html` | 2.0.8 | Cross-platform HTML |

### Utilities

| Package | Version | Purpose |
|---|---|---|
| `uuid` | 3.0.5 | UUID generation |
| `intl` | 0.17.0 | Internationalization / date formatting |
| `supercharged` | 2.1.1 | Dart collection/string extensions |
| `basic_utils` | 3.8.1 | Common utility functions |
| `transparent_image` | 2.0.0 | Transparent placeholder image |
| `pretty_dio_logger` | 1.1.1 | Dio request/response logging |
| `logger` | 1.1.0 | Structured logging |

### Dev Dependencies

| Package | Version | Purpose |
|---|---|---|
| `flutter_launcher_icons` | 0.9.2 | App icon generation |
| `flutter_native_splash` | 1.3.1 | Splash screen generation |
| `mockito` | 5.0.16 | Test mocking |
| `build_runner` | 2.1.4 | Code generation runner |
| `flutter_gen` | 4.1.2 | Asset/color/font/locale code generation |

### Dependency Overrides

| Package | Version | Reason |
|---|---|---|
| `meta` | 1.7.0 | Version conflict resolution |
| `url_launcher_web` | 2.0.6 | Version conflict resolution |

---

## 4. Architecture Overview

### GetX MVC Module Pattern

Every feature follows the same structure under `lib/app/modules/<feature>/`:

```
<feature>/
  bindings/<feature>_binding.dart    — GetX DI registration (Get.lazyPut)
  controllers/<feature>_controller.dart — Business logic + reactive state (.obs)
  views/<feature>_view.dart          — UI widget (extends GetView<Controller>)
  providers/<feature>_provider.dart   — API calls (uses ApiService singleton)
  models/                             — Data models (optional, per-module)
```

### Dependency Injection Flow

1. **Route navigated to** → GetPage triggers Binding
2. **Binding.dependencies()** → `Get.lazyPut(() => Controller())` + `Get.lazyPut(() => Provider())`
3. **Controller.onInit()** → calls Provider methods, sets up Centrifugo subscriptions
4. **View** → `GetView<Controller>` accesses controller via `controller.` getter, uses `Obx(() =>)` for reactive UI

### Global Singleton Pattern

- `GlobalBinding` registers `GlobalController` as **permanent** (never disposed)
- On login success: `loadAuthenticatedControllers()` puts `AuthorizedCentrifugoController`, `UnAuthorizedCentrifugoController`, `TradeController`, `AccountController` as permanent
- On logout: `_purgeTheMemory()` deletes all feature controllers

---

## 5. Directory Structure

```
lib/
  main.dart                              # App entry point
  configure_web.dart                     # Web: PathUrlStrategy
  configure_nonweb.dart                  # Native: no-op
  generated/                             # flutter_gen output (DO NOT EDIT)
    assets.gen.dart                      #   Rive files, images
    colors.gen.dart                      #   From assets/color/colors.xml
    fonts.gen.dart                       #   OpenSans family
    locales.g.dart                       #   English translations
  theme/
    theme.dart                           # darkThemeData + lightThemeData
  services/
    apiService.dart                      # Singleton Dio HTTP client
    constants.dart                       # URLs, topics, version, helpers
    controllerTags.dart                  # GetX controller tags
    storageKeys.dart                     # GetStorage + SecureStorage keys
    localAuthService.dart                # BiometricsService
    localizationService.dart             # i18n (English only)
    interceptors/
      connection_retry_interceptor.dart  # Retry on SocketException
      options.dart                       # RetryOptions (3 retries, 1s interval)
      request_retrier.dart               # DioConnectivityRequestRetrier
      timeout_retry_interceptor.dart     # Retry on connectTimeout errors
  centrifugoClient/
    centrifugo_service.dart              # Centrifugo client wrapper
    src/
      centrifugo_channels.dart           # Channel definitions and helpers
      centrifugo_config.dart             # Connection configuration
  utils/
    basicMath.dart                       # Math helpers
    commonUtils.dart                     # Shared utilities (launchURL etc.)
    computes.dart                        # JSON parsing via compute() (background isolate)
    debounce.dart                        # Debounce helper
    throttle.dart                        # Throttle helper
    emailValidator.dart                  # Email validation regex
    passwordValidator.dart               # Password validation
    numUtil.dart                         # Numeric utilities
    marketUtils.dart                     # Market/trading helpers
    pairAndCurrencyUtils.dart            # PairAndCurrencyUtils singleton (hash maps)
    inputCurrencyFormatter.dart          # TextInputFormatter for currency
    logger.dart                          # Logger instance
    UBDeviceType.dart                    # Device type detection
    environment/
      ubEnv.dart                         # ENV="PRODUCTION" (hardcoded!), VERSION from dart-define
    cryptography/
      encoding.dart                      # RSA encrypt with ub-captcha_ prefix, timestamp
      rsa_encryption.dart                # RsaKeyHelper: key gen, PEM parse, encrypt/decrypt
    extentions/
      basic.dart                         # String.currencyFormat(), .removeComma(), .toDouble()
    middleWares/
      authMiddleware.dart                # GetMiddleware (placeholder, always redirects to /login)
    mixins/
      commonConsts.dart                  # Spacing constants, text styles, border decorations
      filterPopups.dart                  # Order filter UI popups
      formatters.dart                    # currencyFormatter() mixin
      popups.dart                        # Confirmation/info popup helpers
      toast.dart                         # Toaster mixin
  app/
    global/
      binding/
        index.dart                       # GlobalBinding → GlobalController (permanent)
      controller/
        globalController.dart            # App-wide state, connectivity, theme, auth
        authorizedCentrifugoController.dart  # User-private channels
        unAuthorizedCentrifugoController.dart # Public channels
      providers/
        commonDataProvider.dart          # Shared API calls (countries, currencies, pairs, version)
      currency_model.dart
      currency_pairs_model.dart
      name_value_model.dart
      response_model.dart
      authorized_order_event_model.dart
      autocompleteModel.dart
    common/
      components/                        # 70 reusable widgets
      custom/                            # 49 custom widgets/libs
    popups/
      bindings/twoFaPopupBindings.dart
      controllers/twofaPopupController.dart
      veiws/twoFaPopupView.dart
    routes/
      app_routes.dart                    # 35 route constants
      app_pages.dart                     # 33 GetPage definitions
    modules/                             # 23 feature modules (see Section 17)
```

---

## 6. Application Entry Point

**File**: `lib/main.dart`

Startup sequence:
1. `configureApp()` — platform-specific (PathUrlStrategy on web, no-op on native)
2. `GetStorage.init()` — initialize local storage
3. Lock portrait orientation (`SystemChrome.setPreferredOrientations`)
4. Set status bar: black (`ColorName.black2c`)
5. `runApp()` → `GetMaterialApp`:
   - `initialBinding: GlobalBinding()` (permanent)
   - `theme: lightThemeData` / `darkTheme: darkThemeData`
   - `themeMode` from `GetStorage('lightMode')` — **NOTE**: uses `'lightMode'` key directly, not from `StorageKeys` class
   - `fallbackLocale: Locale('en', 'EN')`
   - `translationsKeys: AppTranslation.translations`
   - `initialRoute: Routes.AFTER_SPLASH`
   - `getPages: AppPages.routes`

---

## 7. Routing System

### All 35 Named Route Constants (`app_routes.dart`)

| Route | Path |
|---|---|
| `HOME` | `/home` |
| `LOGIN` | `/login` |
| `LANDING` | `/landing` |
| `SIGNUP` | `/signup` |
| `FORGOT` | `/forgot` |
| `ACCOUNT` | `/account` |
| `TRADE` | `/trade` |
| `OPEN_ORDERS` | `/open-orders` |
| `ORDERS` | `/orders` |
| `ORDER_HISTORY` | `/order-history` |
| `FUNDS` | `/funds` |
| `BALANCE` | `/balance` |
| `DEPOSITS` | `/deposits` |
| `DEPOSIT_DETAILS` | `/depost-details` |
| `WITHDRAWALS` | `/withdrawals` |
| `TRANSACTION_HISTORY` | `/transaction-history` |
| `MARKET` | `/market` |
| `CHANGE_PASSWORD` | `/change-password` |
| `WITHDRAW_ADDRESS_MANAGEMENT` | `/withdraw-address-management` |
| `TWO_FACTOR_AUTHENTICATION` | `/two-factor-authentication` |
| `IDENTITY_VERIFICATION` | `/identity-verification` |
| `IDENTITY_DOCUMENTS` | `/identity-documents` |
| `PHONE_VERIFICATION` | `/phone-verification` |
| `ADD_NEW_ADDRESS` | `/add-new-address` |
| `QR_SCAN` | `/qr-scan` |
| `CHARTS_PAGE` | `/charts-page` |
| `EDIT_FAVORITES` | `/edit-favorites` |
| `AFTER_SPLASH` | `/after-splash` |
| `WEB_VIEW` | `/web-view` |
| `WEB_VIEW_PAGE` | `/web-view-page` |
| `CHECK_YOUR_EMAIL` | `/check-your-email` |
| `SEPARATE_MESSAGE_PAGE` | `/separate-message-page` |
| `EXCHANGE` | `/exchange` |
| `EXCHANGE_SEARCH` | `/exchange-search` |
| `AUTO_EXCHANGE` | `/auto-exchange` |
| `AUTO_EXCHANGE_POPUP` | `/auto-exchange-popup` |

**Initial route**: `/after-splash` — decides `/landing` or `/home` based on auth state.

### GetPage Definitions (`app_pages.dart`)

33 `GetPage` entries, each mapping route → View + Binding. Notable special cases:
- `/funds` uses **multiple bindings**: `[FundsBinding(), BalanceBinding()]`
- `/edit-favorites` uses `MarketBinding` (not a dedicated binding)

---

## 8. Global State & Controllers

### GlobalController (Permanent Singleton)

**File**: `lib/app/global/controller/globalController.dart`
**Registered by**: `GlobalBinding` in `main.dart` (permanent — never disposed)

#### Observable State (`.obs`)

| Variable | Type | Purpose |
|---|---|---|
| `isRedirectContainerDismissed` | `RxBool` | Mobile web → native app redirect banner |
| `doShowRedirect` | `RxBool` | Whether to show redirect banner |
| `currencyPairsArray` | `RxList` | Raw currency pairs list |
| `allCurrencyPairs` | `RxList<CurrencyPairsModel>` | Typed currency pairs |
| `hasConnection` | `RxBool` | Network connectivity status |
| `loggedIn` | `RxBool` | Authentication state |

#### Non-Reactive State

| Variable | Type | Purpose |
|---|---|---|
| `isAppInstalled` | `bool` | Native app installed (for redirect) |
| `deviceType` | `DeviceTypes` | `PHONE` or `TABLET` |
| `isLoggingInWithBiometrics` | `bool` | Biometric login in progress flag |
| `pairsHashMap` | `Map<String, Pairs>` | Quick pair lookup |
| `connectivityResult` | `ConnectivityResult` | Current connection type |

#### Key Methods

| Method | Behavior |
|---|---|
| `enableDarkTheme(bool)` | Write `darkMode` to storage, call `Get.changeThemeMode` |
| `handleLoggedOut({andExitApp})` | Clear tokens, purge all controllers, navigate to `/landing` |
| `checkTokenValidation()` | 28-day token expiry check; biometric re-auth if expired |
| `loadAuthenticatedControllers()` | Put `AuthorizedCentrifugoController`, `UnAuthorizedCentrifugoController`, `TradeController`, `AccountController` as permanent |
| `getPairsCurrenciesCountriesAndVersion()` | Parallel fetch: currencies, pairs, countries |
| `getVersion()` | Check app version, show update popup if outdated |
| `_purgeTheMemory()` | Delete all feature controllers on logout |
| `checkIfRedirectIsNeeded()` | Mobile web → native app redirect logic |
| `setDeviceType()` | `PHONE` if `shortestSide < 600`, else `TABLET` |

### AuthorizedCentrifugoController

**File**: `lib/app/global/controller/authorizedCentrifugoController.dart`

| Property | Details |
|---|---|
| **Connection** | `wss://[dev-]app.unitedbit.com:8800` |
| **Username** | JWT token |
| **Password** | UUID v4 |
| **Topic 1** | `main/trade/user/{channel}/open-orders/` — order events (QoS exactlyOnce) |
| **Topic 2** | `main/trade/user/{channel}/crypto-payments/` — payment events (QoS exactlyOnce) |

#### Reactive State

| Variable | Type | Purpose |
|---|---|---|
| `ordrPayload` | `.obs AuthorizedOrderEventModel` | Latest order event |
| `updateDataSubject` | `GetStream<List<RxUpdateables>>` | Stream of update signals |

**Emits `RxUpdateables`**: `Balances`, `TransactionHistory`, `UserPairBalances`, `OpenOrders`, `OrderHistory`

### UnAuthorizedCentrifugoController

**File**: `lib/app/global/controller/unAuthorizedCentrifugoController.dart`

| Property | Details |
|---|---|
| **Connection** | `wss://[dev-]app.unitedbit.com:8800` |
| **Username** | UUID |
| **Password** | UUID |

Used by `TradeController` to subscribe to:
- `main/trade/ticker` — live price ticker for all pairs
- `main/trade/order-book/{pair}` — live order book for specific pair
- `main/trade/kline/{timeframe}/{pair}` — live OHLC candle updates

---

## 9. API Service & HTTP Layer

**File**: `lib/services/apiService.dart`

### Configuration

| Setting | Value |
|---|---|
| **Pattern** | Singleton |
| **Base URL** | `https://[dev-]app.unitedbit.com/api/v1/` |
| **Content-Type** | `application/json` |
| **Connect Timeout** | 10 seconds |
| **Receive Timeout** | 10 seconds |
| **Platform Header** | `ubandroid-v{version}` (not sent on web) |
| **JSON Decode** | Background isolate via `compute()` |

### Methods

| Method | Signature | Purpose |
|---|---|---|
| `get()` | `get(url, urlGenerator, data, rawUrl)` | GET request with optional URL builder |
| `rawGet()` | `rawGet(rawUrl)` | GET with absolute URL (bypasses base URL) |
| `post()` | `post(url, data)` | POST request |
| `upload()` | `upload(form, stream, url, cancelToken)` | Multipart file upload |

---

## 10. Interceptor Chain

Interceptors execute in this order on every request:

| Order | Interceptor | File | Status | Behavior |
|---|---|---|---|---|
| 1 | **Connection Retry** | `connection_retry_interceptor.dart` | ✅ Active | On `SocketException`, schedule retry via `DioConnectivityRequestRetrier` |
| 2 | **Token Refresh** | (inline in apiService) | ✅ Active | On 403 response, POST `auth/refresh` with refresh token; lock interceptors during refresh |
| 3 | **Auth Header** | (inline in apiService) | ✅ Active | Attach `Bearer {token}` from storage/memory |
| 4 | **Timeout Retry** | `timeout_retry_interceptor.dart` | ❌ Commented out | Would retry on `connectTimeout` errors (1 retry, 5s interval) |
| 5 | **Pretty Dio Logger** | (inline) | ❌ Commented out | Would log requests/responses in debug mode |

### Retry Configuration (`options.dart`)

| Setting | Value |
|---|---|
| **Max retries** | 3 |
| **Retry interval** | 1 second |
| **Evaluator** | `RetryEvaluator` function |

### Request Retrier (`request_retrier.dart`)

`DioConnectivityRequestRetrier` monitors connectivity stream and completes pending retry when connection is restored.

### Timeout Retry (`timeout_retry_interceptor.dart`)

`TimeoutRetryInterceptor` retries on `connectTimeout` errors.

---

## 11. Centrifugo Realtime System

**Location**: `lib/centrifugoClient/`

### Architecture

| File | Class | Purpose |
|---|---|---|
| `centrifugo_service.dart` | `CentrifugoService` | WebSocket client with auto-reconnect, JWT auth, channel subscriptions |
| `centrifugo_channels.dart` | — | Channel name helpers and definitions |
| `centrifugo_config.dart` | — | Connection configuration (URL, tokens) |

### Key Methods

| Method | Purpose |
|---|---|
| `subscribe(channel, callback)` | Subscribe to Centrifugo channel, receive JSON messages |
| `connect(token)` | Connect to Centrifugo with JWT connection token |

---

## 12. Storage System

### GetStorage Keys (`StorageKeys` class in `storageKeys.dart`)

| Key | Purpose |
|---|---|
| `token` | JWT access token |
| `refresh` | JWT refresh token |
| `channel` | User channel ID (for Centrifugo channels) |
| `lastLoginDate` | Timestamp of last login |
| `countries` | Cached country list |
| `currencies` | Cached currency list |
| `favPairs` | User's favorite trading pairs |
| `darkMode` | Theme preference (bool) |
| `pairs` | Cached pairs data |
| `pairsHashMap` | Pairs lookup map |
| `currencyPairsHashMap` | Currency pairs lookup map |
| `coinsHashMap` | Coins lookup map |
| `selectedPair` | Currently selected trading pair |
| `loggedInOnce` | Whether user has logged in before |
| `selectedTimeFrame` | Chart timeframe selection |
| `orderedPairs` | Custom pair ordering |
| `activeMarketTabIndex` | Market tab selection state |
| `savedDepositCoins` | Saved deposit coin selections |
| `savedWithdrawalCoins` | Saved withdrawal coin selections (**BUG**: same key as `savedDepositCoins`!) |
| `lastCancelUpdate` | Timestamp of last cancel action |
| `biometricsActivated` | Whether biometrics are enabled |

### FlutterSecureStorage Keys (`SecureStorageKeys` class)

| Key | Alias | Purpose |
|---|---|---|
| `se` | email | User email (encrypted) |
| `sp` | password | User password (encrypted) |

### ⚠️ Storage Anomaly

`main.dart` reads theme from key `'lightMode'` — this key is **NOT** defined in the `StorageKeys` class. Meanwhile `StorageKeys` has a separate `darkMode` key.

---

## 13. Cryptography

### RSA Encryption (`lib/utils/cryptography/rsa_encryption.dart`)

**Class**: `RsaKeyHelper`

| Method | Purpose |
|---|---|
| `generateKeyPair()` | 2048-bit RSA key pair generation |
| `parsePublicKeyFromPem(String)` | Parse PKCS1 or PKCS8 public key from PEM |
| `parsePrivateKeyFromPem(String)` | Parse private key from PEM |
| `encrypt(String, RSAPublicKey)` | RSA encrypt plaintext |
| `decrypt(String, RSAPrivateKey)` | RSA decrypt ciphertext |
| `generateSignature(String, RSAPrivateKey)` | SHA256 signature generation |

### Encoding (`lib/utils/cryptography/encoding.dart`)

| Function | Purpose |
|---|---|
| `genearateEnc()` | RSA encrypt with public key, prepend timestamp, prefix `"ub-captcha_"` |

**Used by**: signup captcha, forgot password captcha.

---

## 14. Biometrics Service

**File**: `lib/services/localAuthService.dart`
**Package**: `local_auth 1.1.8`

### Flow

1. `isDeviceSupported()` → `canCheckBiometrics()` → `authenticate(biometricOnly: true, stickyAuth: true)`
2. On error: prompts user to settings if no biometrics enrolled
3. `hasBiometrics()` — static capability check

### Consumers

| Consumer | Usage |
|---|---|
| `LoginController.checkForBiometricLogin()` | Auto-login with stored credentials from `FlutterSecureStorage` |
| `GlobalController.checkTokenValidation()` | Re-auth with biometrics if 28-day token expired |
| `GlobalController.canContinueWithBiometrics()` | Check if biometric auth is possible |

### Credential Storage

Credentials stored in `FlutterSecureStorage`: `se` = email, `sp` = password.

---

## 15. Constants & Configuration

**File**: `lib/services/constants.dart`

| Constant | Value | Notes |
|---|---|---|
| `_urlPrefix` | `'dev-'` (DEV) / `''` (PRODUCTION) | Prefix for all URLs |
| `appVersion` | From `VERSION` dart-define | Passed at build time |
| `priceTopic` | `'trade:ticker'` | Centrifugo price ticker channel |
| `orderbookTopic` | `'trade:order-book:'` | Centrifugo orderbook channel prefix |
| `ohlcTopic` | `'trade:kline:'` | Centrifugo candle channel prefix |
| `landingPageAddress` | `'https://www.unitedbit.com'` | Marketing site |
| `cmsAddress` | `'https://content.unitedbit.com'` | CMS/news API |
| `initialPair` | `'BTC-USDT'` | Default trading pair |
| `mainUrl` | `'{prefix}app.unitedbit.com'` | Base domain |
| `centrifugoServer` | `'wss://{prefix}app.unitedbit.com:8800'` | Centrifugo WebSocket URL |
| `baseUrl` / `appUrl` | `'https://{prefix}app.unitedbit.com'` | Base app URL |
| `tradingView` | `'{appUrl}/tv/api/v1/main-route'` | TradingView route |
| `jsAPI` | `'{appUrl}/tv/api/v1/js'` | TradingView JS API |
| `urlPrefix` | `'/api/v1/'` | REST API path prefix |
| `tvUrlPrefix` | `'/tv/api/v1/js/'` | TradingView API path prefix |

### RxUpdateables Enum

```dart
enum RxUpdateables {
  Balances,
  TransactionHistory,
  UserPairBalances,
  OpenOrders,
  OrderHistory,
}
```

### Controller Tags (`controllerTags.dart`)

| Tag | Value |
|---|---|
| `login` | `'login'` |

---

## 16. Environment System

**File**: `lib/utils/environment/ubEnv.dart`

```dart
// ENV is HARDCODED to "PRODUCTION" — does NOT read from dart-define!
const ENV = "PRODUCTION";

// VERSION reads from dart-define correctly:
const VERSION = String.fromEnvironment('VERSION', defaultValue: '1.0.0');
```

### Build-Time Variables

| Variable | Passed Via | Example |
|---|---|---|
| `VERSION` | `--dart-define=VERSION=1.1.7` | App version string |
| `ENV` | `--dart-define=ENV=PRODUCTION` | **IGNORED** — hardcoded in `ubEnv.dart` |

> **BUG**: The `ENV` dart-define is passed by build scripts but never read. The value is hardcoded to `"PRODUCTION"`.

---

## 17. Feature Modules (All 23)

### 17.1 — login

| Property | Details |
|---|---|
| **Purpose** | Email/password + biometric authentication |
| **Provider** | `AuthenticationProvider` |
| **API** | `POST auth/login` |
| **Observable State** | `loginEmail`, `loginPassword`, `isLoggingIn`, `loginEmailError`, `loginPasswordError` |
| **Key Methods** | `login()`, `checkForBiometricLogin()` (reads from FlutterSecureStorage) |
| **Navigation** | Success → `loadAuthenticatedControllers()` → `/home`; 2FA required → `TwoFaPopup` |

### 17.2 — signup

| Property | Details |
|---|---|
| **Purpose** | User registration |
| **Provider** | `SignupProvider` |
| **API** | `POST auth/register` |
| **Observable State** | `isLoading`, `email`, `password`, `repeatPassword`, `emailError`, `passwordError`, `repeatPasswordError`, `step` (signupStep enum) |
| **Key Methods** | `handleSubmitClick()` → RSA encrypt password → signup |
| **Navigation** | Success → `/check-your-email` |

### 17.3 — account

| Property | Details |
|---|---|
| **Purpose** | Profile settings, biometrics toggle, logout, KYC status |
| **Provider** | `AccountProvider` |
| **API** | `POST user/send-verification-email` |
| **Observable State** | `accountData` (UserModel), `requestedForEmail`, `isRequestingForEmail`, `hasBiometrics`, `isBiometricsActivated` |
| **Key Methods** | `getUserData()`, `requestForEmailVerification()`, toggle biometrics |
| **Navigation** | To `/identity-verification`, `/identity-documents`, `/phone-verification`, `/two-factor-authentication`, `/change-password` |

### 17.4 — home

| Property | Details |
|---|---|
| **Purpose** | Dashboard with news, popular pairs, sparkline charts |
| **Provider** | `HomePageProvider` |
| **APIs** | `GET currencies/pairs-statistic?pair_currencies=BTC-USDT\|ETH-USDT\|BCH-USDT\|DASH-USDT`; `GET https://content.unitedbit.com/ubnews?_sort=date:desc&_limit=5` |
| **Observable State** | `isLoadingSparkLine`, `isSilentLoadingSparkLine`, `isLoadingNews`, `isRefreshing`, `isSilentLoadingNews`, `latestNews`, `sparkLinePairs`, `isUserVerified`, `popularPairs` |
| **Centrifugo** | Listens to `tradeController.lastPrice` for live price updates |
| **Navigation** | To `/trade`, `/market`, `/funds`, `/exchange` |

### 17.5 — trade

| Property | Details |
|---|---|
| **Purpose** | Market/limit/stop-limit orders, OHLC charts, order book, ticker |
| **Providers** | `TradeProvider` (`GET user-balance/pair-balance`, `POST order/create`), `FavoritePairsProvider` (`POST currencies/favorite`), `OHLCProvider` (`GET tv/api/v1/js/get-bars`) |
| **Observable State (20+)** | `orderBookData`, `activeChart`, `currentPairName`, `pairBalanceData`, `lastOhlcValue`, `pairs`, `isLoadingPairBalance`, `isCreatingOrder`, `mainActiveIndex`, `subActiveIndex`, `selectedPercentIndex`, `totalValue`, `amountValue`, `priceValue`, `stopValue`, `tradeFee`, `youGet`, `selectedTimeFrame`, `currentPairPrice`, `lastPrice`, `priceArray`, `showLoadingOverlay`, `amountInputLabel`, `priceInputLabel`, `totalInputLabel`, `stopPriceInputLabel` |
| **Centrifugo Subscriptions** | 1) `priceSubscription` → `trade:ticker`; 2) `ohlcSubscription` → `trade:kline:{pair}`; 3) `orderbookSubscription` → `trade:order-book:{pair}`; 4) `updateSubscription` → `authorizedCentrifugoController.updateDataSubject` for `UserPairBalances` |
| **Sub-Controller** | `OHLCChartController` — obs: `isOhlcDetailsOpen`, `isLoadingOhlc`, `chartData`, `mainState`, `secondaryState`, `isLine`, `bids`, `asks`, `isTimeFramePopupOpen`, `selectedTimeFrameButtonIndex` |

### 17.6 — exchange

| Property | Details |
|---|---|
| **Purpose** | Simple currency swap |
| **Provider** | `ExchangeProvider` |
| **APIs** | `GET user-balance/pair-balance`, `POST order/create`, `GET currencies/pairs-statistic` |
| **Observable State** | `sparkLinePairs`, `isLoadingSparkLine`, `totalValue`, `tradeFee`, `inputControllerFrom`, `inputControllerTo`, `pairBalanceData`, `isLoadingBalanceData`, `isLoadingExchangeSubmit`, `possiblePairs`, `savedCoins` |
| **Centrifugo** | Listens to `tradeController.lastPrice` |

### 17.7 — market

| Property | Details |
|---|---|
| **Purpose** | Ticker list, favorites management, search, sorting |
| **Provider** | `MarketProvider` |
| **API** | `GET crypto-payment` (transaction history) |
| **Observable State** | `isPageActive`, `searchComponentParameters`, `pairs`, `favorites`, `orderedPairs`, `sorted`, `tabCurrencies`, `activeTabIndex`, `activeTabString`, `coinSortDirection`, `lastPriceSortDirection`, `changeSortDirection` |
| **Centrifugo** | Listens to `tradeController.lastPrice` |

### 17.8 — funds

| Property | Details |
|---|---|
| **Purpose** | Container for balance, deposits, withdrawals, transaction history, auto-exchange |
| **Observable State** | `activeTabIndex`, `isHeadOpen`, `isUserVerified` |
| **Sub-pages** | balance, deposits, withdrawals, transactionHistory, autoExchange |

### 17.9 — funds/balance

| Property | Details |
|---|---|
| **Provider** | `BalanceProvider` |
| **API** | `GET user-balance/balance?sort=desc` |
| **Observable State** | `isLoading`, `isSilentLoading`, `showSmallBalances`, `showAvailableData`, `isHeadOpen`, `balancesAllData` |
| **Centrifugo** | Listens to `authorizedCentrifugoController.updateDataSubject` for `Balances` |

### 17.10 — funds/deposits

| Property | Details |
|---|---|
| **Provider** | `DepositsProvider` |
| **API** | `GET user-balance/withdraw-deposit` |
| **Observable State** | `isLoadingWithdrawAndDepositData`, `selectedNetworkIndex`, `withdrawAndDepositData`, `selectedNetwork`, `selectedCoin`, `savedCoins` |

### 17.11 — funds/withdrawals

| Property | Details |
|---|---|
| **Provider** | `WithdrawalProvider` |
| **APIs** | `POST crypto-payment/pre-withdraw`, `POST crypto-payment/withdraw`; also uses `DepositsProvider`: `GET user-balance/withdraw-deposit` |
| **Observable State** | `selectedCoin`, `selectedNetwork`, `withdrawAndDepositData`, `isScannerOpen`, `isSubmitting`, `showLoadingInsidePopup`, `isLoadingWithdrawAndDepositData`, `savedCoins`, `address`, `amount`, `fee`, `youGet`, `selectedNetworkIndex`, `selectedPercentIndex` |
| **2FA** | Uses `TwoFaPopup` for withdrawal confirmation |

### 17.12 — funds/transactionHistory

| Property | Details |
|---|---|
| **Provider** | `TransactionHistoryProvider` |
| **APIs** | `GET crypto-payment`, `POST crypto-payment/cancel` |
| **Observable State** | `transactionHistory`, `isLoading`, `showLoadingOverlay`, `isSilentLoading` |
| **Centrifugo** | Listens to `authorizedCentrifugoController.updateDataSubject` for `TransactionHistory` |

### 17.13 — funds/autoExchange

| Property | Details |
|---|---|
| **Provider** | `AutoExchangeProvider` |
| **API** | `POST user-balance/auto-exchange` |
| **Observable State** | `switchValue`, `isSubmitLoading`, `isCoinsListLoading`, `coinsList`, `balances`, `canSubmitAutoExchange`, `searchedCoin` |

### 17.14 — orders

| Property | Details |
|---|---|
| **Purpose** | Container for open orders and order history |
| **Observable State** | `isFullScreen`, `activeTabIndex`, `expanded` |
| **Sub-pages** | openOrders, orderHistory |

### 17.15 — orders/openOrders

| Property | Details |
|---|---|
| **Provider** | `OpenOrdersProvider` |
| **APIs** | `GET order/open-orders`, `POST order/cancel` |
| **Observable State** | `selectedOpenOrderFilterText`, `openOrders`, `isFullScreen`, `isSilentLoading`, `loadingIds`, `loadingData` |
| **Centrifugo** | Listens to `authorizedCentrifugoController.updateDataSubject` for `OpenOrders` |

### 17.16 — orders/orderHistory

| Property | Details |
|---|---|
| **Provider** | `OrderHistoryProvider` |
| **APIs** | `GET order/full-history`, `GET order/detail` |
| **Observable State** | `loadingData`, `filtered`, `silentLoadingData`, `orderHistory`, `loadingId`, `showFilterButton`, `selectedDateButtonIndex`, `selecteTypeButtonIndex`, `filterStartDate`, `filterEndDate`, `showCanceledOrders`, `filterPair` |
| **Centrifugo** | Listens to `authorizedCentrifugoController.updateDataSubject` for `OrderHistory` |

### 17.17 — identityInfo

| Property | Details |
|---|---|
| **Purpose** | KYC personal info (name, DOB, address, country) |
| **Provider** | `IdentityInfoProvider` |
| **API** | `POST user/set-user-profile` |
| **Observable State** | `isLoading`, `userProfile`, `firstName`, `lastName`, `gender`, `birthday`, `country`, `city`, `postalCode`, `address`, `selectedCountry`, `isSubmitting`, `hasAcceptedIdentityImage`, `hasAcceptedResidenceImage` |

### 17.18 — identityDocuments

| Property | Details |
|---|---|
| **Purpose** | KYC document upload (ID, passport, residence proof) |
| **Provider** | `IdentityDocumentsProvider` |
| **APIs** | `GET user/get-user-profile`, `UPLOAD user-profile-image/multiple-upload` |
| **Observable State** | `canChangeIdentityTypeSelect`, `canChangeResidenceTypeSelect`, `isLoading`, `canSubmit`, `activeIdentitySubTypeindex`, `activeAddressSubTypeindex`, `userProfile`, `activeTabIndex`, `identityFrontImage`, `identityBackImage`, `addressFrontImage`, `addressBackImage`, `identityFrontImageUploadFile`, `identityBackImageUploadFile`, `addressFrontImageUploadFile`, `addressBackImageUploadFile`, `hasRejectedIdentity`, `hasRejectedAddress`, `uploadPercent`, `cancelToken` |

### 17.19 — phoneVerification

| Property | Details |
|---|---|
| **Purpose** | Phone number verification via OTP |
| **Provider** | `PhoneVerificationProvider` |
| **APIs** | `POST user/sms-send`, `POST user/sms-enable` |
| **Observable State** | `is2faActivated`, `s2faCode`, `countDownValue`, `canResend`, `step`, `selectedCountry`, `isRequestingForSms`, `isRequestingtoSubmitCode`, `phoneNumber`, `password`, `verificationCode` |

### 17.20 — twoFactorAuthentication

| Property | Details |
|---|---|
| **Purpose** | Google Authenticator 2FA setup |
| **Provider** | `TwoFactorAuthenticationProviderProvider` |
| **APIs** | `GET user/google-2fa-barcode`, `POST user/google-2fa-enable`, `POST user/google-2fa-disable` |
| **Observable State** | `step`, `codeCoppied`, `isLoadingCharCode`, `isFinalSubmitLoading`, `characterCode`, `qrImageAddress`, `code`, `password`, `isEnabled` |

### 17.21 — changePassword

| Property | Details |
|---|---|
| **Provider** | `ChangePasswordProvider` |
| **API** | `POST user/change-password` |
| **Observable State** | `oldPasswordValue`, `oldPasswordError`, `newPasswordValue`, `newPasswordError`, `repeatNewPasswordValue`, `repeatNewPasswordError`, `isSubmitting`, `step` |

### 17.22 — forgot

| Property | Details |
|---|---|
| **Provider** | `ForgotProvider` |
| **API** | `POST auth/forgot-password` |
| **Observable State** | `email`, `emailError`, `isLoading` |
| **Note** | Uses RSA encryption for captcha |

### 17.23 — withdrawAddressManagement

| Property | Details |
|---|---|
| **Provider** | `AddressManagementProvider` |
| **APIs** | `GET withdraw-address`, `POST withdraw-address/delete` |
| **Observable State** | `loadingData`, `isSilentLoading`, `isRefreshing`, `withdrawAddresses` |

### 17.24 — addNewAddress

| Property | Details |
|---|---|
| **Provider** | `AddNewAddressProvider` |
| **API** | `POST withdraw-address/new` |
| **Observable State** | `selectedCoin`, `selectedNetworkIndex`, `selectedNetwork`, `newAddressLabel`, `address`, `isAddingNewAddress`, `networks` |

### 17.25 — afterSplash

| Property | Details |
|---|---|
| **Purpose** | Auth-based route decision (splash screen) |
| **Provider** | None |
| **Logic** | Checks if logged in → `/home`, otherwise → `/landing` |

### 17.26 — landing

| Property | Details |
|---|---|
| **Purpose** | Pre-auth landing page with login/signup buttons |
| **API calls** | None |

### 17.27 — checkYourEmail

| Property | Details |
|---|---|
| **Purpose** | Post-signup email verification notice |
| **API calls** | None |

### 17.28 — qrScan

| Property | Details |
|---|---|
| **Purpose** | QR code scanner for crypto addresses |
| **Package** | `ai_barcode` |

### 17.29 — webViewPage

| Property | Details |
|---|---|
| **Purpose** | Embedded WebView for T&Cs, help, etc. |
| **Package** | `webview_flutter` |

### 17.30 — separateMessagePage

| Property | Details |
|---|---|
| **Purpose** | Generic info/error/success display page |

---

## 18. Popup System

### TwoFaPopup

**Files**:
- `lib/app/popups/bindings/twoFaPopupBindings.dart`
- `lib/app/popups/controllers/twofaPopupController.dart`
- `lib/app/popups/veiws/twoFaPopupView.dart` (note: `veiws` typo in path)

**Controller** (`TwoFaController`):
- Manages 3 input types: email code, phone code, 2FA code
- Observable state: `canSubmit` validation
- Dynamic height based on required input combination

**Used by**: withdrawals (pre-withdraw), any 2FA-gated action.

---

## 19. UI Component Library (70 Widgets)

**Location**: `lib/app/common/components/`

All 70 reusable component widgets:

| Widget | Widget | Widget |
|---|---|---|
| `appbarTextTitle` | `appBarWithBottomBorder` | `autoCompleteList` |
| `bottomCard` | `CenterUBLoading` | `CoinList` |
| `controlledInput` | `countryAutocomplete` | `countrySelectBottomSheet` |
| `currencyAutocomplete` | `currencySelectButtomSheet` | `PaddedWrapper` |
| `pageContainer` | `pairsAutocomplete` | `pairSelectButtomSheet` |
| `roundedCard` | `UBBlackContainer` | `UBBlurryContainer` |
| `UBBorderlessInput` | `UBBottomSheetContainer` | `UBButton` |
| `UBCarousel` | `UBCheckbox` | `UBCircularImage` |
| `UBColoredOverlay` | `UBColumnAnimator` | `UBConnectionLost` |
| `UBCounSearchHistory` | `UBCountUp` | `UBCustomNavBarPaint` |
| `UBDarkOpacityBackgrounded` | `UBDatePicker` | `UBDDMockButton` |
| `UBDottedBorder` | `UBDropdown` | `UBFlipSideSwitcher` |
| `UBFlipSwitcher` | `UBGreyContainer` | `UBHorizontalDivider` |
| `UBInputWithTitleAndPaste` | `UBLi` | `UBLink` |
| `UBLoading` | `UBMessagePage` | `UBNetworkIcon` |
| `UBoops` | `UBPercentSelect` | `UBPlaceholder` |
| `UBPullToRefresh` | `UBRawInput` | `UBRectangle` |
| `UBRoundedButton` | `UBRoundedContainer` | `UBScaleSwitcher` |
| `UBScrollBar` | `UBScrollColumnExpandable` | `UBSection` |
| `UBSelectButton` | `UBShimmer` | `UBSimpleInput` |
| `UBSlideRightSwitcher` | `UBSlideUpSwitcher` | `UBText` |
| `UBToastOnTap` | `UBTooltip` | `UBTopTabs` |
| `UBTwoPartText` | `UBWarningRow` | `UBWrappedButtons` |
| `VCenterText` | | |

---

## 20. Custom Widgets (49 Files)

**Location**: `lib/app/common/custom/`

| Widget/Directory | Purpose |
|---|---|
| `alphabeticListView.dart` | Scrollable alphabetic index sidebar |
| `bezierChart/` | Custom Bezier curve chart library |
| `flutter_advanced_segment/` | Segmented control widget |
| `k_chart/` | Candlestick/OHLC chart (forked library) |
| `refreshHeader.dart` | Custom pull-to-refresh header |
| `rflutter_alert/` | Alert dialog library |
| `sparkline/` | Sparkline mini-chart widgets |
| `toaster/` | Toast notification overlay system |
| `ubSlideToAct/` | Slide-to-confirm action widget |

---

## 21. Theme System

**File**: `lib/theme/theme.dart`

### Dark Theme (`darkThemeData`)

| Property | Value |
|---|---|
| **Font** | OpenSans |
| **Primary Color** | Black |
| **Splash Color** | Blue 100 |
| **Hint Color** | `inputPlaceholderText` |
| **Accent Color** | Yellow Ocher |

### Light Theme (`lightThemeData`)

| Property | Value |
|---|---|
| **Font** | OpenSans |
| **Primary Color** | Blue |
| **Splash Color** | Blue 500 |
| **Highlight Color** | Crimson Red |
| **Accent Color** | Yellow Ocher |

### Font Weights

| Weight | Name |
|---|---|
| 400 | Light |
| Regular | Regular |
| 600 | SemiBold |
| 700 | Bold |
| Italic | Italic |

Colors sourced from `generated/colors.gen.dart` (generated from `assets/color/colors.xml`).

---

## 22. Localization

**English only**. No other languages configured.

| Component | File |
|---|---|
| Translation keys | `generated/locales.g.dart` (generated by `flutter_gen`) |
| Translation map | `AppTranslation.translations` |
| Service | `services/localizationService.dart` |
| Fallback locale | `en_US` |
| Setup | GetX `translationsKeys` in `GetMaterialApp` |

---

## 23. Mixins & Utilities

### Toaster Mixin (`lib/utils/mixins/toast.dart`)

| Method | Purpose |
|---|---|
| `toast()` | Generic toast |
| `toastSuccess()` | Green success toast |
| `toastError()` | Red error toast |
| `toastWarning()` | Warning toast |
| `toastInfo()` | Info toast |
| `toastAction()` | Actionable toast |
| `toastDioError()` | Format and display Dio error |
| `toastAuthorizedEvent()` | Display authorized Centrifugo event toast |

### Popups Mixin (`lib/utils/mixins/popups.dart`)

| Method | Purpose |
|---|---|
| `openConfirmation()` | Confirmation dialog |
| `openUpdatePopup()` | App update dialog |
| Other popup helpers | Various info/action popups |

### Formatter Mixin (`lib/utils/mixins/formatters.dart`)

| Method | Purpose |
|---|---|
| `currencyFormatter()` | Format number as currency string |

### FilterPopups Mixin (`lib/utils/mixins/filterPopups.dart`)

| Method | Purpose |
|---|---|
| `openOrdersFilterSelect()` | Open orders filter popup |
| `orderHistoryFilterButton()` | Order history filter |
| `openOpenOrdersFilter()` | Open orders filter selector |

### CommonConsts Mixin (`lib/utils/mixins/commonConsts.dart`)

| Category | Values |
|---|---|
| **Vertical spacing** | `vspace2` through `vspace48` |
| **Horizontal spacing** | `hspace2` through `hspace24` |
| **Text styles** | Preset text style constants |
| **Borders** | Border decoration presets |
| **Corners** | Rounded corner presets |

### String Extensions (`lib/utils/extentions/basic.dart`)

| Extension | Purpose |
|---|---|
| `String.currencyFormat()` | Format string as currency |
| `String.removeComma()` | Strip commas |
| `String.toDouble()` | Parse to double |
| `Double.toFixedWithoutRounding()` | Fixed decimal without rounding |

### Utility Files

| File | Purpose |
|---|---|
| `basicMath.dart` | Math helper functions |
| `commonUtils.dart` | Shared utilities (`launchURL`, etc.) |
| `computes.dart` | JSON parsing via `compute()` (background isolate) |
| `debounce.dart` | Debounce helper class |
| `throttle.dart` | Throttle helper class |
| `emailValidator.dart` | Email validation regex |
| `passwordValidator.dart` | Password validation rules |
| `numUtil.dart` | Numeric utilities |
| `marketUtils.dart` | Market/trading helpers |
| `pairAndCurrencyUtils.dart` | `PairAndCurrencyUtils` singleton (hash maps for pair/currency lookup) |
| `inputCurrencyFormatter.dart` | `TextInputFormatter` for currency input fields |
| `logger.dart` | Logger instance configuration |
| `UBDeviceType.dart` | Device type detection (`DeviceTypes.PHONE` / `DeviceTypes.TABLET`) |

### Auth Middleware (`lib/utils/middleWares/authMiddleware.dart`)

**⚠️ PLACEHOLDER**: `GetMiddleware` with `isAuthenticated` always `false` — always redirects to `/login`.

---

## 24. Code Generation

**Tool**: `flutter_gen 4.1.2`

### Generated Files (`lib/generated/`)

| File | Source | Content |
|---|---|---|
| `assets.gen.dart` | Asset files (rive, images) | Asset path constants |
| `colors.gen.dart` | `assets/color/colors.xml` | Color constants |
| `fonts.gen.dart` | Font declarations | Font family constants |
| `locales.g.dart` | Translation strings | English translation map |

### Regeneration Command

```bash
flutter pub run build_runner build
```

Configuration in `pubspec.yaml` under `flutter_gen:` section.

---

## 25. Complete API Endpoint Reference

### Authentication

| Method | Endpoint | Purpose | Used By |
|---|---|---|---|
| `POST` | `auth/login` | Login (email, password, recaptcha) | `login` module |
| `POST` | `auth/register` | Register (email, password, recaptcha) | `signup` module |
| `POST` | `auth/refresh` | Refresh JWT token | Interceptor (automatic on 403) |
| `POST` | `auth/forgot-password` | Password reset (email, recaptcha) | `forgot` module |

### User Profile

| Method | Endpoint | Purpose | Used By |
|---|---|---|---|
| `GET` | `user/user-data` | Get user profile | `account` module, `commonDataProvider` |
| `POST` | `user/send-verification-email` | Request email verification | `account` module |
| `POST` | `user/set-user-profile` | Update KYC personal info | `identityInfo` module |
| `GET` | `user/get-user-profile` | Get KYC profile with images | `identityDocuments` module |
| `UPLOAD` | `user-profile-image/multiple-upload` | Upload KYC documents | `identityDocuments` module |
| `POST` | `user/change-password` | Change password | `changePassword` module |

### Phone Verification

| Method | Endpoint | Purpose | Used By |
|---|---|---|---|
| `POST` | `user/sms-send` | Request SMS OTP | `phoneVerification` module |
| `POST` | `user/sms-enable` | Verify phone with OTP | `phoneVerification` module |

### Two-Factor Authentication

| Method | Endpoint | Purpose | Used By |
|---|---|---|---|
| `GET` | `user/google-2fa-barcode` | Get 2FA setup QR code | `twoFactorAuthentication` module |
| `POST` | `user/google-2fa-enable` | Enable Google 2FA | `twoFactorAuthentication` module |
| `POST` | `user/google-2fa-disable` | Disable Google 2FA | `twoFactorAuthentication` module |

### Currencies & Pairs

| Method | Endpoint | Purpose | Used By |
|---|---|---|---|
| `GET` | `currencies` | All currencies (with images, networks) | `commonDataProvider` |
| `GET` | `currencies/pairs-list` | All trading pairs | `commonDataProvider` |
| `GET` | `currencies/pairs-statistic?pair_currencies=...` | Pair price statistics | `home`, `exchange` modules |
| `GET` | `currencies/favorite-pairs` | User's favorite pairs | `commonDataProvider` |
| `POST` | `currencies/favorite` | Toggle pair favorite (add/remove) | `trade` module |

### Trading

| Method | Endpoint | Purpose | Used By |
|---|---|---|---|
| `GET` | `user-balance/pair-balance` | Balance for specific trading pair | `trade`, `exchange` modules |
| `GET` | `user-balance/balance?sort=desc` | All user balances | `balance` module |
| `POST` | `user-balance/auto-exchange` | Auto-exchange toggle | `autoExchange` module |
| `POST` | `order/create` | Create market/limit/stop order | `trade`, `exchange` modules |
| `GET` | `order/open-orders` | List active orders | `openOrders` module |
| `POST` | `order/cancel` | Cancel an order | `openOrders` module |
| `GET` | `order/full-history` | Order history (with filters) | `orderHistory` module |
| `GET` | `order/detail` | Single order details | `orderHistory` module |

### Payments & Transactions

| Method | Endpoint | Purpose | Used By |
|---|---|---|---|
| `GET` | `user-balance/withdraw-deposit` | Deposit/withdraw data for currency | `deposits`, `withdrawals` modules |
| `POST` | `crypto-payment/pre-withdraw` | Initiate withdrawal (needs 2FA) | `withdrawals` module |
| `POST` | `crypto-payment/withdraw` | Confirm withdrawal | `withdrawals` module |
| `GET` | `crypto-payment` | Transaction history | `transactionHistory`, `market` modules |
| `POST` | `crypto-payment/cancel` | Cancel pending withdrawal | `transactionHistory` module |

### Withdrawal Addresses

| Method | Endpoint | Purpose | Used By |
|---|---|---|---|
| `GET` | `withdraw-address` | List user's withdrawal addresses | `withdrawAddressManagement` module |
| `POST` | `withdraw-address/delete` | Delete withdrawal address | `withdrawAddressManagement` module |
| `POST` | `withdraw-address/new` | Add new withdrawal address | `addNewAddress` module |

### Data & System

| Method | Endpoint | Purpose | Used By |
|---|---|---|---|
| `GET` | `main-data/country-list` | List of countries | `commonDataProvider` |
| `GET` | `main-data/version?platform={}&current_version={}` | App version check | `commonDataProvider` |

### TradingView

| Method | Endpoint | Purpose | Used By |
|---|---|---|---|
| `GET` | `tv/api/v1/js/get-bars?symbol=&resolution=&from=&to=` | OHLC candle data | `trade` module (OHLCProvider) |

### External APIs

| Method | URL | Purpose | Used By |
|---|---|---|---|
| `GET` | `https://content.unitedbit.com/ubnews?_sort=date:desc&_limit=5` | News articles | `home` module |

---

## 26. Centrifugo Channel Reference

### Public Channels (UnAuthorizedCentrifugoController)

| Topic | Content | QoS | Consumers |
|---|---|---|---|
| `main/trade/ticker` | Live price ticker for all pairs | Default | `TradeController`, `HomeController`, `MarketController`, `ExchangeController` |
| `main/trade/order-book/{pair}` | Live order book for specific pair | Default | `TradeController` |
| `main/trade/kline/{timeframe}/{pair}` | Live OHLC candle updates | Default | `TradeController` |

### Private Channels (AuthorizedCentrifugoController)

| Topic | Content | QoS | Consumers |
|---|---|---|---|
| `main/trade/user/{channel}/open-orders/` | User's order status events | exactlyOnce | `OpenOrdersController`, `OrderHistoryController`, `BalanceController` |
| `main/trade/user/{channel}/crypto-payments/` | User's payment events | exactlyOnce | `TransactionHistoryController`, `BalanceController` |

### Update Signal Flow

AuthorizedCentrifugoController emits `RxUpdateables` via `updateDataSubject`:

```
Centrifugo event → AuthorizedCentrifugoController → updateDataSubject.add([RxUpdateables.X])
  → BalanceController (Balances)
  → TransactionHistoryController (TransactionHistory)
  → TradeController (UserPairBalances)
  → OpenOrdersController (OpenOrders)
  → OrderHistoryController (OrderHistory)
```

---

## 27. Common Data Provider

**File**: `lib/app/global/providers/commonDataProvider.dart`

| Method | Endpoint | Returns |
|---|---|---|
| `getCountries()` | `GET main-data/country-list` | Country list |
| `getCurrencies()` | `GET currencies` | Currencies with images and networks |
| `getUserData()` | `GET user/user-data` | User profile |
| `getFavoritePairs()` | `GET currencies/favorite-pairs` | Favorite trading pairs |
| `getPairsList()` | `GET currencies/pairs-list` | All trading pairs |
| `getVersion()` | `GET main-data/version?platform={}&current_version={}` | App version check |

---

## 28. Build Scripts

| Script | Purpose | Output | Env |
|---|---|---|---|
| `scripts/buildAPK.sh` | Production Android APK | Split-per-ABI APKs | `ENV=PRODUCTION` |
| `scripts/buildBundle.sh` | Production Android App Bundle | AAB file | `ENV=PRODUCTION` |
| `scripts/buildDevAPK.sh` | Dev Android APK + send to Telegram | Split-per-ABI APKs | `ENV=DEV` |
| `scripts/buildWeb.sh` | Production web build | `build/web/` | `ENV=PRODUCTION` |
| `scripts/buildWeb-dev.sh` | Dev web build | `build/web/` | `ENV=DEV` |
| `scripts/runApp.sh` | Local dev run (device) | Live debug | — |
| `scripts/runWeb.sh` | Local dev run (Chrome) | Port 8000 | `ENV=DEV` |
| `scripts/tag.sh` | Git commit + tag with pubspec version | Git tag | — |

### Common Build Steps

All build scripts follow this pattern:
```bash
flutter pub get
flutter clean
# (splash config if APK)
flutter build <target> \
  --dart-define=VERSION=$version \
  --dart-define=ENV=PRODUCTION|DEV
```

### Web Build Specifics

- Copies `buildConfigs/web/manifest.json` and `index.html` to `web/` before build
- Dev uses `buildConfigs/web/dev/manifest.json`
- Renderer: `--web-renderer canvaskit`

---

## 29. Docker Configuration

| Dockerfile | Base | Flutter | Build | Output |
|---|---|---|---|---|
| `docker/Dockerfiledev` | `debian:bullseye-slim` | 2.10.5 | `scripts/buildWeb-dev.sh` | `nginx:alpine` serving `build/web/` |
| `docker/Dockerfileprod` | `debian:bullseye-slim` | 2.10.5 | `scripts/buildWeb.sh` | `nginx:alpine` serving `build/web/` |
| `docker/Dockerfileapkprod` | `debian:bullseye-slim` | 2.10.5 | `scripts/buildAPK.sh` | APK files (includes Android CLI tools + SDK license acceptance) |

---

## 30. GitLab CI/CD Pipeline

**File**: `.gitlab-ci.yml`

### 6 Stages

| Stage | Trigger | Actions |
|---|---|---|
| `dev-build` | Push to `develop` branch | Build Docker dev image |
| `dev-deploy` | After dev-build | Deploy to `/home/gitlab-runner/dev-m/` |
| `dev-notification` | After dev-deploy | Telegram notification (success/failure) |
| `prod-build` | Push to `master` branch | Build Docker prod image |
| `prod-deploy` | After prod-build | Deploy to `/home/gitlab-runner/prod-m/` |
| `prod-notification` | After prod-deploy | Telegram notification (success/failure) |

---

## 31. Platform Configuration

### Android

| Property | Value |
|---|---|
| `applicationId` | `com.unitedbit.app` |
| `compileSdkVersion` | 33 |
| `minSdkVersion` | 21 |
| `targetSdkVersion` | 33 |
| `Kotlin version` | 1.6.10 |
| `Gradle plugin` | 4.1.2 |
| `Repository` | `mavenCentral()` (jcenter removed) |
| **Permissions** | `INTERNET`, `QUERY_ALL_PACKAGES`, `CAMERA`, `USE_FINGERPRINT`, `USE_BIOMETRIC` |
| **Signing** | Release signing from `key.properties` |

### iOS

| Property | Value |
|---|---|
| **Support** | iOS directory exists |
| **Embedding** | Flutter embedding v2 |
| **Configured plugins** | `local_auth`, `permission_handler` |

### Web

| Property | Value |
|---|---|
| **index.html** | Safari WebGL2 compatibility fix, SPA routing |
| **URL strategy** | PathUrlStrategy (no hash URLs) |
| **Renderer** | CanvasKit |

---

## 32. Testing

### Framework

- `flutter_test` (built-in)
- `mockito 5.0.16`
- Run: `flutter test` (requires Dart 2.x SDK)

### All 15 Test Files

| # | File | Content |
|---|---|---|
| 1 | `test/main_test.dart` | `Get.defaultDialog`, `Get.dialog` widget tests |
| 2 | `test/rx_workers_test.dart` | GetX reactive workers (`once`, `ever`, `debounce`, `interval`, `bindStream`, `trigger`) |
| 3 | `test/app/global/controller/globalController_test.dart` | Placeholder (empty) |
| 4 | `test/app/modules/account/providers/accountProvider_test.dart` | Mock-based provider tests |
| 5 | `test/app/modules/account/providers/accountProvider_test.mocks.dart` | Mockito generated mocks |
| 6 | `test/app/modules/account/views/account_view_test.dart` | Placeholder (empty) |
| 7 | `test/app/modules/trade/controllers/trade_controller_test.dart` | Basic placeholder |
| 8 | `test/utils/benchmark_test.dart` | GetX vs ValueNotifier vs Streams performance benchmark |
| 9 | `test/utils/commonUtils_test.dart` | Placeholder (empty) |
| 10 | `test/utils/dynamic_extentions_test.dart` | Placeholder (empty) |
| 11 | `test/utils/num_extentions_test.dart` | `num.isLowerThan`, `isGreaterThan`, `isEqual` |
| 12 | `test/utils/widget_extentions_test.dart` | Widget `padding`/`margin` extensions |
| 13 | `test/utils/wrapper.dart` | Test wrapper utility with `GetMaterialApp` |
| 14 | `test/utils/cryptography/rsa_encryption_test.dart` | RSA encrypt/decrypt round-trip |
| 15 | `test/utils/extentions/basic_test.dart` | `String.currencyFormat` tests |

### Mock File

`dev-tools/getToken.http` — REST client test request file.

---

## 33. Known Issues & Bugs

> **Last deep audit:** 2025-07-25 (line-by-line source audit)

### 33.1 — Critical Issues

| # | Severity | Issue | Location | Details |
|---|---|---|---|---|
| 1 | **CRITICAL** | Pre-null-safety (Dart 2.11) — entire codebase needs null safety migration | `pubspec.yaml` SDK constraint | Blocks all modern Dart/Flutter features |
| 2 | **CRITICAL** | NPE in `_shouldRefreshToken()` — `err.response` can be null | `lib/services/apiService.dart:177` | Any network error without HTTP response (timeout, DNS failure) crashes here: `err.response.statusCode` on null |
| 3 | **CRITICAL** | Race condition in token refresh interceptor | `lib/services/apiService.dart:82-111` | Concurrent 403 responses both enter the refresh handler; first consumes refresh token, second may fail or corrupt state. Dio locks only queue new requests, not concurrent error handlers |
| 4 | **CRITICAL** | Hardcoded reCAPTCHA API keys in commented-out code | `lib/app/modules/login/views/login_view.dart:145-146` | `apiKey` and `apiSecret` visible in source/version control |
| 5 | **CRITICAL** | Storage key collision: `savedWithdrawalCoins` = `'savedDepositCoins'` | `lib/services/storageKeys.dart:26` | Deposit and withdrawal coin history overwrite each other |

### 33.2 — High Severity Issues

| # | Severity | Issue | Location | Details |
|---|---|---|---|---|
| 6 | **HIGH** | JWT token stored in unencrypted `GetStorage` | `apiService.dart`, `login_controller.dart` | Token readable by any app with file access on rooted devices; should use `FlutterSecureStorage` |
| 7 | **HIGH** | `ENV` hardcoded to `"PRODUCTION"` — ignores `--dart-define=ENV` | `lib/utils/environment/ubEnv.dart:4` | `const ENV = "PRODUCTION";` — build scripts pass ENV but it's never read |
| 8 | **HIGH** | `_timer.cancel()` without null check in `onClose()` | `lib/app/global/controller/authorizedCentrifugoController.dart:79` | `_timer` can be null if `onInit()` fails — crash on dispose |
| 9 | **HIGH** | Platform header always `'ubandroid'` — copy-paste bug | `lib/services/apiService.dart:26` | `(GetPlatform.isAndroid ? 'ubandroid' : 'ubandroid')` — both branches identical |
| 10 | **HIGH** | `async void handleLoggedOut()` — fire-and-forget | `lib/app/global/controller/globalController.dart` | Exceptions in async void are uncatchable |
| 11 | **HIGH** | Uncancelled stream subscriptions in `OHLCChartController` | `lib/app/modules/trade/controllers/ohlcChart_controller.dart` | 3 `.listen()` calls in `onInit()`, empty `onClose()` → memory leak |
| 12 | **HIGH** | Memory leak: `StreamController` never closed in `alphabeticListView` | `lib/app/common/custom/alphabeticListView.dart:70` | `_pixelUpdates` StreamController and `_callback` listener never disposed |
| 13 | **HIGH** | Hardcoded `UserAgentInfo` in order creation | `tradeProvider.dart`, `exchangeProvider.dart` | Always sends `{"browser":"Chrome","device":"web","os":"Win32"}` regardless of platform |
| 14 | **HIGH** | Raw RSA encryption without padding | `lib/utils/cryptography/rsa_encryption.dart:243-250` | `RSAEngine()` used without OAEP/PKCS1 padding — vulnerable to padding oracle attacks |
| 15 | **HIGH** | Race condition in autoExchange stream listener | `lib/app/modules/funds/pages/autoExchange/controllers/auto_exchange_controller.dart:38-42` | Async callback modifies shared state without synchronization |

### 33.3 — Medium Severity Issues

| # | Severity | Issue | Location |
|---|---|---|---|
| 16 | **MEDIUM** | `'lightMode'` key in `main.dart` not in `StorageKeys` class; conflicts with `darkMode` key | `lib/main.dart` |
| 17 | **MEDIUM** | `connectivity` package deprecated → use `connectivity_plus` | `pubspec.yaml` |
| 18 | **MEDIUM** | `flutter_appavailability` unmaintained, no null-safe version | `pubspec.yaml` |
| 19 | **MEDIUM** | `async void onInit()` in multiple controllers | Various controllers |
| 20 | **MEDIUM** | Missing null check on `networksConfigsAndAddresses.length` | `withdrawals_controller.dart` |
| 21 | **MEDIUM** | No `onError` handler on Centrifugo stream subscriptions | `trade_controller.dart` |
| 22 | **MEDIUM** | Mixed `print()` / `debugPrint()` / `log.e()` logging | Various files |
| 23 | **MEDIUM** | No retry limit on token refresh — infinite loop risk | `apiService.dart:82-111` |
| 24 | **MEDIUM** | Zero error handling in all 26 provider files | All providers return raw response; no try/catch |
| 25 | **MEDIUM** | `Future.wait` without `eagerError: false` in balance controller | `balance_controller.dart` |
| 26 | **MEDIUM** | Route typo: `/depost-details` (missing 'i') | `app_routes.dart`, `app_pages.dart` |
| 27 | **MEDIUM** | No certificate pinning for HTTPS/WSS | All network communication |

### 33.4 — Low Severity Issues

| # | Severity | Issue | Location |
|---|---|---|---|
| 28 | **LOW** | Hardcoded Centrifugo/API URLs in `constants.dart` (not environment-driven) | `lib/services/constants.dart` |
| 29 | **LOW** | Docker pinned to Flutter 2.10.5 | `Dockerfile*` |
| 30 | **LOW** | `AuthMiddleware` placeholder (`isAuthenticated` always `false`) | `lib/utils/middleWares/authMiddleware.dart` |
| 31 | **LOW** | Empty test placeholders (64% of test files) | `test/` directory |
| 32 | **LOW** | `genarateEnc()` function name misspelled | `lib/utils/cryptography/encoding.dart` |
| 33 | **LOW** | `veiws/` directory typo in popups | `lib/app/popups/veiws/` |
| 34 | **LOW** | `TwoFactorAuthenticationProviderProvider` — double "Provider" | Provider file |
| 35 | **LOW** | Unused flags: `priceTopicInitialized`, `ohlcTopicInitialized`, `orderBookTopicInitialized` | `trade_controller.dart` |
| 36 | **LOW** | Request retrier has no timeout — hangs indefinitely | `lib/services/interceptors/request_retrier.dart` |
| 37 | **LOW** | Duplicate provider methods across files | `getAppVerion()` in CommonDataProvider + AfterSplashProvider |
| 38 | **LOW** | Commented-out code blocks (10+ locations) | Various files |

---

## 34. Deep Audit Findings Summary

> Added by deep line-by-line audit (2025-07-25). See full report in session artifacts.

### 34.1 — Security Assessment

| Area | Grade | Key Finding |
|---|---|---|
| **Token Storage** | D | JWT in plain `GetStorage`; email/password correctly in `FlutterSecureStorage` |
| **Cryptography** | C | RSA 2048-bit (OK), but no padding scheme, hardcoded public key |
| **Network** | C+ | WSS for Centrifugo (good), HTTPS for API (good), but no cert pinning |
| **Credentials** | D | Hardcoded reCAPTCHA keys in comments; UserAgentInfo hardcoded |
| **Input Validation** | B | Email/password validators exist; currency formatters proper |

### 34.2 — Undocumented Items Found

| Item | Details |
|---|---|
| OHLC topic includes timeframe | Actual: `main/trade/kline/{timeframe}/{pair}`, not `main/trade/kline/{pair}` |
| Hardcoded popular pairs | 9 pairs in `home_controller.dart:39-49`: BTC-USDT, ETH-USDT, BCH-USDT, DASH-USDT, DOGE-USDT, MKR-USDT, LTC-USDT, ETH-BTC, TRX-USDT |
| Duplicate providers | `getAppVerion()` in CommonDataProvider AND AfterSplashProvider; `getTransactionHistory()` in TransactionHistoryProvider AND MarketProvider |
| TimeoutRetryInterceptor commented out | Listed in interceptor chain but actually disabled (apiService.dart:113-125) |
| Total API endpoints | 40 (37 internal + 1 external CMS + 1 TradingView + 1 token refresh) |
| Provider count | 26 files (25 module providers + 1 commonDataProvider) |

### 34.3 — Test Coverage Assessment

| Category | Coverage | Notes |
|---|---|---|
| API Service / Interceptors | **0%** | No tests for token refresh, retry logic, error handling |
| Centrifugo Controllers | **0%** | No tests for message parsing, reconnection, topic management |
| GlobalController | **0%** | Empty placeholder test |
| Trade Module | **~1%** | Basic placeholder only |
| UI Components (119 widgets) | **0%** | Zero tests for UBButton, UBInput, OrderBook, CoinList, etc. |
| Cryptography | **~50%** | RSA round-trip tested; encoding.dart untested |
| Providers (26 files) | **~5%** | Only accountProvider has mock test |
| **Overall Estimated** | **~5%** | Critical gap for a financial trading application |

### 34.4 — Interceptor Chain (Corrected)

| Order | Interceptor | Status | Trigger |
|---|---|---|---|
| 1 | Connection Retry (`RetryOnConnectionChangeInterceptor`) | ✅ Active | `SocketException` |
| 2 | Token Refresh (inline `InterceptorsWrapper`) | ✅ Active | HTTP 403 |
| 3 | Auth Header + Logging (inline `InterceptorsWrapper`) | ✅ Active | Every request |
| 4 | Timeout Retry (`TimeoutRetryInterceptor`) | ❌ **Commented out** | — |
| 5 | Pretty Dio Logger | ❌ **Commented out** | — |

### 34.5 — Centrifugo Channel Reference (Corrected)

| Topic | Exact Pattern | QoS | Controller |
|---|---|---|---|
| Price Ticker | `trade:ticker` | exactlyOnce (2) | UnAuthorized → TradeController |
| Order Book | `trade:order-book:{pair}` | exactlyOnce (2) | UnAuthorized → TradeController |
| **OHLC/Klines** | **`trade:kline:{timeframe}:{pair}`** | exactlyOnce (2) | UnAuthorized → TradeController |
| Open Orders | `user:{channel}:open-orders` | exactlyOnce (2) | AuthorizedCentrifugoController |
| Crypto Payments | `user:{channel}:crypto-payments` | exactlyOnce (2) | AuthorizedCentrifugoController |

> **Correction**: OHLC channel includes `{timeframe}` segment between `kline:` and `{pair}`.

---

## 35. Upgrade Roadmap

### Phase 0 — Critical Bug Fixes (NEW — from deep audit)

- [ ] Fix NPE in `_shouldRefreshToken()` — add null check on `err.response` (`apiService.dart:177`)
- [ ] Remove hardcoded reCAPTCHA keys from `login_view.dart:145-146`
- [ ] Fix `savedWithdrawalCoins` key collision (`storageKeys.dart:26` → change to `'savedWithdrawalCoins'`)
- [ ] Add null check for `_timer.cancel()` in `authorizedCentrifugoController.dart:79`
- [ ] Fix platform header copy-paste bug (`apiService.dart:26`)
- [ ] Add token refresh race condition guard (mutex/flag)
- [ ] Move JWT token to `FlutterSecureStorage`
- [ ] Fix `ENV` to read from `String.fromEnvironment('ENV', defaultValue: 'DEV')`

### Phase 1 — Infrastructure (COMPLETED ✅)

- [x] AGENTS.md comprehensive rewrite
- [x] `.env.example` created
- [x] Docker: pin Flutter 2.10.5
- [x] Android: `jcenter()` → `mavenCentral()`, `lintOptions` → `lint`
- [x] Android: `compileSdkVersion` 31→33, `targetSdkVersion` 30→33

### Phase 2 — Null Safety Migration (HIGH RISK)

- Dart SDK `>=2.11.0 <3.0.0` → `>=2.17.0 <3.0.0` (sound null safety)
- Run `dart migrate` for automated analysis
- Bottom-up migration order: models → services → providers → controllers → views
- Update all packages to null-safe versions
- Each module tested independently

### Phase 3 — Dart 3.x + Flutter 3.x

- SDK constraint `>=3.0.0`
- Material 3, Impeller renderer
- GetX 5.x, Dio 5.x
- Docker update to Flutter 3.x
- Android: Kotlin 1.9, Gradle 8.x
