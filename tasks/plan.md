# Implementation Plan: GoTax Remaining Features

## Overview

Complete all missing modules and gaps in GoTax Vietnamese ERP. From core (tax, warehouse, payroll) to edge (CCDC, budget, cost accounting). Simple → complex. TDD. One slice at a time, end-to-end.

## Current State (after survey)

| Module | Before | Gap |
|--------|--------|-----|
| Tax | ~95% | Rate lookup for TTDB/BVMT/NTNN, SoTienBangChu, buyer/seller info |
| Warehouse | ~30% | Backend 100% complete. Frontend missing 8 sub-pages |
| Payroll | ~50% | SubmitPeriod stub, no GL posting, no tax module link |
| Recurring Entries | 0% | Everything |
| Notifications | ~10% | TaxAlert only. Need general system |
| Cash Flow | ~20% | Static HTML mock. Need backend |
| Data Import/Export | ~25% | 3 imports exist. Need export framework |
| CCDC | 0% | Everything |
| Budget | 0% | Everything |
| Cost Accounting | ~15% | Purchase cost allocation only. Need cost center CRUD |
| Contracts | ~10% | Fields only. Defer — low value for accounting system |
| CRM | 0% | Defer — not core accounting |

## Architecture Decisions

1. **Follow existing patterns exactly** — handler→service→repository, PG + memory dual impl
2. **Domain models in `internal/domain/models_*.go`** — all same package, no sub-packages
3. **Two backends always** — every new repo gets PG + memory impl
4. **Migration per feature** — `{next_version}_{title}.up.sql` + `.down.sql`
5. **Frontend per page** — standalone HTML with Alpine.js, `mountAppShell()`, `apiGet`/`apiPost`
6. **TDD** — write test alongside implementation, `go test -count=1 ./...` before commit
7. **Defer Contracts + CRM** — not core accounting, revisit later

## Task List

### Phase 1: Tax Completion (~95% → 100%)

#### Task 1.1: Tax Rate Lookup for TTDB/BVMT/NTNN
- **What**: Replace hardcoded rate=0 in `resourceTaxDeclarationLines()` and `ntnnDeclarationLines()` with actual rate lookup from `TaxRate` table
- **How**: Service calls `GetActiveRates()` filtering by tax type, matches rate to declaration line items
- **Files**: `tax_service.go:1701-1724`, add rate lookup helper
- **Verify**: `go test -count=1 ./internal/service/ -run TestGenerateDeclaration`
- **Size**: S (1-2 files)

#### Task 1.2: SoTienBangChu (Amount in Words)
- **What**: Vietnamese number-to-words converter for e-invoice TXML
- **How**: Pure function `AmountInWords(amount int64) string` in `internal/einvoice/`, populate `SoTienBangChu` field in `GenerateTXML()`
- **Files**: `internal/einvoice/txml.go:69`, new `internal/einvoice/words.go` + `words_test.go`
- **Verify**: `go test -count=1 ./internal/einvoice/ -run TestAmountInWords`
- **Size**: S (2 files)

#### Task 1.3: Buyer/Seller Info in E-Invoice
- **What**: Wire company profile (seller) and customer (buyer) into TXML generation
- **How**: `GenerateTXML()` takes company + customer params, populates `BKParty` fields
- **Files**: `internal/einvoice/txml.go`, `internal/service/tax_service.go:911`
- **Verify**: `go test -count=1 ./internal/einvoice/ -run TestGenerateTXML`
- **Size**: S (2 files)

**Checkpoint 1**: Tax 100%. `go test -count=1 ./...` passes. Commit.

---

### Phase 2: Warehouse Frontend (Backend done, frontend missing)

All pages follow established pattern: standalone HTML, Alpine.js `x-data`, `mountAppShell()`, `apiGet`/`apiPost`/`apiPut`/`apiDelete`.

#### Task 2.1: Item Categories Page
- **What**: CRUD UI for item categories (code, name, description)
- **Files**: `web/app/warehouse-categories.html`
- **Verify**: Page loads, CRUD works against `/warehouse/categories`
- **Size**: S (1 file)

#### Task 2.2: Items Page
- **What**: CRUD UI for warehouse items (code, name, category, unit, valuation method, min/max stock)
- **Files**: `web/app/warehouse-items.html`
- **Verify**: Page loads, CRUD works against `/warehouse/items`
- **Size**: M (1 file, many fields)

