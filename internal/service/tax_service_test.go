package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/domain"
	"gotax/internal/einvoice"
	"gotax/internal/repository"
	"gotax/internal/xmldsig"
)

func newTaxTestSvc() (*taxService, domain.TaxRepository) {
	repo := repository.NewMemoryTaxRepo()
	return NewTaxService(repo, repository.NewMemoryJournalRepo(), repository.NewMemoryCompanyRepo(), nil, nil).(*taxService), repo
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
		citEntry(vatLine("5111", 0, 100000000)), // 100M < 3B → MICRO 15%
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
		citEntry(vatLine("5111", 0, 5000000000)), // 5B → SMALL 17%
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
		citEntry(vatLine("5111", 0, 60000000000)), // 60B → STANDARD 20%
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

func TestCalculateCIT_IncentiveReduction(t *testing.T) {
	svc, repo := newTaxTestSvc()
	ctx := context.Background()
	// Standard CIT rate 20%
	require.NoError(t, repo.CreateRate(ctx, &domain.TaxRate{
		RateCode: "CIT_STANDARD", TaxType: domain.TaxTypeCIT,
		RateType: domain.RateTypePERCENTAGE, RateValue: 20,
		EffectiveFrom: "2026-01-01", IsActive: true,
	}))
	// Incentive: 50% reduction for new investment project
	require.NoError(t, repo.CreateRate(ctx, &domain.TaxRate{
		RateCode: "CIT_INCENTIVE_NEW_PROJECT", TaxType: domain.TaxTypeCIT,
		RateType: domain.RateTypePERCENTAGE, RateValue: 20,
		IncentiveReducPct: 50,
		ApplicableTo: "INCENTIVE_NEW_PROJECT",
		EffectiveFrom: "2026-01-01", IsActive: true,
	}))

	res, err := svc.CalculateCIT(ctx, "c1", 2026, []domain.JournalEntry{
		citEntry(vatLine("5111", 0, 100000000)), // 100M revenue
		citEntry(vatLine("641", 60000000, 0)),   // 60M expenses
	})
	require.NoError(t, err)
	// Taxable = 100M - 60M = 40M; CIT = 40M * 20% = 8M
	assert.Equal(t, 8000000.0, res.CITPayable)
	// Incentive: 50% reduction → CITFinal = 8M - 4M = 4M
	assert.Equal(t, 4000000.0, res.CITFinal)
	assert.Equal(t, 50.0, res.IncentiveReduc)
}

func TestCalculateCIT_LossCarryForward(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()
	// Year 2025: 100M revenue, 120M expenses → 20M loss
	res2025, err := svc.CalculateCIT(ctx, "c1", 2025, []domain.JournalEntry{
		citEntry(vatLine("5111", 0, 100000000)),
		citEntry(vatLine("641", 120000000, 0)),
	})
	require.NoError(t, err)
	assert.Equal(t, 0.0, res2025.TaxableIncome)
	assert.Equal(t, 0.0, res2025.CITPayable)

	// Year 2026: 100M revenue, 60M expenses → 40M taxable, minus 20M loss = 20M
	res2026, err := svc.CalculateCIT(ctx, "c1", 2026, []domain.JournalEntry{
		citEntry(vatLine("5111", 0, 100000000)),
		citEntry(vatLine("641", 60000000, 0)),
	})
	require.NoError(t, err)
	assert.Equal(t, 20000000.0, res2026.TaxableIncome) // 40M - 20M loss
	assert.Equal(t, 3000000.0, res2026.CITPayable)      // 20M * 15% (MICRO rate)
	assert.Equal(t, 20000000.0, res2026.LossUsed)
}

func TestCalculateCIT_LossCarryForward_Expiry(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()
	// Year 2020: 100M loss
	_, err := svc.CalculateCIT(ctx, "c1", 2020, []domain.JournalEntry{
		citEntry(vatLine("5111", 0, 100000000)),
		citEntry(vatLine("641", 200000000, 0)),
	})
	require.NoError(t, err)

	// Year 2026: 6 years later — loss expired (5-year limit)
	res, err := svc.CalculateCIT(ctx, "c1", 2026, []domain.JournalEntry{
		citEntry(vatLine("5111", 0, 100000000)),
		citEntry(vatLine("641", 60000000, 0)),
	})
	require.NoError(t, err)
	assert.Equal(t, 40000000.0, res.TaxableIncome) // no loss applied (expired)
	assert.Equal(t, 0.0, res.LossUsed)
}

func TestCalculateCIT_ThinCap(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()
	// Revenue 100M, Operating expenses 40M, Interest expense 30M (account 635)
	// EBITDA = 100M - 40M = 60M; 30% limit = 18M; disallowed = 30M - 18M = 12M
	res, err := svc.CalculateCIT(ctx, "c1", 2026, []domain.JournalEntry{
		citEntry(vatLine("5111", 0, 100000000)),
		citEntry(vatLine("641", 40000000, 0)),
		citEntry(vatLine("635", 30000000, 0)),
	})
	require.NoError(t, err)
	assert.Equal(t, 100000000.0, res.Revenue)
	assert.Equal(t, 70000000.0, res.Expenses) // 40M + 30M
	// NonDeductible includes thin cap disallowance: 12M
	assert.Equal(t, 12000000.0, res.NonDeductible)
	// Taxable = 100M - 70M + 12M = 42M
	assert.Equal(t, 42000000.0, res.TaxableIncome)
}

func TestCalculateCIT_RnDSuperDeduction(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()
	// Revenue 100M, Expenses 60M (incl 10M R&D on account 632)
	// R&D gets 200% deduction → add back 10M extra
	res, err := svc.CalculateCIT(ctx, "c1", 2026, []domain.JournalEntry{
		citEntry(vatLine("5111", 0, 100000000)),
		citEntry(vatLine("641", 50000000, 0)),
		citEntry(vatLine("632", 10000000, 0)),
	})
	require.NoError(t, err)
	assert.Equal(t, 100000000.0, res.Revenue)
	assert.Equal(t, 60000000.0, res.Expenses)
	// R&D 10M already in expenses; NonDeductible = -10M (negative = extra deduction)
	// Taxable = 100M - 60M + (-10M) = 30M
	assert.Equal(t, 30000000.0, res.TaxableIncome)
}

func TestCalculateQuarterlyProvisional(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()
	// Q1: 25M revenue, 15M expenses → 10M taxable → CIT 15% MICRO = 1.5M
	entries := []domain.JournalEntry{
		citEntry(vatLine("5111", 0, 25000000)),
		citEntry(vatLine("641", 15000000, 0)),
	}
	prov, err := svc.CalculateQuarterlyProvisional(ctx, "c1", 2026, 1, entries)
	require.NoError(t, err)
	assert.Equal(t, 6000000.0, prov.EstimatedAnnualCIT) // 1.5M / 1 * 4
	assert.Equal(t, 1500000.0, prov.QuarterlyAmount)     // 6M * 1/4

	// Q3: cumulative 75M revenue, 45M expenses → 30M taxable → CIT = 4.5M
	entriesQ3 := []domain.JournalEntry{
		citEntry(vatLine("5111", 0, 75000000)),
		citEntry(vatLine("641", 45000000, 0)),
	}
	provQ3, err := svc.CalculateQuarterlyProvisional(ctx, "c1", 2026, 3, entriesQ3)
	require.NoError(t, err)
	// Extrapolated annual = 4.5M / 3 * 4 = 6M
	assert.Equal(t, 6000000.0, provQ3.EstimatedAnnualCIT)
	// Cumulative required = 6M * 3/4 = 4.5M
	assert.Equal(t, 4500000.0, provQ3.CumulativeRequired)
	// 80% of 6M = 4.8M; 4.5M < 4.8M → NOT compliant
	assert.False(t, provQ3.Compliant)
}

// ─── A4: PIT Engine ─────────────────────────────────────────────────────

func TestCalculatePIT_ResidentProgressive(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()
	res, err := svc.CalculatePIT(ctx, "c1", vatPeriod(), []domain.PITEmployeeInput{
		{GrossMonthly: 25000000, Dependants: 1, Months: 1},
	})
	require.NoError(t, err)
	// insurance 10.5% = 2.625M; taxable = 25 - 2.625 - 11 - 4.4 = 6.975M
	// bracket 5-10M: 6.975 * 10% - 0.25M = 447,500
	assert.Equal(t, 1, res.EmployeeCount)
	assert.Equal(t, 447500.0, res.TotalPIT)
	assert.Equal(t, 25000000.0, res.TotalGross)
	assert.Equal(t, 18025000.0, res.TotalDeductions) // (2.625 + 11 + 4.4)M
}

func TestCalculatePIT_NonResidentFlat(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()
	res, err := svc.CalculatePIT(ctx, "c1", vatPeriod(), []domain.PITEmployeeInput{
		{GrossMonthly: 25000000, Months: 1, NonResident: true},
	})
	require.NoError(t, err)
	assert.Equal(t, 5000000.0, res.TotalPIT) // 25M * 20%
	assert.Equal(t, 0.0, res.TotalDeductions)
}

func TestCalculatePIT_NoTaxableIncome(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()
	res, err := svc.CalculatePIT(ctx, "c1", vatPeriod(), []domain.PITEmployeeInput{
		{GrossMonthly: 10000000, Dependants: 1, Months: 1},
	})
	require.NoError(t, err)
	// taxable = 10 - 1.05 - 11 - 4.4 < 0
	assert.Equal(t, 0.0, res.TotalPIT)
}

func TestCalculatePIT_BracketBoundary(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()
	res, err := svc.CalculatePIT(ctx, "c1", vatPeriod(), []domain.PITEmployeeInput{
		{GrossMonthly: 35000000, Dependants: 0, Months: 1},
	})
	require.NoError(t, err)
	// taxable = 35 - 3.675 - 11 = 20.325M → bracket 18-32M: 20.325*20% - 1.65M = 2.415M
	assert.Equal(t, 2415000.0, res.TotalPIT)
}

func TestCalculatePIT_MultipleMonths(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()
	res, err := svc.CalculatePIT(ctx, "c1", vatPeriod(), []domain.PITEmployeeInput{
		{GrossMonthly: 25000000, Dependants: 1, Months: 12},
	})
	require.NoError(t, err)
	assert.Equal(t, 5370000.0, res.TotalPIT) // 12 * 447,500
	assert.Equal(t, 300000000.0, res.TotalGross)
}

func TestCalculatePIT_MultipleEmployees(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()
	res, err := svc.CalculatePIT(ctx, "c1", vatPeriod(), []domain.PITEmployeeInput{
		{EmployeeID: "e1", GrossMonthly: 25000000, Dependants: 1, Months: 1},
		{EmployeeID: "e2", GrossMonthly: 20000000, Months: 1},
	})
	require.NoError(t, err)
	assert.Equal(t, 2, res.EmployeeCount)
	// e1: 447,500; e2: insurance 2.1M, taxable = 20-2.1-11 = 6.9M → 6.9*10%-0.25 = 440,000
	assert.Equal(t, 887500.0, res.TotalPIT)
	assert.Equal(t, 45000000.0, res.TotalGross)
}

// ─── A5: Declaration Engine ─────────────────────────────────────────────

func newTaxTestSvcWithGL() (*taxService, domain.JournalRepository) {
	repo := repository.NewMemoryTaxRepo()
	jeRepo := repository.NewMemoryJournalRepo()
	return NewTaxService(repo, jeRepo, repository.NewMemoryCompanyRepo(), nil, nil).(*taxService), jeRepo
}

func postedEntry(id string, date string, lines ...domain.JournalLine) domain.JournalEntry {
	d, _ := time.Parse("2006-01-02", date)
	return domain.JournalEntry{
		ID: id, EntryNumber: id, EntryDate: d,
		CompanyID: "c1", Status: domain.JournalEntryPosted, Lines: lines,
	}
}

func TestGenerateDeclaration_GTGT01(t *testing.T) {
	svc, jeRepo := newTaxTestSvcWithGL()
	ctx := context.Background()
	je1 := postedEntry("JE1", "2026-01-10",
		vatLine("5111", 0, 10000000), vatLine("33311", 0, 1000000))
	require.NoError(t, jeRepo.Create(ctx, &je1))
	je2 := postedEntry("JE2", "2026-01-15",
		vatLine("152", 5000000, 0), vatLine("1331", 500000, 0))
	require.NoError(t, jeRepo.Create(ctx, &je2))

	decl, err := svc.GenerateDeclaration(ctx, "c1", domain.DeclTypeGTGT01,
		vatPeriod(), "user-1", nil)
	require.NoError(t, err)
	assert.Equal(t, domain.DeclStatusVALIDATED, decl.Status)
	assert.Equal(t, "c1", decl.CompanyID)

	amounts := map[string]float64{}
	for _, l := range decl.Lines {
		amounts[l.LineCode] = l.Amount
	}
	assert.Equal(t, 500000.0, amounts["14"])  // input VAT goods/services
	assert.Equal(t, 0.0, amounts["15"])       // input VAT FA
	assert.Equal(t, 500000.0, amounts["16"])  // 14 + 15
	assert.Equal(t, 1000000.0, amounts["21"]) // output VAT domestic
	assert.Equal(t, 0.0, amounts["22"])
	assert.Equal(t, 1000000.0, amounts["23"]) // 21 + 22
	assert.Equal(t, 500000.0, amounts["30"])  // 23 - 16
	assert.Equal(t, 500000.0, amounts["31"])  // payable
	assert.Equal(t, 0.0, amounts["32"])       // refundable (XOR)
	for _, l := range decl.Lines {
		assert.Equal(t, domain.SrcTypeFROM_LEDGER, l.SourceType)
	}
}

func TestGenerateDeclaration_TNDN03(t *testing.T) {
	svc, jeRepo := newTaxTestSvcWithGL()
	ctx := context.Background()
	je1 := postedEntry("JE1", "2026-01-10", vatLine("5111", 0, 100000000))
	require.NoError(t, jeRepo.Create(ctx, &je1))
	je2 := postedEntry("JE2", "2026-01-15", vatLine("641", 60000000, 0))
	require.NoError(t, jeRepo.Create(ctx, &je2))
	je3 := postedEntry("JE3", "2026-01-20", vatLine("821", 5000000, 0))
	require.NoError(t, jeRepo.Create(ctx, &je3))

	decl, err := svc.GenerateDeclaration(ctx, "c1", domain.DeclTypeTNDN03,
		domain.TaxPeriod{PeriodType: domain.PeriodTypeAnnual, PeriodYear: 2026, PeriodNumber: 1},
		"user-1", nil)
	require.NoError(t, err)
	assert.Equal(t, domain.DeclStatusVALIDATED, decl.Status)

	amounts := map[string]float64{}
	for _, l := range decl.Lines {
		amounts[l.LineCode] = l.Amount
	}
	assert.Equal(t, 100000000.0, amounts["04"])
	assert.Equal(t, 100000000.0, amounts["06"])
	assert.Equal(t, 45000000.0, amounts["12"]) // 100M - 60M + 5M
	assert.Equal(t, 15.0, amounts["13"])       // MICRO
	assert.Equal(t, 6750000.0, amounts["14"])  // 45M * 15%
}

func TestGenerateDeclaration_Duplicate(t *testing.T) {
	svc, jeRepo := newTaxTestSvcWithGL()
	ctx := context.Background()
	je1 := postedEntry("JE1", "2026-01-10", vatLine("5111", 0, 10000000))
	require.NoError(t, jeRepo.Create(ctx, &je1))
	_, err := svc.GenerateDeclaration(ctx, "c1", domain.DeclTypeGTGT01, vatPeriod(), "u1", nil)
	require.NoError(t, err)
	_, err = svc.GenerateDeclaration(ctx, "c1", domain.DeclTypeGTGT01, vatPeriod(), "u1", nil)
	assert.ErrorIs(t, err, domain.ErrDuplicateDeclaration)
}

func TestGenerateDeclaration_UnsupportedType(t *testing.T) {
	svc, _ := newTaxTestSvcWithGL()
	ctx := context.Background()
	_, err := svc.GenerateDeclaration(ctx, "c1", "INVALID", vatPeriod(), "u1", nil)
	assert.ErrorIs(t, err, domain.ErrDeclarationTypeInvalid)
}

func TestGenerateDeclaration_GTGT02(t *testing.T) {
	svc, jeRepo := newTaxTestSvcWithGL()
	ctx := context.Background()
	// Revenue = 50M (5111 credit), Output VAT = 2.5M (33311 credit)
	je := postedEntry("JE1", "2026-06-15",
		vatLine("5111", 0, 50000000), vatLine("33311", 0, 2500000))
	require.NoError(t, jeRepo.Create(ctx, &je))

	period := domain.TaxPeriod{PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 6}
	decl, err := svc.GenerateDeclaration(ctx, "c1", domain.DeclTypeGTGT02, period, "user-1", nil)
	require.NoError(t, err)
	assert.Equal(t, domain.DeclTypeGTGT02, decl.DeclarationType)
	assert.Equal(t, domain.DeclStatusVALIDATED, decl.Status)

	amounts := map[string]float64{}
	for _, l := range decl.Lines {
		amounts[l.LineCode] = l.Amount
	}
	// GTGT02 direct method: [10]=revenue(SalesTotal), [20]=rate(5%), [30]=tax
	assert.Equal(t, 50000000.0, amounts["10"]) // revenue from SalesTotal
	assert.Equal(t, 5.0, amounts["20"])         // rate
	assert.Equal(t, 2500000.0, amounts["30"])   // 50M * 5%
}

func TestGenerateDeclaration_GTGT03(t *testing.T) {
	svc, jeRepo := newTaxTestSvcWithGL()
	ctx := context.Background()
	je1 := postedEntry("JE1", "2026-04-10",
		vatLine("5111", 0, 10000000), vatLine("33311", 0, 1000000))
	require.NoError(t, jeRepo.Create(ctx, &je1))
	je2 := postedEntry("JE2", "2026-05-15",
		vatLine("152", 5000000, 0), vatLine("1331", 500000, 0))
	require.NoError(t, jeRepo.Create(ctx, &je2))

	period := domain.TaxPeriod{PeriodType: domain.PeriodTypeQuarterly, PeriodYear: 2026, PeriodNumber: 2}
	decl, err := svc.GenerateDeclaration(ctx, "c1", domain.DeclTypeGTGT03, period, "user-1", nil)
	require.NoError(t, err)
	assert.Equal(t, domain.DeclTypeGTGT03, decl.DeclarationType)
	assert.Equal(t, domain.DeclStatusVALIDATED, decl.Status)
	assert.Equal(t, domain.PeriodTypeQuarterly, decl.TaxPeriod.PeriodType)

	amounts := map[string]float64{}
	for _, l := range decl.Lines {
		amounts[l.LineCode] = l.Amount
	}
	assert.Equal(t, 500000.0, amounts["14"])  // input VAT goods/services
	assert.Equal(t, 0.0, amounts["15"])       // input VAT FA
	assert.Equal(t, 500000.0, amounts["16"])  // 14 + 15
	assert.Equal(t, 1000000.0, amounts["21"]) // output VAT domestic
	assert.Equal(t, 0.0, amounts["22"])
	assert.Equal(t, 1000000.0, amounts["23"]) // 21 + 22
	assert.Equal(t, 500000.0, amounts["30"])  // 23 - 16
	assert.Equal(t, 500000.0, amounts["31"])  // payable
	assert.Equal(t, 0.0, amounts["32"])       // refundable (XOR)
}

func TestGenerateDeclaration_ZeroDeclaration(t *testing.T) {
	svc, jeRepo := newTaxTestSvcWithGL()
	ctx := context.Background()
	je1 := postedEntry("JE1", "2026-01-10",
		vatLine("5111", 0, 10000000), vatLine("131", 11000000, 0))
	require.NoError(t, jeRepo.Create(ctx, &je1))
	// all entries cancelled → no posted journals
	jePtr, _ := jeRepo.GetByID(ctx, "JE1")
	jePtr.Status = domain.JournalEntryCancelled
	require.NoError(t, jeRepo.Update(ctx, jePtr))

	decl, err := svc.GenerateDeclaration(ctx, "c1", domain.DeclTypeGTGT01, vatPeriod(), "u1", nil)
	require.NoError(t, err)
	for _, l := range decl.Lines {
		assert.Equal(t, 0.0, l.Amount)
	}
}

func TestGenerateDeclaration_SkipsDraftEntries(t *testing.T) {
	svc, jeRepo := newTaxTestSvcWithGL()
	ctx := context.Background()
	je := postedEntry("JE1", "2026-01-10", vatLine("5111", 0, 10000000))
	je.Status = domain.JournalEntryDraft
	require.NoError(t, jeRepo.Create(ctx, &je))
	decl, err := svc.GenerateDeclaration(ctx, "c1", domain.DeclTypeGTGT01, vatPeriod(), "u1", nil)
	require.NoError(t, err)
	for _, l := range decl.Lines {
		assert.Equal(t, 0.0, l.Amount)
	}
}

func TestGenerateDeclaration_IgnoresOtherCompany(t *testing.T) {
	svc, jeRepo := newTaxTestSvcWithGL()
	ctx := context.Background()
	other := postedEntry("JE-OTHER", "2026-01-10", vatLine("5111", 0, 10000000))
	other.CompanyID = "company-999"
	require.NoError(t, jeRepo.Create(ctx, &other))

	decl, err := svc.GenerateDeclaration(ctx, "c1", domain.DeclTypeGTGT01, vatPeriod(), "u1", nil)
	require.NoError(t, err)
	for _, l := range decl.Lines {
		assert.Equal(t, 0.0, l.Amount)
	}
}

// ─── A6: Payment Automation ─────────────────────────────────────────────

func TestAcknowledgeDeclaration_CreatesPayable(t *testing.T) {
	svc, jeRepo := newTaxTestSvcWithGL()
	ctx := context.Background()
	je1 := postedEntry("JE1", "2026-01-10", vatLine("5111", 0, 10000000))
	require.NoError(t, jeRepo.Create(ctx, &je1))

	decl, err := svc.GenerateDeclaration(ctx, "c1", domain.DeclTypeGTGT01, vatPeriod(), "u1", nil)
	require.NoError(t, err)
	require.NoError(t, svc.SubmitDeclaration(ctx, decl.ID, "u1"))
	require.NoError(t, svc.AcknowledgeDeclaration(ctx, decl.ID, "GDT-ACK-1"))

	payments, err := svc.ListPayments(ctx, domain.PaymentFilter{
		CompanyID: "c1", TaxType: domain.TaxTypeVAT,
		PeriodYear: 2026, PeriodNumber: 1,
	})
	require.NoError(t, err)
	require.Len(t, payments, 1)
	assert.Equal(t, decl.ID, payments[0].DeclarationID)
	assert.Equal(t, 1000000.0, payments[0].DeclaredAmount)
	assert.Equal(t, domain.PayStatusPENDING, payments[0].Status)
	assert.Equal(t, "2026-02-20", payments[0].DueDate[:10])
}

func TestAcknowledgeDeclaration_NoPaymentWhenRefundable(t *testing.T) {
	svc, jeRepo := newTaxTestSvcWithGL()
	ctx := context.Background()
	je1 := postedEntry("JE1", "2026-01-10",
		vatLine("152", 20000000, 0), vatLine("1331", 2000000, 0))
	require.NoError(t, jeRepo.Create(ctx, &je1))

	decl, err := svc.GenerateDeclaration(ctx, "c1", domain.DeclTypeGTGT01, vatPeriod(), "u1", nil)
	require.NoError(t, err)
	require.NoError(t, svc.SubmitDeclaration(ctx, decl.ID, "u1"))
	require.NoError(t, svc.AcknowledgeDeclaration(ctx, decl.ID, "GDT-ACK-2"))

	payments, err := svc.ListPayments(ctx, domain.PaymentFilter{CompanyID: "c1"})
	require.NoError(t, err)
	assert.Empty(t, payments)
}

func TestAcknowledgeDeclaration_CITAnnualDueDate(t *testing.T) {
	svc, jeRepo := newTaxTestSvcWithGL()
	ctx := context.Background()
	je1 := postedEntry("JE1", "2026-06-30", vatLine("5111", 0, 100000000))
	require.NoError(t, jeRepo.Create(ctx, &je1))

	decl, err := svc.GenerateDeclaration(ctx, "c1", domain.DeclTypeTNDN03,
		domain.TaxPeriod{PeriodType: domain.PeriodTypeAnnual, PeriodYear: 2026, PeriodNumber: 1}, "u1", nil)
	require.NoError(t, err)
	require.NoError(t, svc.SubmitDeclaration(ctx, decl.ID, "u1"))
	require.NoError(t, svc.AcknowledgeDeclaration(ctx, decl.ID, "GDT-ACK-3"))

	payments, err := svc.ListPayments(ctx, domain.PaymentFilter{
		CompanyID: "c1", TaxType: domain.TaxTypeCIT, PeriodYear: 2026,
	})
	require.NoError(t, err)
	require.Len(t, payments, 1)
	assert.Equal(t, "2027-03-31", payments[0].DueDate[:10])
}

// ─── B3: E-Invoice Issuance Pipeline ────────────────────────────────────

type stubGDT struct {
	submitResp     *domain.GDTSubmitResponse
	statusResp     *domain.GDTStatusResponse
	declSubmitResp *domain.GDTDeclarationSubmitResponse
	declStatusResp *domain.GDTDeclarationStatusResponse
	submitErr      error
	statusErr      error
	cancelErr      error
	submittedXML   string
	cancelled      bool
}

func (g *stubGDT) SubmitInvoice(_ context.Context, xml, certID string) (*domain.GDTSubmitResponse, error) {
	g.submittedXML = xml
	if g.submitErr != nil {
		return nil, g.submitErr
	}
	return g.submitResp, nil
}
func (g *stubGDT) GetInvoiceStatus(_ context.Context, _ string) (*domain.GDTStatusResponse, error) {
	if g.statusErr != nil {
		return nil, g.statusErr
	}
	return g.statusResp, nil
}
func (g *stubGDT) CancelInvoice(_ context.Context, _, _ string) error {
	g.cancelled = true
	return g.cancelErr
}

func (g *stubGDT) SubmitDeclaration(_ context.Context, xml, certID string) (*domain.GDTDeclarationSubmitResponse, error) {
	g.submittedXML = xml
	if g.submitErr != nil {
		return nil, g.submitErr
	}
	if g.declSubmitResp != nil {
		return g.declSubmitResp, nil
	}
	return &domain.GDTDeclarationSubmitResponse{SubmissionID: "SUB-1", Status: "SUBMITTED"}, nil
}

func (g *stubGDT) QueryDeclarationStatus(_ context.Context, _ string) (*domain.GDTDeclarationStatusResponse, error) {
	if g.statusErr != nil {
		return nil, g.statusErr
	}
	if g.declStatusResp != nil {
		return g.declStatusResp, nil
	}
	if g.statusResp == nil {
		return &domain.GDTDeclarationStatusResponse{Status: "ACKNOWLEDGED", AckRef: "ACK-REF-1"}, nil
	}
	return &domain.GDTDeclarationStatusResponse{Status: g.statusResp.Status, AckRef: "ACK-REF-1"}, nil
}

type stubSigner struct{ err error }

func (s *stubSigner) SignTXML(xmlBody, _ string) (string, error) {
	if s.err != nil {
		return "", s.err
	}
	return "signed:" + xmlBody, nil
}

func (s *stubSigner) SignDocument(xmlBody string) (einvoice.SignResult, error) {
	if s.err != nil {
		return einvoice.SignResult{}, s.err
	}
	return einvoice.SignResult{SignatureBase64: "BASE64SIG:" + xmlBody, SignedAt: "2026-04-15T10:00:00+07:00"}, nil
}

func newTaxTestSvcIssuer(g *stubGDT, signer TXMLSigner) (*taxService, domain.TaxRepository, domain.JournalRepository) {
	repo := repository.NewMemoryTaxRepo()
	jeRepo := repository.NewMemoryJournalRepo()
	return NewTaxService(repo, jeRepo, repository.NewMemoryCompanyRepo(), g, signer).(*taxService), repo, jeRepo
}

func testEInvoice(status domain.EInvLifecycleStatus) *domain.EInvoice {
	return &domain.EInvoice{
		ID: "EINV-1", CompanyID: "c1", Pattern: "01GTKT0/001", Serial: "AA/26E",
		InvoiceType: domain.EInvTypeORIGINAL, BuyerName: "Buyer", CurrencyCode: "VND",
		IssueDate: "2026-04-15", Status: status,
		Subtotal: 1000000, VATAmount: 100000, GrandTotal: 1100000,
		Lines: []domain.EInvoiceLine{{LineNumber: 1, Description: "Svc", Quantity: 1, UnitPrice: 1000000, LineTotal: 1000000, VATRate: 10, VATAmount: 100000}},
	}
}

func TestIssueEInvoice_SubmitsToGDT(t *testing.T) {
	g := &stubGDT{submitResp: &domain.GDTSubmitResponse{TransactionID: "TXN-1", Status: "SUBMITTED", GDTRef: "GDT-1"}}
	svc, repo, jeRepo := newTaxTestSvcIssuer(g, &stubSigner{})
	ctx := context.Background()
	inv := testEInvoice(domain.EInvStatusDRAFT)
	require.NoError(t, repo.CreateEInvoice(ctx, inv))

	require.NoError(t, svc.IssueEInvoice(ctx, inv.ID))

	updated, err := repo.GetEInvoiceByID(ctx, inv.ID)
	require.NoError(t, err)
	assert.Equal(t, domain.EInvStatusSUBMITTED, updated.Status)
	assert.Equal(t, "TXN-1", updated.GDTTransactionID)
	assert.NotEmpty(t, updated.XMLBody)
	assert.NotEmpty(t, updated.SignedXML)
	assert.NotEmpty(t, updated.SigningDate)
	// signer received the generated TXML, GDT received the signed XML
	assert.Equal(t, "signed:"+updated.XMLBody, g.submittedXML)

	// Verify journal entry auto-posted
	entries, _ := jeRepo.GetByStatus(ctx, domain.JournalEntryPosted)
	require.Len(t, entries, 1)
	amounts := map[string]float64{}
	for _, l := range entries[0].Lines {
		amounts[l.AccountCode] = l.DebitAmount - l.CreditAmount
	}
	assert.Equal(t, 1100000.0, amounts["131"])   // Dr AR
	assert.Equal(t, -1000000.0, amounts["5111"])  // Cr Revenue
	assert.Equal(t, -100000.0, amounts["33311"])  // Cr VAT
}

func TestIssueEInvoice_WrongStatus(t *testing.T) {
	g := &stubGDT{submitResp: &domain.GDTSubmitResponse{}}
	svc, repo, _ := newTaxTestSvcIssuer(g, &stubSigner{})
	ctx := context.Background()
	inv := testEInvoice(domain.EInvStatusISSUED)
	require.NoError(t, repo.CreateEInvoice(ctx, inv))

	err := svc.IssueEInvoice(ctx, inv.ID)
	assert.ErrorIs(t, err, domain.ErrInvoiceStatusInvalid)
}

func TestIssueEInvoice_NoGDTConfigured(t *testing.T) {
	svc, repo, _ := newTaxTestSvcIssuer(nil, nil)
	ctx := context.Background()
	inv := testEInvoice(domain.EInvStatusDRAFT)
	require.NoError(t, repo.CreateEInvoice(ctx, inv))

	err := svc.IssueEInvoice(ctx, inv.ID)
	assert.ErrorIs(t, err, domain.ErrGDTUnavailable)
}

func TestIssueEInvoice_GDTErrorMapped(t *testing.T) {
	g := &stubGDT{submitErr: domain.ErrGDTUnauthorized}
	svc, repo, _ := newTaxTestSvcIssuer(g, &stubSigner{})
	ctx := context.Background()
	inv := testEInvoice(domain.EInvStatusDRAFT)
	require.NoError(t, repo.CreateEInvoice(ctx, inv))

	err := svc.IssueEInvoice(ctx, inv.ID)
	assert.ErrorIs(t, err, domain.ErrGDTUnauthorized)
}

func TestCheckInvoiceStatus_Acknowledged(t *testing.T) {
	g := &stubGDT{statusResp: &domain.GDTStatusResponse{Status: "ACKNOWLEDGED"}}
	svc, repo, _ := newTaxTestSvcIssuer(g, &stubSigner{})
	ctx := context.Background()
	inv := testEInvoice(domain.EInvStatusSUBMITTED)
	inv.GDTTransactionID = "TXN-1"
	require.NoError(t, repo.CreateEInvoice(ctx, inv))

	require.NoError(t, svc.CheckInvoiceStatus(ctx, inv.ID))
	updated, _ := repo.GetEInvoiceByID(ctx, inv.ID)
	assert.Equal(t, domain.EInvStatusISSUED, updated.Status)
}

func TestCheckInvoiceStatus_Rejected(t *testing.T) {
	g := &stubGDT{statusResp: &domain.GDTStatusResponse{Status: "REJECTED"}}
	svc, repo, _ := newTaxTestSvcIssuer(g, &stubSigner{})
	ctx := context.Background()
	inv := testEInvoice(domain.EInvStatusSUBMITTED)
	inv.GDTTransactionID = "TXN-1"
	require.NoError(t, repo.CreateEInvoice(ctx, inv))

	require.NoError(t, svc.CheckInvoiceStatus(ctx, inv.ID))
	updated, _ := repo.GetEInvoiceByID(ctx, inv.ID)
	assert.Equal(t, domain.EInvStatusVALIDATED, updated.Status)
}

func TestCheckInvoiceStatus_NotSubmitted(t *testing.T) {
	svc, repo, _ := newTaxTestSvcIssuer(nil, nil)
	ctx := context.Background()
	inv := testEInvoice(domain.EInvStatusDRAFT)
	require.NoError(t, repo.CreateEInvoice(ctx, inv))

	err := svc.CheckInvoiceStatus(ctx, inv.ID)
	assert.ErrorIs(t, err, domain.ErrInvoiceStatusInvalid)
}

func TestCancelEInvoice_NotifiesGDT(t *testing.T) {
	g := &stubGDT{}
	svc, repo, _ := newTaxTestSvcIssuer(g, &stubSigner{})
	ctx := context.Background()
	inv := testEInvoice(domain.EInvStatusISSUED)
	inv.GDTTransactionID = "TXN-1"
	require.NoError(t, repo.CreateEInvoice(ctx, inv))

	require.NoError(t, svc.CancelEInvoice(ctx, inv.ID, "buyer request"))
	assert.True(t, g.cancelled)
	updated, _ := repo.GetEInvoiceByID(ctx, inv.ID)
	assert.Equal(t, domain.EInvStatusCANCELLED, updated.Status)
	assert.Equal(t, "buyer request", updated.CancelReason)
}

// spec §2.2: cancellation only from ISSUED — a SUBMITTED invoice has no CQT
// code yet (Circular 78/2021).
func TestCancelEInvoice_NotCancellableBeforeIssued(t *testing.T) {
	g := &stubGDT{}
	svc, repo, _ := newTaxTestSvcIssuer(g, &stubSigner{})
	ctx := context.Background()
	for _, st := range []domain.EInvLifecycleStatus{
		domain.EInvStatusDRAFT, domain.EInvStatusSIGNED, domain.EInvStatusSUBMITTED, domain.EInvStatusVALIDATED,
	} {
		inv := testEInvoice(st)
		require.NoError(t, repo.CreateEInvoice(ctx, inv))
		err := svc.CancelEInvoice(ctx, inv.ID, "reason")
		assert.ErrorIs(t, err, domain.ErrInvoiceStatusInvalid, "status %s", st)
	}
}

func TestCreateAmendmentInvoice_Adjustment(t *testing.T) {
	g := &stubGDT{}
	svc, invRepo, _ := newTaxTestSvcIssuer(g, &stubSigner{})
	ctx := context.Background()
	original := testEInvoice(domain.EInvStatusISSUED)
	require.NoError(t, invRepo.CreateEInvoice(ctx, original))

	adjLines := []domain.EInvoiceLine{{LineNumber: 1, Description: "Item adjusted", LineTotal: 2000000, VATRate: 10, VATAmount: 200000}}
	adj, err := svc.CreateAmendmentInvoice(ctx, original.ID, domain.EInvTypeADJUSTMENT, adjLines, "user-1")
	require.NoError(t, err)

	assert.Equal(t, domain.EInvTypeADJUSTMENT, adj.InvoiceType)
	assert.Equal(t, original.ID, adj.OriginalInvoiceID)
	assert.Equal(t, domain.EInvStatusDRAFT, adj.Status)
	assert.Equal(t, 2000000.0, adj.Subtotal)
	assert.Equal(t, 200000.0, adj.VATAmount)
	assert.Equal(t, 2200000.0, adj.GrandTotal)
	assert.Equal(t, original.CompanyID, adj.CompanyID)
}

func TestCreateAmendmentInvoice_InvalidType(t *testing.T) {
	g := &stubGDT{}
	svc, invRepo, _ := newTaxTestSvcIssuer(g, &stubSigner{})
	ctx := context.Background()
	original := testEInvoice(domain.EInvStatusISSUED)
	require.NoError(t, invRepo.CreateEInvoice(ctx, original))

	_, err := svc.CreateAmendmentInvoice(ctx, original.ID, domain.EInvTypeORIGINAL, nil, "user-1")
	assert.Error(t, err)
}

func TestCreateAmendmentInvoice_NotIssued(t *testing.T) {
	g := &stubGDT{}
	svc, invRepo, _ := newTaxTestSvcIssuer(g, &stubSigner{})
	ctx := context.Background()
	original := testEInvoice(domain.EInvStatusDRAFT)
	require.NoError(t, invRepo.CreateEInvoice(ctx, original))

	_, err := svc.CreateAmendmentInvoice(ctx, original.ID, domain.EInvTypeREPLACEMENT, nil, "user-1")
	assert.Error(t, err)
}

func TestCreateAmendmentInvoice_Replacement(t *testing.T) {
	g := &stubGDT{}
	svc, invRepo, _ := newTaxTestSvcIssuer(g, &stubSigner{})
	ctx := context.Background()
	original := testEInvoice(domain.EInvStatusISSUED)
	require.NoError(t, invRepo.CreateEInvoice(ctx, original))

	replLines := []domain.EInvoiceLine{
		{LineNumber: 1, Description: "Replaced item", Quantity: 2, UnitPrice: 500000, LineTotal: 1000000, VATRate: 10, VATAmount: 100000},
	}
	repl, err := svc.CreateAmendmentInvoice(ctx, original.ID, domain.EInvTypeREPLACEMENT, replLines, "user-1")
	require.NoError(t, err)

	assert.Equal(t, domain.EInvTypeREPLACEMENT, repl.InvoiceType)
	assert.Equal(t, original.ID, repl.OriginalInvoiceID)
	assert.Equal(t, domain.EInvStatusDRAFT, repl.Status)
	assert.Equal(t, 1000000.0, repl.Subtotal)
	assert.Equal(t, 100000.0, repl.VATAmount)
	assert.Equal(t, 1100000.0, repl.GrandTotal)
	assert.Equal(t, original.CompanyID, repl.CompanyID)
	assert.Equal(t, original.Pattern, repl.Pattern)
	assert.Equal(t, original.Serial, repl.Serial)
}

func TestCreateAmendmentInvoice_CancellationNote(t *testing.T) {
	g := &stubGDT{}
	svc, invRepo, _ := newTaxTestSvcIssuer(g, &stubSigner{})
	ctx := context.Background()
	original := testEInvoice(domain.EInvStatusISSUED)
	require.NoError(t, invRepo.CreateEInvoice(ctx, original))

	cancelLines := []domain.EInvoiceLine{
		{LineNumber: 1, Description: "Cancel all", LineTotal: 1000000, VATRate: 10, VATAmount: 100000},
	}
	cn, err := svc.CreateAmendmentInvoice(ctx, original.ID, domain.EInvTypeCANCELLATION_NOTE, cancelLines, "user-1")
	require.NoError(t, err)

	assert.Equal(t, domain.EInvTypeCANCELLATION_NOTE, cn.InvoiceType)
	assert.Equal(t, original.ID, cn.OriginalInvoiceID)
	assert.Equal(t, domain.EInvStatusDRAFT, cn.Status)
	assert.Equal(t, original.CompanyID, cn.CompanyID)
}

func TestPEMSigner_SignsAndVerifies(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	signer := einvoice.NewPEMSigner(key, "CERT-123", func() time.Time {
		return time.Date(2026, 4, 15, 14, 30, 5, 0, time.FixedZone("+07", 7*3600))
	})
	body, err := einvoice.GenerateTXML(testEInvoice(domain.EInvStatusDRAFT))
	require.NoError(t, err)

	signed, err := signer.SignTXML(string(body), "sig-1")
	require.NoError(t, err)
	assert.Contains(t, signed, "<BK:ChuKySo>")
	assert.Contains(t, signed, "<BK:SerialNumber>CERT-123</BK:SerialNumber>")
	assert.Contains(t, signed, "<BK:ThoiDiemKy>2026-04-15T14:30:05+07:00</BK:ThoiDiemKy>")

	// signature covers the canonical body; verify against it
	canon, err := xmldsig.Canonicalize(body)
	require.NoError(t, err)
	m := regexp.MustCompile(`<BK:DuLieuKy>([^<]+)</BK:DuLieuKy>`).FindStringSubmatch(signed)
	require.Len(t, m, 2)
	require.NoError(t, xmldsig.VerifyBase64(&key.PublicKey, canon, m[1]))
}

func TestPEMSigner_MissingKyThuat(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	signer := einvoice.NewPEMSigner(key, "CERT-123", time.Now)
	_, err = signer.SignTXML("<Invoice/>", "sig-1")
	assert.Error(t, err)
}

// ─── C3: Declaration → GDT Submission ──────────────────────────────────

func newDeclTestDecl() *domain.TaxDeclaration {
	return &domain.TaxDeclaration{
		ID: "DECL-1", CompanyID: "c1", DeclarationType: domain.DeclTypeGTGT01,
		TaxPeriod: domain.TaxPeriod{PeriodType: domain.PeriodTypeQuarterly, PeriodYear: 2026, PeriodNumber: 1},
		Status:    domain.DeclStatusVALIDATED, AdjustmentType: domain.AdjTypeNONE, Version: 1,
		CreatedAt: "2026-04-15T10:00:00+07:00",
		Lines:     []domain.TaxDeclarationLine{{LineCode: "[10]", LineName: "Sales", Amount: 100000000}},
	}
}

func newTaxTestSvcDecl(g *stubGDT, signer TXMLSigner, company *domain.Company) (*taxService, domain.TaxRepository, domain.CompanyRepository) {
	repo := repository.NewMemoryTaxRepo()
	crepo := repository.NewMemoryCompanyRepo()
	if company != nil {
		_ = crepo.Create(context.Background(), company)
	}
	return NewTaxService(repo, repository.NewMemoryJournalRepo(), crepo, g, signer).(*taxService), repo, crepo
}

func TestSubmitDeclaration_FullPipeline(t *testing.T) {
	g := &stubGDT{}
	svc, repo, _ := newTaxTestSvcDecl(g, &stubSigner{}, &domain.Company{ID: "c1", TaxCode: "0100123456", LegalNameVN: "CONG TY ABC"})
	ctx := context.Background()
	d := newDeclTestDecl()
	require.NoError(t, repo.CreateDeclaration(ctx, d))

	require.NoError(t, svc.SubmitDeclaration(ctx, d.ID, "user-1"))

	updated, _ := repo.GetDeclarationByID(ctx, d.ID)
	assert.Equal(t, domain.DeclStatusSUBMITTED, updated.Status)
	assert.Equal(t, "SUB-1", updated.GDTSubmissionID)
	assert.NotEmpty(t, updated.DeclarationXML)
	assert.Contains(t, updated.DeclarationXML, "01/GTGT")
	assert.Contains(t, updated.DeclarationXML, "0100123456")
	assert.Len(t, updated.Signatures, 1)
	assert.Equal(t, "user-1", updated.SubmittedBy)
	assert.NotEmpty(t, updated.GDTResponseXML, "GDT response should be stored")
	assert.Contains(t, updated.GDTResponseXML, "SUB-1")
	// GDT received the signed XML (signer got it)
	assert.True(t, strings.HasPrefix(g.submittedXML, "<BK:BoKe"))
}

func TestSubmitDeclaration_CompanyMissing(t *testing.T) {
	g := &stubGDT{}
	svc, repo, _ := newTaxTestSvcDecl(g, &stubSigner{}, nil) // company never created
	ctx := context.Background()
	d := newDeclTestDecl()
	require.NoError(t, repo.CreateDeclaration(ctx, d))

	err := svc.SubmitDeclaration(ctx, d.ID, "user-1")
	assert.ErrorIs(t, err, domain.ErrCompanyNotFound)
}

func TestSubmitDeclaration_GDTDown(t *testing.T) {
	g := &stubGDT{submitErr: domain.ErrGDTUnavailable}
	svc, repo, _ := newTaxTestSvcDecl(g, &stubSigner{}, &domain.Company{ID: "c1", TaxCode: "0100123456"})
	ctx := context.Background()
	d := newDeclTestDecl()
	require.NoError(t, repo.CreateDeclaration(ctx, d))

	err := svc.SubmitDeclaration(ctx, d.ID, "user-1")
	assert.ErrorIs(t, err, domain.ErrGDTUnavailable)
	updated, _ := repo.GetDeclarationByID(ctx, d.ID)
	assert.Equal(t, domain.DeclStatusVALIDATED, updated.Status) // unchanged
}

func TestSubmitDeclaration_LocalOnlyNoGDT(t *testing.T) {
	svc, repo, _ := newTaxTestSvcDecl(nil, nil, nil)
	ctx := context.Background()
	d := newDeclTestDecl()
	require.NoError(t, repo.CreateDeclaration(ctx, d))

	require.NoError(t, svc.SubmitDeclaration(ctx, d.ID, "user-1"))
	updated, _ := repo.GetDeclarationByID(ctx, d.ID)
	assert.Equal(t, domain.DeclStatusSUBMITTED, updated.Status)
	assert.Empty(t, updated.DeclarationXML) // no XML without GDT wiring
}

func TestCheckDeclarationStatus_AcknowledgedCreatesPayment(t *testing.T) {
	g := &stubGDT{statusResp: &domain.GDTStatusResponse{Status: "ACKNOWLEDGED"}}
	svc, repo, _ := newTaxTestSvcDecl(g, &stubSigner{}, &domain.Company{ID: "c1", TaxCode: "0100123456"})
	ctx := context.Background()
	d := newDeclTestDecl()
	d.Status = domain.DeclStatusSUBMITTED
	d.GDTSubmissionID = "SUB-1"
	require.NoError(t, repo.CreateDeclaration(ctx, d))

	require.NoError(t, svc.CheckDeclarationStatus(ctx, d.ID))

	updated, _ := repo.GetDeclarationByID(ctx, d.ID)
	assert.Equal(t, domain.DeclStatusACKNOWLEDGED, updated.Status)
	assert.Equal(t, "ACK-REF-1", updated.AcknowledgementRef) // from stub statusResp? stub maps status only — see stub
}

func TestCheckDeclarationStatus_Rejected(t *testing.T) {
	g := &stubGDT{statusResp: &domain.GDTStatusResponse{Status: "REJECTED"}}
	svc, repo, _ := newTaxTestSvcDecl(g, &stubSigner{}, nil)
	ctx := context.Background()
	d := newDeclTestDecl()
	d.Status = domain.DeclStatusSUBMITTED
	d.GDTSubmissionID = "SUB-1"
	require.NoError(t, repo.CreateDeclaration(ctx, d))

	require.NoError(t, svc.CheckDeclarationStatus(ctx, d.ID))
	updated, _ := repo.GetDeclarationByID(ctx, d.ID)
	assert.Equal(t, domain.DeclStatusREJECTED, updated.Status)
}

func TestCheckDeclarationStatus_NotSubmitted(t *testing.T) {
	svc, repo, _ := newTaxTestSvcDecl(nil, nil, nil)
	ctx := context.Background()
	d := newDeclTestDecl()
	require.NoError(t, repo.CreateDeclaration(ctx, d))

	err := svc.CheckDeclarationStatus(ctx, d.ID)
	assert.ErrorIs(t, err, domain.ErrDeclarationNotEditable)
}

// GDT response code 02 = duplicate submission — the earlier filing stands;
// acknowledge it and auto-create the payable, never a second one.
func TestCheckDeclarationStatus_DuplicateAcknowledges(t *testing.T) {
	g := &stubGDT{declStatusResp: &domain.GDTDeclarationStatusResponse{Code: "02", Status: "REJECTED", AckRef: "ACK-REF-1"}}
	svc, repo, _ := newTaxTestSvcDecl(g, &stubSigner{}, &domain.Company{ID: "c1", TaxCode: "0100123456"})
	ctx := context.Background()
	d := newDeclTestDecl()
	d.Status = domain.DeclStatusSUBMITTED
	d.GDTSubmissionID = "SUB-1"
	require.NoError(t, repo.CreateDeclaration(ctx, d))

	require.NoError(t, svc.CheckDeclarationStatus(ctx, d.ID))

	updated, _ := repo.GetDeclarationByID(ctx, d.ID)
	assert.Equal(t, domain.DeclStatusACKNOWLEDGED, updated.Status)
	assert.Equal(t, "ACK-REF-1", updated.AcknowledgementRef)
}

// GDT response code 10 = period already declared — filing must be an
// amendment (LanDau=2), never treated as a plain rejection.
func TestCheckDeclarationStatus_AlreadyDeclared(t *testing.T) {
	g := &stubGDT{declStatusResp: &domain.GDTDeclarationStatusResponse{Code: "10", Status: "REJECTED"}}
	svc, repo, _ := newTaxTestSvcDecl(g, &stubSigner{}, nil)
	ctx := context.Background()
	d := newDeclTestDecl()
	d.Status = domain.DeclStatusSUBMITTED
	d.GDTSubmissionID = "SUB-1"
	require.NoError(t, repo.CreateDeclaration(ctx, d))

	err := svc.CheckDeclarationStatus(ctx, d.ID)
	assert.ErrorIs(t, err, domain.ErrDeclarationPeriodAlreadyDeclared)
}

// GDT response code 03 = taxpayer code unknown at GDT — profile problem,
// not a declaration defect.
func TestCheckDeclarationStatus_TaxCodeNotFound(t *testing.T) {
	g := &stubGDT{declStatusResp: &domain.GDTDeclarationStatusResponse{Code: "03", Status: "REJECTED"}}
	svc, repo, _ := newTaxTestSvcDecl(g, &stubSigner{}, nil)
	ctx := context.Background()
	d := newDeclTestDecl()
	d.Status = domain.DeclStatusSUBMITTED
	d.GDTSubmissionID = "SUB-1"
	require.NoError(t, repo.CreateDeclaration(ctx, d))

	err := svc.CheckDeclarationStatus(ctx, d.ID)
	assert.ErrorIs(t, err, domain.ErrGDTInvalidTaxCode)
}

// Synchronous rejection on submit must surface the business error and leave
// the declaration un-submitted (no partial SUBMITTED persist).
func TestSubmitDeclaration_RejectCode(t *testing.T) {
	g := &stubGDT{declSubmitResp: &domain.GDTDeclarationSubmitResponse{Code: "10", Status: "REJECTED"}}
	svc, repo, _ := newTaxTestSvcDecl(g, &stubSigner{}, &domain.Company{ID: "c1", TaxCode: "0100123456"})
	ctx := context.Background()
	d := newDeclTestDecl()
	require.NoError(t, repo.CreateDeclaration(ctx, d))

	err := svc.SubmitDeclaration(ctx, d.ID, "user-1")
	assert.ErrorIs(t, err, domain.ErrDeclarationPeriodAlreadyDeclared)

	updated, _ := repo.GetDeclarationByID(ctx, d.ID)
	assert.Equal(t, domain.DeclStatusVALIDATED, updated.Status)
}

// ─── VAT Reconciliation (BR-VAT-06) ────────────────────────────────────

func TestReconcileVAT_Matched(t *testing.T) {
	svc, repo, _ := newTaxTestSvcIssuer(nil, nil)
	ctx := context.Background()
	period := domain.TaxPeriod{PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 4}
	// declaration with output VAT [23] = 120000
	decl := &domain.TaxDeclaration{
		ID: "DECL-1", CompanyID: "c1", DeclarationType: domain.DeclTypeGTGT01,
		TaxPeriod: period, Status: domain.DeclStatusACKNOWLEDGED, AdjustmentType: domain.AdjTypeNONE,
		Lines: []domain.TaxDeclarationLine{
			{LineCode: "23", LineName: "Tong VAT dau ra", Amount: 120000},
		},
	}
	require.NoError(t, repo.CreateDeclaration(ctx, decl))
	// two issued invoices, VAT 70000 + 50000
	for i, vat := range []float64{70000, 50000} {
		inv := testEInvoice(domain.EInvStatusISSUED)
		inv.ID = fmt.Sprintf("EINV-%d", i+1)
		inv.IssueDate = "2026-04-15"
		inv.VATAmount = vat
		inv.Lines = []domain.EInvoiceLine{{LineNumber: 1, Description: "S", Quantity: 1, UnitPrice: vat * 10, LineTotal: vat * 10, VATRate: 10, VATAmount: vat}}
		require.NoError(t, repo.CreateEInvoice(ctx, inv))
	}

	res, err := svc.ReconcileVAT(ctx, "c1", period)
	require.NoError(t, err)
	assert.Equal(t, "DECL-1", res.DeclarationID)
	assert.Equal(t, 120000.0, res.InvoiceTotal)
	assert.Equal(t, 120000.0, res.DeclarationTotal)
	assert.Equal(t, 0.0, res.Variance)
	assert.True(t, res.Matched)
	assert.Equal(t, 2, res.InvoiceCount)
	require.Len(t, res.ByRate, 1)
	assert.Equal(t, 120000.0, res.ByRate[0].InvoiceVAT)
}

func TestReconcileVAT_Mismatch(t *testing.T) {
	svc, repo, _ := newTaxTestSvcIssuer(nil, nil)
	ctx := context.Background()
	period := domain.TaxPeriod{PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 4}
	decl := &domain.TaxDeclaration{
		ID: "DECL-1", CompanyID: "c1", DeclarationType: domain.DeclTypeGTGT01,
		TaxPeriod: period, Status: domain.DeclStatusACKNOWLEDGED, AdjustmentType: domain.AdjTypeNONE,
		Lines: []domain.TaxDeclarationLine{{LineCode: "23", LineName: "Tong", Amount: 50000}},
	}
	require.NoError(t, repo.CreateDeclaration(ctx, decl))
	inv := testEInvoice(domain.EInvStatusISSUED)
	inv.IssueDate = "2026-04-15"
	inv.VATAmount = 70000
	inv.Lines = []domain.EInvoiceLine{{LineNumber: 1, Quantity: 1, UnitPrice: 700000, LineTotal: 700000, VATRate: 10, VATAmount: 70000}}
	require.NoError(t, repo.CreateEInvoice(ctx, inv))

	res, err := svc.ReconcileVAT(ctx, "c1", period)
	require.NoError(t, err)
	assert.False(t, res.Matched)
	assert.Equal(t, 20000.0, res.Variance)
}

func TestReconcileVAT_ExcludesCancelled(t *testing.T) {
	svc, repo, _ := newTaxTestSvcIssuer(nil, nil)
	ctx := context.Background()
	period := domain.TaxPeriod{PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 4}
	cancelled := testEInvoice(domain.EInvStatusCANCELLED)
	cancelled.IssueDate = "2026-04-15"
	cancelled.VATAmount = 99999
	require.NoError(t, repo.CreateEInvoice(ctx, cancelled))

	res, err := svc.ReconcileVAT(ctx, "c1", period)
	require.NoError(t, err)
	assert.Equal(t, 0.0, res.InvoiceTotal)
	assert.Equal(t, 0, res.InvoiceCount)
	assert.Len(t, res.ExcludedInvoices, 1)
}

func TestReconcileVAT_PeriodFilter(t *testing.T) {
	svc, repo, _ := newTaxTestSvcIssuer(nil, nil)
	ctx := context.Background()
	period := domain.TaxPeriod{PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 4}
	outOfPeriod := testEInvoice(domain.EInvStatusISSUED)
	outOfPeriod.IssueDate = "2026-05-20"
	outOfPeriod.VATAmount = 99999
	require.NoError(t, repo.CreateEInvoice(ctx, outOfPeriod))

	res, err := svc.ReconcileVAT(ctx, "c1", period)
	require.NoError(t, err)
	assert.Equal(t, 0.0, res.InvoiceTotal)
	assert.Equal(t, 0, res.InvoiceCount)
}

func rateBucket(res *domain.VATReconciliationResult, rate float64) *domain.VATRateReconciliation {
	for i := range res.ByRate {
		if res.ByRate[i].Rate == rate {
			return &res.ByRate[i]
		}
	}
	return nil
}

func TestReconcileVAT_PerRateDistinctCounts(t *testing.T) {
	svc, repo, _ := newTaxTestSvcIssuer(nil, nil)
	ctx := context.Background()
	period := domain.TaxPeriod{PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 4}
	decl := &domain.TaxDeclaration{
		ID: "DECL-1", CompanyID: "c1", DeclarationType: domain.DeclTypeGTGT01,
		TaxPeriod: period, Status: domain.DeclStatusACKNOWLEDGED, AdjustmentType: domain.AdjTypeNONE,
		Lines: []domain.TaxDeclarationLine{{LineCode: "23", LineName: "Tong", Amount: 220000}},
	}
	require.NoError(t, repo.CreateDeclaration(ctx, decl))
	// inv A: lines at 10% + 5%; inv B: line at 10% only
	invA := testEInvoice(domain.EInvStatusISSUED)
	invA.ID = "EINV-A"
	invA.VATAmount = 120000
	invA.Lines = []domain.EInvoiceLine{
		{LineNumber: 1, Quantity: 1, UnitPrice: 700000, LineTotal: 700000, VATRate: 10, VATAmount: 70000},
		{LineNumber: 2, Quantity: 1, UnitPrice: 1000000, LineTotal: 1000000, VATRate: 5, VATAmount: 50000},
	}
	invB := testEInvoice(domain.EInvStatusISSUED)
	invB.ID = "EINV-B"
	invB.VATAmount = 100000
	invB.Lines = []domain.EInvoiceLine{
		{LineNumber: 1, Quantity: 1, UnitPrice: 1000000, LineTotal: 1000000, VATRate: 10, VATAmount: 100000},
	}
	require.NoError(t, repo.CreateEInvoice(ctx, invA))
	require.NoError(t, repo.CreateEInvoice(ctx, invB))

	res, err := svc.ReconcileVAT(ctx, "c1", period)
	require.NoError(t, err)
	assert.Equal(t, 2, res.InvoiceCount)
	r10 := rateBucket(res, 10)
	require.NotNil(t, r10)
	assert.Equal(t, 2, r10.InvoiceCount)
	assert.Equal(t, 170000.0, r10.InvoiceVAT)
	r5 := rateBucket(res, 5)
	require.NotNil(t, r5)
	assert.Equal(t, 1, r5.InvoiceCount)
	assert.Equal(t, 50000.0, r5.InvoiceVAT)
}

func TestReconcileVAT_PrefersLine22(t *testing.T) {
	svc, repo, _ := newTaxTestSvcIssuer(nil, nil)
	ctx := context.Background()
	period := domain.TaxPeriod{PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 4}
	// [22] reported, [23] = [21]+[22] includes VAT from non-e-invoice sources
	decl := &domain.TaxDeclaration{
		ID: "DECL-1", CompanyID: "c1", DeclarationType: domain.DeclTypeGTGT01,
		TaxPeriod: period, Status: domain.DeclStatusACKNOWLEDGED, AdjustmentType: domain.AdjTypeNONE,
		Lines: []domain.TaxDeclarationLine{
			{LineCode: "22", LineName: "VAT dau ra (e-invoices)", Amount: 120000},
			{LineCode: "23", LineName: "Tong VAT dau ra", Amount: 140000},
		},
	}
	require.NoError(t, repo.CreateDeclaration(ctx, decl))
	for i, vat := range []float64{70000, 50000} {
		inv := testEInvoice(domain.EInvStatusISSUED)
		inv.ID = fmt.Sprintf("EINV-%d", i+1)
		inv.VATAmount = vat
		inv.Lines = []domain.EInvoiceLine{{LineNumber: 1, Quantity: 1, UnitPrice: vat * 10, LineTotal: vat * 10, VATRate: 10, VATAmount: vat}}
		require.NoError(t, repo.CreateEInvoice(ctx, inv))
	}

	res, err := svc.ReconcileVAT(ctx, "c1", period)
	require.NoError(t, err)
	assert.Equal(t, 120000.0, res.DeclarationTotal)
	assert.True(t, res.Matched)
	assert.Equal(t, 0.0, res.Variance)
}

func TestReconcileVAT_FallbackLine23Notes(t *testing.T) {
	svc, repo, _ := newTaxTestSvcIssuer(nil, nil)
	ctx := context.Background()
	period := domain.TaxPeriod{PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 4}
	decl := &domain.TaxDeclaration{
		ID: "DECL-1", CompanyID: "c1", DeclarationType: domain.DeclTypeGTGT01,
		TaxPeriod: period, Status: domain.DeclStatusACKNOWLEDGED, AdjustmentType: domain.AdjTypeNONE,
		Lines: []domain.TaxDeclarationLine{{LineCode: "23", LineName: "Tong", Amount: 70000}},
	}
	require.NoError(t, repo.CreateDeclaration(ctx, decl))
	inv := testEInvoice(domain.EInvStatusISSUED)
	inv.VATAmount = 70000
	inv.Lines = []domain.EInvoiceLine{{LineNumber: 1, Quantity: 1, UnitPrice: 700000, LineTotal: 700000, VATRate: 10, VATAmount: 70000}}
	require.NoError(t, repo.CreateEInvoice(ctx, inv))

	res, err := svc.ReconcileVAT(ctx, "c1", period)
	require.NoError(t, err)
	assert.True(t, res.Matched)
	assert.Equal(t, 70000.0, res.DeclarationTotal)
	assert.Contains(t, strings.Join(res.Notes, " "), "[22]")
}

// ─── A6: Payment Journal Entry ──────────────────────────────────────────

func TestRecordPayment_CreatesJournalEntry_VAT(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()
	p := &domain.TaxPayment{
		ID: "TP-VAT-1", CompanyID: "c1", TaxType: domain.TaxTypeVAT,
		PeriodYear: 2026, PeriodNumber: 6, DeclaredAmount: 5000000,
		DueDate: "2026-07-30", Status: domain.PayStatusPENDING,
	}
	require.NoError(t, svc.repo.CreatePayment(ctx, p))

	err := svc.RecordPayment(ctx, "TP-VAT-1", 5000000, "2026-07-25", "EFT-001")
	require.NoError(t, err)

	pay, err := svc.repo.GetPaymentByID(ctx, "TP-VAT-1")
	require.NoError(t, err)
	assert.Equal(t, domain.PayStatusPAID, pay.Status)
	assert.Equal(t, 5000000.0, pay.PaidAmount)
	assert.NotEmpty(t, pay.GLJournalID)

	je, err := svc.jeRepo.GetByID(ctx, pay.GLJournalID)
	require.NoError(t, err)
	assert.Equal(t, "Tax payment TP-VAT-1", je.Description)
	assert.Equal(t, domain.VoucherTypePayment, je.VoucherType)
	assert.Len(t, je.Lines, 2)
	// Dr 33311 (VAT payable) / Cr 112 (Bank)
	assert.Equal(t, "33311", je.Lines[0].AccountCode)
	assert.Equal(t, 5000000.0, je.Lines[0].DebitAmount)
	assert.Equal(t, "112", je.Lines[1].AccountCode)
	assert.Equal(t, 5000000.0, je.Lines[1].CreditAmount)
}

func TestRecordPayment_CreatesJournalEntry_CIT(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()
	p := &domain.TaxPayment{
		ID: "TP-CIT-1", CompanyID: "c1", TaxType: domain.TaxTypeCIT,
		PeriodYear: 2026, PeriodNumber: 3, DeclaredAmount: 10000000,
		DueDate: "2026-04-30", Status: domain.PayStatusPENDING,
	}
	require.NoError(t, svc.repo.CreatePayment(ctx, p))

	err := svc.RecordPayment(ctx, "TP-CIT-1", 10000000, "2026-04-28", "BANK-002")
	require.NoError(t, err)

	pay, err := svc.repo.GetPaymentByID(ctx, "TP-CIT-1")
	require.NoError(t, err)
	assert.NotEmpty(t, pay.GLJournalID)

	je, err := svc.jeRepo.GetByID(ctx, pay.GLJournalID)
	require.NoError(t, err)
	assert.Equal(t, "3334", je.Lines[0].AccountCode)
	assert.Equal(t, "112", je.Lines[1].AccountCode)
}

func TestRecordPayment_CreatesJournalEntry_PartialPayment(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()
	p := &domain.TaxPayment{
		ID: "TP-VAT-2", CompanyID: "c1", TaxType: domain.TaxTypeVAT,
		PeriodYear: 2026, PeriodNumber: 6, DeclaredAmount: 5000000,
		DueDate: "2026-07-30", Status: domain.PayStatusPENDING,
	}
	require.NoError(t, svc.repo.CreatePayment(ctx, p))

	err := svc.RecordPayment(ctx, "TP-VAT-2", 3000000, "2026-07-25", "EFT-002")
	require.NoError(t, err)

	pay, err := svc.repo.GetPaymentByID(ctx, "TP-VAT-2")
	require.NoError(t, err)
	assert.Equal(t, domain.PayStatusPARTIAL, pay.Status)
	assert.NotEmpty(t, pay.GLJournalID)

	je, err := svc.jeRepo.GetByID(ctx, pay.GLJournalID)
	require.NoError(t, err)
	assert.Equal(t, 3000000.0, je.Lines[0].DebitAmount)
	assert.Equal(t, 3000000.0, je.Lines[1].CreditAmount)
}

func TestRecordPayment_NoJournalForZeroAmount(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()
	p := &domain.TaxPayment{
		ID: "TP-ZERO2", CompanyID: "c1", TaxType: domain.TaxTypeVAT,
		PeriodYear: 2026, PeriodNumber: 6, DeclaredAmount: 100,
		DueDate: "2026-07-30", Status: domain.PayStatusPENDING,
	}
	require.NoError(t, svc.repo.CreatePayment(ctx, p))
	err := svc.RecordPayment(ctx, "TP-ZERO2", 0, "2026-07-25", "")
	require.NoError(t, err)
	pay, _ := svc.repo.GetPaymentByID(ctx, "TP-ZERO2")
	assert.Equal(t, domain.PayStatusPARTIAL, pay.Status)
	assert.Empty(t, pay.GLJournalID, "no journal for zero amount")
}

func TestRecordPayment_AutoLateDays(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()
	p := &domain.TaxPayment{
		ID: "TP-LATE-1", CompanyID: "c1", TaxType: domain.TaxTypeVAT,
		PeriodYear: 2026, PeriodNumber: 6, DeclaredAmount: 5000000,
		DueDate: "2026-07-20", Status: domain.PayStatusPENDING,
	}
	require.NoError(t, svc.repo.CreatePayment(ctx, p))

	// Pay 5 days late, partial payment
	err := svc.RecordPayment(ctx, "TP-LATE-1", 3000000, "2026-07-25", "EFT-LATE")
	require.NoError(t, err)

	pay, err := svc.repo.GetPaymentByID(ctx, "TP-LATE-1")
	require.NoError(t, err)
	assert.Equal(t, 5, pay.LateDays, "auto-calculated late days from due date")
	underpaid := 5000000.0 - 3000000.0
	expectedInterest := underpaid * 0.0003 * 5.0
	assert.InDelta(t, expectedInterest, pay.LateInterest, 0.01, "late interest = underpaid * 0.03% * late days")
}

func TestRecordPayment_NoLateDays_WhenOnTime(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()
	p := &domain.TaxPayment{
		ID: "TP-ONTIME", CompanyID: "c1", TaxType: domain.TaxTypeVAT,
		PeriodYear: 2026, PeriodNumber: 6, DeclaredAmount: 5000000,
		DueDate: "2026-07-20", Status: domain.PayStatusPENDING,
	}
	require.NoError(t, svc.repo.CreatePayment(ctx, p))

	// Pay on time
	err := svc.RecordPayment(ctx, "TP-ONTIME", 5000000, "2026-07-20", "EFT-ONTIME")
	require.NoError(t, err)

	pay, err := svc.repo.GetPaymentByID(ctx, "TP-ONTIME")
	require.NoError(t, err)
	assert.Equal(t, 0, pay.LateDays)
	assert.Equal(t, 0.0, pay.LateInterest)
}

func TestRecordPayment_LateInterestJournalEntry(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()
	p := &domain.TaxPayment{
		ID: "TP-LATEINT", CompanyID: "c1", TaxType: domain.TaxTypeVAT,
		PeriodYear: 2026, PeriodNumber: 6, DeclaredAmount: 5000000,
		DueDate: "2026-07-20", Status: domain.PayStatusPENDING,
	}
	require.NoError(t, svc.repo.CreatePayment(ctx, p))

	// Pay 10 days late, partial payment
	err := svc.RecordPayment(ctx, "TP-LATEINT", 3000000, "2026-07-30", "EFT-LATEINT")
	require.NoError(t, err)

	pay, err := svc.repo.GetPaymentByID(ctx, "TP-LATEINT")
	require.NoError(t, err)
	assert.Equal(t, 10, pay.LateDays)
	assert.True(t, pay.LateInterest > 0, "late interest should be positive")

	// Should have 2 journal entries: tax payment + late interest
	entries, _ := svc.jeRepo.GetByVoucherType(ctx, domain.VoucherTypePayment)
	taxEntries := 0
	lateEntries := 0
	for _, e := range entries {
		if e.CompanyID != "c1" {
			continue
		}
		if strings.Contains(e.Description, "Tax payment") {
			taxEntries++
			// Tax payment: Dr 33311 / Cr 112
			assert.Equal(t, "33311", e.Lines[0].AccountCode)
			assert.Equal(t, "112", e.Lines[1].AccountCode)
			assert.Equal(t, 3000000.0, e.Lines[0].DebitAmount)
		}
		if strings.Contains(e.Description, "Late interest") {
			lateEntries++
			// Late interest: Dr 6275 / Cr 112
			assert.Equal(t, "6275", e.Lines[0].AccountCode)
			assert.Equal(t, "112", e.Lines[1].AccountCode)
			assert.Equal(t, pay.LateInterest, e.Lines[0].DebitAmount)
		}
	}
	assert.Equal(t, 1, taxEntries, "one tax payment entry")
	assert.Equal(t, 1, lateEntries, "one late interest entry")
}

func TestRecordPayment_OverpaymentStatus(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()
	p := &domain.TaxPayment{
		ID: "TP-OVER", CompanyID: "c1", TaxType: domain.TaxTypeCIT,
		PeriodYear: 2026, PeriodNumber: 3, DeclaredAmount: 10000000,
		DueDate: "2026-04-30", Status: domain.PayStatusPENDING,
	}
	require.NoError(t, svc.repo.CreatePayment(ctx, p))

	// Pay more than declared
	err := svc.RecordPayment(ctx, "TP-OVER", 12000000, "2026-04-28", "BANK-OVER")
	require.NoError(t, err)

	pay, err := svc.repo.GetPaymentByID(ctx, "TP-OVER")
	require.NoError(t, err)
	assert.Equal(t, domain.PayStatusOVERPAID, pay.Status)
	assert.Equal(t, 12000000.0, pay.PaidAmount)
	assert.Equal(t, 0, pay.LateDays, "no late days for overpayment")
	assert.Equal(t, 0.0, pay.LateInterest)
}

func TestRecordPayment_UpdatesCalendarStatus(t *testing.T) {
	svc, repo := newTaxTestSvc()
	ctx := context.Background()

	// Create calendar entry
	cal := &domain.TaxCalendar{
		ID: "CAL-1", CompanyID: "c1", TaxType: domain.TaxTypeVAT,
		PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 6,
		DeclarationDue: "2026-07-20", PaymentDue: "2026-07-20",
		Status: domain.CalStatusPENDING,
	}
	require.NoError(t, repo.CreateCalendarEntry(ctx, cal))

	// Create payment linked to same period
	p := &domain.TaxPayment{
		ID: "TP-CAL-1", CompanyID: "c1", TaxType: domain.TaxTypeVAT,
		PeriodYear: 2026, PeriodNumber: 6, DeclaredAmount: 5000000,
		DueDate: "2026-07-20", Status: domain.PayStatusPENDING,
	}
	require.NoError(t, repo.CreatePayment(ctx, p))

	// Record full payment
	err := svc.RecordPayment(ctx, "TP-CAL-1", 5000000, "2026-07-18", "EFT-CAL")
	require.NoError(t, err)

	// Calendar should be updated to PAID
	entries, err := repo.GetCalendarByPeriod(ctx, "c1", 2026, 6)
	require.NoError(t, err)
	require.Len(t, entries, 1)
	assert.Equal(t, domain.CalStatusPAID, entries[0].Status, "calendar status updated to PAID")
}

func TestRecordPayment_UpdatesCalendarStatus_Partial(t *testing.T) {
	svc, repo := newTaxTestSvc()
	ctx := context.Background()

	cal := &domain.TaxCalendar{
		ID: "CAL-2", CompanyID: "c1", TaxType: domain.TaxTypeVAT,
		PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 6,
		DeclarationDue: "2026-07-20", PaymentDue: "2026-07-20",
		Status: domain.CalStatusPENDING,
	}
	require.NoError(t, repo.CreateCalendarEntry(ctx, cal))

	p := &domain.TaxPayment{
		ID: "TP-CAL-2", CompanyID: "c1", TaxType: domain.TaxTypeVAT,
		PeriodYear: 2026, PeriodNumber: 6, DeclaredAmount: 5000000,
		DueDate: "2026-07-20", Status: domain.PayStatusPENDING,
	}
	require.NoError(t, repo.CreatePayment(ctx, p))

	// Partial payment — calendar stays PENDING
	err := svc.RecordPayment(ctx, "TP-CAL-2", 3000000, "2026-07-18", "EFT-CAL2")
	require.NoError(t, err)

	entries, err := repo.GetCalendarByPeriod(ctx, "c1", 2026, 6)
	require.NoError(t, err)
	require.Len(t, entries, 1)
	assert.Equal(t, domain.CalStatusPENDING, entries[0].Status, "partial payment does not mark calendar PAID")
}

func TestScanOverdueCalendars(t *testing.T) {
	svc, repo := newTaxTestSvc()
	ctx := context.Background()

	// PENDING calendar past due
	cal1 := &domain.TaxCalendar{
		ID: "CAL-OV1", CompanyID: "c1", TaxType: domain.TaxTypeVAT,
		PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 1,
		DeclarationDue: "2026-02-20", PaymentDue: "2026-02-20",
		Status: domain.CalStatusPENDING,
	}
	require.NoError(t, repo.CreateCalendarEntry(ctx, cal1))

	// PENDING calendar still due in future
	cal2 := &domain.TaxCalendar{
		ID: "CAL-OV2", CompanyID: "c1", TaxType: domain.TaxTypeCIT,
		PeriodType: domain.PeriodTypeQuarterly, PeriodYear: 2026, PeriodNumber: 4,
		DeclarationDue: "2026-12-30", PaymentDue: "2026-12-30",
		Status: domain.CalStatusPENDING,
	}
	require.NoError(t, repo.CreateCalendarEntry(ctx, cal2))

	// Already PAID — should not be touched
	cal3 := &domain.TaxCalendar{
		ID: "CAL-OV3", CompanyID: "c1", TaxType: domain.TaxTypeVAT,
		PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 2,
		DeclarationDue: "2026-03-20", PaymentDue: "2026-03-20",
		Status: domain.CalStatusPAID,
	}
	require.NoError(t, repo.CreateCalendarEntry(ctx, cal3))

	// Scan with "today" = 2026-08-05 (simulated via svc.now override)
	svc.now = func() time.Time { return time.Date(2026, 8, 5, 0, 0, 0, 0, time.UTC) }
	count, err := svc.ScanOverdueCalendars(ctx, "c1")
	require.NoError(t, err)
	assert.Equal(t, 1, count, "only cal1 should be marked overdue")

	entries, _ := repo.GetCalendarByCompany(ctx, "c1")
	for _, e := range entries {
		switch e.ID {
		case "CAL-OV1":
			assert.Equal(t, domain.CalStatusOVERDUE, e.Status)
		case "CAL-OV2":
			assert.Equal(t, domain.CalStatusPENDING, e.Status, "future deadline stays PENDING")
		case "CAL-OV3":
			assert.Equal(t, domain.CalStatusPAID, e.Status, "PAID calendar untouched")
		}
	}
}

func TestGetPaymentSummary(t *testing.T) {
	svc, repo := newTaxTestSvc()
	ctx := context.Background()

	// Create payments with different statuses
	payments := []domain.TaxPayment{
		{ID: "TP-S1", CompanyID: "c1", TaxType: domain.TaxTypeVAT, PeriodYear: 2026, PeriodNumber: 1, DeclaredAmount: 5000000, PaidAmount: 5000000, DueDate: "2026-02-20", Status: domain.PayStatusPAID},
		{ID: "TP-S2", CompanyID: "c1", TaxType: domain.TaxTypeVAT, PeriodYear: 2026, PeriodNumber: 2, DeclaredAmount: 3000000, PaidAmount: 1000000, DueDate: "2026-03-20", Status: domain.PayStatusPARTIAL},
		{ID: "TP-S3", CompanyID: "c1", TaxType: domain.TaxTypeCIT, PeriodYear: 2026, PeriodNumber: 1, DeclaredAmount: 10000000, PaidAmount: 0, DueDate: "2026-04-30", Status: domain.PayStatusPENDING},
		{ID: "TP-S4", CompanyID: "c2", TaxType: domain.TaxTypeVAT, PeriodYear: 2026, PeriodNumber: 1, DeclaredAmount: 2000000, PaidAmount: 2000000, DueDate: "2026-02-20", Status: domain.PayStatusPAID},
	}
	for i := range payments {
		require.NoError(t, repo.CreatePayment(ctx, &payments[i]))
	}

	summary, err := svc.GetPaymentSummary(ctx, "c1")
	require.NoError(t, err)
	assert.Equal(t, "c1", summary.CompanyID)
	assert.Equal(t, 3, summary.TotalPayments)
	assert.Equal(t, 18000000.0, summary.TotalDeclared)
	assert.Equal(t, 6000000.0, summary.TotalPaid)
	assert.Equal(t, 12000000.0, summary.TotalOutstanding)
	assert.Equal(t, 1, summary.ByStatus["PAID"])
	assert.Equal(t, 1, summary.ByStatus["PARTIAL"])
	assert.Equal(t, 1, summary.ByStatus["PENDING"])
}

func TestReconcileCIT(t *testing.T) {
	svc, repo := newTaxTestSvc()
	ctx := context.Background()
	period := domain.TaxPeriod{PeriodType: domain.PeriodTypeQuarterly, PeriodYear: 2026, PeriodNumber: 1}

	decl := &domain.TaxDeclaration{
		ID: "DECL-CIT-1", CompanyID: "c1", DeclarationType: domain.DeclTypeTNDN03,
		TaxPeriod: period, Status: domain.DeclStatusACKNOWLEDGED,
		Lines: []domain.TaxDeclarationLine{
			{LineCode: "10", LineName: "Taxable income", Amount: 50000000},
			{LineCode: "14", LineName: "CIT payable", Amount: 10000000},
		},
	}
	require.NoError(t, repo.CreateDeclaration(ctx, decl))

	res, err := svc.ReconcileCIT(ctx, "c1", period)
	require.NoError(t, err)
	assert.Equal(t, "DECL-CIT-1", res.DeclarationID)
	assert.Equal(t, 50000000.0, res.DeclaredTaxable)
	assert.Equal(t, 10000000.0, res.DeclaredTax)
	assert.Equal(t, 20.0, res.TaxRate)
	assert.Equal(t, 10000000.0, res.CalculatedTax)
	assert.Equal(t, 0.0, res.Variance)
	assert.True(t, res.Matched)
}

func TestReconcileCIT_Variance(t *testing.T) {
	svc, repo := newTaxTestSvc()
	ctx := context.Background()
	period := domain.TaxPeriod{PeriodType: domain.PeriodTypeQuarterly, PeriodYear: 2026, PeriodNumber: 1}

	decl := &domain.TaxDeclaration{
		ID: "DECL-CIT-2", CompanyID: "c1", DeclarationType: domain.DeclTypeTNDN03,
		TaxPeriod: period, Status: domain.DeclStatusACKNOWLEDGED,
		Lines: []domain.TaxDeclarationLine{
			{LineCode: "10", LineName: "Taxable income", Amount: 50000000},
			{LineCode: "14", LineName: "CIT payable", Amount: 9500000},
		},
	}
	require.NoError(t, repo.CreateDeclaration(ctx, decl))

	res, err := svc.ReconcileCIT(ctx, "c1", period)
	require.NoError(t, err)
	assert.False(t, res.Matched)
	assert.Equal(t, -500000.0, res.Variance, "declared less than calculated")
}

func TestCalculatePenalty_LateSubmission(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()

	tests := []struct {
		name     string
		daysLate int
		min      float64
		max      float64
	}{
		{"15 days late", 15, 2000000, 5000000},
		{"45 days late", 45, 5000000, 8000000},
		{"75 days late", 75, 8000000, 15000000},
		{"100 days late", 100, 15000000, 25000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			due := time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC)
		_today := due.AddDate(0, 0, tt.daysLate)
			svc.now = func() time.Time { return _today }

			res, err := svc.CalculatePenalty(ctx, due.Format("2006-01-02"))
			require.NoError(t, err)
			assert.Equal(t, tt.daysLate, res.DaysLate)
			assert.Equal(t, tt.min, res.PenaltyMin)
			assert.Equal(t, tt.max, res.PenaltyMax)
		})
	}
}

