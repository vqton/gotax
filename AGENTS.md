# GoTax — Agent Guide

## Mode

- **Caveman always active** — drop articles (a/an/the), filler, pleasantries, hedging. Fragments OK. Code/commits written normal.
- **Karpathy rules** — be a hacker: ship first, polish later. Simple > clever. Write less, ship more. Don't over-abstract. Prefer flat code over indirection.
- **Main theory** — every line must serve the feature. If it doesn't add signal, cut it. No comments, no docstrings unless API surface. No speculative generality.

## Project

Vietnamese tax-compliant General Ledger API. Built to Circular 99/2025/TT-BTC, Decree 123/2020/ND-CP. Multi-tenant, multi-company. Backend for GoTax accounting SaaS.

**Stack:** Go 1.26.5 · Gin v1.12 · pgx/v5 (PostgreSQL) · golang-jwt v5 (RS256) · bcrypt · TOTP · swaggo/swag · testify

## Architecture — Clean Architecture (migrated from flat `internal/gl/`)

```
main.go                     →  entrypoint, DI wiring, backend selection (PG vs memory)
internal/domain/            →  models, repository interfaces, errors. Zero dependencies.
internal/auth/              →  JWT (RS256), TOTP, bcrypt, rate limiter
internal/service/           →  Service (GL) + CompanyService (company). Business rules, validation.
internal/handler/           →  Handler + CompanyHandler + AuthMiddleware + route registration
internal/repository/        →  in-memory + PG impls of all domain repository interfaces
internal/db/                →  PGConfig, NewPool, RunMigrations
cmd/                        →  (empty, unused)
pkg/                        →  (empty, unused)
docs/                       →  generated swagger (DO NOT EDIT)
migrations/                 →  002_gl_schema.sql + 003_company_schema.sql
```

**Two backends** controlled by `DATABASE_URL` env var:
- Set → PostgreSQL (auto-migrates on startup via pgxpool, no version table)
- Unset → in-memory maps + sync.RWMutex (dev/test)

**Request lifecycle:**
```
HTTP → gin.Engine → authMW (JWT verify) → roleMW (RBAC) → Handler → Service → Repository
```

**Layer split:**
- `Handler` — parse request, call service, return JSON. Zero logic.
- `Service` — business rules, validation, orchestration. Pure Go, no HTTP.
- `Repository` — data access. Two impls: `PG*Repo` and `Memory*Repo`.

Every entity has both PG and memory repo. Adding a new entity means implementing the interface twice.

## Domain — Model Files (SOLID: Single Responsibility)

`internal/domain/models*.go` — all `package domain`. Split by bounded context:

| File | Lines | Contents |
|------|-------|----------|
| `models.go` | 344 | GL core (Account, JournalEntry, Period, AccountBalance) + Auth (User, RefreshToken, Session, TokenPair, AuthResult) + AuditEntry + ExchangeRate + ClosingTemplate + Password policy |
| `models_coa.go` | 111 | ApprovalRequest, AccountVersion, AccountMapping, AccountAnalysis, IFRSMapping, AccountUsage, VersionDiff/AccountDiff/Change |
| `models_company.go` | 245 | Company, Branch, FiscalYear, PeriodV2, Department, Employee, BankAccount, EInvoicePattern, DigitalSignature, IntegrationProfile, CompanyContext |
| `models_bank.go` | 327 | Statement, Recon, PaymentOrder, Batch, Loan, TermDeposit, Filters, Reports |
| `models_tax.go` | 482 | Declaration types, TaxRate, EInvoice lifecycle, TaxCalendar, AuditCase, Calc results, Filters |
| `models_cash.go` | 331 | CashReceipt/Payment/Transfer, PettyCash, Inventory, AdvanceRequest/Settlement |
| `models_ob.go` | 117 | OpeningBalance, OpeningBalanceDetail, CarryForwardLog, Circular99Mapping, BalanceMigration |

No sub-packages. Adding a model = add to correct existing file or create new `models_*.go`.

## Auth

```
POST /api/v1/auth/login       →  username+password → JWT pair (access 15m + refresh)
POST /api/v1/auth/refresh      →  rotate refresh token
POST /api/v1/auth/totp/verify  →  2FA challenge after login

authMW (gin middleware)        →  extract Bearer token → verify RS256 → set user_id, username, role in gin.Context
GetUserID(c *gin.Context)      →  helper to read user_id from context
RoleMiddleware(admin, chief)   →  gate routes by role
```

JWT uses RS256. RSA key pair generated at startup from `JWT_SECRET` seed via `crypto/rand`. **Not deterministic** — same seed ≠ same keys. Every restart invalidates all tokens.

## Routes

| Group | Path | Handler | Auth |
|-------|------|---------|------|
| Auth | `/api/v1/auth/{login,refresh,forgot-password,reset-password,totp/verify}` | `Handler` | Public |
| GL | `/api/v1/accounts`, `/journal-entries`, `/periods`, `/exchange-rates`, `/reports`, `/coa/*`, `/audit` | `Handler` | authMW |
| User | `/api/v1/users`, `/me`, `/auth/{change-password,logout,totp/*}` | `Handler` | authMW |
| Company | `/api/v1/companies/**`, `/branches/**`, `/departments/**`, `/employees/**`, `/bank-accounts/**`, `/fiscal-years/**`, `/einvoice-patterns/**`, `/digital-signatures/**`, `/integrations/**` | `CompanyHandler` | authMW |

`RegisterRoutesWithCompany(r, h, ch, authMW, adminMW)` in `handler/handler.go:153`.

