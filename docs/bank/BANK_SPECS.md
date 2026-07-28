# Bank Module — Technical Specifications

**Version:** 1.0
**Date:** July 2026
**Circular:** 99/2025/TT-BTC, Decree 123/2020/ND-CP

---

## 1. Domain Models

### BankStatement

```
BankStatement {
  ID            string          // UUID
  CompanyID     string          // FK -> Company
  BankAccountID string          // FK -> CompanyBankAccount
  StatementDate date            // Ngày sao kê
  FromDate      date            // Từ ngày
  ToDate        date            // Đến ngày
  OpeningBalance float64        // Số dư đầu kỳ
  ClosingBalance float64        // Số dư cuối kỳ
  TotalDebits    float64        // Tổng phát sinh nợ
  TotalCredits   float64        // Tổng phát sinh có
  Currency       string         // VND / USD / EUR
  ImportMethod   string         // CSV, MT940, OFX, MANUAL
  RawFileName    string         // Original import file name
  RawFileHash    string         // SHA-256 of raw file
  LineCount      int            // Number of transaction lines
  Status         BankStatementStatus
  ImportedBy     string         // User who imported
  ImportedAt     time.Time
  Notes          string
}
```

### BankStatementLine

```
BankStatementLine {
  ID              string    // UUID
  StatementID     string    // FK -> BankStatement
  TransactionDate date      // Ngày giao dịch
  ValueDate       date      // Ngày giá trị
  Description     string    // Nội dung giao dịch
  DebitAmount     float64   // Số tiền nợ (chi ra)
  CreditAmount    float64   // Số tiền có (thu vào)
  BalanceAfter    float64   // Số dư sau giao dịch
  ReferenceNo     string    // Số tham chiếu
  BankRef         string    // Mã giao dịch ngân hàng
  Counterparty     string   // Đối tác giao dịch
  CounterpartyAcc  string   // Số TK đối tác
  BankCode         string   // Mã NH đối tác
  RawData          string   // Original import line (audit)
  MatchStatus     MatchStatus
  MatchedLineID   string    // FK -> GL journal line / AR/AP
  MatchedAt       time.Time
  MatchedBy       string
}
```

### BankReconciliation

```
BankReconciliation {
  ID              string    // UUID
  CompanyID       string    // FK -> Company
  BankAccountID   string    // FK -> CompanyBankAccount
  StatementID     string    // FK -> BankStatement
  FromDate        date
  ToDate          date
  OpeningBalance  float64   // Số dư đầu theo sổ kế toán
  ClosingBalance  float64   // Số dư cuối theo sổ kế toán
  StatementBalance float64  // Số dư cuối theo sao kê NH
  Difference      float64   // Chênh lệch
  Status          ReconStatus
  MatchedLines    int       // Number of matched lines
  UnmatchedLines  int       // Number of unmatched lines
  WriteOffAmount  float64   // Total write-off
  CompletedBy     string
  CompletedAt     time.Time
  ReversedAt      time.Time // If reversed
  Notes           string
}
```

### BankReconciliationMatch

```
BankReconciliationMatch {
  ID                string  // UUID
  ReconciliationID  string  // FK -> BankReconciliation
  StatementLineID   string  // FK -> BankStatementLine
  TransactionType   string  // GL_ENTRY, AR_RECEIPT, AP_PAYMENT, CASH_RECEIPT, CASH_PAYMENT, TAX_PAYMENT
  TransactionID     string  // FK to the matched transaction
  TransactionRef    string  // Voucher/reference number
  MatchMethod       string  // AUTO, MANUAL, RULE
  Confidence        float64 // 0.0-1.0 for auto matches
  CreatedAt         time.Time
}

MatchStatus enum(PENDING, MATCHED, UNMATCHED, WRITTEN_OFF)
ReconStatus enum(IN_PROGRESS, COMPLETED, REVERSED)
BankStatementStatus enum(IMPORTED, PROCESSING, RECONCILED, ARCHIVED)
```

### PaymentOrder

