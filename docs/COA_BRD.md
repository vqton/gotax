# BRD: GoTax Chart of Accounts (COA) Module

## Version: 1.0 | Date: 2026-07-27 | Status: DRAFT
## Governing Standard: Circular 99/2025/TT-BTC (effective 01/01/2026)

---

## 1. Executive Summary

### Project
GoTax - Vietnam Tax & Accounting Management System

### Business Need
Chart of Accounts (He thong Tai khoan Ke toan) is the backbone of any accounting system. Current implementation provides basic CRUD within GL module but lacks enterprise-grade COA management features required by Vietnamese enterprises. Circular 99/2025/TT-BTC (effective 01/01/2026) replaced Circular 200/2014/TT-BTC with significant structural changes: 71 level-1 accounts (was 76), flexible enterprise self-modification (no MoF approval), new accounts (TK 215/332/158/171/82112), eliminated accounts (TK 611/631/441/461/466), and strengthened internal control requirements.

Competitors (MISA, Fast, BravoERP) offer mature COA modules with import/export, versioning, budget control, analysis codes, and approval workflows. GoTax must match these capabilities to be viable for enterprise deployment.

### Strategic Alignment
- Circular 99/2025/TT-BTC compliance — mandatory from 01/01/2026
- IFRS convergence roadmap per Decision 345/QD-BTC (Stage 3: 2026+)
- Accounting Strategy to 2030 per Decision 633/QD-TTg
- Enterprise digital transformation demand for cloud-based accounting

---

## 2. Business Objectives

| # | Objective | KPI |
|---|---|---|
| BO1 | Full Circular 99/2025/TT-BTC COA compliance | Zero audit finding on COA structure |
| BO2 | Enterprise-grade COA management at MISA/Bravo parity | Feature coverage >= 80% by Phase 1 |
| BO3 | Seamless COA migration from old Circular 200 | Import accuracy > 99.9% |
| BO4 | Support enterprise COA customization per TT99 Art. 4-5 | Self-service customization |
| BO5 | Enable IFRS dual-reporting capability | IFRS account mapping available |
| BO6 | Audit-ready COA change history | Immutable COA audit trail |

---

## 3. Scope

### In Scope (Phase 1)
- Full Circular 99 COA with all 71 Level-1 accounts + missing accounts (215, 332, 158, 171, 82112)
- Loai 0 (Off-balance sheet) accounts support
- COA CSV/Excel import with validation, mapping, rollback
- COA CSV/Excel/PDF export
- Mass account operations (update, activate, deactivate, freeze)
- Account balance inquiry with drill-down to journal entries
- Account approval workflow (create/edit/delete requires approval)
- COA versioning with effective dating
- Old (TT200) to new (TT99) account mapping table
- Internal Accounting Policy (IGAP) document generator
- Analysis codes: cost center, profit center, project, department
- Audit trail for all COA changes
- Account freeze/lock mechanism (independent of period close)

### In Scope (Phase 2)
- Budget control per account
- IFRS account mapping layer (VAS → IFRS conversion)
- Multi-regime support (TT99, TT133, TT58)
- Bank account integration (direct account-code linking)
- Auto account code generation based on type
- Inter-company account mapping
- Account usage analytics and recommendations

### Out of Scope
- Tax schedule management (covered by Tax module)
- Consolidation elimination rules (covered by Consolidation module)

---

## 4. Key Stakeholders

| Role | Interest |
|---|---|
| Chief Accountant (Ke toan truong) | COA structure compliance, audit trail, IGAP |
| Enterprise Accountant | Daily COA operations, reporting accuracy |
| IT Administrator | COA system configuration, user permissions |
| External Auditor | COA audit trail, VSA compliance |
| CFO | COA flexibility for management reporting, IFRS conversion |
| Tax Authority (GDT) | COA alignment with tax declarations |

---

## 5. Legal & Regulatory Compliance Matrix

| Law/Decree/Circular | Requirement | COA Module Impact |
|---|---|---|
| Circular 99/2025/TT-BTC | Enterprise accounting regime from 01/01/2026 | Primary COA standard, 71 L1 accounts, self-modification rights |
| Circular 133/2016/TT-BTC | SME accounting regime (optional) | Secondary COA for SME tenants |
| Circular 58/2026/TT-BTC | Micro-enterprise regime from 01/07/2026 | Tertiary COA for micro tenants |
| Circular 24/2024/TT-BTC | Administrative/non-business regime | Non-enterprise COA |
| Law on Accounting 88/2015/QH13 | Art. 22 defines account system | Foundational legal basis |
| Law 56/2024/QH15 | Amendments to Accounting Law 2024 | Legal basis for TT99 |
| Decision 345/QD-BTC (2020) | IFRS adoption roadmap | IFRS mapping layer requirement |
| Decision 633/QD-TTg | Accounting Strategy to 2030 | Long-term VFRS convergence |
| VSA (Vietnam Standards on Auditing) | Audit evidence integrity | COA change audit trail |
| Law on Tax Admin 108/2025/QH15 | Tax declaration alignment | COA ↔ tax codes mapping |
| Decree 174/2016/ND-CP | Detailing Accounting Law | COA documentation requirements |

---

## 6. Functional Requirements

### FR1: COA Structure Management
| ID | Requirement | Priority |
|---|---|---|
| FR1.1 | Pre-seed Circular 99 COA (71 L1 accounts + sub-accounts) | P0 |
| FR1.2 | Support Loai 0 (Off-balance sheet: 001-009) | P0 |
| FR1.3 | User can add new accounts (code, name, type, parent) | P0 |
| FR1.4 | User can modify account name, attributes (type locked if used) | P0 |
| FR1.5 | User can deactivate account (soft-delete, no impact on history) | P0 |
| FR1.6 | User can freeze account (prevent new postings, keep balance) | P0 |
| FR1.7 | System validates account code uniqueness, hierarchy, type | P0 |
| FR1.8 | System prevents deletion of account with balance or children | P0 |
| FR1.9 | Support 5+ level hierarchy depth | P0 |
| FR1.10 | System auto-generates account code suggestion by type | P1 |

