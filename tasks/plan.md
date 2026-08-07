# Implementation Plan: Thủ Kho (Warehouse Keeper) Module

## Overview

Add a role-based Warehouse Keeper sub-module to enforce separation of duties between physical goods custody (Thủ kho) and accounting records (Kế toán kho). This is a regulatory blocker — violates Luật Kế toán 2015 Art. 6 and Circular 99/2025/TT-BTC Art. 3.

**Current state:** 0% MISA parity. Warehouse module has CRUD, GRN, DN, transfers, adjustments, takes — all PROD. But zero keeper role concepts.

**Target state:** Full MISA SME 2026 Thủ kho parity. ~10 new files, ~3000 lines, 4 weeks.

---

## Architecture Decisions

### AD-001: Separate service, not extending WarehouseService
- **Decision:** Create `WarehouseKeeperService` as new service, not add methods to `WarehouseService`
- **Rationale:** WarehouseService is already 891 lines. Keeper is a distinct bounded context with its own repository. Keeps single responsibility.
- **Pattern:** Same as how `PayrollService` is separate from `Service`

### AD-002: Separate handler file, not extending WarehouseHandler
- **Decision:** Create `WarehouseKeeperHandler` as new handler file
- **Rationale:** WarehouseHandler is 729 lines. Keeper has distinct routes under `/api/v1/warehouse/keeper/`
- **Pattern:** Same as how `PayrollHandler` is separate from `Handler`

### AD-003: New repository, not extending WarehouseRepository
- **Decision:** Create `WarehouseKeeperRepository` interface + PG + Memory implementations
- **Rationale:** Keeper data (assignments, ledger entries) is distinct from warehouse data. Separate repo keeps interfaces clean.
- **Pattern:** Follows existing per-module repo pattern (`pg_*.go` + `memory_*.go`)

### AD-004: Stock Ledger is read-only overlay
- **Decision:** StockLedgerEntry does NOT trigger GL posting. Source modules (Purchase/Sale) handle GL.
- **Rationale:** Keeper recording is parallel record for cross-referencing only. Avoids double-posting.
- **Verified:** Specs §6.1 confirm this design

### AD-005: Module toggle via config table
- **Decision:** `warehouse_keeper_config` table with `module_enabled` flag
- **Rationale:** Matches MISA toggle pattern. Can be hidden without affecting existing warehouse functionality.
- **Implementation:** Handler checks config before serving keeper routes

---

## Dependency Graph

```
Migration (DB schema)
    │
    ├── Domain models (models_warehouse_keeper.go)
    │       │
    │       ├── Repository interface (interfaces.go — add methods)
    │       │       │
    │       │       ├── PG repository (pg_warehouse_keeper.go)
    │       │       │
    │       │       └── Memory repository (memory_warehouse_keeper.go)
    │       │               │
    │       │               └── Service (warehouse_keeper_service.go)
    │       │                       │
    │       │                       └── Handler (warehouse_keeper_handler.go)
    │       │                               │
    │       │                               └── Routes (handler.go — add to RegisterRoutesWithCompany)
    │       │                                       │
    │       │                                       └── Frontend (7 HTML pages)
    │       │
    │       └── Tests (warehouse_keeper_handler_test.go)
    │
    └── Wire in main.go (PG + memory branches)
```

---

## Task List

### Phase 1: Foundation (Week 1) — MUST HAVE

#### Task 1: Migration — Database Schema
**Description:** Create migration `000031_warehouse_keeper.up.sql` with 4 tables: `warehouse_keeper_assignments`, `stock_ledger_entries`, `keeper_inventory_counts`, `warehouse_keeper_config`. Use `CREATE TABLE IF NOT EXISTS` for idempotency.

**Acceptance criteria:**
- [ ] 4 tables created with proper FKs, indexes, defaults
- [ ] `.up.sql` and `.down.sql` both exist
- [ ] `go migrate` runs without error against PG
- [ ] Memory backend unaffected (no migration needed)

