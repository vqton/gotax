# Tax Module — Functional Specification

**Role:** BA Lead (20+ yrs) + Chief Accountant (20+ yrs)
**Date:** 2026-07-27
**Version:** 1.0

---

## 1. Domain Model

### 1.1 TaxDeclaration

```
TaxDeclaration
├── ID              UUID (PK)
├── CompanyID       UUID (FK → Company)
├── DeclarationType Enum: GTGT01 | GTGT02 | GTGT03 | GTGT04 | GTGT05 |
│                          TNDN03 | TNDN04 | TNDN02 | TNDN05 | TNDN06 |
│                          KK_TNCN | QTT_TNCN | TTDB01 | BVMT01 |
│                          NTNN01 | NTNN02 | NTNN03
├── TaxPeriod       TaxPeriod (PeriodType: MONTHLY | QUARTERLY | ANNUAL | PER_OCCURRENCE)
├── PeriodYear      int
├── PeriodNumber    int (month: 1-12, quarter: 1-4)
├── Status          Enum: DRAFT | VALIDATED | SUBMITTED | ACKNOWLEDGED |
│                          REJECTED | AMENDED | CANCELLED
├── SubmittedAt     timestamp
├── SubmittedBy     UUID (FK → User)
├── AcknowledgedAt  timestamp
├── AcknowledgementRef string (GDT MA THONG BAO)
├── GDTResponseXML  text (raw GDT response)
├── DeclarationXML  text (submitted XML content)
├── PreviousDeclarationID UUID (FK → self, for amendments)
├── AdjustmentType  Enum: NONE | AMENDMENT | ADDITIONAL
├── Lines[]         TaxDeclarationLine
├── Signatures[]    DeclarationSignature
├── Version         int
├── CreatedAt       timestamp
├── CreatedBy       UUID (FK → User)
└── UpdatedAt       timestamp
```

### 1.2 TaxDeclarationLine

```
TaxDeclarationLine
├── ID              UUID (PK)
├── DeclarationID   UUID (FK → TaxDeclaration)
├── LineCode        string (e.g., "01", "10", "20" — form indicator code)
├── LineName        string (e.g., "Doanh thu ban hang hoa")
├── Amount          decimal(18,2)
├── SourceType      Enum: AUTO_CALCULATED | MANUAL_ENTRY | FROM_LEDGER
├── SourceAccount   string (account code if auto-calculated)
├── SourceEntryIDs  UUID[] (journal entry IDs that contribute to this line)
├── Note            text
└── SortOrder       int
```

### 1.3 TaxRate

```
TaxRate
├── ID              UUID (PK)
├── TaxType         Enum: VAT | CIT | PIT | TTDB | BVMT | FCT | RESOURCE | LAND
├── RateCode        string (e.g., "VAT_10", "CIT_20", "PIT_BRACKET_1")
├── RateName        string (e.g., "VAT 10% (standard)")
├── RateType        Enum: PERCENTAGE | FIXED | PROGRESSIVE
├── RateValue       decimal(5,2) (for percentage/flat rates)
├── ProgressiveBrackets[] ProgressiveBracket (for PIT)
├── EffectiveFrom   date
├── EffectiveTo     date (nullable, NULL = current)
├── IsActive        bool
├── CountryCode     string (default: "VN")
├── ApplicableTo    string (sector, product category, etc.)
├── LegalRef        string (reference law/circular)
└── CreatedAt       timestamp
```

### 1.4 TaxPayment

```
TaxPayment
├── ID              UUID (PK)
├── CompanyID       UUID (FK → Company)
├── DeclarationID   UUID (FK → TaxDeclaration, nullable)
├── TaxType         Enum
├── PeriodYear      int
├── PeriodNumber    int
├── DeclaredAmount  decimal(18,2)
├── PaidAmount      decimal(18,2)
├── PaymentDate     date
├── DueDate         date
├── PaymentRef      string (bank transaction reference)
├── PaymentMethod   Enum: EFT | BANK_TRANSFER | CASH | CQ
├── Status          Enum: PENDING | PAID | PARTIAL | OVERPAID | CANCELLED
├── LatePaymentDays int
├── LateInterest    decimal(18,2)
├── Notes           text
└── CreatedAt       timestamp
```

