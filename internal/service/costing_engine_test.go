package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"gotax/internal/domain"
	"gotax/internal/repository"
)

func setupCostingEngineTest(t *testing.T) (*CostingEngine, *repository.MemoryCostingPeriodRepo, *repository.MemoryCostObjectRepo, *repository.MemoryCostPoolRepo, *repository.MemoryCostPoolLineRepo, *repository.MemoryCostingResultRepo, *repository.MemoryCostingResultLineRepo) {
	t.Helper()
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
	engine, periodRepo, costObjectRepo, costPoolRepo, costPoolLineRepo, _, _ := setupCostingEngineTest(t)
	ctx := context.Background()

	createTestPeriod(t, periodRepo, "COMP1", 2026, 8)

	// Standard cost per unit = 850,000 (500k mat + 200k lab + 150k overhead)
	// For 1000 units: total standard = 850,000,000
	obj := &domain.CostObject{
		ID:               "OBJ-STD-001",
		CompanyID:        "COMP1",
		Code:             "SP001",
		Name:             "Product A",
		Type:             "PRODUCT",
		CostingMethod:    "STANDARD",
		StandardMaterial: 500000,
		StandardLabor:    200000,
		StandardOverhead: 150000,
		PlanQuantity:     1000,
		IsActive:         true,
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
	engine, periodRepo, costObjectRepo, costPoolRepo, costPoolLineRepo, _, _ := setupCostingEngineTest(t)
	ctx := context.Background()

	createTestPeriod(t, periodRepo, "COMP1", 2026, 8)

	// Standard per unit = 850k, for 500 units = 425M total standard
	// Actual = 350M total → Variance = -75M (favorable)
	obj := &domain.CostObject{
		ID:               "OBJ-STD-002", CompanyID: "COMP1", Code: "SP002", Name: "Product B",
		Type: "PRODUCT", CostingMethod: "STANDARD", StandardMaterial: 500000, StandardLabor: 200000, StandardOverhead: 150000, PlanQuantity: 500, IsActive: true,
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
	engine, periodRepo, costObjectRepo, costPoolRepo, costPoolLineRepo, _, _ := setupCostingEngineTest(t)
	ctx := context.Background()

	createTestPeriod(t, periodRepo, "COMP1", 2026, 8)

	// Process costing: 500 units started, 400 completed, 100 at 60% completion
	// Equivalent units = 400 + (100 × 0.6) = 460
	obj := &domain.CostObject{
		ID: "OBJ-PRC-001", CompanyID: "COMP1", Code: "PRC001", Name: "Chemical A",
		Type: "PRODUCT", CostingMethod: "PROCESS", CompletedUnits: 400, WIPUnits: 100, CompletionPct: 60, IsActive: true,
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
	engine, periodRepo, costObjectRepo, costPoolRepo, costPoolLineRepo, _, _ := setupCostingEngineTest(t)
	ctx := context.Background()

	createTestPeriod(t, periodRepo, "COMP1", 2026, 8)

	// All units completed, no WIP
	obj := &domain.CostObject{
		ID: "OBJ-PRC-002", CompanyID: "COMP1", Code: "PRC002", Name: "Chemical B",
		Type: "PRODUCT", CostingMethod: "PROCESS", CompletedUnits: 500, WIPUnits: 0, CompletionPct: 0, IsActive: true,
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
	engine, periodRepo, costObjectRepo, costPoolRepo, costPoolLineRepo, _, _ := setupCostingEngineTest(t)
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
	engine, periodRepo, costObjectRepo, costPoolRepo, costPoolLineRepo, resultRepo, _ := setupCostingEngineTest(t)
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

// --- By-product Exclusion Method Tests ---

func TestByProductCosting_DeductsByProductValue(t *testing.T) {
	engine, periodRepo, costObjectRepo, costPoolRepo, costPoolLineRepo, _, _ := setupCostingEngineTest(t)
	ctx := context.Background()

	createTestPeriod(t, periodRepo, "COMP1", 2026, 8)

	// Main product: lumber with by-product NRV deduction
	// StandardCost holds the by-product NRV to deduct
	mainObj := &domain.CostObject{
		ID: "OBJ-BYPROD-001", CompanyID: "COMP1", Code: "LUMBER001", Name: "Lumber",
		Type: "PRODUCT", CostingMethod: "BY_PRODUCT", StandardCost: 50000000, PlanQuantity: 1000, IsActive: true,
	}
	_ = costObjectRepo.Create(ctx, mainObj)

	// Total production cost = 1B
	_ = costPoolRepo.Create(ctx, &domain.CostPool{
		ID: "POOL-MAT", CompanyID: "COMP1", PeriodID: "CP-TEST",
		GLAccountCode: "621", Name: "Materials", Status: "OPEN", TotalAmount: 600000000,
	})
	_ = costPoolLineRepo.Create(ctx, &domain.CostPoolLine{
		ID: "LINE-M-1", PoolID: "POOL-MAT", SourceType: "JOURNAL", Description: "Logs", Amount: 600000000,
	})
	_ = costPoolRepo.Create(ctx, &domain.CostPool{
		ID: "POOL-LAB", CompanyID: "COMP1", PeriodID: "CP-TEST",
		GLAccountCode: "622", Name: "Labor", Status: "OPEN", TotalAmount: 250000000,
	})
	_ = costPoolLineRepo.Create(ctx, &domain.CostPoolLine{
		ID: "LINE-L-1", PoolID: "POOL-LAB", SourceType: "PAYROLL", Description: "Sawmill", Amount: 250000000,
	})
	_ = costPoolRepo.Create(ctx, &domain.CostPool{
		ID: "POOL-OH", CompanyID: "COMP1", PeriodID: "CP-TEST",
		GLAccountCode: "627", Name: "Overhead", Status: "OPEN", TotalAmount: 150000000,
	})
	_ = costPoolLineRepo.Create(ctx, &domain.CostPoolLine{
		ID: "LINE-O-1", PoolID: "POOL-OH", SourceType: "DEPRECIATION", Description: "Machines", Amount: 150000000,
	})

	err := engine.RunCosting(ctx, "COMP1", "CP-TEST")
	assert.NoError(t, err)

	results, _ := engine.resultRepo.ListByPeriod(ctx, "COMP1", "CP-TEST")
	var lumberResult domain.CostingResult
	for _, r := range results {
		if r.CostObjectID == "OBJ-BYPROD-001" {
			lumberResult = r
			break
		}
	}

	// Total cost = 1B, By-product NRV = 50M
	// Main product cost = 1B - 50M = 950M
	// Unit cost = 950M / 1000 = 950,000
	assert.Equal(t, 950000000.0, lumberResult.TotalCost)
	assert.Equal(t, 950000.0, lumberResult.UnitCost)

	// Check by-product deduction line
	lines, _ := engine.resultLineRepo.ListByResult(ctx, lumberResult.ID)
	hasByProdDeduction := false
	for _, l := range lines {
		if l.CostCategory == "BY_PRODUCT_DEDUCTION" {
			hasByProdDeduction = true
			assert.Equal(t, -50000000.0, l.ActualAmount) // negative = deduction
		}
	}
	assert.True(t, hasByProdDeduction, "should have by-product deduction line")
}

// --- Report Tests ---

func TestCostReport_CostCalculationSheet(t *testing.T) {
	engine, periodRepo, costObjectRepo, costPoolRepo, costPoolLineRepo, resultRepo, resultLineRepo := setupCostingEngineTest(t)
	ctx := context.Background()

	createTestPeriod(t, periodRepo, "COMP1", 2026, 8)

	obj := &domain.CostObject{
		ID: "OBJ-RPT-001", CompanyID: "COMP1", Code: "RPT001", Name: "Report Product",
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
	_ = costPoolRepo.Create(ctx, &domain.CostPool{
		ID: "POOL-LAB", CompanyID: "COMP1", PeriodID: "CP-TEST",
		GLAccountCode: "622", Name: "Labor", Status: "OPEN", TotalAmount: 30000000,
	})
	_ = costPoolLineRepo.Create(ctx, &domain.CostPoolLine{
		ID: "LINE-L-1", PoolID: "POOL-LAB", SourceType: "PAYROLL", Description: "Labor", Amount: 30000000,
	})
	_ = costPoolRepo.Create(ctx, &domain.CostPool{
		ID: "POOL-OH", CompanyID: "COMP1", PeriodID: "CP-TEST",
		GLAccountCode: "627", Name: "Overhead", Status: "OPEN", TotalAmount: 20000000,
	})
	_ = costPoolLineRepo.Create(ctx, &domain.CostPoolLine{
		ID: "LINE-O-1", PoolID: "POOL-OH", SourceType: "MANUAL", Description: "OH", Amount: 20000000,
	})

	err := engine.RunCosting(ctx, "COMP1", "CP-TEST")
	assert.NoError(t, err)

	// Test report generation
	svc := NewCostReportService(resultRepo, resultLineRepo, costObjectRepo)
	report, err := svc.GetCostCalculationSheet(ctx, "COMP1", "CP-TEST", "OBJ-RPT-001")
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, "RPT001", report.ObjectCode)
	assert.Equal(t, 100000000.0, report.TotalCost)
	assert.Equal(t, 1000000.0, report.UnitCost)
	assert.Len(t, report.Lines, 3) // mat, lab, oh
}

func TestCostReport_CostSummary(t *testing.T) {
	engine, periodRepo, costObjectRepo, costPoolRepo, costPoolLineRepo, resultRepo, resultLineRepo := setupCostingEngineTest(t)
	ctx := context.Background()

	createTestPeriod(t, periodRepo, "COMP1", 2026, 8)

	obj1 := &domain.CostObject{
		ID: "OBJ-SUM-001", CompanyID: "COMP1", Code: "SUM001", Name: "Product A",
		Type: "PRODUCT", CostingMethod: "SIMPLE", PlanQuantity: 100, IsActive: true,
	}
	obj2 := &domain.CostObject{
		ID: "OBJ-SUM-002", CompanyID: "COMP1", Code: "SUM002", Name: "Product B",
		Type: "PRODUCT", CostingMethod: "SIMPLE", PlanQuantity: 200, IsActive: true,
	}
	_ = costObjectRepo.Create(ctx, obj1)
	_ = costObjectRepo.Create(ctx, obj2)

	_ = costPoolRepo.Create(ctx, &domain.CostPool{
		ID: "POOL-MAT", CompanyID: "COMP1", PeriodID: "CP-TEST",
		GLAccountCode: "621", Name: "Materials", Status: "OPEN", TotalAmount: 100000000,
	})
	_ = costPoolLineRepo.Create(ctx, &domain.CostPoolLine{
		ID: "LINE-M-1", PoolID: "POOL-MAT", SourceType: "JOURNAL", Description: "Mat", Amount: 100000000,
	})

	err := engine.RunCosting(ctx, "COMP1", "CP-TEST")
	assert.NoError(t, err)

	svc := NewCostReportService(resultRepo, resultLineRepo, costObjectRepo)
	summary, err := svc.GetCostSummary(ctx, "COMP1", "CP-TEST")
	assert.NoError(t, err)
	assert.NotNil(t, summary)
	assert.Len(t, summary.Items, 2)
	assert.Equal(t, 200000000.0, summary.TotalAllObjects) // 100M × 2 objects
}

func TestCostReport_WIPValuation(t *testing.T) {
	engine, periodRepo, costObjectRepo, costPoolRepo, costPoolLineRepo, resultRepo, resultLineRepo := setupCostingEngineTest(t)
	ctx := context.Background()

	createTestPeriod(t, periodRepo, "COMP1", 2026, 8)

	obj := &domain.CostObject{
		ID: "OBJ-WIP-RPT", CompanyID: "COMP1", Code: "WIPRPT01", Name: "WIP Product",
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

	svc := NewCostReportService(resultRepo, resultLineRepo, costObjectRepo)
	wipReport, err := svc.GetWIPValuation(ctx, "COMP1", "CP-TEST")
	assert.NoError(t, err)
	assert.NotNil(t, wipReport)
	assert.Len(t, wipReport.Items, 1)
	assert.Equal(t, 50000000.0, wipReport.TotalWIP)
}

func TestCostReport_VarianceAnalysis(t *testing.T) {
	engine, periodRepo, costObjectRepo, costPoolRepo, costPoolLineRepo, resultRepo, resultLineRepo := setupCostingEngineTest(t)
	ctx := context.Background()

	createTestPeriod(t, periodRepo, "COMP1", 2026, 8)

	obj := &domain.CostObject{
		ID: "OBJ-VAR-001", CompanyID: "COMP1", Code: "VAR001", Name: "Variance Product",
		Type: "PRODUCT", CostingMethod: "STANDARD",
		StandardMaterial: 500000, StandardLabor: 200000, StandardOverhead: 150000,
		PlanQuantity: 1000, IsActive: true,
	}
	_ = costObjectRepo.Create(ctx, obj)

	// Actual: 550M + 220M + 160M = 930M
	_ = costPoolRepo.Create(ctx, &domain.CostPool{
		ID: "POOL-MAT", CompanyID: "COMP1", PeriodID: "CP-TEST",
		GLAccountCode: "621", Name: "Materials", Status: "OPEN", TotalAmount: 550000000,
	})
	_ = costPoolLineRepo.Create(ctx, &domain.CostPoolLine{
		ID: "LINE-M-1", PoolID: "POOL-MAT", SourceType: "JOURNAL", Description: "Mat", Amount: 550000000,
	})
	_ = costPoolRepo.Create(ctx, &domain.CostPool{
		ID: "POOL-LAB", CompanyID: "COMP1", PeriodID: "CP-TEST",
		GLAccountCode: "622", Name: "Labor", Status: "OPEN", TotalAmount: 220000000,
	})
	_ = costPoolLineRepo.Create(ctx, &domain.CostPoolLine{
		ID: "LINE-L-1", PoolID: "POOL-LAB", SourceType: "PAYROLL", Description: "Labor", Amount: 220000000,
	})
	_ = costPoolRepo.Create(ctx, &domain.CostPool{
		ID: "POOL-OH", CompanyID: "COMP1", PeriodID: "CP-TEST",
		GLAccountCode: "627", Name: "Overhead", Status: "OPEN", TotalAmount: 160000000,
	})
	_ = costPoolLineRepo.Create(ctx, &domain.CostPoolLine{
		ID: "LINE-O-1", PoolID: "POOL-OH", SourceType: "DEPRECIATION", Description: "OH", Amount: 160000000,
	})

	err := engine.RunCosting(ctx, "COMP1", "CP-TEST")
	assert.NoError(t, err)

	svc := NewCostReportService(resultRepo, resultLineRepo, costObjectRepo)
	varianceReport, err := svc.GetVarianceAnalysis(ctx, "COMP1", "CP-TEST")
	assert.NoError(t, err)
	assert.NotNil(t, varianceReport)
	assert.Len(t, varianceReport.Items, 1)

	item := varianceReport.Items[0]
	assert.Equal(t, 930000000.0, item.ActualCost)
	assert.Equal(t, 850000000.0, item.StandardCost)
	assert.Equal(t, 80000000.0, item.Variance)
	assert.Equal(t, "UNFAVORABLE", item.VarianceType)
}
