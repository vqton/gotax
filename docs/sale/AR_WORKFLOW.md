# AR Module — Order-to-Cash Collection Workflow

**Version:** 1.0
**Date:** July 2026
**Author:** Enterprise Systems Analyst + Chief Accountant (20+ yrs each)
**Regulatory:** Circular 99/2025/TT-BTC (Accts 131, 511, 3331, 521, 229, 642), Decree 123/2020, Decree 254/2026, IFRS 9/VAS 17 (bad debt)
**Benchmarks:** MISA AMIS 2026, FAST Accounting 2026, Bravo ERP 10, Odoo 19 Accounting, SAP FSCM

---

## End-to-End AR Lifecycle

```
 ┌────────────────────────────────────────────────────────────────────────────────────────┐
 │                         ACCOUNTS RECEIVABLE — FULL LIFECYCLE                           │
 └────────────────────────────────────────────────────────────────────────────────────────┘

 Customer Invoice ──→ E-Invoice ──→ GL Post ──→ Aging ──→ Collection ──→ Receipt ──→ Allocation ──→ Settlement
      (Create)      (Sign+Submit)   (Dr131/Cr511) (Buckets)   (Dunning)     (Pay)      (Apply)        (Close)
        │               │              │            │             │           │           │              │
        └── Prepayment / Deposit ──────┘             │             └── Credit Note / Write-Off ──────────┘
                                                      │               (Sales Return / Bad Debt)
                                                      └── Customer Statement (Periodic)
```

---

## Step 1: Invoice Creation

| Dimension | Detail |
|-----------|--------|
| **State & Lifecycle** | Pre-state: SO delivered (DN posted). Trigger: goods delivered/service rendered. Post-state: Invoice created (DRAFT). |
| **Actor** | AR Accountant. Auto-created by system if configured (from signed DN). |
| **Input** | Sales Order (IDs + lines), Delivery Note (qty delivered), Customer master (name, tax code, address, bank), Price/contract terms. |
| **Validations** | R13/R14: Invoice qty ≤ delivered qty per line (tolerance 0%). R17: Unit price variance vs SO price ≤ config threshold (default 2%). R12: No duplicate for same SO+DN combination. R15: Credit limit check — outstanding AR + invoice total ≤ credit limit. |
| **Output** | `CustomerInvoice` in DRAFT status. AR sub-ledger entry (ARTransInvoice, positive amount). |
| **Exception** | Over-invoice blocked. Price variance flagged for AR manager approval. Credit breach: warn (configurable block). Missing SO/DN reference allowed for service/one-off invoices. |

**Benchmark notes:**
- MISA: Invoice auto-created from DN with 2-way match
- SAP: Billing document from delivery — 3-way match (PO-DN-Invoice) for procurement, 2-way (SO-DN) for sales
- Odoo: Invoice from SO or DN — auto-validation on confirm
- Bravo: Allows standalone invoice (service) or from delivery

---

## Step 2: E-Invoice Pipeline (Decree 254/2026)

| Dimension | Detail |
|-----------|--------|
| **State & Lifecycle** | Pre-state: Invoice DRAFT. Trigger: Invoice → Sign. States: DRAFT → SIGNED → SUBMITTED → CODED → ISSUED → POSTED. |
| **Actor** | System (auto-pipeline). AR Accountant (manual re-trigger on failure). |
| **Input** | Invoice data (all fields), Digital certificate (DigitalSignature model), GDT API endpoint, TXML template (Decree 254 format). |
| **Validations** | R4: E-invoice mandatory within 24h of delivery. R21: SLA monitoring. Digital signature valid + unexpired. GDT API connectivity check. XML schema validation before submit. |
| **Output** | Signed XML (XMLDSig), GDT invoice code, QR code, Invoice status updated. |
| **Exception** | GDT rejection → log error, notify AR accountant, allow fix + resubmit. Certificate expired → block pipeline, alert admin. API timeout → retry queue (3 attempts), then escalate. |

