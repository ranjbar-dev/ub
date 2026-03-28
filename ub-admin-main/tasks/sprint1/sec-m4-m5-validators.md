# Task: Fix Email Validator Regex + Strengthen Password Validation

## Priority: 🟡 MEDIUM (SEC-M4 + SEC-M5)
## Files to Modify: `src/app/containers/LoginPage/validators/emailValidator.tsx`, `src/app/containers/LoginPage/validators/passwordValidator.tsx`

## Problem
1. Email regex has unescaped `.` that matches ANY character (e.g., `user@domainXcom` passes)
2. Password validation only requires length ≥ 8 — no complexity requirements for a crypto exchange admin panel

## Current Code

### `src/app/containers/LoginPage/validators/emailValidator.tsx`:
Look for the regex pattern. It should be around line 3:
```typescript
const emailtext = new RegExp(
    '^[a-zA-Z0-9.!#$%&\'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:.[a-zA-Z0-9-]+)*$',
);
```
**Issue:** `(?:.[a-zA-Z0-9-]+)` — the `.` is unescaped, matching ANY character instead of literal dot.

### `src/app/containers/LoginPage/validators/passwordValidator.tsx`:
Look for the length check. Should be around lines 17-20:
```typescript
if (value == null || value.length === 0) {
    sendError(properties.errors[0]);
    return false;
} else if (value.length < 8) {
    sendError(properties.errors[1]);
    return false;
}
```

## Required Changes

### 1. Fix email regex in `emailValidator.tsx`:

Find the unescaped `.` in the regex and replace with `\\.`:
```typescript
const emailtext = new RegExp(
    '^[a-zA-Z0-9.!#$%&\'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-]+)*$',
);
```

The key change is `(?:.)` → `(?:\\.)` — escaping the dot so it only matches literal dots.

### 2. Strengthen password validation in `passwordValidator.tsx`:

Find the `value.length < 8` check and update the minimum to 12 characters. Add a complexity check:

Replace the simple length check with:
```typescript
} else if (value.length < 12) {
    sendError(properties.errors[1]);
    return false;
} else if (!/(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/.test(value)) {
    sendError(properties.errors[2] || 'Password must contain uppercase, lowercase, and a number');
    return false;
}
```

**IMPORTANT:** Read the file first to understand the full structure and the `properties.errors` array. If `properties.errors[2]` doesn't exist, use the inline string. The `sendError` function signature can be found in the same file — it likely calls a message service or sets state.

**NOTE:** If adding a new error message would require changes elsewhere (e.g., i18n translations, error arrays in the parent component), add a comment like `// TODO: Add i18n key for password complexity error` and use the inline English string as fallback.

## Validation
- Build must pass: `$env:NODE_OPTIONS='--openssl-legacy-provider'; npm run build`
- Verify: `user@domainXcom` no longer passes email validation
- Verify: `user@domain.com` still passes
- Verify: `12345678` (8 chars, no complexity) no longer passes password validation
- Verify: `MySecure123!` (12+ chars, mixed) passes
