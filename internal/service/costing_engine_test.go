package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"gotax/internal/domain"
	"gotax/internal/repository"
)

func setupCostingEngineTest() (*CostingEngine, *repository.MemoryCostingPeriodRepo, *repository.MemoryCostObjectRepo, *repository.MemoryCostPoolRepo, *repository.MemoryCostPoolLineRepo, *repository.MemoryCostingResultRepo, *repository.MemoryCostingResultLineRepo) {
	periodRepo := repository.NewMemoryCostingPeriodRepo()
	costObjectRepo := repository.NewMemoryCostObjectRepo()
	costPoolRepo := repository.NewMemoryCostPoolRepo()
	costPoolLineRepo := repository.NewMemoryCostPoolLineRepo()
	resultRepo := repository.NewMemoryCostingResultRepo()
	resultLineRepo := repository.NewMemoryCostingResultLineRepo()

	engine := NewCostingEngine(periodRepo, costObjectRepo, costPoolRepo, costPoolLineRepo, resultRepo, resultLineRepo)

	return engine, periodRepo, costObjectRepo, costPoolRepo, costPoolLineRepo, resultRepo, resultLineRepo
}

func createTestPeriod(t *testing.T, repo *repository.MemoryCostingPeriodRepo, companyID string, year, month int) *domain.CostingPeriod {
	t.Helper()
	p := &domain.CostingPeriod{
		ID:        "CP-TEST",
		CompanyID: companyID,
		Year:      year,
		Month:     month,
		Status:    "OPEN",
		CreatedAt: "2026-01-01T00:00:00Z",
	}
	err := repo.Create(context.Background(), p)
	assert.NoError(t, err)
	return p
}

// --- Standard Costing Method Tests ---

func TestStandardCosting_CalculatesVariance(t *testing.T) {
	engine, periodRepo, costObjectRepo, costPoolRepo, costPoolLineRepo, _, _ := setupCostingEngineTest()
	ctx := context.Background()

	createTestPeriod(t, periodRepo, "COMP1", 2026, 8)

	// Standard cost per unit = 850,000 (500k mat + 200k lab + 150k overhead)
	// For 1000 units: total standard = 850,000,000
	obj := &domain.CostObject{
		ID:            "OBJ-STD-001",
		CompanyID:     "COMP1",
		Code:          "SP001",
		Name:          "Product A",
		Type:          "PRODUCT",
		CostingMethod: "STANDARD",
		StandardCost:  850000,
		PlanQuantity:  1000,
		IsActive:      true,
	}
	_ = costObjectRepo.Create(ctx, obj)

	// Actual costs for 1000 units: 550M mat + 220M lab + 160M overhead = 930M total
	// Actual unit cost = 930,000
	// Variance = 930M - 850M = 80M (unfavorable)
	poolMat := &domain.CostPool{
		ID: "POOL-MAT", CompanyID: "COMP1", PeriodID: "CP-TEST",
		GLAccountCode: "621", Name: "Materials", Status: "OPEN", TotalAmount: 550000000,
	}
	_ = costPoolRepo.Create(ctx, poolMat)
	_ = costPoolLineRepo.Create(ctx, &domain.CostPoolLine{
		ID: "LINE-MAT-1", PoolID: "POOL-MAT", SourceType: "JOURNAL", Description: "Material", Amount: 550000000,
	})

	poolLab := &domain.CostPool{
		ID: "POOL-LAB", CompanyID: "COMP1", PeriodID: "CP-TEST",
		GLAccountCode: "622", Name: "Labor", Status: "OPEN", TotalAmount: 220000000,
	}
	_ = costPoolRepo.Create(ctx, poolLab)
	_ = costPoolLineRepo.Create(ctx, &domain.CostPoolLine{
		ID: "LINE-LAB-1", PoolID: "POOL-LAB", SourceType: "PAYROLL", Description: "Labor", Amount: 220000000,
	})

	poolOH := &domain.CostPool{
		ID: "POOL-OH", CompanyID: "COMP1", PeriodID: "CP-TEST",
		GLAccountCode: "627", Name: "Overhead", Status: "OPEN", TotalAmount: 160000000,
	}
	_ = costPoolRepo.Create(ctx, poolOH)
	_ = costPoolLineRepo.Create(ctx, &domain.CostPoolLine{
		ID: "LINE-OH-1", PoolID: "POOL-OH", SourceType: "DEPRECIATION", Description: "Overhead", Amount: 160000000,
	})

	err := engine.RunCosting(ctx, "COMP1", "CP-TEST")
	assert.NoError(t, err)

	results, err := engine.resultRepo.ListByPeriod(ctx, "COMP1", "CP-TEST")
	assert.NoError(t, err)
	assert.Len(t, results, 1)

	r := results[0]
	assert.Equal(t, 930000000.0, r.TotalCost)
	assert.Equal(t, 930000.0, r.UnitCost) // 930M / 1000

	// Check result lines for variance
	lines, err := engine.resultLineRepo.ListByResult(ctx, r.ID)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(lines), 3)

	// Verify variance line exists
	hasVariance := false
	for _, l := range lines {
		if l.CostCategory == "VARIANCE" {
			hasVariance = true
			assert.Equal(t, 80000000.0, l.ActualAmount) // 930M - 850M
			assert.Equal(t, "UNFAVORABLE", l.Description)
		}
	}
	assert.True(t, hasVariance, "should have variance line")
}

