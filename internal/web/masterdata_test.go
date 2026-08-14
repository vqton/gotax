package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/auth"
	"gotax/internal/domain"
	"gotax/internal/repository"
	"gotax/internal/service"
)

// setupSvc wires the core service facade (periods, exchange rates, cash)
// on in-memory repos, plus the company service.
func setupSvc(t *testing.T) (*gin.Engine, service.CompanyService, *repository.MemoryPeriodRepo, *repository.MemoryExchangeRateRepo, *repository.MemoryCashRepo, *repository.MemoryJournalRepo) {
	t.Helper()
	auth.SetJWTSecret("web-test-secret")

	accRepo := repository.NewMemoryAccountRepo()
	jeRepo := repository.NewMemoryJournalRepo()
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

	svc := service.NewService(accRepo, jeRepo, perRepo, userRepo, auditRepo, rateRepo, templateRepo, approvalRepo, versionRepo, mappingRepo, analysisRepo, ifrsRepo, refreshRepo, resetRepo, obRepo, cashRepo)
	companySvc := service.NewCompanyService(repository.NewMemoryCompanyRepo())
	deps := Deps{Svc: svc, Company: companySvc}

	srv, err := NewServer([]string{"company", "exchange-rates", "periods", "cash-receipts", "cash-payments", "cash-transfers", "cash-book", "cash-flow", "legacy"})
	require.NoError(t, err)

	r := gin.New()
	token, err := auth.GenerateAccessToken(&domain.User{ID: "web-test-user", Username: "admin", Role: domain.UserRoleAdmin})
	require.NoError(t, err)
	// Engine middleware must be registered before routes: gin snapshots the
	// engine handler chain into route nodes at registration time.
	r.Use(func(c *gin.Context) {
		c.Request.Header.Set("Authorization", "Bearer "+token)
	})
	srv.RegisterPages(r, NewPages(deps), srv.NewActions(deps))

	return r, companySvc, perRepo, rateRepo, cashRepo, jeRepo
}