```
PaymentOrder {
  ID              string      // UUID
  CompanyID       string      // FK -> Company
  PaymentDate     date        // Ngày thanh toán
  Amount          float64     // Số tiền
  Currency        string
  ExchangeRate    float64
  BeneficiaryName string      // Người thụ hưởng
  BeneficiaryAcc  string      // Số tài khoản thụ hưởng
  BeneficiaryBank string      // Ngân hàng thụ hưởng
  BeneficiaryBranch string    // Chi nhánh NH thụ hưởng
  BeneficiaryCode string      // Mã số thuế / CMND (optional)
  FromBankAccID   string      // FK -> CompanyBankAccount
  PaymentContent  string      // Nội dung thanh toán
  Urgent          bool        // Urgent (URGENT / NORMAL)
  PaymentType     PaymentOrderType  // SUPPLIER, SALARY, TAX, LOAN, INTERNAL, OTHER
  ReferenceOrders []string    // Reference invoice/purchase order IDs
  Status          PaymentOrderStatus
  CreatedBy       string
  ApprovedBy      string
  ApprovedAt      time.Time
  SubmittedAt     time.Time
  BankRef         string      // Bank transaction reference
  FailureReason   string
  ErrorCode       string
  PrintCount      int
  CreatedAt       time.Time
  UpdatedAt       time.Time
}

PaymentOrderType enum(SUPPLIER, SALARY, TAX, LOAN, INTERNAL, OTHER)
PaymentOrderStatus enum(DRAFT, PENDING_APPROVAL, APPROVED, REJECTED, SUBMITTED, CONFIRMED, FAILED, CANCELLED)
```

### PaymentOrderBatch

```
PaymentOrderBatch {
  ID            string
  CompanyID     string
  BatchName     string       // e.g. "Supplier Payment July 2026"
  BatchDate     date
  TotalAmount   float64
  Currency      string
  OrderCount    int
  Status        BatchStatus  // DRAFT, SUBMITTED, CONFIRMED, PARTIAL
  CreatedBy     string
  SubmittedAt   time.Time
  BankRef       string
  CreatedAt     time.Time
}

BatchStatus enum(DRAFT, SUBMITTED, CONFIRMED, PARTIAL, FAILED)
```

### LoanAgreement

```
LoanAgreement {
  ID              string           // UUID
  CompanyID       string           // FK -> Company
  BankAccountID   string           // FK -> CompanyBankAccount (lending bank)
  ContractNo      string           // Hợp đồng vay số
  LoanType        LoanType         // SHORT_TERM, LONG_TERM, OVERDRAFT
  PrincipalAmount float64          // Số tiền vay
  Currency        string
  InterestRate    float64          // Annual %
  InterestMethod  string           // FIXED, FLOATING (base + margin)
  BaseRate        float64          // Base rate for floating loans
  MarginRate      float64          // Margin over base
  DisbursedAmount float64          // Actual disbursed
  OutstandingBalance float64       // Remaining principal
  StartDate       date
  MaturityDate    date
  RepaymentMethod string           // ANNUITY, STRAIGHT_LINE, BULLET, CUSTOM
  RepaymentFreq   string           // MONTHLY, QUARTERLY, ANNUAL, MATURITY
  Status          LoanStatus
  Notes           string
  CreatedAt       time.Time
  UpdatedAt       time.Time
}

LoanType enum(SHORT_TERM, LONG_TERM, OVERDRAFT)
LoanStatus enum(ACTIVE, FULLY_PAID, OVERDUE, RESTRUCTURED, WRITTEN_OFF)
```

### LoanDisbursement

```
LoanDisbursement {
  ID              string
  LoanID          string    // FK -> LoanAgreement
  DisbursementDate date
  Amount           float64
  ToBankAccountID  string   // FK -> CompanyBankAccount
  ReferenceNo      string
  Notes            string
  CreatedAt        time.Time
}
```

### LoanRepayment

```
LoanRepayment {
  ID              string
  LoanID          string
  RepaymentDate   date
  PrincipalAmount float64
  InterestAmount  float64
  FeeAmount       float64
  TotalAmount     float64
  PaymentOrderID  string    // FK -> PaymentOrder
  Status          RepaymentStatus
  Notes           string
  CreatedAt       time.Time
}

RepaymentStatus enum(SCHEDULED, PAID, PARTIAL, OVERDUE, WAIVED)
```

### TermDeposit

