# Tax Module — Use Cases

**Role:** BA Lead (20+ yrs) + Chief Accountant (20+ yrs)
**Date:** 2026-07-27

---

## UC-01: Monthly VAT Declaration (Deduction Method)

**Actor:** Tax Accountant
**Trigger:** End of tax period (monthly/quarterly)

### Happy Path
1. System detects period end (20th of following month for monthly, 30th for quarterly)
2. Tax Accountant navigates to Tax → VAT Declaration → Create New
3. System selects company, period, form type (01/GTGT)
4. System auto-populates from journal entries:
   - Output VAT from credit side of 3331-linked transactions
   - Input VAT from debit side of 133-linked transactions
   - Purchase/sales ledgers from account code analysis
5. Tax Accountant reviews auto-populated figures
6. Tax Accountant adds manual adjustments if needed (e.g., ineligible input VAT)
7. System validates: debit totals = credit totals, no negative amounts
8. System generates XML file per Circular 80 XSD
9. System displays declaration preview
10. Tax Accountant clicks "Submit to GDT"
11. System prompts for digital signature (USB Token / remote HSM)
12. Tax Accountant signs and confirms
13. System POSTs XML to `thuedientu.gdt.gov.vn`
14. GDT returns acknowledgement receipt (XML with MA THONG BAO)
15. System stores acknowledgement, updates status to SUBMITTED
16. Notification sent to Chief Accountant

### Alternative Paths
- **A1: GDT rejects declaration** → Parse error code, display to user, allow edit + resubmit
- **A2: Auto-population incomplete** → Highlight missing accounts, allow manual entry
- **A3: Digital signature expired** → Block submission, prompt to renew signature
- **A4: Network error during submission** → Queue for retry, notify user

### Exception Paths
- **E1: Period already declared** → Error: "Declaration for period X already submitted on date Y"
- **E2: Period not yet closed in GL** → Warning: "Period not closed. Declarations from open periods may change."
- **E3: No journal entries for period** → Allow zero-declaration (must still submit)
- **E4: GDT system maintenance** → Queue submission, retry when available
- **E5: Certificate revoked** → Block all submissions, urgent notification

### Business Rules
- BR-VAT-01: VAT period is monthly (default) or quarterly (if annual revenue <VND50B)
- BR-VAT-02: Declaration deadline: 20th of following month (monthly)
- BR-VAT-03: Purchase ledger (Bang ke mua vao) must reconcile with input VAT
- BR-VAT-04: Sales ledger (Bang ke ban ra) must reconcile with output VAT
- BR-VAT-05: Zero declaration still required if no transactions
- BR-VAT-06: e-Invoice data and VAT declaration must reconcile per GDT requirement

---

## UC-02: Annual CIT Finalization

**Actor:** Chief Accountant
**Trigger:** Fiscal year end (90 days deadline: 31-Mar following year)

### Happy Path
1. System notifies Chief Accountant of upcoming CIT deadline (Jan notification)
2. Chief Accountant navigates to Tax → CIT → Finalization
3. System selects fiscal year, company
4. System auto-populates:
   - Revenue from relevant account codes (511, 515, 711)
   - Expenses from relevant accounts (6xx, 8xx excluding 821)
   - CIT provisional paid during year (from 3334 payments)
   - Tax losses carried forward (accumulated from prior years)
5. System calculates:
   - Taxable income = Revenue - Deductible expenses + Non-deductible adjustments
   - CIT payable = Taxable income × Rate (after incentives)
   - CIT refundable = Provisional paid - Final CIT
6. Chief Accountant reviews and adjusts:
   - Non-deductible expenses (e.g., fines, penalties)
   - Tax incentive eligibility (sector, location)
   - Related-party transactions
7. System generates 03/TNDN form + appendices
8. System validates against financial statements (B01-DN, B02-DN)
9. Chief Accountant signs with digital signature
10. System submits to GDT
11. GDT returns acknowledgement

