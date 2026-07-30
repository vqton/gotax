package domain

import "time"

type FixedAssetCategoryGORM struct {
	ID                          string    `gorm:"column:id;primaryKey;size:36"`
	CompanyID                   string    `gorm:"column:company_id;not null;size:36;index:idx_facat_company_code,unique"`
	Code                        string    `gorm:"column:code;not null;size:30;index:idx_facat_company_code,unique"`
	Name                        string    `gorm:"column:name;not null;size:255"`
	ParentID                    *string   `gorm:"column:parent_id;size:36;index"`
	Level                       int       `gorm:"column:level;not null;default:1;index"`
	DefaultUsefulLifeMonths     int       `gorm:"column:default_useful_life_months;not null;default:120"`
	DefaultDepreciationMethod   string    `gorm:"column:default_depreciation_method;not null;size:30;default:STRAIGHT_LINE"`
	AssetAccountID              string    `gorm:"column:asset_account_id;not null;size:36"`
	DepreciationAccountID       string    `gorm:"column:depreciation_account_id;not null;size:36"`
	ExpenseAccountID            string    `gorm:"column:expense_account_id;not null;size:36"`
	CreatedAt                   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt                   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (FixedAssetCategoryGORM) TableName() string { return "fixed_asset_categories" }

type FixedAssetGORM struct {
	ID                      string     `gorm:"column:id;primaryKey;size:36"`
	CompanyID               string     `gorm:"column:company_id;not null;size:36;index:idx_fa_company_code,unique"`
	Code                    string     `gorm:"column:code;not null;size:50;index:idx_fa_company_code,unique"`
	Name                    string     `gorm:"column:name;not null;size:255"`
	CategoryID              string     `gorm:"column:category_id;not null;size:36;index"`
	Status                  string     `gorm:"column:status;not null;size:20;default:DRAFT;index"`
	AcquisitionDate         time.Time  `gorm:"column:acquisition_date;not null;type:date"`
	OriginalCost            float64    `gorm:"column:original_cost;not null;default:0"`
	AccumulatedDepreciation float64    `gorm:"column:accumulated_depreciation;not null;default:0"`
	ResidualValue           float64    `gorm:"column:residual_value;not null;default:0"`
	CarryingAmount          float64    `gorm:"column:carrying_amount;not null;default:0"`
	UsefulLifeMonths        int        `gorm:"column:useful_life_months;not null"`
	DepreciationMethod      string     `gorm:"column:depreciation_method;not null;size:30;default:STRAIGHT_LINE"`
	DepreciationStartDate   *time.Time `gorm:"column:depreciation_start_date;type:date"`
	DepreciationEndDate     *time.Time `gorm:"column:depreciation_end_date;type:date"`
	DepartmentID            string     `gorm:"column:department_id;not null;size:36;index"`
	Location                string     `gorm:"column:location;not null;size:255;default:''"`
	UserID                  *string    `gorm:"column:user_id;size:36"`
	SupplierID              *string    `gorm:"column:supplier_id;size:36"`
	ContractNo              *string    `gorm:"column:contract_no;size:50"`
	InvoiceID               *string    `gorm:"column:invoice_id;size:36"`
	SerialNo                *string    `gorm:"column:serial_no;size:100"`
	Manufacturer            *string    `gorm:"column:manufacturer;size:255"`
	ManufactureYear         int        `gorm:"column:manufacture_year;default:0"`
	CountryOfOrigin         *string    `gorm:"column:country_of_origin;size:100"`
	TechnicalSpecs          *string    `gorm:"column:technical_specs;type:text"`
	Notes                   *string    `gorm:"column:notes;type:text"`
	Source                  string     `gorm:"column:source;not null;size:30;default:PURCHASE;index"`
	CIPAccountID            *string    `gorm:"column:cip_account_id;size:36"`
	AssetAccountID          string     `gorm:"column:asset_account_id;not null;size:36"`
	DepreciationAccountID   string     `gorm:"column:depreciation_account_id;not null;size:36"`
	ExpenseAccountID        string     `gorm:"column:expense_account_id;not null;size:36"`
	CreatedAt               time.Time  `gorm:"column:created_at;autoCreateTime"`
	CreatedBy               string     `gorm:"column:created_by;not null;size:36"`
	UpdatedAt               time.Time  `gorm:"column:updated_at;autoUpdateTime"`
	UpdatedBy               string     `gorm:"column:updated_by;not null;size:36"`
}

func (FixedAssetGORM) TableName() string { return "fixed_assets" }

type DepreciationEntryGORM struct {
	ID                  string    `gorm:"column:id;primaryKey;size:36"`
	CompanyID           string    `gorm:"column:company_id;not null;size:36;index"`
	FixedAssetID        string    `gorm:"column:fixed_asset_id;not null;size:36;index:idx_depr_asset_period,unique"`
	PeriodID            string    `gorm:"column:period_id;not null;size:36;index:idx_depr_asset_period,unique;index"`
	PeriodYear          int       `gorm:"column:period_year;not null"`
	PeriodMonth         int       `gorm:"column:period_month;not null"`
	DepreciationAmount  float64   `gorm:"column:depreciation_amount;not null"`
	AccumulatedAfter    float64   `gorm:"column:accumulated_after;not null"`
	CarryingAmountAfter float64   `gorm:"column:carrying_amount_after;not null"`
	GLPosted            bool      `gorm:"column:gl_posted;not null;default:false;index"`
	GLJournalEntryID    *string   `gorm:"column:gl_journal_entry_id;size:36"`
	CreatedAt           time.Time `gorm:"column:created_at;autoCreateTime"`
	CreatedBy           string    `gorm:"column:created_by;not null;size:36"`
}

func (DepreciationEntryGORM) TableName() string { return "depreciation_entries" }

type FixedAssetTransactionGORM struct {
	ID              string    `gorm:"column:id;primaryKey;size:36"`
	CompanyID       string    `gorm:"column:company_id;not null;size:36"`
	FixedAssetID    string    `gorm:"column:fixed_asset_id;not null;size:36;index"`
	TransactionType string    `gorm:"column:transaction_type;not null;size:30;index"`
	TransactionDate time.Time `gorm:"column:transaction_date;not null;type:date;index"`
	Amount          float64   `gorm:"column:amount;not null"`
	OldValue        float64   `gorm:"column:old_value;default:0"`
	NewValue        float64   `gorm:"column:new_value;default:0"`
	Description     *string   `gorm:"column:description;type:text"`
	GLJournalID     *string   `gorm:"column:gl_journal_id;size:36"`
	CreatedBy       string    `gorm:"column:created_by;not null;size:36"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (FixedAssetTransactionGORM) TableName() string { return "fixed_asset_transactions" }

type FixedAssetAllocationGORM struct {
	ID               string `gorm:"column:id;primaryKey;size:36"`
	FixedAssetID     string `gorm:"column:fixed_asset_id;not null;size:36;index"`
	DepartmentID     string `gorm:"column:department_id;not null;size:36;index"`
	AllocationPct    float64 `gorm:"column:allocation_pct;not null"`
	ExpenseAccountID string `gorm:"column:expense_account_id;not null;size:36"`
}

func (FixedAssetAllocationGORM) TableName() string { return "fixed_asset_allocations" }

type FixedAssetInventoryPlanGORM struct {
	ID        string    `gorm:"column:id;primaryKey;size:36"`
	CompanyID string    `gorm:"column:company_id;not null;size:36;index"`
	PlanDate  time.Time `gorm:"column:plan_date;not null;type:date;index"`
	Status    string    `gorm:"column:status;not null;size:20;default:DRAFT;index"`
	Notes     *string   `gorm:"column:notes;type:text"`
	CreatedBy string    `gorm:"column:created_by;not null;size:36"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (FixedAssetInventoryPlanGORM) TableName() string { return "fixed_asset_inventory_plans" }

type FixedAssetInventoryResultGORM struct {
	ID               string  `gorm:"column:id;primaryKey;size:36"`
	PlanID           string  `gorm:"column:plan_id;not null;size:36;index:idx_fair_plan_asset,unique"`
	FixedAssetID     string  `gorm:"column:fixed_asset_id;not null;size:36;index:idx_fair_plan_asset,unique"`
	ExpectedLocation *string `gorm:"column:expected_location;size:255"`
	ActualLocation   *string `gorm:"column:actual_location;size:255"`
	ExpectedStatus   *string `gorm:"column:expected_status;size:30"`
	ActualStatus     *string `gorm:"column:actual_status;size:30"`
	Discrepancy      string  `gorm:"column:discrepancy;not null;size:50;index"`
	Notes            *string `gorm:"column:notes;type:text"`
}

func (FixedAssetInventoryResultGORM) TableName() string { return "fixed_asset_inventory_results" }
