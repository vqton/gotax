package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gotax/internal/auth"
	"gotax/internal/domain"

	"github.com/google/uuid"
)

type Service interface {
	CreateAccount(ctx context.Context, account *domain.Account) error
	GetAccount(ctx context.Context, code string) (*domain.Account, error)
	GetAllAccounts(ctx context.Context, activeOnly bool) ([]domain.Account, error)
	UpdateAccount(ctx context.Context, account *domain.Account) error
	DeleteAccount(ctx context.Context, code string) error
	FreezeAccount(ctx context.Context, code, reason string) error
	UnfreezeAccount(ctx context.Context, code, reason string) error

	CreateEntry(ctx context.Context, entry *domain.JournalEntry, userID string) error
	CreatePostedEntry(ctx context.Context, entry *domain.JournalEntry, userID string) error
	SubmitForReview(ctx context.Context, id, userID string) error
	ReviewEntry(ctx context.Context, id, reviewerID string) error
	ApproveEntry(ctx context.Context, id, approverID string) error
	PostEntry(ctx context.Context, id string) error
	CancelEntry(ctx context.Context, id string) error
	GetEntryByID(ctx context.Context, id string) (*domain.JournalEntry, error)
	GetEntriesByDateRange(ctx context.Context, from, to time.Time) ([]domain.JournalEntry, error)
	GetEntriesByStatus(ctx context.Context, status domain.JournalEntryStatus) ([]domain.JournalEntry, error)
	GetAllEntries(ctx context.Context) ([]domain.JournalEntry, error)

	TrialBalance(ctx context.Context, year, month int) ([]domain.AccountBalance, error)
	BalanceSheet(ctx context.Context, year, month int) ([]domain.AccountBalance, error)
	IncomeStatement(ctx context.Context, year, month int) ([]domain.AccountBalance, error)
	CashFlowStatement(ctx context.Context, companyID string, year, month int) (*CashFlowResult, error)

	CreatePeriod(ctx context.Context, period *domain.Period) error
	GetPeriod(ctx context.Context, id string) (*domain.Period, error)
	GetPeriodByYearMonth(ctx context.Context, year, month int) (*domain.Period, error)
	GetAllPeriods(ctx context.Context) ([]domain.Period, error)
	ClosePeriod(ctx context.Context, id string) error
	ReopenPeriod(ctx context.Context, id string) error

	CreateExchangeRate(ctx context.Context, rate *domain.ExchangeRate) error
	GetExchangeRate(ctx context.Context, currencyCode string, rateDate time.Time) (*domain.ExchangeRate, error)
	ListExchangeRates(ctx context.Context) ([]domain.ExchangeRate, error)

	GetAuditLog(ctx context.Context, limit int) ([]domain.AuditEntry, error)
	GetAuditLogByEntity(ctx context.Context, entityType, entityID string) ([]domain.AuditEntry, error)

	CreateUser(ctx context.Context, user *domain.User, password string) error
	GetUser(ctx context.Context, id string) (*domain.User, error)
	ListUsers(ctx context.Context) ([]domain.User, error)

	Login(ctx context.Context, username, password, ip string) (*domain.AuthResult, error)
	Verify2FA(ctx context.Context, tempToken, code string) (*domain.AuthResult, error)
	RefreshToken(ctx context.Context, refreshTokenStr string) (*domain.TokenPair, error)
	Logout(ctx context.Context, userID, refreshTokenStr string) error
	LogoutAll(ctx context.Context, userID string) error

	AdminResetPassword(ctx context.Context, userID, newPassword string) error
	UpdateUserRole(ctx context.Context, userID string, role domain.UserRole) error
	ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error

	SetupTOTP(ctx context.Context, userID string) (*domain.TOTPSetup, error)
	ConfirmTOTP(ctx context.Context, userID, code string) error
	DisableTOTP(ctx context.Context, userID, currentPassword, code string) error
	GenerateBackupCodes(ctx context.Context, userID string) ([]string, error)

	ListSessions(ctx context.Context, userID string) ([]domain.Session, error)
	RevokeSession(ctx context.Context, userID, sessionID string) error

	GetAccountBalance(ctx context.Context, code string, periodID string) (*domain.AccountBalance, error)
	GetAccountBalanceDrillDown(ctx context.Context, code string, periodID string) ([]domain.JournalEntry, error)
	GetAccountUsage(ctx context.Context, code string) (*domain.AccountUsage, error)

	CreateApprovalRequest(ctx context.Context, req *domain.ApprovalRequest) error
	ApproveRequest(ctx context.Context, id, reviewerID, note string) error
	RejectRequest(ctx context.Context, id, reviewerID, note string) error
	GetApprovalRequests(ctx context.Context, status domain.ApprovalStatus) ([]domain.ApprovalRequest, error)

	CreateAccountVersion(ctx context.Context, reason string) (*domain.AccountVersion, error)
	GetVersion(ctx context.Context, versionNumber string) (*domain.AccountVersion, error)
	ListVersions(ctx context.Context) ([]domain.AccountVersion, error)
	CompareVersions(ctx context.Context, v1, v2 string) (*domain.VersionDiff, error)

	CreateAccountAnalysis(ctx context.Context, analysis *domain.AccountAnalysis) error
	GetAccountAnalysis(ctx context.Context, accountCode string) (*domain.AccountAnalysis, error)
	UpdateAccountAnalysis(ctx context.Context, analysis *domain.AccountAnalysis) error

	CreateAccountMapping(ctx context.Context, mapping *domain.AccountMapping) error
	GetMappingByOldCode(ctx context.Context, sourceRegime, oldCode string) (*domain.AccountMapping, error)
	ListMappings(ctx context.Context, sourceRegime, targetRegime string) ([]domain.AccountMapping, error)

	CreateIFRSMapping(ctx context.Context, mapping *domain.IFRSMapping) error
	GetIFRSMapping(ctx context.Context, vasCode string) (*domain.IFRSMapping, error)
	ListIFRSMappings(ctx context.Context) ([]domain.IFRSMapping, error)

	// Opening Balance
	CreateOpeningBalance(ctx context.Context, ob *domain.OpeningBalance) error
	GetOpeningBalance(ctx context.Context, id string) (*domain.OpeningBalance, error)
	ListOpeningBalances(ctx context.Context, filter domain.OBListFilter) ([]domain.OpeningBalance, error)
	UpdateOpeningBalance(ctx context.Context, ob *domain.OpeningBalance) error
	SubmitOpeningBalance(ctx context.Context, id, userID string) error
	ApproveOpeningBalance(ctx context.Context, id, approverID string) error
	CorrectOpeningBalance(ctx context.Context, id, correctedBy, reason string) (*domain.OpeningBalance, error)
	DeleteOpeningBalance(ctx context.Context, id string) error

	BulkCreateOpeningBalances(ctx context.Context, balances []domain.OpeningBalance) error
	BulkSubmitOpeningBalances(ctx context.Context, ids []string, userID string) error
	BulkApproveOpeningBalances(ctx context.Context, ids []string, approverID string) error

	CreateOpeningBalanceDetail(ctx context.Context, d *domain.OpeningBalanceDetail) error
	GetOpeningBalanceDetails(ctx context.Context, balanceID string) ([]domain.OpeningBalanceDetail, error)
	DeleteOpeningBalanceDetail(ctx context.Context, id string) error

	GetOpeningBalanceTotals(ctx context.Context, companyID, periodID string) (debit, credit float64, err error)
	ValidateOpeningBalancesBalanced(ctx context.Context, companyID, periodID string) (bool, error)

	CarryForward(ctx context.Context, companyID, fromPeriodID, toPeriodID, fromFiscalYear, toFiscalYear, executedBy string) (*domain.CarryForwardLog, error)
	GetCarryForwardLogs(ctx context.Context, companyID string) ([]domain.CarryForwardLog, error)
	GetCarryForwardLogByID(ctx context.Context, id string) (*domain.CarryForwardLog, error)

	YearEndClose(ctx context.Context, companyID, fromPeriodID, toPeriodID, fromYear, toYear, userID string) (*YearEndCloseResult, error)
	ExportYearEndBalances(ctx context.Context, companyID, periodID string) ([]domain.OpeningBalance, error)

	CreateCircular99Mapping(ctx context.Context, m *domain.Circular99Mapping) error
	ListCircular99Mappings(ctx context.Context) ([]domain.Circular99Mapping, error)
	GetCircular99MappingByOldCode(ctx context.Context, oldCode string) (*domain.Circular99Mapping, error)

	CreateBalanceMigration(ctx context.Context, m *domain.BalanceMigration) error
	GetBalanceMigrationByID(ctx context.Context, id string) (*domain.BalanceMigration, error)
	ListBalanceMigrations(ctx context.Context, companyID string) ([]domain.BalanceMigration, error)

	ImportOpeningBalances(ctx context.Context, data []byte, companyID, periodID, createdBy string) (*OBImportResult, error)
	GenerateOpeningBalancePDF(ctx context.Context, companyID, periodID string) ([]byte, error)
	GeneratePrintPDF(ctx context.Context, printType, id string) ([]byte, error)

	// Cash Module
	CreateCashReceipt(ctx context.Context, r *domain.CashReceipt) error
	GetCashReceipt(ctx context.Context, id string) (*domain.CashReceipt, error)
	ListCashReceipts(ctx context.Context, filter domain.CashReceiptFilter) ([]domain.CashReceipt, int, error)
	UpdateCashReceipt(ctx context.Context, r *domain.CashReceipt) error
	DeleteCashReceipt(ctx context.Context, id string) error
	SubmitCashReceipt(ctx context.Context, id, userID string) error
	ApproveCashReceipt(ctx context.Context, id, approverID string) error
	RejectCashReceipt(ctx context.Context, id, reviewerID string) error
	PostCashReceipt(ctx context.Context, id, userID string) error

	CreateCashPayment(ctx context.Context, p *domain.CashPayment) error
	GetCashPayment(ctx context.Context, id string) (*domain.CashPayment, error)
	ListCashPayments(ctx context.Context, filter domain.CashPaymentFilter) ([]domain.CashPayment, int, error)
	UpdateCashPayment(ctx context.Context, p *domain.CashPayment) error
	DeleteCashPayment(ctx context.Context, id string) error
	SubmitCashPayment(ctx context.Context, id, userID string) error
	ApproveCashPayment(ctx context.Context, id, approverID string) error
	RejectCashPayment(ctx context.Context, id, reviewerID string) error
	PostCashPayment(ctx context.Context, id, userID string) error

	CreateCashTransfer(ctx context.Context, t *domain.CashTransfer, userID string) error
	GetCashTransfers(ctx context.Context, companyID string) ([]domain.CashTransfer, error)

	GetCashBook(ctx context.Context, companyID, currency, accountID, fromDate, toDate string) (*domain.CashBook, error)
	GetCashBalance(ctx context.Context, companyID, accountID string) (float64, error)
	GetCashFlowStatement(ctx context.Context, companyID, currency, accountID, fromDate, toDate string) (*CashFlowStatement, error)

	CreatePettyCashFund(ctx context.Context, f *domain.PettyCashFund) error
	GetPettyCashFund(ctx context.Context, id string) (*domain.PettyCashFund, error)
	ListPettyCashFunds(ctx context.Context, companyID string) ([]domain.PettyCashFund, error)

	CreateCashInventory(ctx context.Context, inv *domain.CashInventory) error
	GetCashInventory(ctx context.Context, id string) (*domain.CashInventory, error)
	ListCashInventories(ctx context.Context, companyID string) ([]domain.CashInventory, error)

	CreateAdvance(ctx context.Context, a *domain.AdvanceRequest) error
	GetAdvance(ctx context.Context, id string) (*domain.AdvanceRequest, error)
	ListAdvances(ctx context.Context, companyID string) ([]domain.AdvanceRequest, error)
	UpdateAdvance(ctx context.Context, a *domain.AdvanceRequest) error
	ApproveAdvance(ctx context.Context, id, approverID string) error
	RejectAdvance(ctx context.Context, id, reviewerID string) error
	PayAdvance(ctx context.Context, id, paidBy string) error
	SettleAdvance(ctx context.Context, id, settlementID string) error
	CreateAdvanceSettlement(ctx context.Context, s *domain.AdvanceSettlement) error
	ListAdvanceByStatus(ctx context.Context, companyID string, status domain.AdvanceStatus) ([]domain.AdvanceRequest, error)
}

