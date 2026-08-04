package handler

import (
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

type saleTestSetup struct {
	r    *gin.Engine
	svc  *service.SaleService
	cust *domain.Customer
	so   *domain.SalesOrder
}

func setupSale(t *testing.T) *saleTestSetup {
	t.Helper()
	gin.SetMode(gin.TestMode)

	custRepo := repository.NewMemorySaleRepo()
	soRepo := repository.NewMemorySaleRepo()
	dnRepo := repository.NewMemorySaleRepo()
	invRepo := repository.NewMemorySaleRepo()
	rcptRepo := repository.NewMemorySaleRepo()
	cnRepo := repository.NewMemorySaleRepo()
	artRepo := repository.NewMemorySaleRepo()
	sqRepo := repository.NewMemorySaleRepo()

	svc := service.NewSaleService(custRepo, soRepo, dnRepo, invRepo, rcptRepo, cnRepo, artRepo, sqRepo, nil)
	h := NewSaleHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "sale-user")
		c.Set("username", "saleuser")
		c.Set("role", "admin")
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterSaleRoutes(r, h, noopMW)

	cust := &domain.Customer{
		CompanyID: "CMP001",
		Code:      "CUST001",
		Name:      "Test Customer",
		TaxCode:   "1234567890",
		Currency:  "VND",
		Status:    domain.CustomerActive,
	}
	err := svc.CreateCustomer(nil, cust)
	require.NoError(t, err)

	return &saleTestSetup{r: r, svc: svc, cust: cust}
}

func setupSaleWithSO(t *testing.T, ts *saleTestSetup) *domain.SalesOrder {
	t.Helper()
	so := &domain.SalesOrder{
		CompanyID:  "CMP001",
		SONumber:   "SO-202607-00001",
		CustomerID: ts.cust.ID,
		OrderDate:  time.Now().UTC(),
		Currency:   "VND",
		Status:     domain.SODraft,
		Lines: []domain.SOLine{{
			ID:             "SO-LINE-001",
			ItemName:       "Widget",
			Unit:           "pc",
			Quantity:       10,
			UnitPrice:      100000,
			VATRate:        10,
			VATType:        domain.VAT10,
			RevenueAccount: "511",
			VATAccountID:   "3331",
		}},
	}
	err := ts.svc.CreateSO(nil, so)
	require.NoError(t, err)
	return so
}

func setupSaleWithCustomer(t *testing.T, ts *saleTestSetup) *domain.Customer {
	t.Helper()
	cust := &domain.Customer{
		CompanyID: "CMP001",
		Code:      "CUST002",
		Name:      "Another Customer",
		TaxCode:   "0987654321",
		Currency:  "VND",
		Status:    domain.CustomerActive,
	}
	err := ts.svc.CreateCustomer(nil, cust)
	require.NoError(t, err)
	return cust
}

func postJSON(ts *saleTestSetup, path, body string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	return w
}

func getJSON(ts *saleTestSetup, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", path, nil)
	ts.r.ServeHTTP(w, req)
	return w
}

func putJSON(ts *saleTestSetup, path, body string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	return w
}

func patchJSON(ts *saleTestSetup, path, body string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	return w
}

func delReq(ts *saleTestSetup, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", path, nil)
	ts.r.ServeHTTP(w, req)
	return w
}

// ─── Customer ──────────────────────────────────────────────────────────────

func TestSaleCreateCustomer(t *testing.T) {
	ts := setupSale(t)
	body := `{"code":"NEWCUST","name":"New Customer","tax_code":"1111111111","currency":"VND"}`
	w := postJSON(ts, "/api/v1/sale/customers?company_id=CMP001", body)
	assert.Equal(t, 201, w.Code)
	var cust domain.Customer
	json.Unmarshal(w.Body.Bytes(), &cust)
	assert.NotEmpty(t, cust.ID)
	assert.Equal(t, "NEWCUST", cust.Code)
}

func TestSaleCreateCustomer_EmptyCompany(t *testing.T) {
	ts := setupSale(t)
	body := `{"code":"NEW","name":"No Company","tax_code":"1111111111"}`
	w := postJSON(ts, "/api/v1/sale/customers", body)
	assert.Equal(t, 400, w.Code)
}

