# ADR-001: Purchase/AP Module Architecture

**Status:** Accepted
**Date:** July 2026
**Author:** BA Lead + Chief Accountant (20+ yrs each)
**Regulatory Basis:** Circular 99/2025/TT-BTC, Decree 123/2020/ND-CP (amended by 70/2025), Decree 254/2026/ND-CP

## Context

Purchase/AP module for Vietnamese tax-compliant General Ledger. Needs full Procure-to-Pay (P2P) cycle: supplier management, PO, goods receipt, supplier invoice, AP tracking, payment. Must comply Circular 99/2025/TT-BTC (effective 1 Jan 2026) and Decree 123/2020 e-invoice rules.

Current module is ~70% built with known gaps: 3-way matching, GL auto-posting, e-invoice GDT integration.

## Decision

### Architecture

Follow GoTax MVC + Repository pattern:

```
Handler → Service (business rules) → Repository (PG/Memory)
   ↑                                    ↓
authMW ← Request                      GORM
```

### Key Design Decisions

1. **All-in-one memory repo** — Single `MemoryPurchaseRepo` implements all 6 interfaces, per existing GoTax convention. PG repos are separate per interface.

2. **Two backend support** — PG for production, memory for tests. No mock services needed.

3. **Domain models with self-validation** — `Validate()` methods on all models with enum constants and state machine transitions (`ValidTransition()`).

4. **TEXT columns for dates in migration** — SQLite-compatible migration syntax. GORM handles `time.Time` ↔ TEXT conversion automatically. PG stores dates as ISO format text.

5. **String UUIDs** — All PKs use string UUIDs (36 chars), not auto-increment integers. Consistent with existing GoTax convention.

6. **E-invoice integration deferred** — `ReceiveEInvoice()` method exists in service layer but no GDT XML parser. Build as adapter layer when GDT API stabilizes.

7. **3-way matching deferred** — `ErrInvoice3WayMismatch` defined. Engine will match PO qty × GRN qty × Invoice qty per line with configurable tolerance.

8. **GL auto-posting deferred** — `PostInvoice` sets `gl_posted=true` but creates no journal entry. Post to 331/1331 using existing `JournalEntry` engine.

### Data Model

9 tables in `000006_purchase_schema`:

| Table | Purpose | Key Columns |
|-------|---------|-------------|
| `suppliers` | Vendor master | tax_code, payment_terms, supplier_type |
| `purchase_orders` | PO header | po_number (auto), status state machine |
| `po_lines` | PO lines | account_id (152/156/642), vat_rate |
| `goods_receipt_notes` | GRN header | po_id FK, receipt_date |
| `grn_lines` | GRN lines | quantity_received, quantity_rejected |
| `supplier_invoices` | Invoice header | e_invoice_code, vat_deduction_status |
| `invoice_lines` | Invoice lines | account_id, vat_account_id |
| `ap_transactions` | AP sub-ledger | transaction_type (invoice/credit_note/payment/prepayment/offset) |
| `cost_allocations` | Cost distribution | cost_type, allocation_method, allocated_lines (JSON) |

### API Routes

28 endpoints under `/api/v1/purchase/`:
- `/suppliers/` — CRUD (5)
- `/orders/` — CRUD + approve/cancel/close (7)
- `/receipts/` — CRUD + post/cancel (6)
- `/invoices/` — CRUD + verify/post/cancel/claim-vat (8)
- `/ap/aging` + `/ap/summary` — AP reports (2)

Company-scoped via `company_id` query param. Auth middleware on all routes.

### Status State Machines

**PO:** DRAFT → APPROVED → SENT → PARTIAL → RECEIVED → CLOSED (CANCELLED from any)
**GRN:** DRAFT → POSTED (CANCELLED from DRAFT)
**Invoice:** DRAFT → VERIFIED → POSTED → PAID (CANCELLED from DRAFT/VERIFIED)
**Supplier:** ACTIVE / SUSPENDED / BLACKLISTED

### Vietnamese Tax Compliance

- VAT deduction tracking per invoice (pending → claimed → rejected)
- Supplier tax code (MST) required — 10/13/14 digit format
- Account mapping: 152/156 (inventory), 331 (AP), 1331 (input VAT), 642 (expenses)
- E-invoice data stored as XML string, validation deferred
- AP aging tracked via ap_transactions sub-ledger

## Alternatives Considered

### Separate invoice allocation tables
- **Rejected:** Combined cost allocation table with `allocated_lines` JSON field simpler, matches Vietnamese practice of allocation by qty/value/weight/volume.

### Auto-increment IDs from domain model
- **Rejected:** String UUIDs consistent with existing GoTax, better for distributed systems, prevents ID collision in multi-tenant setup.

### gorm.Model embedding
- **Rejected:** Explicit column tags for full control over migration-matching. gorm.Model adds DeletedAt which we don't use.

## Consequences

1. Migration 000006 must stay in sync with GORM models — all future model changes require both `.up.sql` and GORM struct update.
2. Memory repo permits fast handler/service tests without DB.
3. Module split (purchase vs separate AP module) avoided — Purchase handles both procurement and AP.
4. Future inventory module will consume GRN data from purchase tables.
5. Bank module integration for payment execution.

## Gaps Fixed in v2.0

- GORM models rewritten to match migration schema (table names, column tags, field types)
- PG mapping functions fixed (CostAllocation data corruption, APTransaction field loss)
- Legacy unversioned migration `005_purchase_schema.sql` removed
- AGENTS.md readiness updated from ~0% to ~70%
- READINESS.md rewritten with accurate gap analysis

## References

- `internal/domain/models_purchase.go` — Domain models (417 lines)
- `internal/domain/models_gorm_purchase.go` — GORM models (fix: match migration)
- `internal/repository/pg_purchase.go` — PG repo (fix: correct mapping)
- `internal/repository/memory_purchase.go` — Memory repo (886 lines)
- `internal/service/purchase_service.go` — Business logic (450 lines)
- `internal/handler/purchase_handler.go` — HTTP handlers (428 lines)
- `migrations/000006_purchase_schema.up.sql` — DB schema (209 lines)
- `docs/purchase/PURCHASE_READINESS.md` — Updated gap analysis (v2.0)
