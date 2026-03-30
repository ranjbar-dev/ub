# BUG-010: Centrifugo Token Endpoint Routes Don't Match Frontend Expectations

## Severity: **MEDIUM** — Frontend may call wrong endpoint URLs

## Files
- `ub-server-main/src/Exchange/ApiBundle/Controller/V1/CentrifugoTokenController.php`
- `ub-client-cabinet-main/app/services/CentrifugoAuthService.ts` (verify endpoint URL)
- `ub-exchange-cli-main/internal/api/handler/centrifugo.go` (verify endpoint URL)

## Issue
The PHP CentrifugoTokenController defines routes:
- `POST /auth/centrifugo-connection-token` (name: `api_centrifugo_connection_token`)
- `GET /auth/centrifugo-subscribe-token` (name: `api_centrifugo_subscribe_token`)

The routing.yml registers this controller under prefix `/auth`:
```yaml
exchange_api_centrifugo_token:
  resource: "@ExchangeApiBundle/Controller/V1/CentrifugoTokenController.php"
  type: annotation
  prefix: /auth
```

So full paths are: `/api/v1/auth/centrifugo-connection-token` and `/api/v1/auth/centrifugo-subscribe-token`.

Meanwhile, the Go exchange-cli also likely has Centrifugo token endpoints. Both PHP and Go serve on different ports — need to verify:
1. Which service do frontends call for tokens?
2. Do the endpoint paths match between what PHP serves and what frontends expect?
3. Are the Go token endpoints registered and accessible?

## Verification Needed
- Check React CentrifugoAuthService for the token endpoint URL it calls
- Check Flutter CentrifugoService for the token endpoint URL it calls
- Verify Go handler routes match
- Confirm nginx routes token requests to the correct backend

## Impact
- Authenticated users may not be able to get Centrifugo tokens
- Private channels (open-orders, crypto-payments) will not work without tokens
