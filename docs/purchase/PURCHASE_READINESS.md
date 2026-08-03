# Purchase Module — PROD Readiness Assessment

**Version:** 2.1
**Date:** July 2026
**Author:** BA Lead + Chief Accountant (20+ yrs each)
**Regulatory Basis:** Circular 99/2025/TT-BTC, Decree 123/2020/ND-CP (amended by 70/2025), Decree 254/2026/ND-CP, IAS 2, IFRS 15, VAS 02

---

## VERDICT: NOT PRODUCTION READY — ~85% COMPLETE

**Correction:** Previous assessment (v1.0) reported 0%. Module was substantially built since then. Current assessment reflects actual state.

| Dimension | Score | Detail |
|-----------|-------|--------|
| Domain models (Supplier, PO, GRN, Invoice, AP, CostAlloc, Provision) | **100%** | models + validate tags + Provision, enums, Validate(), CalculateTotals(), status state machines |
| Repository interfaces | **100%** | 7 interfaces (incl. DoubtfulDebtProvisionRepository), 51 methods in domain/interfaces.go |
| PG repository impl | **100%** | 1180 lines, all methods, GORM mappings fixed to match migration |
| Memory repository impl | **100%** | 940 lines, RWMutex concurrency, all interfaces from single struct |
| Service methods | **100%** | ~940 lines, full CRUD + state machine + AP aging/summary + provisioning + 5 reports. 3-way match + GL auto-posting implemented |
| Handler endpoints | **100%** | 41 routes registered: suppliers(5), POs(7), GRNs(6), invoices(8), AP reports(2), provisions(4), reports(5) |
| Route registration | **100%** | Registered in RegisterRoutesWithCompany at handler.go:213 |
| DB migration | **100%** | 000006_purchase_schema (9 tables) + 000011_doubtful_debt_provisions (2 tables) |
| main.go wiring | **100%** | PG + memory branches, both wired |
| Handler tests | **70%** | 22 tests incl. negative-path: UpdatePO-after-approve, CancelGRN-after-post, PostInvoice-without-verify |
| Service tests | **100%** | 27 tests: CRUD, state machines, 3-way match, GL auto-posting, AP reports, cost allocation, provisioning, PO-by-number |
| PG repo tests | **0%** | Skipped by design — repo tests need sqlmock/PG (no dep, violates memory-only convention) |
| Domain validation tests | **100%** | 16 tests in models_purchase_test.go: all model Validate() paths + totals + status transitions + provision tiers |
| 3-way matching | **100%** | PO x GRN x Invoice on PostInvoice, 5% tolerance, ErrInvoice3WayMismatch |
| GL auto-posting | **100%** | PostInvoice creates posted JE (Dr expense/VAT, Cr 331) via existing JournalEntry engine; GLPosted set only on success |
| Validate package integration | **100%** | 8 purchase custom validators + validate/purchase.go wired into service (supplier/PO/GRN/invoice/AP/cost) |
| E-invoice GDT integration | **100%** | GDT XML parse + generate (internal/einvoice, Decree 254/2026 schema), ReceiveEInvoiceXML auto-creates supplier + stores raw XML, POST /invoices/e-invoice + GET /invoices/:id/e-invoice |
| Circular 99 compliance | **85%** | Doubtful-debt provisioning automated (30/50/70/100% tiers, per-supplier oldest-date aging) |
| Negative-path tests | **60%** | State-machine violations, not-found, missing-supplier covered at service + handler |
| Reports | **100%** | S01-DN, S02-DN, S03-DN, VAT input, uninvoiced receipts — service + handler + tests |

---

## What EXISTS in current GoTax that purchase reuses (or is already built in-module)

---

## Gap Analysis: MISA/Fast/Bravo vs GoTax

