package service

import (
	"context"
	"time"

	"gotax/internal/domain"
)

type FinancialAnalysisService struct {
	jeRepo   domain.JournalRepository
	budRepo  domain.BudgetRepository
	compRepo domain.CompanyRepository
}

func NewFinancialAnalysisService(jeRepo domain.JournalRepository, budRepo domain.BudgetRepository, compRepo domain.CompanyRepository) *FinancialAnalysisService {
	return &FinancialAnalysisService{jeRepo: jeRepo, budRepo: budRepo, compRepo: compRepo}
}

// FinancialRatio holds common financial ratios.
type FinancialRatio struct {
	Name        string  `json:"name"`
	Value       float64 `json:"value"`
	Description string  `json:"description"`
}

// CalculateRatios computes financial ratios from GL data via Trial Balance.
func (s *FinancialAnalysisService) CalculateRatios(ctx context.Context, companyID string, year int) ([]FinancialRatio, error) {
	// Get latest period for the year to pull trial balance
	periodID := ""
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)
	entries, err := s.jeRepo.GetByDateRange(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Aggregate by account prefix
	totals := make(map[string]float64)
	for _, e := range entries {
		for _, l := range e.Lines {
			prefix := l.AccountCode[:min(3, len(l.AccountCode))]
			totals[prefix] += l.DebitAmount - l.CreditAmount
		}
	}
	_ = periodID

	revenue := abs(totals["511"]) + abs(totals["512"]) + abs(totals["515"])
	expense := abs(totals["6"]) + abs(totals["8"])
	cogs := abs(totals["632"]) + abs(totals["633"])
	currentAssets := totals["111"] + totals["112"] + totals["113"] + totals["128"]
	currentLiabilities := abs(totals["331"]) + abs(totals["333"]) + abs(totals["341"])
	totalAssets := currentAssets + totals["121"] + totals["122"] + totals["151"] + totals["152"] + totals["153"]
	totalEquity := totals["411"] + totals["412"] - totals["413"]
	totalLiabilities := totalAssets - totalEquity

	var ratios []FinancialRatio
	if revenue > 0 {
		profitMargin := (revenue - expense) / revenue * 100
		ratios = append(ratios, FinancialRatio{"Net Profit Margin", round2(profitMargin), "Net income / Revenue"})
		grossMargin := (revenue - cogs) / revenue * 100
		ratios = append(ratios, FinancialRatio{"Gross Profit Margin", round2(grossMargin), "(Revenue - COGS) / Revenue"})
	}
	if currentLiabilities > 0 {
		ratios = append(ratios, FinancialRatio{"Current Ratio", round2(currentAssets / currentLiabilities), "Current Assets / Current Liabilities"})
		quickRatio := (currentAssets - totals["113"]) / currentLiabilities
		ratios = append(ratios, FinancialRatio{"Quick Ratio", round2(quickRatio), "(Current Assets - Inventory) / Current Liabilities"})
	}
	if totalEquity > 0 {
		ratios = append(ratios, FinancialRatio{"Debt-to-Equity", round2(totalLiabilities / totalEquity), "Total Liabilities / Total Equity"})
	}
	if totalAssets > 0 {
		ratios = append(ratios, FinancialRatio{"Asset Turnover", round2(revenue / totalAssets), "Revenue / Total Assets"})
		roa := (revenue - expense) / totalAssets * 100
		ratios = append(ratios, FinancialRatio{"ROA", round2(roa), "Net Income / Total Assets"})
	}
	if totalEquity > 0 {
		roe := (revenue - expense) / totalEquity * 100
		ratios = append(ratios, FinancialRatio{"ROE", round2(roe), "Net Income / Total Equity"})
	}
	return ratios, nil
}

// BudgetVsActualItem compares budget vs actual for one account.
type BudgetVsActualItem struct {
	AccountCode string  `json:"account_code"`
	Budget      float64 `json:"budget"`
	Actual      float64 `json:"actual"`
	Variance    float64 `json:"variance"`
	VariancePct float64 `json:"variance_pct"`
}

type BudgetVsActualResult struct {
	CompanyID  string               `json:"company_id"`
	PeriodYear int                  `json:"period_year"`
	Items      []BudgetVsActualItem `json:"items"`
}

func (s *FinancialAnalysisService) CompareBudgetVsActual(ctx context.Context, companyID string, periodYear int) (*BudgetVsActualResult, error) {
	budgets, err := s.budRepo.List(ctx, companyID, periodYear)
	if err != nil {
		return nil, err
	}
	startDate := time.Date(periodYear, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(periodYear, 12, 31, 23, 59, 59, 0, time.UTC)
	entries, err := s.jeRepo.GetByDateRange(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}
	actuals := make(map[string]float64)
	for _, e := range entries {
		for _, l := range e.Lines {
			actuals[l.AccountCode] += l.DebitAmount - l.CreditAmount
		}
	}

	result := &BudgetVsActualResult{CompanyID: companyID, PeriodYear: periodYear}
	seen := make(map[string]bool)
	for _, b := range budgets {
		code := b.AccountCode
		if seen[code] {
			continue
		}
		seen[code] = true
		item := BudgetVsActualItem{
			AccountCode: code,
			Budget:      b.Budgeted,
			Actual:      abs(actuals[code]),
		}
		item.Variance = item.Actual - item.Budget
		if item.Budget != 0 {
			item.VariancePct = item.Variance / item.Budget * 100
		}
		result.Items = append(result.Items, item)
	}
	return result, nil
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