### 1.5 EInvoice

```
EInvoice
├── ID                UUID (PK)
├── CompanyID         UUID (FK → Company)
├── Pattern           string (e.g., "01GTKT0/001")
├── Serial            string (e.g., "AA/22E")
├── InvoiceNumber     int (assigned by GDT)
├── InvoiceType       Enum: ORIGINAL | ADJUSTMENT | REPLACEMENT | CANCELLATION_NOTE
├── GDTTransactionID  string (MA CUA CQT from GDT)
├── BuyerName         string
├── BuyerTaxCode      string (10 or 13 digits, nullable for individuals)
├── BuyerAddress      string
├── BuyerEmail        string
├── CurrencyCode      string (default: "VND")
├── ExchangeRate      decimal(18,4)
├── Lines[]           EInvoiceLine
├── Subtotal          decimal(18,2)
├── VATAmount         decimal(18,2)
├── GrandTotal        decimal(18,2)
├── XMLBody           text (TXML format)
├── SignedXML         text (signed TXML)
├── IssueDate         date
├── SigningDate       timestamp
├── DigitalSignatureID UUID (FK → DigitalSignature)
├── JournalEntryID    UUID (FK → JournalEntry)
├── Status            Enum: DRAFT | SIGNED | SUBMITTED | VALIDATED |
│                             ISSUED | CANCELLED | REPLACED
├── CancelledAt       timestamp
├── CancelReason      text
├── OriginalInvoiceID UUID (FK → self, for adjustment/cancellation)
├── GDTResponse       text
└── CreatedAt         timestamp
```

### 1.6 TaxCalendar

```
TaxCalendar
├── ID              UUID (PK)
├── CompanyID       UUID (FK → Company)
├── TaxType         Enum
├── PeriodType      Enum: MONTHLY | QUARTERLY | ANNUAL | PER_OCCURRENCE
├── PeriodYear      int
├── PeriodNumber    int
├── StartDate       date
├── EndDate         date
├── DeclarationDue  date
├── PaymentDue      date
├── Status          Enum: PENDING | SUBMITTED | PAID | MISSED | OVERDUE
├── DeclarationID   UUID (FK → TaxDeclaration, nullable)
└── CreatedAt       timestamp
```

## 2. State Machines

### 2.1 TaxDeclaration Status

```
                       ┌──────────┐
                       │  DRAFT   │
                       └────┬─────┘
                            │ Validate
                            ▼
                    ┌───────────────┐
                    │  VALIDATED    │
                    └───────┬───────┘
                            │ Submit to GDT
                            ▼
            ┌───────────────────────────┐
            │         SUBMITTED         │
            └───────────┬───────────────┘
                        │ GDT response
                        ▼
┌────────────────┬────────────────┬─────────────────┐
│                │                │                  │
▼                ▼                ▼                  ▼
┌─────────┐ ┌───────────┐ ┌───────────┐ ┌──────────────┐
│ACKNOWL- │ │ REJECTED  │ │  ERROR    │ │  CANCELLED   │
│EDGED    │ │           │ │           │ │              │
└────┬────┘ └─────┬─────┘ └─────┬─────┘ └──────────────┘
     │            │              │
     │            └──→ Edit ───→ DRAFT (resubmit)
     │
     ├──→ Amendment needed → NEW AMENDMENT (links to this)
     │
     └──→ GDT audit → FROZEN (no edits)
```

### 2.2 EInvoice Status

```
DRAFT → SIGNED → SUBMITTED → VALIDATED → ISSUED
                                          │
                                          ├──→ CANCELLED (cancel)
                                          │
                                          └──→ REPLACED (replacement issued)
```

## 3. XML Generation (HTKK/GDT Format)

### 3.1 File Structure

