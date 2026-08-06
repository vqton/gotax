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
internal/domain/            →  models, repository interfaces, errors. Zero external deps. 31 files, all package domain.
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

`internal/domain/models*.go` — all `package domain`. Split by bounded context. 32 files, all same package.

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

Route registration entry: `RegisterRoutesWithCompany(r, h, ch, th, cashH, bankH, purchaseH, saleH, whH, faH, pwH, authMW, adminMW)` at `handler/handler.go:201`.

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
| Payroll | `/api/v1/payroll/**` | `PayrollHandler` |

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

~30 versioned files in `migrations/`. **Versioned** (`000001_title.up.sql` + `000001_title.down.sql`) → auto-discovered by golang-migrate and auto-run on PG startup. Versioned files have both `.up.sql` and `.down.sql`. Current latest: `000024_tax_rate_incentive`.

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
| Cash | `pg_cash.go` | (in `memory.go`) |
| Purchase | `pg_purchase.go` | `memory_purchase.go` |
| Sale | `pg_sale.go` | `memory_sale.go` |
| Warehouse | `pg_warehouse.go` | `memory_warehouse.go` |
| FA | `pg_fa.go` | `memory_fa.go` |
| Opening Balance | `pg_opening_balance.go` | `memory_opening_balance.go` |
| Payroll | `pg_payroll.go` | `memory_payroll.go` |

PG repos use GORM (`*gorm.DB`). Memory repos use `sync.RWMutex` + maps.

Memory ID convention: copy struct before mutation, generate ID for copy, write ID back to original pointer. Always use same pointer after Create.

> **Memory repo quirk** — Purchase, Sale, Warehouse memory repos implement all repository interfaces from a single struct (e.g. `MemoryPurchaseRepo` implements Supplier, PurchaseOrder, GRN, SupplierInvoice, APTransaction, CostAllocation repos). This is why `main.go` passes the same repo for all params.

## Service Wiring

`NewService()` in `internal/service/service.go:209` takes **16 repository interfaces** (all GL+Cash related). Each module-level service (Company, Tax, Purchase, Bank, Sale, Warehouse, FA) has its own `New*Service()` with its own repos. No DI framework — manual wiring in `main.go`.

## Handler Files

`handler.go` contains only: `Handler` struct, `NewHandler`, `RegisterRoutes`, `RegisterRoutesWithCompany`. Domain-specific handlers split into separate files:

| File | Domain |
|------|--------|
| `auth_handler.go` | Login, refresh, TOTP, 2FA, sessions |
| `account_handler.go` | Account CRUD, freeze, balance, usage |
| `journal_handler.go` | Journal entry CRUD + workflow |
| `report_handler.go` | Trial balance, balance sheet, income statement |
| `period_handler.go` | Period CRUD, close, reopen |
| `exchange_rate_handler.go` | Exchange rate CRUD |
| `audit_handler.go` | Audit log queries |
| `user_handler.go` | User CRUD, current user |
| `coa_handler.go` | COA approvals, versions, analysis, mappings, IFRS |
| `opening_balance_handler.go` | Opening balance CRUD, carry forward, circular99 mappings, balance migration |
| `payroll_handler.go` | Salary calculation, payslips, PIT, declarations |
| `casbin_register.go` | `RegisterRoutesWithCompanyOpt` — central route registration with extra variadic handlers |

## Validate Package

`internal/validate/` — singleton `go-playground/validator/v10` instance with custom validators for domain enum types: `fastatus`, `damethod`, `fasource`, `fatrtype`, `disposaltype`.

Validation flow: service calls `validate.FixedAsset(a)` or `validate.FixedAssetCategory(c)` which sets defaults, runs struct tag validation, maps `ValidationErrors` to domain errors. Domain still exports `Validate()` methods but service uses validate package instead.

When adding a new module that uses the validator: register custom validators in `internal/validate/validator.go`, add `validate` struct tags to domain models, create module validation function in `internal/validate/`.

## Frontend

**No build step.** Tailwind CSS v4 compiled to `web/static/css/app.css`. Alpine.js bundled at `web/static/js/alpine.min.js`. No webpack, no Vite.

