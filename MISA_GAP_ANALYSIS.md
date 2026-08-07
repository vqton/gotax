# GoTax vs MISA SME 2026 — Module Gap Analysis

**Date:** 2026-08-07
**MISA version:** SME 2026 R6 (latest)
**GoTax status:** In-memory + PG backends, 9 modules complete

## MISA SME 2026 — 20 Business Processes

Source: `helpsme.misa.vn/2026` official knowledge base

| # | MISA Module | Vietnamese | GoTax Status | Notes |
|---|------------|-----------|-------------|-------|
| 1 | Fund/Treasury | Quỹ | ⚠️ Partial | Cash module exists (receipts/payments/transfers). No fund-level mgmt, no multi-cashier allocation |
| 2 | Bank | Ngân hàng | ✅ Done | Statements, reconciliation, loans (backend), deposits |
| 3 | Purchase | Mua hàng | ✅ Done | PO → GRN → Invoice → AP aging |
| 4 | Sales | Bán hàng | ✅ Done | Quote → Order → Delivery → Invoice → AR aging |
| 5 | E-Invoice | Hóa đơn điện tử | ⚠️ Backend only | XML gen (`internal/einvoice/`) + GDT client (`internal/gdt/`) exist. **No UI, no cancel/adjust/replace, no template mgmt** |
| 6 | Warehouse | Kho | ✅ Done | Items, balances, transfers, adjustments, takes, valuations, keeper |
| 7 | Tools & Instruments | CCDC | ✅ Done | Tools/instruments tracking |
| 8 | Fixed Assets | TSCĐ | ✅ Done | Categories, depreciation |
| 9 | Invoice Management | Quản lý hóa đơn | ❌ Missing | Invoice tracking, status mgmt, query, VAT lookup |
| 10 | Tax | Thuế | ⚠️ Partial | Declarations, payments, reconciliation, calendar, VAT report. **No eTax submit, no CIT/PIT XML gen, no BHXH declaration** |
| 11 | GL / Consolidation | Tổng hợp | ⚠️ Partial | GL done (journal entries, periods, COA, exchange rates). **No multi-branch consolidation, no inter-company elimination** |
| 12 | **Cost Accounting** | **Giá thành** | ❌ Missing | BOM costing, WIP, job costing, cost allocation by cost center |
| 13 | Payroll | Tiền lương | ✅ Done | Periods, timekeeping, payslips, declarations, config |
| 14 | **Loan Agreements** | **Khế ước vay** | ⚠️ Backend only | `LoanAgreement` repo exists in `pg_bank.go`. **No handler, no UI, no repayment schedule display** |
| 15 | **Contracts** | **Hợp đồng** | ❌ Missing | Sales/purchase contracts, milestones, linkage to invoices |
| 16 | Budget | Ngân sách | ✅ Done | Budget entries |
| 17 | Warehouse Keeper | Thủ kho | ✅ Done | Keeper role, custody tracking |
| 18 | **Cashier** | **Thủ quỹ** | ❌ Missing | Dedicated cashier role, cash drawer mgmt, daily close |
| 19 | **Financial Analysis** | **Phân tích tài chính** | ❌ Missing | Ratio analysis, trend charts, variance analysis, cash flow forecast |
| 20 | Reports | Báo cáo | ⚠️ Basic | 4 reports (TB, BS, P&L, CF) vs MISA's 200+. **Missing: AP/AR aging detail, tax reports, cost reports, departmental P&L, budget variance** |

## Score

| Status | Count | Modules |
|--------|-------|---------|
| ✅ Done | 9 | Bank, Purchase, Sales, Warehouse, CCDC, Fixed Assets, Payroll, Budget, Warehouse Keeper |
| ⚠️ Partial | 7 | Fund, E-Invoice, Tax, GL/Consolidation, Loan Agreements, Reports, (Cashier concept in Cash) |
| ❌ Missing | 4 | Cost Accounting, Contracts, Cashier, Financial Analysis |

## Priority Build Order

### P0 — High-value, blocks other modules
1. **Cost Accounting (Giá thành)** — required for manufacturing/construction. Needs: cost pools, allocation rules, WIP tracking, job costing, unit cost calc. Builds on existing Cost Centers module.
2. **E-Invoice UI** — backend done. Needs: handler + HTML pages for create/sign/send/cancel/adjust/replace, template management, invoice status tracking, GDT query integration.
3. **Multi-branch Consolidation** — GL done but no branch concept. Needs: Branch entity model, inter-company transactions, elimination entries, consolidated BCTC per Circular 99.

### P1 — Direct accounting value
4. **Loan Agreement UI** — repo exists. Needs: handler + page for disbursement schedule, repayment tracking, interest calc, maturity alerts.
5. **Contracts** — new module. Needs: contract CRUD, link to PO/SO, milestone tracking, expiry alerts, linkage to invoices.
6. **Financial Analysis** — report layer. Needs: ratio calc (liquidity, profitability, leverage), trend comparison, variance analysis, cash flow forecast, dashboard charts.

### P2 — Completeness
7. **Cashier Module** — separate from Cash. Needs: cashier assignment, cash drawer open/close, denomination tracking, daily reconciliation.
8. **Invoice Management** — VAT invoice tracking, status query, lookup by tax code, invalid invoice alerts.
9. **Tax Module Completion** — CIT/PIT XML generation (HTKK format), eTax portal submit, BHXH declaration, PIT withholding vouchers.
10. **Reports Expansion** — AP/AR aging detail, tax summary, cost reports, departmental P&L, budget vs actual, trial balance by cost center.

## Existing Code Assets (Reusable)

| Asset | Location | Status |
|-------|----------|--------|
| E-Invoice XML gen | `internal/einvoice/` | Done — Decree 254/2026 schema |
| GDT API client | `internal/gdt/` | Done — status query, push |
| XML digital signature | `internal/xmldsig/` | Done — RSA signing |
| HTKK tax form XML | `internal/htkk/` | Done — but only basic forms |
| Cost Centers | `internal/service/` + handler | Done — can extend for cost pools |
| Budget module | handler + service | Done — can extend for variance |
| Warehouse Keeper | handler + service | Done — cashier follows same pattern |
| Loan Agreement repo | `internal/repository/pg_bank.go` | Done — needs handler + UI |
| MISA sidebar nav | `web/static/js/app.js` | Done — add new modules to MODULE_GROUPS |

## Specs Written

| Module | Spec Location | Status |
|--------|--------------|--------|
| Cost Accounting (Giá thành) | `docs/COST_ACCOUNTING_SPEC.md` (2270 lines) | ✅ Complete — BRD, 6 methods, 17 use cases, API spec, DB schema, 4-phase plan |

## Notes

- MISA SME 2026 ships with 200+ built-in reports. GoTax has 4. Report expansion is critical for adoption.
- MISA's "Quỹ" (Fund) module is separate from "Thủ quỹ" (Cashier). Fund = allocation/tracking across multiple cash funds. Cashier = daily drawer operations.
- MISA's "Tổng hợp" (Consolidation) supports multi-branch BCTC with inter-company elimination per Circular 99 Appendix 04.
- Cost Accounting in MISA supports: simplified method, coefficient method, job costing, process costing. All need implementing.
- Cost Accounting spec covers all 6 methods per Circular 99/2025/TT-BTC with full GL account mapping (TK 154, 621-627, 632).