func TestSaleCreateCustomer_DuplicateCode(t *testing.T) {
	ts := setupSale(t)
	body := `{"code":"CUST001","name":"Dup","tax_code":"1111111111","currency":"VND"}`
	w := postJSON(ts, "/api/v1/sale/customers?company_id=CMP001", body)
	assert.Equal(t, 400, w.Code)
}

func TestSaleCreateCustomer_ValidationError(t *testing.T) {
	ts := setupSale(t)
	body := `{"code":"","name":"","tax_code":""}`
	w := postJSON(ts, "/api/v1/sale/customers?company_id=CMP001", body)
	assert.Equal(t, 400, w.Code)
}

func TestSaleGetCustomer(t *testing.T) {
	ts := setupSale(t)
	w := getJSON(ts, "/api/v1/sale/customers/"+ts.cust.ID+"?company_id=CMP001")
	assert.Equal(t, 200, w.Code)
	var cust domain.Customer
	json.Unmarshal(w.Body.Bytes(), &cust)
	assert.Equal(t, "CUST001", cust.Code)
}

func TestSaleGetCustomer_CrossTenant(t *testing.T) {
	ts := setupSale(t)
	w := getJSON(ts, "/api/v1/sale/customers/"+ts.cust.ID+"?company_id=OTHER")
	assert.Equal(t, 404, w.Code)
}

func TestSaleGetCustomer_NotFound(t *testing.T) {
	ts := setupSale(t)
	w := getJSON(ts, "/api/v1/sale/customers/nonexistent?company_id=CMP001")
	assert.Equal(t, 404, w.Code)
}

func TestSaleListCustomers(t *testing.T) {
	ts := setupSale(t)
	setupSaleWithCustomer(t, ts)
	w := getJSON(ts, "/api/v1/sale/customers?company_id=CMP001")
	assert.Equal(t, 200, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 2, int(resp["total"].(float64)))
}

func TestSaleUpdateCustomer(t *testing.T) {
	ts := setupSale(t)
	body := `{"name":"Updated Name","tax_code":"9999999999","code":"CUST001","currency":"VND"}`
	w := putJSON(ts, "/api/v1/sale/customers/"+ts.cust.ID+"?company_id=CMP001", body)
	assert.Equal(t, 200, w.Code)
	var cust domain.Customer
	json.Unmarshal(w.Body.Bytes(), &cust)
	assert.Equal(t, "Updated Name", cust.Name)
}

func TestSaleUpdateCustomer_CrossTenant(t *testing.T) {
	ts := setupSale(t)
	body := `{"name":"Hacker","code":"CUST001","tax_code":"0000000000"}`
	w := putJSON(ts, "/api/v1/sale/customers/"+ts.cust.ID+"?company_id=OTHER", body)
	assert.Equal(t, 404, w.Code)
}

func TestSaleDeleteCustomer(t *testing.T) {
	ts := setupSale(t)
	w := delReq(ts, "/api/v1/sale/customers/"+ts.cust.ID+"?company_id=CMP001")
	assert.Equal(t, 204, w.Code)
}

func TestSaleDeleteCustomer_CrossTenant(t *testing.T) {
	ts := setupSale(t)
	w := delReq(ts, "/api/v1/sale/customers/"+ts.cust.ID+"?company_id=OTHER")
	assert.Equal(t, 404, w.Code)
}

// ─── Sales Order ───────────────────────────────────────────────────────────

func TestSaleCreateSO(t *testing.T) {
	ts := setupSale(t)
	body := fmt.Sprintf(`{
		"company_id":"CMP001","so_number":"SO-TEST-001","customer_id":"%s",
		"order_date":"%s","currency":"VND",
		"lines":[{"item_name":"Test Item","unit":"pc","quantity":5,"unit_price":20000,"vat_rate":10,"vat_type":"VAT10","revenue_account_id":"511","vat_account_id":"3331"}]
	}`, ts.cust.ID, time.Now().UTC().Format(time.RFC3339))
	w := postJSON(ts, "/api/v1/sale/orders?company_id=CMP001", body)
	assert.Equal(t, 201, w.Code)
	var so domain.SalesOrder
	json.Unmarshal(w.Body.Bytes(), &so)
	assert.NotEmpty(t, so.ID)
	assert.Equal(t, "SO-TEST-001", so.SONumber)
	assert.Equal(t, 1, len(so.Lines))
}