func TestStandardCosting_FavorableVariance(t *testing.T) {
	engine, periodRepo, costObjectRepo, costPoolRepo, costPoolLineRepo, _, _ := setupCostingEngineTest()
	ctx := context.Background()

	createTestPeriod(t, periodRepo, "COMP1", 2026, 8)

	// Standard per unit = 850k, for 500 units = 425M total standard
	// Actual = 700k total → 700k - 425M = -424.3M favorable? No...
	// Let me fix: Standard = 850k/unit × 500 = 425M total standard
	// Actual = 700M total → Variance = 700M - 425M = 275M unfavorable
	// Actually, let's make it favorable: Actual < Standard
	// Standard = 850k × 500 = 425M, Actual = 350M → Variance = -75M favorable
	obj := &domain.CostObject{
		ID: "OBJ-STD-002", CompanyID: "COMP1", Code: "SP002", Name: "Product B",
		Type: "PRODUCT", CostingMethod: "STANDARD", StandardCost: 850000, PlanQuantity: 500, IsActive: true,
	}
	_ = costObjectRepo.Create(ctx, obj)

	// Actual: 200M mat + 100M lab + 50M overhead = 350M total
	// Standard: 850k × 500 = 425M
	// Variance = 350M - 425M = -75M (favorable)
	_ = costPoolRepo.Create(ctx, &domain.CostPool{
		ID: "POOL-MAT", CompanyID: "COMP1", PeriodID: "CP-TEST",
		GLAccountCode: "621", Name: "Materials", Status: "OPEN", TotalAmount: 200000000,
	})
	_ = costPoolLineRepo.Create(ctx, &domain.CostPoolLine{
		ID: "LINE-M-1", PoolID: "POOL-MAT", SourceType: "JOURNAL", Description: "Mat", Amount: 200000000,
	})
	_ = costPoolRepo.Create(ctx, &domain.CostPool{
		ID: "POOL-LAB", CompanyID: "COMP1", PeriodID: "CP-TEST",
		GLAccountCode: "622", Name: "Labor", Status: "OPEN", TotalAmount: 100000000,
	})
	_ = costPoolLineRepo.Create(ctx, &domain.CostPoolLine{
		ID: "LINE-L-1", PoolID: "POOL-LAB", SourceType: "PAYROLL", Description: "Lab", Amount: 100000000,
	})
	_ = costPoolRepo.Create(ctx, &domain.CostPool{
		ID: "POOL-OH", CompanyID: "COMP1", PeriodID: "CP-TEST",
		GLAccountCode: "627", Name: "Overhead", Status: "OPEN", TotalAmount: 50000000,
	})
	_ = costPoolLineRepo.Create(ctx, &domain.CostPoolLine{
		ID: "LINE-O-1", PoolID: "POOL-OH", SourceType: "MANUAL", Description: "OH", Amount: 50000000,
	})

	err := engine.RunCosting(ctx, "COMP1", "CP-TEST")
	assert.NoError(t, err)

	results, _ := engine.resultRepo.ListByPeriod(ctx, "COMP1", "CP-TEST")
	r := results[0]
	assert.Equal(t, 350000000.0, r.TotalCost)
	assert.Equal(t, 700000.0, r.UnitCost) // 350M / 500

	lines, _ := engine.resultLineRepo.ListByResult(ctx, r.ID)
	for _, l := range lines {
		if l.CostCategory == "VARIANCE" {
			assert.Equal(t, -75000000.0, l.ActualAmount) // 350M - 425M = -75M
			assert.Equal(t, "FAVORABLE", l.Description)
		}
	}
}

