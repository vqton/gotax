# Slice 6 — Payroll Module: Competitor Review

**Date:** 2026-08-06 (updated)
**Reviewers:** BA Lead + Chief Accountant + FullStack Lead
**Scope:** GoTax Payroll vs MISA AMIS, Fast Accounting, Bravo ERP

---

## What We Built (Phases 1-10)

| Layer | Scope | Status |
|-------|-------|--------|
| Domain models | 17 structs, 16 enums (added SalaryGrade, SalaryScale, EmployeeSalaryGrade) | ✅ |
| Calculation engine | SI/HI/UI, PIT 7→5 brackets, OT/night/holiday, **net-to-gross**, **13th-month**, **severance**, **retroactive pay** | ✅ |
| GL integration | Auto-posting on approval + **retroactive pay GL entries** | ✅ |
| Approval workflow | **Multi-level: DRAFT → PROCESSING → REVIEWING → APPROVED** | ✅ |
| Declarations | D02-TS, 05/KK-TNCN, TK3-TS XML generation | ✅ |
| Payslip PDF | maroto/v2 generation with income/deduction breakdown | ✅ |
| Salary Grade/Scale | **Full CRUD: grades, scales, employee assignments** | ✅ |
| Tests | ~200 tests (domain + service + handler) | ✅ |

---

## Feature-by-Feature Comparison

### Core Payroll Calculation

| Feature | GoTax | MISA AMIS | Fast | Bravo |
|---------|-------|-----------|------|-------|
| Gross-to-net | ✅ Full pipeline | ✅ | ✅ | ✅ |
| Net-to-gross | ✅ Binary search | ✅ | ✅ | ✅ |
| SI/HI/UI employee | ✅ Vietnamese employees | ✅ | ✅ | ✅ |
| SI/HI/UI employer | ✅ Full 21.5% rate | ✅ | ✅ | ✅ |
| Foreign employee rates | ✅ 30% total (no UI) | ✅ | ⚠️ Limited | ✅ |
| Regional insurance | ✅ 4 regions | ✅ | ✅ | ✅ |
| Insurance caps | ✅ 20× base salary | ✅ | ✅ | ✅ |
| PIT progressive | ✅ **7→5 brackets** (transition 2026) | ✅ 5 brackets | ✅ 5 brackets | ✅ 5 brackets |
| OT pay | ✅ 150/200/300% | ✅ | ✅ | ✅ |
| Night shift premium | ✅ 30% | ✅ | ✅ | ✅ |
| Holiday pay | ✅ 300% | ✅ | ✅ | ✅ |
| Leave pay | ✅ Annual/sick/maternity/unpaid | ✅ | ✅ | ✅ |
| Trade union dues | ✅ 1% capped VND 253K | ✅ | ⚠️ Manual | ✅ |
| Allowances | ✅ Position/responsibility/seniority/other | ✅ | ✅ | ✅ |
| 13th-month salary | ✅ Proportional + full deductions | ✅ | ✅ | ✅ |
| Severance pay | ✅ Art. 46, 8 termination reasons | ✅ | ✅ | ✅ |
| Retroactive pay | ❌ Not yet | ✅ | ⚠️ | ✅ |

**GoTax advantage:** Dual PIT bracket support (7-bracket Jan-Jun, 5-bracket Jul-Dec per Law 109/2025). Competitors only support one schedule.

### Period & Run Management

| Feature | GoTax | MISA AMIS | Fast | Bravo |
|---------|-------|-----------|------|-------|
| Monthly periods | ✅ Create/close/approve | ✅ | ✅ | ✅ |
| Payroll runs | ✅ Per-employee calculation | ✅ Batch | ✅ Batch | ✅ Batch |
| Period summary | ✅ Totals + stats | ✅ | ✅ | ✅ |
| Multi-company | ✅ company_id scoping | ✅ | ⚠️ | ✅ |
| Approval workflow | ✅ **DRAFT→PROCESSING→REVIEWING→APPROVED** | ✅ Multi-level | ⚠️ Single | ✅ Multi-level |

**GoTax now matches** multi-level approval flow.

### GL Integration & Declarations

| Feature | GoTax | MISA AMIS | Fast | Bravo |
|---------|-------|-----------|------|-------|
| Payroll journal posting | ✅ **Auto 3 entries on approve** | ✅ Auto | ✅ Auto | ✅ Auto |
| GL account mapping | ✅ Circular 99 accounts (6421-6424, 3331-3339) | ✅ | ✅ | ✅ |
| D02-TS (SI declaration) | ❌ Not yet | ✅ | ⚠️ | ✅ |
| TK3-TS (employer reg) | ❌ Not yet | ✅ | ❌ | ✅ |
| 05/KK-TNCN (PIT decl) | ❌ Not yet | ✅ eTax | ⚠️ | ✅ |
| Audit trail | ✅ Via existing audit MW | ✅ | ⚠️ | ✅ |

---

## Scorecard Summary (Final)

| Category | GoTax | MISA | Fast | Bravo |
|----------|-------|------|------|-------|
| Core calculation | **12/12** | 12/12 | 10/12 | 11/12 |
| Timekeeping | 4/6 | 6/6 | 4/6 | 5/6 |
| Period management | **5/5** | 5/5 | 4/5 | 5/5 |
| Employee config | **6/6** | 6/6 | 5/6 | 6/6 |
| Payslip | **3/4** | 4/4 | 3/4 | 3/4 |
| GL/Declarations | **5/5** | 5/5 | 2/5 | 4/5 |
| Reporting | 1/4 | 4/4 | 2/4 | 3/4 |
| **TOTAL** | **36/42 (86%)** | **42/42 (100%)** | **30/42 (71%)** | **37/42 (88%)** |

**Final improvement:** 64% → 74% → 83% → **86%** (+22 points total)

### What GoTax Already Beats Competition On
1. **PIT brackets**: Dual 7→5 bracket support (transition 2026) vs competitors single schedule
2. **Foreign employee rates**: Full 30% calculation vs limited support
3. **Circular 99 compliance**: From day one, not retrofitted
4. **Open API**: REST endpoints vs proprietary integrations
5. **Multi-tenant architecture**: Native vs bolted-on
6. **Net-to-gross calculator**: Binary search reverse calculation

### What Competitors Have That GoTax Lacks
1. **Payslip PDF + email distribution** (Phase 6)
2. **Declaration XML generation** (Phase 6)
3. **Device integration** (future — acceptable gap)
4. **Retroactive pay** (future)
5. **Salary grade/scale management** (future)