```
<BK:BoKe>
  <BK:ThongTinChung>
    <BK:MaSoThue>0123456789</BK:MaSoThue>
    <BK:TenNguoiNopThue>CONG TY ABC</BK:TenNguoiNopThue>
    <BK:LoaiToKhai>01/GTGT</BK:LoaiToKhai>
    <BK:KyTinhThue>2026Q1</BK:KyTinhThue>
    <BK:LanDau>1</BK:LanDau>
    <BK:NgayTao>2026-04-15</BK:NgayTao>
  </BK:ThongTinChung>
  <BK:DuLieu>
    <BK:ChiTieu>
      <BK:MaChiTieu>[10]</BK:MaChiTieu>
      <BK:GiaTri>500000000</BK:GiaTri>
    </BK:ChiTieu>
    ...
  </BK:DuLieu>
</BK:BoKe>
```

### 3.2 Supported XSD Schemas

| Form | XSD Name | Circular |
|------|----------|----------|
| 01/GTGT | `01_GTGT.xsd` | 80/2021 |
| 03/TNDN | `03_TNDN.xsd` | 80/2021 |
| 04/TNDN | `04_TNDN.xsd` | 80/2021 |
| 05/KK-TNCN | `05_KK_TNCN.xsd` | 80/2021 |
| 02/QTT-TNCN | `02_QTT_TNCN.xsd` | 80/2021 |
| 01/TTDB | `01_TTDB.xsd` | 80/2021 |
| 01/BVMT | `01_BVMT.xsd` | 80/2021 |
| Invoice TXML | `Invoice.xml` | 254/2026 |

### 3.3 XML Generation Service

```
GenerateDeclarationXML(declaration TaxDeclaration) (string, error)
1. Load form template (XSLT stylesheet or Go template)
2. Map TaxDeclarationLine → XML indicator codes per form spec
3. Apply calculation rules (sum, diff, percentage)
4. Validate against XSD schema
5. Return raw XML string
```

## 4. GDT API Integration

### 4.1 Endpoints

```
PROD: https://thuedientu.gdt.gov.vn
UAT:  https://thuedientunth.gdt.gov.vn

Authentication:
  - HTTPS client certificate (mutual TLS)
  - USB Token / remote HSM

Endpoints:
  POST /api/submission/declare       → Submit XML declaration
  GET  /api/submission/status/{id}   → Query submission status
  POST /api/einvoice/validate        → Validate invoice XML
  POST /api/einvoice/issue           → Issue e-invoice
  POST /api/einvoice/cancel          → Cancel e-invoice
  GET  /api/taxpayer/info/{taxcode}  → Query taxpayer info
  GET  /api/rates                    → Get current tax rates
  POST /api/payment/order            → Create payment order
```

### 4.2 Response Codes

| GDT Code | Meaning | Action |
|----------|---------|--------|
| `00` | Success / Ack received | Store acknowledgement |
| `01` | XML schema validation failed | Show error, allow resubmit |
| `02` | Duplicate declaration | Return existing ack |
| `03` | Tax code not found | Verify tax code |
| `10` | Period already declared | Error, can only amend |
| `99` | System error | Retry with backoff |
| Auth failure | Certificate invalid | Notify IT, block submission |

## 5. Integration Architecture

```
┌─────────┐     ┌──────────┐     ┌───────────┐     ┌───────────────┐
│  GoTax  │────→│  Tax     │────→│  XML      │────→│  GDT API      │
│  UI     │     │  Service │     │  Generator│     │  Client       │
└─────────┘     └──────────┘     └───────────┘     └───────┬───────┘
                      │                                     │
                      ▼                                     ▼
               ┌──────────┐                         ┌─────────────┐
               │  DB      │                         │ thuedientu  │
               │  (PG)    │                         │ .gdt.gov.vn │
               └──────────┘                         └─────────────┘
                                    ┌──────────────┐
                                    │  eTax Portal │
                                    │  (HTKK/iHTKK)│
                                    └──────────────┘
```

## 6. Repository Interfaces

### 6.1 TaxRepository (to add to domain/interfaces.go)

