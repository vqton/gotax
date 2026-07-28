package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/auth"
	"gotax/internal/domain"
	"gotax/internal/repository"
)

func setupService(t *testing.T) (Service, context.Context) {
	t.Helper()
	accRepo := repository.NewMemoryAccountRepo()
	jeRepo := repository.NewMemoryJournalRepo()
	jeRepo.SetAccounts(accRepo.Accounts())
	perRepo := repository.NewMemoryPeriodRepo()
	userRepo := repository.NewMemoryUserRepo()
	auditRepo := repository.NewMemoryAuditLogRepo()
	rateRepo := repository.NewMemoryExchangeRateRepo()
	templateRepo := repository.NewMemoryClosingTemplateRepo()
	approvalRepo := repository.NewMemoryApprovalRepo()
	versionRepo := repository.NewMemoryAccountVersionRepo()
	mappingRepo := repository.NewMemoryAccountMappingRepo()
	analysisRepo := repository.NewMemoryAccountAnalysisRepo()
	ifrsRepo := repository.NewMemoryIFRSMappingRepo()
	refreshRepo := repository.NewMemoryRefreshTokenRepo()
	resetRepo := repository.NewMemoryPasswordResetTokenRepo()

	obRepo := repository.NewMemoryOpeningBalanceRepo()
	cashRepo := repository.NewMemoryCashRepo()
	svc := NewService(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo,
		approvalRepo, versionRepo, mappingRepo, analysisRepo, ifrsRepo, refreshRepo, resetRepo, obRepo, cashRepo)
	return svc, context.Background()
}

// ─── Accounts ───────────────────────────────────────────────────────

func TestCreateAccount_Success(t *testing.T) {
	svc, ctx := setupService(t)
	a := &domain.Account{Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset}
	assert.NoError(t, svc.CreateAccount(ctx, a))
}

func TestCreateAccount_Duplicate(t *testing.T) {
	svc, ctx := setupService(t)
	a := &domain.Account{Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset}
	require.NoError(t, svc.CreateAccount(ctx, a))
	assert.ErrorIs(t, svc.CreateAccount(ctx, a), domain.ErrAccountCodeExists)
}

func TestCreateAccount_Validation(t *testing.T) {
	svc, ctx := setupService(t)
	assert.ErrorIs(t, svc.CreateAccount(ctx, &domain.Account{}), domain.ErrAccountCodeRequired)
}

func TestGetAccount_Success(t *testing.T) {
	svc, ctx := setupService(t)
	a := &domain.Account{Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset}
	require.NoError(t, svc.CreateAccount(ctx, a))

	got, err := svc.GetAccount(ctx, "1111")
	require.NoError(t, err)
	assert.Equal(t, "Cash", got.Name)
}

func TestGetAccount_NotFound(t *testing.T) {
	svc, ctx := setupService(t)
	_, err := svc.GetAccount(ctx, "9999")
	assert.ErrorIs(t, err, domain.ErrAccountNotFound)
}

func TestGetAllAccounts(t *testing.T) {
	svc, ctx := setupService(t)
	require.NoError(t, svc.CreateAccount(ctx, &domain.Account{Code: "1111", Name: "A", Type: domain.AccountTypeAsset, IsActive: true}))
	require.NoError(t, svc.CreateAccount(ctx, &domain.Account{Code: "2111", Name: "B", Type: domain.AccountTypeLiability, IsActive: true}))
	require.NoError(t, svc.CreateAccount(ctx, &domain.Account{Code: "3111", Name: "C", Type: domain.AccountTypeEquity, IsActive: false}))

	all, err := svc.GetAllAccounts(ctx, false)
	require.NoError(t, err)
	assert.Len(t, all, 3)

	active, err := svc.GetAllAccounts(ctx, true)
	require.NoError(t, err)
	assert.Len(t, active, 2)
}

func TestUpdateAccount(t *testing.T) {
	svc, ctx := setupService(t)
	a := &domain.Account{Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset}
	require.NoError(t, svc.CreateAccount(ctx, a))

	a.Name = "Cash Updated"
	assert.NoError(t, svc.UpdateAccount(ctx, a))

	got, _ := svc.GetAccount(ctx, "1111")
	assert.Equal(t, "Cash Updated", got.Name)
}

func TestDeleteAccount(t *testing.T) {
	svc, ctx := setupService(t)
	require.NoError(t, svc.CreateAccount(ctx, &domain.Account{Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset}))

	// can't delete with children
	require.NoError(t, svc.CreateAccount(ctx, &domain.Account{Code: "11111", Name: "Sub", Type: domain.AccountTypeAsset, ParentCode: "1111"}))
	assert.ErrorIs(t, svc.DeleteAccount(ctx, "1111"), domain.ErrAccountHasChildren)

	assert.NoError(t, svc.DeleteAccount(ctx, "11111"))
	assert.NoError(t, svc.DeleteAccount(ctx, "1111"))
	_, err := svc.GetAccount(ctx, "1111")
	assert.ErrorIs(t, err, domain.ErrAccountNotFound)
}

func TestFreezeUnfreezeAccount(t *testing.T) {
	svc, ctx := setupService(t)
	require.NoError(t, svc.CreateAccount(ctx, &domain.Account{Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset}))

	assert.NoError(t, svc.FreezeAccount(ctx, "1111", "audit"))
	got, _ := svc.GetAccount(ctx, "1111")
	assert.Equal(t, domain.AccountStatusFrozen, got.Status)

	assert.NoError(t, svc.UnfreezeAccount(ctx, "1111", "done"))
	got, _ = svc.GetAccount(ctx, "1111")
	assert.Equal(t, domain.AccountStatusActive, got.Status)
}

// ─── Journal Entries ─────────────────────────────────────────────────

func TestCreateEntry_Success(t *testing.T) {
	svc, ctx := setupService(t)
	require.NoError(t, svc.CreateAccount(ctx, &domain.Account{Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset, Status: domain.AccountStatusActive, IsActive: true}))
	require.NoError(t, svc.CreateAccount(ctx, &domain.Account{Code: "5111", Name: "Expense", Type: domain.AccountTypeExpense, Status: domain.AccountStatusActive, IsActive: true}))

	je := &domain.JournalEntry{
		EntryDate:   time.Now(),
		Description: "Test entry",
		Lines: []domain.JournalLine{
			{LineNumber: 1, AccountCode: "1111", DebitAmount: 100, CreditAmount: 0},
			{LineNumber: 2, AccountCode: "5111", DebitAmount: 0, CreditAmount: 100},
		},
	}
	assert.NoError(t, svc.CreateEntry(ctx, je, "user1"))
	assert.Equal(t, domain.JournalEntryDraft, je.Status)
}

func TestCreateEntry_FrozenAccount(t *testing.T) {
	svc, ctx := setupService(t)
	require.NoError(t, svc.CreateAccount(ctx, &domain.Account{Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset, Status: domain.AccountStatusFrozen, IsActive: true, FreezeReason: "audit"}))

	je := &domain.JournalEntry{
		EntryDate:   time.Now(),
		Description: "Test",
		Lines: []domain.JournalLine{
			{LineNumber: 1, AccountCode: "1111", DebitAmount: 100, CreditAmount: 0},
			{LineNumber: 2, AccountCode: "5111", DebitAmount: 0, CreditAmount: 100},
		},
	}
	assert.Error(t, svc.CreateEntry(ctx, je, "user1"))
}

