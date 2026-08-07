package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"gotax/internal/domain"
	"gotax/internal/repository"
)

func setupCostingJETest(t *testing.T) (*CostingJEService, *repository.MemoryCostingPeriodRepo, *repository.MemoryCostObjectRepo, *repository.MemoryCostPoolRepo, *repository.MemoryCostPoolLineRepo, *repository.MemoryCostingResultRepo, *repository.MemoryCostingResultLineRepo) {
	t.Helper()
	periodRepo := repository.NewMemoryCostingPeriodRepo()
	costObjectRepo := repository.NewMemoryCostObjectRepo()
	costPoolRepo := repository.NewMemoryCostPoolRepo()
	costPoolLineRepo := repository.NewMemoryCostPoolLineRepo()
	resultRepo := repository.NewMemoryCostingResultRepo()
	resultLineRepo := repository.NewMemoryCostingResultLineRepo()

	// Use a simple journal entry creator for tests
	jeCreator := &testJECreator{entries: make([]*domain.JournalEntry, 0)}

	svc := NewCostingJEService(periodRepo, costObjectRepo, costPoolRepo, costPoolLineRepo, resultRepo, resultLineRepo, jeCreator)

	return svc, periodRepo, costObjectRepo, costPoolRepo, costPoolLineRepo, resultRepo, resultLineRepo
}

type testJECreator struct {
	entries []*domain.JournalEntry
}

func (c *testJECreator) CreateEntry(ctx context.Context, entry *domain.JournalEntry, userID string) error {
	c.entries = append(c.entries, entry)
	return nil
}

func createTestCostingPeriod(t *testing.T, repo *repository.MemoryCostingPeriodRepo, companyID string, year, month int) string {
	t.Helper()
	id := "CP-TEST"
	p := &domain.CostingPeriod{
		ID: id, CompanyID: companyID, Year: year, Month: month, Status: "OPEN", CreatedAt: "2026-01-01T00:00:00Z",
	}
	_ = repo.Create(context.Background(), p)
	return id
}

func createTestCostObject(t *testing.T, repo *repository.MemoryCostObjectRepo, id, companyID, code, method string) {
	t.Helper()
	obj := &domain.CostObject{
		ID: id, CompanyID: companyID, Code: code, Name: "Test " + code,
		Type: "PRODUCT", CostingMethod: domain.CostingMethod(method), PlanQuantity: 100, IsActive: true,
	}
	_ = repo.Create(context.Background(), obj)
}

func createTestCostPool(t *testing.T, repo *repository.MemoryCostPoolRepo, id, companyID, periodID, glAcct string, amount float64) {
	t.Helper()
	pool := &domain.CostPool{
		ID: id, CompanyID: companyID, PeriodID: periodID, GLAccountCode: glAcct,
		Name: "Pool " + glAcct, Status: "OPEN", TotalAmount: amount,
	}
	_ = repo.Create(context.Background(), pool)
}

func createTestCostPoolLine(t *testing.T, repo *repository.MemoryCostPoolLineRepo, id, poolID string, amount float64) {
	t.Helper()
	line := &domain.CostPoolLine{
		ID: id, PoolID: poolID, SourceType: "JOURNAL", Description: "Test line", Amount: amount,
	}
	_ = repo.Create(context.Background(), line)
}

// --- Task 34: Cost Pool → GL Entry Tests ---

