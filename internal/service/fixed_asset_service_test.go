package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/domain"
	"gotax/internal/repository"
)

func setupFAService(t *testing.T) (FAServiceInterface, context.Context) {
	t.Helper()
	faRepo := repository.NewMemoryFARepo()
	svc := NewFAService(faRepo)
	return svc, context.Background()
}

// ─── Categories ──────────────────────────────────────────────────────

func TestCreateCategory_Success(t *testing.T) {
	svc, ctx := setupFAService(t)
	cat := &domain.FixedAssetCategory{
		CompanyID:   "CMP001",
		Code:        "MACHINERY",
		Name:        "Máy móc thiết bị",
		Level:       1,
		AssetAccountID:        "2111",
		DepreciationAccountID: "2141",
		ExpenseAccountID:      "6271",
	}
	err := svc.CreateCategory(ctx, cat)
	require.NoError(t, err)
	assert.NotEmpty(t, cat.ID)
}

func TestCreateCategory_Duplicate(t *testing.T) {
	svc, ctx := setupFAService(t)
	cat := &domain.FixedAssetCategory{
		CompanyID:   "CMP001",
		Code:        "MACHINERY",
		Name:        "Máy móc",
		Level:       1,
		AssetAccountID:        "2111",
		DepreciationAccountID: "2141",
		ExpenseAccountID:      "6271",
	}
	require.NoError(t, svc.CreateCategory(ctx, cat))
	err := svc.CreateCategory(ctx, cat)
	assert.ErrorIs(t, err, domain.ErrFACategoryCodeExists)
}

func TestCreateCategory_Validation(t *testing.T) {
	svc, ctx := setupFAService(t)
	err := svc.CreateCategory(ctx, &domain.FixedAssetCategory{})
	assert.ErrorIs(t, err, domain.ErrFACategoryCodeRequired)
}

func TestGetCategory(t *testing.T) {
	svc, ctx := setupFAService(t)
	cat := &domain.FixedAssetCategory{
		CompanyID:   "CMP001", Code: "MACHINERY", Name: "Máy móc", Level: 1,
		AssetAccountID: "2111", DepreciationAccountID: "2141", ExpenseAccountID: "6271",
	}
	require.NoError(t, svc.CreateCategory(ctx, cat))

	got, err := svc.GetCategory(ctx, cat.ID)
	require.NoError(t, err)
	assert.Equal(t, "MACHINERY", got.Code)
}

func TestGetCategoryByCode(t *testing.T) {
	svc, ctx := setupFAService(t)
	cat := &domain.FixedAssetCategory{
		CompanyID: "CMP001", Code: "BUILDING", Name: "Nhà cửa", Level: 1,
		AssetAccountID: "2111", DepreciationAccountID: "2141", ExpenseAccountID: "6271",
	}
	require.NoError(t, svc.CreateCategory(ctx, cat))

	got, err := svc.GetCategoryByCode(ctx, "CMP001", "BUILDING")
	require.NoError(t, err)
	assert.Equal(t, cat.ID, got.ID)
}

func TestListCategories(t *testing.T) {
	svc, ctx := setupFAService(t)
	for _, c := range []string{"MACHINERY", "BUILDING", "VEHICLE"} {
		cat := &domain.FixedAssetCategory{
			CompanyID: "CMP001", Code: c, Name: c, Level: 1,
			AssetAccountID: "2111", DepreciationAccountID: "2141", ExpenseAccountID: "6271",
		}
		require.NoError(t, svc.CreateCategory(ctx, cat))
	}

	list, err := svc.ListCategories(ctx, domain.FACategoryFilter{CompanyID: "CMP001"})
	require.NoError(t, err)
	assert.Len(t, list, 3)
}

func TestUpdateCategory(t *testing.T) {
	svc, ctx := setupFAService(t)
	cat := &domain.FixedAssetCategory{
		CompanyID: "CMP001", Code: "MACHINERY", Name: "Old", Level: 1,
		AssetAccountID: "2111", DepreciationAccountID: "2141", ExpenseAccountID: "6271",
	}
	require.NoError(t, svc.CreateCategory(ctx, cat))
	cat.Name = "New Name"
	require.NoError(t, svc.UpdateCategory(ctx, cat))

	got, _ := svc.GetCategory(ctx, cat.ID)
	assert.Equal(t, "New Name", got.Name)
}

