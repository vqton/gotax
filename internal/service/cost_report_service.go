package service

import (
	"context"

	"gotax/internal/domain"
)

type CostReportService struct {
	resultRepo     domain.CostingResultRepository
	resultLineRepo domain.CostingResultLineRepository
	costObjectRepo domain.CostObjectRepository
}

func NewCostReportService(
	resultRepo domain.CostingResultRepository,
	resultLineRepo domain.CostingResultLineRepository,
	costObjectRepo domain.CostObjectRepository,
) *CostReportService {
	return &CostReportService{
		resultRepo:     resultRepo,
		resultLineRepo: resultLineRepo,
		costObjectRepo: costObjectRepo,
	}
}

type CostCalculationSheet struct {
	PeriodID    string
	ObjectID    string
	ObjectCode  string
	ObjectName  string
	TotalCost   float64
	UnitCost    float64
	Quantity    float64
	Lines       []CostCalculationLine
}

type CostCalculationLine struct {
	CostCategory string
	GLAccount    string
	Description  string
	Amount       float64
}

func (s *CostReportService) GetCostCalculationSheet(ctx context.Context, companyID, periodID, objectID string) (*CostCalculationSheet, error) {
	results, err := s.resultRepo.ListByPeriod(ctx, companyID, periodID)
	if err != nil {
		return nil, err
	}

	var result domain.CostingResult
	found := false
	for _, r := range results {
		if r.CostObjectID == objectID {
			result = r
			found = true
			break
		}
	}
	if !found {
		return nil, domain.ErrCostingResultNotFound
	}

	obj, err := s.costObjectRepo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	lines, err := s.resultLineRepo.ListByResult(ctx, result.ID)
	if err != nil {
		return nil, err
	}

	calcLines := make([]CostCalculationLine, len(lines))
	for i, l := range lines {
		calcLines[i] = CostCalculationLine{
			CostCategory: string(l.CostCategory),
			GLAccount:    l.GLAccountCode,
			Description:  l.Description,
			Amount:       l.ActualAmount,
		}
	}

	return &CostCalculationSheet{
		PeriodID:   periodID,
		ObjectID:   objectID,
		ObjectCode: obj.Code,
		ObjectName: obj.Name,
		TotalCost:  result.TotalCost,
		UnitCost:   result.UnitCost,
		Quantity:   result.OutputQuantity,
		Lines:      calcLines,
	}, nil
}

type CostSummary struct {
	PeriodID         string
	Items            []CostSummaryItem
	TotalAllObjects  float64
	TotalUnitCost    float64
}

type CostSummaryItem struct {
	ObjectID    string
	ObjectCode  string
	ObjectName  string
	TotalCost   float64
	UnitCost    float64
	Quantity    float64
	Method      string
}

func (s *CostReportService) GetCostSummary(ctx context.Context, companyID, periodID string) (*CostSummary, error) {
	results, err := s.resultRepo.ListByPeriod(ctx, companyID, periodID)
	if err != nil {
		return nil, err
	}

	items := make([]CostSummaryItem, len(results))
	totalCost := 0.0
	for i, r := range results {
		obj, err := s.costObjectRepo.GetByID(ctx, r.CostObjectID)
		if err != nil {
			continue
		}
		items[i] = CostSummaryItem{
			ObjectID:   r.CostObjectID,
			ObjectCode: obj.Code,
			ObjectName: obj.Name,
			TotalCost:  r.TotalCost,
			UnitCost:   r.UnitCost,
			Quantity:   r.OutputQuantity,
			Method:     string(r.CostingMethod),
		}
		totalCost += r.TotalCost
	}

	return &CostSummary{
		PeriodID:        periodID,
		Items:           items,
		TotalAllObjects: totalCost,
	}, nil
}

type WIPValuationReport struct {
	PeriodID string
	Items    []WIPValuationItem
	TotalWIP float64
}

type WIPValuationItem struct {
	ObjectID   string
	ObjectCode string
	ObjectName string
	WIPCost    float64
	Quantity   float64
	UnitCost   float64
}

func (s *CostReportService) GetWIPValuation(ctx context.Context, companyID, periodID string) (*WIPValuationReport, error) {
	results, err := s.resultRepo.ListByPeriod(ctx, companyID, periodID)
	if err != nil {
		return nil, err
	}

	items := make([]WIPValuationItem, 0)
	totalWIP := 0.0
	for _, r := range results {
		if r.WIPBegin > 0 || r.WIPEnd > 0 {
			obj, err := s.costObjectRepo.GetByID(ctx, r.CostObjectID)
			if err != nil {
				continue
			}
			wipCost := r.WIPEnd * r.UnitCost
			items = append(items, WIPValuationItem{
				ObjectID:   r.CostObjectID,
				ObjectCode: obj.Code,
				ObjectName: obj.Name,
				WIPCost:    wipCost,
				Quantity:   r.WIPEnd,
				UnitCost:   r.UnitCost,
			})
			totalWIP += wipCost
		}
	}

	// If no WIP-specific data, show all results as WIP (simplified)
	if len(items) == 0 {
		for _, r := range results {
			obj, err := s.costObjectRepo.GetByID(ctx, r.CostObjectID)
			if err != nil {
				continue
			}
			items = append(items, WIPValuationItem{
				ObjectID:   r.CostObjectID,
				ObjectCode: obj.Code,
				ObjectName: obj.Name,
				WIPCost:    r.TotalCost,
				Quantity:   r.OutputQuantity,
				UnitCost:   r.UnitCost,
			})
			totalWIP += r.TotalCost
		}
	}

	return &WIPValuationReport{
		PeriodID: periodID,
		Items:    items,
		TotalWIP: totalWIP,
	}, nil
}

type VarianceAnalysis struct {
	PeriodID string
	Items    []VarianceItem
}

type VarianceItem struct {
	ObjectID     string
	ObjectCode   string
	ObjectName   string
	ActualCost   float64
	StandardCost float64
	Variance     float64
	VarianceType string
}

func (s *CostReportService) GetVarianceAnalysis(ctx context.Context, companyID, periodID string) (*VarianceAnalysis, error) {
	results, err := s.resultRepo.ListByPeriod(ctx, companyID, periodID)
	if err != nil {
		return nil, err
	}

	items := make([]VarianceItem, 0)
	for _, r := range results {
		if r.CostingMethod != "STANDARD" {
			continue
		}

		obj, err := s.costObjectRepo.GetByID(ctx, r.CostObjectID)
		if err != nil {
			continue
		}

		standardCost := (obj.StandardMaterial + obj.StandardLabor + obj.StandardOverhead) * obj.PlanQuantity
		variance := r.TotalCost - standardCost
		varianceType := "UNFAVORABLE"
		if variance < 0 {
			varianceType = "FAVORABLE"
		}

		items = append(items, VarianceItem{
			ObjectID:     r.CostObjectID,
			ObjectCode:   obj.Code,
			ObjectName:   obj.Name,
			ActualCost:   r.TotalCost,
			StandardCost: standardCost,
			Variance:     variance,
			VarianceType: varianceType,
		})
	}

	return &VarianceAnalysis{
		PeriodID: periodID,
		Items:    items,
	}, nil
}
