// Package payment handles payment processing for the exchange platform,
// including:
//
//   - Fiat deposit and withdrawal processing
//   - Crypto payment creation and tracking
//   - Payment callback handling from external payment providers
//   - Payment status management (pending, confirmed, failed)
//
// Admin payment operations (callbacks, status updates) are accessed via
// admin API handlers on port 8001.
package payment
