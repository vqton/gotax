# AR Module — Implementation Roadmap v2

## Overview

Enterprise AR module for GoTax. Target: MISA/FAST/Bravo ERP parity for order-to-cash collection workflow. Phases 0-1 complete (GL auto-posting, aging fix, customer statement). Remaining: collection management (dunning/bad debt/prepayment/FX) and controls (credit limit/recon/KPI/off-BS).

Full lifecycle doc: `docs/sale/AR_WORKFLOW.md`

---

# Purchase Module — P2 Features (TDD)

## Overview

Close the remaining P2 gaps flagged in PURCHASE_READINESS.md v2.1. Five features, each a full vertical slice implemented TDD (RED → GREEN → REFACTOR), committed separately. Order = dependency + risk.

## Order

1. **P2-1 Purchase requisition + approval (M-4)** — pure additive, foundation for PO-from-requisition
2. **P2-2 Purchase return workflow (M-3)** — return GRN + credit note + AP offset + GL reversal
3. **P2-3 Import purchase + landed cost (M-2)** — import duty (3333) + import VAT (33312) on invoice, landed cost breakdown
4. **P2-4 Multi-currency AP FX revaluation (M-5)** — revalue AP at period-end → 515/635
5. **P2-5 E-invoice GDT XML (M-1)** — generate + parse GDT XML for ReceiveEInvoice

## Architecture Decisions

- **Per-feature TDD**: domain validation tests RED first, then service tests, then handler tests. Memory repos used throughout (AGENTS.md convention — no DB).
- **Single interface additions per feature**: new repository interface (Requisition, Return, etc.) added to `domain/interfaces.go`; MemoryPurchaseRepo + PGPurchaseRepo implement.
- **Constructor grows**: `NewPurchaseService` gains one repo param per feature that needs a new repo (9→10 args so far: requisition + FX). All call sites updated: main.go ×2, purchase_service_test.go ×2, purchase_handler_test.go ×1.
- **Migrations sequential**: 000012_requisitions, 000013_returns, 000014_landed_cost, 000015_fx_revaluation. Each with .up + .down.
- **Validate package**: register one custom validator per new enum (`reqstatus`, `fxrstatus`), add `validate.Requisition`/`validate.FXRevaluation` mapper functions.
- **E-invoice XML**: pure-Go `encoding/xml`, no new deps. GDT invoice XML generation + parse, wired into ReceiveEInvoice + a generate endpoint.
- **FX revaluation**: reuse existing ExchangeRate framework + GL CreatePostedEntry. Realized on payment, unrealized at period-end.

## Checkpoints

- After P2-1: `go vet ./... && go test -count=1 ./...` green, commit — ✅ DONE (a9b5c85 + requisition commit)
- After P2-2: same + commit — ✅ DONE (daa69a4 + returns commit)
- After P2-3: same + commit — ✅ DONE (8a8274a + import commit)
- After P2-4: same + commit — ✅ DONE (e0d5cf8 + fx commit)
- After P2-5: same + commit — ✅ DONE (c3068c8 + 8a484a5 + 5074d93 + docs commit)
- Final: docs bumped (readiness → P2 closed), AGENTS.md updated — ✅ DONE (this commit)

## P2-5 Slice Plan — GDT E-Invoice XML (M-1)

Format authority: `docs/purchase/PURCHASE_TEMPLATES.md` §8 (Decree 254/2026 schema). Pure-Go `encoding/xml`, no new deps. New package `internal/einvoice` (adapter layer, per risk mitigation). No migration (EInvoiceData/EInvoiceCode columns exist). No constructor change (reuses invRepo + supRepo).

| Slice | Scope | Deliverable |
|-------|-------|-------------|
| 5.1 | `internal/einvoice/gdt.go` + `gdt_test.go` | Parse([]byte)→SupplierInvoice, Generate(inv)→[]byte. GL account defaults 152/1331 on lines. VAT rate→VATType map. |
| 5.2 | `purchase_service.go` + `purchase_service_test.go` | `ReceiveEInvoiceXML(ctx, companyID, raw)` (parse→ReceiveEInvoice→store raw in EInvoiceData), `GenerateEInvoiceXML(ctx, id)`. |
| 5.3 | `purchase_handler.go` + `purchase_handler_test.go` | `POST /invoices/e-invoice` (raw XML body), `GET /invoices/:id/e-invoice` (XML response). |
| 5.4 | docs + AGENTS.md + plan.md | Readiness rows 0%→100%, AGENTS.md module table, checkpoint. |

Decisions:
- XML→domain line mapping sets AccountID=152, VATAccountID=1331 (GL mapping is app concern, not GDT payload).
- VAT rate to VATType: 0→VAT_0, 5→VAT_5, 8→VAT_8, 10→VAT_10, -1/KCT→NT.
- Generate omits empty Seller/Buyer optional fields; Buyer carries no company data today (only name/tax code available from invoice).
- XML namespace kept as `http://gdt.gov.vn/schemas/einvoice/2026` per template.

## Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| Constructor param explosion | Med | Group into opts struct if 13+; acceptable now |
| Return GL reversal double-posting | High | Unit-test GL entry generation; reuse APGLService |
| FX revaluation math drift | Med | round2 helper, deterministic tests |
| GDT XML format drift | Med | Abstract parser; wire format version in header |


## Architecture Decisions

- **GL posting via JournalEntry.CreatePostedEntry** — skip review cycle for auto-generated entries. Already implemented.
- **Dunning as separate service** — `CollectionService` with own repos/models. Avoids sale_service.go bloat.
- **Bad debt by VAS 17** — aging-bucket % provision + specific identification for known bads.
- **Prepayment as Receipt type** — extend CustomerReceipt with `ReceiptType` (standard/prepayment/refund). No separate table.
- **FX at receipt time** — compare invoice exchange rate vs receipt rate, post realized gain/loss. Month-end batch for unrealized.
- **Credit limit** — check at SO Confirm + Invoice Create. Configurable warn/block per company.
- **AR-GL reconciliation** — compare sub-ledger total vs GL 131 balance. Variance drill-down.
- **E-invoice (Phase 2)** deferred per PI direction. Revisit when GDT API integration scope is clear.

## Task List

### ✅ Phase 0: Critical Fix (COMPLETE)

- [x] Task 1: Fix AR aging — bucket by DueDate (was putting all in Bucket0)

### ✅ Phase 1: Core AR (COMPLETE)

- [x] Task 2: GL auto-posting on invoice post (Dr 131, Cr revenue 5111/3331)
- [x] Task 3: GL auto-posting on receipt post (Dr 1111/1121, Cr 131)
- [x] Task 4: GL auto-posting on credit note post (reverse entries: Dr 5111/3331, Cr 131)
- [x] Task 5: Customer statement endpoint (merged invoice/receipt/CN timeline + running balance)

### ⏸️ Phase 2: E-Invoice Pipeline (DEFERRED)

Phase 2 is deferred—excluded from current AR deliverable. Design doc exists in `docs/sale/AR_WORKFLOW.md`. Revisit after GDT API integration scope is clarified. Task details retained in `todo.md` for reference.

### 🔲 Phase 3: Collection Management (P1) — NEXT

| # | Task | Scope | Depends On |
|---|------|-------|------------|
| 6 | Dunning engine — levels, schedule, auto-reminder, promise-to-pay | L (6-8 files) | Phase 1 ✓ |
| 7 | Bad debt provision calc (VAS 17) + write-off workflow | L (5-7 files) | Task 6 |
| 8 | Prepayment/deposit workflow — type, offset, refund | M (4-5 files) | Phase 1 ✓ |
| 9 | FX revaluation — realized at receipt + month-end unrealized | M (4-5 files) | Phase 1 ✓ |

### 🔲 Phase 4: Controls & Monitoring (P1/P2)

| # | Task | Scope | Depends On |
|---|------|-------|------------|
| 10 | Credit limit enforcement (SO confirm + invoice create) | S (2-3 files) | Phase 1 ✓ |
| 11 | AR-GL month-end reconciliation report | M (3-4 files) | Phase 1 ✓ |
| 12 | DSO + AR KPI dashboard | M (3-4 files) | Phase 1 ✓ |
| 13 | Off-balance-sheet tracking for written-off AR | M (4-5 files) | Task 7 |

## Build Order

```
Phase 3 ──→ 8 (prepayment: standalone, easy win)
         │  9 (FX revaluation: standalone, uses existing exchange rates)
         ├──→ 6 (dunning: biggest scope, start early)
         │         │
         │         └──→ 7 (bad debt: depends on dunning for aging data)
         │
Phase 4 ──→ 10 (credit limit: trivial check)
         │  11 (AR-GL recon: query-based, no new data)
         │  12 (KPI: calc-based, no new data)
         └──→ 13 (off-BS: depends on bad debt write-off)
```

## Recommended Order (by value/complexity ratio)

1. **Task 10** — Credit limit. Smallest effort, regulatory req, protects cash flow.
2. **Task 8** — Prepayment. Common workflow, foundation for deposit handling.
3. **Task 11** — AR-GL recon. Critical for month-end close, banks won't certify without it.
4. **Task 9** — FX revaluation. Required for multi-currency AR. Regulatory requirement per Circular 99.
5. **Task 12** — KPI dashboard. Reporting. Lower code risk.
6. **Task 6** — Dunning. Largest effort. High business value but complex.
7. **Task 7** — Bad debt. Depends on dunning data. High regulatory value.
8. **Task 13** — Off-balance-sheet. Depends on bad debt. Lowest urgency.

## Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| Dunning email/SMS integration complexity | High | Start with status-only tracking, add notification later |
| Bad debt provision % not regulatory-aligned | Medium | Make % configurable, default to VAS 17 rates |
| FX rate source (daily rate from central bank?) | Low | Manual rate input for now; API integration later |
| Credit limit stale (customer pays between check and create) | Low | Race acknowledged — detect at next statement |

---

# Tax Core Foundations — Shared Infrastructure (TDD)

