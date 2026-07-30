package domain

type FixedAssetStatus string
const (
	FADraft          FixedAssetStatus = "DRAFT"
	FACancelled      FixedAssetStatus = "CANCELLED"
	FAActive         FixedAssetStatus = "ACTIVE"
	FADepreciating   FixedAssetStatus = "DEPRECIATING"
	FASuspended      FixedAssetStatus = "SUSPENDED"
	FAFullyDepr      FixedAssetStatus = "FULLY_DEPRECIATED"
	FADisposed       FixedAssetStatus = "DISPOSED"
	FASold           FixedAssetStatus = "SOLD"
)

type DepreciationMethod string
const (
	DepStraightLine     DepreciationMethod = "STRAIGHT_LINE"
	DepDecliningBalance DepreciationMethod = "DECLINING_BALANCE"
	DepProductionBased  DepreciationMethod = "PRODUCTION_BASED"
)

type DisposalType string
const (
	DisposalSale       DisposalType = "SALE"
	DisposalLiquidation DisposalType = "LIQUIDATION"
	DisposalDonation   DisposalType = "DONATION"
	DisposalReturn     DisposalType = "RETURN"
)

type FASource string
const (
	FASourcePurchase          FASource = "PURCHASE"
	FASourceConstruction      FASource = "CONSTRUCTION"
	FASourceLease             FASource = "LEASE"
	FASourceDonation          FASource = "DONATION"
	FASourceCapitalContribution FASource = "CAPITAL_CONTRIBUTION"
	FASourceExchange          FASource = "EXCHANGE"
)

type FATransactionType string
const (
	FATrxAcquisition  FATransactionType = "ACQUISITION"
	FATrxDepreciation FATransactionType = "DEPRECIATION"
	FATrxAdjustment   FATransactionType = "ADJUSTMENT"
	FATrxTransfer     FATransactionType = "TRANSFER"
	FATrxDisposal     FATransactionType = "DISPOSAL"
	FATrxSale         FATransactionType = "SALE"
	FATrxRevaluation  FATransactionType = "REVALUATION"
	FATrxImpairment   FATransactionType = "IMPAIRMENT"
	FATrxCIPTransfer  FATransactionType = "CIP_TRANSFER"
	FATrxSuspend      FATransactionType = "SUSPEND"
	FATrxResume       FATransactionType = "RESUME"
)

type FixedAsset struct {
	ID                     string             `json:"id"`
	CompanyID              string             `json:"company_id"`
	Code                   string             `json:"code"                      validate:"required,max=50"`
	Name                   string             `json:"name"                      validate:"required,max=200"`
	CategoryID             string             `json:"category_id"               validate:"required"`
	Status                 FixedAssetStatus   `json:"status"                    validate:"omitempty,fastatus"`
	AcquisitionDate        string             `json:"acquisition_date"          validate:"required,datetime=2006-01-02"`
	OriginalCost           float64            `json:"original_cost"             validate:"gte=0"`
	AccumulatedDepreciation float64           `json:"accumulated_depreciation"`
	ResidualValue          float64            `json:"residual_value"`
	CarryingAmount         float64            `json:"carrying_amount"`
	UsefulLifeMonths       int                `json:"useful_life_months"        validate:"gt=0"`
	DepreciationMethod     DepreciationMethod `json:"depreciation_method"       validate:"omitempty,damethod"`
	DepreciationStartDate  *string            `json:"depreciation_start_date,omitempty"`
	DepreciationEndDate    *string            `json:"depreciation_end_date,omitempty"`
	DepartmentID           string             `json:"department_id"             validate:"required"`
	Location               string             `json:"location"`
	UserID                 string             `json:"user_id,omitempty"`
	SupplierID             string             `json:"supplier_id,omitempty"`
	ContractNo             string             `json:"contract_no,omitempty"`
	InvoiceID              string             `json:"invoice_id,omitempty"`
	SerialNo               string             `json:"serial_no,omitempty"`
	Manufacturer           string             `json:"manufacturer,omitempty"`
	ManufactureYear        int                `json:"manufacture_year,omitempty"`
	CountryOfOrigin        string             `json:"country_of_origin,omitempty"`
	TechnicalSpecs         string             `json:"technical_specs,omitempty"`
	Notes                  string             `json:"notes,omitempty"`
	Source                 FASource           `json:"source"                    validate:"omitempty,fasource"`
	CIPAccountID           string             `json:"cip_account_id,omitempty"`
	AssetAccountID         string             `json:"asset_account_id"         validate:"required"`
	DepreciationAccountID  string             `json:"depreciation_account_id"  validate:"required"`
	ExpenseAccountID       string             `json:"expense_account_id"       validate:"required"`
	CreatedAt              string             `json:"created_at"`
	CreatedBy              string             `json:"created_by"`
	UpdatedAt              string             `json:"updated_at"`
	UpdatedBy              string             `json:"updated_by"`
}