func TestCostingJE_CostPoolToGLEntry(t *testing.T) {
	svc, periodRepo, costObjectRepo, costPoolRepo, costPoolLineRepo, _, _ := setupCostingJETest(t)
	ctx := context.Background()

	periodID := createTestCostingPeriod(t, periodRepo, "COMP1", 2026, 8)
	createTestCostObject(t, costObjectRepo, "OBJ-001", "COMP1", "SP001", "SIMPLE")
	createTestCostPool(t, costPoolRepo, "POOL-MAT", "COMP1", periodID, "621", 50000000)
	createTestCostPoolLine(t, costPoolLineRepo, "LINE-001", "POOL-MAT", 50000000)
	createTestCostPool(t, costPoolRepo, "POOL-LAB", "COMP1", periodID, "622", 30000000)
	createTestCostPoolLine(t, costPoolLineRepo, "LINE-002", "POOL-LAB", 30000000)
	createTestCostPool(t, costPoolRepo, "POOL-OH", "COMP1", periodID, "627", 20000000)
	createTestCostPoolLine(t, costPoolLineRepo, "LINE-003", "POOL-OH", 20000000)

	err := svc.GenerateCostPoolEntries(ctx, "COMP1", periodID)
	assert.NoError(t, err)

	// Should create 3 journal entries (one per pool)
	entries := svc.jeCreator.(*testJECreator).entries
	assert.Len(t, entries, 3)

	// Find the 621 entry
	var matEntry *domain.JournalEntry
	for _, e := range entries {
		for _, l := range e.Lines {
			if l.AccountCode == "621" && l.CreditAmount > 0 {
				matEntry = e
				break
			}
		}
	}
	assert.NotNil(t, matEntry)

	// Check entry: Dr 154, Cr 621
	assert.Len(t, matEntry.Lines, 2)
	// Debit line: 154 (WIP)
	assert.Equal(t, "154", matEntry.Lines[0].AccountCode)
	assert.Equal(t, 50000000.0, matEntry.Lines[0].DebitAmount)
	assert.Equal(t, 0.0, matEntry.Lines[0].CreditAmount)
	// Credit line: 621 (Direct materials)
	assert.Equal(t, "621", matEntry.Lines[1].AccountCode)
	assert.Equal(t, 0.0, matEntry.Lines[1].DebitAmount)
	assert.Equal(t, 50000000.0, matEntry.Lines[1].CreditAmount)
}

func TestCostingJE_CostPoolToGLEntry_MultiplePools(t *testing.T) {
	svc, periodRepo, costObjectRepo, costPoolRepo, costPoolLineRepo, _, _ := setupCostingJETest(t)
	ctx := context.Background()

	periodID := createTestCostingPeriod(t, periodRepo, "COMP1", 2026, 8)
	createTestCostObject(t, costObjectRepo, "OBJ-001", "COMP1", "SP001", "SIMPLE")
	createTestCostPool(t, costPoolRepo, "POOL-MAT", "COMP1", periodID, "621", 100000000)
	createTestCostPoolLine(t, costPoolLineRepo, "LINE-001", "POOL-MAT", 100000000)
	createTestCostPool(t, costPoolRepo, "POOL-LAB", "COMP1", periodID, "622", 50000000)
	createTestCostPoolLine(t, costPoolLineRepo, "LINE-002", "POOL-LAB", 50000000)

	err := svc.GenerateCostPoolEntries(ctx, "COMP1", periodID)
	assert.NoError(t, err)

	// Should create 2 journal entries
	assert.Len(t, svc.jeCreator.(*testJECreator).entries, 2)

	// Verify entry numbers are unique
	assert.NotEqual(t, svc.jeCreator.(*testJECreator).entries[0].EntryNumber,
		svc.jeCreator.(*testJECreator).entries[1].EntryNumber)
}

// --- Task 35: WIP → Finished Goods Transfer Tests ---

func TestCostingJE_WIPTransfer(t *testing.T) {
	svc, periodRepo, costObjectRepo, costPoolRepo, costPoolLineRepo, resultRepo, resultLineRepo := setupCostingJETest(t)
	ctx := context.Background()

	periodID := createTestCostingPeriod(t, periodRepo, "COMP1", 2026, 8)
	createTestCostObject(t, costObjectRepo, "OBJ-001", "COMP1", "SP001", "SIMPLE")
	createTestCostPool(t, costPoolRepo, "POOL-MAT", "COMP1", periodID, "621", 50000000)
	createTestCostPoolLine(t, costPoolLineRepo, "LINE-001", "POOL-MAT", 50000000)

	// Create a costing result
	result := &domain.CostingResult{
		ID: "CR-001", CompanyID: "COMP1", PeriodID: periodID, CostObjectID: "OBJ-001",
		CostingMethod: "SIMPLE", TotalDirectMat: 50000000, TotalCost: 50000000,
		OutputQuantity: 100, UnitCost: 500000, Status: "FINAL",
	}
	_ = resultRepo.Create(ctx, result)
	_ = resultLineRepo.Create(ctx, &domain.CostingResultLine{
		ID: "CRL-001", ResultID: "CR-001", CostCategory: "DIRECT_MATERIAL",
		GLAccountCode: "621", ActualAmount: 50000000,
	})

	err := svc.GenerateWIPTransferEntry(ctx, "COMP1", periodID, "OBJ-001")
	assert.NoError(t, err)

	// Should create 1 journal entry
	assert.Len(t, svc.jeCreator.(*testJECreator).entries, 1)

	// Check entry: Dr 155 (Finished goods), Cr 154 (WIP)
	entry := svc.jeCreator.(*testJECreator).entries[0]
	assert.Len(t, entry.Lines, 2)
	// Debit: 155
	assert.Equal(t, "155", entry.Lines[0].AccountCode)
	assert.Equal(t, 50000000.0, entry.Lines[0].DebitAmount)
	// Credit: 154
	assert.Equal(t, "154", entry.Lines[1].AccountCode)
	assert.Equal(t, 50000000.0, entry.Lines[1].CreditAmount)
}

