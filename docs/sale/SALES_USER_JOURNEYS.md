# Sale Module — User Journeys

**Version:** 1.0
**Date:** July 2026
**Author:** BA Lead + Chief Accountant (20+ yrs each)

---

## Journey 1: Sales Manager — Create and Manage Sales Orders

**Role:** Sales Manager (phòng kinh doanh)
**Goal:** Process customer orders efficiently, ensure accurate pricing and delivery commitments

### Day-in-the-life

```
08:00 — Check email: customer ABC Corp sends PO for 500 units
  → Login to GoTax
  → Search customer "ABC Corp" (found, active, credit limit OK)
  → Create SO manually (no quotation)
  → Enter 5 line items from customer PO
  → System auto-prices from price list (tier: wholesale)
  → Validate totals match customer PO
  → Submit SO (status: draft)

08:15 — SO auto-routed for approval (value > 50M → manager)
  → Approve SO (status: approved → confirmed)
  → System checks: inventory available (50% stock, balance will be produced)
  → SO now in "processing" state

09:30 — Email customer: order confirmed, expected delivery 5 days

14:00 — Customer calls: changes qty on line 3 (300 → 200)
  → Open SO (status: processing)
  → Amend line 3 qty
  → System recalculates totals
  → Save amendment (version auto-incremented)
  → Email change confirmation to customer
```

**Pain points:**
- Cannot see real-time inventory (P1 feature)
- No auto-pricing from customer contract (P1)
- Price list management manual (P2)

---

## Journey 2: Warehouse Keeper — Process Deliveries

**Role:** Warehouse Keeper (thủ kho)
**Goal:** Pick, pack, and ship goods accurately and on time

### Day-in-the-life

```
07:30 — Check delivery schedule for today
  → Open "Pending Deliveries" list
  → 5 SOs scheduled for today
  → Sort by priority (delivery date oldest first)

07:45 — Process SO-202607-0123 (ABC Corp, 500 units)
  → Select SO → "Create Delivery"
  → System shows SO lines with ordered qty, previously delivered (0), balance
  → Enter actual qty being delivered today: 300 (partial — balance back-ordered)

08:00 — Print picking list → assign to picker

08:30 — Picker confirms goods picked
  → Enter warehouse location: WH-01-Aisle3
  → Enter shipping method: FastShip courier
  → Submit delivery note

08:31 — System auto-posts:
  → Dr 632 (COGS)  Cr 156 (Goods)  (cost = FIFO from inventory)
  → Inventory reduced by 300 units
  → SO status remains "processing" (balance 200 remaining)

08:35 — Print delivery note (2 copies: one for goods, one for customer signature)

09:00 — Courier picks up goods with signed delivery note

14:00 — Process 2 more deliveries (SO-0124 full, SO-0125 partial)
```

**Pain points:**
- Cannot see real-time inventory availability during delivery
- No barcode/RFID scanning (P2)
- No mobile app for warehouse (P2)

---

## Journey 3: AR Accountant — Issue Invoices and Manage AR

**Role:** AR Accountant (kế toán công nợ)
**Goal:** Issue e-invoices on time, track AR, collect payments

### Day-in-the-life

```
08:00 — Review "Ready to Invoice" list (delivered but not invoiced)
  → 3 deliveries ready for invoicing today

08:15 — Process delivery DN-202607-0456 (ABC Corp)
  → Select DN → "Create Invoice"
  → System shows delivered qty, already invoiced qty (0)
  → Select all lines for invoicing
  → Review invoice preview: amounts, VAT, customer info
  → Click "Issue"

08:16 — System processes e-invoice pipeline:
  → Generate TXML XML (Decree 254 format)
  → Sign with company digital certificate
  → Submit to GDT API
  → GDT returns invoice code: 7F8A2B3C
  → Update invoice status → issued
  → Post GL: Dr 131 Cr 5111 + 3331

08:17 — Invoice INV-202607-0123 issued successfully
  → System sends email to customer with PDF + QR code

09:00 — Check AR aging report
  → Total AR: 1.2B VND
  → Current: 800M, 1-30 days: 250M, 31-60: 100M, 61-90: 50M
  → Flag 3 invoices > 60 days for collection follow-up

10:00 — Customer XYZ calls: requests invoice copy
  → Search invoice by number → re-send PDF email

14:00 — Process credit note for DEF Co (returned goods)
  → Select original invoice INV-202607-0050
  → Enter return: 10 units @ 500K = 5M + VAT 0.5M = 5.5M
  → Reason: defective goods
  → Issue credit note → signed → submitted to GDT
  → Post GL: Dr 5111 Cr 131, Dr 3331 Cr 131
```

