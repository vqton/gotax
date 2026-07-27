package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/domain"
	"gotax/internal/repository"
)

func setupOBService(t *testing.T) (Service, context.Context) {
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

	svc := NewService(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo,
		approvalRepo, versionRepo, mappingRepo, analysisRepo, ifrsRepo, refreshRepo, resetRepo, obRepo)
	return svc, context.Background()
}

func validOB() *domain.OpeningBalance {
	return &domain.OpeningBalance{
		CompanyID:    "C001",
		PeriodID:     "P-2025-01",
		AccountCode:  "1111",
		CurrencyCode: "VND",
		DebitAmount:  100_000_000,
		SourceType:   "MANUAL",
		CreatedBy:    "U001",
	}
}

// ─── CRUD ──────────────────────────────────────────────────────────

func TestCreateOpeningBalance_Success(t *testing.T) {
	svc, ctx := setupOBService(t)
	ob := validOB()
	assert.NoError(t, svc.CreateOpeningBalance(ctx, ob))
	assert.NotEmpty(t, ob.ID)
	assert.Equal(t, domain.OBStatusDraft, ob.Status)
}

func TestCreateOpeningBalance_Validation(t *testing.T) {
	svc, ctx := setupOBService(t)
	tests := []struct {
		name string
		ob   *domain.OpeningBalance
		err  error
	}{
		{"empty company", &domain.OpeningBalance{PeriodID: "P1", AccountCode: "1111", DebitAmount: 100}, domain.ErrCompanyIDRequired},
		{"empty period", &domain.OpeningBalance{CompanyID: "C1", AccountCode: "1111", DebitAmount: 100}, domain.ErrPeriodNotFound},
		{"empty account", &domain.OpeningBalance{CompanyID: "C1", PeriodID: "P1", DebitAmount: 100}, domain.ErrAccountCodeRequired},
		{"zero amount", &domain.OpeningBalance{CompanyID: "C1", PeriodID: "P1", AccountCode: "1111"}, domain.ErrAmountRequired},
		{"both debit credit", &domain.OpeningBalance{CompanyID: "C1", PeriodID: "P1", AccountCode: "1111", DebitAmount: 100, CreditAmount: 100}, domain.ErrBothDebitAndCredit},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.ErrorIs(t, svc.CreateOpeningBalance(ctx, tc.ob), tc.err)
		})
	}
}

func TestGetOpeningBalance_Success(t *testing.T) {
	svc, ctx := setupOBService(t)
	ob := validOB()
	require.NoError(t, svc.CreateOpeningBalance(ctx, ob))
	got, err := svc.GetOpeningBalance(ctx, ob.ID)
	require.NoError(t, err)
	assert.Equal(t, ob.AccountCode, got.AccountCode)
}

func TestGetOpeningBalance_NotFound(t *testing.T) {
	svc, ctx := setupOBService(t)
	_, err := svc.GetOpeningBalance(ctx, "nonexistent")
	assert.ErrorIs(t, err, domain.ErrOpeningBalanceNotFound)
}