func (f *FixedAsset) CalcCarryingAmount() float64 {
	return f.OriginalCost - f.AccumulatedDepreciation
}

func (f *FixedAsset) Validate() error {
	if f.Code == "" {
		return ErrFACodeRequired
	}
	if f.Name == "" {
		return ErrFANameRequired
	}
	if f.CategoryID == "" {
		return ErrFACategoryRequired
	}
	if f.OriginalCost < 0 {
		return ErrFAOriginalCostInvalid
	}
	if f.UsefulLifeMonths <= 0 {
		return ErrFAUsefulLifeInvalid
	}
	if f.DepreciationMethod == "" {
		f.DepreciationMethod = DepStraightLine
	}
	switch f.DepreciationMethod {
	case DepStraightLine, DepDecliningBalance, DepProductionBased:
	default:
		return ErrFADepreciationMethodInvalid
	}
	if f.Status == "" {
		f.Status = FADraft
	}
	switch f.Status {
	case FADraft, FACancelled, FAActive, FADepreciating, FASuspended, FAFullyDepr, FADisposed, FASold:
	default:
		return ErrFAStatusInvalid
	}
	if f.Source == "" {
		f.Source = FASourcePurchase
	}
	if f.AssetAccountID == "" {
		return ErrFAAssetAccountRequired
	}
	if f.DepreciationAccountID == "" {
		return ErrFADepreciationAccountRequired
	}
	if f.ExpenseAccountID == "" {
		return ErrFAExpenseAccountRequired
	}
	return nil
}

type FixedAssetCategory struct {
	ID                          string             `json:"id"`
	CompanyID                   string             `json:"company_id"`
	Code                        string             `json:"code"                           validate:"required,max=50"`
	Name                        string             `json:"name"                           validate:"required,max=200"`
	ParentID                    *string            `json:"parent_id,omitempty"`
	Level                       int                `json:"level"                          validate:"min=1,max=3"`
	DefaultUsefulLifeMonths     int                `json:"default_useful_life_months"`
	DefaultDepreciationMethod   DepreciationMethod `json:"default_depreciation_method"    validate:"omitempty,damethod"`
	AssetAccountID              string             `json:"asset_account_id"              validate:"required"`
	DepreciationAccountID       string             `json:"depreciation_account_id"       validate:"required"`
	ExpenseAccountID            string             `json:"expense_account_id"            validate:"required"`
	CreatedAt                   string             `json:"created_at"`
	UpdatedAt                   string             `json:"updated_at"`
}

func (c *FixedAssetCategory) Validate() error {
	if c.Code == "" {
		return ErrFACategoryCodeRequired
	}
	if c.Name == "" {
		return ErrFACategoryNameRequired
	}
	if c.Level < 1 || c.Level > 3 {
		return ErrFACategoryLevelInvalid
	}
	if c.DefaultUsefulLifeMonths <= 0 {
		c.DefaultUsefulLifeMonths = 120
	}
	if c.DefaultDepreciationMethod == "" {
		c.DefaultDepreciationMethod = DepStraightLine
	}
	if c.AssetAccountID == "" {
		return ErrFAAssetAccountRequired
	}
	if c.DepreciationAccountID == "" {
		return ErrFADepreciationAccountRequired
	}
	if c.ExpenseAccountID == "" {
		return ErrFAExpenseAccountRequired
	}
	if c.ParentID != nil && *c.ParentID == c.ID {
		return ErrFACategorySelfParent
	}
	return nil
}

type DepreciationEntry struct {
	ID                  string  `json:"id"`
	CompanyID           string  `json:"company_id"`
	FixedAssetID        string  `json:"fixed_asset_id"`
	PeriodID            string  `json:"period_id"`
	PeriodYear          int     `json:"period_year"`
	PeriodMonth         int     `json:"period_month"`
	DepreciationAmount  float64 `json:"depreciation_amount"`
	AccumulatedAfter    float64 `json:"accumulated_after"`
	CarryingAmountAfter float64 `json:"carrying_amount_after"`
	GLPosted            bool    `json:"gl_posted"`
	GLJournalEntryID    string  `json:"gl_journal_entry_id,omitempty"`
	CreatedAt           string  `json:"created_at"`
	CreatedBy           string  `json:"created_by"`
}

