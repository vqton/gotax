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

	ctx := context.Background()
	purRepo := repository.NewMemoryPurchaseRepo()

	accRepo := repository.NewMemoryAccountRepo()
	jeRepo := repository.NewMemoryJournalRepo()
	jeRepo.SetAccounts(accRepo.Accounts())
	perRepo := repository.NewMemoryPeriodRepo()
	userRepo := repository.NewMemoryUserRepo()
	auditRepo := repository.NewMemoryAuditLogRepo()
	rateRepo := repository.NewMemoryExchangeRateRepo()
	templateRepo := repository.NewMemoryClosingTemplateRepo()
	approvalRepo := repository.NewMemoryApprovalRepo()
	versionRepo := repository.NewMemoryAccountVersionRepo()
	mappingRepo := repository.NewMemoryAccountMappingRepo()
	analysisRepo := repository.NewMemoryAccountAnalysisRepo()
	ifrsRepo := repository.NewMemoryIFRSMappingRepo()
	refreshRepo := repository.NewMemoryRefreshTokenRepo()
	resetRepo := repository.NewMemoryPasswordResetTokenRepo()
	obRepo := repository.NewMemoryOpeningBalanceRepo()
	cashRepo := repository.NewMemoryCashRepo()
	gl := service.NewService(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo,
		approvalRepo, versionRepo, mappingRepo, analysisRepo, ifrsRepo, refreshRepo, resetRepo, obRepo, cashRepo)

	for _, acc := range []domain.Account{
		{Code: "331", Name: "Phai tra nguoi ban", Type: domain.AccountTypeLiability, IsActive: true},
		{Code: "1331", Name: "Thue GTGT duoc khau tru", Type: domain.AccountTypeAsset, IsActive: true},
		{Code: "152", Name: "Nguyen vat lieu", Type: domain.AccountTypeAsset, IsActive: true},
		{Code: "642", Name: "Chi phi quan ly doanh nghiep", Type: domain.AccountTypeExpense, IsActive: true},
		{Code: "3333", Name: "Thue nhap khau", Type: domain.AccountTypeLiability, IsActive: true},
		{Code: "33312", Name: "Thue GTGT hang nhap khau", Type: domain.AccountTypeLiability, IsActive: true},
	} {
		require.NoError(t, gl.CreateAccount(ctx, &acc))
	}

	purSvc := service.NewPurchaseService(purRepo, purRepo, purRepo, purRepo, purRepo, purRepo, purRepo, purRepo, gl)
	purH := NewPurchaseHandler(purSvc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("username", "testuser")
		c.Set("role", "admin")
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterPurchaseRoutes(r, purH, noopMW)

	return r, purSvc, ctx
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

func TestUpdatePO_AfterApprovedFails(t *testing.T) {
	r, svc, ctx := setupPurchaseTest(t)
	sup := makeTestSupplier(t, svc, ctx, "S006B")
	po := &domain.PurchaseOrder{
		CompanyID: "CMP001", PONumber: "PO-UPD", SupplierID: sup.ID, OrderDate: time.Now(),
		Lines: []domain.POItem{{ItemName: "Item", Unit: "pcs", Quantity: 1, UnitPrice: 1000, AccountID: "152", VATAccountID: "1331"}},
	}
	require.NoError(t, svc.CreatePO(ctx, po))
	require.NoError(t, svc.ApprovePO(ctx, po.ID, "user"))

	body := fmt.Sprintf(`{"po_number":"PO-UPD","supplier_id":"%s","order_date":"2026-07-15T00:00:00Z","lines":[{"item_name":"Item","unit":"pcs","quantity":5,"unit_price":2000,"account_id":"152","vat_account_id":"1331"}]}`, sup.ID)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/purchase/orders/"+po.ID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
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

func TestCancelGRN_AfterPostedFails(t *testing.T) {
	r, svc, ctx := setupPurchaseTest(t)
	sup := makeTestSupplier(t, svc, ctx, "S011B")
	po := &domain.PurchaseOrder{
		CompanyID: "CMP001", PONumber: "PO-GRN3", SupplierID: sup.ID, OrderDate: time.Now(),
		Lines: []domain.POItem{{ItemName: "Item", Unit: "pcs", Quantity: 10, UnitPrice: 50000, AccountID: "152", VATAccountID: "1331"}},
	}
	require.NoError(t, svc.CreatePO(ctx, po))
	grn := &domain.GRN{
		CompanyID: "CMP001", GRNNumber: "GRN-002", POID: po.ID, ReceiptDate: time.Now(),
		Lines: []domain.GRNItem{{POLineID: po.Lines[0].ID, ItemName: "Item", Unit: "pcs", QuantityReceived: 10, UnitPrice: 50000, LineTotal: 500000}},
	}
	require.NoError(t, svc.CreateGRN(ctx, grn))
	require.NoError(t, svc.PostGRN(ctx, grn.ID))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/api/v1/purchase/receipts/"+grn.ID+"/cancel", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)

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

func TestPostInvoice_WithoutVerifyFails(t *testing.T) {
	r, svc, ctx := setupPurchaseTest(t)
	sup := makeTestSupplier(t, svc, ctx, "S021")
	inv := &domain.SupplierInvoice{
		CompanyID: "CMP001", InvoiceNumber: "INV-NOVFY", SupplierID: sup.ID,
		SupplierName: "Test", SupplierTaxCode: "S021-TX", InvoiceDate: time.Now(),
		Lines: []domain.SupplierInvoiceLine{{ItemName: "Svc", Unit: "pcs", Quantity: 1, UnitPrice: 1000000, LineTotal: 1000000, AccountID: "642", VATAccountID: "1331", VATRate: 10}},
	}
	require.NoError(t, svc.CreateInvoice(ctx, inv))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/api/v1/purchase/invoices/"+inv.ID+"/post", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)

	loaded, _ := svc.GetInvoice(ctx, inv.ID)
	assert.Equal(t, domain.InvoiceDraft, loaded.Status)
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

// ─── Doubtful Debt Provisions ───────────────────────────────────────────

func TestProvisionCalculateHandler(t *testing.T) {
	r, svc, ctx := setupPurchaseTest(t)
	sup := &domain.Supplier{CompanyID: "CMP001", Code: "SUP001", Name: "Test", TaxCode: "TX001"}
	require.NoError(t, svc.CreateSupplier(ctx, sup))
	asOf := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	require.NoError(t, svc.CreateAPTransaction(ctx, &domain.APTransaction{
		CompanyID: "CMP001", SupplierID: sup.ID, TransactionType: domain.APTransPrepayment,
		TransactionDate: asOf.AddDate(0, -18, 0), Amount: 2000, Currency: "VND",
	}))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/purchase/provisions/calculate?company_id=CMP001&as_of_date=2026-08-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var resp struct {
		AsOfDate string                                   `json:"as_of_date"`
		Lines    []domain.DoubtfulDebtProvisionLine       `json:"lines"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Len(t, resp.Lines, 1)
	assert.Equal(t, 0.50, resp.Lines[0].RatePct)
	assert.Equal(t, 1000.0, resp.Lines[0].ProvisionAmount)
}

func TestProvisionCalculateHandlerNone(t *testing.T) {
	r, _, _ := setupPurchaseTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/purchase/provisions/calculate?company_id=CMP001&as_of_date=2026-08-01", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

func TestProvisionCalculateHandlerBadDate(t *testing.T) {
	r, _, _ := setupPurchaseTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/purchase/provisions/calculate?company_id=CMP001&as_of_date=banana", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestCreateProvisionHandler(t *testing.T) {
	r, svc, ctx := setupPurchaseTest(t)
	sup := &domain.Supplier{CompanyID: "CMP001", Code: "SUP001", Name: "Test", TaxCode: "TX001"}
	require.NoError(t, svc.CreateSupplier(ctx, sup))

	body := fmt.Sprintf(`{
		"as_of_date":"2026-08-01",
		"lines":[{
			"supplier_id":%q,"supplier_name":"Test","outstanding_amount":1000,
			"age_months":18,"rate_pct":0.5
		}]
	}`, sup.ID)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/purchase/provisions?company_id=CMP001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	var prov domain.DoubtfulDebtProvision
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &prov))
	assert.NotEmpty(t, prov.ID)
	assert.Equal(t, 500.0, prov.TotalProvision)

	get := httptest.NewRecorder()
	greq, _ := http.NewRequest("GET", "/api/v1/purchase/provisions/"+prov.ID, nil)
	r.ServeHTTP(get, greq)
	assert.Equal(t, 200, get.Code)
	var got domain.DoubtfulDebtProvision
	require.NoError(t, json.Unmarshal(get.Body.Bytes(), &got))
	assert.Len(t, got.Lines, 1)

	list := httptest.NewRecorder()
	lreq, _ := http.NewRequest("GET", "/api/v1/purchase/provisions?company_id=CMP001", nil)
	r.ServeHTTP(list, lreq)
	assert.Equal(t, 200, list.Code)
}

func TestCreateProvisionHandlerInvalid(t *testing.T) {
	r, _, _ := setupPurchaseTest(t)
	body := `{"as_of_date":"2026-08-01","lines":[]}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/purchase/provisions?company_id=CMP001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

// ─── Regulatory Report Handlers ─────────────────────────────────────────

func TestPurchaseLedgerHandler(t *testing.T) {
	r, svc, ctx := setupPurchaseTest(t)
	sup := &domain.Supplier{CompanyID: "CMP001", Code: "SUP001", Name: "Test", TaxCode: "TX001"}
	require.NoError(t, svc.CreateSupplier(ctx, sup))
	po := &domain.PurchaseOrder{
		CompanyID: "CMP001", PONumber: "PO-1", SupplierID: sup.ID,
		OrderDate: time.Now(), Currency: "VND",
		Lines: []domain.POItem{{
			ItemName: "Widget", Unit: "pcs", Quantity: 10, UnitPrice: 100,
			VATRate: 10, VATType: domain.VAT10, AccountID: "152", VATAccountID: "1331",
		}},
	}
	require.NoError(t, svc.CreatePO(ctx, po))
	grn := &domain.GRN{
		CompanyID: "CMP001", GRNNumber: "GRN-1", POID: po.ID, ReceiptDate: time.Now(),
		Lines: []domain.GRNItem{{
			POLineID: po.Lines[0].ID, ItemName: "Widget", Unit: "pcs", QuantityReceived: 10, UnitPrice: 100, LineTotal: 1000,
		}},
	}
	require.NoError(t, svc.CreateGRN(ctx, grn))
	require.NoError(t, svc.PostGRN(ctx, grn.ID))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/purchase/reports/s01-dn?company_id=CMP001", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var rpt domain.PurchaseLedgerReport
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &rpt))
	require.Len(t, rpt.Rows, 1)
	assert.Equal(t, 1000.0, rpt.Increase)
}

func TestVATInputReportHandler(t *testing.T) {
	r, svc, ctx := setupPurchaseTest(t)
	sup := &domain.Supplier{CompanyID: "CMP001", Code: "SUP001", Name: "Test", TaxCode: "TX001"}
	require.NoError(t, svc.CreateSupplier(ctx, sup))
	inv := &domain.SupplierInvoice{
		CompanyID: "CMP001", InvoiceNumber: "INV-1", SupplierID: sup.ID,
		SupplierName: sup.Name, SupplierTaxCode: sup.TaxCode, InvoiceDate: time.Now(),
		Currency: "VND", VATDeductionStatus: domain.VATPending,
		Lines: []domain.SupplierInvoiceLine{{
			ItemName: "Widget", Unit: "pcs", Quantity: 2, UnitPrice: 1000,
			VATRate: 10, VATType: domain.VAT10, AccountID: "152", VATAccountID: "1331",
		}},
	}
	require.NoError(t, svc.CreateInvoice(ctx, inv))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/purchase/reports/vat-input?company_id=CMP001", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var rpt domain.VATInputReport
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &rpt))
	require.Len(t, rpt.Rows, 1)
	assert.Equal(t, 2000.0, rpt.Rows[0].Subtotal)
	assert.Equal(t, 200.0, rpt.Rows[0].VATAmount)
}

func TestUninvoicedReceiptsHandler(t *testing.T) {
	r, svc, ctx := setupPurchaseTest(t)
	sup := &domain.Supplier{CompanyID: "CMP001", Code: "SUP001", Name: "Test", TaxCode: "TX001"}
	require.NoError(t, svc.CreateSupplier(ctx, sup))
	po := &domain.PurchaseOrder{
		CompanyID: "CMP001", PONumber: "PO-1", SupplierID: sup.ID,
		OrderDate: time.Now(), Currency: "VND",
		Lines: []domain.POItem{{
			ItemName: "Widget", Unit: "pcs", Quantity: 5, UnitPrice: 100,
			VATRate: 10, VATType: domain.VAT10, AccountID: "152", VATAccountID: "1331",
		}},
	}
	require.NoError(t, svc.CreatePO(ctx, po))
	grn := &domain.GRN{
		CompanyID: "CMP001", GRNNumber: "GRN-1", POID: po.ID, ReceiptDate: time.Now(),
		Lines: []domain.GRNItem{{
			POLineID: po.Lines[0].ID, ItemName: "Widget", Unit: "pcs", QuantityReceived: 5, UnitPrice: 100, LineTotal: 500,
		}},
	}
	require.NoError(t, svc.CreateGRN(ctx, grn))
	require.NoError(t, svc.PostGRN(ctx, grn.ID))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/purchase/reports/uninvoiced-receipts?company_id=CMP001", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var rows []domain.UninvoicedReceiptRow
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &rows))
	require.Len(t, rows, 1)
	assert.Equal(t, grn.GRNNumber, rows[0].GRNNumber)
}

