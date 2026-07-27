# BRD: GoTax Company Module

## Version: 1.0 | Date: 2026-07-27 | Status: DRAFT
## Governing Standards: Circular 99/2025/TT-BTC, Law on Enterprise 59/2025/QH15, Decree 01/2021/ND-CP

---

## 1. Executive Summary

### Project
GoTax - Vietnam Tax & Accounting Management System

### Business Need
Company (Doanh nghiep) entity is the root object of any accounting system. Without a proper Company module, multi-tenant operation, tax registration management, and regulatory compliance are impossible. Vietnamese enterprises must register tax codes, manage branches, subsidiaries, e-invoice accounts, digital signatures, and bank accounts — all tied to the Company entity.

Current GoTax GL module has implicit single-tenant assumption. Production deployment requires multi-company support with full Vietnam regulatory compliance.

### Strategic Alignment
- Law on Enterprise 59/2025/QH15 — legal entity registration and governance
- Circular 99/2025/TT-BTC — enterprise accounting regime, effective 01/01/2026
- Circular 58/2026/TT-BTC — micro-enterprise regime, effective 01/07/2026
- Decree 01/2021/ND-CP — business registration
- Decree 23/2025/ND-CP — digital signatures
- Decree 123/2020/ND-CP + Decree 04/2025/ND-CP — e-invoices
- GDT e-tax portal integration mandatory for all enterprises
- VNeID integration per Decision 29/2026/QD-TTg

### Research Context
| Source | Finding |
|--------|---------|
| Tryton Company module v8.0.1 (2026-07-01) | Production/Stable (Class 5). Supports single/multi-company, employee tree structure, extends Party model. Python/PyPI, GPL-3.0. 709 commits. |
| MISA SME 2026 / AMIS | Market leader ~45% SMB. Native GDT integration, e-invoice, multi-company. |
| Fast Accounting 2026 | ~25% share. Strong customization, trading/manufacturing. |
| BravoERP v10.9.4 (2026-05) | Full ERP, ISO 27001:2022 certified. Mid-market focus. |

---

## 2. Business Objectives

| # | Objective | KPI |
|---|-----------|-----|
| BO1 | Support multi-company (parent+subsidiaries) | Unlimited companies per tenant |
| BO2 | Full tax registration management (MST) | 100% GDT-compliant tax info |
| BO3 | E-invoice registration & management | Integration with GDT e-invoice portal |
| BO4 | Digital signature management | Decree 23/2025 compliant |
| BO5 | Bank account integration | Support all VN banks (Napas) |
| BO6 | Branch/subsidiary hierarchy | Tree structure, consolidation ready |
| BO7 | Fiscal year & period config | Support TT99, TT133, TT58 regimes |
| BO8 | Employee & department management | Org chart with cost allocation |
| BO9 | Integration readiness (GDT, Customs, BHXH, DVC) | API gateway ready |
| BO10 | Multi-regime support per company | Enterprise can choose TT99, TT133, or TT58 |

---

## 3. Scope

### In Scope (Phase 1)
- Company CRUD with full Vietnam regulatory fields
- Tax registration info (MST, MST chi nhanh, don vi phu thuoc)
- Business registration (Giay phep DKDN)
- Multi-company hierarchy (parent-child)
- Branch/subsidiary management
- Department/division management
- Employee management (basic: name, title, department)
- Fiscal year & period configuration
- Bank account management (account number, bank code, SWIFT)
- E-invoice registration (serial number, pattern, form)
- Digital signature registration (USB token + remote HSM)
- Integration profiles (GDT, Customs, BHXH, DVC)
- Audit trail for all company changes
- Multi-tenancy isolation

### In Scope (Phase 2)
- Employee payroll integration (BHXH/BHYT/BHTN)
- Company consolidation (intercompany elimination)
- IFRS entity mapping
- Workflow approval for company changes
- Document management (business license, tax certificate scans)
- Company merger/acquisition history
- Transfer pricing documentation

