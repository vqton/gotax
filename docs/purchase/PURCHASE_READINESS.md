# Purchase Module — PROD Readiness Assessment

**Version:** 2.0
**Date:** July 2026
**Author:** BA Lead + Chief Accountant (20+ yrs each)
**Regulatory Basis:** Circular 99/2025/TT-BTC, Decree 123/2020/ND-CP (amended by 70/2025), Decree 254/2026/ND-CP, IAS 2, IFRS 15, VAS 02

---

## VERDICT: NOT PRODUCTION READY — ~70% COMPLETE

**Correction:** Previous assessment (v1.0) reported 0%. Module was substantially built since then. Current assessment reflects actual state.

| Dimension | Score | Detail |
|-----------|-------|--------|
| Domain models (Supplier, PO, GRN, Invoice, AP, CostAlloc) | **100%** | 417 lines, enums, Validate(), CalculateTotals(), status state machines |
| Repository interfaces | **100%** | 6 interfaces, 46 methods in domain/interfaces.go |
| PG repository impl | **100%** | 1127 lines, all 46 methods, GORM mappings fixed to match migration |
| Memory repository impl | **100%** | 886 lines, RWMutex concurrency, all interfaces from single struct |
| Service methods | **95%** | 450 lines, full CRUD + state machine + AP aging/summary. Missing: 3-way match, GL auto-posting |
| Handler endpoints | **100%** | 28 routes registered: suppliers(5), POs(7), GRNs(6), invoices(8), AP reports(2) |
| Route registration | **100%** | Registered in RegisterRoutesWithCompany at handler.go:213 |
| DB migration | **100%** | 000006_purchase_schema: 9 tables, 21 indexes, proper FKs |
| main.go wiring | **100%** | PG + memory branches, both wired |
| Handler tests | **60%** | 19 tests covering supplier, PO, GRN, invoice workflows, AP reports |
| Service tests | **0%** | No purchase_service_test.go |
| PG repo tests | **0%** | No pg_purchase_test.go |
| Domain validation tests | **0%** | No model validation tests |
| 3-way matching | **0%** | ErrInvoice3WayMismatch defined but unused |
| GL auto-posting | **0%** | PostInvoice marks gl_posted=true but creates no journal entry |
| Validate package integration | **0%** | No purchase-specific validators in internal/validate/ |
| E-invoice GDT integration | **0%** | ReceiveEInvoice creates invoice but no XML parse/GDT API |
| Circular 99 compliance | **70%** | Tax rules, FX, provisioning defined but not fully automated |
| Negative-path tests | **20%** | Basic error cases covered, many edge cases untested |

---

## What EXISTS in current GoTax that purchase reuses (or is already built in-module)

---

## Gap Analysis: MISA/Fast/Bravo vs GoTax

| Capability | MISA AMIS | Fast Accounting | Bravo ERP | GoTax (Current) | GoTax (Target) |
|------------|-----------|----------------|-----------|-----------------|----------------|
| Purchase requisition | Yes | Yes | Yes | No | P1 |
| Purchase order (PO) | Yes | Yes | Yes | No | P0 |
| Goods receipt note | Yes | Yes | Yes | No | P0 |
| Supplier invoice recording | Yes | Yes | Yes | No | P0 |
| Return to supplier | Yes | Yes | Yes | No | P1 |
| AP aging by invoice | Yes | Yes | Yes | No | P0 |
| AP aging by due date | Yes | Yes | Yes | No | P0 |
| Payment allocation | Yes | Yes | Yes | No | P0 |
| Prepayment tracking | Yes | Yes | Yes | No | P0 |
| Domestic purchase (VAT) | Yes | Yes | Yes | No | P0 |
| Import purchase (duty) | Yes | Yes | Yes | No | P1 |
| Purchase cost allocation | Yes | Yes | Yes | No | P1 |
| Supplier evaluation | No | Limited | Yes | No | P2 |
| E-invoice receipt (GDT) | Yes | Yes | Via API | No | P0 |
| 3-way matching | Yes | Yes | Yes | No | P0 |
| Multi-currency AP | Yes | Yes | Yes | No | P1 |
| Purchase budgeting | Yes | No | Yes | No | P2 |
| S01-DN (Purchase ledger) | Yes | Yes | Yes | No | P0 |
| S02-DN (Supplier detail) | Yes | Yes | Yes | No | P0 |
| S03-DN (Goods ledger) | Yes | Yes | Yes | No | P0 |

---

## GAP Analysis — Remaining Work for PROD

### Critical (must fix before PROD)

| # | Gap | Impact | Effort | Existing Hook |
|---|-----|--------|--------|---------------|
| G-1 | **3-way matching** (PO × GRN × Invoice) | Manual match → errors, fraud risk | 2-3 days | `ErrInvoice3WayMismatch` defined in domain/errors.go:322 |
| G-2 | **GL auto-posting** (Dr 331/1331 Cr 331) | Manual GL entries → reconciliation burden | 3-5 days | Invoice.PostInvoice → `gl_posted` boolean set but no JE created. Use existing `JournalEntry` engine in service.go |
| G-3 | **3-way match on PostInvoice** | Invoice posts without verifying PO/GRN match | 1 day | Add match check in `PostInvoice` before status update |
| G-4 | **Service tests** | No regression safety for business logic | 3-5 days | 0 tests in internal/service/purchase_service_test.go |
| G-5 | **Circular 99 doubtful debt provisioning** | Manual provisioning → compliance risk | 2 days | Tiered rules (30/50/70/100%) from research ready |

### High (PROD within first 3 months)

| # | Gap | Impact | Effort |
|---|------|--------|--------|
| H-1 | Validate package integration | Repeated validation code | 1 day |
| H-2 | Negative-path handler tests | Untested error scenarios | 2-3 days |
| H-3 | PG repo integration tests | Untested PG backend | 3-5 days |
| H-4 | Prepayment tracking (offset against invoices) | Manual offset tracking | 2-3 days |
| H-5 | Uninvoiced receipt accrual (EOD period-end) | Period-end adjustments manual | 2 days |
| H-6 | S01-DN, S02-DN, S03-DN reports | Regulatory reporting gap | 3-5 days |
| H-7 | VAT input tracking report (Bang ke hoa don VAT) | Monthly VAT return gap | 2-3 days |

### Medium (within 6 months)

| # | Gap | Impact | Effort |
|---|------|--------|--------|
| M-1 | E-invoice XML receipt/parse from GDT | Manual entry for e-invoices | 5-10 days |
| M-2 | Import purchase landed cost (duty, DTT) | Import tracking gap | 3-5 days |
| M-3 | Purchase return workflow | Return handling manual | 2-3 days |
| M-4 | Purchase requisition + approval workflow | Missing upstream doc | 5-7 days |
| M-5 | Multi-currency AP with FX revaluation | Import AP FX risk | 3-5 days |
| M-6 | Supplier portal integration | Supplier self-service | 10+ days |
| M-7 | Domain validation tests | Model rule testing | 1-2 days |

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