**Two UI entry points:**
- `/app/*` — Main accounting UI (40 pages). MISA SME 2026 layout: left sidebar with 12 module groups, top bar, content area.
- `/payroll/*` — Payroll module (8 pages). Separate nav bar (top).
- `/login` — Auth pages at `web/auth/*.html`.

**Shared JS:**
- `web/static/js/app.js` — Global sidebar nav (`MODULE_GROUPS`), API client (`apiGet`/`apiPost`/`apiPut`/`apiDelete` with JWT refresh), formatters, status badges, Alpine store.
- `web/static/js/auth.js` — JWT store, token refresh, alert system.
- `web/static/js/payroll.js` — Payroll-specific nav + API client.

**Static routes in `main.go`:**
```go
r.Static("/assets", "./web/static")   // CSS, JS, images
r.Static("/payroll", "./web/payroll") // Payroll pages
r.Static("/app", "./web/app")         // Main app pages
```

**Page pattern:** Each HTML page is standalone with Alpine.js `x-data` component. Calls `mountAppShell(title, activePath)` on init for sidebar/topbar. Uses `apiGet`/`apiPost` from app.js.

**Wired to backend:** All GL, Auth, Company, Cash, Bank, Purchase, Sale, FA, Tax (partial), Warehouse (partial) pages call real APIs.
**Static/mock data:** Customers, Items, VAT Report, Cash Flow — awaiting backend completion.

## Module Readiness

| Module | Status | Notes |
|--------|--------|-------|
| GL | PROD | Core accounts, journal, periods, reports, COA, opening balances |
| Auth | PROD | Login, JWT (RS256), TOTP, refresh, rate limit, lockout |
| Company | PROD | Company, branches, departments, employees, fiscal years, bank accounts |
| Cash | PROD | Receipts, payments, transfers, petty cash, advances |
| Bank | PROD | Statements, reconciliation, payment orders, loans, term deposits |
| FA | PROD | Full CRUD, depreciation engine (SL/DB), business ops, allocations, inventory |
| Tax | ~90% | Full CRUD + 17 declaration types, payment automation (GL journal, late interest, calendar sync), reconciliation, penalty calc, batch generation. CIT advanced (incentive reduction, loss carry-forward 5yr, thin cap 30% EBITDA, R&D 200% super-deduction, quarterly provisionals). E-invoice lifecycle (replacement, cancellation, auto-post GL). Calendar expanded (PIT, TTDB, BVMT, FCT), deadline alerts. Missing: form XML gen, GDT API push |
| Purchase | PROD (P2 closed) | Full domain models + repos + service + handlers + 55 routes. Missing: GDT API push, supplier portal |
| Sale | PROD (P1 closed) | Full O2C: customers, SQ, SO, DN, invoices, receipts, CNs, AR txn. Missing: e-invoice TXML, GDT push |
| Warehouse | ~30% | Interface + PG + memory repos. Service incomplete |
| Payroll | ~50% | Salary engine, SI/HI/UI calc, PIT, timekeeping, payslips, GL posting, declarations. UI complete (8 pages). Missing: approval workflow |

Full tax/purchase/warehouse/payroll specs at `docs/`.

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
- **Web UI** at `web/app/*.html` (main) and `web/payroll/*.html` (payroll). Served by gin. Not a SPA.
- **PDF** generation with `maroto/v2` for opening balance reports. Not wired to any endpoint yet.
- **Handler error mapping** uses `errors.Is` to switch on domain errors per module. See `internal/handler/fixed_asset_handler.go:79` (the `faError` func) for pattern.
- **GDT env vars**: `GDT_BASE_URL`, `GDT_TOKEN` — e-invoice API client. Without them, `newGDTClient()` returns nil.
- **E-invoice signing**: `EINVOICE_SIGNING_KEY` (RSA PEM), `EINVOICE_CERT_SERIAL`. Without key, ephemeral key generated at startup — signatures invalid after restart.
- **Migration idempotency**: Always use `CREATE TABLE IF NOT EXISTS` / `CREATE INDEX IF NOT EXISTS` in migrations. Legacy `.sql` files (002-007) may have been run manually against PG, creating tables outside migration tracking. Without `IF NOT EXISTS`, golang-migrate fails with `relation already exists`.
- **SQLite rejected**: Adding SQLite support does NOT make sense. Schema deeply PG-specific (`pgcrypto`, `gen_random_uuid()`, `TIMESTAMPTZ`, `plpgsql` triggers, `SUBSTRING FROM` regex, `GREATEST()`). 20+ raw SQL queries in PG repos. 18 migrations would need dual sets. In-memory backend already covers dev/test.
- **No LINQ library needed**: Go has `samber/lo` (19k stars, 120+ functions), `ahmetb/go-linq` (3.6k stars, v4 with iter.Seq), `CreateLab/GLinq` (lazy, fluent). For this codebase, standard `for` loops + `slices.SortFunc` are sufficient. Don't add collection libraries unless profiling shows bottleneck.

