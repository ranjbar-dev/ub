# Task 010: Split sync_kline.go Command

## Priority: LOW
## Risk: LOW (internal refactor of one command)
## Estimated Scope: 1 file → 3-4 files

---

## Problem

`internal/command/sync_kline.go` is 407 lines — the largest CLI command file. It mixes:
1. Flag parsing and parameter setup (~100 lines)
2. Data fetching from external exchange (~100 lines)
3. Kline transformation and normalization (~100 lines)
4. Database persistence and queue publishing (~100 lines)

## Current Structure

```
sync_kline.go (407 lines)
├── syncKlineCmd struct (line 19)
├── Run(ctx, flags) — main entry point (line 40)
├── getParametersFromKlineSync() — DB parameter loading (line 141)
├── setNeededData(flags) — flag parsing (line 193)
├── fetchAndSaveKlines(params) — fetch + transform + save loop (line 264)
└── NewSyncKlineCmd() — constructor (line 397)
```

## Goal

Split into cohesive files by responsibility while keeping them in the `command` package.

## Implementation Plan

### File 1: `sync_kline.go` (keep — struct + constructor + Run)

Keep:
- `syncKlineCmd` struct definition
- `NewSyncKlineCmd()` constructor
- `Run()` method (orchestrates the steps)

### File 2: `sync_kline_params.go` (parameter handling)

Move:
- `setNeededData(flags)` — flag parsing and validation
- `getParametersFromKlineSync()` — load sync parameters from DB
- Any helper types for parameters (e.g., `neededParams` struct if it exists)

### File 3: `sync_kline_fetch.go` (data fetching and persistence)

Move:
- `fetchAndSaveKlines(params)` — the main fetch/transform/save loop
- Any helper functions for kline transformation

## Verification

```bash
go build ./internal/command/
go build ./...
```

## Notes

- All methods are on `*syncKlineCmd` receiver — same package split works naturally
- The `neededParams` type (if defined locally) needs to stay accessible to all files — keep it in the main file or move it to the params file
- This is the only CLI command over 300 lines; others are fine as single files
