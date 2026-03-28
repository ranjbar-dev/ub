# Task: Add JSDoc to All Service Functions

**ID:** p2-jsdoc-services  
**Phase:** 2 — Service Layer Hardening  
**Severity:** 🟡 HIGH  
**Dependencies:** None (can run in parallel with p2-type-services)  

## Problem

Zero service functions have JSDoc comments. AI agents cannot determine what each API call does without reading the URL and request type.

## Scope

All 56 functions across 12 service files in `src/services/`.

## Pattern

### Before
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
/**
 * Fetches paginated list of registered users for admin management.
 * 
 * @param parameters - Filter/pagination params (page, per_page, search, status)
 * @returns Promise with user list in data field
 * @endpoint GET admin/user/list
 */
export const GetUsersAPI = (parameters: any) => {
  return apiService.fetchData({
    data: parameters,
    url: 'admin/user/list',
    requestType: RequestTypes.GET,
  });
};
```

### JSDoc Template

```typescript
/**
 * [Brief description of what this API call does].
 *
 * @param parameters - [What the parameters contain]
 * @returns Promise<StandardResponse> with [what data contains]
 * @endpoint [METHOD] [url]
 */
```

## Service-by-Service Reference

### `user_management_service.ts` (29 functions)
| Function | Description | Endpoint |
|----------|-------------|----------|
| GetUsersAPI | Fetches paginated user list | GET admin/user/list |
| GetUserDetailsAPI | Gets full details for a single user | GET admin/user/detail |
| UpdateUserInfoAPI | Updates user profile fields | POST admin/user/update |
| ChangeUserPasswordAPI | Admin-initiated password change | POST admin/user/change-password |
| ChangeAccountStatusAPI | Changes user account status (active/blocked/etc) | POST admin/user/change-status |
| VerifyDocumentsAPI | Marks user documents as verified | POST admin/user/verify-documents |
| UpdateWhiteAddressAPI | Updates a whitelisted withdrawal address | POST admin/user/white-address/update |
| DeleteWhiteAddressAPI | Removes a whitelisted withdrawal address | POST admin/user/white-address/delete |
| GetWhiteAddressesAPI | Fetches user's whitelisted addresses | GET admin/user/white-addresses |
| GetWalletsAPI | Gets user's crypto wallets | GET admin/user/wallets |
| GetBalancesAPI | Gets user's balance across all currencies | GET admin/user/balances |
| GetUserBalancesHistoryAPI | Balance change history for a user | GET admin/user/balances-history |
| GetDepositHistoryAPI | User's deposit transaction history | GET admin/user/deposit-history |
| GetWithdrawalHistoryAPI | User's withdrawal transaction history | GET admin/user/withdrawal-history |
| GetOpenOrdersAPI | User's currently open trading orders | GET admin/user/open-orders |
| GetFilledOrdersAPI | User's completed (filled) trading orders | GET admin/user/filled-orders |
| GetWithdrawalCommentsAPI | Admin comments on user withdrawals | GET payment/user-comments |
| DeleteCommentAPI | Deletes an admin comment | POST user/admin-comment/delete |
| AddCommentAPI | Adds admin comment to a user | POST user/admin-comment/add |
| EditCommentAPI | Edits existing admin comment | POST user/admin-comment/update |
| DownloadDocumentsAPI | Downloads user verification documents | GET admin/user/download-document |
| ActivateBankAccountAPI | Activates a user's bank account | POST admin/user/activate-bank |
| AddTransferAPI | Creates a manual balance transfer | POST admin/user/transfer |
| UploadUserDocumentAPI | Uploads documents for user verification | POST admin/user/upload-document |
| GetUserImagesAPI | Gets user verification images | GET admin/user/images |
| RejectDocumentAPI | Rejects a verification document | POST admin/user/reject-document |
| ApproveDocumentAPI | Approves a verification document | POST admin/user/approve-document |
| GetUserLoginHistoryAPI | User's login history with IPs | GET admin/user/login-history |
| GetUserNotificationPreferences | User's notification settings | GET admin/user/notification-preferences |

### `admin_reports_service.ts` (7 functions)
| Function | Description | Endpoint |
|----------|-------------|----------|
| AddAdminCommentAPI | Add admin comment to entity | POST user/admin-comment/add |
| DeleteAdminCommentAPI | Delete an admin comment | POST user/admin-comment/delete |
| EditAdminCommentAPI | Edit admin comment text | POST user/admin-comment/update |
| GetWithdrawalCommentsAPI | Get comments on a withdrawal | GET payment/user-comments |
| UpdateFinancialMethodAPI | Update a payment method config | POST currency/update |
| UpdateCurrencyPairAPI | Update trading pair settings | POST currency/update-pair |
| GetCommitionsAPI | Get user commission/fee statistics | GET statistic/user-statistic |

### `external_orders_service.ts` (6 functions)
| Function | Description | Endpoint |
|----------|-------------|----------|
| GetExternalOrdersAPI | List external exchange orders | GET exchange/order |
| GetNetQueueAPI | Get net queue entries | GET exchange/order/queue |
| GetAllQueueAPI | Get all queue entries | GET exchange/order/queue/all |
| ChangeNetQueueStatusAPI | Change aggregation status | POST exchange/aggregation/change-status |
| CancelNetQueueAPI | Cancel a queue order | POST exchange/order/change-status |
| SubmitNetQueueAPI | Submit a queue order | POST exchange/order/change-status |

## Execution Steps

1. For each service file, add JSDoc block above every exported function
2. Use the endpoint tables above as reference
3. Run `npm run checkTs` to verify no syntax errors

## Validation

```bash
npm run checkTs   # Must pass
```
