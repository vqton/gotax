package gl

import (
	"context"
	"errors"
	"time"
)

var (
	ErrAccountNotFound       = errors.New("account not found")
	ErrAccountCodeExists     = errors.New("account code already exists")
	ErrAccountHasChildren    = errors.New("account has children, cannot delete")
	ErrAccountHasBalance     = errors.New("account has balance, cannot delete")
	ErrJournalNotFound       = errors.New("journal entry not found")
	ErrJournalAlreadyPosted  = errors.New("journal entry already posted")
	ErrJournalAlreadyCancelled = errors.New("journal entry already cancelled")
	ErrJournalPeriodClosed   = errors.New("period is closed, cannot post")
	ErrPeriodNotFound        = errors.New("period not found")
	ErrPeriodAlreadyClosed   = errors.New("period already closed")
	ErrAccountInactive       = errors.New("account is inactive")
	ErrJournalAlreadyDraft   = errors.New("journal entry is draft, not posted")
	ErrUserNotFound          = errors.New("user not found")
	ErrUsernameExists        = errors.New("username already exists")
	ErrUserInactive          = errors.New("user is inactive")
	ErrRateNotFound          = errors.New("exchange rate not found")
	ErrTemplateNotFound      = errors.New("closing template not found")
	ErrPeriodHasEntries      = errors.New("period has posted entries, cannot reopen")
	ErrJournalAlreadyReviewed = errors.New("journal entry already reviewed")
	ErrJournalAlreadyApproved = errors.New("journal entry already approved")
)

type AccountRepository interface {
	Create(ctx context.Context, account *Account) error
	GetByCode(ctx context.Context, code string) (*Account, error)
	GetAll(ctx context.Context, activeOnly bool) ([]Account, error)
	Update(ctx context.Context, account *Account) error
	Delete(ctx context.Context, code string) error
	GetChildren(ctx context.Context, parentCode string) ([]Account, error)
}

type JournalRepository interface {
	Create(ctx context.Context, entry *JournalEntry) error
	GetByID(ctx context.Context, id string) (*JournalEntry, error)
	GetByPeriod(ctx context.Context, periodID string) ([]JournalEntry, error)
	GetByDateRange(ctx context.Context, from, to time.Time) ([]JournalEntry, error)
	GetByStatus(ctx context.Context, status JournalEntryStatus) ([]JournalEntry, error)
	GetByVoucherType(ctx context.Context, voucherType VoucherType) ([]JournalEntry, error)
	UpdateStatus(ctx context.Context, id string, status JournalEntryStatus) error
	Update(ctx context.Context, entry *JournalEntry) error
	Approve(ctx context.Context, id, approvedBy string) error
	Review(ctx context.Context, id, reviewedBy string) error
	GetLinesByEntryID(ctx context.Context, entryID string) ([]JournalLine, error)
	GetBalance(ctx context.Context, accountCode string, periodID string) (*AccountBalance, error)
	GetTrialBalance(ctx context.Context, periodID string) ([]AccountBalance, error)
	GetFinancialStatement(ctx context.Context, periodID string, accountTypes []AccountType) ([]AccountBalance, error)
}

type PeriodRepository interface {
	Create(ctx context.Context, period *Period) error
	GetByID(ctx context.Context, id string) (*Period, error)
	GetByYearMonth(ctx context.Context, year, month int) (*Period, error)
	GetAll(ctx context.Context) ([]Period, error)
	UpdateStatus(ctx context.Context, id string, status PeriodStatus) error
	GetOpenPeriod(ctx context.Context) (*Period, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetAll(ctx context.Context) ([]User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}

type AuditLogRepository interface {
	Create(ctx context.Context, entry *AuditEntry) error
	GetByEntity(ctx context.Context, entityType, entityID string) ([]AuditEntry, error)
	GetByUser(ctx context.Context, userID string) ([]AuditEntry, error)
	GetByDateRange(ctx context.Context, from, to time.Time) ([]AuditEntry, error)
	GetAll(ctx context.Context, limit int) ([]AuditEntry, error)
}

type ExchangeRateRepository interface {
	Create(ctx context.Context, rate *ExchangeRate) error
	GetByCurrencyAndDate(ctx context.Context, currencyCode string, rateDate time.Time) (*ExchangeRate, error)
	GetByDateRange(ctx context.Context, from, to time.Time) ([]ExchangeRate, error)
	GetAll(ctx context.Context) ([]ExchangeRate, error)
	Delete(ctx context.Context, id string) error
}

type ClosingTemplateRepository interface {
	Create(ctx context.Context, template *ClosingTemplate) error
	GetByID(ctx context.Context, id string) (*ClosingTemplate, error)
	GetAll(ctx context.Context) ([]ClosingTemplate, error)
	Update(ctx context.Context, template *ClosingTemplate) error
	Delete(ctx context.Context, id string) error
}
