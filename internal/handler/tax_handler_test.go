package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/domain"
	"gotax/internal/einvoice"
	"gotax/internal/repository"
	"gotax/internal/service"
)

type taxTestSetup struct {
	r       *gin.Engine
	svc     service.TaxServiceInterface
	taxRepo domain.TaxRepository
	jeRepo  domain.JournalRepository
	gdt     *stubGDT
	compID  string
}

type stubGDT struct {
	submitResp *domain.GDTSubmitResponse
	submitErr  error
	submitted  bool
}

func (g *stubGDT) SubmitInvoice(_ context.Context, xml, certID string) (*domain.GDTSubmitResponse, error) {
	g.submitted = true
	if g.submitErr != nil {
		return nil, g.submitErr
	}
	return g.submitResp, nil
}
func (g *stubGDT) GetInvoiceStatus(_ context.Context, _ string) (*domain.GDTStatusResponse, error) {
	return &domain.GDTStatusResponse{Status: "ACKNOWLEDGED"}, nil
}
func (g *stubGDT) CancelInvoice(_ context.Context, _, _ string) error { return nil }
func (g *stubGDT) SubmitDeclaration(_ context.Context, xml, certID string) (*domain.GDTDeclarationSubmitResponse, error) {
	g.submitted = true
	if g.submitErr != nil {
		return nil, g.submitErr
	}
	return &domain.GDTDeclarationSubmitResponse{SubmissionID: "SUB-1", Status: "SUBMITTED"}, nil
}
func (g *stubGDT) QueryDeclarationStatus(_ context.Context, _ string) (*domain.GDTDeclarationStatusResponse, error) {
	return &domain.GDTDeclarationStatusResponse{Status: "ACKNOWLEDGED", AckRef: "ACK-REF-1"}, nil
}
func (g *stubGDT) LookupTaxCode(_ context.Context, taxCode string) (*domain.TaxCodeLookupResponse, error) {
	if taxCode == "NOTFOUND" {
		return &domain.TaxCodeLookupResponse{TaxCode: taxCode, Status: "NOT_FOUND", Message: "tax code not found"}, nil
	}
	if taxCode == "INACTIVE" {
		return &domain.TaxCodeLookupResponse{TaxCode: taxCode, Name: "Inactive Corp", Status: "INACTIVE", Message: "tax code inactive"}, nil
	}
	return &domain.TaxCodeLookupResponse{TaxCode: taxCode, Name: "Active Corp", Status: "ACTIVE"}, nil
}

type stubSigner struct{}

func (s *stubSigner) SignTXML(xmlBody, _ string) (string, error) { return "signed:" + xmlBody, nil }
func (s *stubSigner) SignDocument(xmlBody string) (einvoice.SignResult, error) {
	return einvoice.SignResult{SignatureBase64: "BASE64:" + xmlBody, SignedAt: "2026-04-15T10:00:00+07:00"}, nil
}

func newTaxTestSvcIssuerH() (*stubGDT, service.TXMLSigner) {
	return &stubGDT{submitResp: &domain.GDTSubmitResponse{TransactionID: "TXN-1", Status: "SUBMITTED"}}, &stubSigner{}
}

func setupTaxTest(t *testing.T) *taxTestSetup {
	t.Helper()
	gin.SetMode(gin.TestMode)

	taxRepo := repository.NewMemoryTaxRepo()
	jeRepo := repository.NewMemoryJournalRepo()
	gdtStub, signerStub := newTaxTestSvcIssuerH()
	companyRepo := repository.NewMemoryCompanyRepo()
	_ = companyRepo.Create(nil, &domain.Company{
		ID: "company-1", TaxCode: "0100123456", LegalNameVN: "CONG TY TEST",
	})
	taxSvc := service.NewTaxService(taxRepo, jeRepo, companyRepo, gdtStub, signerStub)
	auditRepo := repository.NewMemoryAuditLogRepo()
	taxDeclSvc := service.NewTaxDeclarationService(taxRepo)
	th := NewTaxHandler(taxSvc, taxDeclSvc, auditRepo)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user")
		c.Set("username", "testuser")
		c.Set("role", "admin")
	})
	noopMW := func(c *gin.Context) { c.Next() }
	RegisterTaxRoutes(r, th, noopMW)

	return &taxTestSetup{
		r:       r,
		svc:     taxSvc,
		taxRepo: taxRepo,
		jeRepo:  jeRepo,
		gdt:     gdtStub,
		compID:  "company-1",
	}
}

// ─── Declarations ──────────────────────────────────────────────────────

