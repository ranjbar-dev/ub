# Domain Glossary — UnitedBit Exchange Admin

This glossary defines cryptocurrency exchange domain terminology as used throughout the `ub-admin-main` codebase. It is intended to help AI agents and new developers correctly interpret terms that may otherwise be ambiguous.

---

## Trading Terms

| Term | Definition | Used In |
|------|-----------|---------|
| **Currency Pair** | Two currencies traded against each other (e.g., BTC-USD, ETH-BTC). The base currency is bought/sold against the quote currency. | `CurrencyPairs` container, `GetCurrencyPairsAPI`, `UpdateCurrencyPairAPI` |
| **Pair** | Short form of "currency pair" — the ticker identifier such as `BTC/USD`. Stored as a string field on order records. | `Orders/types.ts` (`Order.pair`), `LocalStorageKeys.PAIRS` |
| **Order** | A request by a user to buy or sell a specific amount of a currency at a price. Contains fields: `amount`, `price`, `pair`, `type`, `side`. | `OpenOrders`, `FilledOrders`, `Orders` containers |
| **Open Order** | An order that has been submitted but not yet fully executed — it is still on the order book waiting to be matched. | `OpenOrders` container, `GetOpenOrdersAPI`, `MessageNames.SET_OPEN_ORDERS_DATA` |
| **Filled Order** | An order that has been completely executed — matched with a counter-party. Also called a "completed trade". | `FilledOrders` container, `MessageNames.SET_FILLED_ORDERS_DATA` |
| **Trade History** | Historical log of all completed trades for a user. | `GetTradeHistoryAPI`, `MessageNames.SET_TRADE_HISTORY_DATA` |
| **Order History** | Historical log of all orders (open and filled) for a user. | `MessageNames.SET_ORDER_HISTORY_DATA` |
| **Side** | Whether an order is a Buy or a Sell. Uses the `Sides` enum. | `OpenOrders/types.ts` `Sides` enum |
| **Maker** | The party who placed the order that rested on the order book. `isMaker: boolean` on order records distinguishes maker from taker. | `Orders/types.ts` (`Order.isMaker`) |
| **Taker** | The party who matched (filled) an existing order. `isMaker: false` on the order record. | `Orders/types.ts` |
| **Bot Trade** | A trade executed automatically by a trading bot. `isTradedByBot: boolean` on order records. | `Orders/types.ts` (`Order.isTradedByBot`) |
| **Market Tick** | A single OHLC (Open, High, Low, Close) candlestick price update for a trading pair. Used for charting. | `MarketTicks` container, `GetMarketTicksAPI`, `SyncTicksAPI` |
| **OHLC** | Open, High, Low, Close — the four price points describing a candle in a candlestick chart. | `MarketTicks` container |
| **Order Book** | The collection of all open buy and sell orders for a trading pair. Maintained in the backend (Go engine); displayed as a summary in admin. | Backend concept; visible in `OpenOrders` |
| **Liquidity Order** | An automated order placed by the platform itself to provide market liquidity and tighter spreads. | `LiquidityOrders` container, `GetLiquidityOrdersAPI` |
| **External Order** | An order routed to an external (third-party) exchange for execution when the internal order book cannot fill it. | `ExternalOrders` container, `GetExternalOrdersAPI` |
| **Net Queue** | Aggregated net positions from multiple orders, waiting to be submitted to an external exchange as a single batch. | `ExternalOrders` (`GetNetQueueAPI`, `ChangeNetQueueStatusAPI`, `SubmitNetQueueAPI`, `CancelNetQueueAPI`) |
| **All Queue** | The complete list of all items in the external order queue, including both net and individual entries. | `GetAllQueueAPI`, `MessageNames.SET_ALL_QUEUE_DATA` |
| **External Exchange** | A third-party cryptocurrency exchange that the platform integrates with for order routing and liquidity. | `ExternalExchange` container, `GetExternalExchangeAPI` |
| **Aggregation Status** | The processing state of a net queue entry — whether it is pending, submitted, or cancelled. | `ChangeNetQueueStatusAPI` |
| **Sync** | The process of pulling the latest market tick (OHLC) data from an external source into the platform. | `SyncTicksAPI`, `GetSyncListAPI`, `MessageNames.SET_SYNC_LIST_DATA` |
| **Fee** | A transaction or commission fee charged on a trade or transfer. Stored on `Order` and `Balance` records. | `Orders/types.ts`, `UserDetails/types.ts` |
| **Price** | The quoted exchange rate for a trading pair at which an order is placed or executed. | `Orders/types.ts` (`Order.price`) |
| **Amount** | The quantity of base currency in an order (how much to buy or sell). | `Orders/types.ts` (`Order.amount`) |
| **Executed** | The portion of an order that has been matched and filled. `Order.executed` tracks partial fills. | `Orders/types.ts` (`Order.executed`) |

