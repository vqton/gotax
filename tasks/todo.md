# Warehouse Module — Implementation Todo

## Phase 0: Foundation — GL Pipeline + Balance Mutation

- [ ] **0.1** `AdjustStockBalance` helper — upsert balance + create `InventoryTransaction` audit log
- [ ] **0.2** `InventoryGLMapper` — reason+type→debit/credit account. Wire `GLService`.
- [ ] **0.3** 30-entry GL posting matrix as code constants (doc §18)
- [ ] **Checkpoint** — StockBalance updates on POST. GL entries created. 22 existing tests pass.

## Phase 1: Core Docs — GRN (Goods Receipt Note)

- [ ] **1.1** GRN domain model (`models_warehouse.go`)
- [ ] **1.2** GRN repo interface + errors
- [ ] **1.3** GRN PG repo + migration (`grns` + `grn_items` tables)
- [ ] **1.4** GRN memory repo
- [ ] **1.5** GRN service (create, post→stock+GL, cancel)
- [ ] **1.6** GRN handler + routes
- [ ] **1.7** GRN tests (lifecycle, over-receipt, cancel)
- [ ] **Checkpoint** — GRN posted→stock+GL. 26+ tests.

## Phase 2: Core Docs — DN (Delivery Note)

- [ ] **2.1** DN domain model (`models_warehouse.go`)
- [ ] **2.2** DN repo interface + errors
- [ ] **2.3** DN PG repo + migration (`dns` + `dn_lines` tables)
- [ ] **2.4** DN memory repo
- [ ] **2.5** DN service (create, pick, pack, post→stock+COGS, cancel)
- [ ] **2.6** DN handler + routes
- [ ] **2.7** DN tests (lifecycle, stock check, short-pick)
- [ ] **Checkpoint** — DN posted→stock- + COGS. 30+ tests.

## Phase 3: Existing Entity Enhancements

- [ ] **3.1** Stock Transfer — PICKED/IN_TRANSIT/RECEIVED states + transit account 151
- [ ] **3.2** Stock Transfer — stock mutation on TRANSFERRED/COMPLETED
- [ ] **3.3** Stock Transfer — tests (committed qty, transit accounting)
- [ ] **3.4** Stock Adjustment — stock mutation + GL on POST
- [ ] **3.5** Stock Adjustment — threshold approval (<5M/5-100M/>100M)
- [ ] **3.6** Stock Take — recount workflow (tolerance, max 3 counts, blind count)
- [ ] **3.7** Stock Take — POST→auto-generate adjustment + GL
- [ ] **3.8** Stock Take — ABC scheduling (A=monthly, B=quarterly, C=yearly)
- [ ] **Checkpoint** — All 9 entities mutate stock + post GL. 40+ tests.

## Phase 4: Special Movements

- [ ] **4.1** Direct Receipt entity (donation, owner_contribution, surplus, sample)
- [ ] **4.2** Direct Receipt — threshold approval
- [ ] **4.3** Sales Return Receipt — credit note→warehouse
- [ ] **4.4** Return to Supplier — pre/post invoice
- [ ] **4.5** Internal Issue — department + cost_center + purpose GL
- [ ] **4.6** Special movement tests (GL assertions per type)
- [ ] **Checkpoint** — All 30 movement types covered. 50+ tests.

## Phase 5: Cost & Valuation Engine

- [ ] **5.1** Valuation — Weighted Average (periodic) calculation
- [ ] **5.2** Valuation — FIFO layer sorting
- [ ] **5.3** Valuation — Specific ID matching
- [ ] **5.4** Cost Revaluation entity
- [ ] **5.5** Valuation — ROLLED_BACK state + reversal
- [ ] **5.6** Valuation engine tests
- [ ] **Checkpoint** — Real cost computation. All 4 methods tested.

## Phase 6: Period-End

- [ ] **6.1** NRV provision 2294 — calculation + GL
- [ ] **6.2** Opening Balance → StockBalance wiring
- [ ] **6.3** Period lock enforcement
- [ ] **6.4** Period-end close workflow (validate→value→provision→lock)
- [ ] **Checkpoint** — Period-end close tested end-to-end.

## Phase 7: Integration Wiring

- [ ] **7.1** Purchase → GRN auto-create + PO received_qty tracking
- [ ] **7.2** Sales → DN auto-create + SO delivered_qty tracking
- [ ] **7.3** Credit Note → Sales Return Receipt
- [ ] **7.4** Debit Note → Supplier Return
- [ ] **7.5** End-to-end integration tests (PO→GRN→Stock→DN→COGS→Invoice)
- [ ] **Checkpoint** — Purchase-to-pay + order-to-cash integrated. Both backends tested.

---

## Test Target Growth

| Phase | Running Total | Cumulative |
|-------|---------------|------------|
| Current | 22 | 22 |
| Phase 0 | +0 | 22 |
| Phase 1 | +7 | 29 |
| Phase 2 | +7 | 36 |
| Phase 3 | +10 | 46 |
| Phase 4 | +10 | 56 |
| Phase 5 | +6 | 62 |
| Phase 6 | +6 | 68 |
| Phase 7 | +8 | 76 |

Target: **76+ tests** at completion. All use in-memory repos, real service, httptest, no mocks.
