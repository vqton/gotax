package gl

import (
	"context"
	"fmt"
	"time"
)

type Service interface {
	CreateAccount(ctx context.Context, account *Account) error
	GetAccount(ctx context.Context, code string) (*Account, error)
	GetAllAccounts(ctx context.Context, activeOnly bool) ([]Account, error)
	UpdateAccount(ctx context.Context, account *Account) error
	DeleteAccount(ctx context.Context, code string) error

	CreateEntry(ctx context.Context, entry *JournalEntry, userID string) error
	SubmitForReview(ctx context.Context, id, userID string) error
	ReviewEntry(ctx context.Context, id, reviewerID string) error
	ApproveEntry(ctx context.Context, id, approverID string) error
	PostEntry(ctx context.Context, id string) error
	CancelEntry(ctx context.Context, id string) error
	GetEntryByID(ctx context.Context, id string) (*JournalEntry, error)
	GetEntriesByDateRange(ctx context.Context, from, to time.Time) ([]JournalEntry, error)
	GetEntriesByStatus(ctx context.Context, status JournalEntryStatus) ([]JournalEntry, error)

	TrialBalance(ctx context.Context, year, month int) ([]AccountBalance, error)
	BalanceSheet(ctx context.Context, year, month int) ([]AccountBalance, error)
	IncomeStatement(ctx context.Context, year, month int) ([]AccountBalance, error)

	CreatePeriod(ctx context.Context, period *Period) error
	GetPeriod(ctx context.Context, id string) (*Period, error)
	GetPeriodByYearMonth(ctx context.Context, year, month int) (*Period, error)
	GetAllPeriods(ctx context.Context) ([]Period, error)
	ClosePeriod(ctx context.Context, id string) error
	ReopenPeriod(ctx context.Context, id string) error

	CreateExchangeRate(ctx context.Context, rate *ExchangeRate) error
	GetExchangeRate(ctx context.Context, currencyCode string, rateDate time.Time) (*ExchangeRate, error)
	ListExchangeRates(ctx context.Context) ([]ExchangeRate, error)

	GetAuditLog(ctx context.Context, limit int) ([]AuditEntry, error)
	GetAuditLogByEntity(ctx context.Context, entityType, entityID string) ([]AuditEntry, error)

	CreateUser(ctx context.Context, user *User, password string) error
	GetUser(ctx context.Context, id string) (*User, error)
	ListUsers(ctx context.Context) ([]User, error)
	Authenticate(ctx context.Context, username, password string) (*User, error)
}

type service struct {
	accounts   AccountRepository
	journals   JournalRepository
	periods    PeriodRepository
	users      UserRepository
	audit      AuditLogRepository
	rates      ExchangeRateRepository
	templates  ClosingTemplateRepository
	now        func() time.Time
}

func NewService(
	accRepo AccountRepository,
	jeRepo JournalRepository,
	perRepo PeriodRepository,
	userRepo UserRepository,
	auditRepo AuditLogRepository,
	rateRepo ExchangeRateRepository,
	templateRepo ClosingTemplateRepository,
) Service {
	return &service{
		accounts:  accRepo,
		journals:  jeRepo,
		periods:   perRepo,
		users:     userRepo,
		audit:     auditRepo,
		rates:     rateRepo,
		templates: templateRepo,
		now:       time.Now,
	}
}

// -- Accounts --

func (s *service) CreateAccount(ctx context.Context, account *Account) error {
	if err := account.Validate(); err != nil {
		return err
	}
	existing, _ := s.accounts.GetByCode(ctx, account.Code)
	if existing != nil {
		return ErrAccountCodeExists
	}
	return s.accounts.Create(ctx, account)
}

func (s *service) GetAccount(ctx context.Context, code string) (*Account, error) {
	return s.accounts.GetByCode(ctx, code)
}

func (s *service) GetAllAccounts(ctx context.Context, activeOnly bool) ([]Account, error) {
	return s.accounts.GetAll(ctx, activeOnly)
}

func (s *service) UpdateAccount(ctx context.Context, account *Account) error {
	existing, err := s.accounts.GetByCode(ctx, account.Code)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrAccountNotFound
	}
	return s.accounts.Update(ctx, account)
}

