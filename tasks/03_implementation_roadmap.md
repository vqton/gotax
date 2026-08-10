# GoTax — Implementation Roadmap & Execution Plan
## Trading Company Module Delivery

**Prepared by:** Lead Business Analyst
**Date:** 2026-08-10
**Scope:** Trading company (mua bán hàng hóa, không sản xuất)
**Target:** MISA SME 2026 feature parity for trading segment

---

## 1. STRATEGIC ASSESSMENT

### 1.1 Current State
- **34 modules implemented** — core GL, tax, e-invoice, bank, purchase, sale, warehouse, payroll, fixed assets, CCDC, costing, contracts, budgets
- **Coverage:** ~56% of MISA SME feature set for trading companies
- **Strengths:** Solid accounting engine, multi-tenant, compliant with Circular 99/2025, e-invoice XML generation
- **Weaknesses:** No invoice lifecycle management, no HTKK integration, no print templates, no data migration tooling

### 1.2 Critical Gaps by Impact

| Gap | Compliance Risk | Operational Impact | Revenue Impact |
|-----|----------------|-------------------|----------------|
| Invoice book management | 🔴 CRITICAL — Nghị định 310/2025 | 🔴 HIGH | 🔴 HIGH |
| HTKK XML export | 🔴 CRITICAL — cannot file taxes | 🔴 HIGH | 🔴 HIGH |
| Prior year data migration | 🔴 CRITICAL — year-end blocker | 🔴 HIGH | 🔴 HIGH |
| Print templates | 🟡 MEDIUM — daily operations | 🔴 HIGH | 🟡 MEDIUM |
| Supplier/customer GDT check | 🟡 MEDIUM — risk management | 🟡 MEDIUM | 🟢 LOW |

### 1.3 Compliance Deadlines (Non-Negotiable)

| Tax Type | Frequency | Deadline | GoTax Status |
|----------|-----------|----------|--------------|
| VAT (01/GTGT) | Monthly | 20th of following month | ⚠️ Can generate data, no HTKK XML |
| CIT (03/TNDN) | Quarterly | Last day of following quarter | ⚠️ Partial — missing 3-tier rates |
| PIT (05/KK-TNCN) | Monthly | 20th of following month | ⚠️ Missing new form structure |
| E-invoice | Continuous | Real-time | ✅ XML generation works |
| Financial statements | Annual | 31st March | ✅ Can generate |

---

## 2. PHASED ROADMAP

### Phase 1: Compliance Foundation (Weeks 1–6)
**Goal:** Enable tax filing and invoice compliance. Without this, company cannot legally operate.

#### Sprint 1.1 — Invoice Books (Weeks 1–2)
**Priority:** 🔴 CRITICAL — blocks all invoice operations

| Task | Description | Effort | Dependencies |
|------|-------------|--------|--------------|
| 1.1.1 | Domain model: `InvoiceBook` (series, from_number, to_number, status, used_count) | 2h | None |
| 1.1.2 | Domain model: `InvoiceNumber` (book_id, number, status, issued_at, invoice_id) | 2h | 1.1.1 |
| 1.1.3 | Repository: PG + Memory for InvoiceBook, InvoiceNumber | 4h | 1.1.1, 1.1.2 |
| 1.1.4 | Service: Create book, allocate number, release number | 4h | 1.1.3 |
| 1.1.5 | Handler: CRUD books, allocate/release endpoints | 3h | 1.1.4 |
| 1.1.6 | Handler: Wire main.go (PG + memory) | 1h | 1.1.5 |
| 1.1.7 | Migration: `invoice_books`, `invoice_numbers` tables | 1h | 1.1.1 |
| 1.1.8 | Tests: Handler + service tests | 4h | 1.1.6 |
| **Total** | | **22h** | |

**Acceptance Criteria:**
- Create invoice book with series and number range
- Allocate next available number from book
- Track which numbers are used/available
- Release unused numbers
- Validate no duplicate numbers across books

#### Sprint 1.2 — Invoice Lifecycle (Week 2–3)
**Priority:** 🔴 CRITICAL — required by Nghị định 310/2025