func TestDeleteCategory(t *testing.T) {
	svc, ctx := setupFAService(t)
	cat := &domain.FixedAssetCategory{
		CompanyID: "CMP001", Code: "MACHINERY", Name: "Test", Level: 1,
		AssetAccountID: "2111", DepreciationAccountID: "2141", ExpenseAccountID: "6271",
	}
	require.NoError(t, svc.CreateCategory(ctx, cat))
	require.NoError(t, svc.DeleteCategory(ctx, cat.ID))
	_, err := svc.GetCategory(ctx, cat.ID)
	assert.ErrorIs(t, err, domain.ErrFACategoryNotFound)
}

// ─── Fixed Assets ────────────────────────────────────────────────────

func makeTestAsset(code string) *domain.FixedAsset {
	return &domain.FixedAsset{
		CompanyID:             "CMP001",
		Code:                  code,
		Name:                  "Test Asset " + code,
		CategoryID:            "cat-1",
		OriginalCost:          100000000,
		UsefulLifeMonths:      60,
		DepreciationMethod:    domain.DepStraightLine,
		DepartmentID:          "dept-1",
		Location:              "Hanoi",
		AssetAccountID:        "2111",
		DepreciationAccountID: "2141",
		ExpenseAccountID:      "6271",
		Source:                domain.FASourcePurchase,
		AcquisitionDate:       "2026-01-15",
	}
}

func TestCreateAsset_Success(t *testing.T) {
	svc, ctx := setupFAService(t)
	a := makeTestAsset("FA001")
	err := svc.CreateAsset(ctx, a)
	require.NoError(t, err)
	assert.NotEmpty(t, a.ID)
	assert.Equal(t, domain.FADraft, a.Status)
	assert.Equal(t, 100000000.0, a.CarryingAmount)
}

func TestCreateAsset_DuplicateCode(t *testing.T) {
	svc, ctx := setupFAService(t)
	a := makeTestAsset("FA001")
	require.NoError(t, svc.CreateAsset(ctx, a))
	err := svc.CreateAsset(ctx, a)
	assert.ErrorIs(t, err, domain.ErrFACodeExists)
}

func TestCreateAsset_Validation(t *testing.T) {
	svc, ctx := setupFAService(t)
	err := svc.CreateAsset(ctx, &domain.FixedAsset{})
	assert.ErrorIs(t, err, domain.ErrFACodeRequired)
}

func TestGetAsset(t *testing.T) {
	svc, ctx := setupFAService(t)
	a := makeTestAsset("FA001")
	require.NoError(t, svc.CreateAsset(ctx, a))

	got, err := svc.GetAsset(ctx, a.ID)
	require.NoError(t, err)
	assert.Equal(t, a.Code, got.Code)
	assert.Equal(t, 100000000.0, got.CarryingAmount)
}

func TestGetAssetByCode(t *testing.T) {
	svc, ctx := setupFAService(t)
	a := makeTestAsset("FA001")
	require.NoError(t, svc.CreateAsset(ctx, a))

	got, err := svc.GetAssetByCode(ctx, "CMP001", "FA001")
	require.NoError(t, err)
	assert.Equal(t, a.ID, got.ID)
}

func TestListAssets(t *testing.T) {
	svc, ctx := setupFAService(t)
	for _, code := range []string{"FA001", "FA002", "FA003"} {
		require.NoError(t, svc.CreateAsset(ctx, makeTestAsset(code)))
	}

	list, total, err := svc.ListAssets(ctx, domain.FAListFilter{CompanyID: "CMP001", Limit: 20, Offset: 0})
	require.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Len(t, list, 3)
}

func TestListAssets_FilterByStatus(t *testing.T) {
	svc, ctx := setupFAService(t)
	a := makeTestAsset("FA001")
	require.NoError(t, svc.CreateAsset(ctx, a))

	active := domain.FAActive
	list, total, err := svc.ListAssets(ctx, domain.FAListFilter{CompanyID: "CMP001", Status: &active})
	require.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.Len(t, list, 0)
}

func TestUpdateAsset(t *testing.T) {
	svc, ctx := setupFAService(t)
	a := makeTestAsset("FA001")
	require.NoError(t, svc.CreateAsset(ctx, a))

	a.Name = "Updated Name"
	require.NoError(t, svc.UpdateAsset(ctx, a))

	got, _ := svc.GetAsset(ctx, a.ID)
	assert.Equal(t, "Updated Name", got.Name)
}

