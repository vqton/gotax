# Implementation Plan: Tax Module → 100%

## Overview

GoTax tax module is ~60% complete. Core engines (VAT/CIT/PIT), declaration engine (GTGT01/TNDN03), GDT pipeline, e-invoice pipeline, and 41 routes exist with 65+ tests. Remaining 40% covers: 13 missing declaration types, CIT advanced logic, e-invoice workflow completion, calendar/alert automation, payment GL posting, and advanced tax modules (FCT, TP, GMT).

## Current State

| Layer | Status | Coverage |
|-------|--------|----------|
| Domain models | DONE | 12 models, 16 enums, 30+ errors |
| Repository | DONE | 26-method interface, PG + Memory |
| Service | DONE | 31 methods, 1245 lines |
| Handler | DONE | 41 routes, 675 lines |
| HTKK XML | PARTIAL | GTGT01 + TNDN03 only |
| GDT Client | DONE | Invoice + declaration submit/status |
| E-Invoice | PARTIAL | TXML + signing, no adjustment workflow |
| Tests | DONE | 65+ tests |
| Migration | DONE | 9 tables |

## Architecture Decisions

1. **Extend existing GenerateDeclaration** — add cases to the switch, not new functions
2. **One declaration type per task** — each form gets its own HTKK XML mapping + validation rules + tests
3. **Vertical slicing** — each task: service method + handler route + HTKK XML + tests
4. **Reuse rate resolver** — all new forms use existing `resolveRate` infrastructure
5. **Reuse payment automation** — `createPaymentForDeclaration` works for any payable form

## Dependency Graph

```
Phase 1: Declaration Types (forms foundation)
    │
    ├── Phase 2: CIT Advanced (provisional, incentives, loss c/f)
    │       │
    │       └── Phase 3: E-Invoice Completion (adjustment, GL posting)
    │
    ├── Phase 4: Calendar & Automation (auto-gen, alerts, late interest)
    │
    └── Phase 5: Advanced Tax (FCT, TP, GMT, consolidation, audit)
            │
            └── Phase 6: Reports & PDF + Phase 7: GDT Production
```

## Task List

### Phase 1: Core Declaration Types (13 tasks)

#### 1.1 VAT Forms

- [ ] **Task 1: GTGT03 quarterly VAT declaration**
  - Extend `GenerateDeclaration` for `DeclTypeGTGT03`
  - HTKK XML: `formCode` mapping, quarterly period key
  - Validation: cross-field rules for quarterly
  - Tests: handler + service

- [ ] **Task 2: GTGT02 monthly VAT (small enterprise)**
  - Extend for `DeclTypeGTGT02`
  - Uses direct method (revenue × rate, no input deduction)
  - HTKK XML mapping
  - Tests

- [ ] **Task 3: GTGT04 per-occurrence VAT**
  - Extend for `DeclTypeGTGT04`
  - For non-continuous business (construction, etc.)
  - HTKK XML mapping
  - Tests

- [ ] **Task 4: GTGT05 non-VAT paysubmit**
  - Extend for `DeclTypeGTGT05`
  - For entities not subject to VAT
  - HTKK XML mapping
  - Tests

#### 1.2 CIT Forms

- [ ] **Task 5: TNDN02 CIT quarterly provisional**
  - Extend for `DeclTypeTNDN02`
  - Q1-Q3 provisional: taxable income × rate, ≥80% of annual estimate
  - HTKK XML mapping
  - Tests

- [ ] **Task 6: TNDN04 CIT annual finalization**
  - Extend for `DeclTypeTNDN04`
  - Full annual CIT with loss carry-forward, incentives
  - HTKK XML mapping
  - Tests

- [ ] **Task 7: TNDN05 CIT restructuring**
  - Extend for `DeclTypeTNDN05`
  - For enterprise mergers/splits
  - HTKK XML mapping
  - Tests

- [ ] **Task 8: TNDN06 CIT petroleum**
  - Extend for `DeclTypeTNDN06`
  - Petroleum-specific CIT rules
  - HTKK XML mapping
  - Tests

#### 1.3 PIT & Other Forms

- [ ] **Task 9: KK_TNCN PIT withholding declaration**
  - Extend for `DeclTypeKKTNCN`
  - Employer PIT withholding for employees
  - Input: payroll data (gross, deductions, dependants)
  - HTKK XML mapping
  - Tests