func TestSaleCreateSO_DuplicateNumber(t *testing.T) {
	ts := setupSale(t)
	so := setupSaleWithSO(t, ts)
	body := fmt.Sprintf(`{
		"company_id":"CMP001","so_number":"%s","customer_id":"%s",
		"order_date":"%s","currency":"VND",
		"lines":[{"item_name":"Dup","unit":"pc","quantity":1,"unit_price":100,"vat_rate":10,"vat_type":"VAT10","revenue_account_id":"511","vat_account_id":"3331"}]
	}`, so.SONumber, ts.cust.ID, time.Now().UTC().Format(time.RFC3339))
	w := postJSON(ts, "/api/v1/sale/orders?company_id=CMP001", body)
	assert.Equal(t, 400, w.Code)
}

func TestSaleCreateSO_ValidationError(t *testing.T) {
	ts := setupSale(t)
	body := `{"company_id":"CMP001","so_number":"","customer_id":""}`
	w := postJSON(ts, "/api/v1/sale/orders?company_id=CMP001", body)
	assert.Equal(t, 400, w.Code)
}

func TestSaleGetSO(t *testing.T) {
	ts := setupSale(t)
	so := setupSaleWithSO(t, ts)
	w := getJSON(ts, "/api/v1/sale/orders/"+so.ID+"?company_id=CMP001")
	assert.Equal(t, 200, w.Code)
	var res domain.SalesOrder
	json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, so.SONumber, res.SONumber)
}

func TestSaleGetSO_CrossTenant(t *testing.T) {
	ts := setupSale(t)
	so := setupSaleWithSO(t, ts)
	w := getJSON(ts, "/api/v1/sale/orders/"+so.ID+"?company_id=OTHER")
	assert.Equal(t, 404, w.Code)
}

func TestSaleListSOs(t *testing.T) {
	ts := setupSale(t)
	setupSaleWithSO(t, ts)
	w := getJSON(ts, "/api/v1/sale/orders?company_id=CMP001")
	assert.Equal(t, 200, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 1, int(resp["total"].(float64)))
}

func TestSaleUpdateSO(t *testing.T) {
	ts := setupSale(t)
	so := setupSaleWithSO(t, ts)
	body := fmt.Sprintf(`{
		"so_number":"SO-UPDATED","customer_id":"%s","order_date":"%s","currency":"VND",
		"lines":[{"item_name":"Updated","unit":"pc","quantity":3,"unit_price":50000,"vat_rate":10,"vat_type":"VAT10","revenue_account_id":"511","vat_account_id":"3331"}]
	}`, ts.cust.ID, time.Now().UTC().Format(time.RFC3339))
	w := putJSON(ts, "/api/v1/sale/orders/"+so.ID+"?company_id=CMP001", body)
	assert.Equal(t, 200, w.Code)
}

func TestSaleUpdateSO_CrossTenant(t *testing.T) {
	ts := setupSale(t)
	so := setupSaleWithSO(t, ts)
	w := putJSON(ts, "/api/v1/sale/orders/"+so.ID+"?company_id=OTHER", `{}`)
	assert.Equal(t, 404, w.Code)
}

func TestSaleApproveSO(t *testing.T) {
	ts := setupSale(t)
	so := setupSaleWithSO(t, ts)
	w := patchJSON(ts, "/api/v1/sale/orders/"+so.ID+"/approve?company_id=CMP001", "")
	assert.Equal(t, 200, w.Code)
}

func TestSaleApproveSO_CrossTenant(t *testing.T) {
	ts := setupSale(t)
	so := setupSaleWithSO(t, ts)
	w := patchJSON(ts, "/api/v1/sale/orders/"+so.ID+"/approve?company_id=OTHER", "")
	assert.Equal(t, 404, w.Code)
}

func TestSaleCancelSO(t *testing.T) {
	ts := setupSale(t)
	so := setupSaleWithSO(t, ts)
	w := patchJSON(ts, "/api/v1/sale/orders/"+so.ID+"/cancel?company_id=CMP001", `{"reason":"test"}`)
	assert.Equal(t, 200, w.Code)
}

