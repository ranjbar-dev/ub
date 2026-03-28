# Unitedbit — Mobile & Web App

Cryptocurrency exchange mobile (Android/iOS) and web application built with Flutter.

## Prerequisites

| Requirement | Version | Notes |
|---|---|---|
| **Dart SDK** | `>=2.11.0 <3.0.0` | **CRITICAL**: Must use Dart 2.x. Dart 3.x will NOT compile this project. |
| **Flutter** | 2.10.5 | Last stable Flutter 2.x release. Flutter 3.x requires Dart 3.x. |
| **Android SDK** | API 33 (compile), API 21+ (min) | Android Studio or `sdkmanager` |
| **Xcode** | Latest compatible with Flutter 2.x | iOS builds only |
| **Chrome** | Any modern version | Web development |

> ⚠️ **This project is pre-null-safety.** If your system Flutter is 3.x, use Docker or FVM to manage SDK versions.

## Quick Start

```bash
# Install dependencies
flutter pub get

# Run on connected device (debug)
flutter run

# Run on Chrome (web)
flutter run -d chrome --web-renderer canvaskit --web-port 8000

# Run tests
flutter test
```

## Build Commands

### Android
```bash
# Dev APK (split per ABI)
./buildDevAPK.sh

# Production APK (split per ABI)
./buildAPK.sh

# Production App Bundle
./buildBundle.sh
```

### Web
```bash
# Dev build (canvaskit renderer)
./buildWeb-dev.sh

# Production build (canvaskit renderer)
./buildWeb.sh
```

### Docker
```bash
# Dev web (nginx)
docker build -f Dockerfiledev -t ub-app-dev .

# Production web (nginx)
docker build -f Dockerfileprod -t ub-app-prod .

# Production APK
docker build -f Dockerfileapkprod -t ub-app-apk .
```

### Local Dev
```bash
./runApp.sh          # flutter run with VERSION dart-define
./runWeb.sh          # flutter run -d chrome on port 8000
```

## Architecture

**GetX MVC module pattern** — each feature lives in `lib/app/modules/<feature>/`:

```
lib/app/modules/<feature>/
  bindings/<feature>_binding.dart     # GetX DI registration
  controllers/<feature>_controller.dart  # Business logic + .obs reactive state
  views/<feature>_view.dart           # UI (extends GetView<Controller>)
  providers/<feature>_provider.dart   # API calls via ApiService
  models/                             # Data models (optional)
```

### Key Architecture Decisions
- **State**: GetX `.obs` reactive variables + `Obx()` widgets (primary), Provider (secondary)
- **HTTP**: Dio 4.0.1 singleton with interceptor chain (retry, token refresh, auth header)
- **Real-time**: MQTT over WebSocket (`wss://`) — dual clients (authorized + unauthorized)
- **Storage**: GetStorage (plain KV) + FlutterSecureStorage (encrypted credentials)
- **Crypto**: RSA/AES via PointyCastle + Encrypt packages
- **Navigation**: 35 named routes via `GetPages`, initial route `/after-splash`

### 23 Feature Modules
| Module | Purpose |
|---|---|
| `login` | Email/password + biometric authentication |
| `signup` | User registration with RSA-encrypted captcha |
| `account` | Profile, settings, biometrics, KYC status |
| `home` | Dashboard with news, popular pairs, sparklines |
| `trade` | Market/limit/stop orders, OHLC charts, order book (MQTT) |
| `exchange` | Simple currency swap |
| `market` | Ticker list, favorites, search, sorting |
| `funds` | Balance, deposits, withdrawals, transaction history, auto-exchange |
| `orders` | Open orders (MQTT live), order history |
| `identityInfo` | KYC personal information |
| `identityDocuments` | KYC document upload |
| `phoneVerification` | Phone OTP verification |
| `twoFactorAuthentication` | Google Authenticator 2FA setup |
| `changePassword` | Password change |
| `forgot` | Password recovery |
| `withdrawAddressManagement` | Manage crypto withdrawal addresses |
| `addNewAddress` | Add new blockchain address |
| `checkYourEmail` | Post-signup email verification |
| `qrScan` | QR code scanner for addresses |
| `webViewPage` | Embedded web content (T&Cs, help) |
| `separateMessagePage` | Generic info/error display |
| `landing` | Pre-auth landing page |
| `afterSplash` | Auth-based route decision |

## Project Structure

```
lib/
  main.dart                    # Entry point (GetMaterialApp)
  app/
    modules/                   # 23 feature modules (GetX MVC)
    global/                    # GlobalController, MQTT controllers, common provider
    common/components/         # 70 reusable UI widgets
    common/custom/             # 49 custom widget libraries
    popups/                    # 2FA popup system
    routes/                    # 35 named routes
  services/                    # ApiService, constants, storage keys, interceptors
  mqttClient/                  # Universal MQTT client (browser + VM)
  utils/                       # Crypto, validators, extensions, mixins
  theme/                       # Dark/light theme definitions
  generated/                   # flutter_gen output (colors, fonts, assets, locales)
```

## Environment Configuration

| Variable | Description | Default |
|---|---|---|
| `VERSION` | App version (via `--dart-define`) | `1.0.0` |
| `ENV` | Environment mode | Hardcoded `"PRODUCTION"` in `ubEnv.dart` |

API base: `https://[dev-]app.unitedbit.com/api/v1/`
MQTT broker: `wss://[dev-]app.unitedbit.com:8443`

## CI/CD

GitLab CI pipeline with 6 stages:
- **Dev**: develop branch → Docker build → deploy to dev-m server → Telegram notification
- **Prod**: master branch → Docker build → deploy to prod-m server → Telegram notification

## Code Generation

```bash
# Regenerate assets, colors, fonts, locales
flutter pub run build_runner build
```

## Documentation

See [AGENTS.md](AGENTS.md) for the complete technical reference (1,400+ lines) covering:
- All dependencies, routes, controllers, providers, API endpoints
- MQTT topics, storage keys, component library
- Build system, platform config, testing, known issues
- Upgrade roadmap