```
TermDeposit {
  ID              string
  CompanyID       string
  BankAccountID   string           // FK -> CompanyBankAccount
  DepositNo       string           // Số chứng chỉ / sổ tiết kiệm
  Amount          float64
  Currency        string
  InterestRate    float64          // Annual %
  Term            int              // Days
  StartDate       date
  MaturityDate    date
  InterestAtMaturity float64       // Calculated
  AutoRenewal     bool
  RenewalTerm     int              // Days (if auto-renewal)
  Status          DepositStatus
  Notes           string
  CreatedAt       time.Time
  MaturedAt       time.Time
}

DepositStatus enum(ACTIVE, MATURED, RENEWED, CLOSED)
```

---

## 2. Repository Interface

```go
type BankRepository interface {
  // Statements
  CreateStatement(ctx context.Context, s *BankStatement) error
  GetStatement(ctx context.Context, id string) (*BankStatement, error)
  ListStatements(ctx context.Context, companyID, bankAccountID string, limit, offset int) ([]BankStatement, int, error)
  
  CreateStatementLines(ctx context.Context, lines []BankStatementLine) error
  GetStatementLines(ctx context.Context, statementID string) ([]BankStatementLine, error)
  UpdateStatementLineMatch(ctx context.Context, lineID string, matchStatus MatchStatus, matchedLineID string) error
  
  // Reconciliation
  CreateReconciliation(ctx context.Context, r *BankReconciliation) error
  GetReconciliation(ctx context.Context, id string) (*BankReconciliation, error)
  ListReconciliations(ctx context.Context, companyID, bankAccountID string) ([]BankReconciliation, error)
  UpdateReconciliation(ctx context.Context, r *BankReconciliation) error
  
  CreateReconciliationMatch(ctx context.Context, m *BankReconciliationMatch) error
  DeleteReconciliationMatch(ctx context.Context, id string) error
  GetReconciliationMatches(ctx context.Context, reconID string) ([]BankReconciliationMatch, error)
  
  // Payment Orders
  CreatePaymentOrder(ctx context.Context, po *PaymentOrder) error
  GetPaymentOrder(ctx context.Context, id string) (*PaymentOrder, error)
  ListPaymentOrders(ctx context.Context, filter PaymentOrderFilter) ([]PaymentOrder, int, error)
  UpdatePaymentOrder(ctx context.Context, po *PaymentOrder) error
  
  CreatePaymentOrderBatch(ctx context.Context, b *PaymentOrderBatch) error
  GetPaymentOrderBatch(ctx context.Context, id string) (*PaymentOrderBatch, error)
  ListPaymentOrderBatches(ctx context.Context, companyID string) ([]PaymentOrderBatch, error)
  UpdatePaymentOrderBatch(ctx context.Context, b *PaymentOrderBatch) error
  
  // Loans
  CreateLoan(ctx context.Context, l *LoanAgreement) error
  GetLoan(ctx context.Context, id string) (*LoanAgreement, error)
  ListLoans(ctx context.Context, companyID string) ([]LoanAgreement, error)
  UpdateLoan(ctx context.Context, l *LoanAgreement) error
  
  CreateDisbursement(ctx context.Context, d *LoanDisbursement) error
  GetDisbursements(ctx context.Context, loanID string) ([]LoanDisbursement, error)
  
  CreateRepayment(ctx context.Context, r *LoanRepayment) error
  GetRepaymentSchedule(ctx context.Context, loanID string) ([]LoanRepayment, error)
  UpdateRepayment(ctx context.Context, r *LoanRepayment) error
  
  // Term Deposits
  CreateDeposit(ctx context.Context, d *TermDeposit) error
  GetDeposit(ctx context.Context, id string) (*TermDeposit, error)
  ListDeposits(ctx context.Context, companyID string) ([]TermDeposit, error)
  UpdateDeposit(ctx context.Context, d *TermDeposit) error
  
  // Reports
  GetBankLedger(ctx context.Context, companyID, bankAccountID, fromDate, toDate string) (*BankLedger, error)
  GetReconciliationReport(ctx context.Context, reconID string) (*ReconciliationReport, error)
}
```

---

## 3. API Endpoints

### Bank Statements
```
POST   /api/v1/bank/statements/import            Import statement file (CSV/MT940)
POST   /api/v1/bank/statements/manual            Create statement manually
GET    /api/v1/bank/statements                    List statements
GET    /api/v1/bank/statements/:id                Get statement detail
GET    /api/v1/bank/statements/:id/lines          Get statement lines
DELETE /api/v1/bank/statements/:id                Delete statement (unmatched only)
```