type service struct {
	accounts  domain.AccountRepository
	journals  domain.JournalRepository
	periods   domain.PeriodRepository
	users     domain.UserRepository
	audit     domain.AuditLogRepository
	rates     domain.ExchangeRateRepository
	templates domain.ClosingTemplateRepository
	approvals domain.ApprovalRepository
	versions  domain.AccountVersionRepository
	mappings  domain.AccountMappingRepository
	analysis  domain.AccountAnalysisRepository
	ifrs      domain.IFRSMappingRepository
	refresh   domain.RefreshTokenRepository
	reset     domain.PasswordResetTokenRepository
	ob        domain.OpeningBalanceRepository
	cash      domain.CashRepository
	limiter   *auth.RateLimiter
	now       func() time.Time
}

func NewService(
	accRepo domain.AccountRepository,
	jeRepo domain.JournalRepository,
	perRepo domain.PeriodRepository,
	userRepo domain.UserRepository,
	auditRepo domain.AuditLogRepository,
	rateRepo domain.ExchangeRateRepository,
	templateRepo domain.ClosingTemplateRepository,
	approvalRepo domain.ApprovalRepository,
	versionRepo domain.AccountVersionRepository,
	mappingRepo domain.AccountMappingRepository,
	analysisRepo domain.AccountAnalysisRepository,
	ifrsRepo domain.IFRSMappingRepository,
	refreshRepo domain.RefreshTokenRepository,
	resetRepo domain.PasswordResetTokenRepository,
	obRepo domain.OpeningBalanceRepository,
	cashRepo domain.CashRepository,
) Service {
	return &service{
		accounts:  accRepo,
		journals:  jeRepo,
		periods:   perRepo,
		users:     userRepo,
		audit:     auditRepo,
		rates:     rateRepo,
		templates: templateRepo,
		approvals: approvalRepo,
		versions:  versionRepo,
		mappings:  mappingRepo,
		analysis:  analysisRepo,
		ifrs:      ifrsRepo,
		refresh:   refreshRepo,
		reset:     resetRepo,
		ob:        obRepo,
		cash:      cashRepo,
		limiter:   auth.NewRateLimiter(5, 15*time.Minute),
		now:       time.Now,
	}
}

// ─── Accounts ──────────────────────────────────────────────────────

func (s *service) CreateAccount(ctx context.Context, account *domain.Account) error {
	if err := account.Validate(); err != nil {
		return err
	}
	existing, err := s.accounts.GetByCode(ctx, account.Code)
	if err != nil && !errors.Is(err, domain.ErrAccountNotFound) {
		return err
	}
	if existing != nil {
		return domain.ErrAccountCodeExists
	}
	if account.ParentCode != "" {
		if _, err := s.accounts.GetByCode(ctx, account.ParentCode); err != nil {
			return err
		}
	}
	return s.accounts.Create(ctx, account)
}

func (s *service) GetAccount(ctx context.Context, code string) (*domain.Account, error) {
	return s.accounts.GetByCode(ctx, code)
}

func (s *service) GetAllAccounts(ctx context.Context, activeOnly bool) ([]domain.Account, error) {
	return s.accounts.GetAll(ctx, activeOnly)
}

func (s *service) UpdateAccount(ctx context.Context, account *domain.Account) error {
	existing, err := s.accounts.GetByCode(ctx, account.Code)
	if err != nil {
		return err
	}
	if existing == nil {
		return domain.ErrAccountNotFound
	}
	if err := account.Validate(); err != nil {
		return err
	}
	// PUT payload carries editable fields only; lifecycle fields are
	// managed by freeze/unfreeze and must survive a plain edit.
	account.IsActive = existing.IsActive
	account.Status = existing.Status
	account.FreezeReason = existing.FreezeReason
	return s.accounts.Update(ctx, account)
}

func (s *service) DeleteAccount(ctx context.Context, code string) error {
	if _, err := s.accounts.GetByCode(ctx, code); err != nil {
		return err
	}
	children, err := s.accounts.GetChildren(ctx, code)
	if err != nil {
		return err
	}
	if len(children) > 0 {
		return domain.ErrAccountHasChildren
	}
	// Block deletion of accounts used in posted journal entries: PG would
	// hit the FK constraint (raw 500), memory backend would leave dangling
	// line references. TT99 keeps accounts with balance/usage open — freeze instead.
	usage, err := s.journals.GetAccountUsage(ctx, code)
	if err != nil {
		return err
	}
	if usage.EntryCount > 0 {
		return domain.ErrAccountHasBalance
	}
	return s.accounts.Delete(ctx, code)
}

// ─── Journal Entries ───────────────────────────────────────────────

func (s *service) CreateEntry(ctx context.Context, entry *domain.JournalEntry, userID string) error {
	if err := entry.Validate(); err != nil {
		return err
	}
	for _, line := range entry.Lines {
		acc, err := s.accounts.GetByCode(ctx, line.AccountCode)
		if err != nil {
			return err
		}
		if err := acc.CanPost(); err != nil {
			return err
		}
	}
	entry.Status = domain.JournalEntryDraft
	entry.CreatedBy = userID
	entry.CurrencyCode = "VND"
	entry.ExchangeRate = 1
	if entry.VoucherType == "" {
		entry.VoucherType = domain.VoucherTypeOther
	}
	if entry.AccountingDate.IsZero() {
		entry.AccountingDate = entry.EntryDate
	}
	// Attach the open period and assign a sequential number when creating
	// manually (API/UI drafts). Entries without an open period stay unnumbered.
	if entry.PeriodID == "" {
		if p, err := s.periods.GetOpenPeriod(ctx); err == nil && p != nil {
			entry.PeriodID = p.ID
		}
	}
	if entry.EntryNumber == "" && entry.PeriodID != "" {
		if p, err := s.periods.GetByID(ctx, entry.PeriodID); err == nil && p != nil {
			if es, err := s.journals.GetByPeriod(ctx, p.ID); err == nil {
				entry.EntryNumber = fmt.Sprintf("%d%02d-%03d", p.Year, p.Month, len(es)+1)
			}
		}
	}
	return s.journals.Create(ctx, entry)
}