### Out of Scope
- Payroll calculation (separate HR module)
- Tax filing engine (separate Tax module)
- Consolidation reporting (separate Consol module)

---

## 4. Key Stakeholders

| Role | Interest |
|------|----------|
| Chief Accountant | Company tax registration, MST, fiscal year config, regulatory filing |
| Enterprise Owner/Director | Multi-company view, consolidation |
| Compliance Officer | Legal registration, licensing, regulatory deadlines |
| Accountant | Daily operations within assigned company |
| IT Admin | Multi-tenant setup, user-company assignment |
| External Auditor | Company legal structure, audit trail |
| Tax Authority (GDT) | MSDN, MST accuracy, e-invoice registration |
| System Integrator | API integration with existing ERP |

---

## 5. Legal & Regulatory Compliance Matrix

| Law/Decree/Circular | Requirement | Company Module Impact |
|---------------------|-------------|----------------------|
| Circular 99/2025/TT-BTC | Enterprise accounting regime from 01/01/2026 | Company accounting regime selection, fiscal year |
| Circular 58/2026/TT-BTC | Micro-enterprise regime from 01/07/2026 | Alternative regime for eligible companies |
| Circular 133/2016/TT-BTC | SME regime (optional) | Alternative regime for SME companies |
| Law on Enterprise 59/2025/QH15 | Legal entity registration, governance | Business registration info, legal form |
| Decree 01/2021/ND-CP | Business registration procedures | MSDN, Giay phep DKDN fields |
| Law on Tax Administration 108/2025/QH15 | Tax code registration, tax filing | MST, MST chi nhanh, tax office code |
| Decree 123/2020/ND-CP | E-invoice regulations | E-invoice registration fields |
| Decree 04/2025/ND-CP | E-invoice amendments | Updated e-invoice requirements |
| Decree 23/2025/ND-CP | Digital signatures | USB token + remote HSM registration |
| Decree 13/2023/ND-CP | Personal data protection | Employee data privacy |
| Circular 39/2014/TT-BTC | Expanded e-invoice (valid until replaced) | Legacy e-invoice support |
| Napas Circular 24/2019 | Bank account verification | Bank account validation |
| Decision 29/2026/QD-TTg | VNeID mandatory auth | Company admin identity verification |
| IFRS per Decision 345/QD-BTC | IFRS convergence roadmap | IFRS entity mapping (Phase 2) |

---

## 6. Functional Requirements

### FR1: Company Registration & Profile
| ID | Requirement | Priority |
|----|-------------|----------|
| FR1.1 | Create company with all mandatory legal fields | P0 |
| FR1.2 | Support legal forms: LLC, JSC, sole proprietorship, partnership, FO, RO | P0 |
| FR1.3 | Store business registration number (MSDN / Giay phep DKDN) | P0 |
| FR1.4 | Store tax code (MST) — primary identifier | P0 |
| FR1.5 | Store branch tax codes (MST chi nhanh) — 1:N | P0 |
| FR1.6 | Store company name (VN + EN) | P0 |
| FR1.7 | Store registered address, head office address | P0 |
| FR1.8 | Store phone, email, website | P0 |
| FR1.9 | Store legal representative (Giam doc / Chu tich) info | P0 |
| FR1.10 | Store chief accountant (Ke toan truong) info | P0 |
| FR1.11 | Store company type: manufacturing, trading, service, construction, etc. | P0 |
| FR1.12 | Store tax office (Chi cuc Thue) code | P0 |
| FR1.13 | Support company status: ACTIVE, SUSPENDED, DISSOLVED, MERGED | P0 |
| FR1.14 | Support company logo upload | P1 |

