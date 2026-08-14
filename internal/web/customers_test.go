package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/auth"
	"gotax/internal/domain"
	"gotax/internal/repository"
	"gotax/internal/service"
)

type customersSetup struct {
	r    *gin.Engine
	repo *repository.MemorySaleRepo
}

// setupCustomers wires the real pages pipeline: memory sale repo → real
// SaleService → Deps → NewServer template sets → RegisterPages, authenticated
// with a real RS256 token through PageAuthMiddleware.
func setupCustomers(t *testing.T) *customersSetup {
	t.Helper()
	auth.SetJWTSecret("web-test-secret")

	repo := repository.NewMemorySaleRepo()
	saleSvc := service.NewSaleService(repo, repo, repo, repo, repo, repo, repo, repo, nil)
	deps := Deps{Sale: saleSvc}

	srv, err := NewServer([]string{"customers", "legacy"})
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

	return &customersSetup{r: r, repo: repo}
}

// Seed one customer through the real service so the page renders repo state,
// not mock data.
func (s *customersSetup) seedCustomer(t *testing.T, code, name string) {
	t.Helper()
	err := s.repo.CreateCustomer(context.Background(), &domain.Customer{
		CompanyID: "CMP001",
		Code:      code,
		Name:      name,
		TaxCode:   "0123456789",
		Currency:  "VND",
		Status:    domain.CustomerActive,
	})
	require.NoError(t, err)
}

func TestCustomersPageRender(t *testing.T) {
	s := setupCustomers(t)
	s.seedCustomer(t, "KH-TEST", "Công ty kiểm thử XYZ")

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/app/customers.html?company_id=CMP001", nil)
	s.r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "Khách hàng")
	assert.Contains(t, body, "KH-TEST")
	assert.Contains(t, body, "Công ty kiểm thử XYZ")
	// Converted pages must not carry Alpine directives anymore.
	assert.NotContains(t, body, "x-data")
}

func TestCustomersCreateAction(t *testing.T) {
	s := setupCustomers(t)

	form := url.Values{
		"code":    {"KH-NEW"},
		"name":    {"Công ty mới"},
		"tax_code": {"0987654321"},
		"phone":   {"0900 000 000"},
		"email":   {"new@co.vn"},
		"address": {"Hà Nội"},
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/app/customers/create?company_id=CMP001", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	s.r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	// Fragment re-render shows the new row.
	assert.Contains(t, body, "KH-NEW")
	assert.Contains(t, body, "Công ty mới")
	// Repo state actually changed (not client-only like the old Alpine page).
	custs, err := s.repo.ListCustomers(context.Background(), "CMP001")
	require.NoError(t, err)
	require.Len(t, custs, 1)
	assert.Equal(t, "KH-NEW", custs[0].Code)
	assert.Equal(t, "new@co.vn", custs[0].Email)
}

func TestCustomersCreateActionValidationError(t *testing.T) {
	s := setupCustomers(t)

	form := url.Values{"code": {"KH-BAD"}} // missing name + tax_code
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/app/customers/create?company_id=CMP001", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	s.r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	// Error surfaced as toast trigger, not a crash.
	assert.Contains(t, w.Header().Get("HX-Trigger"), "error")
	custs, err := s.repo.ListCustomers(context.Background(), "CMP001")
	require.NoError(t, err)
	assert.Empty(t, custs)
}
