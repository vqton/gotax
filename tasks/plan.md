# Payroll Module: Remaining Features Plan

## Overview
Complete payroll module from 74% → 90%+ by implementing remaining features in order of complexity and user value.

## Current State
- Core calculation engine: PROD (net-to-gross, 13th-month, severance, SI/HI/UI, PIT)
- GL auto-posting: PROD (3 balanced journal entries on approval)
- Multi-level approval: PROD (DRAFT→PROCESSING→REVIEWING→APPROVED)
- Competitor parity: 74% (vs MISA 100%, Bravo 88%)

## Task List

### Phase 1: Payslip PDF (Medium complexity, High value)
- [ ] Task 1: Payslip PDF domain model + maroto template
- [ ] Task 2: PDF generation service method
- [ ] Task 3: Handler endpoint + tests
- [ ] Task 4: Email distribution (optional)

### Checkpoint 1: Payslip PDF
- [ ] All tests pass
- [ ] PDF generates with correct data

### Phase 2: Declaration XML (Complex, Medium value)
- [ ] Task 5: D02-TS (SI declaration) XML generation
- [ ] Task 6: 05/KK-TNCN (PIT declaration) XML generation
- [ ] Task 7: TK3-TS (employer registration) XML generation
- [ ] Task 8: Handler endpoints + tests

### Checkpoint 2: Declarations
- [ ] XML validates against form specs
- [ ] Handler tests pass

### Phase 3: Retroactive Pay (Medium complexity)
- [ ] Task 9: Retroactive pay calculation (back-pay for salary changes)
- [ ] Task 10: Service + handler + tests

### Checkpoint 3: Retroactive
- [ ] Correct calculation for mid-month changes
- [ ] Integration with existing payroll run

### Phase 4: Tax Module Link (Coordination)
- [ ] Task 11: KK_TNCN form generation from payroll data
- [ ] Task 12: Service integration with tax module

### Phase 5: Final Review
- [ ] Task 13: code-review-and-quality review
- [ ] Task 14: Update competitor review doc

## Risks
| Risk | Impact | Mitigation |
|------|--------|------------|
| maroto/v2 API changes | Medium | Check existing usage in opening_balance_report.go |
| Declaration XML format changes | Low | XML already generated in internal/payroll/declarations.go |
| Tax module integration complexity | Medium | Keep interface simple, delegate to tax module |