func (s *service) CreatePostedEntry(ctx context.Context, entry *domain.JournalEntry, userID string) error {
	if err := entry.Validate(); err != nil {
		return err
	}
	for _, line := range entry.Lines {
		acc, err := s.accounts.GetByCode(ctx, line.AccountCode)
		if err != nil {
			return err
		}
		if err := acc.CanPost(); err != nil {
			return err
		}
	}
	entry.Status = domain.JournalEntryPosted
	entry.CreatedBy = userID
	now := s.now()
	entry.PostedAt = &now
	entry.CurrencyCode = "VND"
	entry.ExchangeRate = 1
	if entry.VoucherType == "" {
		entry.VoucherType = domain.VoucherTypeOther
	}
	if entry.AccountingDate.IsZero() {
		entry.AccountingDate = entry.EntryDate
	}
	if entry.PeriodID == "" {
		p, err := s.periods.GetByYearMonth(ctx, entry.EntryDate.Year(), int(entry.EntryDate.Month()))
		if err == nil && p != nil {
			entry.PeriodID = p.ID
		}
	}
	return s.journals.Create(ctx, entry)
}

func (s *service) SubmitForReview(ctx context.Context, id, userID string) error {
	entry, err := s.journals.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if entry.Status != domain.JournalEntryDraft {
		return fmt.Errorf("entry must be DRAFT to submit for review")
	}
	return s.journals.Review(ctx, id, userID)
}

func (s *service) ReviewEntry(ctx context.Context, id, reviewerID string) error {
	entry, err := s.journals.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if entry.Status == domain.JournalEntryReviewing {
		return domain.ErrJournalAlreadyReviewed
	}
	if entry.Status != domain.JournalEntryDraft {
		return fmt.Errorf("entry must be DRAFT for review")
	}
	return s.journals.Review(ctx, id, reviewerID)
}

func (s *service) ApproveEntry(ctx context.Context, id, approverID string) error {
	entry, err := s.journals.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if entry.Status == domain.JournalEntryApproved {
		return domain.ErrJournalAlreadyApproved
	}
	if entry.Status != domain.JournalEntryReviewing {
		return fmt.Errorf("entry must be REVIEWING to approve")
	}
	return s.journals.Approve(ctx, id, approverID)
}

func (s *service) PostEntry(ctx context.Context, id string) error {
	entry, err := s.journals.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if entry.Status == domain.JournalEntryPosted {
		return domain.ErrJournalAlreadyPosted
	}
	if entry.Status == domain.JournalEntryCancelled {
		return domain.ErrJournalAlreadyCancelled
	}
	if entry.Status != domain.JournalEntryApproved {
		return fmt.Errorf("entry must be APPROVED to post")
	}
	period, err := s.periods.GetByID(ctx, entry.PeriodID)
	if err != nil {
		return err
	}
	if period == nil {
		return domain.ErrPeriodNotFound
	}
	if period.Status != domain.PeriodOpen {
		return domain.ErrJournalPeriodClosed
	}
	entry.Status = domain.JournalEntryPosted
	now := s.now()
	entry.PostedAt = &now
	return s.journals.Update(ctx, entry)
}

func (s *service) CancelEntry(ctx context.Context, id string) error {
	entry, err := s.journals.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if entry.Status == domain.JournalEntryCancelled {
		return domain.ErrJournalAlreadyCancelled
	}
	if entry.Status == domain.JournalEntryDraft {
		return domain.ErrJournalAlreadyDraft
	}
	return s.journals.UpdateStatus(ctx, id, domain.JournalEntryCancelled)
}

func (s *service) GetEntryByID(ctx context.Context, id string) (*domain.JournalEntry, error) {
	return s.journals.GetByID(ctx, id)
}

func (s *service) GetEntriesByDateRange(ctx context.Context, from, to time.Time) ([]domain.JournalEntry, error) {
	return s.journals.GetByDateRange(ctx, from, to)
}

func (s *service) GetEntriesByStatus(ctx context.Context, status domain.JournalEntryStatus) ([]domain.JournalEntry, error) {
	return s.journals.GetByStatus(ctx, status)
}

func (s *service) GetAllEntries(ctx context.Context) ([]domain.JournalEntry, error) {
	return s.journals.GetAll(ctx)
}

// ─── Reports ───────────────────────────────────────────────────────

func (s *service) TrialBalance(ctx context.Context, year, month int) ([]domain.AccountBalance, error) {
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

func (s *service) BalanceSheet(ctx context.Context, year, month int) ([]domain.AccountBalance, error) {
	period, err := s.periods.GetByYearMonth(ctx, year, month)
	if err != nil {
		return nil, err
	}
	balances, err := s.journals.GetFinancialStatement(ctx, period.ID, []domain.AccountType{
		domain.AccountTypeAsset, domain.AccountTypeLiability, domain.AccountTypeEquity,
	})
	if err != nil {
		return nil, err
	}
	for i := range balances {
		balances[i].Calculate()
	}
	return balances, nil
}

func (s *service) IncomeStatement(ctx context.Context, year, month int) ([]domain.AccountBalance, error) {
	period, err := s.periods.GetByYearMonth(ctx, year, month)
	if err != nil {
		return nil, err
	}
	balances, err := s.journals.GetFinancialStatement(ctx, period.ID, []domain.AccountType{
		domain.AccountTypeRevenue, domain.AccountTypeExpense,
	})
	if err != nil {
		return nil, err
	}
	for i := range balances {
		balances[i].Calculate()
	}
	return balances, nil
}

// ─── Periods ───────────────────────────────────────────────────────

func (s *service) CreatePeriod(ctx context.Context, period *domain.Period) error {
	if err := period.Validate(); err != nil {
		return err
	}
	if period.ID == "" {
		// PG periods.id is VARCHAR(20) — short deterministic ID, unique per (year, month).
		period.ID = fmt.Sprintf("P-%04d%02d", period.Year, period.Month)
	}
	return s.periods.Create(ctx, period)
}

func (s *service) GetPeriod(ctx context.Context, id string) (*domain.Period, error) {
	return s.periods.GetByID(ctx, id)
}

func (s *service) GetPeriodByYearMonth(ctx context.Context, year, month int) (*domain.Period, error) {
	return s.periods.GetByYearMonth(ctx, year, month)
}

func (s *service) GetAllPeriods(ctx context.Context) ([]domain.Period, error) {
	return s.periods.GetAll(ctx)
}

func (s *service) ClosePeriod(ctx context.Context, id string) error {
	period, err := s.periods.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if period.Status == domain.PeriodClosed {
		return domain.ErrPeriodAlreadyClosed
	}
	return s.periods.UpdateStatus(ctx, id, domain.PeriodClosed)
}

func (s *service) ReopenPeriod(ctx context.Context, id string) error {
	period, err := s.periods.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if period.Status == domain.PeriodOpen {
		return nil
	}
	entries, err := s.journals.GetByPeriod(ctx, id)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.Status == domain.JournalEntryPosted {
			return fmt.Errorf("period has posted entries: %w", domain.ErrPeriodHasEntries)
		}
	}
	return s.periods.UpdateStatus(ctx, id, domain.PeriodOpen)
}

// ─── Multi-currency ────────────────────────────────────────────────

func (s *service) CreateExchangeRate(ctx context.Context, rate *domain.ExchangeRate) error {
	if err := rate.Validate(); err != nil {
		return err
	}
	if rate.ID == "" {
		rate.ID = uuid.NewString()
	}
	return s.rates.Create(ctx, rate)
}

func (s *service) GetExchangeRate(ctx context.Context, currencyCode string, rateDate time.Time) (*domain.ExchangeRate, error) {
	return s.rates.GetByCurrencyAndDate(ctx, currencyCode, rateDate)
}

func (s *service) ListExchangeRates(ctx context.Context) ([]domain.ExchangeRate, error) {
	return s.rates.GetAll(ctx)
}

// ─── Audit ─────────────────────────────────────────────────────────

func (s *service) GetAuditLog(ctx context.Context, limit int) ([]domain.AuditEntry, error) {
	return s.audit.GetAll(ctx, limit)
}

func (s *service) GetAuditLogByEntity(ctx context.Context, entityType, entityID string) ([]domain.AuditEntry, error) {
	return s.audit.GetByEntity(ctx, entityType, entityID)
}

// ─── Users / Auth ──────────────────────────────────────────────────

func (s *service) CreateUser(ctx context.Context, user *domain.User, password string) error {
	if err := user.Validate(); err != nil {
		return err
	}
	if password == "" {
		return domain.ErrPasswordRequired
	}
	existing, err := s.users.GetByUsername(ctx, user.Username)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return err
	}
	if existing != nil {
		return domain.ErrUsernameExists
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		return fmt.Errorf("password hashing failed: %w", err)
	}
	user.PasswordHash = hash
	return s.users.Create(ctx, user)
}

func (s *service) GetUser(ctx context.Context, id string) (*domain.User, error) {
	return s.users.GetByID(ctx, id)
}

func (s *service) ListUsers(ctx context.Context) ([]domain.User, error) {
	return s.users.GetAll(ctx)
}

// ─── Full Auth Flow ────────────────────────────────────────────────