func TestCostingJE_WIPTransfer_NotFinal(t *testing.T) {
	svc, periodRepo, costObjectRepo, costPoolRepo, costPoolLineRepo, resultRepo, resultLineRepo := setupCostingJETest(t)
	ctx := context.Background()

	periodID := createTestCostingPeriod(t, periodRepo, "COMP1", 2026, 8)
	createTestCostObject(t, costObjectRepo, "OBJ-001", "COMP1", "SP001", "SIMPLE")
	createTestCostPool(t, costPoolRepo, "POOL-MAT", "COMP1", periodID, "621", 50000000)
	createTestCostPoolLine(t, costPoolLineRepo, "LINE-001", "POOL-MAT", 50000000)

	// Create a DRAFT costing result
	result := &domain.CostingResult{
		ID: "CR-001", CompanyID: "COMP1", PeriodID: periodID, CostObjectID: "OBJ-001",
		CostingMethod: "SIMPLE", TotalCost: 50000000, Status: "DRAFT",
	}
	_ = resultRepo.Create(ctx, result)
	_ = resultLineRepo.Create(ctx, &domain.CostingResultLine{
		ID: "CRL-001", ResultID: "CR-001", CostCategory: "DIRECT_MATERIAL",
		GLAccountCode: "621", ActualAmount: 50000000,
	})

	err := svc.GenerateWIPTransferEntry(ctx, "COMP1", periodID, "OBJ-001")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not finalized")
}

// --- Task 36: Period-End Closing Tests ---

func TestCostingJE_PeriodEndClosing(t *testing.T) {
	svc, periodRepo, costObjectRepo, costPoolRepo, costPoolLineRepo, resultRepo, resultLineRepo := setupCostingJETest(t)
	ctx := context.Background()

	periodID := createTestCostingPeriod(t, periodRepo, "COMP1", 2026, 8)
	createTestCostObject(t, costObjectRepo, "OBJ-001", "COMP1", "SP001", "SIMPLE")
	createTestCostPool(t, costPoolRepo, "POOL-MAT", "COMP1", periodID, "621", 50000000)
	createTestCostPoolLine(t, costPoolLineRepo, "LINE-001", "POOL-MAT", 50000000)

	// Create finalized result
	result := &domain.CostingResult{
		ID: "CR-001", CompanyID: "COMP1", PeriodID: periodID, CostObjectID: "OBJ-001",
		CostingMethod: "SIMPLE", TotalCost: 50000000, OutputQuantity: 100, UnitCost: 500000, Status: "FINAL",
	}
	_ = resultRepo.Create(ctx, result)
	_ = resultLineRepo.Create(ctx, &domain.CostingResultLine{
		ID: "CRL-001", ResultID: "CR-001", CostCategory: "DIRECT_MATERIAL",
		GLAccountCode: "621", ActualAmount: 50000000,
	})

	err := svc.ClosePeriod(ctx, "COMP1", periodID, "admin")
	assert.NoError(t, err)

	// Should create entries for all cost pools + WIP transfer
	// 1 pool entry + 1 WIP transfer = 2 entries
	assert.Len(t, svc.jeCreator.(*testJECreator).entries, 2)

	// Period should be closed
	period, _ := periodRepo.GetByID(ctx, periodID)
	assert.Equal(t, "CLOSED", period.Status)
	assert.Equal(t, "admin", period.ClosedBy)
}

