package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/domain"
	"gotax/internal/repository"
)

func setupYearEndService(t *testing.T) (Service, context.Context) {
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

func TestYearEndClose(t *testing.T) {
	svc, ctx := setupYearEndService(t)

	// Create accounts
	revenueAcct := &domain.Account{Code: "5111", Name: "Revenue", Type: domain.AccountTypeRevenue}
	expenseAcct := &domain.Account{Code: "6321", Name: "COGS", Type: domain.AccountTypeExpense}
	assetAcct := &domain.Account{Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset}
	retainedAcct := &domain.Account{Code: "421", Name: "Retained Earnings", Type: domain.AccountTypeEquity}

	require.NoError(t, svc.CreateAccount(ctx, revenueAcct))
	require.NoError(t, svc.CreateAccount(ctx, expenseAcct))
	require.NoError(t, svc.CreateAccount(ctx, assetAcct))
	require.NoError(t, svc.CreateAccount(ctx, retainedAcct))

	// Create opening balances
	balances := []domain.OpeningBalance{
		{CompanyID: "C001", PeriodID: "P-2024-12", AccountCode: "5111", CurrencyCode: "VND", DebitAmount: 0, CreditAmount: 100_000_000, Status: domain.OBStatusApproved, CreatedBy: "U001"},
		{CompanyID: "C001", PeriodID: "P-2024-12", AccountCode: "6321", CurrencyCode: "VND", DebitAmount: 60_000_000, CreditAmount: 0, Status: domain.OBStatusApproved, CreatedBy: "U001"},
		{CompanyID: "C001", PeriodID: "P-2024-12", AccountCode: "1111", CurrencyCode: "VND", DebitAmount: 40_000_000, CreditAmount: 0, Status: domain.OBStatusApproved, CreatedBy: "U001"},
	}
	require.NoError(t, svc.BulkCreateOpeningBalances(ctx, balances))

	// Run year-end close
	result, err := svc.YearEndClose(ctx, "C001", "P-2024-12", "P-2025-01", "2024", "2025", "U001")
	require.NoError(t, err)

	// Verify results
	assert.Equal(t, 1, result.ClosedRevenueCount)
	assert.Equal(t, 1, result.ClosedExpenseCount)
	assert.Equal(t, 100_000_000.0, result.TotalRevenue)
	assert.Equal(t, 60_000_000.0, result.TotalExpense)
	assert.Equal(t, 40_000_000.0, result.NetIncome)
	assert.Equal(t, 1, result.CarryForwardCount) // Only asset account carried forward
	assert.NotEmpty(t, result.ClosingEntryID)
	assert.NotNil(t, result.CarryForwardLog)
	assert.Equal(t, "COMPLETED", result.CarryForwardLog.Status)
}

func TestYearEndClose_NegativeIncome(t *testing.T) {
	svc, ctx := setupYearEndService(t)

	// Create accounts
	revenueAcct := &domain.Account{Code: "5111", Name: "Revenue", Type: domain.AccountTypeRevenue}
	expenseAcct := &domain.Account{Code: "6321", Name: "COGS", Type: domain.AccountTypeExpense}
	assetAcct := &domain.Account{Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset}

	require.NoError(t, svc.CreateAccount(ctx, revenueAcct))
	require.NoError(t, svc.CreateAccount(ctx, expenseAcct))
	require.NoError(t, svc.CreateAccount(ctx, assetAcct))

	// Create opening balances (loss scenario)
	balances := []domain.OpeningBalance{
		{CompanyID: "C001", PeriodID: "P-2024-12", AccountCode: "5111", CurrencyCode: "VND", DebitAmount: 0, CreditAmount: 30_000_000, Status: domain.OBStatusApproved, CreatedBy: "U001"},
		{CompanyID: "C001", PeriodID: "P-2024-12", AccountCode: "6321", CurrencyCode: "VND", DebitAmount: 50_000_000, CreditAmount: 0, Status: domain.OBStatusApproved, CreatedBy: "U001"},
		{CompanyID: "C001", PeriodID: "P-2024-12", AccountCode: "1111", CurrencyCode: "VND", DebitAmount: 20_000_000, CreditAmount: 0, Status: domain.OBStatusApproved, CreatedBy: "U001"},
	}
	require.NoError(t, svc.BulkCreateOpeningBalances(ctx, balances))

	// Run year-end close
	result, err := svc.YearEndClose(ctx, "C001", "P-2024-12", "P-2025-01", "2024", "2025", "U001")
	require.NoError(t, err)

	// Verify negative income (loss)
	assert.Equal(t, -20_000_000.0, result.NetIncome)
}

func TestYearEndClose_WithAccountMapping(t *testing.T) {
	svc, ctx := setupYearEndService(t)

	// Create accounts (old and new)
	oldAcct := &domain.Account{Code: "1121", Name: "Bank Old", Type: domain.AccountTypeAsset}
	newAcct := &domain.Account{Code: "11211", Name: "Bank New", Type: domain.AccountTypeAsset}
	require.NoError(t, svc.CreateAccount(ctx, oldAcct))
	require.NoError(t, svc.CreateAccount(ctx, newAcct))

	// Create TT200→TT99 mapping
	mapping := &domain.Circular99Mapping{
		OldAccountCode: "1121",
		NewAccountCode: "11211",
		MappingType:    "DIRECT",
		EffectiveDate:  "2025-01-01",
		IsActive:       true,
	}
	require.NoError(t, svc.CreateCircular99Mapping(ctx, mapping))

	// Create opening balance with old account code
	balances := []domain.OpeningBalance{
		{CompanyID: "C001", PeriodID: "P-2024-12", AccountCode: "1121", CurrencyCode: "VND", DebitAmount: 50_000_000, CreditAmount: 0, Status: domain.OBStatusApproved, CreatedBy: "U001"},
	}
	require.NoError(t, svc.BulkCreateOpeningBalances(ctx, balances))

	// Run year-end close
	result, err := svc.YearEndClose(ctx, "C001", "P-2024-12", "P-2025-01", "2024", "2025", "U001")
	require.NoError(t, err)

	// Verify mapping was applied
	assert.Equal(t, 1, result.MappingApplied)
}

func TestExportYearEndBalances(t *testing.T) {
	svc, ctx := setupYearEndService(t)

	// Create opening balances
	balances := []domain.OpeningBalance{
		{CompanyID: "C001", PeriodID: "P-2024-12", AccountCode: "1111", CurrencyCode: "VND", DebitAmount: 40_000_000, Status: domain.OBStatusApproved, CreatedBy: "U001"},
		{CompanyID: "C001", PeriodID: "P-2024-12", AccountCode: "5111", CurrencyCode: "VND", CreditAmount: 100_000_000, Status: domain.OBStatusApproved, CreatedBy: "U001"},
	}
	require.NoError(t, svc.BulkCreateOpeningBalances(ctx, balances))

	// Export
	exported, err := svc.ExportYearEndBalances(ctx, "C001", "P-2024-12")
	require.NoError(t, err)
	assert.Len(t, exported, 2)
}
