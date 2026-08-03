# Tax Module — Production Readiness Assessment

**Role:** BA Lead (20+ yrs) + Chief Accountant (20+ yrs)
**Date:** 2026-07-27
**Version:** 1.0

---

## Executive Summary

GoTax backend has strong accounting foundations (Circular 99/2025/TT-BTC COA, journal lifecycle, period management, basic GL reports) and skeleton tax infrastructure (company tax code, e-invoice pattern CRUD, digital signature CRUD, integration profile CRUD). **Round A (2026-08-03):** rate-table engines for VAT/CIT/PIT, declaration engine (GL→GTGT01/TNDN03 with cross-validation), and declaration→payment automation are now built and tested. **Round B (2026-08-03):** GDT infrastructure — xmldsig (exclusive C14N + RSA-SHA256), GDT API client (submit/status/cancel with retry), and full e-invoice issuance pipeline (TXML generation per Decree 254/2026, digital signing, GDT submit, status polling, cancel notification). **Round C (2026-08-03):** declaration submission — HTKK XML generation (`internal/htkk`, TAX_SPECS §3.1), detached signing, GDT declaration submit + status polling (ack → auto TaxPayment, reject → resubmit). Remaining: GDT real endpoints, adjustment/replacement e-invoices, VAT reconciliation, form PDFs.

**Verdict: NOT PRODUCTION-READY for tax compliance.**

Estimating 30-40% of work remains for a compliant Vietnamese tax module.

---

## Maturity Matrix

| Feature | Maturity | Coverage | PROD Gap |
|---------|----------|----------|----------|
| Circular 99 COA tax accounts | COMPLETE | 20+ tax accounts | None |
| Company tax code mgmt | COMPLETE | 10/13-digit validation, tax office | None |
| E-invoice pattern CRUD | COMPLETE | Schema + CRUD routes | None |
| Digital signature CRUD | PARTIAL | Schema + CRUD routes | Signing via env-key only; no per-signature key mgmt in issue path |
| GDT integration profile | STUB | Schema + CRUD | `TestIntegration` is no-op |
| **Tax declaration module** | **PARTIAL** | Engine: GL→GTGT01/TNDN03 lines + cross-validation + dup rejection + payment automation. CRUD + submit/ack/amend/cancel. HTKK XML gen + GDT submit + status poll (ack→payment, reject→resubmit) | Form PDFs, GDT real endpoint |
| **VAT calc** | **DONE** | Rate-table engine, explicit VAT accounts first (33311/1331/1332), fallback rate × revenue | None (form XML in Round C) |
| **CIT calc** | **DONE** | Size-based rates (MICRO 15% / SMALL 17% / STANDARD 20%), revenue−expenses+other | None (form XML in Round C) |
| **PIT calc** | **DONE** | Progressive brackets, resident/non-resident, dependants, insurance caps | KK-TNCN payroll input, form XML |
| **Tax rate tables** | **DONE** | tax_rates CRUD + `resolveRate` (suffix matching, active window, statutory fallback) | PG seed data for defaults |
| **HTKK/GDT integration** | **PARTIAL** | GDT client: invoice submit/status/cancel + declaration submit/status, retry + error map. HTKK XML (`internal/htkk`) | GDT real endpoint, cert auth |
| **E-invoice pipeline** | **PARTIAL** | TXML gen (Decree 254/2026), XMLDSig sign (C14N+RSA-SHA256), GDT submit/status/cancel, lifecycle DRAFT→SIGNED→SUBMITTED→ISSUED | GDT real endpoint, adjustment invoices, VAT reconciliation |
| **Tax payment tracking** | **PARTIAL** | TaxPayment CRUD, record payment, auto-create on declaration ACKNOWLEDGED + due dates | Payment journal auto-posting, late-interest wiring |
| **Tax audit support** | **MISSING** | Zero code | **Full build needed** |
| **FCT/withholding tax** | **MISSING** | Zero code | **Full build needed** |
| **Global minimum tax** | **MISSING** | Account 82112 exists | **Logic needed** |
| **Tax consolidation** | **MISSING** | ParentCompanyID exists | **Full build needed** |
| **Transfer pricing** | **MISSING** | Zero code | **Full build needed** |

---

## Regulatory Compliance Checklist

