# Purchase Module — Use Cases

**Version:** 2.0
**Date:** July 2026
**Author:** BA Lead + Chief Accountant (20+ yrs each)
**Note:** Handler tests cover UC-1 through UC-6 happy paths. See `internal/handler/purchase_handler_test.go`.

---

## UC-1: Create Supplier

| Element | Detail |
|---------|--------|
| **Actor** | AP Accountant, Purchasing Manager |
| **Precondition** | User authenticated, company selected |
| **Trigger** | New supplier needs to be added to master list |

### Happy Path
1. User fills: code, name, tax code (validated format 10/13/14 digits), address, phone, email
2. User fills payment terms (net30), bank account
3. System validates tax code against GDT database (optional)
4. System saves supplier as active
5. System returns supplier detail with generated ID

### Alternative Paths
- A1: Supplier already exists with same tax code → warn, allow duplicate if different company
- A2: Tax code validation fails → allow save with unverified flag

### Exception Paths
- E1: Required fields missing → validation error
- E2: Tax code format invalid → format error
- E3: Bank account invalid format → warning, allow save

---

## UC-2: Create Purchase Order

| Element | Detail |
|---------|--------|
| **Actor** | Purchasing Manager |
| **Precondition** | Supplier exists, items known |
| **Trigger** | Need to procure goods/services |

### Happy Path
1. User selects supplier
2. User enters one or more line items: description, unit, qty, unit price, VAT rate, GL account
3. System calculates subtotal, VAT, total
4. User enters delivery date, terms
5. User submits PO
6. System validates: supplier active, quantities > 0, prices > 0
7. System generates PO number (PO-YYYYMM-XXXX)
8. System saves PO as draft
9. PO sent for approval (if configured)
10. PO moves to approved → sent status

### Alternative Paths
- A1: User creates PO from purchase requisition → pre-filled from requisition
- A2: PO approval value-based → <10M auto, <100M manager approval, >100M director
- A3: User edits draft PO before approval → version not incremented
- A4: User cancels PO → reason required, inventory received qty checked
- A5: Partial order → user sets expected delivery per line

### Exception Paths
- E1: Supplier is suspended → cannot create PO
- E2: Item GL account incompatible with supplier type → warn
- E3: Budget exceeded (if budget enabled) → blocked or warning

---

## UC-3: Receive Goods Against PO

| Element | Detail |
|---------|--------|
| **Actor** | Warehouse Keeper |
| **Precondition** | PO approved and sent to supplier |
| **Trigger** | Goods arrive at warehouse |

### Happy Path
1. User selects PO
2. System displays PO lines with ordered qty, received qty, balance
3. User enters actual received qty per line
4. User enters warehouse location
5. User enters quality status (ok)
6. System checks: received ≤ ordered × tolerance (105% default)
7. User submits GRN
8. System saves GRN as posted
9. System updates PO status to partial (or received if complete)
10. System posts to inventory (152/153/156) with temporary price if uninvoiced

### Alternative Paths
- A1: Over-receipt within tolerance → allowed, warn
- A2: Over-receipt beyond tolerance → blocked, manager approval required
- A3: Quality reject → user enters rejected qty, system prompts return
- A4: Direct expense (service) → no inventory posting, post to expense account
- A5: Partial receipt → remaining balance open for future receipt

### Exception Paths
- E1: PO not found or cancelled → error
- E2: Quantity received negative → validation error
- E3: Multiple receipts exceed PO quantity → blocked

---

## UC-4: Record Supplier Invoice

| Element | Detail |
|---------|--------|
| **Actor** | AP Accountant |
| **Precondition** | Goods received (or service rendered) |
| **Trigger** | Supplier sends invoice (physical, email, or e-invoice via GDT) |

### Happy Path (3-way match)
1. User selects PO
2. System pre-fills from PO: supplier, items, expected qty, unit price
3. User enters: invoice number, invoice date, due date
4. User adjusts qty from GRN if partial
5. User verifies VAT rate and amounts
6. User submits invoice
7. System performs 3-way match: PO qty ≥ GRN qty ≥ Invoice qty
8. System checks: unit price variance < config threshold (default 5%)
9. Match OK → invoice status = verified
10. User posts invoice → GL update: Dr 152/156, Dr 1331, Cr 331
11. AP transaction created (increase AP)

### Happy Path (direct invoice, no PO)
1. User selects supplier
2. User enters line items manually: description, account, qty, price, VAT
3. Invoice status = draft
4. User submits → verified (no match needed)
5. User posts → GL update

### Alternative Paths
- A1: 3-way match partial (qty OK but price variance) → flag, require AP manager approval
- A2: 3-way match fails (qty mismatch) → hold, require investigation
- A3: Invoice is credit note → flag as negative, link to original invoice
- A4: Invoice is adjustment (price change only) → adjust existing AP balance
- A5: VAT not deductible → Dr full amount to expense/inventory, no 133
- A6: Prepayment exists → offset prepayment against invoice

