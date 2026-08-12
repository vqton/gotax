# GoTax — Agent Guide

## Mode

- **Caveman always active** — drop articles/filler/pleasantries/hedging. Fragments OK. Code/commits normal.
- **Karpathy** — ship first, polish later. Simple > clever. No comments/docstrings unless API surface. No speculative generality.

## Project

Vietnamese tax-compliant General Ledger API. Circular 99/2025/TT-BTC, Decree 123/2020/ND-CP. Multi-tenant, multi-company.

**Stack:** Go 1.26.5 · Gin v1.12 · GORM v1.31 (PostgreSQL via pgx v5) · golang-jwt v5 (RS256) · bcrypt · TOTP · golang-migrate v4 · go-playground/validator v10 · maroto/v2 (PDF) · zap · viper · casbin · go-i18n · swaggo/swag · testify · Alpine.js + htmx (frontend; htmx pages: Flowbite Core v4 JS + hand-rolled app.css)

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
  - tax_service.go          →  tax calculation, e-invoice lifecycle, reconciliation
  - tax_declaration_service.go →  declaration auto-population (3-tier CIT, PIT 25.1)
  - bank_import_service.go  →  bank CSV import (VCB, BIDV, CTG, VTB, ACB parsers)
  - year_end_service.go     →  year-end close (Revenue/Expense → 421, carry-forward, TT200→TT99 mapping)
  - print_service.go        →  PDF generation (Phiếu thu/chi TT99 format)
