# Implementation Plan: Purchase Module (P2P)
## Overview
Build complete Procure-to-Pay module for GoTax: suppliers → purchase orders → goods receipts → supplier invoices → AP + reports. Domain layer exists (models, interfaces, errors). Everything below — repos, services, handlers, routes, migrations — needs building. Mirror `bank/` module pattern.

## Architecture Decisions
- **Separate `PurchaseService` + `PurchaseHandler`** (not merge into main `service.Service`) — matches existing `BankService`/`BankHandler` pattern
- **Slug IDs**: `SUP-` (supplier), `PO-` (PO, but PO number is business key), `GRN-`, `INV-`
- **PO/GRN numbering**: auto-generate `PO-YYYYMM-XXXX` / `GRN-YYYYMM-XXXX` in service
- **No logger dependency**: stdlib `log.Printf` for critical events only
- **No third-party integrations**: skip GDT e-invoice webhook, skip GL auto-posting (P1)
- **Data model**: 5 tables (supplier, purchase_order + po_line, goods_receipt_note + grn_line, supplier_invoice + invoice_line, ap_transaction)

## Phases

### Phase 0: Infrastructure
- Migration schema + pg.go + memory repo stubs

### Phase 1: Suppliers (P0)
- Service CRUD, Handler CRUD, Tests

### Phase 2: Purchase Orders (P0)
- Service: CRUD + approve/cancel/close state machine, Handler, Tests

### Phase 3: GRN (P0)
- Service: CRUD + post/cancel, Handler, Tests

### Phase 4: Supplier Invoices (P0)
- Service: CRUD + verify/post/cancel/claimVAT, 3-way match stub, Handler, Tests

### Phase 5: AP Transactions + Reports (P0)
- Service: AP transaction CRUD + aging reports, Handler (3 endpoints), Tests

### Phase 6: Main.go Wiring
- Register routes, instantiate services, register mocks in tests

### Phase 7: Review
- MISA/Fast/Bravo gap check, lint + vet + test green

## File Plan

| # | File (new) | Purpose |
|---|-----------|---------|
| 1 | `migrations/005_purchase_schema.sql` | 5 tables + indexes |
| 2 | `internal/repository/pg.go` (append) | `RunMigrations` path |
| 3 | `internal/repository/memory_purchase.go` | All 6 repo interfaces |
| 4 | `internal/repository/pg_purchase.go` | All 6 repo PG impls |
| 5 | `internal/service/purchase_service.go` | Business logic |
| 6 | `internal/handler/purchase_handler.go` | HTTP endpoints |
| 7 | `internal/handler/purchase_handler_test.go` | Handler tests |
| 8 | `main.go` (edit) | Route registration + DI |

| # | File (edit) | Change |
|---|------------|--------|
| 1 | `migrations/` | Add `005_purchase_schema.sql` |
| 2 | `internal/repository/pg.go` | Add migration path |
| 3 | `internal/handler/handler_test.go` | Add purchaseSetupTest + register routes |
| 4 | `main.go` | Wire purchase service + handler |

## Dependencies
```
005_purchase_schema.sql
    ↓ pg_purchase.go (needs SQL schema)
    ↓ memory_purchase.go (independent)
purchase_service.go (needs both repos)
    ↓ purchase_handler.go (needs service)
    ↓ handler_test.go (needs handler + repos + service)
main.go (needs everything)
```

## Risk Mitigation
- 3-way match: stub the tolerances (5% qty, 5% price, 10K VND) as configurable constants, skip actual matching logic for v1
- AP aging: simple bucket calc (0, 1-30, 31-60, 61-90, 91-120, 120+)
- GL posting: stub — no auto-journal for now (depends on Decree 99 mapping being finalized)
- E-invoice receipt: stub — store raw XML only, no GDT push webhook
- Purchase requisition: skip (P1)
- Payment: not in purchase scope — handled by bank module

## Implementation Order
Execute in order below, one subagent task per phase.
