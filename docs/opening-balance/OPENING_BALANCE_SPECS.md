# Technical Specifications — Opening Balance Module

**Version:** 1.0 | **Date:** 2026-07-27
**Target:** GoTax Backend (Go 1.26.5, Gin v1.12, pgx/v5)

---

## 1. Data Model

### 1.1 Core Entity

```go
type OpeningBalanceStatus string
const (
    OBStatusDraft     OpeningBalanceStatus = "DRAFT"
    OBStatusPending   OpeningBalanceStatus = "PENDING"    // awaiting approval
    OBStatusApproved  OpeningBalanceStatus = "APPROVED"
    OBStatusRejected  OpeningBalanceStatus = "REJECTED"
    OBStatusCorrected OpeningBalanceStatus = "CORRECTED"  // superseded by correction
)

type OpeningBalance struct {
    ID               string               `json:"id"`
    CompanyID        string               `json:"company_id"`
    PeriodID         string               `json:"period_id"`          // V2 period
    FiscalYearID     string               `json:"fiscal_year_id"`
    AccountCode      string               `json:"account_code"`
    CurrencyCode     string               `json:"currency_code"`      // VND default
    OriginalAmount   float64              `json:"original_amount"`    // in currency
    DebitAmount      float64              `json:"debit_amount"`       // in VND
    CreditAmount     float64              `json:"credit_amount"`      // in VND
    ExchangeRate     float64              `json:"exchange_rate,omitempty"` // if foreign
    Status           OpeningBalanceStatus `json:"status"`
    SourceType       string               `json:"source_type"`        // MANUAL, CARRY_FORWARD, IMPORT, MIGRATION
    BatchID          string               `json:"batch_id,omitempty"` // import batch ref
    Reason           string               `json:"reason,omitempty"`
    Detail           OpeningBalanceDetail `json:"detail,omitempty"`   // per-entity breakdown
    ApprovedBy       string               `json:"approved_by,omitempty"`
    ApprovedAt       *time.Time           `json:"approved_at,omitempty"`
    CorrectedBy      string               `json:"corrected_by,omitempty"`
    CorrectionOf     string               `json:"correction_of,omitempty"` // parent OB ID
    CreatedBy        string               `json:"created_by"`
    CreatedAt        time.Time            `json:"created_at"`
    UpdatedAt        time.Time            `json:"updated_at"`
}
```

### 1.2 Detail Tracking Entity

```go
type DetailEntityType string
const (
    DetailCustomer  DetailEntityType = "CUSTOMER"
    DetailSupplier  DetailEntityType = "SUPPLIER"
    DetailEmployee  DetailEntityType = "EMPLOYEE"
    DetailBank      DetailEntityType = "BANK_ACCOUNT"
    DetailProject   DetailEntityType = "PROJECT"
    DetailContract  DetailEntityType = "CONTRACT"
    DetailDepartment DetailEntityType = "DEPARTMENT"
    DetailFixedAsset DetailEntityType = "FIXED_ASSET"
    DetailInventoryItem DetailEntityType = "INVENTORY_ITEM"
)

type OpeningBalanceDetail struct {
    ID              string           `json:"id"`
    OpeningBalanceID string          `json:"opening_balance_id"`
    EntityType      DetailEntityType `json:"entity_type"`
    EntityID        string           `json:"entity_id"`
    EntityName      string           `json:"entity_name,omitempty"`
    DebitAmount     float64          `json:"debit_amount"`
    CreditAmount    float64          `json:"credit_amount"`
    Quantity        float64          `json:"quantity,omitempty"`       // for inventory
    UnitPrice       float64          `json:"unit_price,omitempty"`     // for inventory
    OriginalCost    float64          `json:"original_cost,omitempty"`  // for FA
    AccDepreciation float64          `json:"acc_depreciation,omitempty"` // for FA
    CounterpartAccount string       `json:"counterpart_account,omitempty"`
    Note            string           `json:"note,omitempty"`
    Status          string           `json:"status,omitempty"`
    CreatedAt       time.Time        `json:"created_at"`
}
```