func TestSaleCancelSO_CrossTenant(t *testing.T) {
	ts := setupSale(t)
	so := setupSaleWithSO(t, ts)
	w := patchJSON(ts, "/api/v1/sale/orders/"+so.ID+"/cancel?company_id=OTHER", `{"reason":"hack"}`)
	assert.Equal(t, 404, w.Code)
}

func dnBody(soID, dnNum, deliveryDate string) string {
	return fmt.Sprintf(`{
		"company_id":"CMP001","dn_number":"%s","so_id":"%s","delivery_date":"%s",
		"lines":[{"so_line_id":"SO-LINE-001","item_name":"Delivered","unit":"pc","quantity_delivered":10,"unit_price":100000}]
	}`, dnNum, soID, deliveryDate)
}

// ─── Delivery Note ─────────────────────────────────────────────────────────

func TestSaleCreateDN(t *testing.T) {
	ts := setupSale(t)
	so := setupSaleWithSO(t, ts)
	w := postJSON(ts, "/api/v1/sale/deliveries?company_id=CMP001", dnBody(so.ID, "DN-TEST-001", time.Now().UTC().Format(time.RFC3339)))
	assert.Equal(t, 201, w.Code)
}

func TestSaleCreateDN_DuplicateNumber(t *testing.T) {
	ts := setupSale(t)
	so := setupSaleWithSO(t, ts)
	b := dnBody(so.ID, "DN-TEST-001", time.Now().UTC().Format(time.RFC3339))
	postJSON(ts, "/api/v1/sale/deliveries?company_id=CMP001", b)
	w := postJSON(ts, "/api/v1/sale/deliveries?company_id=CMP001", b)
	assert.Equal(t, 400, w.Code)
}

func TestSalePostDN(t *testing.T) {
	ts := setupSale(t)
	so := setupSaleWithSO(t, ts)
	var dn domain.DeliveryNote
	w := postJSON(ts, "/api/v1/sale/deliveries?company_id=CMP001", dnBody(so.ID, "DN-POST", time.Now().UTC().Format(time.RFC3339)))
	json.Unmarshal(w.Body.Bytes(), &dn)

	w = patchJSON(ts, "/api/v1/sale/deliveries/"+dn.ID+"/post?company_id=CMP001", "")
	assert.Equal(t, 200, w.Code)
}

// ─── Invoice ───────────────────────────────────────────────────────────────

func TestSaleCreateInvoice(t *testing.T) {
	ts := setupSale(t)
	body := fmt.Sprintf(`{
		"company_id":"CMP001","invoice_number":"INV-TEST-001","customer_id":"%s","customer_name":"Test","customer_tax_code":"1234567890",
		"invoice_date":"%s","currency":"VND","due_date":"%s",
		"lines":[{"item_name":"Inv Item","unit":"pc","quantity":2,"unit_price":100000,"vat_rate":10,"vat_type":"VAT10","revenue_account_id":"511","vat_account_id":"3331"}]
	}`, ts.cust.ID, time.Now().UTC().Format(time.RFC3339), time.Now().UTC().Add(30*24*time.Hour).Format(time.RFC3339))
	w := postJSON(ts, "/api/v1/sale/invoices?company_id=CMP001", body)
	assert.Equal(t, 201, w.Code)
}

func TestSaleCreateInvoice_DuplicateNumber(t *testing.T) {
	ts := setupSale(t)
	body := fmt.Sprintf(`{
		"company_id":"CMP001","invoice_number":"INV-DUP","customer_id":"%s","customer_name":"Test","customer_tax_code":"1234567890",
		"invoice_date":"%s","currency":"VND",
		"lines":[{"item_name":"Inv","unit":"pc","quantity":1,"unit_price":100,"vat_rate":10,"vat_type":"VAT10","revenue_account_id":"511","vat_account_id":"3331"}]
	}`, ts.cust.ID, time.Now().UTC().Format(time.RFC3339))
	postJSON(ts, "/api/v1/sale/invoices?company_id=CMP001", body)
	w := postJSON(ts, "/api/v1/sale/invoices?company_id=CMP001", body)
	assert.Equal(t, 400, w.Code)
}

// ─── Receipt ───────────────────────────────────────────────────────────────

