# Task: Fix AppPages.Deopsits Typo

**ID:** p1-typo  
**Phase:** 1 — Type Safety Foundation  
**Severity:** 🟠 MEDIUM  
**Dependencies:** None  

## Problem

`AppPages.Deopsits` is a typo — should be `Deposits`. This also means the route path is `/Deopsits` instead of `/Deposits`.

## Files to Modify

### 1. `src/app/constants.ts` (line 10)

**Current:**
```typescript
Deopsits = '/Deopsits',
```

**Target:**
```typescript
Deposits = '/Deposits',
```

### 2. Find and update all references

Search the codebase for `Deopsits` and `AppPages.Deopsits`:

```bash
grep -r "Deopsits" src/
grep -r "AppPages.Deopsits" src/
```

Update every occurrence to `Deposits` / `AppPages.Deposits`.

**Likely affected files:**
- `src/app/index.tsx` (route definition)
- `src/app/components/sideNav/index.tsx` or `mainCat.tsx` (navigation link)
- Any container that references this route for navigation

### ⚠️ Important: URL Breaking Change

If the application is currently deployed with `/Deopsits` as a live URL, changing to `/Deposits` will break bookmarks and direct links. Consider:
- Adding a redirect from `/Deopsits` to `/Deposits` in the router
- Or coordinating with backend if the URL is server-routed

## Validation

```bash
npm run checkTs          # Must pass
npm test                 # Must pass
grep -r "Deopsits" src/  # Must return zero results
```