### FR2: Multi-Company Hierarchy
| ID | Requirement | Priority |
|----|-------------|----------|
| FR2.1 | Support parent-child company tree | P0 |
| FR2.2 | Support multiple subsidiaries per parent | P0 |
| FR2.3 | Support branch (Chi nhanh) entity type | P0 |
| FR2.4 | Support representative office (Van phong dai dien) | P0 |
| FR2.5 | Support dependent unit (Don vi phu thuoc) | P0 |
| FR2.6 | Company selection context for all user sessions | P0 |
| FR2.7 | Data isolation: each transaction belongs to one company | P0 |
| FR2.8 | Intercompany transaction tracking (Phase 2) | P2 |
| FR2.9 | Consolidation group configuration (Phase 2) | P2 |

### FR3: Fiscal Year & Period
| ID | Requirement | Priority |
|----|-------------|----------|
| FR3.1 | Configure fiscal year start month (Jan/Jul/Oct/other) | P0 |
| FR3.2 | Auto-generate periods (monthly, quarterly) from fiscal year config | P0 |
| FR3.3 | Support 12-month fiscal year + short period (for setup year) | P0 |
| FR3.4 | Support period state: OPEN, CLOSED, PERMANENTLY_CLOSED | P0 |
| FR3.5 | Validate no postings outside OPEN periods | P0 |
| FR3.6 | Different fiscal year per company in multi-company setup | P0 |
| FR3.7 | Period closing sequence: monthly → quarterly → annual | P0 |
| FR3.8 | Period reopen with audit trail | P1 |

### FR4: Accounting Regime Selection
| ID | Requirement | Priority |
|----|-------------|----------|
| FR4.1 | Company selects accounting regime: TT99 (enterprise), TT133 (SME), TT58 (micro) | P0 |
| FR4.2 | Regime selection determines COA template loaded at setup | P0 |
| FR4.3 | Support regime change at fiscal year boundary | P0 |
| FR4.4 | Regime change logs audit entry | P0 |
| FR4.5 | Circular 99/2025/TT-BTC is default for all new enterprises | P0 |

### FR5: Bank Account Management
| ID | Requirement | Priority |
|----|-------------|----------|
| FR5.1 | Register company bank accounts (1:N) | P0 |
| FR5.2 | Store: bank name, branch, account number, account holder, currency | P0 |
| FR5.3 | Support VN banks + foreign banks in Vietnam | P0 |
| FR5.4 | Support Napas bank code for electronic payments | P0 |
| FR5.5 | Bank account verification (trial deposit) | P1 |
| FR5.6 | Default bank account for cash transactions | P0 |
| FR5.7 | Bank account status: ACTIVE, CLOSED | P0 |

### FR6: E-Invoice Registration
| ID | Requirement | Priority |
|----|-------------|----------|
| FR6.1 | Register e-invoice pattern with GDT | P0 |
| FR6.2 | Store invoice serial numbers (Ky hieu) | P0 |
| FR6.3 | Store invoice form (Mau so) | P0 |
| FR6.4 | Support multiple invoice patterns per company | P0 |
| FR6.5 | Invoice template assignment per company | P0 |
| FR6.6 | E-invoice integration profile (API key, endpoint) | P0 |
| FR6.7 | Status tracking per invoice pattern: REGISTERED, ACTIVE, CANCELLED | P0 |

### FR7: Digital Signature Management
| ID | Requirement | Priority |
|----|-------------|----------|
| FR7.1 | Register USB token digital signature | P0 |
| FR7.2 | Register remote HSM digital signature (CloudCA, ViettelCA, VNPTCA, etc.) | P0 |
| FR7.3 | Store: provider, serial number, valid from/to, owner | P0 |
| FR7.4 | Assign default signature for tax filing | P0 |
| FR7.5 | Signature expiry alert (30 days before) | P1 |
| FR7.6 | Support multiple signatures per company | P0 |

### FR8: Department & Employee
| ID | Requirement | Priority |
|----|-------------|----------|
| FR8.1 | Create department tree under company | P0 |
| FR8.2 | Store: department code, name, manager, parent | P0 |
| FR8.3 | Create employee records linked to company | P0 |
| FR8.4 | Store: employee code, name, title, department, email, phone | P0 |
| FR8.5 | Employee status: ACTIVE, LEAVE, TERMINATED | P0 |
| FR8.6 | Employee tax code (MST ca nhan) | P0 |
| FR8.7 | Employee is also a system user (optional link) | P0 |
| FR8.8 | Employee department assignment with effective dating | P1 |
| FR8.9 | Org chart visualization | P2 |

