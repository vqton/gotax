# Tax Module — Business Requirements Document (BRD)

**Role:** BA Lead (20+ yrs) + Chief Accountant (20+ yrs)
**Date:** 2026-07-27
**Version:** 1.0
**Status:** DRAFT for review

---

## 1. Business Context

### 1.1 Problem Statement

Vietnamese enterprises must comply with a complex and rapidly evolving tax regulatory framework. The current GoTax backend provides general ledger and accounting capabilities under Circular 99/2025/TT-BTC but lacks a dedicated tax module. Users cannot:

- Calculate tax liabilities (VAT, CIT, PIT, TTDB, BVMT, FCT)
- Generate tax declaration forms in GDT-mandated XML format
- Submit declarations electronically to `thuedientu.gdt.gov.vn`
- Issue e-invoices compliant with Decree 254/2026/ND-CP
- Track tax payments and reconciliation
- Manage tax audit workflows

This gap makes GoTax non-viable as a production accounting solution for Vietnamese enterprises.

### 1.2 Target Users

| Persona | Role | Needs |
|---------|------|-------|
| Chief Accountant | Approves declarations, signs submissions | Dashboard, period-end tax review, audit trail |
| Tax Accountant | Prepares declarations, reconciles accounts | Tax calc, form generation, submission |
| AP Accountant | Handles input VAT, purchase invoices | VAT tracking, supplier invoice coding |
| AR Accountant | Handles output VAT, sales invoices | E-invoice issuance, output VAT tracking |
| Payroll Accountant | Handles PIT, social insurance | PIT calculation, declaration |
| CFO/Finance Director | Strategic tax planning, compliance | Tax reports, CIT finalization, risk dashboard |
| External Auditor | Audit tax provisions | Read-only access to tax records |

### 1.3 Business Objectives

1. **OBJ-1:** Enable complete tax declaration lifecycle for all Vietnamese tax types
2. **OBJ-2:** Achieve full GDT electronic submission compliance
3. **OBJ-3:** Enable e-invoice issuance per Decree 254/2026/ND-CP
4. **OBJ-4:** Provide tax payment tracking and reconciliation
5. **OBJ-5:** Support tax audit with full version history and drill-down
6. **OBJ-6:** Automate tax calculation from journal entries

---

## 2. Scope

### 2.1 In Scope

| Module | Tax Types | Forms |
|--------|-----------|-------|
| **VAT** | Output VAT, Input VAT, VAT payable/refundable | 01/GTGT, 02/GTGT, 03/GTGT, 04/GTGT, 05/GTGT, BK Mua vao, BK Ban ra |
| **CIT** | Provisional CIT, Final CIT, CIT incentives | 03/TNDN, 04/TNDN, 02/TNDN, 05/TNDN, Phu luc KHCN |
| **PIT** | Monthly/quarterly PIT, Annual PIT finalization | 05/KK-TNCN, 02/QTT-TNCN, 05-1/BK-TNCN |
| **TTDB** | Special consumption tax | 01/TTDB |
| **BVMT** | Environmental protection tax | 01/BVMT |
| **FCT** | Foreign contractor tax | 01/NTNN, 02/NTNN, 03/NTNN |
| **E-Invoice** | Invoice issuance, replacement, cancellation | TXML format per Decree 254/2026/ND-CP |
| **Tax Calendar** | Deadline tracking, auto-generation | Dashboard alerts |
| **Payments** | Tax payment orders, reconciliation | Payment tracking |

### 2.2 Out of Scope (Phase 2+)

- Global Minimum Tax / Pillar 2 (QDMTT, IIR filing)
- Transfer pricing documentation (TP Report, CbCR)
- Tax consolidation for multi-entity groups
- IFRS tax reporting
- Customs declaration integration
- Social insurance (BHXH) electronic submission

---

## 3. Functional Requirements

### 3.1 Tax Rate Management

