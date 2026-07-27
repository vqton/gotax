package gl

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryAccountRepo_CRUD(t *testing.T) {
	repo := NewMemoryAccountRepo()
	ctx := context.Background()

	t.Run("create and get by code", func(t *testing.T) {
		acc := &Account{Code: "1111", Name: "TM VND", Type: AccountTypeAsset, IsActive: true}
		err := repo.Create(ctx, acc)
		require.NoError(t, err)

		found, err := repo.GetByCode(ctx, "1111")
		require.NoError(t, err)
		assert.Equal(t, "1111", found.Code)
		assert.Equal(t, "TM VND", found.Name)
	})

	t.Run("get not found", func(t *testing.T) {
		_, err := repo.GetByCode(ctx, "NONEXIST")
		assert.ErrorIs(t, err, ErrAccountNotFound)
	})

	t.Run("get all with active filter", func(t *testing.T) {
		repo.Create(ctx, &Account{Code: "1112", Name: "TM USD", Type: AccountTypeAsset, IsActive: true})
		repo.Create(ctx, &Account{Code: "5111", Name: "DT BH", Type: AccountTypeRevenue, IsActive: false})

		all, err := repo.GetAll(ctx, false)
		require.NoError(t, err)
		assert.Len(t, all, 3)

		active, err := repo.GetAll(ctx, true)
		require.NoError(t, err)
		assert.Len(t, active, 2)
	})

	t.Run("update account", func(t *testing.T) {
		err := repo.Update(ctx, &Account{Code: "1111", Name: "TM VND Updated", Type: AccountTypeAsset, IsActive: true})
		require.NoError(t, err)

		found, _ := repo.GetByCode(ctx, "1111")
		assert.Equal(t, "TM VND Updated", found.Name)
	})

	t.Run("get children", func(t *testing.T) {
		repo.Create(ctx, &Account{Code: "111", Name: "Tien mat", Type: AccountTypeAsset, IsParent: true})
		repo.Create(ctx, &Account{Code: "1111", Name: "TM VND", Type: AccountTypeAsset, ParentCode: "111"})
		repo.Create(ctx, &Account{Code: "1112", Name: "TM USD", Type: AccountTypeAsset, ParentCode: "111"})

		children, err := repo.GetChildren(ctx, "111")
		require.NoError(t, err)
		assert.Len(t, children, 2)
	})

	t.Run("delete account", func(t *testing.T) {
		err := repo.Delete(ctx, "1112")
		require.NoError(t, err)

		_, err = repo.GetByCode(ctx, "1112")
		assert.ErrorIs(t, err, ErrAccountNotFound)
	})
}

