# Tax Module — Tax Process Rules & Domain Logic

**Role:** BA Lead (20+ yrs) + Chief Accountant (20+ yrs)
**Date:** 2026-07-27

---

## 1. Core Tax Rules (Vietnamese Tax Law)

### 1.1 VAT Rules

| Rule ID | Rule | Legal Basis | Notes |
|---------|------|-------------|-------|
| VAT-01 | Standard VAT rate: 10% | VAT Law 48/2024/QH15 Art. 8 | |
| VAT-02 | Reduced VAT rate: 8% (Jul 2025 - Dec 2026) | Resolution (2025) | Temporary reduction |
| VAT-03 | Reduced VAT rate: 5% (certain goods/services) | VAT Law Art. 8(2) | Water, agri products, etc. |
| VAT-04 | VAT rate: 0% (export) | VAT Law Art. 8(1) | Exported goods/services |
| VAT-05 | Input VAT deductible only with valid e-invoice | Decree 123/2020 Art. 12 | Validate via GDT |
| VAT-06 | Input VAT ineligible: goods/services for personal use | VAT Law Art. 14 | |
| VAT-07 | Input VAT ineligible: non-cash payment >20M (from 2025: >5M) | Decree 320/2025 | Per Circular 20/2026 |
| VAT-08 | Input VAT on TSCD fully deductible | VAT Law Art. 14 | Including FA used for production |
| VAT-09 | VAT deduction method: mandatory for revenue >1B/year | VAT Law Art. 10 | Or choose direct method |
| VAT-10 | VAT refund: unused input VAT after 12 months | VAT Law Art. 13 | For investment projects |
| VAT-11 | VAT declaration monthly (default) or quarterly (revenue <50B) | Circular 80/2021 | |
| VAT-12 | Deadline: 20th of following month (monthly) | Circular 80/2021 Art. 44 | 30th for quarterly |
| VAT-13 | Zero declaration required even if no activity | Law on Tax Admin Art. 42 | Must not skip periods |

### 1.2 CIT Rules

| Rule ID | Rule | Legal Basis | Notes |
|---------|------|-------------|-------|
| CIT-01 | Standard CIT rate: 20% | CIT Law 67/2025/QH15 Art. 11 | |
| CIT-02 | Micro-enterprise rate: 15% (revenue <3B/year) | CIT Law Art. 11(2) | |
| CIT-03 | Small enterprise rate: 17% (revenue 3-50B/year) | CIT Law Art. 11(2) | |
| CIT-04 | CIT determined on taxable income = revenue - deductible expenses + other income | CIT Law Art. 5 | |
| CIT-05 | Non-deductible expenses include: fines, non-cash >5M penalties, sponsorship (unless specific) | CIT Law Art. 9 | Detailed in Circular 20 |
| CIT-06 | Tax loss carry-forward: max 5 years from loss year | CIT Law Art. 13 | Must have audited statements |
| CIT-07 | CIT incentives: tax holiday (4+9 years) or 10%/15% rate for eligible projects | CIT Law Art. 15-18 | New sectors/locations per 2025 law |
| CIT-08 | Quarterly provisional: 30th of next month | Circular 80/2021 Art. 17 | |
| CIT-09 | Annual finalization: 90 days from FY end | Circular 80/2021 Art. 17 | Usually 31-Mar |
| CIT-10 | Provisional CIT must be ≥80% of final CIT | Decree 320/2025 Art. 11 | Penalty if <80% |
| CIT-11 | Thin capitalization: interest deduction capped at 30% of EBITDA | Decree 132/2020 | Updated per Decree 255/2026 |
| CIT-12 | R&D super-deduction: up to 150% of qualifying expenditure | CIT Law Art. 9 | Per Circular 20/2026 |
| CIT-13 | Related-party transactions must be disclosed | Decree 255/2026 | Transfer pricing filing |
| CIT-14 | Capital transfer tax: 2% of gross proceeds (since 2025) | Decree 320/2025 Art. 11 | Replaces 20% on net gain |

### 1.3 PIT Rules

| Rule ID | Rule | Legal Basis |
|---------|------|-------------|
| PIT-01 | Resident: progressive rates 5-35% | PIT Law Art. 7 |
| PIT-02 | Non-resident: flat 20% | PIT Law Art. 7 |
| PIT-03 | Personal deduction: 11M/month (resident) | PIT Law Art. 9 |
| PIT-04 | Dependant deduction: 4.4M/month/per dependant | PIT Law Art. 9(2) |
| PIT-05 | Social insurance: 8% × gross (capped at 20× base salary) | Law on Social Insurance |
| PIT-06 | Health insurance: 1.5% × gross | Law on Health Insurance |
| PIT-07 | Unemployment insurance: 1% × gross (capped) | Law on Employment |
| PIT-08 | Foreigner: PIT finalization within 45 days of exit | Circular 80/2021 Art. 17 |
| PIT-09 | Employer declares for employee with single income source | PIT Law Art. 24 |
| PIT-10 | Annual finalization: 31-Mar (employer) / 30-Apr (individual) | Circular 80/2021 Art. 17 |