func (s *service) Login(ctx context.Context, username, password, ip string) (*domain.AuthResult, error) {
	if !s.limiter.Allow("login:" + username) {
		return nil, domain.ErrRateLimited
	}

	user, err := s.users.GetByUsername(ctx, username)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if user.IsLocked() {
		return nil, domain.ErrAccountLocked
	}

	if !user.IsActive {
		return nil, domain.ErrAccountInactive
	}

	if err := auth.ComparePassword(user.PasswordHash, password); err != nil {
		user.FailedAttempts++
		if user.FailedAttempts >= domain.MaxLoginAttempts {
			until := s.now().Add(domain.LockoutDuration)
			user.LockedUntil = &until
		}
		s.users.Update(ctx, user)
		return nil, domain.ErrInvalidCredentials
	}

	user.FailedAttempts = 0
	user.LockedUntil = nil
	user.LastLoginAt = &[]time.Time{s.now()}[0]
	user.LastLoginIP = ip
	s.users.Update(ctx, user)

	if user.TOTPEnabled {
		tempToken, err := auth.GenerateRefreshTokenRaw(user.ID)
		if err != nil {
			return nil, err
		}
		return &domain.AuthResult{
			Requires2FA: true,
			TempToken:   tempToken,
			User:        user,
		}, nil
	}

	accessToken, err := auth.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	rawRefresh, err := auth.GenerateRefreshTokenRaw(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken := &domain.RefreshToken{
		UserID:    user.ID,
		TokenHash: auth.HashRefreshToken(rawRefresh),
		IPAddress: ip,
		ExpiresAt: s.now().Add(7 * 24 * time.Hour),
	}
	s.refresh.Create(ctx, refreshToken)

	result := &domain.AuthResult{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
		ExpiresIn:    900,
		TokenType:    "Bearer",
		User:         user,
	}

	if user.IsPasswordExpired() {
		result.PasswordExpired = true
	}

	return result, nil
}

func (s *service) Verify2FA(ctx context.Context, tempToken, code string) (*domain.AuthResult, error) {
	userID := tempToken
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if !auth.VerifyTOTP(user.TOTPSecret, code) {
		return nil, domain.ErrInvalid2FACode
	}

	accessToken, err := auth.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	rawRefresh, err := auth.GenerateRefreshTokenRaw(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken := &domain.RefreshToken{
		UserID:    user.ID,
		TokenHash: auth.HashRefreshToken(rawRefresh),
		ExpiresAt: s.now().Add(7 * 24 * time.Hour),
	}
	s.refresh.Create(ctx, refreshToken)

	return &domain.AuthResult{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
		ExpiresIn:    900,
		TokenType:    "Bearer",
		User:         user,
	}, nil
}

func (s *service) RefreshToken(ctx context.Context, rawRefresh string) (*domain.TokenPair, error) {
	hash := auth.HashRefreshToken(rawRefresh)
	found, err := s.refresh.GetByHash(ctx, hash)
	if err != nil {
		return nil, domain.ErrInvalidRefreshToken
	}
	if found.RevokedAt != nil {
		return nil, domain.ErrRefreshTokenRevoked
	}
	if s.now().After(found.ExpiresAt) {
		return nil, domain.ErrRefreshTokenExpired
	}
	user, err := s.users.GetByID(ctx, found.UserID)
	if err != nil {
		return nil, err
	}
	accessToken, err := auth.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}
	rawNew, err := auth.GenerateRefreshTokenRaw(user.ID)
	if err != nil {
		return nil, err
	}
	refreshToken := &domain.RefreshToken{
		UserID:    user.ID,
		TokenHash: auth.HashRefreshToken(rawNew),
		IPAddress: found.IPAddress,
		ExpiresAt: s.now().Add(7 * 24 * time.Hour),
	}
	if err := s.refresh.Create(ctx, refreshToken); err != nil {
		return nil, err
	}
	if err := s.refresh.Revoke(ctx, found.ID); err != nil {
		return nil, err
	}
	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: rawNew,
		ExpiresIn:    900,
		TokenType:    "Bearer",
	}, nil
}

func (s *service) Logout(ctx context.Context, userID, rawRefresh string) error {
	if rawRefresh == "" {
		return nil
	}
	hash := auth.HashRefreshToken(rawRefresh)
	t, err := s.refresh.GetByHash(ctx, hash)
	if err != nil {
		return nil
	}
	if t.UserID != userID {
		return nil
	}
	return s.refresh.Revoke(ctx, t.ID)
}

func (s *service) LogoutAll(ctx context.Context, userID string) error {
	return s.refresh.RevokeAllByUserID(ctx, userID)
}

// ─── Password Management ──────────────────────────────────────────

func (s *service) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := auth.ComparePassword(user.PasswordHash, currentPassword); err != nil {
		return domain.ErrInvalidCurrentPassword
	}

	if err := domain.ValidatePassword(newPassword); err != nil {
		return err
	}

	newHash, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if domain.IsPasswordInHistory(user.PasswordHistory, newHash) {
		return domain.ErrPasswordReuse
	}

	user.PasswordHistory = append(user.PasswordHistory, user.PasswordHash)
	if len(user.PasswordHistory) > domain.PasswordHistorySize {
		user.PasswordHistory = user.PasswordHistory[len(user.PasswordHistory)-domain.PasswordHistorySize:]
	}

	user.PasswordHash = newHash
	now := s.now()
	user.PasswordChangedAt = &now
	user.FailedAttempts = 0
	user.LockedUntil = nil

	return s.users.Update(ctx, user)
}

func (s *service) ForgotPassword(ctx context.Context, email string) error {
	users, err := s.users.GetAll(ctx)
	if err != nil {
		return err
	}

	for _, u := range users {
		if u.Email == email {
			raw, err := auth.GeneratePasswordResetTokenRaw()
			if err != nil {
				return err
			}
			token := &domain.PasswordResetToken{
				ID:        raw,
				UserID:    u.ID,
				TokenHash: auth.HashRefreshToken(raw),
				ExpiresAt: s.now().Add(1 * time.Hour),
			}
			return s.reset.Create(ctx, token)
		}
	}
	return nil
}

func (s *service) ResetPassword(ctx context.Context, rawToken, newPassword string) error {
	tokens, err := s.reset.GetByID(ctx, rawToken)
	if err != nil {
		return domain.ErrInvalidPasswordResetToken
	}

	if tokens.UsedAt != nil {
		return domain.ErrInvalidPasswordResetToken
	}

	if s.now().After(tokens.ExpiresAt) {
		return domain.ErrPasswordResetTokenExpired
	}

	if err := domain.ValidatePassword(newPassword); err != nil {
		return err
	}

	user, err := s.users.GetByID(ctx, tokens.UserID)
	if err != nil {
		return err
	}

	newHash, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if domain.IsPasswordInHistory(user.PasswordHistory, newHash) {
		return domain.ErrPasswordReuse
	}

	user.PasswordHistory = append(user.PasswordHistory, user.PasswordHash)
	if len(user.PasswordHistory) > domain.PasswordHistorySize {
		user.PasswordHistory = user.PasswordHistory[len(user.PasswordHistory)-domain.PasswordHistorySize:]
	}
	user.PasswordHash = newHash
	now := s.now()
	user.PasswordChangedAt = &now
	user.FailedAttempts = 0
	user.LockedUntil = nil

	if err := s.users.Update(ctx, user); err != nil {
		return err
	}

	return s.reset.MarkUsed(ctx, tokens.ID)
}

// ─── TOTP ─────────────────────────────────────────────────────────

func (s *service) SetupTOTP(ctx context.Context, userID string) (*domain.TOTPSetup, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.TOTPEnabled {
		return nil, domain.Err2FAAlreadyEnabled
	}

	secret := auth.GenerateTOTPSecret()
	qrURL := auth.GenerateTOTPURL(user.Username, secret)

	backupCodes := make([]string, 8)
	for i := range backupCodes {
		b, _ := auth.GenerateRefreshTokenRaw(userID)
		backupCodes[i] = b[:8]
	}

	user.TOTPSecret = secret
	user.BackupCodes = backupCodes
	s.users.Update(ctx, user)

	return &domain.TOTPSetup{
		Secret:      secret,
		QRCodeURL:   qrURL,
		BackupCodes: backupCodes,
	}, nil
}

func (s *service) ConfirmTOTP(ctx context.Context, userID, code string) error {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.TOTPEnabled {
		return domain.Err2FAAlreadyEnabled
	}

	if user.TOTPSecret == "" {
		return domain.Err2FANotSetup
	}

	if !auth.VerifyTOTP(user.TOTPSecret, code) {
		return domain.ErrInvalid2FACode
	}

	user.TOTPEnabled = true
	return s.users.Update(ctx, user)
}

func (s *service) DisableTOTP(ctx context.Context, userID, currentPassword, code string) error {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if !user.TOTPEnabled {
		return domain.Err2FANotSetup
	}

	if err := auth.ComparePassword(user.PasswordHash, currentPassword); err != nil {
		return domain.ErrInvalidCurrentPassword
	}

	if !auth.VerifyTOTP(user.TOTPSecret, code) {
		if !auth.CheckBackupCode(user, code) {
			return domain.ErrInvalid2FACode
		}
	}

	user.TOTPEnabled = false
	user.TOTPSecret = ""
	user.BackupCodes = nil
	return s.users.Update(ctx, user)
}

func (s *service) GenerateBackupCodes(ctx context.Context, userID string) ([]string, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	backupCodes := make([]string, 8)
	for i := range backupCodes {
		b, _ := auth.GenerateRefreshTokenRaw(userID)
		backupCodes[i] = b[:8]
	}

	user.BackupCodes = backupCodes
	s.users.Update(ctx, user)
	return backupCodes, nil
}

// ─── Sessions ──────────────────────────────────────────────────────

