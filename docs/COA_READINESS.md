# GoTax COA Module - Production Readiness Assessment

## Assessment Date: 2026-07-27 | Assessor: BA Lead + Chief Accountant (20+ yrs)

## VERDICT: ❌ NOT PRODUCTION READY

Current COA module is functional for basic GL operations but lacks enterprise-grade COA management features required for production deployment in Vietnamese enterprises competing with MISA/Fast/BravoERP.

---

## Gap Analysis vs Circular 99/2025/TT-BTC

| Requirement (Circular 99) | Current State | Target State | Gap |
|---|---|---|---|
| 71 Level-1 accounts per Appendix II | 71 level-1 seeded correctly | Must match TT99 exactly | OK |
| Loai 0 (Off-balance sheet) accounts | NOT implemented | Must support 001-009 accounts | **CRITICAL** |
| TK 215 (Biological assets) | NOT seeded | Required by TT99 Art. 4 | **HIGH** |
| TK 332 (Dividends payable) | NOT seeded | Required by TT99 | **HIGH** |
| TK 158 (Bonded warehouse) | NOT seeded | Required by TT99 | **HIGH** |
| TK 171 (Govt bond repurchase) | NOT seeded | Required by TT99 | **HIGH** |
| TK 82112 (Global min tax / Pillar 2) | NOT seeded | Required by TT99 Art. 5 | **HIGH** |
| Enterprise self-modify COA (no MoF approval) | Basic CRUD only | Full COA customization IU | **HIGH** |
| Internal Accounting Policy (IGAP) issuance | NOT implemented | Must document customizations | **HIGH** |
| Financial statement line-item preservation | NOT validated | Cross-checks on modification | **HIGH** |
| Multi-regime support (TT99, TT133, TT58) | TT99 only | Toggle between regimes | **HIGH** |

## Gap Analysis vs Enterprise ERP (MISA/Fast/Bravo)

| Feature | MISA | Fast | Bravo | GoTax | Gap |
|---|---|---|---|---|---|
| COA import/export (CSV/Excel) | Yes | Yes | Yes | **No** | **CRITICAL** |
| Mass account update/deactivate | Yes | Yes | Yes | **No** | **CRITICAL** |
| COA versioning & history | Yes | Partial | Yes | **No** | **CRITICAL** |
| Old-to-new account mapping | Yes | Yes | Yes | **No** | **CRITICAL** |
| Bank account linking | Yes | Yes | Yes | **No** | **HIGH** |
| Budget control per account | Yes | Yes | Yes | **No** | **HIGH** |
| Analysis codes (cost/profit center) | Yes | Yes | Yes | **No** | **HIGH** |
| Account approval workflow | Yes | Yes | Yes | **No** | **HIGH** |
| Auto account code generation | Yes | Yes | Yes | **No** | **MEDIUM** |
| Account freeze/lock | Yes | Yes | Yes | **No** | **HIGH** |
| Balance drill-down inquiry | Yes | Yes | Yes | **No** | **HIGH** |
| IFRS account mapping layer | Yes | No | Yes | **No** | **MEDIUM** |
| Account usage history report | Yes | Yes | Yes | **No** | **MEDIUM** |
| Inter-company account mapping | Yes | Partial | Yes | **No** | **MEDIUM** |
| 5-level deep hierarchy support | Yes | Yes | Yes | Unlimited | OK |
| REST API for COA management | Yes | Partial | Yes | Partial | **MEDIUM** |
| Multi-currency accounts | Yes | Yes | Yes | Basic | **MEDIUM** |
| Audit trail for COA changes | Yes | Yes | Yes | Basic | **MEDIUM** |

## Security & Compliance Gaps

| Requirement | Current | Target | Gap |
|---|---|---|---|
| AUTH module integration | Basic JWT | Full RBAC + VNeID | **CRITICAL** (per AUTH assessment) |
| Account modification audit | Basic | Immutable + who/what/when | **MEDIUM** |
| Sensitive account flagging | None | PII/tax data tagging | **HIGH** |
| COA change approval (4-eyes) | None | Dual-control required | **HIGH** |

## Risk Assessment

| Risk | Severity | Description | Mitigation |
|---|---|---|---|
| Data migration failure | HIGH | No COA import = manual data entry errors | Implement CSV import Phase 1 |
| Circular 99 non-compliance | HIGH | Missing TK 215, 332, 158, 171, 82112, Loai 0 | Add missing accounts immediately |
| Enterprise rejection | HIGH | Cannot match MISA/Fast COA flexibility | Feature gap closure roadmap |
| Audit finding | MEDIUM | No IGAP feature = auditor concern | Add IGAP document generator |
| IFRS conversion blocked | MEDIUM | No IFRS mapping layer | Design for Phase 2 |

## Implementation Priority (COA Module)

| Priority | Feature | Est. Effort | Dependencies |
|---|---|---|---|
| P0 | Add missing Circular 99 accounts (215, 332, 158, 171, 82112, Loai 0) | 1 week | DB migration |
| P0 | COA CSV/Excel import | 1 week | File parsing lib |
| P0 | Mass account update/deactivate | 0.5 week | Batch operation |
| P0 | Account freeze/lock (independent of period) | 1 week | Schema change |
| P0 | Account balance inquiry with drill-down | 1.5 weeks | Reporting engine |
| P1 | COA versioning & old/new mapping | 2 weeks | Version table |
| P1 | Account approval workflow (4-eyes) | 2 weeks | Workflow engine |
| P1 | Budget control per account | 2 weeks | Budget module |
| P1 | Analysis codes (cost/profit center) | 1 week | Schema extension |
| P1 | IGAP document generation | 1 week | Template engine |
| P1 | COA export (CSV/Excel/PDF) | 0.5 week | Export engine |
| P2 | IFRS account mapping layer | 2 weeks | Mapping table |
| P2 | Multi-regime toggle (TT99/TT133/TT58) | 2 weeks | Regime engine |
| P2 | Bank account integration | 1 week | External API |
| P2 | Auto account code generation | 0.5 week | Code generator |

**Total estimated effort: 16-22 weeks with 2-3 engineers**

## Recommendation

Phase 1 (P0 + P1): 10-14 weeks → Minimum viable enterprise COA capable of PROD
Phase 2 (P2): 6-8 weeks → Full enterprise parity with MISA/Bravo
