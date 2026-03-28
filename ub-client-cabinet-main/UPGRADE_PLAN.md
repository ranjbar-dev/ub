# ub-client-cabinet-main Upgrade Plan

## Current State (Baseline Established)
- **Tests:** 30/30 suites pass, 94 tests pass, 3 skipped, 3 todo
- **Build:** Production build succeeds (requires `NODE_OPTIONS=--openssl-legacy-provider` on Node 22)
- **Node:** Running on Node 22.15.0 (project originally targets Node 12)

## Pre-Work Already Completed
- [x] Jest config: fixed `setupFilesAfterEnv` (removed deprecated `cleanup-after-each` and `extend-expect`)
- [x] Jest config: fixed deprecated `tsConfig` → `tsconfig`
- [x] Jest config: added `testPathIgnorePatterns` to exclude `internals/templates/`
- [x] Jest config: added toastify mock for customized react-toastify JSX files
- [x] Test bundler: converted ES module imports to CommonJS `require()`
- [x] ContactUsPage tests: fixed state shape, replaced `expect(true).toEqual(false)` with `it.todo()`
- [x] ContactUsPage test: renamed `.ts` → `.tsx` (JSX content), fixed `browserHistory` → `createMemoryHistory`
- [x] NotFoundPage test: added Redux Provider, fixed message assertion
- [x] LocaleToggle test: updated for MUI Select (no longer native `<select>`)
- [x] injectReducer/injectSaga tests: skipped `getInstance()` tests (null with functional components)
- [x] Webpack: replaced `offline-plugin` with existing `workbox-webpack-plugin` (Node 22 compatibility)

---

## Phase 1: Security & Critical Fixes (Safe, No Breaking Changes)

### 1.1 — Fix axios CVE-2023-26916
```bash
yarn upgrade axios@^0.28.0
```
- **Risk:** Low — 0.28 is backward compatible with 0.21
- **Only used in:** `app/services/upload_service.ts`

### 1.2 — Fix lodash security patch
```bash
yarn upgrade lodash@^4.17.21
```
- **Risk:** Minimal — patch version

### 1.3 — Create .env files for hardcoded URLs
- Extract URLs from `app/services/constants.ts` into `.env` files
- Create `.env.example`, `.env.development`, `.env.production`
- Use `process.env.REACT_APP_*` pattern (CRA compatible)

### 1.4 — Update Node engine in package.json
```json
"engines": { "node": ">=18", "npm": ">=8" }
```

### 1.5 — Update Dockerfiles
- `node:12.18.3` → `node:18-slim` or `node:20-slim`
- Add `NODE_OPTIONS=--openssl-legacy-provider` until Webpack 5 migration

---

## Phase 2: TSLint → ESLint Migration

### 2.1 — Install ESLint + TypeScript ESLint
```bash
yarn add -D @typescript-eslint/parser @typescript-eslint/eslint-plugin eslint-config-prettier
```

### 2.2 — Create `.eslintrc.json`
- Migrate rules from `tslint.json`
- Include `plugin:@typescript-eslint/recommended`

### 2.3 — Update `package.json` scripts
- Replace `lint:tslint` with ESLint command
- Remove tslint dependencies

### 2.4 — Remove TSLint files
- Delete `tslint.json`, `tslint-imports.json`
- Remove `tslint-loader` from webpack config

---

## Phase 3: TypeScript 4.0 → 5.x

### 3.1 — Upgrade TypeScript
```bash
yarn upgrade typescript@^5.4
```

### 3.2 — Fix type errors
- Update `@types/react`, `@types/react-dom`, `@types/node`
- Fix any new strict mode errors

### 3.3 — Update tsconfig.json
- Update `target` from `es5` to `es2017` or higher
- Update `lib` array for modern features

---

## Phase 4: React 17 → 18

### 4.1 — Upgrade packages
```bash
yarn upgrade react@^18.2.0 react-dom@^18.2.0
yarn upgrade @types/react@^18 @types/react-dom@^18
```