**Pain points:**
- GDT API occasionally down → manual retry needed
- Large number of invoices → batch issuance would help (P1)
- No auto-reconciliation with bank statement (P1)

---

## Journey 4: Chief Accountant — Month-End Closing

**Role:** Chief Accountant (kế toán trưởng)
**Goal:** Ensure accurate revenue recognition, AR balance, VAT output

### Month-end close

```
Day 28 — Pre-close checks
  → Run "Unbilled Deliveries Report": 2 deliveries not invoiced
  → Accrue revenue: Dr 131 (accrued)  Cr 5111
  → Run "AR sub-ledger vs GL 131" reconciliation
  → Balance matches ✓

  → Run "VAT output vs GDT" reconciliation
  → 1 invoice missing GDT code → investigate
  → Found: GDT submission failed → resubmit → code received ✓

  → Run FX revaluation for foreign currency AR
  → Rate change: 25,000 → 25,200
  → FX gain: 2M → Dr 131  Cr 515

Day 30 — Final close
  → Review AR aging report
  → Calculate bad debt provision (229):
      - 61-90 days: 50M × 50% = 25M
      - 91-120 days: 20M × 70% = 14M
      - 120+ days: 10M × 100% = 10M
      - Total provision: 49M
  → Post provision: Dr 642  Cr 229

  → Lock period for sales
  → Generate S01-BH, S02-BH, S03-BH reports
  → Sign off period closing
```

**Pain points:**
- Period close is manual and time-consuming
- No automated provision calculation (could be improved)
- S01/S02/S03 report format needs to match Circular 99 exactly

---

## Journey 5: CFO — Cash Flow and AR Oversight

**Role:** CFO (giám đốc tài chính)
**Goal:** Manage cash flow, DSO, and credit risk

### Periodic review

```
Monday morning — Review AR snapshot
  → Dashboard:
    - Total AR: 3.5B VND
    - DSO: 42 days (target < 35)
    - Overdue > 60 days: 250M (7.1%)

  → Drill down to top 5 overdue customers
  → Assign collection priorities

  → Approve 2 credit limit increases:
    - ABC Corp: 500M → 800M (good payment history)
    - XYZ Ltd: 200M → 200M (denied, inconsistent payment)

  → Export AR aging to Excel for board report
```

**Pain points:**
- No real-time dashboard (P2)
- DSO calculation not automated (P2)
- Limited drill-down analytics (P2)

---

## Journey 6: Customer — Order and Receive Invoices

**Role:** Customer (khách hàng)
**Goal:** Place orders easily, receive accurate invoices

### Customer experience

```
Step 1 — Customer sends PO via email (no portal, P3)

Step 2 — Sales Manager creates SO in GoTax
  Customer receives: order confirmation email

Step 3 — Goods delivered
  Customer receives: delivery note (signed on receipt)

Step 4 — Invoice issued
  Customer receives: e-invoice PDF via email
  → Contains: invoice number, QR code, GDT invoice code
  → Customer can verify invoice on GDT portal using invoice code

Step 5 — Payment
  Customer transfers to company bank account
  → AR accountant records receipt
  → Customer receives: payment receipt / confirmation

If return needed:
  → Customer requests return
  → Credit note issued (e-invoice negative)
  → Customer receives credit note PDF
```

**Pain points:**
- No customer self-service portal (P3)
- Manual PO submission (email) — no EDI/API (P3)

---

## Journey 7: External Auditor — Verify Sales Transactions

**Role:** External Auditor (kiểm toán độc lập)
**Goal:** Verify revenue recognition, AR existence, VAT compliance

### Audit procedures

```
1. Revenue cut-off testing
  → Pull invoices issued in last 3 days of period + first 3 days of next period
  → Match to delivery dates → verify revenue in correct period
  → Test 100% of period-end accruals

2. AR confirmation
  → Select sample of 20 customers with AR balances
  → Generate confirmation letters from customer statements
  → Verify responses match GoTax records

3. E-invoice compliance
  → Select sample of 30 e-invoices
  → Verify on GDT portal:
    - Invoice code exists
    - Amounts match
    - Signature valid
    - No cancellation after period-end

4. VAT output testing
  → Recalculate VAT output by rate (10%, 8%, 5%, 0%)
  → Verify against tax declaration
  → Test VAT reversal on credit notes

5. Bad debt provision
  → Review aging as at period-end
  → Recalculate provision per company policy
  → Verify provision adequacy
```