func TestDeleteAsset(t *testing.T) {
	svc, ctx := setupFAService(t)
	a := makeTestAsset("FA001")
	require.NoError(t, svc.CreateAsset(ctx, a))
	require.NoError(t, svc.DeleteAsset(ctx, a.ID))
	_, err := svc.GetAsset(ctx, a.ID)
	assert.ErrorIs(t, err, domain.ErrFANotFound)
}

func TestChangeAssetStatus(t *testing.T) {
	svc, ctx := setupFAService(t)
	a := makeTestAsset("FA001")
	require.NoError(t, svc.CreateAsset(ctx, a))
	require.NoError(t, svc.ChangeAssetStatus(ctx, a.ID, domain.FAActive, "user1"))
	got, _ := svc.GetAsset(ctx, a.ID)
	assert.Equal(t, domain.FAActive, got.Status)
}

// ─── Depreciation Engine ────────────────────────────────────────────

func TestDepreciation_StraightLine(t *testing.T) {
	svc, ctx := setupFAService(t)
	a := makeTestAsset("FA001")
	a.DepreciationStartDate = strPtr("2026-01-01")
	a.DepreciationEndDate = strPtr("2030-12-31")
	require.NoError(t, svc.CreateAsset(ctx, a))
	require.NoError(t, svc.ChangeAssetStatus(ctx, a.ID, domain.FAActive, "user1"))

	results, err := svc.RunDepreciation(ctx, domain.FARunDepreciationInput{
		CompanyID: "CMP001",
		Year:      2026,
		Month:     1,
		PeriodID:  "PER-2026-01",
		CreatedBy: "user1",
	})
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.InDelta(t, 1666666.67, results[0].Amount, 0.01)
	assert.Equal(t, 1666666.67, results[0].AccAfter)
	assert.Equal(t, 98333333.33, results[0].CarryAfter)
}

func TestDepreciation_FullLife(t *testing.T) {
	svc, ctx := setupFAService(t)
	a := makeTestAsset("FA100")
	a.OriginalCost = 12000000
	a.UsefulLifeMonths = 12
	a.ResidualValue = 0
	require.NoError(t, svc.CreateAsset(ctx, a))

	for i := 1; i <= 12; i++ {
		pid := "PER-2026-" + fmt.Sprintf("%02d", i)
		_, err := svc.RunDepreciationForAsset(ctx, domain.FARunDepreciationInput{
			CompanyID: "CMP001", Year: 2026, Month: i, PeriodID: pid, CreatedBy: "user1",
		}, a.ID)
		require.NoError(t, err)
	}

	got, _ := svc.GetAsset(ctx, a.ID)
	assert.Equal(t, domain.FAFullyDepr, got.Status)
	assert.Equal(t, 0.0, got.CarryingAmount)
}

func TestDepreciation_DecliningBalance(t *testing.T) {
	svc, ctx := setupFAService(t)
	a := makeTestAsset("FA002")
	a.DepreciationMethod = domain.DepDecliningBalance
	a.OriginalCost = 100000000
	a.UsefulLifeMonths = 60
	require.NoError(t, svc.CreateAsset(ctx, a))

	result, err := svc.RunDepreciationForAsset(ctx, domain.FARunDepreciationInput{
		CompanyID: "CMP001", Year: 2026, Month: 1, PeriodID: "PER-2026-01", CreatedBy: "user1",
	}, a.ID)
	require.NoError(t, err)
	assert.True(t, result.Amount > 0)
}

func TestDepreciation_DuplicatePeriod(t *testing.T) {
	svc, ctx := setupFAService(t)
	a := makeTestAsset("FA003")
	require.NoError(t, svc.CreateAsset(ctx, a))

	input := domain.FARunDepreciationInput{
		CompanyID: "CMP001", Year: 2026, Month: 1, PeriodID: "PER-2026-01", CreatedBy: "user1",
	}
	_, err := svc.RunDepreciationForAsset(ctx, input, a.ID)
	require.NoError(t, err)

	_, err = svc.RunDepreciationForAsset(ctx, input, a.ID)
	assert.ErrorIs(t, err, domain.ErrFADepreciationExists)
}