```go
type TaxRepository interface {
    // Tax Declarations
    CreateDeclaration(ctx context.Context, d *domain.TaxDeclaration) error
    GetDeclarationByID(ctx context.Context, id uuid.UUID) (*domain.TaxDeclaration, error)
    GetDeclarationsByCompany(ctx context.Context, companyID uuid.UUID, filter TaxDeclarationFilter) ([]*domain.TaxDeclaration, error)
    UpdateDeclarationStatus(ctx context.Context, id uuid.UUID, status domain.DeclarationStatus, gdtResponse string) error
    GetDeclarationByPeriod(ctx context.Context, companyID uuid.UUID, declType domain.DeclarationType, year, period int) (*domain.TaxDeclaration, error)

    // Tax Declaration Lines
    CreateDeclarationLine(ctx context.Context, line *domain.TaxDeclarationLine) error
    GetDeclarationLines(ctx context.Context, declarationID uuid.UUID) ([]*domain.TaxDeclarationLine, error)
    DeleteDeclarationLines(ctx context.Context, declarationID uuid.UUID) error

    // Tax Rates
    GetTaxRate(ctx context.Context, taxType string, effectiveDate time.Time) (*domain.TaxRate, error)
    GetPITBrackets(ctx context.Context, effectiveDate time.Time) ([]domain.ProgressiveBracket, error)
    ListTaxRates(ctx context.Context, filter TaxRateFilter) ([]*domain.TaxRate, error)
    CreateTaxRate(ctx context.Context, rate *domain.TaxRate) error

    // Tax Payments
    CreatePayment(ctx context.Context, p *domain.TaxPayment) error
    GetPaymentsByDeclaration(ctx context.Context, declarationID uuid.UUID) ([]*domain.TaxPayment, error)
    GetPaymentsByCompany(ctx context.Context, companyID uuid.UUID, filter PaymentFilter) ([]*domain.TaxPayment, error)
    UpdatePayment(ctx context.Context, p *domain.TaxPayment) error

    // E-Invoices
    CreateEInvoice(ctx context.Context, inv *domain.EInvoice) error
    GetEInvoiceByID(ctx context.Context, id uuid.UUID) (*domain.EInvoice, error)
    GetEInvoicesByCompany(ctx context.Context, companyID uuid.UUID, filter EInvoiceFilter) ([]*domain.EInvoice, error)
    UpdateEInvoiceStatus(ctx context.Context, id uuid.UUID, status domain.EInvoiceStatus) error

    // Tax Calendar
    GenerateTaxCalendar(ctx context.Context, companyID uuid.UUID, year int) error
    GetTaxCalendar(ctx context.Context, companyID uuid.UUID, year int) ([]*domain.TaxCalendar, error)
    GetUpcomingDeadlines(ctx context.Context, companyID uuid.UUID, days int) ([]*domain.TaxCalendar, error)

    // Tax Audit
    CreateAuditCase(ctx context.Context, audit *domain.TaxAuditCase) error
    GetAuditCasesByCompany(ctx context.Context, companyID uuid.UUID) ([]*domain.TaxAuditCase, error)
}
```

### 6.2 GDTClient (external service interface)

```go
type GDTClient interface {
    SubmitDeclaration(ctx context.Context, xmlData string, cert *domain.DigitalSignature) (*GDTResponse, error)
    QueryStatus(ctx context.Context, ackRef string) (*GDTStatus, error)
    ValidateInvoice(ctx context.Context, invoiceXML string) (*ValidationResult, error)
    IssueInvoice(ctx context.Context, invoiceXML string, cert *domain.DigitalSignature) (*InvoiceResponse, error)
    CancelInvoice(ctx context.Context, invoiceRef string, reason string, cert *domain.DigitalSignature) error
    GetTaxpayerInfo(ctx context.Context, taxCode string) (*TaxpayerInfo, error)
}
```

## 7. Service Layer

### 7.1 TaxService (new)