func (s *service) ListSessions(ctx context.Context, userID string) ([]domain.Session, error) {
	tokens, err := s.refresh.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var sessions []domain.Session
	for _, t := range tokens {
		isCurrent := t.RevokedAt == nil && s.now().Before(t.ExpiresAt)
		sessions = append(sessions, domain.Session{
			ID:           t.ID,
			Device:       t.DeviceInfo,
			IP:           t.IPAddress,
			CreatedAt:    t.CreatedAt,
			LastActivity: t.CreatedAt,
			IsCurrent:    isCurrent,
		})
	}
	return sessions, nil
}

func (s *service) RevokeSession(ctx context.Context, userID, sessionID string) error {
	token, err := s.refresh.GetByID(ctx, sessionID)
	if err != nil {
		return err
	}
	if token.UserID != userID {
		return domain.ErrInvalidRefreshToken
	}
	return s.refresh.Revoke(ctx, sessionID)
}

// ─── COA: Freeze ───────────────────────────────────────────────────

func (s *service) FreezeAccount(ctx context.Context, code, reason string) error {
	acc, err := s.accounts.GetByCode(ctx, code)
	if err != nil {
		return err
	}
	if err := acc.Freeze(reason); err != nil {
		return err
	}
	return s.accounts.Update(ctx, acc)
}

func (s *service) UnfreezeAccount(ctx context.Context, code, reason string) error {
	acc, err := s.accounts.GetByCode(ctx, code)
	if err != nil {
		return err
	}
	if err := acc.Unfreeze(reason); err != nil {
		return err
	}
	return s.accounts.Update(ctx, acc)
}

// ─── COA: Balance ──────────────────────────────────────────────────

func (s *service) GetAccountBalance(ctx context.Context, code string, periodID string) (*domain.AccountBalance, error) {
	if _, err := s.accounts.GetByCode(ctx, code); err != nil {
		return nil, err
	}
	return s.journals.GetBalance(ctx, code, periodID)
}

func (s *service) GetAccountBalanceDrillDown(ctx context.Context, code string, periodID string) ([]domain.JournalEntry, error) {
	if _, err := s.accounts.GetByCode(ctx, code); err != nil {
		return nil, err
	}
	return s.journals.GetPostedEntriesByAccount(ctx, periodID, code)
}

func (s *service) GetAccountUsage(ctx context.Context, code string) (*domain.AccountUsage, error) {
	if _, err := s.accounts.GetByCode(ctx, code); err != nil {
		return nil, err
	}
	return s.journals.GetAccountUsage(ctx, code)
}

// ─── COA: Approval ─────────────────────────────────────────────────

func (s *service) CreateApprovalRequest(ctx context.Context, req *domain.ApprovalRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}
	return s.approvals.Create(ctx, req)
}

func (s *service) ApproveRequest(ctx context.Context, id, reviewerID, note string) error {
	req, err := s.approvals.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if req.Status.IsTerminal() {
		return domain.ErrApprovalAlreadyProcessed
	}
	if req.RequestedBy == reviewerID {
		return domain.ErrSelfApproval
	}
	if !req.ExpiresAt.IsZero() && s.now().After(req.ExpiresAt) {
		s.approvals.UpdateStatus(ctx, id, domain.ApprovalExpired, "", "")
		return domain.ErrApprovalExpired
	}
	return s.approvals.UpdateStatus(ctx, id, domain.ApprovalApproved, reviewerID, note)
}

func (s *service) RejectRequest(ctx context.Context, id, reviewerID, note string) error {
	req, err := s.approvals.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if req.Status.IsTerminal() {
		return domain.ErrApprovalAlreadyProcessed
	}
	if note == "" {
		return domain.ErrReviewNoteRequired
	}
	return s.approvals.UpdateStatus(ctx, id, domain.ApprovalRejected, reviewerID, note)
}

func (s *service) GetApprovalRequests(ctx context.Context, status domain.ApprovalStatus) ([]domain.ApprovalRequest, error) {
	if status != "" {
		return s.approvals.GetByStatus(ctx, status)
	}
	return s.approvals.GetAll(ctx)
}

// ─── COA: Versioning ───────────────────────────────────────────────

func (s *service) CreateAccountVersion(ctx context.Context, reason string) (*domain.AccountVersion, error) {
	accounts, err := s.accounts.GetAll(ctx, false)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(accounts)
	if err != nil {
		return nil, fmt.Errorf("snapshot marshal: %w", err)
	}
	lastVer, err := s.versions.GetLatest(ctx)
	nextVer := "v1"
	if err == nil && lastVer != nil {
		parts := strings.Split(lastVer.VersionNumber, ".")
		if len(parts) >= 1 {
			n, _ := strconv.Atoi(strings.TrimPrefix(parts[0], "v"))
			nextVer = fmt.Sprintf("v%d", n+1)
		}
	}
	now := s.now()
	ver := &domain.AccountVersion{
		ID:            fmt.Sprintf("VER-%d", now.UnixNano()),
		VersionNumber: nextVer,
		Snapshot:      string(data),
		ChangeReason:  reason,
		CreatedAt:     now,
	}
	if err := s.versions.Create(ctx, ver); err != nil {
		return nil, err
	}
	return ver, nil
}

func (s *service) GetVersion(ctx context.Context, versionNumber string) (*domain.AccountVersion, error) {
	return s.versions.GetByVersionNumber(ctx, versionNumber)
}

func (s *service) ListVersions(ctx context.Context) ([]domain.AccountVersion, error) {
	return s.versions.GetAll(ctx)
}

func (s *service) CompareVersions(ctx context.Context, v1, v2 string) (*domain.VersionDiff, error) {
	ver1, err := s.versions.GetByVersionNumber(ctx, v1)
	if err != nil {
		return nil, err
	}
	ver2, err := s.versions.GetByVersionNumber(ctx, v2)
	if err != nil {
		return nil, err
	}
	var accs1, accs2 []domain.Account
	if err := json.Unmarshal([]byte(ver1.Snapshot), &accs1); err != nil {
		return nil, fmt.Errorf("unmarshal version %s: %w", v1, err)
	}
	if err := json.Unmarshal([]byte(ver2.Snapshot), &accs2); err != nil {
		return nil, fmt.Errorf("unmarshal version %s: %w", v2, err)
	}
	diff := &domain.VersionDiff{VersionFrom: v1, VersionTo: v2}
	map1 := make(map[string]domain.Account)
	for _, a := range accs1 {
		map1[a.Code] = a
	}
	map2 := make(map[string]domain.Account)
	for _, a := range accs2 {
		map2[a.Code] = a
	}
	for code, a := range map2 {
		if _, ok := map1[code]; !ok {
			diff.Added = append(diff.Added, a)
		}
	}
	for code, a := range map1 {
		if _, ok := map2[code]; !ok {
			diff.Removed = append(diff.Removed, a)
		}
	}
	for code, a1 := range map1 {
		if a2, ok := map2[code]; ok {
			var ad domain.AccountDiff
			changes := make(map[string]domain.Change)
			if a1.Name != a2.Name {
				changes["name"] = domain.Change{OldValue: a1.Name, NewValue: a2.Name}
			}
			if a1.Type != a2.Type {
				changes["type"] = domain.Change{OldValue: a1.Type, NewValue: a2.Type}
			}
			if a1.ParentCode != a2.ParentCode {
				changes["parent_code"] = domain.Change{OldValue: a1.ParentCode, NewValue: a2.ParentCode}
			}
			if a1.Status != a2.Status {
				changes["status"] = domain.Change{OldValue: a1.Status, NewValue: a2.Status}
			}
			if len(changes) > 0 {
				ad.Code = code
				ad.Old = a1
				ad.New = a2
				ad.Changes = changes
				diff.Modified = append(diff.Modified, ad)
			}
		}
	}
	return diff, nil
}

// ─── COA: Analysis ─────────────────────────────────────────────────

func (s *service) CreateAccountAnalysis(ctx context.Context, analysis *domain.AccountAnalysis) error {
	if _, err := s.accounts.GetByCode(ctx, analysis.AccountCode); err != nil {
		return err
	}
	if err := analysis.Validate(); err != nil {
		return err
	}
	if existing, err := s.analysis.GetByAccount(ctx, analysis.AccountCode); err == nil && existing != nil {
		return s.analysis.Update(ctx, analysis)
	}
	return s.analysis.Create(ctx, analysis)
}

func (s *service) GetAccountAnalysis(ctx context.Context, accountCode string) (*domain.AccountAnalysis, error) {
	return s.analysis.GetByAccount(ctx, accountCode)
}

func (s *service) UpdateAccountAnalysis(ctx context.Context, analysis *domain.AccountAnalysis) error {
	if _, err := s.accounts.GetByCode(ctx, analysis.AccountCode); err != nil {
		return err
	}
	return s.analysis.Update(ctx, analysis)
}

// ─── COA: Mappings ─────────────────────────────────────────────────

func (s *service) CreateAccountMapping(ctx context.Context, mapping *domain.AccountMapping) error {
	if err := mapping.Validate(); err != nil {
		return err
	}
	existing, err := s.mappings.GetByOldCode(ctx, mapping.SourceRegime, mapping.OldCode)
	if err == nil && existing != nil {
		return domain.ErrMappingExists
	}
	return s.mappings.Create(ctx, mapping)
}