**Verification:**
- [ ] `ls migrations/000031*` shows both files
- [ ] `DATABASE_URL=... go run .` starts without migration error

**Dependencies:** None
**Files:** `migrations/000031_warehouse_keeper.up.sql`, `migrations/000031_warehouse_keeper.down.sql`
**Estimated scope:** S (2 files)

---

#### Task 2: Domain Models
**Description:** Create `internal/domain/models_warehouse_keeper.go` with 4 structs: `WarehouseKeeperAssignment`, `StockLedgerEntry`, `KeeperInventoryCount`, `WarehouseKeeperConfig`. Plus supporting types: `LedgerFilter`, `KeeperReconciliationItem`, `StockCard`, `KeeperInventorySummaryItem`. All in `package domain`.

**Acceptance criteria:**
- [ ] All structs match specs §1 exactly
- [ ] JSON tags present on all fields
- [ ] Validate tags present (required, min, etc.)
- [ ] Enums as string constants (e.g., `VoucherTypeReceipt = "receipt"`)
- [ ] `go build ./...` passes

**Verification:**
- [ ] `go vet ./internal/domain/` passes
- [ ] Structs compile without errors

**Dependencies:** None
**Files:** `internal/domain/models_warehouse_keeper.go`
**Estimated scope:** S (1 file, ~150 lines)

---

#### Task 3: Repository Interface
**Description:** Add `WarehouseKeeperRepository` interface to `internal/domain/interfaces.go`. Methods: 6 for Assignment, 5 for Stock Ledger, 1 for Reconciliation, 1 for Stock Card, 1 for Keeper Reports. Total: 14 methods.

**Acceptance criteria:**
- [ ] Interface matches specs §2 exactly
- [ ] All method signatures use `context.Context` as first param
- [ ] Filter struct `LedgerFilter` defined
- [ ] `go build ./...` passes

**Verification:**
- [ ] `go vet ./internal/domain/` passes
- [ ] Interface compiles

**Dependencies:** Task 2 (domain models)
**Files:** `internal/domain/interfaces.go` (append to existing file)
**Estimated scope:** XS (add ~50 lines to existing file)

---

#### Task 4: PG Repository
**Description:** Create `internal/repository/pg_warehouse_keeper.go` implementing `WarehouseKeeperRepository`. Follow existing PG repo pattern: struct with `*gorm.DB`, GORM queries, domain↔GORM model conversion.

**Acceptance criteria:**
- [ ] All 14 interface methods implemented
- [ ] Uses `*gorm.DB` (same as `pg_warehouse.go`)
- [ ] Proper error handling (wrap domain errors)
- [ ] Pagination on list methods
- [ ] `go build ./...` passes

**Verification:**
- [ ] `go vet ./internal/repository/` passes
- [ ] No GORM compilation errors

**Dependencies:** Task 2, Task 3
**Files:** `internal/repository/pg_warehouse_keeper.go`
**Estimated scope:** M (1 file, ~400 lines)

---

#### Task 5: Memory Repository
**Description:** Create `internal/repository/memory_warehouse_keeper.go` implementing `WarehouseKeeperRepository`. Follow existing pattern: `sync.RWMutex` + maps, generate IDs with `fmt.Sprintf`.

**Acceptance criteria:**
- [ ] All 14 interface methods implemented
- [ ] Thread-safe with `sync.RWMutex`
- [ ] ID generation follows convention: `KEEPER-{uuid}`, `LEDGER-{uuid}`
- [ ] Copy-before-mutate pattern for updates
- [ ] `go build ./...` passes

**Verification:**
- [ ] `go vet ./internal/repository/` passes
- [ ] Compiles without errors

**Dependencies:** Task 2, Task 3
**Files:** `internal/repository/memory_warehouse_keeper.go`
**Estimated scope:** M (1 file, ~350 lines)