func TestCalculatePenalty_OnTime(t *testing.T) {
	svc, _ := newTaxTestSvc()
	ctx := context.Background()
	svc.now = func() time.Time { return time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC) }

	res, err := svc.CalculatePenalty(ctx, "2026-03-20")
	require.NoError(t, err)
	assert.Equal(t, 0, res.DaysLate)
	assert.Equal(t, 0.0, res.PenaltyMin)
	assert.Equal(t, 0.0, res.PenaltyMax)
}

func TestGenerateCalendarBatch(t *testing.T) {
	svc, repo := newTaxTestSvc()
	ctx := context.Background()

	count, err := svc.GenerateCalendarBatch(ctx, "c1", 2026, []domain.TaxType{domain.TaxTypeVAT, domain.TaxTypeCIT})
	require.NoError(t, err)
	// VAT: 12 monthly + CIT: 4 quarterly = 16 entries
	assert.Equal(t, 16, count)

	entries, err := repo.GetCalendarByCompany(ctx, "c1")
	require.NoError(t, err)
	assert.Len(t, entries, 16)

	// Verify VAT monthly entries exist
	vatCount := 0
	for _, e := range entries {
		if e.TaxType == domain.TaxTypeVAT {
			vatCount++
			assert.Equal(t, domain.PeriodTypeMonthly, e.PeriodType)
			assert.Equal(t, 2026, e.PeriodYear)
			assert.Equal(t, domain.CalStatusPENDING, e.Status)
		}
	}
	assert.Equal(t, 12, vatCount, "12 VAT monthly entries")

	// Verify CIT quarterly entries exist
	citCount := 0
	for _, e := range entries {
		if e.TaxType == domain.TaxTypeCIT {
			citCount++
			assert.Equal(t, domain.PeriodTypeQuarterly, e.PeriodType)
		}
	}
	assert.Equal(t, 4, citCount, "4 CIT quarterly entries")
}