**Benchmark notes:**
- MISA: Built-in e-invoice (MISA e-invoice) — auto-generate TXML, sign, submit
- FAST: Fast e-invoice integration — supports VNPT, Viettel, Hoadon software
- Bravo: Multi-provider e-invoice (BKAV, VNPT, Viettel) via plugin
- Odoo: Viindoo/Toannang e-invoice modules — external API integration
- SAP: SAP Document Compliance for Vietnam — certified GDT integration

---

## Step 3: GL Posting

| Dimension | Detail |
|-----------|--------|
| **State & Lifecycle** | Pre-state: Invoice ISSUED (GDT coded). Trigger: GL posting action. Post-state: Invoice POSTED. |
| **Actor** | AR Accountant. System (auto-post on e-invoice issue if configured). |
| **Input** | Invoice lines with accounts (revenue_account_id, vat_account_id, item line totals), Customer AR account (131 default). |
| **Validations** | R1: Revenue recognition timing correct. R2: Revenue account by type (5111 goods, 5112 services, etc.). R3: VAT output by rate (33311/33312). Account 131 must not be on credit hold. Period must be open. |
| **Output** | Journal Entry: Dr 131 (AR total) Cr 511X (revenue) and Cr 3331 (VAT output). `gl_posted=true`, `gl_posted_at=now()`. |
| **Exception** | GL account invalid → block, require AR manager to fix account mapping. Period closed → cannot post, require period re-open or next period. |

**GL Posting Map (Domestic, goods 100M + VAT 10%):**
```
Dr 131 (AR)            110,000,000
  Cr 5111 (Revenue)    100,000,000
  Cr 3331 (VAT output)  10,000,000
```

**Benchmark notes:**
- MISA: Auto-post to GL on invoice issue — configurable per account template
- SAP: Revenue recognition via revenue account determination (condition technique)
- Odoo: Auto-post on invoice validate — account from product/fiscal position
- All 4: GL period check, account validation, audit trail

---

## Step 4: AR Aging Tracking

| Dimension | Detail |
|-----------|--------|
| **State & Lifecycle** | Continuous. Trigger: Daily batch or on-demand. Invoice is aged from DueDate (not invoice date). |
| **Actor** | System (daily calc). AR Accountant/CFO (on-demand report). |
| **Input** | Invoices with DueDate, BalanceDue, AmountReceived. AR sub-ledger transactions. |
| **Validations** | R8: Aging buckets per Circular 99 guidance. DueDate calc from invoice date + customer payment terms. Foreign currency AR revalued at report date rate. |
| **Output** | ARAgingReport per customer with buckets: Current (0-30d), 31-60d, 61-90d, 91-120d, 120+d. Each bucket = sum of invoice BalanceDue where Days Overdue = ReportDate - DueDate falls in bucket range. |
| **Exception** | Negative aging (payment before due date) → current bucket. Credit balance → separate credit aging report. Zero-balance invoices excluded. |

**Aging Bucket Periods (Circular 99):**

| Bucket | Days Overdue | Provision % (R8) |
|--------|-------------|-------------------|
| Current | ≤ 0 (not yet due) | 0% |
| 1-30 | 1 to 30 | 0% |
| 31-60 | 31 to 60 | 0% (< 6 months) |
| 61-90 | 61 to 90 | 0% |
| 91-120 | 91 to 120 | 0% |
| 121+ | 121 to 180 | 0% |
| 181-365 | 6-12 months | 30% |
| 366-730 | 1-2 years | 50% |
| 731-1095 | 2-3 years | 70% |
| 1096+ | 3+ years | 100% |

**Note:** Current implementation has ALL amounts assigned to Bucket0 regardless of DueDate. Must fix.

**Benchmark notes:**
- MISA: Aging by invoice, automated bucket calc, drill-down to invoice detail
- FAST: Multi-currency aging, provision auto-calc per VAS 17
- Bravo: Aging with customer credit limit comparison, interactive dashboard
- SAP FSCM: Aging with dispute tracking, promise-to-pay integration
- Odoo: Aging by due date, configurable bucket ranges, PDF/Excel export

---

## Step 5: Collection / Dunning Process

