# Tax Module — Data Flows & Integration Architecture

**Role:** BA Lead (20+ yrs) + Chief Accountant (20+ yrs)
**Date:** 2026-07-27

---

## 1. System Context Diagram

```
┌──────────────────────────────────────────────────────────────────────────┐
│                           GOTAX SYSTEM                                    │
│                                                                          │
│  ┌──────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐ │
│  │ GL       │  │ Tax Module   │  │ E-Invoice    │  │ Tax Calendar    │ │
│  │ Module   │◄─┤(NEW)         │  │ Module(NEW)  │  │ Module(NEW)     │ │
│  │(EXISTING)│  │              │  │              │  │                  │ │
│  └────┬─────┘  └──────┬───────┘  └──────┬───────┘  └────────┬─────────┘ │
│       │               │                  │                   │          │
│       └───────────────┴──────────────────┴───────────────────┘          │
│                               │                                         │
└───────────────────────────────┼─────────────────────────────────────────┘
                                │
            ┌───────────────────┼───────────────────┐
            │                   │                   │
            ▼                   ▼                   ▼
    ┌───────────────┐  ┌───────────────┐  ┌───────────────┐
    │  GDT (Tax)    │  │  Bank         │  │  Buyer        │
    │  thuedientu   │  │  EBPP/API     │  │  Email/Portal │
    │  .gdt.gov.vn  │  │               │  │               │
    └───────────────┘  └───────────────┘  └───────────────┘
```

---

## 2. Data Flow: VAT Declaration Generation

```
┌──────────────────────────────────────────────────────────────────────┐
│ VAT DECLARATION GENERATION DATA FLOW                                  │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  GL Journal Entries (POSTED)                                         │
│       │                                                              │
│       ▼                                                              │
│  [Step 1: Extract Output VAT]                                        │
│  Query: SELECT lines FROM journal_entries                            │
│  WHERE period = target AND account_code IN ('511','515','711',etc)   │
│    JOIN account_balances WHERE account LIKE '33311%'                 │
│  → Output VAT = Cr(33311) - Dr(33311)                                │
│       │                                                              │
│       ▼                                                              │
│  [Step 2: Extract Input VAT]                                         │
│  Query: SELECT lines FROM journal_entries                            │
│  WHERE period = target AND account_code IN ('1331','1332',etc)       │
│  → Input VAT = Dr(1331) + Dr(1332)                                   │
│       │                                                              │
│       ▼                                                              │
│  [Step 3: Extract Purchase Ledger]                                   │
│  Query: Purchase-related entries (331, 111, 112, 156, 152)           │
│  WITH VAT info from linked e-invoices or VAT line items              │
│  → Purchase ledger: supplier, invoice, amount, VAT rate, VAT amount  │
│       │                                                              │
│       ▼                                                              │
│  [Step 4: Extract Sales Ledger]                                      │
│  Query: Sales-related entries (131, 111, 112, 511)                   │
│  WITH VAT info from linked e-invoices                                │
│  → Sales ledger: buyer, invoice, amount, VAT rate, VAT amount        │
│       │                                                              │
│       ▼                                                              │
│  [Step 5: Calculate Declaration Lines]                               │
│  Indicator [21] = Output VAT from sales                               │
│  Indicator [22] = Adjustments (manual or correction entries)          │
│  Indicator [23] = [21] + [22]                                        │
│  Indicator [24] = Input VAT from purchases (1331)                    │
│  Indicator [25] = Input VAT from fixed assets (1332)                  │
│  Indicator [30] = [23] - [24] - [25]                                 │
│  Indicator [31] = [30] if > 0 (VAT payable)                          │
│  Indicator [32] = -[30] if < 0 (VAT carried forward/refundable)      │
│       │                                                              │
│       ▼                                                              │
│  [Step 6: Generate XML]                                              │
│  Load 01/GTGT XSLT template                                          │
│  Map indicators → XML indicator elements                              │
│  Generate purchase ledger XML                                         │
│  Generate sales ledger XML                                            │
│  Validate against XSD schema                                          │
│  → Declaration XML ready                                              │
│       │                                                              │
│       ▼                                                              │
│  [Step 7: Store & Submit]                                            │
│  Store XML in TaxDeclaration.DeclarationXML                           │
│  Status → VALIDATED                                                   │
│  Await user review + signing                                          │
│                                                                      │
└──────────────────────────────────────────────────────────────────────┘
```

