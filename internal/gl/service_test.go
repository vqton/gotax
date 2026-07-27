package gl

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAccountRepo struct {
	mock.Mock
}

func (m *mockAccountRepo) Create(ctx context.Context, a *Account) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}
func (m *mockAccountRepo) GetByCode(ctx context.Context, code string) (*Account, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Account), args.Error(1)
}
func (m *mockAccountRepo) GetAll(ctx context.Context, activeOnly bool) ([]Account, error) {
	args := m.Called(ctx, activeOnly)
	return args.Get(0).([]Account), args.Error(1)
}
func (m *mockAccountRepo) Update(ctx context.Context, a *Account) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}
func (m *mockAccountRepo) Delete(ctx context.Context, code string) error {
	args := m.Called(ctx, code)
	return args.Error(0)
}
func (m *mockAccountRepo) GetChildren(ctx context.Context, parentCode string) ([]Account, error) {
	args := m.Called(ctx, parentCode)
	return args.Get(0).([]Account), args.Error(1)
}

type mockJournalRepo struct {
	mock.Mock
}

func (m *mockJournalRepo) Create(ctx context.Context, e *JournalEntry) error {
	args := m.Called(ctx, e)
	return args.Error(0)
}
func (m *mockJournalRepo) GetByID(ctx context.Context, id string) (*JournalEntry, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*JournalEntry), args.Error(1)
}
func (m *mockJournalRepo) GetByPeriod(ctx context.Context, periodID string) ([]JournalEntry, error) {
	args := m.Called(ctx, periodID)
	return args.Get(0).([]JournalEntry), args.Error(1)
}
func (m *mockJournalRepo) GetByDateRange(ctx context.Context, from, to time.Time) ([]JournalEntry, error) {
	args := m.Called(ctx, from, to)
	return args.Get(0).([]JournalEntry), args.Error(1)
}
func (m *mockJournalRepo) GetByStatus(ctx context.Context, status JournalEntryStatus) ([]JournalEntry, error) {
	args := m.Called(ctx, status)
	return args.Get(0).([]JournalEntry), args.Error(1)
}
func (m *mockJournalRepo) GetByVoucherType(ctx context.Context, voucherType VoucherType) ([]JournalEntry, error) {
	args := m.Called(ctx, voucherType)
	return args.Get(0).([]JournalEntry), args.Error(1)
}
func (m *mockJournalRepo) UpdateStatus(ctx context.Context, id string, status JournalEntryStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}
func (m *mockJournalRepo) Update(ctx context.Context, entry *JournalEntry) error {
	args := m.Called(ctx, entry)
	return args.Error(0)
}
func (m *mockJournalRepo) Approve(ctx context.Context, id, approvedBy string) error {
	args := m.Called(ctx, id, approvedBy)
	return args.Error(0)
}
func (m *mockJournalRepo) Review(ctx context.Context, id, reviewedBy string) error {
	args := m.Called(ctx, id, reviewedBy)
	return args.Error(0)
}
func (m *mockJournalRepo) GetLinesByEntryID(ctx context.Context, entryID string) ([]JournalLine, error) {
	args := m.Called(ctx, entryID)
	return args.Get(0).([]JournalLine), args.Error(1)
}
func (m *mockJournalRepo) GetBalance(ctx context.Context, accountCode, periodID string) (*AccountBalance, error) {
	args := m.Called(ctx, accountCode, periodID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AccountBalance), args.Error(1)
}
func (m *mockJournalRepo) GetTrialBalance(ctx context.Context, periodID string) ([]AccountBalance, error) {
	args := m.Called(ctx, periodID)
	return args.Get(0).([]AccountBalance), args.Error(1)
}
func (m *mockJournalRepo) GetFinancialStatement(ctx context.Context, periodID string, accountTypes []AccountType) ([]AccountBalance, error) {
	args := m.Called(ctx, periodID, accountTypes)
	return args.Get(0).([]AccountBalance), args.Error(1)
}

type mockPeriodRepo struct {
	mock.Mock
}

func (m *mockPeriodRepo) Create(ctx context.Context, p *Period) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}
func (m *mockPeriodRepo) GetByID(ctx context.Context, id string) (*Period, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Period), args.Error(1)
}
func (m *mockPeriodRepo) GetByYearMonth(ctx context.Context, year, month int) (*Period, error) {
	args := m.Called(ctx, year, month)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Period), args.Error(1)
}
func (m *mockPeriodRepo) GetAll(ctx context.Context) ([]Period, error) {
	args := m.Called(ctx)
	return args.Get(0).([]Period), args.Error(1)
}
func (m *mockPeriodRepo) UpdateStatus(ctx context.Context, id string, status PeriodStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}
func (m *mockPeriodRepo) GetOpenPeriod(ctx context.Context) (*Period, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Period), args.Error(1)
}

type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) Create(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *mockUserRepo) GetByID(ctx context.Context, id string) (*User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}
func (m *mockUserRepo) GetByUsername(ctx context.Context, username string) (*User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}
func (m *mockUserRepo) GetAll(ctx context.Context) ([]User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]User), args.Error(1)
}
func (m *mockUserRepo) Update(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *mockUserRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type mockAuditRepo struct {
	mock.Mock
}

func (m *mockAuditRepo) Create(ctx context.Context, entry *AuditEntry) error {
	args := m.Called(ctx, entry)
	return args.Error(0)
}
func (m *mockAuditRepo) GetByEntity(ctx context.Context, entityType, entityID string) ([]AuditEntry, error) {
	args := m.Called(ctx, entityType, entityID)
	return args.Get(0).([]AuditEntry), args.Error(1)
}
func (m *mockAuditRepo) GetByUser(ctx context.Context, userID string) ([]AuditEntry, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]AuditEntry), args.Error(1)
}
func (m *mockAuditRepo) GetByDateRange(ctx context.Context, from, to time.Time) ([]AuditEntry, error) {
	args := m.Called(ctx, from, to)
	return args.Get(0).([]AuditEntry), args.Error(1)
}
func (m *mockAuditRepo) GetAll(ctx context.Context, limit int) ([]AuditEntry, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]AuditEntry), args.Error(1)
}

type mockRateRepo struct {
	mock.Mock
}

