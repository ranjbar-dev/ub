// Package command implements the 16 CLI commands available in the exchange-cli
// binary. Each command implements the ConsoleCommand interface:
//
//	type ConsoleCommand interface {
//	    Run(ctx context.Context, flags []string)
//	}
//
// Commands are registered by name in cmd/exchange-cli/main.go and invoked
// via command-line arguments. Most are designed to run as cron jobs:
//
//   - set-user-level: Calculate user trading levels from volume (daily)
//   - initialize-balance: Create balances for new coins (manual)
//   - generate-address: Generate wallet addresses (manual)
//   - retrieve-open-orders: Cache open orders to Redis (every 15min)
//   - submit-bot-orders: Execute bot orders to Binance (every 1min)
//   - sync-kline: Sync OHLC candle data (every 1min)
//   - check-withdrawals: Verify withdrawal status (every 10min)
//   - update-orders-from-external: Sync external order status (every 30min)
//   - retrieve-external-orders: Cache external orders to Redis (periodic)
//   - delete-cache: Clear Redis cache (manual)
//   - ub-captcha-*: Captcha key generation and encryption utilities (manual)
//   - ws-health-check: WebSocket connection health monitoring
//
// All commands obtain dependencies from the DI container (internal/di).
package command
