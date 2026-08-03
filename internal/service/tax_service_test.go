package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/domain"
	"gotax/internal/repository"
)

func newTaxTestSvc() (*taxService, domain.TaxRepository) {
	repo := repository.NewMemoryTaxRepo()
	return NewTaxService(repo).(*taxService), repo
}

// ─── A1: Rate Resolver ──────────────────────────────────────────────────

func TestResolveRate_ByRateCodeSuffix(t *testing.T) {
	svc, repo := newTaxTestSvc()
	ctx := context.Background()
	require.NoError(t, repo.CreateRate(ctx, &domain.TaxRate{
		RateCode: "CIT_STANDARD", TaxType: domain.TaxTypeCIT,
		RateType: domain.RateTypePERCENTAGE, RateValue: 20,
		EffectiveFrom: "2025-01-01", IsActive: true,
	}))

	rate, err := svc.resolveRate(ctx, domain.TaxTypeCIT, "STANDARD", "2026-06-01")
	require.NoError(t, err)
	assert.Equal(t, 20.0, rate.RateValue)
}

func TestResolveRate_ByApplicableToField(t *testing.T) {
	svc, repo := newTaxTestSvc()
	ctx := context.Background()
	require.NoError(t, repo.CreateRate(ctx, &domain.TaxRate{
		RateCode: "VAT8", TaxType: domain.TaxTypeVAT,
		RateType: domain.RateTypePERCENTAGE, RateValue: 8,
		ApplicableTo: "REDUCED", EffectiveFrom: "2025-07-01", IsActive: true,
	}))

	rate, err := svc.resolveRate(ctx, domain.TaxTypeVAT, "REDUCED", "2026-06-01")
	require.NoError(t, err)
	assert.Equal(t, 8.0, rate.RateValue)
}

func TestResolveRate_RespectsEffectiveWindow(t *testing.T) {
	svc, repo := newTaxTestSvc()
	ctx := context.Background()
	// legacy rate, superseded on 2026-01-01
	require.NoError(t, repo.CreateRate(ctx, &domain.TaxRate{
		RateCode: "CIT_STANDARD_LEGACY", TaxType: domain.TaxTypeCIT,
		RateType: domain.RateTypePERCENTAGE, RateValue: 20,
		EffectiveFrom: "2020-01-01", EffectiveTo: "2025-12-31", IsActive: true,
	}))
	require.NoError(t, repo.CreateRate(ctx, &domain.TaxRate{
		RateCode: "CIT_STANDARD", TaxType: domain.TaxTypeCIT,
		RateType: domain.RateTypePERCENTAGE, RateValue: 17,
		EffectiveFrom: "2026-01-01", IsActive: true,
	}))

	rate, err := svc.resolveRate(ctx, domain.TaxTypeCIT, "STANDARD", "2026-06-01")
	require.NoError(t, err)
	assert.Equal(t, 17.0, rate.RateValue)

	legacy, err := svc.resolveRate(ctx, domain.TaxTypeCIT, "STANDARD", "2025-06-01")
	require.NoError(t, err)
	assert.Equal(t, 20.0, legacy.RateValue)
}

func TestResolveRate_SkipsInactive(t *testing.T) {
	svc, repo := newTaxTestSvc()
	ctx := context.Background()
	require.NoError(t, repo.CreateRate(ctx, &domain.TaxRate{
		RateCode: "CIT_STANDARD", TaxType: domain.TaxTypeCIT,
		RateType: domain.RateTypePERCENTAGE, RateValue: 25,
		EffectiveFrom: "2025-01-01", IsActive: false,
	}))

	rate, err := svc.resolveRate(ctx, domain.TaxTypeCIT, "STANDARD", "2026-06-01")
	require.NoError(t, err)
	// statutory fallback for CIT STANDARD
	assert.Equal(t, 20.0, rate.RateValue)
}

