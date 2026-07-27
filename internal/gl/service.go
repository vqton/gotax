package gl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Service interface {
	CreateAccount(ctx context.Context, account *Account) error
	GetAccount(ctx context.Context, code string) (*Account, error)
	GetAllAccounts(ctx context.Context, activeOnly bool) ([]Account, error)
	UpdateAccount(ctx context.Context, account *Account) error
	DeleteAccount(ctx context.Context, code string) error
	FreezeAccount(ctx context.Context, code, reason string) error
	UnfreezeAccount(ctx context.Context, code, reason string) error

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

	// COA Extension Methods
	GetAccountBalance(ctx context.Context, code string, periodID string) (*AccountBalance, error)
	GetAccountBalanceDrillDown(ctx context.Context, code string, periodID string) ([]JournalEntry, error)
	GetAccountUsage(ctx context.Context, code string) (*AccountUsage, error)

	CreateApprovalRequest(ctx context.Context, req *ApprovalRequest) error
	ApproveRequest(ctx context.Context, id, reviewerID, note string) error
	RejectRequest(ctx context.Context, id, reviewerID, note string) error
	GetApprovalRequests(ctx context.Context, status ApprovalStatus) ([]ApprovalRequest, error)

	CreateAccountVersion(ctx context.Context, reason string) (*AccountVersion, error)
	GetVersion(ctx context.Context, versionNumber string) (*AccountVersion, error)
	ListVersions(ctx context.Context) ([]AccountVersion, error)
	CompareVersions(ctx context.Context, v1, v2 string) (*VersionDiff, error)

	CreateAccountAnalysis(ctx context.Context, analysis *AccountAnalysis) error
	GetAccountAnalysis(ctx context.Context, accountCode string) (*AccountAnalysis, error)
	UpdateAccountAnalysis(ctx context.Context, analysis *AccountAnalysis) error

	CreateAccountMapping(ctx context.Context, mapping *AccountMapping) error
	GetMappingByOldCode(ctx context.Context, sourceRegime, oldCode string) (*AccountMapping, error)
	ListMappings(ctx context.Context, sourceRegime, targetRegime string) ([]AccountMapping, error)

	CreateIFRSMapping(ctx context.Context, mapping *IFRSMapping) error
	GetIFRSMapping(ctx context.Context, vasCode string) (*IFRSMapping, error)
	ListIFRSMappings(ctx context.Context) ([]IFRSMapping, error)
}

type service struct {
	accounts   AccountRepository
	journals   JournalRepository
	periods    PeriodRepository
	users      UserRepository
	audit      AuditLogRepository
	rates      ExchangeRateRepository
	templates  ClosingTemplateRepository
	approvals  ApprovalRepository
	versions   AccountVersionRepository
	mappings   AccountMappingRepository
	analysis   AccountAnalysisRepository
	ifrs       IFRSMappingRepository
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
	approvalRepo ApprovalRepository,
	versionRepo AccountVersionRepository,
	mappingRepo AccountMappingRepository,
	analysisRepo AccountAnalysisRepository,
	ifrsRepo IFRSMappingRepository,
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
		now:       time.Now,
	}
}

// -- Accounts --

