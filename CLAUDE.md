# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Usta Top** is a full-stack job marketplace platform (Uzbek job/resume board). The backend is Go with `httprouter`, the frontend is Vue 3 + Vite. Both live in the same repo.

## Commands

### Backend (`main_service/`)

```bash
make run            # go run main.go
make migrate name=<migration_name>   # create a new SQL migration file
```

Run directly: `go run main.go` from `main_service/`.

Swagger docs are auto-generated in `main_service/docs/` — regenerate with `swag init` if you modify route annotations.

### Frontend (`admin-panel/`)

```bash
npm run dev         # Vite dev server
npm run build       # production build → dist/
npm run preview     # preview the production build
```

## Environment Setup

Copy `.env` in `main_service/` and set these required variables:

```
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=
DB_NAME=main_service
DB_SSLMODE=disable
DB_TIMEZONE=Asia/Tashkent
JWT_KEY=<secret>          # required — app exits if missing
APP_PORT=8080             # optional, defaults to 8080
ALLOWED_ORIGINS=*         # optional, defaults to *
```

Frontend: set `VITE_API_BASE_URL` (defaults to `http://localhost:8080/api/v1`).

## Architecture

### Backend (`main_service/`)

**Entry point**: `main.go` — initializes DB (which auto-runs migrations), registers all module routes, applies middleware stack, starts server.

**Middleware stack** (applied globally):
1. `SecurityHeaders()` — standard security headers
2. `CORS()` — configurable via `ALLOWED_ORIGINS`
3. `RateLimit(30, 60)` — token bucket per IP (30 req/s, burst 60)
4. `http.MaxBytesHandler` — 4 MB body limit

**Module layout** — every feature under `module/` follows the same 4-layer pattern:

```
module/<name>_service/
├── cmd.go          # wires routes for this module
├── handler/        # HTTP layer — parse request, call service, write response
├── service/        # business logic (interface + implementation)
├── dto/            # request/response structs with `validate:` tags
└── model/          # DB entity structs
```

Adding a new module: create the directory, implement the 4 layers, then register `module.Routes(router)` in `main.go`.

**Shared helpers** (`helper/`):
- `helper.ENV(key)` — reads env vars (loaded by `LoadEnv()` at startup)
- `helper.WriteJSON`, `WriteError`, `WriteValidation` — uniform response format
- `helper.ValidateStruct` — go-playground/validator with Uzbek error messages
- `helper.GenerateToken` / `helper.ParseToken` — JWT HS256, 2-hour expiry
- `helper.Paginate`, `helper.CursorPaginate` — pagination utilities

**Database**: PostgreSQL via pgx/v5 (no ORM). Pool: min 10 / max 50 connections. All queries use parameterized statements.

**Migrations**: SQL files in `config/migrations/` are embedded at compile time and applied automatically on `DBConnect()`. Name files with a numeric prefix (`001_`, `002_`, …) — they run in alphabetical order and are tracked in the `schema_migrations` table.

**Authentication**: Bearer JWT in `Authorization` header. `middleware/check_role.go` and `middleware/check_role_group.go` enforce role-based access.

**File uploads**: stored in `uploads/` and served at `/uploads/*filepath`. Max 3 MB per file.

### Frontend (`admin-panel/src/`)

**Routing** (`app/router/index.js`): all routes are language-prefixed (`/:lang/*`). Two layouts:
- `AdminLayout` for `/:lang/admin/*`
- `ClientLayout` for public pages

**Module layout** — mirrors backend modules:

```
modules/<feature>/
├── api/            # axios calls to backend
├── pages/          # Vue page components
└── store/          # Pinia store
```

**Shared infra**:
- `app/providers/axios.js` — Axios instance that injects `Authorization: Bearer <token>` on every request
- `shared/stores/langStore` — language state synced to localStorage and URL param
- `shared/layouts/` — AdminLayout, ClientLayout
- `shared/composables/`, `shared/components/`, `shared/utils/`

**Path alias**: `@/` resolves to `src/` (configured in `vite.config.js`).

## Key Patterns

- **Response format**: every backend handler calls `helper.WriteJSON` / `helper.WriteError` — never write raw JSON manually.
- **Validation**: tag DTOs with `validate:` and call `helper.ValidateStruct`; errors are returned in Uzbek via `helper.WriteValidation`.
- **Language support**: translations are stored in the `translations` table (JSONB); category names use a JSONB column (`categories.name`). The frontend language is `/:lang` in the URL.
- **No test files exist** — the project has no automated test suite yet.