func (s *service) GetMappingByOldCode(ctx context.Context, sourceRegime, oldCode string) (*domain.AccountMapping, error) {
	return s.mappings.GetByOldCode(ctx, sourceRegime, oldCode)
}

func (s *service) ListMappings(ctx context.Context, sourceRegime, targetRegime string) ([]domain.AccountMapping, error) {
	if sourceRegime != "" && targetRegime != "" {
		return s.mappings.GetByRegime(ctx, sourceRegime, targetRegime)
	}
	return s.mappings.GetAll(ctx)
}

// ─── COA: IFRS ─────────────────────────────────────────────────────

func (s *service) CreateIFRSMapping(ctx context.Context, mapping *domain.IFRSMapping) error {
	if err := mapping.Validate(); err != nil {
		return err
	}
	existing, err := s.ifrs.GetByVASCode(ctx, mapping.VASCode)
	if err == nil && existing != nil {
		return domain.ErrIFRSMappingExists
	}
	return s.ifrs.Create(ctx, mapping)
}

func (s *service) GetIFRSMapping(ctx context.Context, vasCode string) (*domain.IFRSMapping, error) {
	return s.ifrs.GetByVASCode(ctx, vasCode)
}

func (s *service) ListIFRSMappings(ctx context.Context) ([]domain.IFRSMapping, error) {
	return s.ifrs.GetAll(ctx)
}

// ─── Cash Module ────────────────────────────────────────────────────

// ─── Cash Receipts ──────────────────────────────────────────────────

func (s *service) CreateCashReceipt(ctx context.Context, r *domain.CashReceipt) error {
	if err := r.Validate(); err != nil {
		return err
	}
	if r.Currency == "VND" {
		r.ExchangeRate = 1
		r.AmountVND = r.Amount
	} else {
		r.AmountVND = r.Amount * r.ExchangeRate
	}
	now := s.now()
	nowStr := now.Format("2006-01-02 15:04:05")
	r.Status = domain.CashDraft
	r.CreatedAt = nowStr
	r.UpdatedAt = nowStr

	year := now.Format("2006")
	if len(r.VoucherDate) >= 4 {
		year = r.VoucherDate[:4]
	}
	lastNo, err := s.cash.LastReceiptNo(ctx, r.CompanyID, year)
	var seq int
	if err == nil && lastNo != "" {
		parts := strings.Split(lastNo, "-")
		n := parts[len(parts)-1]
		seq, _ = strconv.Atoi(n)
	}
	seq++
	r.VoucherNo = fmt.Sprintf("R-%s-%04d", year, seq)

	return s.cash.CreateReceipt(ctx, r)
}

func (s *service) GetCashReceipt(ctx context.Context, id string) (*domain.CashReceipt, error) {
	return s.cash.GetReceipt(ctx, id)
}

func (s *service) ListCashReceipts(ctx context.Context, filter domain.CashReceiptFilter) ([]domain.CashReceipt, int, error) {
	return s.cash.ListReceipts(ctx, filter)
}

func (s *service) UpdateCashReceipt(ctx context.Context, r *domain.CashReceipt) error {
	existing, err := s.cash.GetReceipt(ctx, r.ID)
	if err != nil {
		return err
	}
	if existing.Status != domain.CashDraft {
		return domain.ErrCashInvalidStatus
	}
	r.UpdatedAt = s.now().Format("2006-01-02 15:04:05")
	return s.cash.UpdateReceipt(ctx, r)
}

func (s *service) DeleteCashReceipt(ctx context.Context, id string) error {
	existing, err := s.cash.GetReceipt(ctx, id)
	if err != nil {
		return err
	}
	if existing.Status != domain.CashDraft {
		return domain.ErrCashInvalidStatus
	}
	return s.cash.DeleteReceipt(ctx, id)
}

func (s *service) SubmitCashReceipt(ctx context.Context, id, userID string) error {
	r, err := s.cash.GetReceipt(ctx, id)
	if err != nil {
		return err
	}
	if r.Status != domain.CashDraft && r.Status != domain.CashRejected {
		return domain.ErrCashInvalidStatus
	}
	r.Status = domain.CashSubmitted
	r.UpdatedAt = s.now().Format("2006-01-02 15:04:05")
	return s.cash.UpdateReceipt(ctx, r)
}

func (s *service) ApproveCashReceipt(ctx context.Context, id, approverID string) error {
	r, err := s.cash.GetReceipt(ctx, id)
	if err != nil {
		return err
	}
	if !r.Status.ValidTransition(domain.CashApproved) {
		return domain.ErrCashInvalidStatus
	}
	if r.CreatedBy == approverID {
		return domain.ErrSelfCashApproval
	}
	nowStr := s.now().Format("2006-01-02 15:04:05")
	r.Status = domain.CashApproved
	r.ApprovedBy = approverID
	r.ApprovedAt = nowStr
	r.UpdatedAt = nowStr
	return s.cash.UpdateReceipt(ctx, r)
}

func (s *service) RejectCashReceipt(ctx context.Context, id, reviewerID string) error {
	r, err := s.cash.GetReceipt(ctx, id)
	if err != nil {
		return err
	}
	if !r.Status.ValidTransition(domain.CashRejected) {
		return domain.ErrCashInvalidStatus
	}
	nowStr := s.now().Format("2006-01-02 15:04:05")
	r.Status = domain.CashRejected
	r.UpdatedAt = nowStr
	return s.cash.UpdateReceipt(ctx, r)
}

func (s *service) PostCashReceipt(ctx context.Context, id, userID string) error {
	r, err := s.cash.GetReceipt(ctx, id)
	if err != nil {
		return err
	}
	if !r.Status.ValidTransition(domain.CashPosted) {
		return domain.ErrCashInvalidStatus
	}

	entryDate, err := time.Parse("2006-01-02", r.VoucherDate)
	if err != nil {
		return fmt.Errorf("invalid voucher date: %w", err)
	}

	rate := r.ExchangeRate
	if r.Currency == "VND" {
		rate = 1
	}
	entry := &domain.JournalEntry{
		CompanyID:     r.CompanyID,
		CreatedBy:     userID,
		Status:        domain.JournalEntryDraft,
		EntryDate:     entryDate,
		AccountingDate: entryDate,
		Description:   r.Reason,
		CurrencyCode:  r.Currency,
		ExchangeRate:  rate,
		VoucherType:   domain.VoucherTypeReceipt,
		Lines: []domain.JournalLine{
			{AccountCode: r.CashAccountID, DebitAmount: r.AmountVND},
			{AccountCode: r.CreditAccountID, CreditAmount: r.AmountVND},
		},
	}
	s.attachPeriodAndNumber(ctx, entry)
	if err := entry.Validate(); err != nil {
		return err
	}
	if err := s.journals.Create(ctx, entry); err != nil {
		return err
	}

	nowStr := s.now().Format("2006-01-02 15:04:05")
	r.Status = domain.CashPosted
	r.PostedBy = userID
	r.PostedAt = nowStr
	r.GLJournalID = entry.ID
	r.UpdatedAt = nowStr
	return s.cash.UpdateReceipt(ctx, r)
}

// attachPeriodAndNumber attaches the entry to the open period for its entry
// date (if none set) and assigns a sequential entry number within that period,
// same scheme as CreateEntry. Shared by posting flows (cash receipts/payments).
func (s *service) attachPeriodAndNumber(ctx context.Context, entry *domain.JournalEntry) {
	if entry.PeriodID == "" {
		if p, err := s.periods.GetByYearMonth(ctx, entry.EntryDate.Year(), int(entry.EntryDate.Month())); err == nil && p != nil {
			entry.PeriodID = p.ID
		}
	}
	if entry.PeriodID != "" {
		if p, err := s.periods.GetByID(ctx, entry.PeriodID); err == nil && p != nil {
			if es, err := s.journals.GetByPeriod(ctx, p.ID); err == nil {
				entry.EntryNumber = fmt.Sprintf("%d%02d-%03d", p.Year, p.Month, len(es)+1)
			}
		}
	}
}

// ─── Cash Payments ──────────────────────────────────────────────────

func (s *service) CreateCashPayment(ctx context.Context, p *domain.CashPayment) error {
	if err := p.Validate(); err != nil {
		return err
	}
	if p.Currency == "VND" {
		p.ExchangeRate = 1
		p.AmountVND = p.Amount
	} else {
		p.AmountVND = p.Amount * p.ExchangeRate
	}
	now := s.now()
	nowStr := now.Format("2006-01-02 15:04:05")
	p.Status = domain.CashDraft
	p.CreatedAt = nowStr
	p.UpdatedAt = nowStr

	year := now.Format("2006")
	if len(p.VoucherDate) >= 4 {
		year = p.VoucherDate[:4]
	}
	lastNo, err := s.cash.LastPaymentNo(ctx, p.CompanyID, year)
	var seq int
	if err == nil && lastNo != "" {
		parts := strings.Split(lastNo, "-")
		n := parts[len(parts)-1]
		seq, _ = strconv.Atoi(n)
	}
	seq++
	p.VoucherNo = fmt.Sprintf("P-%s-%04d", year, seq)

	return s.cash.CreatePayment(ctx, p)
}

func (s *service) GetCashPayment(ctx context.Context, id string) (*domain.CashPayment, error) {
	return s.cash.GetPayment(ctx, id)
}

func (s *service) ListCashPayments(ctx context.Context, filter domain.CashPaymentFilter) ([]domain.CashPayment, int, error) {
	return s.cash.ListPayments(ctx, filter)
}

