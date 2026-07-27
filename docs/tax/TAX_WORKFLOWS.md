# Tax Module — Workflows & User Journeys

**Role:** BA Lead (20+ yrs) + Chief Accountant (20+ yrs)
**Date:** 2026-07-27

---

## W-01: Complete Tax Month-End Close

**Owners:** Tax Accountant → Chief Accountant

```
┌─────────────────────────────────────────────────────────────────────┐
│ MONTH-END TAX CLOSE WORKFLOW                                        │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  [Period End]                                                       │
│      │                                                              │
│      ▼                                                              │
│  1. GL Period Close Check                                           │
│      │ ── All entries posted? ──→ Warning if open entries           │
│      ▼                                                              │
│  2. VAT Calculation                                                 │
│      ├── Auto-calc output VAT from sales (33311)                    │
│      ├── Auto-calc input VAT from purchases (1331, 1332)            │
│      ├── Generate purchase/sales ledgers                            │
│      └── Reconcile VAT with e-invoice data                          │
│      │                                                              │
│      ▼                                                              │
│  3. PIT Calculation (if applicable)                                 │
│      ├── Calculate per-employee PIT                                 │
│      ├── Generate 05/KK-TNCN                                        │
│      └── Submit to GDT                                              │
│      │                                                              │
│      ▼                                                              │
│  4. TTDB / BVMT Calculation (if applicable)                         │
│      ├── Calculate special consumption tax                          │
│      ├── Calculate environmental tax                                │
│      └── Generate declarations                                      │
│      │                                                              │
│      ▼                                                              │
│  5. VAT Declaration                                                 │
│      ├── Auto-populate 01/GTGT                                      │
│      ├── Review adjustments                                         │
│      ├── Sign with digital certificate                              │
│      └── Submit to GDT                                              │
│      │                                                              │
│      ▼                                                              │
│  6. Payment Processing                                              │
│      ├── Generate payment orders                                    │
│      ├── Pay via bank/EBPP                                          │
│      └── Reconcile payment with declaration                         │
│      │                                                              │
│      ▼                                                              │
│  7. Status Update                                                   │
│      └── Mark tax period as COMPLETED                               │
│                                                                     │
│  Deadlines:                                                         │
│    Monthly: 20th of following month                                 │
│    Quarterly: 30th of following month (or Apr 30 for Q4)            │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

**RACI:**
| Task | Tax Accountant | Chief Accountant | System |
|------|---------------|-----------------|--------|
| GL period check | C | I | R/A |
| VAT calculation | I | A | R |
| PIT calculation | I | A | R |
| TTDB/BVMT calc | I | A | R |
| Declaration review | R | A | C |
| Digital signing | C | R/A | C |
| GDT submission | C | I | R/A |
| Payment | R | A | C |
| Reconciliation | R | A | R |

---

## W-02: CIT Year-End Finalization

**Owners:** Chief Accountant → CFO → External Auditor

```
┌─────────────────────────────────────────────────────────────────────┐
│ CIT YEAR-END FINALIZATION WORKFLOW                                  │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  [Jan - Mar: Post FY Close]                                         │
│      │                                                              │
│      ▼                                                              │
│  1. Financial Statements Preparation                                │
│      ├── Balance Sheet (B01-DN)                                     │
│      ├── Income Statement (B02-DN)                                  │
│      ├── Cash Flow (B03-DN)                                         │
│      ├── Notes (B09-DN)                                             │
│      └── Audit (if mandatory or required)                           │
│      │                                                              │
│      ▼                                                              │
│  2. CIT Provisional Reconciliation                                  │
│      ├── Sum quarterly provisional payments (04/TNDN)               │
│      ├── Check ≥80% of final CIT requirement                        │
│      └── Calculate late interest if <80%                            │
│      │                                                              │
│      ▼                                                              │
│  3. CIT Finalization Calculation                                    │
│      ├── Revenue determination                                      │
│      ├── Deductible expense review                                  │
│      ├── Non-deductible adjustments                                 │
│      ├── Tax incentive calculation                                  │
│      ├── Tax loss carry-forward                                     │
│      ├── Related-party transaction flag                             │
│      └── Global minimum tax (Pillar 2) check                        │
│      │                                                              │
│      ▼                                                              │
│  4. CIT Declaration                                                 │
│      ├── Generate 03/TNDN form                                      │
│      ├── Generate appendices (incentives, KHCN etc.)                │
│      ├── Review by Chief Accountant                                 │
│      ├── Sign with digital certificate                              │
│      └── Submit to GDT                                              │
│      │                                                              │
│      ▼                                                              │
│  5. CIT Payment                                                     │
│      ├── Calculate final amount due (or refund)                     │
│      ├── If refund: Submit refund application                       │
│      ├── If due: Pay remaining balance                              │
│      └── Reconcile                                                  │
│      │                                                              │
│      ▼                                                              │
│  6. Filing Confirmation                                             │
│      ├── Download GDT acknowledgement                               │
│      ├── Store with tax records                                     │
│      └── Notify CFO of completion                                   │
│                                                                     │
│  Deadline: 31-Mar (calendar year)                                   │
│  Penalty for late: 0.03%/day on unpaid amount (Decree 310)          │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## W-03: E-Invoice Issuance Pipeline

