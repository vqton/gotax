# Cash Module — Technical Specifications

**Version:** 1.0
**Date:** July 2026
**Circular:** 99/2025/TT-BTC, Law on Accounting 2015

---

## 1. Domain Model

### CashReceipt

```
CashReceipt {
  ID              UUID
  CompanyID       UUID
  VoucherNo       string        // Auto: R-YYYY-XXXX
  VoucherDate     date
  PostedDate      datetime
  CashAccountID   UUID          // FK -> Account (1111/1112)
  CounterpartID   UUID          // FK -> Customer/Supplier/Employee (nullable)
  CounterpartName string
  CounterpartType enum(customer,supplier,employee,other)
  Currency        string        // VND, USD, EUR
  ExchangeRate    decimal(18,4)
  Amount          decimal(18,2) // Original currency
  AmountVND       decimal(18,2) // VND equivalent
  DebitAccountID  UUID          // FK -> Account (111x)
  CreditAccountID UUID          // FK -> Account (counterpart)
  Reason          text          // Ly do thu
  ReceiptType     enum(customer_payment,loan_recovery,bank_withdrawal,advance_refund,sales,other)
  Status          enum(draft,submitted,approved,rejected,posted)
  Attachments     []File
  ApprovalBy      UUID          // User who approved
  ApprovedAt      datetime
  PostedBy        UUID          // User who posted
  PostedAt        datetime
  CreatedBy       UUID
  CreatedAt       datetime
  UpdatedAt       datetime
  GLJournalID     UUID          // FK -> JournalEntry (nullable, set on posting)
}
```

### CashPayment

```
CashPayment {
  ID              UUID
  CompanyID       UUID
  VoucherNo       string        // Auto: P-YYYY-XXXX
  VoucherDate     date
  PostedDate      datetime
  CashAccountID   UUID          // FK -> Account (1111/1112)
  PayeeID         UUID          // FK -> Customer/Supplier/Employee (nullable)
  PayeeName       string
  PayeeType       enum(supplier,employee,other,government)
  Currency        string
  ExchangeRate    decimal(18,4)
  Amount          decimal(18,2) // Original currency
  AmountVND       decimal(18,2) // VND equivalent
  DebitAccountID  UUID          // FK -> Account (counterpart)
  CreditAccountID UUID          // FK -> Account (111x)
  Reason          text
  PaymentType     enum(supplier_payment,salary,expense,bank_deposit,advance,tax,other)
  Status          enum(draft,submitted,approved,rejected,posted)
  Attachments     []File
  ApprovalBy      UUID
  ApprovedAt      datetime
  PostedBy        UUID
  PostedAt        datetime
  CreatedBy       UUID
  CreatedAt       datetime
  UpdatedAt       datetime
  GLJournalID     UUID          // FK -> JournalEntry
  WHTAmount       decimal(18,2) // Withholding tax (if applicable)
}
```

### CashTransfer

```
CashTransfer {
  ID              UUID
  CompanyID       UUID
  TransferDate    date
  FromAccountID   UUID          // FK -> Account (1111, 1112, 1121)
  ToAccountID     UUID          // FK -> Account (1111, 1112, 1121)
  Amount          decimal(18,2)
  Currency        string
  ExchangeRate    decimal(18,4)
  Reason          text
  TransferType    enum(bank_withdrawal,bank_deposit,currency_conversion)
  Status          enum(draft,posted)
  SourceVoucherID UUID          // FK -> CashPayment (for bank deposit)
  DestVoucherID   UUID          // FK -> CashReceipt (for bank withdrawal)
  CreatedAt       datetime
  PostedAt        datetime
}
```

### CashBook

```
CashBook {
  ID              UUID
  CompanyID       UUID
  Currency        string        // VND, USD, EUR
  AccountID       UUID          // FK -> Account (1111/1112)
  EntryDate       date
  OpeningBalance  decimal(18,2)
  TotalReceipts   decimal(18,2)
  TotalPayments   decimal(18,2)
  ClosingBalance  decimal(18,2)
  Entries         []CashBookEntry
}

CashBookEntry {
  LineNo          int
  VoucherDate     date
  VoucherType     enum(receipt,payment)
  VoucherNo       string
  Description     text
  ReceiptAmount   decimal(18,2)
  PaymentAmount   decimal(18,2)
  RunningBalance  decimal(18,2)
  RefID           UUID          // FK -> CashReceipt or CashPayment
}
```

