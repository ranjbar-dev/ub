# UnitedBit Exchange — Copilot Instructions

> Monorepo root instructions. Sub-projects have their own in
> `ub-server-main/.github/copilot-instructions.md` and
> `ub-exchange-cli-main/.github/copilot-instructions.md`.

## Config File Propagation

Go services (exchange-cli, communicator) have **3 config variants** that must stay in sync:

| File | Used when |
|------|-----------|
| `config.yaml` | Local development |
| `config.docker.yaml` | Docker — **copied over** `config.yaml` at container start |
| `config_test.yaml` | Test suite |

When changing any config value, **update all three files**. Docker compose copies
`config.docker.yaml` over `config.yaml` at runtime, so changes to `config.yaml`
alone are invisible in Docker deployments.

PHP equivalent: `parameters.yml.dist` (template) → `parameters.docker.yml`
(Docker override, mounted as `parameters.yml`). Update both.

## Symfony Cache After Structural Changes

After renaming, moving, or deleting PHP classes in `ub-server-main/`:

```powershell
Remove-Item -Recurse -Force ub-server-main\var\cache -ErrorAction SilentlyContinue
```

Symfony auto-generates container files with hardcoded class paths. Stale cache
will reference old class names and crash the application. `var/cache/` is
gitignored — deleting it is always safe.

## Cross-Service API Verification

This platform has 4 consumers for most API endpoints: PHP backend, Go backend,
React client, Flutter app. When modifying any endpoint, verify **all of**:

1. **Route path** — same across PHP controller annotation, Go router registration, frontend service URL
2. **HTTP method** — PHP `methods={}`, Go `.GET()`/`.POST()`, frontend `fetch()`/`dio.get()`/`dio.post()`
3. **Request/response shape** — both backends return `{ status, message, data }`

Key files to cross-check:
- PHP routes: `ub-server-main/src/Exchange/ApiBundle/Controller/V1/`
- Go routes: `ub-exchange-cli-main/internal/api/routes.go`
- React services: `ub-client-cabinet-main/app/services/`
- Flutter providers: `ub-app-main/lib/app/modules/*/providers/`

## Submodule Workflow

`ub-server-main` is a **git submodule**. After committing inside it:

```powershell
cd C:\Users\root\Desktop\Projects\github\ub
git add ub-server-main
git commit -m "chore: update ub-server-main submodule (description)"
```

The outer repo tracks a pointer to a specific submodule commit. Without
`git add ub-server-main`, the outer repo shows the submodule as "dirty"
even though inner commits are complete.

## Quick Validation Commands

### PHP (ub-server-main)
```powershell
# Syntax check a single file
php -l ub-server-main\src\Exchange\<Bundle>\Services\<File>.php

# Syntax check all PHP files
Get-ChildItem ub-server-main\src -Recurse -Filter *.php | ForEach-Object { php -l $_.FullName } 2>&1 | Select-String "error"

# Run tests (requires Docker)
docker exec -u exchange exchange-app bash -c "cd /home/exchange/project && php vendor/bin/codecept run unit"
```

### Go (ub-exchange-cli-main)
```powershell
cd ub-exchange-cli-main
go build ./cmd/exchange-cli/
go build ./cmd/exchange-httpd/
go vet ./...
```

### Go (ub-communicator-main)
```powershell
cd ub-communicator-main
go build -mod=vendor ./cmd/rabbit-consumer/
```

### TypeScript (ub-admin-main, ub-client-cabinet-main)
```powershell
cd ub-admin-main && npm run lint
cd ub-client-cabinet-main && npm run lint
```

### Flutter (ub-app-main)
```powershell
# Requires Dart SDK 2.12-2.x (NOT 3.x — pre-null-safety project)
cd ub-app-main && flutter analyze
```