| Dimension | Detail |
|-----------|--------|
| **State & Lifecycle** | Pre-state: Invoice overdue. Trigger: DueDate passed + no full payment. States: NOT_DUE → OVERDUE → DUNNING_L1 → DUNNING_L2 → DUNNING_L3 → DUNNING_L4 → COLLECTIONS → WRITE_OFF. |
| **Actor** | AR Accountant (sends reminders). Collection Manager (escalation). System (auto-send L1/L2 emails). |
| **Input** | AR aging report, Customer contact info (email, phone), Invoice detail, Payment history. |
| **Validations** | R8: Bad debt provision age bands. Segmentation: amount thresholds for escalation. Minimum overdue amount to trigger dunning (configurable, default 100K VND). No dunning if dispute/promise-to-pay active. Grace period configurable per customer. |
| **Output** | Dunning letter (email/SMS/portal notification). Collection note. Customer promise-to-pay record. Escalation to supervisor if L3+ or amount > threshold. |
| **Exception** | Customer disputes invoice → suspend dunning, route to dispute resolution. Customer bankrupt → immediate escalation to L4 + legal. Wrong contact → flag for data correction. |

**Dunning Levels:**

| Level | Trigger | Action | Channel |
|-------|---------|--------|---------|
| L1 — Reminder | 1-7 days overdue | Friendly reminder, attach invoice copy | Email |
| L2 — Formal | 8-30 days overdue | Formal demand, include late fee if applicable | Email + SMS |
| L3 — Final | 31-60 days overdue | Final notice, threatened legal action, call | Email + SMS + Phone |
| L4 — Escalation | 60+ days overdue | Legal action, bad debt write-off review | Formal letter |

**Benchmark notes:**
- MISA: Collection support via reminders, aging dashboard
- FAST: Integrated SMS/email reminders, configurable dunning levels
- Bravo: Full dunning module with dunning levels, fees, interest calculation
- SAP FSCM: Collections Management with worklist, promise-to-pay, dispute tracking, auto-dunning
- Odoo: Follow-up levels, automated email reminders, activities for collectors

---

## Step 6: Customer Receipt / Payment Collection

| Dimension | Detail |
|-----------|--------|
| **State & Lifecycle** | Pre-state: Invoice overdue/not due. Trigger: Customer pays. Post-state: Receipt created (DRAFT) → POSTED. |
| **Actor** | AR Accountant (manual entry). Treasurer (bank transfer confirmation). System (auto-import from bank statement). |
| **Input** | Payment notification (bank statement line, cheque, cash count), Customer reference, Amount, Payment method, Bank account. |
| **Validations** | FR-7.1-7.5: Payment method valid. Bank account exists (if bank transfer). Currency match or FX rate provided. Customer exists + not suspended. Receipt number unique. Amount > 0. |
| **Output** | `CustomerReceipt` (DRAFT). Receipt number (RCP-YYYYMM-XXXX). AR sub-ledger entry pending (posted only on POST). |
| **Exception** | Unknown customer → unallocated receipts bucket, AR accountant to identify. Amount mismatch vs remittance → hold, request clarification. Cheque returned → reverse receipt, notify AR manager. |

**Receipt State Machine:**
```
DRAFT ──→ POSTED ──→ RECONCILED ──→ CANCELLED
  │                    (bank rec)
  └────────→ CANCELLED (if never posted)
```

**Cash/Bank GL Entry (on Post):**
```
Dr 111 (Cash) / 112 (Bank)     Amount received
  Cr 131 (AR)                  Amount received
```

**Benchmark notes:**
- MISA: Receipt from collection module, auto-allocate to oldest invoice (configurable)
- FAST: Multi-payment support (cash, bank, cheque, credit card), auto-bank feed
- Bravo: Cash management integrated with AR — receipt, allocation, bank rec in one screen
- SAP: FJAM (Auto-cash application) — AI-based receipt allocation
- Odoo: Auto-reconciliation from bank statements, manual allocation UI

---

## Step 7: Receipt Allocation

