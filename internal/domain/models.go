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
