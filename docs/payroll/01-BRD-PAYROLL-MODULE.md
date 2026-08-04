# Business Requirements Document — Payroll Module

**Document ID:** BRD-PAYROLL-001
**Version:** 1.0
**Date:** 2026-08-04
**Classification:** Internal — Confidential

---

## 1. Business Context

### 1.1 Current State
GoTax is a Vietnamese tax-compliant General Ledger API serving SMEs. It covers GL, Auth, Company, Tax, Cash, Bank, Purchase, Sale, Warehouse, and Fixed Assets modules. The Company module includes basic Employee records (code, name, department, tax code, SI number, bank account).

### 1.2 Problem Statement
Vietnamese SMEs must process payroll monthly, calculating:
- Gross-to-net salary with multiple income components
- Mandatory insurance contributions (SI/HI/UI) at legislated rates
- Personal income tax using progressive 5-bracket schedule
- Social insurance declarations (D02-TS, TK3-TS)
- PIT declarations (05/KK-TNCN quarterly)
- Journal entries to GL

Currently, GoTax users must export employee data to external payroll software (MISA, Fast, Excel), manually create journal entries, and reconcile across systems. This creates:
- **Data duplication** — employee data in GoTax + payroll system
- **Error risk** — manual journal entry posting
- **Compliance risk** — delayed law updates in external systems
- **Time cost** — 2-3 days monthly per company for payroll processing

### 1.3 Business Objective
Build an integrated payroll module that:
1. Calculates Vietnamese-compliant payroll automatically
2. Posts payroll journal entries directly to GoTax GL
3. Generates social insurance and PIT declaration forms
4. Updates automatically when Vietnamese tax/insurance laws change
5. Provides employee self-service for payslips

### 1.4 Success Metrics
| Metric | Target |
|--------|--------|
| Payroll processing time | < 30 minutes for 100 employees |
| Calculation accuracy | 100% (zero manual correction) |
| Law update lag | < 7 days from official gazette |
| User adoption | 60% of GoTax companies within 12 months |
| Support tickets | < 5 per 100 companies per month |

---

## 2. Stakeholders

| Role | Responsibility |
|------|----------------|
| Chief Accountant | Payroll policy, compliance, GL integration |
| HR Manager | Employee data, timekeeping, leave policies |
| Payroll Clerk | Monthly payroll processing, declaration filing |
| IT Admin | System configuration, integration setup |
| Employee | Self-service payslip viewing, leave requests |
| Auditor | Payroll audit trail, compliance verification |
| Tax Authority | Declaration forms, PIT/SI reporting |
| Social Insurance Authority | D02-TS, TK3-TS declarations |

---

## 3. Functional Requirements

### 3.1 Employee Master Data (Extend Existing)

| ID | Requirement | Priority | Status |
|----|-------------|----------|--------|
| EMP-01 | Employee contract type (indefinite/definite/probation) | CRITICAL | ❌ Missing |
| EMP-02 | Salary type (time-based/piece-rate/commission/salary coefficient) | CRITICAL | ❌ Missing |
| EMP-03 | Base salary / salary coefficient | CRITICAL | ❌ Missing |
| EMP-04 | Allowances (transport, meal, housing, phone, responsibility) | CRITICAL | ❌ Missing |
| EMP-05 | Bank account for salary disbursement | HIGH | ⚠️ Exists (BankAccountNo) |
| EMP-06 | Tax code (already exists as PersonalTaxCode) | CRITICAL | ✅ Exists |
| EMP-07 | Social insurance number (already exists) | CRITICAL | ✅ Exists |
| EMP-08 | Dependants for PIT deduction | CRITICAL | ❌ Missing |
| EMP-09 | Insurance registration date | HIGH | ❌ Missing |
| EMP-10 | Insurance base salary (separate from contract salary) | CRITICAL | ❌ Missing |
| EMP-11 | Working region (for minimum wage compliance) | CRITICAL | ❌ Missing |
| EMP-12 | Trade union member status | MEDIUM | ❌ Missing |
| EMP-13 | Foreign employee flag (different insurance rates) | HIGH | ❌ Missing |
| EMP-14 | High-tech talent PIT exemption flag | MEDIUM | ❌ Missing |

### 3.2 Salary Structure

| ID | Requirement | Priority |
|----|-------------|----------|
| SAL-01 | Salary components: base salary, allowances, bonuses | CRITICAL |
| SAL-02 | Salary coefficient system (coefficient × base salary) | HIGH |
| SAL-03 | Salary scale / grade management | MEDIUM |
| SAL-04 | Minimum wage supplement by region | HIGH |
| SAL-05 | Overtime rate configuration (per day type) | CRITICAL |
| SAL-06 | Night shift premium configuration | CRITICAL |
| SAL-07 | Holiday pay rules | CRITICAL |
| SAL-08 | 13th-month salary calculation | HIGH |
| SAL-09 | Bonus calculation (performance, KPI-based) | MEDIUM |
| SAL-10 | Retroactive pay adjustment | MEDIUM |