// ─── Business Operations ─────────────────────────────────────────────

func TestDisposeAsset(t *testing.T) {
	svc, ctx := setupFAService(t)
	a := makeTestAsset("FA001")
	require.NoError(t, svc.CreateAsset(ctx, a))
	require.NoError(t, svc.ChangeAssetStatus(ctx, a.ID, domain.FADepreciating, "user1"))

	err := svc.DisposeAsset(ctx, domain.FADisposalInput{
		FixedAssetID: a.ID,
		DisposalType: domain.DisposalLiquidation,
		DisposalDate: "2026-06-01",
		Proceeds:     50000000,
		CreatedBy:    "user1",
	})
	require.NoError(t, err)

	got, _ := svc.GetAsset(ctx, a.ID)
	assert.Equal(t, domain.FADisposed, got.Status)
}

func TestSellAsset(t *testing.T) {
	svc, ctx := setupFAService(t)
	a := makeTestAsset("FA001")
	require.NoError(t, svc.CreateAsset(ctx, a))
	require.NoError(t, svc.ChangeAssetStatus(ctx, a.ID, domain.FADepreciating, "user1"))

	err := svc.SellAsset(ctx, domain.FASaleInput{
		FixedAssetID: a.ID,
		SaleDate:     "2026-06-01",
		Proceeds:     60000000,
		CreatedBy:    "user1",
	})
	require.NoError(t, err)

	got, _ := svc.GetAsset(ctx, a.ID)
	assert.Equal(t, domain.FASold, got.Status)
}

func TestAdjustAsset(t *testing.T) {
	svc, ctx := setupFAService(t)
	a := makeTestAsset("FA001")
	require.NoError(t, svc.CreateAsset(ctx, a))

	newCost := 120000000.0
	err := svc.AdjustAsset(ctx, domain.FAAdjustmentInput{
		FixedAssetID:    a.ID,
		NewOriginalCost: &newCost,
		AdjustmentDate:  "2026-03-01",
		Reason:          "Revaluation",
		CreatedBy:       "user1",
	})
	require.NoError(t, err)

	got, _ := svc.GetAsset(ctx, a.ID)
	assert.Equal(t, 120000000.0, got.OriginalCost)
}

func TestRevalueAsset(t *testing.T) {
	svc, ctx := setupFAService(t)
	a := makeTestAsset("FA001")
	require.NoError(t, svc.CreateAsset(ctx, a))

	err := svc.RevalueAsset(ctx, domain.FARevaluationInput{
		FixedAssetID:   a.ID,
		FairValue:      150000000,
		RevaluationDate: "2026-06-01",
		CreatedBy:      "user1",
	})
	require.NoError(t, err)

	got, _ := svc.GetAsset(ctx, a.ID)
	assert.Equal(t, 150000000.0, got.CarryingAmount)
}

func TestImpairAsset(t *testing.T) {
	svc, ctx := setupFAService(t)
	a := makeTestAsset("FA001")
	require.NoError(t, svc.CreateAsset(ctx, a))

	err := svc.ImpairAsset(ctx, domain.FAImpairmentInput{
		FixedAssetID:      a.ID,
		ImpairmentAmount:  10000000,
		ImpairmentDate:    "2026-06-01",
		Reason:            "Damage",
		CreatedBy:         "user1",
	})
	require.NoError(t, err)

	got, _ := svc.GetAsset(ctx, a.ID)
	assert.Equal(t, 10000000.0, got.AccumulatedDepreciation)
	assert.Equal(t, 90000000.0, got.CarryingAmount)
}

func TestTransferAsset(t *testing.T) {
	svc, ctx := setupFAService(t)
	a := makeTestAsset("FA001")
	a.DepartmentID = "dept-1"
	require.NoError(t, svc.CreateAsset(ctx, a))

	err := svc.TransferAsset(ctx, domain.FATransferInput{
		FixedAssetID:  a.ID,
		DepartmentID:  "dept-2",
		EffectiveDate: "2026-04-01",
		CreatedBy:     "user1",
	})
	require.NoError(t, err)

	got, _ := svc.GetAsset(ctx, a.ID)
	assert.Equal(t, "dept-2", got.DepartmentID)
}

