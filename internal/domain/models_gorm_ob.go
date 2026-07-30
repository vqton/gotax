package domain

import "time"

type OpeningBalanceGORM struct {
	ID               string     `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID        string     `gorm:"column:company_id;not null;size:36;index:idx_ob_company_period,unique" json:"companyId"`
	PeriodID         string     `gorm:"column:period_id;not null;size:36;index:idx_ob_company_period,unique" json:"periodId"`
	FiscalYearID     *string    `gorm:"column:fiscal_year_id;size:36" json:"fiscalYearId"`
	AccountCode      string     `gorm:"column:account_code;not null;size:20;index:idx_ob_company_period,unique" json:"accountCode"`
	CurrencyCode     string     `gorm:"column:currency_code;not null;size:3" json:"currencyCode"`
	OriginalAmount   float64    `gorm:"column:original_amount;not null" json:"originalAmount"`
	DebitAmount      float64    `gorm:"column:debit_amount;not null" json:"debitAmount"`
	CreditAmount     float64    `gorm:"column:credit_amount;not null" json:"creditAmount"`
	ExchangeRate     float64    `gorm:"column:exchange_rate;default:1" json:"exchangeRate"`
	Status           string     `gorm:"column:status;not null;size:20;default:'DRAFT';index" json:"status"`
	SourceType       string     `gorm:"column:source_type;not null;size:30" json:"sourceType"`
	BatchID          *string    `gorm:"column:batch_id;size:36" json:"batchId"`
	Reason           *string    `gorm:"column:reason;type:text" json:"reason"`
	ApprovedBy       *string    `gorm:"column:approved_by;size:36" json:"approvedBy"`
	ApprovedAt       *time.Time `gorm:"column:approved_at" json:"approvedAt"`
	CorrectedBy      *string    `gorm:"column:corrected_by;size:36" json:"correctedBy"`
	CorrectionOf     *string    `gorm:"column:correction_of;size:36" json:"correctionOf"`
	CorrectionReason *string    `gorm:"column:correction_reason;type:text" json:"correctionReason"`
	CreatedBy        string     `gorm:"column:created_by;not null;size:36" json:"createdBy"`
	CreatedAt        time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt        time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	Details          []OpeningBalanceDetailGORM `gorm:"foreignKey:OpeningBalanceID;constraint:OnDelete:CASCADE" json:"details,omitempty"`
}

func (OpeningBalanceGORM) TableName() string { return "opening_balances" }

type OpeningBalanceDetailGORM struct {
	ID                string     `gorm:"column:id;primaryKey;size:36" json:"id"`
	OpeningBalanceID  string     `gorm:"column:opening_balance_id;not null;size:36;index" json:"openingBalanceId"`
	EntityType        string     `gorm:"column:entity_type;not null;size:30" json:"entityType"`
	EntityID          string     `gorm:"column:entity_id;not null;size:36" json:"entityId"`
	EntityName        *string    `gorm:"column:entity_name;size:255" json:"entityName"`
	DebitAmount       float64    `gorm:"column:debit_amount;not null;default:0" json:"debitAmount"`
	CreditAmount      float64    `gorm:"column:credit_amount;not null;default:0" json:"creditAmount"`
	Quantity          *float64   `gorm:"column:quantity" json:"quantity"`
	UnitPrice         *float64   `gorm:"column:unit_price" json:"unitPrice"`
	OriginalCost      *float64   `gorm:"column:original_cost" json:"originalCost"`
	AccDepreciation   *float64   `gorm:"column:acc_depreciation" json:"accDepreciation"`
	CounterpartAccount *string  `gorm:"column:counterpart_account;size:20" json:"counterpartAccount"`
	Note              *string    `gorm:"column:note;type:text" json:"note"`
	CreatedAt         time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (OpeningBalanceDetailGORM) TableName() string { return "opening_balance_details" }

type CarryForwardLogGORM struct {
	ID               string     `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID        string     `gorm:"column:company_id;not null;size:36;index" json:"companyId"`
	FromPeriodID     string     `gorm:"column:from_period_id;not null;size:36" json:"fromPeriodId"`
	ToPeriodID       string     `gorm:"column:to_period_id;not null;size:36" json:"toPeriodId"`
	FromFiscalYear   int        `gorm:"column:from_fiscal_year;not null" json:"fromFiscalYear"`
	ToFiscalYear     int        `gorm:"column:to_fiscal_year;not null" json:"toFiscalYear"`
	AccountCount     int        `gorm:"column:account_count;not null;default:0" json:"accountCount"`
	TotalDebit       float64    `gorm:"column:total_debit;not null;default:0" json:"totalDebit"`
	TotalCredit      float64    `gorm:"column:total_credit;not null;default:0" json:"totalCredit"`
	ClosingEntryIDs  *string    `gorm:"column:closing_entry_ids;type:text" json:"closingEntryIds"`
	Status           string     `gorm:"column:status;not null;size:20;index" json:"status"`
	ExecutedBy       string     `gorm:"column:executed_by;not null;size:36" json:"executedBy"`
	ExecutedAt       time.Time  `gorm:"column:executed_at;autoCreateTime" json:"executedAt"`
}

func (CarryForwardLogGORM) TableName() string { return "carry_forward_logs" }

type Circular99MappingGORM struct {
	ID                  string     `gorm:"column:id;primaryKey;size:36" json:"id"`
	OldAccountCode      string     `gorm:"column:old_account_code;not null;size:20;uniqueIndex" json:"oldAccountCode"`
	NewAccountCode      string     `gorm:"column:new_account_code;not null;size:20" json:"newAccountCode"`
	MappingType         string     `gorm:"column:mapping_type;not null;size:50" json:"mappingType"`
	SplitRatio          *float64   `gorm:"column:split_ratio" json:"splitRatio"`
	CounterpartAccount  *string    `gorm:"column:counterpart_account;size:20" json:"counterpartAccount"`
	EffectiveDate       time.Time  `gorm:"column:effective_date;not null;type:date" json:"effectiveDate"`
	Note                *string    `gorm:"column:note;type:text" json:"note"`
	IsActive            bool       `gorm:"column:is_active;not null;default:true;index" json:"isActive"`
	CreatedAt           time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt           time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (Circular99MappingGORM) TableName() string { return "circular99_mappings" }

type BalanceMigrationGORM struct {
	ID                string     `gorm:"column:id;primaryKey;size:36" json:"id"`
	CompanyID         string     `gorm:"column:company_id;not null;size:36;index" json:"companyId"`
	FromRegime        string     `gorm:"column:from_regime;not null;size:50" json:"fromRegime"`
	ToRegime          string     `gorm:"column:to_regime;not null;size:50" json:"toRegime"`
	ExecutionDate     time.Time  `gorm:"column:execution_date;not null;type:date" json:"executionDate"`
	Status            string     `gorm:"column:status;not null;size:20;index" json:"status"`
	SourceBalanceID   *string    `gorm:"column:source_balance_id;size:36" json:"sourceBalanceId"`
	TargetBalanceID   *string    `gorm:"column:target_balance_id;size:36" json:"targetBalanceId"`
	JournalEntryID    *string    `gorm:"column:journal_entry_id;size:36" json:"journalEntryId"`
	Summary           *string    `gorm:"column:summary;type:text" json:"summary"`
	ExecutedBy        string     `gorm:"column:executed_by;not null;size:36" json:"executedBy"`
	CreatedAt         time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	ExecutedAt        *time.Time `gorm:"column:executed_at" json:"executedAt"`
}

func (BalanceMigrationGORM) TableName() string { return "balance_migrations" }