| FR-ID | Description | Priority |
|-------|-------------|----------|
| FR-01 | System maintains tax rate table with effective dates | P0 |
| FR-02 | VAT rates: 0%, 5%, 8% (reduced), 10% (standard) | P0 |
| FR-03 | CIT rates: 15% (micro), 17% (small), 20% (standard), incentives | P0 |
| FR-04 | PIT progressive rate table (resident individuals) | P0 |
| FR-05 | PIT flat rate 20% (non-resident individuals) | P0 |
| FR-06 | TTDB rate table by product category | P0 |
| FR-07 | BVMT rate table by pollutant | P0 |
| FR-08 | FCT deemed rates by service/activity type | P0 |
| FR-09 | Rate history preserved for back-filing | P1 |
| FR-10 | Admin UI to update rates with audit trail | P1 |

### 3.2 Tax Declaration

| FR-ID | Description | Priority |
|-------|-------------|----------|
| FR-11 | Create tax declaration with type, period, company | P0 |
| FR-12 | Auto-populate declaration from journal entries | P0 |
| FR-13 | Manual override for declaration line items | P0 |
| FR-14 | Declaration status lifecycle: DRAFT→VALIDATED→SUBMITTED→ACKNOWLEDGED→AMENDED | P0 |
| FR-15 | Validation rules before submission (balance checks, completeness) | P0 |
| FR-16 | Generate XML file per GDT XSD schema | P0 |
| FR-17 | Display declaration preview in browser | P1 |
| FR-18 | Support amended/additional declarations | P1 |
| FR-19 | Declaration version history | P1 |
| FR-20 | Bulk operations (multi-company group declarations) | P2 |

### 3.3 Electronic Submission (GDT Integration)

| FR-ID | Description | Priority |
|-------|-------------|----------|
| FR-21 | HTTPS client for `thuedientu.gdt.gov.vn` APIs | P0 |
| FR-22 | Certificate-based authentication (USB Token / remote HSM) | P0 |
| FR-23 | Submit XML declaration file to GDT | P0 |
| FR-24 | Receive and process GDT acknowledgement | P0 |
| FR-25 | Receive and process GDT rejection/error codes | P0 |
| FR-26 | Retry mechanism for failed submissions | P1 |
| FR-27 | Query submission status from GDT | P1 |
| FR-28 | Download GDT-issued receipt/approval documents | P1 |

### 3.4 VAT Specific

| FR-ID | Description | Priority |
|-------|-------------|----------|
| FR-29 | Calculate output VAT from sales journal entries | P0 |
| FR-30 | Calculate input VAT from purchase journal entries | P0 |
| FR-31 | Generate VAT deduction declaration (01/GTGT) | P0 |
| FR-32 | Generate purchase ledger (Bang ke mua vao) | P0 |
| FR-33 | Generate sales ledger (Bang ke ban ra) | P0 |
| FR-34 | VAT reconciliation report | P0 |
| FR-35 | Track VAT refund eligibility | P1 |
| FR-36 | Support VAT allocation for multi-location enterprises | P1 |
| FR-37 | VAT direct method declaration (02/GTGT) | P1 |

### 3.5 CIT Specific

| FR-ID | Description | Priority |
|-------|-------------|----------|
| FR-38 | Calculate quarterly provisional CIT | P0 |
| FR-39 | Generate annual CIT finalization (03/TNDN) | P0 |
| FR-40 | Generate quarterly CIT provisional (04/TNDN) | P0 |
| FR-41 | Calculate CIT incentive eligibility | P0 |
| FR-42 | Track tax loss carry-forwards | P0 |
| FR-43 | CIT determination working paper | P1 |
| FR-44 | Support capital transfer CIT (02/TNDN, 05/TNDN) | P1 |
| FR-45 | Support CIT for foreign contractors | P1 |

### 3.6 PIT Specific

| FR-ID | Description | Priority |
|-------|-------------|----------|
| FR-46 | Calculate monthly/quarterly PIT per employee | P0 |
| FR-47 | Generate PIT deduction declaration (05/KK-TNCN) | P0 |
| FR-48 | Generate annual PIT finalization (02/QTT-TNCN) | P0 |
| FR-49 | Support dependant deduction tracking | P0 |
| FR-50 | Support tax code registration for employees | P1 |
| FR-51 | Support expatriate PIT finalization on exit | P1 |

