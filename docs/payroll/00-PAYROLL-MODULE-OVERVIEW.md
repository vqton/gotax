# Payroll Module — Overview & Production Readiness Assessment

**Date:** 2026-08-04
**Authors:** BA Lead (20+ years) + Chief Accountant (20+ years)
**Status:** NOT PRODUCTION READY — 0% implemented

---

## Executive Summary

GoTax currently has **no payroll module**. The codebase contains only:
- Basic `Employee` model (Company module) — code, name, department, tax code, SI number, bank account, status
- `PITEmployeeInput` struct for tax calculation — gross salary, dependants, months, SI base
- PIT calculation service (tax module) — progressive 5-bracket schedule per Law 109/2025/QH15

**Missing entirely:**
- Salary calculation engine (gross-to-net, net-to-gross)
- Social insurance / health insurance / unemployment insurance calculation
- Timekeeping integration (attendance, overtime, night shift, leave)
- Payslip generation and distribution
- Payroll journal entry auto-posting to GL
- Social insurance declaration forms (D02-TS, TK3-TS, TK1-TS)
- PIT declaration forms (05/KK-TNCN quarterly from Q2/2026)
- Trade union dues calculation
- 13th-month salary / bonus processing
- Severance / termination pay calculation
- Multi-company / multi-branch payroll consolidation
- Payroll approval workflows
- Employee self-service (payslip viewing, leave requests)
- Banking integration for salary disbursement

**Verdict:** Payroll module requires full build. Estimated 12-16 weeks for PROD-ready MVP.

---

## Regulatory Framework (Current as of August 2026)

### Primary Laws & Decrees

| Document | Effective | Key Impact |
|----------|-----------|------------|
| Law 41/2024/QH15 — Social Insurance | 01 Jul 2025 | SI rates, contribution caps, expanded coverage |
| Law 74/2025/QH15 — Employment (Unemployment Insurance) | 01 Jan 2026 | UI coverage expanded, contract duration reduced to 1 month |
| Law 51/2024/QH15 — Health Insurance (amendment) | 01 Jul 2025 | HI rates, expanded subjects |
| Law 109/2025/QH15 — Personal Income Tax | 01 Jul 2026 (employment income from tax year 2026) | 5-bracket progressive, new deductions |
| Decree 161/2026/ND-CP — Base Salary increase | 01 Jul 2026 | Base salary VND 2,530,000 → SI/HI cap VND 50,600,000 |
| Decree 293/2025/ND-CP — Regional minimum wage | 01 Jan 2026 | Region I: VND 5,310,000; UI cap = 20x regional min |
| Decree 253/2026/ND-CP — PIT implementing decree | 01 Jul 2026 | New deductions, overtime/night shift exempt |
| Resolution 110/2025/UBTVQH15 — Deduction levels | 01 Jan 2026 | Personal: VND 15.5M; Dependant: VND 6.2M |
| Resolution 66.16/NQ-CP — Quarterly PIT declaration | 08 May 2026 | Monthly → quarterly PIT filing |
| Decision 1109/QD-BTC — PIT form update | 08 May 2026 | Form 05/KK-TNCN for quarterly filing |
| Decree 158/2025/ND-CP — SI implementing decree | 01 Jul 2025 | Detailed SI contribution rules |
| Decision 366/QĐ-BHXH — BHXH forms 2026 | 01 Jun 2026 | D02-LT, TK3-TS, TK1-TS updated forms |
| Circular 99/2025/TT-BTC — Accounting standards | 01 Jan 2026 | Payroll GL entries, employee benefit accounting |
| IAS 19 — Employee Benefits (IFRS) | Ongoing | Post-employment benefits, termination benefits |

### Contribution Rates (2026)

#### Vietnamese Employees

