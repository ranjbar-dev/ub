# Task: Remove Duplicate GetWithdrawalCommentsAPI

**ID:** p2-duplicate-apis  
**Phase:** 2 — Service Layer Hardening  
**Severity:** 🟠 MEDIUM  
**Dependencies:** None  

## Problem

`GetWithdrawalCommentsAPI` exists in two files with identical implementations:
- `src/services/user_management_service.ts` (line ~68)
- `src/services/admin_reports_service.ts` (line 25)

Both call the same endpoint: `GET payment/user-comments`.

## Files to Modify

### 1. Determine canonical location

The function relates to withdrawal comments, which is billing/reports domain. Keep it in `admin_reports_service.ts`.

### 2. `src/services/user_management_service.ts` — Remove the duplicate

Delete the duplicate function:
```typescript
// DELETE THIS:
export const GetWithdrawalCommentsAPI = (parameters: any) => {
  return apiService.fetchData({
    data: parameters,
    url: 'payment/user-comments',
    requestType: RequestTypes.GET,
  });
};
```

### 3. Update all imports that used the deleted copy

Search for imports:
```bash
grep -rn "GetWithdrawalCommentsAPI" src/ --include="*.ts" --include="*.tsx"
```

**Expected matches:**
- `src/app/containers/Billing/saga.ts` — likely imports from user_management_service
- `src/app/containers/Reports/saga.ts` — likely imports from admin_reports_service

Update any import that pointed to `user_management_service` to point to `admin_reports_service`:

```typescript
// Before
import { GetWithdrawalCommentsAPI } from 'services/user_management_service';

// After
import { GetWithdrawalCommentsAPI } from 'services/admin_reports_service';
```

## Validation

```bash
npm run checkTs                              # Must pass
npm test                                      # Must pass
grep -rn "GetWithdrawalCommentsAPI" src/services/user_management_service.ts  # Must return nothing
```