- [ ] **Task 10: QTT_TNCN PIT quarterly declaration**
  - Extend for `DeclTypeQTTTNCN`
  - Quarterly PIT finalization
  - HTKK XML mapping
  - Tests

- [ ] **Task 11: TTDB01 resource tax**
  - Extend for `DeclTypeTTDB01`
  - Oil/gas/mineral resource tax
  - Rate: volume or value-based
  - HTKK XML mapping
  - Tests

- [ ] **Task 12: BVMT01 environmental protection tax**
  - Extend for `DeclTypeBVMT01`
  - Fuel, chemicals, batteries, etc.
  - Rate: per-unit (VND/liter, VND/kg)
  - HTKK XML mapping
  - Tests

- [ ] **Task 13: NTNN01-03 non-resident income tax**
  - Extend for `DeclTypeNTNN01`, `NTNN02`, `NTNN03`
  - Foreign contractor withholding
  - Flat 20% on gross (or treaty rate)
  - HTKK XML mapping
  - Tests

### Checkpoint: Phase 1
- [ ] `go test -count=1 ./internal/service/ ./internal/handler/`
- [ ] All 17 declaration types generate valid XML
- [ ] Each form has handler + service tests
- [ ] HTKK XML validates against BK:BoKe structure

---

### Phase 2: CIT Advanced Logic (4 tasks)

- [ ] **Task 14: CIT incentive reduction engine**
  - Tax holidays (4+9 years for new sectors)
  - Rate reductions (10%/15% for eligible projects)
  - Per CIT Law Art. 15-18, Decree 320/2025
  - Store incentive config on TaxRate (ApplicableTo or metadata)
  - Tests

- [ ] **Task 15: CIT loss carry-forward**
  - Max 5 years from loss year
  - Track cumulative losses in separate model or field
  - Apply during annual finalization (TNDN03/TNDN04)
  - Tests

- [ ] **Task 16: CIT non-deductible itemization**
  - Fine/penalty: non-cash >5M penalties (CIT-05)
  - Sponsorship (unless specific deductible)
  - Interest deduction cap (30% EBITDA — thin cap)
  - R&D super-deduction (150%)
  - Tests

- [ ] **Task 17: CIT quarterly provisional payments**
  - Q1-Q3: estimated annual CIT × quarter proportion
  - Must be ≥80% of final (CIT-10)
  - Auto-generate TNDN02 from quarterly entries
  - Tests

### Checkpoint: Phase 2
- [ ] CIT engine handles incentives, loss c/f, thin cap
- [ ] Quarterly provisional ≥80% validation
- [ ] All CIT tests pass

---

### Phase 3: E-Invoice Completion (4 tasks)

- [ ] **Task 18: Adjustment e-invoice workflow**
  - Create linked invoice (type=ADJUSTMENT, OriginalInvoiceID set)
  - Validate: original must be ISSUED
  - Sign → submit → status → ISSUED
  - Update original invoice status to REPLACED
  - Tests

- [ ] **Task 19: Replacement e-invoice workflow**
  - Similar to adjustment but replaces original
  - Original status → REPLACED
  - Tests

- [ ] **Task 20: Cancellation note e-invoice**
  - CancellationNote type
  - Submit to GDT → cancel original
  - Tests

- [ ] **Task 21: E-invoice → journal entry auto-posting**
  - On ISSUED: create journal entry (DR 131, CR 5111, CR 33311)
  - Use existing GL posting infrastructure
  - Link via JournalEntryID
  - Tests

### Checkpoint: Phase 3
- [ ] Full e-invoice lifecycle: create → issue → adjust/replace/cancel → GL posting
- [ ] All e-invoice tests pass

---

### Phase 4: Calendar & Automation (4 tasks)

- [ ] **Task 22: Tax calendar auto-generation**
  - `GenerateCalendar(companyID, year)` — creates entries for all tax types
  - VAT: monthly (20th next month) or quarterly (30th)
  - CIT: quarterly (30th) + annual (31-Mar)
  - PIT: monthly
  - Tests

- [ ] **Task 23: Tax alert auto-generation**
  - On calendar creation or daily cron
  - Alert types: DUE_TODAY (7 days before), WARNING (14 days), CRITICAL (3 days)
  - Channel: IN_APP (email/SMS future)
  - Tests

- [ ] **Task 24: Payment late interest calculation**
  - Daily cron or on-demand
  - Late interest = underpaid × 0.0003 × late_days
  - Update TaxPayment.LateInterest
  - Tests