func TestMemoryJournalRepo_CRUD(t *testing.T) {
	jeRepo := NewMemoryJournalRepo()
	accRepo := NewMemoryAccountRepo()
	ctx := context.Background()

	accRepo.Create(ctx, &Account{Code: "1111", Name: "TM VND", Type: AccountTypeAsset, IsActive: true})
	accRepo.Create(ctx, &Account{Code: "5111", Name: "DT BH", Type: AccountTypeRevenue, IsActive: true})
	accRepo.Create(ctx, &Account{Code: "1561", Name: "Hang hoa", Type: AccountTypeExpense, IsActive: true})

	t.Run("create journal entry", func(t *testing.T) {
		entry := &JournalEntry{
			EntryDate:   time.Date(2026, 7, 24, 0, 0, 0, 0, time.UTC),
			Description: "Mua hang tra tien mat",
			Status:      JournalEntryPosted,
			PeriodID:    "P-2026-07",
			Lines: []JournalLine{
				{AccountCode: "1561", DebitAmount: 10000000, CreditAmount: 0, Description: "Mua hang", LineNumber: 1},
				{AccountCode: "1111", DebitAmount: 0, CreditAmount: 10000000, Description: "Tra tien", LineNumber: 2},
			},
		}

		err := jeRepo.Create(ctx, entry)
		require.NoError(t, err)
		assert.NotEmpty(t, entry.ID)
		assert.Contains(t, entry.ID, "JE-")
		assert.NotZero(t, entry.CreatedAt)
	})

	t.Run("get by id", func(t *testing.T) {
		entry := &JournalEntry{
			EntryDate:   time.Date(2026, 7, 24, 0, 0, 0, 0, time.UTC),
			Description: "Test",
			Status:      JournalEntryDraft,
			Lines:       []JournalLine{{AccountCode: "1111", DebitAmount: 1000, CreditAmount: 0}, {AccountCode: "5111", DebitAmount: 0, CreditAmount: 1000}},
		}
		jeRepo.Create(ctx, entry)

		found, err := jeRepo.GetByID(ctx, entry.ID)
		require.NoError(t, err)
		assert.Equal(t, entry.ID, found.ID)
		assert.Len(t, found.Lines, 2)
	})

	t.Run("get not found", func(t *testing.T) {
		_, err := jeRepo.GetByID(ctx, "NONEXIST")
		assert.ErrorIs(t, err, ErrJournalNotFound)
	})

	t.Run("get by period", func(t *testing.T) {
		entries, err := jeRepo.GetByPeriod(ctx, "P-2026-07")
		require.NoError(t, err)
		assert.NotZero(t, len(entries))
	})

	t.Run("update status", func(t *testing.T) {
		entry := &JournalEntry{
			EntryDate:   time.Date(2026, 7, 24, 0, 0, 0, 0, time.UTC),
			Description: "To cancel",
			Status:      JournalEntryPosted,
			Lines:       []JournalLine{{AccountCode: "1111", DebitAmount: 1000, CreditAmount: 0}, {AccountCode: "5111", DebitAmount: 0, CreditAmount: 1000}},
		}
		jeRepo.Create(ctx, entry)

		err := jeRepo.UpdateStatus(ctx, entry.ID, JournalEntryCancelled)
		require.NoError(t, err)

		found, _ := jeRepo.GetByID(ctx, entry.ID)
		assert.Equal(t, JournalEntryCancelled, found.Status)
	})

	t.Run("get trial balance", func(t *testing.T) {
		// Create a period with a couple posted entries
		jeRepo.Create(ctx, &JournalEntry{
			EntryDate: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC),
			Description: "Sale",
			Status: JournalEntryPosted,
			PeriodID: "P-2026-08",
			Lines: []JournalLine{
				{AccountCode: "1111", DebitAmount: 5000000, CreditAmount: 0},
				{AccountCode: "5111", DebitAmount: 0, CreditAmount: 5000000},
			},
		})
		jeRepo.Create(ctx, &JournalEntry{
			EntryDate: time.Date(2026, 8, 2, 0, 0, 0, 0, time.UTC),
			Description: "Expense",
			Status: JournalEntryPosted,
			PeriodID: "P-2026-08",
			Lines: []JournalLine{
				{AccountCode: "1561", DebitAmount: 2000000, CreditAmount: 0},
				{AccountCode: "1111", DebitAmount: 0, CreditAmount: 2000000},
			},
		})

		balances, err := jeRepo.GetTrialBalance(ctx, "P-2026-08")
		require.NoError(t, err)
		assert.Len(t, balances, 3)

		for _, b := range balances {
			switch b.AccountCode {
			case "1111":
				assert.Equal(t, 5000000.0, b.PeriodDebit)
				assert.Equal(t, 2000000.0, b.PeriodCredit)
			case "5111":
				assert.Equal(t, 0.0, b.PeriodDebit)
				assert.Equal(t, 5000000.0, b.PeriodCredit)
			case "1561":
				assert.Equal(t, 2000000.0, b.PeriodDebit)
				assert.Equal(t, 0.0, b.PeriodCredit)
			}
		}
	})

	t.Run("draft entry not counted in trial balance", func(t *testing.T) {
		jeRepo.Create(ctx, &JournalEntry{
			EntryDate: time.Date(2026, 8, 3, 0, 0, 0, 0, time.UTC),
			Description: "Draft - not counted",
			Status: JournalEntryDraft,
			PeriodID: "P-2026-08",
			Lines: []JournalLine{
				{AccountCode: "1111", DebitAmount: 999999, CreditAmount: 0},
				{AccountCode: "5111", DebitAmount: 0, CreditAmount: 999999},
			},
		})

		balances, _ := jeRepo.GetTrialBalance(ctx, "P-2026-08")
		for _, b := range balances {
			if b.AccountCode == "1111" {
				assert.Equal(t, 5000000.0, b.PeriodDebit)
				break
			}
		}
	})
}