func TestJournalLifecycle(t *testing.T) {
	svc, ctx := setupService(t)
	require.NoError(t, svc.CreateAccount(ctx, &domain.Account{Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset, Status: domain.AccountStatusActive, IsActive: true}))
	require.NoError(t, svc.CreateAccount(ctx, &domain.Account{Code: "5111", Name: "Expense", Type: domain.AccountTypeExpense, Status: domain.AccountStatusActive, IsActive: true}))
	require.NoError(t, svc.CreatePeriod(ctx, &domain.Period{
		ID: "P-2025-06", Year: 2025, Month: 6,
		StartDate: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
		Status:    domain.PeriodOpen,
	}))

	je := &domain.JournalEntry{
		EntryDate:   time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
		PeriodID:    "P-2025-06",
		Description: "Lifecycle test",
		Lines: []domain.JournalLine{
			{LineNumber: 1, AccountCode: "1111", DebitAmount: 100, CreditAmount: 0},
			{LineNumber: 2, AccountCode: "5111", DebitAmount: 0, CreditAmount: 100},
		},
	}
	require.NoError(t, svc.CreateEntry(ctx, je, "user1"))
	id := je.ID

	assert.NoError(t, svc.SubmitForReview(ctx, id, "user1"))
	got, _ := svc.GetEntryByID(ctx, id)
	assert.Equal(t, domain.JournalEntryReviewing, got.Status)

	assert.NoError(t, svc.ApproveEntry(ctx, id, "approver1"))
	got, _ = svc.GetEntryByID(ctx, id)
	assert.Equal(t, domain.JournalEntryApproved, got.Status)

	assert.NoError(t, svc.PostEntry(ctx, id))
	got, _ = svc.GetEntryByID(ctx, id)
	assert.Equal(t, domain.JournalEntryPosted, got.Status)
	assert.NotNil(t, got.PostedAt)
}

func TestCancelEntry(t *testing.T) {
	svc, ctx := setupService(t)
	require.NoError(t, svc.CreateAccount(ctx, &domain.Account{Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset, Status: domain.AccountStatusActive, IsActive: true}))
	require.NoError(t, svc.CreateAccount(ctx, &domain.Account{Code: "5111", Name: "Expense", Type: domain.AccountTypeExpense, Status: domain.AccountStatusActive, IsActive: true}))

	je := &domain.JournalEntry{
		EntryDate:   time.Now(),
		Description: "Cancel test",
		Lines: []domain.JournalLine{
			{LineNumber: 1, AccountCode: "1111", DebitAmount: 100, CreditAmount: 0},
			{LineNumber: 2, AccountCode: "5111", DebitAmount: 0, CreditAmount: 100},
		},
	}
	require.NoError(t, svc.CreateEntry(ctx, je, "user1"))

	assert.ErrorIs(t, svc.CancelEntry(ctx, je.ID), domain.ErrJournalAlreadyDraft)
	require.NoError(t, svc.SubmitForReview(ctx, je.ID, "user1"))
	assert.NoError(t, svc.CancelEntry(ctx, je.ID))

	got, _ := svc.GetEntryByID(ctx, je.ID)
	assert.Equal(t, domain.JournalEntryCancelled, got.Status)
}

func TestGetEntriesByStatus(t *testing.T) {
	svc, ctx := setupService(t)
	require.NoError(t, svc.CreateAccount(ctx, &domain.Account{Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset, Status: domain.AccountStatusActive, IsActive: true}))
	require.NoError(t, svc.CreateAccount(ctx, &domain.Account{Code: "5111", Name: "Expense", Type: domain.AccountTypeExpense, Status: domain.AccountStatusActive, IsActive: true}))

	for i := 0; i < 3; i++ {
		je := &domain.JournalEntry{
			EntryDate:   time.Now(),
			Description: "test",
			Lines: []domain.JournalLine{
				{LineNumber: 1, AccountCode: "1111", DebitAmount: 100, CreditAmount: 0},
				{LineNumber: 2, AccountCode: "5111", DebitAmount: 0, CreditAmount: 100},
			},
		}
		require.NoError(t, svc.CreateEntry(ctx, je, "user1"))
	}

	entries, err := svc.GetEntriesByStatus(ctx, domain.JournalEntryDraft)
	require.NoError(t, err)
	assert.Len(t, entries, 3)
}

// ─── Periods ─────────────────────────────────────────────────────────

func TestCreatePeriod(t *testing.T) {
	svc, ctx := setupService(t)
	p := &domain.Period{
		Year: 2025, Month: 6,
		StartDate: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
		Status:    domain.PeriodOpen,
	}
	assert.NoError(t, svc.CreatePeriod(ctx, p))
}

func TestCloseReopenPeriod(t *testing.T) {
	svc, ctx := setupService(t)
	require.NoError(t, svc.CreatePeriod(ctx, &domain.Period{
		ID: "P-2025-06", Year: 2025, Month: 6,
		StartDate: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
		Status:    domain.PeriodOpen,
	}))

	assert.NoError(t, svc.ClosePeriod(ctx, "P-2025-06"))
	got, _ := svc.GetPeriod(ctx, "P-2025-06")
	assert.Equal(t, domain.PeriodClosed, got.Status)

	assert.ErrorIs(t, svc.ClosePeriod(ctx, "P-2025-06"), domain.ErrPeriodAlreadyClosed)

	assert.NoError(t, svc.ReopenPeriod(ctx, "P-2025-06"))
	got, _ = svc.GetPeriod(ctx, "P-2025-06")
	assert.Equal(t, domain.PeriodOpen, got.Status)
}

// ─── Users / Auth ────────────────────────────────────────────────────

func TestCreateUser_Success(t *testing.T) {
	svc, ctx := setupService(t)
	u := &domain.User{Username: "admin", FullName: "Admin", Role: domain.UserRoleAdmin, IsActive: true}
	assert.NoError(t, svc.CreateUser(ctx, u, "ValidPass123!"))
	assert.NotEmpty(t, u.PasswordHash)
	assert.NotEmpty(t, u.ID)
}

func TestCreateUser_DuplicateUsername(t *testing.T) {
	svc, ctx := setupService(t)
	u := &domain.User{Username: "admin", FullName: "Admin", Role: domain.UserRoleAdmin, IsActive: true}
	require.NoError(t, svc.CreateUser(ctx, u, "ValidPass123!"))
	assert.ErrorIs(t, svc.CreateUser(ctx, u, "ValidPass123!"), domain.ErrUsernameExists)
}

func TestCreateUser_ShortPassword(t *testing.T) {
	svc, ctx := setupService(t)
	assert.Error(t, svc.CreateUser(ctx, &domain.User{Username: "u", FullName: "U", Role: domain.UserRoleAdmin}, ""))
}