func (s *service) DeleteAccount(ctx context.Context, code string) error {
	_, err := s.accounts.GetByCode(ctx, code)
	if err != nil {
		return err
	}
	children, err := s.accounts.GetChildren(ctx, code)
	if err != nil {
		return err
	}
	if len(children) > 0 {
		return ErrAccountHasChildren
	}
	return s.accounts.Delete(ctx, code)
}

// -- Journal Entries (State Machine) --

func (s *service) CreateEntry(ctx context.Context, entry *JournalEntry, userID string) error {
	if err := entry.Validate(); err != nil {
		return err
	}
	for _, line := range entry.Lines {
		acc, err := s.accounts.GetByCode(ctx, line.AccountCode)
		if err != nil {
			return err
		}
		if !acc.IsActive {
			return ErrAccountInactive
		}
	}
	entry.Status = JournalEntryDraft
	entry.CreatedBy = userID
	entry.CurrencyCode = "VND"
	entry.ExchangeRate = 1
	if entry.VoucherType == "" {
		entry.VoucherType = VoucherTypeOther
	}
	if entry.AccountingDate.IsZero() {
		entry.AccountingDate = entry.EntryDate
	}
	return s.journals.Create(ctx, entry)
}

func (s *service) SubmitForReview(ctx context.Context, id, userID string) error {
	entry, err := s.journals.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if entry.Status != JournalEntryDraft {
		return fmt.Errorf("entry must be DRAFT to submit for review")
	}
	return s.journals.Review(ctx, id, userID)
}

func (s *service) ReviewEntry(ctx context.Context, id, reviewerID string) error {
	entry, err := s.journals.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if entry.Status == JournalEntryReviewing {
		return ErrJournalAlreadyReviewed
	}
	if entry.Status != JournalEntryDraft {
		return fmt.Errorf("entry must be DRAFT for review")
	}
	return s.journals.Review(ctx, id, reviewerID)
}

func (s *service) ApproveEntry(ctx context.Context, id, approverID string) error {
	entry, err := s.journals.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if entry.Status == JournalEntryApproved {
		return ErrJournalAlreadyApproved
	}
	if entry.Status != JournalEntryReviewing {
		return fmt.Errorf("entry must be REVIEWING to approve")
	}
	return s.journals.Approve(ctx, id, approverID)
}

func (s *service) PostEntry(ctx context.Context, id string) error {
	entry, err := s.journals.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if entry.Status == JournalEntryPosted {
		return ErrJournalAlreadyPosted
	}
	if entry.Status == JournalEntryCancelled {
		return ErrJournalAlreadyCancelled
	}
	if entry.Status != JournalEntryApproved {
		return fmt.Errorf("entry must be APPROVED to post")
	}
	period, _ := s.periods.GetByID(ctx, entry.PeriodID)
	if period == nil {
		return ErrPeriodNotFound
	}
	if period.Status != PeriodOpen {
		return ErrJournalPeriodClosed
	}
	entry.Status = JournalEntryPosted
	now := s.now()
	entry.PostedAt = &now
	return s.journals.Update(ctx, entry)
}

func (s *service) CancelEntry(ctx context.Context, id string) error {
	entry, err := s.journals.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if entry.Status == JournalEntryCancelled {
		return ErrJournalAlreadyCancelled
	}
	if entry.Status == JournalEntryDraft {
		return ErrJournalAlreadyDraft
	}
	return s.journals.UpdateStatus(ctx, id, JournalEntryCancelled)
}

func (s *service) GetEntryByID(ctx context.Context, id string) (*JournalEntry, error) {
	return s.journals.GetByID(ctx, id)
}

func (s *service) GetEntriesByDateRange(ctx context.Context, from, to time.Time) ([]JournalEntry, error) {
	return s.journals.GetByDateRange(ctx, from, to)
}

func (s *service) GetEntriesByStatus(ctx context.Context, status JournalEntryStatus) ([]JournalEntry, error) {
	return s.journals.GetByStatus(ctx, status)
}

// -- Reports --

func (s *service) TrialBalance(ctx context.Context, year, month int) ([]AccountBalance, error) {
	period, err := s.periods.GetByYearMonth(ctx, year, month)
	if err != nil {
		return nil, err
	}
	balances, err := s.journals.GetTrialBalance(ctx, period.ID)
	if err != nil {
		return nil, err
	}
	for i := range balances {
		balances[i].Calculate()
	}
	return balances, nil
}