type FixedAssetTransaction struct {
	ID              string            `json:"id"`
	CompanyID       string            `json:"company_id"`
	FixedAssetID    string            `json:"fixed_asset_id"`
	TransactionType FATransactionType `json:"transaction_type"`
	TransactionDate string            `json:"transaction_date"`
	Amount          float64           `json:"amount"`
	OldValue        float64           `json:"old_value,omitempty"`
	NewValue        float64           `json:"new_value,omitempty"`
	Description     string            `json:"description,omitempty"`
	GLJournalID     string            `json:"gl_journal_id,omitempty"`
	CreatedBy       string            `json:"created_by"`
	CreatedAt       string            `json:"created_at"`
}

type FixedAssetAllocation struct {
	ID               string  `json:"id"`
	FixedAssetID     string  `json:"fixed_asset_id"`
	DepartmentID     string  `json:"department_id"`
	AllocationPct    float64 `json:"allocation_pct"`
	ExpenseAccountID string  `json:"expense_account_id"`
}

type FixedAssetInventoryPlan struct {
	ID        string `json:"id"`
	CompanyID string `json:"company_id"`
	PlanDate  string `json:"plan_date"`
	Status    string `json:"status"`
	Notes     string `json:"notes,omitempty"`
	CreatedBy string `json:"created_by"`
	CreatedAt string `json:"created_at"`
}

type FixedAssetInventoryResult struct {
	ID               string `json:"id"`
	PlanID           string `json:"plan_id"`
	FixedAssetID     string `json:"fixed_asset_id"`
	ExpectedLocation string `json:"expected_location,omitempty"`
	ActualLocation   string `json:"actual_location,omitempty"`
	ExpectedStatus   string `json:"expected_status,omitempty"`
	ActualStatus     string `json:"actual_status,omitempty"`
	Discrepancy      string `json:"discrepancy"`
	Notes            string `json:"notes,omitempty"`
}

type FACategoryFilter struct {
	CompanyID string
	ParentID  *string
	Level     *int
}

type FAListFilter struct {
	CompanyID    string
	Status       *FixedAssetStatus
	CategoryID   *string
	DepartmentID *string
	Keyword      string
	Limit        int
	Offset       int
}

type FARunDepreciationInput struct {
	CompanyID string
	Year      int
	Month     int
	PeriodID  string
	CreatedBy string
}

type FADisposalInput struct {
	FixedAssetID  string
	DisposalType  DisposalType
	DisposalDate  string
	Proceeds      float64
	Description   string
	CustomerID    string
	InvoiceID     string
	CreatedBy     string
}

type FASaleInput struct {
	FixedAssetID string
	SaleDate     string
	Proceeds     float64
	CustomerID   string
	Description  string
	CreatedBy    string
}

type FAAdjustmentInput struct {
	FixedAssetID         string
	NewOriginalCost      *float64
	NewUsefulLifeMonths  *int
	NewDepreciationMethod *DepreciationMethod
	NewResidualValue     *float64
	AdjustmentDate       string
	Reason               string
	CreatedBy            string
}

type FARevaluationInput struct {
	FixedAssetID   string
	FairValue      float64
	RevaluationDate string
	CreatedBy      string
}

type FAImpairmentInput struct {
	FixedAssetID  string
	ImpairmentAmount float64
	ImpairmentDate   string
	Reason        string
	CreatedBy     string
}

type FATransferInput struct {
	FixedAssetID   string
	DepartmentID   string
	EffectiveDate  string
	CreatedBy      string
}

type FAAllocationInput struct {
	FixedAssetID string
	Allocations  []FAAllocationItem
}

type FAAllocationItem struct {
	DepartmentID     string  `json:"department_id"`
	AllocationPct    float64 `json:"allocation_pct"`
	ExpenseAccountID string  `json:"expense_account_id"`
}

type FACIPTransferInput struct {
	FixedAssetID string
	CIPAccountID string
	TotalCost    float64
	TransferDate string
	CreatedBy    string
}

type FASuspendInput struct {
	FixedAssetID string
	SuspendDate  string
	Reason       string
	CreatedBy    string
}

type FAResumeInput struct {
	FixedAssetID string
	ResumeDate   string
	CreatedBy    string
}
