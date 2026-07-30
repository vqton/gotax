# GoTax — Agent Guide

## Mode

- **Caveman always active** — drop articles/filler/pleasantries/hedging. Fragments OK. Code/commits written normal.
- **Karpathy rules** — ship first, polish later. Simple > clever. Write less, ship more. Don't over-abstract. Prefer flat code.
- **Main theory** — every line serves the feature. No comments/docstrings unless API surface. No speculative generality.

## Project

Vietnamese tax-compliant General Ledger API. Circular 99/2025/TT-BTC, Decree 123/2020/ND-CP. Multi-tenant, multi-company.

**Stack:** Go 1.26.5 · Gin v1.12 · GORM v1.31 (PostgreSQL via pgx v5) · golang-jwt v5 (RS256) · bcrypt · TOTP · golang-migrate v4 · zap · viper · casbin · go-i18n · swaggo/swag · testify

## Architecture

```
main.go                     →  entrypoint, DI wiring, backend selection (PG via GORM vs memory)
internal/domain/            →  models, repository interfaces, errors. Zero external deps.
internal/auth/              →  JWT (RS256), TOTP, bcrypt, rate limiter
internal/service/           →  business rules, validation, orchestration. Pure Go.
internal/handler/           →  HTTP handlers, authMW, route registration
internal/repository/        →  per-module PG + memory impls (pg_*.go, memory_*.go)
internal/db/                →  GORM setup, golang-migrate runner
```

**Two backends** controlled by `DATABASE_URL` env var:
- Set → PostgreSQL via GORM + golang-migrate (auto-runs on startup)
- Unset → in-memory maps + `sync.RWMutex` (dev/test)

**Request lifecycle:**
```
HTTP → gin.Engine → authMW (JWT verify) → roleMW (RBAC) → Handler → Service → Repository
```

**Layer rules:**
- `Handler` — parse request, call service, return JSON. Zero logic.
- `Service` — business rules, validation. No HTTP imports.
- `Repository` — data access. Two impls per module: `PG*Repo` + `Memory*Repo`.

## Domain Models

`internal/domain/models*.go` — all `package domain`. Split by bounded context. 7+ files, all same package.

Adding a model = add to correct existing file or create new `models_*.go`. No sub-packages, no import changes.

## Auth

```
POST /api/v1/auth/login       →  username+password → JWT pair (access 15m + refresh 7d)
POST /api/v1/auth/refresh      →  rotate refresh token
POST /api/v1/auth/totp/verify  →  2FA challenge after login

authMW → extract Bearer → verify RS256 → set user_id, username, role in gin.Context
GetUserID(c) → helper for user_id from context
RoleMiddleware(admin, chief) → RBAC gate
```

JWT: RS256. Key pair generated at startup from `JWT_SECRET` via `crypto/rand`. **Not deterministic** — same seed ≠ same keys. Restart invalidates all tokens.

## Routes

| Group | Path prefix | Handler | Auth |
|-------|-----------|---------|------|
| Auth | `/api/v1/auth/{login,refresh,forgot,reset,totp/verify}` | `Handler` | Public |
| GL | `/api/v1/accounts`, `journal-entries`, `periods`, `exchange-rates`, `reports`, `coa/*`, `audit` | `Handler` | authMW |
| User | `/api/v1/users`, `/me`, `/auth/{change-password,logout,totp/*}` | `Handler` | authMW |
| Company | `/api/v1/companies/**`, `branches/**`, `departments/**`, `employees/**`, `bank-accounts/**`, `fiscal-years/**` | `CompanyHandler` | authMW |
| Tax | `/api/v1/tax/**` | `TaxHandler` | authMW |
| Cash | `/api/v1/cash/**` | `CashHandler` | authMW |
| Bank | `/api/v1/bank/**` | `BankHandler` | authMW |
| Purchase | `/api/v1/purchase/**` | `PurchaseHandler` | authMW |
| Sale | `/api/v1/sale/**` | `SaleHandler` | authMW |
| Warehouse | `/api/v1/warehouse/**` | `WarehouseHandler` | authMW |
| Fixed Asset | `/api/v1/fixed-assets/**` | `FAHandler` | authMW |

Route registration: `RegisterRoutesWithCompany(r, h, ch, th, cashH, bankH, purchaseH, saleH, whH, faH, authMW, adminMW)` at `handler/handler.go:207`.

## Commands