func (s *service) UpdateCashPayment(ctx context.Context, p *domain.CashPayment) error {
	existing, err := s.cash.GetPayment(ctx, p.ID)
	if err != nil {
		return err
	}
	if existing.Status != domain.CashDraft {
		return domain.ErrCashInvalidStatus
	}
	p.UpdatedAt = s.now().Format("2006-01-02 15:04:05")
	return s.cash.UpdatePayment(ctx, p)
}

func (s *service) DeleteCashPayment(ctx context.Context, id string) error {
	existing, err := s.cash.GetPayment(ctx, id)
	if err != nil {
		return err
	}
	if existing.Status != domain.CashDraft {
		return domain.ErrCashInvalidStatus
	}
	return s.cash.DeletePayment(ctx, id)
}

func (s *service) SubmitCashPayment(ctx context.Context, id, userID string) error {
	p, err := s.cash.GetPayment(ctx, id)
	if err != nil {
		return err
	}
	if p.Status != domain.CashDraft && p.Status != domain.CashRejected {
		return domain.ErrCashInvalidStatus
	}
	p.Status = domain.CashSubmitted
	p.UpdatedAt = s.now().Format("2006-01-02 15:04:05")
	return s.cash.UpdatePayment(ctx, p)
}

func (s *service) ApproveCashPayment(ctx context.Context, id, approverID string) error {
	p, err := s.cash.GetPayment(ctx, id)
	if err != nil {
		return err
	}
	if !p.Status.ValidTransition(domain.CashApproved) {
		return domain.ErrCashInvalidStatus
	}
	if p.CreatedBy == approverID {
		return domain.ErrSelfCashApproval
	}
	nowStr := s.now().Format("2006-01-02 15:04:05")
	p.Status = domain.CashApproved
	p.ApprovedBy = approverID
	p.ApprovedAt = nowStr
	p.UpdatedAt = nowStr
	return s.cash.UpdatePayment(ctx, p)
}

func (s *service) RejectCashPayment(ctx context.Context, id, reviewerID string) error {
	p, err := s.cash.GetPayment(ctx, id)
	if err != nil {
		return err
	}
	if !p.Status.ValidTransition(domain.CashRejected) {
		return domain.ErrCashInvalidStatus
	}
	nowStr := s.now().Format("2006-01-02 15:04:05")
	p.Status = domain.CashRejected
	p.UpdatedAt = nowStr
	return s.cash.UpdatePayment(ctx, p)
}

func (s *service) PostCashPayment(ctx context.Context, id, userID string) error {
	p, err := s.cash.GetPayment(ctx, id)
	if err != nil {
		return err
	}
	if !p.Status.ValidTransition(domain.CashPosted) {
		return domain.ErrCashInvalidStatus
	}

	bal, _ := s.cash.GetBalance(ctx, p.CompanyID, p.CashAccountID)
	if bal < p.AmountVND {
		return domain.ErrCashInsufficientBalance
	}

	entryDate, err := time.Parse("2006-01-02", p.VoucherDate)
	if err != nil {
		return fmt.Errorf("invalid voucher date: %w", err)
	}

	rate := p.ExchangeRate
	if p.Currency == "VND" {
		rate = 1
	}
	entry := &domain.JournalEntry{
		CompanyID:     p.CompanyID,
		CreatedBy:     userID,
		Status:        domain.JournalEntryDraft,
		EntryDate:     entryDate,
		AccountingDate: entryDate,
		Description:   p.Reason,
		CurrencyCode:  p.Currency,
		ExchangeRate:  rate,
		VoucherType:   domain.VoucherTypePayment,
		Lines: []domain.JournalLine{
			{AccountCode: p.DebitAccountID, DebitAmount: p.AmountVND},
			{AccountCode: p.CashAccountID, CreditAmount: p.AmountVND},
		},
	}
	s.attachPeriodAndNumber(ctx, entry)
	if err := entry.Validate(); err != nil {
		return err
	}
	if err := s.journals.Create(ctx, entry); err != nil {
		return err
	}

	nowStr := s.now().Format("2006-01-02 15:04:05")
	p.Status = domain.CashPosted
	p.PostedBy = userID
	p.PostedAt = nowStr
	p.GLJournalID = entry.ID
	p.UpdatedAt = nowStr
	return s.cash.UpdatePayment(ctx, p)
}

// ─── Cash Transfers ─────────────────────────────────────────────────

func (s *service) CreateCashTransfer(ctx context.Context, t *domain.CashTransfer, userID string) error {
	if err := t.Validate(); err != nil {
		return err
	}
	rate := t.ExchangeRate
	if t.Currency == "VND" {
		rate = 1
	}

	now := s.now()
	nowStr := now.Format("2006-01-02 15:04:05")
	amountVND := t.Amount * rate

	year := now.Format("2006")
	if len(t.TransferDate) >= 4 {
		year = t.TransferDate[:4]
	}

	// create receipt voucher (sequential per year, same scheme as CreateCashReceipt)
	lastR, err := s.cash.LastReceiptNo(ctx, t.CompanyID, year)
	seqR := 1
	if err == nil && lastR != "" {
		if parts := strings.Split(lastR, "-"); len(parts) > 0 {
			seqR, _ = strconv.Atoi(parts[len(parts)-1])
			seqR++
		}
	}
	receipt := &domain.CashReceipt{
		CompanyID:       t.CompanyID,
		VoucherNo:       fmt.Sprintf("R-%s-%04d", year, seqR),
		VoucherDate:     t.TransferDate,
		CashAccountID:   t.ToAccountID,
		Currency:        t.Currency,
		ExchangeRate:    rate,
		Amount:          t.Amount,
		AmountVND:       amountVND,
		DebitAccountID:  t.ToAccountID,
		CreditAccountID: t.FromAccountID,
		Reason:          t.Reason,
		Status:          domain.CashDraft,
		CreatedAt:       nowStr,
		UpdatedAt:       nowStr,
	}
	if err := s.cash.CreateReceipt(ctx, receipt); err != nil {
		return err
	}

	// create payment voucher (sequential per year, same scheme as CreateCashPayment)
	lastP, err := s.cash.LastPaymentNo(ctx, t.CompanyID, year)
	seqP := 1
	if err == nil && lastP != "" {
		if parts := strings.Split(lastP, "-"); len(parts) > 0 {
			seqP, _ = strconv.Atoi(parts[len(parts)-1])
			seqP++
		}
	}
	payment := &domain.CashPayment{
		CompanyID:       t.CompanyID,
		VoucherNo:       fmt.Sprintf("P-%s-%04d", year, seqP),
		VoucherDate:     t.TransferDate,
		CashAccountID:   t.FromAccountID,
		Currency:        t.Currency,
		ExchangeRate:    rate,
		Amount:          t.Amount,
		AmountVND:       amountVND,
		DebitAccountID:  t.ToAccountID,
		CreditAccountID: t.FromAccountID,
		Reason:          t.Reason,
		Status:          domain.CashDraft,
		CreatedAt:       nowStr,
		UpdatedAt:       nowStr,
	}
	if err := s.cash.CreatePayment(ctx, payment); err != nil {
		return err
	}

	// create single journal entry for the transfer
	entryDate, _ := time.Parse("2006-01-02", t.TransferDate)
	entry := &domain.JournalEntry{
		CompanyID:     t.CompanyID,
		CreatedBy:     userID,
		Status:        domain.JournalEntryDraft,
		EntryDate:     entryDate,
		AccountingDate: entryDate,
		Description:   t.Reason,
		CurrencyCode:  t.Currency,
		ExchangeRate:  rate,
		VoucherType:   domain.VoucherTypeOther,
		Lines: []domain.JournalLine{
			{AccountCode: t.ToAccountID, DebitAmount: amountVND},
			{AccountCode: t.FromAccountID, CreditAmount: amountVND},
		},
	}
	s.attachPeriodAndNumber(ctx, entry)
	if err := entry.Validate(); err != nil {
		return err
	}
	if err := s.journals.Create(ctx, entry); err != nil {
		return err
	}

	// post both vouchers with same journal entry
	receipt.Status = domain.CashPosted
	receipt.PostedAt = nowStr
	receipt.GLJournalID = entry.ID
	if err := s.cash.UpdateReceipt(ctx, receipt); err != nil {
		return err
	}

	payment.Status = domain.CashPosted
	payment.PostedAt = nowStr
	payment.GLJournalID = entry.ID
	if err := s.cash.UpdatePayment(ctx, payment); err != nil {
		return err
	}

	// create transfer record
	t.ID = fmt.Sprintf("TRF-%d", now.UnixNano())
	t.Status = domain.CashPosted
	t.SourceVoucherID = receipt.ID
	t.DestVoucherID = payment.ID
	t.CreatedAt = nowStr
	t.PostedAt = nowStr

	return s.cash.CreateTransfer(ctx, t)
}

func (s *service) GetCashTransfers(ctx context.Context, companyID string) ([]domain.CashTransfer, error) {
	return s.cash.ListTransfers(ctx, companyID)
}

// ─── Cash Book & Balance ────────────────────────────────────────────

func (s *service) GetCashBook(ctx context.Context, companyID, currency, accountID, fromDate, toDate string) (*domain.CashBook, error) {
	return s.cash.GetCashBook(ctx, companyID, currency, accountID, fromDate, toDate)
}