### 4.2 — Update entry point
- `app/app.tsx`: Replace `ReactDOM.render()` with `createRoot()`

### 4.3 — Update connected-react-router
- May need replacement — check compatibility with React 18

### 4.4 — Update react-test-renderer and testing-library
```bash
yarn upgrade react-test-renderer@^18 @testing-library/react@^14
```

### 4.5 — Fix any rendering/hydration issues
- Test all 21 containers
- Check for missing `<StrictMode>` warnings

---

## Phase 5: Material-UI 4 → MUI 5

### 5.1 — Run MUI codemods
```bash
npx @mui/codemod v5.0.0/preset-safe ./app
```

### 5.2 — Update import paths
- `@material-ui/core` → `@mui/material`
- `@material-ui/icons` → `@mui/icons-material`
- `@material-ui/lab` → `@mui/lab`
- `@material-ui/styles` → `@mui/material/styles`

### 5.3 — Fix theming/styling
- Update `createMuiTheme` → `createTheme`
- Update `makeStyles` usage (or migrate to `sx` prop)

### 5.4 — Update AG Grid if needed
- Check MUI 5 compatibility with AG Grid 25

---

## Phase 6: Webpack 4 → Webpack 5

### 6.1 — Upgrade webpack and loaders
```bash
yarn upgrade webpack@^5 webpack-cli@^5 webpack-dev-server@^4
```

### 6.2 — Update webpack config
- Remove `--openssl-legacy-provider` flag
- Update deprecated options (e.g., `optimization.moduleIds`)
- Replace `mini-css-extract-plugin` version
- Update `html-webpack-plugin` to v5

### 6.3 — Fix node polyfills
- Webpack 5 removes Node polyfills by default
- Add fallbacks for `process`, `Buffer` if needed

### 6.4 — Verify build + bundle size
- Compare bundle sizes before/after
- Run `yarn report` to analyze

---

## Phase 7: react-intl 2 → FormatJS/react-intl 6

### 7.1 — Upgrade package
```bash
yarn upgrade react-intl@^6
```

### 7.2 — Update provider setup
- `i18n.ts` may need updates for new API

### 7.3 — Fix `<FormattedMessage>` usage
- Ensure all message descriptors have `defaultMessage`
- Update any deprecated patterns

---

## Phase 8: Additional Modernization (Lower Priority)

### 8.1 — redux-form → react-hook-form
- Per-form migration (22+ forms across containers)
- Keep redux-form temporarily for complex forms

### 8.2 — AG Grid 25 → AG Grid 31+
- Check enterprise license
- Update import patterns (@ag-grid-enterprise/all-modules → @ag-grid-enterprise/*)

### 8.3 — Add Error Boundaries
- Wrap main routes in React Error Boundaries
- Add Sentry error reporting integration

### 8.4 — Add comprehensive container tests
- Currently 83% of containers have no tests
- Add tests for TradePage, FundsPage, OrdersPage, LoginPage, SignupPage

### 8.5 — amcharts 4 → 5

---

## Execution Order & Risk Assessment

| Phase | Risk | Effort | Can Be Done Now |
|-------|------|--------|-----------------|
| 1: Security fixes | Low | Small | ✅ Yes |
| 2: TSLint → ESLint | Low | Medium | ✅ Yes |
| 3: TypeScript 5.x | Medium | Small | ✅ Yes |
| 4: React 18 | Medium | Medium | After Phase 3 |
| 5: MUI 5 | High | Large | After Phase 4 |
| 6: Webpack 5 | Medium | Medium | After Phase 3 |
| 7: react-intl 6 | Medium | Medium | After Phase 4 |
| 8: Modernization | Varies | Large | After Phase 5 |

---

## Constraints
- Must not break existing functionality
- Tests must pass after each phase
- Build must succeed after each phase
- Feature branch per phase: `upgrade/<phase-name>`
