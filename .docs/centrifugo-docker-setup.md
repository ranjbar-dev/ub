# Centrifugo Docker Setup Reference

## Docker Compose Service Definition

```yaml
centrifugo:
  image: centrifugo/centrifugo:v6
  container_name: centrifugo
  command: centrifugo -c config.json
  ports:
    - "8800:8000"    # Centrifugo HTTP API + WebSocket
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

## Centrifugo Configuration File (`centrifugo-config.json`)

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
      "presence": true,
      "history_size": 10,
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

## Network Connectivity

Centrifugo replaces EMQX on the `candle-http` network. Services that previously connected to EMQX on port 1883/8083 now connect to Centrifugo on port 8000 via HTTP API.

### Port Mapping

| Old (EMQX) | New (Centrifugo) | Purpose |
|---|---|---|
| 1883 (MQTT) | 8000 (HTTP API) | Backend publishing |
| 8083 (MQTT WebSocket) | 8000 (WebSocket) | Client connections |
| 8443 (MQTT WSS via nginx) | 8800 (WSS via nginx) | Client connections (production) |

### Nginx Configuration Update

Replace EMQX WebSocket proxy with Centrifugo WebSocket proxy:

```nginx
# Old EMQX config
# location /mqtt {
#     proxy_pass http://emqtt:8083/mqtt;
#     proxy_http_version 1.1;
#     proxy_set_header Upgrade $http_upgrade;
#     proxy_set_header Connection "upgrade";
# }

# New Centrifugo config
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

## Environment Variables

| Variable | Description | Default |
|---|---|---|
| `CENTRIFUGO_TOKEN_HMAC_SECRET_KEY` | JWT signing secret (must match backend) | Required |
| `CENTRIFUGO_API_KEY` | HTTP API authentication key | Required |
| `CENTRIFUGO_ADMIN` | Enable admin web UI | `true` |
| `CENTRIFUGO_ADMIN_PASSWORD` | Admin UI password | Required |
| `CENTRIFUGO_ADMIN_SECRET` | Admin UI secret | Required |

## Health Check

```bash
curl http://localhost:8800/health
```
