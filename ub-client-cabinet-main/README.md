# ub-client-cabinet-main — Client Trading Dashboard

A single-page client trading dashboard for the UnitedBit cryptocurrency exchange platform. Users can trade currency pairs, manage funds (deposits/withdrawals), view order history, complete KYC verification, and manage account security — all with real-time streaming data.

## Prerequisites

| Requirement | Version |
|-------------|---------|
| Node.js | >=18 (`.nvmrc`: lts/hydrogen) |
| Yarn | 1.22.x (`corepack enable && corepack prepare yarn@1.22.22`) |
| npm | >=8 |

> **Note:** Webpack 4 requires `NODE_OPTIONS=--openssl-legacy-provider` on Node 18+. The Docker images set this automatically.

## Quick Start

```bash
# 1. Install dependencies
yarn install

# 2. Start development server (localhost with local flag)
yarn start
# → Opens http://localhost:3000 with hot reload
# → Uses dev-app.unitedbit.com APIs (IS_LOCAL=true)

# 3. Or start pointing to remote dev server
yarn start-dev
```

## Build Commands

| Command | Purpose |
|---------|---------|
| `yarn start` | Dev server (localhost, `IS_LOCAL=true`) |
| `yarn start-dev` | Dev server (remote dev APIs) |
| `yarn build` | Production build (no sourcemaps) |
| `yarn build-dev` | Production build with dev APIs |
| `yarn start:prod` | Serve production build (Express, port 3000) |
| `yarn test` | Run Jest tests with 98% coverage threshold |
| `yarn test:watch` | Jest in watch mode |
| `yarn lint` | ESLint + stylelint + TSLint |
| `yarn typecheck` | TypeScript type checking (`tsc --noEmit`) |
| `yarn report` | Webpack bundle analyzer (port 4200) |
| `yarn generate` | Plop code generator (new container/component) |
| `yarn extract-intl` | Extract i18n messages from source |

## Docker

```bash
# Development image (dev APIs)
docker build -f Dockerfile -t ub-cabinet-dev .

# Production image
docker build -f DockerfileProd -t ub-cabinet-prod .

# Run
docker run -p 3000:3000 ub-cabinet-prod
```

Both images use `node:18-slim` with `NODE_OPTIONS=--openssl-legacy-provider`.

## Architecture Overview

```
React 18 + TypeScript 5.4.5
├── Redux 4 (dynamic reducer/saga injection per container)
├── Redux-Saga 1.1.3 (API side effects)
├── RxJS 6 MessageService (90+ cross-component event types)
├── Centrifugo WebSocket (real-time market data: ticker, order book, trades, kline)
├── Material-UI 4 + styled-components 5
├── AG Grid Enterprise 25 (data tables)
├── TradingView charting_library (vendored)
└── react-intl 2 (i18n: en, de, ar with RTL)
```

**21 containers** (page-level) | **38 components** (reusable UI) | **17 services** (API/business logic)

Each container follows a strict pattern: `index.tsx`, `reducer.ts`, `saga.ts`, `actions.ts`, `constants.ts`, `selectors.ts`, `types.d.ts`, `Loadable.tsx`.

## Key URLs

| Environment | App URL | API Base | Centrifugo WS |
|-------------|---------|----------|------|
| Production | `https://app.unitedbit.com` | `/api/v1/` | `wss://app.unitedbit.com:8800` |
| Development | `https://dev-app.unitedbit.com` | `/api/v1/` | `wss://dev-app.unitedbit.com:8800` |
| Mobile | `https://m.unitedbit.com` | — | — |

## Documentation

- **[AGENTS.md](./AGENTS.md)** — Complete technical reference for AI agents (stack, all containers, services, patterns, conventions)
- **[UPGRADE_PLAN.md](./UPGRADE_PLAN.md)** — 8-phase upgrade roadmap (Phases 1-4 completed)
- **[unused.md](./unused.md)** — Deprecated file tracking (legacy reference)

## Tech Stack

| Layer | Technology | Version |
|-------|-----------|---------|
| UI | React | ^18.3.1 |
| Language | TypeScript | ^5.4.5 |
| State | Redux + Redux-Saga | 4.0.5 / 1.1.3 |
| Design | Material-UI | ^4.11.0 |
| Styling | styled-components | 5.2.0 |
| Data Grid | AG Grid Enterprise | ^25.3.0 |
| Charts | amCharts 4 + TradingView | ^4.10.9 |
| Real-time | Centrifugo WebSocket (centrifuge) | — |
| i18n | react-intl | 2.9.0 |
| HTTP | Fetch API + Axios (uploads) | ^0.21.0 |
| Build | Webpack (CRA-ejected) | 4.44.1 |
| Server | Express | 4.17.1 |
| Test | Jest + ts-jest + Testing Library | ^29.7.0 |

## License

MIT — see [LICENSE.md](./LICENSE.md)