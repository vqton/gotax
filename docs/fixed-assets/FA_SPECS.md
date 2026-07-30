# Fixed Asset Module — Technical Specifications

**Version:** 1.0  
**Date:** 2026-07-30  

---

## 1. Data Models

### 1.1 FixedAsset

```go
type FixedAssetStatus string
const (
    FADraft        FixedAssetStatus = "DRAFT"
    FACancelled    FixedAssetStatus = "CANCELLED"
    FAActive       FixedAssetStatus = "ACTIVE"
    FADepreciating FixedAssetStatus = "DEPRECIATING"
    FASuspended    FixedAssetStatus = "SUSPENDED"
    FAFullyDepr    FixedAssetStatus = "FULLY_DEPRECIATED"
    FADisposed     FixedAssetStatus = "DISPOSED"
    FASold         FixedAssetStatus = "SOLD"
)

type DepreciationMethod string
const (
    DepStraightLine       DepreciationMethod = "STRAIGHT_LINE"
    DepDecliningBalance   DepreciationMethod = "DECLINING_BALANCE"
    DepProductionBased    DepreciationMethod = "PRODUCTION_BASED"
)

type FixedAsset struct {
    ID                  string             `json:"id"`
    CompanyID           string             `json:"company_id"`
    Code                string             `json:"code"`
    Name                string             `json:"name"`
    CategoryID          string             `json:"category_id"`
    Status              FixedAssetStatus   `json:"status"`
    AcquisitionDate     time.Time          `json:"acquisition_date"`
    OriginalCost        float64            `json:"original_cost"`
    AccumulatedDepreciation float64        `json:"accumulated_depreciation"`
    ResidualValue       float64            `json:"residual_value"`
    CarryingAmount      float64            `json:"carrying_amount"`
    UsefulLifeMonths    int                `json:"useful_life_months"`
    DepreciationMethod  DepreciationMethod `json:"depreciation_method"`
    DepreciationStartDate *time.Time       `json:"depreciation_start_date,omitempty"`
    DepreciationEndDate   *time.Time       `json:"depreciation_end_date,omitempty"`
    DepartmentID        string             `json:"department_id"`
    Location            string             `json:"location"`
    UserID              string             `json:"user_id"`
    SupplierID          string             `json:"supplier_id,omitempty"`
    ContractNo          string             `json:"contract_no,omitempty"`
    InvoiceID           string             `json:"invoice_id,omitempty"`
    SerialNo            string             `json:"serial_no,omitempty"`
    Manufacturer        string             `json:"manufacturer,omitempty"`
    ManufactureYear     int                `json:"manufacture_year,omitempty"`
    CountryOfOrigin     string             `json:"country_of_origin,omitempty"`
    TechnicalSpecs      string             `json:"technical_specs,omitempty"`
    Notes               string             `json:"notes,omitempty"`
    Source              string             `json:"source"` // PURCHASE, CONSTRUCTION, LEASE, DONATION, CAPITAL_CONTRIBUTION, EXCHANGE
    CIPAccountID        string             `json:"cip_account_id,omitempty"` // 2411/2412
    AssetAccountID      string             `json:"asset_account_id"`          // 211x/212/213x
    DepreciationAccountID string           `json:"depreciation_account_id"`   // 2141/2142/2143
    ExpenseAccountID    string             `json:"expense_account_id"`        // 6274/6414/6424
    CreatedAt           time.Time          `json:"created_at"`
    CreatedBy           string             `json:"created_by"`
    UpdatedAt           time.Time          `json:"updated_at"`
    UpdatedBy           string             `json:"updated_by"`
}
```

### 1.2 FixedAssetCategory

```go
type FixedAssetCategory struct {
    ID                    string             `json:"id"`
    CompanyID             string             `json:"company_id"`
    Code                  string             `json:"code"`
    Name                  string             `json:"name"`
    ParentID              *string            `json:"parent_id,omitempty"`
    Level                 int                `json:"level"` // 1=Group, 2=Type, 3=Category
    DefaultUsefulLifeMonths int              `json:"default_useful_life_months"`
    DefaultDepreciationMethod DepreciationMethod `json:"default_depreciation_method"`
    AssetAccountID        string             `json:"asset_account_id"`         // Default 211x
    DepreciationAccountID string             `json:"depreciation_account_id"`  // Default 214x
    ExpenseAccountID      string             `json:"expense_account_id"`       // Default 6274/6414/6424
    CreatedAt             time.Time          `json:"created_at"`
    UpdatedAt             time.Time          `json:"updated_at"`
}
```