---

## Financial Terms

| Term | Definition | Used In |
|------|-----------|---------|
| **Deposit** | Funds added to a user's exchange account, either via crypto blockchain transfer or fiat payment. | `Deposits` container, `UpdateDepositAPI`, `MessageNames.SET_DEPOSITS_DATA` |
| **Withdrawal** | Funds removed from a user's exchange account and sent to an external address. | `Withdrawals` container, `UpdateWithdrawAPI`, `MessageNames.SET_WITHDRAWALS_DATA` |
| **Payment** | Generic term for any financial transaction (deposit or withdrawal). The `Billing` container lists payments. | `Billing` container, `GetPaymentAPI`, `Billing/types.ts` |
| **Billing** | The aggregate view of all user financial transactions (deposits + withdrawals + transfers). | `Billing` container, `GetBillingGridDataAPI` |
| **Balance** | The amount of a specific currency held by a user. Broken into `availableAmount`, `inOrderAmount`, and `totalAmount`. | `Balances` container, `UserDetails/types.ts` (`Balance` interface) |
| **Available Amount** | The portion of a balance that is free to trade or withdraw (not locked in open orders). | `UserDetails/types.ts` (`Balance.availableAmount`) |
| **In-Order Amount** | The portion of a balance currently locked in open orders and not available for withdrawal. | `UserDetails/types.ts` (`Balance.inOrderAmount`) |
| **Total Amount** | The sum of available and in-order balances. `totalAmount = availableAmount + inOrderAmount`. | `UserDetails/types.ts` (`Balance.totalAmount`) |
| **BTC Equivalent** | Balance values converted to Bitcoin for a unified summary view. Fields prefixed `btc` (e.g., `btcAvailableSum`). | `UserDetails/types.ts` (`UserBalances`) |
| **Commission** | Fee charged to the user on completed trades, or fee earned by a referrer. Note: spelled `commitions` (typo) throughout the codebase. | `GetCommitionsAPI`, `MessageNames.SET_COMMITIONS_DATA`, `LiquidityOrders` |
| **Finance Method** | A configured payment method used for deposits/withdrawals (e.g., bank transfer, crypto wallet). | `FinanceMethods` container, `GetFinanceMethodsAPI`, `UpdateFinancialMethodAPI` |
| **Transfer** | An internal balance movement between a user's own accounts or between currencies within the platform. | `InternalTransferAPI`, `GetBalanceHistoryAPI`, `MessageNames.SET_BALANCES_HISTORY_DATA` |
| **Transaction (tx)** | A blockchain transaction record. `txId` stores the on-chain transaction hash. | `Billing/types.ts` (`Payment.txId`) |
| **From Address** | The source wallet address of a deposit or withdrawal. | `Billing/types.ts` (`Payment.fromAddress`) |
| **To Address** | The destination wallet address of a deposit or withdrawal. | `Billing/types.ts` (`Payment.toAddress`) |
| **Rejection Reason** | The explanation given when a deposit, withdrawal, or KYC document is rejected. | `Billing/types.ts` (`PaymentDetails.rejectionReason`), `VerificationWindow/types.ts` |
| **Auto Transfer** | Flag indicating whether a deposit should be automatically transferred to the user's trading balance. | `Billing/types.ts` (`PaymentDetails.autoTransfer`) |
| **Should Deposit** | Boolean flag on a deposit record controlling whether the funds should be credited to the user. | `Billing/types.ts` (`Payment.should_deposit`, `DepositSaveData.should_deposit`) |
| **Admin Comment** | A note added by admin staff to a payment transaction or user account for internal tracking. | `Reports` container, `AddAdminCommentAPI`, `AddPaymentCommentAPI`, `Billing/types.ts` (`IAdminComment`) |
| **Wallet** | A cryptocurrency account on the platform — classified as Hot, Cold, Internal, or External. | `Balances` container, `Balances/types.ts` (`IWallet`, `WalletTypes`) |
| **Hot Wallet** | A wallet connected to the internet, used for active trading and quick withdrawals. | `WalletTypes.Hot` |
| **Cold Wallet** | An offline wallet used for long-term storage of platform reserves. | `WalletTypes.Cold` |
| **Total Balance (user)** | The sum of all of a user's balances across all currencies. Displayed on the user profile. | `UserAccounts/types.ts` (`InitialUserDetails.totalBalance`) |
| **Total Deposit** | Cumulative amount deposited by a user across their account lifetime. | `InitialUserDetails.totalDeposit` |
| **Total Withdraw** | Cumulative amount withdrawn by a user across their account lifetime. | `InitialUserDetails.totalWithdraw` |
| **Total On-Trade** | Cumulative amount of user funds currently or historically committed to trades. | `InitialUserDetails.totalOnTrade` |
| **Total Commissions** | Cumulative commission fees charged to or earned by a user. | `InitialUserDetails.totalCommissions` |