### PettyCash

```
PettyCashFund {
  ID              UUID
  CompanyID       UUID
  FundCode         string
  FundName         string
  CustodianID      UUID          // FK -> Employee
  InitialAmount    decimal(18,2)
  CurrentBalance   decimal(18,2)
  Currency         string
  Status           enum(active,frozen,closed)
  CreatedAt        datetime
}

PettyCashTransaction {
  ID              UUID
  FundID           UUID          // FK -> PettyCashFund
  TransactionDate date
  Description     text
  Amount          decimal(18,2) // Positive = topup, Negative = settlement
  Status          enum(pending,approved,settled)
  ReceiptVoucherID UUID         // FK -> CashReceipt (for topup)
  PaymentVoucherID UUID         // FK -> CashPayment (for settlement)
}
```

### CashInventory

```
CashInventory {
  ID              UUID
  CompanyID       UUID
  InventoryDate   date
  CashAccountID   UUID
  Currency        string
  BookBalance     decimal(18,2)
  ActualBalance   decimal(18,2)
  Difference      decimal(18,2)
  DifferenceType  enum(excess,shortage,none)
  Denomination    []DenominationDetail
  Reason          text
  Status          enum(draft,completed)
  ApprovedBy      UUID
  CreatedAt       datetime
}

DenominationDetail {
  Denomination    decimal(18,2) // e.g., 500000, 200000, 100000
  Count           int
  Subtotal        decimal(18,2)
}
```

---

## 2. Repository Interface

```go
type CashRepository interface {
  // Receipts
  CreateReceipt(ctx context.Context, r *domain.CashReceipt) error
  GetReceipt(ctx context.Context, id uuid.UUID) (*domain.CashReceipt, error)
  ListReceipts(ctx context.Context, filter domain.CashReceiptFilter) ([]domain.CashReceipt, int, error)
  UpdateReceipt(ctx context.Context, r *domain.CashReceipt) error
  DeleteReceipt(ctx context.Context, id uuid.UUID) error

  // Payments
  CreatePayment(ctx context.Context, p *domain.CashPayment) error
  GetPayment(ctx context.Context, id uuid.UUID) (*domain.CashPayment, error)
  ListPayments(ctx context.Context, filter domain.CashPaymentFilter) ([]domain.CashPayment, int, error)
  UpdatePayment(ctx context.Context, p *domain.CashPayment) error
  DeletePayment(ctx context.Context, id uuid.UUID) error

  // Cash Book
  GetCashBook(ctx context.Context, companyID uuid.UUID, currency string, from, to time.Time) (*domain.CashBook, error)
  GetBalance(ctx context.Context, companyID uuid.UUID, accountID uuid.UUID) (decimal.Decimal, error)

  // Transfers
  CreateTransfer(ctx context.Context, t *domain.CashTransfer) error
  ListTransfers(ctx context.Context, companyID uuid.UUID) ([]domain.CashTransfer, error)

  // Petty Cash
  CreatePettyCashFund(ctx context.Context, f *domain.PettyCashFund) error
  GetPettyCashFund(ctx context.Context, id uuid.UUID) (*domain.PettyCashFund, error)
  ListPettyCashFunds(ctx context.Context, companyID uuid.UUID) ([]domain.PettyCashFund, error)

  // Inventory
  CreateInventory(ctx context.Context, inv *domain.CashInventory) error
  ListInventories(ctx context.Context, companyID uuid.UUID) ([]domain.CashInventory, error)
}
```

---

## 3. API Endpoints

### Voucher Management