### 3.3 Timekeeping & Attendance

| ID | Requirement | Priority |
|----|-------------|----------|
| TK-01 | Import timekeeping data (CSV, Excel, API) | CRITICAL |
| TK-02 | Standard working hours configuration (8h/day, 40h/week) | CRITICAL |
| TK-03 | Shift management (day/night/rotating) | HIGH |
| TK-04 | Overtime tracking with approval | CRITICAL |
| TK-05 | Night work tracking (22:00-06:00) | HIGH |
| TK-06 | Leave types: annual, sick, maternity, unpaid, compensatory | CRITICAL |
| TK-07 | Leave balance tracking | HIGH |
| TK-08 | Holiday calendar (public holidays, company-specific) | CRITICAL |
| TK-09 | Absent without leave tracking | HIGH |
| TK-10 | Late arrival / early departure penalties | MEDIUM |

### 3.4 Payroll Calculation Engine

| ID | Requirement | Priority |
|----|-------------|----------|
| CALC-01 | Gross-to-net calculation | CRITICAL |
| CALC-02 | Net-to-gross calculation (reverse engineering) | HIGH |
| CALC-03 | Multiple income types: salary, OT, night shift, bonus, commission | CRITICAL |
| CALC-04 | SI contribution: 8% employee (capped at 20× base salary) | CRITICAL |
| CALC-05 | HI contribution: 1.5% employee (capped at 20× base salary) | CRITICAL |
| CALC-06 | UI contribution: 1% employee (capped at 20× regional min wage) | CRITICAL |
| CALC-07 | Employer SI: 17.5% (14% retirement + 3% sickness + 0.5% accident) | CRITICAL |
| CALC-08 | Employer HI: 3% | CRITICAL |
| CALC-09 | Employer UI: 1% (Vietnamese employees only) | CRITICAL |
| CALC-10 | Foreign employee rates (no UI, different total) | HIGH |
| CALC-11 | PIT progressive 5-bracket calculation | CRITICAL |
| CALC-12 | Personal deduction: VND 15,500,000/month | CRITICAL |
| CALC-13 | Dependant deduction: VND 6,200,000/month each | CRITICAL |
| CALC-14 | OT income PIT exemption (from 01 Jul 2026) | CRITICAL |
| CALC-15 | Night shift PIT exemption (from 01 Jul 2026) | CRITICAL |
| CALC-16 | Unused leave PIT exemption (from 01 Jul 2026) | CRITICAL |
| CALC-17 | Supplementary pension insurance deduction (up to VND 3M/month) | MEDIUM |
| CALC-18 | Life insurance premium deduction | MEDIUM |
| CALC-19 | Trade union dues: 1% (members, capped at 253,000 VND/month) | MEDIUM |
| CALC-20 | Salary cap auto-update when base salary changes | CRITICAL |

### 3.5 Payslip & Reporting

| ID | Requirement | Priority |
|----|-------------|----------|
| PAY-01 | Monthly payslip generation (PDF) | CRITICAL |
| PAY-02 | Payslip distribution (email, self-service portal) | HIGH |
| PAY-03 | Payroll summary report (by department, company) | CRITICAL |
| PAY-04 | Salary cost allocation report | HIGH |
| PAY-05 | Insurance contribution summary | CRITICAL |
| PAY-06 | PIT withholding summary | CRITICAL |
| PAY-07 | Overtime summary report | HIGH |
| PAY-08 | Leave balance report | MEDIUM |
| PAY-09 | Year-end tax finalization report | CRITICAL |
| PAY-10 | Payroll comparison (month-over-month) | MEDIUM |

### 3.6 GL Integration

| ID | Requirement | Priority |
|----|-------------|----------|
| GL-01 | Auto-generate payroll journal entries | CRITICAL |
| GL-02 | Journal entry: Dr Salary Expense, Cr Salary Payable | CRITICAL |
| GL-03 | Journal entry: Dr Insurance Expense (employer), Cr Insurance Payable | CRITICAL |
| GL-04 | Journal entry: Cr Employee Insurance Deductions | CRITICAL |
| GL-05 | Journal entry: Cr PIT Withholding Payable | CRITICAL |
| GL-06 | Journal entry: Dr Salary Payable, Cr Bank (payment) | CRITICAL |
| GL-07 | Multi-currency support for expatriate payroll | MEDIUM |
| GL-08 | Departmental cost allocation in journal entries | HIGH |
| GL-09 | Period closing validation (no payroll after period close) | HIGH |