### 1.3 Carry-Forward Entity

```go
type CarryForwardLog struct {
    ID             string    `json:"id"`
    CompanyID      string    `json:"company_id"`
    FromPeriodID   string    `json:"from_period_id"`
    ToPeriodID     string    `json:"to_period_id"`
    FromFiscalYear int       `json:"from_fiscal_year"`
    ToFiscalYear   int       `json:"to_fiscal_year"`
    AccountCount   int       `json:"account_count"`
    TotalDebit     float64   `json:"total_debit"`
    TotalCredit    float64   `json:"total_credit"`
    ClosingEntries []string  `json:"closing_entry_ids,omitempty"` // journal entry IDs
    Status         string    `json:"status"`                     // COMPLETED, ROLLED_BACK
    ExecutedBy     string    `json:"executed_by"`
    ExecutedAt     time.Time `json:"executed_at"`
}
```

### 1.4 Migration Mapping Entity

```go
type Circular99Mapping struct {
    ID                string  `json:"id"`
    OldAccountCode    string  `json:"old_account_code"`     // Circular 200
    NewAccountCode    string  `json:"new_account_code"`     // Circular 99
    MappingType       string  `json:"mapping_type"`         // DIRECT, MERGE, SPLIT, REMAP
    SplitRatio        float64 `json:"split_ratio,omitempty"`
    CounterpartAccount string `json:"counterpart_account,omitempty"`
    EffectiveDate     string  `json:"effective_date"`
    Note              string  `json:"note,omitempty"`
    IsActive          bool    `json:"is_active"`
}
```

### 1.5 Migration Execution

```go
type BalanceMigration struct {
    ID              string    `json:"id"`
    CompanyID       string    `json:"company_id"`
    FromRegime      string    `json:"from_regime"`     // TT200
    ToRegime        string    `json:"to_regime"`       // TT99
    ExecutionDate   string    `json:"execution_date"`  // 2025-12-31
    Status          string    `json:"status"`           // DRAFT, VALIDATED, EXECUTED, ROLLED_BACK
    SourceBalanceID string    `json:"source_balance_id"`
    TargetBalanceID string    `json:"target_balance_id"`
    JournalEntryID  string    `json:"journal_entry_id,omitempty"`
    Summary         string    `json:"summary,omitempty"`
    ExecutedBy      string    `json:"executed_by"`
    CreatedAt       time.Time `json:"created_at"`
    ExecutedAt      *time.Time `json:"executed_at,omitempty"`
}
```

## 2. PostgreSQL Schema