| Regulation | Effective | Required By | GoTax Status |
|-----------|-----------|-------------|--------------|
| Circular 99/2025/TT-BTC | 01-Jan-2026 | All enterprises | COA seeded, rest missing |
| Circular 80/2021/TT-BTC | 01-Jan-2022 | Tax declaration | **Not implemented** |
| CIT Law 67/2025/QH15 | 01-Oct-2025 | All enterprises | **Not implemented** |
| Decree 320/2025/ND-CP | 15-Dec-2025 | CIT guidance | **Not implemented** |
| Circular 20/2026/TT-BTC | 12-Mar-2026 | CIT details | **Not implemented** |
| Decree 236/2025/ND-CP (GMT) | 15-Oct-2025 | MNE groups >EUR750M | **Not implemented** |
| Decree 254/2026/ND-CP (einvoice) | 2026 | All enterprises | **Not implemented** |
| Decree 70/2025/ND-CP | 01-Jun-2025 | E-invoice | **Not implemented** |
| Circular 78/2021/TT-BTC | 01-Jul-2022 | HTKK e-tax | **Not implemented** |
| Decree 310/2025/ND-CP (penalties) | 16-Jan-2026 | Tax penalties ref | **Not implemented** |
| Law on Tax Admin 108/2025/QH15 | 01-Jul-2026 | Tax admin | **Not implemented** |
| PIT Law (new) | 01-Jul-2026 | Personal income tax | **Not implemented** |
| VAT Law 48/2024/QH15 | 01-Jul-2025 | Value-added tax | **Not implemented** |
| Decree 255/2026/ND-CP (TP) | 01-Jul-2026 | Transfer pricing | **Not implemented** |

---

## Key Risk Areas

### 1. No Tax Declaration Engine
Users cannot generate or submit any tax form (GTGT, TNDN, TNCN) from GoTax. This is the core function of any accounting software in Vietnam.

### 2. No E-Invoice Issuance
Since 01-Jul-2022, all enterprises must use e-invoices per Decree 123/2020/ND-CP. GoTax registers invoice patterns but cannot issue a single e-invoice. This is a compliance showstopper.

### 3. No GDT Electronic Submission
GoTax cannot connect to `thuedientu.gdt.gov.vn` to submit declarations, receive acknowledgements, or query taxpayer status.

### 4. No Tax Rate Management
VAT rates (0%, 5%, 8%, 10%), CIT rates (15% micro, 17% small, 20% standard), PIT brackets, TTDB rates are all hard-coded concepts — no configurable tax rate tables exist.

### 5. No Tax Calculation Logic
GoTax has debit/credit journal entries but no engine that computes:
- VAT payable/receivable from purchase/sales journals
- CIT provisional and final amounts
- PIT withholding amounts
- TTDB, BVMT, resource tax liabilities

### 6. Partial COA Only
Circular 99 COA is seeded but Circular 99 also requires:
- New financial statement templates (B01-DN, B02-DN, etc.)
- IFRS convergence disclosures
- Global minimum tax footnote disclosures
- Multi-branch consolidation rules

---

## PROD Readiness Verdict

```
╔══════════════════════════════════════════════════╗
║  TAX MODULE: NOT PRODUCTION READY               ║
║                                                  ║
║  Assessment: PRE-ALPHA (foundation only)         ║
║  Completion: ~20%                                ║
║  Remaining: ~80% (full module build)             ║
║  Timeline estimate: 6-9 months (4-6 devs)       ║
║  Risk level: CRITICAL for Vietnamese market      ║
╚══════════════════════════════════════════════════╝
```

## Recommendation

1. **Phase 1 (Foundation — 2 months):** Build tax domain models, tax rate tables, tax declaration lifecycle, basic VAT/CIT calculation engine.
2. **Phase 2 (Forms — 2 months):** XML generation for all tax forms per Circular 80, HTKK format support, PDF rendering.
3. **Phase 3 (Integration — 1.5 months):** GDT API client (`thuedientu.gdt.gov.vn`), certificate-based auth, submission/acknowledgement handling.
4. **Phase 4 (E-Invoice — 1.5 months):** Invoice data model, TXML generation, digital signing, GDT submission per Decree 254/2026/ND-CP.
5. **Phase 5 (Advanced — 1 month):** Tax reports, payment tracking, audit support, tax consolidation, transfer pricing, GMT/Pillar 2.
