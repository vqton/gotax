package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/domain"
	"gotax/internal/repository"
	"gotax/internal/service"
)

type faTestSetup struct {
	r      *gin.Engine
	svc    service.FAServiceInterface
	faRepo domain.FARepository
	compID string
}

func setupFATest(t *testing.T) *faTestSetup {
	t.Helper()
	gin.SetMode(gin.TestMode)

	faRepo := repository.NewMemoryFARepo()
	faSvc := service.NewFAService(faRepo)
	fh := NewFAHandler(faSvc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("username", "testuser")
		c.Set("role", "admin")
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterFixedAssetRoutes(r, fh, noopMW)

	return &faTestSetup{
		r:      r,
		svc:    faSvc,
		faRepo: faRepo,
		compID: "CMP001",
	}
}

// ─── Categories ──────────────────────────────────────────────────────

func TestFACreateCategory(t *testing.T) {
	ts := setupFATest(t)
	body := `{"code":"MACHINERY","name":"Máy móc","level":1,"asset_account_id":"2111","depreciation_account_id":"2141","expense_account_id":"6271"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/fixed-assets/categories?company_id="+ts.compID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp domain.FixedAssetCategory
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "MACHINERY", resp.Code)
	assert.NotEmpty(t, resp.ID)
}

func TestFACreateCategory_Invalid(t *testing.T) {
	ts := setupFATest(t)
	body := `{}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/fixed-assets/categories?company_id="+ts.compID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFAGetCategory(t *testing.T) {
	ts := setupFATest(t)
	cat := &domain.FixedAssetCategory{
		CompanyID: ts.compID, Code: "MACHINERY", Name: "Máy móc", Level: 1,
		AssetAccountID: "2111", DepreciationAccountID: "2141", ExpenseAccountID: "6271",
	}
	require.NoError(t, ts.svc.CreateCategory(nil, cat))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/fixed-assets/categories/"+cat.ID+"?company_id="+ts.compID, nil)
	ts.r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp domain.FixedAssetCategory
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, cat.ID, resp.ID)
}

func TestFAListCategories(t *testing.T) {
	ts := setupFATest(t)
	for _, code := range []string{"A", "B", "C"} {
		cat := &domain.FixedAssetCategory{
			CompanyID: ts.compID, Code: code, Name: code, Level: 1,
			AssetAccountID: "2111", DepreciationAccountID: "2141", ExpenseAccountID: "6271",
		}
		require.NoError(t, ts.svc.CreateCategory(nil, cat))
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/fixed-assets/categories?company_id="+ts.compID, nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFAUpdateCategory(t *testing.T) {
	ts := setupFATest(t)
	cat := &domain.FixedAssetCategory{
		CompanyID: ts.compID, Code: "MACHINERY", Name: "Old", Level: 1,
		AssetAccountID: "2111", DepreciationAccountID: "2141", ExpenseAccountID: "6271",
	}
	require.NoError(t, ts.svc.CreateCategory(nil, cat))

	body := `{"code":"MACHINERY","name":"New Name","level":1,"asset_account_id":"2111","depreciation_account_id":"2141","expense_account_id":"6271"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/fixed-assets/categories/"+cat.ID+"?company_id="+ts.compID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFADeleteCategory(t *testing.T) {
	ts := setupFATest(t)
	cat := &domain.FixedAssetCategory{
		CompanyID: ts.compID, Code: "MACHINERY", Name: "Test", Level: 1,
		AssetAccountID: "2111", DepreciationAccountID: "2141", ExpenseAccountID: "6271",
	}
	require.NoError(t, ts.svc.CreateCategory(nil, cat))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/fixed-assets/categories/"+cat.ID+"?company_id="+ts.compID, nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

// ─── Fixed Assets ────────────────────────────────────────────────────

func TestFACreateAsset(t *testing.T) {
	ts := setupFATest(t)
	body := `{"code":"FA001","name":"Test Asset","category_id":"cat1","original_cost":100000000,"useful_life_months":60,"depreciation_method":"STRAIGHT_LINE","department_id":"dept1","location":"Hanoi","asset_account_id":"2111","depreciation_account_id":"2141","expense_account_id":"6271","source":"PURCHASE","acquisition_date":"2026-01-15"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/fixed-assets?company_id="+ts.compID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp domain.FixedAsset
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "FA001", resp.Code)
	assert.NotEmpty(t, resp.ID)
}

func TestFACreateAsset_Invalid(t *testing.T) {
	ts := setupFATest(t)
	body := `{}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/fixed-assets?company_id="+ts.compID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFAGetAsset(t *testing.T) {
	ts := setupFATest(t)
	a := &domain.FixedAsset{
		CompanyID: ts.compID, Code: "FA001", Name: "Test", CategoryID: "cat1",
		OriginalCost: 100000000, UsefulLifeMonths: 60, DepreciationMethod: domain.DepStraightLine,
		DepartmentID: "dept1", Location: "Hanoi",
		AssetAccountID: "2111", DepreciationAccountID: "2141", ExpenseAccountID: "6271",
		Source: domain.FASourcePurchase, AcquisitionDate: "2026-01-15",
	}
	require.NoError(t, ts.svc.CreateAsset(nil, a))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/fixed-assets/"+a.ID+"?company_id="+ts.compID, nil)
	ts.r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp domain.FixedAsset
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, a.ID, resp.ID)
}

func TestFAListAssets(t *testing.T) {
	ts := setupFATest(t)
	for _, code := range []string{"FA001", "FA002"} {
		a := &domain.FixedAsset{
			CompanyID: ts.compID, Code: code, Name: code, CategoryID: "cat1",
			OriginalCost: 50000000, UsefulLifeMonths: 60, DepreciationMethod: domain.DepStraightLine,
			DepartmentID: "dept1", Location: "HN",
			AssetAccountID: "2111", DepreciationAccountID: "2141", ExpenseAccountID: "6271",
			Source: domain.FASourcePurchase, AcquisitionDate: "2026-01-15",
		}
		require.NoError(t, ts.svc.CreateAsset(nil, a))
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/fixed-assets?company_id="+ts.compID, nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Data  []domain.FixedAsset `json:"data"`
		Total int                 `json:"total"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 2, resp.Total)
}