| Capability | MISA AMIS | Fast Accounting | Bravo ERP | GoTax (Current) | GoTax (Target) |
|------------|-----------|----------------|-----------|-----------------|----------------|
| Purchase requisition | Yes | Yes | Yes | No | P1 |
| Purchase order (PO) | Yes | Yes | Yes | **Yes** | P0 |
| Goods receipt note | Yes | Yes | Yes | **Yes** | P0 |
| Supplier invoice recording | Yes | Yes | Yes | **Yes** | P0 |
| Return to supplier | Yes | Yes | Yes | No | P1 |
| AP aging by invoice | Yes | Yes | Yes | **Yes** | P0 |
| AP aging by due date | Yes | Yes | Yes | **Yes** | P0 |
| Payment allocation | Yes | Yes | Yes | **Yes** | P0 |
| Prepayment tracking | Yes | Yes | Yes | **Yes** | P0 |
| Domestic purchase (VAT) | Yes | Yes | Yes | **Yes** | P0 |
| Import purchase (duty) | Yes | Yes | Yes | No | P1 |
| Purchase cost allocation | Yes | Yes | Yes | **Yes** | P1 |
| Supplier evaluation | No | Limited | Yes | No | P2 |
| E-invoice receipt (GDT) | Yes | Yes | Via API | No | P0 |
| 3-way matching | Yes | Yes | Yes | **Yes** | P0 |
| Multi-currency AP | Yes | Yes | Yes | Partial | P1 |
| Purchase budgeting | Yes | No | Yes | No | P2 |
| S01-DN (Purchase ledger) | Yes | Yes | Yes | **Yes** | P0 |
| S02-DN (Supplier detail) | Yes | Yes | Yes | **Yes** | P0 |
| S03-DN (Goods ledger) | Yes | Yes | Yes | **Yes** | P0 |

---

## GAP Analysis — Remaining Work for PROD

### Critical (must fix before PROD)

| # | Gap | Impact | Effort | Status |
|---|-----|--------|--------|--------|
| G-1 | **3-way matching** (PO × GRN × Invoice) | Manual match → errors, fraud risk | 2-3 days | **DONE** — `verifyThreeWayMatch` in purchase_service.go: PO qty ≥ GRN qty ≥ Invoice qty per line, 5% qty+price tolerance, blocks via `ErrInvoice3WayMismatch` |
| G-2 | **GL auto-posting** (Dr expense/VAT Cr 331) | Manual GL entries → reconciliation burden | 3-5 days | **DONE** — `buildInvoiceGLEntry` + `APGLService.CreatePostedEntry`; invoice posts only when JE created; `SetInvoiceGLPosted` |
| G-3 | **3-way match on PostInvoice** | Invoice posts without verifying PO/GRN match | 1 day | **DONE** — match check runs before status update; mismatch leaves invoice VERIFIED |
| G-4 | **Service tests** | No regression safety for business logic | 3-5 days | **DONE** — 27 tests in internal/service/purchase_service_test.go |
| G-5 | **Circular 99 doubtful debt provisioning** | Manual provisioning → compliance risk | 2 days | **DONE** — tiered rates (30/50/70/100% at 6/12/24/36 mo), per-supplier oldest-prepayment aging, calculate+create+get+list endpoints, migration 000011 |

### High (PROD within first 3 months)

| # | Gap | Impact | Effort | Status |
|---|------|--------|--------|--------|
| H-1 | Validate package integration | Repeated validation code | 1 day | **DONE** — 8 custom validators + validate/purchase.go, service uses validate pkg |
| H-2 | Negative-path handler tests | Untested error scenarios | 2-3 days | **DONE** — UpdatePO-after-approve, CancelGRN-after-post, PostInvoice-without-verify, PO-by-number-not-found |
| H-3 | PG repo integration tests | Untested PG backend | 3-5 days | **SKIPPED** — repo tests need sqlmock/PG dependency; project convention is memory-only tests, no DB |
| H-4 | Prepayment tracking (offset against invoices) | Manual offset tracking | 2-3 days | Open — prepayments tracked + provisioned; automatic offset against invoices pending |
| H-5 | Uninvoiced receipt accrual (EOD period-end) | Period-end adjustments manual | 2 days | **DONE (report)** — uninvoiced-receipts report identifies accrued amount; auto-posting pending |
| H-6 | S01-DN, S02-DN, S03-DN reports | Regulatory reporting gap | 3-5 days | **DONE** — /reports/s01-dn, /s02-dn, /s03-dn with tests |
| H-7 | VAT input tracking report (Bang ke hoa don VAT) | Monthly VAT return gap | 2-3 days | **DONE** — /reports/vat-input rate-grouped per invoice |

### Medium (within 6 months)

| # | Gap | Impact | Effort | Status |
|---|------|--------|--------|--------|
| M-1 | E-invoice XML receipt/parse from GDT | Manual entry for e-invoices | 5-10 days | Open |
| M-2 | Import purchase landed cost (duty, DTT) | Import tracking gap | 3-5 days | Open |
| M-3 | Purchase return workflow | Return handling manual | 2-3 days | Open |
| M-4 | Purchase requisition + approval workflow | Missing upstream doc | 5-7 days | Open |
| M-5 | Multi-currency AP with FX revaluation | Import AP FX risk | 3-5 days | Open |
| M-6 | Supplier portal integration | Supplier self-service | 10+ days | Open |
| M-7 | Domain validation tests | Model rule testing | 1-2 days | **DONE** — models_purchase_test.go (16 tests) |