| Dimension | Detail |
|-----------|--------|
| **State & Lifecycle** | Pre-state: Receipt POSTED, Invoices outstanding. Trigger: Allocate receipt to invoice(s). Post-state: Allocations recorded, invoice BalanceDue reduced. |
| **Actor** | AR Accountant. System (auto-allocate if single reference match). |
| **Input** | Receipt ID, Invoices list (IDs + amounts), Allocation amounts per invoice, Discount amounts (early payment discount). |
| **Validations** | Sum allocated amount + discount = receipt amount. Allocated amount ≤ invoice BalanceDue. Receipt currency = invoice currency (or FX conversion). Discount account 5213 if early payment discount granted. |
| **Output** | `ReceiptAllocation` records (receipt_id → invoice_id → allocated_amount). Invoice balance_due updated. Invoice status: PAID if BalanceDue ≤ 0. AR sub-ledger entries updated. |
| **Exception** | Over-allocation → block (exceeds receipt amount). Under-allocation → partial payment ok, remainder stays outstanding. Currency mismatch → use receipt FX rate, post FX gain/loss. Short payment with discount → customer entitled? Check discount terms. |

**Allocation GL:**
```
Full payment:     Dr 112   100,000,000    Cr 131   100,000,000
Partial payment:  Dr 112    60,000,000    Cr 131    60,000,000   (40M remains outstanding)
With discount:    Dr 112    98,000,000
                  Dr 5213    2,000,000    Cr 131   100,000,000
With FX gain:     Dr 112   102,000,000    Cr 131   100,000,000    (rate improved)
                                           Cr 515     2,000,000
With FX loss:     Dr 112    98,000,000
                  Dr 635     2,000,000    Cr 131   100,000,000    (rate worsened)
```

**3-Way Matching Invariant:**
```
Invoice.BalanceDue = Invoice.TotalAmount - SUM(ReceiptAllocation.AllocatedAmount) - SUM(CN.TotalAmount for original_invoice=this)
```

**Benchmark notes:**
- MISA: Auto-allocation by FIFO (oldest invoice first), manual override
- FAST: One receipt → multiple invoices, auto-suggest allocation
- Bravo: Allocation with discount auto-calc, FX revaluation at allocation
- SAP: Multiple allocation methods — by invoice number, by aging, by amount
- Odoo: Auto-reconcile, manual reconcile, write-off if difference below threshold

---

## Step 8: Credit Note / Sales Return

| Dimension | Detail |
|-----------|--------|
| **State & Lifecycle** | Pre-state: Invoice POSTED/PAID. Trigger: customer returns goods, price adjustment, discount granted. States: DRAFT → SIGNED → SUBMITTED → POSTED → CANCELLED. |
| **Actor** | AR Accountant. Warehouse Keeper (if goods returned). |
| **Input** | Original Invoice ID, Return reason code, Returned items (qty, condition), Original invoice lines for reference, Goods return DN (if inventory). |
| **Validations** | R5: Must reference original invoice + GDT code. Return qty ≤ original invoice qty (partial return ok). Cannot credit note if invoice is CANCELLED. Original invoice must be POSTED/PAID. Segregation: creator ≠ approver for CN > threshold. |
| **Output** | Credit Note document. Negative e-invoice TXML. AR sub-ledger transaction (ARTransCreditNote, negative). Invoice BalanceDue reduced. |
| **Exception** | Goods damaged → partial credit only (customer responsible). Goods returned after warranty → no credit, customer pays. Price adjustment without return → no inventory reversal. Post-invoice discount → CN type = "discount/allowance", no inventory movement. |

**GL Entries (on Post):**

| Scenario | Debit | Credit |
|----------|-------|--------|
| Revenue reversal | 5111 (Revenue) | 131 (AR) |
| VAT reversal | 3331 (VAT output) | 131 (AR) |
| Inventory return | 156 (Goods) | 632 (COGS) |
| Price adjustment only | 5111 (Revenue) | 131 (AR) |
| Discount/allowance | 5211/5212 | 131 (AR) |