### Bank Reconciliation
```
POST   /api/v1/bank/reconciliations              Start reconciliation
GET    /api/v1/bank/reconciliations               List reconciliations
GET    /api/v1/bank/reconciliations/:id           Get reconciliation
POST   /api/v1/bank/reconciliations/:id/match     Match statement line to transaction
POST   /api/v1/bank/reconciliations/:id/auto-match Run auto-match rules
POST   /api/v1/bank/reconciliations/:id/write-off Create write-off entry
POST   /api/v1/bank/reconciliations/:id/complete  Complete reconciliation
POST   /api/v1/bank/reconciliations/:id/reverse   Reverse reconciliation
GET    /api/v1/bank/reconciliations/:id/report    S09-DN report
```

### Payment Orders
```
POST   /api/v1/bank/payment-orders               Create payment order
GET    /api/v1/bank/payment-orders                List payment orders (filter: date, status, type)
GET    /api/v1/bank/payment-orders/:id            Get payment order
PUT    /api/v1/bank/payment-orders/:id            Update (draft only)
DELETE /api/v1/bank/payment-orders/:id            Delete (draft only)
POST   /api/v1/bank/payment-orders/:id/submit     Submit for approval
POST   /api/v1/bank/payment-orders/:id/approve    Approve
POST   /api/v1/bank/payment-orders/:id/reject     Reject
POST   /api/v1/bank/payment-orders/:id/cancel     Cancel
POST   /api/v1/bank/payment-orders/:id/print      Print UNC form
POST   /api/v1/bank/payment-orders/:id/confirm    Mark as confirmed by bank

POST   /api/v1/bank/payment-batches               Create batch
GET    /api/v1/bank/payment-batches               List batches
GET    /api/v1/bank/payment-batches/:id           Get batch detail
POST   /api/v1/bank/payment-batches/:id/submit    Submit batch for payment
```

### Loans
```
POST   /api/v1/bank/loans                        Create loan agreement
GET    /api/v1/bank/loans                         List loans
GET    /api/v1/bank/loans/:id                     Get loan detail
PUT    /api/v1/bank/loans/:id                     Update loan
POST   /api/v1/bank/loans/:id/disburse            Record disbursement
GET    /api/v1/bank/loans/:id/schedule            Get repayment schedule
POST   /api/v1/bank/loans/:id/pay                 Record repayment
POST   /api/v1/bank/loans/:id/restructure         Restructure loan
```

### Term Deposits
```
POST   /api/v1/bank/term-deposits                Create term deposit
GET    /api/v1/bank/term-deposits                 List term deposits
GET    /api/v1/bank/term-deposits/:id             Get deposit detail
POST   /api/v1/bank/term-deposits/:id/mature      Process maturity
POST   /api/v1/bank/term-deposits/:id/renew       Auto-renew
POST   /api/v1/bank/term-deposits/:id/close       Close before maturity
```

### Reports
```
GET    /api/v1/reports/bank-ledger                S08-DN Bank Deposit Ledger
GET    /api/v1/reports/bank-reconciliation        S09-DN Bank Reconciliation
GET    /api/v1/reports/bank-balance               Current bank balance (all accounts)
GET    /api/v1/reports/bank-balance-history        Balance history chart data
GET    /api/v1/reports/bank-fees                  Bank fee analysis
```

---

## 4. Database Schema (PostgreSQL)