---

#### Task 6: Service
**Description:** Create `internal/service/warehouse_keeper_service.go` with `WarehouseKeeperService` struct. Business logic for: Assignment CRUD, Ledger recording/un-recording, Pending slips query, Reconciliation report, Stock Card, Keeper reports, Config management.

**Acceptance criteria:**
- [ ] `NewWarehouseKeeperService` takes `WarehouseKeeperRepository` + read-only access to warehouse repos (for StockBalance, InventoryTransaction queries)
- [ ] Assignment CRUD with overlap validation
- [ ] Recording: validates keeper assignment, creates ledger entry, updates balance
- [ ] Un-recording: validates period not closed, creates audit trail
- [ ] Reconciliation: joins ledger entries with stock balances
- [ ] Config: get/update with module_enabled check
- [ ] All tests pass: `go test -count=1 ./internal/service/`

**Verification:**
- [ ] `go vet ./internal/service/` passes
- [ ] Service compiles

**Dependencies:** Task 2, Task 3, Task 4, Task 5
**Files:** `internal/service/warehouse_keeper_service.go`
**Estimated scope:** L (1 file, ~500 lines)

---

### Checkpoint: Foundation
- [ ] All tasks 1-6 compile
- [ ] `go vet ./...` passes
- [ ] Domain models match specs
- [ ] Repository interface complete
- [ ] Both PG + Memory repos compile
- [ ] Service compiles with business logic
- [ ] **REVIEW BEFORE PROCEEDING**

---

### Phase 2: API Layer (Week 2) — MUST HAVE

#### Task 7: Handler — Assignment Endpoints
**Description:** Create `internal/handler/warehouse_keeper_handler.go` with `WarehouseKeeperHandler` struct. Implement assignment CRUD endpoints: POST, GET list, GET detail, PUT, DELETE.

**Acceptance criteria:**
- [ ] `NewWarehouseKeeperHandler` takes `*service.WarehouseKeeperService`
- [ ] 5 assignment endpoints implemented
- [ ] Request binding with validation
- [ ] Proper HTTP status codes (201 created, 400 bad request, 404 not found)
- [ ] Error mapping follows `errors.Is` pattern (see `faError` in `fixed_asset_handler.go:79`)
- [ ] `go build ./...` passes

**Verification:**
- [ ] `go vet ./internal/handler/` passes
- [ ] Handler compiles

**Dependencies:** Task 6
**Files:** `internal/handler/warehouse_keeper_handler.go`
**Estimated scope:** M (1 file, ~200 lines for assignments)

---

#### Task 8: Handler — Ledger & Recording Endpoints
**Description:** Add to `warehouse_keeper_handler.go`: Ledger list, Ledger detail, Record slips (bulk), Un-record, Ledger balance. Plus Pending slips list and count.

**Acceptance criteria:**
- [ ] 7 ledger/pending endpoints implemented
- [ ] Record endpoint accepts bulk slip IDs
- [ ] Un-record requires reason in request body
- [ ] Balance endpoint returns per-item balances
- [ ] `go build ./...` passes

**Verification:**
- [ ] `go vet ./internal/handler/` passes

**Dependencies:** Task 7
**Files:** `internal/handler/warehouse_keeper_handler.go` (append)
**Estimated scope:** M (~250 lines)

---

#### Task 9: Handler — Reports & Config Endpoints
**Description:** Add to `warehouse_keeper_handler.go`: Reconciliation report, Stock Card, Inventory summary, Receipt/Issue detail, Count variance, Config get/update.

**Acceptance criteria:**
- [ ] 6 report/config endpoints implemented
- [ ] Reconciliation returns per-item keeper vs accounting qty
- [ ] Stock Card returns opening/closing balance with line items
- [ ] Config update validates boolean fields
- [ ] `go build ./...` passes

**Verification:**
- [ ] `go vet ./internal/handler/` passes