#### Task 2.3: Stock Balances Page
- **What**: Read-only view of stock balances per item per warehouse, with search/filter
- **Files**: `web/app/warehouse-balances.html`
- **Verify**: Page loads, shows balances from `/warehouse/balances`
- **Size**: S (1 file)

#### Task 2.4: Stock Transfers Page
- **What**: Create/submit/approve/transfer/complete/cancel workflow UI
- **Files**: `web/app/warehouse-transfers.html`
- **Verify**: Full workflow works end-to-end
- **Size**: M (1 file, workflow buttons)

#### Task 2.5: Stock Adjustments Page
- **What**: Create/submit/approve/post/reject workflow UI
- **Files**: `web/app/warehouse-adjustments.html`
- **Verify**: Full workflow works end-to-end
- **Size**: M (1 file)

#### Task 2.6: Stock Takes Page
- **What**: Create/start/complete/verify/post workflow UI with count entry
- **Files**: `web/app/warehouse-takes.html`
- **Verify**: Full workflow works end-to-end
- **Size**: M (1 file)

#### Task 2.7: Valuation Runs Page
- **What**: Create/run valuation runs, show results
- **Files**: `web/app/warehouse-valuations.html`
- **Verify**: Valuation run executes and shows results
- **Size**: S (1 file)

#### Task 2.8: Inventory Transactions Page
- **What**: Read-only paginated log of all inventory transactions with filters
- **Files**: `web/app/warehouse-transactions.html`
- **Verify**: Page loads, shows transactions from `/warehouse/transactions`
- **Size**: S (1 file)

**Checkpoint 2**: Warehouse 100% (backend + frontend). All 9 sub-pages functional. Sidebar nav updated. Commit.

---

### Phase 3: Payroll Completion (~50% → 90%)

#### Task 3.1: SubmitPeriod Workflow
- **What**: Implement the stub `SubmitPeriod` handler — transitions period DRAFT→PROCESSING, validates required fields
- **Files**: `internal/handler/payroll_handler.go:242`, `internal/service/payroll_service.go`
- **Verify**: `go test -count=1 ./internal/handler/ -run TestSubmitPeriod`
- **Size**: S (2 files)

#### Task 3.2: Payroll GL Posting
- **What**: When period APPROVED, auto-create journal entries: salary expense (642), insurance expense (642), employer insurance (642), PIT payable (3331), bank payable (331)
- **Files**: `internal/service/payroll_service.go`, needs `JournalService` injection
- **Verify**: `go test -count=1 ./internal/service/ -run TestApprovePeriod`
- **Size**: M (2-3 files, service + test)

#### Task 3.3: Payroll→Tax KK_TNCN Integration
- **What**: Approved payroll period feeds employee data into Tax `GenerateDeclaration` for KK_TNCN/QTT_TNCN
- **Files**: `internal/service/payroll_service.go`, `internal/service/tax_service.go`
- **Verify**: `go test -count=1 ./internal/service/ -run TestPayrollTaxIntegration`
- **Size**: S (2 files)

**Checkpoint 3**: Payroll ~90%. Approval workflow, GL posting, tax link all work. Commit.

---

### Phase 4: Recurring Entries (0% → 100%)

#### Task 4.1: Domain Model + Interfaces
- **What**: `RecurringEntry` model (template JE with frequency, next_run, end_date), repository interface
- **Files**: `internal/domain/models_recurring.go`, `internal/domain/interfaces.go`
- **Verify**: `go build ./...`
- **Size**: S (2 files)

#### Task 4.2: PG + Memory Repos
- **What**: Dual implementation following existing pattern
- **Files**: `internal/repository/pg_recurring.go`, `internal/repository/memory_recurring.go`
- **Verify**: `go test -count=1 ./internal/repository/`
- **Size**: M (2 files)

#### Task 4.3: Migration
- **What**: `000026_recurring_entries.up.sql` + `.down.sql`
- **Files**: `migrations/000026_recurring_entries.*`
- **Verify**: `go build ./...`
- **Size**: XS (2 files)

#### Task 4.4: Service + Handler + Routes
- **What**: CRUD + "run now" endpoint that generates JE from template
- **Files**: `internal/service/recurring_service.go`, `internal/handler/recurring_handler.go`
- **Verify**: `go test -count=1 ./internal/handler/ -run TestRecurringEntry`
- **Size**: M (2-3 files)