```sh
go build ./...
go vet ./...
go test ./...
go test -v -run TestCreateCompany ./internal/handler/
go test -count=1 ./...        # fresh run, no cache

# Run (memory)
JWT_SECRET=devsecret go run .

# Run (PostgreSQL — auto-migrates via golang-migrate)
DATABASE_URL=postgres://... JWT_SECRET=devsecret go run .

# Regenerate swagger (annotations in main.go + handler comments)
swag init --parseDependency --parseInternal
# → docs/docs.go, docs/swagger.json, docs/swagger.yaml — DO NOT EDIT
```

No Makefile, no linter config, no Dockerfile. Lint: `go vet`.

## Testing

- All tests use **in-memory repos** + `httptest.NewRecorder`. No DB, no integration.
- Handler tests: real in-memory repos + real service, gin engine with mock auth middleware (sets user_id, role). No mock services.
- Service tests: real in-memory repos.
- Domain tests: struct validation tests.
- Adding a service method → add to service + both repos + handler + test. No mock setup needed.
- `go test -count=1 ./...` before commit (fresh run, bypasses cache).

## Migration System

Uses `golang-migrate/migrate/v4`. Migration files in `migrations/` named `{version}_{title}.up.sql` / `.down.sql`. Version numbers are sequential (000001-000010 currently). Auto-runs on startup when PG is configured.

Adding a migration: write `.up.sql` + `.down.sql` in `migrations/` with next version number. That's it — golang-migrate discovers them automatically.

Legacy `.sql` files without version numbers (e.g. `001_gl_schema.sql.deprecated`, `003_cash_schema.sql`) are unused — do not reference.

## Repository Files

Per-module naming: `pg_<module>.go` + `memory_<module>.go` in `internal/repository/`. Not monolithic. Adding a module = two new files.

| Module | PG file | Memory file |
|--------|---------|-------------|
| GL | `pg.go` | `memory.go` |
| Company | `pg_company.go` | `memory_company.go` |
| Tax | `pg_tax.go` | `memory_tax.go` |
| Bank | `pg_bank.go` | `memory_bank.go` |
| Cash | `pg_cash.go` | `memory_cash.go` |
| Purchase | `pg_purchase.go` | `memory_purchase.go` |
| Sale | `pg_sale.go` | `memory_sale.go` |
| Warehouse | `pg_warehouse.go` | `memory_warehouse.go` |
| FA | `pg_fa.go` | `memory_fa.go` |

PG repos use GORM (`*gorm.DB`). Memory repos use `sync.RWMutex` + maps + auto-generated IDs.

ID convention: memory repos copy struct before mutation, generate ID for copy, write ID back to original pointer. Always use same pointer after Create.

## Module Readiness

| Module | Status | Notes |
|--------|--------|-------|
| GL | PROD | Core accounts, journal, periods, reports, COA, opening balances |
| Auth | PROD | Login, JWT (RS256), TOTP, refresh, rate limit, lockout |
| Company | PROD | Company, branches, departments, employees, fiscal years, bank accounts |
| Cash | PROD | Receipts, payments, transfers, petty cash, advances |
| Bank | PROD | Statements, reconciliation, payment orders, loans, term deposits |
| Tax | ~20% — NOT PROD | Declaration stubs, e-invoice CRUD, rates CRUD. Missing: declaration engine, XML gen, GDT API, full lifecycle |
| Purchase | ~0% | Interface + PG repos exist. Service layer incomplete |
| Sale | ~0% | Interface + PG repos exist. Service layer incomplete |
| Warehouse | ~0% | Interface + PG repos exist. Service layer incomplete |
| FA | PROD | Full CRUD, depreciation engine (SL/DB), business ops, allocations, inventory |

Full tax/purchase/warehouse specs at `docs/`.

## Adding a Feature — Step Order

1. Interface method in `internal/domain/interfaces.go`
2. Repository impl in `internal/repository/pg_*.go` + `memory_*.go`
3. Service method in `internal/service/service.go` (or `<module>_service.go`)
4. Handler method in `internal/handler/handler.go` (or `<module>_handler.go`) + route registration in `RegisterFixedAssetRoutes` / `RegisterPurchaseRoutes` etc.
5. Tests in `internal/handler/<module>_handler_test.go`
6. Wire in `main.go`
7. `go vet ./... && go test -count=1 ./...`

## Gotchas

- **`JWT_SECRET` required** — server panics if unset.
- **RSA non-deterministic** — same seed ≠ same keys. Tokens die on restart.
- **Go 1.26.5** — `range-over-func` and other modern Go features available.
- **`GenerateToken` (HMAC-SHA256 with hardcoded secret)** — test-only dead code. Production uses `GenerateAccessToken` (RS256).
- **Config** from `config.yaml` + env vars via viper. Env overrides (e.g. `JWT_SECRET`, `DATABASE_URL`).
- **Audit logging** via `internal/handler/audit.go` middleware — logs all state-mutating operations.