### 3.7 E-Invoice

| FR-ID | Description | Priority |
|-------|-------------|----------|
| FR-52 | Create e-invoice from sales entry | P0 |
| FR-53 | Generate TXML invoice XML per GDT format | P0 |
| FR-54 | Sign invoice with digital certificate | P0 |
| FR-55 | Submit invoice to GDT for validation | P0 |
| FR-56 | Issue invoice to buyer (email, portal) | P0 |
| FR-57 | Invoice cancellation workflow | P0 |
| FR-58 | Invoice replacement/adjustment workflow | P0 |
| FR-59 | POS invoice connected to GDT (real-time) | P1 |
| FR-60 | Bulk invoice issuance | P1 |

### 3.8 Tax Payments

| FR-ID | Description | Priority |
|-------|-------------|----------|
| FR-61 | Generate tax payment order from declared amount | P1 |
| FR-62 | Track payment status (pending/paid/overdue) | P1 |
| FR-63 | Reconcile declared vs paid amounts | P1 |
| FR-64 | Calculate late payment interest | P1 |
| FR-65 | Generate tax payment receipt from bank | P2 |

### 3.9 Tax Calendar

| FR-ID | Description | Priority |
|-------|-------------|----------|
| FR-66 | Auto-generate tax deadlines per company regime | P1 |
| FR-67 | Dashboard alerts for upcoming deadlines | P1 |
| FR-68 | Missed deadline detection and escalation | P1 |

### 3.10 Tax Reports

| FR-ID | Description | Priority |
|-------|-------------|----------|
| FR-69 | VAT reconciliation (declared vs GL) | P0 |
| FR-70 | CIT computation worksheet | P0 |
| FR-71 | Tax liability aging report | P1 |
| FR-72 | PIT summary by employee and period | P1 |
| FR-73 | Multi-period tax comparison | P1 |
| FR-74 | Tax audit trail report | P1 |

---

## 4. Non-Functional Requirements

| NFR-ID | Description | Target |
|--------|-------------|--------|
| NFR-01 | Declaration generation from <10k journal entries | <30 seconds |
| NFR-02 | GDT submission round-trip | <60 seconds (incl. network) |
| NFR-03 | Maximum tax rate lookup | <100ms |
| NFR-04 | E-invoice issuance throughput | 100/minute minimum |
| NFR-05 | Multi-tenant data isolation | Complete (no cross-company tax data leak) |
| NFR-06 | Audit trail for all tax operations | Immutable, timestamped |
| NFR-07 | 99.5% uptime during tax season (Jan-Apr) | Target |
| NFR-08 | Responsive UI for declaration editing | <1s per interaction |
| NFR-09 | Secure certificate/key storage | Encrypted at rest |
| NFR-10 | Rate limit for GDT API calls | Configurable, per-company |

---

## 5. Stakeholder Map

| Stakeholder | Interest | Influence | Engagement |
|-------------|----------|-----------|------------|
| Chief Accountant (customer) | Tax compliance, audit readiness | High | Weekly review |
| Tax Accountant (user) | Daily tax operations | High | UX validation |
| CFO (economic buyer) | Cost, compliance risk | High | Monthly steering |
| GDT (regulator) | Submission format, timeliness | External | Via compliance |
| External Auditor | Tax provision correctness | Medium | Quarterly walkthrough |
| Product Manager | Roadmap, competitive parity | High | Daily decisions |
| Engineering | Build feasibility | High | Technical design |

---

## 6. Success Metrics

| KPI | Current | Target (6 months) | Measurement |
|-----|---------|--------------------|-------------|
| Tax declaration types supported | 0 | 15 (all common forms) | Feature count |
| GDT submission success rate | N/A | >99% | Submission logs |
| E-invoice issuance from GoTax | 0 | Enabled | Feature flag |
| Declaration auto-population accuracy | N/A | >95% from journal entries | Spot-check audit |
| Tax report types | 3 (GL only) | 8 (GL + 5 tax) | Feature count |
| User onboarding time for tax module | N/A | <2 hours | Time-to-first-declaration |