internal/validate/          →  go-playground/validator singleton + custom validators
internal/web/               →  htmx server-rendered pages (dashboard/users/journal-entries/coa) + /app/* catch-all
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

Dev seed login: `admin` / `Admin@123456!` (migration 000040). Login rate-limited: 5 attempts / 15 min / username, in-memory — server restart clears it.

## Company ID Pattern

Company-scoped handlers pass `company_id` as **query param**, not from JWT context:
```go
companyID := c.Query("company_id")
```
Same pattern across all module handlers. **Not** `c.GetString("company_id")`.

## Routes

**Two registration patterns:**

**Pattern A — `RegisterRoutesWithCompany`** (handler.go:214): Core modules that share company-scope wiring go through this umbrella function. Signature includes all handlers as params.

**Pattern B — Direct registration** (main.go): Newer modules call `handler.Register*Routes(r, h, authMW)` directly in main.go. Used for modules added after the initial architecture.

**When adding a new module:**
1. Create `internal/handler/<module>_handler.go` with `New*Handler` + `Register*Routes`
2. Wire in **both** PG and memory branches of `main.go` (new repos + new service + new handler)
3. Call `handler.Register*Routes(r, handler, authMW)` directly in main.go (Pattern B)
4. **Do NOT** add params to `RegisterRoutesWithCompany` unless the module truly belongs in the core company-scope umbrella

| Handler | Pattern | Prefix |
|---------|---------|--------|
| `Handler` | A | `/api/v1/{auth,accounts,journal-entries,periods,exchange-rates,reports,coa/*,audit,users,me,print,...}` |
| `CompanyHandler` | A | `/api/v1/companies/**` |
| `TaxHandler` | A | `/api/v1/tax/**` |
| `CashHandler` | A | `/api/v1/cash/**` |
| `BankHandler` | A | `/api/v1/bank/**` |
| `PurchaseHandler` | A | `/api/v1/purchase/**` |
| `SaleHandler` | A | `/api/v1/sale/**` |
| `WarehouseHandler` | A | `/api/v1/warehouse/**` |
| `FAHandler` | A | `/api/v1/fixed-assets/**` |
| `PayrollHandler` | A | `/api/v1/payroll/**` |
| `RecurringHandler` | A | `/api/v1/recurring-entries/**` |
| `BudgetHandler` | A | `/api/v1/budgets/**` |
| `CCDCHandler` | A | `/api/v1/ccdc/**` |
| `CostCenterHandler` | A | `/api/v1/cost-centers/**` |
| `WarehouseKeeperHandler` | A | `/api/v1/warehouse/keeper/**` |
| `CostObjectHandler` | A | `/api/v1/cost-objects/**` |
| `CostPoolHandler` | A | `/api/v1/cost-pools/**` |
| `CostingHandler` | A | `/api/v1/costing/**` |
| `CostReportHandler` | B | `/api/v1/cost-reports/**` |
| `SystemOptionHandler` | B | `/api/v1/system-options/**` |
| `NumberingRuleHandler` | B | `/api/v1/numbering-rules/**` |
| `FiscalYearHandler` | B | `/api/v1/fiscal-years/**` |
| `ReportOptionHandler` | B | `/api/v1/report-options/**` |
| `BackupHandler` | B | `/api/v1/backups/**` |
| `ContractHandler` | B | `/api/v1/contracts/**` |
| `BankImportHandler` | B | `/api/v1/bank-imports/**` |
| `InvoiceBookHandler` | B | `/api/v1/invoice-books/**` |
| `NotificationHandler` | B | `/api/v1/notifications/**` |
| `PriceListHandler` | B | `/api/v1/price-lists/**` |
| `FinancialAnalysisHandler` | B | `/api/v1/financial-analysis/**` |

Auth middleware on all groups except auth endpoints.

**Web pages** (GET `/app/*`, auth-gated via internal/web catch-all): `dashboard`, `users`, `journal-entries`, `coa` are htmx server-rendered; every other file in `web/app/*.html` (e.g. `customers.html`) is served statically and rendered client-side with Alpine. See Frontend.

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

45 versioned files in `migrations/` (000001–000045). **Versioned** (`000001_title.up.sql` + `000001_title.down.sql`) → auto-discovered by golang-migrate and auto-run on PG startup. Current latest: `000045_tt99_coa_loai2`.

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

**Two rendering stacks, no build step:**

1. **htmx server-rendered** (newer, `internal/web/`): converted pages render Go templates from `web/templates/` (base + `_sidebar` + `_topbar` partials, parsed per-page at startup). Mutations via explicit POST routes + htmx fragment swaps. Currently `dashboard`, `users`, `journal-entries`, `coa`. Adding a page = template in `web/app/<page>.html` (defines "content" block) + `Load` func in `internal/web/pages_app.go` + add to `web.NewServer([...])` list in **both** main.go branches.
2. **Alpine.js legacy** (everything else, e.g. `customers.html`): standalone page, `x-data` per page, `mountAppShell(title, activePath)` on init, API via `apiGet`/`apiPost`/`apiPut`/`apiDelete` (JWT refresh). Uses `app-legacy.js` + `auth-legacy.js`; htmx pages use `app.js` + `auth.js`.

`/app/*` catch-all (`internal/web/pages.go`): page in template sets → server-render; otherwise `http.ServeFile` from `web/app/*.html` — served fresh from disk per request, so **HTML edits to static pages need no server restart**.

**CSS: hand-rolled design system, no framework.** Tailwind v4 removed (Aug 2026); its utility classes in HTML are inert. `web/static/css/app.css` (Aug 2026) is the htmx-page design system: tokens, dark slate sidebar, topbar, dropdowns, buttons, cards, KPI grid, badges, tables, status bars, toasts, responsive breakpoints. All styling via semantic classes (`.badge-posted`, `.kpi`, `.btn-primary`…) — **do not use Tailwind utility classes in htmx templates**. Legacy Alpine pages still use their own CSS (e.g. `customers.html` links `customers.css`); only `web/static/css/auth.css` remains from the pre-Tailwind era.

**Flowbite Core (v4):** vendored at `web/static/js/flowbite.min.js` (v4.0.2 — the new framework-agnostic "Core" line; `@flowbite/core` is NOT on npm). Behaviors only (dropdowns, modals, collapses, tooltips, tabs) via `data-*` attributes + auto-init; **no Flowbite CSS classes used** — pages styled by app.css. Global `window.initFlowbite()` re-inits all components (used by app.js `htmx:afterSwap` hook for fragment swaps). Auto-init runs once at script load — components injected by htmx swaps need the re-init hook. Collapse toggles `.hidden` on the target element, NOT a class on the trigger — chevron rotation in app.css uses `.sb-group:has(.sb-links:not(.hidden))`.

**Static assets** served at `/assets/` (→ `web/static/`, `Cache-Control: no-cache` forces revalidation). Script order in base.html: htmx.min.js → flowbite.min.js → app.js. Shell data (`.Shell.*`: Title, NavPath, Username, RoleLabel, AvatarInit, CompanyName, PeriodLabel) built in `internal/web/shell.go`; page data under `.Data.*`.

**UI entry points:**
- `/app/*` — Main accounting UI. MISA SME 2026 layout: left sidebar, top bar, content area.
- `/payroll/*` — Payroll module. Separate nav bar (top).
- `/login` — Auth pages at `web/auth/*.html` (plain gin templates).

## Adding a Feature — Step Order

1. Interface method in `internal/domain/interfaces.go`
2. GORM model in `internal/domain/models_gorm_*.go`
3. Repository impl in `internal/repository/pg_*.go` + `memory_*.go`
4. Service method in `internal/service/<module>_service.go`
5. Handler method in `internal/handler/<module>_handler.go` + `Register*Routes`
6. Wire in `main.go` (both PG + memory branches: new repo, new service, new handler, register route)
7. Tests in `internal/handler/<module>_handler_test.go`
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
- **SQLite rejected**: Schema deeply PG-specific (`pgcrypto`, `gen_random_uuid()`, `TIMESTAMPTZ`, `plpgsql` triggers, `SUBSTRING FROM` regex, `GREATEST()`). 20+ raw SQL queries in PG repos. In-memory backend already covers dev/test.
- **No LINQ library needed**: Standard `for` loops + `slices.SortFunc` sufficient. Don't add collection libraries unless profiling shows bottleneck.
- **Server start**: `setsid nohup go run . & disown` — plain `nohup &` dies when the launching shell/command times out (process group killed). Templates resolve relative to cwd → start from repo root. Needs `JWT_SECRET` (+ `DATABASE_URL` for PG branch).
- **COA structure**: accounts = 9 loại (1-char grouping headers, NOT postable), cấp 1 = 3-digit. `NormalBalance` derived, not stored: ASSET/EXPENSE→DEBIT, else CREDIT. 221 accounts, 71 cấp 1. Service guards: `CreateAccount` requires parent exists, `UpdateAccount` runs `Validate()`, `DeleteAccount` blocked when account used in POSTED entries (`GetAccountUsage().EntryCount`).
- **Playwright (this host)**: `chrome` channel uninstallable on Kali — symlink cached chromium: `ln -sf ~/.cache/ms-playwright/chromium-1234/chrome-linux64/chrome /opt/google/chrome/chrome`. Screenshots must live under repo root (`.playwright-mcp/` allowed, `/tmp` denied).

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
