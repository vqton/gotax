# Bank Module — Readiness Assessment

**Status: NOT PRODUCTION READY (~15% complete)**

**Assessment Date:** July 2026
**Circular:** 99/2025/TT-BTC (effective Jan 1, 2026)
**Reference:** Circular 200/2014/TT-BTC (replaced), Decree 123/2020/ND-CP, Law on Accounting 2015, Law on Credit Institutions 2024

---

## Executive Summary

Bank module ~15% complete. Foundation (CompanyBankAccount CRUD) exists. Core banking operations (reconciliation, statement import, payment initiation, loan tracking, deposit management) absent.

Cannot operate in PROD. Gap = ~85% of module.

---

## What Exists (Foundation Only)

| Artifact | Location | Status |
|----------|----------|--------|
| `CompanyBankAccount` model | `internal/domain/models.go:593` | Implemented |
| `BankAccountStatus` enum | `internal/domain/models.go:448` | Implemented |
| Bank account CRUD interface | `internal/domain/interfaces.go:167-171` | Implemented |
| Bank account CRUD (PG) | `internal/repository/pg_company.go:453-500` | Implemented |
| Bank account CRUD (Memory) | `internal/repository/memory_company.go:456-515` | Implemented |
| Bank account service | `internal/service/company.go:350-377` | Implemented |
| Bank account handler | `internal/handler/company.go:488-546` | Implemented |
| Bank account routes | `internal/handler/company.go:67-70, 107-109` | Implemented |
| Bank account tests | `internal/handler/company_handler_test.go:271-286` | Implemented |
| `CashTransfer` bank ops | `internal/domain/models.go:1485-1510` | Implemented |
| `PaymentMethod` bank types (EFT, BANK_TRANSFER) | `internal/domain/models.go:784-785` | Implemented |
| Account 1121 "Tien gui NH VND" seed | `migrations/002_gl_schema_circular99.sql` | Seed data |
| Account 1122 "Tien gui NH USD" seed | `migrations/002_gl_schema_circular99.sql` | Seed data |
| `company_bank_accounts` table | `migrations/003_company_schema.sql` | Migration exists |

---

## Gap Analysis — Full Matrix

### Domain Models — 2/14

| Model | Status | Priority | Notes |
|-------|--------|----------|-------|
| `CompanyBankAccount` (TKNH) | DONE | P0 | Basic fields only |
| `BankStatement` (Sao kê NH) | MISSING | P0 | Import from bank |
| `BankStatementLine` (Dòng sao kê) | MISSING | P0 | Individual transactions |
| `BankReconciliation` (Đối chiếu NH) | MISSING | P0 | Matching engine |
| `PaymentOrder` (UNC/UNC/TT) | MISSING | P0 | Electronic payment order |
| `PaymentOrderBatch` (Lệnh thanh toán lô) | MISSING | P1 | Batch payment |
| `LoanAgreement` (Hợp đồng vay) | MISSING | P1 | Account 341 tracking |
| `LoanDisbursement` (Giải ngân) | MISSING | P1 | Loan drawdown |
| `TermDeposit` (Tiền gửi có kỳ hạn) | MISSING | P1 | Account 1281 |
| `BankGuarantee` (Bảo lãnh NH) | MISSING | P2 | Letter of guarantee |
| `BankFeeSchedule` (Biểu phí NH) | MISSING | P2 | Fee tracking |
| `Cheque` (Séc) | MISSING | P2 | Cheque management |
| `CashFlowForecast` (Dự báo dòng tiền) | MISSING | P2 | Forecasting |
| `BankIntegrationProfile` (Tích hợp NH) | MISSING | P1 | Per-bank connection config |

### Repository — 2/16

| Implementation | Status | Priority |
|----------------|--------|----------|
| `CompanyBankAccount` CRUD | DONE | P0 |
| `BankStatement` import & query | MISSING | P0 |
| `BankReconciliation` engine | MISSING | P0 |
| `PaymentOrder` CRUD | MISSING | P0 |
| `PaymentOrderBatch` CRUD | MISSING | P1 |
| `LoanAgreement` CRUD | MISSING | P1 |
| `LoanDisbursement` CRUD | MISSING | P1 |
| `TermDeposit` CRUD | MISSING | P1 |
| `BankGuarantee` CRUD | MISSING | P2 |
| Bank statement import (MT940) | MISSING | P0 |
| Bank statement import (CSV) | MISSING | P0 |
| Bank statement import (OFX) | MISSING | P1 |
| Bank statement import (AEB43) | MISSING | P1 |
| Bank statement import (CODA) | MISSING | P1 |
| EBICS client | MISSING | P1 |
| Vietnamese banking API integration | MISSING | P1 |

