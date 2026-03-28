# Task: Test Formatter Utilities

## Goal
Create `src/utils/__tests__/formatters.test.ts` with comprehensive tests for all formatter functions.

## Context
- File to test: `src/utils/formatters.ts`
- Functions: queryStringer, CurrencyFormater, Format, FormatDate, DatePrefixer, censor, safeFinancialAdd
- safeFinancialAdd is NEW (Sprint 2) and uses fixed-point math to avoid float errors
- queryStringer was updated in Sprint 1 to use encodeURIComponent
- No existing tests for any of these functions

## File to Create: `src/utils/__tests__/formatters.test.ts`

### Test Cases Required

#### queryStringer
- Empty params object → empty string
- Single param → `?key=value`
- Multiple params → `?key1=value1&key2=value2`
- Special characters encoded → `encodeURIComponent` applied to keys and values
- Null/undefined values skipped (check actual implementation behavior)

#### safeFinancialAdd (CRITICAL — financial precision)
- `safeFinancialAdd('0.1', '0.2')` → `'0.30000000'` (NOT `'0.30000000000000004'`)
- `safeFinancialAdd('1', '2')` → `'3.00000000'`
- `safeFinancialAdd('0', '0')` → `'0.00000000'`
- `safeFinancialAdd('999999.99999999', '0.00000001')` → `'1000000.00000000'`
- Custom decimals: `safeFinancialAdd('0.1', '0.2', 2)` → `'0.30'`
- String number inputs: `safeFinancialAdd('100.50', '200.25')` → `'300.75000000'`
- Negative numbers if applicable

#### CurrencyFormater
- Zero → appropriate output
- Positive number → comma-separated
- Number with decimals → preserves decimals
- Very large number → proper comma formatting

#### Format
- Positive number → German locale format
- Zero → '0'
- Negative number → formatted correctly

#### FormatDate
- 'YYYYMMDD' → 'YYYY-MM-DD'
- Edge cases (null, empty string)

#### DatePrefixer
- Single digit (5) → '05'
- Double digit (12) → '12'

#### censor
- String → masked characters (check implementation for exact behavior)

## Validation
- Run: `npx react-scripts test --watchAll=false --testPathPattern="formatters.test" --verbose`
- All tests pass