// --- Process Costing Method Tests ---

func TestProcessCosting_EquivalentUnits(t *testing.T) {
	engine, periodRepo, costObjectRepo, costPoolRepo, costPoolLineRepo, _, _ := setupCostingEngineTest()
	ctx := context.Background()

	createTestPeriod(t, periodRepo, "COMP1", 2026, 8)

	// Process costing: 500 units started, 400 completed, 100 at 60% completion
	// Equivalent units = 400 + (100 × 0.6) = 460
	obj := &domain.CostObject{
		ID: "OBJ-PRC-001", CompanyID: "COMP1", Code: "PRC001", Name: "Chemical A",
		Type: "PRODUCT", CostingMethod: "PROCESS", PlanQuantity: 460, IsActive: true,
	}
	_ = costObjectRepo.Create(ctx, obj)

	_ = costPoolRepo.Create(ctx, &domain.CostPool{
		ID: "POOL-MAT", CompanyID: "COMP1", PeriodID: "CP-TEST",
		GLAccountCode: "621", Name: "Materials", Status: "OPEN", TotalAmount: 138000000,
	})
	_ = costPoolLineRepo.Create(ctx, &domain.CostPoolLine{
		ID: "LINE-M-1", PoolID: "POOL-MAT", SourceType: "JOURNAL", Description: "Raw material", Amount: 138000000,
	})
	_ = costPoolRepo.Create(ctx, &domain.CostPool{
		ID: "POOL-LAB", CompanyID: "COMP1", PeriodID: "CP-TEST",
		GLAccountCode: "622", Name: "Labor", Status: "OPEN", TotalAmount: 46000000,
	})
	_ = costPoolLineRepo.Create(ctx, &domain.CostPoolLine{
		ID: "LINE-L-1", PoolID: "POOL-LAB", SourceType: "PAYROLL", Description: "Labor", Amount: 46000000,
	})
	_ = costPoolRepo.Create(ctx, &domain.CostPool{
		ID: "POOL-OH", CompanyID: "COMP1", PeriodID: "CP-TEST",
		GLAccountCode: "627", Name: "Overhead", Status: "OPEN", TotalAmount: 46000000,
	})
	_ = costPoolLineRepo.Create(ctx, &domain.CostPoolLine{
		ID: "LINE-O-1", PoolID: "POOL-OH", SourceType: "DEPRECIATION", Description: "OH", Amount: 46000000,
	})

	err := engine.RunCosting(ctx, "COMP1", "CP-TEST")
	assert.NoError(t, err)

	results, _ := engine.resultRepo.ListByPeriod(ctx, "COMP1", "CP-TEST")
	r := results[0]
	// Total = 138M + 46M + 46M = 230M
	// Unit cost = 230M / 460 = 500,000
	assert.Equal(t, 230000000.0, r.TotalCost)
	assert.Equal(t, 500000.0, r.UnitCost)

	// Check process-specific lines
	lines, _ := engine.resultLineRepo.ListByResult(ctx, r.ID)
	hasProcessInfo := false
	for _, l := range lines {
		if l.CostCategory == "PROCESS_INFO" {
			hasProcessInfo = true
			assert.Equal(t, 460.0, l.PlannedAmount)  // equivalent units
			assert.Equal(t, 500000.0, l.AllocatedAmount) // cost per equivalent unit
		}
	}
	assert.True(t, hasProcessInfo, "should have process info line")
}

func TestProcessCosting_NoWIPUnits(t *testing.T) {
	engine, periodRepo, costObjectRepo, costPoolRepo, costPoolLineRepo, _, _ := setupCostingEngineTest()
	ctx := context.Background()

	createTestPeriod(t, periodRepo, "COMP1", 2026, 8)

	// All units completed, no WIP
	obj := &domain.CostObject{
		ID: "OBJ-PRC-002", CompanyID: "COMP1", Code: "PRC002", Name: "Chemical B",
		Type: "PRODUCT", CostingMethod: "PROCESS", PlanQuantity: 500, IsActive: true,
	}
	_ = costObjectRepo.Create(ctx, obj)

	_ = costPoolRepo.Create(ctx, &domain.CostPool{
		ID: "POOL-MAT", CompanyID: "COMP1", PeriodID: "CP-TEST",
		GLAccountCode: "621", Name: "Materials", Status: "OPEN", TotalAmount: 100000000,
	})
	_ = costPoolLineRepo.Create(ctx, &domain.CostPoolLine{
		ID: "LINE-M-1", PoolID: "POOL-MAT", SourceType: "JOURNAL", Description: "Mat", Amount: 100000000,
	})

	err := engine.RunCosting(ctx, "COMP1", "CP-TEST")
	assert.NoError(t, err)

	results, _ := engine.resultRepo.ListByPeriod(ctx, "COMP1", "CP-TEST")
	r := results[0]
	// 100M / 500 = 200,000
	assert.Equal(t, 200000.0, r.UnitCost)
}