### 1.4 E-Invoice Rules

| Rule ID | Rule | Legal Basis |
|---------|------|-------------|
| EINV-01 | E-invoice mandatory since 01-Jul-2022 | Decree 123/2020 Art. 1 |
| EINV-02 | Invoice within 24h of delivery (B2B) | Decree 123/2020 Art. 8 |
| EINV-03 | Invoice immediately (B2C, retail, POS) | Decree 70/2025 Art. 1 |
| EINV-04 | POS connected to GDT (real-time data) | Decree 70/2025 Art. 2 |
| EINV-05 | Foreign suppliers must register via GDT portal | Decree 70/2025 Art. 4 |
| EINV-06 | Adjustment invoice must reference original GDT ID | Circular 78/2021 Art. 5 |
| EINV-07 | Cancellation requires buyer agreement (documented) | Decree 254/2026 Art. 4 |
| EINV-08 | Invoice data retained 10 years | Law on Tax Admin Art. 33 |
| EINV-09 | TXML format per GDT specification | GDT Decision |
| EINV-10 | Digital signature must match registered certificate | Decree 23/2025 Art. 24 |

### 1.5 Penalty Rules (Decree 310/2025/ND-CP)

| Violation | Penalty |
|-----------|---------|
| Late submission <30 days | 2-5M VND |
| Late submission 30-60 days | 5-8M VND |
| Late submission 60-90 days | 8-15M VND |
| Late submission >90 days | 15-25M VND |
| Incorrect declaration (self-corrected) | Warning |
| Incorrect declaration (GDT detected) | 10-20% of underpaid tax |
| Non-submission of declaration | 15-25M + forced assessment |
| Late payment interest | 0.03%/day of unpaid amount |

---

## 2. Domain Validation Rules

### 2.1 Tax Declaration Validation

```
ValidateDeclaration(decl):
  FOR each line:
    - Amount >= 0
    - LineCode exists in form definition
    - If source = AUTO_CALCULATED: SourceAccount must be valid account
    - If source = FROM_LEDGER: SourceEntryIDs must exist and belong to period

  IF decl.DeclarationType == GTGT01:
    - [16] = [14] + [15]
    - [23] = [21] + [22]
    - [30] = [23] - [16]
    - XOR: [31] > 0 OR [32] > 0 (never both)
    - Purchase ledger total VAT = [16]
    - Sales ledger total VAT = [23]

  IF decl.DeclarationType == TNDN03:
    - [04] >= [06] (revenue >= taxable revenue)
    - [14] <= [12] × [13] (tax <= taxable income × rate)
    - [17] > 0: provisional payments made
    - IF [18] > 0 AND [17] < [16] × 0.8: Flag late interest

  IF decl.DeclarationType == KK_TNCN:
    - [03] >= [09] (taxable income >= after deductions)
    - [11] <= [05] / (4.4M × months) (dependant count consistent)
    - [10] matches progressive bracket result

  GENERAL:
    - No duplicate declaration for same period/type
    - Company must have valid tax code
    - Period must be within open fiscal year
    - If amending: PreviousDeclarationID must reference ACKNOWLEDGED declaration
```

### 2.2 E-Invoice Validation

```
ValidateEInvoice(invoice):
  IF invoice.BuyerTaxCode != "":
    Must be 10 or 13 digits (or 12 digits for new format)
    Must pass Luhn check (for validation)

  IF invoice.InvoiceType == ADJUSTMENT:
    OriginalInvoiceID required
    Must reference ISSUED invoice
    Line amounts can be positive (increase) or negative (decrease)

  IF invoice.InvoiceType == REPLACEMENT:
    OriginalInvoiceID required
    Total must differ from original

  IF invoice.InvoiceType == CANCELLATION_NOTE:
    OriginalInvoiceID required
    All amounts must be negative of original

  For all invoices:
    Sum of line VAT amounts = Total VAT
    Grand Total = Subtotal + Total VAT
    Exchange rate required if currency ≠ VND
    Buyer name required (not blank)
    At least 1 line item
```

### 2.3 Payment Validation

```
ValidatePayment(payment):
  IF payment.DeclaredAmount < payment.PaidAmount:
    Flag: Overpayment (overpaid amount qualifies for refund)
    Status = OVERPAID

  IF payment.DeclaredAmount > payment.PaidAmount:
    Flag: Underpayment
    Calculate late days = Today - payment.DueDate
    LateInterest = UnderpaidAmount × 0.03% × LateDays
    Status = PARTIAL (if paid > 0) or PENDING (if paid = 0)

  IF payment.DeclaredAmount == payment.PaidAmount:
    Status = PAID
```

---

## 3. Accounting Integration Rules

### 3.1 VAT Journal Entry Generation