```sql
-- Bank Statements
CREATE TABLE bank_statements (
    id VARCHAR(50) PRIMARY KEY,
    company_id VARCHAR(50) NOT NULL REFERENCES companies(id),
    bank_account_id VARCHAR(50) NOT NULL REFERENCES company_bank_accounts(id),
    statement_date DATE NOT NULL,
    from_date DATE NOT NULL,
    to_date DATE NOT NULL,
    opening_balance NUMERIC(18,2) NOT NULL DEFAULT 0,
    closing_balance NUMERIC(18,2) NOT NULL DEFAULT 0,
    total_debits NUMERIC(18,2) NOT NULL DEFAULT 0,
    total_credits NUMERIC(18,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'VND',
    import_method VARCHAR(20) NOT NULL,
    raw_file_name VARCHAR(255),
    raw_file_hash VARCHAR(64),
    line_count INT NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'IMPORTED',
    imported_by VARCHAR(50) NOT NULL,
    imported_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    notes TEXT,
    UNIQUE(company_id, bank_account_id, from_date, to_date)
);

CREATE TABLE bank_statement_lines (
    id VARCHAR(50) PRIMARY KEY,
    statement_id VARCHAR(50) NOT NULL REFERENCES bank_statements(id),
    transaction_date DATE NOT NULL,
    value_date DATE,
    description TEXT,
    debit_amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    credit_amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    balance_after NUMERIC(18,2),
    reference_no VARCHAR(100),
    bank_ref VARCHAR(100),
    counterparty VARCHAR(255),
    counterparty_acc VARCHAR(50),
    bank_code VARCHAR(20),
    raw_data TEXT,
    match_status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    matched_line_id VARCHAR(50),
    matched_at TIMESTAMPTZ,
    matched_by VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Bank Reconciliation
CREATE TABLE bank_reconciliations (
    id VARCHAR(50) PRIMARY KEY,
    company_id VARCHAR(50) NOT NULL REFERENCES companies(id),
    bank_account_id VARCHAR(50) NOT NULL REFERENCES company_bank_accounts(id),
    statement_id VARCHAR(50) REFERENCES bank_statements(id),
    from_date DATE NOT NULL,
    to_date DATE NOT NULL,
    opening_balance NUMERIC(18,2) NOT NULL DEFAULT 0,
    closing_balance NUMERIC(18,2) NOT NULL DEFAULT 0,
    statement_balance NUMERIC(18,2) NOT NULL DEFAULT 0,
    difference NUMERIC(18,2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'IN_PROGRESS',
    matched_lines INT NOT NULL DEFAULT 0,
    unmatched_lines INT NOT NULL DEFAULT 0,
    write_off_amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    completed_by VARCHAR(50),
    completed_at TIMESTAMPTZ,
    reversed_at TIMESTAMPTZ,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, bank_account_id, from_date, to_date)
);

CREATE TABLE bank_reconciliation_matches (
    id VARCHAR(50) PRIMARY KEY,
    reconciliation_id VARCHAR(50) NOT NULL REFERENCES bank_reconciliations(id),
    statement_line_id VARCHAR(50) NOT NULL REFERENCES bank_statement_lines(id),
    transaction_type VARCHAR(30) NOT NULL,
    transaction_id VARCHAR(50) NOT NULL,
    transaction_ref VARCHAR(100),
    match_method VARCHAR(20) NOT NULL DEFAULT 'MANUAL',
    confidence NUMERIC(5,4) DEFAULT 1.0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(reconciliation_id, statement_line_id)
);

-- Payment Orders
CREATE TABLE payment_orders (
    id VARCHAR(50) PRIMARY KEY,
    company_id VARCHAR(50) NOT NULL REFERENCES companies(id),
    payment_date DATE NOT NULL,
    amount NUMERIC(18,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'VND',
    exchange_rate NUMERIC(18,4) NOT NULL DEFAULT 1,
    beneficiary_name VARCHAR(255) NOT NULL,
    beneficiary_acc VARCHAR(50) NOT NULL,
    beneficiary_bank VARCHAR(255) NOT NULL,
    beneficiary_branch VARCHAR(255),
    beneficiary_code VARCHAR(50),
    from_bank_acc_id VARCHAR(50) NOT NULL REFERENCES company_bank_accounts(id),
    payment_content TEXT,
    urgent BOOLEAN NOT NULL DEFAULT FALSE,
    payment_type VARCHAR(30) NOT NULL DEFAULT 'OTHER',
    reference_orders TEXT[], -- Array of invoice/purchase order IDs
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    created_by VARCHAR(50) NOT NULL,
    approved_by VARCHAR(50),
    approved_at TIMESTAMPTZ,
    submitted_at TIMESTAMPTZ,
    bank_ref VARCHAR(100),
    failure_reason TEXT,
    error_code VARCHAR(20),
    print_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE payment_order_batches (
    id VARCHAR(50) PRIMARY KEY,
    company_id VARCHAR(50) NOT NULL REFERENCES companies(id),
    batch_name VARCHAR(255) NOT NULL,
    batch_date DATE NOT NULL,
    total_amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'VND',
    order_count INT NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    created_by VARCHAR(50) NOT NULL,
    submitted_at TIMESTAMPTZ,
    bank_ref VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE batch_payment_orders (
    batch_id VARCHAR(50) NOT NULL REFERENCES payment_order_batches(id),
    payment_order_id VARCHAR(50) NOT NULL REFERENCES payment_orders(id),
    PRIMARY KEY (batch_id, payment_order_id)
);

-- Loans
CREATE TABLE loan_agreements (
    id VARCHAR(50) PRIMARY KEY,
    company_id VARCHAR(50) NOT NULL REFERENCES companies(id),
    bank_account_id VARCHAR(50) NOT NULL REFERENCES company_bank_accounts(id),
    contract_no VARCHAR(50) NOT NULL,
    loan_type VARCHAR(20) NOT NULL,
    principal_amount NUMERIC(18,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'VND',
    interest_rate NUMERIC(9,4) NOT NULL,
    interest_method VARCHAR(20) NOT NULL DEFAULT 'FIXED',
    base_rate NUMERIC(9,4),
    margin_rate NUMERIC(9,4),
    disbursed_amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    outstanding_balance NUMERIC(18,2) NOT NULL DEFAULT 0,
    start_date DATE NOT NULL,
    maturity_date DATE NOT NULL,
    repayment_method VARCHAR(20) NOT NULL,
    repayment_freq VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, contract_no)
);

CREATE TABLE loan_disbursements (
    id VARCHAR(50) PRIMARY KEY,
    loan_id VARCHAR(50) NOT NULL REFERENCES loan_agreements(id),
    disbursement_date DATE NOT NULL,
    amount NUMERIC(18,2) NOT NULL,
    to_bank_account_id VARCHAR(50) REFERENCES company_bank_accounts(id),
    reference_no VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE loan_repayments (
    id VARCHAR(50) PRIMARY KEY,
    loan_id VARCHAR(50) NOT NULL REFERENCES loan_agreements(id),
    repayment_date DATE NOT NULL,
    principal_amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    interest_amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    fee_amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    total_amount NUMERIC(18,2) NOT NULL,
    payment_order_id VARCHAR(50) REFERENCES payment_orders(id),
    status VARCHAR(20) NOT NULL DEFAULT 'SCHEDULED',
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Term Deposits
CREATE TABLE term_deposits (
    id VARCHAR(50) PRIMARY KEY,
    company_id VARCHAR(50) NOT NULL REFERENCES companies(id),
    bank_account_id VARCHAR(50) NOT NULL REFERENCES company_bank_accounts(id),
    deposit_no VARCHAR(50) NOT NULL,
    amount NUMERIC(18,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'VND',
    interest_rate NUMERIC(9,4) NOT NULL,
    term_days INT NOT NULL,
    start_date DATE NOT NULL,
    maturity_date DATE NOT NULL,
    interest_at_maturity NUMERIC(18,2) NOT NULL DEFAULT 0,
    auto_renewal BOOLEAN NOT NULL DEFAULT FALSE,
    renewal_term_days INT,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    matured_at TIMESTAMPTZ,
    UNIQUE(company_id, deposit_no)
);

-- Indexes
CREATE INDEX idx_bank_statements_company ON bank_statements(company_id);
CREATE INDEX idx_bank_statements_account ON bank_statements(bank_account_id);
CREATE INDEX idx_bank_statement_lines_stmt ON bank_statement_lines(statement_id);
CREATE INDEX idx_bank_statement_lines_match ON bank_statement_lines(match_status);
CREATE INDEX idx_bank_recon_company ON bank_reconciliations(company_id);
CREATE INDEX idx_bank_recon_account ON bank_reconciliations(bank_account_id);
CREATE INDEX idx_payment_orders_company ON payment_orders(company_id);
CREATE INDEX idx_payment_orders_status ON payment_orders(status);
CREATE INDEX idx_payment_orders_date ON payment_orders(payment_date);
CREATE INDEX idx_payment_batches_company ON payment_order_batches(company_id);
CREATE INDEX idx_loans_company ON loan_agreements(company_id);
CREATE INDEX idx_loans_status ON loan_agreements(status);
CREATE INDEX idx_loan_disbursements_loan ON loan_disbursements(loan_id);
CREATE INDEX idx_loan_repayments_loan ON loan_repayments(loan_id);
CREATE INDEX idx_term_deposits_company ON term_deposits(company_id);
CREATE INDEX idx_term_deposits_status ON term_deposits(status);
```