func TestSaleCreateReceipt(t *testing.T) {
	ts := setupSale(t)
	body := fmt.Sprintf(`{
		"company_id":"CMP001","receipt_number":"RCP-TEST-001","customer_id":"%s",
		"receipt_date":"%s","amount":500000,"currency":"VND"
	}`, ts.cust.ID, time.Now().UTC().Format(time.RFC3339))
	w := postJSON(ts, "/api/v1/sale/receipts?company_id=CMP001", body)
	assert.Equal(t, 201, w.Code)
}

func TestSaleCreateReceipt_DuplicateNumber(t *testing.T) {
	ts := setupSale(t)
	body := fmt.Sprintf(`{
		"company_id":"CMP001","receipt_number":"RCP-DUP","customer_id":"%s",
		"receipt_date":"%s","amount":500000,"currency":"VND"
	}`, ts.cust.ID, time.Now().UTC().Format(time.RFC3339))
	postJSON(ts, "/api/v1/sale/receipts?company_id=CMP001", body)
	w := postJSON(ts, "/api/v1/sale/receipts?company_id=CMP001", body)
	assert.Equal(t, 400, w.Code)
}

// ─── Credit Note ───────────────────────────────────────────────────────────

func TestSaleCreateCN(t *testing.T) {
	ts := setupSale(t)
	invBody := fmt.Sprintf(`{
		"company_id":"CMP001","invoice_number":"INV-FOR-CN","customer_id":"%s","customer_name":"Test","customer_tax_code":"1234567890",
		"invoice_date":"%s","currency":"VND",
		"lines":[{"item_name":"CN Item","unit":"pc","quantity":1,"unit_price":100000,"vat_rate":10,"vat_type":"VAT10","revenue_account_id":"511","vat_account_id":"3331"}]
	}`, ts.cust.ID, time.Now().UTC().Format(time.RFC3339))
	var inv domain.CustomerInvoice
	w := postJSON(ts, "/api/v1/sale/invoices?company_id=CMP001", invBody)
	json.Unmarshal(w.Body.Bytes(), &inv)

	cnBody := fmt.Sprintf(`{
		"company_id":"CMP001","cn_number":"CN-TEST-001","original_invoice_id":"%s","customer_id":"%s",
		"return_date":"%s","subtotal":100000,"tax_amount":10000,"total_amount":110000,
		"lines":[{"item_name":"Returned","unit":"pc","quantity":1,"unit_price":100000,"vat_rate":10,"line_total":100000,"line_vat_amount":10000}]
	}`, inv.ID, ts.cust.ID, time.Now().UTC().Format(time.RFC3339))
	w = postJSON(ts, "/api/v1/sale/credit-notes?company_id=CMP001", cnBody)
	assert.Equal(t, 201, w.Code)
}

func TestSaleCreateCN_DuplicateNumber(t *testing.T) {
	ts := setupSale(t)
	invBody := fmt.Sprintf(`{
		"company_id":"CMP001","invoice_number":"INV-FOR-CN2","customer_id":"%s","customer_name":"Test","customer_tax_code":"1234567890",
		"invoice_date":"%s","currency":"VND",
		"lines":[{"item_name":"CN2","unit":"pc","quantity":1,"unit_price":100,"vat_rate":10,"vat_type":"VAT10","revenue_account_id":"511","vat_account_id":"3331"}]
	}`, ts.cust.ID, time.Now().UTC().Format(time.RFC3339))
	var inv domain.CustomerInvoice
	w := postJSON(ts, "/api/v1/sale/invoices?company_id=CMP001", invBody)
	json.Unmarshal(w.Body.Bytes(), &inv)

	cnBody := fmt.Sprintf(`{
		"company_id":"CMP001","cn_number":"CN-DUP","original_invoice_id":"%s","customer_id":"%s",
		"return_date":"%s","subtotal":100,"total_amount":110,
		"lines":[{"item_name":"Dup","unit":"pc","quantity":1,"unit_price":100,"vat_rate":10,"line_total":100,"line_vat_amount":10}]
	}`, inv.ID, ts.cust.ID, time.Now().UTC().Format(time.RFC3339))
	postJSON(ts, "/api/v1/sale/credit-notes?company_id=CMP001", cnBody)
	w = postJSON(ts, "/api/v1/sale/credit-notes?company_id=CMP001", cnBody)
	assert.Equal(t, 400, w.Code)
}

// ─── Sales Quotation ───────────────────────────────────────────────────────