---

## User / Account Terms

| Term | Definition | Used In |
|------|-----------|---------|
| **KYC** | Know Your Customer — the legal process of verifying a user's identity using government-issued documents. | `VerificationWindow` container, `GetUserImagesAPI` |
| **AML** | Anti-Money Laundering — compliance requirements to prevent money laundering. Enforced via KYC and transaction monitoring. | Backend enforcement; referenced in `VerificationWindow` |
| **User Verification** | The admin-facing process of reviewing and approving/rejecting a user's submitted identity documents. | `VerificationWindow` container |
| **Confirmation Status** | The KYC verification state for a specific document category. Uses the `ConfirmationStatus` enum. | `UserAccounts/types.ts` (`ConfirmationStatus`), `User` interface fields |
| **Identity Confirmation Status** | Verification status of the user's ID document (passport, national ID). | `User.identityConfirmationStatus`, `InitialUserDetails.identityConfirmationStatus` |
| **Address Confirmation Status** | Verification status of the user's proof-of-address document. | `User.addressConfirmationStatus` |
| **Phone Confirmation Status** | Verification status of the user's phone number. | `User.phoneConfirmationStatus` |
| **Profile Status** | The overall status/state of a user's profile (active, pending, blocked, etc.). | `User.profileStatus`, `InitialUserDetails.profileStatus` |
| **Account Status** | The operational status of a user's exchange account. | `InitialUserDetails.accountStatus` |
| **Profile Image Data** | Metadata for a KYC document image uploaded by the user, including document type and verification status. | `VerificationWindow/types.ts` (`ProfileImageData`) |
| **Image Status** | The verification state of a submitted KYC document image: `CONFIRMED` or `REJECTED`. | `VerificationWindow/types.ts` (`ImageStatusStrings`) |
| **Identity Type** | The category of KYC document: `identity` (ID/passport) or `address` (proof of address). | `VerificationWindow/constants.ts` (`IdentityTypes`) |
| **Sub Type** | A sub-classification of a document image (e.g., front vs back, document type variant). | `VerificationWindow/types.ts` (`ProfileImageData.subType`) |
| **Is Back** | Boolean flag indicating whether a document image is the back side of a two-sided ID. | `ProfileImageData.isBack` |
| **White Address** | A pre-approved cryptocurrency withdrawal address added to a user's allowlist. Withdrawals to non-white addresses may be blocked. | `GetUserWhiteAddressesAPI`, `MessageNames.SET_WHITEADDRESSES_DATA`, `UserDetails/types.ts` (`Address`) |
| **Is Favorite** | Flag on a white address indicating it is a preferred/starred address. | `UserDetails/types.ts` (`Address.isFavorite`) |
| **Permission** | A capability or access right that can be granted or revoked for a user account. | `GetUserPermissionsAPI`, `UpdateUserPermissionsAPI`, `UserDetails/types.ts` (`Permission`) |
| **Trust Level** | A numeric score indicating the credibility/risk level of a user. Affects withdrawal limits and permissions. | `InitialUserDetails.trustLevel` |
| **User Level** | A tier classification for users (e.g., Bronze, Silver, Gold) affecting fees, limits, and features. | `InitialUserDetails.userLevelId`, `userLevelName` |
| **User Group** | A classification category for grouping users with similar properties or treatment. | `InitialUserDetails.groupId`, `groupName` |
| **Manager** | An admin-level user assigned to manage/oversee a specific client account. | `InitialUserDetails.managerId`, `managerName`, `GetManagersAPI` |
| **Refer Key** | A unique referral code assigned to a user for the referral program. | `User.referKey` |
| **Referral ID** | The ID of the user who referred this user to the platform. | `User.referralId` |
| **Registered IP** | The IP address used when the user first registered their account. | `User.registeredIP` |
| **Login History** | An audit log of user authentication events, including IP address, device, and success/failure status. | `LoginHistory` container, `GetLoginHistoryAPI` |
| **Login State** | The outcome of a login attempt: `SUCCESSFUL` or `FAILED`. Uses `StateStrings` enum. | `LoginHistory/types.ts` (`StateStrings`) |
| **Admin** | A privileged platform user with access to the admin panel. | `Admins` container |
| **System ID** | An internal system-assigned identifier for a user, separate from the primary `id`. | `InitialUserDetails.systemId` |

