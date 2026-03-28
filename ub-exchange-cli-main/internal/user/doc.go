// Package user manages user accounts and profiles for the exchange platform:
//
//   - User CRUD operations (create, read, update)
//   - User profile management (name, contact, KYC documents)
//   - KYC level calculation based on trading volume and verification status
//   - Permission management (feature access control)
//   - User configuration preferences
//
// User entities are the central identity model referenced by orders, balances,
// payments, and authentication throughout the system.
package user
