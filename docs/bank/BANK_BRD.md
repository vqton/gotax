# Bank Module — Business Requirements Document (BRD)

**Version:** 1.0
**Date:** July 2026
**Author:** BA Lead + Chief Accountant (20+ yrs each)
**Regulatory Basis:** Circular 99/2025/TT-BTC, Decree 123/2020/ND-CP, Law on Credit Institutions 2024, Law on Accounting 2015

---

## 1. Executive Summary

Bank module manages enterprise bank deposits (Account 112), bank transactions, statements, reconciliation, payment orders, loans, and term deposits. Circular 99/2025/TT-BTC renames Account 112 to "Demand deposits" and expands foreign currency revaluation requirements.

Current GoTax has CompanyBankAccount CRUD only. Full banking operations — statement import, reconciliation, payment initiation, loan tracking — required for PROD.

## 2. Business Objectives

| # | Objective | Success Metric | Priority |
|---|-----------|----------------|----------|
| OBJ-1 | Track all bank accounts per company in VND + foreign currencies | Zero untracked bank accounts | P0 |
| OBJ-2 | Import bank statements automatically (CSV/MT940/OFX) | < 5 min per import | P0 |
| OBJ-3 | Reconcile bank transactions against GL entries | Match rate > 90% auto | P0 |
| OBJ-4 | Generate and submit electronic payment orders to banks | UNC/UNC issuance within system | P0 |
| OBJ-5 | Track loan agreements, disbursements, and repayment schedules | Complete loan lifecycle | P1 |
| OBJ-6 | Manage term deposits and calculate interest | Auto interest calculation | P1 |
| OBJ-7 | Produce S08-DN (Bank Deposit Ledger) and S09-DN (Reconciliation) | Regulatory compliant | P0 |
| OBJ-8 | Support multi-currency bank accounts with auto revaluation | VAS/Circular 99 compliant | P1 |
| OBJ-9 | Track bank fees by transaction and period | Fee analysis reports | P2 |

## 3. Stakeholders

| Stakeholder | Role | Key Concern |
|-------------|------|-------------|
| Chief Accountant | Oversees bank reconciliation | Accuracy, completeness, audit trail |
| Treasury Manager | Manages cash position | Real-time balance, forecast |
| Accountant (AR) | Posts customer payments via bank | Ease of use, auto-match |
| Accountant (AP) | Posts supplier payments | Payment scheduling, approval |
| CFO/Treasurer | Strategic cash management | Cash flow forecast, liquidity |
| External Auditor | Verifies bank balances | Audit trail, reconciliation evidence |
| Tax Authority (GDT) | Reviews bank transactions | Tax compliance, undeclared income |

## 4. Functional Requirements

### FR-1: Bank Account Management

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-1.1 | Create bank account with bank code, name, branch, account number, holder, currency | P0 |
| FR-1.2 | Mark default bank account per currency | P0 |
| FR-1.3 | Deactivate/close bank account (no new transactions) | P0 |
| FR-1.4 | Support 20+ Vietnamese bank codes (VCB, CTG, BIDV, ACB, etc.) | P1 |
| FR-1.5 | Validate account number format per bank | P1 |
| FR-1.6 | Store IBAN/SWIFT code for international transfers | P1 |

### FR-2: Bank Statement Import

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-2.1 | Import bank statements via CSV file upload | P0 |
| FR-2.2 | Import MT940 (SWIFT) statement format | P1 |
| FR-2.3 | Import OFX/QFX format | P1 |
| FR-2.4 | Import Vietnamese banks' proprietary CSV formats (VCB, CTG, BIDV, etc.) | P0 |
| FR-2.5 | Manual statement line entry | P0 |
| FR-2.6 | Validate import: duplicate detection, balance continuity check | P0 |
| FR-2.7 | Store raw import file for audit trail | P0 |

### FR-3: Bank Reconciliation

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-3.1 | Start new reconciliation for a bank account + date range | P0 |
| FR-3.2 | Display statement lines vs GL entries side-by-side | P0 |
| FR-3.3 | Auto-match rules: amount + date + reference fuzzy match | P0 |
| FR-3.4 | Manual match / unmatch | P0 |
| FR-3.5 | Create write-off for minor differences (bank fees, rounding) | P0 |
| FR-3.6 | Complete reconciliation — lock period for bank account | P0 |
| FR-3.7 | Print S09-DN Bank Reconciliation Report | P0 |
| FR-3.8 | Historical reconciliation view | P1 |