---

## 3. Data Flow: E-Invoice Issuance

```
┌──────────────────────────────────────────────────────────────────────┐
│ E-INVOICE ISSUANCE DATA FLOW                                          │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Sales Order / Manual Entry                                          │
│       │                                                              │
│       ▼                                                              │
│  [Step 1: Create Invoice Data]                                       │
│  Data entry: buyer info, line items, VAT rates                       │
│  Select invoice pattern (from company's EInvoicePattern)              │
│  Select digital signature (from company's DigitalSignature)           │
│  Calculate: subtotal, VAT by line, totals                             │
│  Default currency: VND, exchange rate if foreign                     │
│       │                                                              │
│       ▼                                                              │
│  [Step 2: Validate]                                                  │
│  - Buyer tax code format (10/13 digits)                              │
│  - Line totals match VAT calculation                                 │
│  - Valid invoice pattern + serial                                    │
│  - Digital signature not expired                                     │
│       │                                                              │
│       ▼                                                              │
│  [Step 3: Generate TXML]                                             │
│  Template: TXML invoice per GDT spec                                 │
│  Fields: header, buyer, line items, totals, payment terms             │
│  Metadata: pattern, serial, date/time                                │
│  → Raw XML                                                           │
│       │                                                              │
│       ▼                                                              │
│  [Step 4: Sign XML]                                                  │
│  Load certificate (USB Token / remote HSM)                           │
│  Hash XML with SHA-256                                                │
│  Sign hash with RSA private key                                      │
│  Embed signature in XML: <ChuKySo>...</ChuKySo>                      │
│  → Signed XML                                                        │
│       │                                                              │
│       ▼                                                              │
│  [Step 5: Submit to GDT]                                             │
│  HTTP POST to GDT e-invoice API                                      │
│  mTLS with client certificate                                        │
│  GDT validates and assigns:                                          │
│    - Invoice number (if approved)                                    │
│    - MA CUA CQT (GDT transaction ID)                                 │
│  → GDT Response XML                                                  │
│       │                                                              │
│       ▼                                                              │
│  [Step 6: Post-Issuance]                                             │
│  Store issued invoice in DB                                          │
│  Generate journal entry:                                              │
│    Dr 131 (AR) = Grand Total                                         │
│      Cr 511 (Revenue) = Subtotal                                     │
│      Cr 33311 (Output VAT) = VAT Amount                              │
│  Send invoice PDF to buyer (email)                                   │
│  Update invoice status = ISSUED                                      │
│       │                                                              │
│       ▼                                                              │
│  [Step 7: VAT Ledger Update]                                         │
│  Update sales ledger (Bang ke ban ra)                                │
│  Auto-include in next VAT declaration                                │
│                                                                      │
└──────────────────────────────────────────────────────────────────────┘
```

---

## 4. Data Flow: GDT Submission + Acknowledgement

```
┌──────────────────────────────────────────────────────────────────────┐
│ GDT SUBMISSION & ACKNOWLEDGEMENT DATA FLOW                            │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  TaxAccountant                                                       │
│       │                                                              │
│       ▼                                                              │
│  [Submit Declaration]                                                │
│       │                                                              │
│       ▼                                                              │
│  ┌─────────────────────────────────────────────────────────┐        │
│  │                GDT API Client                            │        │
│  │                                                          │        │
│  │  1. Load company certificate & GDT endpoint URL          │        │
│  │  2. Establish mTLS connection                            │        │
│  │  3. POST /api/submission/declare                        │        │
│  │     Headers:                                             │        │
│  │       Content-Type: application/xml                     │        │
│  │       X-GDT-Cert-Serial: <cert serial>                  │        │
│  │     Body: Declaration XML                                │        │
│  │                                                          │        │
│  │  4. Receive Response                                     │        │
│  │     HTTP 200:                                            │        │
│  │       <BK:KetQua>                                        │        │
│  │         <BK:MaSoThue>0123456789</BK:MaSoThue>           │        │
│  │         <BK:MaThongBao>GDT20260415ABC123</BK:MaThongBao>│        │
│  │         <BK:MaKetQua>00</BK:MaKetQua>                   │        │
│  │         <BK:NgayNhan>2026-04-15</BK:NgayNhan>           │        │
│  │       </BK:KetQua>                                       │        │
│  │                                                          │        │
│  │     HTTP 200 (rejection):                                │        │
│  │       <BK:KetQua>                                        │        │
│  │         <BK:MaKetQua>01</BK:MaKetQua>                   │        │
│  │         <BK:MaLoi>ERR-015</BK:MaLoi>                    │        │
│  │         <BK:NoiDung>XML khong dung dinh dang</BK:NoiDung>│        │
│  │       </BK:KetQua>                                       │        │
│  │                                                          │        │
│  └─────────────────────────────────────────────────────────┘        │
│       │                                                              │
│       ├── Success ──→ Update declaration:                            │
│       │                Status = ACKNOWLEDGED                         │
│       │                AcknowledgementRef = MaThongBao               │
│       │                GDTResponseXML = full response                │
│       │                Notify user: green success, show ref #        │
│       │                                                              │
│       └── Rejected ──→ Update declaration:                           │
│                         Status = REJECTED                            │
│                         GDTResponseXML = full response               │
│                         Notify user: red error, show message         │
│                         Allow edit + resubmit                        │
│                                                                      │
│  [Error Handling]                                                    │
│  HTTP 401/403: Certificate invalid → Notify IT, block                │
│  HTTP 408/502/503: GDT system busy → Retry 3x (1s, 5s, 30s)         │
│    if all fail: Queue for manual retry, notify user                  │
│  Network timeout: Same retry logic                                   │
│                                                                      │
└──────────────────────────────────────────────────────────────────────┘
```