**Benchmark notes:**
- MISA: CN from return DN — auto-reverse revenue + VAT + inventory
- FAST: CN with reason code, auto-e-invoice generation for negative invoice
- Bravo: Multi-type credit notes (return, discount, allowance, bad debt)
- SAP: Credit memo request → credit memo — approval workflow
- Odoo: Credit note wizard — full refund, partial refund, price adjustment

---

## Step 9: Prepayment / Customer Deposit

| Dimension | Detail |
|-----------|--------|
| **State & Lifecycle** | Pre-state: SO approved, no goods delivered. Trigger: Customer pays deposit before delivery. Post-state: Prepayment recorded, AR credit balance. |
| **Actor** | AR Accountant. Sales Manager (request deposit). |
| **Input** | Sales Order or contract, Payment receipt, Deposit amount, Optional deposit e-invoice request. |
| **Validations** | R10: Deposit amount ≤ order total. Customer exists + not blacklisted. Deposit e-invoice required per Decree 254 if customer requests. |
| **Output** | CustomerReceipt (prepayment type). AR sub-ledger ARTransPrepayment (negative = credit balance). Optional deposit e-invoice to GDT. |
| **Exception** | Customer overpays deposit → track excess as credit balance. Customer wants deposit refund → create refund receipt, offset against deposit. Deposit e-invoice issued but order cancelled → must cancel e-invoice on GDT. |

**GL Entries:**

| Step | Debit | Credit |
|------|-------|--------|
| Receive deposit | 112 (Bank) 100,000,000 | 131 (AR) 100,000,000 |
| Offset against final invoice | 131 (AR credit bal) 100,000,000 | 131 (AR invoice) 100,000,000 |
| Final invoice (remaining) | 131 (AR) 10,000,000 | 5111 (Revenue) + 3331 (VAT) |

Deposit e-invoice (if issued at deposit time):
```
Dr — no GL entry (deposit is advance, not revenue)
  Simply advances showing on AR aging as negative/credit
```

**Benchmark notes:**
- MISA: Deposit receipt with advance invoice, offset on final invoice
- FAST: Prepayment tracking, auto-offset on invoice
- Bravo: Advance payment module, offset via clearing account
- SAP: Customer down payment — special GL indicator, auto-clearing
- Odoo: Down payment invoice workflow — percentage or fixed amount

---

## Step 10: Bad Debt Write-Off

| Dimension | Detail |
|-----------|--------|
| **State & Lifecycle** | Pre-state: Invoice overdue 180+ days, collection exhausted. Trigger: Director approval for write-off. Post-state: Invoice written off, tracked off-balance-sheet. |
| **Actor** | Chief Accountant (propose). Director (approve). AR Accountant (execute). |
| **Input** | Invoice detail, Aging report (120+ days), Collection history, Dunning trail, Legal action outcome (if any), Director approval document. |
| **Validations** | R8: Write-off per Circular 99 — all collection attempts documented. Approval required (segregation: proposer ≠ approver). Amount ≤ approved write-off limit. Account 131 balance verified before write-off. Must not write off if customer has credit balance (offset first). |
| **Output** | Write-off journal entry. Invoice status: WRITTEN_OFF. Off-balance-sheet tracking record. AR sub-ledger entry (ARTransOffset?). Tax implications: VAT adjustment may be required. |
| **Exception** | Customer pays after write-off → reinstate AR, record payment. Partial recovery → record collected amount, remaining stays written off. Fraud suspicion → escalate to legal + audit. |

**GL Entries:**

| Method | Debit | Credit | Condition |
|--------|-------|--------|-----------|
| Direct write-off | 642 (Admin expense) | 131 (AR) | Small amounts, no provision |
| Use provision | 229 (Provision) | 131 (AR) | If provision previously booked |
| Recovery | 131 (AR) | 711 (Other income) | Reinstatement |
| Recovery (cash) | 112 (Bank) | 131 (AR) | Cash received |

**Off-balance-sheet tracking (10 years per Accounting Law Art 13):**
```
Off-balance-sheet: 00X — Written-off AR (monitoring account)
```

