# Cash Module — Business Requirements Document (BRD)

**Version:** 1.0
**Date:** July 2026
**Author:** BA Lead & Chief Accountant (20+ yrs)
**Status:** Draft
**Circular:** 99/2025/TT-BTC, Law on Accounting 2015, Decree 123/2020/ND-CP

---

## 1. Executive Summary

Cash management is core financial operation. Every business processes daily cash receipts, payments, bank transfers, petty cash, and cash inventory. GoTax Cash module must handle full lifecycle: voucher creation → approval → execution → posting → reporting — compliant with Circular 99/2025/TT-BTC.

Current state: 0% complete. No code. Need full build.

---

## 2. Business Objectives

| # | Objective | Success Metric | Priority |
|---|-----------|----------------|----------|
| BO1 | Digital cash receipt/payment vouchers replace paper | Zero paper voucher loss | P0 |
| BO2 | Real-time cash balance visibility | Balance refresh < 1s | P0 |
| BO3 | Circular 99 compliant cash book (S07-DN) | Audit-ready reports | P0 |
| BO4 | Multi-currency cash (VND, USD, EUR) | All currencies supported | P0 |
| BO5 | Petty cash fund management | Every fund tracked to VND | P1 |
| BO6 | Periodic cash inventory with discrepancy handling | Inventory within 24h of period end | P1 |
| BO7 | Cash flow forecast | 30/60/90 day projection | P2 |
| BO8 | Bank integration (auto reconciliation) | 80% auto-match | P3 |

---

## 3. Scope

### In Scope (P0)
- Cash receipt voucher (Phiếu thu) — create, approve, post
- Cash payment voucher (Phiếu chi) — create, approve, post
- Cash book (Sổ quỹ tiền mặt S07-DN / S04a-DNN)
- Cash detail ledger (Sổ chi tiết quỹ S07a-DN / S04b-DNN)
- Cash transfer (rút tiền NH / gửi tiền NH / chuyển tiền)
- Bank account management
- Multi-currency (VND, USD, EUR, etc.)
- Cash balance inquiry (real-time)
- Approval workflow (kế toán trưởng → giám đốc)

### In Scope (P1)
- Petty cash fund (tạm ứng) — request, approve, settle
- Cash inventory (kiểm kê quỹ) — 08a-TT form
- Bank reconciliation (đối chiếu NH)
- Cash discrepancy handling (thừa/thiếu quỹ)

### In Scope (P2)
- Cash flow forecast (dự báo dòng tiền)
- Cash denomination tracking
- Integration with AR (phải thu) — auto receipt
- Integration with AP (phải trả) — auto payment

### Out of Scope
- POS hardware integration
- EBICS bank protocol
- ATM cash management
- Physical cash counting hardware

---

## 4. Functional Requirements

### FR1: Cash Receipt Voucher (Phiếu thu)

| ID | Requirement | Priority |
|----|-------------|----------|
| FR1.1 | User creates receipt voucher with: date, amount, currency, payer, reason, account code (Nợ 111/Có đối ứng), attachments | P0 |
| FR1.2 | System auto-generates sequential voucher number per cash book per year | P0 |
| FR1.3 | System validates: debit = credit, balance > 0, cash account is 1111/1112 | P0 |
| FR1.4 | Receipt types: customer payment, loan recovery, bank withdrawal, advance refund, other | P0 |
| FR1.5 | Workflow: Draft → Submitted → Approved → Posted | P0 |
| FR1.6 | Posted receipt updates cash book + GL immediately | P0 |
| FR1.7 | Print receipt on Mẫu 01-TT (Circular 99 format) | P0 |
| FR1.8 | Foreign currency receipt: record original currency + exchange rate + VND equivalent | P0 |

### FR2: Cash Payment Voucher (Phiếu chi)

| ID | Requirement | Priority |
|----|-------------|----------|
| FR2.1 | User creates payment voucher with: date, amount, currency, payee, reason, account code (Nợ đối ứng/Có 111), attachments | P0 |
| FR2.2 | System validates: sufficient cash balance (optional configurable) | P0 |
| FR2.3 | Same workflow: Draft → Submitted → Approved → Posted | P0 |
| FR2.4 | Payment types: supplier payment, salary, expense, bank deposit, advance, tax payment, other | P0 |
| FR2.5 | Print payment on Mẫu 02-TT (Circular 99 format) | P0 |
| FR2.6 | Withholding tax handling for payments subject to WHT | P1 |

### FR3: Cash Book (Sổ quỹ tiền mặt)