| Task | Description | Effort | Dependencies |
|------|-------------|--------|--------------|
| 1.2.1 | Extend `CustomerInvoice` model: original_invoice_id, adjustment_reason, replacement_for | 2h | Existing |
| 1.2.2 | Service: Create adjustment invoice (links to original) | 3h | 1.2.1 |
| 1.2.3 | Service: Create replacement invoice (voids original) | 3h | 1.2.1 |
| 1.2.4 | Service: Cancel invoice with reason (Decree 310) | 2h | 1.2.1 |
| 1.2.5 | Handler: POST /adjust, /replace, /cancel endpoints | 3h | 1.2.2–1.2.4 |
| 1.2.6 | Integration: Link invoice number allocation to lifecycle | 2h | 1.1.4, 1.2.2–1.2.4 |
| 1.2.7 | Missing/damaged invoice report handler | 2h | 1.1.1 |
| 1.2.8 | Tests: Full lifecycle test suite | 4h | 1.2.5 |
| **Total** | | **21h** | |

**Acceptance Criteria:**
- Adjustment invoice references original, calculates difference
- Replacement invoice creates new number, voids original
- Cancel records reason, updates status, cannot reverse
- Audit trail for all lifecycle operations

#### Sprint 1.3 — HTKK XML Export (Weeks 3–4)
**Priority:** 🔴 CRITICAL — cannot file taxes without this

| Task | Description | Effort | Dependencies |
|------|-------------|--------|--------------|
| 1.3.1 | Research HTKK 5.6.0 XML schema (01/GTGT, 03/TNDN, 05/KK-TNCN) | 4h | None |
| 1.3.2 | XML templates for VAT declaration (01/GTGT) | 6h | 1.3.1 |
| 1.3.3 | XML templates for CIT finalization (03/TNDN) — 3-tier rates | 6h | 1.3.1 |
| 1.3.4 | XML templates for PIT return (05/KK-TNCN) — item 25.1 | 4h | 1.3.1 |
| 1.3.5 | Service: Auto-compile GL data into declaration data structures | 6h | 1.3.2–1.3.4 |
| 1.3.6 | Handler: POST /export-xml endpoints per declaration type | 3h | 1.3.5 |
| 1.3.7 | Handler: Tax penalty calculator (late filing/payment) | 4h | 1.3.5 |
| 1.3.8 | Tests: XML output validation against HTKK samples | 6h | 1.3.6 |
| **Total** | | **39h** | |

**Acceptance Criteria:**
- Generate valid HTKK 5.6.0 XML for 01/GTGT, 03/TNDN, 05/KK-TNCN
- Auto-populate from GL journal entries
- Correctly apply 3-tier CIT rates (15%/17%/20%)
- Include item 25.1 exempt income in PIT
- Penalty calculation per Circular 10/2022/TT-BTC

#### Sprint 1.4 — Prior Year Migration (Weeks 4–5)
**Priority:** 🔴 CRITICAL — year-end blocker

| Task | Description | Effort | Dependencies |
|------|-------------|--------|--------------|
| 1.4.1 | Service: Export year-end balances (all account types) | 4h | Existing GL |
| 1.4.2 | Service: Import opening balances for new year | 4h | 1.4.1 |
| 1.4.3 | Service: Auto-carry CCDC, FA, prepaid, deferred revenue | 6h | 1.4.2 |
| 1.4.4 | Service: TT200 → TT99 account mapping migration | 8h | Existing mappings |
| 1.4.5 | Handler: POST /migrate-year endpoint | 2h | 1.4.3, 1.4.4 |
| 1.4.6 | Tests: Round-trip migration validation | 4h | 1.4.5 |
| **Total** | | **28h** | |

**Acceptance Criteria:**
- One-click year-end close creates new fiscal year with balances
- CCDC, FA, prepaid expenses, deferred revenue carried forward
- TT200 accounts mapped to TT99 equivalents
- Opening balance report generated automatically

#### Sprint 1.5 — Print Templates (Weeks 5–6)
**Priority:** 🟡 MEDIUM — daily operations

