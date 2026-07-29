# Implementation Plan: Inventory (Warehouse) Module

## Overview

Complete inventory module from 40% foundation to 100% production-ready. Current state: 9 entity CRUD built, 0% GL posting pipeline, 0% stock balance mutation on POST, stub valuation engine.

Goal per Circular 99/2025/TT-BTC (direct posting to 151-158, no TK 611). 30 movement types must post double-entry GL.

## Architecture Decisions

- **GL posting as service-layer concern** — `WarehouseService` calls `GLService` (existing) after state transitions. No direct DB.
- **StockBalance mutation explicit** — POST transitions in adjustment/take/transfer service methods update balances. No hidden side effects.
- **GRN/DN as new domain entities** — follow existing pattern: model → interfaces → PG+memo repos → service → handler → routes → test.
- **Existing state machines extend** — add PICKED/IN_TRANSIT/RECEIVED to TransferStatus; add ROLLED_BACK to ValuationRunStatus.
- **Cost engine separate service** — `ValuationService` (or method on `WarehouseService`) with weighted_avg/FIFO/specific_id computation. No external lib.
- **No new external dependencies** — pure Go math for cost calc.
- **Vertical slices per entity** — each task delivers one complete workflow (model→DB→service→handler→test→GL).

## Dependency Graph

```
Phase 0: GL Pipeline
  ├── StockBalance mutation helpers
  └── GL posting interface
  
Phase 1: Core Docs (GRN + DN)
  ├── GRN entity + lifecycle + GL posting
  └── DN entity + lifecycle + GL posting
  
Phase 2: Special Movements
  ├── Direct Receipt (non-PO)
  ├── Sales Return Receipt
  ├── Supplier Return
  └── Internal Issue
  
Phase 3: Cost & Period-End
  ├── Valuation Engine (real calc)
  ├── Cost Revaluation entity
  ├── FIFO layer tracking
  └── NRV Provision (2294)
  
Phase 4: Integration
  ├── OB→StockBalance wiring
  ├── Purchase module→GRN wiring
  ├── Sales module→DN wiring
  └── Period lock enforcement
```

## In-Scope vs Out-of-Scope

| In-Scope | Out-of-Scope |
|----------|-------------|
| GL posting for all 30 movement types | Barcode/RFID scanner hardware integration |
| GRN + DN domain entities + lifecycle | Physical warehouse layout (bin/rack management) |
| Real valuation engine (WAC, FIFO, Specific ID) | Serial number tracking engine |
| StockBalance mutation on POST/POSTED | Lot expiry date alerts |
| Stock Transfer PICKED/IN_TRANSIT states | Demand forecasting / reorder suggestions |
| NRV provision 2294 | ABC analysis engine |
| Cost Revaluation entity | Manufacturing BOM/routing |
| Opening Balance → StockBalance wiring | Customs declaration tracking |
| Commit stocks for pending transfers | E-invoice integration (purchase) |

## Phases & Tasks

### Phase 0: Foundation — GL Pipeline + Balance Mutation

Base for everything. Without this, no inventory movement changes financials.

- [ ] **Task 0.1: StockBalance mutation helpers** — `AdjustStockBalance(companyID, warehouseID, itemID, period, deltaQty, deltaValue)`. Called by POST transitions. Upserts balance record. Creates `InventoryTransaction` audit log.
- [ ] **Task 0.2: GL posting integration** — Define `InventoryGLMapper` (reason+type→debit/credit account). Wire `GLService` into `WarehouseService`. Post journal entries on every POST transition.
- [ ] **Task 0.3: Reason-based GL account map** — Realize the 30-entry matrix (doc §18) as code constants. Map AdjType+Reason, TransferStatus, TakeItem variance, etc. to Debit/Credit pairs.
- [ ] **Checkpoint: Phase 0** — `StockBalance` updates on adjustment/take POST. GL entries created and retrievable. All existing 22 tests still pass.

### Phase 1: Core Docs — GRN (Goods Receipt Note)

Purchase receipt workflow. Links Purchase Order → Warehouse.

