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

	svc := NewService(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo,
		approvalRepo, versionRepo, mappingRepo, analysisRepo, ifrsRepo, refreshRepo, resetRepo)
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