func TestCreateDeclaration(t *testing.T) {
	ts := setupTaxTest(t)
	body := `{
		"company_id":"` + ts.compID + `",
		"declaration_type":"GTGT01",
		"tax_period":{"period_type":"MONTHLY","period_year":2026,"period_number":1},
		"lines":[{"line_code":"01","line_name":"Test","amount":1000000,"source_type":"MANUAL_ENTRY","sort_order":1}]
	}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tax/declarations", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp domain.TaxDeclaration
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, domain.DeclTypeGTGT01, resp.DeclarationType)
	assert.Equal(t, domain.DeclStatusDRAFT, resp.Status)
	assert.NotEmpty(t, resp.ID)
}

func TestGenerateDeclaration(t *testing.T) {
	ts := setupTaxTest(t)
	je := &domain.JournalEntry{
		ID: "JE1", EntryNumber: "JE1", CompanyID: ts.compID,
		EntryDate: time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC),
		Status:    domain.JournalEntryPosted,
		Lines: []domain.JournalLine{
			{AccountCode: "5111", CreditAmount: 10000000},
			{AccountCode: "33311", CreditAmount: 1000000},
		},
	}
	require.NoError(t, ts.jeRepo.Create(context.Background(), je))

	body := `{
		"company_id":"` + ts.compID + `",
		"declaration_type":"GTGT01",
		"tax_period":{"period_type":"MONTHLY","period_year":2026,"period_number":1}
	}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tax/declarations/generate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp domain.TaxDeclaration
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, domain.DeclStatusVALIDATED, resp.Status)
	assert.Len(t, resp.Lines, 9)
	amounts := map[string]float64{}
	for _, l := range resp.Lines {
		amounts[l.LineCode] = l.Amount
	}
	assert.Equal(t, 1000000.0, amounts["23"])
	assert.Equal(t, 1000000.0, amounts["31"])
}

func TestGenerateDeclaration_Conflict(t *testing.T) {
	ts := setupTaxTest(t)
	body := `{
		"company_id":"` + ts.compID + `",
		"declaration_type":"GTGT01",
		"tax_period":{"period_type":"MONTHLY","period_year":2026,"period_number":1}
	}`
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/tax/declarations/generate", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		ts.r.ServeHTTP(w, req)
		if i == 0 {
			assert.Equal(t, http.StatusCreated, w.Code)
		} else {
			assert.Equal(t, http.StatusConflict, w.Code)
		}
	}
}

