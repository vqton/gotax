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

// Compile checks: ensure mock types implement interfaces
var _ AccountRepository = (*mockAccountRepo)(nil)
var _ JournalRepository = (*mockJournalRepo)(nil)
var _ PeriodRepository = (*mockPeriodRepo)(nil)
var _ UserRepository = (*mockUserRepo)(nil)
var _ AuditLogRepository = (*mockAuditRepo)(nil)
var _ ExchangeRateRepository = (*mockRateRepo)(nil)
var _ ClosingTemplateRepository = (*mockTemplateRepo)(nil)

func TestAccountService_Create(t *testing.T) {
	accRepo := new(mockAccountRepo)
	jeRepo := new(mockJournalRepo)
	perRepo := new(mockPeriodRepo)
	userRepo := new(mockUserRepo)
	auditRepo := new(mockAuditRepo)
	rateRepo := new(mockRateRepo)
	templateRepo := new(mockTemplateRepo)
	svc := NewService(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo)

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
	svc := NewService(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo)

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
	svc := NewService(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo)

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
	svc := NewService(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo)

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
	svc := NewService(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo)

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
	svc := NewService(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo)

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
	svc := NewService(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo)

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
	svc := NewService(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo)
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
	svc := NewService(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo)

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
	svc := NewService(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo)

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
