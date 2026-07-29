package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/domain"
	"gotax/internal/repository"
	"gotax/internal/service"
)

func setupPurchaseTest(t *testing.T) (*gin.Engine, *service.PurchaseService, context.Context) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	purRepo := repository.NewMemoryPurchaseRepo()
	purSvc := service.NewPurchaseService(purRepo, purRepo, purRepo, purRepo, purRepo, purRepo)
	purH := NewPurchaseHandler(purSvc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("username", "testuser")
		c.Set("role", "admin")
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterPurchaseRoutes(r, purH, noopMW)

	return r, purSvc, context.Background()
}

// ─── Supplier ────────────────────────────────────────────────────────────

func TestCreateSupplier(t *testing.T) {
	r, _, _ := setupPurchaseTest(t)

	body := `{"code":"SUP001","name":"Test Supplier","tax_code":"TX001"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/purchase/suppliers?company_id=CMP001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	var sup domain.Supplier
	err := json.Unmarshal(w.Body.Bytes(), &sup)
	require.NoError(t, err)
	assert.NotEmpty(t, sup.ID)
	assert.Equal(t, domain.SupplierActive, sup.Status)
}

func TestCreateSupplierDuplicateCode(t *testing.T) {
	r, svc, ctx := setupPurchaseTest(t)
	require.NoError(t, svc.CreateSupplier(ctx, &domain.Supplier{
		CompanyID: "CMP001", Code: "SUP001", Name: "Existing", TaxCode: "TX001",
	}))

	body := `{"code":"SUP001","name":"Duplicate","tax_code":"TX002"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/purchase/suppliers?company_id=CMP001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestGetSupplier(t *testing.T) {
	r, svc, ctx := setupPurchaseTest(t)
	sup := &domain.Supplier{CompanyID: "CMP001", Code: "SUP001", Name: "Test", TaxCode: "TX001"}
	require.NoError(t, svc.CreateSupplier(ctx, sup))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/purchase/suppliers/"+sup.ID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var got domain.Supplier
	json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, sup.ID, got.ID)
}

func TestGetSupplierNotFound(t *testing.T) {
	r, _, _ := setupPurchaseTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/purchase/suppliers/nonexistent", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}

func TestListSuppliers(t *testing.T) {
	r, svc, ctx := setupPurchaseTest(t)
	svc.CreateSupplier(ctx, &domain.Supplier{CompanyID: "CMP001", Code: "S1", Name: "Sup1", TaxCode: "T1"})
	svc.CreateSupplier(ctx, &domain.Supplier{CompanyID: "CMP001", Code: "S2", Name: "Sup2", TaxCode: "T2"})
	svc.CreateSupplier(ctx, &domain.Supplier{CompanyID: "CMP002", Code: "S3", Name: "Sup3", TaxCode: "T3"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/purchase/suppliers?company_id=CMP001", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var resp struct {
		Data  []domain.Supplier `json:"data"`
		Total int               `json:"total"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 2, resp.Total)
}

func TestUpdateSupplier(t *testing.T) {
	r, svc, ctx := setupPurchaseTest(t)
	sup := &domain.Supplier{CompanyID: "CMP001", Code: "SUP001", Name: "Old Name", TaxCode: "TX001"}
	require.NoError(t, svc.CreateSupplier(ctx, sup))

	body := fmt.Sprintf(`{"code":"SUP001","name":"New Name","tax_code":"TX001","company_id":"CMP001"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/purchase/suppliers/"+sup.ID+"?company_id=CMP001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	svc.GetSupplier(ctx, sup.ID)
}

func TestDeleteSupplier(t *testing.T) {
	r, svc, ctx := setupPurchaseTest(t)
	sup := &domain.Supplier{CompanyID: "CMP001", Code: "SUP001", Name: "Test", TaxCode: "TX001"}
	require.NoError(t, svc.CreateSupplier(ctx, sup))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/purchase/suppliers/"+sup.ID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)

	_, err := svc.GetSupplier(ctx, sup.ID)
	assert.Error(t, err)
}

// ─── Purchase Order ──────────────────────────────────────────────────────

func makeTestSupplier(t *testing.T, svc *service.PurchaseService, ctx context.Context, code string) *domain.Supplier {
	t.Helper()
	sup := &domain.Supplier{CompanyID: "CMP001", Code: code, Name: code + " Name", TaxCode: code + "-TX"}
	require.NoError(t, svc.CreateSupplier(ctx, sup))
	return sup
}

func TestCreatePO(t *testing.T) {
	r, svc, ctx := setupPurchaseTest(t)
	sup := makeTestSupplier(t, svc, ctx, "S001")

	body := fmt.Sprintf(`{"po_number":"PO-202607-0001","supplier_id":"%s","order_date":"2026-07-15T00:00:00Z","lines":[{"item_name":"Widget","unit":"pcs","quantity":10,"unit_price":50000,"account_id":"152","vat_account_id":"1331","vat_rate":10}]}`, sup.ID)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/purchase/orders?company_id=CMP001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	var po domain.PurchaseOrder
	err := json.Unmarshal(w.Body.Bytes(), &po)
	require.NoError(t, err)
	assert.NotEmpty(t, po.ID)
	assert.Equal(t, domain.POStatusDraft, po.Status)
	assert.Equal(t, 550000.0, po.TotalAmount)
}

func TestGetPO(t *testing.T) {
	r, svc, ctx := setupPurchaseTest(t)
	sup := makeTestSupplier(t, svc, ctx, "S002")
	po := &domain.PurchaseOrder{
		CompanyID: "CMP001", PONumber: "PO-001", SupplierID: sup.ID, OrderDate: time.Now(),
		Lines: []domain.POItem{{ItemName: "Item", Unit: "pcs", Quantity: 1, UnitPrice: 1000, AccountID: "152", VATAccountID: "1331"}},
	}
	require.NoError(t, svc.CreatePO(ctx, po))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/purchase/orders/"+po.ID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestListPOs(t *testing.T) {
	r, svc, ctx := setupPurchaseTest(t)
	sup := makeTestSupplier(t, svc, ctx, "S003")
	po := &domain.PurchaseOrder{
		CompanyID: "CMP001", PONumber: "PO-002", SupplierID: sup.ID, OrderDate: time.Now(),
		Lines: []domain.POItem{{ItemName: "Item", Unit: "pcs", Quantity: 1, UnitPrice: 1000, AccountID: "152", VATAccountID: "1331"}},
	}
	require.NoError(t, svc.CreatePO(ctx, po))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/purchase/orders?company_id=CMP001", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var resp struct {
		Data  []domain.PurchaseOrder `json:"data"`
		Total int                    `json:"total"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 1, resp.Total)
}

func TestApprovePO(t *testing.T) {
	r, svc, ctx := setupPurchaseTest(t)
	sup := makeTestSupplier(t, svc, ctx, "S004")
	po := &domain.PurchaseOrder{
		CompanyID: "CMP001", PONumber: "PO-003", SupplierID: sup.ID, OrderDate: time.Now(),
		Lines: []domain.POItem{{ItemName: "Item", Unit: "pcs", Quantity: 1, UnitPrice: 1000, AccountID: "152", VATAccountID: "1331"}},
	}
	require.NoError(t, svc.CreatePO(ctx, po))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/api/v1/purchase/orders/"+po.ID+"/approve", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	loaded, _ := svc.GetPO(ctx, po.ID)
	assert.Equal(t, domain.POStatusApproved, loaded.Status)
}

func TestCancelPO(t *testing.T) {
	r, svc, ctx := setupPurchaseTest(t)
	sup := makeTestSupplier(t, svc, ctx, "S005")
	po := &domain.PurchaseOrder{
		CompanyID: "CMP001", PONumber: "PO-004", SupplierID: sup.ID, OrderDate: time.Now(),
		Lines: []domain.POItem{{ItemName: "Item", Unit: "pcs", Quantity: 1, UnitPrice: 1000, AccountID: "152", VATAccountID: "1331"}},
	}
	require.NoError(t, svc.CreatePO(ctx, po))

	body := `{"reason":"changed mind"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/api/v1/purchase/orders/"+po.ID+"/cancel", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	loaded, _ := svc.GetPO(ctx, po.ID)
	assert.Equal(t, domain.POStatusCancelled, loaded.Status)
}

func TestClosePOInvalidTransition(t *testing.T) {
	r, svc, ctx := setupPurchaseTest(t)
	sup := makeTestSupplier(t, svc, ctx, "S006")
	po := &domain.PurchaseOrder{
		CompanyID: "CMP001", PONumber: "PO-CLOSE", SupplierID: sup.ID, OrderDate: time.Now(),
		Lines: []domain.POItem{{ItemName: "Item", Unit: "pcs", Quantity: 1, UnitPrice: 1000, AccountID: "152", VATAccountID: "1331"}},
	}
	require.NoError(t, svc.CreatePO(ctx, po))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/api/v1/purchase/orders/"+po.ID+"/close", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

// ─── GRN ─────────────────────────────────────────────────────────────────

func TestCreateGRN(t *testing.T) {
	r, svc, ctx := setupPurchaseTest(t)
	sup := makeTestSupplier(t, svc, ctx, "S010")
	po := &domain.PurchaseOrder{
		CompanyID: "CMP001", PONumber: "PO-GRN", SupplierID: sup.ID, OrderDate: time.Now(),
		Lines: []domain.POItem{{ItemName: "Item", Unit: "pcs", Quantity: 10, UnitPrice: 50000, AccountID: "152", VATAccountID: "1331"}},
	}
	require.NoError(t, svc.CreatePO(ctx, po))

	body := fmt.Sprintf(`{"grn_number":"GRN-202607-0001","po_id":"%s","receipt_date":"2026-07-16T00:00:00Z","lines":[{"po_line_id":"%s","item_name":"Item","unit":"pcs","quantity_received":10,"unit_price":50000,"line_total":500000}]}`, po.ID, po.Lines[0].ID)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/purchase/receipts?company_id=CMP001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

func TestPostGRN(t *testing.T) {
	r, svc, ctx := setupPurchaseTest(t)
	sup := makeTestSupplier(t, svc, ctx, "S011")
	po := &domain.PurchaseOrder{
		CompanyID: "CMP001", PONumber: "PO-GRN2", SupplierID: sup.ID, OrderDate: time.Now(),
		Lines: []domain.POItem{{ItemName: "Item", Unit: "pcs", Quantity: 10, UnitPrice: 50000, AccountID: "152", VATAccountID: "1331"}},
	}
	require.NoError(t, svc.CreatePO(ctx, po))
	grn := &domain.GRN{
		CompanyID: "CMP001", GRNNumber: "GRN-001", POID: po.ID, ReceiptDate: time.Now(),
		Lines: []domain.GRNItem{{POLineID: po.Lines[0].ID, ItemName: "Item", Unit: "pcs", QuantityReceived: 10, UnitPrice: 50000, LineTotal: 500000}},
	}
	require.NoError(t, svc.CreateGRN(ctx, grn))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/api/v1/purchase/receipts/"+grn.ID+"/post", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	loaded, _ := svc.GetGRN(ctx, grn.ID)
	assert.Equal(t, domain.GRNPosted, loaded.Status)
}

// ─── Invoice ─────────────────────────────────────────────────────────────

func TestCreateInvoice(t *testing.T) {
	r, svc, ctx := setupPurchaseTest(t)
	sup := makeTestSupplier(t, svc, ctx, "S020")

	body := fmt.Sprintf(`{"invoice_number":"INV-202607-0001","supplier_id":"%s","supplier_name":"Test Supplier","supplier_tax_code":"S020-TX","invoice_date":"2026-07-20T00:00:00Z","lines":[{"item_name":"Service","unit":"pcs","quantity":1,"unit_price":1000000,"account_id":"642","vat_account_id":"1331","vat_rate":10}]}`, sup.ID)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/purchase/invoices?company_id=CMP001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

func TestInvoiceWorkflow(t *testing.T) {
	r, svc, ctx := setupPurchaseTest(t)
	sup := makeTestSupplier(t, svc, ctx, "S030")
	inv := &domain.SupplierInvoice{
		CompanyID: "CMP001", InvoiceNumber: "INV-WF", SupplierID: sup.ID,
		SupplierName: "Test", SupplierTaxCode: "S030-TX", InvoiceDate: time.Now(),
		Lines: []domain.SupplierInvoiceLine{{ItemName: "Svc", Unit: "pcs", Quantity: 1, UnitPrice: 1000000, LineTotal: 1000000, AccountID: "642", VATAccountID: "1331", VATRate: 10}},
	}
	require.NoError(t, svc.CreateInvoice(ctx, inv))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/api/v1/purchase/invoices/"+inv.ID+"/verify", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PATCH", "/api/v1/purchase/invoices/"+inv.ID+"/post", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PATCH", "/api/v1/purchase/invoices/"+inv.ID+"/claim-vat", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	loaded, _ := svc.GetInvoice(ctx, inv.ID)
	assert.Equal(t, domain.InvoicePosted, loaded.Status)
	assert.Equal(t, domain.VATClaimed, loaded.VATDeductionStatus)
}

func TestPostInvoice(t *testing.T) {
	r, _, _ := setupPurchaseTest(t)

	// 1. Create supplier via API
	body := `{"code":"SUP100","name":"PostInvoice Supplier","tax_code":"TX100"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/purchase/suppliers?company_id=CMP001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
	var sup domain.Supplier
	err := json.Unmarshal(w.Body.Bytes(), &sup)
	require.NoError(t, err)

	// 2. Create PO via API
	poBody := fmt.Sprintf(`{"po_number":"PO-INV-001","supplier_id":"%s","order_date":"2026-07-15T00:00:00Z","lines":[{"item_name":"Widget","unit":"pcs","quantity":10,"unit_price":50000,"account_id":"152","vat_account_id":"1331","vat_rate":10}]}`, sup.ID)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/purchase/orders?company_id=CMP001", strings.NewReader(poBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
	var po domain.PurchaseOrder
	err = json.Unmarshal(w.Body.Bytes(), &po)
	require.NoError(t, err)

	// 3. Create GRN via API
	grnBody := fmt.Sprintf(`{"grn_number":"GRN-INV-001","po_id":"%s","receipt_date":"2026-07-16T00:00:00Z","lines":[{"po_line_id":"%s","item_name":"Widget","unit":"pcs","quantity_received":10,"unit_price":50000,"line_total":500000}]}`, po.ID, po.Lines[0].ID)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/purchase/receipts?company_id=CMP001", strings.NewReader(grnBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
	var grn domain.GRN
	err = json.Unmarshal(w.Body.Bytes(), &grn)
	require.NoError(t, err)

	// 4. Create invoice in draft via API
	invBody := fmt.Sprintf(`{"invoice_number":"INV-POST-001","supplier_id":"%s","supplier_name":"PostInvoice Supplier","supplier_tax_code":"TX100","invoice_date":"2026-07-20T00:00:00Z","lines":[{"item_name":"Service","unit":"pcs","quantity":1,"unit_price":1000000,"account_id":"642","vat_account_id":"1331","vat_rate":10}]}`, sup.ID)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/purchase/invoices?company_id=CMP001", strings.NewReader(invBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
	var inv domain.SupplierInvoice
	err = json.Unmarshal(w.Body.Bytes(), &inv)
	require.NoError(t, err)

	// 5. Verify invoice via API
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PATCH", "/api/v1/purchase/invoices/"+inv.ID+"/verify", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// 6. Post invoice via API
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PATCH", "/api/v1/purchase/invoices/"+inv.ID+"/post", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// 7. Get invoice and assert final state
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/purchase/invoices/"+inv.ID, nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	var loaded domain.SupplierInvoice
	err = json.Unmarshal(w.Body.Bytes(), &loaded)
	require.NoError(t, err)
	assert.Equal(t, domain.InvoicePosted, loaded.Status)
	assert.True(t, loaded.GLPosted)
	assert.NotNil(t, loaded.GLPostedAt)
}

// ─── AP Reports ──────────────────────────────────────────────────────────

func TestAPAgingReport(t *testing.T) {
	r, svc, ctx := setupPurchaseTest(t)
	sup := makeTestSupplier(t, svc, ctx, "S040")
	svc.CreateAPTransaction(ctx, &domain.APTransaction{
		CompanyID: "CMP001", SupplierID: sup.ID, TransactionType: domain.APTransInvoice,
		TransactionDate: time.Now(), Amount: 1000000, Currency: "VND",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/purchase/ap/aging?company_id=CMP001", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var report []domain.APAgingReport
	json.Unmarshal(w.Body.Bytes(), &report)
	assert.NotEmpty(t, report)
}

func TestAPSummary(t *testing.T) {
	r, svc, ctx := setupPurchaseTest(t)
	sup := makeTestSupplier(t, svc, ctx, "S050")
	svc.CreateAPTransaction(ctx, &domain.APTransaction{
		CompanyID: "CMP001", SupplierID: sup.ID, TransactionType: domain.APTransInvoice,
		TransactionDate: time.Now(), Amount: 5000000, Currency: "VND",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/purchase/ap/summary?company_id=CMP001", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var summary []domain.APSummary
	json.Unmarshal(w.Body.Bytes(), &summary)
	assert.NotEmpty(t, summary)
}