**Benchmark notes:**
- MISA: Bad debt provision auto-calc, write-off with approval, off-balance-sheet tracking
- FAST: Provision by aging %, write-off with document reference
- Bravo: Multi-method provision (% by age, individual assessment), write-off workflow
- SAP: FSCM Bad Debt — provision calculation, write-off, recovery, credit insurance
- Odoo: Write-off in reconciliation, off-balance-sheet tracking via accounting entries

---

## Step 11: Customer Statement

| Dimension | Detail |
|-----------|--------|
| **State & Lifecycle** | Periodic (monthly/quarterly). Trigger: Statement generation for all active customers with balance ≠ 0. |
| **Actor** | AR Accountant (generate). System (auto-generate and email). Customer (receive). |
| **Input** | Customer invoices (period), Receipts (period), Credit notes (period), Opening balance (prior period), Payment terms. |
| **Validations** | Opening balance = prior period closing. All transactions in period included. Statement date = period end. No unposted invoices/receipts (must post first). |
| **Output** | Customer Statement of Account: Opening balance, invoices, payments, credit notes, adjustments, closing balance. PDF (email/portal). |
| **Exception** | Disputed invoice → mark on statement as "disputed". Zero-balance customer → skip (configurable). Foreign currency → show in original currency + VND equivalent. |

**Statement Format:**
```
CÔNG TY TNHH GOTAX                                          GIAY BAO CONG NO
MST: 0123456789                                              So: STM-202607-0001
                                                             Ngay: 31/07/2026

Khach hang: Cong ty TNHH ABC
MST: 0987654321
Dia chi: 123 Nguyen Hue, Q1, HCMC

KY TRUOC:                                               100,000,000
  Ngay    Chung tu       Noi dung           Phat sinh     Phat sinh   So du
                                            No            Co
  01/07   INV-202607-01  Hang hoa A        110,000,000                210,000,000
  05/07   RCP-202607-01  Thanh toan INV-01                50,000,000  160,000,000
  10/07   INV-202607-02  Hang hoa B        220,000,000                380,000,000
  15/07   CN-202607-01   Tra lai INV-01                   110,000,000 270,000,000
  20/07   RCP-202607-02  Thanh toan INV-02               220,000,000   50,000,000
KY NAY:                                                  50,000,000

QUA HAN:                                                  0
DEN HAN:                                                  50,000,000  (INV-202607-03, due 15/08)
```

**Benchmark notes:**
- MISA: Customer statement with aging, multi-currency, email auto-send
- FAST: Statement with payment reminders, overdue highlight
- Bravo: AR statement with credit limit, aging, due dates — PDF/Excel/HTML
- SAP: F.27 (Customer Statement) — configurable format, multiple variants
- Odoo: Statement from accounting — filter by date, include detail lines

---

## Step 12: Month-End Closing

| Dimension | Detail |
|-----------|--------|
| **State & Lifecycle** | Pre-state: Month end. Trigger: Period close procedure. Post-state: Period locked, AR reconciled, reports generated. |
| **Actor** | Chief Accountant (execute). AR Accountant (prepare). |
| **Input** | All AR transactions for period, GL balance for 131, Bank statements, Customer confirmations. |
| **Validations** | R19: AR sub-ledger total = GL 131 balance (zero variance). R20: VAT output from invoices = GDT registered total. R18: All deliveries invoiced (unbilled deliveries report). No DRAFT invoices/receipts unposted. FX revaluation completed for foreign currency AR. Bad debt provision calculated and posted. |
| **Output** | AR closing checklist: (1) reconciliation report, (2) VAT output reconciliation, (3) unbilled deliveries report, (4) aging report, (5) provision calculation, (6) FX revaluation entries, (7) customer confirmations sent. |
| **Exception** | AR-GL variance → hold close, investigate, post adjustment. Missing e-invoice codes → investigate before VAT declaration. Unbilled deliveries accrue revenue (Dr 131 accrued Cr 511 — reverse next period). Unreconciled bank receipts → follow up in next period. |

**Period-End Checklist:**

