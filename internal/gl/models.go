package gl

import (
	"errors"
	"math"
	"time"
	"unicode"
)

type AccountType string

const (
	AccountTypeAsset      AccountType = "ASSET"
	AccountTypeLiability  AccountType = "LIABILITY"
	AccountTypeEquity     AccountType = "EQUITY"
	AccountTypeRevenue    AccountType = "REVENUE"
	AccountTypeExpense    AccountType = "EXPENSE"
)

type NormalBalance string

const (
	NormalBalanceDebit  NormalBalance = "DEBIT"
	NormalBalanceCredit NormalBalance = "CREDIT"
)

type VoucherType string

const (
	VoucherTypeReceipt    VoucherType = "THU"
	VoucherTypePayment    VoucherType = "CHI"
	VoucherTypeSale       VoucherType = "BAN"
	VoucherTypePurchase   VoucherType = "MUA"
	VoucherTypeOther      VoucherType = "KHAC"
	VoucherTypeClosing    VoucherType = "KC"
)

type DetailBy string

const (
	DetailByNone       DetailBy = ""
	DetailByObject     DetailBy = "OBJECT"
	DetailByProject    DetailBy = "PROJECT"
	DetailByContract   DetailBy = "CONTRACT"
	DetailByCostItem   DetailBy = "COST_ITEM"
	DetailByDepartment DetailBy = "DEPARTMENT"
)

type UserRole string

const (
	UserRoleAdmin          UserRole = "admin"
	UserRoleChiefAccountant UserRole = "chief_accountant"
	UserRoleAccountant     UserRole = "accountant"
	UserRoleViewer         UserRole = "viewer"
)

type AuditAction string

const (
	AuditActionCreate  AuditAction = "CREATE"
	AuditActionUpdate  AuditAction = "UPDATE"
	AuditActionDelete  AuditAction = "DELETE"
	AuditActionApprove AuditAction = "APPROVE"
	AuditActionPost    AuditAction = "POST"
	AuditActionCancel  AuditAction = "CANCEL"
	AuditActionClose   AuditAction = "CLOSE"
	AuditActionReview  AuditAction = "REVIEW"
	AuditActionLogin   AuditAction = "LOGIN"
)

var (
	ErrAccountCodeRequired       = errors.New("account code is required")
	ErrAccountNameRequired       = errors.New("account name is required")
	ErrAccountInvalidType        = errors.New("account type is invalid")
	ErrAccountCodeInvalid        = errors.New("account code must be digits, min 3 chars")
	ErrJournalEntryNoLines       = errors.New("journal entry must have at least one line")
	ErrJournalEntryUnbalanced    = errors.New("total debit must equal total credit")
	ErrJournalEntryInvalidDate   = errors.New("entry date is required")
	ErrJournalEntryNoDescription = errors.New("description is required")
	ErrJournalLineInvalidAmount  = errors.New("line amount must be positive")
	ErrInvalidCurrencyCode       = errors.New("currency code must be 3 chars")
	ErrInvalidExchangeRate       = errors.New("exchange rate must be positive")
	ErrInvalidUserRole           = errors.New("invalid user role")
	ErrUsernameRequired          = errors.New("username is required")
	ErrPasswordRequired          = errors.New("password is required")
)

type Account struct {
	Code        string      `json:"code"`
	Name        string      `json:"name"`
	Name2       string      `json:"name2,omitempty"`
	Type        AccountType `json:"type"`
	ParentCode  string      `json:"parent_code,omitempty"`
	IsActive    bool        `json:"is_active"`
	IsForeign   bool        `json:"is_foreign"`
	DetailBy    DetailBy    `json:"detail_by,omitempty"`
	IsParent    bool        `json:"is_parent"`
	ArrearsDays int         `json:"arrears_days,omitempty"`
	Note        string      `json:"note,omitempty"`
}

func (a *Account) Validate() error {
	if a.Code == "" {
		return ErrAccountCodeRequired
	}
	if len(a.Code) < 3 {
		return ErrAccountCodeInvalid
	}
	for _, c := range a.Code {
		if !unicode.IsDigit(c) {
			return ErrAccountCodeInvalid
		}
	}
	if a.Name == "" {
		return ErrAccountNameRequired
	}
	switch a.Type {
	case AccountTypeAsset, AccountTypeLiability, AccountTypeEquity, AccountTypeRevenue, AccountTypeExpense:
	default:
		return ErrAccountInvalidType
	}
	return nil
}