---

## Technical Terms

| Term | Definition | Used In |
|------|-----------|---------|
| **StandardResponse** | The standard API response envelope: `{ status: boolean, message: string, data: T, token?: string }`. All API calls return this shape. | `src/services/constants.ts` |
| **RequestParameters** | The shape of an outbound API request: `{ requestType, url, data, isRawUrl?, requestName? }`. | `src/services/constants.ts` |
| **MessageService** | A singleton RxJS pub/sub service for inter-component communication. Components publish and subscribe to named events. | `src/services/message_service.ts` |
| **MessageNames** | An enum of 67+ named event types for the `MessageService`. Each name corresponds to a specific data update or UI action. | `src/services/message_service.ts` |
| **BroadcastMessage** | The message envelope published through `MessageService`: `{ name: MessageNames, data: any }`. | `src/services/message_service.ts` |
| **GridNames** | An enum identifying specific AG Grid instances for targeted updates via `MessageService`. | `src/services/message_service.ts` |
| **Container** | A page-level React component bundled with its own Redux slice, saga, selectors, and types. Each container in `app/containers/` is a self-contained feature module. | `src/app/containers/` |
| **Slice** | A Redux Toolkit state slice for a container, defining the initial state and reducers. Named `slice.ts` in each container. | `src/app/containers/*/slice.ts` |
| **Saga** | A Redux-Saga generator function handling async side effects (API calls). Named `saga.ts` in each container. | `src/app/containers/*/saga.ts` |
| **Selector** | A memoized function (via `createSelector`) that derives data from the Redux store. Named `selectors.ts`. | `src/app/containers/*/selectors.ts` |
| **SimpleGrid** | An AG Grid wrapper component providing a consistent data table UI across all admin pages. | `src/app/components/` |
| **AppPages** | An enum mapping page names to their React Router route paths (e.g., `AppPages.UserAccounts = '/userAccounts/'`). | `src/app/constants.ts` |
| **WindowTypes** | An enum of modal/popup window types that can be opened via `MessageNames.OPEN_WINDOW`. | `src/app/constants.ts` (`WindowTypes`) |
| **BaseUrl** | The configured API server base URL. Set per environment via `.env` files. | `src/services/constants.ts` |
| **LocalStorageKeys** | An enum of all `localStorage` key names used by the admin app for persisting session data and preferences. | `src/services/constants.ts` (`LocalStorageKeys`) |
| **UploadUrls** | An enum of relative URL paths for file upload endpoints (e.g., KYC image upload). | `src/services/constants.ts` (`UploadUrls`) |
| **RootState** | The TypeScript interface representing the complete Redux store shape, composed of all container states. | `src/types/RootState.ts` |
| **ApiService** | The singleton HTTP client class used for all API communication. Always accessed via `ApiService.getInstance()`. | `src/services/` |
| **Loadable** | A lazy-loading wrapper component for each container page, using `React.lazy` for code splitting. | `src/app/containers/*/Loadable.tsx` |
| **ALLPY_PARAMS_TO_GRID** | MessageService event that applies filter/sort parameters to a grid. (Note: typo for "APPLY".) | `MessageNames.ALLPY_PARAMS_TO_GRID` |

---

## Enum Reference

### `Sides` — Order Direction
Defined in `src/app/containers/OpenOrders/types.ts`

| Value | Meaning |
|-------|---------|
| `Buy = 'buy'` | A purchase order (buying base currency) |
| `BUY = 'BUY'` | Uppercase variant, used in some API responses |
| `Sell = 'sell'` | A sell order (selling base currency) |
| `SELL = 'SELL'` | Uppercase variant, used in some API responses |

---

### `DepositStatusStrings` — Deposit Lifecycle
Defined in `src/app/containers/Deposits/types.ts`

| Value | Meaning |
|-------|---------|
| `Created = 'created'` | Deposit initiated but not yet confirmed on-chain |
| `InProgress = 'in progress'` | Deposit is being processed |
| `Confirmed = 'CONFIRMED'` | Deposit confirmed on the blockchain |
| `Completed = 'completed'` | Deposit fully processed and credited to user balance |
| `COMPLETED = 'COMPLETED'` | Uppercase variant of completed |
| `Rejected = 'reject'` | Deposit rejected by admin (lowercase variant) |
| `Rejectedd = 'REJECTED'` | Deposit rejected — uppercase variant (note: enum key has a typo) |