func (m *mockRateRepo) Create(ctx context.Context, rate *ExchangeRate) error {
	args := m.Called(ctx, rate)
	return args.Error(0)
}
func (m *mockRateRepo) GetByCurrencyAndDate(ctx context.Context, currencyCode string, rateDate time.Time) (*ExchangeRate, error) {
	args := m.Called(ctx, currencyCode, rateDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ExchangeRate), args.Error(1)
}
func (m *mockRateRepo) GetByDateRange(ctx context.Context, from, to time.Time) ([]ExchangeRate, error) {
	args := m.Called(ctx, from, to)
	return args.Get(0).([]ExchangeRate), args.Error(1)
}
func (m *mockRateRepo) GetAll(ctx context.Context) ([]ExchangeRate, error) {
	args := m.Called(ctx)
	return args.Get(0).([]ExchangeRate), args.Error(1)
}
func (m *mockRateRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type mockTemplateRepo struct {
	mock.Mock
}

func (m *mockTemplateRepo) Create(ctx context.Context, template *ClosingTemplate) error {
	args := m.Called(ctx, template)
	return args.Error(0)
}
func (m *mockTemplateRepo) GetByID(ctx context.Context, id string) (*ClosingTemplate, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ClosingTemplate), args.Error(1)
}
func (m *mockTemplateRepo) GetAll(ctx context.Context) ([]ClosingTemplate, error) {
	args := m.Called(ctx)
	return args.Get(0).([]ClosingTemplate), args.Error(1)
}
func (m *mockTemplateRepo) Update(ctx context.Context, template *ClosingTemplate) error {
	args := m.Called(ctx, template)
	return args.Error(0)
}
func (m *mockTemplateRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ─── COA Mock Repos ──────────────────────────────────────────────

type mockApprovalRepo struct {
	mock.Mock
}

func (m *mockApprovalRepo) Create(ctx context.Context, req *ApprovalRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}
func (m *mockApprovalRepo) GetByID(ctx context.Context, id string) (*ApprovalRequest, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ApprovalRequest), args.Error(1)
}
func (m *mockApprovalRepo) GetByStatus(ctx context.Context, status ApprovalStatus) ([]ApprovalRequest, error) {
	args := m.Called(ctx, status)
	return args.Get(0).([]ApprovalRequest), args.Error(1)
}
func (m *mockApprovalRepo) GetByEntity(ctx context.Context, entityType, entityID string) ([]ApprovalRequest, error) {
	args := m.Called(ctx, entityType, entityID)
	return args.Get(0).([]ApprovalRequest), args.Error(1)
}
func (m *mockApprovalRepo) UpdateStatus(ctx context.Context, id string, status ApprovalStatus, reviewedBy, reviewNote string) error {
	args := m.Called(ctx, id, status, reviewedBy, reviewNote)
	return args.Error(0)
}
func (m *mockApprovalRepo) GetAll(ctx context.Context) ([]ApprovalRequest, error) {
	args := m.Called(ctx)
	return args.Get(0).([]ApprovalRequest), args.Error(1)
}

type mockAccountVersionRepo struct {
	mock.Mock
}

func (m *mockAccountVersionRepo) Create(ctx context.Context, ver *AccountVersion) error {
	args := m.Called(ctx, ver)
	return args.Error(0)
}
func (m *mockAccountVersionRepo) GetByVersionNumber(ctx context.Context, versionNumber string) (*AccountVersion, error) {
	args := m.Called(ctx, versionNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AccountVersion), args.Error(1)
}
func (m *mockAccountVersionRepo) GetLatest(ctx context.Context) (*AccountVersion, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AccountVersion), args.Error(1)
}
func (m *mockAccountVersionRepo) GetAll(ctx context.Context) ([]AccountVersion, error) {
	args := m.Called(ctx)
	return args.Get(0).([]AccountVersion), args.Error(1)
}

type mockAccountMappingRepo struct {
	mock.Mock
}

func (m *mockAccountMappingRepo) Create(ctx context.Context, mp *AccountMapping) error {
	args := m.Called(ctx, mp)
	return args.Error(0)
}
func (m *mockAccountMappingRepo) GetByOldCode(ctx context.Context, sourceRegime, oldCode string) (*AccountMapping, error) {
	args := m.Called(ctx, sourceRegime, oldCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AccountMapping), args.Error(1)
}
func (m *mockAccountMappingRepo) GetByRegime(ctx context.Context, sourceRegime, targetRegime string) ([]AccountMapping, error) {
	args := m.Called(ctx, sourceRegime, targetRegime)
	return args.Get(0).([]AccountMapping), args.Error(1)
}
func (m *mockAccountMappingRepo) GetAll(ctx context.Context) ([]AccountMapping, error) {
	args := m.Called(ctx)
	return args.Get(0).([]AccountMapping), args.Error(1)
}

type mockAccountAnalysisRepo struct {
	mock.Mock
}

func (m *mockAccountAnalysisRepo) Create(ctx context.Context, a *AccountAnalysis) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}
func (m *mockAccountAnalysisRepo) GetByAccount(ctx context.Context, accountCode string) (*AccountAnalysis, error) {
	args := m.Called(ctx, accountCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AccountAnalysis), args.Error(1)
}
func (m *mockAccountAnalysisRepo) Update(ctx context.Context, a *AccountAnalysis) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

type mockIFRSMappingRepo struct {
	mock.Mock
}

func (m *mockIFRSMappingRepo) Create(ctx context.Context, mp *IFRSMapping) error {
	args := m.Called(ctx, mp)
	return args.Error(0)
}
func (m *mockIFRSMappingRepo) GetByVASCode(ctx context.Context, vasCode string) (*IFRSMapping, error) {
	args := m.Called(ctx, vasCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*IFRSMapping), args.Error(1)
}
func (m *mockIFRSMappingRepo) GetAll(ctx context.Context) ([]IFRSMapping, error) {
	args := m.Called(ctx)
	return args.Get(0).([]IFRSMapping), args.Error(1)
}
func (m *mockIFRSMappingRepo) Update(ctx context.Context, mp *IFRSMapping) error {
	args := m.Called(ctx, mp)
	return args.Error(0)
}

// ─── NewService helper ───────────────────────────────────────────

type serviceOption func(*service)

func newTestServiceWithCOA(
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
) *service {
	return NewService(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo, approvalRepo, versionRepo, mappingRepo, analysisRepo, ifrsRepo).(*service)
}