```
On output VAT (from sales):
  Dr 131/111/112  Total
    Cr 511          Subtotal (before VAT)
    Cr 33311        VAT amount (output)

On input VAT (from purchases):
  Dr 156/152/211  Subtotal
  Dr 1331/1332    VAT amount (input)
    Cr 331/111/112  Total

VAT Accounts Mapping:
  Output VAT → 33311 (credit balance = liability to tax authority)
  Input VAT  → 1331 (debit balance = receivable from tax authority)
  VAT payable = Cr(33311) - Dr(1331) → if positive, journal to 333
```

### 3.2 CIT Journal Entry Generation

```
Quarterly provisional:
  Dr 8211     CIT expense (current)
    Cr 3334   CIT payable

Annual finalization adjustment:
  If additional tax due:
    Dr 8211     Additional expense
      Cr 3334   Additional payable

  If overpaid (refund):
    Dr 3334     Overpayment
      Cr 8211   Reduce expense

Global minimum tax (Pillar 2):
  Dr 82112   Top-up tax expense
    Cr 3334   Top-up tax payable
```

### 3.3 Period Closing Rules

```
Before closing tax period:
  - All journal entries for period must be POSTED
  - VAT declaration must be SUBMITTED (or flagged as late)
  - PIT declaration must be SUBMITTED (if applicable)
  - All e-invoices for period must be ISSUED (or CANCELLED)

Period close validation:
  CAN_CLOSE(PERIOD):
    FOREACH tax type applicable to company:
      IF declaration exists AND status in (SUBMITTED, ACKNOWLEDGED): OK
      IF no declaration: WARNING (zero declaration required)
    FOREACH sales entry with invoice not issued: WARNING

  Period cannot be closed if:
    - Any tax deadline has passed without declaration
    - GDT submission pending (if electronic submission enabled)
```

---

## 4. Tax Calendar Generation Rules

```
GenerateTaxCalendar(company, year):
  FOR each tax type applicable to company:
    IF tax type == VAT:
      IF company.revenue < 50B:
        period = QUARTERLY
        deadline = last day of month following quarter
      ELSE:
        period = MONTHLY
        deadline = 20th of following month

    IF tax type == CIT:
      period = QUARTERLY (provisional)
      deadline = 30th of month following quarter (Q1-Q3)
      annual deadline = 31-Mar (or 90 days from FY end)
      annual payment deadline = same as filing

    IF tax type == PIT:
      IF total PIT/month < 50M:
        period = QUARTERLY
      ELSE:
        period = MONTHLY
      deadline = 20th of following month (monthly)
      annual deadline = 31-Mar (employer) / 30-Apr (individual)

    IF tax type == TTDB:
      period = MONTHLY
      deadline = 20th of following month

    IF tax type == BVMT:
      period = MONTHLY
      deadline = 20th of following month

    IF tax type == FCT:
      period = PER_OCCURRENCE
      deadline = 20th of month following payment

  FOR each generated period:
    IF deadline falls on weekend/holiday:
      deadline = next working day (per Law on Tax Admin Art. 44)

  HANDLE special cases:
    Q4 deadline for quarterly VAT: 30-Apr (not 30-Jan)
    CIT provisional Q4 overdue if not paid by 30-Jan of next year
    PIT annual deadline for calendar year: 31-Mar (employer), 30-Apr (individual)
    Different FY: recalculate based on FY end date
```

---

## 5. GDT Submission Rules

```
SubmitDeclaration(declaration, certificate):
  PRE-CHECK:
    - Certificate not expired (validFrom < now < validTo)
    - Certificate not revoked (Status == ACTIVE)
    - Declaration status == VALIDATED
    - Company integration profile == ACTIVE

  For USB Token certificates:
    - Require user to physically insert token
    - Prompt for PIN
    - PIN retry limit: 3, then lock token

  For Remote HSM:
    - Use stored credential or SSO
    - Rate limit: 5 submissions/minute per company

  XML generation:
    - Load form template (form-specific XSLT)
    - Map line items → indicator codes
    - Validate against GDT XSD
    - Wrap in SOAP envelope if required (HTKK protocol)

  Submission:
    - POST to GDT endpoint with mTLS
    - Timeout: 120 seconds
    - Retry: 3 times with exponential backoff (1s, 5s, 30s)
    - On HTTP 200: Parse response XML
    - On HTTP 4xx: Parse error code, return user-friendly message
    - On HTTP 5xx: Retry, if persistent → queue for manual intervention

  Response handling:
    SUCCESS (acknowledgement received):
      - Extract MaThongBao (acknowledgement number)
      - Update declaration status = ACKNOWLEDGED
      - Store response XML
      - Notify user

    FAILURE (rejected):
      - Extract MaLoi (error code) + NoiDung (error message)
      - Map error code to user-friendly message
      - Update declaration status = REJECTED
      - Allow user to edit and resubmit
```
