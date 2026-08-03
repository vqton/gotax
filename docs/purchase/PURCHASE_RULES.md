# Purchase Module — Business Rules & Compliance

**Version:** 2.0
**Date:** July 2026
**Author:** BA Lead + Chief Accountant (20+ yrs each)
**Regulatory Basis:** Circular 99/2025/TT-BTC, Decree 123/2020/ND-CP, Decree 70/2025, Decree 254/2026, IAS 2, VAS 02, VAS 17
**Note:** Accounting rules implemented partially. Cost allocation engine pending (defined in domain).

---

## 1. Accounting Rules

### R1: Purchase Cost Recognition

**Rule:** Purchase cost = purchase price + import duties + non-reclaimable taxes + transport + handling + other directly attributable costs minus trade discounts and rebates (IAS 2.11, Circular 99 - Account 152/156 guidance)

**Source:** IAS 2.11, Circular 99/2025/TT-BTC

**Implementation:**
- `purchase_cost` field for landed cost tracking
- Cost allocation engine for import (by_qty, by_value, by_weight, by_volume)

### R2: Inventory Valuation Methods

**Rule:** Support 4 methods per Circular 99 (added Standard Cost to existing Specific ID, Weighted Average, FIFO). LIFO prohibited per IAS 2.

**Source:** Circular 99/2025/TT-BTC, IAS 2.23-25

**Implementation:**
- `cost_method` config per item or company
- Standard Cost: record variance to COGS
- Periodic cost recalculation engine

### R3: Goods in Transit

**Rule:** If goods invoiced but not received by period-end → record to Account 151 (Goods in transit). Next period → reverse 151 and record to 152/156.

**Source:** Circular 99, Account 151 guidance

**Implementation:**
- Auto-detect end-of-period uninvoiced receipts
- Create 151 entry automatically
- Auto-reverse next period

### R4: Temporary Price (Goods Received Not Invoiced)

**Rule:** If GRN posted but no invoice yet → record at PO price (temporary). When invoice arrives → adjust to actual price. Circular 99 allows accrual.

**Source:** Circular 99/2025/TT-BTC

**Implementation:**
- `temp_price` flag on GRN posting
- Price adjustment GL entry when invoice posted
- Variance tracked and reported

### R5: Prepayment Tracking

**Rule:** Prepayment to supplier → Dr 331 (debit balance = advance). Upon invoice receipt → offset prepayment against AP. Perpetual tracking per supplier.

**Source:** Circular 99, Account 331 guidance

**Implementation:**
- AP sub-ledger tracks invoice-by-invoice balance
- Prepayment auto-allocated per user preference (FIFO or specific)

### R6: VAT Input Deduction

**Rule:** VAT input deductible only if:
1. Valid e-invoice from supplier (Decree 123 Art 56)
2. Goods/services used for taxable business activities
3. Payment via bank for invoices > 20M VND (Decree 123 requirements)
4. VAT invoice matches PO VAT rate

**Source:** Law on VAT 13/2024/QH15, Decree 123/2020, Circular 99 Account 1331

**Implementation:**
- `vat_deduction_status` per invoice: pending, claimed, rejected
- Invoice > 20M → require bank payment evidence
- VAT rate matching validation

### R7: Discount and Rebate Treatment

**Rule:**
- Trade discount: reduce purchase cost (deduct from 152/156)
- Early payment discount: record as financial income (515)
- Volume rebate (contractual): accrue and reduce cost

**Source:** IAS 2.11, Circular 99 (Account 515)

**Implementation:**
- Separate discount type fields on invoice
- GL: trade discount → reduce inventory cost; payment discount → 515

### R8: Import Purchase Accounting

**Rule:** Import cost = CIF value + import duty (3333) + special consumption tax (3332) + environmental tax + import VAT (33312). All non-reclaimable taxes → add to inventory cost.

**Source:** Circular 99, Accounts 3332, 3333, 33381, 33312

**Implementation:**
- Import purchase flow with separate duty/VAT lines
- `landed_cost` breakdown fields

### R9: Purchase Return

**Rule:** Return goods → reverse purchase:
- If posted but uninvoiced: reverse GRN
- If posted and invoiced: must receive credit note → reverse invoice and adjust AP

