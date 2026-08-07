package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"gotax/internal/domain"
)

type CostingEngine struct {
	periodRepo      domain.CostingPeriodRepository
	costObjectRepo  domain.CostObjectRepository
	costPoolRepo    domain.CostPoolRepository
	costPoolLineRepo domain.CostPoolLineRepository
	resultRepo      domain.CostingResultRepository
	resultLineRepo  domain.CostingResultLineRepository
}

func NewCostingEngine(
	periodRepo domain.CostingPeriodRepository,
	costObjectRepo domain.CostObjectRepository,
	costPoolRepo domain.CostPoolRepository,
	costPoolLineRepo domain.CostPoolLineRepository,
	resultRepo domain.CostingResultRepository,
	resultLineRepo domain.CostingResultLineRepository,
) *CostingEngine {
	return &CostingEngine{
		periodRepo:      periodRepo,
		costObjectRepo:  costObjectRepo,
		costPoolRepo:    costPoolRepo,
		costPoolLineRepo: costPoolLineRepo,
		resultRepo:      resultRepo,
		resultLineRepo:  resultLineRepo,
	}
}

func (e *CostingEngine) RunCosting(ctx context.Context, companyID, periodID string) error {
	period, err := e.periodRepo.GetByID(ctx, periodID)
	if err != nil {
		return domain.ErrCostingPeriodNotFound
	}
	if period.Status == "CLOSED" {
		return domain.ErrCostingPeriodAlreadyClosed
	}

	pools, err := e.costPoolRepo.ListByPeriod(ctx, companyID, periodID)
	if err != nil {
		return err
	}
	if len(pools) == 0 {
		return domain.ErrCostPoolNoLines
	}

	objects, err := e.costObjectRepo.List(ctx, companyID)
	if err != nil {
		return err
	}
	if len(objects) == 0 {
		return fmt.Errorf("no cost objects found for company")
	}

	for i := range objects {
		if !objects[i].IsActive {
			continue
		}
		if err := e.calculateForPeriod(ctx, companyID, periodID, &objects[i], pools); err != nil {
			return fmt.Errorf("costing failed for object %s: %w", objects[i].Code, err)
		}
	}

	return nil
}

func costingGenID(prefix string) string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%s-%x", prefix, b)
}

func (e *CostingEngine) calculateForPeriod(ctx context.Context, companyID, periodID string, obj *domain.CostObject, pools []domain.CostPool) error {
	totalDirectMat := 0.0
	totalDirectLab := 0.0
	totalOverhead := 0.0

	for _, pool := range pools {
		lines, err := e.costPoolLineRepo.ListByPool(ctx, pool.ID)
		if err != nil {
			return err
		}

		poolAmount := 0.0
		for _, line := range lines {
			poolAmount += line.Amount
		}

		switch pool.GLAccountCode {
		case "621":
			totalDirectMat += poolAmount
		case "622":
			totalDirectLab += poolAmount
		case "623", "627":
			totalOverhead += poolAmount
		}
	}

	totalCost := totalDirectMat + totalDirectLab + totalOverhead

	var unitCost float64
	if obj.PlanQuantity > 0 {
		unitCost = totalCost / obj.PlanQuantity
	}

	result := &domain.CostingResult{
		ID:             costingGenID("CR"),
		CompanyID:      companyID,
		PeriodID:       periodID,
		CostObjectID:   obj.ID,
		CostingMethod:  obj.CostingMethod,
		TotalDirectMat: totalDirectMat,
		TotalDirectLab: totalDirectLab,
		TotalOverhead:  totalOverhead,
		TotalCost:      totalCost,
		OutputQuantity: obj.PlanQuantity,
		UnitCost:       unitCost,
		Status:         "DRAFT",
		CreatedAt:      time.Now().Format("2006-01-02T15:04:05Z"),
		UpdatedAt:      time.Now().Format("2006-01-02T15:04:05Z"),
	}

	if err := e.resultRepo.Create(ctx, result); err != nil {
		return err
	}

	line1 := &domain.CostingResultLine{
		ID:            costingGenID("CRL"),
		ResultID:      result.ID,
		CostCategory:  "DIRECT_MATERIAL",
		GLAccountCode: "621",
		Description:   "Direct materials",
		PlannedAmount: 0,
		ActualAmount:  totalDirectMat,
		AllocatedAmount: totalDirectMat,
		CreatedAt:     time.Now().Format("2006-01-02T15:04:05Z"),
	}
	if err := e.resultLineRepo.Create(ctx, line1); err != nil {
		return err
	}

	line2 := &domain.CostingResultLine{
		ID:            costingGenID("CRL"),
		ResultID:      result.ID,
		CostCategory:  "DIRECT_LABOR",
		GLAccountCode: "622",
		Description:   "Direct labor",
		PlannedAmount: 0,
		ActualAmount:  totalDirectLab,
		AllocatedAmount: totalDirectLab,
		CreatedAt:     time.Now().Format("2006-01-02T15:04:05Z"),
	}
	if err := e.resultLineRepo.Create(ctx, line2); err != nil {
		return err
	}

	line3 := &domain.CostingResultLine{
		ID:            costingGenID("CRL"),
		ResultID:      result.ID,
		CostCategory:  "OVERHEAD",
		GLAccountCode: "627",
		Description:   "Manufacturing overhead",
		PlannedAmount: 0,
		ActualAmount:  totalOverhead,
		AllocatedAmount: totalOverhead,
		CreatedAt:     time.Now().Format("2006-01-02T15:04:05Z"),
	}
	if err := e.resultLineRepo.Create(ctx, line3); err != nil {
		return err
	}

	return nil
}
