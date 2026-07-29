# Sale Module — Use Cases

**Version:** 1.0
**Date:** July 2026
**Author:** BA Lead + Chief Accountant (20+ yrs each)

---

## UC-1: Create Customer

| Element | Detail |
|---------|--------|
| **Actor** | Sales Manager, AR Accountant |
| **Precondition** | User authenticated, company selected |
| **Trigger** | New customer needs to be added to master list |

### Happy Path
1. User fills: code, name, tax code (validated format 10/13/14 digits), address, phone, email
2. User fills payment terms (net30), credit limit, bank account
3. User assigns customer type (domestic/export) and group (retail/wholesale)
4. System validates tax code format
5. System checks for duplicate tax code within company
6. System saves customer as active
7. System returns customer detail with generated ID

### Alternative Paths
- A1: Customer already exists with same tax code → warn, allow if different customer name
- A2: Credit limit set to 0 → interpreted as no limit

### Exception Paths
- E1: Required fields missing → validation error
- E2: Tax code format invalid → format error
- E3: Duplicate customer code within company → error, suggest next available code

---

## UC-2: Create Sales Order

| Element | Detail |
|---------|--------|
| **Actor** | Sales Manager |
| **Precondition** | Customer exists, items/services known |
| **Trigger** | Customer places order |

### Happy Path
1. User selects customer
2. User enters one or more line items: description, unit, qty, unit price, discount %, VAT rate, revenue account
3. System calculates subtotal, VAT, total
4. User enters expected delivery date, delivery terms
5. User submits SO
6. System validates: customer active, quantities > 0, prices > 0
7. System generates SO number (SO-YYYYMM-XXXX)
8. System saves SO as draft
9. SO sent for approval workflow (if configured)
10. SO moves to approved → confirmed status

### Alternative Paths
- A1: User creates SO from quotation → pre-filled from quotation
- A2: SO approval value-based → <10M auto, <100M manager, >100M director
- A3: User edits draft SO before approval → version not incremented
- A4: User cancels SO → reason required, delivered qty checked
- A5: Customer has credit limit → credit check on confirm
- A6: Inventory check on confirm → warn if insufficient stock (P1)

### Exception Paths
- E1: Customer is suspended → cannot create SO
- E2: Revenue account incompatible with customer type → warn
- E3: Credit limit exceeded → blocked (or manager override)
- E4: VAT rate incompatible with customer type (export = 0%) → auto-set

---

## UC-3: Deliver Goods Against SO

| Element | Detail |
|---------|--------|
| **Actor** | Warehouse Keeper |
| **Precondition** | SO approved and confirmed |
| **Trigger** | Goods ready for delivery to customer |

### Happy Path
1. User selects SO
2. System displays SO lines with ordered qty, delivered qty, balance
3. User enters actual delivered qty per line
4. User enters warehouse location, delivery address, shipping method
5. System checks: delivered ≤ ordered × tolerance (105% default)
6. User submits delivery note
7. System saves DN as posted
8. System updates SO status to processing (or delivered if complete)
9. System posts COGS: Dr 632 Cr 152/156
10. System reduces inventory (152/156)

### Alternative Paths
- A1: Over-delivery within tolerance → allowed, warn
- A2: Over-delivery beyond tolerance → blocked, manager approval required
- A3: Service delivery → no inventory posting, post directly
- A4: Partial delivery → remaining balance open for future delivery
- A5: Direct drop-ship (supplier → customer) → special handling

### Exception Paths
- E1: SO not found or cancelled → error
- E2: Quantity delivered negative → validation error
- E3: Multiple deliveries exceed SO quantity → blocked

---

## UC-4: Issue Customer Invoice

| Element | Detail |
|---------|--------|
| **Actor** | AR Accountant |
| **Precondition** | Delivery note posted (goods) or order confirmed (services) |
| **Trigger** | Need to bill customer |

### Happy Path
1. User selects SO (optionally specific DN lines)
2. System displays deliverable lines with delivered qty, already invoiced qty, balance
3. User selects lines and quantities to invoice
4. System calculates amounts, VAT
5. User enters invoice notes (printed on invoice)
6. User reviews invoice preview
7. User submits invoice
8. System saves invoice as draft
9. System generates TXML e-invoice XML
10. System signs XML with registered digital certificate
11. System submits to GDT portal
12. GDT returns invoice code
13. System updates invoice as issued
14. System posts GL: Dr 131 Cr 5111 + Cr 3331
15. System updates SO invoiced qty

### Alternative Paths
- A1: Consolidated invoice → user selects multiple DNs against same customer
- A2: Invoice for services (no delivery) → invoice directly from SO
- A3: GDT submission fails → system saves error, user can retry
- A4: Invoice with prepayment offset → system shows prepayment balance, auto-offset
- A5: Export invoice → VAT 0%, no VAT output posting