## Module Inventory (22 modules, ~498 routes)

| Module | Prefix | Routes | Handler Lines | Status |
|--------|--------|--------|---------------|--------|
| GL/Accounts | `/api/v1/accounts` | 5 | 1087 (shared) | PROD |
| Journal Entries | `/api/v1/journal-entries` | 8 | ↑ | PROD |
| Reports | `/api/v1/reports` | 3 | ↑ | PROD |
| Opening Balances | `/api/v1/opening-balances` | 13 | 306 | PROD |
| Carry Forward | `/api/v1/carry-forward` | 3 | ↑ | PROD |
| Circular 99 Mappings | `/api/v1/circular99-mappings` | 3 | ↑ | PROD |
| Balance Migrations | `/api/v1/balance-migrations` | 3 | ↑ | PROD |
| Periods | `/api/v1/periods` | 5 | ↑ | PROD |
| Exchange Rates | `/api/v1/exchange-rates` | 2 | ↑ | PROD |
| Audit Log | `/api/v1/audit` | 2 | ↑ | PROD |
| Users | `/api/v1/users` | 3 | ↑ | PROD |
| COA Management | `/api/v1/coa/*` | 14 | ↑ | PROD |
| User Auth (2FA/sessions) | `/api/v1/auth/*` | 10 | ↑ | PROD |
| Company | `/api/v1/companies` | 48 | 714 | PROD |
| Tax | `/api/v1/tax` | 53 | 1656 | ~90% |
| Cash | `/api/v1/cash` | 40 | 639 | PROD |
| Bank | `/api/v1/bank` | 37 | 563 | PROD |
| Purchase (P2P) | `/api/v1/purchase` | 55 | 870 | PROD |
| Sale (O2C) | `/api/v1/sale` | 52 | 890 | PROD |
| Warehouse | `/api/v1/warehouse` | 48 | 729 | ~30% |
| Fixed Assets | `/api/v1/fixed-assets` | 34 | 541 | PROD |

## Missing Modules (SME Gap Analysis)

Compared to MISA SME (19 subsystems, Vietnam market leader) and standard SME ERP requirements:

| Module | Priority | Why Missing |
|--------|----------|-------------|
| **Payroll** | CRITICAL | Spec complete at `docs/payroll/`. Full build needed: salary engine, SI/HI/UI, PIT, timekeeping, payslips, GL posting, D02-TS/05/KK-TNCN declarations, approval workflow. Est. 12-16 weeks. |
| **Tools & Equipment (CCDC)** | HIGH | Vietnamese accounting standard requires separate tracking from fixed assets. Account code 153. |
| **Cost Accounting** | HIGH | Product/job costing for manufacturing/construction SMEs. |
| **Budget Management** | HIGH | Budget planning by department, actual vs. budget variance. |
| **Contracts** | MEDIUM | Sales/purchase contract tracking, renewal alerts. |
| **CRM** | MEDIUM | Customer interaction history, pipeline. |
| **Notifications** | MEDIUM | Invoice due dates, approval reminders, low stock alerts. |
| **Recurring Entries** | MEDIUM | Auto-generate monthly rent, depreciation entries. |
| **Cash Flow Forecasting** | MEDIUM | Projected cash position from receivables/payables. |
| **Data Import/Export** | MEDIUM | Bulk Excel import for initial setup. |

**Verdict**: GoTax covers ~65% of Vietnamese SME needs. Without payroll, cannot be standalone.