func TestLogin_Success(t *testing.T) {
	svc, ctx := setupService(t)
	auth.SetJWTSecret("test-secret-login")
	u := &domain.User{Username: "admin", FullName: "Admin", Role: domain.UserRoleAdmin, IsActive: true}
	require.NoError(t, svc.CreateUser(ctx, u, "ValidPass123!"))

	result, err := svc.Login(ctx, "admin", "ValidPass123!", "127.0.0.1")
	require.NoError(t, err)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
	assert.Equal(t, "Bearer", result.TokenType)
	assert.Equal(t, 900, result.ExpiresIn)
}

func TestLogin_WrongPassword(t *testing.T) {
	svc, ctx := setupService(t)
	u := &domain.User{Username: "admin", FullName: "Admin", Role: domain.UserRoleAdmin, IsActive: true}
	require.NoError(t, svc.CreateUser(ctx, u, "ValidPass123!"))

	_, err := svc.Login(ctx, "admin", "wrong", "127.0.0.1")
	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}

func TestLogin_AccountLocked(t *testing.T) {
	svc, ctx := setupService(t)
	original := domain.MaxLoginAttempts
	domain.MaxLoginAttempts = 3
	defer func() { domain.MaxLoginAttempts = original }()

	u := &domain.User{Username: "admin", FullName: "Admin", Role: domain.UserRoleAdmin, IsActive: true}
	require.NoError(t, svc.CreateUser(ctx, u, "ValidPass123!"))

	for i := 0; i < 3; i++ {
		svc.Login(ctx, "admin", "wrong", "127.0.0.1")
	}
	// after 3 failures, user is locked
	_, err := svc.Login(ctx, "admin", "wrong", "127.0.0.1")
	assert.ErrorIs(t, err, domain.ErrAccountLocked)

	_, err = svc.Login(ctx, "admin", "ValidPass123!", "127.0.0.1")
	assert.ErrorIs(t, err, domain.ErrAccountLocked)
}

func TestLogin_Requires2FA(t *testing.T) {
	svc, ctx := setupService(t)
	auth.SetJWTSecret("test-secret-for-jwt")

	u := &domain.User{Username: "admin", FullName: "Admin", Role: domain.UserRoleAdmin, IsActive: true, TOTPEnabled: true, TOTPSecret: "test"}
	require.NoError(t, svc.CreateUser(ctx, u, "ValidPass123!"))

	result, err := svc.Login(ctx, "admin", "ValidPass123!", "127.0.0.1")
	require.NoError(t, err)
	assert.True(t, result.Requires2FA)
	assert.NotEmpty(t, result.TempToken)
}

func TestRefreshToken(t *testing.T) {
	svc, ctx := setupService(t)
	auth.SetJWTSecret("test-secret-2")

	u := &domain.User{Username: "admin", FullName: "Admin", Role: domain.UserRoleAdmin, IsActive: true}
	require.NoError(t, svc.CreateUser(ctx, u, "ValidPass123!"))

	result, err := svc.Login(ctx, "admin", "ValidPass123!", "127.0.0.1")
	require.NoError(t, err)

	pair, err := svc.RefreshToken(ctx, result.RefreshToken)
	require.NoError(t, err)
	assert.NotEmpty(t, pair.AccessToken)
}

func TestChangePassword(t *testing.T) {
	svc, ctx := setupService(t)
	u := &domain.User{Username: "admin", FullName: "Admin", Role: domain.UserRoleAdmin, IsActive: true}
	require.NoError(t, svc.CreateUser(ctx, u, "ValidPass123!"))

	assert.NoError(t, svc.ChangePassword(ctx, u.ID, "ValidPass123!", "NewValidPass456!"))
	// old password should no longer work
	_, err := svc.Login(ctx, "admin", "ValidPass123!", "127.0.0.1")
	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}

func TestForgotPassword(t *testing.T) {
	svc, ctx := setupService(t)
	u := &domain.User{Username: "admin", FullName: "Admin", Email: "admin@test.com", Role: domain.UserRoleAdmin, IsActive: true}
	require.NoError(t, svc.CreateUser(ctx, u, "ValidPass123!"))

	assert.NoError(t, svc.ForgotPassword(ctx, "admin@test.com"))
}

// ─── Exchange Rates ──────────────────────────────────────────────────

func TestCreateExchangeRate(t *testing.T) {
	svc, ctx := setupService(t)
	rate := &domain.ExchangeRate{CurrencyCode: "USD", AverageRate: 23000, RateDate: time.Now()}
	assert.NoError(t, svc.CreateExchangeRate(ctx, rate))
}

func TestGetExchangeRate_NotFound(t *testing.T) {
	svc, ctx := setupService(t)
	_, err := svc.GetExchangeRate(ctx, "USD", time.Now())
	assert.ErrorIs(t, err, domain.ErrRateNotFound)
}

// ─── Reports ─────────────────────────────────────────────────────────

func TestTrialBalance(t *testing.T) {
	svc, ctx := setupService(t)
	require.NoError(t, svc.CreateAccount(ctx, &domain.Account{Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset, Status: domain.AccountStatusActive, IsActive: true}))
	require.NoError(t, svc.CreateAccount(ctx, &domain.Account{Code: "5111", Name: "Expense", Type: domain.AccountTypeExpense, Status: domain.AccountStatusActive, IsActive: true}))
	require.NoError(t, svc.CreatePeriod(ctx, &domain.Period{
		ID: "P-2025-06", Year: 2025, Month: 6,
		StartDate: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
		Status:    domain.PeriodOpen,
	}))

	je := &domain.JournalEntry{
		EntryDate:   time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
		PeriodID:    "P-2025-06",
		Description: "test",
		Lines: []domain.JournalLine{
			{LineNumber: 1, AccountCode: "1111", DebitAmount: 500, CreditAmount: 0},
			{LineNumber: 2, AccountCode: "5111", DebitAmount: 0, CreditAmount: 500},
		},
	}
	require.NoError(t, svc.CreateEntry(ctx, je, "user1"))
	require.NoError(t, svc.SubmitForReview(ctx, je.ID, "user1"))
	require.NoError(t, svc.ApproveEntry(ctx, je.ID, "approver1"))
	require.NoError(t, svc.PostEntry(ctx, je.ID))

	balances, err := svc.TrialBalance(ctx, 2025, 6)
	require.NoError(t, err)
	assert.Len(t, balances, 2)
}

// ─── COA Operations ──────────────────────────────────────────────────

func TestAccountBalance(t *testing.T) {
	svc, ctx := setupService(t)
	require.NoError(t, svc.CreateAccount(ctx, &domain.Account{Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset, Status: domain.AccountStatusActive, IsActive: true}))
	require.NoError(t, svc.CreateAccount(ctx, &domain.Account{Code: "5111", Name: "Expense", Type: domain.AccountTypeExpense, Status: domain.AccountStatusActive, IsActive: true}))
	require.NoError(t, svc.CreatePeriod(ctx, &domain.Period{
		ID: "P-2025-06", Year: 2025, Month: 6,
		StartDate: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
		Status:    domain.PeriodOpen,
	}))

	je := &domain.JournalEntry{
		EntryDate:   time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
		PeriodID:    "P-2025-06",
		Description: "test",
		Lines: []domain.JournalLine{
			{LineNumber: 1, AccountCode: "1111", DebitAmount: 1000, CreditAmount: 0},
			{LineNumber: 2, AccountCode: "5111", DebitAmount: 0, CreditAmount: 1000},
		},
	}
	require.NoError(t, svc.CreateEntry(ctx, je, "user1"))
	require.NoError(t, svc.SubmitForReview(ctx, je.ID, "user1"))
	require.NoError(t, svc.ApproveEntry(ctx, je.ID, "approver1"))
	require.NoError(t, svc.PostEntry(ctx, je.ID))

	b, err := svc.GetAccountBalance(ctx, "1111", "P-2025-06")
	require.NoError(t, err)
	assert.Equal(t, 1000.0, b.PeriodDebit)
	assert.Equal(t, 0.0, b.PeriodCredit)
}

