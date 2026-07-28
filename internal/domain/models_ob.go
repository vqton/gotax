package domain

import "time"

type OpeningBalanceStatus string
const (
	OBStatusDraft     OpeningBalanceStatus = "DRAFT"
	OBStatusPending   OpeningBalanceStatus = "PENDING"
	OBStatusApproved  OpeningBalanceStatus = "APPROVED"
	OBStatusRejected  OpeningBalanceStatus = "REJECTED"
	OBStatusCorrected OpeningBalanceStatus = "CORRECTED"
)
func (s OpeningBalanceStatus) Valid() bool {
	switch s {
	case OBStatusDraft, OBStatusPending, OBStatusApproved, OBStatusRejected, OBStatusCorrected:
		return true
	}
	return false
}
func (s OpeningBalanceStatus) IsTerminal() bool { return s == OBStatusApproved || s == OBStatusCorrected }
func (s OpeningBalanceStatus) CanSubmit() bool  { return s == OBStatusDraft || s == OBStatusRejected }
func (s OpeningBalanceStatus) CanEdit() bool    { return s == OBStatusDraft || s == OBStatusRejected }

type DetailEntityType string
const (
	DetailCustomer      DetailEntityType = "CUSTOMER"
	DetailSupplier      DetailEntityType = "SUPPLIER"
	DetailEmployee      DetailEntityType = "EMPLOYEE"
	DetailBank          DetailEntityType = "BANK_ACCOUNT"
	DetailProject       DetailEntityType = "PROJECT"
	DetailContract      DetailEntityType = "CONTRACT"
	DetailDepartment    DetailEntityType = "DEPARTMENT"
	DetailFixedAsset    DetailEntityType = "FIXED_ASSET"
	DetailInventoryItem DetailEntityType = "INVENTORY_ITEM"
)
func (e DetailEntityType) Valid() bool {
	switch e {
	case DetailCustomer, DetailSupplier, DetailEmployee, DetailBank,
		DetailProject, DetailContract, DetailDepartment, DetailFixedAsset, DetailInventoryItem:
		return true
	}
	return false
}

type OBListFilter struct {
	CompanyID  string              `json:"company_id"`
	PeriodID   string              `json:"period_id"`
	Status     OpeningBalanceStatus `json:"status,omitempty"`
	AccountCode string             `json:"account_code,omitempty"`
	Currency   string              `json:"currency,omitempty"`
	SourceType string              `json:"source_type,omitempty"`
}

type OpeningBalance struct {
	ID               string               `json:"id"`
	CompanyID        string               `json:"company_id"`
	PeriodID         string               `json:"period_id"`
	FiscalYearID     string               `json:"fiscal_year_id,omitempty"`
	AccountCode      string               `json:"account_code"`
	CurrencyCode     string               `json:"currency_code"`
	OriginalAmount   float64              `json:"original_amount,omitempty"`
	DebitAmount      float64              `json:"debit_amount"`
	CreditAmount     float64              `json:"credit_amount"`
	ExchangeRate     float64              `json:"exchange_rate,omitempty"`
	Status           OpeningBalanceStatus  `json:"status"`
	SourceType       string               `json:"source_type"`
	BatchID          string               `json:"batch_id,omitempty"`
	Reason           string               `json:"reason,omitempty"`
	ApprovedBy       string               `json:"approved_by,omitempty"`
	ApprovedAt       *time.Time           `json:"approved_at,omitempty"`
	CorrectedBy      string               `json:"corrected_by,omitempty"`
	CorrectionOf     string               `json:"correction_of,omitempty"`
	CorrectionReason string               `json:"correction_reason,omitempty"`
	CreatedBy        string               `json:"created_by"`
	CreatedAt        time.Time            `json:"created_at"`
	UpdatedAt        time.Time            `json:"updated_at"`
}