func TestCostingJE_PeriodEndClosing_AlreadyClosed(t *testing.T) {
	svc, periodRepo, _, _, _, _, _ := setupCostingJETest(t)
	ctx := context.Background()

	periodID := createTestCostingPeriod(t, periodRepo, "COMP1", 2026, 8)
	// Close it first
	period, _ := periodRepo.GetByID(ctx, periodID)
	period.Status = "CLOSED"
	_ = periodRepo.Update(ctx, period)

	err := svc.ClosePeriod(ctx, "COMP1", periodID, "admin")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already closed")
}

// --- Task 37: Opening Balance Carry-Forward Tests ---

func TestCostingJE_OpeningBalanceCarryForward(t *testing.T) {
	svc, periodRepo, _, _, _, resultRepo, _ := setupCostingJETest(t)
	ctx := context.Background()

	// Create previous period with different ID
	prevPeriodID := "CP-PREV"
	p := &domain.CostingPeriod{
		ID: prevPeriodID, CompanyID: "COMP1", Year: 2026, Month: 7, Status: "CLOSED", CreatedAt: "2026-01-01T00:00:00Z",
	}
	_ = periodRepo.Create(ctx, p)

	// Create previous period result
	prevResult := &domain.CostingResult{
		ID: "CR-PREV", CompanyID: "COMP1", PeriodID: prevPeriodID, CostObjectID: "OBJ-001",
		CostingMethod: "SIMPLE", TotalCost: 50000000, OutputQuantity: 100, UnitCost: 500000,
		WIPEnd: 20, Status: "FINAL",
	}
	_ = resultRepo.Create(ctx, prevResult)

	// Create current period
	currPeriodID := "CP-CURR"
	p2 := &domain.CostingPeriod{
		ID: currPeriodID, CompanyID: "COMP1", Year: 2026, Month: 8, Status: "OPEN", CreatedAt: "2026-01-01T00:00:00Z",
	}
	_ = periodRepo.Create(ctx, p2)

	err := svc.CarryForwardOpeningBalance(ctx, "COMP1", currPeriodID)
	assert.NoError(t, err)

	// Should create a new result with WIPBegin = previous WIPEnd
	results, _ := resultRepo.ListByPeriod(ctx, "COMP1", currPeriodID)
	assert.Len(t, results, 1)
	assert.Equal(t, 20.0, results[0].WIPBegin) // Carried forward from previous period
}

func TestCostingJE_OpeningBalanceCarryForward_NoPrevious(t *testing.T) {
	svc, periodRepo, _, _, _, resultRepo, _ := setupCostingJETest(t)
	ctx := context.Background()

	// Only create current period (no previous)
	currPeriodID := createTestCostingPeriod(t, periodRepo, "COMP1", 2026, 8)

	err := svc.CarryForwardOpeningBalance(ctx, "COMP1", currPeriodID)
	assert.NoError(t, err)

	// Should not create any results (no previous data)
	results, _ := resultRepo.ListByPeriod(ctx, "COMP1", currPeriodID)
	assert.Len(t, results, 0)
}

func TestCostingJE_COGSEntry(t *testing.T) {
	svc, periodRepo, _, _, _, resultRepo, _ := setupCostingJETest(t)
	ctx := context.Background()

	periodID := createTestCostingPeriod(t, periodRepo, "COMP1", 2026, 8)
	result := &domain.CostingResult{
		ID: "CR-001", CompanyID: "COMP1", PeriodID: periodID, CostObjectID: "OBJ-001",
		CostingMethod: "SIMPLE", TotalCost: 50000000, OutputQuantity: 100, UnitCost: 500000,
		Status: "FINAL",
	}
	_ = resultRepo.Create(ctx, result)

	err := svc.GenerateCOGSEntry(ctx, "COMP1", periodID, "OBJ-001", 50)
	assert.NoError(t, err)

	entries := svc.jeCreator.(*testJECreator).entries
	assert.Len(t, entries, 1)

	entry := entries[0]
	assert.Equal(t, "632", entry.Lines[0].AccountCode)
	assert.Equal(t, 25000000.0, entry.Lines[0].DebitAmount)
	assert.Equal(t, "155", entry.Lines[1].AccountCode)
	assert.Equal(t, 25000000.0, entry.Lines[1].CreditAmount)
}

