# ub-admin-main

Admin Panel SPA for the UnitedBit cryptocurrency exchange platform. Provides back-office management of users, orders, deposits, withdrawals, balances, currencies, and KYC verification workflows.

## Tech Stack

- **React 17.0.2** / **TypeScript 5.4.5** / CRA 5.0.1 (Webpack 5)
- **Redux Toolkit 1.3.6** + **Redux-Saga 1.1.3** + redux-injectors 1.3.0
- **Material-UI 4.12.4** + styled-components 5.1.1
- **AG Grid 23.2.0** for all data tables
- **i18next** (English / German)
- **RxJS** Subject-based pub/sub event bus (`MessageService`)

## Prerequisites

- **Node.js** ≥ 18 (see `.nvmrc`)
- **Yarn** 1.22+ (preferred — `npm ci` has known arborist bugs)

## Getting Started

```bash
# Install dependencies
yarn install

# Start development server (port 3000)
npm start

# Run tests
npm test

# Type check
npm run checkTs

# Lint
npm run lint

# Production build
npm run build
```

### Environment Setup

Copy `.env.example` to `.env.local` and configure:

```bash
REACT_APP_API_BASE_URL=https://admin.unitedbit.com/api/v1/
REACT_APP_API_URL=https://admin.unitedbit.com
REACT_APP_WEB_APP_URL=https://dev-app.unitedbit.com
```

## Project Structure

```
src/
├── app/
│   ├── components/         # 34 shared reusable UI components
│   ├── containers/         # 27 page-level containers (Redux-Saga pattern)
│   ├── constants.ts        # AppPages enum, WindowTypes, Buttons
│   ├── index.tsx           # App root: router, routes, nav categories
│   ├── ForceStyles.tsx     # Global CSS injection
│   └── NewWindowContainer.tsx  # Multi-window support
├── services/               # API service layer (11 files)
│   ├── apiService.ts       # Singleton HTTP client (fetch, JWT, CSRF, retry)
│   ├── messageService.ts   # RxJS pub/sub event bus (67 events)
│   └── ...                 # Domain-specific API services
├── store/                  # Redux store configuration
├── types/                  # Root state types
├── styles/                 # Theme, global styles
├── locales/                # i18n translations (en/de)
├── utils/                  # Formatters, hooks, grid utilities
└── index.tsx               # Entry point
```

### Container Pattern (6-file convention)

Each page container follows:
```
containers/<PageName>/
  index.tsx      — Component (useInjectReducer, useInjectSaga, JSX)
  slice.ts       — Redux Toolkit slice (action triggers for sagas)
  saga.ts        — Side effects (API calls → MessageService.send())
  selectors.ts   — Memoized selectors
  types.ts       — TypeScript interfaces
  Loadable.tsx   — Code splitting wrapper
```

## Key Commands

| Command | Purpose |
|---------|---------|
| `npm start` | Dev server (port 3000) |
| `npm run build` | Production build |
| `npm test` | Jest tests (90% coverage threshold) |
| `npm run checkTs` | TypeScript type-check |
| `npm run lint` | ESLint check |
| `npm run lint:fix` | ESLint auto-fix |
| `npm run generate` | Scaffold new container via plop |

## Documentation

- **[AGENTS.md](./AGENTS.md)** — AI agent guide with how-to recipes, full API reference, and architecture details
- **[ARCHITECTURE.md](./ARCHITECTURE.md)** — Comprehensive architecture overview
- **[DEEP-AUDIT.md](./DEEP-AUDIT.md)** — Code quality audit report and scorecard
- **[PLAN.md](./PLAN.md)** — Master improvement plan with phased execution
- **[docs/GLOSSARY.md](./docs/GLOSSARY.md)** — Domain terminology (trading, financial, technical terms)

## CI/CD

GitLab CI pipeline (`.gitlab-ci.yml`):
- **Build**: Docker image with `node:18-slim`
- **Deploy**: Copy build artifacts to host
- **Notify**: Telegram bot on success/failure
- **Trigger**: `master` branch only