---

## 5. Data Flow: Tax Calendar & Alerts

```
┌──────────────────────────────────────────────────────────────────────┐
│ TAX CALENDAR & ALERT DATA FLOW                                        │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  [System Cron: Daily at 00:00]                                       │
│       │                                                              │
│       ▼                                                              │
│  For each active company:                                            │
│       │                                                              │
│       ▼                                                              │
│  [Check Calendar Status]                                             │
│  SELECT * FROM tax_calendar                                          │
│  WHERE company_id = ? AND year = ? AND status = 'PENDING'            │
│       │                                                              │
│       ▼                                                              │
│  [Calculate Days Until Deadline]                                     │
│  days_remaining = DATEDIFF(day, TODAY, declaration_due)              │
│       │                                                              │
│       ├── days == 14 → CREATE ALERT                                  │
│       │   Type: INFO, Channel: EMAIL                                 │
│       │   Template: "VAT Q1/2026 due in 14 days"                    │
│       │                                                              │
│       ├── days == 7 → CREATE ALERT                                   │
│       │   Type: WARNING, Channel: EMAIL + IN_APP                     │
│       │                                                              │
│       ├── days == 3 → CREATE ALERT                                   │
│       │   Type: CRITICAL, Channel: EMAIL + IN_APP + SMS              │
│       │                                                              │
│       ├── days == 0 → CREATE ALERT                                   │
│       │   Type: DUE_TODAY, Channel: ALL                              │
│       │                                                              │
│       ├── days < 0 → UPDATE STATUS = OVERDUE                         │
│       │   CREATE ESCALATION ALERT                                    │
│       │   Escalate to: Chief Accountant                              │
│       │   If no action in 7 days: Escalate to CFO                    │
│       │                                                              │
│       └── Declaration submitted → UPDATE STATUS = SUBMITTED          │
│           CLEAR all alerts for this period                           │
│                                                                      │
│  [SMS Integration]                                                   │
│  Via SMS provider API (e.g., Vietguys, Infobip)                      │
│  Template: "GoTax alert: [TaxType] [Period] due [Date]"              │
│                                                                      │
└──────────────────────────────────────────────────────────────────────┘
```

---

## 6. Integration: External Systems

### 6.1 GDT (thuedientu.gdt.gov.vn)

```
Protocol: HTTPS + mTLS (mutual TLS with client certificate)
Certificate: USB Token (physical) or Remote HSM (cloud)
Authentication: Certificate-based (no username/password)

Environments:
  PROD: https://thuedientu.gdt.gov.vn
  UAT:  https://thuedientunth.gdt.gov.vn
  DEV:  https://thuedientudd.gdt.gov.vn (if available)

Rate Limits:
  - Declaration submission: 10/minute per company
  - Invoice issuance: 50/minute per company
  - Information queries: 100/minute per company

Timeouts:
  - Connection: 30 seconds
  - Read: 60 seconds
  - Total: 120 seconds

Retry Policy:
  - Network errors: 3 attempts (1s, 5s, 30s exponential backoff)
  - Server errors (5xx): 3 attempts (same backoff)
  - Client errors (4xx): No retry, user intervention required

Error Codes:
  See TAX_SPECS.md Section 4.2
```

