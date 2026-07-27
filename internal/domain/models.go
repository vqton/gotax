package domain

import (
	"math"
	"time"
	"unicode"
)

type AccountType string
const (AccountTypeAsset AccountType="ASSET"; AccountTypeLiability AccountType="LIABILITY"; AccountTypeEquity AccountType="EQUITY"; AccountTypeRevenue AccountType="REVENUE"; AccountTypeExpense AccountType="EXPENSE")
func (at AccountType) NormalBalance() NormalBalance {
	switch at {
	case AccountTypeAsset, AccountTypeExpense: return NormalBalanceDebit
	case AccountTypeLiability, AccountTypeEquity, AccountTypeRevenue: return NormalBalanceCredit
	default: return NormalBalanceDebit
	}
}

type NormalBalance string
const (NormalBalanceDebit NormalBalance="DEBIT"; NormalBalanceCredit NormalBalance="CREDIT")

type VoucherType string
const (VoucherTypeReceipt VoucherType="THU"; VoucherTypePayment VoucherType="CHI"; VoucherTypeSale VoucherType="BAN"; VoucherTypePurchase VoucherType="MUA"; VoucherTypeOther VoucherType="KHAC"; VoucherTypeClosing VoucherType="KC")

type DetailBy string
const (DetailByNone DetailBy=""; DetailByObject DetailBy="OBJECT"; DetailByProject DetailBy="PROJECT"; DetailByContract DetailBy="CONTRACT"; DetailByCostItem DetailBy="COST_ITEM"; DetailByDepartment DetailBy="DEPARTMENT")

type UserRole string
const (UserRoleAdmin UserRole="admin"; UserRoleChiefAccountant UserRole="chief_accountant"; UserRoleAccountant UserRole="accountant"; UserRoleViewer UserRole="viewer")

type AuditAction string
const (AuditActionCreate AuditAction="CREATE"; AuditActionUpdate AuditAction="UPDATE"; AuditActionDelete AuditAction="DELETE"; AuditActionApprove AuditAction="APPROVE"; AuditActionPost AuditAction="POST"; AuditActionCancel AuditAction="CANCEL"; AuditActionClose AuditAction="CLOSE"; AuditActionReview AuditAction="REVIEW"; AuditActionLogin AuditAction="LOGIN")

type AccountStatus string
const (AccountStatusActive AccountStatus="ACTIVE"; AccountStatusFrozen AccountStatus="FROZEN"; AccountStatusInactive AccountStatus="INACTIVE")

type ApprovalStatus string
const (ApprovalPending ApprovalStatus="PENDING"; ApprovalApproved ApprovalStatus="APPROVED"; ApprovalRejected ApprovalStatus="REJECTED"; ApprovalCancelled ApprovalStatus="CANCELLED"; ApprovalExpired ApprovalStatus="EXPIRED")
func (s ApprovalStatus) IsTerminal() bool { return s==ApprovalApproved||s==ApprovalRejected||s==ApprovalCancelled||s==ApprovalExpired }

type Account struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Name2 string `json:"name2,omitempty"`
	Type AccountType `json:"type"`
	ParentCode string `json:"parent_code,omitempty"`
	IsActive bool `json:"is_active"`
	Status AccountStatus `json:"status,omitempty"`
	FreezeReason string `json:"freeze_reason,omitempty"`
	IsForeign bool `json:"is_foreign"`
	DetailBy DetailBy `json:"detail_by,omitempty"`
	IsParent bool `json:"is_parent"`
	ArrearsDays int `json:"arrears_days,omitempty"`
	Note string `json:"note,omitempty"`
}

func (a *Account) Validate() error {
	if a.Code=="" { return ErrAccountCodeRequired }
	if len(a.Code)<3 { return ErrAccountCodeInvalid }
	for _,c:=range a.Code{if !unicode.IsDigit(c){return ErrAccountCodeInvalid}}
	if a.Name=="" { return ErrAccountNameRequired }
	switch a.Type{
	case AccountTypeAsset,AccountTypeLiability,AccountTypeEquity,AccountTypeRevenue,AccountTypeExpense:
	default: return ErrAccountInvalidType
	}
	if a.Status=="" { a.Status=AccountStatusActive }
	switch a.Status{
	case AccountStatusActive,AccountStatusFrozen,AccountStatusInactive:
	default: return ErrAccountStatusInvalid
	}
	if a.Status==AccountStatusFrozen&&a.FreezeReason=="" { return ErrFreezeReasonRequired }
	return nil
}
func (a *Account) Freeze(reason string) error {
	if a.Status==AccountStatusFrozen { return ErrAccountAlreadyFrozen }
	if reason=="" { return ErrFreezeReasonRequired }
	a.Status=AccountStatusFrozen; a.FreezeReason=reason; return nil
}
func (a *Account) Unfreeze(reason string) error {
	if a.Status!=AccountStatusFrozen { return ErrAccountNotFrozen }
	a.Status=AccountStatusActive; a.FreezeReason=""; return nil
}
func (a *Account) CanPost() error {
	if a.Status==AccountStatusFrozen { return ErrAccountFrozen }
	if !a.IsActive||a.Status==AccountStatusInactive { return ErrAccountInactive }
	return nil
}

type JournalLine struct {
	ID string `json:"id,omitempty"`
	EntryID string `json:"entry_id,omitempty"`
	LineNumber int `json:"line_number"`
	AccountCode string `json:"account_code"`
	DebitAmount float64 `json:"debit_amount"`
	CreditAmount float64 `json:"credit_amount"`
	Description string `json:"description,omitempty"`
	CurrencyCode string `json:"currency_code,omitempty"`
	ForeignAmount float64 `json:"foreign_amount,omitempty"`
	ExchangeRate float64 `json:"exchange_rate,omitempty"`
	ObjectID string `json:"object_id,omitempty"`
	ProjectID string `json:"project_id,omitempty"`
	ContractID string `json:"contract_id,omitempty"`
	CostItemID string `json:"cost_item_id,omitempty"`
	DepartmentID string `json:"department_id,omitempty"`
}

type JournalEntryStatus string
const (JournalEntryDraft JournalEntryStatus="DRAFT"; JournalEntryReviewing JournalEntryStatus="REVIEWING"; JournalEntryApproved JournalEntryStatus="APPROVED"; JournalEntryPosted JournalEntryStatus="POSTED"; JournalEntryCancelled JournalEntryStatus="CANCELLED")

