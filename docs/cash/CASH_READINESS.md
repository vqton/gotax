# Cash Module — Readiness Assessment

**Status: NOT PRODUCTION READY (0% complete)**

**Assessment Date:** July 2026
**Circular:** 99/2025/TT-BTC (effective Jan 1, 2026)
**Reference:** Decree 123/2020/ND-CP, Law on Accounting 2015

---

## Executive Summary

Cash module 0% complete. No domain models, no repository, no service, no handler, no routes, no migrations, no docs. Only COA seed data (accounts 111, 1111, 1112, 112, 113) and voucher-type enums (`THU`/`CHI`) exist as foundation.

Cannot operate in PROD. Gap = entire module.

---

## What Exists (Foundation Only)

| Artifact | Location | Status |
|----------|----------|--------|
| Account 1111 "Tien mat VND" | `migrations/002_gl_schema_circular99.sql:227` | Seed data |
| Account 1112 "Tien mat USD" | `migrations/002_gl_schema_circular99.sql:229` | Seed data |
| Account 1121 "Tien gui NH VND" | `migrations/002_gl_schema_circular99.sql:231` | Seed data |
| Account 113 "Tien dang chuyen" | `migrations/002_gl_schema_circular99.sql:233` | Seed data |
| `VoucherTypeReceipt` / `VoucherTypePayment` | `internal/domain/models.go:23` | Enum only, no logic |
| `PayMethodCASH` / `PayMethodCQ` | `internal/domain/models.go:786-787` | Tax module only |
| B03-DN Cash Flow Statement requirement | `docs/archive/BRD_GL_MODULE.md:238` | P0 requirement, not implemented |

---

## Gap Analysis — Full Matrix

### Domain Models — 0/10

| Model | Status | Priority |
|-------|--------|----------|
| `CashBook` (Sổ quỹ) | MISSING | P0 |
| `CashReceipt` (Phiếu thu) | MISSING | P0 |
| `CashPayment` (Phiếu chi) | MISSING | P0 |
| `CashTransfer` (Rút/gửi tiền NH) | MISSING | P0 |
| `PettyCash` (Quỹ tạm ứng) | MISSING | P1 |
| `CashForecast` (Dự báo dòng tiền) | MISSING | P2 |
| `CashInventory` (Kiểm kê quỹ) | MISSING | P1 |
| `BankAccount` (TK ngân hàng) | MISSING | P0 |
| `CurrencyHolding` (Tiền theo loại) | MISSING | P1 |
| `CashDenomination` (Mệnh giá) | MISSING | P2 |

### Repository — 0/2

| Implementation | Status | Priority |
|----------------|--------|----------|
| `CashRepository` interface | MISSING | P0 |
| `PGCashRepo` | MISSING | P0 |
| `MemoryCashRepo` | MISSING | P0 |

### Service — 0/1

| Method Group | Status | Priority |
|--------------|--------|----------|
| `CashService` | MISSING | P0 |
| Receipt processing | MISSING | P0 |
| Payment processing | MISSING | P0 |
| Cash balance management | MISSING | P0 |
| Petty cash management | MISSING | P1 |
| Cash inventory | MISSING | P1 |
| Cash flow forecast | MISSING | P2 |
| Bank reconciliation | MISSING | P1 |

### Handler/API — 0/15+

| Endpoint | Method | Status |
|----------|--------|--------|
| `POST /api/v1/cash/receipts` | Create receipt | MISSING |
| `GET /api/v1/cash/receipts` | List receipts | MISSING |
| `GET /api/v1/cash/receipts/:id` | Get receipt | MISSING |
| `POST /api/v1/cash/payments` | Create payment | MISSING |
| `GET /api/v1/cash/payments` | List payments | MISSING |
| `GET /api/v1/cash/payments/:id` | Get payment | MISSING |
| `POST /api/v1/cash/transfers` | Bank transfer | MISSING |
| `GET /api/v1/cash/cash-book` | Cash book report | MISSING |
| `POST /api/v1/cash/inventory` | Cash inventory | MISSING |
| `GET /api/v1/cash/balance` | Current balance | MISSING |
| `GET /api/v1/cash/forecast` | Cash flow forecast | MISSING |
| `POST /api/v1/cash/petty-cash` | Create petty cash | MISSING |
| `GET /api/v1/cash/bank-accounts` | List bank accounts | MISSING |
| `GET /api/v1/reports/cash-flow` | B03-DN statement | MISSING |
| `GET /api/v1/reports/cash-book` | Sổ quỹ (S07-DN) | MISSING |