| Component | Employee | Employer | Total |
|-----------|----------|----------|-------|
| Social Insurance (Retirement & Survivorship) | 8% | 14% | 22% |
| Social Insurance (Sickness & Maternity) | 0% | 3% | 3% |
| Social Insurance (Occupational Accident & Disease) | 0% | 0.5% | 0.5% |
| Health Insurance | 1.5% | 3% | 4.5% |
| Unemployment Insurance | 1% | 1% | 2% |
| **Total** | **10.5%** | **21.5%** | **32%** |

#### Foreign Employees (12+ month contracts)

| Component | Employee | Employer | Total |
|-----------|----------|----------|-------|
| Social Insurance (all components) | 8% | 17.5% | 25.5% |
| Health Insurance | 1.5% | 3% | 4.5% |
| Unemployment Insurance | N/A | N/A | N/A |
| **Total** | **9.5%** | **20.5%** | **30%** |

### Salary Caps (2026)

| Cap Type | Period | Calculation | Amount (VND/month) |
|----------|--------|-------------|---------------------|
| SI/HI contribution base | 01 Jan – 30 Jun 2026 | 20 × reference level (VND 2,340,000) | 46,800,000 |
| SI/HI contribution base | 01 Jul – 31 Dec 2026 | 20 × base salary (VND 2,530,000) | 50,600,000 |
| UI contribution base (Region I) | All year | 20 × VND 5,310,000 | 106,200,000 |
| Trade union dues (members) | All year | 1% × 20 × reference level/base salary | 234,000 / 253,000 |

**Note:** Reference level (mức tham chiếu) = VND 2,340,000 from 01 Jan. Base salary (mức lương cơ sở) increases to VND 2,530,000 from 01 Jul per Decree 161/2026/NĐ-CP.

### PIT Progressive Schedule

**⚠️ TRANSITION PERIOD:** From 01 Jan to 30 Jun 2026, the OLD 7-bracket schedule applies (Law 103/2016/QH13 as amended). From 01 Jul 2026, the NEW 5-bracket schedule applies (Law 109/2025/QH15). Both apply to tax year 2026.

#### New Schedule (From 01 Jul 2026 — Law 109/2025/QH15)

| Bracket | Monthly Taxable Income (VND) | Rate |
|---------|------------------------------|------|
| 1 | Up to 10,000,000 | 5% |
| 2 | Over 10,000,000 to 30,000,000 | 10% |
| 3 | Over 30,000,000 to 60,000,000 | 20% |
| 4 | Over 60,000,000 to 100,000,000 | 30% |
| 5 | Over 100,000,000 | 35% |

#### Old Schedule (01 Jan – 30 Jun 2026 — Law 103/2016/QH13)

| Bracket | Monthly Taxable Income (VND) | Rate |
|---------|------------------------------|------|
| 1 | Up to 5,000,000 | 5% |
| 2 | Over 5,000,000 to 10,000,000 | 10% |
| 3 | Over 10,000,000 to 18,000,000 | 15% |
| 4 | Over 18,000,000 to 32,000,000 | 20% |
| 5 | Over 32,000,000 to 52,000,000 | 25% |
| 6 | Over 52,000,000 to 80,000,000 | 30% |
| 7 | Over 80,000,000 | 35% |

**Deductions:**
- Personal: VND 15,500,000/month (from 01 Jan 2026 — Resolution 110/2025/UBTVQH15)
- Per dependant: VND 6,200,000/month (from 01 Jan 2026)
- Overtime income: FULLY EXEMPT from PIT (from 01 Jul 2026 — Law 109/2025 Art. 4)
- Night shift income: FULLY EXEMPT from PIT (from 01 Jul 2026)
- Unused paid leave: FULLY EXEMPT from PIT (from 01 Jul 2026)

**⚠️ Note:** From 01 Jan to 30 Jun 2026, OT/night/leave income IS subject to PIT under the old law. Only the new deductions (15.5M personal, 6.2M dependent) apply from 01 Jan.

### Overtime Rules (Labour Code 2019)

