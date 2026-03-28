# Flutter App — ub-app-main

## Stack
- **Dart SDK >=2.11.0 <3.0.0** (pre-null-safety — CRITICAL upgrade needed)
- **Flutter 2.x** (last compatible with Dart 2.x)
- **GetX 4.3.4** — routing, state management, DI (bindings, controllers, obs)
- **Dio 4.0.1** — HTTP client with interceptors (token refresh on 403, retry on connection loss)
- **MQTT Client 9.6.1** — dual controllers (authorized + unauthorized)
- **RxDart 0.27.2** — reactive streams (BehaviorSubject)
- **Provider 6.0.1** — secondary state management
- **GetStorage 2.0.3** — local KV persistence
- **FlutterSecureStorage 4.2.1** — encrypted credential storage
- **PointyCastle 3.3.5 + Encrypt 5.0.1** — RSA/AES cryptography
- **Decimal 1.3.0** — precise monetary arithmetic

## Architecture

### Module Pattern (GetX MVC)
Each feature: `lib/app/modules/<feature>/`
- `bindings/<feature>_binding.dart` — GetX DI registration
- `controllers/<feature>_controller.dart` — business logic + reactive state (.obs)
- `views/<feature>_view.dart` — UI (extends GetView<Controller>)
- `providers/<feature>_provider.dart` — API calls (uses ApiService singleton)
- `models/` — data models (optional)

### 23 Feature Modules
| Module | Purpose |
|--------|---------|
| `account` | Profile, settings, biometrics |
| `home` | Dashboard, news, popular pairs |
| `landing` | Pre-auth landing page |
| `login` | Email/password + biometric auth |
| `signup` | Registration flow |
| `trade` | Market/limit orders, OHLC charts, order book (MQTT) |
| `exchange` | Currency exchange/swap |
| `market` | Ticker, favorites, price changes |
| `funds` | Balance, deposits, withdrawals, transaction history, auto-exchange |
| `orders` | Open orders (MQTT live), order history |
| `identityInfo` | KYC info (DoB, address) |
| `identityDocuments` | KYC document upload (ID, passport) |
| `phoneVerification` | Phone OTP verification |
| `twoFactorAuthentication` | 2FA setup (Google Authenticator) |
| `changePassword` | Password change |
| `withdrawAddressManagement` | Crypto withdrawal addresses |
| `addNewAddress` | Add blockchain address |
| `forgot` | Password recovery |
| `checkYourEmail` | Email verification status |
| `qrScan` | QR scanner for addresses |
| `webViewPage` | Embedded web content (T&Cs, help) |
| `separateMessagePage` | Generic info/error display |
| `afterSplash` | Auth-based route decision |

### 42+ Named Routes
Initial route: `/after-splash` (decides `/landing` or `/home` based on auth)
- Pre-auth: `/landing`, `/login`, `/signup`, `/forgot`, `/check-your-email`, `/two-factor-authentication`
- Main: `/home`, `/market`, `/trade`, `/account`, `/funds`, `/orders`, `/exchange`
- Sub: `/balance`, `/deposits`, `/withdrawals`, `/transaction-history`, `/open-orders`, `/order-history`
- KYC: `/identity-verification`, `/identity-documents`, `/phone-verification`
- Settings: `/change-password`, `/withdraw-address-management`, `/add-new-address`
- Utility: `/qr-scan`, `/web-view-page`, `/separate-message-page`, `/exchange-search`, `/auto-exchange`, `/edit-favorites`

## Key Services

### ApiService (`lib/services/apiService.dart`)
Singleton HTTP client wrapping Dio:
- Base URL: `https://{prefix}app.unitedbit.com/api/v1/` (prefix = `dev-` in DEV mode)
- Token refresh on 403 via interceptor chain
- Connection retry on SocketException
- Platform header: `ubandroid-v{version}`
- Timeouts: 10s connect, 10s receive
- Methods: `get()`, `post()`, `rawGet()`, `upload()`

### MQTT Controllers
- **AuthorizedMqttController**: User-private topics (`main/trade/user/{channel}/open-orders/`, `crypto-payments/`)
- **UnAuthorizedMqttController**: Public topics (`main/trade/ticker`, `order-book/{pair}`, `kline/{pair}`)
- Broker: `wss://app.unitedbit.com:8443`
- Auth: JWT token as username, UUID as password

### GlobalController (Permanent Singleton)
- Connectivity monitoring, device type detection
- Theme management (dark/light persisted via GetStorage)
- Pair/currency hash maps
- Auth state (`loggedIn` observable)