func TestFAUpdateAsset(t *testing.T) {
	ts := setupFATest(t)
	a := &domain.FixedAsset{
		CompanyID: ts.compID, Code: "FA001", Name: "Old", CategoryID: "cat1",
		OriginalCost: 100000000, UsefulLifeMonths: 60, DepreciationMethod: domain.DepStraightLine,
		DepartmentID: "dept1", Location: "HN",
		AssetAccountID: "2111", DepreciationAccountID: "2141", ExpenseAccountID: "6271",
		Source: domain.FASourcePurchase, AcquisitionDate: "2026-01-15",
	}
	require.NoError(t, ts.svc.CreateAsset(nil, a))

	body := `{"code":"FA001","name":"Updated","category_id":"cat1","original_cost":100000000,"useful_life_months":60,"depreciation_method":"STRAIGHT_LINE","department_id":"dept1","location":"HN","asset_account_id":"2111","depreciation_account_id":"2141","expense_account_id":"6271","source":"PURCHASE","acquisition_date":"2026-01-15"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/fixed-assets/"+a.ID+"?company_id="+ts.compID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFADeleteAsset(t *testing.T) {
	ts := setupFATest(t)
	a := &domain.FixedAsset{
		CompanyID: ts.compID, Code: "FA001", Name: "Test", CategoryID: "cat1",
		OriginalCost: 100000000, UsefulLifeMonths: 60, DepreciationMethod: domain.DepStraightLine,
		DepartmentID: "dept1", Location: "HN",
		AssetAccountID: "2111", DepreciationAccountID: "2141", ExpenseAccountID: "6271",
		Source: domain.FASourcePurchase, AcquisitionDate: "2026-01-15",
	}
	require.NoError(t, ts.svc.CreateAsset(nil, a))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/fixed-assets/"+a.ID+"?company_id="+ts.compID, nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

// ─── Depreciation ────────────────────────────────────────────────────

func TestFARunDepreciation(t *testing.T) {
	ts := setupFATest(t)
	a := &domain.FixedAsset{
		CompanyID: ts.compID, Code: "FA001", Name: "Test", CategoryID: "cat1",
		OriginalCost: 100000000, UsefulLifeMonths: 60, DepreciationMethod: domain.DepStraightLine,
		DepartmentID: "dept1", Location: "HN",
		AssetAccountID: "2111", DepreciationAccountID: "2141", ExpenseAccountID: "6271",
		Source: domain.FASourcePurchase, AcquisitionDate: "2026-01-15",
	}
	require.NoError(t, ts.svc.CreateAsset(nil, a))

	body := `{"year":2026,"month":1,"period_id":"PER-2026-01"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/fixed-assets/depreciation/run?company_id="+ts.compID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFARunDepreciationForAsset(t *testing.T) {
	ts := setupFATest(t)
	a := &domain.FixedAsset{
		CompanyID: ts.compID, Code: "FA001", Name: "Test", CategoryID: "cat1",
		OriginalCost: 100000000, UsefulLifeMonths: 60, DepreciationMethod: domain.DepStraightLine,
		DepartmentID: "dept1", Location: "HN",
		AssetAccountID: "2111", DepreciationAccountID: "2141", ExpenseAccountID: "6271",
		Source: domain.FASourcePurchase, AcquisitionDate: "2026-01-15",
	}
	require.NoError(t, ts.svc.CreateAsset(nil, a))

	body := `{"year":2026,"month":1,"period_id":"PER-2026-01"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/fixed-assets/"+a.ID+"/depreciate?company_id="+ts.compID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// ─── Business Operations ─────────────────────────────────────────────

func TestFADisposeAsset(t *testing.T) {
	ts := setupFATest(t)
	a := &domain.FixedAsset{
		CompanyID: ts.compID, Code: "FA001", Name: "Test", CategoryID: "cat1",
		OriginalCost: 100000000, UsefulLifeMonths: 60, DepreciationMethod: domain.DepStraightLine,
		DepartmentID: "dept1", Location: "HN",
		AssetAccountID: "2111", DepreciationAccountID: "2141", ExpenseAccountID: "6271",
		Source: domain.FASourcePurchase, AcquisitionDate: "2026-01-15",
	}
	require.NoError(t, ts.svc.CreateAsset(nil, a))
	require.NoError(t, ts.svc.ChangeAssetStatus(nil, a.ID, domain.FADepreciating, "user1"))

	body := `{"disposal_type":"LIQUIDATION","disposal_date":"2026-06-01","proceeds":50000000}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/fixed-assets/"+a.ID+"/dispose?company_id="+ts.compID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFASellAsset(t *testing.T) {
	ts := setupFATest(t)
	a := &domain.FixedAsset{
		CompanyID: ts.compID, Code: "FA002", Name: "Test", CategoryID: "cat1",
		OriginalCost: 100000000, UsefulLifeMonths: 60, DepreciationMethod: domain.DepStraightLine,
		DepartmentID: "dept1", Location: "HN",
		AssetAccountID: "2111", DepreciationAccountID: "2141", ExpenseAccountID: "6271",
		Source: domain.FASourcePurchase, AcquisitionDate: "2026-01-15",
	}
	require.NoError(t, ts.svc.CreateAsset(nil, a))
	require.NoError(t, ts.svc.ChangeAssetStatus(nil, a.ID, domain.FADepreciating, "user1"))

	body := `{"sale_date":"2026-06-01","proceeds":60000000}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/fixed-assets/"+a.ID+"/sell?company_id="+ts.compID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// ─── Inventory ───────────────────────────────────────────────────────

func TestFACreateInventoryPlan(t *testing.T) {
	ts := setupFATest(t)
	body := `{"plan_date":"2026-12-31","status":"DRAFT"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/fixed-assets/inventory/plans?company_id="+ts.compID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestFAListInventoryPlans(t *testing.T) {
	ts := setupFATest(t)
	p := &domain.FixedAssetInventoryPlan{
		CompanyID: ts.compID, PlanDate: "2026-12-31", Status: "DRAFT", CreatedBy: "user1",
	}
	require.NoError(t, ts.svc.CreateInventoryPlan(nil, p))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/fixed-assets/inventory/plans?company_id="+ts.compID, nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