### 1.3 DepreciationEntry

```go
type DepreciationEntry struct {
    ID                string    `json:"id"`
    CompanyID         string    `json:"company_id"`
    FixedAssetID      string    `json:"fixed_asset_id"`
    PeriodID          string    `json:"period_id"`
    PeriodYear        int       `json:"period_year"`
    PeriodMonth       int       `json:"period_month"`
    DepreciationAmount float64  `json:"depreciation_amount"`
    AccumulatedAfter  float64   `json:"accumulated_after"`
    CarryingAmountAfter float64 `json:"carrying_amount_after"`
    GLPosted          bool      `json:"gl_posted"`
    GLJournalEntryID  string    `json:"gl_journal_entry_id,omitempty"`
    CreatedAt         time.Time `json:"created_at"`
    CreatedBy         string    `json:"created_by"`
}
```

### 1.4 FixedAssetTransaction

```go
type FATransactionType string
const (
    FATrxAcquisition     FATransactionType = "ACQUISITION"
    FATrxDepreciation    FATransactionType = "DEPRECIATION"
    FATrxAdjustment      FATransactionType = "ADJUSTMENT"
    FATrxTransfer       FATransactionType = "TRANSFER"
    FATrxDisposal       FATransactionType = "DISPOSAL"
    FATrxSale           FATransactionType = "SALE"
    FATrxRevaluation    FATransactionType = "REVALUATION"
    FATrxImpairment     FATransactionType = "IMPAIRMENT"
    FATrxCIPTransfer    FATransactionType = "CIP_TRANSFER"
)

type FixedAssetTransaction struct {
    ID              string             `json:"id"`
    CompanyID       string             `json:"company_id"`
    FixedAssetID    string             `json:"fixed_asset_id"`
    TransactionType FATransactionType  `json:"transaction_type"`
    TransactionDate time.Time          `json:"transaction_date"`
    Amount          float64            `json:"amount"`
    OldValue        float64            `json:"old_value"`
    NewValue        float64            `json:"new_value"`
    Description     string             `json:"description"`
    GLJournalID     string             `json:"gl_journal_id,omitempty"`
    CreatedBy       string             `json:"created_by"`
    CreatedAt       time.Time          `json:"created_at"`
}
```

### 1.5 FixedAssetAllocation

```go
type FixedAssetAllocation struct {
    ID             string  `json:"id"`
    FixedAssetID   string  `json:"fixed_asset_id"`
    DepartmentID   string  `json:"department_id"`
    AllocationPct  float64 `json:"allocation_pct"` // Sum = 100%
    ExpenseAccountID string `json:"expense_account_id"`
}
```

### 1.6 FixedAssetInventory

```go
type FixedAssetInventoryPlan struct {
    ID          string    `json:"id"`
    CompanyID   string    `json:"company_id"`
    PlanDate    time.Time `json:"plan_date"`
    Status      string    `json:"status"` // DRAFT, IN_PROGRESS, COMPLETED
    Notes       string    `json:"notes"`
    CreatedBy   string    `json:"created_by"`
    CreatedAt   time.Time `json:"created_at"`
}

type FixedAssetInventoryResult struct {
    ID               string    `json:"id"`
    PlanID           string    `json:"plan_id"`
    FixedAssetID     string    `json:"fixed_asset_id"`
    ExpectedLocation string    `json:"expected_location"`
    ActualLocation   string    `json:"actual_location"`
    ExpectedStatus   string    `json:"expected_status"`
    ActualStatus     string    `json:"actual_status"`
    Discrepancy      string    `json:"discrepancy"` // OK, MISSING, DAMAGED, UNREGISTERED
    Notes            string    `json:"notes"`
}
```

---

## 2. Database Schema

### 2.1 Migration: 010_fixed_asset_schema.up.sql