func TestResolveRate_StatutoryFallback(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()

	rate, err := svc.resolveRate(ctx, domain.TaxTypeVAT, "STANDARD", "2026-06-01")
	require.NoError(t, err)
	assert.Equal(t, 10.0, rate.RateValue)

	cit, err := svc.resolveRate(ctx, domain.TaxTypeCIT, "SMALL", "2026-06-01")
	require.NoError(t, err)
	assert.Equal(t, 17.0, cit.RateValue)
}

// ─── A2: VAT Engine ─────────────────────────────────────────────────────

func vatPeriod() domain.TaxPeriod {
	return domain.TaxPeriod{PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 1}
}

func vatEntry(num string, lines []domain.JournalLine) domain.JournalEntry {
	return domain.JournalEntry{EntryNumber: num, Lines: lines}
}

func vatLine(account string, dr, cr float64) domain.JournalLine {
	return domain.JournalLine{AccountCode: account, DebitAmount: dr, CreditAmount: cr}
}

func TestCalculateVAT_FallbackUsesRateTable(t *testing.T) {
	svc, repo := newTaxTestSvc()
	ctx := context.Background()
	require.NoError(t, repo.CreateRate(ctx, &domain.TaxRate{
		RateCode: "VAT_STANDARD", TaxType: domain.TaxTypeVAT,
		RateType: domain.RateTypePERCENTAGE, RateValue: 8,
		ApplicableTo: "STANDARD", EffectiveFrom: "2026-01-01", IsActive: true,
	}))

	res, err := svc.CalculateVAT(ctx, "c1", vatPeriod(), []domain.JournalEntry{
		vatEntry("E1", []domain.JournalLine{vatLine("5111", 0, 10000000)}),
	})
	require.NoError(t, err)
	assert.Equal(t, 800000.0, res.OutputVAT) // 10M * 8% from rate table
	assert.Equal(t, 800000.0, res.VATPayable)
}

func TestCalculateVAT_ExplicitAccountsWin(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()
	res, err := svc.CalculateVAT(ctx, "c1", vatPeriod(), []domain.JournalEntry{
		vatEntry("E1", []domain.JournalLine{
			vatLine("152", 5000000, 0),
			vatLine("1331", 500000, 0),
			vatLine("331", 0, 5500000),
		}),
	})
	require.NoError(t, err)
	// explicit 1331 used; no double-count via 152 fallback
	assert.Equal(t, 500000.0, res.InputVAT)
	assert.Equal(t, 5000000.0, res.PurchaseTotal)
}

func TestCalculateVAT_ExplicitOutput(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()
	res, err := svc.CalculateVAT(ctx, "c1", vatPeriod(), []domain.JournalEntry{
		vatEntry("E1", []domain.JournalLine{
			vatLine("131", 11000000, 0),
			vatLine("5111", 0, 10000000),
			vatLine("33311", 0, 1000000),
		}),
	})
	require.NoError(t, err)
	assert.Equal(t, 1000000.0, res.OutputVAT) // explicit 33311, no rate applied
	assert.Equal(t, 10000000.0, res.SalesTotal)
}

func TestCalculateVAT_FAInput(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()
	res, err := svc.CalculateVAT(ctx, "c1", vatPeriod(), []domain.JournalEntry{
		vatEntry("E1", []domain.JournalLine{vatLine("2111", 10000000, 0)}),
	})
	require.NoError(t, err)
	assert.Equal(t, 1000000.0, res.InputVATFA) // 10M * 10% default
	assert.Equal(t, 1000000.0, res.TotalInputVAT)
}

func TestCalculateVAT_Refundable(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()
	res, err := svc.CalculateVAT(ctx, "c1", vatPeriod(), []domain.JournalEntry{
		vatEntry("E1", []domain.JournalLine{vatLine("5111", 0, 2000000)}),
		vatEntry("E2", []domain.JournalLine{vatLine("152", 3000000, 0)}),
	})
	require.NoError(t, err)
	assert.Equal(t, 200000.0, res.OutputVAT)
	assert.Equal(t, 300000.0, res.InputVAT)
	assert.Equal(t, 0.0, res.VATPayable)
	assert.Equal(t, 100000.0, res.VATRefundable)
}