**Status: Round A COMPLETE (2026-08-03, commits 3e84268→c6c3a86). Round B next.**

## Overview

Build the shared tax foundations flagged in TAX_READINESS.md: configurable rate engine, real VAT/CIT/PIT calculation, declaration engine, declaration→payment automation, GDT client + XMLDSig. All TDD (RED→GREEN), memory repos, per-slice commits. Order = dependency + risk.

Sources: `docs/tax/TAX_RULES.md` (rules), `docs/tax/TAX_SPECS.md` (forms), `docs/tax/TAX_READINESS.md` (gaps).

## Current State (verified)

- `TaxServiceInterface` + `taxService` exist. CRUD for declarations/rates/payments/e-invoices/calendar/alerts/audit DONE.
- `CalculateVAT`/`CalculateCIT`: **hardcoded** rates (10%/20%), naive account classification.
- `CalculatePIT`: **stub** — returns empty `PITResult{}`.
- **Zero** `tax_service_test.go` (handler tests only, 610 lines).
- `MemoryTaxRepo` in `internal/repository/memory.go:1152`; PG in `pg_tax.go`.
- `TaxRepository` interface in `interfaces.go:194` — no rate-lookup-by-date method.
- No HTTP client, no XMLDSig anywhere in codebase. `internal/einvoice` = parse/generate only.
- `DigitalSignature` model = metadata only (no key storage).

## Architecture Decisions

- **Rate table drives everything.** Rates stored in `tax_rates` (TaxRate model exists: tax_type, rate_code, rate_value, brackets, effective_from/to, is_active, applicable_to). New service helper `resolveRate(taxType, applicableTo, onDate)` reads via `TaxRepository.GetRates`. No new migration (table exists — verify 000009_tax_schema).
- **VAT engine reads explicit VAT accounts first** (Cr 33311 output, Dr 1331/1332 input), falls back to rate × revenue for accounts without explicit VAT line. Replaces hardcoded lists.
- **CIT rate by company size** — rate table keys: `STANDARD` 20%, `SMALL` 17% (3-50B), `MICRO` 15% (<3B). Company revenue input passed in (no company repo dep in engine).
- **PIT signature changes** — current `CalculatePIT(ctx, companyID, period, employeeIDs)` carries no salary data. New input struct `PITEmployeeInput` (gross, dependants, residence, months). Breaking change: interface + handler route body. Accepted — stub unusable as-is.
- **Declaration engine pulls posted journals** — taxService gains journalRepo + periodRepo deps (constructor change, wire in main.go + tests). Generate → compute → map to form lines → validate cross-sums → DRAFT(VALIDATED).
- **GDT client isolated** — new package `internal/gdtclient` (or `internal/gdt`): interface `GDTClient`, mock HTTP server in tests, retry 1s/5s/30s, timeout 120s. XMLDSig separate package `internal/xmldsig` (C14N + RSA-SHA256). Wiring into IssueEInvoice = Round B.
- **No new migrations Round A** — all on existing tables.

## Rounds (each = several slices, checkpoint, commit)

### Round A: Calculation + Declaration core (THIS ROUND)
| # | Slice | Scope | Deps |
|---|-------|-------|------|
| A1 | Rate resolver — `resolveRate` helper + tests | S | — |
| A2 | VAT engine rewrite (rate-table driven, 33311/1331 first) | M | A1 |
| A3 | CIT engine rewrite (size-based rate, provisional/80% rule) | M | A1 |
| A4 | PIT engine (brackets, deductions, insurances) + signature change | M | A1 |
| A5 | Declaration engine — GenerateDeclaration from posted journals + form-line mapping + cross-validation | L | A2, A3 |
| A6 | Declaration→payment automation (create TaxPayment, due date, late flag) | M | A5 |

### Round B: GDT infrastructure (next)
| # | Slice | Scope | Deps |
|---|-------|-------|------|
| B1 | `internal/xmldsig` — C14N + RSA-SHA256 sign/verify | S | — |
| B2 | `internal/gdt` client — submit/acknowledge, retry, error map, mock server | M | — |
| B3 | Wire IssueEInvoice → sign → GDT submit → status transitions | L | B1, B2 |

### Round C: Tax form XML (later)
GTGT-01/TNDN-01/KK-TNCN per Circular 80 + HTKK wrapper. Depends on A5.

## Checkpoints

- A1-A2: `go vet ./... && go test -count=1 ./...` green, commit each
- A3-A4: same + commit
- A5-A6: same + commit — Round A done
- B: same + commit — Round B done
- Final: docs + AGENTS.md + todo.md updated

## Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| PIT signature change breaks handler | Med | Update handler + tests in same slice A4 |
| NewService constructor growth (journalRepo, periodRepo) | Med | Add 2 params; update main.go ×2 + tests |
| GDT XML/SOAP format drift | High | Mock server + isolate format in gdt package |
| VAT account classification edge cases | Med | Explicit-account-first, fallback documented, tests per rule |
| Rate table empty → calc fails | Med | Seed defaults in memory repo tests; PG: seed migration (Round A checkpoint verifies) |
