package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/domain"
	"gotax/internal/repository"
)

func setupExportTest(t *testing.T) (*ExportService, context.Context) {
	t.Helper()
	accRepo := repository.NewMemoryAccountRepo()
	jeRepo := repository.NewMemoryJournalRepo()
	jeRepo.SetAccounts(accRepo.Accounts())
	perRepo := repository.NewMemoryPeriodRepo()
	ctx := context.Background()

	// Create a period for 2025-01
	period := &domain.Period{
		ID:        "period-2025-01",
		Year:      2025,
		Month:     1,
		StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC),
		Status:    domain.PeriodOpen,
	}
	err := perRepo.Create(ctx, period)
	require.NoError(t, err)

	// Create accounts
	acc1 := &domain.Account{Code: "1111", Name: "Tiền mặt", Type: domain.AccountTypeAsset, IsActive: true}
	acc2 := &domain.Account{Code: "5111", Name: "Doanh thu bán hàng", Type: domain.AccountTypeRevenue, IsActive: true}
	require.NoError(t, accRepo.Create(ctx, acc1))
	require.NoError(t, accRepo.Create(ctx, acc2))

	// Create a posted journal entry
	entry := &domain.JournalEntry{
		ID:          "je-001",
		CompanyID:   "comp-1",
		EntryNumber: "PBT/2025/01/001",
		VoucherType: domain.VoucherTypeOther,
		EntryDate:   time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
		PeriodID:    "period-2025-01",
		Description: "Test entry",
		Status:      domain.JournalEntryPosted,
		Lines: []domain.JournalLine{
			{LineNumber: 1, AccountCode: "1111", DebitAmount: 1000000, CreditAmount: 0, Description: "Debit line"},
			{LineNumber: 2, AccountCode: "5111", DebitAmount: 0, CreditAmount: 1000000, Description: "Credit line", ObjectID: "OBJ-001"},
		},
	}
	require.NoError(t, jeRepo.Create(ctx, entry))

	svc := NewExportService(jeRepo, accRepo, perRepo)
	return svc, ctx
}

func TestExportJournalEntries(t *testing.T) {
	svc, ctx := setupExportTest(t)

	data, err := svc.ExportJournalEntries(ctx, "comp-1", 2025, 1)
	require.NoError(t, err)
	assert.NotEmpty(t, data)
	// xlsx files start with PK zip signature
	assert.Equal(t, byte('P'), data[0])
	assert.Equal(t, byte('K'), data[1])
}

func TestExportJournalEntries_NoData(t *testing.T) {
	svc, ctx := setupExportTest(t)

	// Period 2025-02 doesn't exist → error
	_, err := svc.ExportJournalEntries(ctx, "comp-1", 2025, 2)
	assert.Error(t, err)
}

func TestExportTrialBalance(t *testing.T) {
	svc, ctx := setupExportTest(t)

	data, err := svc.ExportTrialBalance(ctx, "comp-1", 2025, 1)
	require.NoError(t, err)
	assert.NotEmpty(t, data)
	assert.Equal(t, byte('P'), data[0])
	assert.Equal(t, byte('K'), data[1])
}

func TestExportTrialBalance_NoPeriod(t *testing.T) {
	svc, ctx := setupExportTest(t)

	_, err := svc.ExportTrialBalance(ctx, "comp-1", 2099, 12)
	// GetTrialBalance returns empty if period not found, so this should succeed with empty data
	assert.NoError(t, err)
}
