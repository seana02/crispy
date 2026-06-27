# Crispy: Agent Instructions

## Project Overview

Crispy is a personal finance desktop app: Go backend + SolidJS frontend, built with Wails, local SQLite storage. Early active development; core accounting features (accounts, transactions, postings, recurring, rules) are partially implemented.

## Build & Dev Commands

**Frontend build:**
```bash
cd frontend && npm install && npm run build
```
Outputs to `frontend/dist`. This **must exist** before Go build.

**Full build:**
```bash
wails build --clean
```
Embeds `frontend/dist` into Go binary. Cross-platform: see `scripts/build-*.sh` for darwin/arm64, darwin/universal, windows/amd64.

**Linux dev (hot reload):**
```bash
wails dev -tags webkit2_41
```

**Tests:**
```bash
go test ./...
```

## Architecture

- **`domain/`** – pure business logic (structs, enums, validation)
- **`db/`** – SQLite repository pattern + migrations (in `db/migrations/`)
- **`service/`** – business operations (wraps repo, called by Wails bindings)
- **`scheduler/`** – recurring transaction job runner
- **`automation/`** – rule engine for auto-tagging and postings
- **`currency/`** – multi-currency conversion (partially stubbed)
- **`frontend/src/`** – SolidJS app, routes in `App.tsx`, components in `components/`, pages in `pages/`, stores in `stores/`

Dependency wiring: `main.go` → `app.go:NewApp()` → `dependencies.go:BuildDependencies()`. Wails binds Go methods on `App` struct for frontend RPC.

## Key Quirks

1. **Frontend must build before Go build.** `wails build` embeds `frontend/dist`; if missing, build fails.
2. **Database path:** `~/.config/crispy/data.db` (auto-created, migrations run on startup).
3. **Wails auto-generates TypeScript bindings** in `frontend/wailsjs/` after Go changes. Do **not** edit manually; regenerate instead.
4. **Frontend framework is SolidJS** (not React). Different reactivity model; use `createSignal`, `createStore`, not hooks.
5. **Tests use in-memory SQLite** with `SetMaxOpenConns(1)` (see `db/sqlite_test.go:28`). Required for test isolation.
6. **Enum bindings** (Go→TypeScript): Status = Pending/Posted/Cleared; AccountType = Asset/Liability/Revenue/Expense/Equity. Defined in `dependencies.go:WailsEnumBindings()`.

## Dependencies & Setup

- Go 1.23+ (toolchain 1.24.9 in `go.mod`)
- Node.js (for frontend)
- Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- Frontend: SolidJS, Vite, TypeScript, Ark UI, Lucide icons, @solidjs/router

## Key Files

- `main.go` – Wails app entry, window config
- `app.go` – App struct lifecycle (startup, domReady, beforeClose, shutdown)
- `dependencies.go` – dependency injection setup
- `frontend/package.json` – frontend scripts (dev, build, serve)
- `frontend/vite.config.ts` – Vite config (path aliases: `src/`, `wailsjs/`)
- `req.md` – detailed functional & technical requirements, data model, future roadmap