### Service — 2/20+

| Method Group | Status | Priority |
|--------------|--------|----------|
| Bank account CRUD | DONE | P0 |
| Bank statement import engine | MISSING | P0 |
| Bank reconciliation engine | MISSING | P0 |
| Auto-reconciliation rules | MISSING | P0 |
| Manual reconciliation | MISSING | P0 |
| Payment order creation | MISSING | P0 |
| Payment order approval workflow | MISSING | P0 |
| Payment order submission to bank | MISSING | P1 |
| Payment batch processing | MISSING | P1 |
| Loan agreement management | MISSING | P1 |
| Loan interest calculation | MISSING | P1 |
| Disbursement processing | MISSING | P1 |
| Term deposit management | MISSING | P1 |
| Bank fee tracking | MISSING | P2 |
| Bank balance history | MISSING | P1 |
| Cheque management | MISSING | P2 |
| Draft S07-DN (Cash Book for bank) | MISSING | P0 |
| Draft S08-DN (Bank Deposit Ledger) | MISSING | P0 |
| Draft S09-DN (Bank Reconciliation) | MISSING | P0 |
| Cash flow forecast | MISSING | P2 |

### Handler/API — 5/30+

| Endpoint | Method | Status |
|----------|--------|--------|
| `POST /api/v1/companies/:id/bank-accounts` | Create bank account | DONE |
| `GET /api/v1/companies/:id/bank-accounts` | List bank accounts | DONE |
| `GET /api/v1/companies/:id/bank-accounts/:id` | Get bank account | DONE |
| `PUT /api/v1/companies/:id/bank-accounts/:id` | Update bank account | DONE |
| `DELETE /api/v1/companies/:id/bank-accounts/:id` | Deactivate | DONE |
| `POST /api/v1/bank/statements/import` | Import statement | MISSING |
| `GET /api/v1/bank/statements` | List statements | MISSING |
| `GET /api/v1/bank/statements/:id/lines` | Statement lines | MISSING |
| `POST /api/v1/bank/reconciliations` | Start reconciliation | MISSING |
| `POST /api/v1/bank/reconciliations/:id/match` | Match lines | MISSING |
| `POST /api/v1/bank/reconciliations/:id/complete` | Complete reconciliation | MISSING |
| `GET /api/v1/bank/reconciliations` | List reconciliations | MISSING |
| `POST /api/v1/bank/payment-orders` | Create payment order | MISSING |
| `GET /api/v1/bank/payment-orders` | List payment orders | MISSING |
| `GET /api/v1/bank/payment-orders/:id` | Get payment order | MISSING |
| `POST /api/v1/bank/payment-orders/:id/submit` | Submit to bank | MISSING |
| `POST /api/v1/bank/payment-orders/:id/approve` | Approve | MISSING |
| `POST /api/v1/bank/payment-batches` | Create batch | MISSING |
| `GET /api/v1/bank/payment-batches` | List batches | MISSING |
| `POST /api/v1/bank/payment-batches/:id/submit` | Submit batch | MISSING |
| `GET /api/v1/bank/loans` | List loans | MISSING |
| `POST /api/v1/bank/loans` | Create loan | MISSING |
| `GET /api/v1/bank/loans/:id/schedule` | Payment schedule | MISSING |
| `POST /api/v1/bank/loans/:id/disburse` | Disburse | MISSING |
| `POST /api/v1/bank/term-deposits` | Create deposit | MISSING |
| `GET /api/v1/bank/term-deposits` | List deposits | MISSING |
| `POST /api/v1/bank/term-deposits/:id/mature` | Process maturity | MISSING |
| `GET /api/v1/bank/balance` | Current balance | MISSING |
| `GET /api/v1/bank/balance/history` | Balance history | MISSING |
| `GET /api/v1/reports/bank-book` | S08-DN Bank ledger | MISSING |
| `GET /api/v1/reports/bank-reconciliation` | S09-DN reconciliation | MISSING |

### Migrations — 2/10

| Table | Status |
|-------|--------|
| `company_bank_accounts` | DONE |
| `bank_statements` | MISSING |
| `bank_statement_lines` | MISSING |
| `bank_reconciliations` | MISSING |
| `bank_reconciliation_matches` | MISSING |
| `payment_orders` | MISSING |
| `payment_order_batches` | MISSING |
| `loan_agreements` | MISSING |
| `loan_disbursements` | MISSING |
| `term_deposits` | MISSING |

### Reports — 0/6

| Report | Form | Status |
|--------|------|--------|
| Bank Deposit Ledger (Sổ tiền gửi NH) | S08-DN | MISSING |
| Bank Reconciliation Report | S09-DN | MISSING |
| Payment Order Register | — | MISSING |
| Loan Schedule | — | MISSING |
| Cash Flow Statement (B03-DN) | B03-DN | MISSING |
| Bank Balance History | — | MISSING |