---

## 5. Statement Import Format (CSV)

### Vietnamese Bank CSV Import (Recommended Standard)

```csv
TransactionDate,ValueDate,Description,DebitAmount,CreditAmount,Balance,ReferenceNo,BankRef,Counterparty,CounterpartyAcc
2026-01-15,2026-01-15,CT KH THANH TOAN HD123,0,50000000,150000000,REF001,BIDV240115001,CONG TY ABC,123456789
2026-01-16,2026-01-16,UNC THANH TOAN NHA CC,X,30000000,120000000,REF002,VCB240116002,CONG TY XYZ,987654321
```

### MT940 Format Support

- Tag 61: Statement line (value date, entry date, debit/credit, amount, reference)
- Tag 86: Description / details
- Mapping: 61 → BankStatementLine, 86 → Description

---

## 6. Auto-Match Rules

| # | Rule | Description | Confidence |
|---|------|-------------|------------|
| R1 | Exact amount + date match | Statement line amount matches GL entry amount on same date | High (0.95) |
| R2 | Amount match + 1 day tolerance | Amount matches, date off by 1 day | High (0.90) |
| R3 | Reference match | Statement reference matches voucher no | High (0.95) |
| R4 | Fuzzy reference match | Partial reference match (e.g., "HD123" in description matches invoice "HD123") | Medium (0.80) |
| R5 | Counterparty amount match | Counterparty name fuzzy match + amount match | Medium (0.75) |
| R6 | Rounding tolerance | Amount differs by < 500 VND (bank fees, rounding) | Low (0.50) |
| R7 | Customer AR payment match | Customer name in description + amount matches AR receipt | Medium (0.80) |
| R8 | Supplier AP payment match | Supplier name in description + amount matches AP payment | Medium (0.80) |