func TestAccountUsage(t *testing.T) {
	svc, ctx := setupService(t)
	require.NoError(t, svc.CreateAccount(ctx, &domain.Account{Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset, Status: domain.AccountStatusActive, IsActive: true}))
	require.NoError(t, svc.CreateAccount(ctx, &domain.Account{Code: "5111", Name: "Expense", Type: domain.AccountTypeExpense, Status: domain.AccountStatusActive, IsActive: true}))
	require.NoError(t, svc.CreatePeriod(ctx, &domain.Period{
		ID: "P-2025-06", Year: 2025, Month: 6,
		StartDate: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
		Status:    domain.PeriodOpen,
	}))

	je := &domain.JournalEntry{
		EntryDate:   time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
		PeriodID:    "P-2025-06",
		Description: "test",
		Lines: []domain.JournalLine{
			{LineNumber: 1, AccountCode: "1111", DebitAmount: 100, CreditAmount: 0},
			{LineNumber: 2, AccountCode: "5111", DebitAmount: 0, CreditAmount: 100},
		},
	}
	require.NoError(t, svc.CreateEntry(ctx, je, "user1"))
	require.NoError(t, svc.SubmitForReview(ctx, je.ID, "user1"))
	require.NoError(t, svc.ApproveEntry(ctx, je.ID, "approver1"))
	require.NoError(t, svc.PostEntry(ctx, je.ID))

	usage, err := svc.GetAccountUsage(ctx, "1111")
	require.NoError(t, err)
	assert.Equal(t, 1, usage.EntryCount)
	assert.Equal(t, 100.0, usage.TotalDebit)
}

func TestApprovalRequestLifecycle(t *testing.T) {
	svc, ctx := setupService(t)
	req := &domain.ApprovalRequest{
		EntityType:  "account",
		EntityID:    "1111",
		RequestType: "change",
		Reason:      "restructure",
		RequestedBy: "user1",
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}
	require.NoError(t, svc.CreateApprovalRequest(ctx, req))

	assert.ErrorIs(t, svc.ApproveRequest(ctx, req.ID, "user1", "ok"), domain.ErrSelfApproval)
	assert.NoError(t, svc.ApproveRequest(ctx, req.ID, "approver1", "ok"))
	assert.ErrorIs(t, svc.ApproveRequest(ctx, req.ID, "approver2", "ok"), domain.ErrApprovalAlreadyProcessed)
}

func TestCreateAccountVersion(t *testing.T) {
	svc, ctx := setupService(t)
	require.NoError(t, svc.CreateAccount(ctx, &domain.Account{Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset}))

	ver, err := svc.CreateAccountVersion(ctx, "initial snapshot")
	require.NoError(t, err)
	assert.Equal(t, "v1", ver.VersionNumber)

	ver2, err := svc.CreateAccountVersion(ctx, "second snapshot")
	require.NoError(t, err)
	assert.Equal(t, "v2", ver2.VersionNumber)
}

func TestAccountMapping(t *testing.T) {
	svc, ctx := setupService(t)
	m := &domain.AccountMapping{
		SourceRegime: "VAS",
		TargetRegime: "IFRS",
		OldCode:      "1111",
		NewCode:      "IFRS-1111",
		MappingType:  "DIRECT",
	}
	assert.NoError(t, svc.CreateAccountMapping(ctx, m))
	assert.ErrorIs(t, svc.CreateAccountMapping(ctx, m), domain.ErrMappingExists)

	got, err := svc.GetMappingByOldCode(ctx, "VAS", "1111")
	require.NoError(t, err)
	assert.Equal(t, "IFRS-1111", got.NewCode)
}

func TestIFRSMapping(t *testing.T) {
	svc, ctx := setupService(t)
	m := &domain.IFRSMapping{VASCode: "1111", IFRSCode: "IFRS-1111"}
	require.NoError(t, svc.CreateIFRSMapping(ctx, m))

	assert.ErrorIs(t, svc.CreateIFRSMapping(ctx, m), domain.ErrIFRSMappingExists)

	got, err := svc.GetIFRSMapping(ctx, "1111")
	require.NoError(t, err)
	assert.Equal(t, "IFRS-1111", got.IFRSCode)

	list, err := svc.ListIFRSMappings(ctx)
	require.NoError(t, err)
	assert.Len(t, list, 1)
}

func TestAccountAnalysis(t *testing.T) {
	svc, ctx := setupService(t)
	require.NoError(t, svc.CreateAccount(ctx, &domain.Account{Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset}))

	a := &domain.AccountAnalysis{AccountCode: "1111", CostCenterID: "CC-01"}
	assert.NoError(t, svc.CreateAccountAnalysis(ctx, a))

	got, err := svc.GetAccountAnalysis(ctx, "1111")
	require.NoError(t, err)
	assert.Equal(t, "CC-01", got.CostCenterID)
}

func TestLogout(t *testing.T) {
	svc, ctx := setupService(t)
	auth.SetJWTSecret("test-secret-logout")

	u := &domain.User{Username: "admin", FullName: "Admin", Role: domain.UserRoleAdmin, IsActive: true}
	require.NoError(t, svc.CreateUser(ctx, u, "ValidPass123!"))

	result, err := svc.Login(ctx, "admin", "ValidPass123!", "127.0.0.1")
	require.NoError(t, err)

	assert.NoError(t, svc.Logout(ctx, u.ID, result.RefreshToken))

	_, err = svc.RefreshToken(ctx, result.RefreshToken)
	assert.ErrorIs(t, err, domain.ErrRefreshTokenRevoked)
}

func TestLogoutAll(t *testing.T) {
	svc, ctx := setupService(t)
	auth.SetJWTSecret("test-secret-logout-all")

	u := &domain.User{Username: "admin", FullName: "Admin", Role: domain.UserRoleAdmin, IsActive: true}
	require.NoError(t, svc.CreateUser(ctx, u, "ValidPass123!"))

	result, err := svc.Login(ctx, "admin", "ValidPass123!", "127.0.0.1")
	require.NoError(t, err)

	assert.NoError(t, svc.LogoutAll(ctx, u.ID))
	_, err = svc.RefreshToken(ctx, result.RefreshToken)
	assert.ErrorIs(t, err, domain.ErrRefreshTokenRevoked)
}

