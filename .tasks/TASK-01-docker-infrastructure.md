# TASK-01: Replace EMQX Docker Container with Centrifugo

## Objective
Remove all EMQX/emqtt Docker service definitions and replace with Centrifugo across all docker-compose files. Create a Centrifugo config file.

## Files to Modify

### 1. `docker-compose.dev.yml` (Root)
- **Current**: Contains `emqtt` service using `emqx/emqx:v3.0.0` image
- **Current config**:
  ```yaml
  emqtt:
    image: emqx/emqx:v3.0.0
    environment:
      EMQX_LOADED_PLUGINS: "emqx_management,emqx_auth_http,emqx_recon,emqx_retainer,emqx_dashboard"
      EMQX_AUTH__HTTP__AUTH_REQ: "http://nginx/api/v1/emqtt/login"
      EMQX_AUTH__HTTP__SUPER_REQ: "http://nginx/api/v1/emqtt/superuser"
      EMQX_AUTH__HTTP__ACL_REQ: "http://nginx/api/v1/emqtt/acl"
    ports:
      - "1883:1883"
      - "8083:8083"
  ```
- **Action**: Replace `emqtt` service with `centrifugo` service:
  ```yaml
  centrifugo:
    image: centrifugo/centrifugo:v6
    container_name: centrifugo
    command: centrifugo -c config.json
    ports:
      - "8800:8000"
    environment:
      - CENTRIFUGO_TOKEN_HMAC_SECRET_KEY=${CENTRIFUGO_TOKEN_HMAC_SECRET_KEY:-your-secret-key}
      - CENTRIFUGO_API_KEY=${CENTRIFUGO_API_KEY:-your-api-key}
      - CENTRIFUGO_ADMIN=true
      - CENTRIFUGO_ADMIN_PASSWORD=${CENTRIFUGO_ADMIN_PASSWORD:-admin}
      - CENTRIFUGO_ADMIN_SECRET=${CENTRIFUGO_ADMIN_SECRET:-admin-secret}
    volumes:
      - ./centrifugo-config.json:/centrifugo/config.json
    networks:
      - candle-http
    restart: unless-stopped
  ```

### 2. `ub-server-main/docker-compose.yml` (Local dev)
- **Current**: Contains `emqtt` service with EMQX v3.0.0, auth webhooks, ports 1883/8083
- **Action**: Replace `emqtt` service with `centrifugo` service (same as above). Remove EMQX auth webhook environment variables. Keep same network.

### 3. `ub-server-main/docker-compose-dev.yml` (Dev server)
- **Current**: Contains `emqtt` service with EMQX v4.0.0
- **Action**: Replace with `centrifugo` service

### 4. `ub-server-main/docker-compose-prod.yml` (Production)
- **Current**: Contains `emqtt` service with EMQX v3.0.0 and HTTPS auth callbacks
- **Action**: Replace with `centrifugo` service with production-ready environment variables

### 5. Create `centrifugo-config.json` (Root level)
- **Action**: Create new file:
  ```json
  {
    "client": {
      "token": {
        "hmac_secret_key": "your-secret-key"
      },
      "allowed_origins": [
        "http://localhost:3000",
        "http://localhost:3001",
        "https://app.unitedbit.com",
        "https://admin.unitedbit.com"
      ]
    },
    "http_api": {
      "key": "your-api-key"
    },
    "admin": {
      "enabled": true,
      "password": "admin",
      "secret": "admin-secret"
    },
    "channel_namespaces": [
      {
        "name": "trade",
        "presence": false,
        "history_size": 1,
        "history_ttl": "60s",
        "allow_subscribe_for_client": true
      },
      {
        "name": "user",
        "presence": true,
        "history_size": 5,
        "history_ttl": "300s",
        "allow_subscribe_for_client": true
      }
    ],
    "log_level": "info"
  }
  ```

### 6. Also create `ub-server-main/centrifugo-config.json` if docker-compose files there reference it locally

## Channel Naming Convention (for reference)
- MQTT `main/trade/ticker/{pair}` → Centrifugo `trade:ticker:{pair}`
- MQTT `main/trade/order-book/{pair}` → Centrifugo `trade:order-book:{pair}`
- MQTT `main/trade/trade-book/{pair}` → Centrifugo `trade:trade-book:{pair}`
- MQTT `main/trade/kline/{timeFrame}/{pair}` → Centrifugo `trade:kline:{timeFrame}:{pair}`
- MQTT `main/trade/chart/{timeFrame}/{pair}` → Centrifugo `trade:chart:{timeFrame}:{pair}`
- MQTT `main/trade/user/{channel}/open-orders/` → Centrifugo `user:{channel}:open-orders`
- MQTT `main/trade/user/{channel}/crypto-payments/` → Centrifugo `user:{channel}:crypto-payments`

## Nginx Config Notes
- If any nginx config proxies EMQX WebSocket (port 8083/8443), update to proxy Centrifugo WebSocket at `/connection/websocket` on port 8000
- Look in `ub-server-main/.docker/nginx/` for nginx configs

## Validation
- `docker-compose config` should pass on all modified files
- Centrifugo container should start and respond to `curl http://localhost:8800/health`

## Reference Docs
- See `.docs/centrifugo-docker-setup.md` for full Docker setup reference
- See `.docs/centrifugo-overview.md` for architecture overview
