# AR Module — Implementation Roadmap

## Overview

Ship AR module to enterprise parity. 12-step lifecycle documented in `docs/sale/AR_WORKFLOW.md`. Current: core sub-ledger CRUD exists (~20% complete). Target: AR aging, GL auto-posting, e-invoice pipeline, customer statement, collection tools.

## Architecture Decisions

- **GL posting via existing JournalEntry service** — reuse service.Service.CreateEntry/PostEntry, not new GL engine
- **E-invoice via pipeline interface** — TXMLGenerator + Signer + GDTSubmitter interfaces, not monolithic
- **Aging by DueDate** — new DB query grouping invoices by age buckets, not in-memory calc
- **Vertical slicing per feature** — each phase delivers working, testable endpoint(s)

## Task List

### Phase 0: GOTAX-AR-001 — Fix AR Aging (Critical Bug)

- [ ] Task 1: Fix AR aging report — bucket by DueDate not total
- [ ] Checkpoint: aging shows correct 30/60/90/120+ buckets

### Phase 1: Core AR (P0)

- [ ] Task 2: GL auto-posting on invoice post (Dr 131 Cr 511/3331)
- [ ] Task 3: GL auto-posting on receipt post (Dr 111/112 Cr 131)
- [ ] Task 4: GL auto-posting on credit note post (reverse entries)
- [ ] Task 5: Customer statement endpoint + report
- [ ] Checkpoint: invoice→GL→statement flow works end-to-end

### Phase 2: Revenue & Compliance (P0)

- [ ] Task 6: E-invoice pipeline interface + TXML generator (Decree 254)
- [ ] Task 7: E-invoice digital signer (XMLDSig via DigitalSignature model)
- [ ] Task 8: GDT API client (submit, cancel, query)
- [ ] Task 9: E-invoice auto-pipeline (hook into PostInvoice)
- [ ] Checkpoint: e-invoice issues end-to-end (mock GDT)

### Phase 3: Collection Management (P1)

- [ ] Task 10: Dunning engine — levels, schedule, auto-email trigger
- [ ] Task 11: Bad debt provision calc (229) + write-off workflow
- [ ] Task 12: Prepayment/deposit workflow (offset on invoice)
- [ ] Task 13: FX revaluation at receipt time + month-end
- [ ] Checkpoint: dunning→bad-debt→prepayment→FX flow works

### Phase 4: Controls & Monitoring (P1/P2)

- [ ] Task 14: Credit limit enforcement (check at SO confirm + invoice)
- [ ] Task 15: AR-GL month-end reconciliation report
- [ ] Task 16: DSO + AR KPI dashboard endpoint
- [ ] Task 17: Off-balance-sheet tracking for written-off AR
- [ ] Checkpoint: all P1 features complete + regression pass

## Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| GDT API integration complexity | High | Mock GDT for dev, contract tests for prod |
| E-invoice XML schema changes (Decree 254) | Medium | Schema-versioned TXML generator |
| Foreign currency AR edge cases | Medium | Test with 3 currencies, reval calc verification |
| Period-end close timing with AR-GL rec | Medium | Block close if variance, allow override |