| Task | Description | Effort | Dependencies |
|------|-------------|--------|--------------|
| 1.5.1 | PDF templates: Phiếu thu, Phiếu chi (TT99 format) | 6h | None |
| 1.5.2 | PDF templates: Sổ nhật ký thu tiền, Sổ chi tiết quỹ | 6h | None |
| 1.5.3 | PDF templates: Bank UNC (ACB, VietABank formats) | 4h | None |
| 1.5.4 | PDF templates: Financial statements (TT99 format) | 6h | Existing reports |
| 1.5.5 | Handler: GET /print/:type/:id endpoints | 3h | 1.5.1–1.5.4 |
| 1.5.6 | Tests: PDF output validation | 3h | 1.5.5 |
| **Total** | | **28h** | |

**Phase 1 Total: ~138h (~3.5 weeks with 1 developer, ~2 weeks with 2 developers)**

---

### Phase 2: Operational Efficiency (Weeks 7–12)
**Goal:** Reduce manual data entry, improve accuracy. Features that save time per transaction.

#### Sprint 2.1 — GDT Tax Code Validation (Week 7)
**Priority:** 🟡 HIGH VALUE — simple API, big ROI

| Task | Description | Effort | Dependencies |
|------|-------------|--------|--------------|
| 2.1.1 | GDT API client: lookup tax code (active/inactive) | 4h | Existing GDT client |
| 2.1.2 | Service: Validate supplier/customer tax code | 2h | 2.1.1 |
| 2.1.3 | Handler: GET /validate-tax-code endpoint | 1h | 2.1.2 |
| 2.1.4 | Integration: Warn on purchase/sales if supplier/buyer inactive | 2h | 2.1.2 |
| 2.1.5 | Tests | 2h | 2.1.3 |
| **Total** | | **11h** | |

#### Sprint 2.2 — Warehouse Automation (Weeks 7–9)
**Priority:** 🟡 MEDIUM — reduces double entry

| Task | Description | Effort | Dependencies |
|------|-------------|--------|--------------|
| 2.2.1 | Service: Auto-create GRN from supplier invoice | 4h | Existing purchase |
| 2.2.2 | Service: Auto-create DN from customer invoice | 4h | Existing sale |
| 2.2.3 | Domain: Lot tracking fields (lot_number, expiry_date) on items | 3h | Existing warehouse |
| 2.2.4 | Domain: Min/max stock levels on warehouse items | 2h | 2.2.3 |
| 2.2.5 | Service: Stock level warning check on transaction | 3h | 2.2.4 |
| 2.2.6 | Handler: Stock warning endpoint | 1h | 2.2.5 |
| 2.2.7 | Migration: Add columns to existing tables | 1h | 2.2.3, 2.2.4 |
| 2.2.8 | Tests | 4h | 2.2.6 |
| **Total** | | **22h** | |

#### Sprint 2.3 — Payroll Updates (Weeks 8–9)
**Priority:** 🟡 MEDIUM — compliance + efficiency

| Task | Description | Effort | Dependencies |
|------|-------------|--------|--------------|
| 2.3.1 | Update PIT progressive table (Law 109/2025) | 2h | Existing payroll |
| 2.3.2 | Update personal deduction levels (2026) | 1h | 2.3.1 |
| 2.3.3 | Service: Salary allocation by department | 3h | Existing payroll |
| 2.3.4 | Service: SI/HI/UI declaration XML generation | 6h | 2.3.1 |
| 2.3.5 | Handler: POST /declarations endpoint | 2h | 2.3.4 |
| 2.3.6 | Tests | 3h | 2.3.5 |
| **Total** | | **17h** | |

#### Sprint 2.4 — Bank Real-time + Auto-Match (Weeks 9–11)
**Priority:** 🟡 MEDIUM — high frequency use

| Task | Description | Effort | Dependencies |
|------|-------------|--------|--------------|
| 2.4.1 | Bank API client interface (BIDV, VCB, ACB adapters) | 6h | None |
| 2.4.2 | Service: Fetch transactions via API | 4h | 2.4.1 |
| 2.4.3 | Service: Auto-match (amount + reference fuzzy) | 6h | 2.4.2 |
| 2.4.4 | Handler: GET /sync, POST /match endpoints | 3h | 2.4.3 |
| 2.4.5 | Tests: Mock bank API responses | 4h | 2.4.4 |
| **Total** | | **23h** | |