#### Task 4.5: Frontend
- **What**: Recurring entries list + create/edit modal
- **Files**: `web/app/recurring-entries.html`
- **Verify**: Page loads, CRUD + run now works
- **Size**: S (1 file)

#### Task 4.6: Wire in main.go
- **What**: PG + memory branches, route registration
- **Files**: `main.go`
- **Verify**: `go test -count=1 ./...`
- **Size**: S (1 file)

**Checkpoint 4**: Recurring entries module complete. Commit.

---

### Phase 5: Notifications (10% → 80%)

#### Task 5.1: Domain Model + Interfaces
- **What**: `Notification` model (entity_type, entity_id, user_id, title, message, read_at, link), repository interface
- **Files**: `internal/domain/models_notification.go`, `internal/domain/interfaces.go`
- **Verify**: `go build ./...`
- **Size**: S (2 files)

#### Task 5.2: PG + Memory Repos
- **What**: Dual implementation
- **Files**: `internal/repository/pg_notification.go`, `internal/repository/memory_notification.go`
- **Verify**: `go test -count=1 ./internal/repository/`
- **Size**: M (2 files)

#### Task 5.3: Migration
- **What**: `000027_notifications.up.sql` + `.down.sql`
- **Files**: `migrations/000027_notifications.*`
- **Verify**: `go build ./...`
- **Size**: XS (2 files)

#### Task 5.4: Service + Handler + Routes
- **What**: CRUD + mark-read + unread-count + fire-on-event (invoice due, approval pending, low stock)
- **Files**: `internal/service/notification_service.go`, `internal/handler/notification_handler.go`
- **Verify**: `go test -count=1 ./internal/handler/ -run TestNotification`
- **Size**: M (2-3 files)

#### Task 5.5: Frontend
- **What**: Notification bell in topbar with dropdown, unread count badge
- **Files**: `web/static/js/app.js` (add bell), `web/app/notifications.html` (full list page)
- **Verify**: Bell shows count, clicking opens dropdown, full page shows all
- **Size**: S (2 files)

#### Task 5.6: Wire in main.go
- **Files**: `main.go`
- **Verify**: `go test -count=1 ./...`
- **Size**: S (1 file)

**Checkpoint 5**: Notifications module complete. Commit.

---

### Phase 6: Cash Flow (20% → 80%)

#### Task 6.1: Backend API
- **What**: `GET /api/v1/reports/cash-flow?company_id=&year=&month=` — aggregate journal entries into operating/investing/financing sections using Vietnamese account code ranges (111, 112, 128, 228, 331, 333, 511, 642, etc.)
- **Files**: `internal/service/report_service.go` (add `CashFlowStatement`), `internal/handler/report_handler.go` (add route)
- **Verify**: `go test -count=1 ./internal/handler/ -run TestCashFlow`
- **Size**: M (2 files)

#### Task 6.2: Frontend Integration
- **What**: Replace hardcoded mock data in `web/app/cash-flow.html` with real API calls
- **Files**: `web/app/cash-flow.html`
- **Verify**: Page shows real data from journal entries
- **Size**: S (1 file)

**Checkpoint 6**: Cash flow statement works end-to-end with real data. Commit.

---

### Phase 7: Data Import/Export (25% → 70%)

#### Task 7.1: Journal Entry Excel Export
- **What**: `GET /api/v1/journal-entries/export?company_id=&period=` — export filtered journal entries to .xlsx
- **Files**: `internal/service/export_service.go`, `internal/handler/report_handler.go`
- **Verify**: `go test -count=1 ./internal/service/ -run TestExportJournalEntries`
- **Size**: M (2 files)

#### Task 7.2: Trial Balance Excel Export
- **What**: `GET /api/v1/reports/trial-balance/export?company_id=&period=`
- **Files**: `internal/service/export_service.go`, `internal/handler/report_handler.go`
- **Verify**: `go test -count=1 ./internal/service/ -run TestExportTrialBalance`
- **Size**: S (1 file, extends export_service)

**Checkpoint 7**: GL export works. Commit.

---

### Phase 8: CCDC (Tools & Equipment) (0% → 100%)

Similar to Fixed Assets module but for tools/equipment (account 153). Simpler — no depreciation engine needed initially.

#### Task 8.1: Domain Model + Interfaces
- **What**: `ToolEquipment` model (code, name, category, purchase_date, cost, status, warehouse_id), repository interface
- **Files**: `internal/domain/models_ccdc.go`, `internal/domain/interfaces.go`
- **Verify**: `go build ./...`
- **Size**: S (2 files)