```sql
-- 004_opening_balance_schema.sql

CREATE TABLE opening_balances (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id      UUID NOT NULL REFERENCES companies(id),
    period_id       UUID NOT NULL REFERENCES period_v2(id),
    fiscal_year_id  UUID NOT NULL REFERENCES fiscal_years(id),
    account_code    VARCHAR(20) NOT NULL REFERENCES accounts(code),
    currency_code   VARCHAR(3) NOT NULL DEFAULT 'VND',
    original_amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    debit_amount    NUMERIC(18,2) NOT NULL DEFAULT 0,
    credit_amount   NUMERIC(18,2) NOT NULL DEFAULT 0,
    exchange_rate   NUMERIC(12,4),
    status          VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    source_type     VARCHAR(20) NOT NULL DEFAULT 'MANUAL',
    batch_id        UUID,
    reason          TEXT,
    approved_by     UUID REFERENCES users(id),
    approved_at     TIMESTAMPTZ,
    corrected_by    UUID REFERENCES users(id),
    correction_of   UUID REFERENCES opening_balances(id),
    created_by      UUID REFERENCES users(id) NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT chk_amount CHECK (debit_amount >= 0 AND credit_amount >= 0),
    CONSTRAINT chk_one_nonzero CHECK (debit_amount > 0 OR credit_amount > 0),
    CONSTRAINT chk_status CHECK (status IN ('DRAFT','PENDING','APPROVED','REJECTED','CORRECTED')),
    CONSTRAINT chk_source CHECK (source_type IN ('MANUAL','CARRY_FORWARD','IMPORT','MIGRATION')),
    CONSTRAINT fk_correction FOREIGN KEY (correction_of) REFERENCES opening_balances(id)
);

CREATE TABLE opening_balance_details (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    opening_balance_id  UUID NOT NULL REFERENCES opening_balances(id) ON DELETE CASCADE,
    entity_type         VARCHAR(20) NOT NULL,
    entity_id           VARCHAR(50) NOT NULL,
    entity_name         VARCHAR(255),
    debit_amount        NUMERIC(18,2) NOT NULL DEFAULT 0,
    credit_amount       NUMERIC(18,2) NOT NULL DEFAULT 0,
    quantity            NUMERIC(18,4),
    unit_price          NUMERIC(18,4),
    original_cost       NUMERIC(18,2),
    acc_depreciation    NUMERIC(18,2),
    counterpart_account VARCHAR(20),
    note               TEXT,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE carry_forward_logs (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id        UUID NOT NULL REFERENCES companies(id),
    from_period_id    UUID NOT NULL REFERENCES period_v2(id),
    to_period_id      UUID NOT NULL REFERENCES period_v2(id),
    from_fiscal_year  INT NOT NULL,
    to_fiscal_year    INT NOT NULL,
    account_count     INT NOT NULL,
    total_debit       NUMERIC(18,2) NOT NULL DEFAULT 0,
    total_credit      NUMERIC(18,2) NOT NULL DEFAULT 0,
    closing_entry_ids UUID[],
    status            VARCHAR(20) NOT NULL DEFAULT 'COMPLETED',
    executed_by       UUID REFERENCES users(id),
    executed_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE circular99_mappings (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    old_account_code    VARCHAR(20) NOT NULL,
    new_account_code    VARCHAR(20) NOT NULL,
    mapping_type        VARCHAR(20) NOT NULL,
    split_ratio         NUMERIC(5,4),
    counterpart_account VARCHAR(20),
    effective_date      DATE NOT NULL,
    note               TEXT,
    is_active          BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE balance_migrations (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id        UUID NOT NULL REFERENCES companies(id),
    from_regime       VARCHAR(10) NOT NULL,
    to_regime         VARCHAR(10) NOT NULL,
    execution_date    DATE NOT NULL,
    status            VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    source_balance_id UUID REFERENCES opening_balances(id),
    target_balance_id UUID REFERENCES opening_balances(id),
    journal_entry_id  UUID REFERENCES journal_entries(id),
    summary           TEXT,
    executed_by       UUID REFERENCES users(id),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    executed_at       TIMESTAMPTZ
);

-- Indexes
CREATE INDEX idx_ob_company_period ON opening_balances(company_id, period_id);
CREATE INDEX idx_ob_account ON opening_balances(account_code);
CREATE INDEX idx_ob_status ON opening_balances(status);
CREATE INDEX idx_ob_source ON opening_balances(source_type);
CREATE INDEX idx_ob_detail_entity ON opening_balance_details(entity_type, entity_id);
CREATE UNIQUE INDEX idx_ob_unique_account_period ON opening_balances(company_id, period_id, account_code, currency_code) WHERE status = 'APPROVED';
```

## 3. Repository Interface