// Compile checks: ensure mock types implement interfaces
var _ AccountRepository = (*mockAccountRepo)(nil)
var _ JournalRepository = (*mockJournalRepo)(nil)
var _ PeriodRepository = (*mockPeriodRepo)(nil)
var _ UserRepository = (*mockUserRepo)(nil)
var _ AuditLogRepository = (*mockAuditRepo)(nil)
var _ ExchangeRateRepository = (*mockRateRepo)(nil)
var _ ClosingTemplateRepository = (*mockTemplateRepo)(nil)
var _ ApprovalRepository = (*mockApprovalRepo)(nil)
var _ AccountVersionRepository = (*mockAccountVersionRepo)(nil)
var _ AccountMappingRepository = (*mockAccountMappingRepo)(nil)
var _ AccountAnalysisRepository = (*mockAccountAnalysisRepo)(nil)
var _ IFRSMappingRepository = (*mockIFRSMappingRepo)(nil)

func TestAccountService_Create(t *testing.T) {
	accRepo := new(mockAccountRepo)
	jeRepo := new(mockJournalRepo)
	perRepo := new(mockPeriodRepo)
	userRepo := new(mockUserRepo)
	auditRepo := new(mockAuditRepo)
	rateRepo := new(mockRateRepo)
	templateRepo := new(mockTemplateRepo)
	svc := newTestServiceWithCOA(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo, nil, nil, nil, nil, nil)

	t.Run("success create account", func(t *testing.T) {
		acc := &Account{Code: "1111", Name: "TM VND", Type: AccountTypeAsset, IsActive: true}
		accRepo.On("GetByCode", mock.Anything, "1111").Return(nil, ErrAccountNotFound).Once()
		accRepo.On("Create", mock.Anything, acc).Return(nil).Once()

		err := svc.CreateAccount(context.Background(), acc)
		assert.NoError(t, err)
		accRepo.AssertExpectations(t)
	})

	t.Run("fail - duplicate code", func(t *testing.T) {
		existing := &Account{Code: "1111", Name: "Existing"}
		acc := &Account{Code: "1111", Name: "New", Type: AccountTypeAsset, IsActive: true}
		accRepo.On("GetByCode", mock.Anything, "1111").Return(existing, nil).Once()

		err := svc.CreateAccount(context.Background(), acc)
		assert.ErrorIs(t, err, ErrAccountCodeExists)
		accRepo.AssertExpectations(t)
	})

	t.Run("fail - invalid account data", func(t *testing.T) {
		acc := &Account{Code: "", Name: "", Type: AccountTypeAsset, IsActive: true}
		err := svc.CreateAccount(context.Background(), acc)
		assert.Error(t, err)
	})
}

func TestAccountService_Delete(t *testing.T) {
	accRepo := new(mockAccountRepo)
	jeRepo := new(mockJournalRepo)
	perRepo := new(mockPeriodRepo)
	userRepo := new(mockUserRepo)
	auditRepo := new(mockAuditRepo)
	rateRepo := new(mockRateRepo)
	templateRepo := new(mockTemplateRepo)
	svc := newTestServiceWithCOA(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo, nil, nil, nil, nil, nil)

	t.Run("success delete leaf account", func(t *testing.T) {
		acc := &Account{Code: "1112", Name: "Test", Type: AccountTypeAsset, IsActive: true}
		accRepo.On("GetByCode", mock.Anything, "1112").Return(acc, nil).Once()
		accRepo.On("GetChildren", mock.Anything, "1112").Return([]Account{}, nil).Once()
		accRepo.On("Delete", mock.Anything, "1112").Return(nil).Once()

		err := svc.DeleteAccount(context.Background(), "1112")
		assert.NoError(t, err)
		accRepo.AssertExpectations(t)
	})

	t.Run("fail - account not found", func(t *testing.T) {
		accRepo.On("GetByCode", mock.Anything, "9999").Return(nil, ErrAccountNotFound).Once()
		err := svc.DeleteAccount(context.Background(), "9999")
		assert.ErrorIs(t, err, ErrAccountNotFound)
	})

	t.Run("fail - has children", func(t *testing.T) {
		acc := &Account{Code: "111", Name: "Parent", Type: AccountTypeAsset, IsActive: true}
		accRepo.On("GetByCode", mock.Anything, "111").Return(acc, nil).Once()
		accRepo.On("GetChildren", mock.Anything, "111").Return([]Account{{Code: "1111"}}, nil).Once()

		err := svc.DeleteAccount(context.Background(), "111")
		assert.ErrorIs(t, err, ErrAccountHasChildren)
	})
}

func TestJournalService_CreateEntry(t *testing.T) {
	accRepo := new(mockAccountRepo)
	jeRepo := new(mockJournalRepo)
	perRepo := new(mockPeriodRepo)
	userRepo := new(mockUserRepo)
	auditRepo := new(mockAuditRepo)
	rateRepo := new(mockRateRepo)
	templateRepo := new(mockTemplateRepo)
	now := time.Now()
	svc := newTestServiceWithCOA(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo, nil, nil, nil, nil, nil)

	t.Run("success create draft entry", func(t *testing.T) {
		entry := &JournalEntry{
			EntryDate:   now,
			Description: "Test entry",
			Lines: []JournalLine{
				{AccountCode: "1111", DebitAmount: 100000, CreditAmount: 0, Description: "Debit"},
				{AccountCode: "5111", DebitAmount: 0, CreditAmount: 100000, Description: "Credit"},
			},
		}

		accRepo.On("GetByCode", mock.Anything, "1111").Return(&Account{Code: "1111", Type: AccountTypeAsset, IsActive: true}, nil).Once()
		accRepo.On("GetByCode", mock.Anything, "5111").Return(&Account{Code: "5111", Type: AccountTypeRevenue, IsActive: true}, nil).Once()
		jeRepo.On("Create", mock.Anything, entry).Return(nil).Once()

		err := svc.CreateEntry(context.Background(), entry, "user-1")
		assert.NoError(t, err)
		assert.Equal(t, JournalEntryDraft, entry.Status)
		assert.Equal(t, "user-1", entry.CreatedBy)
		accRepo.AssertExpectations(t)
		jeRepo.AssertExpectations(t)
	})

	t.Run("fail - inactive account", func(t *testing.T) {
		entry := &JournalEntry{
			EntryDate:   now,
			Description: "Test",
			Lines: []JournalLine{
				{AccountCode: "1111", DebitAmount: 100000, CreditAmount: 0},
				{AccountCode: "5111", DebitAmount: 0, CreditAmount: 100000},
			},
		}

		accRepo.On("GetByCode", mock.Anything, "1111").Return(&Account{Code: "1111", Type: AccountTypeAsset, IsActive: false}, nil).Once()

		err := svc.CreateEntry(context.Background(), entry, "user-1")
		assert.ErrorIs(t, err, ErrAccountInactive)
	})

	t.Run("fail - unknown account code", func(t *testing.T) {
		entry := &JournalEntry{
			EntryDate:   now,
			Description: "Test",
			Lines: []JournalLine{
				{AccountCode: "9999", DebitAmount: 100000, CreditAmount: 0},
				{AccountCode: "5111", DebitAmount: 0, CreditAmount: 100000},
			},
		}

		accRepo.On("GetByCode", mock.Anything, "9999").Return(nil, ErrAccountNotFound).Once()

		err := svc.CreateEntry(context.Background(), entry, "user-1")
		assert.ErrorIs(t, err, ErrAccountNotFound)
	})

	t.Run("fail - invalid entry data", func(t *testing.T) {
		entry := &JournalEntry{
			EntryDate:   now,
			Description: "",
			Lines: []JournalLine{
				{AccountCode: "1111", DebitAmount: 100000, CreditAmount: 0},
				{AccountCode: "5111", DebitAmount: 0, CreditAmount: 100000},
			},
		}

		err := svc.CreateEntry(context.Background(), entry, "user-1")
		assert.Error(t, err)
	})
}