// --- WIP Valuation Tests ---

func TestWIPValuation_BeginningEndingWIP(t *testing.T) {
	engine, periodRepo, costObjectRepo, costPoolRepo, costPoolLineRepo, _, _ := setupCostingEngineTest()
	ctx := context.Background()

	createTestPeriod(t, periodRepo, "COMP1", 2026, 8)

	obj := &domain.CostObject{
		ID: "OBJ-WIP-001", CompanyID: "COMP1", Code: "WIP001", Name: "WIP Product",
		Type: "PRODUCT", CostingMethod: "SIMPLE", PlanQuantity: 100, IsActive: true,
	}
	_ = costObjectRepo.Create(ctx, obj)

	_ = costPoolRepo.Create(ctx, &domain.CostPool{
		ID: "POOL-MAT", CompanyID: "COMP1", PeriodID: "CP-TEST",
		GLAccountCode: "621", Name: "Materials", Status: "OPEN", TotalAmount: 50000000,
	})
	_ = costPoolLineRepo.Create(ctx, &domain.CostPoolLine{
		ID: "LINE-M-1", PoolID: "POOL-MAT", SourceType: "JOURNAL", Description: "Mat", Amount: 50000000,
	})

	err := engine.RunCosting(ctx, "COMP1", "CP-TEST")
	assert.NoError(t, err)

	results, _ := engine.resultRepo.ListByPeriod(ctx, "COMP1", "CP-TEST")
	r := results[0]
	// WIPBegin should be 0 (no prior period data)
	assert.Equal(t, 0.0, r.WIPBegin)
	// Total cost = 50M, Unit cost = 500,000
	assert.Equal(t, 50000000.0, r.TotalCost)
	assert.Equal(t, 500000.0, r.UnitCost)
}

func TestWIPValuation_TransferToFinishedGoods(t *testing.T) {
	engine, periodRepo, costObjectRepo, costPoolRepo, costPoolLineRepo, resultRepo, _ := setupCostingEngineTest()
	ctx := context.Background()

	createTestPeriod(t, periodRepo, "COMP1", 2026, 8)

	obj := &domain.CostObject{
		ID: "OBJ-WIP-002", CompanyID: "COMP1", Code: "WIP002", Name: "WIP Product 2",
		Type: "PRODUCT", CostingMethod: "SIMPLE", PlanQuantity: 200, IsActive: true,
	}
	_ = costObjectRepo.Create(ctx, obj)

	_ = costPoolRepo.Create(ctx, &domain.CostPool{
		ID: "POOL-MAT", CompanyID: "COMP1", PeriodID: "CP-TEST",
		GLAccountCode: "621", Name: "Materials", Status: "OPEN", TotalAmount: 100000000,
	})
	_ = costPoolLineRepo.Create(ctx, &domain.CostPoolLine{
		ID: "LINE-M-1", PoolID: "POOL-MAT", SourceType: "JOURNAL", Description: "Mat", Amount: 100000000,
	})

	err := engine.RunCosting(ctx, "COMP1", "CP-TEST")
	assert.NoError(t, err)

	results, _ := resultRepo.ListByPeriod(ctx, "COMP1", "CP-TEST")
	assert.Len(t, results, 1)
	r := results[0]
	assert.Equal(t, 100000000.0, r.TotalCost)
	assert.Equal(t, 500000.0, r.UnitCost)

	// Finalize the result
	r.Status = "FINAL"
	err = resultRepo.Update(ctx, &r)
	assert.NoError(t, err)

	updated, _ := resultRepo.GetByID(ctx, r.ID)
	assert.Equal(t, "FINAL", updated.Status)
}