func TestGenerateCalendarBatch_SkipsDuplicates(t *testing.T) {
	svc, repo := newTaxTestSvc()
	ctx := context.Background()

	// Generate once
	_, err := svc.GenerateCalendarBatch(ctx, "c1", 2026, []domain.TaxType{domain.TaxTypeVAT})
	require.NoError(t, err)

	// Generate again — should skip existing
	count, err := svc.GenerateCalendarBatch(ctx, "c1", 2026, []domain.TaxType{domain.TaxTypeVAT})
	require.NoError(t, err)
	assert.Equal(t, 0, count, "no new entries created")

	entries, _ := repo.GetCalendarByCompany(ctx, "c1")
	assert.Len(t, entries, 12, "still only 12 entries")
}

func TestGenerateCalendarBatch_AllTaxTypes(t *testing.T) {
	svc, repo := newTaxTestSvc()
	ctx := context.Background()

	count, err := svc.GenerateCalendarBatch(ctx, "c1", 2026, []domain.TaxType{
		domain.TaxTypeVAT, domain.TaxTypeCIT, domain.TaxTypePIT,
		domain.TaxTypeTTDB, domain.TaxTypeBVMT,
	})
	require.NoError(t, err)
	// VAT:12 + CIT:4 + PIT:12 + TTDB:4 + BVMT:4 = 36
	assert.Equal(t, 36, count)

	entries, err := repo.GetCalendarByCompany(ctx, "c1")
	require.NoError(t, err)
	assert.Len(t, entries, 36)

	byType := map[domain.TaxType]int{}
	for _, e := range entries {
		byType[e.TaxType]++
	}
	assert.Equal(t, 12, byType[domain.TaxTypeVAT])
	assert.Equal(t, 4, byType[domain.TaxTypeCIT])
	assert.Equal(t, 12, byType[domain.TaxTypePIT])
	assert.Equal(t, 4, byType[domain.TaxTypeTTDB])
	assert.Equal(t, 4, byType[domain.TaxTypeBVMT])
}

