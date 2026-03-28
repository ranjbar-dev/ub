# Task: Type RequestParameters.data

**ID:** p1-request-params  
**Phase:** 1 — Type Safety Foundation  
**Severity:** 🔴 CRITICAL  
**Dependencies:** None  
**Blocks:** p2-type-services  

## Problem

`RequestParameters.data` is typed as `any`, meaning all request payloads sent to the API have zero type checking.

## File to Modify

**`src/services/constants.ts`**

### Current Code (lines 7–13)
```typescript
export interface RequestParameters {
	requestType: RequestTypes;
	url: string;
	data: any;
	isRawUrl?: boolean;
	requestName?: string;
}
```

### Target Code
```typescript
export interface RequestParameters<T = Record<string, unknown>> {
	requestType: RequestTypes;
	url: string;
	data: T;
	isRawUrl?: boolean;
	requestName?: string;
}
```

## Why This Is Safe

The default `T = Record<string, unknown>` means existing callers that pass `RequestParameters` without a type argument continue to compile. All service functions pass object literals to `data:`, which satisfy `Record<string, unknown>`.

## Affected Consumers

The main consumer is `ApiService.fetchData()` in `src/services/api_service.ts`:

```typescript
// Current (line 24)
public async fetchData(params: RequestParameters) {

// After this change — no change needed, default T applies
public async fetchData(params: RequestParameters) {
```

All 12 service files pass `{ data: parameters, url: '...', requestType: RequestTypes.X }` which continues to work.

## Validation

```bash
npm run checkTs   # Must pass
npm test          # Must pass
```
