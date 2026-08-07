# Task List: Thủ Kho (Warehouse Keeper) Module

Generated from `tasks/plan.md`. Check off as tasks complete.

## Phase 1: Foundation (Week 1)

- [ ] **Task 1:** Migration — Database Schema (4 tables)
  - Acceptance: 4 tables with FKs, indexes; `.up.sql` + `.down.sql`; migration runs
  - Verify: `ls migrations/000031*`; `DATABASE_URL=... go run .`
  - Files: `migrations/000031_warehouse_keeper.{up,down}.sql`
  - Size: S

- [ ] **Task 2:** Domain Models
  - Acceptance: 4 structs + filter/types match specs §1; JSON + validate tags; `go build` passes
  - Verify: `go vet ./internal/domain/`
  - Files: `internal/domain/models_warehouse_keeper.go`
  - Size: S

- [ ] **Task 3:** Repository Interface
  - Acceptance: 14-method interface in `interfaces.go`; matches specs §2
  - Verify: `go vet ./internal/domain/`
  - Files: `internal/domain/interfaces.go` (append)
  - Size: XS

- [ ] **Task 4:** PG Repository
  - Acceptance: All 14 methods; GORM; error handling; pagination
  - Verify: `go vet ./internal/repository/`
  - Files: `internal/repository/pg_warehouse_keeper.go`
  - Size: M

- [ ] **Task 5:** Memory Repository
  - Acceptance: All 14 methods; RWMutex; ID generation; copy-before-mutate
  - Verify: `go vet ./internal/repository/`
  - Files: `internal/repository/memory_warehouse_keeper.go`
  - Size: M

- [ ] **Task 6:** Service
  - Acceptance: Assignment CRUD + overlap validation; recording/un-recording; reconciliation; config
  - Verify: `go vet ./internal/service/`
  - Files: `internal/service/warehouse_keeper_service.go`
  - Size: L

### Checkpoint: Foundation
- [ ] `go vet ./...` passes
- [ ] All domain models compile
- [ ] Both repos compile
- [ ] Service compiles
- [ ] **GATE: Review before Phase 2**

---

## Phase 2: API Layer (Week 2)

- [ ] **Task 7:** Handler — Assignment Endpoints
  - Acceptance: 5 CRUD endpoints; proper HTTP codes; error mapping
  - Verify: `go vet ./internal/handler/`
  - Files: `internal/handler/warehouse_keeper_handler.go`
  - Size: M

- [ ] **Task 8:** Handler — Ledger & Recording Endpoints
  - Acceptance: 7 endpoints; bulk record; un-record with reason
  - Verify: `go vet ./internal/handler/`
  - Files: `internal/handler/warehouse_keeper_handler.go` (append)
  - Size: M

- [ ] **Task 9:** Handler — Reports & Config Endpoints
  - Acceptance: 6 endpoints; reconciliation; stock card; config
  - Verify: `go vet ./internal/handler/`
  - Files: `internal/handler/warehouse_keeper_handler.go` (append)
  - Size: M

- [ ] **Task 10:** Route Registration + main.go Wiring
  - Acceptance: Routes under `/api/v1/warehouse/keeper/`; auth MW; main.go wired
  - Verify: `JWT_SECRET=devsecret go run .` starts; curl returns 401
  - Files: `warehouse_keeper_handler.go`, `handler.go`, `main.go`
  - Size: M

- [ ] **Task 11:** Handler Tests — Assignment
  - Acceptance: 5 tests; in-memory repos; mock auth MW
  - Verify: `go test -count=1 ./internal/handler/ -run TestKeeper`
  - Files: `internal/handler/warehouse_keeper_handler_test.go`
  - Size: M

- [ ] **Task 12:** Handler Tests — Ledger & Recording
  - Acceptance: 4 tests; record/un-record/pending
  - Verify: `go test -count=1 ./internal/handler/ -run TestKeeper`
  - Files: `internal/handler/warehouse_keeper_handler_test.go` (append)
  - Size: M

