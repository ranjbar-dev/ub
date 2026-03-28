# Task: Create Domain Glossary

**ID:** p7-glossary  
**Phase:** 7 — Documentation for AI Agents  
**Severity:** 🟡 HIGH  
**Dependencies:** None  

## Problem

The app uses cryptocurrency exchange domain terminology that AI agents may not correctly interpret. Terms like "KYC", "AML", "order book", "fill", "pair", "net queue" appear throughout the codebase without definition.

## File to Create

**`docs/GLOSSARY.md`** (create `docs/` directory if needed)

### Content

```markdown
# Domain Glossary — UnitedBit Exchange

## Trading Terms

| Term | Definition | Used In |
|------|-----------|---------|
| **Currency Pair** | Two currencies traded against each other (e.g., BTC-USD, ETH-BTC) | CurrencyPairs container |
| **Order** | A request to buy or sell a specific amount of currency at a price | OpenOrders, FilledOrders |
| **Open Order** | An order that hasn't been fully executed yet | OpenOrders container |
| **Filled Order** | An order that has been completely executed (matched with counter-party) | FilledOrders container |
| **Side** | Whether an order is a Buy or Sell | OpenOrders/types.ts Sides enum |
| **Market Tick** | A single price update/trade event on a trading pair | MarketTicks container |
| **Order Book** | The collection of all open buy and sell orders for a pair | (backend concept, displayed in admin) |
| **Liquidity Order** | An automated order placed to provide market liquidity | LiquidityOrders container |
| **External Order** | An order routed to an external exchange for execution | ExternalOrders container |
| **Net Queue** | Aggregated net positions waiting for external exchange execution | ExternalOrders (GetNetQueueAPI) |

## Financial Terms

| Term | Definition | Used In |
|------|-----------|---------|
| **Deposit** | Funds added to a user's exchange account (crypto or fiat) | Deposits container |
| **Withdrawal** | Funds removed from a user's exchange account | Withdrawals container |
| **Balance** | User's available funds per currency | Balances container |
| **Commission** | Fee charged on trades (note: code uses "commitions" — typo) | GetCommitionsAPI |
| **Finance Method** | A payment method configuration (bank, crypto wallet, etc.) | FinanceMethods container |
| **Transfer** | Internal balance transfer between accounts or currencies | UserDetails (AddTransferAPI) |

## User/Account Terms

| Term | Definition | Used In |
|------|-----------|---------|
| **KYC** | Know Your Customer — identity verification process | VerificationWindow |
| **AML** | Anti-Money Laundering compliance | (backend enforcement) |
| **User Verification** | Process of reviewing uploaded ID documents | VerificationWindow |
| **White Address** | Pre-approved cryptocurrency withdrawal address | UserDetails (GetWhiteAddressesAPI) |
| **Account Status** | Active, Blocked, Pending, etc. | UserAccounts, ConfirmationStatus enum |
| **Admin Comment** | Notes added by admin staff to user accounts or payments | Reports (AddAdminCommentAPI) |
| **Login History** | Audit log of user authentication events with IPs | LoginHistory container |

## Technical Terms

| Term | Definition | Used In |
|------|-----------|---------|
| **StandardResponse** | API response format: `{ status: boolean, message: string, data: T }` | services/constants.ts |
| **MessageService** | RxJS pub/sub system for inter-component communication | services/message_service.ts |
| **MessageNames** | Enum of 67 event types for MessageService | services/message_service.ts |
| **BroadcastMessage** | Message shape for MessageService events | services/message_service.ts |
| **Container** | A page-level React component with its own slice/saga/selectors | app/containers/ |
| **SimpleGrid** | AG Grid wrapper component for data tables | app/components/ |
| **BaseUrl** | The API server base URL | services/constants.ts |
| **AppPages** | Enum mapping page names to route paths | app/constants.ts |

## Enum Reference

### DepositStatusStrings
- `Completed` / `COMPLETED` — deposit fully processed
- `Created` — deposit initiated but not confirmed
- `InProgress` — deposit being processed
- `Rejected` / `REJECTED` — deposit rejected by admin
- `Confirmed` / `CONFIRMED` — deposit confirmed on blockchain

### Sides
- `Buy` / `BUY` — purchase order
- `Sell` / `SELL` — sell order

### ConfirmationStatus
- `Confirmed` — user/document verified
- `NotConfirmed` — pending verification
- `Incomplete` — partial verification
- `Rejected` — verification denied

### WindowTypes
- `DEPOSIT_DETAILS`, `WITHDRAW_DETAILS`, etc. — popup window types
```

## Validation

No compilation needed — this is documentation.