| Condition | Minimum Rate | Total Pay |
|-----------|-------------|-----------|
| Normal working days | +50% | 150% |
| Weekly rest days | +100% | 200% |
| Public holidays / Tet | +200% | 300% |
| Night work premium | +30% | 30% of hourly rate |
| Night overtime | +20% on top of OT rate | Cumulative |

**Limits:**
- Max 8 hours/day, 48 hours/week (State encourages 40-hour week)
- Max 40 hours overtime/month
- Max 200 hours overtime/year (300 for specific industries)
- Max 12 hours total work/day (ordinary + overtime)

### Public Holidays (2026)

Per Labour Code 2019 Art. 112 + Decree 128/2025/NĐ-CP + Official Letter 9859/VPCP-KGVX:

| # | Holiday | Official Date | Days Off | Notes |
|---|---------|--------------|----------|-------|
| 1 | New Year's Day | 01 Jan | 1 | |
| 2 | Tet Nguyen Dan | 16-20 Feb | 5 | 1 day before Tet + 4 days after. Full break: 14-22 Feb (9 days incl. weekends) |
| 3 | Hung Kings Commemoration | 26 Apr | 1 | 10/3 lunar calendar |
| 4 | Reunification Day | 30 Apr | 1 | |
| 5 | International Labour Day | 01 May | 1 | |
| 6 | National Day | 01-02 Sep | 2 | Aug 31 (Mon) swapped to Aug 22 (Sat). Break: 29 Aug - 02 Sep |

**Total statutory paid holidays: 11 days** (5 Tet + 1 New Year + 1 Hung Kings + 1 Reunification + 1 Labour Day + 2 National Day)

**Note:** Employers may choose alternative Tet/National Day arrangements per Decree 128/2025 Art. 7. Private sector gets 5 Tet days (1 before + 4 after) or alternative split.

---

## Gap Analysis: GoTax vs. Market Leaders

### MISA AMIS HRM (Market Leader, ~45% SMB Vietnam)

| Feature | MISA | GoTax Current |
|---------|------|---------------|
| Timekeeping integration | ✅ Auto-sync from devices | ❌ None |
| Salary calculation | ✅ Flexible formulas, AI-assisted | ❌ None |
| Insurance auto-calc | ✅ Auto, law-updated | ❌ Only PIT input struct |
| Payslip generation | ✅ Mobile + email | ❌ None |
| PIT declaration | ✅ eTax integrated | ⚠️ Basic PIT calc only |
| SI declaration | ✅ D02-TS, TK3-TS | ❌ None |
| GL auto-posting | ✅ Integrated | ❌ None |
| Employee self-service | ✅ Mobile app | ❌ None |
| Approval workflows | ✅ Multi-level | ❌ None |
| Multi-company | ✅ | ⚠️ Company module exists |
| AI features | ✅ MISA AVA assistant | ❌ None |

### BRAVO ERP (Enterprise, ~10% market)

| Feature | BRAVO | GoTax Current |
|---------|-------|---------------|
| Salary parameters | ✅ Per department/position | ❌ None |
| Timekeeping connection | ✅ Device integration | ❌ None |
| Overtime tracking | ✅ Registration + stats | ❌ None |
| Insurance calculation | ✅ Auto, configurable rates | ❌ None |
| GL auto-posting | ✅ Journal entries | ❌ None |
| Minimum wage supplement | ✅ Regional auto-calc | ❌ None |
| Salary scale management | ✅ Grade/coefficient | ❌ None |
| Payslip email | ✅ Per employee | ❌ None |
| PIT declaration | ✅ Form 05/TNCN | ❌ None |

### Fast Accounting (SMB, ~25% market)

| Feature | Fast | GoTax Current |
|---------|------|---------------|
| Payroll processing | ✅ | ❌ None |
| Insurance calculation | ✅ | ❌ None |
| PIT calculation | ✅ | ⚠️ Basic only |
| Journal entry posting | ✅ | ❌ None |
| Excel import/export | ✅ | ❌ None |
| Multi-company | ✅ | ⚠️ Partial |