### Exception Paths
- E1: No deliverable lines found → error
- E2: GDT service unavailable → queue for retry, notify admin
- E3: Digital signature expired → require new signature before issuance
- E4: Invoice qty > delivered qty → blocked

---

## UC-5: Record Customer Payment

| Element | Detail |
|---------|--------|
| **Actor** | AR Accountant, Treasurer |
| **Precondition** | Invoice issued (or prepayment situation) |
| **Trigger** | Customer sends payment |

### Happy Path
1. User selects customer
2. System displays outstanding invoices with balance due
3. User enters payment: amount, date, payment method (bank/cash/cheque), reference
4. User allocates payment to specific invoice(s)
5. System checks: allocated total ≤ payment amount
6. User submits receipt
7. System saves receipt as posted
8. System posts GL: Dr 112/111 Cr 131
9. System updates invoice: amount_received, balance_due
10. If balance_due = 0 → mark invoice as paid
11. System adjusts AR total for customer

### Alternative Paths
- A1: Partial payment → partially reduce invoice balance, leave remainder
- A2: Overpayment → create credit balance on customer AR
- A3: Prepayment (no invoice yet) → record as deposit (credit AR)
- A4: Early payment with discount → allocate discount to 5213
- A5: Foreign currency payment → record FX gain/loss (515/635)
- A6: Bank transfer with fee → handle collection fee separately

### Exception Paths
- E1: Payment amount negative → validation error
- E2: Customer not found → error
- E3: Duplicate payment reference → warn, allow if different amount
- E4: Payment allocated to already-paid invoice → error

---

## UC-6: Process Sales Return / Credit Note

| Element | Detail |
|---------|--------|
| **Actor** | AR Accountant, Warehouse Keeper |
| **Precondition** | Invoice issued and posted |
| **Trigger** | Customer returns goods or requests adjustment |

### Happy Path
1. User selects original invoice
2. System displays invoice lines
3. User selects lines and quantities being returned
4. User enters return reason (defect/wrong item/excess/other)
5. User enters return type: full return, partial return, price adjustment
6. If goods returned: user enters returned qty and warehouse location
7. System calculates credit amounts (revenue reversal + VAT reversal)
8. User submits credit note
9. System generates negative e-invoice TXML
10. System signs + submits to GDT
11. GDT returns credit note code
12. System posts credit note to GL: Dr 5111 Cr 131, Dr 3331 Cr 131
13. If goods returned: Dr 152/156 Cr 632 (COGS reversal)
14. System updates original invoice: credit_note_amount, balance_due

### Alternative Paths
- A1: Price adjustment (no return) → only revenue/VAT adjustment, no inventory reversal
- A2: Return with replacement → credit note + new delivery note
- A3: Customer has no more outstanding balance → refund required

### Exception Paths
- E1: Original invoice not found → error
- E2: Return qty > original invoice qty → blocked
- E3: Credit note for already-credited line → blocked
- E4: GDT rejects credit note → fix and resubmit

---

## UC-7: Manage AR Aging

| Element | Detail |
|---------|--------|
| **Actor** | AR Accountant, CFO |
| **Precondition** | Invoices exist and are posted |
| **Trigger** | Periodic review of outstanding AR |

### Happy Path
1. User requests AR aging report
2. System calculates for each invoice: balance = total - sum(payments)
3. System computes days overdue = today - due_date
4. System assigns aging bucket: current, 1-30, 31-60, 61-90, 91-120, 120+
5. System groups by customer
6. System returns aged trial balance
7. User reviews overdue items and initiates collection actions

### Alternative Paths
- A1: Filter by customer → single customer aging
- A2: Filter by aging bucket → focus on critical items
- A3: Export to Excel for further analysis

### Exception Paths
- E1: No data → empty report

---

## UC-8: Month-End Closing for Sales

| Element | Detail |
|---------|--------|
| **Actor** | Chief Accountant, AR Accountant |
| **Precondition** | All period transactions entered |
| **Trigger** | Period close |

### Happy Path
1. Run unbilled deliveries report → accrue revenue if needed
2. Reconcile AR sub-ledger vs GL 131 balance
3. Reconcile VAT output vs GDT portal data
4. Revalue foreign currency AR at month-end rate
5. Review AR aging and calculate bad debt provision
6. Verify all e-invoices have valid GDT codes
7. Close period for sales (lock further posting)

### Alternative Paths
- A1: Discrepancy found → investigate, adjust, re-run
- A2: Unbilled deliveries → auto-create accrual journal

### Exception Paths
- E1: AR sub-ledger ≠ GL 131 → cannot close period
- E2: VAT output mismatch with GDT → flag and resolve
