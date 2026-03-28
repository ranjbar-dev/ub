// Package adminhandler implements admin-only REST API endpoint handlers
// served on port 8001. These endpoints are called by the PHP backend
// (ub-server) for privileged operations:
//
//   - Order fulfillment (manual trade execution)
//   - Payment callbacks and status updates (deposit/withdrawal)
//   - User balance adjustments
//
// Admin endpoints use a separate authentication middleware that validates
// admin JWT tokens.
package adminhandler