| ID | Requirement | Priority |
|----|-------------|----------|
| FR3.1 | Real-time cash book: opening balance, receipts, payments, closing balance per day | P0 |
| FR3.2 | Separate cash book per currency (VND, USD, EUR) | P0 |
| FR3.3 | Cash book maintains running balance after each transaction | P0 |
| FR3.4 | Print Sổ quỹ tiền mặt (S07-DN / S04a-DNN) with Circular 99 format | P0 |
| FR3.5 | Cash book must match GL account 111 balance at any time | P0 |

### FR4: Cash Transfer

| ID | Requirement | Priority |
|----|-------------|----------|
| FR4.1 | Bank withdrawal (rút NH về quỹ): Debit 111, Credit 112 | P0 |
| FR4.2 | Bank deposit (gửi tiền NH): Debit 112, Credit 111 | P0 |
| FR4.3 | Cash transfer between currencies (VND→USD): Debit 1112, Credit 1111 | P0 |
| FR4.4 | Transfer generates both a payment (source) and receipt (destination) | P0 |

### FR5: Petty Cash

| ID | Requirement | Priority |
|----|-------------|----------|
| FR5.1 | Create petty cash fund: amount, custodian, purpose, period | P1 |
| FR5.2 | Top-up fund when low | P1 |
| FR5.3 | Settlement: receipts against fund, return surplus | P1 |
| FR5.4 | Fund status: Active, Frozen, Closed | P1 |

### FR6: Cash Inventory

| ID | Requirement | Priority |
|----|-------------|----------|
| FR6.1 | Initiate inventory: all cash counted by denomination | P1 |
| FR6.2 | System generates inventory report (Mẫu 08a-TT) | P1 |
| FR6.3 | Discrepancy handling: book vs actual difference recorded | P1 |
| FR6.4 | Excess cash: Debit 111, Credit 338 | P1 |
| FR6.5 | Shortage: determine responsible party, Debit 138/334, Credit 111 | P1 |
| FR6.6 | Mandatory at period-end (monthly at minimum per Decree 99) | P1 |

### FR7: Cash Flow Statement (B03-DN)

| ID | Requirement | Priority |
|----|-------------|----------|
| FR7.1 | Direct method cash flow from operations | P0 |
| FR7.2 | Indirect method cash flow from operations | P1 |
| FR7.3 | Cash flow from investing activities | P0 |
| FR7.4 | Cash flow from financing activities | P0 |
| FR7.5 | Net cash increase/decrease = opening - closing | P0 |
| FR7.6 | B03-DN format per Circular 99 Appendix II | P0 |

### FR8: Multi-Currency

| ID | Requirement | Priority |
|----|-------------|----------|
| FR8.1 | Each transaction records original currency amount + exchange rate | P0 |
| FR8.2 | Period-end revaluation of foreign currency cash (Art. 52 Circular 99) | P1 |
| FR8.3 | FX gain/loss: Debit/Credit 515/635 | P1 |

### FR9: Approval Workflow

| ID | Requirement | Priority |
|----|-------------|----------|
| FR9.1 | Configurable approval rules per amount threshold | P0 |
| FR9.2 | Small amount (< configurable): kế toán trưởng approves | P0 |
| FR9.3 | Large amount: giám đốc approves | P0 |
| FR9.4 | Email/notification on approval request | P2 |
| FR9.5 | Rejection with reason | P0 |

---

## 5. Non-Functional Requirements

| ID | Requirement | Target |
|----|-------------|--------|
| NFR1 | Response time for any API | < 500ms p95 |
| NFR2 | Concurrent users | 100 |
| NFR3 | Audit trail for all cash transactions | Immutable log |
| NFR4 | Data retention | 10 years (per Law on Accounting) |
| NFR5 | Cash book accuracy | Balance always correct |
| NFR6 | Export formats | PDF, Excel, CSV |

---

## 6. Stakeholders

| Role | Interest |
|------|----------|
| Thủ quỹ (Cashier) | Daily receipt/payment execution, cash book |
| Kế toán thanh toán (AP/AR Accountant) | Voucher creation |
| Kế toán trưởng (Chief Accountant) | Approval, reporting |
| Giám đốc (Director) | Large payment approval |
| Kiểm toán (Auditor) | Cash audit trail |
| Ngân hàng (Bank) | Bank reconciliation |

---

## 7. Assumptions

- Cash accounts 1111 (VND), 1112 (foreign currency) exist in COA
- User authentication via existing JWT authMW
- Multi-tenant ready (company_id scoped)
- GL posting via existing JournalEntry mechanism
- PDF generation via existing template engine

---

## 8. Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| Cash balance not matching GL | Financial report errors | Mandatory periodic reconciliation |
| Unauthorized cash payment | Fraud | Multi-level approval workflow |
| Foreign currency FX loss | Financial loss | Automated revaluation |
| Data loss | Compliance violation | Immutable audit log + backup |
