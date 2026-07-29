# Sale Module — Business Rules & Compliance

**Version:** 1.0
**Date:** July 2026
**Author:** BA Lead + Chief Accountant (20+ yrs each)
**Regulatory Basis:** Circular 99/2025/TT-BTC, Decree 123/2020/ND-CP, Decree 70/2025, Decree 254/2026, IFRS 15, VAS 14, VAS 08, Law on VAT 13/2024/QH15

---

## 1. Accounting Rules

### R1: Revenue Recognition

**Rule:** Revenue recognized when control of goods/services transfers to customer. Per IFRS 15 5-step model: (1) identify contract, (2) identify performance obligations, (3) determine transaction price, (4) allocate price, (5) recognize revenue when obligation satisfied.

**Source:** IFRS 15.9, VAS 14, Circular 99 Account 511 guidance

**Implementation:**
- Revenue recognized on delivery (DN posted) for goods
- Revenue recognized on service completion or % completion for services
- Deferred revenue (3387) if payment received before obligation satisfied
- `revenue_recognition_date` field on invoice

### R2: Revenue Account Mapping

**Rule:** Revenue recorded to correct sub-account based on transaction type:
- 5111: Revenue from goods sales (ban hang hoa)
- 5112: Revenue from services (cung cap dich vu)
- 5113: Revenue from leasing (cho thue tai san)
- 5114: Revenue from construction (hop dong xay dung)
- 5117: Other revenue (doanh thu khac)
- 5118: Revenue from internal sales (ban noi bo)

**Source:** Circular 99/2025/TT-BTC, Account 511 guidance

**Implementation:**
- `revenue_account_id` per line item configurable
- Auto-mapping by item type if not specified
- Segregation by product category

### R3: VAT Output Tracking

**Rule:** Every taxable invoice must record:
1. VAT rate applied (0%, 5%, 8%, 10%, KCT, NT)
2. VAT amount per rate
3. E-invoice code from GDT
4. Signature verification status

VAT output must match total VAT collected from customers.

**Source:** Law on VAT 13/2024/QH15, Decree 123/2020 Art 48

**Implementation:**
- `vat_type` per invoice line
- VAT output aggregated by rate for tax return
- E-invoice code mandatory for all issued invoices
- VAT output account 3331 segregated by rate (33311/33312)

### R4: E-Invoice Mandatory Issuance

**Rule:** E-invoice must be issued within 24 hours of delivery (goods) or payment receipt (services). E-invoice must be signed by registered digital certificate and submitted to GDT for coding.

**Source:** Decree 123/2020 Art 48, Decree 254/2026

**Implementation:**
- Auto-trigger e-invoice pipeline on invoice post
- SLA monitoring: time from delivery to invoice issuance
- Penalty warning if >24h threshold approached
- Digital signature validation before GDT submission

### R5: Credit Note / Sales Return

**Rule:** Credit note issued when:
1. Goods returned (full or partial)
2. Price adjustment after invoice
3. Discount/allowance granted after invoice
4. Invoice cancellation per regulatory process

Credit note must reference original invoice number + GDT invoice code.

**Source:** Decree 123/2020 Art 56, VAS 08

**Implementation:**
- Credit note lines reference original invoice lines
- Credit note revenue reversal matches original revenue account
- VAT reversal only if original VAT was charged
- Credit note e-invoice submitted to GDT (negative amount)

### R6: Corrective Invoice Rules

**Rule:** Three methods per Decree 70/2025:
1. **Adjustment increase**: New invoice with additional amount, referencing original
2. **Adjustment decrease**: New invoice with negative differential, referencing original
3. **Replacement**: Cancel original, issue new invoice with original reference

**Source:** Decree 70/2025/ND-CP, Decree 254/2026

**Implementation:**
- `adjustment_type` field: increase, decrease, replacement
- `original_invoice_id` reference for all correction types
- Full audit trail: original → correction → reason
- GL reversal + re-posting for replacement

### R7: Trade Discount and Allowance

**Rule:**
- Trade discount (before invoice): reduce invoice line price, not separate entry
- Commercial discount (5211): post-sale discount granted
- Sales allowance (5212): price reduction after sale
- Early payment discount (5213): discount for prompt payment
- All discounts tracked per customer for reporting

**Source:** Circular 99, Account 521 guidance, IFRS 15.46-53

**Implementation:**
- Pre-invoice discount → reduced revenue on invoice
- Post-invoice discount/allowance → separate credit note or adjustment
- Early payment discount → recorded at receipt time (5213)
- Discount GL: 5211/5212/5213 as contra-revenue accounts

### R8: AR Aging and Bad Debt

**Rule:** AR classified by aging buckets. Bad debt provision (Account 229) based on aging:
- < 6 months: 0% provision
- 6-12 months: 30% provision
- 1-2 years: 50% provision
- 2-3 years: 70% provision
- 3+ years: 100% provision

Bad debt write-off: Dr 642 (expense) Cr 131 after all collection attempts exhausted.

**Source:** Circular 99, Account 229/642 guidance, VAS 17

**Implementation:**
- Auto-calculate aging buckets
- Provision calculation per customer per invoice
- Write-off requires director approval (segregation of duty)
- Written-off AR tracked off-balance-sheet for 10 years

### R9: Foreign Currency AR

**Rule:** AR in foreign currency recorded at transaction date rate. At period-end, revalue at month-end rate. FX gain/loss to 515 or 635.