func TestJournalService_PostEntry(t *testing.T) {
	accRepo := new(mockAccountRepo)
	jeRepo := new(mockJournalRepo)
	perRepo := new(mockPeriodRepo)
	userRepo := new(mockUserRepo)
	auditRepo := new(mockAuditRepo)
	rateRepo := new(mockRateRepo)
	templateRepo := new(mockTemplateRepo)
	now := time.Now()
	svc := newTestServiceWithCOA(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo, nil, nil, nil, nil, nil)

	t.Run("success full lifecycle", func(t *testing.T) {
		entry := &JournalEntry{
			EntryDate:   now,
			Description: "Test entry",
			PeriodID:    "P-2026-07",
			Lines: []JournalLine{
				{AccountCode: "1111", DebitAmount: 100000, CreditAmount: 0, Description: "Debit"},
				{AccountCode: "5111", DebitAmount: 0, CreditAmount: 100000, Description: "Credit"},
			},
		}

		// CreateEntry: validate accounts, create, set ID
		accRepo.On("GetByCode", mock.Anything, "1111").Return(&Account{Code: "1111", Type: AccountTypeAsset, IsActive: true}, nil).Once()
		accRepo.On("GetByCode", mock.Anything, "5111").Return(&Account{Code: "5111", Type: AccountTypeRevenue, IsActive: true}, nil).Once()
		jeRepo.On("Create", mock.Anything, entry).Return(nil).Run(func(args mock.Arguments) {
			e := args.Get(1).(*JournalEntry)
			e.ID = "JE-001"
		}).Once()

		err := svc.CreateEntry(context.Background(), entry, "user-1")
		assert.NoError(t, err)
		assert.Equal(t, JournalEntryDraft, entry.Status)
		assert.Equal(t, "JE-001", entry.ID)

		// SubmitForReview
		jeRepo.On("GetByID", mock.Anything, "JE-001").Return(entry, nil).Once()
		jeRepo.On("Review", mock.Anything, "JE-001", "user-1").Return(nil).Once()

		err = svc.SubmitForReview(context.Background(), "JE-001", "user-1")
		assert.NoError(t, err)

		// Simulate repo state change
		entry.Status = JournalEntryReviewing

		// ApproveEntry
		jeRepo.On("GetByID", mock.Anything, "JE-001").Return(entry, nil).Once()
		jeRepo.On("Approve", mock.Anything, "JE-001", "approver-1").Return(nil).Once()

		err = svc.ApproveEntry(context.Background(), "JE-001", "approver-1")
		assert.NoError(t, err)

		// Simulate repo state change
		entry.Status = JournalEntryApproved

		// PostEntry
		jeRepo.On("GetByID", mock.Anything, "JE-001").Return(entry, nil).Once()
		perRepo.On("GetByID", mock.Anything, "P-2026-07").Return(&Period{ID: "P-2026-07", Status: PeriodOpen}, nil).Once()
		jeRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()

		err = svc.PostEntry(context.Background(), "JE-001")
		assert.NoError(t, err)
		assert.Equal(t, JournalEntryPosted, entry.Status)
		assert.NotNil(t, entry.PostedAt)

		accRepo.AssertExpectations(t)
		jeRepo.AssertExpectations(t)
		perRepo.AssertExpectations(t)
	})

	t.Run("fail - already posted", func(t *testing.T) {
		entry := &JournalEntry{
			ID:     "JE-002",
			Status: JournalEntryPosted,
		}

		jeRepo.On("GetByID", mock.Anything, "JE-002").Return(entry, nil).Once()

		err := svc.PostEntry(context.Background(), "JE-002")
		assert.ErrorIs(t, err, ErrJournalAlreadyPosted)
	})

	t.Run("fail - cancelled entry", func(t *testing.T) {
		entry := &JournalEntry{
			ID:     "JE-003",
			Status: JournalEntryCancelled,
		}

		jeRepo.On("GetByID", mock.Anything, "JE-003").Return(entry, nil).Once()

		err := svc.PostEntry(context.Background(), "JE-003")
		assert.ErrorIs(t, err, ErrJournalAlreadyCancelled)
	})

	t.Run("fail - draft cannot be posted directly", func(t *testing.T) {
		entry := &JournalEntry{
			ID:     "JE-004",
			Status: JournalEntryDraft,
		}

		jeRepo.On("GetByID", mock.Anything, "JE-004").Return(entry, nil).Once()

		err := svc.PostEntry(context.Background(), "JE-004")
		assert.ErrorContains(t, err, "must be APPROVED")
	})

	t.Run("fail - period not found", func(t *testing.T) {
		entry := &JournalEntry{
			ID:       "JE-005",
			Status:   JournalEntryApproved,
			PeriodID: "P-2099-01",
		}

		jeRepo.On("GetByID", mock.Anything, "JE-005").Return(entry, nil).Once()
		perRepo.On("GetByID", mock.Anything, "P-2099-01").Return(nil, ErrPeriodNotFound).Once()

		err := svc.PostEntry(context.Background(), "JE-005")
		assert.ErrorIs(t, err, ErrPeriodNotFound)
	})

	t.Run("fail - period closed", func(t *testing.T) {
		entry := &JournalEntry{
			ID:       "JE-006",
			Status:   JournalEntryApproved,
			PeriodID: "P-2026-06",
		}

		jeRepo.On("GetByID", mock.Anything, "JE-006").Return(entry, nil).Once()
		perRepo.On("GetByID", mock.Anything, "P-2026-06").Return(&Period{ID: "P-2026-06", Status: PeriodClosed}, nil).Once()

		err := svc.PostEntry(context.Background(), "JE-006")
		assert.ErrorIs(t, err, ErrJournalPeriodClosed)
	})
}