---

## Production Readiness Checklist

| Requirement | Status | Priority |
|-------------|--------|----------|
| Gross-to-net salary calculation | ❌ Missing | CRITICAL |
| Net-to-gross salary calculation | ❌ Missing | CRITICAL |
| SI/HI/UI contribution calculation | ❌ Missing | CRITICAL |
| PIT progressive calculation (5 brackets) | ⚠️ Basic exists | CRITICAL |
| PIT deduction management (personal + dependants) | ❌ Missing | CRITICAL |
| Overtime pay calculation | ❌ Missing | CRITICAL |
| Night shift pay calculation | ❌ Missing | HIGH |
| Leave types and pay (annual, sick, maternity) | ❌ Missing | HIGH |
| Holiday pay calculation | ❌ Missing | HIGH |
| 13th-month salary | ❌ Missing | HIGH |
| Severance / termination pay | ❌ Missing | HIGH |
| Timekeeping data import | ❌ Missing | CRITICAL |
| Payslip generation (PDF) | ❌ Missing | CRITICAL |
| Payslip email distribution | ❌ Missing | MEDIUM |
| Payroll journal entry to GL | ❌ Missing | CRITICAL |
| D02-TS social insurance declaration | ❌ Missing | CRITICAL |
| TK3-TS employer registration | ❌ Missing | HIGH |
| TK1-TS employee registration | ❌ Missing | HIGH |
| 05/KK-TNCN quarterly PIT declaration | ❌ Missing | CRITICAL |
| Trade union dues calculation | ❌ Missing | MEDIUM |
| Multi-company payroll | ❌ Missing | HIGH |
| Approval workflow | ❌ Missing | HIGH |
| Employee self-service portal | ❌ Missing | MEDIUM |
| Banking integration (salary file) | ❌ Missing | MEDIUM |
| Retroactive pay calculation | ❌ Missing | MEDIUM |
| Arrears calculation | ❌ Missing | MEDIUM |
| Payroll audit trail | ❌ Missing | HIGH |
| Salary scale / grade management | ❌ Missing | MEDIUM |
| Allowance management | ❌ Missing | HIGH |
| Deduction management | ❌ Missing | HIGH |

**Overall Status: 0/30 requirements met. Module requires full implementation.**

---

## Recommended Implementation Phases

### Phase 1: Core Payroll Engine (Weeks 1-4)
- Domain models (SalaryStructure, PayrollRun, Payslip, etc.)
- Gross-to-net calculation engine
- SI/HI/UI contribution calculation
- PIT progressive calculation with new deductions
- Basic payslip generation

### Phase 2: Timekeeping & Attendance (Weeks 5-8)
- Timekeeping data import (CSV/API)
- Overtime calculation
- Night shift calculation
- Leave management (annual, sick, maternity, unpaid)
- Holiday pay

### Phase 3: GL Integration & Declarations (Weeks 9-12)
- Payroll journal entry auto-posting
- D02-TS social insurance declaration generation
- 05/KK-TNCN quarterly PIT declaration
- TK3-TS / TK1-TS forms

### Phase 4: Advanced Features (Weeks 13-16)
- Approval workflows
- Employee self-service
- 13th-month salary / bonus processing
- Severance / termination pay
- Multi-company consolidation
- Banking integration

---

## Competitive Positioning

GoTax payroll module, when complete, will:
1. **Match MISA AMIS** on core payroll + insurance + PIT automation
2. **Match BRAVO** on GL integration + audit trail
3. **Exceed both** on open API architecture + multi-tenant design
4. **Unique advantage**: Integrated with GoTax GL → seamless journal posting, no double entry
5. **Unique advantage**: Circular 99/2025/TT-BTC compliant from day one
6. **Unique advantage**: IFRS/IAS 19 employee benefits ready

**Target market:** SMEs (50-500 employees) who already use GoTax for accounting and want integrated payroll.