func TestCreateDeclaration_Invalid(t *testing.T) {
	ts := setupTaxTest(t)
	body := `{"company_id":"","declaration_type":"INVALID"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tax/declarations", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetDeclaration(t *testing.T) {
	ts := setupTaxTest(t)
	d := &domain.TaxDeclaration{
		CompanyID:       ts.compID,
		DeclarationType: domain.DeclTypeGTGT01,
		TaxPeriod:       domain.TaxPeriod{PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 1},
	}
	require.NoError(t, ts.svc.CreateDeclaration(nil, d))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tax/declarations/"+d.ID, nil)
	ts.r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp domain.TaxDeclaration
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, d.ID, resp.ID)
}

func TestGetDeclaration_NotFound(t *testing.T) {
	ts := setupTaxTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tax/declarations/nonexistent", nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestListDeclarations(t *testing.T) {
	ts := setupTaxTest(t)
	for i := 1; i <= 3; i++ {
		d := &domain.TaxDeclaration{
			CompanyID:       ts.compID,
			DeclarationType: domain.DeclTypeGTGT01,
			TaxPeriod:       domain.TaxPeriod{PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: i},
		}
		ts.svc.CreateDeclaration(nil, d)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tax/declarations?company_id="+ts.compID, nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp []domain.TaxDeclaration
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp, 3)
}

func TestSubmitDeclaration(t *testing.T) {
	ts := setupTaxTest(t)
	d := &domain.TaxDeclaration{
		CompanyID:       ts.compID,
		DeclarationType: domain.DeclTypeGTGT01,
		TaxPeriod:       domain.TaxPeriod{PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 1},
	}
	require.NoError(t, ts.svc.CreateDeclaration(nil, d))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tax/declarations/"+d.ID+"/submit", nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	updated, _ := ts.svc.GetDeclaration(nil, d.ID)
	assert.Equal(t, domain.DeclStatusSUBMITTED, updated.Status)
}

func TestSubmitDeclaration_NonDraft(t *testing.T) {
	ts := setupTaxTest(t)
	d := &domain.TaxDeclaration{
		CompanyID:       ts.compID,
		DeclarationType: domain.DeclTypeGTGT01,
		TaxPeriod:       domain.TaxPeriod{PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 1},
	}
	require.NoError(t, ts.svc.CreateDeclaration(nil, d))
	ts.svc.SubmitDeclaration(nil, d.ID, "user")

	// second submit should fail
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tax/declarations/"+d.ID+"/submit", nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestCheckDeclarationStatus(t *testing.T) {
	ts := setupTaxTest(t)
	d := &domain.TaxDeclaration{
		ID:              "DECL-1",
		CompanyID:       ts.compID,
		DeclarationType: domain.DeclTypeGTGT01,
		TaxPeriod:       domain.TaxPeriod{PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 1},
		Status:          domain.DeclStatusSUBMITTED,
		GDTSubmissionID: "SUB-1",
	}
	require.NoError(t, ts.svc.CreateDeclaration(nil, d))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tax/declarations/"+d.ID+"/check-status", nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	updated, _ := ts.svc.GetDeclaration(nil, d.ID)
	assert.Equal(t, domain.DeclStatusACKNOWLEDGED, updated.Status)
}

func TestAcknowledgeDeclaration(t *testing.T) {
	ts := setupTaxTest(t)
	d := &domain.TaxDeclaration{
		CompanyID:       ts.compID,
		DeclarationType: domain.DeclTypeGTGT01,
		TaxPeriod:       domain.TaxPeriod{PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 1},
	}
	require.NoError(t, ts.svc.CreateDeclaration(nil, d))
	ts.svc.SubmitDeclaration(nil, d.ID, "user")

	w := httptest.NewRecorder()
	body := `{"reference":"GDT-REF-001"}`
	req, _ := http.NewRequest("POST", "/api/v1/tax/declarations/"+d.ID+"/acknowledge", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	updated, _ := ts.svc.GetDeclaration(nil, d.ID)
	assert.Equal(t, domain.DeclStatusACKNOWLEDGED, updated.Status)
	assert.Equal(t, "GDT-REF-001", updated.AcknowledgementRef)
}

func TestRejectDeclaration(t *testing.T) {
	ts := setupTaxTest(t)
	d := &domain.TaxDeclaration{
		CompanyID:       ts.compID,
		DeclarationType: domain.DeclTypeGTGT01,
		TaxPeriod:       domain.TaxPeriod{PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 1},
	}
	require.NoError(t, ts.svc.CreateDeclaration(nil, d))
	ts.svc.SubmitDeclaration(nil, d.ID, "user")

	w := httptest.NewRecorder()
	body := `{"reason":"data mismatch"}`
	req, _ := http.NewRequest("POST", "/api/v1/tax/declarations/"+d.ID+"/reject", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	updated, _ := ts.svc.GetDeclaration(nil, d.ID)
	assert.Equal(t, domain.DeclStatusREJECTED, updated.Status)
}

func TestCancelDeclaration(t *testing.T) {
	ts := setupTaxTest(t)
	d := &domain.TaxDeclaration{
		CompanyID:       ts.compID,
		DeclarationType: domain.DeclTypeGTGT01,
		TaxPeriod:       domain.TaxPeriod{PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 1},
	}
	require.NoError(t, ts.svc.CreateDeclaration(nil, d))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tax/declarations/"+d.ID+"/cancel", nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	updated, _ := ts.svc.GetDeclaration(nil, d.ID)
	assert.Equal(t, domain.DeclStatusCANCELLED, updated.Status)
}

func TestAmendDeclaration(t *testing.T) {
	ts := setupTaxTest(t)
	d := &domain.TaxDeclaration{
		CompanyID:       ts.compID,
		DeclarationType: domain.DeclTypeGTGT01,
		TaxPeriod:       domain.TaxPeriod{PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 1},
	}
	require.NoError(t, ts.svc.CreateDeclaration(nil, d))
	ts.svc.SubmitDeclaration(nil, d.ID, "user")
	ts.svc.AcknowledgeDeclaration(nil, d.ID, "ref-1")

	w := httptest.NewRecorder()
	body := `{"lines":[{"line_code":"01","line_name":"Amended","amount":2000000,"source_type":"MANUAL_ENTRY","sort_order":1}]}`
	req, _ := http.NewRequest("POST", "/api/v1/tax/declarations/"+d.ID+"/amend", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var amended domain.TaxDeclaration
	json.Unmarshal(w.Body.Bytes(), &amended)
	assert.Equal(t, domain.AdjTypeAMENDMENT, amended.AdjustmentType)
	assert.Equal(t, d.ID, amended.PreviousDeclID)
	assert.Equal(t, domain.DeclStatusDRAFT, amended.Status)
}

// ─── Tax Rates ─────────────────────────────────────────────────────────

func TestCreateRate(t *testing.T) {
	ts := setupTaxTest(t)
	body := `{"tax_type":"VAT","rate_code":"VAT10","rate_name":"VAT 10%","rate_type":"PERCENTAGE","rate_value":10,"effective_from":"2026-01-01","is_active":true}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tax/rates", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var resp domain.TaxRate
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "VAT10", resp.RateCode)
}

func TestListRates(t *testing.T) {
	ts := setupTaxTest(t)
	ts.svc.CreateRate(nil, &domain.TaxRate{TaxType: domain.TaxTypeVAT, RateCode: "VAT8", RateName: "VAT 8%", RateType: domain.RateTypePERCENTAGE, RateValue: 8, EffectiveFrom: "2026-01-01", IsActive: true})
	ts.svc.CreateRate(nil, &domain.TaxRate{TaxType: domain.TaxTypeVAT, RateCode: "VAT10", RateName: "VAT 10%", RateType: domain.RateTypePERCENTAGE, RateValue: 10, EffectiveFrom: "2026-01-01", IsActive: true})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tax/rates?tax_type=VAT", nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var rates []domain.TaxRate
	json.Unmarshal(w.Body.Bytes(), &rates)
	assert.Len(t, rates, 2)
}