func TestGenerateDeadlineAlerts(t *testing.T) {
	svc, repo := newTaxTestSvc()
	ctx := context.Background()

	// Create calendar entry due in 3 days
	require.NoError(t, repo.CreateCalendarEntry(ctx, &domain.TaxCalendar{
		ID: "CAL-1", CompanyID: "c1", TaxType: domain.TaxTypeVAT,
		PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 1,
		StartDate: "2026-01-01", EndDate: "2026-01-31",
		DeclarationDue: "2026-02-20", PaymentDue: "2026-02-20",
		Status: domain.CalStatusPENDING,
	}))

	// Create calendar entry due today
	require.NoError(t, repo.CreateCalendarEntry(ctx, &domain.TaxCalendar{
		ID: "CAL-2", CompanyID: "c1", TaxType: domain.TaxTypeCIT,
		PeriodType: domain.PeriodTypeQuarterly, PeriodYear: 2026, PeriodNumber: 1,
		StartDate: "2026-01-01", EndDate: "2026-03-31",
		DeclarationDue: "2026-04-30", PaymentDue: "2026-04-30",
		Status: domain.CalStatusPENDING,
	}))

	// Set "now" to 2026-02-17 (3 days before CAL-1 due, 72 days before CAL-2)
	svc.now = func() time.Time { return time.Date(2026, 2, 17, 10, 0, 0, 0, time.UTC) }

	count, err := svc.GenerateDeadlineAlerts(ctx, "c1", 7)
	require.NoError(t, err)
	assert.Equal(t, 1, count, "only CAL-1 is within 7 days")

	alerts, err := repo.GetAlerts(ctx, "c1", 10)
	require.NoError(t, err)
	require.Len(t, alerts, 1)
	assert.Equal(t, "CAL-1", alerts[0].CalendarID)
	assert.Equal(t, domain.AlertTypeWARNING, alerts[0].AlertType)
	assert.Contains(t, alerts[0].Message, "VAT")
}