func TestGetEntriesByDateRange(t *testing.T) {
	svc, ctx := setupService(t)
	require.NoError(t, svc.CreateAccount(ctx, &domain.Account{Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset, Status: domain.AccountStatusActive, IsActive: true}))
	require.NoError(t, svc.CreateAccount(ctx, &domain.Account{Code: "5111", Name: "Expense", Type: domain.AccountTypeExpense, Status: domain.AccountStatusActive, IsActive: true}))

	je := &domain.JournalEntry{
		EntryDate:   time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
		Description: "test",
		Lines: []domain.JournalLine{
			{LineNumber: 1, AccountCode: "1111", DebitAmount: 100, CreditAmount: 0},
			{LineNumber: 2, AccountCode: "5111", DebitAmount: 0, CreditAmount: 100},
		},
	}
	require.NoError(t, svc.CreateEntry(ctx, je, "user1"))

	entries, err := svc.GetEntriesByDateRange(ctx, time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	assert.Len(t, entries, 1)
}

func TestDrillDown(t *testing.T) {
	svc, ctx := setupService(t)
	require.NoError(t, svc.CreateAccount(ctx, &domain.Account{Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset, Status: domain.AccountStatusActive, IsActive: true}))
	require.NoError(t, svc.CreateAccount(ctx, &domain.Account{Code: "5111", Name: "Expense", Type: domain.AccountTypeExpense, Status: domain.AccountStatusActive, IsActive: true}))
	require.NoError(t, svc.CreatePeriod(ctx, &domain.Period{
		ID: "P-2025-06", Year: 2025, Month: 6,
		StartDate: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
		Status:    domain.PeriodOpen,
	}))

	je := &domain.JournalEntry{
		EntryDate:   time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
		PeriodID:    "P-2025-06",
		Description: "test",
		Lines: []domain.JournalLine{
			{LineNumber: 1, AccountCode: "1111", DebitAmount: 100, CreditAmount: 0},
			{LineNumber: 2, AccountCode: "5111", DebitAmount: 0, CreditAmount: 100},
		},
	}
	require.NoError(t, svc.CreateEntry(ctx, je, "user1"))
	require.NoError(t, svc.SubmitForReview(ctx, je.ID, "user1"))
	require.NoError(t, svc.ApproveEntry(ctx, je.ID, "approver1"))
	require.NoError(t, svc.PostEntry(ctx, je.ID))

	entries, err := svc.GetAccountBalanceDrillDown(ctx, "1111", "P-2025-06")
	require.NoError(t, err)
	assert.Len(t, entries, 1)
}

func TestSessions(t *testing.T) {
	svc, ctx := setupService(t)
	auth.SetJWTSecret("test-secret-sessions")

	u := &domain.User{Username: "admin", FullName: "Admin", Role: domain.UserRoleAdmin, IsActive: true}
	require.NoError(t, svc.CreateUser(ctx, u, "ValidPass123!"))

	svc.Login(ctx, "admin", "ValidPass123!", "127.0.0.1")

	sessions, err := svc.ListSessions(ctx, u.ID)
	require.NoError(t, err)
	assert.Len(t, sessions, 1)
	assert.True(t, sessions[0].IsCurrent)
}

func TestGetPeriodByYearMonth(t *testing.T) {
	svc, ctx := setupService(t)
	require.NoError(t, svc.CreatePeriod(ctx, &domain.Period{
		ID: "P-2025-06", Year: 2025, Month: 6,
		StartDate: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
		Status:    domain.PeriodOpen,
	}))

	p, err := svc.GetPeriodByYearMonth(ctx, 2025, 6)
	require.NoError(t, err)
	assert.Equal(t, "P-2025-06", p.ID)
}

// ─── Cash Module Tests ──────────────────────────────────────────────

func TestCreateCashReceipt_Success(t *testing.T) {
	svc, ctx := setupService(t)
	createTestAccount(t, svc, ctx, "1111", domain.AccountTypeAsset)
	createTestAccount(t, svc, ctx, "5111", domain.AccountTypeRevenue)

	r := &domain.CashReceipt{
		CompanyID: "CMP001", CashAccountID: "1111",
		DebitAccountID: "1111", CreditAccountID: "5111",
		Amount: 1000000, AmountVND: 1000000,
		Currency: "VND", ExchangeRate: 1,
		VoucherDate: "2026-07-01", Reason: "Cash sales",
		ReceiptType: domain.ReceiptSales, CounterpartType: domain.CounterpartCustomer,
		CounterpartName: "Nguyen Van A",
		Status: domain.CashDraft, CreatedBy: "user1",
	}
	err := svc.CreateCashReceipt(ctx, r)
	require.NoError(t, err)
	assert.NotEmpty(t, r.ID)
	assert.NotEmpty(t, r.VoucherNo)
	assert.Equal(t, domain.CashDraft, r.Status)
}

func TestCreateCashReceipt_ValidationError(t *testing.T) {
	svc, ctx := setupService(t)
	r := &domain.CashReceipt{Amount: 0}
	err := svc.CreateCashReceipt(ctx, r)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "positive")
}

func TestSubmitCashReceipt(t *testing.T) {
	svc, ctx := setupService(t)
	createTestAccount(t, svc, ctx, "1111", domain.AccountTypeAsset)
	createTestAccount(t, svc, ctx, "5111", domain.AccountTypeRevenue)

	r := &domain.CashReceipt{
		CompanyID: "CMP001", CashAccountID: "1111",
		DebitAccountID: "1111", CreditAccountID: "5111",
		Amount: 1000000, AmountVND: 1000000,
		Currency: "VND", ExchangeRate: 1,
		VoucherDate: "2026-07-01", Reason: "Test",
		ReceiptType: domain.ReceiptOther, CounterpartType: domain.CounterpartOther,
		Status: domain.CashDraft, CreatedBy: "user1",
	}
	require.NoError(t, svc.CreateCashReceipt(ctx, r))

	err := svc.SubmitCashReceipt(ctx, r.ID, "user1")
	require.NoError(t, err)
	got, _ := svc.GetCashReceipt(ctx, r.ID)
	assert.Equal(t, domain.CashSubmitted, got.Status)
}

func TestSubmitCashReceipt_AlreadyPosted(t *testing.T) {
	svc, ctx := setupService(t)
	createTestAccount(t, svc, ctx, "1111", domain.AccountTypeAsset)
	createTestAccount(t, svc, ctx, "5111", domain.AccountTypeRevenue)

	r := &domain.CashReceipt{
		CompanyID: "CMP001", CashAccountID: "1111",
		DebitAccountID: "1111", CreditAccountID: "5111",
		Amount: 1000000, AmountVND: 1000000,
		Currency: "VND", ExchangeRate: 1,
		VoucherDate: "2026-07-01", Reason: "Test",
		ReceiptType: domain.ReceiptOther, CounterpartType: domain.CounterpartOther,
		Status: domain.CashDraft, CreatedBy: "user1",
	}
	require.NoError(t, svc.CreateCashReceipt(ctx, r))
	require.NoError(t, svc.SubmitCashReceipt(ctx, r.ID, "user1"))
	require.NoError(t, svc.ApproveCashReceipt(ctx, r.ID, "user2"))
	require.NoError(t, svc.PostCashReceipt(ctx, r.ID, "user2"))

	err := svc.SubmitCashReceipt(ctx, r.ID, "user1")
	require.Error(t, err)
}