### FR-4: Payment Orders (UNC/Lệnh chuyển tiền)

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-4.1 | Create single payment order (beneficiary, amount, bank, reference) | P0 |
| FR-4.2 | Create batch payment orders (salary, supplier batch) | P1 |
| FR-4.3 | Approval workflow for payment orders (maker-checker) | P0 |
| FR-4.4 | Print UNC form per Vietnamese bank template | P0 |
| FR-4.5 | Export payment file per bank format | P1 |
| FR-4.6 | Submit payment order status tracking (pending, sent, confirmed, failed) | P0 |
| FR-4.7 | Recurring payment orders | P2 |
| FR-4.8 | Integration with Vietnamese e-banking portals (optional) | P2 |

### FR-5: Loan Management

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-5.1 | Create loan agreement (bank, amount, interest rate, term, repayment schedule) | P1 |
| FR-5.2 | Track disbursements (partial/full) | P1 |
| FR-5.3 | Auto-generate repayment schedule (straight-line, annuity, bullet) | P1 |
| FR-5.4 | Record interest payments and principal payments | P1 |
| FR-5.5 | Loan balance tracking (Account 341) | P1 |
| FR-5.6 | Loan covenant tracking / early repayment penalty | P2 |

### FR-6: Term Deposit Management

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-6.1 | Create term deposit contract (bank, amount, term, interest rate, auto-renewal) | P1 |
| FR-6.2 | Track deposit balance (Account 1281) | P1 |
| FR-6.3 | Auto-calculate interest at maturity | P1 |
| FR-6.4 | Process maturity: principal + interest → bank account | P1 |
| FR-6.5 | Auto-renewal with new interest rate | P2 |

### FR-7: Bank Reports

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-7.1 | S07-DN Cash Book (combined cash + bank) | P0 |
| FR-7.2 | S08-DN Bank Deposit Ledger (Sổ tiền gửi ngân hàng) | P0 |
| FR-7.3 | S09-DN Bank Reconciliation Report (Bảng đối chiếu NH) | P0 |
| FR-7.4 | Bank Balance History Report | P1 |
| FR-7.5 | Bank Fee Analysis Report | P2 |
| FR-7.6 | Cash Flow Statement (B03-DN) bank component | P0 |

### FR-8: Bank Integration

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-8.1 | EBICS (Electronic Banking Internet Communication Standard) client | P2 |
| FR-8.2 | Vietnamese bank API integration (VCB B@NKS, CTG eBank, etc.) | P2 |
| FR-8.3 | Auto balance inquiry | P2 |
| FR-8.4 | Auto statement pull | P2 |
| FR-8.5 | Real-time payment status | P2 |

## 5. Non-Functional Requirements

| ID | Requirement | Target |
|----|-------------|--------|
| NFR-1 | Reconciliation performance | < 5s for 5000-line match |
| NFR-2 | Statement import throughput | 100K lines < 30s |
| NFR-3 | Audit trail for all bank ops | All mutations logged |
| NFR-4 | Data retention | 10 years (Circular 99 Art. 41) |
| NFR-5 | Security — payment approval | 2-person rule (maker-checker) |
| NFR-6 | Multi-tenant data isolation | Bank data per company only |
| NFR-7 | Bank account number encryption at rest | AES-256 for sensitive fields |
| NFR-8 | Availability | 99.9% uptime |

## 6. Assumptions & Constraints

| # | Description |
|---|-------------|
| A-1 | Vietnamese banks provide statement exports in CSV/MT940 formats |
| A-2 | Each company has at least one bank account |
| A-3 | GL journal entries already posted for bank-related transactions |
| A-4 | Circular 99/2025/TT-BTC is the governing accounting standard |
| A-5 | Bank integration via EBICS requires partnership with each bank |
| A-6 | UNC printing follows bank-specific form templates |
| C-1 | Must support both PG and in-memory backends |
| C-2 | Must follow existing GoTax architecture (Handler → Service → Repository) |
| C-3 | No external banking SDK dependencies (build from spec) |
| C-4 | Must support Vietnamese language throughout UI and reports |

## 7. Glossary

| Term | Definition |
|------|------------|
| UNC | Uỷ nhiệm chi — Payment order for bank transfer |
| UNT | Uỷ nhiệm thu — Collection order |
| TK 112 | Account 112 — Demand deposits (Circular 99) |
| S08-DN | Sổ tiền gửi ngân hàng — Bank deposit ledger |
| S09-DN | Bảng đối chiếu ngân hàng — Bank reconciliation report |
| EBICS | Electronic Banking Internet Communication Standard |
| MT940 | SWIFT bank statement message format |
| OFX | Open Financial Exchange format |
| B03-DN | Báo cáo lưu chuyển tiền tệ — Cash flow statement |
| Match rate | Percentage of statement lines automatically matched to GL entries |
| Maker-checker | Two-person approval rule (one creates, one approves) |