```sql
CREATE TABLE fixed_asset_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    code VARCHAR(20) NOT NULL,
    name VARCHAR(255) NOT NULL,
    parent_id UUID REFERENCES fixed_asset_categories(id),
    level SMALLINT NOT NULL CHECK (level BETWEEN 1 AND 3),
    default_useful_life_months INTEGER NOT NULL DEFAULT 120,
    default_depreciation_method VARCHAR(30) NOT NULL DEFAULT 'STRAIGHT_LINE',
    asset_account_id VARCHAR(20) NOT NULL,
    depreciation_account_id VARCHAR(20) NOT NULL,
    expense_account_id VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, code)
);

CREATE TABLE fixed_assets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    code VARCHAR(50) NOT NULL,
    name VARCHAR(500) NOT NULL,
    category_id UUID NOT NULL REFERENCES fixed_asset_categories(id),
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    acquisition_date DATE NOT NULL,
    original_cost DECIMAL(20,2) NOT NULL CHECK (original_cost >= 0),
    accumulated_depreciation DECIMAL(20,2) NOT NULL DEFAULT 0,
    residual_value DECIMAL(20,2) NOT NULL DEFAULT 0,
    carrying_amount DECIMAL(20,2) GENERATED ALWAYS AS (original_cost - accumulated_depreciation) STORED,
    useful_life_months INTEGER NOT NULL CHECK (useful_life_months > 0),
    depreciation_method VARCHAR(30) NOT NULL DEFAULT 'STRAIGHT_LINE',
    depreciation_start_date DATE,
    depreciation_end_date DATE,
    department_id UUID REFERENCES departments(id),
    location VARCHAR(255),
    user_id UUID REFERENCES employees(id),
    supplier_id UUID REFERENCES suppliers(id),
    contract_no VARCHAR(100),
    invoice_id UUID REFERENCES supplier_invoices(id),
    serial_no VARCHAR(100),
    manufacturer VARCHAR(255),
    manufacture_year INTEGER,
    country_of_origin VARCHAR(100),
    technical_specs TEXT,
    notes TEXT,
    source VARCHAR(30) NOT NULL DEFAULT 'PURCHASE',
    cip_account_id VARCHAR(20),
    asset_account_id VARCHAR(20) NOT NULL,
    depreciation_account_id VARCHAR(20) NOT NULL,
    expense_account_id VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL REFERENCES users(id),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID NOT NULL REFERENCES users(id),
    UNIQUE(company_id, code)
);

CREATE TABLE fixed_asset_allocations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fixed_asset_id UUID NOT NULL REFERENCES fixed_assets(id) ON DELETE CASCADE,
    department_id UUID NOT NULL REFERENCES departments(id),
    allocation_pct DECIMAL(5,2) NOT NULL CHECK (allocation_pct > 0 AND allocation_pct <= 100),
    expense_account_id VARCHAR(20) NOT NULL,
    UNIQUE(fixed_asset_id, department_id)
);

CREATE TABLE depreciation_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    fixed_asset_id UUID NOT NULL REFERENCES fixed_assets(id),
    period_id UUID NOT NULL REFERENCES periods(id),
    period_year INTEGER NOT NULL,
    period_month INTEGER NOT NULL CHECK (period_month BETWEEN 1 AND 12),
    depreciation_amount DECIMAL(20,2) NOT NULL CHECK (depreciation_amount >= 0),
    accumulated_after DECIMAL(20,2) NOT NULL,
    carrying_amount_after DECIMAL(20,2) NOT NULL,
    gl_posted BOOLEAN NOT NULL DEFAULT FALSE,
    gl_journal_entry_id UUID REFERENCES journal_entries(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL REFERENCES users(id),
    UNIQUE(fixed_asset_id, period_id)
);

CREATE TABLE fixed_asset_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    fixed_asset_id UUID NOT NULL REFERENCES fixed_assets(id),
    transaction_type VARCHAR(30) NOT NULL,
    transaction_date DATE NOT NULL,
    amount DECIMAL(20,2) NOT NULL,
    old_value DECIMAL(20,2),
    new_value DECIMAL(20,2),
    description TEXT,
    gl_journal_id UUID REFERENCES journal_entries(id),
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE fixed_asset_inventory_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    plan_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    notes TEXT,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE fixed_asset_inventory_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id UUID NOT NULL REFERENCES fixed_asset_inventory_plans(id) ON DELETE CASCADE,
    fixed_asset_id UUID NOT NULL REFERENCES fixed_assets(id),
    expected_location VARCHAR(255),
    actual_location VARCHAR(255),
    expected_status VARCHAR(30),
    actual_status VARCHAR(30),
    discrepancy VARCHAR(20) NOT NULL DEFAULT 'OK',
    notes TEXT
);

CREATE INDEX idx_fixed_assets_company ON fixed_assets(company_id);
CREATE INDEX idx_fixed_assets_status ON fixed_assets(status);
CREATE INDEX idx_fixed_assets_category ON fixed_assets(category_id);
CREATE INDEX idx_fixed_assets_department ON fixed_assets(department_id);
CREATE INDEX idx_depreciation_entries_asset ON depreciation_entries(fixed_asset_id);
CREATE INDEX idx_depreciation_entries_period ON depreciation_entries(period_id);
CREATE INDEX idx_depreciation_entries_gl ON depreciation_entries(gl_posted);
CREATE INDEX idx_fa_transactions_asset ON fixed_asset_transactions(fixed_asset_id);
```

