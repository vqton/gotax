# Bank Module — User Journeys

---

## Journey 1: Monthly Bank Closing — Accountant

**Persona:** Trang, Accountant (3 yrs experience)
**Goal:** Complete monthly bank reconciliation for Vietcombank account
**Pain:** Manual Excel reconciliation takes 2 days

### Steps

| Step | Action | System Response | Emotion |
|------|--------|----------------|---------|
| 1 | Login to GoTax | Dashboard shows pending tasks | 😐 Neutral |
| 2 | Navigate to Bank → Statements | Statement list (empty), import button visible | 😐 |
| 3 | Download CSV from VCB e-banking | CSV file saved locally | 😐 |
| 4 | Click "Import Statement" → select CSV | File upload dialog | 😐 |
| 5 | Select bank account (VCB VND), upload file | Progress bar, parsing starts | 😐 |
| 6 | View preview: 245 lines, opening balance matches | Preview shows summary, balance OK | 🙂 Good |
| 7 | Click "Confirm Import" | Statement imported, lines stored | 🙂 |
| 8 | Navigate to Bank → Reconciliation → "Auto-Match" | 215/245 lines auto-matched (87.7%) | 🙂 |
| 9 | Review unmatched lines (30 items) | Side-by-side view: statement (left) vs GL (right) | 😐 |
| 10 | Manually match 25 lines by selecting pairs | Matches confirmed, moving to matched list | 🙂 |
| 11 | Investigate 5 remaining items | 3 bank fees (write-off), 2 unknown transactions | 😟 Confused |
| 12 | Create write-off for 3 bank fees (22,000 VND total) | Write-off posted, Nợ 635/Có 112 | 🙂 |
| 13 | Mark 2 unknown items as "Needs Investigation" | Flagged for Chief Accountant | 😐 |
| 14 | Check reconciliation summary: difference = 0 | All items resolved | 😊 Satisfied |
| 15 | Complete reconciliation | Period locked, S09-DN generated | 😊 |
| 16 | Export S08-DN + S09-DN as PDF | PDF downloaded | 😊 |

**Time saved:** 2 days → 30 minutes
**Pain points:** CSV format differs per bank (VCB vs BIDV vs CTG)

---

## Journey 2: Supplier Payment — Accountant + Chief Accountant

**Persona:** Minh, AP Accountant (maker) & Mr. Hai, Chief Accountant (checker)
**Goal:** Pay supplier invoice VND 150,000,000 due tomorrow

### Steps (Maker)

| Step | Action | System Response | Emotion |
|------|--------|----------------|---------|
| 1 | Open AP → Due invoices list | Filters: due date ≤ tomorrow | 😐 |
| 2 | Select invoice #INV-2026-0789 for VND 150M | Invoice detail shown | 😐 |
| 3 | Click "Pay via Bank" | Payment order form pre-filled | 🙂 Smart |
| 4 | Verify beneficiary, amount, bank account | Validations pass | 🙂 |
| 5 | Add payment content: "TT HD 0789" | Content validated | 😐 |
| 6 | Click "Submit for Approval" | Status → PENDING_APPROVAL | 🙂 |

### Steps (Checker)

| Step | Action | System Response | Emotion |
|------|--------|----------------|---------|
| 1 | Receive notification: "Payment pending approval" | Dashboard badge + email | 😐 |
| 2 | Open payment order review | Detail: beneficiary, amount, supporting docs | 😐 |
| 3 | Verify against PO #PO-2026-0456 | Match found, delivery confirmed | 🙂 |
| 4 | Click "Approve" | Status → APPROVED, maker notified | 🙂 |

### Steps (Maker — continued)

| Step | Action | System Response | Emotion |
|------|--------|----------------|---------|
| 5 | Print UNC form | UNC-VCB template rendered | 🙂 |
| 6 | Submit payment via VCB e-banking (external) | Wait for confirmation | 😐 |
| 7 | Return to system → Mark as Confirmed | Status → CONFIRMED, GL posted | 😊 |
| 8 | Attach bank confirmation screenshot | Evidence stored | 🙂 |

**Total time:** 15 minutes (versus 1 hour manual)
**Key insight:** UNC printing with correct bank format saves manual form filling

---

## Journey 3: Loan Management — Treasury Manager

**Persona:** Mrs. Lan, Treasury Manager (15 yrs)
**Goal:** Track new short-term loan from BIDV, set up repayment schedule

### Steps

