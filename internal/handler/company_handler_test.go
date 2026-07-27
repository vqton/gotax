package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gotax/internal/domain"
	"gotax/internal/repository"
	"gotax/internal/service"
)

type companyTestSetup struct {
	r    *gin.Engine
	svc  service.CompanyService
	comp *domain.Company
}

func setupCompanyWithTestData(t *testing.T) *companyTestSetup {
	t.Helper()
	gin.SetMode(gin.TestMode)

	repo := repository.NewMemoryCompanyRepo()
	svc := service.NewCompanyService(repo)
	ch := NewCompanyHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "tenant-1")
		c.Set("username", "tenantuser")
		c.Set("role", "admin")
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterCompanyRoutes(r, ch, noopMW, noopMW)

	comp := &domain.Company{
		TenantID:        "tenant-1",
		TaxCode:         "1234567890",
		LegalNameVN:     "Test Co",
		LegalForm:       domain.LegalFormJSC,
		AccountingRegime: domain.RegimeTT99,
	}
	svc.CreateCompany(nil, comp)

	return &companyTestSetup{r: r, svc: svc, comp: comp}
}

// ─── Company ────────────────────────────────────────────────────────

func TestCreateCompany(t *testing.T) {
	ts := setupCompanyWithTestData(t)
	body := `{"tax_code":"0987654321","legal_name_vn":"New Co","legal_form":"JSC","accounting_regime":"TT99"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/companies", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	var result domain.Company
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.NotEmpty(t, result.ID)
}

func TestCreateCompany_BadRequest(t *testing.T) {
	ts := setupCompanyWithTestData(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/companies", strings.NewReader(`not-json`))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestListCompanies(t *testing.T) {
	ts := setupCompanyWithTestData(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/companies", nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestGetCompany(t *testing.T) {
	ts := setupCompanyWithTestData(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/companies/"+ts.comp.ID, nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	var result domain.Company
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, "Test Co", result.LegalNameVN)
}

func TestGetCompany_NotFound(t *testing.T) {
	ts := setupCompanyWithTestData(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/companies/nonexistent", nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

func TestGetCompanyByTaxCode(t *testing.T) {
	ts := setupCompanyWithTestData(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/companies/tax-code/1234567890", nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestUpdateCompany(t *testing.T) {
	ts := setupCompanyWithTestData(t)
	body := `{"legal_name_vn":"Updated Name","tax_code":"1234567890","legal_form":"JSC","accounting_regime":"TT99"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/companies/"+ts.comp.ID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestDeactivateCompany(t *testing.T) {
	ts := setupCompanyWithTestData(t)
	body := `{"reason":"closed"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/companies/"+ts.comp.ID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestGetCompanyHierarchy(t *testing.T) {
	ts := setupCompanyWithTestData(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/companies/"+ts.comp.ID+"/hierarchy", nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestSwitchCompany(t *testing.T) {
	ts := setupCompanyWithTestData(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/companies/"+ts.comp.ID+"/switch", nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

// ─── Branch ─────────────────────────────────────────────────────────

func TestCreateBranch(t *testing.T) {
	ts := setupCompanyWithTestData(t)
	body := `{"branch_tax_code":"BR-001","branch_name":"Branch 1"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/companies/"+ts.comp.ID+"/branches", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
}

func TestListBranches(t *testing.T) {
	ts := setupCompanyWithTestData(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/companies/"+ts.comp.ID+"/branches", nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

// ─── Fiscal Year ────────────────────────────────────────────────────

func TestCreateFiscalYear(t *testing.T) {
	ts := setupCompanyWithTestData(t)
	body := `{"year":2026,"start_month":1}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/companies/"+ts.comp.ID+"/fiscal-years", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
}

func TestListFiscalYears(t *testing.T) {
	ts := setupCompanyWithTestData(t)
	ts.svc.CreateFiscalYear(nil, &domain.FiscalYear{CompanyID: ts.comp.ID, Year: 2026, StartMonth: 1})
	ts.svc.CreateFiscalYear(nil, &domain.FiscalYear{CompanyID: ts.comp.ID, Year: 2027, StartMonth: 1})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/companies/"+ts.comp.ID+"/fiscal-years", nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

// ─── Period V2 ──────────────────────────────────────────────────────

func TestGeneratePeriods(t *testing.T) {
	ts := setupCompanyWithTestData(t)
	fy := &domain.FiscalYear{CompanyID: ts.comp.ID, Year: time.Now().Year(), StartMonth: 1}
	ts.svc.CreateFiscalYear(nil, fy)

	body := `{"year":` + itoa(time.Now().Year()) + `,"start_month":1}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/companies/"+ts.comp.ID+"/periods/generate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
}

func TestClosePeriod(t *testing.T) {
	ts := setupCompanyWithTestData(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/companies/"+ts.comp.ID+"/periods/period-1/close", nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestGetCurrentPeriod_NotFound(t *testing.T) {
	ts := setupCompanyWithTestData(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/companies/"+ts.comp.ID+"/periods/current", nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

// ─── Department ─────────────────────────────────────────────────────

func TestCreateDepartment(t *testing.T) {
	ts := setupCompanyWithTestData(t)
	body := `{"code":"D001","name":"Engineering"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/companies/"+ts.comp.ID+"/departments", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
}

func TestListDepartments(t *testing.T) {
	ts := setupCompanyWithTestData(t)
	ts.svc.CreateDepartment(nil, &domain.Department{CompanyID: ts.comp.ID, Code: "D001", Name: "Eng"})
	ts.svc.CreateDepartment(nil, &domain.Department{CompanyID: ts.comp.ID, Code: "D002", Name: "HR"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/companies/"+ts.comp.ID+"/departments", nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

// ─── Employee ───────────────────────────────────────────────────────

func TestCreateEmployee(t *testing.T) {
	ts := setupCompanyWithTestData(t)
	body := `{"employee_code":"E001","full_name":"Alice"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/companies/"+ts.comp.ID+"/employees", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
}

func TestListEmployees(t *testing.T) {
	ts := setupCompanyWithTestData(t)
	ts.svc.CreateEmployee(nil, &domain.Employee{CompanyID: ts.comp.ID, EmployeeCode: "E001", FullName: "Alice"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/companies/"+ts.comp.ID+"/employees", nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

// ─── Bank Account ───────────────────────────────────────────────────

func TestCreateBankAccount(t *testing.T) {
	ts := setupCompanyWithTestData(t)
	body := `{"account_number":"123456789","bank_name":"Vietcombank","account_holder":"Test Co"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/companies/"+ts.comp.ID+"/bank-accounts", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
}

func TestListBankAccounts(t *testing.T) {
	ts := setupCompanyWithTestData(t)
	ts.svc.CreateBankAccount(nil, &domain.CompanyBankAccount{CompanyID: ts.comp.ID, AccountNumber: "123", BankName: "VCB"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/companies/"+ts.comp.ID+"/bank-accounts", nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

// ─── E-Invoice ──────────────────────────────────────────────────────

func TestRegisterEInvoicePattern(t *testing.T) {
	ts := setupCompanyWithTestData(t)
	body := `{"pattern_code":"01GTKT0/001","serial":"AA/25E"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/companies/"+ts.comp.ID+"/einvoice-patterns", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
}

func TestListEInvoicePatterns(t *testing.T) {
	ts := setupCompanyWithTestData(t)
	ts.svc.RegisterEInvoicePattern(nil, &domain.EInvoicePattern{
		CompanyID: ts.comp.ID, PatternCode: "01GTKT0/001", Serial: "AA/25E",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/companies/"+ts.comp.ID+"/einvoice-patterns", nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

// ─── Digital Signature ──────────────────────────────────────────────

func TestRegisterDigitalSignature(t *testing.T) {
	ts := setupCompanyWithTestData(t)
	body := `{"owner_name":"Alice","serial_number":"SN-001","provider":"VIETTEL","signature_type":"USB_TOKEN","valid_from":"2026-01-01","valid_to":"2027-01-01"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/companies/"+ts.comp.ID+"/digital-signatures", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
}

func TestListDigitalSignatures(t *testing.T) {
	ts := setupCompanyWithTestData(t)
	ts.svc.RegisterDigitalSignature(nil, &domain.DigitalSignature{
		CompanyID: ts.comp.ID, OwnerName: "Alice", SerialNumber: "SN-001", Provider: "VIETTEL",
		SignatureType: domain.SignatureUSBToken,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/companies/"+ts.comp.ID+"/digital-signatures", nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

// ─── Integration ────────────────────────────────────────────────────

func TestCreateIntegrationProfile(t *testing.T) {
	ts := setupCompanyWithTestData(t)
	body := `{"integration_type":"GDT","endpoint_url":"https://api.example.com"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/companies/"+ts.comp.ID+"/integrations", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
}

func TestListIntegrationProfiles(t *testing.T) {
	ts := setupCompanyWithTestData(t)
	ts.svc.CreateIntegrationProfile(nil, &domain.IntegrationProfile{
		CompanyID: ts.comp.ID, IntegrationType: domain.IntegrationGDT,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/companies/"+ts.comp.ID+"/integrations", nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}