func TestJournalService_CancelEntry(t *testing.T) {
	accRepo := new(mockAccountRepo)
	jeRepo := new(mockJournalRepo)
	perRepo := new(mockPeriodRepo)
	userRepo := new(mockUserRepo)
	auditRepo := new(mockAuditRepo)
	rateRepo := new(mockRateRepo)
	templateRepo := new(mockTemplateRepo)
	svc := newTestServiceWithCOA(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo, nil, nil, nil, nil, nil)

	t.Run("success cancel posted entry", func(t *testing.T) {
		entry := &JournalEntry{ID: "JE-001", Status: JournalEntryPosted}

		jeRepo.On("GetByID", mock.Anything, "JE-001").Return(entry, nil).Once()
		jeRepo.On("UpdateStatus", mock.Anything, "JE-001", JournalEntryCancelled).Return(nil).Once()

		err := svc.CancelEntry(context.Background(), "JE-001")
		assert.NoError(t, err)
		jeRepo.AssertExpectations(t)
	})

	t.Run("fail - already cancelled", func(t *testing.T) {
		entry := &JournalEntry{ID: "JE-002", Status: JournalEntryCancelled}
		jeRepo.On("GetByID", mock.Anything, "JE-002").Return(entry, nil).Once()

		err := svc.CancelEntry(context.Background(), "JE-002")
		assert.ErrorIs(t, err, ErrJournalAlreadyCancelled)
	})

	t.Run("fail - draft cannot cancel", func(t *testing.T) {
		entry := &JournalEntry{ID: "JE-003", Status: JournalEntryDraft}
		jeRepo.On("GetByID", mock.Anything, "JE-003").Return(entry, nil).Once()

		err := svc.CancelEntry(context.Background(), "JE-003")
		assert.ErrorIs(t, err, ErrJournalAlreadyDraft)
	})
}

func TestJournalService_TrialBalance(t *testing.T) {
	accRepo := new(mockAccountRepo)
	jeRepo := new(mockJournalRepo)
	perRepo := new(mockPeriodRepo)
	userRepo := new(mockUserRepo)
	auditRepo := new(mockAuditRepo)
	rateRepo := new(mockRateRepo)
	templateRepo := new(mockTemplateRepo)
	svc := newTestServiceWithCOA(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo, nil, nil, nil, nil, nil)

	t.Run("success get trial balance", func(t *testing.T) {
		balances := []AccountBalance{
			{AccountCode: "1111", AccountType: AccountTypeAsset, PeriodDebit: 1000000, PeriodCredit: 300000},
			{AccountCode: "5111", AccountType: AccountTypeRevenue, PeriodDebit: 0, PeriodCredit: 700000},
		}

		perRepo.On("GetByYearMonth", mock.Anything, 2026, 7).Return(&Period{ID: "P-2026-07", Status: PeriodOpen}, nil).Once()
		jeRepo.On("GetTrialBalance", mock.Anything, "P-2026-07").Return(balances, nil).Once()

		result, err := svc.TrialBalance(context.Background(), 2026, 7)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, 700000.0, result[0].ClosingBalance) // 1000000 - 300000 = 700000
		perRepo.AssertExpectations(t)
		jeRepo.AssertExpectations(t)
	})

	t.Run("fail - period not found", func(t *testing.T) {
		perRepo.On("GetByYearMonth", mock.Anything, 2099, 1).Return(nil, ErrPeriodNotFound).Once()

		_, err := svc.TrialBalance(context.Background(), 2099, 1)
		assert.ErrorIs(t, err, ErrPeriodNotFound)
	})
}

func TestAccountService_Update(t *testing.T) {
	accRepo := new(mockAccountRepo)
	jeRepo := new(mockJournalRepo)
	perRepo := new(mockPeriodRepo)
	userRepo := new(mockUserRepo)
	auditRepo := new(mockAuditRepo)
	rateRepo := new(mockRateRepo)
	templateRepo := new(mockTemplateRepo)
	svc := newTestServiceWithCOA(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo, nil, nil, nil, nil, nil)

	t.Run("success update account", func(t *testing.T) {
		existing := &Account{Code: "1111", Name: "Old Name", Type: AccountTypeAsset, IsActive: true}
		updated := &Account{Code: "1111", Name: "New Name", Type: AccountTypeAsset, IsActive: true}

		accRepo.On("GetByCode", mock.Anything, "1111").Return(existing, nil).Once()
		accRepo.On("Update", mock.Anything, updated).Return(nil).Once()

		err := svc.UpdateAccount(context.Background(), updated)
		assert.NoError(t, err)
		accRepo.AssertExpectations(t)
	})

	t.Run("fail - account not found", func(t *testing.T) {
		updated := &Account{Code: "9999", Name: "Ghost", Type: AccountTypeAsset, IsActive: true}
		accRepo.On("GetByCode", mock.Anything, "9999").Return(nil, ErrAccountNotFound).Once()

		err := svc.UpdateAccount(context.Background(), updated)
		assert.ErrorIs(t, err, ErrAccountNotFound)
	})
}

func TestJournalService_GetByDateRange(t *testing.T) {
	accRepo := new(mockAccountRepo)
	jeRepo := new(mockJournalRepo)
	perRepo := new(mockPeriodRepo)
	userRepo := new(mockUserRepo)
	auditRepo := new(mockAuditRepo)
	rateRepo := new(mockRateRepo)
	templateRepo := new(mockTemplateRepo)
	svc := newTestServiceWithCOA(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo, nil, nil, nil, nil, nil)
	now := time.Now()

	t.Run("success get entries by date range", func(t *testing.T) {
		from := now.Add(-7 * 24 * time.Hour)
		to := now
		entries := []JournalEntry{
			{ID: "JE-001", Description: "Entry 1"},
			{ID: "JE-002", Description: "Entry 2"},
		}

		jeRepo.On("GetByDateRange", mock.Anything, from, to).Return(entries, nil).Once()

		result, err := svc.GetEntriesByDateRange(context.Background(), from, to)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		jeRepo.AssertExpectations(t)
	})
}