```go
type OpeningBalanceRepository interface {
    // Core CRUD
    Create(ctx context.Context, ob *OpeningBalance) error
    GetByID(ctx context.Context, id string) (*OpeningBalance, error)
    GetByPeriod(ctx context.Context, companyID, periodID string) ([]OpeningBalance, error)
    GetByAccount(ctx context.Context, companyID, periodID, accountCode string) (*OpeningBalance, error)
    GetAll(ctx context.Context, companyID, periodID string) ([]OpeningBalance, error)
    UpdateStatus(ctx context.Context, id string, status OpeningBalanceStatus, approvedBy string) error
    Delete(ctx context.Context, id string) error

    // Bulk operations
    BulkCreate(ctx context.Context, balances []OpeningBalance) error
    BulkUpdateStatus(ctx context.Context, ids []string, status OpeningBalanceStatus) error

    // Detail tracking
    CreateDetail(ctx context.Context, detail *OpeningBalanceDetail) error
    GetDetailsByBalance(ctx context.Context, balanceID string) ([]OpeningBalanceDetail, error)
    GetDetailsByEntity(ctx context.Context, entityType DetailEntityType, entityID string) ([]OpeningBalanceDetail, error)

    // Carry-forward
    CreateCarryForward(ctx context.Context, log *CarryForwardLog, balances []OpeningBalance) error
    GetCarryForwardByCompany(ctx context.Context, companyID string) ([]CarryForwardLog, error)

    // Balance reconciliation
    GetTotals(ctx context.Context, companyID, periodID string) (totalDebit, totalCredit float64, err error)
    ValidateBalance(ctx context.Context, companyID, periodID string) (bool, error)

    // Migration
    CreateMigration(ctx context.Context, m *BalanceMigration) error
    GetMigrationByID(ctx context.Context, id string) (*BalanceMigration, error)
    GetMigrationsByCompany(ctx context.Context, companyID string) ([]BalanceMigration, error)

    // Import
    GetByBatch(ctx context.Context, batchID string) ([]OpeningBalance, error)
}
```

## 4. Service Interface

```go
type OpeningBalanceService interface {
    // CRUD with business logic
    Create(ctx context.Context, ob *OpeningBalance, userID string) error
    GetByPeriod(ctx context.Context, companyID, periodID string) ([]OpeningBalance, error)
    GetByAccount(ctx context.Context, companyID, periodID, accountCode string) (*OpeningBalance, error)
    Update(ctx context.Context, ob *OpeningBalance, userID string) error
    Delete(ctx context.Context, id, userID string) error

    // Approval
    SubmitForApproval(ctx context.Context, id, userID string) error
    Approve(ctx context.Context, id, approverID string) error
    Reject(ctx context.Context, id, approverID, reason string) error

    // Detail management
    SetDetail(ctx context.Context, balanceID string, details []OpeningBalanceDetail) error
    GetDetail(ctx context.Context, balanceID string) ([]OpeningBalanceDetail, error)
    GetBalancesByEntity(ctx context.Context, companyID string, entityType DetailEntityType, entityID string) ([]OpeningBalance, error)

    // Fiscal year operations
    CarryForward(ctx context.Context, companyID, fromPeriodID, toPeriodID, userID string) error
    GetCarryForwardHistory(ctx context.Context, companyID string) ([]CarryForwardLog, error)

    // Bulk import
    ImportFromExcel(ctx context.Context, companyID, periodID, filePath string, encoding string) (*ImportResult, error)
    ImportPreview(ctx context.Context, filePath string) (*ImportPreview, error)

    // Correction
    CorrectBalance(ctx context.Context, balanceID string, newDebit, newCredit float64, reason, userID string) error

    // Circular 99 migration
    GetCircular99Mapping(ctx context.Context) ([]Circular99Mapping, error)
    ExecuteMigration(ctx context.Context, companyID, userID string) (*MigrationResult, error)
    ValidateOpeningBalances(ctx context.Context, companyID, periodID string) (*ValidationResult, error)

    // Reports
    GetOpeningBalanceReport(ctx context.Context, companyID, periodID string) (*OpeningBalanceReport, error)
    GetOpeningBalanceAuditTrail(ctx context.Context, companyID string, from, to time.Time) ([]AuditEntry, error)
}

type ImportResult struct {
    TotalRows    int      `json:"total_rows"`
    SuccessRows  int      `json:"success_rows"`
    ErrorRows    int      `json:"error_rows"`
    Errors       []ImportError `json:"errors,omitempty"`
    BatchID      string   `json:"batch_id"`
}

type ImportPreview struct {
    Columns  []string        `json:"columns"`
    Rows     []ImportRow     `json:"rows"`
    Warnings []string        `json:"warnings,omitempty"`
}

type ImportRow struct {
    RowNumber    int               `json:"row_number"`
    AccountCode  string            `json:"account_code"`
    DebitAmount  float64           `json:"debit_amount"`
    CreditAmount float64           `json:"credit_amount"`
    Currency     string            `json:"currency"`
    Valid        bool              `json:"valid"`
    Errors       []string          `json:"errors,omitempty"`
}

type ValidationResult struct {
    CompanyID     string `json:"company_id"`
    PeriodID      string `json:"period_id"`
    TotalAccounts int    `json:"total_accounts"`
    Balanced      bool   `json:"balanced"`
    TotalDebit    float64 `json:"total_debit"`
    TotalCredit   float64 `json:"total_credit"`
    Diff          float64 `json:"diff"`
    Warnings      []ValidationWarning `json:"warnings,omitempty"`
}

type ValidationWarning struct {
    AccountCode string `json:"account_code"`
    Message     string `json:"message"`
}

type MigrationResult struct {
    TotalAccounts    int     `json:"total_accounts"`
    DirectMapped     int     `json:"direct_mapped"`
    SplitMerged      int     `json:"split_merged"`
    ManualRequired   int     `json:"manual_required"`
    ZeroBalances     int     `json:"zero_balances"`
    JournalEntryID   string  `json:"journal_entry_id,omitempty"`
}

type OpeningBalanceReport struct {
    Period          string                   `json:"period"`
    GeneratedAt     string                   `json:"generated_at"`
    Summary         BalanceSummary            `json:"summary"`
    AccountBalances []AccountBalanceReportItem `json:"account_balances"`
}

type AccountBalanceReportItem struct {
    AccountCode    string  `json:"account_code"`
    AccountName    string  `json:"account_name"`
    AccountType    string  `json:"account_type"`
    NormalBalance  string  `json:"normal_balance"`
    DebitAmount    float64 `json:"debit_amount"`
    CreditAmount   float64 `json:"credit_amount"`
    DetailCount    int     `json:"detail_count,omitempty"`
    Status         string  `json:"status"`
}

type BalanceSummary struct {
    TotalCredit   float64 `json:"total_credit"`
    Diff          float64 `json:"diff"`
    AccountCount  int     `json:"account_count"`
    Balanced      bool    `json:"balanced"`
}
```

