# GoTax Company Module — Production Readiness Analysis

## Version: 1.0 | Date: 2026-07-27 | Status: DRAFT

---

## 1. Executive Verdict: NOT READY for PRODUCTION

**Current state:** GoTax has no Company module. GL module operates on implicit single-company assumption. Multi-tenant isolation, tax registration, fiscal year configuration, employee management, e-invoice registration, digital signature management, bank account integration, and government portal integration — all absent.

**Tryton Company Module v8.0.1 Assessment:** Production/Stable (PyPI Classifier 5). Mature with 18 years of development (2008-2026). However, Tryton is a Python-based monolithic ERP, not a Go microservice. Tryton's Company module cannot be directly used or adapted for GoTax architecture.

---

## 2. Gap Analysis: Current vs. Required

| Capability | Current (GoTax) | Required (Production) | Gap |
|------------|----------------|----------------------|-----|
| Multi-company | None | Full isolation, company context | CRITICAL |
| Tax registration (MST) | None | GDT-compliant MST fields | CRITICAL |
| Business registration | None | MSDN, Giay phep DKDN | CRITICAL |
| Fiscal year config | Implicit calendar year | Configurable per company | HIGH |
| Period management | Basic period model | Open/close/permanent close | HIGH |
| Accounting regime | TT99 only | TT99, TT133, TT58 selection | HIGH |
| Bank accounts | None | Multi-account, Napas codes | HIGH |
| E-invoice registration | None | GDT pattern, serial, form | HIGH |
| Digital signature | None | USB token, HSM registration | HIGH |
| Department/Employee | None | Tree structure, MST ca nhan | HIGH |
| Company profile | None | 15+ legal fields | CRITICAL |
| Audit trail for company changes | AuditLog exists but no company scope | Company-scoped immutable log | MEDIUM |
| Government portal integration | None | GDT, Customs, BHXH, DVC | HIGH |
| Data isolation guarantee | None | Row-level security | CRITICAL |
| Multi-tenancy | None | Tenant-per-company or tenant-per-group | CRITICAL |

---

## 3. Tryton Company Module v8.0.1 — Detailed Review

### Metadata
| Attribute | Value |
|-----------|-------|
| Latest version | 8.0.1 (2026-07-01) |
| PyPI Status | Production/Stable (5) |
| License | GPL-3.0-or-later |
| Python version | >= 3.10 |
| Author | B2CK SRL (Cedric Krier) |
| Maintainer | Tryton Foundation |
| Repository | https://code.tryton.org/tryton |
| Commits | 709 |
| Release history | 58 releases (v1.0.0 in 2008 to v8.0.1 in 2026) |

### Architecture
- Extends Party model (party → company)
- MultiValue fields for per-company property storage
- Employee model with company link
- Tree structure for parent-child company hierarchy
- PostgreSQL schema-based multi-tenancy option
- MultiValue concept: cost price, payable account stored per-company
- Referential data shared across companies by default
- Transactional documents linked to single company

### Key Models
1. **Company**: currency, employees, header/footer texts, parent_company
2. **Employee**: party, company, department, internal/external
3. **User**: current_company context, current_employee context

### Strengths
- 18 years mature, active development
- Multi-company support with data constraints
- MultiValue fields for per-company configuration
- PostgreSQL schema isolation capability
- Strong foreign-key constraint system to prevent cross-company data mixing
- Consolidation module available separately

### Weaknesses
- Python/PyPI — incompatible with Go stack
- Referential data shared by default — privacy concern
- No automatic intercompany transactions (manual booking)
- User documentation "rather scarce" per Capterra reviews
- No built-in Vietnamese regulatory fields (MST, MSDN, etc.)
- No e-invoice or digital signature management
- No built-in GDT/Customs/BHXH integration
- Steep learning curve for advanced features
- No Vietnamese accounting regime concept (TT99/TT133/TT58)

### G2 Review Snippets
- "Rock solid core system, comparable with Enterprise solutions" (5/5)
- "User documentation is rather scarce" (4/5)
- "Flexibility, rock solid, speed, simple UI" (4/5)
- "Does not allow accounting and support of business assets" (3/5)