// ─── Payments ──────────────────────────────────────────────────────────

func TestCreatePayment(t *testing.T) {
	ts := setupTaxTest(t)
	body := `{"company_id":"` + ts.compID + `","tax_type":"VAT","period_year":2026,"period_number":1,"declared_amount":5000000,"paid_amount":0,"due_date":"2026-04-20"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tax/payments", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var p domain.TaxPayment
	json.Unmarshal(w.Body.Bytes(), &p)
	assert.Equal(t, domain.PayStatusPENDING, p.Status)
}

func TestRecordPayment(t *testing.T) {
	ts := setupTaxTest(t)
	p := &domain.TaxPayment{
		CompanyID:      ts.compID,
		TaxType:        domain.TaxTypeVAT,
		PeriodYear:     2026,
		PeriodNumber:   1,
		DeclaredAmount: 5000000,
		DueDate:        "2026-04-20",
	}
	require.NoError(t, ts.svc.CreatePayment(nil, p))

	w := httptest.NewRecorder()
	body := `{"amount":5000000,"date":"2026-04-15","reference":"PAY-REF-001"}`
	req, _ := http.NewRequest("POST", "/api/v1/tax/payments/"+p.ID+"/record", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	updated, _ := ts.svc.GetPayment(nil, p.ID)
	assert.Equal(t, domain.PayStatusPAID, updated.Status)
	assert.Equal(t, 5000000.0, updated.PaidAmount)
	assert.NotEmpty(t, updated.GLJournalID, "journal entry should be created")

	je, err := ts.jeRepo.GetByID(nil, updated.GLJournalID)
	require.NoError(t, err)
	assert.Equal(t, "33311", je.Lines[0].AccountCode, "Dr VAT payable")
	assert.Equal(t, "112", je.Lines[1].AccountCode, "Cr bank")
}

// ─── E-Invoices ────────────────────────────────────────────────────────

func TestCreateEInvoice(t *testing.T) {
	ts := setupTaxTest(t)
	body := `{
		"company_id":"` + ts.compID + `",
		"pattern":"01GTKT0/001",
		"serial":"AA/25E",
		"invoice_type":"ORIGINAL",
		"buyer_name":"Test Buyer",
		"buyer_tax_code":"9876543210",
		"currency_code":"VND",
		"issue_date":"2026-03-15",
		"lines":[
			{"description":"Service","unit":"pc","quantity":1,"unit_price":1000000,"line_total":1000000,"vat_rate":10}
		]
	}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tax/e-invoices", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var inv domain.EInvoice
	json.Unmarshal(w.Body.Bytes(), &inv)
	assert.Equal(t, domain.EInvStatusDRAFT, inv.Status)
	assert.Equal(t, 1000000.0, inv.Subtotal)
	assert.Equal(t, 100000.0, inv.VATAmount)
	assert.Equal(t, 1100000.0, inv.GrandTotal)
}

func TestIssueEInvoice(t *testing.T) {
	ts := setupTaxTest(t)
	inv := &domain.EInvoice{
		CompanyID:   ts.compID,
		Pattern:     "01GTKT0/001",
		Serial:      "AA/25E",
		InvoiceType: domain.EInvTypeORIGINAL,
		BuyerName:   "Test Buyer",
		CurrencyCode: "VND",
		IssueDate:   "2026-03-15",
		Status:      domain.EInvStatusDRAFT,
		Subtotal:    1000000,
		VATAmount:   100000,
		GrandTotal:  1100000,
		Lines:       []domain.EInvoiceLine{{LineNumber: 1, Description: "Service", Quantity: 1, UnitPrice: 1000000, LineTotal: 1000000, VATRate: 10, VATAmount: 100000}},
	}
	require.NoError(t, ts.svc.CreateEInvoice(nil, inv))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tax/e-invoices/"+inv.ID+"/issue", nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	updated, _ := ts.svc.GetEInvoice(nil, inv.ID)
	assert.Equal(t, domain.EInvStatusSUBMITTED, updated.Status)
	assert.Equal(t, "TXN-1", updated.GDTTransactionID)
	assert.True(t, ts.gdt.submitted)
}