**Owner:** AR Accountant

```
┌────────────────────────────────────────────────────────────────────────┐
│ E-INVOICE ISSUANCE PIPELINE                                             │
├────────────────────────────────────────────────────────────────────────┤
│                                                                        │
│  [Sale Transaction Occurs]                                             │
│      │                                                                 │
│      ▼                                                                 │
│  1. Data Entry                                                         │
│      ├── From sales order (auto) or manual entry                       │
│      ├── Buyer info: name, tax code, address, email                    │
│      ├── Line items: description, qty, unit price, VAT%                │
│      └── Currency: VND or foreign (with exchange rate)                 │
│      │                                                                 │
│      ▼                                                                 │
│  2. Invoice Creation                                                   │
│      ├── Select invoice pattern (from company's registered patterns)   │
│      ├── Generate TXML XML                                             │
│      ├── Validate: balance, tax code format, mandatory fields          │
│      └── Display preview                                               │
│      │                                                                 │
│      ▼                                                                 │
│  3. Digital Signing                                                    │
│      ├── Load digital certificate (default company signature)          │
│      ├── Sign XML with SHA-256-RSA (per Decree 254/2026)               │
│      ├── Verify signature                                               │
│      └── Status → SIGNED                                               │
│      │                                                                 │
│      ▼                                                                 │
│  4. GDT Submission                                                     │
│      ├── POST to GDT e-invoice API                                     │
│      ├── GDT validates and assigns invoice number                      │
│      ├── Receive GDT transaction ID (MA CUA CQT)                      │
│      └── Status → ISSUED                                               │
│      │                                                                 │
│      ▼                                                                 │
│  5. Post-Issuance                                                      │
│      ├── Journal entry creation (Dr AR, Cr Revenue, Cr VAT)            │
│      ├── Email invoice PDF to buyer                                    │
│      ├── Update sales ledger                                           │
│      └── Status → COMPLETED                                            │
│                                                                        │
│  Invoice types:                                                        │
│    - 01GTKT: Original                                                   │
│    - 02GTKT: Adjustment/substitution                                   │
│    - 04GTKT: Replacement                                              │
│    - 07GTKT: Cancel note                                               │
│                                                                        │
│  Timing: B2B → within 24h of delivery                                  │
│          B2C → immediately or at POS                                   │
│                                                                        │
└────────────────────────────────────────────────────────────────────────┘
```

---

## W-04: E-Invoice Cancellation/Adjustment

**Owner:** AR Accountant → Chief Accountant (approval for adjustment)

