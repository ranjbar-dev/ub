# BUG-006: Centrifugo API Endpoint Path May Be Incorrect for v5+

## Severity: **MEDIUM** — API calls may fail depending on Centrifugo version

## Files
- `ub-exchange-cli-main/internal/platform/centrifugo.go` — uses `/api/publish` and `/api/broadcast`
- `ub-server-main/src/Exchange/CoreBundle/Services/CentrifugoClient.php` — uses phpcent library

## Issue
The Go CentrifugoClient constructs API endpoints as:
```go
c.baseURL + "/api/publish"   // e.g., http://centrifugo:8000/api/publish
c.baseURL + "/api/broadcast" // e.g., http://centrifugo:8000/api/broadcast
```

But if `baseURL` is already set to `http://centrifugo:8000/api` (as configured in config.yaml), the actual URL becomes:
```
http://centrifugo:8000/api/api/publish  ← DOUBLE /api/
```

This depends on whether `centrifugo.api_url` in config is `http://centrifugo:8000` or `http://centrifugo:8000/api`.

Current config values:
- Go `config.yaml`: needs verification
- PHP `parameters.yml.dist`: `centrifugo.api_url: "http://centrifugo:8000/api"`

If PHP's `api_url` includes `/api` and Go's `doRequest` also appends `/api/publish`, the paths will differ between services.

## Verification Needed
1. Check Go config value for `centrifugo.api_url`
2. Verify phpcent library's expected URL format
3. Ensure both services construct the same final API URL

## Impact
- One or both services may fail to publish to Centrifugo
- Cannot be verified without running Docker services