// ─── A3: CIT Engine ─────────────────────────────────────────────────────

func citEntry(lines ...domain.JournalLine) domain.JournalEntry {
	return domain.JournalEntry{Lines: lines}
}

func TestCalculateCIT_MicroRate(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()
	res, err := svc.CalculateCIT(ctx, "c1", 2026, []domain.JournalEntry{
		citEntry(vatLine("5111", 0, 100000000)),  // 100M < 3B → MICRO 15%
		citEntry(vatLine("641", 60000000, 0)),
	})
	require.NoError(t, err)
	assert.Equal(t, 15.0, res.TaxRate)
	assert.Equal(t, 40000000.0, res.TaxableIncome)
	assert.Equal(t, 6000000.0, res.CITPayable)
}

func TestCalculateCIT_SmallRate(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()
	res, err := svc.CalculateCIT(ctx, "c1", 2026, []domain.JournalEntry{
		citEntry(vatLine("5111", 0, 5000000000)),  // 5B → SMALL 17%
		citEntry(vatLine("641", 2000000000, 0)),
	})
	require.NoError(t, err)
	assert.Equal(t, 17.0, res.TaxRate)
	assert.Equal(t, 510000000.0, res.CITPayable) // 3B * 17%
}

func TestCalculateCIT_StandardRate(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()
	res, err := svc.CalculateCIT(ctx, "c1", 2026, []domain.JournalEntry{
		citEntry(vatLine("5111", 0, 60000000000)),  // 60B → STANDARD 20%
		citEntry(vatLine("641", 30000000000, 0)),
	})
	require.NoError(t, err)
	assert.Equal(t, 20.0, res.TaxRate)
	assert.Equal(t, 6000000000.0, res.CITPayable) // 30B * 20%
}

func TestCalculateCIT_NonDeductibleAddedBack(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()
	res, err := svc.CalculateCIT(ctx, "c1", 2026, []domain.JournalEntry{
		citEntry(vatLine("5111", 0, 100000000)),
		citEntry(vatLine("641", 60000000, 0)),
		citEntry(vatLine("821", 5000000, 0)),
	})
	require.NoError(t, err)
	assert.Equal(t, 45000000.0, res.TaxableIncome) // 100M - 60M + 5M
	assert.Equal(t, 6750000.0, res.CITPayable)     // 45M * 15%
}

func TestCalculateCIT_Loss(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()
	res, err := svc.CalculateCIT(ctx, "c1", 2026, []domain.JournalEntry{
		citEntry(vatLine("5111", 0, 50000000)),
		citEntry(vatLine("641", 90000000, 0)),
	})
	require.NoError(t, err)
	assert.Equal(t, 0.0, res.TaxableIncome)
	assert.Equal(t, 0.0, res.CITPayable)
}

func TestCalculateCIT_RateTableOverridesSize(t *testing.T) {
	svc, repo := newTaxTestSvc()
	ctx := context.Background()
	require.NoError(t, repo.CreateRate(ctx, &domain.TaxRate{
		RateCode: "CIT_MICRO", TaxType: domain.TaxTypeCIT,
		RateType: domain.RateTypePERCENTAGE, RateValue: 12.5,
		ApplicableTo: "MICRO", EffectiveFrom: "2026-01-01", IsActive: true,
	}))
	res, err := svc.CalculateCIT(ctx, "c1", 2026, []domain.JournalEntry{
		citEntry(vatLine("5111", 0, 100000000)),
		citEntry(vatLine("641", 60000000, 0)),
	})
	require.NoError(t, err)
	assert.Equal(t, 12.5, res.TaxRate)
	assert.Equal(t, 5000000.0, res.CITPayable) // 40M * 12.5%
}
