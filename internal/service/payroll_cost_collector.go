package service

import (
	"context"

	"gotax/internal/domain"
)

// PayrollCostCollector implements CostDataCollector for payroll direct labor costs.
type PayrollCostCollector struct {
	payrollRepo interface {
		ListRunsByPeriod(ctx context.Context, periodID string) ([]domain.PayrollRun, error)
	}
}

func NewPayrollCostCollector(payrollRepo interface {
	ListRunsByPeriod(ctx context.Context, periodID string) ([]domain.PayrollRun, error)
}) *PayrollCostCollector {
	return &PayrollCostCollector{payrollRepo: payrollRepo}
}

func (c *PayrollCostCollector) CollectMaterialCosts(ctx context.Context, companyID, periodID string) ([]domain.CostPoolLineInput, error) {
	return nil, nil
}

func (c *PayrollCostCollector) CollectLaborCosts(ctx context.Context, companyID, periodID string) ([]domain.CostPoolLineInput, error) {
	runs, err := c.payrollRepo.ListRunsByPeriod(ctx, periodID)
	if err != nil {
		return nil, err
	}

	var lines []domain.CostPoolLineInput
	for _, run := range runs {
		if run.CompanyID != companyID {
			continue
		}
		if run.GrossSalary <= 0 {
			continue
		}
		lines = append(lines, domain.CostPoolLineInput{
			SourceID:    run.ID,
			Description: "Direct labor - " + run.EmployeeID,
			Amount:      run.GrossSalary,
		})
	}

	return lines, nil
}

func (c *PayrollCostCollector) CollectOverheadCosts(ctx context.Context, companyID, periodID string) ([]domain.CostPoolLineInput, error) {
	return nil, nil
}
