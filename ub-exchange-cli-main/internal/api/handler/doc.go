// Package handler implements the public-facing REST API endpoint handlers
// for the exchange platform. Each handler is a Gin handler factory that
// accepts a service dependency and returns a gin.HandlerFunc.
//
// Handlers follow a consistent pattern:
//  1. Bind and validate request parameters (JSON body or query)
//  2. Extract the authenticated user from Gin context (if required)
//  3. Delegate to the appropriate service method
//  4. Return the service response as JSON
//
// Endpoint groups:
//   - Auth: Login, register, 2FA, password reset, Centrifugo tokens
//   - Orders: Create, cancel, list open/traded orders, stop orders
//   - Payments: List payments, crypto payment creation
//   - User: Profile, KYC, device management, settings
//   - Market Data: Currencies, pairs, order book, trade book
//   - Balance: User balance queries
//   - Withdraw Address: CRUD for withdrawal addresses
package handler