func (s *service) BalanceSheet(ctx context.Context, year, month int) ([]AccountBalance, error) {
	period, err := s.periods.GetByYearMonth(ctx, year, month)
	if err != nil {
		return nil, err
	}
	balances, err := s.journals.GetFinancialStatement(ctx, period.ID, []AccountType{
		AccountTypeAsset, AccountTypeLiability, AccountTypeEquity,
	})
	if err != nil {
		return nil, err
	}
	for i := range balances {
		balances[i].Calculate()
	}
	return balances, nil
}

func (s *service) IncomeStatement(ctx context.Context, year, month int) ([]AccountBalance, error) {
	period, err := s.periods.GetByYearMonth(ctx, year, month)
	if err != nil {
		return nil, err
	}
	balances, err := s.journals.GetFinancialStatement(ctx, period.ID, []AccountType{
		AccountTypeRevenue, AccountTypeExpense,
	})
	if err != nil {
		return nil, err
	}
	for i := range balances {
		balances[i].Calculate()
	}
	return balances, nil
}

// -- Periods --

func (s *service) CreatePeriod(ctx context.Context, period *Period) error {
	return s.periods.Create(ctx, period)
}

func (s *service) GetPeriod(ctx context.Context, id string) (*Period, error) {
	return s.periods.GetByID(ctx, id)
}

func (s *service) GetPeriodByYearMonth(ctx context.Context, year, month int) (*Period, error) {
	return s.periods.GetByYearMonth(ctx, year, month)
}

func (s *service) GetAllPeriods(ctx context.Context) ([]Period, error) {
	return s.periods.GetAll(ctx)
}

func (s *service) ClosePeriod(ctx context.Context, id string) error {
	period, err := s.periods.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if period.Status == PeriodClosed {
		return ErrPeriodAlreadyClosed
	}
	return s.periods.UpdateStatus(ctx, id, PeriodClosed)
}

func (s *service) ReopenPeriod(ctx context.Context, id string) error {
	period, err := s.periods.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if period.Status == PeriodOpen {
		return nil
	}
	entries, _ := s.journals.GetByPeriod(ctx, id)
	for _, e := range entries {
		if e.Status == JournalEntryPosted {
			return fmt.Errorf("period has posted entries: %w", ErrPeriodHasEntries)
		}
	}
	return s.periods.UpdateStatus(ctx, id, PeriodOpen)
}

// -- Multi-currency --

func (s *service) CreateExchangeRate(ctx context.Context, rate *ExchangeRate) error {
	if err := rate.Validate(); err != nil {
		return err
	}
	return s.rates.Create(ctx, rate)
}

func (s *service) GetExchangeRate(ctx context.Context, currencyCode string, rateDate time.Time) (*ExchangeRate, error) {
	return s.rates.GetByCurrencyAndDate(ctx, currencyCode, rateDate)
}

func (s *service) ListExchangeRates(ctx context.Context) ([]ExchangeRate, error) {
	return s.rates.GetAll(ctx)
}

// -- Audit --

func (s *service) GetAuditLog(ctx context.Context, limit int) ([]AuditEntry, error) {
	return s.audit.GetAll(ctx, limit)
}

func (s *service) GetAuditLogByEntity(ctx context.Context, entityType, entityID string) ([]AuditEntry, error) {
	return s.audit.GetByEntity(ctx, entityType, entityID)
}

// -- Users / Auth --

func (s *service) CreateUser(ctx context.Context, user *User, password string) error {
	if err := user.Validate(); err != nil {
		return err
	}
	if password == "" {
		return ErrPasswordRequired
	}
	existing, _ := s.users.GetByUsername(ctx, user.Username)
	if existing != nil {
		return ErrUsernameExists
	}
	hash, err := hashPassword(password)
	if err != nil {
		return fmt.Errorf("password hashing failed: %w", err)
	}
	user.PasswordHash = hash
	return s.users.Create(ctx, user)
}

func (s *service) GetUser(ctx context.Context, id string) (*User, error) {
	return s.users.GetByID(ctx, id)
}

func (s *service) ListUsers(ctx context.Context) ([]User, error) {
	return s.users.GetAll(ctx)
}

func (s *service) Authenticate(ctx context.Context, username, password string) (*User, error) {
	user, err := s.users.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if !user.IsActive {
		return nil, ErrUserInactive
	}
	if err := comparePassword(user.PasswordHash, password); err != nil {
		return nil, err
	}
	return user, nil
}