### Alternative Paths
- **A1: Tax loss carry-forward** → Auto-detect loss year, apply per BR-CIT-04
- **A2: CIT incentive claim** → Prompt for investment documentation, verify eligibility
- **A3: Transfer pricing disclosure** → Flag if related-party transactions > threshold
- **A4: Global minimum tax (Pillar 2)** → Generate GMT supplementary form (account 82112)

### Exception Paths
- **E1: Financial statements not audited** → Warning for mandatory-audit companies
- **E2: Provisional CIT <80% of final CIT** → Late payment interest calculation per Decree 310
- **E3: Negative CIT (refund case)** → Flag for tax refund procedure

### Business Rules
- BR-CIT-01: CIT rate: 20% standard, 15% micro (revenue <VND3B), 17% small (VND3-50B)
- BR-CIT-02: Quarterly provisional deadline: 30th of month following quarter
- BR-CIT-03: Annual finalization deadline: 90 days from FY end
- BR-CIT-04: Tax losses carry forward max 5 years (per CIT Law 2025 Art. 13)
- BR-CIT-05: Total provisional CIT ≥ 80% of final CIT (penalty if <80%)
- BR-CIT-06: CIT incentives require registration and documentation per Circular 20

---

## UC-03: Monthly PIT Declaration

**Actor:** Payroll Accountant
**Trigger:** End of tax period (20th of following month)

### Happy Path
1. Payroll Accountant navigates to Tax → PIT → Monthly Declaration
2. System selects company, period, form type (05/KK-TNCN)
3. System loads employee list with PIT data from payroll
4. Tax Accountant reviews per-employee:
   - Gross income, deductions, dependants
   - PIT calculated per progressive rate table
   - Previous period cumulative amounts
5. System generates 05/KK-TNCN XML
6. Submit to GDT per UC-01 step 10-14

### Alternative Paths
- **A1: Employee has personal tax code** → Validate MST before submission
- **A2: Employee without tax code** → Flag, registration required before filing
- **A3: Expatriate finalization on exit** → Trigger special finalization flow

### Business Rules
- BR-PIT-01: PIT deadline: 20th of following month (monthly), 30th (quarterly)
- BR-PIT-02: Employee with single income source → employer declares
- BR-PIT-03: Employee with multiple income sources → self-declare at year-end
- BR-PIT-04: Dependant deduction: VND4.4M/month/dependant (2024 rate)
- BR-PIT-05: Personal deduction: VND11M/month (resident individuals)
- BR-PIT-06: PIT finalization for expatriates: within 45 days of exit

---

## UC-04: E-Invoice Issuance

**Actor:** AR Accountant
**Trigger:** Sale transaction recorded (goods shipped/service delivered)