**Dependencies:** Task 8
**Files:** `internal/handler/warehouse_keeper_handler.go` (append)
**Estimated scope:** M (~200 lines)

---

#### Task 10: Route Registration
**Description:** Add `RegisterKeeperRoutes` function and wire into `handler.go`. Add `WarehouseKeeperHandler` parameter to `RegisterRoutesWithCompany`.

**Acceptance criteria:**
- [ ] Routes registered under `/api/v1/warehouse/keeper/`
- [ ] Auth middleware applied (same as warehouse routes)
- [ ] `RegisterRoutesWithCompany` signature updated (new param)
- [ ] `main.go` updated to create and pass `WarehouseKeeperHandler`
- [ ] `go build ./...` passes
- [ ] Server starts without error

**Verification:**
- [ ] `JWT_SECRET=devsecret go run .` starts
- [ ] `curl localhost:8080/api/v1/warehouse/keeper/assignments` returns 401 (auth required)

**Dependencies:** Task 7, Task 8, Task 9
**Files:** `internal/handler/warehouse_keeper_handler.go` (route registration), `internal/handler/handler.go` (update signature), `main.go` (wire)
**Estimated scope:** M (3 files)

---

#### Task 11: Handler Tests — Assignment
**Description:** Create `internal/handler/warehouse_keeper_handler_test.go` with tests for assignment CRUD. Follow existing pattern: in-memory repos, real service, mock auth middleware.

**Acceptance criteria:**
- [ ] TestCreateAssignment, TestListAssignments, TestGetAssignment, TestUpdateAssignment, TestDeleteAssignment
- [ ] Uses in-memory repo + real service
- [ ] Mock auth middleware (sets user_id, role)
- [ ] `go test -count=1 ./internal/handler/ -run TestKeeper` passes

**Verification:**
- [ ] All keeper tests pass
- [ ] No existing tests broken

**Dependencies:** Task 10
**Files:** `internal/handler/warehouse_keeper_handler_test.go`
**Estimated scope:** M (~200 lines)

---

#### Task 12: Handler Tests — Ledger & Recording
**Description:** Add tests for ledger endpoints: TestRecordSlips, TestUnrecordSlip, TestListLedgerEntries, TestGetPendingSlips.

**Acceptance criteria:**
- [ ] Recording creates ledger entry with correct fields
- [ ] Un-recording requires reason
- [ ] Cannot un-record in closed period
- [ ] Pending slips shows unrecorded items only
- [ ] `go test -count=1 ./internal/handler/ -run TestKeeper` passes

**Verification:**
- [ ] All keeper tests pass

**Dependencies:** Task 11
**Files:** `internal/handler/warehouse_keeper_handler_test.go` (append)
**Estimated scope:** M (~200 lines)

---

### Checkpoint: API Layer
- [ ] All endpoints implemented (18 routes)
- [ ] Route registration works
- [ ] Server starts, auth required on all keeper routes
- [ ] Assignment CRUD tests pass
- [ ] Ledger recording tests pass
- [ ] `go test -count=1 ./...` passes (no regressions)
- [ ] **REVIEW BEFORE PROCEEDING**

---

### Phase 3: Reports & Polish (Week 3) — SHOULD HAVE

#### Task 13: Reconciliation Report
**Description:** Implement reconciliation logic in service + handler. Joins StockLedgerEntry (keeper records) with StockBalance (accounting records) on item_id + warehouse_id.

**Acceptance criteria:**
- [ ] Returns per-item: keeper qty, accounting qty, variance qty, variance value
- [ ] Filters by warehouse and date range
- [ ] Highlights variances > $1M VND
- [ ] Export to Excel endpoint

**Verification:**
- [ ] `go test -count=1 ./internal/handler/ -run TestKeeperReconciliation` passes

**Dependencies:** Task 12
**Files:** `internal/service/warehouse_keeper_service.go`, `internal/handler/warehouse_keeper_handler.go`
**Estimated scope:** S (~100 lines)