func (s *service) GetCashBalance(ctx context.Context, companyID, accountID string) (float64, error) {
	return s.cash.GetBalance(ctx, companyID, accountID)
}

type CashFlowStatement struct {
	CompanyID        string  `json:"company_id"`
	FromDate         string  `json:"from_date"`
	ToDate           string  `json:"to_date"`
	OpeningBalance   float64 `json:"opening_balance"`
	OperatingInflow  float64 `json:"operating_inflow"`
	OperatingOutflow float64 `json:"operating_outflow"`
	OperatingNet     float64 `json:"operating_net"`
	InvestingInflow  float64 `json:"investing_inflow"`
	InvestingOutflow float64 `json:"investing_outflow"`
	InvestingNet     float64 `json:"investing_net"`
	FinancingInflow  float64 `json:"financing_inflow"`
	FinancingOutflow float64 `json:"financing_outflow"`
	FinancingNet     float64 `json:"financing_net"`
	NetCashFlow      float64 `json:"net_cash_flow"`
	ClosingBalance   float64 `json:"closing_balance"`
}

func (s *service) GetCashFlowStatement(ctx context.Context, companyID, currency, accountID, fromDate, toDate string) (*CashFlowStatement, error) {
	cb, err := s.cash.GetCashBook(ctx, companyID, currency, accountID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	stmt := &CashFlowStatement{
		CompanyID:        companyID,
		FromDate:         fromDate,
		ToDate:           toDate,
		OpeningBalance:   cb.OpeningBalance,
		OperatingInflow:  cb.TotalReceipts,
		OperatingOutflow: cb.TotalPayments,
	}
	stmt.OperatingNet = stmt.OperatingInflow - stmt.OperatingOutflow
	stmt.NetCashFlow = stmt.OperatingNet
	stmt.ClosingBalance = cb.ClosingBalance
	return stmt, nil
}

// ─── Petty Cash ─────────────────────────────────────────────────────

func (s *service) CreatePettyCashFund(ctx context.Context, f *domain.PettyCashFund) error {
	if f.FundCode == "" {
		return fmt.Errorf("fund code required")
	}
	if f.FundName == "" {
		return fmt.Errorf("fund name required")
	}
	if f.CustodianID == "" {
		return fmt.Errorf("custodian required")
	}
	if f.InitialAmount < 0 {
		return fmt.Errorf("initial amount must be non-negative")
	}
	if f.Currency == "" {
		f.Currency = "VND"
	}
	f.CurrentBalance = f.InitialAmount
	f.Status = domain.PettyCashActive
	f.CreatedAt = s.now().Format("2006-01-02 15:04:05")
	return s.cash.CreatePettyCashFund(ctx, f)
}

func (s *service) GetPettyCashFund(ctx context.Context, id string) (*domain.PettyCashFund, error) {
	return s.cash.GetPettyCashFund(ctx, id)
}

func (s *service) ListPettyCashFunds(ctx context.Context, companyID string) ([]domain.PettyCashFund, error) {
	return s.cash.ListPettyCashFunds(ctx, companyID)
}

// ─── Cash Inventory ─────────────────────────────────────────────────

func (s *service) CreateCashInventory(ctx context.Context, inv *domain.CashInventory) error {
	if inv.InventoryDate == "" {
		return fmt.Errorf("inventory date required")
	}
	if inv.CashAccountID == "" {
		return fmt.Errorf("cash account required")
	}
	if inv.Currency == "" {
		inv.Currency = "VND"
	}
	inv.Difference = inv.ActualBalance - inv.BookBalance
	inv.Status = domain.CashInventoryDraft
	inv.CreatedAt = s.now().Format("2006-01-02 15:04:05")
	return s.cash.CreateInventory(ctx, inv)
}

func (s *service) GetCashInventory(ctx context.Context, id string) (*domain.CashInventory, error) {
	return s.cash.GetInventory(ctx, id)
}

func (s *service) ListCashInventories(ctx context.Context, companyID string) ([]domain.CashInventory, error) {
	return s.cash.ListInventories(ctx, companyID)
}

// ─── Advance Request / Settlement ─────────────────────────────────────

func (s *service) CreateAdvance(ctx context.Context, a *domain.AdvanceRequest) error {
	if err := a.Validate(); err != nil {
		return err
	}
	now := s.now().Format("2006-01-02 15:04:05")
	a.ID = "" // let repo generate
	a.Status = domain.AdvanceDraft
	a.CreatedAt = now
	a.UpdatedAt = now
	return s.cash.CreateAdvance(ctx, a)
}

func (s *service) GetAdvance(ctx context.Context, id string) (*domain.AdvanceRequest, error) {
	return s.cash.GetAdvance(ctx, id)
}

func (s *service) ListAdvances(ctx context.Context, companyID string) ([]domain.AdvanceRequest, error) {
	return s.cash.ListAdvances(ctx, companyID)
}

func (s *service) UpdateAdvance(ctx context.Context, a *domain.AdvanceRequest) error {
	existing, err := s.cash.GetAdvance(ctx, a.ID)
	if err != nil {
		return err
	}
	if existing.Status != domain.AdvanceDraft && existing.Status != domain.AdvanceRejected {
		return fmt.Errorf("can only update draft or rejected advances")
	}
	if err := a.Validate(); err != nil {
		return err
	}
	a.Status = existing.Status
	a.CreatedAt = existing.CreatedAt
	a.UpdatedAt = s.now().Format("2006-01-02 15:04:05")
	return s.cash.UpdateAdvance(ctx, a)
}

func (s *service) ApproveAdvance(ctx context.Context, id, approverID string) error {
	a, err := s.cash.GetAdvance(ctx, id)
	if err != nil {
		return err
	}
	if !a.Status.ValidTransition(domain.AdvanceApproved) {
		return fmt.Errorf("cannot approve from status %s", a.Status)
	}
	now := s.now().Format("2006-01-02 15:04:05")
	a.Status = domain.AdvanceApproved
	a.ApprovedBy = approverID
	a.ApprovedAt = now
	a.UpdatedAt = now
	return s.cash.UpdateAdvance(ctx, a)
}

func (s *service) RejectAdvance(ctx context.Context, id, reviewerID string) error {
	a, err := s.cash.GetAdvance(ctx, id)
	if err != nil {
		return err
	}
	if !a.Status.ValidTransition(domain.AdvanceRejected) {
		return fmt.Errorf("cannot reject from status %s", a.Status)
	}
	now := s.now().Format("2006-01-02 15:04:05")
	a.Status = domain.AdvanceRejected
	a.UpdatedAt = now
	return s.cash.UpdateAdvance(ctx, a)
}

func (s *service) PayAdvance(ctx context.Context, id, paidBy string) error {
	a, err := s.cash.GetAdvance(ctx, id)
	if err != nil {
		return err
	}
	if !a.Status.ValidTransition(domain.AdvancePaid) {
		return fmt.Errorf("cannot pay from status %s", a.Status)
	}
	if a.Amount <= 0 {
		return domain.ErrCashAmountRequired
	}
	now := s.now().Format("2006-01-02 15:04:05")
	a.Status = domain.AdvancePaid
	a.PaidBy = paidBy
	a.PaidAt = now
	a.UpdatedAt = now
	return s.cash.UpdateAdvance(ctx, a)
}

func (s *service) SettleAdvance(ctx context.Context, id, settlementID string) error {
	a, err := s.cash.GetAdvance(ctx, id)
	if err != nil {
		return err
	}
	if !a.Status.ValidTransition(domain.AdvanceSettled) {
		return fmt.Errorf("cannot settle from status %s", a.Status)
	}
	now := s.now().Format("2006-01-02 15:04:05")
	a.Status = domain.AdvanceSettled
	a.UpdatedAt = now
	return s.cash.UpdateAdvance(ctx, a)
}

func (s *service) CreateAdvanceSettlement(ctx context.Context, set *domain.AdvanceSettlement) error {
	if set.AdvanceID == "" {
		return fmt.Errorf("advance_id required")
	}
	if set.TotalSpent < 0 {
		return fmt.Errorf("total_spent cannot be negative")
	}
	now := s.now().Format("2006-01-02 15:04:05")
	set.ID = ""
	set.Status = "DRAFT"
	set.CreatedAt = now
	return s.cash.CreateAdvanceSettlement(ctx, set)
}

func (s *service) ListAdvanceByStatus(ctx context.Context, companyID string, status domain.AdvanceStatus) ([]domain.AdvanceRequest, error) {
	return s.cash.ListAdvancesByStatus(ctx, companyID, status)
}

func (s *service) AdminResetPassword(ctx context.Context, userID, newPassword string) error {
	if err := domain.ValidatePassword(newPassword); err != nil {
		return err
	}
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	newHash, err := auth.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("hashing failed: %w", err)
	}
	user.PasswordHistory = append(user.PasswordHistory, user.PasswordHash)
	if len(user.PasswordHistory) > domain.PasswordHistorySize {
		user.PasswordHistory = user.PasswordHistory[len(user.PasswordHistory)-domain.PasswordHistorySize:]
	}
	user.PasswordHash = newHash
	now := s.now()
	user.PasswordChangedAt = &now
	user.FailedAttempts = 0
	user.LockedUntil = nil
	return s.users.Update(ctx, user)
}

func (s *service) UpdateUserRole(ctx context.Context, userID string, role domain.UserRole) error {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	user.Role = role
	return s.users.Update(ctx, user)
}