#### Task 8.2: PG + Memory Repos
- **Files**: `internal/repository/pg_ccdc.go`, `internal/repository/memory_ccdc.go`
- **Verify**: `go test -count=1 ./internal/repository/`
- **Size**: M (2 files)

#### Task 8.3: Migration
- **Files**: `migrations/000028_ccdc.up.sql`, `migrations/000028_ccdc.down.sql`
- **Verify**: `go build ./...`
- **Size**: XS (2 files)

#### Task 8.4: Service + Handler + Routes
- **Files**: `internal/service/ccdc_service.go`, `internal/handler/ccdc_handler.go`
- **Verify**: `go test -count=1 ./internal/handler/ -run TestCCDC`
- **Size**: M (2-3 files)

#### Task 8.5: Frontend
- **Files**: `web/app/ccdc.html`
- **Verify**: CRUD works
- **Size**: S (1 file)

#### Task 8.6: Wire in main.go
- **Files**: `main.go`
- **Verify**: `go test -count=1 ./...`
- **Size**: S (1 file)

**Checkpoint 8**: CCDC module complete. Commit.

---

### Phase 9: Budget Management (0% → 80%)

#### Task 9.1: Domain Model + Interfaces
- **What**: `Budget` model (account_id, period, amount, actual, variance), `BudgetLine` model, repository interface
- **Files**: `internal/domain/models_budget.go`, `internal/domain/interfaces.go`
- **Verify**: `go build ./...`
- **Size**: S (2 files)

#### Task 9.2: PG + Memory Repos
- **Files**: `internal/repository/pg_budget.go`, `internal/repository/memory_budget.go`
- **Verify**: `go test -count=1 ./internal/repository/`
- **Size**: M (2 files)

#### Task 9.3: Migration
- **Files**: `migrations/000029_budget.up.sql`, `migrations/000029_budget.down.sql`
- **Verify**: `go build ./...`
- **Size**: XS (2 files)

#### Task 9.4: Service + Handler + Routes
- **What**: CRUD + budget vs actual comparison report
- **Files**: `internal/service/budget_service.go`, `internal/handler/budget_handler.go`
- **Verify**: `go test -count=1 ./internal/handler/ -run TestBudget`
- **Size**: M (2-3 files)

#### Task 9.5: Frontend
- **Files**: `web/app/budget.html`
- **Verify**: CRUD + variance report works
- **Size**: M (1 file, report view)

#### Task 9.6: Wire in main.go
- **Files**: `main.go`
- **Verify**: `go test -count=1 ./...`
- **Size**: S (1 file)

**Checkpoint 9**: Budget module complete. Commit.

---

### Phase 10: Cost Accounting (15% → 60%)

#### Task 10.1: Cost Center CRUD
- **What**: `CostCenter` model (code, name, parent_id for hierarchy), CRUD service/handler/routes
- **Files**: `internal/domain/models_cost.go`, `internal/repository/pg_cost.go`, `internal/repository/memory_cost.go`, `internal/service/cost_service.go`, `internal/handler/cost_handler.go`, migration, `main.go`
- **Verify**: `go test -count=1 ./...`
- **Size**: L (7+ files, but each is small)

#### Task 10.2: Cost Center Frontend
- **Files**: `web/app/cost-centers.html`
- **Verify**: CRUD + hierarchy view works
- **Size**: S (1 file)

**Checkpoint 10**: Cost center management works. Commit.

---

### Phase 11: Final Integration + AGENTS.md Update

#### Task 11.1: Update AGENTS.md
- **What**: Update all module statuses, route counts, missing modules table
- **Files**: `AGENTS.md`
- **Verify**: Accurate reflection of current state
- **Size**: S (1 file)

#### Task 11.2: Full Test Suite
- **What**: `go vet ./... && go test -count=1 ./...` — everything green
- **Verify**: All tests pass, no regressions
- **Size**: verification only

**Checkpoint 11**: All complete. Final commit.

---

## Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Payroll GL posting complexity | Med | Follow existing `postGLEntry()` pattern from warehouse |
| Cash flow account code mapping | Med | Use Circular 99 account ranges, test with sample data |
| Migration conflicts | Low | Always use `IF NOT EXISTS`, check latest version first |
| Frontend consistency | Low | Copy pattern from existing pages verbatim |

## Open Questions

- None — all patterns well-established in codebase