#### Sprint 2.5 — Purchase/Sales Enhancements (Weeks 10–11)
**Priority:** 🟡 MEDIUM

| Task | Description | Effort | Dependencies |
|------|-------------|--------|--------------|
| 2.5.1 | Service: Import purchase data from e-invoice XML | 4h | Existing einvoice |
| 2.5.2 | Service: Price list management + auto-calculate selling price | 4h | Existing sale |
| 2.5.3 | Handler: POST /import-einvoice, GET /price-lists | 3h | 2.5.1, 2.5.2 |
| 2.5.4 | Tests | 3h | 2.5.3 |
| **Total** | | **14h** | |

#### Sprint 2.6 — Financial Analysis (Week 11–12)
**Priority:** 🟡 LOW-MEDIUM

| Task | Description | Effort | Dependencies |
|------|-------------|--------|--------------|
| 2.6.1 | Service: Financial ratio calculation (10+ ratios) | 4h | Existing reports |
| 2.6.2 | Service: Budget vs Actual comparison | 3h | Existing budgets |
| 2.6.3 | Handler: GET /ratios, GET /budget-vs-actual | 2h | 2.6.1, 2.6.2 |
| 2.6.4 | Frontend: Financial analysis dashboard | 4h | 2.6.3 |
| 2.6.5 | Tests | 2h | 2.6.3 |
| **Total** | | **15h** | |

**Phase 2 Total: ~102h (~2.5 weeks with 2 developers)**

---

### Phase 3: Infrastructure & Polish (Weeks 13–16)
**Goal:** Product polish, user experience, operational convenience.

#### Sprint 3.1 — Excel Import/Export (Week 13)
| Task | Description | Effort |
|------|-------------|--------|
| 3.1.1 | Excel library setup (excelize) | 2h |
| 3.1.2 | Import: Suppliers, Customers, Items from Excel | 6h |
| 3.1.3 | Export: Journal entries, trial balance, BCTC to Excel | 4h |
| 3.1.4 | Tests | 3h |
| **Total** | | **15h** |

#### Sprint 3.2 — Batch Operations (Week 13–14)
| Task | Description | Effort |
|------|-------------|--------|
| 3.2.1 | Service: Batch approve/post/cancel for vouchers | 4h |
| 3.2.2 | Handler: POST /batch-approve, /batch-post, /batch-cancel | 2h |
| 3.2.3 | Frontend: Multi-select checkbox UI | 3h |
| 3.2.4 | Tests | 2h |
| **Total** | | **11h** |

#### Sprint 3.3 — Document Attachment (Week 14)
| Task | Description | Effort |
|------|-------------|--------|
| 3.3.1 | Domain: Attachment model (entity_type, entity_id, file_path, mime_type) | 2h |
| 3.3.2 | Service: Upload, download, list attachments | 3h |
| 3.3.3 | Handler: POST /upload, GET /download endpoints | 2h |
| 3.3.4 | Migration: attachments table | 1h |
| 3.3.5 | Tests | 2h |
| **Total** | | **10h** |

#### Sprint 3.4 — Audit Trail (Week 14–15)
| Task | Description | Effort |
|------|-------------|--------|
| 3.4.1 | Domain: AuditChange model (entity, field, old_value, new_value, user, timestamp) | 2h |
| 3.4.2 | Middleware: Auto-capture changes on PUT/PATCH | 4h |
| 3.4.3 | Handler: GET /audit-trail/:entity/:id | 2h |
| 3.4.4 | Tests | 2h |
| **Total** | | **10h** |

#### Sprint 3.5 — Period-End Automation (Week 15–16)
| Task | Description | Effort |
|------|-------------|--------|
| 3.5.1 | Service: One-click period-end (depreciation + allocation + FX + tax) | 6h |
| 3.5.2 | Handler: POST /period-end/:period_id | 2h |
| 3.5.3 | Frontend: Period-end wizard UI | 3h |
| 3.5.4 | Tests | 3h |
| **Total** | | **14h** |

#### Sprint 3.6 — Multi-Currency Revaluation (Week 16)
| Task | Description | Effort |
|------|-------------|--------|
| 3.6.1 | Service: FX revaluation for all foreign currency accounts | 4h |
| 3.6.2 | Service: Auto-post revaluation journal entries | 3h |
| 3.6.3 | Tests | 2h |
| **Total** | | **9h** |