```
┌────────────────────────────────────────────────────────────────────────┐
│ E-INVOICE CANCELLATION/ADJUSTMENT                                       │
├────────────────────────────────────────────────────────────────────────┤
│                                                                        │
│  [Need to change issued invoice]                                       │
│      │                                                                 │
│      ├── [Cancel & Replace: Wrong buyer, amount, or tax rate]          │
│      │   │                                                             │
│      │   ▼                                                             │
│      │  1. Draft cancellation invoice (07GTKT)                         │
│      │  2. Obtain buyer confirmation (email/document)                  │
│      │  3. Submit cancellation to GDT                                  │
│      │  4. GDT approves cancellation                                   │
│      │  5. Issue new invoice (04GTKT replacement)                      │
│      │  6. Reverse original journal entry                              │
│      │  7. Create new journal entry                                    │
│      │                                                                 │
│      ├── [Adjustment: Partial correction]                              │
│      │   │                                                             │
│      │   ▼                                                             │
│      │  1. Draft adjustment invoice (02GTKT)                           │
│      │  2. Reference original invoice + GDT ID                         │
│      │  3. Show only the difference                                    │
│      │  4. Submit to GDT                                               │
│      │  5. Create adjustment journal entry                             │
│      │                                                                 │
│      └── [Buyer rejects cancellation]                                  │
│          └── Must contact GDT for resolution                           │
│                                                                        │
└────────────────────────────────────────────────────────────────────────┘
```

---

## W-05: Tax Audit from GDT

**Owner:** Chief Accountant + External support

```
┌────────────────────────────────────────────────────────────────────────┐
│ TAX AUDIT WORKFLOW                                                      │
├────────────────────────────────────────────────────────────────────────┤
│                                                                        │
│  [GDT Issues Audit Decision]                                           │
│      │                                                                 │
│      ▼                                                                 │
│  1. Notification Received                                              │
│      ├── Audit notice number, scope, period                            │
│      ├── Audit team info, scheduled dates                              │
│      └── Create audit case in GoTax                                    │
│      │                                                                 │
│      ▼                                                                 │
│  2. Document Preparation                                               │
│      ├── Export all declarations for audited period                    │
│      ├── Export underlying journal entries                             │
│      ├── Export e-invoices                                              │
│      ├── Export source documents (contracts, PO, receipts)             │
│      └── Prepare reconciliation reports                                │
│      │                                                                 │
│      ▼                                                                 │
│  3. Freeze Period (if not already)                                     │
│      ├── Lock declarations from editing                                │
│      ├── Lock period from new entries (if not done)                    │
│      └── System audit trail export                                     │
│      │                                                                 │
│      ▼                                                                 │
│  4. On-Site Audit Support                                              │
│      ├── Provide system access (read-only) to auditor                  │
│      ├── Generate tax audit trail report                               │
│      └── Respond to auditor queries                                    │
│      │                                                                 │
│      ▼                                                                 │
│  5. Audit Result                                                       │
│      ├── Receive audit record (Bien ban kiem tra)                      │
│      ├── If no violation → Close audit case                            │
│      ├── If tax adjustment →                                           │
│      │     Create amended declaration                                  │
│      │     Pay additional tax + interest                               │
│      └── Request re-inspection (if disagree)                           │
│      │                                                                 │
│      ▼                                                                 │
│  6. Close Audit Case                                                   │
│      ├── Store audit documentation                                     │
│      ├── Update company tax risk score                                 │
│      └── Notify management                                             │
│                                                                        │
└────────────────────────────────────────────────────────────────────────┘
```

---

## W-06: Tax Calendar & Alert Lifecycle

**Owner:** System (automated) + Tax Accountant