## Commands

```sh
go build ./...
go test ./...
go test -v -run TestCreateCompany ./internal/handler/
go vet ./...

# Run (memory)
JWT_SECRET=devsecret go run .

# Run (PostgreSQL — auto-migrates)
DATABASE_URL=postgres://... JWT_SECRET=devsecret go run .

# Regenerate swagger (annotations in main.go + handler comments)
swag init --parseDependency --parseInternal
# → docs/docs.go, docs/swagger.json, docs/swagger.yaml — DO NOT EDIT
```

No Makefile, no linter config, no CI, no Dockerfile.

## Quality

Code quality standard at `docs/standards/CODE_QUALITY.md`: static analysis toolchain, 8-dimension review checklist, layer rules, PRR gate, quality metrics. Run `go vet ./... && go test ./...` before commit; full standard applies before merge.

## Testing

Test strategy at `docs/standards/TEST_STRATEGY.md`: test pyramid, Go idioms, coverage targets, race detection, fuzzing, CI gates, outdated practices.

- All tests use **in-memory repos** + `httptest.NewRecorder` — no DB, no integration.
- `handler_test.go` — creates real in-memory repos + real `service.Service`, gin engine with mock auth middleware (sets user_id, role). No mock services.
- `company_handler_test.go` — same pattern, uses real `service.CompanyService` + `MemoryCompanyRepo`.
- `service_test.go` — real in-memory repos, tests full journal lifecycle (draft→submit→approve→post→cancel), login/2FA/refresh/lockout, exchange rates, reports (trial balance, account balance, drill-down), COA ops.
- `domain/models_test.go` — 22 struct validation tests for all domain models.
- Test count: 160 across domain (22) + handler (72) + service (66) — all green.
- Bank module: Statements, Recon, Payment Orders, Batches, Loans, Term Deposits, Reports. 11 tables in `004_bank_module.sql`. PG + memory repos. `BankService` 35+ methods. `BankHandler` 16 endpoints. 16 handler tests.

- Adding a new service method: add a field to `mockService` in `handler_test.go`? **No** — current tests use real service. Just add the method to service, repos, and write endpoint test in handler_test.go.

## Gotchas

- **`JWT_SECRET` required** — server panics if unset.
- **RSA non-deterministic** — same `JWT_SECRET` ≠ same keys. All tokens die on restart.
- **Migrations every startup** — raw SQL exec via `db.RunMigrations`, no version table. Adding a migration = write `.sql` in `migrations/` and append path to `db/pg.go` `migrations` slice.
- **`001_gl_schema.sql.deprecated`** in `migrations/` — unused file, do not reference.
- **`GenerateToken` uses HMAC-SHA256** with hardcoded fallback secret — test-only dead code. Production uses `GenerateAccessToken` (RS256).
- **Go 1.26.5** — `range-over-func` and other modern Go features available.
- **`CompanyService` uses separate `CompanyRepository`** (not the GL service interface). Wired independently in main.go.
- **Memory repos auto-generate IDs** — they copy the struct and generate an ID for the copy, then write it back to the original pointer (`c.ID = cp.ID`). Always use the same pointer after Create.
- **`pg_company.go`** is a separate file from `pg.go` — both in `internal/repository/`. All PG repos are in those two files.
- **Domain models split by bounded context** — 7 `models_*.go` files in `internal/domain/`. All same package, zero import changes when adding a model. Add to correct file or create new `models_*.go`.

## Tax Module — NOT PROD READY

Tax documentation at `docs/tax/` — comprehensive analysis from BA Lead + Chief Accountant perspective (20+ yrs each). Researched: MISA, Fast, BravoERP, Tryton, GDT, EY, PwC, KPMG, Deloitte, gov portals.

**Verdict: Tax module ~20% complete. NOT production-ready.**

What exists (foundation only):
- Circular 99 COA tax accounts (20+ tax accounts)
- Company tax code/office fields + validation
- E-invoice pattern CRUD, digital signature CRUD, GDT integration profile stub
- `TestIntegration` is no-op

What's MISSING (~80%):
- Tax declaration engine (VAT/CIT/PIT/TTDB/BVMT/FCT forms)
- Tax rate tables (configurable with effective dates)
- Tax calculation from journal entries
- XML generation (01/GTGT, 03/TNDN, 05/KK-TNCN, etc.)
- GDT API client for `thuedientu.gdt.gov.vn`
- E-invoice issuance pipeline (TXML → sign → submit → issue)
- Tax payment tracking + reconciliation
- Tax calendar + deadline alerts
- Tax audit support
- Global minimum tax (Pillar 2) — account 82112 exists, no logic

See `docs/tax/TAX_READINESS.md` for full matrix. See `docs/tax/TAX_BRD.md`, `TAX_SPECS.md`, `TAX_USE_CASES.md`, `TAX_WORKFLOWS.md`, `TAX_TEMPLATES.md`, `TAX_RULES.md`, `TAX_DATA_FLOWS.md` for full design.

**Adding a Feature — Step Order**

1. Interface method in `internal/domain/interfaces.go`
2. Repository impl in `internal/repository/pg.go` (or `pg_company.go`) + `internal/repository/memory.go` (or `memory_company.go`)
3. Service method in `internal/service/service.go` (or `company.go`)
4. Handler method in `internal/handler/handler.go` (or `company.go`) + route registration
5. Test in `internal/handler/handler_test.go` (or `company_handler_test.go`)
6. Wire in `main.go`
7. `go vet ./... && go test ./...`