func (s *service) CreateAccount(ctx context.Context, account *Account) error {
	if err := account.Validate(); err != nil {
		return err
	}
	existing, err := s.accounts.GetByCode(ctx, account.Code)
	if err != nil && !errors.Is(err, ErrAccountNotFound) {
		return err
	}
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
		if err := acc.CanPost(); err != nil {
			return err
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
	period, err := s.periods.GetByID(ctx, entry.PeriodID)
	if err != nil {
		return err
	}
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
	if err := period.Validate(); err != nil {
		return err
	}
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
	entries, err := s.journals.GetByPeriod(ctx, id)
	if err != nil {
		return err
	}
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
	existing, err := s.users.GetByUsername(ctx, user.Username)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return err
	}
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

// ─── COA: Account Freeze/Unfreeze ───────────────────────────────────────────

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

// ─── COA: Account Balance & Drill-down ──────────────────────────────────────

func (s *service) GetAccountBalance(ctx context.Context, code string, periodID string) (*AccountBalance, error) {
	if _, err := s.accounts.GetByCode(ctx, code); err != nil {
		return nil, err
	}
	return s.journals.GetBalance(ctx, code, periodID)
}

func (s *service) GetAccountBalanceDrillDown(ctx context.Context, code string, periodID string) ([]JournalEntry, error) {
	if _, err := s.accounts.GetByCode(ctx, code); err != nil {
		return nil, err
	}
	all, err := s.journals.GetByPeriod(ctx, periodID)
	if err != nil {
		return nil, err
	}
	var result []JournalEntry
	for _, e := range all {
		if e.Status != JournalEntryPosted {
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

func (s *service) GetAccountUsage(ctx context.Context, code string) (*AccountUsage, error) {
	if _, err := s.accounts.GetByCode(ctx, code); err != nil {
		return nil, err
	}
	all, err := s.journals.GetByDateRange(ctx, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), s.now())
	if err != nil {
		return nil, err
	}
	usage := &AccountUsage{AccountCode: code}
	for _, e := range all {
		if e.Status != JournalEntryPosted {
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

// ─── COA: Approval Workflow ─────────────────────────────────────────────────

func (s *service) CreateApprovalRequest(ctx context.Context, req *ApprovalRequest) error {
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
		return ErrApprovalAlreadyProcessed
	}
	if req.RequestedBy == reviewerID {
		return ErrSelfApproval
	}
	if !req.ExpiresAt.IsZero() && s.now().After(req.ExpiresAt) {
		s.approvals.UpdateStatus(ctx, id, ApprovalExpired, "", "")
		return ErrApprovalExpired
	}
	return s.approvals.UpdateStatus(ctx, id, ApprovalApproved, reviewerID, note)
}

func (s *service) RejectRequest(ctx context.Context, id, reviewerID, note string) error {
	req, err := s.approvals.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if req.Status.IsTerminal() {
		return ErrApprovalAlreadyProcessed
	}
	if note == "" {
		return ErrReviewNoteRequired
	}
	return s.approvals.UpdateStatus(ctx, id, ApprovalRejected, reviewerID, note)
}

func (s *service) GetApprovalRequests(ctx context.Context, status ApprovalStatus) ([]ApprovalRequest, error) {
	if status != "" {
		return s.approvals.GetByStatus(ctx, status)
	}
	return s.approvals.GetAll(ctx)
}

// ─── COA: Versioning ────────────────────────────────────────────────────────

func (s *service) CreateAccountVersion(ctx context.Context, reason string) (*AccountVersion, error) {
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
	ver := &AccountVersion{
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

func (s *service) GetVersion(ctx context.Context, versionNumber string) (*AccountVersion, error) {
	return s.versions.GetByVersionNumber(ctx, versionNumber)
}

func (s *service) ListVersions(ctx context.Context) ([]AccountVersion, error) {
	return s.versions.GetAll(ctx)
}

func (s *service) CompareVersions(ctx context.Context, v1, v2 string) (*VersionDiff, error) {
	ver1, err := s.versions.GetByVersionNumber(ctx, v1)
	if err != nil {
		return nil, err
	}
	ver2, err := s.versions.GetByVersionNumber(ctx, v2)
	if err != nil {
		return nil, err
	}
	var accs1, accs2 []Account
	if err := json.Unmarshal([]byte(ver1.Snapshot), &accs1); err != nil {
		return nil, fmt.Errorf("unmarshal version %s: %w", v1, err)
	}
	if err := json.Unmarshal([]byte(ver2.Snapshot), &accs2); err != nil {
		return nil, fmt.Errorf("unmarshal version %s: %w", v2, err)
	}
	diff := &VersionDiff{VersionFrom: v1, VersionTo: v2}
	map1 := make(map[string]Account)
	for _, a := range accs1 {
		map1[a.Code] = a
	}
	map2 := make(map[string]Account)
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
			var ad AccountDiff
			changes := make(map[string]Change)
			if a1.Name != a2.Name {
				changes["name"] = Change{OldValue: a1.Name, NewValue: a2.Name}
			}
			if a1.Type != a2.Type {
				changes["type"] = Change{OldValue: a1.Type, NewValue: a2.Type}
			}
			if a1.ParentCode != a2.ParentCode {
				changes["parent_code"] = Change{OldValue: a1.ParentCode, NewValue: a2.ParentCode}
			}
			if a1.Status != a2.Status {
				changes["status"] = Change{OldValue: a1.Status, NewValue: a2.Status}
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

// ─── COA: Account Analysis ─────────────────────────────────────────────────

func (s *service) CreateAccountAnalysis(ctx context.Context, analysis *AccountAnalysis) error {
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

func (s *service) GetAccountAnalysis(ctx context.Context, accountCode string) (*AccountAnalysis, error) {
	return s.analysis.GetByAccount(ctx, accountCode)
}

func (s *service) UpdateAccountAnalysis(ctx context.Context, analysis *AccountAnalysis) error {
	if _, err := s.accounts.GetByCode(ctx, analysis.AccountCode); err != nil {
		return err
	}
	return s.analysis.Update(ctx, analysis)
}

// ─── COA: Account Mappings ─────────────────────────────────────────────────

func (s *service) CreateAccountMapping(ctx context.Context, mapping *AccountMapping) error {
	if err := mapping.Validate(); err != nil {
		return err
	}
	existing, err := s.mappings.GetByOldCode(ctx, mapping.SourceRegime, mapping.OldCode)
	if err == nil && existing != nil {
		return ErrMappingExists
	}
	return s.mappings.Create(ctx, mapping)
}

func (s *service) GetMappingByOldCode(ctx context.Context, sourceRegime, oldCode string) (*AccountMapping, error) {
	return s.mappings.GetByOldCode(ctx, sourceRegime, oldCode)
}

func (s *service) ListMappings(ctx context.Context, sourceRegime, targetRegime string) ([]AccountMapping, error) {
	if sourceRegime != "" && targetRegime != "" {
		return s.mappings.GetByRegime(ctx, sourceRegime, targetRegime)
	}
	return s.mappings.GetAll(ctx)
}

// ─── COA: IFRS Mapping ────────────────────────────────────────────────────

func (s *service) CreateIFRSMapping(ctx context.Context, mapping *IFRSMapping) error {
	if err := mapping.Validate(); err != nil {
		return err
	}
	existing, err := s.ifrs.GetByVASCode(ctx, mapping.VASCode)
	if err == nil && existing != nil {
		return ErrIFRSMappingExists
	}
	return s.ifrs.Create(ctx, mapping)
}

func (s *service) GetIFRSMapping(ctx context.Context, vasCode string) (*IFRSMapping, error) {
	return s.ifrs.GetByVASCode(ctx, vasCode)
}

func (s *service) ListIFRSMappings(ctx context.Context) ([]IFRSMapping, error) {
	return s.ifrs.GetAll(ctx)
}