func TestCostingJE_COGSEntry_NotFinal(t *testing.T) {
	svc, periodRepo, _, _, _, resultRepo, _ := setupCostingJETest(t)
	ctx := context.Background()

	periodID := createTestCostingPeriod(t, periodRepo, "COMP1", 2026, 8)
	result := &domain.CostingResult{
		ID: "CR-001", CompanyID: "COMP1", PeriodID: periodID, CostObjectID: "OBJ-001",
		CostingMethod: "SIMPLE", TotalCost: 50000000, OutputQuantity: 100, UnitCost: 500000,
		Status: "DRAFT",
	}
	_ = resultRepo.Create(ctx, result)

	err := svc.GenerateCOGSEntry(ctx, "COMP1", periodID, "OBJ-001", 50)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not finalized")
}

func TestCostingJE_CollectMaterialCosts(t *testing.T) {
	svc, periodRepo, _, costPoolRepo, costPoolLineRepo, _, _ := setupCostingJETest(t)
	ctx := context.Background()

	periodID := createTestCostingPeriod(t, periodRepo, "COMP1", 2026, 8)

	lines := []CostPoolLineInput{
		{SourceID: "WH-001", Description: "Raw material A", Amount: 30000000},
		{SourceID: "WH-002", Description: "Raw material B", Amount: 20000000},
	}

	err := svc.CollectMaterialCosts(ctx, "COMP1", periodID, lines)
	assert.NoError(t, err)

	pools, _ := costPoolRepo.ListByPeriod(ctx, "COMP1", periodID)
	assert.Len(t, pools, 1)
	assert.Equal(t, "621", pools[0].GLAccountCode)
	assert.Equal(t, 50000000.0, pools[0].TotalAmount)

	poolLines, _ := costPoolLineRepo.ListByPool(ctx, pools[0].ID)
	assert.Len(t, poolLines, 2)
}

func TestCostingJE_CollectLaborCosts(t *testing.T) {
	svc, periodRepo, _, costPoolRepo, _, _, _ := setupCostingJETest(t)
	ctx := context.Background()

	periodID := createTestCostingPeriod(t, periodRepo, "COMP1", 2026, 8)

	lines := []CostPoolLineInput{
		{SourceID: "PAY-001", Description: "Direct labor line 1", Amount: 40000000},
	}

	err := svc.CollectLaborCosts(ctx, "COMP1", periodID, lines)
	assert.NoError(t, err)

	pools, _ := costPoolRepo.ListByPeriod(ctx, "COMP1", periodID)
	assert.Len(t, pools, 1)
	assert.Equal(t, "622", pools[0].GLAccountCode)
	assert.Equal(t, 40000000.0, pools[0].TotalAmount)
}

func TestCostingJE_ReopenPeriod(t *testing.T) {
	svc, periodRepo, _, costPoolRepo, _, resultRepo, _ := setupCostingJETest(t)
	ctx := context.Background()

	periodID := createTestCostingPeriod(t, periodRepo, "COMP1", 2026, 8)
	period, _ := periodRepo.GetByID(ctx, periodID)
	period.Status = "CLOSED"
	_ = periodRepo.Update(ctx, period)

	result := &domain.CostingResult{
		ID: "CR-001", CompanyID: "COMP1", PeriodID: periodID, CostObjectID: "OBJ-001",
		CostingMethod: "SIMPLE", TotalCost: 50000000, Status: "FINAL",
	}
	_ = resultRepo.Create(ctx, result)

	pool := &domain.CostPool{
		ID: "POOL-001", CompanyID: "COMP1", PeriodID: periodID,
		GLAccountCode: "621", Name: "Direct materials", Status: "OPEN", TotalAmount: 30000000,
	}
	_ = costPoolRepo.Create(ctx, pool)

	err := svc.ReopenPeriod(ctx, "COMP1", periodID)
	assert.NoError(t, err)

	periodAfter, _ := periodRepo.GetByID(ctx, periodID)
	assert.Equal(t, "OPEN", periodAfter.Status)

	entries := svc.jeCreator.(*testJECreator).entries
	assert.Len(t, entries, 2)
}