func TestMemoryPeriodRepo_CRUD(t *testing.T) {
	repo := NewMemoryPeriodRepo()
	ctx := context.Background()

	t.Run("create period", func(t *testing.T) {
		p := &Period{Year: 2026, Month: 7, StartDate: time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2026, 7, 31, 0, 0, 0, 0, time.UTC), Status: PeriodOpen}
		err := repo.Create(ctx, p)
		require.NoError(t, err)
		assert.Equal(t, "P-2026-07", p.ID)
	})

	t.Run("get by id", func(t *testing.T) {
		p, err := repo.GetByID(ctx, "P-2026-07")
		require.NoError(t, err)
		assert.Equal(t, 2026, p.Year)
		assert.Equal(t, 7, p.Month)
	})

	t.Run("get by year month", func(t *testing.T) {
		p, err := repo.GetByYearMonth(ctx, 2026, 7)
		require.NoError(t, err)
		assert.Equal(t, "P-2026-07", p.ID)
	})

	t.Run("get all", func(t *testing.T) {
		repo.Create(ctx, &Period{Year: 2026, Month: 1, StartDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC), Status: PeriodOpen})
		repo.Create(ctx, &Period{Year: 2026, Month: 12, StartDate: time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC), Status: PeriodOpen})

		all, err := repo.GetAll(ctx)
		require.NoError(t, err)
		assert.Len(t, all, 3)
		assert.Equal(t, 1, all[0].Month)
		assert.Equal(t, 7, all[1].Month)
		assert.Equal(t, 12, all[2].Month)
	})

	t.Run("update status", func(t *testing.T) {
		err := repo.UpdateStatus(ctx, "P-2026-07", PeriodClosed)
		require.NoError(t, err)

		p, _ := repo.GetByID(ctx, "P-2026-07")
		assert.Equal(t, PeriodClosed, p.Status)
	})

	t.Run("period not found", func(t *testing.T) {
		_, err := repo.GetByID(ctx, "P-2030-01")
		assert.ErrorIs(t, err, ErrPeriodNotFound)
	})
}

func TestFullWorkflow(t *testing.T) {
	accRepo := NewMemoryAccountRepo()
	jeRepo := NewMemoryJournalRepo()
	perRepo := NewMemoryPeriodRepo()
	userRepo := NewMemoryUserRepo()
	auditRepo := NewMemoryAuditLogRepo()
	rateRepo := NewMemoryExchangeRateRepo()
	templateRepo := NewMemoryClosingTemplateRepo()
	approvalRepo := NewMemoryApprovalRepo()
	verRepo := NewMemoryAccountVersionRepo()
	mappingRepo := NewMemoryAccountMappingRepo()
	analysisRepo := NewMemoryAccountAnalysisRepo()
	ifrsRepo := NewMemoryIFRSMappingRepo()
	svc := NewService(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo, approvalRepo, verRepo, mappingRepo, analysisRepo, ifrsRepo, nil, nil)
	ctx := context.Background()

	perRepo.Create(ctx, &Period{Year: 2026, Month: 7, StartDate: time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2026, 7, 31, 0, 0, 0, 0, time.UTC), Status: PeriodOpen})

	svc.CreateAccount(ctx, &Account{Code: "1111", Name: "TM VND", Type: AccountTypeAsset, IsActive: true})
	svc.CreateAccount(ctx, &Account{Code: "5111", Name: "DT BH", Type: AccountTypeRevenue, IsActive: true})
	svc.CreateAccount(ctx, &Account{Code: "1561", Name: "Hang hoa", Type: AccountTypeExpense, IsActive: true})

	entry := &JournalEntry{
		EntryDate:   time.Date(2026, 7, 24, 0, 0, 0, 0, time.UTC),
		Description: "Mua hang tra tien mat",
		PeriodID:    "P-2026-07",
		Lines: []JournalLine{
			{AccountCode: "1561", DebitAmount: 10000000, CreditAmount: 0, Description: "Mua hang"},
			{AccountCode: "1111", DebitAmount: 0, CreditAmount: 10000000, Description: "Tra tien"},
		},
	}

	err := svc.CreateEntry(ctx, entry, "user-1")
	require.NoError(t, err)
	assert.Equal(t, JournalEntryDraft, entry.Status)
	assert.NotEmpty(t, entry.ID)

	err = svc.SubmitForReview(ctx, entry.ID, "user-1")
	require.NoError(t, err)

	err = svc.ApproveEntry(ctx, entry.ID, "admin")
	require.NoError(t, err)

	err = svc.PostEntry(ctx, entry.ID)
	require.NoError(t, err)
	assert.Equal(t, JournalEntryPosted, entry.Status)
	assert.NotNil(t, entry.PostedAt)

	tb, err := svc.TrialBalance(ctx, 2026, 7)
	require.NoError(t, err)
	assert.Len(t, tb, 2)
}