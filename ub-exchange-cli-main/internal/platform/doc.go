// Package platform provides infrastructure abstractions that decouple business
// logic from external systems. Each file defines an interface and its concrete
// implementation for a single infrastructure concern:
//
//   - RedisClient: Redis key-value and sorted-set operations
//   - Cache: Application-level caching with TTL (backed by go-redis/cache)
//   - Configs: Configuration access via Viper
//   - Logger: Structured logging via uber/zap with Sentry integration
//   - HTTPClient: Outbound HTTP requests
//   - JwtHandler: JWT token creation and validation
//   - PasswordEncoder: bcrypt password hashing
//   - MqttClient: MQTT publish/subscribe for real-time data push
//   - RabbitMqClient: RabbitMQ message publishing for async communication
//   - WsClient: WebSocket connection management
//
// All concrete implementations are registered in the DI container (internal/di)
// and should be consumed via their interface types, never directly.
package platform
