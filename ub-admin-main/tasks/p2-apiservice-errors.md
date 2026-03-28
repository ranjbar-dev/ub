# Task: Fix ApiService Error Handling

**ID:** p2-apiservice-errors  
**Phase:** 2 — Service Layer Hardening  
**Severity:** 🔴 CRITICAL  
**Dependencies:** p1-generic-response  
**Blocks:** p3-saga-utils  

## Problem

`ApiService.fetchData()` has fundamentally broken error handling:
1. On 422, it returns raw JSON but callers don't check for errors
2. On 401, it sends a MessageService event but returns `undefined` to the caller
3. On 500, it sends a toast message but returns `undefined` to the caller
4. No try-catch around `fetch()` — network errors crash the saga
5. In production mode, errors are completely silenced

## File to Modify

**`src/services/api_service.ts`**

### Current Code (lines 59–124)
```typescript
public async fetchData(params: RequestParameters) {
  const requestType = params.requestType;
  const baseUrl = params.isRawUrl ? BaseUrl : BaseUrl + 'admin/';
  const url = requestType === RequestTypes.GET
    ? baseUrl + params.url + queryStringer(params.data)
    : baseUrl + params.url;
  
  let content: RequestInit = {
    method: requestType,
    headers: this.setHeaders(),
    body: requestType === RequestTypes.POST ? JSON.stringify(params.data) : undefined,
  };

  let response = await fetch(url, content);
  
  if (response.status === 200) {
    let json = await response.json();
    return json;
  } else if (response.status === 422) {
    let json = await response.json();
    return json;
  } else if (response.status === 401) {
    MessageService.send({
      name: MessageNames.AUTH_ERROR_EVENT,
      value: 'auth error',
    });
  } else {
    if (process.env.NODE_ENV !== 'production') {
      let json = await response.json();
      MessageService.send({
        name: MessageNames.TOAST,
        value: json,
        type: 'error',
      });
    }
  }
}
```

### Target Code
```typescript
/** Represents an API error with status code and optional validation errors. */
export class ApiError extends Error {
  constructor(
    message: string,
    public readonly statusCode: number,
    public readonly errors?: Record<string, string[]>,
    public readonly rawResponse?: unknown,
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

public async fetchData<T = unknown>(
  params: RequestParameters
): Promise<StandardResponse<T>> {
  const requestType = params.requestType;
  const baseUrl = params.isRawUrl ? BaseUrl : BaseUrl + 'admin/';
  const url = requestType === RequestTypes.GET
    ? baseUrl + params.url + queryStringer(params.data)
    : baseUrl + params.url;

  const content: RequestInit = {
    method: requestType,
    headers: this.setHeaders(),
    body: requestType === RequestTypes.POST
      ? JSON.stringify(params.data)
      : undefined,
  };

  let response: Response;
  try {
    response = await fetch(url, content);
  } catch (networkError) {
    throw new ApiError(
      `Network error calling ${requestType} ${params.url}`,
      0,
    );
  }

  const json = await response.json().catch(() => null);

  if (response.status === 200) {
    return json as StandardResponse<T>;
  }

  if (response.status === 401) {
    MessageService.send({
      name: MessageNames.AUTH_ERROR_EVENT,
      value: 'auth error',
    });
    throw new ApiError('Authentication required', 401);
  }

  if (response.status === 422) {
    throw new ApiError(
      json?.message || 'Validation failed',
      422,
      json?.errors,
      json,
    );
  }

  // All other error codes (500, 403, etc.)
  throw new ApiError(
    json?.message || `API error: ${response.status}`,
    response.statusCode,
    undefined,
    json,
  );
}
```

### Key Changes
1. **Network errors caught** — `fetch()` wrapped in try-catch
2. **All non-200 responses throw** — callers always know about failures
3. **Custom `ApiError` class** — includes statusCode, validation errors
4. **422 errors carry validation details** — useful for form validation
5. **401 still fires MessageService** — for global auth redirect
6. **Generic return type** — `fetchData<T>` returns `StandardResponse<T>`
7. **No silent error swallowing** — production and dev behave the same

### Export the ApiError class

Add to the module exports so sagas can catch specific error types:
```typescript
export { ApiError };
```

## Consumer Impact

All saga files currently have no try-catch around API calls. After this change:
- Sagas that previously received `undefined` on error will now get a thrown `ApiError`
- This is actually BETTER because the saga will stop execution instead of continuing with undefined data
- Phase 3 (p3-saga-trycatch) will add proper try-catch to all sagas

## Validation

```bash
npm run checkTs   # Must pass
npm test          # Must pass (login saga test may need adjustment)
```
