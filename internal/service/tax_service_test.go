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

func newTaxTestSvc() (*taxService, domain.TaxRepository) {
	repo := repository.NewMemoryTaxRepo()
	return NewTaxService(repo, repository.NewMemoryJournalRepo()).(*taxService), repo
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
	return NewTaxService(repo, jeRepo).(*taxService), jeRepo
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
		vatPeriod(), "user-1")
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
		"user-1")
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
	_, err := svc.GenerateDeclaration(ctx, "c1", domain.DeclTypeGTGT01, vatPeriod(), "u1")
	require.NoError(t, err)
	_, err = svc.GenerateDeclaration(ctx, "c1", domain.DeclTypeGTGT01, vatPeriod(), "u1")
	assert.ErrorIs(t, err, domain.ErrDuplicateDeclaration)
}

func TestGenerateDeclaration_UnsupportedType(t *testing.T) {
	svc, _ := newTaxTestSvcWithGL()
	ctx := context.Background()
	_, err := svc.GenerateDeclaration(ctx, "c1", domain.DeclTypeKKTNCN, vatPeriod(), "u1")
	assert.ErrorIs(t, err, domain.ErrDeclarationTypeInvalid)
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

	decl, err := svc.GenerateDeclaration(ctx, "c1", domain.DeclTypeGTGT01, vatPeriod(), "u1")
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
	decl, err := svc.GenerateDeclaration(ctx, "c1", domain.DeclTypeGTGT01, vatPeriod(), "u1")
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

	decl, err := svc.GenerateDeclaration(ctx, "c1", domain.DeclTypeGTGT01, vatPeriod(), "u1")
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

	decl, err := svc.GenerateDeclaration(ctx, "c1", domain.DeclTypeGTGT01, vatPeriod(), "u1")
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

	decl, err := svc.GenerateDeclaration(ctx, "c1", domain.DeclTypeGTGT01, vatPeriod(), "u1")
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
		domain.TaxPeriod{PeriodType: domain.PeriodTypeAnnual, PeriodYear: 2026, PeriodNumber: 1}, "u1")
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