func TestCheckEInvoiceStatus(t *testing.T) {
	ts := setupTaxTest(t)
	inv := &domain.EInvoice{
		CompanyID:   ts.compID,
		Pattern:     "01GTKT0/001",
		Serial:      "AA/25E",
		InvoiceType: domain.EInvTypeORIGINAL,
		BuyerName:   "Test Buyer",
		CurrencyCode: "VND",
		IssueDate:   "2026-03-15",
		Status:      domain.EInvStatusSUBMITTED,
		GDTTransactionID: "TXN-1",
		Subtotal:    1000000,
		VATAmount:   100000,
		GrandTotal:  1100000,
		Lines:       []domain.EInvoiceLine{{LineNumber: 1, Description: "Service", Quantity: 1, UnitPrice: 1000000, LineTotal: 1000000, VATRate: 10, VATAmount: 100000}},
	}
	require.NoError(t, ts.svc.CreateEInvoice(nil, inv))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tax/e-invoices/"+inv.ID+"/status", nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	updated, _ := ts.svc.GetEInvoice(nil, inv.ID)
	assert.Equal(t, domain.EInvStatusISSUED, updated.Status)
}

func TestIssueEInvoice_GDTDown(t *testing.T) {
	ts := setupTaxTest(t)
	ts.gdt.submitErr = domain.ErrGDTUnavailable
	inv := &domain.EInvoice{
		CompanyID:   ts.compID,
		Pattern:     "01GTKT0/001",
		Serial:      "AA/25E",
		InvoiceType: domain.EInvTypeORIGINAL,
		BuyerName:   "Test Buyer",
		CurrencyCode: "VND",
		IssueDate:   "2026-03-15",
		Status:      domain.EInvStatusDRAFT,
		Subtotal:    1000000,
		VATAmount:   100000,
		GrandTotal:  1100000,
		Lines:       []domain.EInvoiceLine{{LineNumber: 1, Description: "Service", Quantity: 1, UnitPrice: 1000000, LineTotal: 1000000, VATRate: 10, VATAmount: 100000}},
	}
	require.NoError(t, ts.svc.CreateEInvoice(nil, inv))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tax/e-invoices/"+inv.ID+"/issue", nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadGateway, w.Code)

	updated, _ := ts.svc.GetEInvoice(nil, inv.ID)
	assert.Equal(t, domain.EInvStatusSIGNED, updated.Status)
	assert.Empty(t, updated.GDTTransactionID)
}

