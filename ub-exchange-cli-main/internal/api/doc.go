// Package api provides the HTTP server, route registration, and middleware
// for the exchange REST API. It runs on two ports:
//
//   - Port 8000: Public API endpoints (client-facing)
//   - Port 8001: Admin API endpoints (internal tools)
//
// Sub-packages:
//   - handler/: Public endpoint handlers (auth, orders, payments, etc.)
//   - adminhandler/: Admin endpoint handlers (order fulfillment, balance updates)
//   - middleware/: Authentication and recovery middleware
//
// All routes are registered under /api/v1/ using the Gin framework.
// Authentication is handled via JWT Bearer tokens in the Authorization header.
//
// Response format: { "status": bool, "message": string, "data": any }
package api