func TestSaleCreateSQ(t *testing.T) {
	ts := setupSale(t)
	body := `{"company_id":"CMP001","qn_number":"QN-TEST-001"}`
	w := postJSON(ts, "/api/v1/sale/quotations?company_id=CMP001", body)
	assert.Equal(t, 201, w.Code)
}

func TestSaleCreateSQ_EmptyNumber(t *testing.T) {
	ts := setupSale(t)
	body := `{"company_id":"CMP001","qn_number":""}`
	w := postJSON(ts, "/api/v1/sale/quotations?company_id=CMP001", body)
	assert.Equal(t, 400, w.Code)
}

func TestSaleListSQs(t *testing.T) {
	ts := setupSale(t)
	body := `{"company_id":"CMP001","qn_number":"QN-LIST"}`
	postJSON(ts, "/api/v1/sale/quotations?company_id=CMP001", body)
	w := getJSON(ts, "/api/v1/sale/quotations?company_id=CMP001")
	assert.Equal(t, 200, w.Code)
}

// ─── AR Reports ────────────────────────────────────────────────────────────

func TestSaleGetARAging(t *testing.T) {
	ts := setupSale(t)
	w := getJSON(ts, "/api/v1/sale/ar/aging?company_id=CMP001")
	assert.Equal(t, 200, w.Code)
}

func TestSaleGetARSummary(t *testing.T) {
	ts := setupSale(t)
	w := getJSON(ts, "/api/v1/sale/ar/summary?company_id=CMP001")
	assert.Equal(t, 200, w.Code)
}

// ─── Cross-tenant ──────────────────────────────────────────────────────────

func TestSaleCrossTenant_GetDN(t *testing.T) {
	ts := setupSale(t)
	so := setupSaleWithSO(t, ts)
	var dn domain.DeliveryNote
	w := postJSON(ts, "/api/v1/sale/deliveries?company_id=CMP001", dnBody(so.ID, "DN-XT", time.Now().UTC().Format(time.RFC3339)))
	json.Unmarshal(w.Body.Bytes(), &dn)
	require.NotEmpty(t, dn.ID)

	w = getJSON(ts, "/api/v1/sale/deliveries/"+dn.ID+"?company_id=OTHER")
	assert.Equal(t, 404, w.Code)
}

func TestSaleCrossTenant_GetInvoice(t *testing.T) {
	ts := setupSale(t)
	body := fmt.Sprintf(`{"company_id":"CMP001","invoice_number":"INV-XT","customer_id":"%s","customer_name":"T","customer_tax_code":"123","invoice_date":"%s","currency":"VND","lines":[{"item_name":"X","unit":"pc","quantity":1,"unit_price":10,"vat_rate":10,"vat_type":"VAT10","revenue_account_id":"511","vat_account_id":"3331"}]}`, ts.cust.ID, time.Now().UTC().Format(time.RFC3339))
	var inv domain.CustomerInvoice
	w := postJSON(ts, "/api/v1/sale/invoices?company_id=CMP001", body)
	json.Unmarshal(w.Body.Bytes(), &inv)

	w = getJSON(ts, "/api/v1/sale/invoices/"+inv.ID+"?company_id=OTHER")
	assert.Equal(t, 404, w.Code)
}

func TestSaleCrossTenant_GetReceipt(t *testing.T) {
	ts := setupSale(t)
	body := fmt.Sprintf(`{"company_id":"CMP001","receipt_number":"RCP-XT","customer_id":"%s","receipt_date":"%s","amount":100000,"currency":"VND"}`, ts.cust.ID, time.Now().UTC().Format(time.RFC3339))
	var rcpt domain.CustomerReceipt
	w := postJSON(ts, "/api/v1/sale/receipts?company_id=CMP001", body)
	json.Unmarshal(w.Body.Bytes(), &rcpt)

	w = getJSON(ts, "/api/v1/sale/receipts/"+rcpt.ID+"?company_id=OTHER")
	assert.Equal(t, 404, w.Code)
}