---

#### Task 14: Stock Card Report
**Description:** Implement stock card generation. Queries StockLedgerEntry for a specific item+warehouse+period, returns ordered list with running balance.

**Acceptance criteria:**
- [ ] Returns opening balance, line items, closing balance
- [ ] Supports period filter (YYYYMM format)
- [ ] Format matches Template 4 from specs

**Verification:**
- [ ] `go test -count=1 ./internal/handler/ -run TestKeeperStockCard` passes

**Dependencies:** Task 12
**Files:** `internal/service/warehouse_keeper_service.go`, `internal/handler/warehouse_keeper_handler.go`
**Estimated scope:** S (~80 lines)

---

#### Task 15: Keeper Reports
**Description:** Implement inventory summary, receipt/issue detail, count variance reports.

**Acceptance criteria:**
- [ ] Inventory summary: per-item qty + value per warehouse
- [ ] Receipt/Issue detail: grouped by voucher type
- [ ] Count variance: from StockTake items with variance

**Verification:**
- [ ] `go test -count=1 ./internal/handler/ -run TestKeeperReport` passes

**Dependencies:** Task 12
**Files:** `internal/service/warehouse_keeper_service.go`, `internal/handler/warehouse_keeper_handler.go`
**Estimated scope:** S (~100 lines)

---

#### Task 16: Module Toggle
**Description:** Implement config get/update. Handler checks `module_enabled` before serving keeper endpoints.

**Acceptance criteria:**
- [ ] GET /config returns current config
- [ ] PUT /config updates config
- [ ] When module disabled, keeper routes return 404 or hidden
- [ ] Config changes logged

**Verification:**
- [ ] `go test -count=1 ./internal/handler/ -run TestKeeperConfig` passes

**Dependencies:** Task 12
**Files:** `internal/service/warehouse_keeper_service.go`, `internal/handler/warehouse_keeper_handler.go`
**Estimated scope:** XS (~50 lines)

---

### Checkpoint: Reports
- [ ] Reconciliation report works
- [ ] Stock Card works
- [ ] Keeper reports work
- [ ] Module toggle works
- [ ] `go test -count=1 ./...` passes
- [ ] **REVIEW BEFORE PROCEEDING**

---

### Phase 4: Frontend (Week 4) — MUST HAVE

#### Task 17: Keeper Assignment Page
**Description:** Create `web/app/warehouse-keeper-assignments.html` with Alpine.js component. CRUD for keeper assignments. Table view + create/edit modal.

**Acceptance criteria:**
- [ ] Lists all assignments for company
- [ ] Create form: warehouse dropdown, user dropdown, role, effective dates
- [ ] Edit/delete actions
- [ ] Uses `apiGet`/`apiPost`/`apiPut`/`apiDelete` from app.js
- [ ] Calls `mountAppShell(title, activePath)` on init

**Verification:**
- [ ] Page loads in browser
- [ ] Can create assignment
- [ ] Can list assignments

**Dependencies:** Task 10
**Files:** `web/app/warehouse-keeper-assignments.html`
**Estimated scope:** M (~300 lines HTML/JS)

---

#### Task 18: Pending Slips Page
**Description:** Create `web/app/warehouse-keeper-pending.html`. Shows pending slips with checkboxes, bulk "Ghi sổ" button.

**Acceptance criteria:**
- [ ] Lists unrecorded slips for assigned warehouse
- [ ] Checkbox selection (individual + select all)
- [ ] "Ghi sổ" button records selected slips
- [ ] Slip count badge
- [ ] Expandable rows showing item details

**Verification:**
- [ ] Page loads
- [ ] Can select and record slips
- [ ] Slips move to recorded after recording

**Dependencies:** Task 10
**Files:** `web/app/warehouse-keeper-pending.html`
**Estimated scope:** M (~350 lines)

---