**Phase 3 Total: ~69h (~1.7 weeks with 2 developers)**

---

## 3. EXECUTION TIMELINE

```
Week  1  2  3  4  5  6  7  8  9  10  11  12  13  14  15  16
      ├──┴──┤
      │ 1.1  │ Invoice Books
      │  ├───┴───┤
      │  │  1.2   │ Invoice Lifecycle
      │  │  ├─────┴─────┤
      │  │  │    1.3     │ HTKK XML
      │  │  │  ├────┴────┤
      │  │  │  │   1.4    │ Data Migration
      │  │  │  │  ├───┴───┤
      │  │  │  │  │  1.5  │ Print Templates
      ├──┴──┴──┴──┴──┴────┤  PHASE 1 COMPLETE
                          ├──────┴──┤
                          │  2.1     │ GDT Validation
                          │  ├───┴───┴───┤
                          │  │    2.2     │ Warehouse
                          │  │  ├────┴───┤
                          │  │  │  2.3    │ Payroll
                          │  │  │  ├─────┴─────┤
                          │  │  │  │    2.4     │ Bank API
                          │  │  │  │  ├───┴───┤
                          │  │  │  │  │ 2.5   │ Purchase/Sales
                          │  │  │  │  │  ├───┴───┤
                          │  │  │  │  │  │  2.6  │ Analysis
                          ├──┴──┴──┴──┴──┴──┴────┤  PHASE 2 COMPLETE
                                                 ├───┴──┤
                                                 │ 3.1  │ Excel
                                                 │ ├─┴──┤
                                                 │ │ 3.2│ Batch
                                                 │ │ ├─┴─┤
                                                 │ │ │3.3│ Attach
                                                 │ │ │ ├─┴──┤
                                                 │ │ │ │ 3.4 │ Audit
                                                 │ │ │ │ ├──┴──┤
                                                 │ │ │ │ │  3.5 │ Period-End
                                                 │ │ │ │ │  ├─┴─┤
                                                 │ │ │ │ │  │3.6│ FX
                                                 ├──┴──┴──┴──┴──┴──┤  PHASE 3 COMPLETE
```

---

## 4. RESOURCE PLAN

### 4.1 Team Composition

| Role | Count | Responsibility |
|------|-------|----------------|
| Backend Developer (Go) | 2 | Domain, service, handler, repository, tests |
| Frontend Developer | 1 | HTML/Alpine.js pages, PDF templates |
| BA/QA (part-time) | 0.5 | Requirements, test cases, UAT |

### 4.2 Effort Summary

| Phase | Backend (h) | Frontend (h) | Total (h) | Duration |
|-------|-------------|--------------|-----------|----------|
| Phase 1: Compliance | 120 | 18 | 138 | Weeks 1–6 |
| Phase 2: Operations | 85 | 17 | 102 | Weeks 7–12 |
| Phase 3: Polish | 50 | 19 | 69 | Weeks 13–16 |
| **TOTAL** | **255** | **54** | **309** | **16 weeks** |

### 4.3 Critical Path

```
Invoice Books → Invoice Lifecycle → E-Invoice Integration
     ↓
HTKK XML Export → VAT/CIT/PIT Reports
     ↓
Prior Year Migration → Year-End Close
     ↓
Print Templates → Daily Operations
```

**Any delay on the critical path delays go-live.**

---

## 5. RISK REGISTER

| # | Risk | Probability | Impact | Mitigation |
|---|------|-------------|--------|------------|
| R1 | HTKK XML schema changes mid-development | Medium | High | Lock HTKK 5.6.0 spec early, build versioned XML generator |
| R2 | GDT API downtime during testing | High | Medium | Build mock GDT client, test with sample responses |
| R3 | Decree 310/2025 interpretation ambiguity | Medium | High | Consult tax advisor, build flexible invoice lifecycle |
| R4 | Prior year migration data quality issues | High | High | Build validation report before migration, dry-run mode |
| R5 | PDF template layout mismatches MISA | Medium | Low | Use MISA samples as reference, iterate with accountant |
| R6 | Two-tier CIT rate calculation errors | Medium | High | Unit tests with known values, edge cases (revenue thresholds) |
| R7 | Scope creep from "nice to have" requests | High | Medium | Strict phase gating, defer non-critical items |

