# Task: Standardize File Naming Convention

**ID:** p5-file-naming  
**Phase:** 5 — Code Organization  
**Severity:** 🟠 MEDIUM  
**Dependencies:** None  

## Problem

Inconsistent file naming across the project:
- Services: `snake_case` (e.g., `user_management_service.ts`, `api_service.ts`)
- Containers: `PascalCase` (e.g., `UserAccounts/`, `Billing/`)
- Components: mixed (`mainCat.tsx`, `SimpleGrid.tsx`, `sideNav/`)
- Utils: `camelCase` (e.g., `fileDownload.ts`, `formatters.ts`)

## Target Convention

Use the convention already dominant in each directory:
- **Containers:** `PascalCase` directories (already consistent — no change)
- **Services:** `camelCase` (rename from snake_case)
- **Utils:** `camelCase` (already mostly consistent)
- **Components:** `PascalCase` for component files, `camelCase` for utility files

## Files to Rename

### Services (snake_case → camelCase)
| Current | Target |
|---------|--------|
| `api_service.ts` | `apiService.ts` |
| `security_service.ts` | `securityService.ts` |
| `user_management_service.ts` | `userManagementService.ts` |
| `admin_reports_service.ts` | `adminReportsService.ts` |
| `external_orders_service.ts` | `externalOrdersService.ts` |
| `external_exchange_service.ts` | `externalExchangeService.ts` |
| `markets_service.ts` | `marketsService.ts` |
| `deposits_service.ts` | `depositsService.ts` |
| `withdrawals_service.ts` | `withdrawalsService.ts` |
| `billing_service.ts` | `billingService.ts` |
| `balances_service.ts` | `balancesService.ts` |
| `admins_service.ts` | `adminsService.ts` |
| `message_service.ts` | `messageService.ts` |
| `orders_service.ts` | `ordersService.ts` |

### Update ALL imports

After renaming each file, update every import across the codebase:
```bash
# Example: find all imports of the old name
grep -rn "from 'services/api_service'" src/ --include="*.ts" --include="*.tsx"
grep -rn "from 'services/user_management_service'" src/ --include="*.ts" --include="*.tsx"
# ... etc for each service
```

## Execution Steps

⚠️ **This is a high-risk refactor.** Do one file at a time and verify after each:

1. Rename the file (git mv if in git)
2. Find all imports of the old name
3. Update all imports to new name
4. Run `npm run checkTs`
5. Repeat for next file

## Alternative: Configure ESLint Rule

Instead of renaming, add an ESLint rule to enforce naming going forward:
```json
{
  "rules": {
    "unicorn/filename-case": ["error", { "case": "camelCase" }]
  }
}
```

## Validation

```bash
npm run checkTs   # Must pass after each rename
npm test          # Must pass
npm run build     # Must build
```