### 3.7 Government Declarations

| ID | Requirement | Priority |
|----|-------------|----------|
| DEC-01 | D02-TS: Employee SI/HI/UI registration list | CRITICAL |
| DEC-02 | TK3-TS: Employer SI/HI registration | HIGH |
| DEC-03 | TK1-TS: Employee SI/HI information | HIGH |
| DEC-04 | 05/KK-TNCN: Quarterly PIT declaration (from Q2/2026) | CRITICAL |
| DEC-05 | Tax withholding certificate (per employee) | HIGH |
| DEC-06 | Annual PIT finalization support | CRITICAL |
| DEC-07 | XML generation for electronic submission | HIGH |
| DEC-08 | I-VAN integration for SI declaration submission | MEDIUM |

### 3.8 Workflow & Approval

| ID | Requirement | Priority |
|----|-------------|----------|
| WF-01 | Payroll preparation → review → approval → payment workflow | CRITICAL |
| WF-02 | Multi-level approval (preparer → reviewer → approver) | HIGH |
| WF-03 | Payroll lock after approval (prevent re-processing) | CRITICAL |
| WF-04 | Ad-hoc adjustment approval | HIGH |
| WF-05 | Retroactive pay approval | MEDIUM |

### 3.9 Employee Self-Service

| ID | Requirement | Priority |
|----|-------------|----------|
| ESS-01 | View payslips (current and historical) | HIGH |
| ESS-02 | Download payslip PDF | HIGH |
| ESS-03 | View leave balances | MEDIUM |
| ESS-04 | Submit leave requests | MEDIUM |
| ESS-05 | Update personal information (bank account, address) | MEDIUM |
| ESS-06 | View insurance contribution history | MEDIUM |

---

## 4. Non-Functional Requirements

| ID | Requirement | Priority |
|----|-------------|----------|
| NFR-01 | Performance: Process 100 employees payroll in < 30 seconds | HIGH |
| NFR-02 | Accuracy: 100% calculation accuracy (zero tolerance) | CRITICAL |
| NFR-03 | Security: Payroll data encrypted at rest and in transit | CRITICAL |
| NFR-04 | Audit: Full audit trail for all payroll changes | CRITICAL |
| NFR-05 | Compliance: Auto-update when laws change | HIGH |
| NFR-06 | Availability: 99.9% uptime during payroll processing windows | HIGH |
| NFR-07 | Backup: Daily payroll data backup | HIGH |
| NFR-08 | Multi-tenant: Company data isolation | CRITICAL |

---

## 5. Integration Points

| System | Integration Type | Priority |
|--------|-----------------|----------|
| GoTax GL | Direct API (same codebase) | CRITICAL |
| GoTax Company (Employee) | Direct API (same codebase) | CRITICAL |
| GoTax Tax (PIT) | Direct API (same codebase) | CRITICAL |
| Timekeeping devices | CSV/Excel import | HIGH |
| MISA AMIS | API (if migrating from MISA) | MEDIUM |
| Banking systems | Salary file generation (CSV/XML) | MEDIUM |
| I-VAN (SI declaration) | API | MEDIUM |
| eTax (PIT declaration) | XML generation | HIGH |
| Email (payslip distribution) | SMTP | MEDIUM |

---

## 6. Assumptions

1. Employee master data exists in GoTax Company module
2. Timekeeping data is imported manually (no real-time device integration in MVP)
3. Vietnamese language is primary; English secondary
4. VND is the primary currency (multi-currency in later phase)
5. Monthly payroll cycle (weekly/bi-weekly in later phase)
6. Single legal entity per company (multi-entity consolidation in later phase)

---

## 7. Dependencies

1. Employee model must be extended with payroll-specific fields
2. GL module must support payroll account codes (3331, 3332, 3333, etc.)
3. Tax module PIT calculation must be updated for new deductions
4. PDF generation capability (maroto/v2 already in codebase)

---

## 8. Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| Law changes during development | HIGH | Parameterized rates, config-driven |
| Complex overtime rules | MEDIUM | Rule engine with configurable rates |
| Multi-region minimum wage | MEDIUM | Region-based configuration |
| Foreign employee special rules | MEDIUM | Flag-based branching |
| I-VAN integration complexity | LOW | Defer to Phase 2 |
| Data migration from existing systems | MEDIUM | CSV import templates |

---

## 9. Approval

| Role | Name | Date | Signature |
|------|------|------|-----------|
| Product Owner | | | |
| Technical Lead | | | |
| Chief Accountant | | | |
| BA Lead | | | |
