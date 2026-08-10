package service

import (
	"context"

	"gotax/internal/domain"
)

type TaxDeclarationService struct {
	declRepo domain.TaxRepository
}

func NewTaxDeclarationService(
	declRepo domain.TaxRepository,
) *TaxDeclarationService {
	return &TaxDeclarationService{
		declRepo: declRepo,
	}
}

// PopulateVAT01 auto-fills VAT declaration (01/GTGT) lines from invoice data.
// Per Circular 99/2025 Appendix I:
//
//	[10] = Tổng doanh thu hàng hóa, dịch vụ bán ra (revenue)
//	[20] = Tổng số thuế GTGT hàng hóa, dịch vụ bán ra (output VAT)
//	[25] = Tổng số thuế GTGT được khấu trừ (input VAT deductible)
//	[30] = Số thuế GTGT phải nộp = [20] - [25]
func (s *TaxDeclarationService) PopulateVAT01(ctx context.Context, companyID string, periodYear int, periodNumber int, revenue, outputVAT, inputVAT float64) (*domain.TaxDeclaration, error) {
	vatPayable := outputVAT - inputVAT
	if vatPayable < 0 {
		vatPayable = 0
	}

	decl := &domain.TaxDeclaration{
		CompanyID:       companyID,
		DeclarationType: domain.DeclTypeGTGT01,
		TaxPeriod: domain.TaxPeriod{
			PeriodType:   domain.PeriodTypeMonthly,
			PeriodYear:   periodYear,
			PeriodNumber: periodNumber,
		},
		AdjustmentType: domain.AdjTypeNONE,
		Lines: []domain.TaxDeclarationLine{
			{LineCode: "10", Amount: revenue, LineName: "Tổng doanh thu HĐ, DV bán ra"},
			{LineCode: "20", Amount: outputVAT, LineName: "Tổng thuế GTGT HĐ, DV bán ra"},
			{LineCode: "25", Amount: inputVAT, LineName: "Tổng thuế GTGT được khấu trừ"},
			{LineCode: "30", Amount: vatPayable, LineName: "Thuế GTGT phải nộp"},
		},
	}
	return decl, nil
}

// PopulateCIT03 auto-fills CIT finalization (03/TNDN) per Law 67/2025.
// 3-tier rates: 15% (≤2B VND), 17% (2-18B), 20% (>18B)
func (s *TaxDeclarationService) PopulateCIT03(ctx context.Context, companyID string, year int, revenue, deductibleExpenses float64) (*domain.TaxDeclaration, error) {
	taxableIncome := revenue - deductibleExpenses
	if taxableIncome < 0 {
		taxableIncome = 0
	}

	// 3-tier calculation per Law 67/2025/QH15
	var tax15, tax17, tax20 float64
	remaining := taxableIncome

	const (
		tier1Limit = 2_000_000_000  // 2 billion VND
		tier2Limit = 18_000_000_000 // 18 billion VND
	)

	if remaining > 0 {
		tier1 := remaining
		if tier1 > tier1Limit {
			tier1 = tier1Limit
		}
		tax15 = tier1 * 0.15
		remaining -= tier1
	}
	if remaining > 0 {
		tier2 := remaining
		if tier2 > (tier2Limit - tier1Limit) {
			tier2 = tier2Limit - tier1Limit
		}
		tax17 = tier2 * 0.17
		remaining -= tier2
	}
	if remaining > 0 {
		tax20 = remaining * 0.20
	}

	totalTax := tax15 + tax17 + tax20

	decl := &domain.TaxDeclaration{
		CompanyID:       companyID,
		DeclarationType: domain.DeclTypeTNDN03,
		TaxPeriod: domain.TaxPeriod{
			PeriodType:   domain.PeriodTypeAnnual,
			PeriodYear:   year,
			PeriodNumber: 1,
		},
		AdjustmentType: domain.AdjTypeNONE,
		Lines: []domain.TaxDeclarationLine{
			{LineCode: "10", Amount: revenue, LineName: "Tổng doanh thu tính thuế"},
			{LineCode: "20", Amount: deductibleExpenses, LineName: "Tổng chi phí được trừ"},
			{LineCode: "30", Amount: taxableIncome, LineName: "Thu nhập tính thuế"},
			{LineCode: "31a", Amount: tax15, LineName: "Thuế TNDN 15% (≤2 tỷ)"},
			{LineCode: "31b", Amount: tax17, LineName: "Thuế TNDN 17% (2-18 tỷ)"},
			{LineCode: "31c", Amount: tax20, LineName: "Thuế TNDN 20% (>18 tỷ)"},
			{LineCode: "40", Amount: totalTax, LineName: "Thuế TNDN phải nộp"},
		},
	}
	return decl, nil
}

// PopulatePIT05 auto-fills PIT declaration (05/KK-TNCN) per Law 109/2025.
func (s *TaxDeclarationService) PopulatePIT05(ctx context.Context, companyID string, periodYear int, periodNumber int, totalIncome, deductions, exemptIncome float64) (*domain.TaxDeclaration, error) {
	taxableIncome := totalIncome - deductions - exemptIncome
	if taxableIncome < 0 {
		taxableIncome = 0
	}

	decl := &domain.TaxDeclaration{
		CompanyID:       companyID,
		DeclarationType: domain.DeclTypeKKTNCN,
		TaxPeriod: domain.TaxPeriod{
			PeriodType:   domain.PeriodTypeMonthly,
			PeriodYear:   periodYear,
			PeriodNumber: periodNumber,
		},
		AdjustmentType: domain.AdjTypeNONE,
		Lines: []domain.TaxDeclarationLine{
			{LineCode: "10", Amount: totalIncome, LineName: "Tổng thu nhập chịu thuế"},
			{LineCode: "20", Amount: deductions, LineName: "Tổng giảm trừ gia cảnh"},
			{LineCode: "25.1", Amount: exemptIncome, LineName: "Thu nhập miễn thuế theo nghị quyết"},
			{LineCode: "30", Amount: taxableIncome, LineName: "Thu nhập tính thuế"},
			{LineCode: "40", Amount: 0, LineName: "Thuế TNCN phải nộp"},
		},
	}
	return decl, nil
}

// CreateDeclaration creates and persists a tax declaration.
func (s *TaxDeclarationService) CreateDeclaration(ctx context.Context, decl *domain.TaxDeclaration) error {
	if decl.CompanyID == "" {
		return domain.ErrCompanyIDRequired
	}
	if !decl.DeclarationType.Valid() {
		return domain.ErrDeclarationTypeInvalid
	}
	return s.declRepo.CreateDeclaration(ctx, decl)
}