**Source:** VAS 10 (Foreign Exchange), Circular 99 Account 515/635

**Implementation:**
- `currency` + `exchange_rate` on invoice and receipt
- Month-end auto-revaluation engine
- Realized FX on receipt allocation (rate at receipt vs rate at invoice)
- Unrealized FX on period-end revaluation

### R10: Prepayment/Deposit

**Rule:** Customer prepayment recorded as credit balance on AR (Account 131). If deposit e-invoice issued per Decree 254, VAT must be declared at deposit time. Final invoice issued for remaining balance, referencing deposit invoice.

**Source:** Decree 123/2020, Circular 99 Account 131 guidance

**Implementation:**
- Prepayment receipt: Dr 112 Cr 131 (credit balance)
- Optional deposit e-invoice if customer requires
- Prepayment offset: Dr 131 (credit bal) Cr 131 (AR invoice)
- Final invoice amount = total - deposit

---

## 2. Validation Rules

### R11: Customer Tax Code Validation

**Rule:** Customer tax code must be valid Vietnamese format (10 digits, 13 digits, or 13+1 digits for dependent units). Optionally validate against GDT database.

**Source:** Decree 123/2020

**Implementation:**
- Regex validation on save: `^\d{10}$|^\d{13}$|^\d{13}-\d{3}$`
- Optional GDT API lookup
- Warning (not block) if validation fails

### R12: Duplicate Invoice Prevention

**Rule:** No duplicate invoice for same SO + delivery combination. Same customer + items + amounts within 24h flagged as potential duplicate.

**Implementation:**
- Unique constraint on (so_id, dn_id) for invoicing
- Fuzzy duplicate detection on create

### R13: Over-Delivery Tolerance

**Rule:** Delivery qty may exceed SO qty up to configurable tolerance (default 5%). Above tolerance requires manager approval. Over-delivery proportion invoiced at same unit price.

**Source:** Standard business practice

**Implementation:**
- `over_delivery_tolerance_pct` company config
- Warn if exceeded, block if > double tolerance
- Manager override via PATCH

### R14: Over-Invoice Prevention

**Rule:** Invoice qty cannot exceed delivered qty (by line). Tolerance = 0% (strict match). Partial invoicing allowed (invoice less than delivered).

**Implementation:**
- `invoiced_qty ≤ delivered_qty` per line item
- Block if exceeded
- Allow partial billing

### R15: Credit Limit Enforcement

**Rule:** If customer credit limit is set, SO confirmation must check: total outstanding AR + SO total ≤ credit limit. Override requires manager approval.

**Source:** Risk management best practice

**Implementation:**
- Check at SO confirm step
- Configurable action: warn vs block
- Manager override with reason + audit trail

### R16: Segregation of Duties

**Rule:**
- Sales person who creates SO cannot approve same SO
- AR accountant who records receipt cannot reconcile bank statement
- Person who creates credit note cannot void credit note
- All invoice cancellations require supervisor approval

**Source:** Internal control best practice, Circular 99

**Implementation:**
- Role-based access on approval endpoints
- created_by ≠ approved_by validation
- Audit log for override actions

### R17: Price Integrity

**Rule:** Unit price on invoice must match unit price on SO (± configurable tolerance, default 2%). Variance above threshold requires manager approval.

**Source:** Standard business practice

**Implementation:**
- SO price cached at confirm time
- Invoice price validated against SO price
- Tolerance configurable per company

---

## 3. Period-End Rules

### R18: Revenue Accrual

**Rule:** If goods delivered in period but not yet invoiced, accrue revenue at period-end. Auto-reverse next period.

**Source:** Matching principle (IFRS 15)

**Implementation:**
- "Unbilled Deliveries Report" identifies candidates
- Auto-create accrual journal entry
- Auto-reverse on first day of next period

### R19: Month-End AR Reconciliation

**Rule:** AR sub-ledger (by customer by invoice) must equal GL Account 131 balance. Any variance must be investigated before period close.

**Source:** Circular 99

**Implementation:**
- Reconciliation report: sub-ledger total vs GL balance
- Variance items listed for investigation
- Period close blocked if variance > threshold

### R20: VAT Output vs GDT Reconciliation

**Rule:** Total VAT output from sales invoices must match GDT-registered e-invoice VAT total. Discrepancy = penalty risk.

**Source:** Decree 123/2020

**Implementation:**
- Monthly reconciliation report
- Compare GoTax VAT output vs GDT portal data
- Flag unmatched invoices

---

## 4. Compliance Rules

### R21: E-Invoice Timeliness SLA

**Rule:** E-invoice issued within 24h of delivery. Monitoring dashboard shows compliance rate.

**Source:** Decree 123/2020 Art 48

**Implementation:**
- `delivery_time` → `invoice_time` tracking
- Dashboard: compliance rate, breached invoices
- Warning at 18h threshold

### R22: Audit Trail

**Rule:** Every state change on SO, DN, invoice, receipt, credit note logged with: timestamp, user_id, action, old_status, new_status, reason.

**Source:** Accounting Law 2015

**Implementation:**
- Use existing AuditEntry domain model
- Auto-log on every status transition

### R23: Data Retention

**Rule:** All sales documents retained for minimum 10 years per Accounting Law. E-invoice XML retained in original format for legal validity.

**Source:** Accounting Law 2015 Art 13

**Implementation:**
- Soft-delete only (no hard delete)
- E-invoice raw XML stored in DB
- Archival process for documents > 10 years
