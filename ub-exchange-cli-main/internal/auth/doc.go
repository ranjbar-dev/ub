// Package auth implements authentication and authorization for the exchange
// platform. It handles:
//
//   - User login/registration with email and phone
//   - JWT token issuance and validation
//   - Two-factor authentication (TOTP via pquerna/otp)
//   - Password reset flows
//   - MQTT ACL authorization (controls which topics users can subscribe to)
//   - Device fingerprinting for login security
//
// The auth service is consumed by API handlers (internal/api/handler/auth.go)
// and middleware (internal/api/middleware/).
package auth