### FR2: COA Import/Export
| ID | Requirement | Priority |
|---|---|---|
| FR2.1 | Import COA from CSV with header validation | P0 |
| FR2.2 | Import COA from Excel (.xlsx) with multi-sheet support | P0 |
| FR2.3 | Import validation: duplicate codes, invalid types, missing parents | P0 |
| FR2.4 | Import preview: show changes before committing | P0 |
| FR2.5 | Import rollback on failure (transactional) | P0 |
| FR2.6 | Import old (TT200) to new (TT99) account mapping | P0 |
| FR2.7 | Export COA to CSV with all attributes | P0 |
| FR2.8 | Export COA to Excel with formatted columns, filters | P1 |
| FR2.9 | Export COA to PDF (official chart document) | P1 |
| FR2.10 | Import audit log: record every imported change | P0 |

### FR3: COA Versioning & History
| ID | Requirement | Priority |
|---|---|---|
| FR3.1 | System snapshots COA on each change (versioned) | P0 |
| FR3.2 | User can view COA as-of any date | P0 |
| FR3.3 | User can compare two COA versions (diff view) | P0 |
| FR3.4 | System tracks effective dating for account changes | P0 |
| FR3.5 | User can revert to previous COA version (with audit) | P1 |

### FR4: COA Approval Workflow
| ID | Requirement | Priority |
|---|---|---|
| FR4.1 | Account creation requires 4-eyes approval (accountant + chief accountant) | P0 |
| FR4.2 | Account modification requires approval if account has balance | P0 |
| FR4.3 | Account deactivation requires approval | P0 |
| FR4.4 | Account freeze can be direct (no approval) | P0 |
| FR4.5 | Approval notification via system + email | P1 |
| FR4.6 | Bulk operations (import) require single approval | P0 |

### FR5: Analysis Codes
| ID | Requirement | Priority |
|---|---|---|
| FR5.1 | Cost center assignment per account | P1 |
| FR5.2 | Profit center assignment per account | P1 |
| FR5.3 | Department assignment per account | P0 (exists) |
| FR5.4 | Project assignment per account | P0 (exists) |
| FR5.5 | Contract assignment per account | P0 (exists) |
| FR5.6 | User-defined analysis dimensions | P2 |

### FR6: COA Reporting & Inquiry
| ID | Requirement | Priority |
|---|---|---|
| FR6.1 | Account balance inquiry with drill-down to journal lines | P0 |
| FR6.2 | Account usage report (which journal entries reference this account) | P0 |
| FR6.3 | Account activity report (period debit/credit, running balance) | P0 |
| FR6.4 | COA change history report | P0 |
| FR6.5 | Account statistics: parent/child count, depth level | P1 |
| FR6.6 | Chart view of accounts by type (treemap/sunburst) | P2 |

### FR7: IGAP (Internal Accounting Policy)
| ID | Requirement | Priority |
|---|---|---|
| FR7.1 | Auto-generate IGAP document from current COA configuration | P1 |
| FR7.2 | IGAP includes: COA, account descriptions, accounting methods | P1 |
| FR7.3 | IGAP export to PDF with company letterhead | P1 |
| FR7.4 | IGAP versioning aligned with COA versions | P1 |

### FR8: IFRS Mapping (Phase 2)
| ID | Requirement | Priority |
|---|---|---|
| FR8.1 | VAS account ↔ IFRS account mapping table | P2 |
| FR8.2 | IFRS reclassification rules per account | P2 |
| FR8.3 | IFRS adjustment journal auto-generation | P2 |
| FR8.4 | Dual reporting: VAS balance sheet + IFRS balance sheet | P2 |

---

## 7. Non-Functional Requirements

| # | Requirement | Target |
|---|---|---|
| NFR1 | COA load time (full chart, 500+ accounts) | < 1 second |
| NFR2 | Import throughput | 10,000 accounts/minute |
| NFR3 | COA query response (single account lookup) | < 50ms |
| NFR4 | Version storage | Keep all versions, unlimited retention |
| NFR5 | Concurrent COA editing | Lock-based, one editor at a time |
| NFR6 | Export file size | CSV < 5MB, Excel < 2MB |
| NFR7 | Audit log retention | 5 years (Law on Accounting Art. 13) |

---

## 8. Risks & Mitigation

| Risk | Impact | Probability | Mitigation |
|---|---|---|---|
| Circular 99 misinterpretation | High | Medium | Consult VAA/MoF guidelines, cross-check with KPMG/EY analysis |
| Data loss during COA migration | Critical | Low | Transactional import with rollback |
| Enterprise user resistance to cloud COA | Medium | Medium | Feature parity with MISA, offline fallback |
| IFRS convergence timeline changes | Medium | Low | Modular IFRS layer, decouple from core COA |
| User creates conflicting accounts | Medium | Medium | Validation rules, hierarchy checker, approval gate |

---

## 9. Budget Estimate (Phase 1)

| Item | Cost (VND) |
|---|---|
| Development (3 months, 3 engineers) | ~540M |
| COA migration consulting (VAA/MoF expert) | ~50M |
| Security audit | ~50M |
| Data migration for existing clients | ~30M |
| Infrastructure (cloud) | ~50M |
| **Total** | **~720M** |
