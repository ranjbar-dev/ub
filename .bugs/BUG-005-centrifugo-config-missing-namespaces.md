# BUG-005: Centrifugo `user` Namespace Missing Authentication Requirement

## Severity: **HIGH** — Private user data may be accessible without authorization

## Files
- `centrifugo-config.json` (root)
- `ub-server-main/centrifugo-config.json`

## Issue
The `user` namespace in `centrifugo-config.json` has `allow_subscribe_for_client: true` but does NOT enforce subscription token validation. Current config:

```json
{
  "name": "user",
  "presence": true,
  "history_size": 5,
  "history_ttl": "300s",
  "allow_subscribe_for_client": true
}
```

This means **any authenticated connection** (with a valid connection token) can subscribe to ANY user's private channels (`user:{channel}:open-orders`, `user:{channel}:crypto-payments`) without a subscription-specific token.

The old EMQX setup had ACL rules that verified the user's `privateChannelName` matched their JWT. The new Centrifugo config does not replicate this.

## Fix
The `user` namespace should require subscription tokens:
```json
{
  "name": "user",
  "presence": true,
  "history_size": 5,
  "history_ttl": "300s",
  "allow_subscribe_for_client": false,
  "allow_subscribe_for_client_with_token": true
}
```

This forces clients to obtain a subscription token (from `/api/v1/auth/centrifugo-subscribe-token`) that includes the specific channel name, preventing unauthorized cross-user access.

## Impact
- **Security vulnerability**: User A could subscribe to User B's open-orders and crypto-payments channels
- Cannot be verified without running Docker services, but config analysis confirms the gap
