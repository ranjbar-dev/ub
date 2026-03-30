// Package auth implements authentication and authorization for the exchange
// platform. It handles:
//
//   - User login/registration with email and phone
//   - JWT token issuance and validation
//   - Two-factor authentication (TOTP via pquerna/otp)
//   - Password reset flows
//   - Centrifugo token generation (JWT for real-time connections)
//   - Device fingerprinting for login security
//
// The auth service is consumed by API handlers (internal/api/handler/auth.go)
// and middleware (internal/api/middleware/).
package auth