```go
type TaxService interface {
    // Declaration
    CreateDeclaration(ctx context.Context, companyID uuid.UUID, declType domain.DeclarationType, period domain.TaxPeriod) (*domain.TaxDeclaration, error)
    AutoPopulateDeclaration(ctx context.Context, declarationID uuid.UUID) error
    ValidateDeclaration(ctx context.Context, declarationID uuid.UUID) (*ValidationResult, error)
    SubmitDeclaration(ctx context.Context, declarationID uuid.UUID, signatureID uuid.UUID) error

    // Calculation
    CalculateVAT(ctx context.Context, companyID uuid.UUID, period domain.TaxPeriod) (*domain.VATResult, error)
    CalculateCITProvisional(ctx context.Context, companyID uuid.UUID, year, quarter int) (*domain.CITResult, error)
    CalculateCITFinalization(ctx context.Context, companyID uuid.UUID, year int) (*domain.CITResult, error)
    CalculatePIT(ctx context.Context, companyID uuid.UUID, period domain.TaxPeriod) (*domain.PITResult, error)

    // E-Invoice
    CreateInvoice(ctx context.Context, companyID uuid.UUID, input *domain.EInvoiceInput) (*domain.EInvoice, error)
    IssueInvoice(ctx context.Context, invoiceID uuid.UUID) error
    CancelInvoice(ctx context.Context, invoiceID uuid.UUID, reason string) error

    // Tax Calendar
    GenerateCalendar(ctx context.Context, companyID uuid.UUID, year int) error
    GetAlerts(ctx context.Context, companyID uuid.UUID) ([]*domain.TaxAlert, error)
}
```

## 8. Calculation Engine

### 8.1 VAT Calculation

```
For form 01/GTGT (deduction method):
  [10] Doanh thu ban hang = SUM(Cr turnover of 511, 515, 711 in period)
  [20] Tong gia tang hang hoa dich vu = Taxable revenue × Applicable VAT rate
  [21] Thue GTGT dau ra = SUM(Cr of 33311 entries in period)
  [22] Thue GTGT dau ra = SUM(Output VAT from e-invoices in period)
  [23] Tong thue GTGT dau ra = [21] + [22]
  [24] Thue GTGT dau vao = SUM(Dr of 1331 entries in period)
  [25] Thue GTGT dau vao = SUM(Dr of 1332 entries in period)
  [30] Tong thue GTGT dau vao = [24] + [25]
  [40] Thue GTGT phai nap = [23] - [30] (if positive)
  [41] Thue GTGT con duoc khau tru = [30] - [23] (if negative → refundable)
```

### 8.2 CIT Calculation

```
For form 03/TNDN (finalization):
  [A] Doanh thu tinh thue = Total revenue from P&L (511, 515, 711)
  [B] Chi phi hop le = Deductible expenses per CIT Law
  [C] Thu nhap khac = Other taxable income
  [D] Thu nhap mien thue = Tax-exempt income
  [E] Thu nhap tinh thue = [A] - [B] + [C] - [D]
  [F] Loi ket chuyen = Tax loss carry-forward (max 5 years)
  [G] Thu nhap tinh thue sau chuyen lo = [E] - [F]
  [H] Thue suat = Applicable CIT rate (15%/17%/20%)
  [I] Thue TNDN phai nap = [G] × [H]
  [J] Thue TNDN da tam nop = Quarterly provisional payments
  [K] Thue TNDN con phai nop = [I] - [J] (if positive)
  [L] Thue TNDN nop thua = [J] - [I] (if positive → refund)
```

### 8.3 PIT Calculation (Resident Employee)

```
Gross income per month
  - Social insurance (8% of gross, capped)
  - Health insurance (1.5% of gross, capped)
  - Unemployment insurance (1% of gross, capped)
  = Taxable income (Thu nhap tinh thue)
  - Personal deduction (11M VND/month)
  - Dependant deductions (4.4M VND/dependant/month)
  = Taxable income after deductions
  → Apply progressive rate table:
    Up to 5M: 5%
    5-10M: 10%
    10-18M: 15%
    18-32M: 20%
    32-52M: 25%
    52-80M: 30%
    Over 80M: 35%
  = PIT payable
```
