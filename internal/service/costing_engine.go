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

func (e *CostingEngine) createResultLine(ctx context.Context, resultID string, category domain.CostCategory, glAcct, desc string, planned, actual, allocated float64) error {
	line := &domain.CostingResultLine{
		ID:              costingGenID("CRL"),
		ResultID:        resultID,
		CostCategory:    category,
		GLAccountCode:   glAcct,
		Description:     desc,
		PlannedAmount:   planned,
		ActualAmount:    actual,
		AllocatedAmount: allocated,
		CreatedAt:       time.Now().Format("2006-01-02T15:04:05Z"),
	}
	return e.resultLineRepo.Create(ctx, line)
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

	// Process costing: compute equivalent units from completed + WIP
	outputQuantity := obj.PlanQuantity
	if obj.CostingMethod == "PROCESS" {
		outputQuantity = obj.CompletedUnits + (obj.WIPUnits * obj.CompletionPct / 100)
	}

	var unitCost float64
	if outputQuantity > 0 {
		unitCost = totalCost / outputQuantity
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
		OutputQuantity: outputQuantity,
		UnitCost:       unitCost,
		Status:         "DRAFT",
		CreatedAt:      time.Now().Format("2006-01-02T15:04:05Z"),
		UpdatedAt:      time.Now().Format("2006-01-02T15:04:05Z"),
	}

	if err := e.resultRepo.Create(ctx, result); err != nil {
		return err
	}

	if err := e.createResultLine(ctx, result.ID, domain.CostCategoryDirectMaterial, "621", "Direct materials", 0, totalDirectMat, totalDirectMat); err != nil {
		return err
	}
	if err := e.createResultLine(ctx, result.ID, domain.CostCategoryDirectLabor, "622", "Direct labor", 0, totalDirectLab, totalDirectLab); err != nil {
		return err
	}
	if err := e.createResultLine(ctx, result.ID, domain.CostCategoryOverhead, "627", "Manufacturing overhead", 0, totalOverhead, totalOverhead); err != nil {
		return err
	}

	// Standard costing: variance = actual - (std_mat + std_lab + std_oh) × quantity
	if obj.CostingMethod == "STANDARD" {
		standardTotal := (obj.StandardMaterial + obj.StandardLabor + obj.StandardOverhead) * obj.PlanQuantity
		variance := totalCost - standardTotal
		varianceType := "UNFAVORABLE"
		varGlAcct := "627" // unfavorable → debit overhead
		if variance < 0 {
			varianceType = "FAVORABLE"
			varGlAcct = "627" // favorable → credit overhead (same account, opposite entry)
		}

		if err := e.createResultLine(ctx, result.ID, domain.CostCategoryVariance, varGlAcct, varianceType, standardTotal, variance, variance); err != nil {
			return err
		}
	}

	// Process costing: record equivalent units info
	if obj.CostingMethod == "PROCESS" {
		if err := e.createResultLine(ctx, result.ID, domain.CostCategoryProcessInfo, "", "Equivalent units", outputQuantity, totalCost, unitCost); err != nil {
			return err
		}
	}

	return nil
}