type JournalEntry struct {
	ID string `json:"id,omitempty"`
	CompanyID string `json:"company_id,omitempty"`
	EntryNumber string `json:"entry_number"`
	VoucherType VoucherType `json:"voucher_type,omitempty"`
	EntryDate time.Time `json:"entry_date"`
	AccountingDate time.Time `json:"accounting_date,omitempty"`
	PeriodID string `json:"period_id,omitempty"`
	Description string `json:"description"`
	Status JournalEntryStatus `json:"status"`
	CurrencyCode string `json:"currency_code,omitempty"`
	ExchangeRate float64 `json:"exchange_rate,omitempty"`
	Lines []JournalLine `json:"lines"`
	CreatedBy string `json:"created_by,omitempty"`
	ReviewedBy string `json:"reviewed_by,omitempty"`
	ApprovedBy string `json:"approved_by,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	PostedAt *time.Time `json:"posted_at,omitempty"`
	ApprovedAt *time.Time `json:"approved_at,omitempty"`
}

func (je *JournalEntry) TotalDebit() float64 {
	var t float64; for _,l:=range je.Lines{t+=l.DebitAmount}; return t
}
func (je *JournalEntry) TotalCredit() float64 {
	var t float64; for _,l:=range je.Lines{t+=l.CreditAmount}; return t
}
func (je *JournalEntry) HasDebit() bool { return je.TotalDebit()>0 }
func (je *JournalEntry) HasCredit() bool { return je.TotalCredit()>0 }
func (je *JournalEntry) Validate() error {
	if je.EntryDate.IsZero() { return ErrJournalEntryInvalidDate }
	if je.Description=="" { return ErrJournalEntryNoDescription }
	if je.CurrencyCode=="" { je.CurrencyCode="VND" }
	if je.CurrencyCode!="VND"&&je.ExchangeRate<=0 { return ErrInvalidExchangeRate }
	if len(je.Lines)==0 { return ErrJournalEntryNoLines }
	for _,l:=range je.Lines{
		if l.DebitAmount<0||l.CreditAmount<0 { return ErrJournalLineInvalidAmount }
		if l.DebitAmount==0&&l.CreditAmount==0 { return ErrJournalLineInvalidAmount }
	}
	td:=je.TotalDebit(); tc:=je.TotalCredit()
	if math.Abs(td-tc)>0.001 { return ErrJournalEntryUnbalanced }
	return nil
}

type PeriodStatus string
const (PeriodOpen PeriodStatus="OPEN"; PeriodClosing PeriodStatus="CLOSING"; PeriodClosed PeriodStatus="CLOSED"; PeriodLocked PeriodStatus="LOCKED")

type Period struct {
	ID string `json:"id"`
	Year int `json:"year"`
	Month int `json:"month"`
	StartDate time.Time `json:"start_date"`
	EndDate time.Time `json:"end_date"`
	Status PeriodStatus `json:"status"`
}
func (p *Period) Validate() error {
	if p.Year<2000||p.Year>2100 { return ErrPeriodYearOutOfRange }
	if p.Month<1||p.Month>12 { return ErrPeriodMonthInvalid }
	if p.StartDate.IsZero()||p.EndDate.IsZero() { return ErrPeriodDateRequired }
	if p.EndDate.Before(p.StartDate) { return ErrPeriodEndBeforeStart }
	switch p.Status{
	case PeriodOpen,PeriodClosing,PeriodClosed,PeriodLocked:
	default: return ErrPeriodStatusInvalid
	}
	return nil
}

type AccountBalance struct {
	AccountCode string `json:"account_code"`
	AccountType AccountType `json:"account_type"`
	PeriodID string `json:"period_id"`
	OpenBalanceDebit float64 `json:"open_balance_debit"`
	OpenBalanceCredit float64 `json:"open_balance_credit"`
	PeriodDebit float64 `json:"period_debit"`
	PeriodCredit float64 `json:"period_credit"`
	TotalDebit float64 `json:"total_debit"`
	TotalCredit float64 `json:"total_credit"`
	ClosingBalance float64 `json:"closing_balance"`
}
func (b *AccountBalance) Calculate() {
	b.TotalDebit=b.OpenBalanceDebit+b.PeriodDebit; b.TotalCredit=b.OpenBalanceCredit+b.PeriodCredit
	b.ClosingBalance=b.TotalDebit-b.TotalCredit
}

type User struct {
	ID string `json:"id"`
	Username string `json:"username"`
	PasswordHash string `json:"-"`
	FullName string `json:"full_name"`
	Email string `json:"email,omitempty"`
	Role UserRole `json:"role"`
	IsActive bool `json:"is_active"`
	FailedAttempts int `json:"-"`
	LockedUntil *time.Time `json:"-"`
	PasswordChangedAt *time.Time `json:"password_changed_at"`
	PasswordHistory []string `json:"-"`
	TOTPSecret string `json:"-"`
	TOTPEnabled bool `json:"totp_enabled"`
	BackupCodes []string `json:"-"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	LastLoginIP string `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
func (u *User) Validate() error {
	if u.Username=="" { return ErrUsernameRequired }
	switch u.Role{case UserRoleAdmin,UserRoleChiefAccountant,UserRoleAccountant,UserRoleViewer:default: return ErrInvalidUserRole}
	return nil
}
func (u *User) IsLocked() bool { return u.LockedUntil!=nil&&time.Now().Before(*u.LockedUntil) }
func (u *User) IsPasswordExpired() bool { return u.PasswordChangedAt==nil||time.Now().After(u.PasswordChangedAt.Add(PasswordMaxAge)) }

type RefreshToken struct {
	ID string `json:"id"`
	UserID string `json:"user_id"`
	TokenHash string `json:"-"`
	DeviceInfo string `json:"device_info,omitempty"`
	IPAddress string `json:"ip_address,omitempty"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
}
type PasswordResetToken struct {
	ID string `json:"id"`
	UserID string `json:"user_id"`
	TokenHash string `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	UsedAt *time.Time `json:"used_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
type TOTPSetup struct {
	Secret string `json:"secret"`
	QRCodeURL string `json:"qr_code_url"`
	BackupCodes []string `json:"backup_codes"`
}
type Session struct {
	ID string `json:"id"`
	Device string `json:"device"`
	IP string `json:"ip"`
	CreatedAt time.Time `json:"created_at"`
	LastActivity time.Time `json:"last_activity"`
	IsCurrent bool `json:"is_current"`
}
type TokenPair struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn int `json:"expires_in"`
	TokenType string `json:"token_type"`
}
type AuthResult struct {
	AccessToken string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn int `json:"expires_in,omitempty"`
	TokenType string `json:"token_type,omitempty"`
	User *User `json:"user,omitempty"`
	Requires2FA bool `json:"requires_2fa,omitempty"`
	TempToken string `json:"temp_token,omitempty"`
	PasswordExpired bool `json:"password_expired,omitempty"`
}

