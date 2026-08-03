# GoTax — Agent Guide

## Mode

- **Caveman always active** — drop articles/filler/pleasantries/hedging. Fragments OK. Code/commits normal.
- **Karpathy** — ship first, polish later. Simple > clever. No comments/docstrings unless API surface. No speculative generality.

## Project

Vietnamese tax-compliant General Ledger API. Circular 99/2025/TT-BTC, Decree 123/2020/ND-CP. Multi-tenant, multi-company.

**Stack:** Go 1.26.5 · Gin v1.12 · GORM v1.31 (PostgreSQL via pgx v5) · golang-jwt v5 (RS256) · bcrypt · TOTP · golang-migrate v4 · go-playground/validator v10 · maroto/v2 (PDF) · zap · viper · casbin · go-i18n · swaggo/swag · testify

## Architecture

```
main.go                     →  entrypoint, DI wiring, backend selection (PG via GORM vs memory)
internal/domain/            →  models, repository interfaces, errors. Zero external deps.
internal/auth/              →  JWT (RS256), TOTP, bcrypt, rate limiter
internal/service/           →  business rules, validation, orchestration. Pure Go.
internal/handler/           →  HTTP handlers, authMW, route registration
internal/repository/        →  per-module PG + memory impls (pg_*.go, memory_*.go)
internal/db/                →  GORM setup, golang-migrate runner
internal/validate/          →  go-playground/validator singleton + custom validators
internal/einvoice/          →  GDT e-invoice XML parse + generate (Decree 254/2026 schema). Pure encoding/xml.
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

`internal/domain/models*.go` — all `package domain`. Split by bounded context. 28 files, all same package.

Adding a model = add to correct existing file or create new `models_*.go`. No sub-packages, no import changes.

No separate DTO types — handler binds JSON directly into domain structs. `validate` struct tags on domain models checked by `internal/validate/` functions.

## Auth

JWT: RS256. Key pair generated at startup from `JWT_SECRET` via `crypto/rand`. **Not deterministic** — same seed ≠ same keys. Restart invalidates all tokens.

```
POST /api/v1/auth/login       →  username+password → JWT pair (access 15m + refresh 7d)
POST /api/v1/auth/refresh      →  rotate refresh token
POST /api/v1/auth/totp/verify  →  2FA challenge after login