### Exception Paths
- E1: Duplicate invoice number from same supplier → warn, block if exact match
- E2: Invoice date before last PO close → warn
- E3: Supplier not found → require create supplier first
- E4: Invoice total mismatch (line sum ≠ header total) → flag, require manual override
- E5: E-invoice XML parse error → reject, log error details

---

## UC-5: Pay Supplier

| Element | Detail |
|---------|--------|
| **Actor** | AP Accountant, Treasurer |
| **Precondition** | Invoice posted, AP balance due |
| **Trigger** | Payment due or early payment for discount |

### Happy Path
1. User selects supplier
2. System displays open invoices with amounts, due dates, aging
3. User selects invoices to pay (one or more)
4. User enters: payment date, payment amount, bank account
5. User enters reference (payment order number)
6. System allocates payment to selected invoices (FIFO or user-specified)
7. If early payment discount available → calculate and apply
8. System posts payment: Dr 331, Cr 112
9. AP transaction created (decrease AP)
10. Invoice status updated to paid (or partial)

### Alternative Paths
- A1: Partial payment of invoice → remaining balance tracked for next payment
- A2: Prepayment (advance) → no invoice, Dr 331 (debit balance)
- A3: Offset prepayment against invoice → auto-match
- A4: Multi-currency payment → FX gain/loss calculation to 515/635

### Exception Paths
- E1: Insufficient bank balance → warn, allow if user confirms
- E2: Payment > invoice balance → block
- E3: Supplier bank account missing → warn
- E4: Invoice already paid in full → block

---

## UC-6: Return Goods to Supplier

| Element | Detail |
|---------|--------|
| **Actor** | Warehouse Keeper, AP Accountant |
| **Precondition** | Goods received, possibly invoiced |
| **Trigger** | Quality issue, wrong item, excess quantity |

### Happy Path (pre-invoice return)
1. User selects GRN
2. User enters items/quantities to return
3. User selects return reason
4. System creates return GRN (negative)
5. System reverses inventory: Dr 331 (temp), Cr 152/156
6. PO status adjusted

### Alternative Paths
- A1: Post-invoice return → must receive credit note from supplier first
- A2: Credit note received → reduce AP balance
- A3: Return with replacement → track replacement PO

### Exception Paths
- E1: Goods consumed/used → cannot return, must sell or write off
- E2: Return quantity > received quantity → block

---

## UC-7: Generate AP Aging Report

| Element | Detail |
|---------|--------|
| **Actor** | AP Accountant, CFO |
| **Precondition** | Supplier invoices posted |
| **Trigger** | Periodic (month-end) or on demand |

### Happy Path
1. User selects report date
2. Optionally filters by supplier, currency
3. System calculates: current (not due), 1-30, 31-60, 61-90, 91-120, 120+ days overdue
4. System groups by supplier
5. System shows: invoice#, date, due date, days overdue, amount, balance
6. System shows totals per aging bucket
7. User can drill down to invoice detail

### Exception Paths
- E1: No AP data → empty report
- E2: Date range invalid → validation error

---

## UC-8: Receive E-Invoice from GDT

| Element | Detail |
|---------|--------|
| **Actor** | GDT System (automated API call) |
| **Precondition** | Supplier issued e-invoice on GDT portal, Decree 254/2026 |
| **Trigger** | Supplier submits e-invoice → GDT → pushes to buyer |

### Happy Path
1. GDT sends XML payload to GoTax webhook
2. System parses XML: validate structure, extract fields
3. System validates: supplier tax code exists in suppliers
4. System matches invoice to PO (by item description or manual link)
5. System creates SupplierInvoice in draft status
6. System stores raw XML in e_invoice_data JSONB
7. System sends notification to AP accountant: "New e-invoice received"
8. AP accountant reviews, verifies, posts

### Alternative Paths
- A1: Supplier not in master → create as "auto-created" supplier, flag for review
- A2: No PO match → create as direct invoice
- A3: Duplicate invoice detection → same supplier + invoice# + amount = duplicate

### Exception Paths
- E1: XML parse error → log error, notify admin
- E2: XML signature invalid → reject, flag for investigation
- E3: GDT service unavailable → queue for retry

---

## UC-9: Track VAT Input Deduction

| Element | Detail |
|---------|--------|
| **Actor** | Tax Accountant |
| **Precondition** | Supplier invoice with VAT posted |
| **Trigger** | Monthly VAT declaration preparation |

### Happy Path
1. User views VAT input tracking report
2. System shows all invoices with VAT, by period
3. Status: pending (unclaimed), claimed (submitted in VAT return), rejected (by GDT)
4. User marks invoices as claimed in this period
5. System links to VAT declaration module

### Exception Paths
- E1: Invoice VAT > declared revenue × 10% → flag for GDT audit risk
- E2: Non-deductible VAT (e.g., personal use) → mark as non-deductible
- E3: VAT rate mismatch (PO said 8%, invoice says 10%) → flag for review