func TestCancelEInvoice(t *testing.T) {
	ts := setupTaxTest(t)
	inv := &domain.EInvoice{
		CompanyID:   ts.compID,
		Pattern:     "01GTKT0/001",
		Serial:      "AA/25E",
		InvoiceType: domain.EInvTypeORIGINAL,
		BuyerName:   "Test Buyer",
		CurrencyCode: "VND",
		IssueDate:   "2026-03-15",
		Status:      domain.EInvStatusISSUED,
		Subtotal:    1000000,
		VATAmount:   100000,
		GrandTotal:  1100000,
		Lines:       []domain.EInvoiceLine{{LineNumber: 1, Description: "Service", Quantity: 1, UnitPrice: 1000000, LineTotal: 1000000, VATRate: 10, VATAmount: 100000}},
	}
	require.NoError(t, ts.svc.CreateEInvoice(nil, inv))

	w := httptest.NewRecorder()
	body := `{"reason":"buyer returned goods"}`
	req, _ := http.NewRequest("POST", "/api/v1/tax/e-invoices/"+inv.ID+"/cancel", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	updated, _ := ts.svc.GetEInvoice(nil, inv.ID)
	assert.Equal(t, domain.EInvStatusCANCELLED, updated.Status)
	assert.Equal(t, "buyer returned goods", updated.CancelReason)
}

// ─── Calendar ──────────────────────────────────────────────────────────

func TestCreateCalendarEntry(t *testing.T) {
	ts := setupTaxTest(t)
	body := `{"company_id":"` + ts.compID + `","tax_type":"VAT","period_type":"MONTHLY","period_year":2026,"period_number":4,"start_date":"2026-04-01","end_date":"2026-04-30","declaration_due":"2026-05-20","payment_due":"2026-05-20"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tax/calendar", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetCalendarByPeriod(t *testing.T) {
	ts := setupTaxTest(t)
	cal := &domain.TaxCalendar{
		CompanyID:      ts.compID,
		TaxType:        domain.TaxTypeVAT,
		PeriodType:     domain.PeriodTypeMonthly,
		PeriodYear:     2026,
		PeriodNumber:   1,
		StartDate:      "2026-01-01",
		EndDate:        "2026-01-31",
		DeclarationDue: "2026-02-20",
		Status:         domain.CalStatusPENDING,
	}
	require.NoError(t, ts.svc.CreateCalendarEntry(nil, cal))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tax/calendar/period/"+ts.compID+"/2026/1", nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var entries []domain.TaxCalendar
	json.Unmarshal(w.Body.Bytes(), &entries)
	assert.Len(t, entries, 1)
}

// ─── Audit Cases ───────────────────────────────────────────────────────

func TestCreateAuditCase(t *testing.T) {
	ts := setupTaxTest(t)
	body := `{"company_id":"` + ts.compID + `","audit_period_start":"2025-01-01","audit_period_end":"2025-12-31","audit_decision_number":"QD-001/2025","auditor_name":"Tax Dept"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tax/audit-cases", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCloseAuditCase(t *testing.T) {
	ts := setupTaxTest(t)
	a := &domain.TaxAuditCase{
		CompanyID:        ts.compID,
		AuditPeriodStart: "2025-01-01",
		AuditPeriodEnd:   "2025-12-31",
		AuditDecNumber:   "QD-001",
		AuditorName:      "Tax Dept",
		Status:           domain.AuditCaseOPEN,
	}
	require.NoError(t, ts.svc.CreateAuditCase(nil, a))

	w := httptest.NewRecorder()
	body := `{"findings":"no issues","penalty_amount":0}`
	req, _ := http.NewRequest("POST", "/api/v1/tax/audit-cases/"+a.ID+"/close", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	updated, _ := ts.svc.GetAuditCase(nil, a.ID)
	assert.Equal(t, domain.AuditCaseCLOSED, updated.Status)
}

// ─── Calculations ──────────────────────────────────────────────────────

func TestCalculateVAT(t *testing.T) {
	ts := setupTaxTest(t)
	body := `{
		"company_id":"` + ts.compID + `",
		"period":{"period_type":"MONTHLY","period_year":2026,"period_number":1},
		"entries":[
			{"entry_number":"E001","description":"Sales","lines":[
				{"account_code":"5111","debit_amount":0,"credit_amount":10000000}
			]},
			{"entry_number":"E002","description":"Purchase","lines":[
				{"account_code":"152","debit_amount":5000000,"credit_amount":0}
			]}
		]
	}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tax/calculate/vat", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var result domain.VATResult
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, ts.compID, result.CompanyID)
	assert.Equal(t, 1000000.0, result.OutputVAT)   // 10M * 10%
	assert.Equal(t, 500000.0, result.InputVAT)     // 5M * 10%
	assert.Equal(t, 500000.0, result.VATPayable)   // 1M - 500K
}

func TestCalculateCIT(t *testing.T) {
	ts := setupTaxTest(t)
	body := `{
		"company_id":"` + ts.compID + `",
		"year":2026,
		"entries":[
			{"entry_number":"E001","description":"Revenue","lines":[
				{"account_code":"5111","debit_amount":0,"credit_amount":100000000}
			]},
			{"entry_number":"E002","description":"Expenses","lines":[
				{"account_code":"641","debit_amount":60000000,"credit_amount":0}
			]},
			{"entry_number":"E003","description":"Non-deductible","lines":[
				{"account_code":"821","debit_amount":5000000,"credit_amount":0}
			]}
		]
	}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tax/calculate/cit", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var result domain.CITResult
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, ts.compID, result.CompanyID)
	assert.Equal(t, 100000000.0, result.Revenue)
	assert.Equal(t, 60000000.0, result.Expenses)
	assert.Equal(t, 5000000.0, result.NonDeductible)
	assert.Equal(t, 45000000.0, result.TaxableIncome) // 100M - 60M + 5M
	assert.Equal(t, 15.0, result.TaxRate)             // 100M revenue → MICRO
	assert.Equal(t, 6750000.0, result.CITPayable)     // 45M * 15%
}

func TestCalculateCIT_IncentiveReduction(t *testing.T) {
	ts := setupTaxTest(t)
	// Seed incentive rate: 50% reduction
	ctx := context.Background()
	require.NoError(t, ts.taxRepo.CreateRate(ctx, &domain.TaxRate{
		RateCode: "CIT_INCENTIVE_NEW_PROJECT", TaxType: domain.TaxTypeCIT,
		RateType: domain.RateTypePERCENTAGE, RateValue: 20,
		IncentiveReducPct: 50,
		ApplicableTo: "INCENTIVE_NEW_PROJECT",
		EffectiveFrom: "2026-01-01", IsActive: true,
	}))

	body := `{
		"company_id":"` + ts.compID + `",
		"year":2026,
		"entries":[
			{"entry_number":"E001","description":"Revenue","lines":[
				{"account_code":"5111","debit_amount":0,"credit_amount":100000000}
			]},
			{"entry_number":"E002","description":"Expenses","lines":[
				{"account_code":"641","debit_amount":60000000,"credit_amount":0}
			]}
		]
	}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tax/calculate/cit", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var result domain.CITResult
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, 40000000.0, result.TaxableIncome)
	assert.Equal(t, 8000000.0, result.CITPayable)  // 40M * 20%
	assert.Equal(t, 4000000.0, result.CITFinal)     // 50% reduction
	assert.Equal(t, 50.0, result.IncentiveReduc)
}

func TestCalculateCIT_LossCarryForward(t *testing.T) {
	ts := setupTaxTest(t)
	ctx := context.Background()

	// Year 2025: create a loss
	body2025 := `{
		"company_id":"` + ts.compID + `",
		"year":2025,
		"entries":[
			{"entry_number":"E001","description":"Revenue","lines":[
				{"account_code":"5111","debit_amount":0,"credit_amount":100000000}
			]},
			{"entry_number":"E002","description":"Expenses","lines":[
				{"account_code":"641","debit_amount":120000000,"credit_amount":0}
			]}
		]
	}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tax/calculate/cit", strings.NewReader(body2025))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Year 2026: profit, loss should reduce taxable income
	body2026 := `{
		"company_id":"` + ts.compID + `",
		"year":2026,
		"entries":[
			{"entry_number":"E003","description":"Revenue","lines":[
				{"account_code":"5111","debit_amount":0,"credit_amount":100000000}
			]},
			{"entry_number":"E004","description":"Expenses","lines":[
				{"account_code":"641","debit_amount":60000000,"credit_amount":0}
			]}
		]
	}`
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/v1/tax/calculate/cit", strings.NewReader(body2026))
	req2.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	var result domain.CITResult
	json.Unmarshal(w2.Body.Bytes(), &result)
	assert.Equal(t, 20000000.0, result.TaxableIncome) // 40M - 20M loss
	assert.Equal(t, 20000000.0, result.LossUsed)

	_ = ctx // used indirectly
}

func TestCalculateCIT_ThinCap(t *testing.T) {
	ts := setupTaxTest(t)
	body := `{
		"company_id":"` + ts.compID + `",
		"year":2026,
		"entries":[
			{"entry_number":"E001","lines":[
				{"account_code":"5111","debit_amount":0,"credit_amount":100000000}
			]},
			{"entry_number":"E002","lines":[
				{"account_code":"641","debit_amount":40000000,"credit_amount":0}
			]},
			{"entry_number":"E003","lines":[
				{"account_code":"635","debit_amount":30000000,"credit_amount":0}
			]}
		]
	}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tax/calculate/cit", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var result domain.CITResult
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, 12000000.0, result.NonDeductible) // thin cap disallowance
	assert.Equal(t, 42000000.0, result.TaxableIncome)
}

func TestCalculateQuarterlyProvisional(t *testing.T) {
	ts := setupTaxTest(t)
	body := `{
		"company_id":"` + ts.compID + `",
		"year":2026,
		"quarter":3,
		"entries":[
			{"entry_number":"E001","lines":[
				{"account_code":"5111","debit_amount":0,"credit_amount":75000000}
			]},
			{"entry_number":"E002","lines":[
				{"account_code":"641","debit_amount":45000000,"credit_amount":0}
			]}
		]
	}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tax/calculate/cit/provisional", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var result domain.QuarterlyProvisionalResult
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, 3, result.Quarter)
	assert.Equal(t, 6000000.0, result.EstimatedAnnualCIT)
}

func TestCalculateVAT_ZeroInput(t *testing.T) {
	ts := setupTaxTest(t)
	body := `{
		"company_id":"` + ts.compID + `",
		"period":{"period_type":"MONTHLY","period_year":2026,"period_number":1},
		"entries":[
			{"entry_number":"E001","description":"Sales","lines":[
				{"account_code":"5111","debit_amount":0,"credit_amount":10000000}
			]}
		]
	}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tax/calculate/vat", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var result domain.VATResult
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, 1000000.0, result.OutputVAT)
	assert.Equal(t, 0.0, result.InputVAT)
	assert.Equal(t, 1000000.0, result.VATPayable)
}

func TestCalculatePIT(t *testing.T) {
	ts := setupTaxTest(t)
	body := `{
		"company_id":"` + ts.compID + `",
		"period":{"period_type":"MONTHLY","period_year":2026,"period_number":1},
		"employees":[
			{"employee_id":"e1","gross_monthly":25000000,"dependants":1,"months":1}
		]
	}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tax/calculate/pit", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var result domain.PITResult
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, ts.compID, result.CompanyID)
	assert.Equal(t, 1, result.EmployeeCount)
	assert.Equal(t, 33750.0, result.TotalPIT)
}