```
POST   /api/v1/cash/receipts           Create receipt
GET    /api/v1/cash/receipts           List receipts (filter: date, type, status, currency)
GET    /api/v1/cash/receipts/:id       Get receipt detail
PUT    /api/v1/cash/receipts/:id       Update receipt (draft only)
DELETE /api/v1/cash/receipts/:id       Delete receipt (draft only)
POST   /api/v1/cash/receipts/:id/submit    Submit for approval
POST   /api/v1/cash/receipts/:id/approve   Approve
POST   /api/v1/cash/receipts/:id/reject    Reject
POST   /api/v1/cash/receipts/:id/post      Post to GL

POST   /api/v1/cash/payments           Create payment
GET    /api/v1/cash/payments           List payments (filter: date, type, status, currency)
GET    /api/v1/cash/payments/:id       Get payment detail
PUT    /api/v1/cash/payments/:id       Update payment (draft only)
DELETE /api/v1/cash/payments/:id       Delete payment (draft only)
POST   /api/v1/cash/payments/:id/submit    Submit for approval
POST   /api/v1/cash/payments/:id/approve   Approve
POST   /api/v1/cash/payments/:id/reject    Reject
POST   /api/v1/cash/payments/:id/post      Post to GL
```

### Transfers & Balance

```
POST   /api/v1/cash/transfers          Create bank transfer
GET    /api/v1/cash/transfers          List transfers
GET    /api/v1/cash/cash-book          Get cash book (query: currency, from, to)
GET    /api/v1/cash/balance            Get current balance (query: currency)
```

### Petty Cash

```
POST   /api/v1/cash/petty-cash         Create petty cash fund
GET    /api/v1/cash/petty-cash         List funds
POST   /api/v1/cash/petty-cash/:id/topup   Top up fund
POST   /api/v1/cash/petty-cash/:id/settle  Settle fund
```

### Inventory

```
POST   /api/v1/cash/inventory          Create inventory
GET    /api/v1/cash/inventory          List inventories
GET    /api/v1/cash/inventory/:id      Get inventory detail
```

### Reports

```
GET    /api/v1/reports/cash-flow       Cash Flow Statement (B03-DN)
GET    /api/v1/reports/cash-book       Cash Book (S07-DN / S04a-DNN)
GET    /api/v1/reports/cash-ledger     Cash Detail Ledger (S07a-DN / S04b-DNN)
```

---

## 4. Database Schema (PostgreSQL)