- [ ] **Task 25: Payment → GL auto-posting**
  - On payment record: DR 3331, CR 111/112
  - Use existing journal posting infrastructure
  - Tests

### Checkpoint: Phase 4
- [ ] Calendar auto-generates for full year
- [ ] Alerts fire before deadlines
- [ ] Late interest calculates correctly
- [ ] Payments post to GL

---

### Phase 5: Advanced Tax (5 tasks)

- [ ] **Task 26: FCT (Foreign Contractor Tax)**
  - Domain model: FCTDeclaration
  - Withholding on payments to foreign contractors
  - Rate: 1% (goods) or 5% (services) on gross
  - Declaration + payment
  - HTKK XML
  - Tests

- [ ] **Task 27: Transfer pricing (Decree 255/2026)**
  - Related-party transaction tracking
  - Disclosure report (Form 01/TNDN)
  - Arm's-length principle validation
  - Tests

- [ ] **Task 28: Global Minimum Tax (Pillar 2)**
  - MNE groups >EUR750M revenue
  - Top-up tax calculation
  - Account 82112 posting
  - Tests

- [ ] **Task 29: Tax consolidation**
  - Multi-entity parent-subsidiary
  - Consolidated tax return
  - Intercompany elimination
  - Tests

- [ ] **Task 30: Tax audit workflow**
  - Status machine: OPEN → IN_REVIEW → DOCS_REQUESTED → FINDINGS → CLOSED
  - Document tracking (request/submit/track)
  - Penalty calculation (Decree 310/2025)
  - Tests

### Checkpoint: Phase 5
- [ ] All advanced tax modules have domain + service + handler + tests
- [ ] FCT declaration generates valid XML
- [ ] Transfer pricing disclosure complete

---

### Phase 6: Reports & PDF (2 tasks)

- [ ] **Task 31: Tax declaration PDF rendering**
  - Generate PDF from TaxDeclaration + lines
  - Use existing `maroto/v2` infrastructure
  - Support all form types
  - Tests

- [ ] **Task 32: Tax summary reports**
  - VAT summary by period
  - CIT summary by year
  - Payment history
  - Outstanding balance
  - Tests

### Checkpoint: Phase 6
- [ ] PDF renders for all declaration types
- [ ] Summary reports return correct data

---

### Phase 7: GDT Production (2 tasks)

- [ ] **Task 33: GDT real endpoint integration**
  - Replace stub URLs with real GDT endpoints
  - Certificate-based mTLS auth (replace bearer token)
  - Error handling for real GDT responses
  - Integration tests against sandbox

- [ ] **Task 34: GDT sandbox testing**
  - Set up GDT sandbox environment
  - Test full submission flow
  - Validate XML against real GDT schemas

### Checkpoint: Phase 7
- [ ] GDT submission works against sandbox
- [ ] mTLS auth established
- [ ] Full flow: declaration → XML → sign → submit → ack

---

## Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| GDT sandbox access | HIGH | Request early, test with stubs until available |
| HTKK XML schema validation | HIGH | Test against published XSD schemas |
| CIT incentive complexity | MEDIUM | Start with simplest incentive, extend incrementally |
| Transfer pricing rules | MEDIUM | Defer to Phase 5, not blocking core compliance |
| PILar 2 (GMT) | LOW | Only applies to MNEs >EUR750M, most SMEs exempt |
| Tax rate changes | LOW | Rate resolver already supports effective windows |

## Open Questions

1. Should GTGT02 (small enterprise) use direct method or simplified deduction?
2. FCT: does the codebase need a foreign contractor model, or is it a declaration-only feature?
3. Transfer pricing: is the related-party tracking part of Tax or Company module?
4. GDT sandbox: who has credentials? Need to request from GDT if not available.

## Timeline Estimate

| Phase | Duration | Cumulative |
|-------|----------|------------|
| Phase 1: Declaration Types | 2-3 weeks | 2-3 weeks |
| Phase 2: CIT Advanced | 1-2 weeks | 3-5 weeks |
| Phase 3: E-Invoice | 1-2 weeks | 4-7 weeks |
| Phase 4: Calendar & Auto | 1 week | 5-8 weeks |
| Phase 5: Advanced Tax | 2-3 weeks | 7-11 weeks |
| Phase 6: Reports & PDF | 1 week | 8-12 weeks |
| Phase 7: GDT Production | 1 week | 9-13 weeks |

**Total: ~9-13 weeks to 100%**
