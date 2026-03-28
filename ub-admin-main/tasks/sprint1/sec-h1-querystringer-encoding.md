# Task: Fix queryStringer URL Encoding

## Priority: 🟠 HIGH (SEC-H1)
## Files to Modify: `src/utils/formatters.ts`

## Problem
`queryStringer()` builds query strings by directly concatenating param values without `encodeURIComponent()`. This allows parameter injection, URL structure manipulation, and potential XSS when user-controlled data flows into GET parameters (search filters, sort fields, pagination).

## Current Code

### `src/utils/formatters.ts` lines 1-13:
```typescript
export function queryStringer(params: any): string {
	let qs = '?';
	let counter = 0;
	for (let key in params) {
		let prefix = counter === 0 ? '' : '&';
		qs += prefix + key + '=' + params[key];    // ⚠️ Line 6 - NO ENCODING!
		counter++;
	}
	if(counter===0){
		qs=''
	}
	return qs;
}
```

### Called from `src/services/apiService.ts` line 47:
```typescript
const url = params.requestType === RequestTypes.GET
    ? baseUrl + params.url + queryStringer(params.data)
    : baseUrl + params.url;
```

## Required Changes

Replace the entire `queryStringer` function in `src/utils/formatters.ts`:
```typescript
/**
 * Builds a URL query string from an object of key-value pairs.
 * All keys and values are URL-encoded to prevent injection.
 *
 * @param params - Object with string key-value pairs
 * @returns Query string starting with '?' or empty string if no params
 */
export function queryStringer(params: Record<string, unknown>): string {
	const entries = Object.entries(params).filter(
		([, value]) => value !== undefined && value !== null,
	);
	if (entries.length === 0) {
		return '';
	}
	const qs = entries
		.map(([key, value]) => `${encodeURIComponent(key)}=${encodeURIComponent(String(value))}`)
		.join('&');
	return `?${qs}`;
}
```

Changes:
- Added `encodeURIComponent()` on both keys AND values
- Typed params as `Record<string, unknown>` instead of `any`
- Filters out null/undefined values
- Uses modern `Object.entries` + `.map().join()` pattern
- Added JSDoc

## Validation
- Build must pass: `$env:NODE_OPTIONS='--openssl-legacy-provider'; npm run build`
- Verify: `queryStringer({ name: 'test&inject=true', page: 1 })` returns `?name=test%26inject%3Dtrue&page=1`