func TestSaleCrossTenant_GetCN(t *testing.T) {
	ts := setupSale(t)
	invBody := fmt.Sprintf(`{"company_id":"CMP001","invoice_number":"INV-CN-XT","customer_id":"%s","customer_name":"T","customer_tax_code":"123","invoice_date":"%s","currency":"VND","lines":[{"item_name":"X","unit":"pc","quantity":1,"unit_price":10,"vat_rate":10,"vat_type":"VAT10","revenue_account_id":"511","vat_account_id":"3331"}]}`, ts.cust.ID, time.Now().UTC().Format(time.RFC3339))
	var inv domain.CustomerInvoice
	w := postJSON(ts, "/api/v1/sale/invoices?company_id=CMP001", invBody)
	json.Unmarshal(w.Body.Bytes(), &inv)
	cnBody := fmt.Sprintf(`{"company_id":"CMP001","cn_number":"CN-XT","original_invoice_id":"%s","customer_id":"%s","return_date":"%s","subtotal":10,"total_amount":11,"lines":[{"item_name":"R","unit":"pc","quantity":1,"unit_price":10,"vat_rate":10,"line_total":10,"line_vat_amount":1}]}`, inv.ID, ts.cust.ID, time.Now().UTC().Format(time.RFC3339))
	var cn domain.CreditNote
	w = postJSON(ts, "/api/v1/sale/credit-notes?company_id=CMP001", cnBody)
	json.Unmarshal(w.Body.Bytes(), &cn)

	w = getJSON(ts, "/api/v1/sale/credit-notes/"+cn.ID+"?company_id=OTHER")
	assert.Equal(t, 404, w.Code)
}

func TestSaleCrossTenant_GetSQ(t *testing.T) {
	ts := setupSale(t)
	body := `{"company_id":"CMP001","qn_number":"QN-XT"}`
	var sq domain.SalesQuotation
	w := postJSON(ts, "/api/v1/sale/quotations?company_id=CMP001", body)
	json.Unmarshal(w.Body.Bytes(), &sq)

	w = getJSON(ts, "/api/v1/sale/quotations/"+sq.ID+"?company_id=OTHER")
	assert.Equal(t, 404, w.Code)
}

func TestGetCustomerStatement(t *testing.T) {
	ts := setupSale(t)

	// post an invoice to have AR activity
	invBody := fmt.Sprintf(`{"company_id":"CMP001","invoice_number":"STMT-INV","customer_id":"%s","customer_name":"Test Customer","customer_tax_code":"1234567890","customer_address":"addr","invoice_type":"domestic","currency":"VND","invoice_date":"%s","lines":[{"item_name":"svc","quantity":1,"unit_price":200,"vat_rate":10,"vat_type":"VAT_10","revenue_account_id":"5111"}]}`,
		ts.cust.ID, time.Now().Add(-48*time.Hour).Format(time.RFC3339))
	w := postJSON(ts, "/api/v1/sale/invoices?company_id=CMP001", invBody)
	assert.Equal(t, 201, w.Code)
	var inv domain.CustomerInvoice
	json.Unmarshal(w.Body.Bytes(), &inv)

	// post the invoice
	w = patchJSON(ts, "/api/v1/sale/invoices/"+inv.ID+"/post?company_id=CMP001", "")
	assert.Equal(t, 200, w.Code)

	// get statement
	w = getJSON(ts, "/api/v1/sale/ar/statement?customer_id="+ts.cust.ID+"&company_id=CMP001")
	assert.Equal(t, 200, w.Code)
	var stmt domain.CustomerStatement
	err := json.Unmarshal(w.Body.Bytes(), &stmt)
	require.NoError(t, err)
	assert.Equal(t, ts.cust.ID, stmt.Customer.ID)
	assert.Equal(t, 220.0, stmt.ClosingBal) // 200 + 20 VAT = 220
}