func TestCIPTransfer(t *testing.T) {
	svc, ctx := setupFAService(t)
	a := makeTestAsset("FA001")
	a.Source = domain.FASourceConstruction
	a.OriginalCost = 0
	require.NoError(t, svc.CreateAsset(ctx, a))

	err := svc.CIPTransfer(ctx, domain.FACIPTransferInput{
		FixedAssetID: a.ID,
		CIPAccountID: "2411",
		TotalCost:    200000000,
		TransferDate: "2026-06-01",
		CreatedBy:    "user1",
	})
	require.NoError(t, err)

	got, _ := svc.GetAsset(ctx, a.ID)
	assert.Equal(t, 200000000.0, got.OriginalCost)
	assert.Equal(t, domain.FAActive, got.Status)
}

func TestSuspendResumeDepreciation(t *testing.T) {
	svc, ctx := setupFAService(t)
	a := makeTestAsset("FA001")
	require.NoError(t, svc.CreateAsset(ctx, a))
	require.NoError(t, svc.ChangeAssetStatus(ctx, a.ID, domain.FADepreciating, "user1"))

	require.NoError(t, svc.SuspendDepreciation(ctx, domain.FASuspendInput{
		FixedAssetID: a.ID,
		SuspendDate:  "2026-05-01",
		Reason:       "Under maintenance",
		CreatedBy:    "user1",
	}))
	got, _ := svc.GetAsset(ctx, a.ID)
	assert.Equal(t, domain.FASuspended, got.Status)

	require.NoError(t, svc.ResumeDepreciation(ctx, domain.FAResumeInput{
		FixedAssetID: a.ID,
		ResumeDate:   "2026-07-01",
		CreatedBy:    "user1",
	}))
	got, _ = svc.GetAsset(ctx, a.ID)
	assert.Equal(t, domain.FADepreciating, got.Status)
}

func TestTransferAsset_SameDepartment(t *testing.T) {
	svc, ctx := setupFAService(t)
	a := makeTestAsset("FA001")
	a.DepartmentID = "dept-1"
	require.NoError(t, svc.CreateAsset(ctx, a))

	err := svc.TransferAsset(ctx, domain.FATransferInput{
		FixedAssetID: a.ID, DepartmentID: "dept-1",
		EffectiveDate: "2026-06-01", CreatedBy: "user1",
	})
	assert.ErrorIs(t, err, domain.ErrFATransferSameDepartment)
}

func TestCIPTransfer_NotApplicable(t *testing.T) {
	svc, ctx := setupFAService(t)
	a := makeTestAsset("FA001")
	require.NoError(t, svc.CreateAsset(ctx, a))

	err := svc.CIPTransfer(ctx, domain.FACIPTransferInput{
		FixedAssetID: a.ID, CIPAccountID: "2411",
		TotalCost: 200000000, TransferDate: "2026-06-01", CreatedBy: "user1",
	})
	assert.ErrorIs(t, err, domain.ErrFACIPTransferNotApplicable)
}

func TestSuspend_NotApplicable(t *testing.T) {
	svc, ctx := setupFAService(t)
	a := makeTestAsset("FA001")
	require.NoError(t, svc.CreateAsset(ctx, a))

	err := svc.SuspendDepreciation(ctx, domain.FASuspendInput{
		FixedAssetID: a.ID, SuspendDate: "2026-06-01", CreatedBy: "user1",
	})
	assert.ErrorIs(t, err, domain.ErrFASuspensionNotApplicable)
}

func TestResume_NotApplicable(t *testing.T) {
	svc, ctx := setupFAService(t)
	a := makeTestAsset("FA001")
	require.NoError(t, svc.CreateAsset(ctx, a))

	err := svc.ResumeDepreciation(ctx, domain.FAResumeInput{
		FixedAssetID: a.ID, ResumeDate: "2026-07-01", CreatedBy: "user1",
	})
	assert.ErrorIs(t, err, domain.ErrFAResumeNotApplicable)
}

// ─── Allocations ─────────────────────────────────────────────────────