func (at AccountType) NormalBalance() NormalBalance {
	switch at {
	case AccountTypeAsset, AccountTypeExpense:
		return NormalBalanceDebit
	case AccountTypeLiability, AccountTypeEquity, AccountTypeRevenue:
		return NormalBalanceCredit
	default:
		return NormalBalanceDebit
	}
}

type JournalLine struct {
	ID           string  `json:"id,omitempty"`
	EntryID      string  `json:"entry_id,omitempty"`
	LineNumber   int     `json:"line_number"`
	AccountCode  string  `json:"account_code"`
	DebitAmount  float64 `json:"debit_amount"`
	CreditAmount float64 `json:"credit_amount"`
	Description  string  `json:"description,omitempty"`
	CurrencyCode string  `json:"currency_code,omitempty"`
	ForeignAmount float64 `json:"foreign_amount,omitempty"`
	ExchangeRate float64 `json:"exchange_rate,omitempty"`
	ObjectID     string  `json:"object_id,omitempty"`
	ProjectID    string  `json:"project_id,omitempty"`
	ContractID   string  `json:"contract_id,omitempty"`
	CostItemID   string  `json:"cost_item_id,omitempty"`
	DepartmentID string  `json:"department_id,omitempty"`
}

type JournalEntryStatus string

const (
	JournalEntryDraft     JournalEntryStatus = "DRAFT"
	JournalEntryReviewing JournalEntryStatus = "REVIEWING"
	JournalEntryApproved  JournalEntryStatus = "APPROVED"
	JournalEntryPosted    JournalEntryStatus = "POSTED"
	JournalEntryCancelled JournalEntryStatus = "CANCELLED"
)

type JournalEntry struct {
	ID             string             `json:"id,omitempty"`
	CompanyID      string             `json:"company_id,omitempty"`
	EntryNumber    string             `json:"entry_number"`
	VoucherType    VoucherType        `json:"voucher_type,omitempty"`
	EntryDate      time.Time          `json:"entry_date"`
	AccountingDate time.Time          `json:"accounting_date,omitempty"`
	PeriodID       string             `json:"period_id,omitempty"`
	Description    string             `json:"description"`
	Status         JournalEntryStatus `json:"status"`
	CurrencyCode   string             `json:"currency_code,omitempty"`
	ExchangeRate   float64            `json:"exchange_rate,omitempty"`
	Lines          []JournalLine      `json:"lines"`
	CreatedBy      string             `json:"created_by,omitempty"`
	ReviewedBy     string             `json:"reviewed_by,omitempty"`
	ApprovedBy     string             `json:"approved_by,omitempty"`
	CreatedAt      time.Time          `json:"created_at,omitempty"`
	PostedAt       *time.Time         `json:"posted_at,omitempty"`
	ApprovedAt     *time.Time         `json:"approved_at,omitempty"`
}

func (je *JournalEntry) TotalDebit() float64 {
	var total float64
	for _, l := range je.Lines {
		total += l.DebitAmount
	}
	return total
}

func (je *JournalEntry) TotalCredit() float64 {
	var total float64
	for _, l := range je.Lines {
		total += l.CreditAmount
	}
	return total
}

func (je *JournalEntry) HasDebit() bool {
	return je.TotalDebit() > 0
}

func (je *JournalEntry) HasCredit() bool {
	return je.TotalCredit() > 0
}

func (je *JournalEntry) Validate() error {
	if je.EntryDate.IsZero() {
		return ErrJournalEntryInvalidDate
	}
	if je.Description == "" {
		return ErrJournalEntryNoDescription
	}
	if je.CurrencyCode == "" {
		je.CurrencyCode = "VND"
	}
	if je.CurrencyCode != "VND" && je.ExchangeRate <= 0 {
		return ErrInvalidExchangeRate
	}
	if len(je.Lines) == 0 {
		return ErrJournalEntryNoLines
	}
	for _, l := range je.Lines {
		if l.DebitAmount < 0 || l.CreditAmount < 0 {
			return ErrJournalLineInvalidAmount
		}
		if l.DebitAmount == 0 && l.CreditAmount == 0 {
			return ErrJournalLineInvalidAmount
		}
	}
	td := je.TotalDebit()
	tc := je.TotalCredit()
	if math.Abs(td-tc) > 0.001 {
		return ErrJournalEntryUnbalanced
	}
	return nil
}