func TestTimeDependentOperations(t *testing.T) {
	ts := setupTaxTest(t)

	// Create draft → submit → reject → resubmit → acknowledge → amend
	d := &domain.TaxDeclaration{
		CompanyID:       ts.compID,
		DeclarationType: domain.DeclTypeGTGT01,
		TaxPeriod:       domain.TaxPeriod{PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 1},
	}
	require.NoError(t, ts.svc.CreateDeclaration(nil, d))
	assert.Equal(t, domain.DeclStatusDRAFT, d.Status)

	// Draft → Submitted
	require.NoError(t, ts.svc.SubmitDeclaration(nil, d.ID, "user-1"))
	updated1, _ := ts.svc.GetDeclaration(nil, d.ID)
	assert.Equal(t, domain.DeclStatusSUBMITTED, updated1.Status)

	// Submitted → Rejected
	require.NoError(t, ts.svc.RejectDeclaration(nil, d.ID, "fix data"))
	updated2, _ := ts.svc.GetDeclaration(nil, d.ID)
	assert.Equal(t, domain.DeclStatusREJECTED, updated2.Status)

	// Rejected → Resubmit → Acknowledge
	require.NoError(t, ts.svc.SubmitDeclaration(nil, d.ID, "user-1"))
	require.NoError(t, ts.svc.AcknowledgeDeclaration(nil, d.ID, "GDT-ACK-001"))
	updated3, _ := ts.svc.GetDeclaration(nil, d.ID)
	assert.Equal(t, domain.DeclStatusACKNOWLEDGED, updated3.Status)

	// Cancel draft declaration
	d2 := &domain.TaxDeclaration{
		CompanyID:       ts.compID,
		DeclarationType: domain.DeclTypeGTGT01,
		TaxPeriod:       domain.TaxPeriod{PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 2},
	}
	require.NoError(t, ts.svc.CreateDeclaration(nil, d2))
	require.NoError(t, ts.svc.CancelDeclaration(nil, d2.ID))
	updated4, _ := ts.svc.GetDeclaration(nil, d2.ID)
	assert.Equal(t, domain.DeclStatusCANCELLED, updated4.Status)
}