---

## 3. API Endpoints

### 3.1 FA Categories

| Method | Path | Handler | Auth |
|--------|------|---------|------|
| POST | `/api/v1/fixed-assets/categories` | `CreateFACategory` | admin, chief |
| GET | `/api/v1/fixed-assets/categories` | `ListFACategories` | auth |
| GET | `/api/v1/fixed-assets/categories/:id` | `GetFACategory` | auth |
| PUT | `/api/v1/fixed-assets/categories/:id` | `UpdateFACategory` | admin, chief |
| DELETE | `/api/v1/fixed-assets/categories/:id` | `DeleteFACategory` | admin |

### 3.2 Fixed Assets

| Method | Path | Handler | Auth |
|--------|------|---------|------|
| POST | `/api/v1/fixed-assets` | `CreateFixedAsset` | auth |
| GET | `/api/v1/fixed-assets` | `ListFixedAssets` | auth |
| GET | `/api/v1/fixed-assets/:id` | `GetFixedAsset` | auth |
| PUT | `/api/v1/fixed-assets/:id` | `UpdateFixedAsset` | auth |
| PATCH | `/api/v1/fixed-assets/:id/activate` | `ActivateFixedAsset` | chief |
| PATCH | `/api/v1/fixed-assets/:id/suspend` | `SuspendFixedAsset` | auth |
| PATCH | `/api/v1/fixed-assets/:id/resume` | `ResumeFixedAsset` | auth |
| PATCH | `/api/v1/fixed-assets/:id/dispose` | `DisposeFixedAsset` | chief |
| PATCH | `/api/v1/fixed-assets/:id/sell` | `SellFixedAsset` | chief |
| POST | `/api/v1/fixed-assets/:id/adjust` | `AdjustFixedAsset` | chief |
| POST | `/api/v1/fixed-assets/:id/revalue` | `RevalueFixedAsset` | admin |
| POST | `/api/v1/fixed-assets/:id/impair` | `ImpairFixedAsset` | chief |
| POST | `/api/v1/fixed-assets/:id/transfer` | `TransferFixedAsset` | auth |
| POST | `/api/v1/fixed-assets/:id/allocate` | `SetAllocation` | auth |
| POST | `/api/v1/fixed-assets/:id/cip-to-fa` | `CIPToFA` | chief |

### 3.3 Depreciation

| Method | Path | Handler | Auth |
|--------|------|---------|------|
| POST | `/api/v1/fixed-assets/depreciation/run` | `RunDepreciation` | chief |
| POST | `/api/v1/fixed-assets/depreciation/post` | `PostDepreciation` | chief |
| POST | `/api/v1/fixed-assets/depreciation/unpost` | `UnpostDepreciation` | admin |
| GET | `/api/v1/fixed-assets/:id/depreciation` | `GetDepreciationSchedule` | auth |
| GET | `/api/v1/fixed-assets/depreciation/entries` | `ListDepreciationEntries` | auth |

### 3.4 FA Inventory