### Migrations — 0/1

| Table | Status |
|-------|--------|
| `cash_receipts` | MISSING |
| `cash_payments` | MISSING |
| `cash_books` | MISSING |
| `cash_transfers` | MISSING |
| `cash_inventories` | MISSING |
| `petty_cash_funds` | MISSING |
| `cash_forecasts` | MISSING |
| `bank_accounts` | MISSING |

### Reports — 0/5

| Report | Form | Status |
|--------|------|--------|
| Cash Book (Sổ quỹ tiền mặt) | S07-DN / S04a-DNN | MISSING |
| Cash Detail Ledger (Sổ chi tiết quỹ) | S07a-DN / S04b-DNN | MISSING |
| Cash Inventory Report | Mẫu 08a-TT | MISSING |
| Cash Flow Statement (Báo cáo LCTT) | B03-DN | MISSING |
| Cash Forecast Report | — | MISSING |

---

## Legal Compliance Gaps

### Circular 99/2025/TT-BTC — Cash Articles

| Article | Requirement | Compliance |
|---------|-------------|------------|
| Art. 12 (TK 111) | Cash receipt/payment must have vouchers with signatures | NOT IMPLEMENTED |
| Art. 12 (TK 112) | Bank deposit tracking | NOT IMPLEMENTED |
| Art. 12 (TK 113) | Cash in transit | NOT IMPLEMENTED |
| Art. 14 | Financial statements include cash flows | NOT IMPLEMENTED (B03-DN) |
| Art. 3 | Internal governance for cash transactions | NOT IMPLEMENTED |
| Appendix III | Accounting book templates (S07-DN) | NOT IMPLEMENTED |
| Art. 52 | Foreign currency cash revaluation | NOT IMPLEMENTED |

### Decree 123/2020/ND-CP — E-Invoice for Cash

| Requirement | Compliance |
|-------------|------------|
| E-invoice for cash sales > 200K VND | NOT IMPLEMENTED |
| Cash payment threshold for VAT deduction | NOT IMPLEMENTED |
| POS integration for cash receipts | NOT IMPLEMENTED |

---

## ERP Comparable Analysis

### MISA AMIS "Tiền mặt" Module
- **Full cash receipt/payment voucher management**
- **Cash book (Sổ quỹ tiền mặt) with real-time balance**
- **Bank transfer (rút/gửi tiền NH về nhập quỹ)**
- **Petty cash (tạm ứng)**
- **Cash flow forecasting**
- **Cash inventory (kiểm kê quỹ)**
- **Integration with AR/AP modules**
- **E-invoice for cash transactions**

### Fast Accounting "Tiền mặt, tiền gửi, tiền vay"
- **Cash/bank/loan management**
- **Receipt/payment vouchers**
- **Supplier payment processing**
- **Cash deposit/withdrawal**

### Bravo ERP "Quản lý vốn bằng tiền"
- **Multi-currency cash management**
- **Cash flow planning**
- **Receipt/payment control**
- **Bank integration**

### Tryton "Cash Management"
- **Cash journals**
- **Cash statements**
- **Petty cash**
- **Bank statement import**

---

## Build Estimate

| Phase | Components | Est. Effort |
|-------|-----------|-------------|
| P0 | Domain models, Repository interface, Cash API CRUD, Cash book report | 2 weeks |
| P1 | Petty cash, Cash inventory, Bank reconciliation | 1 week |
| P2 | Cash flow forecast, Cash forecast, Denomination tracking | 1 week |
| P3 | Bank integration (EBICS, API), POS integration | 2 weeks |
| **Total** | | **~6 weeks** |

---

## Recommendation

**DO NOT DEPLOY.** Cash module at 0% — no code, no schema, no API.

Build P0 first: domain models + receipt/payment CRUD + cash book report + B03-DN statement.