func TestListOpeningBalances(t *testing.T) {
	svc, ctx := setupOBService(t)
	require.NoError(t, svc.CreateOpeningBalance(ctx, validOB()))
	ob2 := validOB()
	ob2.AccountCode = "2111"
	ob2.CreditAmount = 50_000_000
	ob2.DebitAmount = 0
	require.NoError(t, svc.CreateOpeningBalance(ctx, ob2))
	list, err := svc.ListOpeningBalances(ctx, domain.OBListFilter{CompanyID: "C001", PeriodID: "P-2025-01"})
	require.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestUpdateOpeningBalance_Success(t *testing.T) {
	svc, ctx := setupOBService(t)
	ob := validOB()
	require.NoError(t, svc.CreateOpeningBalance(ctx, ob))
	ob.DebitAmount = 200_000_000
	assert.NoError(t, svc.UpdateOpeningBalance(ctx, ob))
	got, _ := svc.GetOpeningBalance(ctx, ob.ID)
	assert.Equal(t, 200_000_000.0, got.DebitAmount)
}

func TestUpdateOpeningBalance_RejectedAfterApproval(t *testing.T) {
	svc, ctx := setupOBService(t)
	ob := validOB()
	require.NoError(t, svc.CreateOpeningBalance(ctx, ob))
	require.NoError(t, svc.SubmitOpeningBalance(ctx, ob.ID, "U001"))
	require.NoError(t, svc.ApproveOpeningBalance(ctx, ob.ID, "U002"))
	ob.DebitAmount = 999_999
	assert.Error(t, svc.UpdateOpeningBalance(ctx, ob))
}

func TestDeleteOpeningBalance_Success(t *testing.T) {
	svc, ctx := setupOBService(t)
	ob := validOB()
	require.NoError(t, svc.CreateOpeningBalance(ctx, ob))
	assert.NoError(t, svc.DeleteOpeningBalance(ctx, ob.ID))
	_, err := svc.GetOpeningBalance(ctx, ob.ID)
	assert.ErrorIs(t, err, domain.ErrOpeningBalanceNotFound)
}

// ─── Status lifecycle ──────────────────────────────────────────────

func TestSubmitAndApproveOpeningBalance(t *testing.T) {
	svc, ctx := setupOBService(t)
	ob := validOB()
	require.NoError(t, svc.CreateOpeningBalance(ctx, ob))
	assert.NoError(t, svc.SubmitOpeningBalance(ctx, ob.ID, "U001"))

	got, _ := svc.GetOpeningBalance(ctx, ob.ID)
	assert.Equal(t, domain.OBStatusPending, got.Status)
	assert.NoError(t, svc.ApproveOpeningBalance(ctx, ob.ID, "U002"))
	got, _ = svc.GetOpeningBalance(ctx, ob.ID)
	assert.Equal(t, domain.OBStatusApproved, got.Status)
}

func TestApproveWithoutSubmit_Fails(t *testing.T) {
	svc, ctx := setupOBService(t)
	ob := validOB()
	require.NoError(t, svc.CreateOpeningBalance(ctx, ob))
	assert.Error(t, svc.ApproveOpeningBalance(ctx, ob.ID, "U002"))
}

// ─── Correction ────────────────────────────────────────────────────

func TestCorrectOpeningBalance(t *testing.T) {
	svc, ctx := setupOBService(t)
	ob := validOB()
	require.NoError(t, svc.CreateOpeningBalance(ctx, ob))
	require.NoError(t, svc.SubmitOpeningBalance(ctx, ob.ID, "U001"))
	require.NoError(t, svc.ApproveOpeningBalance(ctx, ob.ID, "U002"))

	corrected, err := svc.CorrectOpeningBalance(ctx, ob.ID, "U003", "Wrong amount")
	require.NoError(t, err)
	assert.Equal(t, domain.OBStatusDraft, corrected.Status)
	assert.Equal(t, ob.ID, corrected.CorrectionOf)
	assert.Equal(t, "Wrong amount", corrected.CorrectionReason)

	original, _ := svc.GetOpeningBalance(ctx, ob.ID)
	assert.Equal(t, domain.OBStatusCorrected, original.Status)
}

// ─── Details ───────────────────────────────────────────────────────

func TestOpeningBalanceDetailLifecycle(t *testing.T) {
	svc, ctx := setupOBService(t)
	ob := validOB()
	require.NoError(t, svc.CreateOpeningBalance(ctx, ob))

	det := &domain.OpeningBalanceDetail{
		OpeningBalanceID: ob.ID,
		EntityType:       domain.DetailFixedAsset,
		EntityID:         "FA001",
		EntityName:       "Server",
		DebitAmount:      100_000_000,
		OriginalCost:     100_000_000,
	}
	require.NoError(t, svc.CreateOpeningBalanceDetail(ctx, det))
	assert.NotEmpty(t, det.ID)

	details, err := svc.GetOpeningBalanceDetails(ctx, ob.ID)
	require.NoError(t, err)
	assert.Len(t, details, 1)

	assert.NoError(t, svc.DeleteOpeningBalanceDetail(ctx, det.ID))
	details, _ = svc.GetOpeningBalanceDetails(ctx, ob.ID)
	assert.Len(t, details, 0)
}

// ─── Totals & balance ──────────────────────────────────────────────

func TestGetTotals_AfterApproval(t *testing.T) {
	svc, ctx := setupOBService(t)
	ob1 := validOB()
	require.NoError(t, svc.CreateOpeningBalance(ctx, ob1))
	require.NoError(t, svc.SubmitOpeningBalance(ctx, ob1.ID, "U001"))
	require.NoError(t, svc.ApproveOpeningBalance(ctx, ob1.ID, "U002"))

	ob2 := validOB()
	ob2.AccountCode = "2111"
	ob2.CreditAmount = 100_000_000
	ob2.DebitAmount = 0
	require.NoError(t, svc.CreateOpeningBalance(ctx, ob2))
	require.NoError(t, svc.SubmitOpeningBalance(ctx, ob2.ID, "U001"))
	require.NoError(t, svc.ApproveOpeningBalance(ctx, ob2.ID, "U002"))

	debit, credit, err := svc.GetOpeningBalanceTotals(ctx, "C001", "P-2025-01")
	require.NoError(t, err)
	assert.Equal(t, 100_000_000.0, debit)
	assert.Equal(t, 100_000_000.0, credit)

	balanced, err := svc.ValidateOpeningBalancesBalanced(ctx, "C001", "P-2025-01")
	require.NoError(t, err)
	assert.True(t, balanced)
}

// ─── Carry Forward ─────────────────────────────────────────────────

func TestCarryForwardAndLogs(t *testing.T) {
	svc, ctx := setupOBService(t)
	ob := validOB()
	ob.PeriodID = "P-2024-12"
	require.NoError(t, svc.CreateOpeningBalance(ctx, ob))
	require.NoError(t, svc.SubmitOpeningBalance(ctx, ob.ID, "U001"))
	require.NoError(t, svc.ApproveOpeningBalance(ctx, ob.ID, "U002"))

	log, err := svc.CarryForward(ctx, "C001", "P-2024-12", "P-2025-01", "2024", "2025", "U001")
	require.NoError(t, err)
	assert.Equal(t, "COMPLETED", log.Status)
	assert.Equal(t, 1, log.AccountCount)

	logs, err := svc.GetCarryForwardLogs(ctx, "C001")
	require.NoError(t, err)
	assert.Len(t, logs, 1)

	got, err := svc.GetCarryForwardLogByID(ctx, log.ID)
	require.NoError(t, err)
	assert.Equal(t, log.ID, got.ID)
}

// ─── Circular 99 Mapping ───────────────────────────────────────────

func TestCircular99MappingLifecycle(t *testing.T) {
	svc, ctx := setupOBService(t)
	m := &domain.Circular99Mapping{
		OldAccountCode: "1121",
		NewAccountCode: "11211",
		MappingType:    "SPLIT",
		SplitRatio:     0.7,
		EffectiveDate:  "2025-01-01",
	}
	require.NoError(t, svc.CreateCircular99Mapping(ctx, m))
	assert.NotEmpty(t, m.ID)

	got, err := svc.GetCircular99MappingByOldCode(ctx, "1121")
	require.NoError(t, err)
	assert.Equal(t, "11211", got.NewAccountCode)

	list, err := svc.ListCircular99Mappings(ctx)
	require.NoError(t, err)
	assert.Len(t, list, 1)
}

// ─── Balance Migration ─────────────────────────────────────────────

func TestBalanceMigrationLifecycle(t *testing.T) {
	svc, ctx := setupOBService(t)
	m := &domain.BalanceMigration{
		CompanyID:     "C001",
		FromRegime:    "VAS",
		ToRegime:      "IFRS",
		ExecutionDate: "2025-06-30",
		ExecutedBy:    "U001",
	}
	require.NoError(t, svc.CreateBalanceMigration(ctx, m))
	assert.NotEmpty(t, m.ID)

	got, err := svc.GetBalanceMigrationByID(ctx, m.ID)
	require.NoError(t, err)
	assert.Equal(t, "VAS", got.FromRegime)

	list, err := svc.ListBalanceMigrations(ctx, "C001")
	require.NoError(t, err)
	assert.Len(t, list, 1)
}

// ─── Bulk Operations ───────────────────────────────────────────────

func TestBulkCreateOpeningBalances(t *testing.T) {
	svc, ctx := setupOBService(t)
	obs := []domain.OpeningBalance{
		{CompanyID: "C001", PeriodID: "P-2025-01", AccountCode: "1111", DebitAmount: 100_000_000, CreatedBy: "U001"},
		{CompanyID: "C001", PeriodID: "P-2025-01", AccountCode: "2111", CreditAmount: 100_000_000, CreatedBy: "U001"},
	}
	assert.NoError(t, svc.BulkCreateOpeningBalances(ctx, obs))
	assert.NotEmpty(t, obs[0].ID)
	assert.NotEmpty(t, obs[1].ID)
}

func TestBulkApproveOpeningBalances(t *testing.T) {
	svc, ctx := setupOBService(t)
	obs := []domain.OpeningBalance{
		{CompanyID: "C001", PeriodID: "P-2025-01", AccountCode: "1111", DebitAmount: 100_000_000, CreatedBy: "U001"},
		{CompanyID: "C001", PeriodID: "P-2025-01", AccountCode: "2111", CreditAmount: 100_000_000, CreatedBy: "U001"},
	}
	require.NoError(t, svc.BulkCreateOpeningBalances(ctx, obs))
	ids := []string{obs[0].ID, obs[1].ID}
	require.NoError(t, svc.BulkSubmitOpeningBalances(ctx, ids, "U002"))
	assert.NoError(t, svc.BulkApproveOpeningBalances(ctx, ids, "U003"))

	got, _ := svc.GetOpeningBalance(ctx, obs[0].ID)
	assert.Equal(t, domain.OBStatusApproved, got.Status)
}