func TestGenerateDeclaration_GTGT03(t *testing.T) {
	ts := setupTaxTest(t)
	je := &domain.JournalEntry{
		ID: "JE1", EntryNumber: "JE1", CompanyID: ts.compID,
		EntryDate: time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC),
		Status:    domain.JournalEntryPosted,
		Lines: []domain.JournalLine{
			{AccountCode: "5111", CreditAmount: 10000000},
			{AccountCode: "33311", CreditAmount: 1000000},
		},
	}
	require.NoError(t, ts.jeRepo.Create(context.Background(), je))

	body := `{
		"company_id":"` + ts.compID + `",
		"declaration_type":"GTGT03",
		"tax_period":{"period_type":"QUARTERLY","period_year":2026,"period_number":2}
	}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tax/declarations/generate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp domain.TaxDeclaration
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, domain.DeclTypeGTGT03, resp.DeclarationType)
	assert.Equal(t, domain.DeclStatusVALIDATED, resp.Status)
	assert.Equal(t, "QUARTERLY", string(resp.TaxPeriod.PeriodType))
}

func TestReconcileVAT(t *testing.T) {
	ts := setupTaxTest(t)
	decl := &domain.TaxDeclaration{
		ID: "DECL-1", CompanyID: ts.compID, DeclarationType: domain.DeclTypeGTGT01,
		TaxPeriod: domain.TaxPeriod{PeriodType: domain.PeriodTypeMonthly, PeriodYear: 2026, PeriodNumber: 4},
		Status:    domain.DeclStatusACKNOWLEDGED, AdjustmentType: domain.AdjTypeNONE,
		Lines: []domain.TaxDeclarationLine{{LineCode: "23", LineName: "Tong", Amount: 100000}},
	}
	require.NoError(t, ts.svc.CreateDeclaration(nil, decl))
	inv := &domain.EInvoice{
		CompanyID:   ts.compID,
		Pattern:     "01GTKT0/001",
		Serial:      "AA/25E",
		InvoiceType: domain.EInvTypeORIGINAL,
		BuyerName:   "Buyer",
		CurrencyCode: "VND",
		IssueDate:   "2026-04-15",
		Status:      domain.EInvStatusISSUED,
		Subtotal:    1000000,
		VATAmount:   100000,
		GrandTotal:  1100000,
		Lines:       []domain.EInvoiceLine{{LineNumber: 1, Description: "S", Quantity: 1, UnitPrice: 1000000, LineTotal: 1000000, VATRate: 10, VATAmount: 100000}},
	}
	require.NoError(t, ts.svc.CreateEInvoice(nil, inv))

	body := `{"company_id":"` + ts.compID + `","period":{"period_type":"MONTHLY","period_year":2026,"period_number":4}}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tax/reconcile/vat", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var res domain.VATReconciliationResult
	json.Unmarshal(w.Body.Bytes(), &res)
	assert.True(t, res.Matched)
	assert.Equal(t, 100000.0, res.InvoiceTotal)
}

func TestValidateTaxCode_Active(t *testing.T) {
	ts := setupTaxTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tax/validate-tax-code?tax_code=0100123456", nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp domain.TaxCodeLookupResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "ACTIVE", resp.Status)
	assert.Equal(t, "Active Corp", resp.Name)
}

func TestValidateTaxCode_NotFound(t *testing.T) {
	ts := setupTaxTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tax/validate-tax-code?tax_code=NOTFOUND", nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp domain.TaxCodeLookupResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "NOT_FOUND", resp.Status)
}

func TestValidateTaxCode_Inactive(t *testing.T) {
	ts := setupTaxTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tax/validate-tax-code?tax_code=INACTIVE", nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp domain.TaxCodeLookupResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "INACTIVE", resp.Status)
	assert.Equal(t, "Inactive Corp", resp.Name)
}

func TestValidateTaxCode_MissingParam(t *testing.T) {
	ts := setupTaxTest(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tax/validate-tax-code", nil)
	ts.r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