func TestCompanyPageRenderAndSave(t *testing.T) {
	r, companySvc, _, _, _, _ := setupSvc(t)
	companySvc.CreateCompany(context.Background(), &domain.Company{
		ID:               "CMP001",
		LegalForm:        domain.LegalFormLLC1Member,
		AccountingRegime: domain.RegimeTT99,
		LegalNameVN:      "Công ty TNHH Kế toán GOTAX",
		TaxCode:          "0100000001",
		RegAddress:       "Số 1 Tràng Tiền, Hoàn Kiếm, Hà Nội",
		Phone:            "024 0000 0001",
		Email:            "info@gotax.vn",
		LegalRepName:     "Nguyễn Văn A",
		LegalRepTitle:    "Giám đốc",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/app/company.html?company_id=CMP001", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "Thông tin công ty")
	assert.Contains(t, body, "Công ty TNHH Kế toán GOTAX")
	assert.Contains(t, body, "0100000001")
	assert.NotContains(t, body, "x-data")

	// Save edits legal name + tax code.
	form := url.Values{
		"legal_name_vn": {"Công ty TNHH Kế toán GOTAX (đổi tên)"},
		"tax_code":      {"0100000001"},
		"reg_address":   {"Số 1 Tràng Tiền, Hoàn Kiếm, Hà Nội"},
	}
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/app/company/save?company_id=CMP001", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("HX-Trigger"), "success")
	got, err := companySvc.GetCompany(context.Background(), "CMP001")
	require.NoError(t, err)
	assert.Equal(t, "Công ty TNHH Kế toán GOTAX (đổi tên)", got.LegalNameVN)
	assert.NotEmpty(t, got.ID)
}

func TestCompanyPageSaveMissingFields(t *testing.T) {
	r, companySvc, _, _, _, _ := setupSvc(t)
	companySvc.CreateCompany(context.Background(), &domain.Company{
		ID: "CMP001", LegalNameVN: "Công ty TNHH A", TaxCode: "0100000001",
		LegalForm: domain.LegalFormLLC1Member, AccountingRegime: domain.RegimeTT99,
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/app/company/save?company_id=CMP001",
		strings.NewReader(url.Values{"tax_code": {"0100000001"}}.Encode())) // missing legal_name_vn
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	assert.Contains(t, w.Header().Get("HX-Trigger"), "error")
	got, err := companySvc.GetCompany(context.Background(), "CMP001")
	require.NoError(t, err)
	assert.Equal(t, "Công ty TNHH A", got.LegalNameVN)
}

func TestExchangeRatesPageRenderAndCreate(t *testing.T) {
	r, _, _, rateRepo, _, _ := setupSvc(t)
	rateRepo.Create(context.Background(), &domain.ExchangeRate{
		CurrencyCode: "USD", RateDate: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC),
		BuyRate: 25400, SellRate: 25500, AverageRate: 25450,
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/app/exchange-rates.html", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "Tỷ giá")
	assert.Contains(t, body, "USD")
	assert.Contains(t, body, "25.450")
	assert.NotContains(t, body, "x-data")

	form := url.Values{
		"currency_code": {"EUR"},
		"rate_date":     {"2026-08-14"},
		"buy_rate":      {"26400"},
		"sell_rate":     {"26550"},
		"average_rate":  {"26475"},
	}
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/app/exchange-rates/create", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "EUR")
	assert.Contains(t, w.Body.String(), "26.400")
	rates, err := rateRepo.GetAll(context.Background())
	require.NoError(t, err)
	require.Len(t, rates, 2)
	var eur *domain.ExchangeRate
	for i := range rates {
		if rates[i].CurrencyCode == "EUR" {
			eur = &rates[i]
		}
	}
	require.NotNil(t, eur)
	assert.Equal(t, 26475.0, eur.AverageRate)
	assert.NotEmpty(t, eur.ID)
}

func TestExchangeRatesCreateValidationError(t *testing.T) {
	r, _, _, rateRepo, _, _ := setupSvc(t)

	form := url.Values{
		"currency_code": {"USD"},
		"rate_date":     {"2026-08-14"},
		// missing average_rate → Validate fails
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/app/exchange-rates/create", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	assert.Contains(t, w.Header().Get("HX-Trigger"), "error")
	rates, err := rateRepo.GetAll(context.Background())
	require.NoError(t, err)
	assert.Empty(t, rates)
}

func TestPeriodsPageRenderAndCreate(t *testing.T) {
	r, _, perRepo, _, _, _ := setupSvc(t)
	perRepo.Create(context.Background(), &domain.Period{
		Year: 2026, Month: 7, Status: domain.PeriodOpen,
		StartDate: time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 7, 31, 0, 0, 0, 0, time.UTC),
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/app/periods.html", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "Kỳ kế toán")
	assert.Contains(t, body, "7")
	assert.Contains(t, body, "2026")
	assert.NotContains(t, body, "x-data")

	form := url.Values{"year": {"2026"}, "month": {"8"}}
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/app/periods/create", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	got, err := perRepo.GetByYearMonth(context.Background(), 2026, 8)
	require.NoError(t, err)
	assert.Equal(t, 2026, got.Year)
	assert.Equal(t, 8, got.Month)
	assert.Equal(t, domain.PeriodOpen, got.Status)
	assert.False(t, got.StartDate.IsZero())
	assert.False(t, got.EndDate.IsZero())
}

func TestPeriodsCreateValidationError(t *testing.T) {
	r, _, perRepo, _, _, _ := setupSvc(t)

	// invalid month → facade Validate() rejects
	form := url.Values{"year": {"2026"}, "month": {"13"}}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/app/periods/create", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	assert.Contains(t, w.Header().Get("HX-Trigger"), "error")
	all, err := perRepo.GetAll(context.Background())
	require.NoError(t, err)
	assert.Empty(t, all)
}

func TestPeriodsCloseAndReopen(t *testing.T) {
	r, _, perRepo, _, _, _ := setupSvc(t)
	perRepo.Create(context.Background(), &domain.Period{
		ID: "P202608", Year: 2026, Month: 8, Status: domain.PeriodOpen,
		StartDate: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 8, 31, 0, 0, 0, 0, time.UTC),
	})

	closeForm := url.Values{"id": {"P202608"}}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/app/periods/close", strings.NewReader(closeForm.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("HX-Trigger"), "success")
	p, err := perRepo.GetByID(context.Background(), "P202608")
	require.NoError(t, err)
	assert.Equal(t, domain.PeriodClosed, p.Status)

	reopenForm := url.Values{"id": {"P202608"}}
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/app/periods/reopen", strings.NewReader(reopenForm.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("HX-Trigger"), "success")
	p, err = perRepo.GetByID(context.Background(), "P202608")
	require.NoError(t, err)
	assert.Equal(t, domain.PeriodOpen, p.Status)
}

func TestPeriodsCloseValidationError(t *testing.T) {
	r, _, _, _, _, _ := setupSvc(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/app/periods/close",
		strings.NewReader(url.Values{"id": {"P-NOPE"}}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("HX-Trigger"), "error")
}