authMW → extract Bearer → verify RS256 → set user_id, username, role in gin.Context
GetUserID(c) → helper for user_id from context
RoleMiddleware(admin, chief) → RBAC gate
```

## Company ID Pattern

Company-scoped handlers (purchase, warehouse, FA) pass `company_id` as **query param**, not from JWT context:
```go
companyID := c.Query("company_id")
```
Same pattern across all module handlers.

## Routes

Route registration entry: `RegisterRoutesWithCompany(r, h, ch, th, cashH, bankH, purchaseH, saleH, whH, faH, authMW, adminMW)` at `handler/handler.go:207`.

| Group | Prefix | Handler |
|-------|--------|---------|
| Auth | `/api/v1/auth/{login,refresh,forgot,reset,totp/verify}` | `Handler` |
| GL | `/api/v1/accounts`, `journal-entries`, `periods`, `exchange-rates`, `reports`, `coa/*`, `audit` | `Handler` |
| User | `/api/v1/users`, `/me`, `/auth/{change-password,logout,totp/*}` | `Handler` |
| Company | `/api/v1/companies/**` | `CompanyHandler` |
| Tax | `/api/v1/tax/**` | `TaxHandler` |
| Cash | `/api/v1/cash/**` | `CashHandler` |
| Bank | `/api/v1/bank/**` | `BankHandler` |
| Purchase | `/api/v1/purchase/**` | `PurchaseHandler` |
| Sale | `/api/v1/sale/**` | `SaleHandler` |
| Warehouse | `/api/v1/warehouse/**` | `WarehouseHandler` |
| Fixed Asset | `/api/v1/fixed-assets/**` | `FAHandler` |

Auth middleware on all groups except auth endpoints.

## Commands

```sh
JWT_SECRET=devsecret go run .                           # memory backend
DATABASE_URL=postgres://... JWT_SECRET=devsecret go .   # PG backend (auto-migrates)
go build ./... && go vet ./...
go test -count=1 ./...                                  # fresh run, no cache
go test -v -run TestCreateCompany ./internal/handler/   # single test

# Regenerate swagger (annotations in main.go + handler comments)
swag init --parseDependency --parseInternal
# → docs/docs.go, docs/swagger.json, docs/swagger.yaml — DO NOT EDIT
```

No Makefile, no Dockerfile, no linter config. Lint: `go vet`.

## Testing

- All tests use **in-memory repos** + `httptest.NewRecorder`. No DB, no integration.
- Handler tests: real in-memory repos + real service, gin engine with mock auth middleware (sets user_id, role). No mock services.
- Service tests: real in-memory repos.
- Domain tests: struct validation tests.
- Adding a service method → add to service + both repos + handler + test. No mock setup needed.
- `go test -count=1 ./...` before commit.

## Migration System

28 files in `migrations/`. **Versioned** (`000001_title.up.sql` + `000001_title.down.sql`) → auto-discovered by golang-migrate and auto-run on PG startup. Versioned files have both `.up.sql` and `.down.sql`. Current latest: `000010_fa_schema`.

**Legacy** (`.sql` only, no version prefix) — UNUSED. Do not reference: `002_gl_schema_circular99.sql`, `003_company_schema.sql`, `003_cash_schema.sql`, `004_bank_module.sql`, `004_advance_schema.sql`, `006_sale_schema.sql`, `007_warehouse_schema.sql`.

Adding a migration: write `{next_version}_{title}.up.sql` + `.down.sql` in `migrations/`.

## Repository Files

Per-module naming: `pg_<module>.go` + `memory_<module>.go` in `internal/repository/`. Adding a module = two new files.

| Module | PG | Memory |
|--------|----|--------|
| GL | `pg.go` | `memory.go` |
| Company | `pg_company.go` | `memory_company.go` |
| Tax | `pg_tax.go` | `memory_tax.go` |
| Bank | `pg_bank.go` | `memory_bank.go` |
| Cash | `pg_cash.go` | `memory_cash.go` |
| Purchase | `pg_purchase.go` | `memory_purchase.go` |
| Sale | `pg_sale.go` | `memory_sale.go` |
| Warehouse | `pg_warehouse.go` | `memory_warehouse.go` |
| FA | `pg_fa.go` | `memory_fa.go` |

PG repos use GORM (`*gorm.DB`). Memory repos use `sync.RWMutex` + maps.

Memory ID convention: copy struct before mutation, generate ID for copy, write ID back to original pointer. Always use same pointer after Create.

> **Memory repo quirk** — Purchase, Sale, Warehouse memory repos implement all repository interfaces from a single struct (e.g. `MemoryPurchaseRepo` implements Supplier, PurchaseOrder, GRN, SupplierInvoice, APTransaction, CostAllocation repos). This is why `main.go` passes the same repo for all params.

## Service Wiring

`NewService()` in `internal/service/service.go:209` takes **16 repository interfaces** (all GL+Cash related). Each module-level service (Company, Tax, Purchase, Bank, Sale, Warehouse, FA) has its own `New*Service()` with its own repos. No DI framework — manual wiring in `main.go`.

## Validate Package

`internal/validate/` — singleton `go-playground/validator/v10` instance with custom validators for domain enum types: `fastatus`, `damethod`, `fasource`, `fatrtype`, `disposaltype`.

Validation flow: service calls `validate.FixedAsset(a)` or `validate.FixedAssetCategory(c)` which sets defaults, runs struct tag validation, maps `ValidationErrors` to domain errors. Domain still exports `Validate()` methods but service uses validate package instead.

When adding a new module that uses the validator: register custom validators in `internal/validate/validator.go`, add `validate` struct tags to domain models, create module validation function in `internal/validate/`.

## Module Readiness

| Module | Status | Notes |
|--------|--------|-------|
| GL | PROD | Core accounts, journal, periods, reports, COA, opening balances |
| Auth | PROD | Login, JWT (RS256), TOTP, refresh, rate limit, lockout |
| Company | PROD | Company, branches, departments, employees, fiscal years, bank accounts |
| Cash | PROD | Receipts, payments, transfers, petty cash, advances |
| Bank | PROD | Statements, reconciliation, payment orders, loans, term deposits |
| FA | PROD | Full CRUD, depreciation engine (SL/DB), business ops, allocations, inventory |
| Tax | ~40% | Rate resolver + VAT/CIT/PIT engines (rate-table), declaration engine (GL→GTGT01/TNDN03, cross-validation), declaration→payment automation (due dates). Missing: form XML gen, GDT API push, e-invoice engine |
| Purchase | PROD (P2 closed) | Full domain models + repos + service + handlers + 55 routes + 47 handler tests + 66 service tests + 28 domain tests + 7 einvoice tests. Requisition+approval, returns (return GRN + credit note), import + landed cost, AP FX revaluation (515/635), GDT e-invoice XML (parse/generate, internal/einvoice), 3-way matching, GL auto-posting, doubtful-debt provisioning (Circular 99), 5 reports. Missing: GDT API push, supplier portal |
| Sale | ~0% | Interface + PG + memory repos. Service incomplete |
| Warehouse | ~0% | Interface + PG + memory repos. Service incomplete |

Full tax/purchase/warehouse specs at `docs/`.

## Adding a Feature — Step Order

1. Interface method in `internal/domain/interfaces.go`
2. Repository impl in `internal/repository/pg_*.go` + `memory_*.go`
3. Service method in `internal/service/<module>_service.go`
4. Handler method in `internal/handler/<module>_handler.go` + route registration
5. Validation in `internal/validate/` (add tags + validate function)
6. Tests in `internal/handler/<module>_handler_test.go`
7. Wire in `main.go` (PG + memory branches)
8. `go vet ./... && go test -count=1 ./...`

## Gotchas

- **`JWT_SECRET` required** — server panics if unset.
- **RSA non-deterministic** — same seed ≠ same keys. Tokens die on restart.
- **Go 1.26.5** — `range-over-func` and other modern Go features available.
- **`GenerateToken` (HMAC-SHA256 with hardcoded secret)** — test-only dead code. Production uses `GenerateAccessToken` (RS256).
- **Config** from `config.yaml` + env vars via viper. Env overrides (e.g. `JWT_SECRET`, `DATABASE_URL`).
- **Audit logging** via `internal/handler/audit.go` middleware — logs all state-mutating operations.
- **Web UI** at `web/auth/*.html` (login, 2FA, forgot/reset password). Served by gin. Not a SPA.
- **PDF** generation with `maroto/v2` for opening balance reports. Not wired to any endpoint yet.
- **Handler error mapping** uses `errors.Is` to switch on domain errors per module. See `internal/handler/fixed_asset_handler.go:475` for pattern.