### 6.2 HTKK Software Integration

```
Legacy support for HTKK (Hỗ trợ kê khai) file format:
  - Export: XML file in HTKK-compatible format
  - Import: HTKK-generated XML declaration files
  - User downloads from GoTax, opens in HTKK, signs, submits via HTKK

This is fallback only. Primary channel is direct GDT API.
```

### 6.3 Bank Integration (Tax Payment)

```
Integration method: EBPP (Electronic Bill Presentment and Payment)
  or: Corporate banking API (individual bank integrations)

Supported banks (target):
  - Vietcombank (VCB)
  - VietinBank (CTG)
  - BIDV (BID)
  - Techcombank (TCB)
  - ACB
  - MB Bank (MB)

Data sent:
  - Payment order with NSNN (State Budget) account
  - Tax type revenue code (Ma NDKT)
  - Amount, reference declaration number
  - Company name, tax code

Data received:
  - Payment confirmation reference
  - Payment timestamp
  - Bank fee (if any)
```

### 6.4 SMS / Notification Integration

```
Email: SMTP (configurable, e.g., SendGrid, AWS SES)
SMS: Third-party SMS API (e.g., Vietguys, Twilio, Infobip)
In-app: WebSocket or server-sent events for real-time dashboard

Alert channels per severity:
  INFO:     Email only
  WARNING:  Email + In-app notification
  CRITICAL: Email + In-app + SMS
  OVERDUE:  Email + In-app + SMS + escalation
```

---

## 7. Database Schema Migration Plan

### 7.1 New Tables

```
tax_rates
  - id, tax_type, rate_code, rate_name, rate_type, rate_value
  - effective_from, effective_to, is_active, legal_ref
  - created_at

tax_rate_brackets
  - id, tax_rate_id, min_amount, max_amount, rate_percent, flat_amount

tax_declarations
  - id, company_id, declaration_type, tax_period, period_year, period_number
  - status, submitted_at, submitted_by, acknowledged_at
  - acknowledgement_ref, gdt_response_xml, declaration_xml
  - previous_declaration_id, adjustment_type, version
  - created_at, created_by, updated_at

tax_declaration_lines
  - id, declaration_id, line_code, line_name, amount
  - source_type, source_account, source_entry_ids, note, sort_order

tax_payments
  - id, company_id, declaration_id, tax_type
  - period_year, period_number, declared_amount, paid_amount
  - payment_date, due_date, payment_ref, payment_method
  - status, late_payment_days, late_interest, notes

e_invoices
  - id, company_id, pattern, serial, invoice_number
  - invoice_type, gdt_transaction_id, buyer_* (name, tax_code, address, email)
  - currency_code, exchange_rate, subtotal, vat_amount, grand_total
  - xml_body, signed_xml, issue_date, signing_date
  - digital_signature_id, journal_entry_id
  - status, cancelled_at, cancel_reason, original_invoice_id
  - gdt_response, created_at

e_invoice_lines
  - id, e_invoice_id, line_number, description, unit, quantity
  - unit_price, line_total, vat_rate, vat_amount

tax_calendar
  - id, company_id, tax_type, period_type
  - period_year, period_number, start_date, end_date
  - declaration_due, payment_due, status
  - declaration_id, created_at

tax_alerts
  - id, company_id, tax_calendar_id, alert_type, channel
  - message, sent_at, acknowledged_at, acknowledged_by

tax_audit_cases
  - id, company_id, audit_period_start, audit_period_end
  - audit_decision_number, auditor_name, auditor_contact
  - status, findings, penalty_amount, created_at, closed_at
```

### 7.2 Existing Table Modifications

```sql
-- Add tax fields to companies
ALTER TABLE companies ADD COLUMN IF NOT EXISTS
  tax_declaration_frequency VARCHAR(20) DEFAULT 'MONTHLY';

-- Add tax fields to account_balances
ALTER TABLE account_balances ADD COLUMN IF NOT EXISTS
  tax_code VARCHAR(20);  -- Link to tax type for aggregation

-- Add tax reference to journal_entries
ALTER TABLE journal_entries ADD COLUMN IF NOT EXISTS
  vat_rate DECIMAL(5,2);  -- VAT rate applicable to this entry
```
