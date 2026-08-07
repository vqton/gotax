package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gotax/internal/domain"
	"gotax/internal/repository"
)

func TestWarehouseCostCollector_CollectMaterialCostsByDateRange(t *testing.T) {
	invRepo := repository.NewMemoryInventoryTransactionRepo()
	collector := NewWarehouseCostCollector(invRepo)
	ctx := context.Background()

	now := time.Now()
	from := now.AddDate(0, 0, -1)
	to := now.AddDate(0, 0, 1)

	txn := &domain.InventoryTransaction{
		ID:        "TXN-001",
		CompanyID: "COMP1",
		TransType: domain.TransIssue,
		TotalCost: 50000000,
		CreatedAt: now,
	}
	_ = invRepo.CreateInventoryTransaction(ctx, txn)

	lines, err := collector.CollectMaterialCostsByDateRange(ctx, "COMP1", from, to)
	assert.NoError(t, err)
	assert.Len(t, lines, 1)
	assert.Equal(t, 50000000.0, lines[0].Amount)
}

func TestPayrollCostCollector_CollectLaborCosts(t *testing.T) {
	payrollRepo := repository.NewMemoryPayrollRepo()
	collector := NewPayrollCostCollector(payrollRepo)
	ctx := context.Background()

	run := &domain.PayrollRun{
		ID:          "PR-001",
		PeriodID:    "PERIOD-001",
		CompanyID:   "COMP1",
		GrossSalary: 20000000,
	}
	_ = payrollRepo.CreateRun(ctx, run)

	lines, err := collector.CollectLaborCosts(ctx, "COMP1", "PERIOD-001")
	assert.NoError(t, err)
	assert.Len(t, lines, 1)
	assert.Equal(t, 20000000.0, lines[0].Amount)
}
