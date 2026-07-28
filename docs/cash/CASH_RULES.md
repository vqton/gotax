# Cash Module — Business Rules

**Sources:** Circular 99/2025/TT-BTC, Law on Accounting 2015, Decree 123/2020/ND-CP, VAS

---

## 1. Transaction Rules

### R1: Receipt Voucher (Phiếu thu — Mẫu 01-TT)

| Rule | Description |
|------|-------------|
| R1.1 | Every cash receipt must have Phiếu thu with unique sequential number |
| R1.2 | Phiếu thu requires: date, amount (words + digits), payer, reason, account coding (Nợ/Có) |
| R1.3 | Signatures required: Người lập phiếu, Kế toán trưởng, Thủ quỹ, Người nộp tiền, Giám đốc |
| R1.4 | Cash amount > 20M VND requires Giám đốc signature (per company policy, configurable) |
| R1.5 | Posting updates cash book + GL simultaneously |

### R2: Payment Voucher (Phiếu chi — Mẫu 02-TT)

| Rule | Description |
|------|-------------|
| R2.1 | Every cash payment must have Phiếu chi with unique sequential number |
| R2.2 | Phiếu chi requires: date, amount (words + digits), payee, reason, account coding |
| R2.3 | Signatures: Người lập phiếu, Kế toán trưởng, Thủ quỹ, Người nhận tiền, Giám đốc |
| R2.4 | Cannot post payment > available cash balance (configurable override) |
| R2.5 | Payment to related parties > 100M VND requires special approval (Decree 123) |

### R3: Cash Book (Sổ quỹ tiền mặt — S07-DN / S04a-DNN)

| Rule | Description |
|------|-------------|
| R3.1 | Cash book per currency per account |
| R3.2 | Running balance after every transaction |
| R3.3 | Daily closing balance must equal physical cash count |
| R3.4 | Cashier and accountant maintain parallel books, periodic cross-check |
| R3.5 | Monthly closing: print, sign, archive |

### R4: Cash Inventory (Biên bản kiểm kê quỹ — Mẫu 08a-TT)

| Rule | Description |
|------|-------------|
| R4.1 | Mandatory at least monthly (Circular 99 Art. 32) |
| R4.2 | Unannounced inventory permitted |
| R4.3 | Inventory committee: Kế toán trưởng + Thủ quỹ + independent witness |
| R4.4 | Count by denomination (VND: 500K, 200K, 100K, 50K, 20K, 10K, 5K, 2K, 1K) |
| R4.5 | Excess: Debit 111, Credit 338 (chờ xử lý) |
| R4.6 | Shortage: Debit 138 (cá nhân bồi thường) or 334 (trừ lương), Credit 111 |
| R4.7 | Signed by all committee members, filed for audit |

### R5: Dual Custody

| Rule | Description |
|------|-------------|
| R5.1 | Cashier (Thủ quỹ) holds physical cash — sole custodian |
| R5.2 | Accountant records cash transactions — no physical access |
| R5.3 | Cashier cannot create vouchers; Accountant cannot access cash |
| R5.4 | Violation: immediate escalation |

### R6: Approval Tiers

| Threshold | Approver |
|-----------|----------|
| < 10M VND | Kế toán trưởng (Chief Accountant) |
| 10M - 100M VND | Kế toán trưởng + Giám đốc |
| > 100M VND | Kế toán trưởng + Giám đốc + HĐQT (if required) |
| *Configurable per company policy* | |

---

## 2. Legal Rules

### L1: Law on Accounting 2015 (Art. 12-17)

| Rule | Description |
|------|-------------|
| L1.1 | All cash transactions documented with original vouchers |
| L1.2 | Electronic vouchers = same legal standing as paper |
| L1.3 | Voucher retention: minimum 10 years |
| L1.4 | Accounting books: printed, bound, archived annually |

### L2: Circular 99/2025/TT-BTC Cash Provisions

| Article | Rule |
|---------|------|
| Art. 12 | TK 111: only actual cash receipt/payment recorded |
| Art. 12(1b) | Vouchers require full signatures per chế độ chứng từ |
| Art. 12(1c) | Cashier must maintain daily cash book |
| Art. 12(1d) | Periodic cash inventory mandatory |
| Art. 14 | Cash flow information in financial statements |
| Art. 52 | Foreign currency cash: period-end revaluation |

### L3: Decree 123/2020/ND-CP (E-Invoice)

| Rule | Description |
|------|-------------|
| L3.1 | Cash sales > 200K VND must issue e-invoice |
| L3.2 | E-invoice data flows to GDT via XML API |
| L3.3 | Invoice cancellation requires digital signature |

### L4: VAT Deduction Rules (Non-Cash Payment)

| Rule | Description |
|------|-------------|
| L4.1 | Input VAT deductible only if payment via non-cash method for transactions > 20M VND |
| L4.2 | Exception: wages, utilities, small-value purchases |
| L4.3 | Cash payment > 20M VND → VAT deduction disallowed |
| L4.4 | System must flag VAT deduction risk for cash payments > 20M VND |

---

## 3. Validation Rules

| Field | Rule | Error |
|-------|------|-------|
| Amount | > 0 | "Số tiền phải lớn hơn 0" |
| Voucher date | Not in future | "Ngày chứng từ không được lớn hơn ngày hiện tại" |
| Currency | In allowed list | "Đơn vị tiền tệ không hợp lệ" |
| Cash account | Must be 1111/1112 | "Tài khoản tiền mặt không đúng" |
| Debit = Credit | Must balance | "Tổng Nợ và Tổng Có không khớp" |
| Exchange rate | > 0 for FCY | "Tỷ giá phải lớn hơn 0" |
| Status transition | No skip | "Không thể chuyển trạng thái" |
| Approval | Different user | "Người duyệt phải khác người lập" |
| Voucher number | Unique per year | "Số chứng từ đã tồn tại" |

---

## 4. Periodic Rules

| Rule | Frequency | Action |
|------|-----------|--------|
| Cash inventory | Monthly minimum | Create inventory record |
| Cash book closing | Daily | Print, sign |
| Cash book print | End of period | PDF archive |
| GL reconciliation | Monthly | Match cash book vs GL 111 |
| Bank reconciliation | Monthly | Match cash book vs bank statement |
| FCY revaluation | Monthly/Quarterly | Adjust FX gain/loss |
| B03-DN preparation | Quarterly/Annual | Cash flow statement |
| Voucher archiving | Annual | Bind, seal, store 10 yrs |

---

## 5. Security Rules

| Rule | Description |
|------|-------------|
| S1 | Cashier role cannot create vouchers |
| S2 | Accountant cannot execute payments |
| S3 | Posted vouchers immutable — reversal creates new voucher |
| S4 | Audit log for every status change |
| S5 | Balance inquiry requires authentication |
| S6 | Bulk export requires Chief Accountant approval |
