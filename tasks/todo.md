# Tax Module → 100% — Task Checklist

## Phase 1: Core Declaration Types

### VAT Forms
- [ ] Task 1: GTGT03 quarterly VAT declaration
- [ ] Task 2: GTGT02 monthly VAT (small enterprise / direct method)
- [ ] Task 3: GTGT04 per-occurrence VAT
- [ ] Task 4: GTGT05 non-VAT paysubmit

### CIT Forms
- [ ] Task 5: TNDN02 CIT quarterly provisional
- [ ] Task 6: TNDN04 CIT annual finalization
- [ ] Task 7: TNDN05 CIT restructuring
- [ ] Task 8: TNDN06 CIT petroleum

### PIT & Other Forms
- [ ] Task 9: KK_TNCN PIT withholding declaration
- [ ] Task 10: QTT_TNCN PIT quarterly declaration
- [ ] Task 11: TTDB01 resource tax
- [ ] Task 12: BVMT01 environmental protection tax
- [ ] Task 13: NTNN01-03 non-resident income tax

### Checkpoint: Phase 1
- [ ] All 17 declaration types generate valid HTKK XML
- [ ] Each form has service + handler tests
- [ ] `go test -count=1 ./internal/service/ ./internal/handler/`

---

## Phase 2: CIT Advanced Logic

- [ ] Task 14: CIT incentive reduction engine
- [ ] Task 15: CIT loss carry-forward (5-year limit)
- [ ] Task 16: CIT non-deductible itemization (thin cap, R&D)
- [ ] Task 17: CIT quarterly provisional payments (≥80% rule)

### Checkpoint: Phase 2
- [ ] CIT handles incentives, loss c/f, thin cap
- [ ] Quarterly provisional ≥80% validation works
- [ ] All CIT tests pass

---

## Phase 3: E-Invoice Completion

- [ ] Task 18: Adjustment e-invoice workflow
- [ ] Task 19: Replacement e-invoice workflow
- [ ] Task 20: Cancellation note e-invoice workflow
- [ ] Task 21: E-invoice → journal entry auto-posting

### Checkpoint: Phase 3
- [ ] Full e-invoice lifecycle: create → issue → adjust/replace/cancel → GL
- [ ] All e-invoice tests pass

---

## Phase 4: Calendar & Automation

- [ ] Task 22: Tax calendar auto-generation (annual, per company)
- [ ] Task 23: Tax alert auto-generation (deadline-based)
- [ ] Task 24: Payment late interest scheduled calculation
- [ ] Task 25: Payment → GL auto-posting

### Checkpoint: Phase 4
- [ ] Calendar auto-generates for full year
- [ ] Alerts fire before deadlines
- [ ] Late interest calculates correctly
- [ ] Payments post to GL

---

## Phase 5: Advanced Tax

- [ ] Task 26: FCT (Foreign Contractor Tax)
- [ ] Task 27: Transfer pricing (Decree 255/2026)
- [ ] Task 28: Global Minimum Tax (Pillar 2)
- [ ] Task 29: Tax consolidation (multi-entity)
- [ ] Task 30: Tax audit workflow

### Checkpoint: Phase 5
- [ ] All advanced tax modules have domain + service + handler + tests
- [ ] FCT declaration generates valid XML

---

## Phase 6: Reports & PDF

- [ ] Task 31: Tax declaration PDF rendering
- [ ] Task 32: Tax summary reports

### Checkpoint: Phase 6
- [ ] PDF renders for all declaration types
- [ ] Summary reports return correct data

---

## Phase 7: GDT Production

- [ ] Task 33: GDT real endpoint integration + mTLS auth
- [ ] Task 34: GDT sandbox testing

### Checkpoint: Phase 7
- [ ] GDT submission works against sandbox
- [ ] Full flow: declaration → XML → sign → submit → ack

---

## Summary

| Phase | Tasks | Est. Duration |
|-------|-------|---------------|
| Phase 1 | 13 | 2-3 weeks |
| Phase 2 | 4 | 1-2 weeks |
| Phase 3 | 4 | 1-2 weeks |
| Phase 4 | 4 | 1 week |
| Phase 5 | 5 | 2-3 weeks |
| Phase 6 | 2 | 1 week |
| Phase 7 | 2 | 1 week |
| **Total** | **34** | **9-13 weeks** |
