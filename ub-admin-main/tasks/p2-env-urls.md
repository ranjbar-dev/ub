# Task: Use Environment Variables for API URLs

**ID:** p2-env-urls  
**Phase:** 2 — Service Layer Hardening  
**Severity:** 🟠 MEDIUM  
**Dependencies:** None  

## Problem

`BaseUrl` in `src/services/constants.ts` is hardcoded. Different environments (dev, staging, production) need different API URLs.

## File to Modify

**`src/services/constants.ts`** (line 48)

### Current Code
```typescript
export const BaseUrl = 'https://api.example.com/api/v1/';
```

### Target Code
```typescript
export const BaseUrl: string =
  process.env.REACT_APP_API_BASE_URL || 'https://api.example.com/api/v1/';
```

### Also Update `.env.example` (create if not exists)

Create `ub-admin-main/.env.example`:
```env
# API base URL — include trailing slash
REACT_APP_API_BASE_URL=https://api.example.com/api/v1/
```

> **Note:** Create React App (CRA) automatically loads variables prefixed with `REACT_APP_` from `.env` files. No additional setup needed.

### Also Check `.env` and `.gitignore`

Ensure `.env` is in `.gitignore` (likely already is for CRA projects):
```
.env.local
.env.development.local
.env.test.local
.env.production.local
```

## Validation

```bash
npm run checkTs   # Must pass
npm run build     # Must build (uses env vars at build time)
```