---

## 6. MVP DEFINITION

### MVP = Phase 1 Only (Weeks 1–6)
**What the company can do with MVP:**
- ✅ Issue e-invoices with proper book/number management
- ✅ Adjust, replace, cancel invoices per Decree 310
- ✅ File VAT, CIT, PIT declarations via HTKK XML
- ✅ Close fiscal year and migrate opening balances
- ✅ Print all vouchers and reports in TT99 format
- ❌ Cannot auto-import from bank API (manual CSV works)
- ❌ Cannot auto-create GRN/DN from invoices
- ❌ Cannot validate supplier/buyer tax codes in real-time

### MVP Go-Live Criteria
1. All 16 MUST HAVE features implemented and tested
2. HTKK XML validated against 3 sample declarations
3. Invoice lifecycle tested with adjustment, replacement, cancellation scenarios
4. Prior year migration tested with sample data
5. Print templates match MISA SME output
6. `go test -count=1 ./...` passes
7. UAT sign-off from accountant

---

## 7. TESTING STRATEGY

### 7.1 Unit Tests (per sprint)
- Every service method: happy path + error cases
- Every handler: HTTP request → response validation
- Repository: in-memory tests (no DB required)

### 7.2 Integration Tests (per phase)
- Phase 1: Invoice book → number allocation → lifecycle → XML export
- Phase 1: Prior year migration end-to-end
- Phase 2: Bank API → auto-match → GL entry
- Phase 2: E-invoice import → purchase invoice creation

### 7.3 Compliance Tests (Phase 1)
- HTKK XML schema validation (XSD if available)
- Invoice number sequence integrity
- Tax calculation accuracy (VAT, CIT 3-tier, PIT progressive)
- Decree 310 compliance checklist

### 7.4 UAT (End of each phase)
- Accountant runs real scenarios
- Compare output with MISA SME for same data
- Sign-off before phase completion

---

## 8. GO-LIVE CHECKLIST

| # | Item | Owner | Status |
|---|------|-------|--------|
| 1 | All Phase 1 features implemented | Dev | ⬜ |
| 2 | HTKK XML validated with tax authority sample | BA/QA | ⬜ |
| 3 | Invoice lifecycle tested per Decree 310 | QA | ⬜ |
| 4 | Prior year migration dry-run successful | Dev | ✅ |
| 5 | Print templates approved by accountant | BA | ✅ |
| 6 | All tests passing | Dev | ✅ |
| 7 | Security review (auth, RBAC, data isolation) | Dev | ⬜ |
| 8 | Performance baseline (100 concurrent users) | Dev | ⬜ |
| 9 | Deployment script tested | Dev | ⬜ |
| 10 | Rollback plan documented | Dev | ⬜ |

---

## 9. SUCCESS METRICS

| Metric | Target | Measurement |
|--------|--------|-------------|
| Tax filing time | < 30 min (vs 2+ hours manual) | Time from "generate" to "filed" |
| Invoice processing time | < 2 min per invoice | Time from "create" to "issued" |
| Data migration time | < 1 hour for full year | Time from "start migration" to "verified" |
| Error rate | < 0.1% on tax calculations | Defects per 1000 transactions |
| User satisfaction | > 4.0/5.0 | Post-go-live survey |

---

## 10. NEXT STEPS

1. **Immediate (this week):**
   - Review this plan with development team
   - Lock HTKK 5.6.0 XML schema reference
   - Set up project board with sprint tasks

2. **Week 1:**
   - Begin Sprint 1.1 (Invoice Books)
   - Start HTKK XML research in parallel

3. **Ongoing:**
   - Weekly sprint reviews with BA
   - Bi-weekly stakeholder demos
   - Phase gate reviews before proceeding

---

*This plan is based on MISA SME 2026 feature set, Vietnamese tax regulations (Circular 99/2025, Decree 123/2020, Decree 310/2025, Law 67/2025, Law 109/2025), and standard trading company operations in Vietnam.*