### Verdict: Reference Architecture Only
Tryton Company module's **design patterns** (MultiValue fields, company context, employee-company link, party extension) are excellent reference. However, the module **cannot be used** in GoTax because:
1. Language mismatch (Python vs Go)
2. Missing Vietnam-specific fields (MST, MSDN, tax office code, etc.)
3. Missing Vietnamese regulatory compliance (TT99, e-invoice, digital signature)
4. No government portal integration (GDT, Customs, BHXH)
5. Monolithic vs microservice architecture mismatch

---

## 4. Competitor Benchmarking

### MISA SME 2026 / AMIS
| Feature | MISA | GoTax Target |
|---------|------|-------------|
| Multi-company | Yes (separate DBs) | Yes (same DB, isolated) |
| Tax registration | Full MST, MSDN | Full MST, MSDN |
| E-invoice | Native (meInvoice) | GDT integration |
| Digital signature | Supported | Supported |
| Bank integration | Major VN banks | Napas API |
| Employee/Dept | Full | Phase 1 basic |
| GDT integration | Built-in | API adapter |
| Regime support | TT99, TT133, TT58 | TT99 default, others configurable |
| AI assistant | AVA (AI helper) | Phase 3 |
| Market share | ~45% | Target |
| Price | ~9M/year (SaaS) | TBD |

### Fast Accounting 2026
| Feature | Fast | GoTax Target |
|---------|------|-------------|
| Multi-company | Yes (separate) | Yes (isolated) |
| Inventory | Strong (multi-warehouse) | Separate module |
| Customization | High (certified partners) | API-first |
| Tax compliance | Full VAS | Full TT99 |
| Market share | ~25% | Target |
| Price | ~15-25M (one-time) | Subscription |

### BravoERP v10.9.4
| Feature | Bravo | GoTax Target |
|---------|-------|-------------|
| Multi-company | Yes | Yes |
| Full ERP | Accounting + CRM + Mfg | Accounting + Tax first |
| ISO 27001:2022 | Certified | Target |
| Mobile app | Bravo 10 (v10.9.4) | Phase 2 |
| Market share | ~10% | Target |
| Price | ~30-200M | TBD |

---

## 5. Production Readiness Checklist

### Critical (Must Have for Go-Live)
- [ ] Company CRUD with 15+ legal fields
- [ ] Multi-company isolation (row-level security)
- [ ] Tax code (MST) validation per GDT format
- [ ] Fiscal year + period config per company
- [ ] Accounting regime selection (TT99 default)
- [ ] Bank account registration
- [ ] Company context in every API call
- [ ] Audit trail for company changes
- [ ] Data isolation guarantee (automated tests)
- [ ] Role-based access scoped to company

### High (Must Have within 30 Days)
- [ ] E-invoice registration (pattern, serial, form)
- [ ] Digital signature management
- [ ] Department tree
- [ ] Employee records
- [ ] GDT integration profile
- [ ] MST chi nhanh (branch tax codes)
- [ ] Multi-regime COA loading by company

### Medium (Phase 2)
- [ ] Org chart visualization
- [ ] Document management (scanned licenses)
- [ ] Intercompany transaction tracking
- [ ] Consolidation group config
- [ ] IFRS entity mapping
- [ ] Workflow approval for changes

---

## 6. Key Risks for PROD Deployment

1. **Data isolation breach**: Most critical risk. Need company_id on every table, RLS policies, and comprehensive integration test suite
2. **MST format changes**: GDT may change format. Design validation as configurable regex
3. **Multi-company performance**: Company context switch must be fast. Consider Redis cache for company config
4. **Regulatory overlap**: Companies eligible for multiple regimes (TT99 vs TT133). Clear guidance needed
5. **Migration from existing systems**: Excel/CSV import for company master data essential
6. **GDT API instability**: Adapter pattern with retry, circuit breaker, and manual fallback

---

## 7. Recommendation

**Do NOT deploy current GoTax to production for multi-company use cases.**

Minimum viable path to production:
1. **Week 1-2**: Company table design, CRUD API, row-level security
2. **Week 3-4**: Fiscal year + period management, regime selection
3. **Week 5-6**: Tax registration, bank account, e-invoice, digital signature
4. **Week 7**: Employee + department basic
5. **Week 8**: Integration profiles, GDT adapter
6. **Week 9-10**: Testing (unit, integration, security), documentation
7. **Week 11-12**: Pilot with 1-2 enterprises, bug fixes

Total: ~3 months to production-ready Company module.