func TestReportRoutes(t *testing.T) {
	ts := setupSale(t)

	// invoice + post
	invBody := fmt.Sprintf(`{"company_id":"CMP001","invoice_number":"RPT-INV","customer_id":"%s","customer_name":"Test Customer","customer_tax_code":"1234567890","customer_address":"addr","invoice_type":"domestic","currency":"VND","invoice_date":"%s","lines":[{"item_name":"Widget","unit":"pcs","quantity":2,"unit_price":100,"vat_rate":10,"vat_type":"VAT_10","revenue_account_id":"5111"}]}`,
		ts.cust.ID, time.Now().Add(-48*time.Hour).Format(time.RFC3339))
	w := postJSON(ts, "/api/v1/sale/invoices?company_id=CMP001", invBody)
	require.Equal(t, 201, w.Code)
	var inv domain.CustomerInvoice
	json.Unmarshal(w.Body.Bytes(), &inv)
	w = patchJSON(ts, "/api/v1/sale/invoices/"+inv.ID+"/post?company_id=CMP001", "")
	require.Equal(t, 200, w.Code)

	t.Run("s01-bh", func(t *testing.T) {
		w = getJSON(ts, "/api/v1/sale/reports/s01-bh?company_id=CMP001&customer_id="+ts.cust.ID)
		assert.Equal(t, 200, w.Code)
		var rpt domain.SalesLedgerReport
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &rpt))
		require.Len(t, rpt.Rows, 1)
		assert.Equal(t, 220.0, rpt.Total)
	})

	t.Run("s02-bh", func(t *testing.T) {
		w = getJSON(ts, "/api/v1/sale/reports/s02-bh?company_id=CMP001&customer_id="+ts.cust.ID)
		assert.Equal(t, 200, w.Code)
		var rpt domain.CustomerLedgerReport
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &rpt))
		assert.Equal(t, 220.0, rpt.ClosingBalance)
	})

	t.Run("s03-bh", func(t *testing.T) {
		w = getJSON(ts, "/api/v1/sale/reports/s03-bh?company_id=CMP001")
		assert.Equal(t, 200, w.Code)
		var rpt domain.GoodsLedgerReport
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &rpt))
		require.Len(t, rpt.Rows, 1)
		assert.Equal(t, 200.0, rpt.TotalRevenue)
	})

	t.Run("vat-output", func(t *testing.T) {
		w = getJSON(ts, "/api/v1/sale/reports/vat-output?company_id=CMP001")
		assert.Equal(t, 200, w.Code)
		var rpt domain.VATOutputReport
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &rpt))
		require.Len(t, rpt.Rows, 1)
		assert.Equal(t, 10.0, rpt.Rows[0].VatRate)
		assert.Equal(t, 20.0, rpt.Rows[0].VatAmount)
	})

	t.Run("unbilled-deliveries", func(t *testing.T) {
		w = getJSON(ts, "/api/v1/sale/reports/unbilled-deliveries?company_id=CMP001")
		assert.Equal(t, 200, w.Code)
		var rpt domain.UnbilledDeliveryReport
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &rpt))
		assert.Empty(t, rpt.Rows)
	})

	t.Run("ar-recon", func(t *testing.T) {
		w = getJSON(ts, "/api/v1/sale/ar/recon?company_id=CMP001")
		assert.Equal(t, 200, w.Code)
		var rpt domain.ARGLReconciliation
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &rpt))
		assert.Equal(t, 220.0, rpt.SubledgerTotal)
	})
}

func TestNextNumberEndpoints(t *testing.T) {
	ts := setupSale(t)

	type nextNum struct {
		NextNumber string `json:"next_number"`
	}

	t.Run("orders", func(t *testing.T) {
		w := getJSON(ts, "/api/v1/sale/orders/next-number?company_id=CMP001&yyyymm=202607")
		assert.Equal(t, 200, w.Code)
		var r nextNum
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &r))
		assert.Equal(t, "SO-202607-00001", r.NextNumber)
	})

	t.Run("deliveries", func(t *testing.T) {
		w := getJSON(ts, "/api/v1/sale/deliveries/next-number?company_id=CMP001&yyyymm=202607")
		assert.Equal(t, 200, w.Code)
		var r nextNum
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &r))
		assert.Equal(t, "DN-202607-00001", r.NextNumber)
	})

	t.Run("invoices", func(t *testing.T) {
		w := getJSON(ts, "/api/v1/sale/invoices/next-number?company_id=CMP001&yyyymm=202607")
		assert.Equal(t, 200, w.Code)
		var r nextNum
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &r))
		assert.Equal(t, "INV-202607-00001", r.NextNumber)
	})

	t.Run("default yyyymm is current month", func(t *testing.T) {
		w := getJSON(ts, "/api/v1/sale/orders/next-number?company_id=CMP001")
		assert.Equal(t, 200, w.Code)
		var r nextNum
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &r))
		assert.Regexp(t, `^SO-\d{6}-\d{5}$`, r.NextNumber)
	})
}