## 5. API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET    | /api/v1/companies/:companyId/opening-balances | List OB by period |
| POST   | /api/v1/companies/:companyId/opening-balances | Create OB |
| GET    | /api/v1/companies/:companyId/opening-balances/:id | Get OB by ID |
| PUT    | /api/v1/companies/:companyId/opening-balances/:id | Update OB |
| DELETE | /api/v1/companies/:companyId/opening-balances/:id | Delete OB (DRAFT only) |
| POST   | /api/v1/companies/:companyId/opening-balances/:id/submit | Submit for approval |
| POST   | /api/v1/companies/:companyId/opening-balances/:id/approve | Approve |
| POST   | /api/v1/companies/:companyId/opening-balances/:id/reject | Reject |
| GET    | /api/v1/companies/:companyId/opening-balances/:id/details | Get detail breakdown |
| PUT    | /api/v1/companies/:companyId/opening-balances/:id/details | Set detail breakdown |
| POST   | /api/v1/companies/:companyId/opening-balances/import | Excel import |
| POST   | /api/v1/companies/:companyId/opening-balances/import/preview | Import preview |
| GET    | /api/v1/companies/:companyId/opening-balances/validate | Validate balancing |
| GET    | /api/v1/companies/:companyId/opening-balances/report | Opening balance report |
| POST   | /api/v1/companies/:companyId/opening-balances/carry-forward | Execute carry-forward |
| GET    | /api/v1/companies/:companyId/opening-balances/carry-forward/history | CF history |
| POST   | /api/v1/companies/:companyId/opening-balances/:id/correct | Correct balance |
| GET    | /api/v1/companies/:companyId/opening-balances/migrations | List migrations |
| POST   | /api/v1/companies/:companyId/opening-balances/migrations/execute | Execute C99 migration |