func TestSetAllocations(t *testing.T) {
	svc, ctx := setupFAService(t)
	a := makeTestAsset("FA001")
	require.NoError(t, svc.CreateAsset(ctx, a))

	allocs := []domain.FixedAssetAllocation{
		{DepartmentID: "dept-1", AllocationPct: 60, ExpenseAccountID: "6271"},
		{DepartmentID: "dept-2", AllocationPct: 40, ExpenseAccountID: "6272"},
	}
	require.NoError(t, svc.SetAllocations(ctx, a.ID, allocs))

	got, err := svc.GetAllocations(ctx, a.ID)
	require.NoError(t, err)
	assert.Len(t, got, 2)
}

func TestSetAllocations_InvalidSum(t *testing.T) {
	svc, ctx := setupFAService(t)
	a := makeTestAsset("FA001")
	require.NoError(t, svc.CreateAsset(ctx, a))

	allocs := []domain.FixedAssetAllocation{
		{DepartmentID: "dept-1", AllocationPct: 50, ExpenseAccountID: "6271"},
	}
	err := svc.SetAllocations(ctx, a.ID, allocs)
	assert.ErrorIs(t, err, domain.ErrFAAllocationPctSum)
}

// ─── Depreciation Queries ───────────────────────────────────────────

func TestGetDepreciationByPeriod(t *testing.T) {
	svc, ctx := setupFAService(t)
	a := makeTestAsset("FA001")
	require.NoError(t, svc.CreateAsset(ctx, a))
	require.NoError(t, svc.ChangeAssetStatus(ctx, a.ID, domain.FAActive, "user1"))

	_, err := svc.RunDepreciationForAsset(ctx, domain.FARunDepreciationInput{
		CompanyID: "CMP001", Year: 2026, Month: 1, PeriodID: "PER-2026-01", CreatedBy: "user1",
	}, a.ID)
	require.NoError(t, err)

	entries, err := svc.GetDepreciationByPeriod(ctx, "PER-2026-01")
	require.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.Equal(t, a.ID, entries[0].FixedAssetID)
	assert.NotEmpty(t, entries[0].CreatedBy)
}

// ─── Inventory ───────────────────────────────────────────────────────

func TestInventoryPlans(t *testing.T) {
	svc, ctx := setupFAService(t)
	p := &domain.FixedAssetInventoryPlan{
		CompanyID: "CMP001",
		PlanDate:  "2026-12-31",
		Status:    "DRAFT",
		CreatedBy: "user1",
	}
	require.NoError(t, svc.CreateInventoryPlan(ctx, p))
	assert.NotEmpty(t, p.ID)

	got, err := svc.GetInventoryPlan(ctx, p.ID)
	require.NoError(t, err)
	assert.Equal(t, p.ID, got.ID)

	list, err := svc.ListInventoryPlans(ctx, "CMP001")
	require.NoError(t, err)
	assert.Len(t, list, 1)
}

func TestInventoryResults(t *testing.T) {
	svc, ctx := setupFAService(t)
	a := makeTestAsset("FA001")
	require.NoError(t, svc.CreateAsset(ctx, a))

	p := &domain.FixedAssetInventoryPlan{
		CompanyID: "CMP001", PlanDate: "2026-12-31", CreatedBy: "user1",
	}
	require.NoError(t, svc.CreateInventoryPlan(ctx, p))

	r := &domain.FixedAssetInventoryResult{
		PlanID:       p.ID,
		FixedAssetID: a.ID,
		Discrepancy:  "MATCH",
	}
	require.NoError(t, svc.CreateInventoryResult(ctx, r))
	assert.NotEmpty(t, r.ID)

	results, err := svc.GetInventoryResults(ctx, p.ID)
	require.NoError(t, err)
	assert.Len(t, results, 1)
}

// ─── Transactions ───────────────────────────────────────────────────

func TestTransactions(t *testing.T) {
	svc, ctx := setupFAService(t)
	a := makeTestAsset("FA001")
	require.NoError(t, svc.CreateAsset(ctx, a))

	txn := &domain.FixedAssetTransaction{
		CompanyID:       "CMP001",
		FixedAssetID:    a.ID,
		TransactionType: domain.FATrxAcquisition,
		TransactionDate: "2026-01-15",
		Amount:          100000000,
		CreatedBy:       "user1",
	}
	require.NoError(t, svc.RecordTransaction(ctx, txn))
	assert.NotEmpty(t, txn.ID)

	list, err := svc.GetTransactions(ctx, a.ID)
	require.NoError(t, err)
	assert.Len(t, list, 1)
}

func strPtr(s string) *string { return &s }
