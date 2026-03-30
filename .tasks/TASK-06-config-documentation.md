# TASK-06: Update Configuration & Documentation

## Objective
Update all configuration files, documentation, and nginx configs to reflect the EMQX → Centrifugo migration.

## Files to Modify

### 1. UPDATE: Root `AGENTS.md`
- Replace all EMQX/MQTT references with Centrifugo
- Update architecture diagram (EMQX → Centrifugo)
- Update MQTT topic references → Centrifugo channel names
- Update port mapping (1883/8083/8443 → 8000/8800)
- Update authentication flow description
- Update Gotchas section (remove EMQX-specific gotchas, add Centrifugo notes)
- Keep the section structure but update content

### 2. UPDATE: Root `README.md`
- Replace EMQX references with Centrifugo

### 3. UPDATE: `ub-server-main/AGENTS.md`
- Remove EMQX section
- Add Centrifugo section
- Update real-time data flow description
- Update credential references (remove mqtt_abbas, add centrifugo API key)

### 4. UPDATE: `ub-exchange-cli-main/AGENTS.md`
- Remove MQTT auth endpoint references
- Replace MQTT broker references with Centrifugo
- Update topic/channel descriptions

### 5. UPDATE: `ub-exchange-cli-main/ARCHITECTURE.md`
- Update MQTT pub source → Centrifugo HTTP API publish
- Update topic references → channel references

### 6. UPDATE: `ub-app-main/AGENTS.md`
- Replace mqtt_client reference with centrifuge-dart
- Update topic subscriptions → channel subscriptions

### 7. UPDATE: Nginx configs (if present)
- Check `ub-server-main/.docker/nginx/` for any nginx configuration files
- If EMQX WebSocket proxy exists (port 8083/8443), replace with Centrifugo WebSocket proxy:
  ```nginx
  location /connection/websocket {
      proxy_pass http://centrifugo:8000/connection/websocket;
      proxy_http_version 1.1;
      proxy_set_header Upgrade $http_upgrade;
      proxy_set_header Connection "upgrade";
      proxy_set_header X-Real-IP $remote_addr;
      proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
      proxy_read_timeout 3600s;
  }
  ```
- Remove any EMQX-specific proxy rules (port 1883, 8083)
- Update SSL/WSS configuration if port 8443 was used for EMQX WSS

### 8. UPDATE: `.env` files and `.env.example` files
- Remove EMQX-related environment variables
- Add Centrifugo environment variables:
  ```
  CENTRIFUGO_TOKEN_HMAC_SECRET_KEY=your-secret-key
  CENTRIFUGO_API_KEY=your-api-key
  CENTRIFUGO_ADMIN_PASSWORD=admin
  CENTRIFUGO_ADMIN_SECRET=admin-secret
  ```

### 9. UPDATE: Frontend `.env` files
- `ub-client-cabinet-main/.env` — Update WebSocket URL (remove port 8443 MQTT WSS)
- `ub-app-main/.env` or equivalent — Update WebSocket URL

## Key Search Patterns
```bash
grep -r "emqx\|emqtt\|EMQX\|EMQTT\|mqtt\|MQTT" --include="*.md" --include="*.yml" --include="*.yaml" --include="*.env*" --include="*.nginx" --include="*.conf" .
```

## Reference
- See `.docs/centrifugo-overview.md` for the full architecture and naming convention
- See `.docs/centrifugo-docker-setup.md` for nginx config reference