### BiometricsService (`lib/services/localAuthService.dart`)
- Fingerprint/face authentication via `local_auth`
- Stores encrypted credentials in FlutterSecureStorage

## Storage
### GetStorage (plain KV)
`token`, `refresh`, `channel`, `darkMode`, `favPairs`, `orderedPairs`, `biometricsActivated`, `currencies`, `countries`, `pairs`, `lastLoginDate`

### FlutterSecureStorage (encrypted)
`se` (email), `sp` (password)

## Build & Dev Commands
```bash
flutter pub get              # Install dependencies (requires Dart 2.x)
flutter run                  # Debug on device
flutter test                 # Run tests

# Build scripts (Linux/Docker)
./buildDevAPK.sh             # Android DEV APK
./buildAPK.sh                # Android PROD APK (--split-per-abi)
./buildWeb-dev.sh            # Web DEV (canvaskit renderer)
./buildWeb.sh                # Web PROD (canvaskit renderer)
./runApp.sh                  # Local flutter run
./runWeb.sh                  # Local flutter run -d chrome
```

## Environment
- `ENV` constant in `lib/utils/environment/ubEnv.dart`: `"PRODUCTION"` or `"DEV"`
- Version passed via `--dart-define=VERSION=$version --dart-define=ENV=PRODUCTION`

## Android Config
- `applicationId`: `com.unitedbit.app`
- `compileSdkVersion`: 33, `minSdkVersion`: 21, `targetSdkVersion`: 33
- Kotlin: 1.6.10, Gradle plugin: 4.1.2
- `mavenCentral()` (jcenter removed)

## Code Generation
- `flutter_gen` generates assets, colors, fonts, locales into `lib/generated/`
- Run: `flutter pub run build_runner build`
- Config in pubspec.yaml under `flutter_gen:`

## Testing
- **Framework**: flutter_test + mockito 5.0.16
- **13 test files**: widget smoke tests, reactive worker tests, util tests, crypto tests, controller tests
- **Run**: `flutter test` (requires Dart 2.x SDK)

## UI Component Library
100+ custom widgets in `lib/app/common/components/`:
- Form: UBSimpleInput, UBRawInput, UBInputWithTitleAndPaste
- Buttons: UBButton, UBRoundedButton
- Text: UBText, UBTwoPartText
- Loading: UBLoading, UBShimmer
- Navigation: UBTopTabs, bottomCard, sideMenu
- Data: CoinList, UBCountUp, percent_indicator

## Critical Issues

### ⚠️ Pre-Null-Safety
The entire codebase is pre-null-safety (Dart 2.11). Every `.dart` file needs `?` nullable annotations, `late` keywords, and null checks. This is the largest migration effort in the monorepo.

### Other Issues
- `ENV` hardcoded to `"PRODUCTION"` in `ubEnv.dart` (bypasses dart-define)
- `connectivity` package deprecated → use `connectivity_plus`
- `flutter_appavailability` — unmaintained
- Hardcoded MQTT/API URLs in `constants.dart`
- Docker pinned to Flutter 2.10.5 — cannot use Flutter 3.x until null safety migration is complete
- Local Flutter 3.x SDK cannot build/test this project (requires Dart <3.0)

## Upgrade Roadmap

### Phase 1 — Infrastructure (safe, no Dart changes)
- [x] AGENTS.md comprehensive rewrite
- [x] .env.example created
- [x] Docker: pin Flutter 2.10.5 (fix broken builds)
- [x] Android: jcenter() → mavenCentral(), lintOptions → lint
- [x] Android: compileSdkVersion 31→33, targetSdkVersion 30→33

### Phase 2 — Null Safety Migration (HIGH RISK, HIGH EFFORT)
- Dart SDK `>=2.11.0 <3.0.0` → `>=2.17.0 <3.0.0` (sound null safety)
- Run `dart migrate` for automated analysis
- Bottom-up: models → services → providers → controllers → views
- Update all packages to null-safe versions
- Each module tested independently

### Phase 3 — Dart 3.x + Flutter 3.x
- SDK constraint `>=3.0.0`
- Material 3, Impeller renderer
- GetX 4.x → 5.x (if ready)
- Dio 4 → Dio 5 (interceptor API changes)
- Docker: update to Flutter 3.x
- Android: Kotlin 1.6→1.9, Gradle 4.1→8.x, AGP 7.x+