#### Task 19: Stock Ledger Page
**Description:** Create `web/app/warehouse-keeper-ledger.html`. Shows stock ledger entries with filters.

**Acceptance criteria:**
- [ ] Filter: warehouse, item, date range, status
- [ ] Table: date, voucher type, voucher no, description, receipt qty, issue qty, balance
- [ ] Optional cost columns (based on config)
- [ ] Print button (generates PDF-ready view)
- [ ] Pagination

**Verification:**
- [ ] Page loads
- [ ] Filters work
- [ ] Data displays correctly

**Dependencies:** Task 10
**Files:** `web/app/warehouse-keeper-ledger.html`
**Estimated scope:** M (~300 lines)

---

#### Task 20: Stock Card & Reconciliation Pages
**Description:** Create `web/app/warehouse-keeper-stock-card.html` and `web/app/warehouse-keeper-reconciliation.html`.

**Acceptance criteria:**
- [ ] Stock Card: warehouse + item selector, period picker, card display
- [ ] Reconciliation: warehouse + date range, table with variance highlighting
- [ ] Export buttons

**Verification:**
- [ ] Both pages load
- [ ] Stock Card shows correct format
- [ ] Reconciliation shows variances

**Dependencies:** Task 10
**Files:** `web/app/warehouse-keeper-stock-card.html`, `web/app/warehouse-keeper-reconciliation.html`
**Estimated scope:** M (~400 lines total)

---

#### Task 21: Count & Config Pages
**Description:** Create `web/app/warehouse-keeper-count.html` and `web/app/warehouse-keeper-config.html`.

**Acceptance criteria:**
- [ ] Count page: shows count sheet, enter physical qty, submit
- [ ] Config page: toggle module on/off, toggle cost visibility

**Verification:**
- [ ] Both pages load
- [ ] Config toggle works

**Dependencies:** Task 10
**Files:** `web/app/warehouse-keeper-count.html`, `web/app/warehouse-keeper-config.html`
**Estimated scope:** M (~350 lines total)

---

### Checkpoint: Complete
- [ ] All 7 frontend pages created
- [ ] All backend tests pass
- [ ] `go build ./...` passes
- [ ] `go test -count=1 ./...` passes (no regressions)
- [ ] All success criteria from BRD §10 met
- [ ] **FINAL REVIEW**

---

## Risks and Mitigations

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| `RegisterRoutesWithCompany` signature change breaks other handlers | Medium | High | Add new param at end with zero-value default, or create `RegisterRoutesWithCompanyKeeper` variant |
| Memory repo complexity for ledger balance tracking | Low | Medium | Keep simple: iterate entries, sum receipts/issues |
| Frontend Alpine.js state management complexity | Medium | Low | Follow existing page patterns exactly |
| Circular 99 interpretation changes mid-implementation | Low | High | Monitor MOF, keep business rules in service layer |

## Open Questions

1. **Route registration:** Should we modify `RegisterRoutesWithCompany` signature (breaking change) or create separate `RegisterKeeperRoutes` called independently? **Recommendation:** Separate `RegisterKeeperRoutes` called from `main.go` — zero impact on existing handlers.

2. **Cost visibility:** Default hidden or visible? **Recommendation:** Hidden by default (safer), toggleable per company.

3. **Auto-recording:** Should GRN/DN auto-create ledger entries? **Recommendation:** No — keep manual workflow for separation of duties. Can add as optional config later.

## Summary

| Phase | Tasks | Duration | Deliverable |
|-------|-------|----------|-------------|
| Foundation | 1-6 | Week 1 | Domain, repos, service, migration |
| API Layer | 7-12 | Week 2 | Handler, routes, tests |
| Reports | 13-16 | Week 3 | Reconciliation, stock card, reports, toggle |
| Frontend | 17-21 | Week 4 | 7 HTML pages |
| **Total** | **21 tasks** | **4 weeks** | **Full MISA parity** |
