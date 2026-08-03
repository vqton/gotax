# Purchase Module — User Journeys

**Version:** 2.0
**Date:** July 2026
**Author:** BA Lead + Chief Accountant (20+ yrs each)
**Note:** Journeys describe target state. Current implementation covers supplier/PO/GRN/invoice CRUD and basic AP reports.

---

## UJ-1: AP Accountant — Daily Invoice Processing

**Persona:** Trang, AP Accountant at a trading company, 3 years experience.
**Goal:** Process 30-50 supplier invoices per day, ensure accurate AP tracking.
**Frustration:** Manual data entry, chasing missing POs, reconciling supplier statements.

### Journey

| Time | Step | Touch Point | Emotion |
|------|------|-------------|---------|
| 08:00 | Check email for e-invoice notifications from GoTax | Email/Notification | Neutral |
| 08:15 | Open GoTax → Purchase → Invoices → "Pending Review" | System lists invoices auto-received from GDT | Satisfied (no data entry) |
| 08:20 | Review invoice from XYZ Corp: click to open detail | Shows parsed XML, linked PO (auto-matched) | Happy |
| 08:25 | Verify 3-way match: PO qty=10, GRN qty=10, Inv qty=10 ✓ | Green checkmark | Confident |
| 08:26 | Click "Verify & Post" → GL entry displayed for review | Shows Dr 156/1331, Cr 331 | Attentive |
| 08:27 | Click "Confirm Post" → Invoice posted, AP updated | Success message | Satisfied |
| 09:00 | New invoice without PO match (service fee) | Manual entry form | Neutral |
| 09:05 | Find supplier "ABC Maintenance" in supplier list | Search works fast | OK |
| 09:10 | Enter invoice: service line, GL account 642, no VAT | Auto-calculates, GL preview | Efficient |
| 09:12 | Click "Post" → Done | Success | Productive |
| 14:00 | Supplier calls to check payment status of INV-12300 | Phone call | Stressed |
| 14:05 | Open AP aging → search supplier → see "Due 05/08" | Quick lookup | Relieved |
| 14:06 | Inform supplier "Will be paid on 05/08 as scheduled" | — | Professional |
| 16:00 | Month-end: run "Uninvoiced Receipts Report" | 3 items found | Proactive |
| 16:10 | Create accrual entries for goods received not invoiced | Auto-generates GL | Efficient |
| 16:30 | Run "Supplier Ledger vs GL" reconciliation → balanced ✓ | Zero variance | Satisfied |

**Outcome:** 45 invoices processed, no errors, end-of-day AP aging accurate. Ready for month-end close.

---

## UJ-2: Warehouse Keeper — Goods Receipt

**Persona:** Minh, Warehouse Keeper, 5 years experience.
**Goal:** Receive goods quickly, ensure stock accuracy, minimize paperwork.
**Frustration:** Paper GRNs, manual re-typing, delayed PO references.

### Journey

| Time | Step | Touch Point | Emotion |
|------|------|-------------|---------|
| 09:00 | Delivery arrives from XYZ Corp: 10 computers + 10 monitors | Physical delivery | Neutral |
| 09:05 | Check delivery note against PO PO-202607-0001 on phone | Mobile-friendly view | Good |
| 09:10 | Open GoTax → Purchase → Receipts → "New from PO" | Scan PO number barcode | Fast |
| 09:12 | System pre-fills PO lines with ordered quantities | No typing | Happy |
| 09:15 | Count items: computers OK (10), monitors OK (10) | Enter received qty = 10 each | Focused |
| 09:16 | Quality check: all boxes sealed, no damage | Mark quality OK | Trusting |
| 09:17 | Click "Submit" → GRN created, inventory updated | Receipt posted | Satisfied |
| 09:18 | Printer auto-prints GRN for driver signature | Physical copy | Complete |
| 09:30 | Enter 2nd delivery: partial receipt (5 of 10 ordered) | Enter qty=5, balance open | Efficient |
| 09:35 | Reject 1 damaged unit → enter qty_ok=4, qty_reject=1 | Auto-flags for return | Proactive |

**Outcome:** Goods received in < 15 min per delivery, real-time inventory update, zero paperwork errors.

---

## UJ-3: Chief Accountant — Month-End AP Review

**Persona:** Ms. Hoa, Chief Accountant, 20 years experience.
**Goal:** Ensure AP accuracy before month-end close, comply with Circular 99.
**Frustration:** Manual spreadsheets, last-minute corrections, audit findings.

### Journey

| Time | Step | Touch Point | Emotion |
|------|------|-------------|---------|
| Day -5 | Run "Pre-Close Checklist" | System audit report | Organized |
| Day -5 | Review AP aging: total 886M, 36% current, 41% 1-30d, 11% 30-60d, 6% 60-90d, 6% 90+ | Clear visualization | Informed |
| Day -5 | Drill down on 90+ overdue: 50M to ABC Ltd → flag for follow-up | Action required | Concerned |
| Day -4 | Review "Uninvoiced Receipts": 3 items, total 85M | Accrual needed | Thorough |
| Day -4 | Approve accrual entries: Dr 156 85M, Cr 338 85M | One-click approval | Efficient |
| Day -3 | Approve invoices needing override (price variance >5%) | 2 items, reviewed, approved | Controlled |
| Day -2 | Run "AP Sub-ledger vs GL 331" → variance = 0 ✓ | Green | Confident |
| Day -1 | Approve bulk payment proposal: 320M to key suppliers | Select invoices, approve | Strategic |
| Day 0 | Period lock: close purchase period, no further postings | Action: lock period | Complete |
| Day +1 | Sign month-end AP package: S01-DN, S02-DN, aging report | E-sign | Final |