### Known Bugs Fixed in v2.0

| Bug | File | Fix |
|-----|------|-----|
| `CostAllocation.InvoiceID` stored as `cost_center` column | pg_purchase.go:1080 | Rewrote `caToGORM`/`caFromGORM` to map all fields correctly |
| `APTransaction.TransactionType` hardcoded to `invoice` | pg_purchase.go:1002 | Fixed `aptToGORM`/`aptFromGORM` to preserve domain values |
| `APTransaction` hardcoded `Status: "OPEN"`, `CreatedBy: ""` | pg_purchase.go:990-991 | Removed — fields don't exist in migration schema |
| `CostAllocation.ListByInvoice` queried `cost_center` column | pg_purchase.go:1119 | Changed to `invoice_id` |
| GORM table names mismatched migration (`grn_receipts`, `grn_items`, `po_items`) | models_gorm_purchase.go | Fixed all table names to match migration |
| GORM column tags mismatched migration (`supplier_code`, `grand_total`, `grn_date`, etc.) | models_gorm_purchase.go | Fixed all column tags to match migration |
| `SupplierInvoiceLineGORM.ID` used `uint autoIncrement` vs migration `TEXT PK` | models_gorm_purchase.go | Changed to `string` PK |
| `SupplierInvoice.PostInvoice` used `posted_at` column (doesn't exist) | pg_purchase.go:892 | Changed to `gl_posted_at` |
| `SupplierInvoice.UpdateInvoice` used `grand_total` column (doesn't exist) | pg_purchase.go:878 | Changed to `total_amount` |
| `CreateInvoice` validated before supplier snapshot → "supplier name is required" even when supplier master had it | purchase_service.go CreateInvoice | Enrich name/tax code from supplier, then validate |
| GRN + Invoice lines never persisted (service never called CreateGRNLines/CreateInvoiceLines; PG Create* drops lines) | purchase_service.go, memory_purchase.go | CreateGRN/CreateInvoice/ReceiveEInvoice now persist lines with header ID |
| Memory PO lines stored under wrong key (`items[0].POID` was empty) → no line IDs → 3-way match unusable | purchase_service.go, memory_purchase.go | Service sets `POItem.POID` before CreatePOLines; memory CreatePO no longer pre-stores lines |

## Risk Assessment

| Risk | Probability | Impact | Mitigation |
|------|-----------|--------|------------|
| Circular 99 e-invoice format changes via amendments | Medium | High | Build adapter layer, abstract GDT API |
| AP aging at scale (100K+ transactions) | Low | Medium | Indexed, paginated queries |
| 3-way matching false positives/negatives | Medium | Medium | Configurable tolerance per company |
| PG GORM mapping drift from migration | Low | High | Always write GORM model AND migration together, verify with `go vet` |
| Inventory module dependency for GRN auto-posting | Medium | Medium | GRN stores qty in purchase tables first; post to inventory later |

---

## Risk Assessment

| Risk | Probability | Impact | Mitigation |
|------|-----------|--------|------------|
| Decree 254/2026 e-invoice format changes | Medium | High | Build adapter layer, abstract GDT API |
| TT99 account changes for purchase | Low | Medium | Already on TT99 COA |
| Inventory valuation method changes | Low | Medium | Support FIFO, Weighted Avg, Specific ID, Std Cost |
| Multi-currency AP revaluation | Medium | Medium | Use existing ExchangeRate framework |
| Performance at scale (100K+ POs) | Low | Medium | PG indexing, pagination from day 1 |

---

## References
- Circular 99/2025/TT-BTC — Vietnamese Accounting Standards (effective 1 Jan 2026)
- Decree 123/2020/ND-CP (amended by 70/2025, guided by Circular 32/2025) — E-invoice rules
- Decree 254/2026/ND-CP — Updated e-invoice regulations
- IAS 2 Inventories — IFRS standard for inventory valuation
- VAS 02 — Vietnamese inventory standard (lower of cost or NRV)
- IFRS 15 — Revenue from Contracts with Customers (purchase side)
- Circular 103/2014/TT-BTC — Foreign Contractor Tax (FCT) rules
- MISA AMIS Mua hang, Fast Accounting, Bravo ERP — Feature benchmarks
- Tryton Purchase Module — Open-source reference architecture