| Method | Path | Handler | Auth |
|--------|------|---------|------|
| POST | `/api/v1/fixed-assets/inventory/plans` | `CreateInventoryPlan` | chief |
| GET | `/api/v1/fixed-assets/inventory/plans` | `ListInventoryPlans` | auth |
| GET | `/api/v1/fixed-assets/inventory/plans/:id` | `GetInventoryPlan` | auth |
| POST | `/api/v1/fixed-assets/inventory/plans/:id/start` | `StartInventory` | chief |
| POST | `/api/v1/fixed-assets/inventory/results` | `RecordInventoryResult` | auth |
| POST | `/api/v1/fixed-assets/inventory/plans/:id/complete` | `CompleteInventory` | chief |

### 3.5 Reports

| Method | Path | Handler | Auth |
|--------|------|---------|------|
| GET | `/api/v1/reports/fa-register` | `FARegister` | auth |
| GET | `/api/v1/reports/fa-depreciation-schedule` | `FADepreciationSchedule` | auth |
| GET | `/api/v1/reports/fa-movement` | `FAMovement` | auth |
| GET | `/api/v1/reports/fa-inventory` | `FAInventoryReport` | auth |
| GET | `/api/v1/reports/fa-aging` | `FAAging` | auth |
| GET | `/api/v1/reports/fa-by-department` | `FAByDepartment` | auth |

### 3.6 Import/Export

| Method | Path | Handler | Auth |
|--------|------|---------|------|
| POST | `/api/v1/fixed-assets/import` | `ImportFixedAssets` | admin |
| GET | `/api/v1/fixed-assets/export` | `ExportFixedAssets` | admin |

---

## 4. Depreciation Calculation Engine

### 4.1 Straight-Line Method

```
Monthly depreciation = (OriginalCost - ResidualValue) / UsefulLifeMonths
First period prorata = MonthlyDepreciation × (DaysRemainingInMonth / TotalDaysInMonth)
Last period adjustment = RemainingCarryingAmount - ResidualValue
```

### 4.2 Declining Balance Method

```
Depreciation rate = (1 / UsefulLifeYears) × DecliningFactor
Monthly depreciation = CarryingAmount × DepreciationRate / 12
DecliningFactor = 2 (double-declining) or configurable
Switch to straight-line when SL depreciation > Declining depreciation (optional per IAS 16)
```

### 4.3 Production-Based Method

```
Depreciation per unit = (OriginalCost - ResidualValue) / TotalEstimatedProduction
Monthly depreciation = DepreciationPerUnit × ActualProductionInMonth
```

### 4.4 Algorithm

```
Input: FixedAsset, Period (year, month)
Output: DepreciationEntry

1. Validate asset status = DEPRECIATING
2. Validate period not already calculated (check depreciation_entries)
3. Validate period is after DepreciationStartDate, before DepreciationEndDate
4. Calculate monthly amount based on method
5. If first period: apply prorata temporis
6. If last period: adjust to make carrying amount = residual value
7. If fully depreciated: set entry amount = remaining carrying amount - residual value
8. Create DepreciationEntry with status unposted
```

---

## 5. State Machine

```
Valid transitions:

DRAFT → CANCELLED      (cancel registration)
DRAFT → ACTIVE          (activate asset, set dep start date)
ACTIVE → DEPRECIATING   (start depreciation — first period run)
DEPRECIATING → SUSPENDED   (temporary stop)
SUSPENDED → DEPRECIATING   (resume depreciation)
DEPRECIATING → FULLY_DEPRECIATED  (carrying amount = 0)
FULLY_DEPRECIATED → DISPOSED  (scrap/liquidate)
FULLY_DEPRECIATED → SOLD      (sell with proceeds)
DEPRECIATING → DISPOSED   (early disposal)
DEPRECIATING → SOLD       (early sale)
```

---

## 6. Integration Architecture

```
┌─────────────┐     ┌──────────────────┐     ┌──────────────┐
│   Handler   │────→│    Service       │────→│  Repository  │
│ (FAHandler) │     │ (FAService)      │     │  (PG + Mem)  │
└─────────────┘     └──────────────────┘     └──────────────┘
                          │
                          │ calls
                          v
                    ┌──────────────────┐
                    │  Depreciation    │
                    │  Engine          │
                    │  (pure Go calc)  │
                    └──────────────────┘
                          │
                          │ posts to
                          v
                    ┌──────────────────┐
                    │  GL Service      │
                    │  (journal_entry) │
                    └──────────────────┘
```