### FR9: Integration Profiles
| ID | Requirement | Priority |
|----|-------------|----------|
| FR9.1 | GDT e-tax integration profile (HTKK/iHTKK credentials) | P0 |
| FR9.2 | Customs (Hai quan) integration profile | P1 |
| FR9.3 | Social Insurance (BHXH) integration profile | P1 |
| FR9.4 | DVC (Cong Dich vu Cong Quoc gia) integration | P1 |
| FR9.5 | API key management per integration endpoint | P0 |
| FR9.6 | Integration status: CONNECTED, DISCONNECTED, ERROR | P0 |

### FR10: Company Settings & Preferences
| ID | Requirement | Priority |
|----|-------------|----------|
| FR10.1 | Default currency (VND mandatory, foreign currency optional) | P0 |
| FR10.2 | Rounding precision (0/1/2 decimals) | P0 |
| FR10.3 | Date format (DD/MM/YYYY) | P0 |
| FR10.4 | Language (VN primary, EN secondary) | P0 |
| FR10.5 | Report header/footer templates | P1 |
| FR10.6 | Logo position in reports | P1 |

---

## 7. Non-Functional Requirements

| # | Requirement | Target |
|---|-------------|--------|
| NFR1 | Company profile load time | < 200ms |
| NFR2 | Company switch context (session switch) | < 500ms |
| NFR3 | Multi-company: max supported companies per tenant | Unlimited (tested to 1000+) |
| NFR4 | Concurrent users per company | 500 concurrent |
| NFR5 | Employee records per company | 10,000+ |
| NFR6 | Audit log retention | 5 years (Law on Accounting Art. 13) |
| NFR7 | Data isolation guarantee | Zero cross-company data leak |
| NFR8 | API response time (company CRUD) | < 100ms p95 |
| NFR9 | Uptime SLA | 99.9% (production) |

---

## 8. Risks & Mitigation

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Multi-company data leak | Critical | Low | Strict SQL row-level security, company_id on every table, test coverage |
| Tax code (MST) format change by GDT | Medium | Medium | Configurable MST validation regex, quick update path |
| Company regime change mid-year | High | Low | Enforce at fiscal year boundary only, migration wizard |
| GDT portal API changes | High | Medium | Adapter pattern per integration, retry + circuit breaker |
| Regulation overlap (TT99 vs TT133 vs TT58) | Medium | Low | Clear eligibility rules, regime selection guidance |
| Employee PII data breach | Critical | Low | Field-level encryption, Decree 13/2023 compliance |
| Legacy Circular 200 to TT99 migration | High | Medium | Transactional migration, dry-run mode, rollback |

---

## 9. Budget Estimate (Phase 1)

| Item | Cost (VND) |
|------|------------|
| Development (4 months, 3 engineers) | ~720M |
| Compliance consulting (VAA/MoF expert) | ~60M |
| GDT integration testing | ~40M |
| Security audit + penetration test | ~60M |
| Infrastructure (multi-tenant, HA) | ~80M |
| Legal review of regulatory compliance | ~30M |
| **Total Phase 1** | **~990M** |

---

## 10. Success Criteria

1. Enterprise can create and operate 10+ companies on single GoTax instance
2. Tax registration data passes GDT validation format
3. Company switch < 500ms with zero data cross-contamination
4. Fiscal year + period close works correctly for all regimes (TT99, TT133, TT58)
5. E-invoice registration syncs with GDT portal successfully
6. Digital signature integration signs test document successfully
7. Bank account reconciliation matches external bank statement
8. Zero data leak in security audit
9. Audit trail captures 100% of company changes
10. Employee-company-department hierarchy supports org structure accurately