### Checkpoint: API Layer
- [ ] 18 routes registered
- [ ] Auth required on all keeper routes
- [ ] Assignment tests pass
- [ ] Ledger tests pass
- [ ] `go test -count=1 ./...` passes (no regressions)
- [ ] **GATE: Review before Phase 3**

---

## Phase 3: Reports & Polish (Week 3)

- [ ] **Task 13:** Reconciliation Report
  - Acceptance: Per-item keeper vs accounting qty; variance highlighting; Excel export
  - Verify: `go test -count=1 ./internal/handler/ -run TestKeeperReconciliation`
  - Files: `warehouse_keeper_service.go`, `warehouse_keeper_handler.go`
  - Size: S

- [ ] **Task 14:** Stock Card Report
  - Acceptance: Opening/closing balance; line items; period filter
  - Verify: `go test -count=1 ./internal/handler/ -run TestKeeperStockCard`
  - Files: `warehouse_keeper_service.go`, `warehouse_keeper_handler.go`
  - Size: S

- [ ] **Task 15:** Keeper Reports (Summary, Detail, Variance)
  - Acceptance: 3 report types; proper formatting
  - Verify: `go test -count=1 ./internal/handler/ -run TestKeeperReport`
  - Files: `warehouse_keeper_service.go`, `warehouse_keeper_handler.go`
  - Size: S

- [ ] **Task 16:** Module Toggle
  - Acceptance: Config get/update; module_enabled gates keeper routes
  - Verify: `go test -count=1 ./internal/handler/ -run TestKeeperConfig`
  - Files: `warehouse_keeper_service.go`, `warehouse_keeper_handler.go`
  - Size: XS

### Checkpoint: Reports
- [ ] Reconciliation works
- [ ] Stock Card works
- [ ] Reports work
- [ ] Toggle works
- [ ] `go test -count=1 ./...` passes
- [ ] **GATE: Review before Phase 4**

---

## Phase 4: Frontend (Week 4)

- [ ] **Task 17:** Keeper Assignment Page
  - Acceptance: CRUD UI; Alpine.js; uses apiGet/apiPost
  - Verify: Page loads in browser
  - Files: `web/app/warehouse-keeper-assignments.html`
  - Size: M

- [ ] **Task 18:** Pending Slips Page
  - Acceptance: Checkbox selection; bulk Ghi sổ; slip count
  - Verify: Can record slips from UI
  - Files: `web/app/warehouse-keeper-pending.html`
  - Size: M

- [ ] **Task 19:** Stock Ledger Page
  - Acceptance: Filters; table; print button; pagination
  - Verify: Ledger displays correctly
  - Files: `web/app/warehouse-keeper-ledger.html`
  - Size: M

- [ ] **Task 20:** Stock Card + Reconciliation Pages
  - Acceptance: Both pages functional; export buttons
  - Verify: Both pages load and display data
  - Files: `web/app/warehouse-keeper-stock-card.html`, `web/app/warehouse-keeper-reconciliation.html`
  - Size: M

- [ ] **Task 21:** Count + Config Pages
  - Acceptance: Count sheet; config toggles
  - Verify: Both pages load
  - Files: `web/app/warehouse-keeper-count.html`, `web/app/warehouse-keeper-config.html`
  - Size: M

### Checkpoint: Complete
- [ ] All 7 frontend pages created
- [ ] All backend tests pass
- [ ] `go build ./...` passes
- [ ] `go test -count=1 ./...` passes (no regressions)
- [ ] All BRD §10 success criteria met
- [ ] **FINAL REVIEW**

---

## Progress Summary

| Phase | Tasks | Done | Status |
|-------|-------|------|--------|
| Foundation | 6 | 0/6 | NOT STARTED |
| API Layer | 6 | 0/6 | NOT STARTED |
| Reports | 4 | 0/4 | NOT STARTED |
| Frontend | 5 | 0/5 | NOT STARTED |
| **Total** | **21** | **0/21** | **NOT STARTED** |
