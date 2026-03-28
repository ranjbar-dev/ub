// Package userbalance manages cryptocurrency and fiat balances for exchange
// users. It handles:
//
//   - Balance queries (available, frozen, total per currency)
//   - Balance updates during trade settlement (atomic credit/debit)
//   - Balance initialization for new currencies
//   - Balance synchronization with external wallet services
//
// All balance operations use shopspring/decimal for precision arithmetic.
// Balance updates during trade matching are performed within database
// transactions to ensure consistency.
package userbalance