**Source:** Decree 70/2025 (amending Decree 123), Circular 99

**Implementation:**
- Return GRN (negative receipt)
- Credit note invoice (negative amount)
- GL reversal entries

---

## 2. Validation Rules

### R10: PO Approval Thresholds

| PO Total (VND) | Required Approval | Auto-Approve |
|----------------|-------------------|--------------|
| < 10,000,000 | None | Yes |
| 10,000,000 - 100,000,000 | Purchasing Manager | No |
| 100,000,000 - 1,000,000,000 | + Chief Accountant | No |
| > 1,000,000,000 | + CFO/Director | No |

**Source:** Enterprise internal control policy (configurable per company)

### R11: 3-Way Match Tolerances

| Tolerance Type | Default | Configurable |
|----------------|---------|--------------|
| Quantity over-receipt | 5% of PO qty | Yes |
| Unit price variance | 5% of PO price | Yes |
| Tax amount variance | Absolute: 10,000 VND | Yes |

**Source:** Industry practice (MISA, Fast, Bravo consistency)

### R12: Supplier Tax Code Validation

- Format: 10 digits (standard) or 13 digits (subsidiary) or 14 digits (new format)
- Validate against GDT tax code database (optional, requires API)
- Cannot have 2 suppliers with same tax code in same company (configurable)

**Source:** Law on Tax Administration 108/2025/QH15

### R13: Invoice Number Uniqueness

- Per supplier × invoice number → unique constraint
- Exception: corrective/adjustment invoices with same number but different date
- Detection: same supplier + same invoice# + same amount = likely duplicate

**Source:** Decree 123/2020

---

## 3. Security & Audit Rules

### R14: Audit Trail

| Event | Fields Logged |
|-------|---------------|
| PO created | user, timestamp, initial state |
| PO approved | user, timestamp, previous state→new state |
| PO cancelled | user, timestamp, reason |
| GRN posted | user, timestamp, PO reference |
| Invoice posted | user, timestamp, GL entries |
| Invoice paid | user, timestamp, payment reference |
| Any edit | user, timestamp, field old→new |

**Source:** Accounting Law Article 12

### R15: Access Control

| Role | Permissions |
|------|------------|
| AP Clerk | View suppliers, create/edit draft POs, record GRN, record invoices (draft), view reports |
| AP Manager | Approve POs, verify invoices, approve payments, all AP Clerk permissions |
| Chief Accountant | Post invoices to GL, adjust AP, period-end closing, override blocked transactions |
| CFO | Approve large POs, approve large payments, view all reports |
| Auditor | Read-only access to all purchase data + audit log |
| Warehouse Keeper | Create GRN (post), view PO (read only) |

**Source:** Internal control standards

---

## 4. Period-End Rules

### R16: Period Locking

- Closed period → no new postings allowed
- Exception: only Chief Accountant can reverse-lock for corrections
- PO/GRN dates validated against open period calendar

**Source:** Circular 99, Period management

### R17: AP Sub-ledger = GL Reconciliation

- Monthly: sum(AP sub-ledger balance by supplier) = GL 331 balance
- Any variance > configurable threshold → flagged for investigation
- Reconciliation report required before period close

**Source:** Accounting Law, audit requirements

---

## 5. Compliance Rules (Decree 254/2026 E-Invoice)

### R18: E-Invoice Receipt Requirements

- System must accept XML format per GDT standard (Decree 254, effective 01/07/2026)
- Validate digital signature (Decree 23/2025)
- Store raw XML immutably
- Provide QR code lookup for buyer verification

### R19: Invoice Archive Retention

- Purchase invoices: 10 years (Accounting Law)
- E-invoice XML: 10 years in original format
- Archived in compressed + encrypted storage
- Monthly backup to separate geographic location

---

## 6. Multi-Company Rules

### R20: Inter-Company Purchase

- Company A sells to Company B (same group)
- Must have valid contract and e-invoice
- Eliminated on consolidation per Circular 99

### R21: Data Isolation

- Supplier data per company_id
- No cross-company supplier visibility
- Reports filtered by company context