```
┌────────────────────────────────────────────────────────────────────────┐
│ TAX CALENDAR & ALERT LIFECYCLE                                          │
├────────────────────────────────────────────────────────────────────────┤
│                                                                        │
│  [Company Setup / Fiscal Year Start]                                   │
│      │                                                                 │
│      ▼                                                                 │
│  1. Calendar Generation                                                │
│      ├── Determine tax regime (TT99/TT133/TT58)                        │
│      ├── Determine declaration frequency per tax type                  │
│      ├── Generate annual calendar (deadlines per period)               │
│      └── Store in TaxCalendar table                                    │
│      │                                                                 │
│      ▼                                                                 │
│  2. Declaration Tracking                                               │
│      ├── Monitor period end dates                                      │
│      ├── Detect period without declaration                             │
│      └── Update calendar status per period                             │
│      │                                                                 │
│      ▼                                                                 │
│  3. Alert Triggers                                                     │
│      ├── D-14: "VAT declaration due in 14 days" (email)               │
│      ├── D-7:  "VAT declaration due in 7 days" (email + in-app)       │
│      ├── D-3:  "VAT declaration due in 3 days" (email + in-app + SMS) │
│      ├── D-0:  "VAT declaration DUE TODAY" (urgent)                    │
│      ├── D+1:  "VAT declaration OVERDUE" (escalation)                  │
│      └── D+7:  "VAT declaration overdue 7 days" (CFO notification)     │
│      │                                                                 │
│      ▼                                                                 │
│  4. Declaration Submitted                                              │
│      └── Status → SUBMITTED, clear alerts                              │
│      │                                                                 │
│      ▼                                                                 │
│  5. Payment Tracking                                                   │
│      ├── Payment due = Declaration deadline + (if payment separate)    │
│      ├── Alert if payment not made within N days of submission         │
│      └── Calculate late payment interest                               │
│                                                                        │
└────────────────────────────────────────────────────────────────────────┘
```

---

## UJ-01: Day-in-the-Life of a Tax Accountant

**Persona:** Ms. Lan, Tax Accountant at ABC Corp

### Morning (8:30-10:00)
1. Login to GoTax → Dashboard
2. Check tax calendar: 1 declaration due in 7 days (VAT Q1)
3. Review system-generated VAT calculation for Q1
4. Flag: Input VAT on fixed asset purchase may be ineligible
5. Manually adjust the amount, add note for audit trail

### Mid-day (10:00-11:30)
1. Open CIT provisional declaration for Q4
2. Auto-populate from year-end trial balance
3. Review: Revenue = 15.2B, Expenses = 12.8B, Profit = 2.4B
4. Check: Non-deductible expenses identified = 120M (fines + sponsorship)
5. Adjusted taxable income = 2.52B
6. CIT at 20% = 504M
7. Compare: Already paid 380M across Q1-Q3
8. Final Q4 due = 124M
9. Sign and submit to GDT

### Afternoon (13:30-15:00)
1. Create sales invoice for customer X (new export order)
2. Invoice: 10 items, export goods (VAT 0%), total $50,000 USD
3. System auto-calculates: VAT = 0, exchange rate = 25,450
4. Sign with digital certificate
5. Submit to GDT → Validated → Issued
6. Email PDF to buyer
7. Reconcile: Journal entry generated automatically

### End of Day (15:00-16:00)
1. Check GDT acknowledgement inbox
2. Yesterday's VAT declaration acknowledged (success)
3. Print acknowledgement and file
4. Review payment due dates for upcoming week
5. Generate payment order for CIT balance: 124M
6. Submit to bank via integrated EBPP

---

## UJ-02: Fiscal Year-End with Chief Accountant

**Persona:** Mr. Hung, Chief Accountant at ABC Corp

### Week 1 (Jan)
1. System notifies: CIT finalization period open (90 days remaining)
2. Instruct team to prepare financial statements
3. Review: CIT provisional payments = 380M (need ≥80% of final)

### Week 6 (Feb)
1. Financial statements ready (B01-DN, B02-DN, B03-DN, B09-DN)
2. External audit completed (mandatory per Decree 90/2025)
3. System auto-calculates final CIT: 504M
4. Provisionals (380M) < 80% of 504M (403M) → Late interest applies
5. Calculate: Underpayment = 504M - 380M = 124M
6. Late interest at 0.03%/day from Q4 payment due (30-Jan) to now (15-Feb) = 16 days = 595,200 VND

### Week 10 (Mar)
1. Generate 03/TNDN declaration
2. Review all adjustments with team
3. Sign with digital certificate
4. Submit to GDT → Acknowledged
5. Payment: 124M + interest = 124.595M
6. Complete CIT year-end checklist
7. Archive all documentation per Circular 99 Art. 28 (5-year retention)

### Lessons for System:
- Auto-alert for <80% provisional threshold mid-year (not just at year-end)
- Real-time monitoring would have saved 595K in interest