---

## Legal Compliance Gaps

### Circular 99/2025/TT-BTC — Bank Articles

| Article | Requirement | Compliance |
|---------|-------------|------------|
| Art. 12 (TK 112) | Demand deposits tracking in VND & foreign currency | PARTIAL (account seed, no bank ops) |
| Art. 12 (TK 113) | Cash in transit | NOT IMPLEMENTED |
| Art. 12 (TK 128) | Held-to-maturity investments (term deposits 1281) | NOT IMPLEMENTED |
| Art. 12 (TK 341) | Loans and borrowings tracking | NOT IMPLEMENTED |
| Art. 14 | Financial statements including cash flows B03-DN | NOT IMPLEMENTED |
| Art. 15 | Cash flow statement (direct/indirect method) | NOT IMPLEMENTED |
| Art. 29 | Balance conversion for 112 accounts | NOT IMPLEMENTED |
| Appendix III | Accounting book templates (S08-DN, S09-DN) | NOT IMPLEMENTED |
| Art. 52 | Foreign currency bank account revaluation | NOT IMPLEMENTED |

### Decree 123/2020/ND-CP — E-Invoice & Payment

| Requirement | Compliance |
|-------------|------------|
| E-invoice for bank transfers > 200K VND | NOT IMPLEMENTED |
| Electronic payment threshold for VAT deduction | NOT IMPLEMENTED |
| UNC/UNC e-submission to tax authority | NOT IMPLEMENTED |

### SBV Regulations

| Regulation | Requirement | Compliance |
|------------|-------------|------------|
| Law on Credit Institutions 2024 | Payment service provider compliance | NOT IMPLEMENTED |
| Circular 32/2024/TT-NHNN | Electronic banking security | NOT IMPLEMENTED |
| Circular 17/2024/TT-NHNN | Payment intermediary services | NOT IMPLEMENTED |

---

## ERP Comparable Analysis

### MISA AMIS "Tiền gửi" Module
- Full bank deposit management with bank statement import
- Automated bank reconciliation with AI (MISA AVA) — R72 release
- UNC printing for 15+ Vietnamese banks (ACB, SHB, VIB, VPBank, etc.)
- E-banking integration for balance inquiry & payment
- Payment QR code generation
- Bank fee tracking
- Bank balance history
- S07-DN / S08-DN / S09-DN reports

### Fast Accounting "Tiền gửi NH"
- Bank deposit tracking with statement import
- UNC/UNC payment order management
- Bank transaction matching
- Multi-currency support
- Cash flow reporting

### BravoERP "Vốn bằng tiền"
- Full cash/bank/deposit/loan management
- Term deposit contract management (Khế ước tiền gửi)
- Loan agreement tracking (Hợp đồng vay)
- Interest calculation engine
- Cash flow planning & forecasting
- Fund management
- Multiple bank account tracking per currency

### Tryton "Account Statement"
- Bank statement import (MT940, OFX, AEB43, CODA)
- Account reconciliation with write-off
- Payment module (SEPA, Stripe, direct debit)
- Statement rules engine
- Party reconciliation

---

## Build Estimate

| Phase | Components | Est. Effort |
|-------|-----------|-------------|
| P0.1 | Domain models (BankStatement, BankStatementLine, BankReconciliation, PaymentOrder) | 3 days |
| P0.2 | Bank statement import engine (CSV, manual entry) | 5 days |
| P0.3 | Bank reconciliation engine (manual + auto-match rules) | 5 days |
| P0.4 | Payment order workflow (create → approve → submit) | 3 days |
| P1.1 | Loan agreement + disbursement + interest | 5 days |
| P1.2 | Term deposit management | 3 days |
| P1.3 | Bank integration (EBICS client, Vietnamese bank API) | 8 days |
| P1.4 | Reports (S08-DN, S09-DN, balance history) | 3 days |
| P2.1 | Cheque management | 3 days |
| P2.2 | Cash flow forecast | 3 days |
| P2.3 | MI/MT940/OFX/CODA import formats | 5 days |
| **Total** | | **~46 days (~9 weeks)** |

---

## Recommendation

**DO NOT DEPLOY.** Bank module ~15% complete.

Build P0 first (3 weeks): statement import + reconciliation + payment order workflow.

Priority order:
1. Statement import engine (CSV → normalized bank statement lines)
2. Bank reconciliation (manual match + auto-rule match engine)
3. Payment order management (UNC workflow)
4. Loan/deposit management
5. Bank integration (Vietnamese banks API/EBICS)
6. Reports (S08-DN, S09-DN)