type AuditEntry struct {
	ID string `json:"id"`
	UserID string `json:"user_id,omitempty"`
	Username string `json:"username"`
	IPAddress string `json:"ip_address,omitempty"`
	Action AuditAction `json:"action"`
	EntityType string `json:"entity_type"`
	EntityID string `json:"entity_id,omitempty"`
	OldValue string `json:"old_value,omitempty"`
	NewValue string `json:"new_value,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type ExchangeRate struct {
	ID string `json:"id,omitempty"`
	CurrencyCode string `json:"currency_code"`
	RateDate time.Time `json:"rate_date"`
	BuyRate float64 `json:"buy_rate,omitempty"`
	SellRate float64 `json:"sell_rate,omitempty"`
	AverageRate float64 `json:"average_rate"`
	Source string `json:"source,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
func (r *ExchangeRate) Validate() error {
	if len(r.CurrencyCode)!=3 { return ErrInvalidCurrencyCode }
	if r.AverageRate<=0 { return ErrInvalidExchangeRate }
	return nil
}

type ClosingTemplate struct {
	ID string `json:"id,omitempty"`
	Name string `json:"name"`
	Description string `json:"description,omitempty"`
	SequenceOrder int `json:"sequence_order"`
	IsActive bool `json:"is_active"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
type ClosingTemplateLine struct {
	ID string `json:"id,omitempty"`
	TemplateID string `json:"template_id,omitempty"`
	LineNumber int `json:"line_number"`
	DebitAccount string `json:"debit_account"`
	CreditAccount string `json:"credit_account"`
	Formula string `json:"formula"`
	Description string `json:"description,omitempty"`
}

// COA models
type ApprovalRequest struct {
	ID string `json:"id"`
	TenantID string `json:"tenant_id,omitempty"`
	EntityType string `json:"entity_type"`
	EntityID string `json:"entity_id"`
	RequestType string `json:"request_type"`
	OldValue string `json:"old_value,omitempty"`
	NewValue string `json:"new_value"`
	Reason string `json:"reason"`
	Status ApprovalStatus `json:"status"`
	RequestedBy string `json:"requested_by"`
	ReviewedBy string `json:"reviewed_by,omitempty"`
	ReviewNote string `json:"review_note,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	ReviewedAt *time.Time `json:"reviewed_at,omitempty"`
}
func (r *ApprovalRequest) Validate() error {
	if r.EntityType=="" { return ErrEntityTypeRequired }
	if r.EntityID=="" { return ErrEntityIDRequired }
	if r.RequestType=="" { return ErrRequestTypeRequired }
	if r.Reason=="" { return ErrApprovalReasonRequired }
	if r.RequestedBy=="" { return ErrRequesterRequired }
	if r.Status=="" { r.Status=ApprovalPending }
	return nil
}

type AccountVersion struct {
	ID string `json:"id"`
	VersionNumber string `json:"version_number"`
	Snapshot string `json:"snapshot"`
	ChangeSummary string `json:"change_summary,omitempty"`
	ChangeReason string `json:"change_reason,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
type AccountMapping struct {
	ID string `json:"id"`
	SourceRegime string `json:"source_regime"`
	TargetRegime string `json:"target_regime"`
	OldCode string `json:"old_code"`
	NewCode string `json:"new_code"`
	MappingType string `json:"mapping_type"`
	SplitRatio float64 `json:"split_ratio,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
func (m *AccountMapping) Validate() error {
	if m.SourceRegime=="" { return ErrSourceRegimeRequired }
	if m.TargetRegime=="" { return ErrTargetRegimeRequired }
	if m.OldCode=="" { return ErrOldCodeRequired }
	if m.NewCode=="" { return ErrNewCodeRequired }
	switch m.MappingType{case "DIRECT","MERGE","SPLIT","CUSTOM":default: return ErrMappingTypeInvalid}
	if m.MappingType=="SPLIT"&&m.SplitRatio<=0 { return ErrSplitRatioRequired }
	return nil
}
type AccountAnalysis struct {
	AccountCode string `json:"account_code"`
	CostCenterID string `json:"cost_center_id,omitempty"`
	ProfitCenterID string `json:"profit_center_id,omitempty"`
	DepartmentID string `json:"department_id,omitempty"`
	ProjectID string `json:"project_id,omitempty"`
	CustomDim1 string `json:"custom_dim1,omitempty"`
	CustomDim2 string `json:"custom_dim2,omitempty"`
}
func (a *AccountAnalysis) Validate() error {
	if a.AccountCode=="" { return ErrAccountCodeRequired }
	return nil
}
type IFRSMapping struct {
	ID string `json:"id"`
	VASCode string `json:"vas_code"`
	IFRSCode string `json:"ifrs_code"`
	IFRSName string `json:"ifrs_name,omitempty"`
	ReclassificationRule string `json:"reclassification_rule,omitempty"`
	AdjustmentType string `json:"adjustment_type,omitempty"`
	IsActive bool `json:"is_active"`
}
func (m *IFRSMapping) Validate() error {
	if m.VASCode=="" { return ErrVASCodeRequired }
	if m.IFRSCode=="" { return ErrIFRSCodeRequired }
	return nil
}
type AccountUsage struct {
	AccountCode string `json:"account_code"`
	EntryCount int `json:"entry_count"`
	TotalDebit float64 `json:"total_debit"`
	TotalCredit float64 `json:"total_credit"`
	LastUsedDate string `json:"last_used_date,omitempty"`
}
type VersionDiff struct {
	VersionFrom string `json:"version_from"`
	VersionTo string `json:"version_to"`
	Added []Account `json:"added,omitempty"`
	Removed []Account `json:"removed,omitempty"`
	Modified []AccountDiff `json:"modified,omitempty"`
}
type AccountDiff struct {
	Code string `json:"code"`
	Old Account `json:"old"`
	New Account `json:"new"`
	Changes map[string]Change `json:"changes"`
}
type Change struct{ OldValue interface{} `json:"old_value"`; NewValue interface{} `json:"new_value"` }

// ─── Company models ──────────────────────────────────────────────

type CompanyStatus string
const (CompanyStatusActive CompanyStatus="ACTIVE"; CompanyStatusSuspended CompanyStatus="SUSPENDED"; CompanyStatusDissolved CompanyStatus="DISSOLVED"; CompanyStatusMerged CompanyStatus="MERGED")

type LegalForm string
const (LegalFormLLC1Member LegalForm="LLC_1MEMBER"; LegalFormLLC2Members LegalForm="LLC_2MEMBERS"; LegalFormJSC LegalForm="JSC"; LegalFormSoleProp LegalForm="SOLE_PROP"; LegalFormPartnership LegalForm="PARTNERSHIP"; LegalFormFO LegalForm="FO"; LegalFormRO LegalForm="RO")

type CompanyType string
const (CompanyTypeManufacturing CompanyType="MANUFACTURING"; CompanyTypeTrading CompanyType="TRADING"; CompanyTypeService CompanyType="SERVICE"; CompanyTypeConstruction CompanyType="CONSTRUCTION"; CompanyTypeAgriculture CompanyType="AGRICULTURE"; CompanyTypeOther CompanyType="OTHER")

type CompanySize string
const (CompanySizeMicro CompanySize="MICRO"; CompanySizeSmall CompanySize="SMALL"; CompanySizeMedium CompanySize="MEDIUM"; CompanySizeLarge CompanySize="LARGE")

type AccountingRegime string
const (RegimeTT99 AccountingRegime="TT99"; RegimeTT133 AccountingRegime="TT133"; RegimeTT58 AccountingRegime="TT58")

type BranchType string
const (BranchTypeBranch BranchType="BRANCH"; BranchTypeRepOffice BranchType="REP_OFFICE"; BranchTypeDependentUnit BranchType="DEPENDENT_UNIT")

type PeriodStatusV2 string
const (PeriodV2Open PeriodStatusV2="OPEN"; PeriodV2Closed PeriodStatusV2="CLOSED"; PeriodV2PermanentlyClosed PeriodStatusV2="PERMANENTLY_CLOSED"; PeriodV2Future PeriodStatusV2="FUTURE")

type EmployeeStatus string
const (EmployeeActive EmployeeStatus="ACTIVE"; EmployeeLeave EmployeeStatus="LEAVE"; EmployeeTerminated EmployeeStatus="TERMINATED")

type BankAccountStatus string
const (BankAccountActive BankAccountStatus="ACTIVE"; BankAccountClosed BankAccountStatus="CLOSED")

type EInvoiceStatus string
const (EInvoiceRegistered EInvoiceStatus="REGISTERED"; EInvoiceActive EInvoiceStatus="ACTIVE"; EInvoiceCancelled EInvoiceStatus="CANCELLED")

type SignatureType string
const (SignatureUSBToken SignatureType="USB_TOKEN"; SignatureRemoteHSM SignatureType="REMOTE_HSM")

type SignatureStatus string
const (SignatureActive SignatureStatus="ACTIVE"; SignatureExpired SignatureStatus="EXPIRED"; SignatureRevoked SignatureStatus="REVOKED")

type IntegrationType string
const (IntegrationGDT IntegrationType="GDT"; IntegrationCustoms IntegrationType="CUSTOMS"; IntegrationBHXH IntegrationType="BHXH"; IntegrationDVC IntegrationType="DVC")

type IntegrationStatus string
const (IntegrationConnected IntegrationStatus="CONNECTED"; IntegrationDisconnected IntegrationStatus="DISCONNECTED"; IntegrationError IntegrationStatus="ERROR")

type PeriodTypeV2 string
const (PeriodTypeMonthly PeriodTypeV2="MONTHLY"; PeriodTypeQuarterly PeriodTypeV2="QUARTERLY"; PeriodTypeAnnual PeriodTypeV2="ANNUAL"; PeriodTypePerOccurrence PeriodTypeV2="PER_OCCURRENCE")

type Company struct {
	ID string `json:"id"`
	TenantID string `json:"tenant_id,omitempty"`
	LegalNameVN string `json:"legal_name_vn"`
	LegalNameEN string `json:"legal_name_en,omitempty"`
	ShortName string `json:"short_name,omitempty"`
	LegalForm LegalForm `json:"legal_form"`
	TaxCode string `json:"tax_code"`
	BusinessRegNo string `json:"business_reg_no,omitempty"`
	BusinessRegDate string `json:"business_reg_date,omitempty"`
	BusinessRegPlace string `json:"business_reg_place,omitempty"`
	RegAddress string `json:"reg_address"`
	RegProvince string `json:"reg_province,omitempty"`
	RegDistrict string `json:"reg_district,omitempty"`
	HeadOfficeAddress string `json:"head_office_address,omitempty"`
	HeadOfficeProvince string `json:"head_office_province,omitempty"`
	HeadOfficeDistrict string `json:"head_office_district,omitempty"`
	Phone string `json:"phone,omitempty"`
	Email string `json:"email,omitempty"`
	Website string `json:"website,omitempty"`
	LegalRepName string `json:"legal_rep_name,omitempty"`
	LegalRepTitle string `json:"legal_rep_title,omitempty"`
	LegalRepIDNumber string `json:"legal_rep_id_number,omitempty"`
	ChiefAccountant string `json:"chief_accountant,omitempty"`
	ChiefAccountantEmail string `json:"chief_accountant_email,omitempty"`
	TaxOfficeCode string `json:"tax_office_code,omitempty"`
	TaxOfficeName string `json:"tax_office_name,omitempty"`
	AccountingRegime AccountingRegime `json:"accounting_regime"`
	FiscalYearStartMonth int `json:"fiscal_year_start_month"`
	DefaultCurrency string `json:"default_currency"`
	SecondaryCurrency string `json:"secondary_currency,omitempty"`
	CompanyType CompanyType `json:"company_type,omitempty"`
	CompanySize CompanySize `json:"company_size,omitempty"`
	Status CompanyStatus `json:"status"`
	ParentCompanyID string `json:"parent_company_id,omitempty"`
	LogoURL string `json:"logo_url,omitempty"`
	Settings map[string]any `json:"settings,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Company) Validate() error {
	if c.TaxCode=="" { return ErrCompanyCodeRequired }
	if len(c.TaxCode)!=10&&len(c.TaxCode)!=13 { return ErrInvalidTaxCodeFormat }
	if c.LegalNameVN=="" { return ErrCompanyNameRequired }
	switch c.LegalForm{case LegalFormLLC1Member,LegalFormLLC2Members,LegalFormJSC,LegalFormSoleProp,LegalFormPartnership,LegalFormFO,LegalFormRO:default: return ErrCompanyInvalidLegalForm}
	switch c.AccountingRegime{case RegimeTT99,RegimeTT133,RegimeTT58:default: return ErrCompanyInvalidRegime}
	if c.Status=="" { c.Status=CompanyStatusActive }
	if c.DefaultCurrency=="" { c.DefaultCurrency="VND" }
	if c.FiscalYearStartMonth==0 { c.FiscalYearStartMonth=1 }
	return nil
}

type CompanyBranch struct {
	ID string `json:"id"`
	CompanyID string `json:"company_id"`
	BranchName string `json:"branch_name"`
	BranchTaxCode string `json:"branch_tax_code"`
	BranchType BranchType `json:"branch_type"`
	Address string `json:"address,omitempty"`
	Phone string `json:"phone,omitempty"`
	Email string `json:"email,omitempty"`
	ManagerName string `json:"manager_name,omitempty"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type FiscalYear struct {
	ID string `json:"id"`
	CompanyID string `json:"company_id"`
	Year int `json:"year"`
	StartMonth int `json:"start_month"`
	IsShortYear bool `json:"is_short_year"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
	Status string `json:"status"`
	Periods []PeriodV2 `json:"periods,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
type PeriodV2 struct {
	ID string `json:"id"`
	CompanyID string `json:"company_id"`
	FiscalYearID string `json:"fiscal_year_id"`
	PeriodType PeriodTypeV2 `json:"period_type"`
	PeriodNumber int `json:"period_number"`
	Label string `json:"label"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
	Status PeriodStatusV2 `json:"status"`
	OpenedAt string `json:"opened_at,omitempty"`
	ClosedAt string `json:"closed_at,omitempty"`
	ClosedBy string `json:"closed_by,omitempty"`
	ReopenedCount int `json:"reopened_count,omitempty"`
}
type Department struct {
	ID string `json:"id"`
	CompanyID string `json:"company_id"`
	Code string `json:"code"`
	Name string `json:"name"`
	ParentID string `json:"parent_id,omitempty"`
	ManagerID string `json:"manager_id,omitempty"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type Employee struct {
	ID string `json:"id"`
	CompanyID string `json:"company_id"`
	EmployeeCode string `json:"employee_code"`
	FullName string `json:"full_name"`
	Title string `json:"title,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
	DepartmentID string `json:"department_id,omitempty"`
	PersonalTaxCode string `json:"personal_tax_code,omitempty"`
	SocialInsuranceNo string `json:"social_insurance_no,omitempty"`
	BankAccountNo string `json:"bank_account_no,omitempty"`
	UserID string `json:"user_id,omitempty"`
	Status EmployeeStatus `json:"status"`
	HireDate string `json:"hire_date,omitempty"`
	TerminationDate string `json:"termination_date,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type CompanyBankAccount struct {
	ID string `json:"id"`
	CompanyID string `json:"company_id"`
	BankCode string `json:"bank_code,omitempty"`
	BankName string `json:"bank_name"`
	BranchName string `json:"branch_name,omitempty"`
	AccountNumber string `json:"account_number"`
	AccountHolder string `json:"account_holder"`
	Currency string `json:"currency"`
	IsDefault bool `json:"is_default"`
	IsVerified bool `json:"is_verified"`
	Status BankAccountStatus `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type EInvoicePattern struct {
	ID string `json:"id"`
	CompanyID string `json:"company_id"`
	PatternCode string `json:"pattern_code"`
	Serial string `json:"serial"`
	Form string `json:"form,omitempty"`
	InvoiceType string `json:"invoice_type,omitempty"`
	Status EInvoiceStatus `json:"status"`
	GDTStatus string `json:"gdt_status,omitempty"`
	Description string `json:"description,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type DigitalSignature struct {
	ID string `json:"id"`
	CompanyID string `json:"company_id"`
	SignatureType SignatureType `json:"signature_type"`
	Provider string `json:"provider,omitempty"`
	SerialNumber string `json:"serial_number"`
	OwnerName string `json:"owner_name,omitempty"`
	CertificateSubject string `json:"certificate_subject,omitempty"`
	CertificateIssuer string `json:"certificate_issuer,omitempty"`
	ValidFrom string `json:"valid_from"`
	ValidTo string `json:"valid_to"`
	Status SignatureStatus `json:"status"`
	IsDefault bool `json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type IntegrationProfile struct {
	ID string `json:"id"`
	CompanyID string `json:"company_id"`
	IntegrationType IntegrationType `json:"integration_type"`
	EndpointURL string `json:"endpoint_url,omitempty"`
	Status IntegrationStatus `json:"status"`
	LastConnectedAt string `json:"last_connected_at,omitempty"`
	LastErrorAt string `json:"last_error_at,omitempty"`
	LastErrorMsg string `json:"last_error_msg,omitempty"`
	Config map[string]any `json:"config,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type CompanyContext struct {
	CompanyID string `json:"company_id"`
	LegalNameVN string `json:"legal_name_vn"`
	TaxCode string `json:"tax_code"`
	AccountingRegime string `json:"accounting_regime"`
	DefaultCurrency string `json:"default_currency"`
	FiscalYear int `json:"fiscal_year"`
	CurrentPeriodID string `json:"current_period_id"`
	CurrentPeriodLabel string `json:"current_period_label"`
	Permissions []string `json:"permissions"`
}

// ─── Tax Module Types ─────────────────────────────────────────────

type DeclarationType string
const (
	DeclTypeGTGT01    DeclarationType = "GTGT01"
	DeclTypeGTGT02    DeclarationType = "GTGT02"
	DeclTypeGTGT03    DeclarationType = "GTGT03"
	DeclTypeGTGT04    DeclarationType = "GTGT04"
	DeclTypeGTGT05    DeclarationType = "GTGT05"
	DeclTypeTNDN03    DeclarationType = "TNDN03"
	DeclTypeTNDN04    DeclarationType = "TNDN04"
	DeclTypeTNDN02    DeclarationType = "TNDN02"
	DeclTypeTNDN05    DeclarationType = "TNDN05"
	DeclTypeTNDN06    DeclarationType = "TNDN06"
	DeclTypeKKTNCN    DeclarationType = "KK_TNCN"
	DeclTypeQTTTNCN   DeclarationType = "QTT_TNCN"
	DeclTypeTTDB01    DeclarationType = "TTDB01"
	DeclTypeBVMT01    DeclarationType = "BVMT01"
	DeclTypeNTNN01    DeclarationType = "NTNN01"
	DeclTypeNTNN02    DeclarationType = "NTNN02"
	DeclTypeNTNN03    DeclarationType = "NTNN03"
)
func (dt DeclarationType) Valid() bool {
	switch dt {
	case DeclTypeGTGT01, DeclTypeGTGT02, DeclTypeGTGT03, DeclTypeGTGT04, DeclTypeGTGT05,
		DeclTypeTNDN03, DeclTypeTNDN04, DeclTypeTNDN02, DeclTypeTNDN05, DeclTypeTNDN06,
		DeclTypeKKTNCN, DeclTypeQTTTNCN,
		DeclTypeTTDB01, DeclTypeBVMT01,
		DeclTypeNTNN01, DeclTypeNTNN02, DeclTypeNTNN03:
		return true
	}
	return false
}

type DeclarationStatus string
const (
	DeclStatusDRAFT         DeclarationStatus = "DRAFT"
	DeclStatusVALIDATED     DeclarationStatus = "VALIDATED"
	DeclStatusSUBMITTED     DeclarationStatus = "SUBMITTED"
	DeclStatusACKNOWLEDGED  DeclarationStatus = "ACKNOWLEDGED"
	DeclStatusREJECTED      DeclarationStatus = "REJECTED"
	DeclStatusAMENDED       DeclarationStatus = "AMENDED"
	DeclStatusCANCELLED     DeclarationStatus = "CANCELLED"
	DeclStatusFROZEN        DeclarationStatus = "FROZEN"
)
func (s DeclarationStatus) Valid() bool {
	switch s {
	case DeclStatusDRAFT, DeclStatusVALIDATED, DeclStatusSUBMITTED,
		DeclStatusACKNOWLEDGED, DeclStatusREJECTED,
		DeclStatusAMENDED, DeclStatusCANCELLED, DeclStatusFROZEN:
		return true
	}
	return false
}
func (s DeclarationStatus) IsTerminal() bool {
	return s == DeclStatusACKNOWLEDGED || s == DeclStatusCANCELLED || s == DeclStatusFROZEN
}
func (s DeclarationStatus) CanSubmit() bool {
	return s == DeclStatusDRAFT || s == DeclStatusVALIDATED || s == DeclStatusREJECTED
}

type TaxType string
const (
	TaxTypeVAT      TaxType = "VAT"
	TaxTypeCIT      TaxType = "CIT"
	TaxTypePIT      TaxType = "PIT"
	TaxTypeTTDB     TaxType = "TTDB"
	TaxTypeBVMT     TaxType = "BVMT"
	TaxTypeFCT      TaxType = "FCT"
	TaxTypeRESOURCE TaxType = "RESOURCE"
	TaxTypeLAND     TaxType = "LAND"
)

type RateType string
const (
	RateTypePERCENTAGE    RateType = "PERCENTAGE"
	RateTypeFIXED         RateType = "FIXED"
	RateTypePROGRESSIVE   RateType = "PROGRESSIVE"
)

type EInvoiceType string
const (
	EInvTypeORIGINAL        EInvoiceType = "ORIGINAL"
	EInvTypeADJUSTMENT      EInvoiceType = "ADJUSTMENT"
	EInvTypeREPLACEMENT     EInvoiceType = "REPLACEMENT"
	EInvTypeCANCELLATION_NOTE EInvoiceType = "CANCELLATION_NOTE"
)

type EInvLifecycleStatus string
const (
	EInvStatusDRAFT      EInvLifecycleStatus = "DRAFT"
	EInvStatusSIGNED     EInvLifecycleStatus = "SIGNED"
	EInvStatusSUBMITTED  EInvLifecycleStatus = "SUBMITTED"
	EInvStatusVALIDATED  EInvLifecycleStatus = "VALIDATED"
	EInvStatusISSUED     EInvLifecycleStatus = "ISSUED"
	EInvStatusCANCELLED  EInvLifecycleStatus = "CANCELLED"
	EInvStatusREPLACED   EInvLifecycleStatus = "REPLACED"
)
func (s EInvLifecycleStatus) Valid() bool {
	switch s {
	case EInvStatusDRAFT, EInvStatusSIGNED, EInvStatusSUBMITTED,
		EInvStatusVALIDATED, EInvStatusISSUED,
		EInvStatusCANCELLED, EInvStatusREPLACED:
		return true
	}
	return false
}
func (s EInvLifecycleStatus) CanCancel() bool {
	return s == EInvStatusISSUED
}

type CalendarStatus string
const (
	CalStatusPENDING   CalendarStatus = "PENDING"
	CalStatusSUBMITTED CalendarStatus = "SUBMITTED"
	CalStatusPAID      CalendarStatus = "PAID"
	CalStatusMISSED    CalendarStatus = "MISSED"
	CalStatusOVERDUE   CalendarStatus = "OVERDUE"
)

type PaymentMethod string
const (
	PayMethodEFT          PaymentMethod = "EFT"
	PayMethodBANK_TRANSFER PaymentMethod = "BANK_TRANSFER"
	PayMethodCASH         PaymentMethod = "CASH"
	PayMethodCQ           PaymentMethod = "CQ"
)

type PaymentStatus string
const (
	PayStatusPENDING  PaymentStatus = "PENDING"
	PayStatusPAID     PaymentStatus = "PAID"
	PayStatusPARTIAL  PaymentStatus = "PARTIAL"
	PayStatusOVERPAID PaymentStatus = "OVERPAID"
)

type AdjustmentType string
const (
	AdjTypeNONE       AdjustmentType = "NONE"
	AdjTypeAMENDMENT  AdjustmentType = "AMENDMENT"
	AdjTypeADDITIONAL AdjustmentType = "ADDITIONAL"
)

type SourceType string
const (
	SrcTypeAUTO       SourceType = "AUTO_CALCULATED"
	SrcTypeMANUAL     SourceType = "MANUAL_ENTRY"
	SrcTypeFROM_LEDGER SourceType = "FROM_LEDGER"
)

type AlertType string
const (
	AlertTypeINFO      AlertType = "INFO"
	AlertTypeWARNING   AlertType = "WARNING"
	AlertTypeCRITICAL  AlertType = "CRITICAL"
	AlertTypeDUE_TODAY AlertType = "DUE_TODAY"
)

type AlertChannel string
const (
	AlertChanEMAIL AlertChannel = "EMAIL"
	AlertChanINAPP AlertChannel = "IN_APP"
	AlertChanSMS   AlertChannel = "SMS"
	AlertChanALL   AlertChannel = "ALL"
)

type AuditCaseStatus string
const (
	AuditCaseOPEN       AuditCaseStatus = "OPEN"
	AuditCaseINPROGRESS AuditCaseStatus = "IN_PROGRESS"
	AuditCaseCLOSED     AuditCaseStatus = "CLOSED"
)

// ─── Tax Module Structs ──────────────────────────────────────────

type TaxPeriod struct {
	PeriodType   PeriodTypeV2 `json:"period_type"`
	PeriodYear   int         `json:"period_year"`
	PeriodNumber int         `json:"period_number"`
}

type ProgressiveBracket struct {
	MinAmount   float64 `json:"min_amount"`
	MaxAmount   float64 `json:"max_amount"`
	RatePercent float64 `json:"rate_percent"`
	FlatAmount  float64 `json:"flat_amount,omitempty"`
}

type TaxRate struct {
	ID              string               `json:"id"`
	TaxType         TaxType              `json:"tax_type"`
	RateCode        string               `json:"rate_code"`
	RateName        string               `json:"rate_name"`
	RateType        RateType             `json:"rate_type"`
	RateValue       float64              `json:"rate_value,omitempty"`
	Brackets        []ProgressiveBracket `json:"brackets,omitempty"`
	EffectiveFrom   string               `json:"effective_from"`
	EffectiveTo     string               `json:"effective_to,omitempty"`
	IsActive        bool                 `json:"is_active"`
	ApplicableTo    string               `json:"applicable_to,omitempty"`
	LegalRef        string               `json:"legal_ref,omitempty"`
	CreatedAt       string               `json:"created_at,omitempty"`
}

type DeclarationSignature struct {
	SignatureID string `json:"signature_id"`
	SignedAt    string `json:"signed_at"`
	SignedBy    string `json:"signed_by"`
	SignatureData string `json:"signature_data,omitempty"`
}

type TaxDeclaration struct {
	ID                  string               `json:"id"`
	CompanyID           string               `json:"company_id"`
	DeclarationType     DeclarationType      `json:"declaration_type"`
	TaxPeriod           TaxPeriod            `json:"tax_period"`
	Status              DeclarationStatus    `json:"status"`
	SubmittedAt         string               `json:"submitted_at,omitempty"`
	SubmittedBy         string               `json:"submitted_by,omitempty"`
	AcknowledgedAt      string               `json:"acknowledged_at,omitempty"`
	AcknowledgementRef  string               `json:"acknowledgement_ref,omitempty"`
	GDTResponseXML      string               `json:"gdt_response_xml,omitempty"`
	DeclarationXML      string               `json:"declaration_xml,omitempty"`
	PreviousDeclID      string               `json:"previous_declaration_id,omitempty"`
	AdjustmentType      AdjustmentType       `json:"adjustment_type"`
	Lines               []TaxDeclarationLine `json:"lines,omitempty"`
	Signatures          []DeclarationSignature `json:"signatures,omitempty"`
	Version             int                  `json:"version"`
	CreatedAt           string               `json:"created_at"`
	CreatedBy           string               `json:"created_by"`
	UpdatedAt           string               `json:"updated_at"`
}

func (d *TaxDeclaration) Validate() error {
	if d.CompanyID == "" { return ErrCompanyIDRequired }
	if !d.DeclarationType.Valid() { return ErrDeclarationTypeInvalid }
	if d.TaxPeriod.PeriodYear < 2000 || d.TaxPeriod.PeriodYear > 2100 { return ErrPeriodYearOutOfRange }
	if d.TaxPeriod.PeriodNumber < 1 { return ErrPeriodNumberInvalid }
	switch d.TaxPeriod.PeriodType {
	case PeriodTypeMonthly:
		if d.TaxPeriod.PeriodNumber > 12 { return ErrPeriodNumberInvalid }
	case PeriodTypeQuarterly:
		if d.TaxPeriod.PeriodNumber > 4 { return ErrPeriodNumberInvalid }
	}
	if d.Status == "" { d.Status = DeclStatusDRAFT }
	if !d.Status.Valid() { return ErrDeclarationStatusInvalid }
	return nil
}

func (d *TaxDeclaration) CanSubmit() bool {
	return d.Status.CanSubmit()
}

func (d *TaxDeclaration) CanAmend() bool {
	return d.Status == DeclStatusACKNOWLEDGED
}

type TaxDeclarationLine struct {
	ID              string     `json:"id,omitempty"`
	DeclarationID   string     `json:"declaration_id,omitempty"`
	LineCode        string     `json:"line_code"`
	LineName        string     `json:"line_name"`
	Amount          float64    `json:"amount"`
	SourceType      SourceType `json:"source_type"`
	SourceAccount   string     `json:"source_account,omitempty"`
	SourceEntryIDs  []string   `json:"source_entry_ids,omitempty"`
	Note            string     `json:"note,omitempty"`
	SortOrder       int        `json:"sort_order"`
}

type TaxPayment struct {
	ID              string        `json:"id"`
	CompanyID       string        `json:"company_id"`
	DeclarationID   string        `json:"declaration_id,omitempty"`
	TaxType         TaxType       `json:"tax_type"`
	PeriodYear      int           `json:"period_year"`
	PeriodNumber    int           `json:"period_number"`
	DeclaredAmount  float64       `json:"declared_amount"`
	PaidAmount      float64       `json:"paid_amount"`
	PaymentDate     string        `json:"payment_date,omitempty"`
	DueDate         string        `json:"due_date"`
	PaymentRef      string        `json:"payment_ref,omitempty"`
	PaymentMethod   PaymentMethod `json:"payment_method,omitempty"`
	Status          PaymentStatus `json:"status"`
	LateDays        int           `json:"late_days,omitempty"`
	LateInterest    float64       `json:"late_interest,omitempty"`
	Notes           string        `json:"notes,omitempty"`
	CreatedAt       string        `json:"created_at"`
}

func (p *TaxPayment) Validate() error {
	if p.CompanyID == "" { return ErrCompanyIDRequired }
	if p.DeclaredAmount <= 0 { return ErrPaymentAmountRequired }
	if p.DueDate == "" { return ErrPaymentDueDateRequired }
	if p.Status == "" { p.Status = PayStatusPENDING }
	return nil
}

func (p *TaxPayment) CalculateLateInterest() {
	if p.PaidAmount >= p.DeclaredAmount {
		p.LateDays = 0
		p.LateInterest = 0
		return
	}
	underpaid := p.DeclaredAmount - p.PaidAmount
	p.LateInterest = underpaid * 0.0003 * float64(p.LateDays)
}

type EInvoiceLine struct {
	ID          string  `json:"id,omitempty"`
	EInvoiceID  string  `json:"e_invoice_id,omitempty"`
	LineNumber  int     `json:"line_number"`
	Description string  `json:"description"`
	Unit        string  `json:"unit,omitempty"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	LineTotal   float64 `json:"line_total"`
	VATRate     float64 `json:"vat_rate"`
	VATAmount   float64 `json:"vat_amount"`
}

type EInvoice struct {
	ID                string         `json:"id"`
	CompanyID         string         `json:"company_id"`
	Pattern           string         `json:"pattern"`
	Serial            string         `json:"serial"`
	InvoiceNumber     int64          `json:"invoice_number,omitempty"`
	InvoiceType       EInvoiceType   `json:"invoice_type"`
	GDTTransactionID  string         `json:"gdt_transaction_id,omitempty"`
	BuyerName         string         `json:"buyer_name"`
	BuyerTaxCode      string         `json:"buyer_tax_code,omitempty"`
	BuyerAddress      string         `json:"buyer_address,omitempty"`
	BuyerEmail        string         `json:"buyer_email,omitempty"`
	CurrencyCode      string         `json:"currency_code"`
	ExchangeRate      float64        `json:"exchange_rate,omitempty"`
	Lines             []EInvoiceLine `json:"lines"`
	Subtotal          float64        `json:"subtotal"`
	VATAmount         float64        `json:"vat_amount"`
	GrandTotal        float64        `json:"grand_total"`
	XMLBody           string         `json:"xml_body,omitempty"`
	SignedXML         string         `json:"signed_xml,omitempty"`
	IssueDate         string         `json:"issue_date"`
	SigningDate       string         `json:"signing_date,omitempty"`
	DigitalSignatureID string        `json:"digital_signature_id,omitempty"`
	JournalEntryID    string         `json:"journal_entry_id,omitempty"`
	Status            EInvLifecycleStatus `json:"status"`
	CancelledAt       string         `json:"cancelled_at,omitempty"`
	CancelReason      string         `json:"cancel_reason,omitempty"`
	OriginalInvoiceID string         `json:"original_invoice_id,omitempty"`
	GDTResponse       string         `json:"gdt_response,omitempty"`
	CreatedAt         string         `json:"created_at"`
}

func (inv *EInvoice) Validate() error {
	if inv.CompanyID == "" { return ErrCompanyIDRequired }
	if inv.BuyerName == "" { return ErrBuyerNameRequired }
	if inv.Pattern == "" { return ErrInvoicePatternRequired }
	if len(inv.Lines) == 0 { return ErrInvoiceNoLines }
	totalVAT := 0.0
	totalLines := 0.0
	for _, line := range inv.Lines {
		computedVAT := line.LineTotal * line.VATRate / 100.0
		if computedVAT != line.VATAmount {
			return ErrInvoiceVATMismatch
		}
		totalVAT += line.VATAmount
		totalLines += line.LineTotal
	}
	if totalLines != inv.Subtotal { return ErrInvoiceSubtotalMismatch }
	if totalVAT != inv.VATAmount { return ErrInvoiceVATTotalMismatch }
	if inv.Subtotal + inv.VATAmount != inv.GrandTotal { return ErrInvoiceGrandTotalMismatch }
	if inv.Status == "" { inv.Status = EInvStatusDRAFT }
	if !inv.Status.Valid() { return ErrInvoiceStatusInvalid }
	if inv.CurrencyCode == "" { inv.CurrencyCode = "VND" }
	switch inv.InvoiceType {
	case EInvTypeADJUSTMENT, EInvTypeREPLACEMENT, EInvTypeCANCELLATION_NOTE:
		if inv.OriginalInvoiceID == "" { return ErrOriginalInvoiceRequired }
	}
	return nil
}

type TaxCalendar struct {
	ID              string         `json:"id"`
	CompanyID       string         `json:"company_id"`
	TaxType         TaxType        `json:"tax_type"`
	PeriodType      PeriodTypeV2   `json:"period_type"`
	PeriodYear      int            `json:"period_year"`
	PeriodNumber    int            `json:"period_number"`
	StartDate       string         `json:"start_date"`
	EndDate         string         `json:"end_date"`
	DeclarationDue  string         `json:"declaration_due"`
	PaymentDue      string         `json:"payment_due,omitempty"`
	Status          CalendarStatus `json:"status"`
	DeclarationID   string         `json:"declaration_id,omitempty"`
	CreatedAt       string         `json:"created_at"`
}

type TaxAlert struct {
	ID            string       `json:"id"`
	CompanyID     string       `json:"company_id"`
	CalendarID    string       `json:"calendar_id,omitempty"`
	AlertType     AlertType    `json:"alert_type"`
	Channel       AlertChannel `json:"channel"`
	Message       string       `json:"message"`
	SentAt        string       `json:"sent_at"`
	AcknowledgedAt string      `json:"acknowledged_at,omitempty"`
	AcknowledgedBy string      `json:"acknowledged_by,omitempty"`
}

type TaxAuditCase struct {
	ID               string         `json:"id"`
	CompanyID        string         `json:"company_id"`
	AuditPeriodStart string         `json:"audit_period_start"`
	AuditPeriodEnd   string         `json:"audit_period_end"`
	AuditDecNumber   string         `json:"audit_decision_number"`
	AuditorName      string         `json:"auditor_name"`
	AuditorContact   string         `json:"auditor_contact,omitempty"`
	Status           AuditCaseStatus `json:"status"`
	Findings         string         `json:"findings,omitempty"`
	PenaltyAmount    float64        `json:"penalty_amount,omitempty"`
	CreatedAt        string         `json:"created_at"`
	ClosedAt         string         `json:"closed_at,omitempty"`
}

// ─── Calculation Result Structs ──────────────────────────────────

type VATResult struct {
	CompanyID        string  `json:"company_id"`
	Period           TaxPeriod `json:"period"`
	OutputVAT        float64 `json:"output_vat"`
	InputVAT         float64 `json:"input_vat"`
	InputVATFA       float64 `json:"input_vat_fa"`
	TotalInputVAT    float64 `json:"total_input_vat"`
	VATPayable       float64 `json:"vat_payable"`
	VATRefundable    float64 `json:"vat_refundable"`
	PurchaseTotal    float64 `json:"purchase_total"`
	SalesTotal       float64 `json:"sales_total"`
}

type CITResult struct {
	CompanyID        string    `json:"company_id"`
	PeriodYear       int       `json:"period_year"`
	PeriodType       string    `json:"period_type"`
	Revenue          float64   `json:"revenue"`
	Expenses         float64   `json:"expenses"`
	NonDeductible    float64   `json:"non_deductible"`
	TaxableIncome    float64   `json:"taxable_income"`
	TaxRate          float64   `json:"tax_rate"`
	CITPayable       float64   `json:"cit_payable"`
	IncentiveReduc   float64   `json:"incentive_reduction"`
	CITFinal         float64   `json:"cit_final"`
	Provisionals     float64   `json:"provisionals"`
	BalanceDue       float64   `json:"balance_due"`
	Refund           float64   `json:"refund"`
	LateInterest     float64   `json:"late_interest"`
}

type PITResult struct {
	CompanyID       string  `json:"company_id"`
	Period          TaxPeriod `json:"period"`
	EmployeeCount   int     `json:"employee_count"`
	TotalGross      float64 `json:"total_gross"`
	TotalDeductions float64 `json:"total_deductions"`
	TotalPIT        float64 `json:"total_pit"`
}

// ─── Filter Structs ─────────────────────────────────────────────

type TaxDeclarationFilter struct {
	CompanyID       string           `json:"company_id,omitempty"`
	DeclarationType DeclarationType  `json:"declaration_type,omitempty"`
	Status          DeclarationStatus `json:"status,omitempty"`
	PeriodYear      int              `json:"period_year,omitempty"`
	PeriodNumber    int              `json:"period_number,omitempty"`
	FromDate        string           `json:"from_date,omitempty"`
	ToDate          string           `json:"to_date,omitempty"`
}

type TaxRateFilter struct {
	TaxType     TaxType `json:"tax_type,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
	EffectiveOn string  `json:"effective_on,omitempty"`
}

type PaymentFilter struct {
	CompanyID    string        `json:"company_id,omitempty"`
	TaxType      TaxType       `json:"tax_type,omitempty"`
	Status       PaymentStatus `json:"status,omitempty"`
	PeriodYear   int           `json:"period_year,omitempty"`
	PeriodNumber int           `json:"period_number,omitempty"`
}

type EInvoiceFilter struct {
	CompanyID string         `json:"company_id,omitempty"`
	Status    EInvLifecycleStatus `json:"status,omitempty"`
	FromDate  string         `json:"from_date,omitempty"`
	ToDate    string         `json:"to_date,omitempty"`
}

type EInvoiceInput struct {
	CompanyID         string         `json:"company_id"`
	Pattern           string         `json:"pattern"`
	Serial            string         `json:"serial"`
	InvoiceType       EInvoiceType   `json:"invoice_type"`
	BuyerName         string         `json:"buyer_name"`
	BuyerTaxCode      string         `json:"buyer_tax_code,omitempty"`
	BuyerAddress      string         `json:"buyer_address,omitempty"`
	BuyerEmail        string         `json:"buyer_email,omitempty"`
	CurrencyCode      string         `json:"currency_code,omitempty"`
	ExchangeRate      float64        `json:"exchange_rate,omitempty"`
	Lines             []EInvoiceLine `json:"lines"`
	DigitalSignatureID string        `json:"digital_signature_id,omitempty"`
	IssueDate         string         `json:"issue_date"`
}

// ─── Password Policy ──────────────────────────────────────────────

var PasswordMaxAge = 90 * 24 * time.Hour
var PasswordMinLength = 12
var PasswordHistorySize = 10
var PasswordRequireUpper = true
var PasswordRequireLower = true
var PasswordRequireDigit = true
var PasswordRequireSpecial = true
var MaxLoginAttempts = 5
var LockoutDuration = 30 * time.Minute

func ValidatePassword(password string) error {
	if len(password)<PasswordMinLength { return ErrPasswordTooShort }
	var hasUpper,hasLower,hasDigit,hasSpecial bool
	for _,c:=range password{
		switch{
		case unicode.IsUpper(c): hasUpper=true
		case unicode.IsLower(c): hasLower=true
		case unicode.IsDigit(c): hasDigit=true
		case unicode.IsPunct(c)||unicode.IsSymbol(c): hasSpecial=true
		}
	}
	if PasswordRequireUpper&&!hasUpper { return ErrPasswordNoUpper }
	if PasswordRequireLower&&!hasLower { return ErrPasswordNoLower }
	if PasswordRequireDigit&&!hasDigit { return ErrPasswordNoDigit }
	if PasswordRequireSpecial&&!hasSpecial { return ErrPasswordNoSpecial }
	return nil
}

func IsPasswordInHistory(history []string, newHash string) bool {
	for _,h:=range history{if h==newHash{return true}}; return false
}
