# Task: Add TypeScript Types to All Service Functions

**ID:** p2-type-services  
**Phase:** 2 — Service Layer Hardening  
**Severity:** 🔴 CRITICAL  
**Dependencies:** p1-generic-response, p1-request-params  

## Problem

All 56 service functions across 12 files use `parameters: any`. No function has typed params or return types. AI agents cannot determine what data each API endpoint expects or returns.

## Scope

12 service files, 56 functions total:

| File | Functions | Count |
|------|-----------|-------|
| `src/services/user_management_service.ts` | GetUsersAPI, GetUserDetailsAPI, UpdateUserInfoAPI, ChangeUserPasswordAPI, ChangeAccountStatusAPI, VerifyDocumentsAPI, UpdateWhiteAddressAPI, DeleteWhiteAddressAPI, GetWhiteAddressesAPI, GetWalletsAPI, GetBalancesAPI, GetUserBalancesHistoryAPI, GetDepositHistoryAPI, GetWithdrawalHistoryAPI, GetOpenOrdersAPI, GetFilledOrdersAPI, GetWithdrawalCommentsAPI, DeleteCommentAPI, AddCommentAPI, EditCommentAPI, DownloadDocumentsAPI, ActivateBankAccountAPI, AddTransferAPI, UploadUserDocumentAPI, GetUserImagesAPI, RejectDocumentAPI, ApproveDocumentAPI, GetUserLoginHistoryAPI, GetUserNotificationPreferences | 29 |
| `src/services/external_orders_service.ts` | GetExternalOrdersAPI, GetNetQueueAPI, GetAllQueueAPI, ChangeNetQueueStatusAPI, CancelNetQueueAPI, SubmitNetQueueAPI | 6 |
| `src/services/admin_reports_service.ts` | AddAdminCommentAPI, DeleteAdminCommentAPI, EditAdminCommentAPI, GetWithdrawalCommentsAPI, UpdateFinancialMethodAPI, UpdateCurrencyPairAPI, GetCommitionsAPI | 7 |
| `src/services/security_service.ts` | loginAPI | 1 |
| `src/services/markets_service.ts` | (market-related APIs) | ~5 |
| `src/services/deposits_service.ts` | (deposit-related APIs) | ~3 |
| `src/services/withdrawals_service.ts` | (withdrawal-related APIs) | ~4 |
| `src/services/billing_service.ts` | (billing-related APIs) | ~3 |
| `src/services/balances_service.ts` | (balance-related APIs) | ~3 |
| `src/services/admins_service.ts` | (admin management APIs) | ~3 |
| `src/services/orders_service.ts` | (order-related APIs) | ~4 |
| `src/services/external_exchange_service.ts` | (external exchange APIs) | ~3 |

## Pattern

### Before (every function looks like this)
```typescript
export const GetUsersAPI = (parameters: any) => {
  return apiService.fetchData({
    data: parameters,
    url: 'admin/user/list',
    requestType: RequestTypes.GET,
  });
};
```

### After
```typescript
/** Fetches paginated list of users for the admin panel. */
export const GetUsersAPI = (
  parameters: GetUsersParams
): Promise<StandardResponse<UserListResponse>> => {
  return apiService.fetchData({
    data: parameters,
    url: 'admin/user/list',
    requestType: RequestTypes.GET,
  });
};
```

### Type Definitions

Create types alongside each service file or in a shared `src/services/types/` directory:

```typescript
// src/services/types/user_management.types.ts

export interface PaginationParams {
  page?: number;
  per_page?: number;
  sort_by?: string;
  sort_dir?: 'asc' | 'desc';
}

export interface GetUsersParams extends PaginationParams {
  search?: string;
  status?: string;
}

export interface UserListResponse {
  items: User[];
  total: number;
  page: number;
  per_page: number;
}

export interface GetUserDetailsParams {
  user_id: number;
}

export interface ChangeAccountStatusParams {
  user_id: number;
  status: string;
  reason?: string;
}
```

## Execution Steps

1. Create `src/services/types/` directory
2. For each service file:
   a. Read the file and identify all function names
   b. Read the saga(s) that call these functions to understand the data shapes
   c. Check the API URL patterns for clues about expected params
   d. Create parameter interfaces (InputParams) and response interfaces
   e. Apply types to the function signature
   f. Add JSDoc comment with brief description
3. Run `npm run checkTs` after each file

## Example: `external_orders_service.ts` Full Conversion

```typescript
import { apiService } from './api_service';
import { RequestTypes, StandardResponse } from './constants';

export interface ExternalOrdersParams {
  page?: number;
  per_page?: number;
  pair?: string;
  status?: string;
}

export interface NetQueueParams {
  page?: number;
  per_page?: number;
}

export interface ChangeQueueStatusParams {
  id: number;
  status: string;
}

/** Fetches external orders list. */
export const GetExternalOrdersAPI = (
  parameters: ExternalOrdersParams
): Promise<StandardResponse> => {
  return apiService.fetchData({
    data: parameters,
    url: 'exchange/order',
    requestType: RequestTypes.GET,
  });
};

/** Fetches net queue entries. */
export const GetNetQueueAPI = (
  parameters: NetQueueParams
): Promise<StandardResponse> => {
  return apiService.fetchData({
    data: parameters,
    url: 'exchange/order/queue',
    requestType: RequestTypes.GET,
  });
};
// ... same pattern for remaining functions
```

## Validation

```bash
npm run checkTs   # Must pass after each file
npm test          # Must pass
```
