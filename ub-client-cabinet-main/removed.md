# Removed Files Log — ub-client-cabinet-main

This file tracks files that were removed or reorganized during project cleanup.

## Removed Files

| File | Date Removed | Reason |
|------|-------------|--------|
| `typescript-5.4.5.tgz` | 2026-03-30 | Bundled npm package (5.56 MB) — should use npm registry |
| `node_modules_fresh/` | 2026-03-30 | Debugging copy of node_modules — not needed |
| `node_modules_new/` | 2026-03-30 | Debugging copy of node_modules — not needed |
| `tslint.json` | 2026-03-30 | Deprecated — ESLint is now primary linter (TSLint migration complete) |
| `tslint-imports.json` | 2026-03-30 | Deprecated — part of TSLint (no longer used) |
| `yarn-error.log` | 2026-03-30 | Stale error log from 2024 |
| `unused.md` | 2026-03-30 | Legacy file tracking from old repository — no ongoing value |
| `dist/` | 2026-03-30 | Unused alternative build output directory |
| `stats.json` | 2026-03-30 | Webpack bundle analysis artifact — auto-generated |
| `tsconfig.tsbuildinfo` | 2026-03-30 | TypeScript incremental build cache — auto-generated |

## Reorganized Files

| File | Moved From | Moved To | Reason |
|------|-----------|----------|--------|
| `version.sh` | root | `scripts/version.sh` | Scripts organized in scripts/ |