// ─── Tax Rate Lookup for TTDB/BVMT/NTNN ──────────────────────────────────

func TestGenerateDeclaration_TTDB_UsesRateFromTable(t *testing.T) {
	taxRepo := repository.NewMemoryTaxRepo()
	jeRepo := repository.NewMemoryJournalRepo()
	svc := NewTaxService(taxRepo, jeRepo, repository.NewMemoryCompanyRepo(), nil, nil).(*taxService)
	ctx := context.Background()

	// Create active TTDB rate
	require.NoError(t, taxRepo.CreateRate(ctx, &domain.TaxRate{
		RateCode: "TTDB_STANDARD", TaxType: domain.TaxTypeTTDB,
		RateType: domain.RateTypePERCENTAGE, RateValue: 15,
		EffectiveFrom: "2025-01-01", IsActive: true,
	}))

	// Post journal entry with 3332 credit (special consumption tax base)
	je := postedEntry("JE-TTDB-1", "2026-03-15",
		domain.JournalLine{AccountCode: "3332", CreditAmount: 1000000})
	require.NoError(t, jeRepo.Create(ctx, &je))

	decl, err := svc.GenerateDeclaration(ctx, "c1", domain.DeclTypeTTDB01,
		domain.TaxPeriod{PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 3}, "u1", nil)
	require.NoError(t, err)
	require.Len(t, decl.Lines, 3)

	// Line 10: base = 1,000,000
	assert.Equal(t, "10", decl.Lines[0].LineCode)
	assert.Equal(t, 1000000.0, decl.Lines[0].Amount)

	// Line 20: rate = 15 (from table, not default 10)
	assert.Equal(t, "20", decl.Lines[1].LineCode)
	assert.Equal(t, 15.0, decl.Lines[1].Amount)

	// Line 30: tax payable = 1,000,000 * 15 / 100 = 150,000
	assert.Equal(t, "30", decl.Lines[2].LineCode)
	assert.Equal(t, 150000.0, decl.Lines[2].Amount)
}

func TestGenerateDeclaration_NTNN_FallbackToDefault(t *testing.T) {
	taxRepo := repository.NewMemoryTaxRepo()
	jeRepo := repository.NewMemoryJournalRepo()
	svc := NewTaxService(taxRepo, jeRepo, repository.NewMemoryCompanyRepo(), nil, nil).(*taxService)
	ctx := context.Background()

	// No rate in table — should use default 5%
	je := postedEntry("JE-NTNN-1", "2026-03-15",
		domain.JournalLine{AccountCode: "511", CreditAmount: 2000000})
	require.NoError(t, jeRepo.Create(ctx, &je))

	decl, err := svc.GenerateDeclaration(ctx, "c1", domain.DeclTypeNTNN01,
		domain.TaxPeriod{PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 3}, "u1", nil)
	require.NoError(t, err)
	require.Len(t, decl.Lines, 3)

	// Line 10: income = 2,000,000
	assert.Equal(t, 2000000.0, decl.Lines[0].Amount)

	// Line 20: rate = 5 (default)
	assert.Equal(t, 5.0, decl.Lines[1].Amount)

	// Line 30: tax = 2,000,000 * 5 / 100 = 100,000
	assert.Equal(t, 100000.0, decl.Lines[2].Amount)
}
