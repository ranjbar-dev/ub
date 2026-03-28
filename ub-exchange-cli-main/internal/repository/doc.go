// Package repository provides GORM-based data access repositories for all
// domain entities. Each repository is a thin wrapper around *gorm.DB that
// encapsulates SQL queries for a single entity type.
//
// Naming convention: <Entity>Repository (e.g., OrderRepository, UserRepository).
// Constructor convention: New<Entity>Repository(db *gorm.DB).
//
// Repositories are registered in the DI container (internal/di) and consumed
// by service layers via their interface types defined in each domain package
// (e.g., order.Repository, user.Repository).
//
// All financial values use shopspring/decimal — never float64.
// Transactions are managed by the caller using *gorm.DB.Begin().
package repository