type PeriodStatus string

const (
	PeriodOpen    PeriodStatus = "OPEN"
	PeriodClosing PeriodStatus = "CLOSING"
	PeriodClosed  PeriodStatus = "CLOSED"
	PeriodLocked  PeriodStatus = "LOCKED"
)

type Period struct {
	ID        string       `json:"id"`
	Year      int          `json:"year"`
	Month     int          `json:"month"`
	StartDate time.Time    `json:"start_date"`
	EndDate   time.Time    `json:"end_date"`
	Status    PeriodStatus `json:"status"`
}

type AccountBalance struct {
	AccountCode       string      `json:"account_code"`
	AccountType       AccountType `json:"account_type"`
	PeriodID          string      `json:"period_id"`
	OpenBalanceDebit  float64     `json:"open_balance_debit"`
	OpenBalanceCredit float64     `json:"open_balance_credit"`
	PeriodDebit       float64     `json:"period_debit"`
	PeriodCredit      float64     `json:"period_credit"`
	TotalDebit        float64     `json:"total_debit"`
	TotalCredit       float64     `json:"total_credit"`
	ClosingBalance    float64     `json:"closing_balance"`
}

func (b *AccountBalance) Calculate() {
	b.TotalDebit = b.OpenBalanceDebit + b.PeriodDebit
	b.TotalCredit = b.OpenBalanceCredit + b.PeriodCredit
	b.ClosingBalance = b.TotalDebit - b.TotalCredit
}

type User struct {
	ID           string   `json:"id"`
	Username     string   `json:"username"`
	PasswordHash string   `json:"-"`
	FullName     string   `json:"full_name"`
	Email        string   `json:"email,omitempty"`
	Role         UserRole `json:"role"`
	IsActive     bool     `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (u *User) Validate() error {
	if u.Username == "" {
		return ErrUsernameRequired
	}
	switch u.Role {
	case UserRoleAdmin, UserRoleChiefAccountant, UserRoleAccountant, UserRoleViewer:
	default:
		return ErrInvalidUserRole
	}
	return nil
}

type AuditEntry struct {
	ID         string      `json:"id"`
	UserID     string      `json:"user_id,omitempty"`
	Username   string      `json:"username"`
	IPAddress  string      `json:"ip_address,omitempty"`
	Action     AuditAction `json:"action"`
	EntityType string      `json:"entity_type"`
	EntityID   string      `json:"entity_id,omitempty"`
	OldValue   string      `json:"old_value,omitempty"`
	NewValue   string      `json:"new_value,omitempty"`
	CreatedAt  time.Time   `json:"created_at"`
}

type ExchangeRate struct {
	ID           string    `json:"id,omitempty"`
	CurrencyCode string    `json:"currency_code"`
	RateDate     time.Time `json:"rate_date"`
	BuyRate      float64   `json:"buy_rate,omitempty"`
	SellRate     float64   `json:"sell_rate,omitempty"`
	AverageRate  float64   `json:"average_rate"`
	Source       string    `json:"source,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
}

func (r *ExchangeRate) Validate() error {
	if len(r.CurrencyCode) != 3 {
		return ErrInvalidCurrencyCode
	}
	if r.AverageRate <= 0 {
		return ErrInvalidExchangeRate
	}
	return nil
}

type ClosingTemplate struct {
	ID            string    `json:"id,omitempty"`
	Name          string    `json:"name"`
	Description   string    `json:"description,omitempty"`
	SequenceOrder int       `json:"sequence_order"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
}

type ClosingTemplateLine struct {
	ID            string `json:"id,omitempty"`
	TemplateID    string `json:"template_id,omitempty"`
	LineNumber    int    `json:"line_number"`
	DebitAccount  string `json:"debit_account"`
	CreditAccount string `json:"credit_account"`
	Formula       string `json:"formula"`
	Description   string `json:"description,omitempty"`
}