### Happy Path
1. AR Accountant creates sales invoice from sales order or standalone
2. System auto-populates:
   - Buyer info (name, tax code, address)
   - Line items (description, quantity, unit price, VAT rate, VAT amount)
   - Totals (subtotal, VAT, grand total)
   - Invoice pattern (from company's registered patterns)
3. AR Accountant selects invoice type:
   - Original (01GTKT)
   - Adjustment (02GTKT)
   - Replacement (04GTKT)
   - Cancellation (07GTKT - note)
4. System validates:
   - Buyer tax code format (10 or 13 digits)
   - Line items total matches VAT calculation
   - Invoice serial from registered pattern
5. AR Accountant clicks "Issue"
6. System generates TXML XML per GDT format (Decree 254/2026/ND-CP)
7. System signs with company's default digital signature
8. System submits to GDT
9. GDT validates and returns:
   - Invoice number (GDT-assigned)
   - GDT transaction ID (MA CUA CQT)
   - Validation timestamp
10. System stores issued invoice
11. System sends to buyer (email with PDF attachment, or portal link)
12. Journal entry auto-created (Dr AR, Cr Revenue, Cr Output VAT)

### Alternative Paths
- **A1: Buyer is individual** → No tax code required, use personal ID
- **A2: Buyer requests invoice later** → Create from delivery note
- **A3: Exchange rate for foreign currency** → Apply GDT-published rate

### Exception Paths
- **E1: GDT rejects invoice (duplicate reference)** → Generate new reference, retry
- **E2: GDT system unavailable** → Queue invoice, retry with timestamp
- **E3: Buyer tax code invalid in GDT system** → Verify with ETax portal, allow correction
- **E4: Digital signature error** → Block issuance, notify IT

### Business Rules
- BR-EINV-01: E-invoice mandatory for all B2B/B2C since 01-Jul-2022
- BR-EINV-02: Invoice must be issued within 24h of delivery (B2B) / immediately (B2C)
- BR-EINV-03: POS invoices connected to GDT in real-time (Decree 70/2025)
- BR-EINV-04: Adjustment invoice references original invoice number + GDT transaction ID
- BR-EINV-05: Cancel invoice requires buyer agreement and GDT notification
- BR-EINV-06: Foreign suppliers must register and issue via GDT portal

---

## UC-05: Tax Payment & Reconciliation

**Actor:** Tax Accountant
**Trigger:** After declaration submission AND payment due date

### Happy Path
1. System retrieves declared amount from submitted declaration
2. System marks payment due date per tax type
3. Tax Accountant navigates to Tax → Payment → Create Payment Order
4. System generates payment order with:
   - Payee: Kho bac Nha nuoc (State Treasury)
   - Amount: Declared tax amount
   - Reference: Declaration acknowledgement number
   - Tax account: Correct revenue code (NDKT)
5. Tax Accountant submits payment order to bank (via integrated EBPP)
6. Bank returns payment confirmation
7. System records payment, reconciles against declaration
8. Status updates: PAID

### Alternative Paths
- **A1: Partial payment** → Track underpayment, calculate late interest
- **A2: Overpayment** → Flag for refund or offset against future liability
- **A3: Payment after deadline** → Auto-calculate late payment interest (0.03%/day per Decree 310)

---

## UC-06: Tax Audit Support

**Actor:** Chief Accountant / External Auditor
**Trigger:** Tax audit notification from GDT

### Happy Path
1. Chief Accountant navigates to Tax → Audit → New Audit Case
2. Records audit information (audit period, auditor, notice number)
3. System freezes declaration edits for audited period
4. Auditor reviews:
   - All submitted declarations with version history
   - Underlying journal entries (drill-down from declaration line)
   - Supporting documents (e-invoices, contracts)
   - Audit trail (who submitted, when, changes)
5. Auditor records audit findings
6. If adjustment needed → create amended declaration
7. Audit case closed

### Business Rules
- BR-AUD-01: Tax audit can cover up to 10 years (violations) or 5 years (no violations)
- BR-AUD-02: Declaration data frozen during audit period
- BR-AUD-03: Amended declaration after audit requires special approval

---

## UC-07: Tax Calendar & Deadline Management

**Actor:** Tax Accountant
**Trigger:** New company setup / Fiscal year start

### Happy Path
1. System detects company's tax regime (TT99/TT133/TT58)
2. System determines declaration frequency per tax type:
   - VAT: Monthly (default) or Quarterly (if revenue <VND50B)
   - CIT: Quarterly provisional, annual finalization
   - PIT: Monthly (default) or Quarterly (<VND50M/month)
   - TTDB: Monthly
   - BVMT: Monthly
   - FCT: Per payment or quarterly
3. System generates annual tax calendar per company
4. System sends dashboard alerts at configurable intervals:
   - 14 days before deadline (email)
   - 7 days before deadline (email + in-app)
   - 3 days before deadline (email + in-app + SMS)
   - Deadline day (urgent notification)
5. Missed deadline: escalation chain (Tax Accountant → Chief Accountant → CFO)

### Business Rules
- BR-CAL-01: Monthly declarations due by 20th of following month
- BR-CAL-02: Quarterly declarations due by 30th of following quarter (or 30th in Apr)
- BR-CAL-03: Annual CIT finalization: 31-Mar (calendar year), 90 days from FY end
- BR-CAL-04: Annual PIT finalization: 31-Mar (employer), 30-Apr (individual)
- BR-CAL-05: Deadline falling on holiday → next working day (per tax admin law)
