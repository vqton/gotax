# GoTax — Agent Guide

## Mode

- **Caveman always active** — drop articles/filler/pleasantries/hedging. Fragments OK. Code/commits normal.
- **Karpathy** — ship first, polish later. Simple > clever. No comments/docstrings unless API surface. No speculative generality.

## Project

Vietnamese tax-compliant General Ledger API. Circular 99/2025/TT-BTC, Decree 123/2020/ND-CP. Multi-tenant, multi-company.

**Stack:** Go 1.26.5 · Gin v1.12 · GORM v1.31 (PostgreSQL via pgx v5) · golang-jwt v5 (RS256) · bcrypt · TOTP · golang-migrate v4 · go-playground/validator v10 · maroto/v2 (PDF) · zap · viper · casbin · go-i18n · swaggo/swag · testify · Alpine.js + Tailwind CSS v4 (frontend)

## Architecture

```
main.go                     →  entrypoint, DI wiring, backend selection (PG via GORM vs memory)
internal/domain/            →  models, repository interfaces, errors. Zero external deps. ~44 files, all package domain.
internal/auth/              →  JWT (RS256), TOTP, bcrypt, rate limiter
internal/authz/             →  Casbin RBAC policies
internal/config/            →  viper config loader
internal/db/                →  GORM setup, golang-migrate runner
internal/einvoice/          →  GDT e-invoice XML parse + generate (Decree 254/2026 schema). Pure encoding/xml.
internal/gdt/               →  GDT API client (e-invoice status, push)
internal/handler/           →  HTTP handlers, authMW, route registration
internal/htkk/              →  HTKK tax form XML generation
internal/i18n/              →  go-i18n translation bundle
internal/logger/            →  zap logger + gin middleware
internal/repository/        →  per-module PG + memory impls (pg_*.go, memory_*.go)
internal/service/           →  business rules, validation, orchestration. Pure Go.
internal/validate/          →  go-playground/validator singleton + custom validators
internal/xmldsig/           →  XML digital signature (RSA, for e-invoice)
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

`internal/domain/models*.go` — all `package domain`. Split by bounded context. ~44 files, all same package.

Adding a model = add to correct existing file or create new `models_*.go`. No sub-packages, no import changes.

GORM models live in `internal/domain/models_gorm_*.go` (one per module). Table name set via `TableName()` method.

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

Company-scoped handlers pass `company_id` as **query param**, not from JWT context:
```go
companyID := c.Query("company_id")
```
Same pattern across all module handlers. **Not** `c.GetString("company_id")`.

## Routes

Route registration entry: `RegisterRoutesWithCompany(r, h, ch, th, cashH, bankH, purchaseH, saleH, whH, faH, pwH, recH, budH, ccdcH, ccH, keeperH, authMW, adminMW)` at `handler/handler.go:208`.

Adding a new handler requires:
1. Create `internal/handler/<module>_handler.go` with `Register<R>Routes(r, h, authMW)`
2. Add handler param to `RegisterRoutesWithCompany` signature
3. Wire in both PG and memory branches of `main.go`

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
| Warehouse Keeper | `/api/v1/warehouse/keeper/**` | `WarehouseKeeperHandler` |
| Fixed Asset | `/api/v1/fixed-assets/**` | `FAHandler` |
| Payroll | `/api/v1/payroll/**` | `PayrollHandler` |
| Recurring | `/api/v1/recurring-entries/**` | `RecurringHandler` |
| Budget | `/api/v1/budgets/**` | `BudgetHandler` |
| CCDC | `/api/v1/ccdc/**` | `CCDCHandler` |
| Cost Centers | `/api/v1/cost-centers/**` | `CostCenterHandler` |

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

31+ versioned files in `migrations/`. **Versioned** (`000001_title.up.sql` + `000001_title.down.sql`) → auto-discovered by golang-migrate and auto-run on PG startup. Current latest: `000031_warehouse_keeper`.

**Legacy** (`.sql` only, no version prefix) — UNUSED. Do not reference: `002_gl_schema_circular99.sql`, `003_company_schema.sql`, `003_cash_schema.sql`, `004_bank_module.sql`, `004_advance_schema.sql`, `006_sale_schema.sql`, `007_warehouse_schema.sql`.

Adding a migration: write `{next_version}_{title}.up.sql` + `.down.sql` in `migrations/`.

## Repository Files

Per-module naming: `pg_<module>.go` + `memory_<module>.go` in `internal/repository/`. Adding a module = two new files.

PG repos use GORM (`*gorm.DB`). Memory repos use `sync.RWMutex` + maps.

Memory ID convention: copy struct before mutation, generate ID for copy, write ID back to original pointer. Always use same pointer after Create.

> **Memory repo quirk** — Purchase, Sale, Warehouse memory repos implement all repository interfaces from a single struct (e.g. `MemoryPurchaseRepo` implements Supplier, PurchaseOrder, GRN, SupplierInvoice, APTransaction, CostAllocation repos). This is why `main.go` passes the same repo for all params.

## Service Wiring

Each module-level service has its own `New*Service()` with its own repos. No DI framework — manual wiring in `main.go` (both PG and memory branches).

## Handler Files

`handler.go` contains only: `Handler` struct, `NewHandler`, `RegisterRoutes`, `RegisterRoutesWithCompany`. Domain-specific handlers split into separate files. Each has its own `New*Handler` and `Register*Routes` function.

## Validate Package

`internal/validate/` — singleton `go-playground/validator/v10` instance with custom validators for domain enum types.

When adding a new module that uses the validator: register custom validators in `internal/validate/validator.go`, add `validate` struct tags to domain models, create module validation function in `internal/validate/`.

## Frontend

**No build step.** Tailwind CSS v4 compiled to `web/static/css/app.css`. Alpine.js bundled at `web/static/js/alpine.min.js`. No webpack, no Vite.

**UI entry points:**
- `/app/*` — Main accounting UI. MISA SME 2026 layout: left sidebar, top bar, content area.
- `/payroll/*` — Payroll module. Separate nav bar (top).
- `/login` — Auth pages at `web/auth/*.html`.

**Shared JS:** `web/static/js/app.js` — sidebar nav, API client (`apiGet`/`apiPost`/`apiPut`/`apiDelete` with JWT refresh), formatters, Alpine store.

**Page pattern:** Each HTML page is standalone with Alpine.js `x-data`. Calls `mountAppShell(title, activePath)` on init. Uses `apiGet`/`apiPost` from app.js.

## Adding a Feature — Step Order

1. Interface method in `internal/domain/interfaces.go`
2. GORM model in `internal/domain/models_gorm_*.go`
3. Repository impl in `internal/repository/pg_*.go` + `memory_*.go`
4. Service method in `internal/service/<module>_service.go`
5. Handler method in `internal/handler/<module>_handler.go` + `Register*Routes`
6. Add handler param to `RegisterRoutesWithCompany` in `handler.go`
7. Add handler param to `RegisterRoutesWithCompanyOpt` in `casbin_register.go`
8. Wire in `main.go` (PG branch + memory branch)
9. Tests in `internal/handler/<module>_handler_test.go`
10. `go vet ./... && go test -count=1 ./...`

## Gotchas

- **`JWT_SECRET` required** — server panics if unset.
- **RSA non-deterministic** — same seed ≠ same keys. Tokens die on restart.
- **Go 1.26.5** — `range-over-func` and other modern Go features available.
- **`GenerateToken` (HMAC-SHA256 with hardcoded secret)** — test-only dead code. Production uses `GenerateAccessToken` (RS256).
- **Config** from `config.yaml` + env vars via viper. Env overrides (e.g. `JWT_SECRET`, `DATABASE_URL`).
- **Audit logging** via `internal/handler/audit.go` middleware — logs all state-mutating operations.
- **Web UI** at `web/app/*.html` (main) and `web/payroll/*.html` (payroll). Served by gin. Not a SPA.
- **PDF** generation with `maroto/v2` for opening balance reports. Not wired to any endpoint yet.
- **Handler error mapping** uses `errors.Is` to switch on domain errors per module. See `internal/handler/fixed_asset_handler.go:79` (the `faError` func) for pattern.
- **GDT env vars**: `GDT_BASE_URL`, `GDT_TOKEN` — e-invoice API client. Without them, `newGDTClient()` returns nil.
- **E-invoice signing**: `EINVOICE_SIGNING_KEY` (RSA PEM), `EINVOICE_CERT_SERIAL`. Without key, ephemeral key generated at startup — signatures invalid after restart.
- **Migration idempotency**: Always use `CREATE TABLE IF NOT EXISTS` / `CREATE INDEX IF NOT EXISTS` in migrations. Legacy `.sql` files (002-007) may have been run manually against PG, creating tables outside migration tracking. Without `IF NOT EXISTS`, golang-migrate fails with `relation already exists`.
- **SQLite rejected**: Schema deeply PG-specific (`pgcrypto`, `gen_random_uuid()`, `TIMESTAMPTZ`, `plpgsql` triggers, `SUBSTRING FROM` regex, `GREATEST()`). 20+ raw SQL queries in PG repos. In-memory backend already covers dev/test.
- **No LINQ library needed**: Standard `for` loops + `slices.SortFunc` sufficient. Don't add collection libraries unless profiling shows bottleneck.

## Commenting & Documentation Standards

Karpathy: no comments unless API surface. ERP rules below override for business-logic-heavy code.

**When to comment (mandatory):**
- Business rule driven by regulation (Circular 99, Decree 123, Labor Code 2019 Art. 46)
- Non-obvious accounting logic (net-to-gross search, dual-half-year PIT brackets)
- Account mapping decisions (which GL account why)
- Data exclusion reasons in migrations

**When NOT to comment:**
- What syntax does (`// iterate over slice` — delete)
- Obvious flow (`// call service` — delete)
- Self-documenting names (`CalculateNetToGross` — name is enough)

**Standard task tags:**
```
// TODO: sync with master CRM before go-live
// FIX: temporary patch during UAT
// DEPRECATED: old table, kept for migration rollback
```

**GL account mapping comments:**
```
// Account 6421 — Salary expense per Circular 99/2025, Appendix 01
// Account 3331 — Payable to employee (salary withheld)
// Account 3335 — SI payable (employer portion, Art. 86 Social Insurance Law)
```

**Revision history (complex business logic files only):**
```go
// REVISION HISTORY:
// 2026-08-07   dev   PAYROLL-001   Initial: net-to-gross with dual-half-year PIT
// 2026-08-07   dev   PAYROLL-002   Add: severance per Labor Code 2019 Art. 46
```
