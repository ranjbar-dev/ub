# Task: Remove Dead Code from security_service.ts

**ID:** p2-dead-code  
**Phase:** 2 — Service Layer Hardening  
**Severity:** 🟠 MEDIUM  
**Dependencies:** None  

## Problem

`security_service.ts` has only 1 active function but contains 55 lines of commented-out dead code (lines 12–67). This confuses AI agents and bloats the file.

## File to Modify

**`src/services/security_service.ts`**

### Current Content (67 lines)
```typescript
import { apiService } from './api_service';
import { RequestTypes } from './constants';

export const loginAPI = (parameters: any) => {
  return apiService.fetchData({
    data: parameters,
    url: 'security/admin-sign-in',
    requestType: RequestTypes.POST,
    isRawUrl: true,
  });
};

// export const getCurrentUserDetailsAPI = () => {
// 	return apiService.fetchData({
// 		data: '',
// 		url: 'admin/current-admin',
// 		requestType: RequestTypes.GET,
// 		isRawUrl: true,
// 	});
// };
// ...  (55 more lines of commented-out functions)
```

### Target Content
```typescript
import { apiService } from './api_service';
import { RequestTypes } from './constants';

export interface LoginParams {
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: {
    id: number;
    email: string;
    name: string;
  };
}

/**
 * Authenticates an admin user and returns a JWT token.
 *
 * @param parameters - Login credentials (email, password)
 * @returns Promise with JWT token and user info
 * @endpoint POST security/admin-sign-in
 */
export const loginAPI = (parameters: LoginParams) => {
  return apiService.fetchData({
    data: parameters,
    url: 'security/admin-sign-in',
    requestType: RequestTypes.POST,
    isRawUrl: true,
  });
};
```

### What Was Removed
The following commented-out functions were deleted (they can be recovered from git history if needed):
- `getCurrentUserDetailsAPI` — GET admin/current-admin
- `logoutAPI` — POST security/logout
- Various other commented security endpoints

## Validation

```bash
npm run checkTs   # Must pass
npm test          # Must pass
```