---

### `ConfirmationStatus` — KYC Verification States
Defined in `src/app/containers/UserAccounts/types.ts`

| Value | Meaning |
|-------|---------|
| `NotConfirmed = 'not_confirmed'` | Document/identity not yet submitted or reviewed |
| `Confirmed = 'confirmed'` | Document/identity verified and approved |
| `Incomplete = 'incomplete'` | Partial submission — some documents still missing |
| `Rejected = 'rejected'` | Document/identity rejected by admin; user must resubmit |

---

### `ImageStatusStrings` — KYC Document Image Status
Defined in `src/app/containers/VerificationWindow/types.ts`

| Value | Meaning |
|-------|---------|
| `Confirmed = 'CONFIRMED'` | Document image accepted and verified |
| `Rejected = 'REJECTED'` | Document image rejected; reason stored in `rejectionReason` |

---

### `IdentityTypes` — KYC Document Categories
Defined in `src/app/containers/VerificationWindow/constants.ts`

| Value | Meaning |
|-------|---------|
| `Identity = 'identity'` | Government-issued ID document (passport, national ID, driver's license) |
| `Address = 'address'` | Proof-of-address document (utility bill, bank statement) |

---

### `StateStrings` — Login Attempt Result
Defined in `src/app/containers/LoginHistory/types.ts`

| Value | Meaning |
|-------|---------|
| `Successful = 'SUCCESSFUL'` | Login attempt succeeded |
| `Failed = 'FAILED'` | Login attempt failed (wrong password, blocked, etc.) |

---

### `WalletTypes` — Platform Wallet Classification
Defined in `src/app/containers/Balances/types.ts`

| Value | Meaning |
|-------|---------|
| `Hot = 'hot'` | Internet-connected wallet used for active trading |
| `Cold = 'cold'` | Offline wallet used for long-term reserve storage |
| `Internal = 'internal'` | Internal platform wallet (not user-facing) |
| `External = 'external'` | Wallet on an external exchange |

---

### `WindowTypes` — Modal Window Types
Defined in `src/app/constants.ts`

| Value | Meaning |
|-------|---------|
| `User = 'user'` | User detail popup window |
| `Verification = 'verification'` | KYC document verification popup |

---

### `AppPages` — Route Paths
Defined in `src/app/constants.ts`

| Enum Key | Path | Feature |
|----------|------|---------|
| `RootPage` | `/` | Root redirect |
| `LoginPage` | `/login` | Admin authentication |
| `HomePage` | `/home` | Dashboard |
| `UserAccounts` | `/userAccounts/` | User management list |
| `LoginHistory` | `/loginHistory` | Login audit trail |
| `OpenOrders` | `/OpenOrders` | Active trading orders |
| `FilledOrders` | `/FilledOrders` | Completed trades |
| `ExternalOrders` | `/ExternalOrders` | External exchange orders |
| `ExternalExchange` | `/ExternalExchange` | Exchange integration config |
| `MarketTicks` | `/MarketTicks` | OHLC market data management |
| `Deposits` | `/Deposits` | Deposit transactions |
| `Withdrawals` | `/Withdrawals` | Withdrawal transactions |
| `FinanceMethods` | `/FinanceMethods` | Payment method configuration |
| `Balances` | `/Balances` | Crypto wallet balances |
| `ScanBlock` | `/ScanBlock` | Blockchain scanner |
| `Admins` | `/Admins` | Admin user management |
| `CurrencyPairs` | `/CurrencyPairs` | Trading pair configuration |
| `LiquidityOrders` | `/LiquidityOrders` | Commission/liquidity reports |

---

## Known Typos in the Codebase

AI agents should be aware of the following misspellings present in the codebase. Do not "correct" them in code — they are part of the existing API contract.

| Misspelling | Correct Form | Where |
|-------------|-------------|-------|
| `commitions` | commissions | `GetCommitionsAPI`, `MessageNames.SET_COMMITIONS_DATA`, `LiquidityOrders` |
| `Rejectedd` | Rejected | `DepositStatusStrings.Rejectedd` enum key |
| `ALLPY_PARAMS_TO_GRID` | APPLY_PARAMS_TO_GRID | `MessageNames.ALLPY_PARAMS_TO_GRID` |