**Outcome:** Clean month-end close, accurate AP, zero audit findings. CFO receives AP aging with full confidence.

---

## UJ-4: Purchasing Manager — PO Management

**Persona:** Mr. Tuan, Purchasing Manager, 8 years experience.
**Goal:** Maintain optimal inventory levels, control purchase costs, manage suppliers.
**Frustration:** No visibility into PO status, delayed approvals, budget overruns.

### Journey

| Time | Step | Touch Point | Emotion |
|------|------|-------------|---------|
| 08:00 | Dashboard: "Open POs" widget shows 15 POs pending | Overview | Informed |
| 08:05 | 5 POs pending approval: filter by ">100M" → 2 items | Shows ABC Ltd 150M, XYZ 200M | Focused |
| 08:10 | Review ABC Ltd PO: supplier history good, price OK | Click approve | Fast |
| 08:12 | Review XYZ PO: price seems high vs last purchase → hold | Flag for negotiation | Cautious |
| 09:00 | Create new PO for DEF Co: select from approved requisition | Pre-filled | Efficient |
| 09:05 | Generate PO, print, email to supplier | Auto-email with PDF | Productive |
| 14:00 | Check PO status: PO-0005 → partial (7/10 received) | Visual progress | Transparent |
| 14:30 | Create PO amendment: increase qty from 10 to 12 | Change order with version | Flexible |
| 16:00 | Run "Purchase Summary" report: this month 1.2B, budget 1.5B → 80% used | On track | Confident |

**Outcome:** All urgent POs processed by noon, spending within budget, suppliers managed effectively.

---

## UJ-5: CFO — AP Cash Flow Planning

**Persona:** Mr. Duc, CFO, 15 years experience.
**Goal:** Optimize cash flow, minimize late payment penalties, negotiate early payment discounts.
**Frustration:** Can't predict AP outflows, manual cash forecasting.

### Journey

| Time | Step | Touch Point | Emotion |
|------|------|-------------|---------|
| Monday AM | Open GoTax → Reports → AP Aging Dashboard | Visual: aging distribution | Overview |
| Monday AM | See 886M total AP, 366M due within 30 days | Payment pressure | Concerned |
| Monday AM | Switch to Cash Flow View: next week -120M, next month -366M | Prediction | Informed |
| Monday AM | Identify 50M overdue 90+ days → risk of supplier stop | Action needed | Alerted |
| 10:00 | Review "Early Payment Discount Opportunities" | 2 invoices: save 2M if paid in 10 days | Strategic |
| 10:05 | Approve early payment: 100M with 2% discount → save 2M | Approve | Satisfied |
| 11:00 | Approve bulk payment batch: 320M to top 5 suppliers | One-click approval | Efficient |
| Monday PM | Download AP forecast for board meeting | Excel export | Prepared |

**Outcome:** Cash flow optimized, early discounts captured, supplier relationships maintained.

---

## UJ-6: Tax Accountant — VAT Input Declaration

**Persona:** Ms. Lan, Tax Accountant, 6 years experience.
**Goal:** Accurate VAT input declaration, maximize deduction, avoid GDT audit.
**Frustration:** Matching purchase invoices to VAT declaration, data inconsistencies.

### Journey

| Time | Step | Touch Point | Emotion |
|------|------|-------------|---------|
| Day 5 | Open GoTax → Purchase → VAT Input Report for July | List of 45 invoices | Ready |
| Day 5 | Verify: total VAT input = 75M from purchase module | Auto-calculated | Confident |
| Day 5 | Check: 3 invoices >20M, all have bank payment evidence ✓ | Auto-validated | Compliant |
| Day 5 | Flag: 1 invoice with 10% VAT rate but PO says 8% | Investigation needed | Cautious |
| Day 5 | Contact supplier → confirmed error, will issue adjustment | Create note on invoice | Proactive |
| Day 5 | Exclude 1 invoice from claim: personal use goods → not deductible | Mark as non-deductible | Accurate |
| Day 6 | Generate "VAT Input Summary" for tax return | Export to HTKK format | Efficient |
| Day 6 | Submit VAT return: VAT input 74M | Done | Satisfied |

**Outcome:** VAT deduction optimized, compliant with Decree 123, audit-ready records.

---

## UJ-7: External Auditor — AP Audit

**Persona:** Mr. Phuc, External Auditor from Big4, 10 years experience.
**Goal:** Verify AP existence, accuracy, and cut-off; assess internal controls.
**Frustration:** Manual sampling, paper trails, slow document retrieval.

### Journey

| Time | Step | Touch Point | Emotion |
|------|------|-------------|---------|
| Day 1 | Request: AP aging as at 31/07/2026 | Export from system | Fast |
| Day 1 | Sample selection: top 5 suppliers (85% of balance) | Automated sampling | Efficient |
| Day 1 | For each: request PO, GRN, invoice, payment evidence | All in one view per invoice | Impressed |
| Day 1 | Verify 3-way match for 5 sampled invoices | Auto-match history visible | Confident |
| Day 1 | Check cut-off: GRNs before 31/07 posted, invoices after 31/07 properly dated | Automated | Compliant |
| Day 2 | Re-perform AP aging calculation | System matches manual calc | Verified |
| Day 2 | Test AP sub-ledger = GL 331 balance | Zero variance | Clean |
| Day 2 | Review VAT input deduction for sampled invoices | Valid e-invoices, bank payments | Conform |
| Day 2 | Conclusion: AP controls effective, no material misstatement | Clean audit opinion | Satisfied |

**Outcome:** Audit completed in 2 days (vs 5 days manual), zero findings, recommendation for control effectiveness.