func TestSupplierLedgerHandlerRequiresSupplier(t *testing.T) {
	r, _, _ := setupPurchaseTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/purchase/reports/s02-dn?company_id=CMP001", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

// ─── Requisition ─────────────────────────────────────────────────────────

func TestCreateRequisition(t *testing.T) {
	r, _, _ := setupPurchaseTest(t)
	body := `{"requisition_number":"REQ-H1","requester_id":"u1","requester_name":"Alice","lines":[{"item_name":"Widget","unit":"pcs","quantity":10,"estimated_price":1000,"account_id":"152"}]}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/purchase/requisitions?company_id=CMP001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
	var out domain.PurchaseRequisition
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &out))
	require.NotEmpty(t, out.ID)
	assert.InDelta(t, 10000, out.TotalEstimated, 0.001)
}

func TestRequisitionWorkflow(t *testing.T) {
	r, svc, ctx := setupPurchaseTest(t)
	sup := makeTestSupplier(t, svc, ctx, "SRQ1")
	req := &domain.PurchaseRequisition{
		CompanyID: "CMP001", RequisitionNumber: "REQ-H2", RequesterID: "u1",
		RequesterName: "Alice", CreatedBy: "u1",
		Lines: []domain.RequisitionItem{{ItemName: "W", Unit: "pcs", Quantity: 1, EstimatedPrice: 100, AccountID: "152"}},
	}
	require.NoError(t, svc.CreateRequisition(ctx, req))

	for _, tc := range []struct {
		method string
		path   string
	}{
		{"PATCH", "/api/v1/purchase/requisitions/" + req.ID + "/submit"},
		{"PATCH", "/api/v1/purchase/requisitions/" + req.ID + "/approve"},
		{"POST", "/api/v1/purchase/requisitions/" + req.ID + "/convert-to-po?supplier_id=" + sup.ID},
	} {
		w := httptest.NewRecorder()
		rq, _ := http.NewRequest(tc.method, tc.path, nil)
		r.ServeHTTP(w, rq)
		assert.Equal(t, 200, w.Code, "%s %s", tc.method, tc.path)
	}

	loaded, _ := svc.GetRequisition(ctx, req.ID)
	assert.Equal(t, domain.ReqOrdered, loaded.Status)
}

func TestGetRequisitionNotFound(t *testing.T) {
	r, _, _ := setupPurchaseTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/purchase/requisitions/nope", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

func TestRequisitionApproveWithoutSubmitFails(t *testing.T) {
	r, svc, ctx := setupPurchaseTest(t)
	req := &domain.PurchaseRequisition{
		CompanyID: "CMP001", RequisitionNumber: "REQ-H3", RequesterID: "u1",
		RequesterName: "Alice", CreatedBy: "u1",
		Lines: []domain.RequisitionItem{{ItemName: "W", Unit: "pcs", Quantity: 1, EstimatedPrice: 100, AccountID: "152"}},
	}
	require.NoError(t, svc.CreateRequisition(ctx, req))

	w := httptest.NewRecorder()
	rq, _ := http.NewRequest("PATCH", "/api/v1/purchase/requisitions/"+req.ID+"/approve", nil)
	r.ServeHTTP(w, rq)
	assert.Equal(t, 400, w.Code)
}

// ─── Return / Credit Note (P2-2) ─────────────────────────────────────────

func TestCreateReturnGRNHandler(t *testing.T) {
	r, svc, ctx := setupPurchaseTest(t)
	sup := &domain.Supplier{CompanyID: "CMP001", Code: "R-SUP", Name: "R Sup", TaxCode: "R-TX"}
	require.NoError(t, svc.CreateSupplier(ctx, sup))
	po := &domain.PurchaseOrder{
		CompanyID: "CMP001", PONumber: "PO-R1", SupplierID: sup.ID, OrderDate: time.Now(), Currency: "VND",
		Lines: []domain.POItem{{ItemName: "W", Unit: "pcs", Quantity: 10, UnitPrice: 1000, VATRate: 10, VATType: domain.VAT10, AccountID: "152", VATAccountID: "1331"}},
	}
	require.NoError(t, svc.CreatePO(ctx, po))
	grn := &domain.GRN{
		CompanyID: "CMP001", GRNNumber: "GRN-R1", POID: po.ID, ReceiptDate: time.Now(),
		Lines: []domain.GRNItem{{POLineID: po.Lines[0].ID, ItemName: "W", Unit: "pcs", QuantityReceived: 10}},
	}
	require.NoError(t, svc.CreateGRN(ctx, grn))
	require.NoError(t, svc.PostGRN(ctx, grn.ID))

	body := fmt.Sprintf(`{"grn_number":"GRN-RR1","return_of_grn_id":"%s","receipt_date":"%s","lines":[{"po_line_id":"%s","item_name":"W","unit":"pcs","quantity_received":2}]}`,
		grn.ID, time.Now().Format("2006-01-02T15:04:05Z07:00"), po.Lines[0].ID)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/purchase/returns?company_id=CMP001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
	var out domain.GRN
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &out))
	assert.Equal(t, domain.GRNPosted, out.Status)
}

func TestCreateCreditNoteHandler(t *testing.T) {
	r, svc, ctx := setupPurchaseTest(t)
	sup := &domain.Supplier{CompanyID: "CMP001", Code: "C-SUP", Name: "C Sup", TaxCode: "C-TX"}
	require.NoError(t, svc.CreateSupplier(ctx, sup))
	po := &domain.PurchaseOrder{
		CompanyID: "CMP001", PONumber: "PO-C1", SupplierID: sup.ID, OrderDate: time.Now(), Currency: "VND",
		Lines: []domain.POItem{{ItemName: "W", Unit: "pcs", Quantity: 10, UnitPrice: 1000, VATRate: 10, VATType: domain.VAT10, AccountID: "152", VATAccountID: "1331"}},
	}
	require.NoError(t, svc.CreatePO(ctx, po))
	grn := &domain.GRN{
		CompanyID: "CMP001", GRNNumber: "GRN-C1", POID: po.ID, ReceiptDate: time.Now(),
		Lines: []domain.GRNItem{{POLineID: po.Lines[0].ID, ItemName: "W", Unit: "pcs", QuantityReceived: 10}},
	}
	require.NoError(t, svc.CreateGRN(ctx, grn))
	require.NoError(t, svc.PostGRN(ctx, grn.ID))
	inv := &domain.SupplierInvoice{
		CompanyID: "CMP001", InvoiceNumber: "INV-C1", SupplierID: sup.ID, POID: po.ID, GRNID: grn.ID, InvoiceDate: time.Now(),
		Lines: []domain.SupplierInvoiceLine{{POLineID: po.Lines[0].ID, ItemName: "W", Unit: "pcs", Quantity: 10, UnitPrice: 1000, VATRate: 10, VATType: domain.VAT10, AccountID: "152", VATAccountID: "1331"}},
	}
	require.NoError(t, svc.CreateInvoice(ctx, inv))
	require.NoError(t, svc.VerifyInvoice(ctx, inv.ID))
	require.NoError(t, svc.PostInvoice(ctx, inv.ID))

	body := fmt.Sprintf(`{"invoice_number":"CN-C1","original_invoice_id":"%s","invoice_type":"credit_note","invoice_date":"%s","lines":[{"item_name":"W","unit":"pcs","quantity":2,"unit_price":1000,"vat_rate":10,"vat_type":"VAT_10","account_id":"152","vat_account_id":"1331"}]}`,
		inv.ID, time.Now().Format("2006-01-02T15:04:05Z07:00"))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/purchase/credit-notes?company_id=CMP001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
	var out domain.SupplierInvoice
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &out))
	assert.True(t, out.TotalAmount < 0)
	assert.Equal(t, inv.ID, out.OriginalInvoiceID)
}

func TestCreateImportInvoiceHandler(t *testing.T) {
	r, svc, ctx := setupPurchaseTest(t)
	sup := &domain.Supplier{CompanyID: "CMP001", Code: "I-SUP", Name: "I Sup", TaxCode: "I-TX", SupplierType: domain.SupplierTypeImport}
	require.NoError(t, svc.CreateSupplier(ctx, sup))
	po := &domain.PurchaseOrder{
		CompanyID: "CMP001", PONumber: "PO-I1", SupplierID: sup.ID, OrderDate: time.Now(), Currency: "VND",
		Lines: []domain.POItem{{ItemName: "W", Unit: "pcs", Quantity: 10, UnitPrice: 1000, VATRate: 10, VATType: domain.VAT10, AccountID: "152", VATAccountID: "1331"}},
	}
	require.NoError(t, svc.CreatePO(ctx, po))
	grn := &domain.GRN{
		CompanyID: "CMP001", GRNNumber: "GRN-I1", POID: po.ID, ReceiptDate: time.Now(),
		Lines: []domain.GRNItem{{POLineID: po.Lines[0].ID, ItemName: "W", Unit: "pcs", QuantityReceived: 10}},
	}
	require.NoError(t, svc.CreateGRN(ctx, grn))
	require.NoError(t, svc.PostGRN(ctx, grn.ID))

	body := fmt.Sprintf(`{"invoice_number":"IMP-I1","supplier_id":"%s","invoice_type":"import","po_id":"%s","grn_id":"%s","invoice_date":"%s","import_duty":500000,"import_vat":1050000,"customs_declaration_number":"CD-2025-1","lines":[{"po_line_id":"%s","item_name":"W","unit":"pcs","quantity":10,"unit_price":1000,"vat_rate":10,"vat_type":"VAT_10","account_id":"152","vat_account_id":"1331"}]}`,
		sup.ID, po.ID, grn.ID, time.Now().Format("2006-01-02T15:04:05Z07:00"), po.Lines[0].ID)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/purchase/invoices?company_id=CMP001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
	var inv domain.SupplierInvoice
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &inv))
	assert.InDelta(t, 500000, inv.ImportDuty, 0.001)

	for _, path := range []string{
		"/api/v1/purchase/invoices/" + inv.ID + "/verify",
		"/api/v1/purchase/invoices/" + inv.ID + "/post",
	} {
		w := httptest.NewRecorder()
		rq, _ := http.NewRequest("PATCH", path, nil)
		r.ServeHTTP(w, rq)
		assert.Equal(t, 200, w.Code, path)
	}
}
