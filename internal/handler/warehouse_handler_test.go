package handler

import (
	"context"
	"encoding/json"
	"fmt"
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

func setupWarehouseTest(t *testing.T) (*gin.Engine, *service.WarehouseService, context.Context) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	whRepo := repository.NewMemoryWarehouseRepo()
	catRepo := repository.NewMemoryItemCategoryRepo()
	itemRepo := repository.NewMemoryItemRepo()
	balRepo := repository.NewMemoryStockBalanceRepo()
	txnRepo := repository.NewMemoryInventoryTransactionRepo()
	trfRepo := repository.NewMemoryStockTransferRepo()
	adjRepo := repository.NewMemoryStockAdjustmentRepo()
	takeRepo := repository.NewMemoryStockTakeRepo()
	valRepo := repository.NewMemoryInventoryValuationRunRepo()
	whSvc := service.NewWarehouseService(whRepo, catRepo, itemRepo, balRepo, txnRepo, trfRepo, adjRepo, takeRepo, valRepo)
	whH := NewWarehouseHandler(whSvc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("username", "testuser")
		c.Set("role", "admin")
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterWarehouseRoutes(r, whH, noopMW)

	return r, whSvc, context.Background()
}

// ─── Warehouse ────────────────────────────────────────────────────────

func TestCreateWarehouse(t *testing.T) {
	r, _, _ := setupWarehouseTest(t)
	body := `{"code":"WH001","name":"Main Warehouse"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/warehouse/warehouses?company_id=CMP001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
	var wh domain.Warehouse
	json.Unmarshal(w.Body.Bytes(), &wh)
	assert.NotEmpty(t, wh.ID)
	assert.Equal(t, "WH001", wh.Code)
}

func TestCreateWarehouseDuplicateCode(t *testing.T) {
	r, svc, ctx := setupWarehouseTest(t)
	require.NoError(t, svc.CreateWarehouse(ctx, &domain.Warehouse{CompanyID: "CMP001", Code: "WH001", Name: "Existing"}))
	body := `{"code":"WH001","name":"Duplicate"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/warehouse/warehouses?company_id=CMP001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestGetWarehouse(t *testing.T) {
	r, svc, ctx := setupWarehouseTest(t)
	wh := &domain.Warehouse{CompanyID: "CMP001", Code: "WH001", Name: "Main"}
	require.NoError(t, svc.CreateWarehouse(ctx, wh))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/warehouse/warehouses/"+wh.ID, nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	var got domain.Warehouse
	json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, wh.ID, got.ID)
}

func TestGetWarehouseNotFound(t *testing.T) {
	r, _, _ := setupWarehouseTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/warehouse/warehouses/nonexistent", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

func TestListWarehouses(t *testing.T) {
	r, svc, ctx := setupWarehouseTest(t)
	svc.CreateWarehouse(ctx, &domain.Warehouse{CompanyID: "CMP001", Code: "WH01", Name: "WH1"})
	svc.CreateWarehouse(ctx, &domain.Warehouse{CompanyID: "CMP001", Code: "WH02", Name: "WH2"})
	svc.CreateWarehouse(ctx, &domain.Warehouse{CompanyID: "CMP002", Code: "WH03", Name: "WH3"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/warehouse/warehouses?company_id=CMP001", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	var list []domain.Warehouse
	json.Unmarshal(w.Body.Bytes(), &list)
	assert.Equal(t, 2, len(list))
}

func TestUpdateWarehouse(t *testing.T) {
	r, svc, ctx := setupWarehouseTest(t)
	wh := &domain.Warehouse{CompanyID: "CMP001", Code: "WH001", Name: "Old"}
	require.NoError(t, svc.CreateWarehouse(ctx, wh))

	body := `{"code":"WH001","name":"New Name"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/warehouse/warehouses/"+wh.ID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestDeleteWarehouse(t *testing.T) {
	r, svc, ctx := setupWarehouseTest(t)
	wh := &domain.Warehouse{CompanyID: "CMP001", Code: "WH001", Name: "Test"}
	require.NoError(t, svc.CreateWarehouse(ctx, wh))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/warehouse/warehouses/"+wh.ID, nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 204, w.Code)
}

// ─── Category ─────────────────────────────────────────────────────────

func TestCreateCategory(t *testing.T) {
	r, _, _ := setupWarehouseTest(t)
	body := `{"code":"RAW","name":"Raw Materials"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/warehouse/categories?company_id=CMP001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
	var cat domain.ItemCategory
	json.Unmarshal(w.Body.Bytes(), &cat)
	assert.NotEmpty(t, cat.ID)
}

func TestCreateCategoryDuplicateCode(t *testing.T) {
	r, svc, ctx := setupWarehouseTest(t)
	require.NoError(t, svc.CreateCategory(ctx, &domain.ItemCategory{CompanyID: "CMP001", Code: "RAW", Name: "Raw"}))
	body := `{"code":"RAW","name":"Duplicate"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/warehouse/categories?company_id=CMP001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestGetCategory(t *testing.T) {
	r, svc, ctx := setupWarehouseTest(t)
	cat := &domain.ItemCategory{CompanyID: "CMP001", Code: "RAW", Name: "Raw Materials"}
	require.NoError(t, svc.CreateCategory(ctx, cat))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/warehouse/categories/"+cat.ID, nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

// ─── Item ─────────────────────────────────────────────────────────────

func TestCreateItem(t *testing.T) {
	r, _, _ := setupWarehouseTest(t)
	body := `{"code":"ITEM001","name":"Test Item","unit":"pcs","valuation_method":"weighted_average"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/warehouse/items?company_id=CMP001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
	var item domain.Item
	json.Unmarshal(w.Body.Bytes(), &item)
	assert.NotEmpty(t, item.ID)
}

func TestCreateItemValidationError(t *testing.T) {
	r, _, _ := setupWarehouseTest(t)
	body := `{"name":"No Code"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/warehouse/items?company_id=CMP001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestGetItem(t *testing.T) {
	r, svc, ctx := setupWarehouseTest(t)
	item := &domain.Item{CompanyID: "CMP001", Code: "ITEM001", Name: "Test", Unit: "pcs"}
	require.NoError(t, svc.CreateItem(ctx, item))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/warehouse/items/"+item.ID, nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestDeleteItem(t *testing.T) {
	r, svc, ctx := setupWarehouseTest(t)
	item := &domain.Item{CompanyID: "CMP001", Code: "ITEM001", Name: "Test", Unit: "pcs"}
	require.NoError(t, svc.CreateItem(ctx, item))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/warehouse/items/"+item.ID, nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 204, w.Code)
}

// ─── Stock Transfer ───────────────────────────────────────────────────

func TestCreateStockTransfer(t *testing.T) {
	r, _, _ := setupWarehouseTest(t)
	body := `{"transfer_number":"TRF001","from_warehouse_id":"WH-A","to_warehouse_id":"WH-B","items":[{"item_id":"ITEM001","quantity":10,"unit_cost":5000}]}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/warehouse/transfers?company_id=CMP001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
	var trf domain.StockTransfer
	json.Unmarshal(w.Body.Bytes(), &trf)
	assert.NotEmpty(t, trf.ID)
	assert.Equal(t, domain.TransStatusDraft, trf.Status)
}

func TestCreateStockTransferSameWarehouse(t *testing.T) {
	r, _, _ := setupWarehouseTest(t)
	body := `{"transfer_number":"TRF001","from_warehouse_id":"WH-A","to_warehouse_id":"WH-A","items":[{"item_id":"ITEM001","quantity":10,"unit_cost":5000}]}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/warehouse/transfers?company_id=CMP001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestStockTransferLifecycle(t *testing.T) {
	r, svc, ctx := setupWarehouseTest(t)
	trf := &domain.StockTransfer{
		CompanyID: "CMP001", TransferNumber: "TRF001",
		FromWarehouseID: "WH-A", ToWarehouseID: "WH-B",
		Items: []domain.TransferItem{{ItemID: "ITM001", Quantity: 10, UnitCost: 5000}},
	}
	require.NoError(t, svc.CreateStockTransfer(ctx, trf))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("PATCH", "/api/v1/warehouse/transfers/"+trf.ID+"/submit", nil))
	assert.Equal(t, 200, w.Code, "submit should succeed")

	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("PATCH", "/api/v1/warehouse/transfers/"+trf.ID+"/approve", nil))
	assert.Equal(t, 200, w.Code, "approve should succeed")

	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("PATCH", "/api/v1/warehouse/transfers/"+trf.ID+"/transfer", nil))
	assert.Equal(t, 200, w.Code, "transfer should succeed")

	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("PATCH", "/api/v1/warehouse/transfers/"+trf.ID+"/complete", nil))
	assert.Equal(t, 200, w.Code, "complete should succeed")

	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("PATCH", "/api/v1/warehouse/transfers/"+trf.ID+"/submit", nil))
	assert.Equal(t, 400, w.Code, "submit from completed should fail")
}

func TestCancelStockTransfer(t *testing.T) {
	r, svc, ctx := setupWarehouseTest(t)
	trf := &domain.StockTransfer{
		CompanyID: "CMP001", TransferNumber: "TRF001",
		FromWarehouseID: "WH-A", ToWarehouseID: "WH-B",
		Items: []domain.TransferItem{{ItemID: "ITM001", Quantity: 10, UnitCost: 5000}},
	}
	require.NoError(t, svc.CreateStockTransfer(ctx, trf))

	body := `{"reason":"no longer needed"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/api/v1/warehouse/transfers/"+trf.ID+"/cancel", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

// ─── Stock Adjustment ─────────────────────────────────────────────────

func TestCreateStockAdjustment(t *testing.T) {
	r, _, _ := setupWarehouseTest(t)
	body := `{"adjustment_number":"ADJ001","warehouse_id":"WH-A","adj_type":"INCREASE","items":[{"item_id":"ITEM001","qty_before":10,"qty_after":15,"unit_cost":5000,"reason":"found"}]}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/warehouse/adjustments?company_id=CMP001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
}

func TestStockAdjustmentLifecycle(t *testing.T) {
	r, svc, ctx := setupWarehouseTest(t)
	a := &domain.StockAdjustment{
		CompanyID: "CMP001", AdjustmentNumber: "ADJ001",
		WarehouseID: "WH-A", AdjType: domain.AdjIncrease,
		Items: []domain.AdjItem{{ItemID: "ITM001", QtyBefore: 10, QtyAfter: 15, UnitCost: 5000}},
	}
	require.NoError(t, svc.CreateStockAdjustment(ctx, a))

	id := a.ID

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("PATCH", "/api/v1/warehouse/adjustments/"+id+"/submit", nil))
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("PATCH", "/api/v1/warehouse/adjustments/"+id+"/approve", nil))
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("PATCH", "/api/v1/warehouse/adjustments/"+id+"/post", nil))
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("PATCH", "/api/v1/warehouse/adjustments/"+id+"/submit", nil))
	assert.Equal(t, 400, w.Code)
}

func TestRejectStockAdjustment(t *testing.T) {
	r, svc, ctx := setupWarehouseTest(t)
	a := &domain.StockAdjustment{
		CompanyID: "CMP001", AdjustmentNumber: "ADJ001",
		WarehouseID: "WH-A", AdjType: domain.AdjDecrease,
		Items: []domain.AdjItem{{ItemID: "ITM001", QtyBefore: 10, QtyAfter: 8, UnitCost: 5000}},
	}
	require.NoError(t, svc.CreateStockAdjustment(ctx, a))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("PATCH", "/api/v1/warehouse/adjustments/"+a.ID+"/reject?reason=bad", nil))
	assert.Equal(t, 200, w.Code, "draft can be rejected")
}

// ─── Stock Take ───────────────────────────────────────────────────────

func TestCreateStockTake(t *testing.T) {
	r, _, _ := setupWarehouseTest(t)
	body := `{"take_number":"TAKE001","warehouse_id":"WH-A","take_date":"2026-07-29","items":[{"item_id":"ITEM001","expected_qty":10,"actual_qty":10,"unit_cost":5000}]}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/warehouse/takes?company_id=CMP001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
}

func TestStockTakeLifecycle(t *testing.T) {
	r, svc, ctx := setupWarehouseTest(t)
	tk := &domain.StockTake{
		CompanyID: "CMP001", TakeNumber: "TAKE001",
		WarehouseID: "WH-A", TakeDate: "2026-07-29",
		Items: []domain.TakeItem{{ItemID: "ITM001", ExpectedQty: 10, ActualQty: 12, UnitCost: 5000}},
	}
	require.NoError(t, svc.CreateStockTake(ctx, tk))

	id := tk.ID

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("PATCH", "/api/v1/warehouse/takes/"+id+"/start", nil))
	assert.Equal(t, 200, w.Code, "start should succeed")

	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("PATCH", "/api/v1/warehouse/takes/"+id+"/complete", nil))
	assert.Equal(t, 200, w.Code, "complete should succeed")

	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("PATCH", "/api/v1/warehouse/takes/"+id+"/verify", nil))
	assert.Equal(t, 200, w.Code, "verify should succeed")

	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("PATCH", "/api/v1/warehouse/takes/"+id+"/post", nil))
	assert.Equal(t, 200, w.Code, "post should succeed")

	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("PATCH", "/api/v1/warehouse/takes/"+id+"/start", nil))
	assert.Equal(t, 400, w.Code, "start from posted should fail")
}

// ─── Valuation Run ────────────────────────────────────────────────────

func TestCreateValuationRun(t *testing.T) {
	r, _, _ := setupWarehouseTest(t)
	body := `{"valuation_date":"2026-07-29","method":"weighted_average"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/warehouse/valuations?company_id=CMP001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
	var v domain.InventoryValuationRun
	json.Unmarshal(w.Body.Bytes(), &v)
	assert.NotEmpty(t, v.ID)
	assert.Equal(t, domain.ValRunPending, v.Status)
}

func TestRunValuation(t *testing.T) {
	r, svc, ctx := setupWarehouseTest(t)
	v := &domain.InventoryValuationRun{CompanyID: "CMP001", ValuationDate: "2026-07-29", Method: domain.ValuationWeightedAvg}
	require.NoError(t, svc.CreateValuationRun(ctx, v))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/warehouse/valuations/"+v.ID+"/run", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	updated, _ := svc.GetValuationRun(ctx, v.ID)
	assert.Equal(t, domain.ValRunCompleted, updated.Status)
}

// ─── Stock Balance ────────────────────────────────────────────────────

func TestListStockBalances(t *testing.T) {
	r, svc, ctx := setupWarehouseTest(t)
	svc.UpsertStockBalance(ctx, &domain.StockBalance{
		CompanyID: "CMP001", WarehouseID: "WH-A", ItemID: "ITM001",
		Period: "202607", Quantity: 100, UnitCost: 5000, TotalCost: 500000,
	})
	svc.UpsertStockBalance(ctx, &domain.StockBalance{
		CompanyID: "CMP001", WarehouseID: "WH-A", ItemID: "ITM002",
		Period: "202607", Quantity: 50, UnitCost: 10000, TotalCost: 500000,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/warehouse/balances?company_id=CMP001", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	var list []domain.StockBalance
	json.Unmarshal(w.Body.Bytes(), &list)
	assert.Equal(t, 2, len(list))
}

// ─── Transactions ────────────────────────────────────────────────────

func TestListInventoryTransactions(t *testing.T) {
	r, svc, ctx := setupWarehouseTest(t)
	for i := 0; i < 3; i++ {
		svc.CreateInventoryTransaction(ctx, &domain.InventoryTransaction{
			CompanyID: "CMP001", WarehouseID: "WH-A", ItemID: fmt.Sprintf("ITM%03d", i+1),
			TransType: domain.TransReceipt, Quantity: 10, QtyAfter: 10, UnitCost: 5000, TotalCost: 50000,
			CreatedBy: "test-user",
		})
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/warehouse/transactions?company_id=CMP001", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	var resp struct {
		Data  []domain.InventoryTransaction `json:"data"`
		Total int                           `json:"total"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 3, resp.Total)
}