func TestAccountService_GetAllAccounts(t *testing.T) {
	accRepo := new(mockAccountRepo)
	jeRepo := new(mockJournalRepo)
	perRepo := new(mockPeriodRepo)
	userRepo := new(mockUserRepo)
	auditRepo := new(mockAuditRepo)
	rateRepo := new(mockRateRepo)
	templateRepo := new(mockTemplateRepo)
	svc := newTestServiceWithCOA(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo, nil, nil, nil, nil, nil)

	t.Run("success get all active", func(t *testing.T) {
		accounts := []Account{
			{Code: "1111", Name: "TM VND", Type: AccountTypeAsset, IsActive: true},
			{Code: "5111", Name: "DT BH", Type: AccountTypeRevenue, IsActive: true},
		}

		accRepo.On("GetAll", mock.Anything, true).Return(accounts, nil).Once()

		result, err := svc.GetAllAccounts(context.Background(), true)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		accRepo.AssertExpectations(t)
	})
}

func TestAccountService_GetAccount(t *testing.T) {
	accRepo := new(mockAccountRepo)
	jeRepo := new(mockJournalRepo)
	perRepo := new(mockPeriodRepo)
	userRepo := new(mockUserRepo)
	auditRepo := new(mockAuditRepo)
	rateRepo := new(mockRateRepo)
	templateRepo := new(mockTemplateRepo)
	svc := newTestServiceWithCOA(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo, nil, nil, nil, nil, nil)

	t.Run("success get account by code", func(t *testing.T) {
		acc := &Account{Code: "1111", Name: "TM VND", Type: AccountTypeAsset}
		accRepo.On("GetByCode", mock.Anything, "1111").Return(acc, nil).Once()

		result, err := svc.GetAccount(context.Background(), "1111")
		assert.NoError(t, err)
		assert.Equal(t, "1111", result.Code)
		accRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		accRepo.On("GetByCode", mock.Anything, "9999").Return(nil, ErrAccountNotFound).Once()
		_, err := svc.GetAccount(context.Background(), "9999")
		assert.ErrorIs(t, err, ErrAccountNotFound)
	})
}

// ─── COA: Freeze/Unfreeze Tests ─────────────────────────────────