func TestApproveCashReceipt(t *testing.T) {
	svc, ctx := setupService(t)
	createTestAccount(t, svc, ctx, "1111", domain.AccountTypeAsset)
	createTestAccount(t, svc, ctx, "5111", domain.AccountTypeRevenue)

	r := &domain.CashReceipt{
		CompanyID: "CMP001", CashAccountID: "1111",
		DebitAccountID: "1111", CreditAccountID: "5111",
		Amount: 1000000, AmountVND: 1000000,
		Currency: "VND", ExchangeRate: 1,
		VoucherDate: "2026-07-01", Reason: "Test",
		ReceiptType: domain.ReceiptOther, CounterpartType: domain.CounterpartOther,
		Status: domain.CashDraft, CreatedBy: "user1",
	}
	require.NoError(t, svc.CreateCashReceipt(ctx, r))
	require.NoError(t, svc.SubmitCashReceipt(ctx, r.ID, "user1"))

	err := svc.ApproveCashReceipt(ctx, r.ID, "user2")
	require.NoError(t, err)
	got, _ := svc.GetCashReceipt(ctx, r.ID)
	assert.Equal(t, domain.CashApproved, got.Status)
}

func TestApproveCashReceipt_SelfApproval(t *testing.T) {
	svc, ctx := setupService(t)
	createTestAccount(t, svc, ctx, "1111", domain.AccountTypeAsset)
	createTestAccount(t, svc, ctx, "5111", domain.AccountTypeRevenue)

	r := &domain.CashReceipt{
		CompanyID: "CMP001", CashAccountID: "1111",
		DebitAccountID: "1111", CreditAccountID: "5111",
		Amount: 1000000, AmountVND: 1000000,
		Currency: "VND", ExchangeRate: 1,
		VoucherDate: "2026-07-01", Reason: "Test",
		ReceiptType: domain.ReceiptOther, CounterpartType: domain.CounterpartOther,
		Status: domain.CashDraft, CreatedBy: "user1",
	}
	require.NoError(t, svc.CreateCashReceipt(ctx, r))
	require.NoError(t, svc.SubmitCashReceipt(ctx, r.ID, "user1"))

	err := svc.ApproveCashReceipt(ctx, r.ID, "user1")
	require.ErrorIs(t, err, domain.ErrSelfCashApproval)
}

func TestRejectCashReceipt(t *testing.T) {
	svc, ctx := setupService(t)
	createTestAccount(t, svc, ctx, "1111", domain.AccountTypeAsset)
	createTestAccount(t, svc, ctx, "5111", domain.AccountTypeRevenue)

	r := &domain.CashReceipt{
		CompanyID: "CMP001", CashAccountID: "1111",
		DebitAccountID: "1111", CreditAccountID: "5111",
		Amount: 1000000, AmountVND: 1000000,
		Currency: "VND", ExchangeRate: 1,
		VoucherDate: "2026-07-01", Reason: "Test",
		ReceiptType: domain.ReceiptOther, CounterpartType: domain.CounterpartOther,
		Status: domain.CashDraft, CreatedBy: "user1",
	}
	require.NoError(t, svc.CreateCashReceipt(ctx, r))
	require.NoError(t, svc.SubmitCashReceipt(ctx, r.ID, "user1"))

	err := svc.RejectCashReceipt(ctx, r.ID, "user2")
	require.NoError(t, err)
	got, _ := svc.GetCashReceipt(ctx, r.ID)
	assert.Equal(t, domain.CashRejected, got.Status)

	err = svc.SubmitCashReceipt(ctx, r.ID, "user1")
	require.NoError(t, err)
	got, _ = svc.GetCashReceipt(ctx, r.ID)
	assert.Equal(t, domain.CashSubmitted, got.Status)
}

func TestCashFullLifecycle(t *testing.T) {
	svc, ctx := setupService(t)
	createTestAccount(t, svc, ctx, "1111", domain.AccountTypeAsset)
	createTestAccount(t, svc, ctx, "5111", domain.AccountTypeRevenue)

	r := &domain.CashReceipt{
		CompanyID: "CMP001", CashAccountID: "1111",
		DebitAccountID: "1111", CreditAccountID: "5111",
		Amount: 2000000, AmountVND: 2000000,
		Currency: "VND", ExchangeRate: 1,
		VoucherDate: "2026-07-01", Reason: "Full life cycle",
		ReceiptType: domain.ReceiptSales, CounterpartType: domain.CounterpartCustomer,
		Status: domain.CashDraft, CreatedBy: "user1",
	}
	require.NoError(t, svc.CreateCashReceipt(ctx, r))

	require.NoError(t, svc.SubmitCashReceipt(ctx, r.ID, "user1"))
	require.NoError(t, svc.ApproveCashReceipt(ctx, r.ID, "user2"))
	require.NoError(t, svc.PostCashReceipt(ctx, r.ID, "user2"))

	got, err := svc.GetCashReceipt(ctx, r.ID)
	require.NoError(t, err)
	assert.Equal(t, domain.CashPosted, got.Status)
	assert.NotEmpty(t, got.GLJournalID)

	bal, err := svc.GetCashBalance(ctx, "CMP001", "1111")
	require.NoError(t, err)
	assert.Equal(t, 2000000.0, bal)
}

func TestCreateCashPayment_Success(t *testing.T) {
	svc, ctx := setupService(t)
	createTestAccount(t, svc, ctx, "1111", domain.AccountTypeAsset)
	createTestAccount(t, svc, ctx, "3311", domain.AccountTypeLiability)

	p := &domain.CashPayment{
		CompanyID: "CMP001", CashAccountID: "1111",
		DebitAccountID: "3311", CreditAccountID: "1111",
		Amount: 500000, AmountVND: 500000,
		Currency: "VND", ExchangeRate: 1,
		VoucherDate: "2026-07-01", Reason: "Supplier payment",
		PaymentType: domain.PaymentSupplier, PayeeType: domain.CounterpartSupplier,
		Status: domain.CashDraft, CreatedBy: "user1",
	}
	err := svc.CreateCashPayment(ctx, p)
	require.NoError(t, err)
	assert.NotEmpty(t, p.ID)
	assert.Equal(t, domain.CashDraft, p.Status)
}