| Step | Action | System Response | Emotion |
|------|--------|----------------|---------|
| 1 | Navigate to Bank → Loans → "New Loan" | Form: contract info | 😐 |
| 2 | Enter contract #LD-BIDV-0726: VND 5B, 8%/yr, 12 months, annuity | Loan created ACTIVE | 🙂 |
| 3 | Click "Record Disbursement": date today, full amount | Disbursement recorded | 🙂 |
| 4 | View auto-generated repayment schedule | 12 monthly payments of ~435M displayed | 😊 Clear |
| 5 | Compare with bank's amortization table | Matches exactly (difference = 0) | 😊 Confident |
| 6 | First payment due next month | System will auto-remind | 🙂 |

**One month later:**

| Step | Action | System Response | Emotion |
|------|--------|----------------|---------|
| 7 | Dashboard shows "Loan payment due in 5 days" | Alert bell | 😐 |
| 8 | Click alert → opens loan detail | Outstanding balance, due amount shown | 😐 |
| 9 | Click "Pay Installment" | Creates payment order pre-filled | 🙂 Smart |
| 10 | Provide payment content | Content: "TT đợt 1 - HĐ LD-BIDV-0726" | 😐 |
| 11 | Submit → Approve → Confirm (same flow as Journey 2) | Payment posted, loan balance updated | 😊 |

**Ongoing:** Monthly payments auto-scheduled, principal decreases each month

---

## Journey 4: Term Deposit — Treasury Manager

**Persona:** Mrs. Lan, Treasury Manager
**Goal:** Park excess cash VND 1B in 3-month term deposit, Vietcombank

### Steps

| Step | Action | System Response | Emotion |
|------|--------|----------------|---------|
| 1 | Navigate to Bank → Term Deposits → "New" | Deposit creation form | 😐 |
| 2 | Select bank account (VCB VND), amount 1B | Auto-suggest from cash balance | 😐 |
| 3 | Term: 92 days, rate: 5.5%/yr, auto-renewal: ON | Interest preview: 1B × 5.5% × 92/365 ≈ 13.86M | 🙂 |
| 4 | Confirm | GL posted: Nợ 1281 / Có 112 (1B) | 🙂 |
| 5 | Dashboard shows active deposit | CD period indicator | 😊 Visible |

**92 days later:**

| Step | Action | System Response | Emotion |
|------|--------|----------------|---------|
| 6 | System notification: "Term deposit maturing tomorrow" | Auto-triggered | 🙂 |
| 7 | Open deposit → Review maturity | Interest = 13.86M (correct) | 🙂 |
| 8 | Click "Process Maturity" | Principal + Interest → 112 | 🙂 |
| 9 | Auto-renewal triggered: new 92-day deposit created | New ACTIVE deposit at current rate (5.2%) | 😊 |
| 10 | Interest income recorded: Nợ 112 / Có 515 (13.86M) | FYI income recognized | 🙂 |

---

## Journey 5: Year-End FX Revaluation — Chief Accountant

**Persona:** Mr. Hai, Chief Accountant
**Goal:** Revalue USD bank account at year-end per Circular 99

### Steps

| Step | Action | System Response | Emotion |
|------|--------|----------------|---------|
| 1 | Navigate to Period-End → FX Revaluation | FCY accounts listed | 😐 |
| 2 | Select USD bank account (balance: $50,000) | Current book rate: 25,450 VND/USD | 😐 |
| 3 | System fetches VCB USD TTM rate: 25,620 | Rate refresh | 🙂 |
| 4 | Preview: Gain = $50,000 × (25,620-25,450) = 8,500,000 VND | Entry: Nợ 1122 / Có 515 | 🙂 |
| 5 | Confirm revaluation | GL posted, log entry created | 🙂 |
| 6 | Print FX revaluation report for auditor | Evidence for audit | 😊 Audit-ready |

---

## Journey 6: Payment Batch — Payroll Accountant

**Persona:** Ms. Thu, Payroll Accountant
**Goal:** Process monthly salary payment for 85 employees via bank

### Steps

| Step | Action | System Response | Emotion |
|------|--------|----------------|---------|
| 1 | Navigate to Bank → Payment Batches → "Create from Salary" | Batch form with salary data pre-loaded | 🙂 Smart |
| 2 | Select salary period (July 2026) | 85 employees, total VND 1.2B loaded | 🙂 |
| 3 | Default from-bank account (VCB) applied | Each line: employee name, amount, bank account | 😐 |
| 4 | Click "Validate" | Check: all account numbers valid, sufficient balance | 🙂 |
| 5 | Click "Submit Batch for Approval" | Batch → PENDING_APPROVAL | 🙂 |
| 6 | Checker reviews batch summary (not each line) | Batch-level: total 1.2B, 85 orders | 🙂 |
| 7 | Approve → Maker submits to bank | Bank file exported (format configured per bank) | 😊 |
| 8 | Confirm batch: 85/85 confirmed successful | Batch → CONFIRMED, GL postings created | 😊 |

**Time saved:** 4 hours → 20 minutes
**Key insight:** Batch salary prevents 85 individual approval workflows