```sql
CREATE TABLE cash_receipts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  company_id UUID NOT NULL REFERENCES companies(id),
  voucher_no VARCHAR(20) NOT NULL,
  voucher_date DATE NOT NULL,
  posted_date TIMESTAMPTZ,
  cash_account_id UUID NOT NULL REFERENCES accounts(id),
  counterpart_id UUID,
  counterpart_name VARCHAR(255),
  counterpart_type VARCHAR(20) NOT NULL DEFAULT 'other',
  currency VARCHAR(3) NOT NULL DEFAULT 'VND',
  exchange_rate NUMERIC(18,4) NOT NULL DEFAULT 1,
  amount NUMERIC(18,2) NOT NULL,
  amount_vnd NUMERIC(18,2) NOT NULL,
  debit_account_id UUID NOT NULL REFERENCES accounts(id),
  credit_account_id UUID NOT NULL REFERENCES accounts(id),
  reason TEXT,
  receipt_type VARCHAR(30) NOT NULL DEFAULT 'other',
  status VARCHAR(20) NOT NULL DEFAULT 'draft',
  approved_by UUID REFERENCES users(id),
  approved_at TIMESTAMPTZ,
  posted_by UUID REFERENCES users(id),
  posted_at TIMESTAMPTZ,
  gl_journal_id UUID,
  created_by UUID NOT NULL REFERENCES users(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(company_id, voucher_no)
);

CREATE TABLE cash_payments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  company_id UUID NOT NULL REFERENCES companies(id),
  voucher_no VARCHAR(20) NOT NULL,
  voucher_date DATE NOT NULL,
  posted_date TIMESTAMPTZ,
  cash_account_id UUID NOT NULL REFERENCES accounts(id),
  payee_id UUID,
  payee_name VARCHAR(255),
  payee_type VARCHAR(20) NOT NULL DEFAULT 'other',
  currency VARCHAR(3) NOT NULL DEFAULT 'VND',
  exchange_rate NUMERIC(18,4) NOT NULL DEFAULT 1,
  amount NUMERIC(18,2) NOT NULL,
  amount_vnd NUMERIC(18,2) NOT NULL,
  debit_account_id UUID NOT NULL REFERENCES accounts(id),
  credit_account_id UUID NOT NULL REFERENCES accounts(id),
  reason TEXT,
  payment_type VARCHAR(30) NOT NULL DEFAULT 'other',
  status VARCHAR(20) NOT NULL DEFAULT 'draft',
  approved_by UUID REFERENCES users(id),
  approved_at TIMESTAMPTZ,
  posted_by UUID REFERENCES users(id),
  posted_at TIMESTAMPTZ,
  gl_journal_id UUID,
  wht_amount NUMERIC(18,2) DEFAULT 0,
  created_by UUID NOT NULL REFERENCES users(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(company_id, voucher_no)
);

CREATE TABLE cash_transfers (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  company_id UUID NOT NULL REFERENCES companies(id),
  transfer_date DATE NOT NULL,
  from_account_id UUID NOT NULL REFERENCES accounts(id),
  to_account_id UUID NOT NULL REFERENCES accounts(id),
  amount NUMERIC(18,2) NOT NULL,
  currency VARCHAR(3) NOT NULL DEFAULT 'VND',
  exchange_rate NUMERIC(18,4) NOT NULL DEFAULT 1,
  reason TEXT,
  transfer_type VARCHAR(30) NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'draft',
  source_voucher_id UUID,
  dest_voucher_id UUID,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  posted_at TIMESTAMPTZ
);

CREATE TABLE petty_cash_funds (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  company_id UUID NOT NULL REFERENCES companies(id),
  fund_code VARCHAR(20) NOT NULL,
  fund_name VARCHAR(255) NOT NULL,
  custodian_id UUID NOT NULL REFERENCES employees(id),
  initial_amount NUMERIC(18,2) NOT NULL,
  current_balance NUMERIC(18,2) NOT NULL,
  currency VARCHAR(3) NOT NULL DEFAULT 'VND',
  status VARCHAR(20) NOT NULL DEFAULT 'active',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(company_id, fund_code)
);

CREATE TABLE cash_inventories (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  company_id UUID NOT NULL REFERENCES companies(id),
  inventory_date DATE NOT NULL,
  cash_account_id UUID NOT NULL REFERENCES accounts(id),
  currency VARCHAR(3) NOT NULL DEFAULT 'VND',
  book_balance NUMERIC(18,2) NOT NULL,
  actual_balance NUMERIC(18,2) NOT NULL,
  difference NUMERIC(18,2) NOT NULL DEFAULT 0,
  difference_type VARCHAR(10) NOT NULL DEFAULT 'none',
  reason TEXT,
  status VARCHAR(20) NOT NULL DEFAULT 'draft',
  approved_by UUID REFERENCES users(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE cash_inventory_details (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  inventory_id UUID NOT NULL REFERENCES cash_inventories(id),
  denomination NUMERIC(18,2) NOT NULL,
  count INT NOT NULL DEFAULT 0,
  subtotal NUMERIC(18,2) NOT NULL
);
```

---

## 5. Validation Rules

| Field | Rule |
|-------|------|
| Amount | Must be > 0 |
| Voucher date | Cannot be in future (unless configured) |
| Currency | Must be VND, USD, EUR, or configured currency |
| Cash account | Must be 1111 (VND) or 1112 (foreign currency) |
| Debit = Credit | Double-entry must balance |
| Balance check | Optional: cannot exceed cash balance |
| Status transition | draft → submitted → approved → posted (no skip) |
| Approval | Requires different user than creator |
| Foreign currency | Exchange rate required, must be > 0 |

---

## 6. GL Posting Integration

On posting, system creates JournalEntry:

**Receipt posting:**
```
Nợ 111x (CashAccountID)   — AmountVND
Có xxxx (CreditAccountID) — AmountVND
```

**Payment posting:**
```
Nợ xxxx (DebitAccountID)  — AmountVND
Có 111x (CashAccountID)   — AmountVND
```

GL journal reference stored in `gl_journal_id` field. Reversal creates offsetting entry.
