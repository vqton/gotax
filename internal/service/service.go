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
	SubmitForReview(ctx context.Context, id, userID string) error
	ReviewEntry(ctx context.Context, id, reviewerID string) error
	ApproveEntry(ctx context.Context, id, approverID string) error
	PostEntry(ctx context.Context, id string) error
	CancelEntry(ctx context.Context, id string) error
	GetEntryByID(ctx context.Context, id string) (*domain.JournalEntry, error)
	GetEntriesByDateRange(ctx context.Context, from, to time.Time) ([]domain.JournalEntry, error)
	GetEntriesByStatus(ctx context.Context, status domain.JournalEntryStatus) ([]domain.JournalEntry, error)

	TrialBalance(ctx context.Context, year, month int) ([]domain.AccountBalance, error)
	BalanceSheet(ctx context.Context, year, month int) ([]domain.AccountBalance, error)
	IncomeStatement(ctx context.Context, year, month int) ([]domain.AccountBalance, error)

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

	CreateCircular99Mapping(ctx context.Context, m *domain.Circular99Mapping) error
	ListCircular99Mappings(ctx context.Context) ([]domain.Circular99Mapping, error)
	GetCircular99MappingByOldCode(ctx context.Context, oldCode string) (*domain.Circular99Mapping, error)

	CreateBalanceMigration(ctx context.Context, m *domain.BalanceMigration) error
	GetBalanceMigrationByID(ctx context.Context, id string) (*domain.BalanceMigration, error)
	ListBalanceMigrations(ctx context.Context, companyID string) ([]domain.BalanceMigration, error)

	ImportOpeningBalances(ctx context.Context, data []byte, companyID, periodID, createdBy string) (*OBImportResult, error)
	GenerateOpeningBalancePDF(ctx context.Context, companyID, periodID string) ([]byte, error)
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
		return domain.ErrAccountHasChildren
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
	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
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

	user.PasswordHash = newHash
	now := s.now()
	user.PasswordChangedAt = &now
	user.FailedAttempts = 0
	user.LockedUntil = nil
	user.PasswordHistory = nil

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
	all, err := s.journals.GetByPeriod(ctx, periodID)
	if err != nil {
		return nil, err
	}
	var result []domain.JournalEntry
	for _, e := range all {
		if e.Status != domain.JournalEntryPosted {
			continue
		}
		for _, l := range e.Lines {
			if l.AccountCode == code {
				result = append(result, e)
				break
			}
		}
	}
	return result, nil
}

func (s *service) GetAccountUsage(ctx context.Context, code string) (*domain.AccountUsage, error) {
	if _, err := s.accounts.GetByCode(ctx, code); err != nil {
		return nil, err
	}
	all, err := s.journals.GetByDateRange(ctx, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), s.now())
	if err != nil {
		return nil, err
	}
	usage := &domain.AccountUsage{AccountCode: code}
	for _, e := range all {
		if e.Status != domain.JournalEntryPosted {
			continue
		}
		for _, l := range e.Lines {
			if l.AccountCode == code {
				usage.EntryCount++
				usage.TotalDebit += l.DebitAmount
				usage.TotalCredit += l.CreditAmount
				if e.EntryDate.Format("2006-01-02") > usage.LastUsedDate {
					usage.LastUsedDate = e.EntryDate.Format("2006-01-02")
				}
			}
		}
	}
	return usage, nil
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
