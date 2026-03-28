# Task: Add JSDoc to Shared Components

**ID:** p4-component-jsdoc  
**Phase:** 4 — Component Props & Documentation  
**Severity:** 🟠 MEDIUM  
**Dependencies:** None  

## Problem

No shared component has JSDoc documentation. AI agents can't determine the purpose or usage of components without reading the entire implementation.

## Scope

All components in `src/app/components/` directory.

## Pattern

### Before
```typescript
const SimpleGrid = (props) => {
  // 200 lines of AG Grid config...
};
```

### After
```typescript
/**
 * Wrapper around AG Grid providing standard admin panel data table features.
 * Handles pagination, sorting, filtering, and row selection.
 *
 * @example
 * ```tsx
 * <SimpleGrid
 *   columnDefs={columns}
 *   rowData={users}
 *   onGridReady={handleGridReady}
 *   pagination={true}
 * />
 * ```
 */
const SimpleGrid: React.FC<SimpleGridProps> = (props) => {
  // ...
};
```

### JSDoc Template for Components

```typescript
/**
 * [One-line description of what the component does].
 * [Additional context about when/where it's used].
 *
 * @example
 * ```tsx
 * <ComponentName prop1="value" prop2={data} />
 * ```
 */
```

## Components to Document

| Component | Description |
|-----------|-------------|
| SimpleGrid | AG Grid wrapper for data tables |
| sideNav | Left sidebar navigation with route links |
| mainCat | Main category navigation items |
| DateFilter | Date range picker for grid filtering |
| RejectPopup | Confirmation modal for rejection actions |
| ConfirmPopup | Confirmation modal for approval actions |
| Toast | Notification snackbar component |
| WindowWrapper | Floating window/panel container |
| PriceInput | Numeric input formatted for currency values |
| StatusBadge | Colored status indicator chip |

## Execution Steps

1. List all component files: `find src/app/components/ -name "*.tsx"`
2. For each component file:
   a. Read the component to understand its purpose
   b. Add JSDoc block above the component definition
   c. Include at least one `@example` usage
3. Run `npm run checkTs` to verify no issues

## Validation

```bash
npm run checkTs   # Must pass
```