---

## 7. GL Posting Integration

### Payment Order → GL
```
Upon payment order confirmation:
Nợ 331/642/... (Payable/Expense account)    — Amount
Có 112x (Bank account)                       — Amount
```

### Reconciliation Write-off → GL
```
Upon write-off creation:
Nợ 635 (Chi phí tài chính)                  — Amount (bank fee)
Có 112x (Bank account)                       — Amount
```

### Loan Disbursement → GL
```
Nợ 112x (Bank account)                       — Disbursed amount
Có 341 (Vay và nợ thuê tài chính)           — Disbursed amount
```

### Loan Repayment → GL
```
Nợ 341 (Vay và nợ thuê tài chính)          — Principal
Nợ 635 (Chi phí tài chính)                  — Interest
Có 112x (Bank account)                       — Total
```

### Term Deposit → GL
```
Opening:
Nợ 1281 (Tiền gửi có kỳ hạn)               — Amount
Có 112x (Bank account)                       — Amount

Maturity:
Nợ 112x (Bank account)                       — Principal + Interest
Có 1281 (Tiền gửi có kỳ hạn)               — Principal
Có 515 (Doanh thu tài chính)                — Interest
```

---

## 8. Revaluation (Circular 99 Art. 6)

Foreign currency bank accounts (1122) must be revalued at period-end:

- Rate: Average transfer buying/selling rate of bank where deposit account is opened
- Difference: 1122 → 515 (gain) or 635 (loss)

```
Nợ 1122 (Nếu tỷ giá tăng)                    — Exchange gain
Có 515                                        — Exchange gain

Nợ 635                                       — Exchange loss
Có 1122 (Nếu tỷ giá giảm)                    — Exchange loss
```

---

## 9. UNC Forms

Bank-specific UNC templates must be implemented for printable forms. Key Vietnamese banks:

| Bank | Template Code | Fields |
|------|--------------|--------|
| Vietcombank | UNC-VCB-01 | Số UNC, Ngày, Người thụ hưởng, Số TK, NH thụ hưởng, Số tiền (chữ+số), Nội dung |
| BIDV | UNC-BIDV-01 | Same structure, different layout |
| VietinBank | UNC-CTG-01 | Same structure, different layout |
| ACB | UNC-ACB-01 | Includes QR code field |
| VPBank | UNC-VPB-01 | Includes additional reference fields |
| MB Bank | UNC-MB-01 | Includes fee allocation field |
| Techcombank | UNC-TCB-01 | Digital signature field |
| HDBank | UNC-HDB-01 | Standard format |
| SHB | UNC-SHB-01 | Custom layout |
| VIB | UNC-VIB-01 | Includes purpose code |