func TestPostCashPayment_UpdatesBalance(t *testing.T) {
	svc, ctx := setupService(t)
	createTestAccount(t, svc, ctx, "1111", domain.AccountTypeAsset)
	createTestAccount(t, svc, ctx, "3311", domain.AccountTypeLiability)

	// first add cash
	r := &domain.CashReceipt{
		CompanyID: "CMP001", CashAccountID: "1111",
		DebitAccountID: "1111", CreditAccountID: "5111",
		Amount: 2000000, AmountVND: 2000000,
		Currency: "VND", ExchangeRate: 1,
		VoucherDate: "2026-07-01", Reason: "Add cash",
		ReceiptType: domain.ReceiptOther, CounterpartType: domain.CounterpartOther,
		Status: domain.CashDraft, CreatedBy: "user1",
	}
	require.NoError(t, svc.CreateCashReceipt(ctx, r))
	require.NoError(t, svc.SubmitCashReceipt(ctx, r.ID, "user1"))
	require.NoError(t, svc.ApproveCashReceipt(ctx, r.ID, "user2"))
	require.NoError(t, svc.PostCashReceipt(ctx, r.ID, "user2"))

	p := &domain.CashPayment{
		CompanyID: "CMP001", CashAccountID: "1111",
		DebitAccountID: "3311", CreditAccountID: "1111",
		Amount: 500000, AmountVND: 500000,
		Currency: "VND", ExchangeRate: 1,
		VoucherDate: "2026-07-01", Reason: "Pay supplier",
		PaymentType: domain.PaymentSupplier, PayeeType: domain.CounterpartSupplier,
		Status: domain.CashDraft, CreatedBy: "user1",
	}
	require.NoError(t, svc.CreateCashPayment(ctx, p))
	require.NoError(t, svc.SubmitCashPayment(ctx, p.ID, "user1"))
	require.NoError(t, svc.ApproveCashPayment(ctx, p.ID, "user2"))
	require.NoError(t, svc.PostCashPayment(ctx, p.ID, "user2"))

	got, _ := svc.GetCashPayment(ctx, p.ID)
	assert.Equal(t, domain.CashPosted, got.Status)

	bal, _ := svc.GetCashBalance(ctx, "CMP001", "1111")
	assert.Equal(t, 1500000.0, bal)
}

func TestListCashReceipts(t *testing.T) {
	svc, ctx := setupService(t)
	createTestAccount(t, svc, ctx, "1111", domain.AccountTypeAsset)
	createTestAccount(t, svc, ctx, "5111", domain.AccountTypeRevenue)

	r1 := &domain.CashReceipt{
		CompanyID: "CMP001", CashAccountID: "1111",
		DebitAccountID: "1111", CreditAccountID: "5111",
		Amount: 100000, AmountVND: 100000,
		Currency: "VND", ExchangeRate: 1,
		VoucherDate: "2026-07-01", Reason: "R1",
		ReceiptType: domain.ReceiptSales, CounterpartType: domain.CounterpartCustomer,
		Status: domain.CashDraft, CreatedBy: "user1",
	}
	r2 := &domain.CashReceipt{
		CompanyID: "CMP001", CashAccountID: "1111",
		DebitAccountID: "1111", CreditAccountID: "5111",
		Amount: 200000, AmountVND: 200000,
		Currency: "VND", ExchangeRate: 1,
		VoucherDate: "2026-07-02", Reason: "R2",
		ReceiptType: domain.ReceiptOther, CounterpartType: domain.CounterpartOther,
		Status: domain.CashDraft, CreatedBy: "user1",
	}
	require.NoError(t, svc.CreateCashReceipt(ctx, r1))
	require.NoError(t, svc.CreateCashReceipt(ctx, r2))

	list, total, err := svc.ListCashReceipts(ctx, domain.CashReceiptFilter{CompanyID: "CMP001", Limit: 10})
	require.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, list, 2)
}

func TestCreateCashTransfer(t *testing.T) {
	svc, ctx := setupService(t)
	createTestAccount(t, svc, ctx, "1111", domain.AccountTypeAsset)
	createTestAccount(t, svc, ctx, "1121", domain.AccountTypeAsset)

	tf := &domain.CashTransfer{
		CompanyID: "CMP001", TransferDate: "2026-07-01",
		FromAccountID: "1121", ToAccountID: "1111",
		Amount: 10000000, Currency: "VND", ExchangeRate: 1,
		Reason: "Rut tien NH ve quy", TransferType: domain.TransferBankWithdrawal,
	}
	err := svc.CreateCashTransfer(ctx, tf)
	require.NoError(t, err)
	assert.NotEmpty(t, tf.ID)
}

func TestCashBookReport(t *testing.T) {
	svc, ctx := setupService(t)
	createTestAccount(t, svc, ctx, "1111", domain.AccountTypeAsset)
	createTestAccount(t, svc, ctx, "5111", domain.AccountTypeRevenue)

	r := &domain.CashReceipt{
		CompanyID: "CMP001", CashAccountID: "1111",
		DebitAccountID: "1111", CreditAccountID: "5111",
		Amount: 5000000, AmountVND: 5000000,
		Currency: "VND", ExchangeRate: 1,
		VoucherDate: "2026-07-01", Reason: "Sales",
		ReceiptType: domain.ReceiptSales, CounterpartType: domain.CounterpartCustomer,
		Status: domain.CashDraft, CreatedBy: "user1",
	}
	require.NoError(t, svc.CreateCashReceipt(ctx, r))
	require.NoError(t, svc.SubmitCashReceipt(ctx, r.ID, "user1"))
	require.NoError(t, svc.ApproveCashReceipt(ctx, r.ID, "user2"))
	require.NoError(t, svc.PostCashReceipt(ctx, r.ID, "user2"))

	cb, err := svc.GetCashBook(ctx, "CMP001", "VND", "1111", "2026-07-01", "2026-07-31")
	require.NoError(t, err)
	assert.Equal(t, 5000000.0, cb.TotalReceipts)
	assert.Equal(t, 5000000.0, cb.ClosingBalance)
	assert.Equal(t, 0.0, cb.OpeningBalance)
	assert.Len(t, cb.Entries, 1)
}

func TestCreatePettyCashFund(t *testing.T) {
	svc, ctx := setupService(t)
	f := &domain.PettyCashFund{
		CompanyID: "CMP001", FundCode: "PC001", FundName: "Quy tam ung van phong",
		CustodianID: "EMP001", InitialAmount: 5000000, CurrentBalance: 5000000,
		Currency: "VND", Status: domain.PettyCashActive,
	}
	err := svc.CreatePettyCashFund(ctx, f)
	require.NoError(t, err)
	assert.NotEmpty(t, f.ID)
}

func TestCreateCashInventory(t *testing.T) {
	svc, ctx := setupService(t)
	createTestAccount(t, svc, ctx, "1111", domain.AccountTypeAsset)
	createTestAccount(t, svc, ctx, "5111", domain.AccountTypeRevenue)

	r := &domain.CashReceipt{
		CompanyID: "CMP001", CashAccountID: "1111",
		DebitAccountID: "1111", CreditAccountID: "5111",
		Amount: 10000000, AmountVND: 10000000,
		Currency: "VND", ExchangeRate: 1,
		VoucherDate: "2026-07-01", Reason: "Sales",
		ReceiptType: domain.ReceiptSales, CounterpartType: domain.CounterpartCustomer,
		Status: domain.CashDraft, CreatedBy: "user1",
	}
	require.NoError(t, svc.CreateCashReceipt(ctx, r))
	require.NoError(t, svc.SubmitCashReceipt(ctx, r.ID, "user1"))
	require.NoError(t, svc.ApproveCashReceipt(ctx, r.ID, "user2"))
	require.NoError(t, svc.PostCashReceipt(ctx, r.ID, "user2"))

	inv := &domain.CashInventory{
		CompanyID: "CMP001", InventoryDate: "2026-07-15",
		CashAccountID: "1111", Currency: "VND",
		BookBalance: 10000000, ActualBalance: 10000000,
		Difference: 0, DifferenceType: "none",
		Denominations: []domain.DenominationDetail{
			{Denomination: 500000, Count: 20, Subtotal: 10000000},
		},
		Status: domain.CashInventoryDraft,
	}
	err := svc.CreateCashInventory(ctx, inv)
	require.NoError(t, err)
	assert.NotEmpty(t, inv.ID)
}