```
□ 1. Verify all invoices posted to GL (no unposted invoices)
□ 2. AR sub-ledger = GL 131 balance
□ 3. VAT output = GDT invoice codes (reconciliation)
□ 4. All deliveries invoiced (accrue unbilled)
□ 5. FX revaluation (foreign currency AR)
□ 6. Bad debt provision (Account 229)
□ 7. Aging report generated + reviewed
□ 8. Customer statements sent
□ 9. AR aging reviewed by CFO
□10. Period closed (no new AR transactions in closed period)
```

**Benchmark notes:**
- MISA: Period-lock by module, AR-GL reconciliation report, month-end wizard
- FAST: Auto-period check, AR-GL reconcile button, blocking period indicators
- Bravo: Multi-step period close with checklist, variance warnings
- SAP: F.07 (Balance Carryforward), OB52 (Period control), AR-GL rec (F.13)
- Odoo: Period lock by date, lock date on journal, reconciliation dashboard

---

## Benchmark Summary: Feature Coverage

| Feature | GoTax | MISA | FAST | Bravo | Odoo | SAP |
|---------|-------|------|------|-------|------|-----|
| Customer master | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Credit limit check | 🟡 field exists, no logic | ✅ | ✅ | ✅ | ✅ | ✅ |
| Sales Order → Invoice | ✅ 2-way match | ✅ | ✅ | ✅ | ✅ | ✅ |
| E-invoice (TXML→sign→GDT) | 🟡 fields exist, no pipeline | ✅ | ✅ | ✅ | 🟡 module | ✅ |
| GL auto-posting | 🟡 GLPosted flag, no JE | ✅ | ✅ | ✅ | ✅ | ✅ |
| AR sub-ledger | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| AR aging (by invoice date) | 🟡 all in Bucket0 | ✅ | ✅ | ✅ | ✅ | ✅ |
| AR aging (by DueDate) | ❌ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Customer statement | ❌ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Dunning / Collection | ❌ | 🟡 basic | ✅ | ✅ | ✅ | ✅ |
| Bad debt provision | ❌ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Bad debt write-off | ❌ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Prepayment/deposit | 🟡 type exists, no workflow | ✅ | ✅ | ✅ | ✅ | ✅ |
| Credit note (return) | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Credit note (discount) | 🟡 model supports | ✅ | ✅ | ✅ | ✅ | ✅ |
| FX revaluation | ❌ | ✅ | ✅ | ✅ | ✅ | ✅ |
| DSO calculation | ❌ | 🟡 | 🟡 | ✅ | 🟡 | ✅ |
| Auto-bank reconciliation | ❌ via Bank module | ✅ | ✅ | ✅ | ✅ | ✅ |
| Month-end AR-GL rec | ❌ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Off-balance-sheet tracking | ❌ | ✅ | ✅ | ✅ | 🟡 | ✅ |

**Legend:** ✅ = exists/complete, 🟡 = partial/field-exists-no-logic, ❌ = missing

---

## Gap Analysis vs Enterprise ERPs

### Top 5 Missing Features (P0 Priority)

| # | Feature | ERP Benchmark | Effort | Business Impact |
|---|---------|---------------|--------|----------------|
| 1 | **AR aging by DueDate** | All ERPs do this | 1-2d | Aging report useless without it — all amounts show current |
| 2 | **GL auto-posting** | All ERPs auto-post on invoice/receipt | 3-5d | Revenue not in GL = manual entry, reconciliation nightmare |
| 3 | **E-invoice pipeline (TXML→sign→GDT)** | MISA/FAST/Bravo native; Odoo/SAP via addon | 4-6w | Regulatory blocker — must issue e-invoice within 24h |
| 4 | **Customer statement** | Every ERP generates this | 3-5d | Customers need statements for their records |
| 5 | **Credit limit enforcement** | SAP/MISA enforce at SO/invoice | 2-3d | Overdue risk without credit control |

### Top 5 Missing Features (P1 Priority)