## 6. Integration with Existing Reports

### Trial Balance (`GetTrialBalance`)
Modified to include opening balances:
```go
func (s *service) GetTrialBalance(ctx context.Context, year, month int) ([]domain.AccountBalance, error) {
    period := s.getPeriod(year, month)
    balances := s.journals.GetTrialBalance(ctx, period.ID)
    obs := s.openingBalances.GetByPeriod(ctx, period.ID)
    for i, b := range balances {
        if ob := findOB(obs, b.AccountCode); ob != nil {
            balances[i].OpenBalanceDebit = ob.DebitAmount
            balances[i].OpenBalanceCredit = ob.CreditAmount
        }
        // else: carry forward from previous period
    }
    return balances, nil
}
```

### Balance Sheet (B01-DN)
Column "Số đầu kỳ" populated from opening balances.

### Income Statement
Opening balance adjustments routed through Account 421.

## 7. Validation Engine

All validation rules centralized in `internal/service/opening_balance.go`:

```go
func (s *openingBalanceService) validate(ctx context.Context, ob *OpeningBalance) error {
    // Account must exist and be active
    acc, err := s.accounts.GetByCode(ctx, ob.AccountCode)
    if err != nil { return ErrAccountNotFound }
    if acc.IsParent { return ErrParentAccountBalance }
    if !acc.IsActive { return ErrAccountInactive }
    
    // Period must be open
    period, err := s.periods.GetPeriodByID(ctx, ob.PeriodID)
    if err != nil { return ErrPeriodNotFound }
    if period.Status != PeriodV2Open { return ErrPeriodNotOpen }
    
    // Debit XOR Credit
    if ob.DebitAmount == 0 && ob.CreditAmount == 0 {
        return ErrAmountRequired
    }
    if ob.DebitAmount > 0 && ob.CreditAmount > 0 {
        return ErrBothDebitAndCredit
    }
    
    // Normal balance warning (non-blocking for correction scenarios)
    normalBalance := acc.Type.NormalBalance()
    if (normalBalance == NormalBalanceDebit && ob.CreditAmount > 0) ||
       (normalBalance == NormalBalanceCredit && ob.DebitAmount > 0) {
        // Log warning but allow
    }
    
    // Unique constraint: one approved OB per account per period per currency
    existing, err := s.repo.GetByAccount(ctx, ob.CompanyID, ob.PeriodID, ob.AccountCode)
    if err == nil && existing != nil && existing.Status == OBStatusApproved {
        if ob.ID == "" || ob.ID != existing.ID {
            return ErrDuplicateOpeningBalance
        }
    }
    
    return nil
}
```

## 8. Carry-Forward Algorithm

```
1. Lock: prevent new entries in closing period
2. Identify: accounts with closing balances
3. Zero revenue/expense: create closing entries (5xx,6xx,7xx,8xx → 421)
4. Carry-forward: asset (1xx,2xx), liability (3xx), equity (4xx)
5. Create OB records: status=APPROVED, source=CARRY_FORWARD
6. Audit trail: log every carried-forward balance
7. Unlock: open new period for posting
```

## 9. Excel Import Template

Columns for `opening_balance_import_template.xlsx`:

| Column | Required | Format | Example |
|--------|----------|--------|---------|
| Account Code | Yes | String (3-20 digits) | 1111 |
| Currency | No | ISO 4217 | VND |
| Debit Amount | Yes* | Number (>=0) | 125000000 |
| Credit Amount | Yes* | Number (>=0) | 0 |
| Exchange Rate | No | Number (>0) | 25450 |
| Entity Type | No | Enum | CUSTOMER |
| Entity ID | No | String | KH001 |
| Entity Name | No | String | Company A |
| Quantity | No | Number | 100 |
| Unit Price | No | Number | 50000 |
| Note | No | Text | Opening balance |

*At least one of Debit or Credit must be > 0.