func TestCashInventoryDiscrepancy(t *testing.T) {
	svc, ctx := setupService(t)
	createTestAccount(t, svc, ctx, "1111", domain.AccountTypeAsset)
	createTestAccount(t, svc, ctx, "5111", domain.AccountTypeRevenue)

	r := &domain.CashReceipt{
		CompanyID: "CMP001", CashAccountID: "1111",
		DebitAccountID: "1111", CreditAccountID: "5111",
		Amount: 10000000, AmountVND: 10000000,
		Currency: "VND", ExchangeRate: 1,
		VoucherDate: "2026-07-01", Reason: "Sales",
		ReceiptType: domain.ReceiptSales, CounterpartType: domain.CounterpartCustomer,
		Status: domain.CashDraft, CreatedBy: "user1",
	}
	require.NoError(t, svc.CreateCashReceipt(ctx, r))
	require.NoError(t, svc.SubmitCashReceipt(ctx, r.ID, "user1"))
	require.NoError(t, svc.ApproveCashReceipt(ctx, r.ID, "user2"))
	require.NoError(t, svc.PostCashReceipt(ctx, r.ID, "user2"))

	inv := &domain.CashInventory{
		CompanyID: "CMP001", InventoryDate: "2026-07-15",
		CashAccountID: "1111", Currency: "VND",
		BookBalance: 10000000, ActualBalance: 9500000,
		Difference: -500000, DifferenceType: "shortage",
		Status: domain.CashInventoryDraft,
	}
	err := svc.CreateCashInventory(ctx, inv)
	require.NoError(t, err)
	assert.Equal(t, "shortage", inv.DifferenceType)
}

func createTestAccount(t *testing.T, svc Service, ctx context.Context, code string, atype domain.AccountType) {
	t.Helper()
	require.NoError(t, svc.CreateAccount(ctx, &domain.Account{
		Code: code, Name: "Test "+code, Type: atype, IsActive: true,
	}))
}

// ─── Advance Request / Settlement Tests ───────────────────────────────

func TestCreateAdvance_Success(t *testing.T) {
	svc, ctx := setupService(t)
	a := &domain.AdvanceRequest{
		CompanyID: "CMP001", RequestorID: "user1",
		Amount: 5000000, AmountVND: 5000000,
		Currency: "VND", ExchangeRate: 1,
		Purpose: "Travel advance",
	}
	err := svc.CreateAdvance(ctx, a)
	require.NoError(t, err)
	assert.NotEmpty(t, a.ID)
	assert.Equal(t, domain.AdvanceDraft, a.Status)
}

func TestCreateAdvance_ValidationError(t *testing.T) {
	svc, ctx := setupService(t)
	a := &domain.AdvanceRequest{CompanyID: "CMP001", Amount: 0, Purpose: ""}
	err := svc.CreateAdvance(ctx, a)
	assert.Error(t, err)
}

func TestAdvanceFullLifecycle(t *testing.T) {
	svc, ctx := setupService(t)
	a := &domain.AdvanceRequest{
		CompanyID: "CMP001", RequestorID: "user1",
		Amount: 5000000, AmountVND: 5000000,
		Currency: "VND", ExchangeRate: 1,
		Purpose: "Travel advance",
	}
	require.NoError(t, svc.CreateAdvance(ctx, a))

	require.NoError(t, svc.ApproveAdvance(ctx, a.ID, "manager1"))
	got, err := svc.GetAdvance(ctx, a.ID)
	require.NoError(t, err)
	assert.Equal(t, domain.AdvanceApproved, got.Status)

	require.NoError(t, svc.PayAdvance(ctx, a.ID, "cashier1"))
	got, err = svc.GetAdvance(ctx, a.ID)
	require.NoError(t, err)
	assert.Equal(t, domain.AdvancePaid, got.Status)
	assert.Equal(t, "cashier1", got.PaidBy)

	require.NoError(t, svc.SettleAdvance(ctx, a.ID, "settle-1"))
	got, err = svc.GetAdvance(ctx, a.ID)
	require.NoError(t, err)
	assert.Equal(t, domain.AdvanceSettled, got.Status)
}

func TestAdvanceRejectThenUpdate(t *testing.T) {
	svc, ctx := setupService(t)
	a := &domain.AdvanceRequest{
		CompanyID: "CMP001", RequestorID: "user1",
		Amount: 1000000, AmountVND: 1000000,
		Currency: "VND", ExchangeRate: 1,
		Purpose: "advance",
	}
	require.NoError(t, svc.CreateAdvance(ctx, a))

	require.NoError(t, svc.RejectAdvance(ctx, a.ID, "manager1"))
	got, err := svc.GetAdvance(ctx, a.ID)
	require.NoError(t, err)
	assert.Equal(t, domain.AdvanceRejected, got.Status)

	got.Amount = 2000000
	require.NoError(t, svc.UpdateAdvance(ctx, got))
	got2, err := svc.GetAdvance(ctx, a.ID)
	require.NoError(t, err)
	assert.Equal(t, 2000000.0, got2.Amount)
}

func TestAdvanceInvalidTransition(t *testing.T) {
	svc, ctx := setupService(t)
	a := &domain.AdvanceRequest{
		CompanyID: "CMP001", RequestorID: "user1",
		Amount: 1000000, AmountVND: 1000000,
		Currency: "VND", ExchangeRate: 1,
		Purpose: "advance",
	}
	require.NoError(t, svc.CreateAdvance(ctx, a))

	// can't pay without approving
	err := svc.PayAdvance(ctx, a.ID, "cashier1")
	assert.Error(t, err)
}

func TestListAdvances(t *testing.T) {
	svc, ctx := setupService(t)
	for i := 0; i < 3; i++ {
		a := &domain.AdvanceRequest{
			CompanyID: "CMP001", RequestorID: "user1",
			Amount: 1000000, AmountVND: 1000000,
			Currency: "VND", ExchangeRate: 1,
			Purpose: "advance test",
		}
		require.NoError(t, svc.CreateAdvance(ctx, a))
	}
	advances, err := svc.ListAdvances(ctx, "CMP001")
	require.NoError(t, err)
	assert.Len(t, advances, 3)
}

func TestGetAllPeriods(t *testing.T) {
	svc, ctx := setupService(t)
	require.NoError(t, svc.CreatePeriod(ctx, &domain.Period{
		ID: "P-2025-06", Year: 2025, Month: 6,
		StartDate: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
		Status:    domain.PeriodOpen,
	}))
	require.NoError(t, svc.CreatePeriod(ctx, &domain.Period{
		ID: "P-2025-07", Year: 2025, Month: 7,
		StartDate: time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 7, 31, 0, 0, 0, 0, time.UTC),
		Status:    domain.PeriodOpen,
	}))

	periods, err := svc.GetAllPeriods(ctx)
	require.NoError(t, err)
	assert.Len(t, periods, 2)
}
