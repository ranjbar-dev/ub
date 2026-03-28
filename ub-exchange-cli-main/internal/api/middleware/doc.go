// Package middleware provides Gin middleware for the exchange REST API:
//
//   - AuthMiddleware: Validates JWT Bearer tokens and sets the authenticated
//     user in the Gin context under the UserKey. Returns 401 on invalid tokens.
//
//   - AdminAuthMiddleware: Validates admin-level JWT tokens for the admin API
//     (port 8001). Similar to AuthMiddleware but checks admin privileges.
//
//   - NonRequiredAuthMiddleware: Marks a route as optionally authenticated.
//     If a valid token is present the user is set in context; otherwise the
//     request proceeds with a nil user.
//
// Context keys:
//   - UserKey: Stores *user.User for the authenticated user
//   - NonRequiredAuthKey: Flag indicating optional auth mode
package middleware