func (ob *OpeningBalance) Validate() error {
	if ob.CompanyID == "" { return ErrCompanyIDRequired }
	if ob.PeriodID == "" { return ErrPeriodNotFound }
	if ob.AccountCode == "" { return ErrAccountCodeRequired }
	if ob.CurrencyCode == "" { ob.CurrencyCode = "VND" }
	if ob.DebitAmount < 0 || ob.CreditAmount < 0 { return ErrJournalLineInvalidAmount }
	if ob.DebitAmount == 0 && ob.CreditAmount == 0 { return ErrAmountRequired }
	if ob.DebitAmount > 0 && ob.CreditAmount > 0 { return ErrBothDebitAndCredit }
	if ob.CurrencyCode != "VND" && ob.ExchangeRate <= 0 { return ErrInvalidExchangeRate }
	if ob.Status == "" { ob.Status = OBStatusDraft }
	if !ob.Status.Valid() { return ErrDeclarationStatusInvalid }
	switch ob.SourceType {
	case "MANUAL", "CARRY_FORWARD", "IMPORT", "MIGRATION", "":
	default: return ErrMappingTypeInvalid
	}
	if ob.SourceType == "" { ob.SourceType = "MANUAL" }
	return nil
}

func (ob *OpeningBalance) CanSubmit() bool  { return ob.Status.CanSubmit() }
func (ob *OpeningBalance) CanEdit() bool    { return ob.Status.CanEdit() }

type OpeningBalanceDetail struct {
	ID               string           `json:"id"`
	OpeningBalanceID string           `json:"opening_balance_id,omitempty"`
	EntityType       DetailEntityType `json:"entity_type"`
	EntityID         string           `json:"entity_id"`
	EntityName       string           `json:"entity_name,omitempty"`
	DebitAmount      float64          `json:"debit_amount"`
	CreditAmount     float64          `json:"credit_amount"`
	Quantity         float64          `json:"quantity,omitempty"`
	UnitPrice        float64          `json:"unit_price,omitempty"`
	OriginalCost     float64          `json:"original_cost,omitempty"`
	AccDepreciation  float64          `json:"acc_depreciation,omitempty"`
	CounterpartAccount string        `json:"counterpart_account,omitempty"`
	Note             string           `json:"note,omitempty"`
	CreatedAt        time.Time        `json:"created_at"`
}

func (d *OpeningBalanceDetail) Validate() error {
	if !d.EntityType.Valid() { return ErrEntityTypeRequired }
	if d.EntityID == "" { return ErrEntityIDRequired }
	if d.DebitAmount < 0 || d.CreditAmount < 0 { return ErrJournalLineInvalidAmount }
	if d.DebitAmount == 0 && d.CreditAmount == 0 { return ErrAmountRequired }
	if d.DebitAmount > 0 && d.CreditAmount > 0 { return ErrBothDebitAndCredit }
	return nil
}

type CarryForwardLog struct {
	ID              string    `json:"id"`
	CompanyID       string    `json:"company_id"`
	FromPeriodID    string    `json:"from_period_id"`
	ToPeriodID      string    `json:"to_period_id"`
	FromFiscalYear  int       `json:"from_fiscal_year"`
	ToFiscalYear    int       `json:"to_fiscal_year"`
	AccountCount    int       `json:"account_count"`
	TotalDebit      float64   `json:"total_debit"`
	TotalCredit     float64   `json:"total_credit"`
	ClosingEntryIDs []string  `json:"closing_entry_ids,omitempty"`
	Status          string    `json:"status"`
	ExecutedBy      string    `json:"executed_by"`
	ExecutedAt      time.Time `json:"executed_at"`
}

type Circular99Mapping struct {
	ID                  string  `json:"id"`
	OldAccountCode      string  `json:"old_account_code"`
	NewAccountCode      string  `json:"new_account_code"`
	MappingType         string  `json:"mapping_type"`
	SplitRatio          float64 `json:"split_ratio,omitempty"`
	CounterpartAccount  string  `json:"counterpart_account,omitempty"`
	EffectiveDate       string  `json:"effective_date"`
	Note                string  `json:"note,omitempty"`
	IsActive            bool    `json:"is_active"`
}

type BalanceMigration struct {
	ID              string    `json:"id"`
	CompanyID       string    `json:"company_id"`
	FromRegime      string    `json:"from_regime"`
	ToRegime        string    `json:"to_regime"`
	ExecutionDate   string    `json:"execution_date"`
	Status          string    `json:"status"`
	SourceBalanceID string    `json:"source_balance_id,omitempty"`
	TargetBalanceID string    `json:"target_balance_id,omitempty"`
	JournalEntryID  string    `json:"journal_entry_id,omitempty"`
	Summary         string    `json:"summary,omitempty"`
	ExecutedBy      string    `json:"executed_by"`
	CreatedAt       time.Time  `json:"created_at"`
	ExecutedAt      *time.Time `json:"executed_at,omitempty"`
}