- [ ] **Task 1.1: GRN domain model** — `GRN` struct with header (PO ref, warehouse, receipt_date, status) + `GRNItem` lines (item, qty_received, qty_rejected, unit_price, lot#, serial#). Status: DRAFT→POSTED→CANCELLED. Add to `models_warehouse.go`.
- [ ] **Task 1.2: GRN repository interfaces + errors** — `GRNRepository` (Create, GetByID, List, Update, UpdateStatus, GetItems) in `interfaces.go`. Errors in `errors.go`.
- [ ] **Task 1.3: GRN PG repo** — `pg_warehouse.go`: `CREATE TABLE IF NOT EXISTS grns` + `grn_items`. Add to `007_warehouse_schema.sql`.
- [ ] **Task 1.4: GRN memory repo** — `memory_warehouse.go`: map-based, concurrent-safe.
- [ ] **Task 1.5: GRN service** — `warehouse_service.go`: Create, Submit, Post (→ StockBalance + GL: Dr 152/156 Cr 331), Cancel (reverse).
- [ ] **Task 1.6: GRN handler + routes** — `warehouse_handler.go`: POST/GET/list/GET/:id/PATCH post/cancel.
- [ ] **Task 1.7: GRN tests** — `warehouse_handler_test.go`: create from PO, post, cancel, over-receipt reject, duplicate supplier delivery warn.
- [ ] **Checkpoint: Phase 1** — Full GRN lifecycle tested. Stock increases on POST. GL Dr 152/156 Cr 331 created. PO integration (received_qty tracking) deferred to Phase 4.

### Phase 2: Core Docs — DN (Delivery Note)

Sales delivery workflow. Links Sales Order → Warehouse.

- [ ] **Task 2.1: DN domain model** — `DN` struct + `DNLine`. Status: DRAFT→[PICKED→PACKED→]POSTED→CANCELLED. Add to `models_warehouse.go`.
- [ ] **Task 2.2: DN repository interfaces + errors** — `DNRepository` in `interfaces.go`.
- [ ] **Task 2.3: DN PG repo** — `pg_warehouse.go` + schema migration.
- [ ] **Task 2.4: DN memory repo** — `memory_warehouse.go`.
- [ ] **Task 2.5: DN service** — Create, Submit, Pick, Pack (opt), Post (→ StockBalance -qty, GL: Dr 632 Cr 152/156), Cancel (reverse).
- [ ] **Task 2.6: DN handler + routes** — POST/GET/list/GET/:id/PATCH transitions.
- [ ] **Task 2.7: DN tests** — Full lifecycle, stock check, over-delivery tolerance, serial scan.
- [ ] **Checkpoint: Phase 2** — Full DN lifecycle tested. Stock decreases on POST. COGS posted. Short-pick → partial delivery + backorder.

### Phase 3: Existing Entity Enhancements

Upgrade stock transfer/adjustment/take to spec-level state machines with stock + GL effects.

- [ ] **Task 3.1: Stock Transfer — add PICKED/IN_TRANSIT/RECEIVED states** — Extend `TransferStatus` enum. Add transit account (151) GL posting logic. `committed_qty` field on `StockBalance`.
- [ ] **Task 3.2: Stock Transfer — stock mutation** — On TRANSFERRED: source stock -= qty, transit += qty. On COMPLETED: dest stock += qty, transit -= qty. GL Dr 151 Cr 152/156 → Dr 152/156 Cr 151.
- [ ] **Task 3.3: Stock Transfer — tests** — Committed qty, transit accounting, short receipt with carrier claim.
- [ ] **Task 3.4: Stock Adjustment — stock mutation + GL** — On POST: adjust StockBalance ±qty. GL per AdjType+Reason: increase→Dr 152 Cr 3381/632/711, decrease→Dr 632/642/811 Cr 152.
- [ ] **Task 3.5: Stock Adjustment — threshold approval** — Configurable approval levels (<5M auto, 5-100M Mgr, >100M CA). Escalate if same item adjusted 3+/month.
- [ ] **Task 3.6: Stock Take — recount workflow** — Variance > tolerance→trigger recount. Max 3 counts. Blind count option. Auto-accept within tolerance.
- [ ] **Task 3.7: Stock Take — POST→adjustment** — On POST: auto-generate `StockAdjustment` for variance lines. Gain→Dr 152 Cr 632/711. Loss→Dr 632/811 Cr 152.
- [ ] **Task 3.8: Stock Take — ABC scheduling** — A=monthly, B=quarterly, C=yearly. Prevent overlapping counts for same warehouse.
- [ ] **Checkpoint: Phase 3** — All 9 existing entities now mutate stock + post GL. State machines match spec (or justified divergence). 30+ tests pass.

### Phase 4: Special Movements

- [ ] **Task 4.1: Direct Receipt entity** — Non-PO receipt (donation/sample/owner_contribution/surplus). Model + repo + service. Reason→GL map: donation→Cr 711, owner→Cr 4111, surplus→Cr 3381. Handler + routes.
- [ ] **Task 4.2: Direct Receipt — threshold approval** — <10M auto, 10-100M WH Mgr, >100M CA.
- [ ] **Task 4.3: Sales Return Receipt** — Return from customer. Ties to Credit Note. GL Dr 152/156 Cr 632 (reverse COGS at original cost). Handler + routes.
- [ ] **Task 4.4: Return to Supplier** — Outbound return to vendor. Pre-invoice: reverse GRN GL. Post-invoice: Dr 331 Cr 152 + 1331 (via credit note).
- [ ] **Task 4.5: Internal Issue** — Samples/marketing/charity/admin. Department + cost_center tracking. GL per purpose: samples→Dr 6417, charity→Dr 811, admin→Dr 642.
- [ ] **Task 4.6: Special movement tests** — Full lifecycle + GL assertion for each type.
- [ ] **Checkpoint: Phase 4** — All 30 movement types from §18 matrix covered by code. 50+ tests pass.

### Phase 5: Cost & Valuation Engine

- [ ] **Task 5.1: Valuation engine — Weighted Average (periodic)** — Real calculation: `(open_value + receipt_value) / (open_qty + receipt_qty)`. Apply to all issues in period. Post adjustment Dr/Cr 632.
- [ ] **Task 5.2: Valuation engine — FIFO** — Sort receipt layers by date. Issue from earliest layer. Track remaining layers for closing stock. Post adjustment.
- [ ] **Task 5.3: Valuation engine — Specific ID** — Match issue to receipt via item+warehouse+date (serial# integration deferred). Post adjustment.
- [ ] **Task 5.4: Cost Revaluation entity** — Ad-hoc unit cost change (Chief Accountant only). GL: Dr 152 Cr 611 (increase) or reverse. Model + repo + service + handler.
- [ ] **Task 5.5: Valuation engine — Rollback** — `ROLLED_BACK` state. Reverse all cost adjustment entries. Reopen period for correction.
- [ ] **Task 5.6: Valuation engine tests** — WAC calculation verified. FIFO layers correct. Specific ID matches.
- [ ] **Checkpoint: Phase 5** — Valuation engine computes real costs. All 4 methods tested. Cost adjustments created.

### Phase 6: Period-End

- [ ] **Task 6.1: NRV provision 2294** — Calculate NRV per item: selling price - completion cost - selling cost. Compare cost. Post: Dr 632 Cr 2294 (increase) or reverse. Deferred if no selling price data.
- [ ] **Task 6.2: Opening Balance → StockBalance** — Wire OB module POST hook to create `StockBalance` records. GL Dr 152 Cr 3388/419.
- [ ] **Task 6.3: Period lock enforcement** — Before any POST: check period OPEN. After valuation complete: lock new movements. Enforce throughout service layer.
- [ ] **Task 6.4: Period-end close workflow** — Validate all docs POSTED. Run valuation. Post provision. Lock period. Integration test.
- [ ] **Checkpoint: Phase 6** — Period-end close tested end-to-end. Provision 2294 correct. OB integration works.

### Phase 7: Integration Wiring

- [ ] **Task 7.1: Purchase → GRN** — Purchase order POSTED hook auto-creates GRN. PO line `received_qty` tracking.
- [ ] **Task 7.2: Sales → DN** — Sales order CONFIRMED hook enables DN creation. SO line `delivered_qty` tracking.
- [ ] **Task 7.3: Credit Note → Sales Return** — Credit note POSTED triggers return receipt. Reverse COGS.
- [ ] **Task 7.4: Debit Note → Supplier Return** — Debit note POSTED triggers supplier return.
- [ ] **Task 7.5: Service layer integration tests** — Full PO→GRN→Stock→DN→COGS→Invoice cycle with GL assertion.
- [ ] **Checkpoint: Phase 7** — Purchase-to-pay and order-to-cash cycles integrated with warehouse. Both backends (PG + memory) tested.

## Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| GL pipeline changes affect existing modules | High | Define `InventoryGLMapper` as new interface; existing GL unchanged |
| FIFO layer tracking perf at scale | Medium | Layer store as ordered slice in `StockBalance`; cap layers (max 12mo) |
| GRN/DN duplicate with purchase/sale module existing fields | Medium | New entities own warehouse lifecycle; purchase/sale ref fields stay as simple TEXT |
| Period lock conflicts with concurrent ops | Low | Advisory lock per period; service-level mutex for memory backend |
| Valuation engine deadline slip | Low | Weighted_avg first (simplest); FIFO + Specific ID deferred to Phase 5.2-5.3 |

## Open Questions

- GL posting: reuse existing `GLService` interface or new `InventoryGLService`? Recommend reuse — `GLService.CreateJournalEntry(ctx, entries)` already exists.
- Stock balance granularity: per-period (current) or running balance (perpetual)? Current per-period is fine; perpetual adds complexity without regulatory requirement.
- Serial# tracking: build now or defer? Defer — no regulatory mandate in TT 99 for serial tracking. Items with `is_serialized=true` get validation on GRN/DN but full engine in v2.1.