func TestCOA_FreezeAccount(t *testing.T) {
	accRepo := new(mockAccountRepo)
	svc := newTestServiceWithCOA(accRepo, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	t.Run("success freeze", func(t *testing.T) {
		acc := &Account{Code: "1111", Name: "TM", Type: AccountTypeAsset, Status: AccountStatusActive}
		accRepo.On("GetByCode", mock.Anything, "1111").Return(acc, nil).Once()
		accRepo.On("Update", mock.Anything, mock.MatchedBy(func(a *Account) bool {
			return a.Status == AccountStatusFrozen && a.FreezeReason == "audit hold"
		})).Return(nil).Once()

		err := svc.FreezeAccount(context.Background(), "1111", "audit hold")
		assert.NoError(t, err)
		accRepo.AssertExpectations(t)
	})

	t.Run("fail - already frozen", func(t *testing.T) {
		acc := &Account{Code: "1111", Name: "TM", Type: AccountTypeAsset, Status: AccountStatusFrozen}
		accRepo.On("GetByCode", mock.Anything, "1111").Return(acc, nil).Once()

		err := svc.FreezeAccount(context.Background(), "1111", "reason")
		assert.ErrorIs(t, err, ErrAccountAlreadyFrozen)
	})
}

func TestCOA_UnfreezeAccount(t *testing.T) {
	accRepo := new(mockAccountRepo)
	svc := newTestServiceWithCOA(accRepo, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	t.Run("success unfreeze", func(t *testing.T) {
		acc := &Account{Code: "1111", Name: "TM", Type: AccountTypeAsset, Status: AccountStatusFrozen, FreezeReason: "hold"}
		accRepo.On("GetByCode", mock.Anything, "1111").Return(acc, nil).Once()
		accRepo.On("Update", mock.Anything, mock.MatchedBy(func(a *Account) bool {
			return a.Status == AccountStatusActive && a.FreezeReason == ""
		})).Return(nil).Once()

		err := svc.UnfreezeAccount(context.Background(), "1111", "resolved")
		assert.NoError(t, err)
		accRepo.AssertExpectations(t)
	})

	t.Run("fail - not frozen", func(t *testing.T) {
		acc := &Account{Code: "1111", Name: "TM", Type: AccountTypeAsset, Status: AccountStatusActive}
		accRepo.On("GetByCode", mock.Anything, "1111").Return(acc, nil).Once()

		err := svc.UnfreezeAccount(context.Background(), "1111", "reason")
		assert.ErrorIs(t, err, ErrAccountNotFrozen)
	})
}

// ─── COA: Approval Tests ─────────────────────────────────────────

func TestCOA_ApproveRequest(t *testing.T) {
	approvalRepo := new(mockApprovalRepo)
	svc := newTestServiceWithCOA(nil, nil, nil, nil, nil, nil, nil, approvalRepo, nil, nil, nil, nil)

	t.Run("success approve", func(t *testing.T) {
		req := &ApprovalRequest{
			ID: "APPR-001", Status: ApprovalPending, RequestedBy: "user-1",
		}
		approvalRepo.On("GetByID", mock.Anything, "APPR-001").Return(req, nil).Once()
		approvalRepo.On("UpdateStatus", mock.Anything, "APPR-001", ApprovalApproved, "reviewer-1", "looks good").Return(nil).Once()

		err := svc.ApproveRequest(context.Background(), "APPR-001", "reviewer-1", "looks good")
		assert.NoError(t, err)
		approvalRepo.AssertExpectations(t)
	})

	t.Run("fail - self approval", func(t *testing.T) {
		req := &ApprovalRequest{
			ID: "APPR-002", Status: ApprovalPending, RequestedBy: "same-user",
		}
		approvalRepo.On("GetByID", mock.Anything, "APPR-002").Return(req, nil).Once()

		err := svc.ApproveRequest(context.Background(), "APPR-002", "same-user", "ok")
		assert.ErrorIs(t, err, ErrSelfApproval)
	})

	t.Run("fail - already processed", func(t *testing.T) {
		req := &ApprovalRequest{
			ID: "APPR-003", Status: ApprovalApproved, RequestedBy: "user-1",
		}
		approvalRepo.On("GetByID", mock.Anything, "APPR-003").Return(req, nil).Once()

		err := svc.ApproveRequest(context.Background(), "APPR-003", "reviewer-1", "ok")
		assert.ErrorIs(t, err, ErrApprovalAlreadyProcessed)
	})
}

func TestCOA_RejectRequest(t *testing.T) {
	approvalRepo := new(mockApprovalRepo)
	svc := newTestServiceWithCOA(nil, nil, nil, nil, nil, nil, nil, approvalRepo, nil, nil, nil, nil)

	t.Run("success reject", func(t *testing.T) {
		req := &ApprovalRequest{
			ID: "APPR-001", Status: ApprovalPending, RequestedBy: "user-1",
		}
		approvalRepo.On("GetByID", mock.Anything, "APPR-001").Return(req, nil).Once()
		approvalRepo.On("UpdateStatus", mock.Anything, "APPR-001", ApprovalRejected, "reviewer-1", "needs more info").Return(nil).Once()

		err := svc.RejectRequest(context.Background(), "APPR-001", "reviewer-1", "needs more info")
		assert.NoError(t, err)
	})

	t.Run("fail - empty note", func(t *testing.T) {
		req := &ApprovalRequest{
			ID: "APPR-002", Status: ApprovalPending, RequestedBy: "user-1",
		}
		approvalRepo.On("GetByID", mock.Anything, "APPR-002").Return(req, nil).Once()

		err := svc.RejectRequest(context.Background(), "APPR-002", "reviewer-1", "")
		assert.ErrorIs(t, err, ErrReviewNoteRequired)
	})

	t.Run("fail - already processed", func(t *testing.T) {
		req := &ApprovalRequest{
			ID: "APPR-003", Status: ApprovalRejected, RequestedBy: "user-1",
		}
		approvalRepo.On("GetByID", mock.Anything, "APPR-003").Return(req, nil).Once()

		err := svc.RejectRequest(context.Background(), "APPR-003", "reviewer-1", "note")
		assert.ErrorIs(t, err, ErrApprovalAlreadyProcessed)
	})
}

func TestCOA_CreateApprovalRequest(t *testing.T) {
	approvalRepo := new(mockApprovalRepo)
	svc := newTestServiceWithCOA(nil, nil, nil, nil, nil, nil, nil, approvalRepo, nil, nil, nil, nil)

	t.Run("success create", func(t *testing.T) {
		req := &ApprovalRequest{
			EntityType: "ACCOUNT", EntityID: "1111", RequestType: "FREEZE",
			Reason: "audit", RequestedBy: "user-1",
		}
		approvalRepo.On("Create", mock.Anything, req).Return(nil).Once()

		err := svc.CreateApprovalRequest(context.Background(), req)
		assert.NoError(t, err)
	})

	t.Run("fail - missing reason", func(t *testing.T) {
		req := &ApprovalRequest{
			EntityType: "ACCOUNT", EntityID: "1111", RequestType: "FREEZE",
			Reason: "", RequestedBy: "user-1",
		}
		err := svc.CreateApprovalRequest(context.Background(), req)
		assert.ErrorIs(t, err, ErrApprovalReasonRequired)
	})
}

// ─── COA: Versioning Tests ───────────────────────────────────────

func TestCOA_CreateAccountVersion(t *testing.T) {
	accRepo := new(mockAccountRepo)
	verRepo := new(mockAccountVersionRepo)
	svc := newTestServiceWithCOA(accRepo, nil, nil, nil, nil, nil, nil, nil, verRepo, nil, nil, nil)

	t.Run("success create version v1", func(t *testing.T) {
		accounts := []Account{
			{Code: "1111", Name: "TM", Type: AccountTypeAsset, Status: AccountStatusActive},
		}
		accRepo.On("GetAll", mock.Anything, false).Return(accounts, nil).Once()
		verRepo.On("GetLatest", mock.Anything).Return((*AccountVersion)(nil), ErrVersionNotFound).Once()
		verRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()

		ver, err := svc.CreateAccountVersion(context.Background(), "initial version")
		assert.NoError(t, err)
		assert.Equal(t, "v1", ver.VersionNumber)
		assert.Equal(t, "initial version", ver.ChangeReason)
	})

	t.Run("success increment version", func(t *testing.T) {
		accounts := []Account{{Code: "1111", Name: "TM", Type: AccountTypeAsset, Status: AccountStatusActive}}
		existing := &AccountVersion{VersionNumber: "v1", CreatedAt: time.Now()}

		accRepo.On("GetAll", mock.Anything, false).Return(accounts, nil).Once()
		verRepo.On("GetLatest", mock.Anything).Return(existing, nil).Once()
		verRepo.On("Create", mock.Anything, mock.MatchedBy(func(v *AccountVersion) bool {
			return v.VersionNumber == "v2"
		})).Return(nil).Once()

		ver, err := svc.CreateAccountVersion(context.Background(), "version 2")
		assert.NoError(t, err)
		assert.Equal(t, "v2", ver.VersionNumber)
	})
}

func TestCOA_CompareVersions(t *testing.T) {
	verRepo := new(mockAccountVersionRepo)
	svc := newTestServiceWithCOA(nil, nil, nil, nil, nil, nil, nil, nil, verRepo, nil, nil, nil)

	t.Run("detect added/removed/modified", func(t *testing.T) {
		v1 := &AccountVersion{
			VersionNumber: "v1",
			Snapshot:      `[{"code":"1111","name":"TM","type":"ASSET","status":"ACTIVE"}]`,
		}
		v2 := &AccountVersion{
			VersionNumber: "v2",
			Snapshot:      `[{"code":"1111","name":"Cash","type":"ASSET","status":"ACTIVE"},{"code":"5111","name":"Revenue","type":"REVENUE","status":"ACTIVE"}]`,
		}

		verRepo.On("GetByVersionNumber", mock.Anything, "v1").Return(v1, nil).Once()
		verRepo.On("GetByVersionNumber", mock.Anything, "v2").Return(v2, nil).Once()

		diff, err := svc.CompareVersions(context.Background(), "v1", "v2")
		assert.NoError(t, err)
		assert.Equal(t, "v1", diff.VersionFrom)
		assert.Equal(t, "v2", diff.VersionTo)
		assert.Len(t, diff.Added, 1)
		assert.Equal(t, "5111", diff.Added[0].Code)
		assert.Len(t, diff.Modified, 1)
		assert.Equal(t, "1111", diff.Modified[0].Code)
		assert.Contains(t, diff.Modified[0].Changes, "name")
	})
}

// ─── COA: Account Mapping Tests ──────────────────────────────────

func TestCOA_AccountMapping(t *testing.T) {
	mappingRepo := new(mockAccountMappingRepo)
	svc := newTestServiceWithCOA(nil, nil, nil, nil, nil, nil, nil, nil, nil, mappingRepo, nil, nil)

	t.Run("success create mapping", func(t *testing.T) {
		m := &AccountMapping{
			SourceRegime: "TT200", TargetRegime: "TT99", OldCode: "611", NewCode: "632",
			MappingType: "DIRECT",
		}
		mappingRepo.On("GetByOldCode", mock.Anything, "TT200", "611").Return((*AccountMapping)(nil), ErrMappingNotFound).Once()
		mappingRepo.On("Create", mock.Anything, m).Return(nil).Once()

		err := svc.CreateAccountMapping(context.Background(), m)
		assert.NoError(t, err)
	})

	t.Run("fail - duplicate", func(t *testing.T) {
		existing := &AccountMapping{OldCode: "611"}
		m := &AccountMapping{
			SourceRegime: "TT200", TargetRegime: "TT99", OldCode: "611", NewCode: "632",
			MappingType: "DIRECT",
		}
		mappingRepo.On("GetByOldCode", mock.Anything, "TT200", "611").Return(existing, nil).Once()

		err := svc.CreateAccountMapping(context.Background(), m)
		assert.ErrorIs(t, err, ErrMappingExists)
	})

	t.Run("success get mapping", func(t *testing.T) {
		m := &AccountMapping{OldCode: "611", NewCode: "632"}
		mappingRepo.On("GetByOldCode", mock.Anything, "TT200", "611").Return(m, nil).Once()

		result, err := svc.GetMappingByOldCode(context.Background(), "TT200", "611")
		assert.NoError(t, err)
		assert.Equal(t, "611", result.OldCode)
	})
}

// ─── COA: Account Analysis Tests ─────────────────────────────────

func TestCOA_AccountAnalysis(t *testing.T) {
	accRepo := new(mockAccountRepo)
	analysisRepo := new(mockAccountAnalysisRepo)
	svc := newTestServiceWithCOA(accRepo, nil, nil, nil, nil, nil, nil, nil, nil, nil, analysisRepo, nil)

	t.Run("success create analysis", func(t *testing.T) {
		a := &AccountAnalysis{AccountCode: "1111", CostCenterID: "CC-001"}
		accRepo.On("GetByCode", mock.Anything, "1111").Return(&Account{Code: "1111"}, nil).Once()
		analysisRepo.On("GetByAccount", mock.Anything, "1111").Return((*AccountAnalysis)(nil), ErrAnalysisNotFound).Once()
		analysisRepo.On("Create", mock.Anything, a).Return(nil).Once()

		err := svc.CreateAccountAnalysis(context.Background(), a)
		assert.NoError(t, err)
	})

	t.Run("upsert on existing", func(t *testing.T) {
		existing := &AccountAnalysis{AccountCode: "1111", CostCenterID: "CC-001"}
		updated := &AccountAnalysis{AccountCode: "1111", CostCenterID: "CC-002"}

		accRepo.On("GetByCode", mock.Anything, "1111").Return(&Account{Code: "1111"}, nil).Once()
		analysisRepo.On("GetByAccount", mock.Anything, "1111").Return(existing, nil).Once()
		analysisRepo.On("Update", mock.Anything, updated).Return(nil).Once()

		err := svc.CreateAccountAnalysis(context.Background(), updated)
		assert.NoError(t, err)
	})
}

// ─── COA: IFRS Mapping Tests ─────────────────────────────────────

func TestCOA_IFRSMapping(t *testing.T) {
	ifrsRepo := new(mockIFRSMappingRepo)
	svc := newTestServiceWithCOA(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, ifrsRepo)

	t.Run("success create IFRS mapping", func(t *testing.T) {
		m := &IFRSMapping{VASCode: "1111", IFRSCode: "IFRS-1100", AdjustmentType: "RECLASSIFY"}
		ifrsRepo.On("GetByVASCode", mock.Anything, "1111").Return((*IFRSMapping)(nil), ErrIFRSMappingNotFound).Once()
		ifrsRepo.On("Create", mock.Anything, m).Return(nil).Once()

		err := svc.CreateIFRSMapping(context.Background(), m)
		assert.NoError(t, err)
	})

	t.Run("fail - duplicate", func(t *testing.T) {
		existing := &IFRSMapping{VASCode: "1111"}
		m := &IFRSMapping{VASCode: "1111", IFRSCode: "IFRS-1100"}
		ifrsRepo.On("GetByVASCode", mock.Anything, "1111").Return(existing, nil).Once()

		err := svc.CreateIFRSMapping(context.Background(), m)
		assert.ErrorIs(t, err, ErrIFRSMappingExists)
	})

	t.Run("success get", func(t *testing.T) {
		m := &IFRSMapping{VASCode: "1111", IFRSCode: "IFRS-1100"}
		ifrsRepo.On("GetByVASCode", mock.Anything, "1111").Return(m, nil).Once()

		result, err := svc.GetIFRSMapping(context.Background(), "1111")
		assert.NoError(t, err)
		assert.Equal(t, "IFRS-1100", result.IFRSCode)
	})
}

// ─── COA: Balance Drill-down Tests ───────────────────────────────

func TestCOA_GetAccountBalance(t *testing.T) {
	accRepo := new(mockAccountRepo)
	jeRepo := new(mockJournalRepo)
	svc := newTestServiceWithCOA(accRepo, jeRepo, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	t.Run("success get balance", func(t *testing.T) {
		accRepo.On("GetByCode", mock.Anything, "1111").Return(&Account{Code: "1111"}, nil).Once()
		jeRepo.On("GetBalance", mock.Anything, "1111", "P-2026-07").Return(&AccountBalance{AccountCode: "1111", ClosingBalance: 500000}, nil).Once()

		b, err := svc.GetAccountBalance(context.Background(), "1111", "P-2026-07")
		assert.NoError(t, err)
		assert.Equal(t, 500000.0, b.ClosingBalance)
	})

	t.Run("fail - account not found", func(t *testing.T) {
		accRepo.On("GetByCode", mock.Anything, "9999").Return(nil, ErrAccountNotFound).Once()
		_, err := svc.GetAccountBalance(context.Background(), "9999", "P-2026-07")
		assert.ErrorIs(t, err, ErrAccountNotFound)
	})
}