| # | Feature | ERP Benchmark | Effort |
|---|---------|---------------|--------|
| 6 | Dunning/collection workflow | Bravo/SAP full suite | 2-3w |
| 7 | Bad debt provision + write-off | All do this | 1-2w |
| 8 | FX revaluation at period-end | All do this | 3-5d |
| 9 | Prepayment/deposit workflow | All do this | 1-2w |
| 10 | AR-GL month-end reconciliation | All have this report | 2-3d |

---

## Appendix A: AR Invoice State Machine (Detailed)

```
                              ┌──────────────────────────────────────────────┐
                              │              AR INVOICE STATES              │
                              └──────────────────────────────────────────────┘

          ┌──→ CANCELLED (pre-GDT)
          │
DRAFT ────┼──→ SIGNED ──→ SUBMITTED ──→ CODED ──→ ISSUED ──→ POSTED ──→ PAID ──→ CLOSED
          │       │            │            │          │          │
          └──→ REPLACED ───→ (new invoice referencing this)
                  │
                  └──→ CANCELLED (post-GDT cancellation via GDT API)

State Transition Rules:
  DRAFT      → any: editable, no GL impact
  SIGNED     → digital signature applied
  SUBMITTED  → sent to GDT, awaiting response
  CODED      → GDT assigned invoice code (legally valid from this point)
  ISSUED     → legally issued, customer notified
  POSTED     → GL entry created (Dr 131 Cr 511/3331)
  PAID       → BalanceDue = 0 (via receipt allocation or credit note)
  CLOSED     → fully settled, no further action possible
  CANCELLED  → voided (GL reversal if posted)
  REPLACED   → superseded by corrective invoice

  State guards:
    - DRAFT → SIGNED:     TXML generated + signed
    - SIGNED → SUBMITTED: GDT API call
    - SUBMITTED → CODED:  GDT returns code (else → SUBMITTED_FAILED)
    - CODED → ISSUED:     invoice visible to customer
    - ISSUED → POSTED:    GL posting period must be open
    - POSTED → PAID:      all allocations sum = TotalAmount
    - POSTED → CANCELLED: requires GL reversal + GDT cancellation
    - Any → CANCELLED:    only if not POSTED/PAID (pre-GDT cancel is simpler)
```

## Appendix B: AR Account 131 — GL Integration

```
Account 131 (Phai thu khach hang):

  DEBIT SIDE (AR increases):               CREDIT SIDE (AR decreases):
    Invoice issued                           Customer receipt
    Adjustment increase (corrective)         Credit note (sales return)
    FX revaluation gain                      Prepayment/deposit
    Accrued revenue (unbilled delivery)       Discount/allowance granted
    Write-off recovery                        Bad debt write-off
                                              FX revaluation loss
                                              Offset to other AR

  Normal balance: DEBIT
  Credit balance = customer prepayment/deposit (advance)

  Sub-accounts per Circular 99:
    1311 — Phai thu khach hang trong nuoc (domestic)
    1312 — Phai thu khach hang nuoc ngoai (foreign)
    1313 — Phai thu noi bo (inter-company)
    1318 — Phai thu khac (other receivables)

  AR sub-ledger = Σ customer balances
  GL 131 balance = Σ AR sub-ledger
  Monthly: AR sub-ledger MUST = GL 131 (R19)
```

## Appendix C: AR Documents — Numbering & Retention

| Document | Pattern | Reset | Retention | Legal Basis |
|----------|---------|-------|-----------|-------------|
| Customer Invoice | INV-YYYYMM-XXXX | Monthly | 10 years | Accounting Law Art 13 |
| Customer Receipt | RCP-YYYYMM-XXXX | Monthly | 10 years | Accounting Law Art 13 |
| Credit Note | CN-YYYYMM-XXXX | Monthly | 10 years | Decree 123 Art 56 |
| Customer Statement | STM-YYYYMM-XXXX | Monthly | 5 years | Standard practice |
| Dunning Letter | DUN-YYYYMM-XXXX | Monthly | 5 years | Standard practice |
| Bad Debt Write-Off | WOF-YYYYMM-XXXX | Yearly | 10 years | Accounting Law Art 13 |
