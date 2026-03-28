// Package jwt provides JWT token helper utilities for creating and parsing
// JSON Web Tokens used in API authentication. It wraps golang-jwt/jwt/v5
// with exchange-specific claims and key management.
//
// JWT keys are loaded from config/jwt/ (RSA public/private key pair).
package jwt
